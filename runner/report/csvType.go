package report

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"weblive/common"
	"weblive/common/mlogger"
)

func WriteCsv(filename string, result chan *common.DataRespResult) {
	resultLength := 0
	currFile := filename + ".csv"
	nfs, err := os.OpenFile(currFile, os.O_RDWR|os.O_CREATE, 0755)
	// BOM header, solve excel garbled problem
	nfs.WriteString("\xEF\xBB\xBF")
	if err != nil {
		log.Fatalf("can not create file, err is %+v", err)
	}
	defer func() {
		nfs.Close()
		if resultLength == 0 {
			os.Remove(currFile)
			mlogger.Info("No data to export")
		} else {
			mlogger.Info(fmt.Sprintf("Write csv: %s", currFile))
		}
	}()
	w := csv.NewWriter(nfs)
	defer w.Flush()
	w.Write([]string{"URL", "Redirect", "Title", "Status_Code", "IP", "CDN", "Finger"})
	for v := range result {
		resultLength++
		mlogger.Info(fmt.Sprintf("%s %s %s %s", v.Result.Target, v.Result.StatusCode, v.Result.Title, v.Result.Finger))
		w.Write([]string{v.RespContent.Target, v.RespContent.Redirect, v.Result.Title, strconv.Itoa(v.RespContent.StatusCode), v.Result.IP, v.Result.CDN, v.Result.Finger})
	}
}
