package feature_deploy

import (
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"testing"
)

//Daniel - This just makes sure we don't panic for now.
//All functions should be updated to return errors we can test in a separate pr.
func TestClient_DeployBranch(t *testing.T) {
	client, r := testFeatureDeploysClient("_fixtures/create_new_stack")
	defer r.Stop()
	client.DeployBranch("_fixtures/template_parameters.json", "_fixtures/template_cloudformation.yml", "master", "12345678912345678912345678912345678912", "12345678912345678912345678912345678912")
}

func testFeatureDeploysClient(cassetteName string) (*Client, *recorder.Recorder) {
	client := New()
	r, _ := recorder.New(cassetteName)
	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		delete(i.Request.Headers, "X-Amz-Security-Token")
		delete(i.Response.Headers, "X-Amzn-Requestid")
		return nil
	})

	client.setTransport(r)
	return client, r
}
