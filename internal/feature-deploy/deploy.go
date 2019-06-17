package feature_deploy

import (
	"github.com/Pocket/ops-cli/internal/aws/cloudformation"
)

func DeployBranch(prefix string, branchName string, imageName string) {
	cloudformationClient := cloudformation.New()

	stackName := stackNameFromBranchName(prefix, branchName)

	if cloudformationClient.StackExists(stackName) {
		deployECS()
	} else {
		deployStack()
	}
}

func deployECS() {
}

func deployStack() {

}
