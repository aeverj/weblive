package screenshot

import (
	"testing"
	"time"
)

func TestDoScreen(t *testing.T) {
	go func() {
		go func() {
			go DoScreen("http://www.baidu.com", "test")
		}()
	}()
	time.Sleep(10 * time.Second)
}
