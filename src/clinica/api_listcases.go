package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type caseBrief struct {
	ID      int       `json:"id"`
	Patient string    `json:"patient"`
	Summary string    `json:"summary"`
	Opened  time.Time `json:"opened"`
	Status  int       `json:"status"`
	Updated time.Time `json:"updated"`
}

func listCases(w http.ResponseWriter, r *http.Request) {
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
	var cbs []caseBrief
	err := cf.dbx.Select(&cbs, `SELECT cases.id,patients.name AS patient,summary,opened,status,updated FROM
	    cases,patients WHERE cases.patient_id=patients.id ORDER BY updated DESC LIMIT 100`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cbs)
}
