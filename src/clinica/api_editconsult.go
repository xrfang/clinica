package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func apiEditConsult(w http.ResponseWriter, r *http.Request) {
	s := sessions.Get(w, r)
	var u user
	if s.Unmarshal("user", &u) != nil {
		http.Error(w, "访问令牌失效，请重新登录", http.StatusUnauthorized)
		return
	}
	if u.Role == RoleDisabled {
		http.Error(w, "您的使用权限被冻结，请联系管理员", http.StatusForbidden)
		return
	}
	var args []interface{}
	var props []string

	mode := r.PostFormValue("mode")
	args = append(args, mode)
	props = append(props, "mode=?")

	status := r.PostFormValue("status")
	args = append(args, status)
	props = append(props, "status=?")

	at := r.PostFormValue("time")
	at, err := fmtDateTime(at, "Y-m-d H:i")
	if err == nil {
		args = append(args, at)
		props = append(props, "time=?")
	}

	args = append(args, time.Now().Format("2006-01-02 15:04:05"))
	props = append(props, "updated=?")

	id := r.PostFormValue("id")
	args = append(args, id)

	cmd := fmt.Sprintf(`UPDATE consults SET %s WHERE id=?`, strings.Join(props, ","))
	_, err = cf.dbx.Exec(cmd, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
