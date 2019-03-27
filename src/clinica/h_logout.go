package main

import "net/http"

func logout(w http.ResponseWriter, r *http.Request) {
	//TODO: delete cookie/session
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}
