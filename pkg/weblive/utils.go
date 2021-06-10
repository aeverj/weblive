package weblive

import (
	"bufio"
	"github.com/antchfx/htmlquery"
	"github.com/projectdiscovery/mapcidr"
	"golang.org/x/net/html"
	"net"
	"os"
	"path/filepath"
	"strings"
)



func GetUrl(path string) (*bufio.Scanner, *os.File) {
	file, _ := os.Open(path)
	reader := bufio.NewScanner(file)
	reader.Split(bufio.ScanLines)
	return reader, file
}

// returns all the targets within a cidr range or the single target
func Targets(target string, ports string) chan string {
	results := make(chan string)
	go func() {
		defer close(results)

		// A valid target does not contain:
		// *
		// spaces
		if strings.ContainsAny(target, " *") {
			return
		}

		// test if the target is a cidr
		if IsCidr(target) {
			cidrIps, err := mapcidr.IPAddresses(target)
			if err != nil {
				return
			}
			for _, ip := range cidrIps {
				if len(ports) > 0 {
					for _, port := range strings.Split(ports, ",") {
						results <- ip + ":" + port
					}
				}else {
					results <- ip
				}
			}
		} else if strings.Index(target, "http") < 0 {
			if len(ports) > 0 {
				for _, port := range strings.Split(ports, ",") {
					results <- target + ":" + port
				}
			}else {
				results <- target
			}
		} else {
			results <- target
		}
	}()
	return results
}
// IsCidr determines if the given ip is a cidr range
func IsCidr(ip string) bool {
	_, _, err := net.ParseCIDR(ip)
	return err == nil
}

func FileExists(path string)bool {
	f,err := filepath.Glob(path)
	if err == nil {
		if len(f)>0{
			return true
		}
	}
	return false

}

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
