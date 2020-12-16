package xlsxwriter

import (
	"bufio"
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/x/x/pkg/weblive"
	"io"
	"log"
	"os"
	"strings"
)

func OutputXlse(infoChan [][]string) {
	createResultPath()
	file := xlsx.NewFile()
	ipChan := make(chan string)
	sheet,err := file.AddSheet("sheet1")
	row := sheet.AddRow()
	row.AddCell().Value = "URL"
	row.AddCell().Value = "Redirect"
	row.AddCell().Value = "Title"
	row.AddCell().Value = "Status_Code"
	row.AddCell().Value = "IP"
	row.AddCell().Value = "CDN"
	row.AddCell().Value = "Finger"
	sheet.Cols[0].Width = 32
	sheet.Cols[1].Width = 32
	sheet.Cols[2].Width = 20
	sheet.Cols[3].Width = 5
	sheet.Cols[4].Width = 20
	sheet.Cols[5].Width = 5
	sheet.Cols[6].Width = 32
	if err != nil  {
		panic(err)
	}
	fmt.Printf("共发现%d个网站，保存于目录 %s\n",len(infoChan),weblive.OutputPath)
	go writeToFile(ipChan,weblive.OutputPath + "/ip.txt")
	var iplist []string
	for _,info := range infoChan{
		row := sheet.AddRow()
		for _,v := range info{
			row.AddCell().Value = v
		}
		if info[5] == "false"{
			for _,ip := range strings.Split(info[4],","){
				if !weblive.IsContainStr(ip,iplist){
					iplist = append(iplist, ip)
					ipChan <- ip
				}
			}
		}
	}
	close(ipChan)
	//row.SetHeightCM(1) //设置每行的高度
	err = file.Save(weblive.OutputPath + "/result.xlsx")
	if err != nil {
		panic(err)
	}
}

func createResultPath() {
	err:=os.MkdirAll(weblive.OutputPath,os.ModePerm)
	if err != nil {
		log.Panic("创建文件失败！")
	}
}

func writerSink(writer io.Writer, in <-chan string) {
	for v := range in {
		writer.Write([]byte(fmt.Sprintln(v)))
	}
}

func writeToFile(c <-chan string, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	writerSink(writer, c)

}