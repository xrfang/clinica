package main

import (
	"net/http"
	"regexp"
	"strings"
)

func patients(w http.ResponseWriter, r *http.Request) {
	s := sessions.Get(w, r)
	var u user
	if s.Unmarshal("user", &u) != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	if u.Role < RoleEditor {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	switch r.Method {
	case "GET":
		qry := strings.TrimSpace(r.URL.Query().Get("q"))
		r := regexp.MustCompile(`\s+`)
		var total int
		var ps []patient
		ts := r.Split(qry, -1)
		if len(ts) > 0 {
			for _, t := range ts {
				cnt, p := getPatients(t)
				total = cnt
				ps = append(ps, p...)
			}
		} else {
			total, ps = getPatients("")
		}
		renderTemplate(w, "patients.html", struct {
			Query    string
			Editor   string
			Total    int
			Patients []patient
		}{
			Query:    qry,
			Editor:   u.Caption(),
			Total:    total,
			Patients: ps,
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
