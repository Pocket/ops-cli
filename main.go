package main

import (
	"fmt"
	"github.com/Pocket/ops-cli/internal/aws"
	featureDeploy "github.com/Pocket/ops-cli/internal/feature-deploy"
	"github.com/Pocket/ops-cli/internal/git"
	"github.com/pkg/errors"
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
		//If you work on this add your name!
	}
	app.Copyright = "(c) 2019 Read It Later, Inc."
}

func addCommands(app *cli.App) {

	app.Commands = []cli.Command{
		{
			Name:    "active-branches",
			Aliases: []string{"ab"},
			Usage:   "Get a list of all the branches with commits in the last 8 days",
			Action: func(c *cli.Context) error {
				activeBranches, unactiveBranches := git.GetActiveAndUnactiveBranchNames()
				fmt.Println("----------------Active Branches----------------")
				for _, branch := range activeBranches {
					fmt.Println(branch)
				}
				fmt.Println()
				fmt.Println()
				fmt.Println()
				fmt.Println()
				fmt.Println("----------------UnActive Branches----------------")
				for _, branch := range unactiveBranches {
					fmt.Println(branch)
				}

				return nil
			},
		},
		{
			Name:    "feature-cleanup",
			Aliases: []string{"ac"},
			Usage:   "Cleanup all unactive stacks with the prefix",
			Action: func(c *cli.Context) error {
				featureDeploy.CleanUpBranches("WebFeatureDeploy-")
				return nil
			},
		},
		{
			Name:    "stack-exists",
			Aliases: []string{"se"},
			Usage:   "Check if a stack exists",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "stack-name, s",
					Usage:  "The stack to check for",
					EnvVar: "STACK_NAME",
				},
			},
			Action: func(c *cli.Context) error {
				if !aws.StackExists(c.String("stack-name")) {
					return errors.New("Stack not found")
				}
				fmt.Println("Stack found")
				return nil
			},
		},
		{
			Name:    "stack-deploy",
			Aliases: []string{"s"},
			Usage:   "Deploy a stack",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "stack-name, s",
					Usage:  "The stack name to deploy",
					EnvVar: "STACK_NAME",
				},
				cli.StringFlag{
					Name:   "param-file, p",
					Usage:  "The parameter file",
					EnvVar: "PARAM_FILE",
				},
				cli.StringFlag{
					Name:   "template-file, t",
					Usage:  "The template file",
					EnvVar: "TEMPLATE_FILE",
				},
			},
			Action: func(c *cli.Context) error {
				stackName := c.String("stack-name")
				aws.CreateStackParams(c.String("param-file"), &stackName, c.String("template-file"))
				return nil
			},
		},
		{
			Name:    "feature-deploy",
			Aliases: []string{"fd"},
			Usage:   "Deploy a feature branch",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "stack-name, s",
					Usage:  "The stack name to deploy",
					EnvVar: "STACK_NAME",
				},
				cli.StringFlag{
					Name:   "param-file, p",
					Usage:  "The parameter file",
					EnvVar: "PARAM_FILE",
				},
				cli.StringFlag{
					Name:   "template-file, t",
					Usage:  "The template file",
					EnvVar: "TEMPLATE_FILE",
				},
				cli.StringFlag{
					Name:   "git-sha, g",
					Usage:  "The git sha",
					EnvVar: "GIT_SHA",
				},
			},
			Action: func(c *cli.Context) error {

				return nil
			},
		},
		{
			Name:    "stacks-to-delete",
			Aliases: []string{"sd"},
			Usage:   "Get a list of stacks to delete",
			Action: func(c *cli.Context) error {
				stackBranchNames := featureDeploy.BranchesToDelete("WebFeatureDeploy-")

				for _, stackBranchName := range stackBranchNames {
					fmt.Println(stackBranchName)
				}

				return nil
			},
		},
		{
			Name:    "up-to-date",
			Aliases: []string{"ud"},
			Usage:   "Up to date",
			Action: func(c *cli.Context) error {
				git.UpToDateWithOriginMaster()

				return nil
			},
		},
	}
}
