package runner

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"weblive/pkg/weblive"
)

type CustomHeaders []string

// String returns just a label
func (c *CustomHeaders) String() string {
	return "Custom Global Headers"
}

// Set a new global header
func (c *CustomHeaders) Set(value string) error {
	*c = append(*c, value)
	return nil
}

type scanOptions struct {
	Methods             string
	RequestPATH         string
	RequestBody         string
	CustomHeaders       CustomHeaders
	FollowHostRedirects bool
	Verbose             bool
	Timeout             int
	Ports               string
}

type options struct {
	OutputFile  string
	Timeout     int
	Threads     int
	InputFile   string
	ScanOptions *scanOptions
	Ctx         *context.Context
}

func ParseOptions() *options {
	scoption := &scanOptions{}
	options := &options{ScanOptions: scoption}
	flag.StringVar(&options.InputFile, "iF", "", "Load urls from file")
	flag.IntVar(&options.Timeout, "timeout", 3, "Timeout in seconds")
	flag.IntVar(&options.Threads, "threads", 50, "Number of threads")
	flag.Var(&scoption.CustomHeaders, "H", "Custom Header")
	flag.StringVar(&scoption.Methods, "M", "GET", "Request Method")
	flag.StringVar(&scoption.RequestPATH, "path", "/", "Request Path")
	flag.StringVar(&scoption.RequestBody, "dataFile", "", "The Post data file path")
	flag.BoolVar(&scoption.FollowHostRedirects, "follow_redirects", false, "Follow Redirects")
	//flag.BoolVar(&scoption.Verbose, "v", false, "Verbose")
	flag.StringVar(&scoption.Ports, "ports", "", "Custom ports")
	flag.StringVar(&options.OutputFile, "output", "", "Output file")
	flag.Parse()
	validateOptions(options)
	return options
}

func validateOptions(opt *options) {
	if opt.InputFile == "" || !weblive.FileExists(opt.InputFile) {
		filename := filepath.Base(os.Args[0])
		log.Fatalf("The -iF parameter is missing or file does not exist.\nExample:\n  %s -iF input.txt", filename)
	}
	if opt.ScanOptions.RequestBody != "" {
		if !weblive.FileExists(opt.ScanOptions.RequestBody) {
			log.Fatalf("File %s does not exist!\n", opt.ScanOptions.RequestBody)
		}
		f, err := os.Open(opt.ScanOptions.RequestBody)
		if err != nil {
			log.Fatalf("File %s does not open!\n", opt.ScanOptions.RequestBody)
		} else {
			f.Close()
		}
	}
}
