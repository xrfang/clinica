package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	audit "github.com/xrfang/go-audit"
	yaml "gopkg.in/yaml.v2"
)

type Configuration struct {
	LogFile string `yaml:"log_file"`
	Port    string `yaml:"port"`
	WebRoot string `yaml:"web_root"`
	DBPath  string `yaml:"dbPpath"`
	TLSKey  string `yaml:"tls_key"`
	TLSCrt  string `yaml:"tls_crt"`
	binPath string
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

var cf Configuration

func loadConfig(fn string) {
	cf.binPath = path.Dir(os.Args[0])
	cf.Port = "8080"
	cf.WebRoot = "../webroot"
	cf.LogFile = "../log/log"
	cf.DBPath = "../conf/config.db"
	if err := cf.load(fn); err != nil {
		fmt.Printf("[ERROR]cf.load(%s): %v\n", fn, err)
		os.Exit(1)
	}
}
