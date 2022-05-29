package model

import (
	"time"
)

type JsonTime time.Time

func (j JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(j).Format("2006-01-02 15:04:05") + `"`), nil
}

// Hospital storing hospital info
type Hospital struct {
	ID                  int64  `json:"id"`
	Address             string `json:"address"`
	CellPhones          string `json:"cell_phones"`
	ChiefMedicalOfficer string `json:"chief_medical_officer"`
	OfficialName        string `json:"official_name"`
	Title               string `json:"title"`
	URL                 string `json:"url"`
}
