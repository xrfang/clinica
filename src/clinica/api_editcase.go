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
	var (
		args []interface{}
		prop []string
		upds []string
		stmt string
	)

	summary := r.PostFormValue("summary")
	args = append(args, summary)
	prop = append(prop, "summary")
	upds = append(upds, "summary=?")

	status := r.PostFormValue("status")
	args = append(args, status)
	prop = append(prop, "status")
	upds = append(upds, "status=?")

	opened := r.PostFormValue("opened")
	opened, _ = fmtDateTime(opened, "Y-m-d")

	args = append(args, time.Now().Format("2006-01-02 15:04:05"))
	prop = append(prop, "updated")
	upds = append(upds, "updated=?")

	id := r.PostFormValue("id")
	if id == "" {
		if opened == "" {
			opened = time.Now().Format("2006-01-02")
		}
		args = append(args, opened)
		prop = append(prop, "opened")
		patientID := r.PostFormValue("patient_id")
		args = append(args, patientID)
		prop = append(prop, "patient_id")
		stmt = `INSERT INTO cases (` + strings.Join(prop, ",") + `) VALUES (?` +
			strings.Repeat(`,?`, len(prop)-1) + `)`
	} else {
		if opened != "" {
			args = append(args, opened)
			upds = append(upds, "opened=?")
		}
		args = append(args, id)
		stmt = `UPDATE cases SET ` + strings.Join(upds, ",") + ` WHERE id=?`
	}
	res, err := cf.dbx.Exec(stmt, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if id != "" {
		fmt.Fprintln(w, id)
	} else {
		id, _ := res.LastInsertId()
		fmt.Fprintln(w, id)
	}
}
