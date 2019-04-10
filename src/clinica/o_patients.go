package main

import (
	"database/sql"
	"encoding/json"
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
	Updated  sql.NullString
}

func (p patient) String() string {
	buf, _ := json.Marshal(map[string]interface{}{
		"id":       p.ID,
		"name":     p.Name,
		"gender":   p.Gender,
		"birthday": p.Birthday.String,
		"contact":  p.Contact.String,
		"memo":     p.Memo.String,
		"updated":  p.Updated.String,
	})
	return string(buf)
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

func getPatients(term string) (int, []patient) {
	var vals []interface{}
	qry := "SELECT patients.*,MAX(cases.updated) AS updated FROM patients LEFT JOIN cases ON cases.patient_id=patients.id"
	cnt := "SELECT COUNT(*) FROM patients"
	if term != "" {
		term = "%" + term + "%"
		qry += " WHERE name LIKE ? OR contact LIKE ? OR memo LIKE ?"
		vals = []interface{}{term, term, term}
	}
	qry = fmt.Sprintf("%s GROUP BY patients.id ORDER BY updated DESC,name LIMIT 100", qry)
	var c int
	audit.Assert(cf.dbx.Get(&c, cnt))
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
	id := arg("id")
	var cmd string
	if id == "" || id == "0" { //INSERT
		cmd = fmt.Sprintf(`INSERT INTO patients (%s) VALUES (%s)`, strings.Join(keys, ","), strings.Repeat("?", len(vals)))
	} else { //UPDATE
		var sets []string
		for _, k := range keys {
			sets = append(sets, k+"=?")
		}
		vals = append(vals, id)
		cmd = fmt.Sprintf(`UPDATE patients SET %s WHERE id=?`, strings.Join(sets, ","))
	}
	fmt.Println(cmd)
	fmt.Println(vals)
	//TODO: exec sql
	return http.StatusOK, "OK"
}

func delPatient(id string) (int, string) {
	_, err := cf.dbx.Exec(`DELETE FROM patients WHERE id=?`, id)
	if err == nil {
		return http.StatusOK, "OK"
	}
	return http.StatusInternalServerError, "内部错误: " + err.Error()
}
