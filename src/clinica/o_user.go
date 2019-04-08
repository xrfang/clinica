package main

import (
	"database/sql"

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
