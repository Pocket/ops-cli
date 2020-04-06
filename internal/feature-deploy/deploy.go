package feature_deploy

import (
	"github.com/Pocket/ops-cli/internal/github"
	"github.com/Pocket/ops-cli/internal/settings"
	"github.com/Pocket/ops-cli/internal/slack"
	"github.com/Pocket/ops-cli/internal/util"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

func (c *Client) DeployBranch(parametersFile, templateFile, branchName, gitSHA, imageName string, waitStable bool) {
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
		}, waitStable)
	} else {
		c.ecsClient.DeployUpdate(createdSettings.ECSCluster, &stackNameSuffix, &[]string{imageName}, waitStable)
	}
}

func (c *Client) NotifyDeployBranch(parametersFile, templateFile, branchName, gitSHA, slackWebhook string, githubUsername string, compareURL string, githubParams github.Params) error {
	stackNameSuffix := util.DomainSafeString(branchName)

	createdSettings := settings.NewSettingsParams(parametersFile, &templateFile, &gitSHA, &branchName, &stackNameSuffix)
	err := c.NotifyGithubDeployBranch(*createdSettings, githubParams)
	if err != nil {
		return err
	}
	return c.NotifySlack(*createdSettings, slackWebhook, githubUsername, compareURL)
}

func (c *Client) NotifySlack(createdSettings settings.Settings, slackWebHook string, githubUsername string, compareURL string) error {
	text := "Completed deploy of <" + compareURL + "|" + *createdSettings.GitSHA + "> - *" + *createdSettings.BranchName + "* by <https://github.com/" + githubUsername + "|" + githubUsername + ">"
	return slack.NewSlackRequest(
		createdSettings.SlackDeploySettings.Username,
		createdSettings.SlackDeploySettings.Channel,
		createdSettings.SlackDeploySettings.Icon,
		text,
		"#36a64f",
		*createdSettings.GetDeployUrl(),
		*createdSettings.GetDeployUrl(),
		*createdSettings.GetDeployUrl(),
	).SendSlackNotification(slackWebHook)
}

func (c *Client) NotifyGithubDeployBranch(createdSettings settings.Settings, githubParams github.Params) error {
	return github.New(&githubParams, nil).NotifyGitHubDeploy(
		*createdSettings.BranchName,
		false,
		*createdSettings.GetBaseUrl(),
		*createdSettings.GetDeployUrl(),
		*createdSettings.GetLogUrl(),
	)
}
