package feature_deploy

import (
	"github.com/Pocket/ops-cli/internal/aws"
	"github.com/Pocket/ops-cli/internal/aws/cloudformation"
	"github.com/Pocket/ops-cli/internal/aws/ecs"
	"github.com/Pocket/ops-cli/internal/slack"
	"github.com/Pocket/ops-cli/internal/util"
)

func DeployBranch(parametersFile, templateFile, branchName, gitSHA, imageName string, slackWebhook string) {
	cloudformationClient := cloudformation.New()

	stackNameSuffix := util.DomainSafeString(branchName)

	settings := aws.NewSettingsParams(parametersFile, &stackNameSuffix, templateFile, &gitSHA)

	if !cloudformationClient.StackExists(*settings.StackName) {
		cloudformationClient.CreateStack(settings)
	} else {
		ecsClient := ecs.New()
		ecsClient.DeployUpdate(settings.ECSCluster, &stackNameSuffix, &[]string{imageName})
	}

	text := "Completed deploy of " + gitSHA + " - *" + branchName + "*"

	slackRequest := slack.NewSlackRequest("Buster", "#log-feature-deploys", ":ship:", text, "#36a64f", *settings.GetDeployUrl(), *settings.GetDeployUrl(), *settings.GetDeployUrl())
	err := slackRequest.SendSlackNotification(slackWebhook)
	if err != nil {
		panic("Error notifying slack: " + err.Error())
	}
}
