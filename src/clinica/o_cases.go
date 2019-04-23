package main

import (
	"strconv"
	"time"
)

type caseBrief struct {
	ID          int       `json:"id"`
	PatientID   int       `json:"patient_id"`
	PatientName string    `json:"patient_name"`
	Patient     patient   `json:"patient,omitempty"`
	Summary     string    `json:"summary"`
	Opened      time.Time `json:"opened"`
	Status      int       `json:"status"`
	Updated     time.Time `json:"updated"`
}

func getCases(id int) ([]caseBrief, error) {
	var cbs []caseBrief
	qry := `SELECT cases.id,patients.id AS patientid,patients.name AS patientname,summary,
		opened,status,updated FROM cases,patients WHERE cases.patient_id=patients.id`
	if id != 0 {
		qry += ` AND cases.id=` + strconv.Itoa(id)
	} else {
		qry += ` ORDER BY updated DESC LIMIT 100`
	}
	err := cf.dbx.Select(&cbs, qry)
	return cbs, err
}

type consult struct {
	ID      int       `json:"id"`
	Case_ID int       `json:"case_id"`
	Mode    int       `json:"mode"`
	Time    time.Time `json:"time"`
	Status  int       `json:"status"`
	Updated time.Time `json:"updated"`
	Details []details `json:"details"`
}

func getConsults(id int) ([]consult, error) {
	var cs []consult
	err := cf.dbx.Select(&cs, `SELECT * FROM consults WHERE case_id=? ORDER BY updated DESC LIMIT 100`, id)
	return cs, err
}

type details struct {
}
