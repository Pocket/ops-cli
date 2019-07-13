package github

import (
	"context"
	"errors"
	"github.com/google/go-github/v26/github"
	"golang.org/x/oauth2"
	"net/http"
)

type Client struct {
	client        *github.Client
	clientContext context.Context
}

func New(accessToken string, transport http.RoundTripper) *Client {
	clientContext := context.TODO()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: accessToken,
		},
	)

	var tc *http.Client
	if transport != nil {
		tc = &http.Client{
			Transport: transport,
		}
		clientContext = context.WithValue(clientContext, oauth2.HTTPClient, tc)
	}

	tc = oauth2.NewClient(clientContext, ts)

	client := github.NewClient(tc)

	return &Client{
		client:        client,
		clientContext: clientContext,
	}
}

func (c *Client) CreateDeployment(owner string, repo string, branchName string, productionEnvironment bool, environment string, environmentURL string) error {
	autoMerge := false
	transientEnvironment := !productionEnvironment
	requiredContexts := []string{}

	deployment, _, err := c.client.Repositories.CreateDeployment(c.clientContext, owner, repo, &github.DeploymentRequest{
		Ref:                   &branchName,
		AutoMerge:             &autoMerge,
		Environment:           &environment,
		TransientEnvironment:  &transientEnvironment,
		ProductionEnvironment: &productionEnvironment,
		RequiredContexts:      &requiredContexts,
	})

	if err != nil {
		return errors.New("Error creating github deployment: " + err.Error())
	}

	status := "success"
	_, _, err = c.client.Repositories.CreateDeploymentStatus(c.clientContext, owner, repo, *deployment.ID, &github.DeploymentStatusRequest{
		State:          &status,
		Environment:    &environment,
		EnvironmentURL: &environmentURL,
	})

	if err != nil {
		return errors.New("Error creating active github deployment status: " + err.Error())
	}

	return nil
}

func (c *Client) DeleteDeployment(owner string, repo string, ref string, environment string) error {
	return c.UpdateDeploymentStatusForAllMatchingDeploys(owner, repo, ref, environment, "inactive")
}

func (c *Client) GetDeployments(owner string, repo string, ref string, environment string) ([]*github.Deployment, error) {
	deployments, _, err := c.client.Repositories.ListDeployments(c.clientContext, owner, repo, &github.DeploymentsListOptions{
		Ref:         ref,
		Environment: environment,
	})

	if err != nil {
		return nil, errors.New("Error listing deployments: " + err.Error())
	}

	return deployments, nil
}

func (c *Client) UpdateDeploymentStatus(owner string, repo string, status string, deploymentId int64) error {
	_, _, err := c.client.Repositories.CreateDeploymentStatus(c.clientContext, owner, repo, deploymentId, &github.DeploymentStatusRequest{
		State: &status,
	})
	if err != nil {
		return errors.New("Error creating github deployment status: " + err.Error())
	}

	return nil
}

func (c *Client) UpdateDeploymentStatusForAllMatchingDeploys(owner string, repo string, ref string, environment string, status string) error {
	deployments, err := c.GetDeployments(owner, repo, ref, environment)

	for _, deployment := range deployments {
		status := "inactive"
		err = c.UpdateDeploymentStatus(owner, repo, status, *deployment.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) NotifyGitHubDeploy(owner string, repo string, branchName string, productionEnvironment bool, environment string, environmentURL string) error {
	err := c.DeleteDeployment(owner, repo, branchName, environment)
	if err != nil {
		return err
	}

	err = c.CreateDeployment(owner, repo, branchName, productionEnvironment, environment, environmentURL)
	if err != nil {
		return err
	}

	return nil
}
