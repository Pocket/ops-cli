package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/dnaeon/go-vcr/recorder"
	"strings"
)

func ActiveCloudFormationStackBranchesWithPrefix(prefix string) []string {
	stacks := activeCloudFormationStacks()

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

	return activeCloudFormationBranches
}

func DeleteStack(stackName string) {
	_, err := cloudFormationClient().DeleteStackRequest(&cloudformation.DeleteStackInput{StackName: &stackName}).Send(context.TODO())
	if err != nil {
		panic("error deleting stack, " + err.Error())
	}
}

func activeCloudFormationStacks() ([]cloudformation.Stack) {
	stacks, err := cloudFormationClient().DescribeStacksRequest(&cloudformation.DescribeStacksInput{}).Send(context.TODO())
	//TODO: Do we need to paginate? (aws cli doesn't)
	if err != nil {
		panic("error getting stacks, " + err.Error())
	}

	return stacks.Stacks
}

func cloudFormationClient() *cloudformation.Client {
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	cfg.HTTPClient.Transport = awsRecorder()
	return cloudformation.New(cfg)
}

func awsRecorder() *recorder.Recorder {
	r, err := recorder.New("fixtures/cloudformation")
	if err != nil {
		panic(err)
	}
	defer r.Stop() // Make sure recorder is stopped once done with it

	return r
}
