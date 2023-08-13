package utils

import (
	"encoding/csv"
	"fmt"
)

func GetResult(result chan []string, w *csv.Writer) {
	w.Write([]string{"URL", "Redirect", "Title", "Status_Code", "IP", "CDN", "Finger"})
	for v := range result {
		if len(v) == 7 {
			w.Write(v)
			fmt.Printf("%v %v %v\n", v[1], v[2], v[3])
		} else {
			fmt.Printf("%v %v\n", v[0], v[1])
		}
	}
}
