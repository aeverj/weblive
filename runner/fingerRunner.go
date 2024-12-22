package runner

import (
	"strings"
	"weblive/common"
	"weblive/core/fingerprint"
)

func (r *Runner) runAsyncFinger(dataRespResults chan *common.DataRespResult) chan *common.DataRespResult {
	ch := make(chan *common.DataRespResult)
	go func() {
		for dataRespResult := range dataRespResults {
			fingerResult := fingerprint.DoHttpFingerPrint(fingerprint.HttpFingerInfo{
				Headers:    dataRespResult.RespContent.Headers,
				StatusCode: dataRespResult.RespContent.StatusCode,
				Html:       dataRespResult.RespContent.Html,
			})
			dataRespResult.Result.Finger = strings.Join(fingerResult, ", ")
			ch <- dataRespResult
		}
		close(ch)
	}()
	return ch
}
