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
	var cbs []caseBrief
	cid := r.URL.Query().Get("id")
	if cid != "" {
		id, _ := strconv.Atoi(cid)
		if id <= 0 {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		var err error
		cbs, err = getCases(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(cbs) == 0 {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
	}
	renderTemplate(w, "edit.html", struct {
		Caption string
		Case    caseBrief
	}{
		Caption: u.Caption(),
		Case:    cbs[0],
	})
}
