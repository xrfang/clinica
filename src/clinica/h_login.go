package main

import (
	"net/http"

	audit "github.com/xrfang/go-audit"
)

func login(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			http.Error(w, e.(error).Error(), http.StatusInternalServerError)
		}
	}()
	s := sessions.Get(w, r)
	audit.Assert(r.ParseForm())
	user := r.Form.Get("user")
	pass := r.Form.Get("pass")
	var mesg string
	if user != "" && pass != "" {
		u := getUser(user, pass)
		if u != nil {
			setCookie(w, "user", user, 86400*30)
			s.Marshal("user", u)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		mesg = "用户名或密码错误"
	}
	if user == "" {
		user = getCookie(r, "user")
	}
	renderTemplate(w, "login.html", struct {
		User string
		Err  string
	}{user, mesg})
}
