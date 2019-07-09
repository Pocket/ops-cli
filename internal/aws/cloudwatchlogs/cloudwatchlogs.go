package cloudwatchlogs

import (
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"time"
)

type Client struct {
	client        *cloudwatchlogs.Client
	clientContext context.Context
}

func New() *Client {
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	return &Client{
		client:        cloudwatchlogs.New(cfg),
		clientContext: context.Background(),
	}
}

func (c *Client) ExportLogGroupAndWait(logGroupName string, s3BucketName string) error {
	exportTaskId := c.ExportLogGroup(logGroupName, s3BucketName)

	isExporting := true
	var errorExporting error

	for isExporting {
		isExporting, errorExporting = c.IsExportTaskRunning(*exportTaskId)
		if !isExporting && errorExporting != nil {
			return errorExporting
		}
		time.Sleep(5 * time.Second)
	}

	return nil
}

//NOTE: Each account can only have one active (RUNNING or PENDING) export task at a time.
func (c *Client) ExportLogGroup(logGroupName string, s3BucketName string) *string {

	from := int64(300)
	to := time.Now().AddDate(0, 0, 1).Unix()*1000

	response, err := c.client.CreateExportTaskRequest(&cloudwatchlogs.CreateExportTaskInput{
		Destination:       &s3BucketName,
		DestinationPrefix: &logGroupName,
		LogGroupName:      &logGroupName,
		From:              &from,
		To:                &to,
	}).Send(c.clientContext)

	if err != nil {
		panic("error creating the log group export, " + err.Error())
	}

	return response.TaskId
}

func (c *Client) IsExportTaskRunning(exportTaskId string) (bool, error) {
	response, err := c.client.DescribeExportTasksRequest(&cloudwatchlogs.DescribeExportTasksInput{
		TaskId: &exportTaskId,
	}).Send(c.clientContext)

	if err != nil {
		panic("There was an error getting the export task: " + err.Error())
	}

	if len(response.ExportTasks) == 0 {
		panic("The export task does not exist")
	}

	exportTask := response.ExportTasks[0]

	if exportTask.Status.Code == cloudwatchlogs.ExportTaskStatusCodeCancelled ||
		exportTask.Status.Code == cloudwatchlogs.ExportTaskStatusCodeFailed ||
		exportTask.Status.Code == cloudwatchlogs.ExportTaskStatusCodePendingCancel {
		return false, errors.New("The export task failed")
	}

	if exportTask.Status.Code == cloudwatchlogs.ExportTaskStatusCodeCompleted {
		return false, nil
	}

	return true, nil
}
