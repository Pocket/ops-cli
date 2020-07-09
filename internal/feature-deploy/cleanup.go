package feature_deploy

import (
	"github.com/Pocket/ops-cli/internal/git"
	"github.com/Pocket/ops-cli/internal/github"
	"github.com/Pocket/ops-cli/internal/settings"
	"github.com/Pocket/ops-cli/internal/slack"
	"github.com/Pocket/ops-cli/internal/util"
	"time"
)

func (c *Client) CleanUpBranches(paramFilePath string, slackWebHook string, olderThanDate time.Time, githubParams github.Params, mainBranch string) {
	stackPrefix := *settings.NewSettingsParams(paramFilePath, nil, nil, nil, nil).StackPrefix

	branchesToDelete := c.StacksToDelete(stackPrefix, olderThanDate, &mainBranch)
	for _, branchName := range branchesToDelete {
		formattedBranchName := util.DomainSafeString(branchName)
		branchSettings := settings.NewSettingsParams(paramFilePath, nil, nil, &branchName, &formattedBranchName)
		c.CleanUpBranch(branchSettings, &slackWebHook, &githubParams)
	}
}

func (c *Client) CleanUpBranch(settings *settings.Settings, slackWebHook *string, githubParams *github.Params) {
	err := c.cloudWatchLogsClient.ExportLogGroupAndWait(*settings.LogGroupPrefix+*settings.FormattedBranchName, *settings.ArchiveLogsBucketName)
	if err != nil {
		panic("There was an error backing up the feature logs: " + err.Error())
	}

	c.cloudFormationClient.DeleteStack(*settings.StackName)

	if githubParams != nil {
		c.GithubNotify(settings, githubParams)
	}

	if slackWebHook != nil {
		c.SlackNotify(settings, slackWebHook)
	}
}

func (c *Client) SlackNotify(settings *settings.Settings, slackWebHook *string) {
	text := "Cleaned up *" + *settings.BranchName + "*"

	slackRequest := slack.NewSlackRequestText(settings.SlackCleanUpSettings.Username, settings.SlackCleanUpSettings.Channel, settings.SlackCleanUpSettings.Icon, text)
	err := slackRequest.SendSlackNotification(*slackWebHook)
	if err != nil {
		panic("Error notifying slack: " + err.Error())
	}
}

func (c *Client) GithubNotify(settings *settings.Settings, githubParams *github.Params) {
	err := github.New(githubParams, nil).DeleteDeployment(*settings.BranchName, *settings.GetBaseUrl())
	if err != nil {
		panic("Error notifying github: " + err.Error())
	}
}

func (c *Client) StacksToDelete(prefix string, olderThanDate time.Time, mainBranch *string) []string {
	stackBranchNames := c.cloudFormationClient.ActiveCloudFormationStackBranchesWithPrefix(prefix)

	activeBranchNames, unactiveBranchNames := git.GetActiveAndUnactiveBranchNames(olderThanDate)

	var branchesToClean []string

	//If an unactive branch has a stack lets set it up to delete
	for _, branchName := range unactiveBranchNames {
		if util.StringInSlice(branchName, stackBranchNames) {
			branchesToClean = append(branchesToClean, branchName)
		}
	}

	//If a stack is active, and is not in the activeBranches and not in the unactive branches lets set it up to delete
	//The branch was deleted
	for _, branchName := range stackBranchNames {
		if !util.StringInSlice(branchName, activeBranchNames) && !util.StringInSlice(branchName, unactiveBranchNames) {
			branchesToClean = append(branchesToClean, branchName)
		}
	}

	branchesToClean = util.RemoveDuplicatesFromSlice(branchesToClean)

	// If we provide a main branch as a param, exclude that branch from the branches to clean
	if mainBranch != nil {
		branchesToClean = util.ExcludeMainBranchFromSlice(branchesToClean, mainBranch)
	}

	return branchesToClean
}
