package git

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"time"
)

func GetActiveBranches() ([]*plumbing.Reference, error) {
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
