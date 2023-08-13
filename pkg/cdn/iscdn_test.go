package cdn

import (
	"fmt"
	"testing"
)

func TestIscdn(t *testing.T) {
	//cdnJson := Iscdn("crm.cscec.com")
	//Mdns()
	resolve, i, err := Resolve("crm.cscec.com")
	if err != nil {
		return
	}
	fmt.Printf("%#v\n", resolve)
	fmt.Printf("%#v", i)
}
