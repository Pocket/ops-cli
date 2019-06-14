package main

import (
	"fmt"
	"github.com/Pocket/ops-cli/internal/aws"
	feature_deploy "github.com/Pocket/ops-cli/internal/feature-deploy"
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
				activeBranches, unactiveBranches, err := git.GetActiveAndUnactiveBranchNames()
				if err != nil {
					return err
				}

				for _, branch := range activeBranches {
					fmt.Println(branch)
				}
				fmt.Println()
				fmt.Println()
				fmt.Println()
				fmt.Println()
				for _, branch := range unactiveBranches {
					fmt.Println(branch)
				}

				return nil
			},
		},
		{
			Name:    "active-cloudformation",
			Aliases: []string{"ac"},
			Usage:   "Get a list of all the active cloudformation stacks",
			Action:  func(c *cli.Context) error {
				stacks, err := aws.GetActiveCloudFormationStackBranchesWithPrefix("WebFeatureDeploy-")
				if err != nil {
					return err
				}

				for _, stack := range stacks {
					fmt.Println(stack)
				}

				return nil
			},
		},
		{
			Name:    "stacks-to-delete",
			Aliases: []string{"sd"},
			Usage:   "Get a list of stacks to delete",
			Action:  func(c *cli.Context) error {
				stackBranchNames, err := feature_deploy.GetStacksToDelete()
				if err != nil {
					return err
				}

				for _, stackBranchName := range stackBranchNames {
					fmt.Println(stackBranchName)
				}

				return nil
			},
		},
	}
}

