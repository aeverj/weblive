package runner

import (
	"bytes"
	"crypto/tls"
	"github.com/antchfx/htmlquery"
	"github.com/x/x/pkg/weblive"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type SimpleRunner struct {
	client *http.Client
	config *weblive.Config
}

func NewSimpleRunner(conf *weblive.Config) *SimpleRunner {
	var simplerunner SimpleRunner
	proxyURL := http.ProxyFromEnvironment
	customProxy := conf.ProxyURL

	if len(customProxy) > 0 {
		pu, err := url.Parse(customProxy)
		if err == nil {
			proxyURL = http.ProxyURL(pu)
		}
	}

	simplerunner.config = conf
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Panicln("cookie 设置错误")
	}
	simplerunner.client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return http.ErrUseLastResponse
			}
			return nil
		},
		Jar:     jar,
		Timeout: time.Duration(time.Duration(conf.Timeout) * time.Second),
		Transport: &http.Transport{
			Proxy:               proxyURL,
			MaxIdleConns:        1000,
			MaxIdleConnsPerHost: 500,
			MaxConnsPerHost:     500,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				Renegotiation:      tls.RenegotiateOnceAsClient,
			},
		}}

	simplerunner.client.CheckRedirect = nil

	return &simplerunner
}

func (r *SimpleRunner) Prepare(url string) *http.Request {
	user_agent := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/82.0.4083.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.3538.77 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:75.0) Gecko/20100101 Firefox/75.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/82.0.4080.0 Safari/537.36 Edg/82.0.453.0",
		"Mozilla/5.0 (Windows NT 10.0; rv:76.0) Gecko/20100101 Firefox/76.0"}
	rand.Seed(time.Now().UnixNano())
	node := rand.Intn(6)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("url：" + url + " 错误！!!！")
		return nil
	}
	req.Header.Add("User-Agent", user_agent[node])
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	req.Header.Add("Accept-Encoding", "text/html")
	req.Header.Add("Referer", "https://www.google.com/")
	req.Header.Add("Cache-Control", "Cache-Control")
	req = req.WithContext(r.config.Context)
	req.Proto = "HTTP/1.1"
	return req
}

func determineEncoding(r io.Reader) (encoding.Encoding, []byte) {
	content, _ := ioutil.ReadAll(r)
	if len(content) == 0 {
		return nil, nil
	}
	start := bytes.Index(content,[]byte("title"))
	if start == -1{
		start = 0
	}
	e, _, _ := charset.DetermineEncoding(content[start:], "")
	return e, content
}

func (r *SimpleRunner) Execute(req *http.Request, ahost map[string]bool) (*weblive.CollyData, error) {
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	weblive.Lock.Lock()
	if _, ok := ahost[resp.Request.URL.Host+resp.Request.URL.Scheme]; ok {
		weblive.Lock.Unlock()
		return nil, nil
	} else {
		ahost[resp.Request.URL.Host+resp.Request.URL.Scheme] = true
	}
	weblive.Lock.Unlock()
	cData := &weblive.CollyData{}
	cData.Url = req.URL
	cData.Redirect = resp.Request.URL
	cData.StatusCode = resp.StatusCode

	e, byteht := determineEncoding(resp.Body)
	if e == nil && byteht == nil {
		cData.Html = ""
	} else {
		ht, _, _ := transform.Bytes(e.NewDecoder(), byteht)
		defer resp.Body.Close()
		cData.Html = string(ht)
	}

	cData.Cookies = make(map[string]string)
	for _, cookie := range resp.Header["Set-Cookie"] {
		keyValues := strings.Split(cookie, ";")
		for _, keyValueString := range keyValues {
			keyValueSlice := strings.Split(keyValueString, "=")
			if len(keyValueSlice) > 1 {
				if weblive.IsContainStr(strings.ToLower(strings.Trim(keyValueSlice[0], " ")), []string{"expires", "domain", "path", "samesite", "max-age", "version"}) {
					continue
				}
				key, value := strings.ToLower(strings.Trim(keyValueSlice[0], " ")), keyValueSlice[1]
				cData.Cookies[key] = value
			}

		}

	}
	cData.Headers = make(map[string][]string)
	resp.Header.Del("Set-Cookie")
	for k, v := range resp.Header {
		lowerCaseKey := strings.ToLower(k)
		cData.Headers[lowerCaseKey] = v
	}
	if cData.Html != "" {
		doc, err := htmlquery.Parse(strings.NewReader(cData.Html))
		if err == nil {
			cData.Title = strings.Replace(weblive.GetTitle(doc), "\n", "", -1)
			scriptNode := htmlquery.Find(doc, "//script")
			for _, value := range scriptNode {
				src := htmlquery.SelectAttr(value, "src")
				if src != "" {
					cData.Scripts = append(cData.Scripts, src)
				}
			}
			metaNode := htmlquery.Find(doc, "//meta")
			cData.Meta = make(map[string]string)
			for _, value := range metaNode {
				name := htmlquery.SelectAttr(value, "name")
				content := htmlquery.SelectAttr(value, "content")
				if name != "" && content != "" {
					cData.Meta[strings.ToLower(name)] = content
				}
			}
		}
	}

	return cData, err
}
