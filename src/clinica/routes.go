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
	http.HandleFunc("/api/listcases", apiListCases)
	http.HandleFunc("/api/editcase", apiEditCase)
	http.HandleFunc("/api/editconsult", apiEditConsult)
	http.HandleFunc("/editcase", editCase)
	http.HandleFunc("/patients", patients)
}
