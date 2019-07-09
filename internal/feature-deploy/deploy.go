package feature_deploy

import (
	settings2 "github.com/Pocket/ops-cli/internal/settings"
	"github.com/Pocket/ops-cli/internal/slack"
	"github.com/Pocket/ops-cli/internal/util"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

func (c *Client) DeployBranch(parametersFile, templateFile, branchName, gitSHA, imageName string) {
	stackNameSuffix := util.DomainSafeString(branchName)

	settings := settings2.NewSettingsParams(parametersFile, &templateFile, &gitSHA, &branchName, &stackNameSuffix)

	if !c.cloudFormationClient.StackExists(*settings.StackName) {
		c.cloudFormationClient.CreateStack(&cloudformation.CreateStackInput{
			StackName:    settings.StackName,
			Tags:         settings.Tags,
			Parameters:   settings.Parameters,
			TemplateBody: settings.TemplateBody,
			OnFailure:    settings.OnFailure,
			Capabilities: settings.Capabilities,
		})
	} else {
		c.ecsClient.DeployUpdate(settings.ECSCluster, &stackNameSuffix, &[]string{imageName})
	}
}

func (c *Client) NotifyDeployBranch(parametersFile, templateFile, branchName, gitSHA, slackWebhook string, githubUsername string, compareURL string) {
	stackNameSuffix := util.DomainSafeString(branchName)

	settings := settings2.NewSettingsParams(parametersFile, &templateFile, &gitSHA, &branchName, &stackNameSuffix)

	text := "Completed deploy of <" + compareURL + "|" + gitSHA + "> - *" + branchName + "* by <https://github.com/" + githubUsername + "|" + githubUsername + ">"

	slackRequest := slack.NewSlackRequest(settings.SlackDeploySettings.Username, settings.SlackDeploySettings.Username, settings.SlackDeploySettings.Icon, text, "#36a64f", *settings.GetDeployUrl(), *settings.GetDeployUrl(), *settings.GetDeployUrl())
	err := slackRequest.SendSlackNotification(slackWebhook)
	if err != nil {
		panic("Error notifying slack: " + err.Error())
	}
}
