package main

import (
	"time"
)

type caseBrief struct {
	ID      int       `json:"id"`
	Patient string    `json:"patient"`
	Summary string    `json:"summary"`
	Opened  time.Time `json:"opened"`
	Status  int       `json:"status"`
	Updated time.Time `json:"updated"`
}

func getCases() ([]caseBrief, error) {
	var cbs []caseBrief
	err := cf.dbx.Select(&cbs, `SELECT cases.id,patients.name AS patient,summary,opened,status,updated FROM
		cases,patients WHERE cases.patient_id=patients.id ORDER BY updated DESC LIMIT 100`)
	return cbs, err
}

type consult struct {
	ID      int       `json:"id"`
	CaseID  int       `json:"case_id"`
	Mode    int       `json:"mode"`
	Time    time.Time `json:"time"`
	Status  int       `json:"status"`
	Updated time.Time `json:"updated"`
	Details []details `json:"details"`
}

type details struct {
}
