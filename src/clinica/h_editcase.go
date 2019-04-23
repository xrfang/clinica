package main

import (
	"net/http"
	"strconv"
)

func editCase(w http.ResponseWriter, r *http.Request) {
	s := sessions.Get(w, r)
	var u user
	if s.Unmarshal("user", &u) != nil {
		http.Error(w, "访问令牌失效，请重新登录", http.StatusUnauthorized)
		return
	}
	if u.Role < RoleEditor {
		http.Error(w, "您没有编辑医案的权限", http.StatusForbidden)
		return
	}
	var consults []consult
	var patient string
	cid := r.URL.Query().Get("id")
	if cid != "" {
		id, _ := strconv.Atoi(cid)
		if id <= 0 {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		c, err := getCases(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(c) == 0 {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		patient = c[0].PatientName
		consults, err = getConsults(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	renderTemplate(w, "edit.html", struct {
		Caption  string
		Patient  string
		Consults []consult
	}{
		Caption:  u.Caption(),
		Patient:  patient,
		Consults: consults,
	})
}
