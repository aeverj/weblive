package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/x/x/pkg/runner"
	"github.com/x/x/pkg/wappalyzer"
	"github.com/x/x/pkg/weblive"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

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

func getResult(result chan []string, w *csv.Writer) {
	w.Write([]string{"URL", "Redirect", "Title", "Status_Code", "IP", "CDN", "Finger"})
	for v := range result {
		if len(v) == 7 {
			w.Write(v)
			fmt.Printf("%v %v %v\n", v[1], v[2], v[3])
		} else {
			fmt.Printf("%v %v\n", v[0], v[1])
		}
	}
}

func main() {
	// Parse the command line flags and verify
	options := runner.ParseOptions()
	output := make(chan []string)

	// Init GeoLite for analyze web page
	wapp := wappalyzer.Init()
	defer wapp.Geodb.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	options.Ctx = &ctx

	// Create output file
	if options.OutputFile == "" {
		t := time.Now().Format("2006-01-02_150405")
		options.OutputFile = fmt.Sprintf("./result/%s.csv", t)
	}
	defer fmt.Println("result save for " + options.OutputFile)

	// Check if the folder named "result" exist
	os.Mkdir("result", 0666)
	nfs, err := os.OpenFile(options.OutputFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("can not create file, err is %+v", err)
	}
	defer nfs.Close()

	// Create a Writer for output file
	w := csv.NewWriter(nfs)
	defer func() {
		time.Sleep(500 * time.Millisecond)
		w.Flush()
	}()

	// Create http client
	simplerun := runner.NewSimpleRunner(options)
	urlchan := make(chan string)
	var wg sync.WaitGroup

	// Gather web page info from chanel
	go getResult(output, w)

	// Read the target from input file
	reader, file := weblive.GetUrl(options.InputFile)
	go func() {
		defer file.Close()
		defer close(urlchan)
		for reader.Scan() {
			urlchan <- reader.Text()
		}
	}()
	// Init thread count from option and start scan
	cht := make(chan interface{}, options.Threads)
	for target := range urlchan {
		for url := range weblive.Targets(target, options.ScanOptions.Ports) {
			cht <- struct{}{}
			wg.Add(1)
			go func(url string) {
				defer func() {
					wg.Done()
					<-cht
				}()
				cdata, err := simplerun.Execute(url)
				if err == nil {
					info := getInfoList(wapp.Analyze(cdata))
					output <- info
				} else {
					output <- []string{url, "Failed !!", err.Error()}
				}
			}(url)
		}
	}
	wg.Wait()
	close(output)
}
