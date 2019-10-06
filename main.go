package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

var (
	version = "0.0.0"
	build   = "0"
)

func main() {
	app := cli.NewApp()
	app.Name = "nomad plugin"
	app.Usage = "nomad plugin"
	app.Action = run
	app.Version = fmt.Sprintf("%s+%s", version, build)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "addr",
			Usage:  "nomad addr",
			EnvVar: "PLUGIN_ADDR",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "nomad token",
			EnvVar: "PLUGIN_TOKEN",
		},
		cli.StringFlag{
			Name:   "region",
			Usage:  "nomad region",
			EnvVar: "PLUGIN_REGION",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{}

	return plugin.Exec()
}
