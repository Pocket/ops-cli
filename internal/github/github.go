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
	params        *Params
}

type Params struct {
	Owner string
	Repo  string
	AccessToken string
}

func New(params *Params, transport http.RoundTripper) *Client {
	clientContext := context.TODO()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: params.AccessToken,
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
		params: params,
	}
}

func (c *Client) SetTransport(transport http.RoundTripper)  {
	tc := &http.Client{
		Transport: transport,
	}
	c.clientContext = context.WithValue(c.clientContext, oauth2.HTTPClient, tc)
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: c.params.AccessToken,
		},
	)
	tc = oauth2.NewClient(c.clientContext, ts)
	c.client = github.NewClient(tc)
}

func (c *Client) CreateDeployment(branchName string, productionEnvironment bool, environment string, environmentURL string) error {
	autoMerge := false
	transientEnvironment := !productionEnvironment
	requiredContexts := []string{}

	deployment, _, err := c.client.Repositories.CreateDeployment(c.clientContext, c.params.Owner, c.params.Repo, &github.DeploymentRequest{
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
	_, _, err = c.client.Repositories.CreateDeploymentStatus(c.clientContext, c.params.Owner, c.params.Repo, *deployment.ID, &github.DeploymentStatusRequest{
		State:          &status,
		Environment:    &environment,
		EnvironmentURL: &environmentURL,
	})

	if err != nil {
		return errors.New("Error creating active github deployment status: " + err.Error())
	}

	return nil
}

func (c *Client) DeleteDeployment(ref string, environment string) error {
	return c.UpdateDeploymentStatusForAllMatchingDeploys(ref, environment, "inactive")
}

func (c *Client) GetDeployments(ref string, environment string) ([]*github.Deployment, error) {
	deployments, _, err := c.client.Repositories.ListDeployments(c.clientContext, c.params.Owner, c.params.Repo, &github.DeploymentsListOptions{
		Ref:         ref,
		Environment: environment,
	})

	if err != nil {
		return nil, errors.New("Error listing deployments: " + err.Error())
	}

	return deployments, nil
}

func (c *Client) UpdateDeploymentStatus(status string, deploymentId int64) error {
	_, _, err := c.client.Repositories.CreateDeploymentStatus(c.clientContext, c.params.Owner, c.params.Repo, deploymentId, &github.DeploymentStatusRequest{
		State: &status,
	})
	if err != nil {
		return errors.New("Error creating github deployment status: " + err.Error())
	}

	return nil
}

func (c *Client) UpdateDeploymentStatusForAllMatchingDeploys(ref string, environment string, status string) error {
	deployments, err := c.GetDeployments(ref, environment)

	for _, deployment := range deployments {
		status := "inactive"
		err = c.UpdateDeploymentStatus(status, *deployment.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) NotifyGitHubDeploy(branchName string, productionEnvironment bool, environment string, environmentURL string) error {
	err := c.DeleteDeployment(branchName, environment)
	if err != nil {
		return err
	}

	err = c.CreateDeployment(branchName, productionEnvironment, environment, environmentURL)
	if err != nil {
		return err
	}

	return nil
}
