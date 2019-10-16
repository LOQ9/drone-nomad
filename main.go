package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

var (
	version = "0.1.0"
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
			Name:   "repo",
			Usage:  "docker repository",
			EnvVar: "PLUGIN_REPO",
		},
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
		cli.StringFlag{
			Name:   "template",
			Usage:  "nomad template",
			EnvVar: "PLUGIN_TEMPLATE",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Repo: Repo{
			Owner: c.String("repo.owner"),
			Name:  c.String("repo.name"),
		},
		Build: Build{
			Tag:     c.String("build.tag"),
			Number:  c.Int("build.number"),
			Event:   c.String("build.event"),
			Status:  c.String("build.status"),
			Commit:  c.String("commit.sha"),
			Ref:     c.String("commit.ref"),
			Branch:  c.String("commit.branch"),
			Author:  c.String("commit.author"),
			Link:    c.String("build.link"),
			Started: c.Int64("build.started"),
			Created: c.Int64("build.created"),
		},
		Job: Job{
			Started: c.Int64("job.started"),
		},
		Config: Config{
			Address:  c.String("addr"),
			Token:    c.String("token"),
			Region:   c.String("region"),
			Template: c.String("template"),
		},
	}

	return plugin.Exec()
}
