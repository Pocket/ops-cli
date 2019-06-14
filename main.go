package main

import (
	"fmt"
	"github.com/urfave/cli"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
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
				branches, err := getActiveBranches()
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

func getActiveBranches() ([]*plumbing.Reference, error) {
	r, err := git.PlainOpen(".")
	if err != nil {
		return nil, err
	}

	refs, err := r.References()
	if err != nil {
		return nil, err
	}

	//We use "refs/remotes/origin/master" because plumbing.Master refers to the local master
	masterReference, err := r.Reference("refs/remotes/origin/master", false)
	if err != nil {
		return nil, err
	}

	eightDaysAgo := time.Now().AddDate(0,0,-8)

	var activeBranches []*plumbing.Reference

	refs.ForEach(func(ref *plumbing.Reference) error {
		// The HEAD is omitted in a `git show-ref` so we ignore the symbolic
		// references, the HEAD
		if ref.Type() == plumbing.SymbolicReference {
			return nil
		}

		commit, err := r.CommitObject(ref.Hash())
		if err != nil {
			return nil
		}

		if commit.Hash == masterReference.Hash() {
			return nil
		}

		lastCommitTime := commit.Author.When
		if lastCommitTime.After(eightDaysAgo) {
			activeBranches = append(activeBranches, ref)
		}

		return nil
	})

	return activeBranches, nil
}