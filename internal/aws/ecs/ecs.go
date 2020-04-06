package ecs

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"net/http"
	"strings"
)

type Client struct {
	client        *ecs.Client
	clientContext context.Context
}

func (c *Client) setTransport(transport *http.Transport) {
	c.client.Config.HTTPClient.Transport = transport
}

func New() *Client {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	return &Client{
		client:        ecs.New(cfg),
		clientContext: context.Background(),
	}
}

func (c *Client) SetTransport(transport http.RoundTripper) {
	c.client.Config.HTTPClient.Transport = transport
}

func (c *Client) DeployUpdate(clusterName *string, serviceName *string, imageNames *[]string, waitStable bool) {

	deployment, err := c.getLatestDeployment(clusterName, serviceName)
	if err != nil {
		panic("error getting the latest deployment, " + err.Error())
	}
	taskDefinition, err := c.registerTaskDefinition(deployment.TaskDefinition, imageNames)

	if err != nil {
		panic("error registering the latest deployment, " + err.Error())
	}

	err = c.updateService(clusterName, serviceName, nil, &taskDefinition)

	if err != nil {
		panic("error updating the service, " + err.Error())
	}

	if waitStable {
		err = c.wait(clusterName, serviceName, &taskDefinition)
		if err != nil {
			panic("error waiting for the service, " + err.Error())
		}
	}
}

// RegisterTaskDefinition updates the existing task definition's image.
func (c *Client) registerTaskDefinition(task *string, images *[]string) (string, error) {
	output, err := c.client.DescribeTaskDefinitionRequest(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: task,
	}).Send(c.clientContext)

	if err != nil {
		return "", err
	}

	var definitions []ecs.ContainerDefinition

	for _, d := range output.TaskDefinition.ContainerDefinitions {
		for _, image := range *images {
			imageName := strings.Split(image, ":")[0]
			if strings.HasPrefix(*d.Image, imageName) {
				d.Image = &image
			}
			definitions = append(definitions, d)
		}
	}
	input := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: definitions,
		Family:               output.TaskDefinition.Family,
		NetworkMode:          output.TaskDefinition.NetworkMode,
		PlacementConstraints: output.TaskDefinition.PlacementConstraints,
		TaskRoleArn:          output.TaskDefinition.TaskRoleArn,
		Volumes:              output.TaskDefinition.Volumes,
		Memory:               output.TaskDefinition.Memory,
		Cpu:                  output.TaskDefinition.Cpu,
		ExecutionRoleArn:     output.TaskDefinition.ExecutionRoleArn,
		Tags:                 output.Tags,
	}

	resp, err := c.client.RegisterTaskDefinitionRequest(input).Send(c.clientContext)
	if err != nil {
		return "", err
	}
	return *resp.TaskDefinition.TaskDefinitionArn, nil
}

// UpdateService updates the service to use the new task definition.
func (c *Client) updateService(cluster, service *string, count *int64, arn *string) error {
	input := &ecs.UpdateServiceInput{
		Cluster: cluster,
		Service: service,
	}
	if count != nil && *count != -1 {
		input.DesiredCount = count
	}
	if arn != nil {
		input.TaskDefinition = arn
	}
	_, err := c.client.UpdateServiceRequest(input).Send(c.clientContext)
	return err
}

// Wait waits for the service to finish being updated.
func (c *Client) wait(cluster, service, arn *string) error {
	return c.client.WaitUntilServicesStable(c.clientContext, &ecs.DescribeServicesInput{
		Cluster:  cluster,
		Services: []string{*service},
	}, aws.WithWaiterMaxAttempts(80))
}

func (c *Client) getLatestDeployment(cluster, service *string) (*ecs.Deployment, error) {
	input := &ecs.DescribeServicesInput{
		Cluster:  cluster,
		Services: []string{*service},
	}
	output, err := c.client.DescribeServicesRequest(input).Send(c.clientContext)
	if err != nil {
		return nil, err
	}

	if len(output.Services) == 0 {
		return nil, errors.New("No active ecs services")
	}

	ds := output.Services[0].Deployments
	for _, d := range ds {
		if *d.TaskDefinition == *output.Services[0].TaskDefinition {
			return &d, nil
		}
	}
	return nil, nil
}
