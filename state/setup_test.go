package state

import (
	"flag"
	"os"
	"testing"

	"github.com/canyanio/rating-agent-hep/config"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if !testing.Short() {
		config.Init("")
	}
	result := m.Run()
	os.Exit(result)
}
