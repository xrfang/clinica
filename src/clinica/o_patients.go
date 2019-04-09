package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	audit "github.com/xrfang/go-audit"
)

type patient struct {
	ID       int
	Name     string
	Gender   int
	Birthday sql.NullString
	Contact  sql.NullString
	Memo     sql.NullString
}

func (p patient) Caption() string {
	tag := "女"
	if p.Gender != 0 {
		tag = "男"
	}
	if len(p.Birthday.String) >= 4 {
		by, _ := strconv.Atoi(p.Birthday.String[:4])
		if by > 0 {
			age := time.Now().Year() - by
			tag = fmt.Sprintf("%s %d岁", tag, age)
		}
	}
	return p.Name + " (" + tag + ")"
}

func getPatients(args map[string]string, limit, offset int) (int, []patient) {
	var cond []string
	var vals []interface{}
	for k, v := range args {
		switch k {
		case "id":
			id, _ := strconv.Atoi(v)
			if id <= 0 {
				return 0, nil
			}
			cond = append(cond, fmt.Sprintf("id=%d", id))
		case "gender":
			g, _ := strconv.Atoi(v)
			if g < 0 || g > 1 {
				return 0, nil
			}
			cond = append(cond, fmt.Sprint("gender=%d", g))
		case "birthday":
			cond = append(cond, "birthday LIKE ?%")
			vals = append(vals, v)
		case "name", "contact", "memo":
			cond = append(cond, k+" LIKE %?%")
			vals = append(vals, v)
		default:
			panic(fmt.Errorf("invalid argument: %s", k))
		}
	}
	qry := "SELECT * FROM patients"
	cnt := "SELECT COUNT(*) FROM patients"
	if len(cond) > 0 {
		where := " WHERE " + strings.Join(cond, " AND ")
		qry += where
		cnt += where
	}
	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	qry = fmt.Sprintf("%s ORDER BY name,contact LIMIT %d OFFSET %d", qry, limit, offset)
	var c int
	audit.Assert(cf.dbx.Get(&c, cnt, vals...))
	var ps []patient
	audit.Assert(cf.dbx.Select(&ps, qry, vals...))
	return c, ps
}

func setPatient(args url.Values) (int, string) {
	arg := func(name string) string {
		return strings.TrimSpace(args.Get(name))
	}
	var keys []string
	var vals []interface{}
	name := arg("name")
	if len(name) < 1 || len(name) > 32 {
		return http.StatusBadRequest, "姓名不能为空或超过10个汉字"
	}
	keys = append(keys, "name")
	vals = append(vals, name)
	gender := arg("gender")
	if gender != "0" && gender != "1" {
		return http.StatusBadRequest, "性别设置不正确"
	}
	keys = append(keys, "gender")
	vals = append(vals, gender)
	birthday := arg("birthday")
	if len(birthday) > 10 {
		return http.StatusBadRequest, "生日设置不正确"
	}
	keys = append(keys, "birthday")
	vals = append(vals, birthday)
	contact := arg("contact")
	if len(contact) > 64 {
		return http.StatusBadRequest, "联系方式过长"
	}
	keys = append(keys, "contact")
	vals = append(vals, contact)
	memo := arg("memo")
	if len(memo) > 1024 {
		return http.StatusBadRequest, "备注过长"
	}
	keys = append(keys, "memo")
	vals = append(vals, memo)
	//TODO 判断insert/update
	//TODO: exec sql
	return http.StatusOK, "OK"
}

func delPatient(id int) (int, string) {
	_, err := cf.dbx.Exec(`DELETE FROM patients WHERE id=?`, id)
	if err == nil {
		return http.StatusOK, "OK"
	}
	return http.StatusInternalServerError, "内部错误: " + err.Error()
}
