package main

import (
	"os"

	"github.com/mendersoftware/go-lib-micro/log"
	"github.com/urfave/cli"

	"github.com/canyanio/rating-agent-hep/config"
	"github.com/canyanio/rating-agent-hep/server"
)

func main() {
	doMain(os.Args)
}

func doMain(args []string) {
	var configPath string
	var configDebug bool

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Usage:       "Configuration `FILE`. Supports JSON, TOML, YAML and HCL formatted configs.",
				Value:       "config.yaml",
				Destination: &configPath,
			},
			&cli.BoolFlag{
				Name:        "debug",
				Usage:       "Enable debug mode and verbose logging",
				Destination: &configDebug,
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
		err := config.Init(configPath)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		log.Setup(configDebug)

		return nil
	}

	err := app.Run(args)
	if err != nil {
		cli.NewExitError(err.Error(), 1)
	}
}

func cmdAgent(args *cli.Context) error {
	srv := server.NewUDPServer()
	err := srv.Start()
	return err
}
