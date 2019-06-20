package feature_deploy

import (
	"github.com/Pocket/ops-cli/internal/aws"
	"github.com/Pocket/ops-cli/internal/aws/cloudformation"
	"github.com/Pocket/ops-cli/internal/aws/ecs"
	"github.com/Pocket/ops-cli/internal/util"
)

func DeployBranch(parametersFile, templateFile, branchName, gitSHA, imageName string) {
	cloudformationClient := cloudformation.New()

	stackNameSuffix := util.DomainSafeString(branchName)

	settings := aws.NewSettingsParams(parametersFile, &stackNameSuffix, templateFile, &gitSHA, &branchName)

	if !cloudformationClient.StackExists(*settings.StackName) {
		cloudformationClient.CreateStack(settings)
	} else {
		ecsClient := ecs.New()
		ecsClient.DeployUpdate(settings.ECSCluster, &stackNameSuffix, &[]string{imageName})
	}
}
