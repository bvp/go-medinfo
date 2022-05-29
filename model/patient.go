package model

import (
	"time"
)

// Patient storing patient info
type Patient struct {
	ID               int64     `json:"id"`
	CreatedDate      time.Time `json:"created_date"`
	LastModifiedDate time.Time `json:"last_modified_date"`
	Description      string    `json:"description"`
	HospitalID       string    `json:"hospital_id"`
	Name             string    `json:"name"`
	Pin              string    `json:"pin"`
	Version          int64     `json:"version"`
}

type Registration struct {
	Speciality string
	DoctorName string
	Room       string
	Date       time.Time
	Status     string
	Reject     string
}
