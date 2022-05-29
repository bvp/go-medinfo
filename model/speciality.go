package model

import (
	"time"
)

type Speciality struct {
	ID             int64     `json:"id"`
	Title          string    `json:"title"`
	CreateDateTime time.Time `json:"create_date_time"`
	UpdateDateTime time.Time `json:"update_date_time"`
}
