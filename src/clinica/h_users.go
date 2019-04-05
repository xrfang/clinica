package main

import (
	"net/http"
	"strconv"

	"github.com/xrfang/go-audit"
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
	}
	switch r.Method {
	case "GET":
		var us []user
		audit.Assert(cf.dbx.Select(&us, "SELECT name,login,role FROM users ORDER BY role DESC,login"))
		var ul []map[string]string
		for _, u := range us {
			var style string
			switch u.Role {
			case RoleDisabled:
				style = "secondary"
			case RoleReader:
				style = "info"
			case RoleEditor:
				style = "primary"
			case RoleAdmin:
				style = "danger"
			}
			ul = append(ul, map[string]string{
				"name":  u.Name,
				"login": u.Login,
				"role":  strconv.Itoa(u.Role),
				"style": style,
			})
		}
		renderTemplate(w, "users.html", ul)
	case "POST":

	case "DELETE":
		login := r.URL.Query().Get("login")
		cf.dbx.MustExec(`DELETE FROM users WHERE login=?`, login)
	}
}
