package main

import (
	"net/http"
	"path"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.ServeFile(w, r, path.Join(cf.WebRoot, r.URL.Path))
		return
	}
	s := sessions.Get(w, r)
	var u *user
	if !s.Get("user", &u) || u == nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	renderTemplate(w, "home.html", struct {
		IsAdmin bool
	}{
		IsAdmin: u.Role == RoleAdmin,
	})
}
