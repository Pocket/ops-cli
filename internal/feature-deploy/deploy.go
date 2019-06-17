package feature_deploy

import (
	"github.com/Pocket/ops-cli/internal/aws"
	"github.com/Pocket/ops-cli/internal/aws/cloudformation"
	"github.com/Pocket/ops-cli/internal/util"
)

func DeployBranch(parametersFile, templateFile, branchName, gitSHA, imageName string) {
	cloudformationClient := cloudformation.New()

	stackNameSuffix := util.DomainSafeString(branchName)

	settings := aws.NewSettingsParams(parametersFile, &stackNameSuffix, templateFile, &gitSHA)

	if cloudformationClient.StackExists(*settings.StackName) {
		cloudformationClient.CreateStack(settings)
	} else {
		deployECS(settings, imageName)
	}
}

func deployECS(settings *aws.Settings, imageName string) {

}

func deployStack() {

}
