package runner

import (
	"fmt"
	"github.com/remeh/sizedwaitgroup"
	"weblive/common"
	"weblive/common/mlogger"
	"weblive/core/screenshot"
)

func (r *Runner) runAsyncScreen(dataRespResults chan *common.DataRespResult) chan *common.DataRespResult {
	ch := make(chan *common.DataRespResult)
	go func() {
		swg := sizedwaitgroup.New(r.headlessThread)
		for dataRespResult := range dataRespResults {
			swg.Add()
			go func(dataRespResult *common.DataRespResult, swg *sizedwaitgroup.SizedWaitGroup) {
				defer swg.Done()
				screenPath, err := screenshot.DoScreen(dataRespResult.RespContent.Target, r.resultPath)
				if err != nil {
					mlogger.Warn(fmt.Sprintf("An exception occurred while accessing the website %s", dataRespResult.RespContent.Target))
				}
				dataRespResult.Result.ScreenPath = screenPath
				ch <- dataRespResult
			}(dataRespResult, &swg)
		}
		swg.Wait()
		close(ch)
	}()
	return ch
}
