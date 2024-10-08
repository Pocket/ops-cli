package git

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func repo(repoPath string) (*git.Repository) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		panic(err)
	}

	return r
}

func main(r *git.Repository, mainBranch string) *plumbing.Reference {
	//We use "refs/remotes/origin/master" because plumbing.Master refers to the local master
	mainRefName := plumbing.ReferenceName("refs/remotes/origin/" + mainBranch)
	mainReference, err := r.Reference(mainRefName, false)
    if err != nil {
        panic(err)
    }

	return mainReference
}
