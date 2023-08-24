package runner

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
	"weblive/pkg/utils"
	"weblive/pkg/wappalyzer"
	"weblive/pkg/weblive"
)

func Run() {
	// Parse the command line flags and verify
	options := ParseOptions()
	output := make(chan []string)

	// Init GeoLite for analyze web page
	wapp := wappalyzer.Init()
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
	os.Mkdir("result", 0766)
	nfs, err := os.OpenFile(options.OutputFile, os.O_RDWR|os.O_CREATE, 0666)
	// BOM header, solve excel garbled problem
	nfs.WriteString("\xEF\xBB\xBF")
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
	simplerun := NewSimpleRunner(options)
	urlchan := make(chan string)
	var wg sync.WaitGroup

	// Gather web page info from chanel
	go utils.GetResult(output, w)

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
					info := utils.GetInfoList(wapp.Analyze(cdata))
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
