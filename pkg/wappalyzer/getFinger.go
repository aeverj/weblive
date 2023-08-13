package wappalyzer

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"weblive/pkg/cdn"
	"weblive/pkg/weblive"
)

// 以字符串形式嵌入 assets/hello.txt
//
//go:embed apps.json
var apps []byte

type application struct {
	Name       string   `json:"name,ompitempty"`
	Version    string   `json:"version"`
	Categories []string `json:"categories,omitempty"`

	Cats     []int                  `json:"cats,omitempty"`
	Cookies  interface{}            `json:"cookies,omitempty"`
	Js       interface{}            `json:"js,omitempty"`
	Headers  interface{}            `json:"headers,omitempty"`
	HTML     interface{}            `json:"html,omitempty"`
	Excludes interface{}            `json:"excludes,omitempty"`
	Implies  interface{}            `json:"implies,omitempty"`
	Meta     map[string]interface{} `json:"meta,omitempty"`
	Scripts  interface{}            `json:"script,omitempty"`
	URL      string                 `json:"url,omitempty"`
	Website  string                 `json:"website,omitempty"`
}

type category struct {
	Name     string `json:"name,omitempty"`
	Priority int    `json:"priority,omitempty"`
}

type temp struct {
	Apps       map[string]*json.RawMessage `json:"apps"`
	Categories map[string]*json.RawMessage `json:"categories"`
}

type wappalyzer struct {
	Apps       map[string]*application
	Categories map[string]*category
}

type pattern struct {
	str        string
	regex      *regexp.Regexp
	version    string
	confidence string
}

func appendResult(app *application, result []string) []string {
	imp := parseImplies(app.Implies)
	for _, v := range imp {
		if !weblive.IsContainStr(v, result) {
			result = append(result, v)
		}
	}
	if !weblive.IsContainStr(app.Name, result) {
		result = append(result, app.Name)
	}
	return result
}

func parsePatterns(patterns interface{}) (result map[string][]*pattern) {
	parsed := make(map[string][]string)
	switch ptrn := patterns.(type) {
	case string:
		parsed["main"] = append(parsed["main"], ptrn)
	case map[string]interface{}:
		for k, v := range ptrn {
			parsed[k] = append(parsed[k], v.(string))
		}
	case []interface{}:
		var slice []string
		for _, v := range ptrn {
			slice = append(slice, v.(string))
		}
		parsed["main"] = slice
	default:
		return nil
	}
	result = make(map[string][]*pattern)
	for k, v := range parsed {
		for _, str := range v {
			appPattern := &pattern{}
			slice := strings.Split(str, "\\;")
			for i, item := range slice {
				if item == "" {
					continue
				}
				if i > 0 {
					additional := strings.Split(item, ":")
					if len(additional) > 1 {
						if additional[0] == "version" {
							appPattern.version = additional[1]
						} else {
							appPattern.confidence = additional[1]
						}
					}
				} else {
					appPattern.str = item
					first := strings.Replace(item, `\/`, `/`, -1)
					second := strings.Replace(first, `\\`, `\`, -1)
					reg, err := regexp.Compile(fmt.Sprintf("%s%s", "(?i)", strings.Replace(second, `/`, `\/`, -1)))
					if err == nil {
						appPattern.regex = reg
					}
				}
			}
			result[k] = append(result[k], appPattern)
		}
	}
	return result
}

func parseImplies(imp interface{}) []string {
	var result []string
	switch Implies := imp.(type) {
	case string:
		result = append(result, Implies)
	case []interface{}:
		for _, v := range Implies {
			result = append(result, v.(string))
		}
	default:
		result = nil
	}
	return result
}

func hasapp(cdata *weblive.CollyData, wapp *wappalyzer) []string {
	var result []string
	for _, app := range wapp.Apps {
		if app.HTML != nil {
			patterns := parsePatterns(app.HTML)
			for _, v := range patterns {
				for _, patten := range v {
					if patten.regex != nil && patten.regex.Match([]byte(cdata.Html)) {
						result = appendResult(app, result)
					}
				}

			}
		}
		if app.Cookies != nil {
			patterns := parsePatterns(app.Cookies)
			for cookieName, _ := range patterns {
				cookieNameLowerCase := strings.ToLower(cookieName)
				if _, ok := cdata.Cookies[cookieNameLowerCase]; ok {
					result = appendResult(app, result)
				}
			}
		}
		if app.Headers != nil {
			patterns := parsePatterns(app.Headers)
			for headerName, v := range patterns {
				headerNameLowerCase := strings.ToLower(headerName)
				for _, pattrn := range v {
					if headersSlice, ok := cdata.Headers[headerNameLowerCase]; ok {
						for _, header := range headersSlice {
							if pattrn.regex != nil && pattrn.regex.Match([]byte(header)) {
								result = appendResult(app, result)
							} else if pattrn.regex == nil {
								result = appendResult(app, result)
							}
						}
					}
				}
			}
		}
		if app.Scripts != nil {
			patterns := parsePatterns(app.Scripts)
			for _, v := range patterns {
				for _, pattrn := range v {
					if pattrn.regex != nil {
						for _, script := range cdata.Scripts {
							if pattrn.regex.Match([]byte(script)) {
								result = appendResult(app, result)
							}
						}
					}
				}
			}
		}
		if app.Meta != nil {
			patterns := parsePatterns(app.Meta)
			for metaName, v := range patterns {
				metaNameLowerCase := strings.ToLower(metaName)
				for _, patten := range v {
					if value, ok := cdata.Meta[metaNameLowerCase]; ok && patten.regex.Match([]byte(value)) {
						result = appendResult(app, result)
					}
				}

			}
		}
	}
	return result
}

func Init() *wappalyzer {
	temporary := &temp{}
	err := json.Unmarshal(apps, &temporary)
	if err != nil {
		log.Fatalln("file `apps.json` not found!")
	}
	wapp := &wappalyzer{}
	wapp.Apps = make(map[string]*application)
	wapp.Categories = make(map[string]*category)
	for k, v := range temporary.Apps {
		app := &application{}
		app.Name = k
		if err = json.Unmarshal(*v, app); err != nil {
			log.Println(err)
		}
		wapp.Apps[k] = app
	}
	for k, v := range temporary.Categories {
		catg := &category{}
		if err = json.Unmarshal(*v, catg); err != nil {
			log.Println(err)
		}
		wapp.Categories[k] = catg
	}
	return wapp
}

func (wapp *wappalyzer) Analyze(cData *weblive.CollyData) (result *weblive.Website) {

	website := &weblive.Website{}
	website.Url = cData.Url
	website.StatusCode = cData.StatusCode
	website.Redirect = cData.Redirect
	website.Title = cData.Title
	website.Finger = hasapp(cData, wapp)
	website.Cdn, website.Ip, _ = cdn.Resolve(strings.Split(cData.Url.Host, ":")[0])
	result = website
	return
}
