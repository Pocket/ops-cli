package main

import (
	"fmt"
	"github.com/Pocket/ops-cli/internal/git"
	"github.com/urfave/cli"
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
		//If you work on this add your name!
	}
	app.Copyright = "(c) 2019 Read It Later, Inc."
}

func addCommands(app *cli.App)  {
	app.Commands = []cli.Command{
		{
			Name:    "active-branches",
			Aliases: []string{"ab"},
			Usage:   "Get a list of all the branches with commits in the last 8 days",
			Action:  func(c *cli.Context) error {
				branches, err := git.GetActiveBranches()
				if err != nil {
					return err
				}

				for _, branch := range branches {
					fmt.Println(branch.Name())
				}

				return nil
			},
		},
	}
}

