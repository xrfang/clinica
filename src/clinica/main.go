package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	audit "github.com/xrfang/go-audit"
	res "github.com/xrfang/go-res"
)

func main() {
	ver := flag.Bool("version", false, "show version info")
	pkg := flag.String("pack", "", "pack resources under directory")
	dbg := flag.Bool("debug", false, "debug mode")
	flag.Parse()
	if *ver {
		fmt.Println(verinfo())
		return
	}
	if *pkg != "" {
		audit.Assert(res.Pack(*pkg))
		fmt.Printf("resources under '%s' packed.\n", *pkg)
		return
	}
	loadConfig()
	if !*dbg {
		audit.Assert(res.Extract(cf.WebRoot, res.OverwriteIfNewer))
	}
	audit.ExpVars(map[string]interface{}{
		"config":  cf,
		"version": _G_REVS + "." + _G_HASH,
	})
	audit.SetLogFile(cf.LogFile)
	audit.ExpLogs()
	audit.SetDebugging(*dbg)
	setupRoutes()
	svr := http.Server{
		Addr:         ":" + cf.Port,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
	audit.Assert(svr.ListenAndServe())
}
