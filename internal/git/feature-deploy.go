package git

import (
	"github.com/Pocket/ops-cli/internal/util"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"strings"
	"time"
)

func GetActiveAndUnactiveBranchNames(olderThanDate time.Time, mainBranch *string) ([]string, []string) {
	activeBranches, unactiveBranches := getActiveAndUnactiveBranches(olderThanDate, *mainBranch)
	return getBranchShortNames(activeBranches), getBranchShortNames(unactiveBranches)
}

/**
 * Gets the active and unactive branch refs
 */
func getActiveAndUnactiveBranches(olderThanDate time.Time, mainBranch *string) ([]*plumbing.Reference, []*plumbing.Reference) {
	r := repo(".")

	refs, err := r.References()
	if err != nil {
		panic(err)
	}

	mainReference := main(r, *mainBranch)

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

		if commit.Hash == mainReference.Hash() {
			return nil
		}

		lastCommitTime := commit.Author.When
		if lastCommitTime.After(olderThanDate) {
			activeBranches = append(activeBranches, ref)
		} else {
			unactiveBranches = append(unactiveBranches, ref)
		}

		return nil
	})

	return activeBranches, unactiveBranches
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
