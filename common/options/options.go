package options

import "context"

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

type RequestOptions struct {
	Methods             string
	RequestPATH         string
	RequestBody         string
	CustomHeaders       CustomHeaders
	FollowHostRedirects bool
	Timeout             int
	Ports               string
}

type Options struct {
	BrowserThread int
	Threads       int
	InputArg      string
	OutputType    string
	Verbose       bool
	BrowserPath   string
	Target        chan string
	ReqOption     *RequestOptions
	Ctx           *context.Context
}

var CurrentOption *Options

func init() {
	CurrentOption = &Options{
		BrowserThread: 6,
		Threads:       30,
		Target:        make(chan string),
		OutputType:    "html",
		ReqOption: &RequestOptions{
			Methods: "GET",
			Timeout: 10,
			Ports:   "80,443",
		},
	}
}
