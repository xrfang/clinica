package main

import (
	"strconv"
	"time"

	"github.com/xrfang/go-audit"
)

type caseBrief struct {
	ID          int       `json:"id"`
	PatientID   int       `json:"patient_id"`
	PatientName string    `json:"patient_name"`
	Patient     patient   `json:"patient,omitempty"`
	Summary     string    `json:"summary"`
	Opened      time.Time `json:"opened"`
	Status      int       `json:"status"`
	Consults    []consult `json:"consult"`
	Updated     time.Time `json:"updated"`
}

func (c *caseBrief) GetConsults() (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	audit.Assert(cf.dbx.Select(&c.Consults, `SELECT * FROM consults WHERE case_id=? ORDER BY 
	    time DESC LIMIT 100`, c.ID))
	for i := 0; i < len(c.Consults); i++ {
		audit.Assert(c.Consults[i].GetRecords())
	}
	return
}

func getCases(id int) ([]caseBrief, error) {
	var cbs []caseBrief
	qry := `SELECT cases.id,patients.id AS patient_id,patients.name AS patient_name,summary,
		opened,status,updated FROM cases,patients WHERE cases.patient_id=patients.id`
	if id != 0 {
		qry += ` AND cases.id=` + strconv.Itoa(id)
	} else {
		qry += ` ORDER BY updated DESC LIMIT 100`
	}
	err := cf.dbx.Select(&cbs, qry)
	if err == nil && id != 0 {
		err = cbs[0].GetConsults()
	}
	return cbs, err
}

type consult struct {
	ID      int       `json:"id"`
	CaseID  int       `json:"case_id"`
	Mode    int       `json:"mode"`
	Time    time.Time `json:"time"`
	Status  int       `json:"status"`
	Records []record  `json:"records"`
	Updated time.Time `json:"updated"`
}

type record struct {
	ID        int       `json:"id"`
	ConsultID int       `json:"consult_id"`
	Type      int       `json:"type"`
	ClassID   int       `json:"class_id"`
	Category  string    `json:"category"`
	Caption   string    `json:"caption"`
	Details   string    `json:"details"`
	Updated   time.Time `json:"updated"`
}

func (c *consult) GetRecords() (err error) {
	err = cf.dbx.Select(&c.Records, `SELECT * FROM records WHERE consult_id=? ORDER BY type,class_id`, c.ID)
	if err != nil {
		return
	}
	for i, r := range c.Records {
		r.Category = cref.String(r.Type, r.ClassID)
		c.Records[i] = r
	}
	return
}
