package main

import (
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
	var (
		args []interface{}
		prop []string
		upds []string
		stmt string
	)
	mode := r.PostFormValue("mode")
	args = append(args, mode)
	prop = append(prop, "mode")
	upds = append(upds, "mode=?")

	status := r.PostFormValue("status")
	args = append(args, status)
	prop = append(prop, "status")
	upds = append(upds, "status=?")

	at := r.PostFormValue("time")
	at, err := fmtDateTime(at, "Y-m-d H:i")
	if err == nil {
		args = append(args, at)
		prop = append(prop, "time")
		upds = append(upds, "time=?")
	}

	args = append(args, time.Now().Format("2006-01-02 15:04:05"))
	prop = append(prop, "updated")
	upds = append(upds, "updated=?")

	id := r.PostFormValue("id")
	if id == "" {
		caseID := r.PostFormValue("case_id")
		args = append(args, caseID)
		prop = append(prop, "case_id")
		stmt = `INSERT INTO consults (` + strings.Join(prop, ",") + `) VALUES (?` +
			strings.Repeat(`,?`, len(prop)-1) + `)`
	} else {
		args = append(args, id)
		stmt = `UPDATE consults SET ` + strings.Join(upds, ",") + ` WHERE id=?`
	}
	_, err = cf.dbx.Exec(stmt, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
