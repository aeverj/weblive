package weblive

import (
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"net"
	"strings"
)

func Set(items []net.IP) []net.IP {
	var ips []net.IP
	tmp := make(map[string]interface{})
	for _, eachitem := range items {
		l := len(tmp)
		tmp[eachitem.String()] = 1
		if len(tmp) > l {
			ips = append(ips, eachitem)
		}
	}
	return ips
}

func IsContainIP(item net.IP, items []string) bool {
	for _, eachitem := range items {
		_, cdir, _ := net.ParseCIDR(eachitem)
		if cdir.Contains(item) {
			return true
		}
	}
	return false
}

func IsContainInt(item uint, items []uint) bool {
	for _, eachitem := range items {
		if item == eachitem {
			return true
		}
	}
	return false
}

func IsContainStr(item string, items []string) bool {
	for _, eachItem := range items {
		if eachItem == strings.Trim(item, " ") {
			return true
		}
	}
	return false
}
func GetTitle(node *html.Node) string {
	title := ""

	if titlenode := htmlquery.FindOne(node, "//title"); titlenode != nil {
		title = htmlquery.InnerText(titlenode)
		if title != "" {
			return title
		}
	}
	if titlenode := htmlquery.FindOne(node, "//h1"); titlenode != nil {
		title = htmlquery.InnerText(titlenode)
		if title != "" {
			return title
		}
	}
	if titlenode := htmlquery.FindOne(node, "//h2"); titlenode != nil {
		title = htmlquery.InnerText(titlenode)
		if title != "" {
			return title
		}
	}
	if titlenode := htmlquery.FindOne(node, "//h3"); titlenode != nil {
		title = htmlquery.InnerText(titlenode)
		if title != "" {
			return title
		}
	}
	return title

}
