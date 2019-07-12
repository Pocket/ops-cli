package feature_deploy

import (
	"github.com/Pocket/ops-cli/internal/settings"
	"github.com/Pocket/ops-cli/internal/slack"
	"github.com/Pocket/ops-cli/internal/util"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

func (c *Client) DeployBranch(parametersFile, templateFile, branchName, gitSHA, imageName string) {
	stackNameSuffix := util.DomainSafeString(branchName)

	createdSettings := settings.NewSettingsParams(parametersFile, &templateFile, &gitSHA, &branchName, &stackNameSuffix)

	if !c.cloudFormationClient.StackExists(*createdSettings.StackName) {
		c.cloudFormationClient.CreateStack(&cloudformation.CreateStackInput{
			StackName:    createdSettings.StackName,
			Tags:         createdSettings.Tags,
			Parameters:   createdSettings.Parameters,
			TemplateBody: createdSettings.TemplateBody,
			OnFailure:    createdSettings.OnFailure,
			Capabilities: createdSettings.Capabilities,
		})
	} else {
		c.ecsClient.DeployUpdate(createdSettings.ECSCluster, &stackNameSuffix, &[]string{imageName})
	}
}

func (c *Client) NotifyDeployBranch(parametersFile, templateFile, branchName, gitSHA, slackWebHook string, githubUsername string, compareURL string) {
	stackNameSuffix := util.DomainSafeString(branchName)

	createdSettings := settings.NewSettingsParams(parametersFile, &templateFile, &gitSHA, &branchName, &stackNameSuffix)

	text := "Completed deploy of <" + compareURL + "|" + gitSHA + "> - *" + branchName + "* by <https://github.com/" + githubUsername + "|" + githubUsername + ">"

	slackRequest := slack.NewSlackRequest(createdSettings.SlackDeploySettings.Username, createdSettings.SlackDeploySettings.Channel, createdSettings.SlackDeploySettings.Icon, text, "#36a64f", *createdSettings.GetDeployUrl(), *createdSettings.GetDeployUrl(), *createdSettings.GetDeployUrl())
	err := slackRequest.SendSlackNotification(slackWebHook)
	if err != nil {
		panic("Error notifying slack: " + err.Error())
	}
}
