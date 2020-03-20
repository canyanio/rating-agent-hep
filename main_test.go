package main

import (
	"flag"
	"os"
	"strings"
	"testing"
)

var runAcceptanceTests bool

// used for parsing '-cli-args' for urfave/cli when running acceptance tests
// this is because of a conflict between urfave/cli and regular go flags required for testing (can't mix the two)
var cliArgsRaw string

var _ = func() bool {
	testing.Init()
	return true
}()

func init() {
	flag.BoolVar(&runAcceptanceTests, "acceptance-tests", false, "set flag when running acceptance tests")
	flag.StringVar(&cliArgsRaw, "cli-args", "", "for passing urfave/cli args (single string) when golang flags are specified (avoids conflict)")

	flag.Parse()
}

func TestRunMain(t *testing.T) {
	if !runAcceptanceTests {
		t.Skip()
	}

	// parse '-cli-args', remember about binary name at idx 0
	var cliArgs []string

	if cliArgsRaw != "" {
		cliArgs = []string{os.Args[0]}
		splitArgs := strings.Split(cliArgsRaw, " ")
		cliArgs = append(cliArgs, splitArgs...)
	}

	doMain(cliArgs)
}
