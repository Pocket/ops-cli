package feature_deploy

import (
	"github.com/Pocket/ops-cli/internal/github"
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

func (c *Client) NotifyDeployBranch(parametersFile, templateFile, branchName, gitSHA, slackWebhook string, githubUsername string, compareURL string, githubAccessToken string, githubOwner string, githubRepo string) error {
	stackNameSuffix := util.DomainSafeString(branchName)

	createdSettings := settings.NewSettingsParams(parametersFile, &templateFile, &gitSHA, &branchName, &stackNameSuffix)
	err := c.NotifyGithubDeployBranch(*createdSettings, githubAccessToken, githubOwner, githubRepo)
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

func (c *Client) NotifyGithubDeployBranch(createdSettings settings.Settings, accessToken string, owner string, repo string) error {
	return github.New(
		accessToken,
		nil,
	).NotifyGitHubDeploy(
		owner,
		repo,
		*createdSettings.BranchName,
		false,
		*createdSettings.GetBaseUrl(),
		*createdSettings.GetDeployUrl(),
	)
}
