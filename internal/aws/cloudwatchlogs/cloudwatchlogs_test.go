package cloudwatchlogs

import (
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"gotest.tools/assert"
	"testing"
)

func TestExportLogGroupAndWait(t *testing.T) {
	client, r := testCloudWatchLogsClient("export_log_group_and_wait")
	defer r.Stop()
	err := client.ExportLogGroupAndWait("/ecs/webFeature/WebFeatureDeploy-testfeature-deploys-go", "pocket-web-feature-archived-logs")
	assert.Assert(t, err == nil)
}

func testCloudWatchLogsClient(cassetteName string) (*Client, *recorder.Recorder) {
	client := New()
	r, _ := recorder.New(cassetteName)
	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		delete(i.Request.Headers, "X-Amz-Security-Token")
		delete(i.Response.Headers, "X-Amzn-Requestid")
		return nil
	})

	client.client.Config.HTTPClient.Transport = r
	return client, r
}
