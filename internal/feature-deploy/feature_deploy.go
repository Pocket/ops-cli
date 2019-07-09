package feature_deploy

import (
	"github.com/Pocket/ops-cli/internal/aws/cloudformation"
	"github.com/Pocket/ops-cli/internal/aws/cloudwatchlogs"
	"github.com/Pocket/ops-cli/internal/aws/ecs"
)

type Client struct {
	cloudFormationClient *cloudformation.Client
	cloudWatchLogsClient *cloudwatchlogs.Client
	ecsClient            *ecs.Client
}

func New() *Client {
	return &Client{
		cloudFormationClient: cloudformation.New(),
		cloudWatchLogsClient: cloudwatchlogs.New(),
		ecsClient:            ecs.New(),
	}
}
