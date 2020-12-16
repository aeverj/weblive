package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/x/x/pkg/runner"
	"github.com/x/x/pkg/wappalyzer"
	"github.com/x/x/pkg/weblive"
	"github.com/x/x/pkg/xlsxwriter"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func getUrl(path string) (*bufio.Scanner, *os.File) {
	file, err := os.Open(path)
	if err != nil {
		panic("获取url失败，无法找到文件url.txt")
	}
	reader := bufio.NewScanner(file)
	reader.Split(bufio.ScanLines)
	return reader, file
}

func getInfoList(info *weblive.Website) []string {
	var infoList []string
	infoList = append(infoList, info.Url.String())
	infoList = append(infoList, info.Redirect.String())
	infoList = append(infoList, info.Title)
	infoList = append(infoList, strconv.Itoa(info.StatusCode))
	var ips []string
	for _, ip := range info.Ip {
		ips = append(ips, fmt.Sprint(ip))
	}
	infoList = append(infoList, strings.Join(ips, ","))
	infoList = append(infoList, fmt.Sprint(info.Cdn))
	infoList = append(infoList, strings.Join(info.Finger, ","))
	return infoList
}

func main() {
	var wapp = wappalyzer.Init()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	t := time.Now().Format("2006-01-02_150405")
	path := fmt.Sprintf("./result/%s", t)
	weblive.OutputPath = path
	c := &weblive.Config{
		OutputFile: path,
		Timeout:    20,
		Threads:    40,
		Context:    ctx,
		ProxyURL:   "",
		Path:       "url.txt",
	}
	simplerun := runner.NewSimpleRunner(c)

	allhost := make(map[string]bool)
	urlchan := make(chan string)
	var webInfo [][]string
	var wg sync.WaitGroup
	reader, file := getUrl(c.Path)
	for i := 0; i < c.Threads; i++ {
		wg.Add(1)
	}
	go func() {
		defer file.Close()
		defer close(urlchan)
		for reader.Scan() {
			urlchan <- reader.Text()
		}

	}()

	for i := 0; i < c.Threads; i++ {
		go func() {
			defer wg.Done()
			for url := range urlchan {
				req := simplerun.Prepare(url)
				if req==nil{
					continue
				}
				cdata, _ := simplerun.Execute(req, allhost)
				if cdata != nil {
					info := getInfoList(wapp.Analyze(cdata))
					log.Println(strings.Join(info, " | "))
					webInfo = append(webInfo, info)
				}
			}
		}()
	}
	wg.Wait()
	xlsxwriter.OutputXlse(webInfo)
}
