package runner

import (
	"testing"
	"weblive/common/options"
)

func TestPrepareTarget(t *testing.T) {
	options.CurrentOption.InputArg = "E:\\gitrepo\\weblive\\input.txt"
	//options.CurrentOption.Verbose = true
	options.CurrentOption.OutputType = "csv"
	runner := New()
	runner.Run()
}
