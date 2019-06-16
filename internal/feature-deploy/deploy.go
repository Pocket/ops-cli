package feature_deploy

import (
	"github.com/Pocket/ops-cli/internal/aws"
)

func DeployBranch(prefix string, branchName string, imageName string) {
	stackName := stackNameFromBranchName(prefix, branchName)

	if aws.StackExists(stackName) {
		deployECS()
	} else {
		deployStack()
	}
}

func deployECS() {

}

func deployStack() {

}
