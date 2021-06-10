package weblive

import (
	"net"
	"net/url"
)
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
