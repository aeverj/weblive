package cdn

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/miekg/dns"
	"net"
	"strings"
	"time"
	"weblive/common/mlogger"
)

//go:embed cdn.json
var apps []byte

type CDN struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

var cdnMap map[string]CDN

func init() {
	err := json.Unmarshal(apps, &cdnMap)
	if err != nil {
		mlogger.Error(fmt.Sprintf("cdn.json is error: %s", err))
	}
}

func ResolveIP(host string) []net.IP {
	ns, err := net.LookupIP(strings.Split(host, ":")[0])
	if err != nil {
		mlogger.Warn(fmt.Sprintf("ResolveIP err: %s", err))
		return nil
	}
	return ns
}

func Resolve(src string) (cdn string, dstIP string, err error) {
	var lastErr error
	var cnameList []string
	var tmpIP []string
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
					tmpIP = append(tmpIP, record.A.String())
				}
			}
		}
		lastErr = nil
		break
	}
	err = lastErr
	dstIP = strings.Join(tmpIP, ", ")
	for _, v := range cnameList {
		for cdnk, cdnv := range cdnMap {
			if strings.Contains(v, cdnk) {
				cdn = cdnv.Name
				return
			}
		}
	}
	if len(dstIP) > 1 {
		cdn = "unknown"
	} else {
		cdn = "no"
	}
	return
}
