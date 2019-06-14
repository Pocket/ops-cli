package git

import (
	"github.com/Pocket/ops-cli/internal/util"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"strings"
	"time"
)

func GetActiveAndUnactiveBranchNames() ([]string, []string, error) {
	activeBranches, unactiveBranches, err := getActiveAndUnactiveBranches()
	if err != nil {
		return nil, nil, err
	}
	return getBranchShortNames(activeBranches), getBranchShortNames(unactiveBranches), nil
}

/*
 * Gets the short names of the branches
 */
func getBranchShortNames(branches []*plumbing.Reference) []string {
	var branchShortNames []string
	for _, branch := range branches {
		branchShortNames = append(branchShortNames, strings.TrimPrefix(branch.Name().Short(), "origin/"))
	}
	return util.RemoveDuplicatesFromSlice(branchShortNames)
}

/**
 * Gets the active and unactive branch refs
 */
func getActiveAndUnactiveBranches() ([]*plumbing.Reference, []*plumbing.Reference, error) {
	r, err := git.PlainOpen(".")
	if err != nil {
		return nil, nil, err
	}

	refs, err := r.References()
	if err != nil {
		return nil, nil, err
	}

	//We use "refs/remotes/origin/master" because plumbing.Master refers to the local master
	masterReference, err := r.Reference("refs/remotes/origin/master", false)
	if err != nil {
		return nil, nil, err
	}

	eightDaysAgo := time.Now().AddDate(0,0,-8)

	var activeBranches []*plumbing.Reference
	var unactiveBranches []*plumbing.Reference

	err = refs.ForEach(func(ref *plumbing.Reference) error {
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
		} else {
			unactiveBranches = append(unactiveBranches, ref)
		}

		return nil
	})

	return activeBranches, unactiveBranches, err
}
