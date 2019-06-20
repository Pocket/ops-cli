package cloudformation

import (
	"context"
	"github.com/Pocket/ops-cli/internal/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/dnaeon/go-vcr/recorder"
	"strings"
)

type Client struct {
	client        *cloudformation.Client
	clientContext context.Context
}

func New() *Client {
	client, clientContext := cloudFormationClient()
	return &Client{
		client:        client,
		clientContext: clientContext,
	}
}

func (c *Client) ActiveCloudFormationStackBranchesWithPrefix(prefix string) []string {
	stacks := c.activeCloudFormationStacks()

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

func (c *Client) DeleteStack(stackName string) {
	_, err := c.client.DeleteStackRequest(&cloudformation.DeleteStackInput{
		StackName: &stackName,
	}).Send(c.clientContext)
	if err != nil {
		panic("error deleting stack, " + err.Error())
	}
}

func (c *Client) StackExists(stackName string) bool {
	_, err := c.client.DescribeStacksRequest(&cloudformation.DescribeStacksInput{
		StackName: &stackName,
	}).Send(c.clientContext)

	if err != nil {
		return false
	}

	return true
}

func (c *Client) CreateStack(settings *aws.Settings) *string {
	createResponse, err := c.client.CreateStackRequest(&cloudformation.CreateStackInput{
		StackName:    settings.StackName,
		Tags:         settings.Tags,
		Parameters:   settings.Parameters,
		TemplateBody: settings.TemplateBody,
		OnFailure:    settings.OnFailure,
		Capabilities: settings.Capabilities,
	}).Send(c.clientContext)

	if err != nil {
		panic("error creating stack," + err.Error())
	}

	err = c.client.WaitUntilStackCreateComplete(c.clientContext, &cloudformation.DescribeStacksInput{
		StackName: createResponse.StackId,
	})

	if err != nil {
		panic("error waiting for stack complete, " + err.Error())
	}

	return createResponse.StackId
}

func (c *Client) CreateStackParams(paramFilePath string, stackName *string, templatefilePath string) *string {
	settings := aws.NewSettingsParams(paramFilePath, stackName, templatefilePath, nil)
	stackId := c.CreateStack(settings)
	return stackId
}

func (c *Client) activeCloudFormationStacks() []cloudformation.Stack {
	stacks, err := c.client.DescribeStacksRequest(&cloudformation.DescribeStacksInput{}).Send(c.clientContext)
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
