package main

import (
	"github.com/Pocket/ops-cli/internal/commands"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"time"
)

func main() {
	var app = cli.NewApp()
	addInfo(app)
	addCommands(app)
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func addInfo(app *cli.App) {
	app.Name = "Pocket DevOps CLI"
	app.Usage = "An tool for all of Pockets DevOps Commands"
	app.Version = "0.0.1"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Daniel Brooks",
			Email: "dbrooks@getpocket.com",
		},
		cli.Author{
			Name:  "Kaiser Shahid",
			Email: "kshahid@getpocket.com",
		},
		//If you work on this add your name!
	}
	app.Copyright = "(c) 2019 Read It Later, Inc."
}

func addCommands(app *cli.App) {
	app.Commands = []cli.Command{
		commands.FeatureCleanup(),
		commands.FeatureDeploy(),
		commands.FeatureDeployNotify(),
		commands.GithubDeployNotify(),
		commands.StackExists(),
		commands.CreateStack(),
		commands.EcsDeploy(),
		commands.UpToDate(),
	}
}
