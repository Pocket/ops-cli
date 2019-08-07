package commands

import (
	"errors"
	"fmt"
	"time"
	"github.com/Pocket/ops-cli/internal/aws/cloudformation"
	"github.com/Pocket/ops-cli/internal/aws/ecs"
	"github.com/Pocket/ops-cli/internal/git"
	"github.com/Pocket/ops-cli/internal/github"
	"gopkg.in/urfave/cli.v1"
	featureDeploy "github.com/Pocket/ops-cli/internal/feature-deploy"
)

func FeatureCleanup() cli.Command {
	return cli.Command{
		Name:    "feature-cleanup",
		Aliases: []string{"ac"},
		Usage:   "Cleanup all unactive stacks with the prefix",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "param-file, p",
				Usage:  "The parameter file",
				EnvVar: "PARAM_FILE",
			},
			cli.BoolTFlag{
				Name:   "dry-run, d",
				Usage:  "Perform a dry run, true by default",
				EnvVar: "DRY_RUN",
			},
			cli.StringFlag{
				Name:   "slack-webhook, s",
				Usage:  "The slack webhook",
				EnvVar: "SLACK_WEBHOOK",
			},
			cli.IntFlag{
				Name:   "days-old, do",
				Usage:  "The days old",
				EnvVar: "DAYS_OLD",
				Value:  8,
			},
			cli.StringFlag{
				Name:   "github-token, ght",
				Usage:  "The github token",
				EnvVar: "GITHUB_TOKEN",
			},
			cli.StringFlag{
				Name:   "github-owner, gho",
				Usage:  "The github owner",
				EnvVar: "GITHUB_OWNER",
			},
			cli.StringFlag{
				Name:   "github-repo, ghr",
				Usage:  "The github repo",
				EnvVar: "GITHUB_REPO",
			},
		},
		Action: func(c *cli.Context) error {
			stackPrefix := c.String("stack-prefix")
			olderThanDate := time.Now().AddDate(0, 0, -c.Int("days-old"))
			client := featureDeploy.New()

			if c.BoolT("dry-run") {
				stackBranchNames := client.StacksToDelete(stackPrefix, olderThanDate)

				for _, stackBranchName := range stackBranchNames {
					fmt.Println(stackBranchName)
				}

				return nil
			}
			client.CleanUpBranches(
				c.String("param-file"),
				c.String("slack-webhook"),
				olderThanDate,
				github.Params{
					AccessToken: c.String("github-token"),
					Owner:       c.String("github-owner"),
					Repo:        c.String("github-repo"),
				},
			)
			return nil
		},
	}
}

func FeatureDeploy() cli.Command {
	return cli.Command{
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
			cli.StringFlag{
				Name:   "image-name, i",
				Usage:  "The image name",
				EnvVar: "IMAGE_NAME",
			},
		},
		Action: func(c *cli.Context) error {
			featureDeploy.New().DeployBranch(c.String("param-file"), c.String("template-file"), c.String("branch-name"), c.String("git-sha"), c.String("image-name"))
			return nil
		},
	}
}

func FeatureDeployNotify() cli.Command {
	return cli.Command{
		Name:    "feature-deploy-notify",
		Aliases: []string{"fdn"},
		Usage:   "Notify about a feature branch",
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
			cli.StringFlag{
				Name:   "slack-webhook, s",
				Usage:  "The image name",
				EnvVar: "SLACK_WEBHOOK",
			},
			cli.StringFlag{
				Name:   "github-username, u",
				Usage:  "The github username",
				EnvVar: "GITHUB_USERNAME",
			},
			cli.StringFlag{
				Name:   "github-compare-url, cu",
				Usage:  "The github compare-url",
				EnvVar: "GITHUB_COMPARE_URL",
			},
			cli.StringFlag{
				Name:   "github-token, ght",
				Usage:  "The github token",
				EnvVar: "GITHUB_TOKEN",
			},
			cli.StringFlag{
				Name:   "github-owner, gho",
				Usage:  "The github owner",
				EnvVar: "GITHUB_OWNER",
			},
			cli.StringFlag{
				Name:   "github-repo, ghr",
				Usage:  "The github repo",
				EnvVar: "GITHUB_REPO",
			},
		},
		Action: func(c *cli.Context) error {
			return featureDeploy.New().NotifyDeployBranch(
				c.String("param-file"),
				c.String("template-file"),
				c.String("branch-name"),
				c.String("git-sha"),
				c.String("slack-webhook"),
				c.String("github-username"),
				c.String("github-compare-url"),
				github.Params{
					AccessToken: c.String("github-token"),
					Owner:       c.String("github-owner"),
					Repo:        c.String("github-repo"),
				},
			)
		},
	}
}

func GithubDeployNotify() cli.Command {
	return cli.Command{
		Name:    "github-deploy-notify",
		Aliases: []string{"ghdn"},
		Usage:   "Notify github of a deployment",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "branch-name, b",
				Usage:  "The branch name",
				EnvVar: "BRANCH_NAME",
			},
			cli.StringFlag{
				Name:   "environment, env",
				Usage:  "The environment",
				EnvVar: "ENVIRONMENT",
			},
			cli.StringFlag{
				Name:   "url, u",
				Usage:  "The url of the deploy",
				EnvVar: "URL",
			},
			cli.BoolFlag{
				Name:   "production, prod",
				Usage:  "Production environment",
				EnvVar: "PRODUCTION_ENVIRONMENT",
			},
			cli.StringFlag{
				Name:   "github-token, ght",
				Usage:  "The github token",
				EnvVar: "GITHUB_TOKEN",
			},
			cli.StringFlag{
				Name:   "github-owner, gho",
				Usage:  "The github owner",
				EnvVar: "GITHUB_OWNER",
			},
			cli.StringFlag{
				Name:   "github-repo, ghr",
				Usage:  "The github repo",
				EnvVar: "GITHUB_REPO",
			},
			cli.StringFlag{
				Name: "log-url, lurl",
				Usage: "The AWS log group url",
				EnvVar: "LOG_URL",
			},
		},
		Action: func(c *cli.Context) error {
			return github.New(&github.Params{
				AccessToken: c.String("github-token"),
				Owner:       c.String("github-owner"),
				Repo:        c.String("github-repo"),
			},
				nil,
			).NotifyGitHubDeploy(
				c.String("branch-name"),
				c.Bool("production"),
				c.String("environment"),
				c.String("url"),
				c.String("log-url"),
			)
		},
	}
}

func StackExists() cli.Command {
	return cli.Command{
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
	}
}

func CreateStack() cli.Command {
	return cli.Command{
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
	}
}

func EcsDeploy() cli.Command {
	return cli.Command{
		Name:    "ecs-deploy",
		Aliases: []string{"ed"},
		Usage:   "ECS Deploy",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "cluster-name, c",
				Usage:  "The cluster name",
				EnvVar: "CLUSTER_NAME",
			},
			cli.StringFlag{
				Name:   "service-name, s",
				Usage:  "The service name",
				EnvVar: "SERVICE_NAME",
			},
			cli.StringSliceFlag{
				Name:   "image-names, i",
				Usage:  "The image names",
				EnvVar: "IMAGE_NAMES",
			},
		},
		Action: func(c *cli.Context) error {
			ecsClient := ecs.New()
			clusterName := c.String("cluster-name")
			serviceName := c.String("service-name")
			imageNames := c.StringSlice("image-names")
			ecsClient.DeployUpdate(&clusterName, &serviceName, &imageNames)
			return nil
		},
	}
}

func UpToDate() cli.Command {
	return cli.Command{
		Name:    "up-to-date",
		Aliases: []string{"ud"},
		Usage:   "Up to date",
		Action: func(c *cli.Context) error {
			git.UpToDateWithOriginMaster()

			return nil
		},
	}
}
