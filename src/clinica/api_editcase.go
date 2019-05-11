package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func apiEditCase(w http.ResponseWriter, r *http.Request) {
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

	summary := r.PostFormValue("summary")
	args = append(args, summary)
	props = append(props, "summary=?")

	status := r.PostFormValue("status")
	args = append(args, status)
	props = append(props, "status=?")

	opened := r.PostFormValue("opened")
	_, err := time.Parse("2006-01-02", opened)
	if err == nil {
		args = append(args, opened)
		props = append(props, "opened=?")
	}

	id := r.PostFormValue("id")
	args = append(args, id)

	cmd := fmt.Sprintf(`UPDATE cases SET %s WHERE id=?`, strings.Join(props, ","))
	_, err = cf.dbx.Exec(cmd, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
