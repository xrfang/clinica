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
	var u user
	if s.Unmarshal("user", &u) != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	caption := u.Login
	if u.Name.String != "" {
		caption += " (" + u.Name.String + ")"
	}
	renderTemplate(w, "home.html", struct {
		Caption  string
		IsEditor bool
		IsAdmin  bool
	}{
		Caption:  u.Caption(),
		IsAdmin:  u.Role == RoleAdmin,
		IsEditor: u.Role >= RoleEditor,
	})
}
