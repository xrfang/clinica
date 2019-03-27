package main

import (
	"math/rand"
	"net/http"

	audit "github.com/xrfang/go-audit"
)

func login(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			http.Error(w, e.(error).Error(), http.StatusInternalServerError)
		}
	}()
	audit.Assert(r.ParseForm())
	user := r.Form.Get("user")
	pass := r.Form.Get("pass")
	var mesg string
	if user != "" && pass != "" {
		mesg = "incorrect username or password"
		//TODO: authenticate user
		if rand.Intn(100) > 49 {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
	}
	if user == "" {
		user = getCookie(r, "user")
	}
	renderTemplate(w, "login.html", struct {
		User string
		Err  string
	}{user, mesg})
}
