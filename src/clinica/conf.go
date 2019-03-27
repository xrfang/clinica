package main

import (
	"os"
	"path"
	"path/filepath"

	audit "github.com/xrfang/go-audit"
)

type Configuration struct {
	LogFile string
	Port    string
	WebRoot string
	//TODO: define configuration items
	binPath string
	cfgFile string
	cfgPath string
}

func (c Configuration) abs(fn string) string {
	if fn == "" || path.IsAbs(fn) {
		return fn
	}
	p, _ := filepath.Abs(path.Join(c.binPath, fn))
	return p
}

func (c *Configuration) Load(fn string) {
	f, err := os.Open(fn)
	audit.Assert(err)
	defer f.Close()
	//TODO: load configuration from f
	c.cfgFile = c.abs(fn)
	c.cfgPath = path.Dir(c.cfgFile)
}

var cf Configuration

func loadConfig() {
	cf.binPath = path.Dir(os.Args[0])
	cf.Port = "8080"
	cf.WebRoot = "../webroot"
	cf.LogFile = "../log/log"
	//TODO: load configuration from file, e.g. cf.Load(...)
	cf.WebRoot = cf.abs(cf.WebRoot)
	cf.LogFile = cf.abs(cf.LogFile)
}
