package github

import (
	"context"
	"github.com/google/go-github/v26/github"
	"golang.org/x/oauth2"
)

type Client struct {
	client        *github.Client
	clientContext context.Context
}

func New(accessToken string) *Client {
	clientContext := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: accessToken,
		},
	)
	tc := oauth2.NewClient(clientContext, ts)

	client := github.NewClient(tc)

	return &Client{
		client:        client,
		clientContext: clientContext,
	}
}

func (c *Client) CreateDeployment(owner string, repo string, gitSHA string, productionEnvironment bool, environment string, environmentURL string) {
	autoMerge := false
	transientEnvironment := !productionEnvironment

	deployment, _, err := c.client.Repositories.CreateDeployment(c.clientContext, owner, repo, &github.DeploymentRequest{
		Ref:                   &gitSHA,
		AutoMerge:             &autoMerge,
		Environment:           &environment,
		TransientEnvironment:  &transientEnvironment,
		ProductionEnvironment: &productionEnvironment,
	})

	if err != nil {
		panic("Error creating github deployment: " + err.Error())
	}

	status := "success"
	c.client.Repositories.CreateDeploymentStatus(c.clientContext, owner, repo, *deployment.ID, &github.DeploymentStatusRequest{
		State:          &status,
		Environment:    &environment,
		EnvironmentURL: &environmentURL,
	})

}
