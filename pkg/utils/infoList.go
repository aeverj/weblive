package utils

import (
	"fmt"
	"strconv"
	"strings"
	"weblive/pkg/weblive"
)

func GetInfoList(info *weblive.Website) []string {
	var infoList []string
	infoList = append(infoList, info.Url.String())
	infoList = append(infoList, info.Redirect.String())
	infoList = append(infoList, info.Title)
	infoList = append(infoList, strconv.Itoa(info.StatusCode))
	var ips []string
	for _, ip := range info.Ip {
		ips = append(ips, fmt.Sprint(ip))
	}
	infoList = append(infoList, strings.Join(ips, ","))
	infoList = append(infoList, fmt.Sprint(info.Cdn))
	infoList = append(infoList, strings.Join(info.Finger, ","))
	return infoList
}
