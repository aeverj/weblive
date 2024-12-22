package main

import (
	"weblive/runner"
	"weblive/runner/parameter"
)

func main() {
	parameter.ParseOptions()
	r := runner.New()
	r.Run()
}
