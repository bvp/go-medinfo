package model

import (
	"fmt"
	"time"
)

// Schedule storing schedule info
type Schedule struct {
	HospitalID     int       `json:"hospital_id"`
	DoctorFullName string    `json:"doctor_full_name"`
	SpecID         int       `json:"spec_id"`
	SpecDesc       string    `json:"spec_desc"`
	Room           string    `json:"room"`
	Info           string    `json:"info"`
	Req            int       `json:"req"`
	ReqURL         string    `json:"req_url"`
	FreeReqs       []FreeReq `json:"freeReqs`
	Date           time.Time `json:"date"`
	TimeStart      time.Time `json:"time_start"`
	TimeEnd        time.Time `json:"time_end"`
}

func (sch Schedule) String() string {
	return fmt.Sprintf("HospitalID: %d, DoctorFullName: %s, SpecID: %d, Spec: %s, Room: %s, Info: %s, Req: %d, ReqURL: %s, Date: %s, TimeStart: %s, TimeEnd: %s",
		sch.HospitalID,
		sch.DoctorFullName,
		sch.SpecID,
		sch.SpecDesc,
		sch.Room,
		sch.Info,
		sch.Req,
		sch.ReqURL,
		sch.Date.Format("2006-01-02"),
		sch.TimeStart.Format("2006-01-02 15:04"),
		sch.TimeEnd.Format("2006-01-02 15:04"),
	)
}

type FreeReq struct {
	LuNum   int
	Time    string
	Text    string
	CuponId string
}
