package main

import (
	"fmt"
	"github.com/Pocket/ops-cli/internal/aws/cloudformation"
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
			Name:    "feature-cleanup",
			Aliases: []string{"ac"},
			Usage:   "Cleanup all unactive stacks with the prefix",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "stack-prefix, p",
					Usage:  "The stack prefix for these deployments (WebFeatureDeploy-)",
					EnvVar: "STACK_PREFIX",
				},
				cli.BoolTFlag{
					Name:   "dry-run, d",
					Usage:  "Perform a dry run, true by default",
					EnvVar: "DRY_RUN",
				},
			},
			Action: func(c *cli.Context) error {
				stackPrefix := c.String("stack-prefix")
				if c.BoolT("dry-run") {
					stackBranchNames := featureDeploy.BranchesToDelete(stackPrefix)

					for _, stackBranchName := range stackBranchNames {
						fmt.Println(stackBranchName)
					}

					return nil
				}

				featureDeploy.CleanUpBranches(stackPrefix)
				return nil
			},
		},
		{
			Name:    "feature-deploy",
			Aliases: []string{"fd"},
			Usage:   "Deploy a feature branch",
			Flags: []cli.Flag{
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
				cli.StringFlag{
					Name:   "branch-name, b",
					Usage:  "The branch name",
					EnvVar: "BRANCH_NAME",
				},
			},
			Action: func(c *cli.Context) error {

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
				cloudformationClient := cloudformation.New()

				if !cloudformationClient.StackExists(c.String("stack-name")) {
					return errors.New("Stack not found")
				}
				fmt.Println("Stack found")
				return nil
			},
		},
		{
			Name:    "create-stack",
			Aliases: []string{"cs"},
			Usage:   "Create a stack",
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
				cloudformationClient := cloudformation.New()

				stackName := c.String("stack-name")
				cloudformationClient.CreateStackParams(c.String("param-file"), &stackName, c.String("template-file"))
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
