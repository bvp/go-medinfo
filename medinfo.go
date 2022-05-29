package medinfo

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bvp/go-medinfo/config"
	"github.com/bvp/go-medinfo/model"
	"github.com/bvp/go-medinfo/util"
)

type Client struct {
	LName   string
	FName   string
	MName   string
	BDate   string
	Snils   string
	httpCli *http.Client
}

func NewClient(lname string, fname string, mname string, bdate string, snils string, httpCli *http.Client) *Client {
	if httpCli == nil {
		httpCli = util.GetHttpClient()
	}
	cli := &Client{
		LName:   lname,
		FName:   fname,
		MName:   mname,
		BDate:   bdate,
		Snils:   snils,
		httpCli: httpCli,
	}
	return cli
}

func (cli *Client) GetHospitals() {
	hospitals := []model.Hospital{}

	surl, _ := url.Parse(config.UrlRasp)
	resp, _ := cli.httpRequest(&Data{
		URL:     surl,
		Method:  "GET",
		Payload: nil,
	})

	if doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp)); err == nil {
		doc.Find("select#lu_ch_select > optgroup option").Each(func(i int, s *goquery.Selection) {
			title := strings.TrimSpace(s.Text())
			value, _ := s.Attr("value")
			id, _ := strconv.ParseInt(value, 10, 64)
			h := model.Hospital{}
			h.ID = id
			h.Title = title
			hospitals = append(hospitals, h)
		})
	} else {
		fmt.Printf("Can't fetch hospitals - %s\n", err.Error())
	}
	for _, h := range hospitals {
		fmt.Println(h)
	}
}

func (cli *Client) GetSpecialities() {
	specialities := []model.Speciality{}

	surl, _ := url.Parse(config.UrlRasp)
	resp, _ := cli.httpRequest(&Data{
		URL:     surl,
		Method:  "GET",
		Payload: nil,
	})

	if doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp)); err == nil {
		doc.Find("select.PopupDoct > option").Each(func(i int, s *goquery.Selection) {
			title := strings.TrimSpace(s.Text())
			value, _ := s.Attr("value")
			id, _ := strconv.ParseInt(value, 10, 64)
			spec := model.Speciality{}
			spec.ID = id
			spec.Title = title
			specialities = append(specialities, spec)
		})
	} else {
		fmt.Printf("Can't fetch hospitals - %s\n", err.Error())
	}
	for _, s := range specialities {
		fmt.Println(s)
	}
}

func (cli *Client) GetRegistrations() string {
	surl, _ := url.Parse(config.UrlReqPac)
	payload := map[string]string{
		"egis3_go_type":   "0",
		"no_yandex":       "09051945",
		"cupon_id":        "0",
		"k_egis_row":      "1",
		"pc_apnt_str":     "",
		"rec_pac_type":    "1",
		"pc_apnt_num":     "-1945",
		"req_request":     "1",
		"panel_view_type": "1",
		"lu_num":          "400",
		"reg_type":        "2",
		"pc_lname":        cli.LName,
		"pc_fname":        cli.FName,
		"pc_mname":        cli.MName,
		"pc_date":         cli.BDate,
		"pc_snils":        cli.Snils,
		"pc_amb":          "0",
		"req_out_type":    "1",
		"k_dsave":         "0",
	}

	resp, _ := cli.httpRequest(&Data{
		URL:     surl,
		Method:  "POST",
		Payload: payload,
	})

	if doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp)); err == nil {
		doc.Find("table#tbl_obj > tbody > tr").Each(func(i int, s *goquery.Selection) {
			// html, _ := s.Html()
			// fmt.Printf("%s\n", html)
			speciality := s.Find("td:nth-child(4)").Text()
			doctor := strings.TrimSpace(s.Find("td:nth-child(5)").Text())
			room := strings.TrimSpace(s.Find("td:nth-child(6)").Text())
			reg := regexp.MustCompile("[^0-9]+")
			_date := reg.ReplaceAllString(s.Find("td:nth-child(7)").Text(), "")
			date, _ := time.Parse("020120061504", _date)
			status := s.Find("td:nth-child(8) > span:nth-child(1)").Text()
			additional := s.Find("td:nth-child(8) > span:nth-child(3)").Text()
			created := s.Find("td:nth-child(9)").Text()
			rawReject := s.Find("td:nth-child(10) > span").AttrOr("onclick", "")
			reject := strings.TrimSuffix(strings.TrimPrefix(rawReject, "js_online_reject_go("), ");")
			fmt.Printf("%s, %s, %s, %s, %s (%s) [%s] - %s\n", speciality, doctor, room, date.Format("2006-01-02 15:04"), status, additional, created, reject)
		})
	}
	return resp
}

func (cli *Client) checkPagesCount(_url string) (pagesCount int) {
	surl, _ := url.Parse(_url)
	resp, _ := cli.httpRequest(&Data{
		URL:     surl,
		Method:  "GETT",
		Payload: nil,
	})
	if doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp)); err == nil {
		headL := doc.Find("td.tdHeadL").First()
		headLChildren := headL.Children()
		if headLChildren.Size() > 1 {
			if headLChildren.Find("td").Is(".pp_td_Nav") {
				pagesCount = headLChildren.Find("a.pp_NavRef").Size() + 1
			} else {
				pagesCount = 1
			}
		} else {
			pagesCount = 1
		}
	}
	return
}

func (cli *Client) GetEgisUpdate(schedId int) {
	params := url.Values{
		"q_type":  {"10"},
		"shed_id": {fmt.Sprintf("%d", schedId)},
	}
	surl, _ := url.Parse(config.UrlApp + params.Encode())
	/* resp, _ :=  */ cli.httpRequest(&Data{
		URL:     surl,
		Method:  "GET",
		Payload: nil,
	})
	// fmt.Printf("egis update resp - %s\n", resp)
}

func (cli *Client) GetFreeReqFor(luNum int, schedId int) (freeReqs []model.FreeReq) {
	cli.GetEgisUpdate(schedId)
	surl, _ := url.Parse(fmt.Sprintf("%s?lu_num=%d&shed_idx=%d", config.UrlRaspDay, luNum, schedId))
	resp, _ := cli.httpRequest(&Data{
		URL:     surl,
		Method:  "GET",
		Payload: nil,
	})
	if doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp)); err == nil {
		dataTable := doc.Find(".table_grid")
		dataTable.Find("tbody > tr").Each(func(i int, row *goquery.Selection) {
			link := row.Find("td:nth-child(1) > span > a").AttrOr("href", "")
			stime := strings.TrimSpace(row.Find("td:nth-child(1)").Text())
			stitle := row.Find("td:nth-child(2)").Text()
			if stitle == "Явка свободна" {
				cuponId := strings.ReplaceAll(link, fmt.Sprintf("/index.php/req_form?lu_num=%d&cupon_id=", luNum), "")
				freeReq := model.FreeReq{LuNum: luNum, Time: stime, Text: stitle, CuponId: cuponId}
				freeReqs = append(freeReqs, freeReq)
			}
		})
	}
	return
}

// GetSchedules getting schedule information for specific data
func (cli *Client) GetSchedules(luNum int, dt string, spec int, fio string) (doctors []string, schedules []model.Schedule) {
	efio := ""
	if fio != "" {
		efio = url.QueryEscape(string(util.EncodeWindows1251(fio)))
	}

	_url := fmt.Sprintf("%s%d&shed_dt=%s&spec_num=%d&x=%d&y=%d&fio=%s&sort=1", config.UrlSchedule, luNum, dt, spec, rand.Intn(47), rand.Intn(47), efio)

	pagesCount := cli.checkPagesCount(_url)
	for pn := 1; pn <= pagesCount; pn++ {
		surl, _ := url.Parse(fmt.Sprintf("%s&pn=%d", _url, pn))
		resp, _ := cli.httpRequest(&Data{
			URL:     surl,
			Method:  "GET",
			Payload: nil,
		})
		if doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp)); err == nil {
			dataTable := doc.Find("body > table > tbody > tr:nth-child(3) > td:nth-child(1) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(3) > table > tbody > tr > td:nth-child(2) > table:nth-child(2)") // Headers
			rows := dataTable.Find("tbody > tr")                                                                                                                                                                                        // Headers

			rows0 := rows.First().Find("td") // Headers 1 - Дни недели
			headers1 := make([]string, rows0.Size())
			rows0.Each(func(i int, s *goquery.Selection) {
				value := s.Text()
				headers1[i] = value
			})

			rows1 := rows.First().Next().Find("td") // Headers 2 - Даты
			headers2 := make([]string, rows1.Size())
			rows1.Each(func(i int, s *goquery.Selection) {
				value := s.Text()
				headers2[i] = value
			})

			for i := 2; i < rows.Size(); i++ {
				row := rows.Eq(i)
				cols := row.Find("td")
				if cols.Size() == 10 {
					d := cols.Eq(0).Text()
					doctors = append(doctors, d)
					s, a := util.GetOwnText(cols.Eq(1), "a")
					if a != "" {
						s = fmt.Sprintf("%s %s", s, a)
					}
					r := cols.Eq(2).Text()

					for day := 3; day < 9; day++ {
						sch := &model.Schedule{}
						sch.HospitalID = luNum
						sch.SpecID = spec
						sch.SpecDesc = s
						sch.DoctorFullName = d
						sch.Room = r
						sch.Date, _ = time.Parse("02-01-2006", headers2[day-3])
						sch.Info = strings.TrimSpace(cols.Eq(day).AttrOr("title", ""))

						var re = "0"
						stime, removed := util.GetOwnText(cols.Eq(day), "span")

						_schLink := strings.TrimRight(strings.ReplaceAll(cols.Eq(day).Find("a").AttrOr("href", ""), "Javascript:td_sel_click(", ""), ")")
						schlData := strings.Split(_schLink, ",")
						if len(schlData) == 3 {
							sch.ReqURL = fmt.Sprintf("%s?lu_num=%s&shed_idx=%s", config.UrlRaspDay, schlData[1], schlData[2])
							schedId, _ := strconv.Atoi(schlData[2])
							sch.FreeReqs = cli.GetFreeReqFor(luNum, schedId)
						}

						if stime != "" {
							_stime := strings.Split(strings.ReplaceAll(stime, " - ", "-"), "-")
							if len(_stime) > 0 {
								st, et := _stime[0], _stime[1]
								pst, _ := time.Parse("15:04", st)
								pet, _ := time.Parse("15:04", et)
								sch.TimeStart = sch.Date.Add(time.Hour*time.Duration(pst.Hour()) + time.Minute*time.Duration(pst.Minute()))
								sch.TimeEnd = sch.Date.Add(time.Hour*time.Duration(pet.Hour()) + time.Minute*time.Duration(pet.Minute()))
							}
							if removed != "" {
								re = removed
							}

							sch.Req, _ = strconv.Atoi(re)
						}

						schedules = append(schedules, *sch)
					}
				}
			}
		} else {
			fmt.Printf("ERROR GetSchedulesFor - %s\n", err.Error())
		}
	}
	return
}

func (cli *Client) GetHospitalInfo(luNum int64) (h model.Hospital) {
	surl, _ := url.Parse(fmt.Sprintf("%s%d&ptab_id=-1", config.UrlHospitalCard, luNum))
	resp, _ := cli.httpRequest(&Data{
		URL:     surl,
		Method:  "GET",
		Payload: nil,
	})
	if doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp)); err == nil {
		doc.Find("body > table > tbody > tr:nth-child(3) > td > table > tbody > tr > td:nth-child(2) > table:nth-child(3) > tbody > tr > td:nth-child(1) > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
			switch strings.TrimSpace(s.Find("td:nth-child(1).tdContR").Text()) {
			case "Официальное наименование":
				h.OfficialName = s.Find("td:nth-child(2).tdContRD").Text()
			case "Адрес":
				h.Address = s.Find("td:nth-child(2).tdContRD").Text()
			case "Контактные телефоны":
				h.CellPhones = s.Find("td:nth-child(2).tdContRD").Text()
			case "Главный врач":
				h.ChiefMedicalOfficer = s.Find("td:nth-child(2).tdContRD").Text()
			case "Официальный сайт":
				h.URL = s.Find("td:nth-child(2).tdContRD").Text()
			}
		})
	}
	return
}
