package parameter

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"weblive/common/mlogger"
	"weblive/common/options"
	"weblive/utils"
)

func validateOptions() {
	opt := options.CurrentOption
	if opt.InputArg == "" {
		filename := filepath.Base(os.Args[0])
		mlogger.Error(fmt.Sprintf("The -i parameter is missing.\nExample:\n  %s -i input.txt \n\tor\n %s -i domain", filename, filename))
	}

	if opt.ReqOption.RequestBody != "" {
		if !utils.FileExists(opt.ReqOption.RequestBody) {
			mlogger.Error(fmt.Sprintf("File %s does not exist!\n", opt.ReqOption.RequestBody))
		}
		f, err := os.Open(opt.ReqOption.RequestBody)
		if err != nil {
			mlogger.Error(fmt.Sprintf("File %s does not open!\n", opt.ReqOption.RequestBody))
		} else {
			err := f.Close()
			if err != nil {
				mlogger.Error(fmt.Sprintf("File [%s] closed exception", opt.ReqOption.RequestBody))
			}
		}
	}
	if opt.ReqOption.Ports == "" {
		mlogger.Warn("You should set the port parameters, I've set 80,443 for you.")
		opt.ReqOption.Ports = "80,443"
	}
}

func ParseOptions() {
	flag.IntVar(&options.CurrentOption.Threads, "th", 30, "Number of threads")
	flag.StringVar(&options.CurrentOption.BrowserPath, "e", "", "Chrome executable file path")
	flag.IntVar(&options.CurrentOption.BrowserThread, "bth", 6, "Number of browser threads")
	flag.StringVar(&options.CurrentOption.InputArg, "i", "", "Input file path")
	flag.StringVar(&options.CurrentOption.OutputType, "o", "html", "Output type; html or csv")
	flag.IntVar(&options.CurrentOption.ReqOption.Timeout, "timeout", 10, "Number of timeout")
	flag.Var(&options.CurrentOption.ReqOption.CustomHeaders, "H", "Custom Header")
	flag.StringVar(&options.CurrentOption.ReqOption.Methods, "M", "GET", "Request Method")
	flag.StringVar(&options.CurrentOption.ReqOption.RequestPATH, "body_path", "/", "Url path")
	flag.StringVar(&options.CurrentOption.ReqOption.RequestBody, "data_file", "", "The Post data file path")
	flag.BoolVar(&options.CurrentOption.ReqOption.FollowHostRedirects, "follow_redirects", false, "Follow Redirects")
	flag.BoolVar(&options.CurrentOption.Verbose, "v", false, "Verbose")
	flag.StringVar(&options.CurrentOption.ReqOption.Ports, "ports", "80,443", "Custom ports")
	flag.Parse()
	validateOptions()
}
