package main

import (
	"net/http"
)

func chpass(w http.ResponseWriter, r *http.Request) {
	s := sessions.Get(w, r)
	var u user
	if s.Unmarshal("user", &u) != nil {
		http.Error(w, "访问令牌失效，请重新登录", http.StatusUnauthorized)
		return
	}
	op := r.FormValue("old")
	if !CheckPasswordHash(op, u.Passwd.String) {
		http.Error(w, "原密码不正确", http.StatusForbidden)
		return
	}
	np := r.FormValue("new")
	_, err := cf.dbx.Exec(`UPDATE users SET passwd=? WHERE login=?`, HashPassword(np), u.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.Passwd.String = HashPassword(np)
	s.Marshal("user", u)
}
