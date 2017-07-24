package main

import (
	"fmt"
	"os"

	"log"

	"github.com/honestbee/drone-helm-repo/pkg/util"
	"gopkg.in/urfave/cli.v1"
)

var build = "0" // build number set at compile-time
var logger *util.Logger

func main() {
	app := initApp(run)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func initApp(runAction cli.ActionFunc) *cli.App {
	app := cli.NewApp()
	app.Name = "helm repo plugin"
	app.Usage = "Package and upload Helm charts to storage provider"
	app.Action = runAction
	app.Version = fmt.Sprintf("1.0.%s", build)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "source-dir",
			Usage:  "`PATH` to recursively search for Charts",
			EnvVar: "PLUGIN_SOURCE_DIR,SOURCE_DIR",
			Value:  ".",
		},
		cli.StringSliceFlag{
			Name:   "exclude",
			Usage:  "`LIST` of excluded directories in source-dir to exclude from Chart search",
			EnvVar: "PLUGIN_EXCLUDE,EXCLUDE",
		},
		cli.StringFlag{
			Name:   "repo-url",
			Usage:  "`BASE_URL` for the helm repository",
			EnvVar: "PLUGIN_REPO_URL,REPO_URL",
		},
		cli.StringFlag{
			Name:   "storage-url",
			Usage:  "`URL` of the container to store charts to (i.e s3://my-bucket/prefix)",
			EnvVar: "PLUGIN_STORAGE_URL,STORAGE_URL",
		},
		cli.StringFlag{
			Name:   "aws-access-key",
			Usage:  "AWS Access Key `AWS_ACCESS_KEY`",
			EnvVar: "AWS_ACCESS_KEY_ID,AWS_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "aws-secret-key",
			Usage:  "AWS Secret Key `AWS_SECRET_KEY`",
			EnvVar: "AWS_SECRET_ACCESS_KEY,AWS_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   "aws-region",
			Usage:  "AWS Region `AWS_REGION`",
			EnvVar: "PLUGIN_AWS_REGION, AWS_REGION",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "show debug logs",
			EnvVar: "PLUGIN_DEBUG,DEBUG",
		},
	}
	return app
}

func run(c *cli.Context) error {
	logger = util.NewLogger(c.Bool("debug"), false)
	plugin := Plugin{
		Config: Config{
			SourceDir:    c.String("source-dir"),
			Exclude:      c.StringSlice("exclude"),
			StorageURL:   c.String("storage-url"),
			RepoURL:      c.String("repo-url"),
			Debug:        c.Bool("debug"),
			AWSAccessKey: c.String("aws-acces-key"),
			AWSSecretKey: c.String("aws-secret-key"),
			AWSRegion:    c.String("aws-region"),
		},
	}
	return plugin.Exec()
}
