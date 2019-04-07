package main

import (
	"net/http"
)

func setupRoutes() {
	http.HandleFunc("/", home)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/users", users)
	http.HandleFunc("/chpass", chpass)
	http.HandleFunc("/listcases", listCases)
}
