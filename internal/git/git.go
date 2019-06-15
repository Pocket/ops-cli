package git

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func getRepo(repoPath string) (*git.Repository) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		panic(err)
	}

	return r
}

func getMaster(r *git.Repository) *plumbing.Reference {
	//We use "refs/remotes/origin/master" because plumbing.Master refers to the local master
	masterReference, err := r.Reference("refs/remotes/origin/master", false)
	if err != nil {
		panic(err)
	}
	return masterReference
}
