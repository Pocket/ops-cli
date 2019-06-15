package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"strings"
)

func GetActiveCloudFormationStackBranchesWithPrefix(prefix string) ([]string, error) {
	stacks, err := getActiveCloudFormationStacks()
	if err != nil {
		return nil, err
	}

	var activeCloudFormationBranches []string
	for _, stack := range stacks {
		if strings.HasPrefix(*stack.StackName, prefix) {
			for _, tag := range stack.Tags {
				if *tag.Key == "BranchName" {
					activeCloudFormationBranches = append(activeCloudFormationBranches, *tag.Value)
				}
			}
		}
	}

	return activeCloudFormationBranches, err
}

func getActiveCloudFormationStacks() ([]*cloudformation.Stack, error) {
	sess, err := session.NewSession(&aws.Config{})
	if err != nil {
		return nil, err
	}

	client := cloudformation.New(sess)
	stacks, err := client.DescribeStacks(&cloudformation.DescribeStacksInput{})
	if err != nil {
		return nil, err
	}

	return stacks.Stacks, nil
}
