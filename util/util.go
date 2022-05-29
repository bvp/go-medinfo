// util
package util

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/charmap"
)

var (
	Log *log.Logger
)

func GetHttpClient() *http.Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Jar: jar,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
		},
		Timeout: 1 * time.Minute,
	}

	return client
}

func NewLogToStdout() {
	Log = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}

func NewLogToFile(logpath string) {
	file, err := os.Create(logpath)
	if err != nil {
		panic(err)
	}
	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
}

func GetOwnText(qs *goquery.Selection, removeSelector string) (result string, removed string) {
	t := strings.TrimSpace(qs.Text())
	removed = strings.TrimSpace(qs.Children().Find(removeSelector).Text())
	result = strings.TrimSpace(strings.TrimRight(t, removed))
	return
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func DecodeWindows1251(ba []uint8) []uint8 {
	dec := charmap.Windows1251.NewDecoder()
	out, _ := dec.Bytes(ba)
	return out
}

func EncodeWindows1251(ba string) string {
	enc := charmap.Windows1251.NewEncoder()
	out, _ := enc.String(ba)
	return out
}

func UniqStrings(s []string) []string {
	uniq := make(map[string]bool, len(s))
	uniqSlice := make([]string, len(uniq))
	for _, elem := range s {
		if !uniq[elem] {
			uniqSlice = append(uniqSlice, elem)
			uniq[elem] = true
		}
	}
	return uniqSlice
}
