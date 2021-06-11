package runner

import (
	"bytes"
	"crypto/tls"
	"github.com/antchfx/htmlquery"
	retryablehttp "github.com/projectdiscovery/retryablehttp-go"
	"github.com/x/x/pkg/weblive"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"
)

type SimpleRunner struct {
	client *retryablehttp.Client
	config *options
}

func NewSimpleRunner(options *options) *SimpleRunner {
	var simplerunner SimpleRunner
	simplerunner.config = options
	var retryablehttpOptions = retryablehttp.DefaultOptionsSpraying
	retryablehttpOptions.Timeout = time.Duration(time.Duration(options.Timeout) * time.Second)
	retryablehttpOptions.RetryMax = 0

	var redirectFunc = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse // Tell the http client to not follow redirect
	}
	if options.ScanOptions.FollowHostRedirects {
		// Only follow redirects on the same host
		redirectFunc = func(redirectedRequest *http.Request, previousRequest []*http.Request) error { // timo
			// Check if we get a redirect to a differen host
			var newHost = redirectedRequest.URL.Host
			var oldHost = previousRequest[0].URL.Host
			if newHost != oldHost {
				return http.ErrUseLastResponse // Tell the http client to not follow redirect
			}
			return nil
		}
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln("Failed to set cookies!")
	}
	simplerunner.client = retryablehttp.NewWithHTTPClient(&http.Client{
		CheckRedirect: redirectFunc,
		Jar:           jar,
		Timeout:       time.Duration(time.Duration(options.Timeout) * time.Second),
		Transport: &http.Transport{
			MaxIdleConnsPerHost: -1,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				Renegotiation:      tls.RenegotiateOnceAsClient,
			},
			DisableKeepAlives: true,
		}}, retryablehttpOptions)
	return &simplerunner
}

func (r *SimpleRunner) newRequest(targetURL string) (req *retryablehttp.Request, err error) {
	if r.config.ScanOptions.RequestBody == "" {
		req, err = retryablehttp.NewRequest(r.config.ScanOptions.Methods, targetURL, nil)
	}else {
		file, _ := os.Open(r.config.ScanOptions.RequestBody)
		req, err = retryablehttp.NewRequest(r.config.ScanOptions.Methods, targetURL,file)
	}
	if err != nil {
		return nil,err
	}
	// set default user agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/82.0.4080.0 Safari/537.36 Edg/82.0.453.0")
	// set default encoding to accept utf8
	req.Header.Add("Accept-Charset", "utf-8")

	if len(r.config.ScanOptions.CustomHeaders) > 0 {
		for _,c := range r.config.ScanOptions.CustomHeaders{
			headers := strings.Split(c,":")
			if len(headers) == 2 {
				req.Header.Set(strings.TrimSpace(headers[0]),strings.TrimSpace(headers[1]))
			}
		}
	}
	req.WithContext(*r.config.Ctx)
	return
}

func determineEncoding(r io.Reader) (encoding.Encoding, []byte) {
	content, _ := ioutil.ReadAll(r)
	if len(content) == 0 {
		return nil, nil
	}
	start := bytes.Index(content, []byte("title"))
	if start == -1 {
		start = 0
	}
	e, _, _ := charset.DetermineEncoding(content[start:], "")
	return e, content
}

func (r *SimpleRunner) Execute(targetURL string) (*weblive.CollyData, error) {
	protocol := ""
	if strings.Index(targetURL, "http") < 0 {
		protocol = "https://"
	}
retry:
	req,err := r.newRequest(protocol + targetURL)
	if err != nil {
		log.Fatalln(err)
		return nil,err
	}
	resp, err := r.client.Do(req)
	if err != nil {
		if protocol == "https://" {
			protocol = "http://"
			goto retry
		}
		return nil, err
	}
	defer resp.Body.Close()
	cData := &weblive.CollyData{}
	cData.Url = req.URL
	cData.Redirect = resp.Request.URL
	cData.StatusCode = resp.StatusCode

	e, byteht := determineEncoding(resp.Body)
	if e == nil && byteht == nil {
		cData.Html = ""
	} else {
		ht, _, _ := transform.Bytes(e.NewDecoder(), byteht)
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
