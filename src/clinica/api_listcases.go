package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type caseBrief struct {
	ID      int            `json:"id"`
	Name    string         `json:"name"`
	Gender  int            `json:"gender"`
	Age     int            `json:"age"`
	Contact sql.NullString `json:"contact"`
	Memo    sql.NullString `json:"memo"`
	Summary sql.NullString `json:"summary"`
	Date    time.Time      `json:"date"`
	Updated time.Time      `json:"updated"`
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
	err := cf.dbx.Select(&cbs, `SELECT * FROM cases ORDER BY updated DESC,date DESC LIMIT 10`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var cl []map[string]string
	for _, cb := range cbs {
		var gender string
		switch cb.Gender {
		case 1:
			gender = "女"
		case 2:
			gender = "男"
		default:
			gender = "未知"
		}
		cl = append(cl, map[string]string{
			"id":      strconv.Itoa(cb.ID),
			"name":    cb.Name,
			"gender":  gender,
			"age":     strconv.Itoa(cb.Age),
			"summary": cb.Summary.String,
			"date":    cb.Date.Format("2006-01-02"),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cl)
}
