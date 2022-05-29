package config

import (
	"fmt"
)

var (
	UrlBase         = "https://www.medinfo-yar.ru"
	UrlApp          = fmt.Sprintf("%s/app/m2_ajax.php", UrlBase)
	UrlRasp         = fmt.Sprintf("%s/index.php/shed", UrlBase)
	UrlRaspDay      = fmt.Sprintf("%s/index.php/req_day", UrlBase)
	UrlReqPac       = fmt.Sprintf("%s/index.php/req_pac", UrlBase)
	UrlSchedule     = fmt.Sprintf("%s/index.php/shed?lu_num=", UrlBase)    // lu_num is id of hospital
	UrlHospitalCard = fmt.Sprintf("%s/index.php/lu_info?lu_num=", UrlBase) // lu_num is id of hospital
	//UrlScheduleXy   = fmt.Sprintf("&x=%d&y=%d", rand.Intn(24), rand.Intn(24)) // x and y is a click coordinates
)
