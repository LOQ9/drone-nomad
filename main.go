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
			Name:   "consul_token",
			Usage:  "consul token",
			EnvVar: "PLUGIN_CONSUL_TOKEN",
		},
		cli.StringFlag{
			Name:   "vault_token",
			Usage:  "vault token",
			EnvVar: "PLUGIN_VAULT_TOKEN",
		},
		cli.StringFlag{
			Name:   "region",
			Usage:  "nomad region",
			EnvVar: "PLUGIN_REGION",
		},
		cli.StringFlag{
			Name:   "namespace",
			Usage:  "nomad namespace",
			EnvVar: "PLUGIN_NAMESPACE",
		},
		cli.StringFlag{
			Name:   "template",
			Usage:  "nomad template",
			EnvVar: "PLUGIN_TEMPLATE",
		},
		cli.BoolFlag{
			Name:   "preserve_counts",
			Usage:  "nomad preserve_counts",
			EnvVar: "PLUGIN_PRESERVE_COUNTS",
		},
		cli.BoolFlag{
			Name:   "watch_deployment",
			Usage:  "nomad watch_deployment",
			EnvVar: "PLUGIN_WATCH_DEPLOYMENT",
		},
		cli.DurationFlag{
			Name:   "watch_deployment_timeout",
			Usage:  "nomad watch_deployment_timeout",
			EnvVar: "PLUGIN_WATCH_DEPLOYMENT_TIMEOUT",
		},
		cli.StringFlag{
			Name:   "tls_ca_cert",
			Usage:  "nomad tls ca certificate file",
			EnvVar: "PLUGIN_TLS_CA_CERT",
		},
		cli.StringFlag{
			Name:   "tls_ca_path",
			Usage:  "nomad tls ca certificate file path",
			EnvVar: "PLUGIN_TLS_CA_PATH",
		},
		cli.StringFlag{
			Name:   "tls_ca_cert_pem",
			Usage:  "nomad tls ca certificate pem",
			EnvVar: "PLUGIN_TLS_CA_CERT_PEM",
		},
		cli.StringFlag{
			Name:   "tls_client_cert",
			Usage:  "nomad tls client certificate",
			EnvVar: "PLUGIN_TLS_CLIENT_CERT",
		},
		cli.StringFlag{
			Name:   "tls_client_cert_pem",
			Usage:  "nomad tls client certificate pem",
			EnvVar: "PLUGIN_TLS_CLIENT_CERT_PEM",
		},
		cli.StringFlag{
			Name:   "tls_client_key",
			Usage:  "nomad tls client private key",
			EnvVar: "PLUGIN_TLS_CLIENT_KEY",
		},
		cli.StringFlag{
			Name:   "tls_client_key_pem",
			Usage:  "nomad tls client private key pem",
			EnvVar: "PLUGIN_TLS_CLIENT_KEY_PEM",
		},
		cli.StringFlag{
			Name:   "tls_server_name",
			Usage:  "nomad tls server name",
			EnvVar: "PLUGIN_TLS_SERVERNAME",
		},
		cli.BoolFlag{
			Name:   "tls_insecure",
			Usage:  "nomad tls insecure",
			EnvVar: "PLUGIN_TLS_INSECURE",
		},
		cli.StringFlag{
			Name:   "debug",
			Usage:  "nomad debug",
			EnvVar: "PLUGIN_DEBUG",
		},
		cli.StringFlag{
			Name:   "dry_run",
			Usage:  "nomad dry run",
			EnvVar: "PLUGIN_DRY_RUN",
		},
		cli.StringFlag{
			Name:   "repo.owner",
			Usage:  "repository owner",
			EnvVar: "DRONE_REPO_OWNER",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA",
		},
		cli.StringFlag{
			Name:   "commit.ref",
			Value:  "refs/heads/master",
			Usage:  "git commit ref",
			EnvVar: "DRONE_COMMIT_REF",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "git commit branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},
		cli.StringFlag{
			Name:   "commit.author",
			Usage:  "git author name",
			EnvVar: "DRONE_COMMIT_AUTHOR",
		},
		cli.StringFlag{
			Name:   "commit.message",
			Usage:  "commit message",
			EnvVar: "DRONE_COMMIT_MESSAGE",
		},
		cli.StringFlag{
			Name:   "build.event",
			Value:  "push",
			Usage:  "build event",
			EnvVar: "DRONE_BUILD_EVENT",
		},
		cli.IntFlag{
			Name:   "build.number",
			Usage:  "build number",
			EnvVar: "DRONE_BUILD_NUMBER",
		},
		cli.IntFlag{
			Name:   "build.parent",
			Usage:  "build parent",
			EnvVar: "DRONE_BUILD_PARENT",
		},
		cli.StringFlag{
			Name:   "build.status",
			Usage:  "build status",
			Value:  "success",
			EnvVar: "DRONE_BUILD_STATUS",
		},
		cli.StringFlag{
			Name:   "build.link",
			Usage:  "build link",
			EnvVar: "DRONE_BUILD_LINK",
		},
		cli.Int64Flag{
			Name:   "build.started",
			Usage:  "build started",
			EnvVar: "DRONE_BUILD_STARTED",
		},
		cli.Int64Flag{
			Name:   "build.created",
			Usage:  "build created",
			EnvVar: "DRONE_BUILD_CREATED",
		},
		cli.StringFlag{
			Name:   "build.tag",
			Usage:  "build tag",
			EnvVar: "DRONE_TAG",
		},
		cli.Int64Flag{
			Name:   "job.started",
			Usage:  "job started",
			EnvVar: "DRONE_JOB_STARTED",
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
			Parent:  c.Int("build.parent"),
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
			Address:                c.String("addr"),
			Token:                  c.String("token"),
			ConsulToken:            c.String("consul_token"),
			VaultToken:             c.String("vault_token"),
			Region:                 c.String("region"),
			Namespace:              c.String("namespace"),
			Template:               c.String("template"),
			PreserveCounts:         c.Bool("preserve_counts"),
			WatchDeployment:        c.Bool("watch_deployment"),
			WatchDeploymentTimeout: c.Duration("watch_deployment_timeout"),
			TLSCACert:              c.String("tls_ca_cert"),
			TLSCAPath:              c.String("tls_ca_path"),
			TLSCACertPem:           c.String("tls_ca_cert_pem"),
			TLSClientCert:          c.String("tls_client_cert"),
			TLSClientCertPem:       c.String("tls_client_cert_pem"),
			TLSClientKey:           c.String("tls_client_key"),
			TLSClientKeyPem:        c.String("tls_client_key_pem"),
			TLSServerName:          c.String("tls_server_name"),
			TLSInsecure:            c.Bool("tls_insecure"),
			Debug:                  c.Bool("debug"),
			DryRun:                 c.Bool("dry_run"),
		},
	}

	return plugin.Exec()
}
