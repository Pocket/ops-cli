package cloudwatchlogs

import (
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"gotest.tools/assert"
	"testing"
)

func TestExportLogGroupAndWait(t *testing.T) {
	client, r := testCloudWatchLogsClient("_fixtures/export_log_group_and_wait")
	defer r.Stop()
	err := client.ExportLogGroupAndWait("/ecs/webFeature/WebFeatureDeploy-testfeature-deploys-go", "pocket-web-feature-archived-logs")
	assert.Assert(t, err == nil)
}

func TestExportLogGroup(t *testing.T) {
	client, r := testCloudWatchLogsClient("_fixtures/export_log_group")
	defer r.Stop()
	taskId := client.ExportLogGroup("/ecs/webFeature/WebFeatureDeploy-testfeature-deploys-go", "pocket-web-feature-archived-logs")
	assert.Assert(t, *taskId == "793e51c8-8313-4f7b-988f-a5273ed59090")
}

func TestIsExportTaskPending(t *testing.T) {
	client, r := testCloudWatchLogsClient("_fixtures/export_task_pending")
	defer r.Stop()
	running, err := client.IsExportTaskRunning("790e51c8-8313-1234-988f-a5273ed59090")
	assert.Assert(t, running)
	assert.Assert(t, err == nil)
}

func TestIsExportTaskRunning(t *testing.T) {
	client, r := testCloudWatchLogsClient("_fixtures/export_task_running")
	defer r.Stop()
	running, err := client.IsExportTaskRunning("790e51c8-8313-1234-988f-a5273ed59090")
	assert.Assert(t, running)
	assert.Assert(t, err == nil)
}

func TestIsExportTaskPendingCancel(t *testing.T) {
	client, r := testCloudWatchLogsClient("_fixtures/export_task_pending_cancel")
	defer r.Stop()
	running, err := client.IsExportTaskRunning("790e51c8-8313-1234-988f-a5273ed59090")
	assert.Assert(t, !running)
	assert.Assert(t, err != nil)
}

func TestIsExportTaskCompleted(t *testing.T) {
	client, r := testCloudWatchLogsClient("_fixtures/export_task_completed")
	defer r.Stop()
	running, err := client.IsExportTaskRunning("790e51c8-8313-1234-988f-a5273ed59090")
	assert.Assert(t, !running)
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
