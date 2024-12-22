package runner

import (
	"fmt"
	"github.com/remeh/sizedwaitgroup"
	"strings"
	"weblive/common"
	"weblive/common/mlogger"
	"weblive/core/cdn"
)

func (r *Runner) runAsyncDNS(dataRespResults chan *common.DataRespResult) chan *common.DataRespResult {
	ch := make(chan *common.DataRespResult)
	go func() {
		swg := sizedwaitgroup.New(r.httpThread)
		for dataRespResult := range dataRespResults {
			swg.Add()
			go func(dataRespResult *common.DataRespResult, swg *sizedwaitgroup.SizedWaitGroup) {
				defer swg.Done()
				url := dataRespResult.RespContent.Redirect
				host := ""
				switch url {
				case "":
					host = strings.Split(strings.Split(dataRespResult.RespContent.Target, "://")[1], ":")[0]
				default:
					host = strings.Split(strings.Split(dataRespResult.RespContent.Redirect, "://")[1], ":")[0]
				}
				cdn, ip, err := cdn.Resolve(host)
				if err != nil {
					mlogger.Warn(fmt.Sprintf("Resolve err: %s", err))
				}
				dataRespResult.Result.CDN = cdn
				dataRespResult.Result.IP = ip
				ch <- dataRespResult
			}(dataRespResult, &swg)
		}
		swg.Wait()
		close(ch)
	}()
	return ch
}
