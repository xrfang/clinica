package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	audit "github.com/xrfang/go-audit"
	yaml "gopkg.in/yaml.v2"
)

const (
	RoleDisabled = -1
	RoleReader   = 0
	RoleEditor   = 1
	RoleAdmin    = 2
)

type Configuration struct {
	LogFile string `yaml:"log_file"`
	Port    string `yaml:"port"`
	WebRoot string `yaml:"web_root"`
	DBPath  string `yaml:"dbPpath"`
	TLSKey  string `yaml:"tls_key"`
	TLSCrt  string `yaml:"tls_crt"`
	binPath string
	dbx     *sqlx.DB
}

func (c Configuration) abs(fn string) string {
	if fn == "" || path.IsAbs(fn) {
		return fn
	}
	p, _ := filepath.Abs(path.Join(c.binPath, fn))
	return p
}

func (c *Configuration) load(fn string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	f, err := os.Open(fn)
	audit.Assert(err)
	defer f.Close()
	audit.Assert(yaml.NewDecoder(f).Decode(c))
	c.WebRoot = c.abs(c.WebRoot)
	c.LogFile = c.abs(c.LogFile)
	c.DBPath = c.abs(c.DBPath)
	c.TLSCrt = c.abs(c.TLSCrt)
	c.TLSKey = c.abs(c.TLSKey)
	return
}

func (c *Configuration) initDB() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(audit.Trace("initDB: %v", e).Error())
			os.Exit(1)
		}
	}()
	err := os.MkdirAll(path.Dir(c.DBPath), 0755)
	audit.Assert(err)
	dsn := fmt.Sprintf("file:%s?cache=shared", c.DBPath)
	c.dbx = sqlx.MustConnect("sqlite3", dsn)
	c.dbx.SetMaxOpenConns(1)
	//用户表
	c.dbx.MustExec(`CREATE TABLE IF NOT EXISTS users ( 
		login TEXT NOT NULL,
		passwd TEXT,
		name TEXT,
		role INTEGER NOT NULL DEFAULT 0,
		PRIMARY KEY(login)
	)`)
	c.dbx.MustExec(`INSERT OR IGNORE INTO users (login,passwd,name,role) VALUES (?,?,?,?)`,
		"admin", HashPassword("Password01!"), "管理员", RoleAdmin)
	//病案
	c.dbx.MustExec(`CREATE TABLE IF NOT EXISTS cases (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		gender INTEGER NOT NULL,
		age INTEGER NOT NULL,
		contact TEXT,
		memo TEXT,
		summary TEXT,
		date DATE NOT NULL,
		updated DATETIME NOT NULL
	)`)
	//条目
	c.dbx.MustExec(`CREATE TABLE IF NOT EXISTS items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		case_id INTEGER NOT NULL,
		type INTEGER NOT NULL,
		caption TEXT NOT NULL DEFAULT '',
		details TEXT NOT NULL DEFAULT '',
		date DATE NOT NULL,
		updated DATETIME NOT NULL,
		FOREIGN KEY(case_id) REFERENCES cases(id)
	)`)
}

var cf Configuration

func loadConfig(fn string) {
	cf.binPath = path.Dir(os.Args[0])
	cf.Port = "2562"
	cf.WebRoot = "../webroot"
	cf.LogFile = "../log/log"
	cf.DBPath = "../conf/clinica.db"
	if err := cf.load(fn); err != nil {
		fmt.Printf("[ERROR]cf.load(%s): %v\n", fn, err)
		os.Exit(1)
	}
	cf.initDB()
}
