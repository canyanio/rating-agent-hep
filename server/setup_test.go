package server

import (
	"flag"
	"os"
	"testing"

	"github.com/canyanio/rating-agent-hep/config"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if !testing.Short() {
		config.Init("config.test.yml")
	}
	result := m.Run()
	os.Exit(result)
}
