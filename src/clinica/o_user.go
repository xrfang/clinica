package main

import (
	"database/sql"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	audit "github.com/xrfang/go-audit"
)

type user struct {
	Login  string
	Passwd sql.NullString
	Name   sql.NullString
	Role   int
}

func (u user) Caption() string {
	if u.Name.String != "" {
		return u.Login + " (" + u.Name.String + ")"
	}
	return u.Login
}

func listUsers() []user {
	var us []user
	audit.Assert(cf.dbx.Select(&us, "SELECT * FROM users ORDER BY role DESC,login"))
	return us
}

func getUser(login, passwd string) *user {
	var u user
	err := cf.dbx.Get(&u, "SELECT * FROM users WHERE login=?", login)
	if err != nil || u.Role == RoleDisabled {
		return nil
	}
	if CheckPasswordHash(passwd, u.Passwd.String) {
		return &u
	}
	return nil
}

func setUser(args url.Values) (int, string) {
	arg := func(name string) string {
		return strings.TrimSpace(args.Get(name))
	}
	login := arg("login")
	if login == "" {
		return http.StatusBadRequest, "登录名(login)不能为空"
	}
	rl := regexp.MustCompile(`(?i)^[a-z0-9]{1,16}$`)
	if !rl.MatchString(login) {
		return http.StatusBadRequest, "登录名(login)格式不符合要求"
	}
	role, err := strconv.Atoi(arg("role"))
	if err != nil || role < -1 || role > 2 {
		return http.StatusBadRequest, "权限设置(role)错误"
	}
	keys := []string{"login", "role"}
	vals := []interface{}{login, role}
	upds := `) ON CONFLICT(login) DO UPDATE SET role=?`
	name := arg("name")
	if name != "" {
		if len(name) > 32 {
			return http.StatusBadRequest, "姓名长度超过限制"
		}
		keys = append(keys, "name")
		vals = append(vals, name)
		upds += `,name=?`
	}
	passwd := args.Get("passwd")
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
		return http.StatusInternalServerError, "内部错误: " + err.Error()
	}
	return http.StatusOK, "OK"
}

func delUser(login string) (int, string) {
	_, err := cf.dbx.Exec(`DELETE FROM users WHERE login=?`, login)
	if err == nil {
		return http.StatusOK, "OK"
	}
	return http.StatusInternalServerError, "内部错误: " + err.Error()
}
