package main

import (
	"html/template"
	"net/http"
	"path"
	"path/filepath"

	audit "github.com/xrfang/go-audit"
)

func getCookie(r *http.Request, name string) string {
	c, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return c.Value
}

func setCookie(w http.ResponseWriter, name, value string, age int) {
	http.SetCookie(w, &http.Cookie{
		Name:   name,
		Value:  value,
		Path:   "/",
		MaxAge: age,
		Secure: false,
	})
}

func renderTemplate(w http.ResponseWriter, tpl string, args interface{}) {
	defer func() {
		if e := recover(); e != nil {
			http.Error(w, e.(error).Error(), http.StatusInternalServerError)
		}
	}()
	helper := template.FuncMap{
		"ver": func() string {
			return "V" + _G_REVS + "." + _G_HASH
		},
	}
	tDir := path.Join(cf.WebRoot, "templates")
	t, err := template.New("body").Funcs(helper).ParseFiles(path.Join(tDir, tpl))
	audit.Assert(err)
	sfs, err := filepath.Glob(path.Join(tDir, "shared/*"))
	if len(sfs) > 0 {
		t, err = t.ParseFiles(sfs...)
		audit.Assert(err)
	}
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	audit.Assert(t.Execute(w, args))
}
