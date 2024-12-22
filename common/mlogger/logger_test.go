package mlogger

import "testing"

func TestError(t *testing.T) {
	Info("测试正常消息")
	Warn("测试警告消息")
	Error("测试错误消息")
}
