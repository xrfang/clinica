package main

import (
	"net/http"

	audit "github.com/xrfang/go-audit"
)

type user struct {
	Login  string
	Passwd string
	Name   string
	Role   int
}

func getUser(login, passwd string) *user {
	var u user
	err := cf.dbx.Get(&u, "SELECT * FROM users WHERE login=?", login)
	if err != nil || u.Role == RoleDisabled {
		return nil
	}
	if CheckPasswordHash(passwd, u.Passwd) {
		return &u
	}
	return nil
}

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
			s.Set("user", u)
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
