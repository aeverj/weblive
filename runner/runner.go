package runner

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"github.com/remeh/sizedwaitgroup"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"weblive/common"
	"weblive/common/mlogger"
	"weblive/common/options"
	"weblive/runner/report"
	"weblive/utils"
)

type Runner struct {
	httpThread     int
	headlessThread int
	timeout        int
	resultPath     string
	urlChan        chan []string
	reqBody        io.Reader
}

func New() *Runner {
	runner := &Runner{
		httpThread:     options.CurrentOption.Threads,
		headlessThread: 6,
		timeout:        options.CurrentOption.ReqOption.Timeout,
		urlChan:        make(chan []string),
		resultPath:     time.Now().Format("2006-01-02_150405"),
	}
	runner.readBody()
	runner.prepareTarget()
	return runner
}

func readTarget() {
	opt := options.CurrentOption
	if utils.FileExists(opt.InputArg) {
		go func() {
			f, err := os.Open(opt.InputArg)
			defer f.Close()
			if err != nil {
				mlogger.Error(fmt.Sprintf("Open file %s error", opt.InputArg))
			}
			inputScanner := bufio.NewScanner(f)
			for inputScanner.Scan() {
				opt.Target <- inputScanner.Text()
			}
			close(opt.Target)
		}()
	} else {
		go func() {
			opt.Target <- opt.InputArg
			close(opt.Target)
		}()
	}
}

func (r *Runner) readBody() {
	opt := options.CurrentOption
	if opt.ReqOption.RequestBody == "" {
		return
	}
	if utils.FileExists(opt.ReqOption.RequestBody) {
		f, err := os.Open(opt.ReqOption.RequestBody)
		defer f.Close()
		if err != nil {
			mlogger.Error(fmt.Sprintf("Open file %s error", opt.ReqOption.RequestBody))
		}
		r.reqBody = f
	}
}

func (r *Runner) prepareTarget() {
	readTarget()
	go func(urlChan chan []string) {
		defer close(urlChan)
		for target := range options.CurrentOption.Target {
			switch {
			case strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://"):
				urlChan <- []string{target}
			case strings.Contains(target, ":"):
				urlChan <- []string{fmt.Sprintf("https://%s", target), fmt.Sprintf("http://%s", target)}
			case !strings.Contains(target, ":") && !strings.Contains(target, "/"):
				for _, v1 := range strings.Split(options.CurrentOption.ReqOption.Ports, ",") {
					if strings.Contains(v1, "-") {
						ipArr := strings.Split(v1, "-")
						start, _ := strconv.Atoi(ipArr[0])
						end, _ := strconv.Atoi(ipArr[1])
						for i := start; i < end+1; i++ {
							urlChan <- []string{fmt.Sprintf("https://%s:%d", target, i), fmt.Sprintf("http://%s:%d", target, i)}
						}
					} else {
						port, _ := strconv.Atoi(v1)
						urlChan <- []string{fmt.Sprintf("https://%s:%d", target, port), fmt.Sprintf("http://%s:%d", target, port)}
					}
				}
			default:
				mlogger.Warn(fmt.Sprintf("Target %s is invalid", target))
			}
		}
	}(r.urlChan)
}

func (r *Runner) Requests(url string) (respContent *common.RespContent, err error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: time.Duration(options.CurrentOption.ReqOption.Timeout) * time.Second,
	}
	request, err := http.NewRequest(options.CurrentOption.ReqOption.Methods, url, r.reqBody)
	if err != nil {
		mlogger.Warn(fmt.Sprintf("NewRequest error url is: %s", url))
		return
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36")
	request.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	for _, header := range options.CurrentOption.ReqOption.CustomHeaders {
		request.Header.Set(strings.Trim(strings.Split(header, ":")[0], " "), strings.Trim(strings.Split(header, ":")[1], " "))
	}
	resp, err := client.Do(request)
	if err != nil {
		mlogger.Warn(fmt.Sprintf("An exception occurred while accessing the website %s", url))
		return
	}
	Redirect := ""
	if resp.Request.URL.String() != url {
		Redirect = resp.Request.URL.String()
	}
	respContent = &common.RespContent{
		Target:     url,
		Redirect:   Redirect,
		Headers:    resp.Header,
		StatusCode: resp.StatusCode,
		Html:       utils.ReadBody(resp.Body),
	}
	return
}

func (r *Runner) asyncRequest() (respContentChan chan *common.RespContent) {
	respContentChan = make(chan *common.RespContent)
	go func() {
		swg := sizedwaitgroup.New(r.httpThread)
		for urlList := range r.urlChan {
			swg.Add()
			go func(urlList []string, swg *sizedwaitgroup.SizedWaitGroup) {
				defer swg.Done()
				for _, u := range urlList {
					respContent, err := r.Requests(u)
					if err != nil {
						mlogger.Warn(err.Error())
						continue
					}
					respContentChan <- respContent
					break
				}
			}(urlList, &swg)
		}
		swg.Wait()
		close(respContentChan)
	}()
	return
}

func (r *Runner) Run() {
	dataChan := make(chan *common.DataRespResult)
	go func() {
		for respContent := range r.asyncRequest() {
			headers := ""
			for name, value := range respContent.Headers {
				headers += fmt.Sprintf("%s: %s\n", name, value[0])
			}
			data := &common.DataRespResult{
				RespContent: respContent,
				Result: &common.Result{
					Target:     respContent.Target,
					Redirect:   respContent.Redirect,
					Title:      utils.GetTitle(respContent.Html),
					StatusCode: respContent.StatusCode,
					Header:     headers,
				},
			}
			dataChan <- data
		}
		close(dataChan)
	}()

	resultChan := make(chan *common.DataRespResult)
	switch options.CurrentOption.OutputType {
	case "html":
		fingerChan := r.runAsyncFinger(dataChan)
		dnsChan := r.runAsyncDNS(fingerChan)
		resultChan = r.runAsyncScreen(dnsChan)
		report.WriteHtml(r.resultPath, resultChan)
	case "csv":
		fingerChan := r.runAsyncFinger(dataChan)
		resultChan = r.runAsyncDNS(fingerChan)
		report.WriteCsv(r.resultPath, resultChan)
	}
}
