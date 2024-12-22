package report

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"weblive/common"
	"weblive/common/mlogger"
	"weblive/utils"
)

//go:embed template.html
var template string

func WriteHtml(filename string, result chan *common.DataRespResult) {
	content := make([]interface{}, 0)

	// 创建并启动进度显示
	spinner := utils.NewSpinner()
	spinner.Start()

	// 处理任务
	for v := range result {
		content = append(content, v.Result)
		spinner.Increment()
	}

	// 停止进度显示
	spinner.Stop()

	jsonData, err := json.Marshal(content)
	if err != nil {
		mlogger.Error("JSON marshal failed: " + err.Error())
		return
	}
	jsAssignment := []byte("window.assetData = " + string(jsonData))

	if len(jsonData) > 1 {
		err := os.WriteFile(fmt.Sprintf("%s/index.html", filename), []byte(template), 0755)
		if err != nil {
			mlogger.Error(err.Error())
		}
		err = os.WriteFile(fmt.Sprintf("%s/data.js", filename), jsAssignment, 0755)
		if err != nil {
			mlogger.Error(err.Error())
		}
	} else {
		mlogger.Info("No data to export")
		os.Remove(filename)
	}
}
