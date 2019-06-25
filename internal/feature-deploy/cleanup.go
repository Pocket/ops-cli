package feature_deploy

import (
	"github.com/Pocket/ops-cli/internal/aws/cloudformation"
	"github.com/Pocket/ops-cli/internal/git"
	"github.com/Pocket/ops-cli/internal/slack"
	"github.com/Pocket/ops-cli/internal/util"
)

func CleanUpBranches(prefix string, slackWebhook string) {
	client := cloudformation.New()

	branchesToDelete := BranchesToDelete(prefix)

	for _, branchName := range branchesToDelete {
		client.DeleteStack(stackNameFromBranchName(prefix, branchName))

		text := "Cleaned up " + branchName + "*"

		slackRequest := slack.NewSlackRequestText("Damage Control", "#log-feature-deploys", ":cleanup:", text)
		err := slackRequest.SendSlackNotification(slackWebhook)
		if err != nil {
			panic("Error notifying slack: " + err.Error())
		}
	}
}

func BranchesToDelete(prefix string) []string {
	client := cloudformation.New()

	stackBranchNames := client.ActiveCloudFormationStackBranchesWithPrefix(prefix)

	activeBranchNames, unactiveBranchNames := git.GetActiveAndUnactiveBranchNames()

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

	return branchesToClean
}
