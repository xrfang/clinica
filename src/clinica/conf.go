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
	c.dbx.MustExec(`CREATE TABLE IF NOT EXISTS users --本系统用户表
	( 
		login  TEXT NOT NULL,    --登录用户名
		passwd TEXT,             --密码（使用BCrypt加密保存）
		name   TEXT,             --姓名
		role   INTEGER NOT NULL, --权限（-1=禁止登录；0=读者；1=编辑；2=管理员）
		PRIMARY KEY(login)
	)`)
	c.dbx.MustExec(`INSERT OR IGNORE INTO users (login,passwd,name,role) VALUES (?,?,?,?)`,
		"admin", HashPassword("Password01!"), "管理员", RoleAdmin)
	c.dbx.MustExec(`CREATE TABLE IF NOT EXISTS patients --患者注册表
	(
		id       INTEGER PRIMARY KEY AUTOINCREMENT,
		name     TEXT NOT NULL,     --姓名
		gender   INTEGER NOT NULL,  --性别
		birthday INTEGER NOT NULL,  --生日（格式：yyyymmdd）
		contact  TEXT,              --联系方式（一般为手机号）
		memo     TEXT               --备注
	)`)
	c.dbx.MustExec(`CREATE UNIQUE INDEX IF NOT EXISTS ident ON patients (name, contact)`)
	c.dbx.MustExec(`CREATE TABLE IF NOT EXISTS cases --医案表（一病一案，一个患者可有多个医案）
	(
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		patient_id INTEGER NOT NULL,  --患者ID
		summary    TEXT,              --简述（由医生填写，而非病人主诉）
		opened     DATE NOT NULL,     --首诊日期（格式：yyyymmdd）
		status     INTEGER NOT NULL,  --状态（0=尚未结束；1=痊愈/显效；2=失败；3=无反馈）
		updated    DATETIME NOT NULL, --最后编辑时间
		FOREIGN KEY(patient_id) REFERENCES patients(id)
	)`)
	c.dbx.MustExec(`CREATE TABLE IF NOT EXISTS sessions -- 就诊记录表（一个医案可有多个就诊记录）
	(
		id      INTEGER PRIMARY KEY AUTOINCREMENT,
		case_id INTEGER NOT NULL,  --医案ID
		mode    INTEGER NOT NULL,  --就诊方式（0=当面；1=远程直接沟通；2=他人代述）
		time    DATETIME NOT NULL, --就诊时间
		status  INTEGER NOT NULL,  --状态（0=就诊完成；1=预约中；2=未赴约；3=取消）
		updated DATETIME NOT NULL, --最后编辑时间
		FOREIGN KEY(case_id) REFERENCES cases(id)
	)`)
	c.dbx.MustExec(`CREATE TABLE IF NOT EXISTS records --诊疗记录表（一次就诊可有多个诊疗记录）
	(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id INTEGER NOT NULL,  --就诊记录ID
		type       INTEGER NOT NULL,  --记录类型（0=主诉；1=诊断；2=辩证；3=思路；4=开方）
		class_id   INTEGER,           --子类型
		caption    TEXT,              --标题
		details    TEXT,              --内容
		updated    DATETIME NOT NULL, --最后编辑时间
		FOREIGN    KEY(session_id) REFERENCES sessions(id),
		FOREIGN    KEY(class_id) REFERENCES classes(id)
	)`)
	c.dbx.MustExec(`CREATE TABLE IF NOT EXISTS classes --类型表
	(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type_id INTEGER NOT NULL, --元类型（即records.type）
		caption TEXT NOT NULL     --类型名称
	)`)
	c.dbx.MustExec(`CREATE UNIQUE INDEX IF NOT EXISTS cls ON classes (type_id, caption)`)
	c.dbx.Exec(`INSERT INTO classes (type_id,caption) VALUES (1, '望诊')`)
	c.dbx.Exec(`INSERT INTO classes (type_id,caption) VALUES (1, '闻声')`)
	c.dbx.Exec(`INSERT INTO classes (type_id,caption) VALUES (1, '闻味')`)
	c.dbx.Exec(`INSERT INTO classes (type_id,caption) VALUES (1, '问诊')`)
	c.dbx.Exec(`INSERT INTO classes (type_id,caption) VALUES (1, '脉诊')`)
	c.dbx.Exec(`INSERT INTO classes (type_id,caption) VALUES (1, '腹诊')`)
	c.dbx.Exec(`INSERT INTO classes (type_id,caption) VALUES (1, '病灶触诊')`)
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
