package main

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/xrfang/go-audit"
)

func users(w http.ResponseWriter, r *http.Request) {
	s := sessions.Get(w, r)
	var u user
	if s.Unmarshal("user", &u) != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	if u.Role != RoleAdmin {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	switch r.Method {
	case "GET":
		var us []user
		audit.Assert(cf.dbx.Select(&us, "SELECT name,login,role FROM users ORDER BY role DESC,login"))
		var ul []map[string]string
		for _, u := range us {
			var style string
			switch u.Role {
			case RoleDisabled:
				style = "secondary"
			case RoleReader:
				style = "info"
			case RoleEditor:
				style = "primary"
			case RoleAdmin:
				style = "danger"
			}
			ul = append(ul, map[string]string{
				"caption": u.Caption(),
				"name":    u.Name.String,
				"login":   u.Login,
				"role":    strconv.Itoa(u.Role),
				"style":   style,
			})
		}
		renderTemplate(w, "users.html", struct {
			Admin string
			Users []map[string]string
		}{
			Admin: u.Caption(),
			Users: ul,
		})
	case "POST":
		login := strings.TrimSpace(r.FormValue("login"))
		if login == "" {
			http.Error(w, "登录名(login)不能为空", http.StatusBadRequest)
			return
		}
		rl := regexp.MustCompile(`(?i)^[a-z0-9]{1,16}$`)
		if !rl.MatchString(login) {
			http.Error(w, "登录名(login)格式不符合要求", http.StatusBadRequest)
			return
		}
		role, err := strconv.Atoi(r.FormValue("role"))
		if err != nil || role < -1 || role > 2 {
			http.Error(w, "权限设置(role)错误", http.StatusBadRequest)
			return
		}
		keys := []string{"login", "role"}
		vals := []interface{}{login, role}
		upds := `) ON CONFLICT(login) DO UPDATE SET role=?`
		name := strings.TrimSpace(r.FormValue("name"))
		if name != "" {
			if len(name) > 32 {
				http.Error(w, "姓名长度超过限制", http.StatusBadRequest)
				return
			}
			keys = append(keys, "name")
			vals = append(vals, name)
			upds += `,name=?`
		}
		passwd := r.FormValue("passwd")
		if passwd != "" {
			keys = append(keys, "passwd")
			vals = append(vals, HashPassword(passwd))
			upds += `,passwd=?`
		}
		vals = append(vals, vals[1:]...)
		cmd := `INSERT INTO users (` + strings.Join(keys, ",") + `) VALUES (?` +
			strings.Repeat(",?", len(keys)-1) + upds
		_, err = cf.dbx.Exec(cmd, vals...)
		if err != nil {
			http.Error(w, "内部错误: "+err.Error(), http.StatusInternalServerError)
		}
	case "DELETE":
		login := r.URL.Query().Get("login")
		if login == u.Login {
			http.Error(w, "不能删除当前用户", http.StatusForbidden)
			return
		}
		cf.dbx.MustExec(`DELETE FROM users WHERE login=?`, login)
	}
}
