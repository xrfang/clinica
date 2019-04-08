package main

import (
	"net/http"
)

func users(w http.ResponseWriter, r *http.Request) {
	s := sessions.Get(w, r)
	var u user
	if s.Unmarshal("user", &u) != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	if u.Role != RoleAdmin {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	switch r.Method {
	case "GET":
		renderTemplate(w, "users.html", struct {
			Admin string
			Users []user
		}{
			Admin: u.Caption(),
			Users: listUsers(),
		})
	case "POST":
		r.ParseForm()
		code, mesg := setUser(r.Form)
		http.Error(w, mesg, code)
	case "DELETE":
		login := r.URL.Query().Get("login")
		if login == u.Login {
			http.Error(w, "不能删除当前用户", http.StatusForbidden)
			return
		}
		code, mesg := delUser(login)
		http.Error(w, mesg, code)
	}
}
