package weblive

import (
	"context"
	"net"
	"net/url"
	"sync"
)
type Config struct {
	OutputFile string          `json:"outputfile"`
	Timeout    int             `json:"timeout"`
	Threads    int             `json:"threads"`
	Context    context.Context `json:"-"`
	ProxyURL   string          `json:"proxyurl"`
	Path       string          `json:"path"`
}
var OutputPath = ""
var Lock sync.Mutex
type Website struct {
	Url        *url.URL
	Redirect   *url.URL
	Title      string
	StatusCode int
	Ip         []net.IP
	Cdn        string
	Finger     []string
}
type CollyData struct {
	Url        *url.URL
	Redirect   *url.URL
	Html       string
	Headers    map[string][]string
	Scripts    []string
	Cookies    map[string]string
	Meta       map[string]string
	StatusCode int
	Title string
}
