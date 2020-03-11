package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mendersoftware/go-lib-micro/config"
	"github.com/urfave/cli"

	dconfig "github.com/canyanio/rating-agent-hep/config"
)

func main() {
	doMain(os.Args)
}

func doMain(args []string) {
	var configPath string

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Usage:       "Configuration `FILE`. Supports JSON, TOML, YAML and HCL formatted configs.",
				Value:       "config.yaml",
				Destination: &configPath,
			},
		},
		Commands: []cli.Command{
			{
				Name:   "agent",
				Usage:  "Run the HEP agent",
				Action: cmdAgent,
				Flags:  []cli.Flag{},
			},
		},
	}
	app.Usage = "rating-agent-hep"
	app.Version = "1.0.0"
	app.Action = cmdAgent

	app.Before = func(args *cli.Context) error {
		err := config.FromConfigFile(configPath, dconfig.Defaults)
		if err != nil {
			return cli.NewExitError(
				fmt.Sprintf("error loading configuration: %s", err),
				1)
		}

		// Enable setting config values by environment variables
		config.Config.SetEnvPrefix("RATING_HEP_AGENT")
		config.Config.AutomaticEnv()
		config.Config.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

		return nil
	}

	err := app.Run(args)
	if err != nil {
		log.Fatal(err)
	}
}

func cmdAgent(args *cli.Context) error {
	return nil
}
