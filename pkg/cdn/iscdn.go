package cdn

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"net"
	"os"
	"strings"
	"time"
	"weblive/pkg/weblive"
)

//go:embed cdn.json
var apps []byte

//type CDN struct {
//	Item map[string]*Info
//}

type CDN struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

var cdnMap map[string]CDN

func init() {
	err := json.Unmarshal(apps, &cdnMap)
	if err != nil {
		log.Fatalf("apps.json is error: %s", err)
	}
}

func ResolveIP(host string) []net.IP {
	ns, err := net.LookupIP(strings.Split(host, ":")[0])
	if err != nil {
		return nil
	}
	ns = weblive.Set(ns)
	return ns
}

func Resolve(src string) (cdn string, dstIP []net.IP, err error) {
	var lastErr error
	var cnameList []string
	for i := 0; i < 3; i++ {
		c := new(dns.Client)
		m := new(dns.Msg)
		m.SetQuestion(dns.Fqdn(src), dns.TypeA)
		m.RecursionDesired = true
		r, _, err := c.Exchange(m, "114.114.114.114:53")
		if err != nil {
			lastErr = err
			time.Sleep(1 * time.Second * time.Duration(i+1))
			continue
		}
		for _, ans := range r.Answer {
			switch ans.Header().Rrtype {
			case dns.TypeCNAME:
				record, isType := ans.(*dns.CNAME)
				if isType {
					cnameList = append(cnameList, record.Target)
				}
			case dns.TypeA:
				record, isType := ans.(*dns.A)
				if isType {
					dstIP = append(dstIP, record.A)
				}
			}
		}
		lastErr = nil
		break
	}
	err = lastErr
	for _, v := range cnameList {
		for cdnk, cdnv := range cdnMap {
			if strings.Contains(v, cdnk) {
				cdn = cdnv.Name
				return
			}
		}
	}
	if len(dstIP) > 1 {
		cdn = "无法判断"
	} else {
		cdn = "无"
	}
	return
}

func Mdns() {
	//config, _ := dns.ClientConfigFromFile("/etc/resolv.conf")
	c := new(dns.Client)

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn("crm.cscec.com"), dns.TypeA)
	m.RecursionDesired = true

	r, _, err := c.Exchange(m, "223.5.5.5:53")
	if r == nil {
		log.Fatalf("*** error: %s\n", err.Error())
	}

	if r.Rcode != dns.RcodeSuccess {
		log.Fatalf(" *** invalid answer name %s after MX query for %s\n", os.Args[1], os.Args[1])
	}
	// Stuff must be in the answer section
	for _, a := range r.Answer {
		fmt.Printf("%v\n", a)
	}
}
