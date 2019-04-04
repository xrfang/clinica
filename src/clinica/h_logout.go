package main

import "net/http"

func logout(w http.ResponseWriter, r *http.Request) {
	s := sessions.Get(w, r)
	s.Del("user")
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}
