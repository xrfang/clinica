package main

import (
	"net/http"
	"strconv"
	"time"
)

func apiEditConsultRecord(w http.ResponseWriter, r *http.Request) {
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
	_type, err := strconv.Atoi(r.PostFormValue("type"))
	if err != nil || _type < 0 || _type > 4 {
		http.Error(w, "invalid [type]", http.StatusBadRequest)
		return
	}
	classID, err := strconv.Atoi(r.PostFormValue("class_id"))
	if err != nil || classID < 0 {
		http.Error(w, "invalid [class_id]", http.StatusBadRequest)
		return
	}
	if _type != 1 {
		classID = 0
	}
	id := r.PostFormValue("id")
	consultID := r.PostFormValue("consult_id")
	caption := r.PostFormValue("caption")
	details := r.PostFormValue("details")
	updated := time.Now().Format(time.RFC3339)

	if id == "" {
		_, err = cf.dbx.Exec(`INSERT INTO records (consult_id,type,class_id,caption,details,updated) 
            VALUES (?,?,?,?,?,?)`, consultID, _type, classID, caption, details, updated)
	} else {
		_, err = cf.dbx.Exec(`UPDATE records SET caption=?,details=?,updated=? WHERE id=?`,
			caption, details, updated, id)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
