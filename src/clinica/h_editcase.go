package main

import "net/http"

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
	renderTemplate(w, "edit.html", struct {
		Caption string
	}{
		Caption: u.Caption(),
	})
}
