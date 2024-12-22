package mlogger

import (
	"fmt"
	"os"
	"weblive/common/options"
)

const (
	infoPrefix  = "[+]"
	warnPrefix  = "[*]"
	errorPrefix = "[x]"
	red         = "\033[31m"
	blue        = "\033[34m"
	yellow      = "\033[33m"
	reset       = "\033[0m"
)

func Error(msg string) {
	fmt.Printf("%s%s %s%s\n", red, errorPrefix, msg, reset)
	os.Exit(1)
}
func Info(msg string) {
	fmt.Printf("%s%s %s%s\n", blue, infoPrefix, msg, reset)
}
func Warn(msg string) {
	if options.CurrentOption.Verbose {
		fmt.Printf("%s%s %s%s\n", yellow, warnPrefix, msg, reset)
	}
}
