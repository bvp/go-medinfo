package model

// Doctor storing doctor info
type Doctor struct {
	ID             int64  `json:"id"`
	FullName       string `json:"full_name"`
	HospitalID     int64  `json:"hospital_id"`
	Room           int64  `json:"room"`
	SpecialityID   int64  `json:"speciality_id"`
	SpecialityDesc string `json:"speciality_desc"`
}
