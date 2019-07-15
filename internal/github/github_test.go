package github

import (
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"gotest.tools/assert"
	"testing"
)

func testGithubClient(cassetteName string) (*Client, *recorder.Recorder) {
	r, _ := recorder.New(cassetteName)
	// Add a filter which removes Authorization headers from all requests:
	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		return nil
	})

	client := New("", "Pocket", "Web", r)
	return client, r
}

func TestClient_CreateDeployment(t *testing.T) {
	client, r := testGithubClient("_fixtures/create_deployment_branch")
	defer r.Stop()
	err := client.CreateDeployment( "subscriptions-ending-soon", false, "web-feature", "https://feature.test.com")
	assert.NilError(t, err)
}

func TestClient_DeleteDeployment(t *testing.T) {
	client, r := testGithubClient("_fixtures/delete_deployment")
	defer r.Stop()
	err := client.DeleteDeployment("subscriptions-ending-soon", "web-feature")
	assert.NilError(t, err)
}

func TestClient_GetDeployments(t *testing.T) {
	client, r := testGithubClient("_fixtures/get_deployments")
	defer r.Stop()
	deployments, err := client.GetDeployments( "subscriptions-ending-soon", "web-feature")
	assert.NilError(t, err)
	assert.Assert(t, len(deployments) == 1)
	for _, deployment := range deployments {
		assert.Assert(t, *deployment.Environment == "web-feature")
		assert.Assert(t, *deployment.Ref == "subscriptions-ending-soon")
	}
}

func TestClient_UpdateDeploymentStatusForAllMatchingDeploys(t *testing.T) {
	client, r := testGithubClient("_fixtures/update_deployment_status")
	defer r.Stop()
	err := client.UpdateDeploymentStatusForAllMatchingDeploys( "subscriptions-ending-soon", "web-feature", "pending")
	assert.NilError(t, err)
}

func TestClient_NotifyGitHubDeploy_Initial(t *testing.T) {
	client, r := testGithubClient("_fixtures/notify_github_deploy_initial")
	defer r.Stop()
	err := client.NotifyGitHubDeploy( "subscriptions-ending-soon", false, "feature.com", "https://feature.com/subscriptions-ending-soon")
	assert.NilError(t, err)
}


func TestClient_NotifyGitHubDeploy_Update(t *testing.T) {
	client, r := testGithubClient("_fixtures/notify_github_deploy_update")
	defer r.Stop()
	err := client.NotifyGitHubDeploy("Pocket", "Web", "subscriptions-ending-soon", false, "feature.com", "https://feature.com/subscriptions-ending-soon")
	assert.NilError(t, err)
}