package cdn

import (
	"fmt"
	"testing"
)

func TestResolve(t *testing.T) {
	cdn, ip, err := Resolve("epbf.bitautoimg.com")
	if err != nil {
		return
	}
	fmt.Println(cdn)
	fmt.Println(ip)
}
