package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func apiListCases(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Type", "application/json")
	id := 0
	caseID := r.URL.Query().Get("id")
	if caseID != "" {
		id, _ = strconv.Atoi(caseID)
		if id == 0 {
			http.Error(w, "case-id (if given) must be postive integer", http.StatusBadRequest)
			return
		}
	}
	cbs, err := getCases(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(cbs)
}
