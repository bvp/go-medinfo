package medinfo

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/bvp/go-medinfo/util"
)

const (
	PayloadForm = "application/x-www-form-urlencoded"
	PayloadJson = "application/json"
	PayloadXml  = "text/xml"
)

type Data struct {
	Method  string
	URL     *url.URL
	Headers map[string]string
	Params  map[string]string
	Accept  string
	Type    string
	Payload interface{}
}

func getType(T interface{}) string {
	if reflect.TypeOf(T) != nil {
		xType := reflect.TypeOf(T).Kind()
		switch xType {
		case reflect.Struct:
			return "struct"
		case reflect.Map:
			return "map"
		case reflect.Ptr:
			iType := reflect.Indirect(reflect.ValueOf(T)).Kind()
			fmt.Printf("itype - %s\n", iType)
			switch iType {
			case reflect.Struct:
				return "struct"
			case reflect.Map:
				return "map"
			}
		}
	}
	return ""
}

func normalizeData(data *Data) (err error) {
	if data == nil {
		return errors.New("data is nil")
	}
	switch getType(data.Payload) {
	case "struct":
		if data.Type == "xml" {
			data.Type = PayloadXml
			payload, err := xml.Marshal(data.Payload)
			if err != nil {
				return err
			}
			data.Payload = payload
		} else {
			data.Type = PayloadJson
			payload, err := json.Marshal(data.Payload)
			if err != nil {
				return err
			}
			data.Payload = payload
		}
	case "map":
		data.Type = PayloadForm
		form := url.Values{}
		data.Headers = map[string]string{"Content-Type": data.Type}
		for key, value := range data.Payload.(map[string]string) {
			form.Set(key, util.EncodeWindows1251(value))
		}
		data.Payload = form.Encode()
	}

	if len(data.Method) == 0 {
		if data.Payload != nil {
			data.Method = http.MethodPost
		} else {
			data.Method = http.MethodGet
		}
	}
	if data.Params != nil {
		q := data.URL.Query()
		for key, value := range data.Params {
			q.Add(key, value)
		}
		data.URL.RawQuery = q.Encode()
	}
	return
}

func (cli *Client) httpRequest(data *Data) (response string, err error) {
	if err = normalizeData(data); err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}

	var req (*http.Request)
	if data.Payload != nil {
		payload := strings.NewReader(data.Payload.(string))
		req, err = http.NewRequest(data.Method, data.URL.String(), payload)
		if err != nil {
			return "", err
		}
	} else {
		req, err = http.NewRequest(data.Method, data.URL.String(), nil)
		if err != nil {
			return "", err
		}
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:100.0) Gecko/20100101 Firefox/100.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")

	for key, value := range data.Headers {
		req.Header.Set(key, value)
	}

	res, err := cli.httpCli.Do(req)

	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = errors.New("res.StatusCode: " + strconv.Itoa(res.StatusCode))
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	response = string(util.DecodeWindows1251(body))

	return
}
