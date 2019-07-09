package cloudformation

import (
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"gotest.tools/assert"
	"testing"
)

func testCloudFormationClient(cassetteName string) (*Client, *recorder.Recorder) {
	client := New()
	r, _ := recorder.New(cassetteName)
	// Add a filter which removes Authorization headers from all requests:
	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		delete(i.Request.Headers, "X-Amz-Security-Token")
		delete(i.Response.Headers, "X-Amzn-Requestid")
		return nil
	})

	client.client.Config.HTTPClient.Transport = r
	return client, r
}

func TestActiveCloudFormationStackBranchesWithPrefix(t *testing.T) {
	client, r := testCloudFormationClient("_fixtures/active_stacks")
	defer r.Stop()

	strings := []string{
		"expire-ecr",
		"admin/syndication-article-v2-1-1",
		"sqs-cleanup",
		"queue-deprecated",
	}
	stacks := client.ActiveCloudFormationStackBranchesWithPrefix("WebFeatureDeploy-")
	assert.DeepEqual(t, stacks, strings)
}

func TestStackExists(t *testing.T) {
	client, r := testCloudFormationClient("_fixtures/active_stacks")
	defer r.Stop()

	assert.Assert(t, client.StackExists("expire-ecr"))
	assert.Assert(t, !client.StackExists("do-i-exist"))
}
