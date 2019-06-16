package feature_deploy

import (
	"github.com/Pocket/ops-cli/internal/aws"
	"github.com/Pocket/ops-cli/internal/git"
	"github.com/Pocket/ops-cli/internal/util"
)

func CleanUpBranches(prefix string) {
	branchesToDelete := BranchesToDelete(prefix)

	for _, branchName := range branchesToDelete {
		aws.DeleteStack(stackNameFromBranchName(prefix, branchName))
		//TODO: Notify Slack
	}
}

func BranchesToDelete(prefix string) []string {
	stackBranchNames := aws.ActiveCloudFormationStackBranchesWithPrefix(prefix)

	activeBranchNames, unactiveBranchNames := git.GetActiveAndUnactiveBranchNames()

	var branchesToClean []string

	//If an unactive branch has a stack lets set it up to delete
	for _, branchName := range unactiveBranchNames {
		if util.StringInSlice(branchName, stackBranchNames) {
			branchesToClean = append(branchesToClean, branchName)
		}
	}

	//If a stack is active, and is not in the activeBranches and not in the unactive branches lets set it up to delete
	//The branch was deleted
	for _, branchName := range stackBranchNames {
		if !util.StringInSlice(branchName, activeBranchNames) && !util.StringInSlice(branchName, unactiveBranchNames) {
			branchesToClean = append(branchesToClean, branchName)
		}
	}

	branchesToClean = util.RemoveDuplicatesFromSlice(branchesToClean)

	return branchesToClean
}
