package feature_deploy

import (
	"github.com/Pocket/ops-cli/internal/aws"
	"github.com/Pocket/ops-cli/internal/git"
	"github.com/Pocket/ops-cli/internal/util"
)

func GetStacksToDelete() ([]string, error) {
	stackBranchNames, err := aws.GetActiveCloudFormationStackBranchesWithPrefix("WebFeatureDeploy-")
	if err != nil {
		return nil, err
	}

	activeBranchNames, unactiveBranchNames := git.GetActiveAndUnactiveBranchNames()

	var branchesToClean []string

	for _, branchName := range unactiveBranchNames {
		if util.StringInSlice(branchName, stackBranchNames) {
			branchesToClean = append(branchesToClean, branchName)
		}
	}

	for _, branchName := range stackBranchNames {
		if !util.StringInSlice(branchName, activeBranchNames) && !util.StringInSlice(branchName, unactiveBranchNames) {
			branchesToClean = append(branchesToClean, branchName)
		}
	}

	branchesToClean = util.RemoveDuplicatesFromSlice(branchesToClean)

	return branchesToClean, nil
}
