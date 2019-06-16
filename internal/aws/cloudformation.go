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
	client, clientContext := cloudFormationClient()

	_, err := client.DeleteStackRequest(&cloudformation.DeleteStackInput{
		StackName: &stackName,
	}).Send(clientContext)
	if err != nil {
		panic("error deleting stack, " + err.Error())
	}
}

func StackExists(stackName string) bool {
	client, clientContext := cloudFormationClient()

	stack, err := client.DescribeStacksRequest(&cloudformation.DescribeStacksInput{
		StackName: &stackName,
	}).Send(clientContext)

	if err != nil {
		panic("error getting stack, " + err.Error())
	}

	return len(stack.Stacks) > 0
}

func CreateStack(settings *Settings) *string {
	client, clientContext := cloudFormationClient()

	createResponse, err := client.CreateStackRequest(&cloudformation.CreateStackInput{
		StackName:    settings.StackName,
		Tags:         settings.Tags,
		Parameters:   settings.Parameters,
		TemplateBody: settings.TemplateBody,
		OnFailure:    settings.OnFailure,
		Capabilities: settings.Capabilities,
	}).Send(clientContext)

	if err != nil {
		panic("error creating stack," + err.Error())
	}

	err = client.WaitUntilStackCreateComplete(clientContext, &cloudformation.DescribeStacksInput{
		StackName: createResponse.StackId,
	})

	if err != nil {
		panic("error waiting for stack complete, " + err.Error())
	}

	return createResponse.StackId
}

func CreateStackParams(paramFilePath string, stackName *string, templatefilePath string) *string {
	settings := NewSettingsParams(paramFilePath, stackName, templatefilePath, nil)
	stackId := CreateStack(settings)
	return stackId
}

func activeCloudFormationStacks() ([]cloudformation.Stack) {
	client, clientContext := cloudFormationClient()

	stacks, err := client.DescribeStacksRequest(&cloudformation.DescribeStacksInput{}).Send(clientContext)
	//TODO: Do we need to paginate? (aws cli doesn't)
	if err != nil {
		panic("error getting stacks, " + err.Error())
	}

	return stacks.Stacks
}

func cloudFormationClient() (*cloudformation.Client, context.Context) {
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	cfg.HTTPClient.Transport = awsCloudformationRecorder()
	return cloudformation.New(cfg), context.Background()
}

func awsCloudformationRecorder() *recorder.Recorder {
	r, err := recorder.New("_fixtures/cloudformation")
	if err != nil {
		panic(err)
	}
	defer r.Stop() // Make sure recorder is stopped once done with it

	return r
}
