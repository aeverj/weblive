package fingerprint

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed web_fingerprint_v3.json
var fingerContent []byte

type WebFingerPrint struct {
	Name           string            `json:"name"`
	Path           string            `json:"path"`
	RequestMethod  string            `json:"request_method"`
	RequestHeaders map[string]string `json:"request_headers"`
	RequestData    string            `json:"request_data"`
	StatusCode     int               `json:"status_code"`
	Headers        map[string]string `json:"headers"`
	Keyword        []string          `json:"keyword"`
	FaviconHash    []string          `json:"favicon_hash"`
	Priority       int               `json:"priority"`
}

type FingerHub struct {
	webFingerPrint []WebFingerPrint
}

type HttpFingerInfo struct {
	Headers    map[string][]string
	StatusCode int
	Html       string
}

var h FingerHub

func init() {
	if err := json.Unmarshal(fingerContent, &h.webFingerPrint); err != nil {
		fmt.Println("Fingerprint file parsing error")
	}
}

func DoHttpFingerPrint(hInfo HttpFingerInfo) (result []string) {
	for _, finger := range h.webFingerPrint {
		all := strings.Join(result, ", ") + ", "
		if strings.Contains(all, finger.Name+", ") {
			continue
		}

		cflag := false
		for k, v := range finger.Headers {
			if len(hInfo.Headers[strings.ToLower(k)]) > 0 {
				if strings.Contains(hInfo.Headers[strings.ToLower(k)][0], v) {
					cflag = true
					goto addFinger
				}
			}
		}
		for _, v := range finger.Keyword {
			if !strings.Contains(strings.ToLower(hInfo.Html), strings.ToLower(v)) {
				cflag = false
				break
			}
			cflag = true
		}
	addFinger:
		if cflag == true {
			result = append(result, finger.Name)
		}
	}
	return
}
