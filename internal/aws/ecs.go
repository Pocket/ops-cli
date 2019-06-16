package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/dnaeon/go-vcr/recorder"
)

func DeployUpdate(clusterName *string, serviceName *string, imageName *string) {
	//client, clientContext := ecsClient()

	//TODO: get current service
	//TODO: get current service task definition
	//TODO: create a new task definition based on the current one and replace the image like here: https://github.com/fabfuel/ecs-deploy/blob/238a6715339f55b05cfe46338a8bc71d54aa4d93/ecs_deploy/ecs.py#L231
	//TODO: deploy new task definition

	//client.DescribeServicesRequest(&ecs.DescribeServicesInput{
	//	Cluster: clusterName,
	//	Services: []string{
	//		*serviceName,
	//	},
	//})
	//
	//client.DescribeTaskDefinitionRequest(&ecs.DescribeTaskDefinitionInput{
	//
	//})
	//
	//client.UpdateServiceRequest(&ecs.UpdateServiceInput{
	//	Cluster:
	//})

}

func ecsClient() (*ecs.Client, context.Context) {
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	cfg.HTTPClient.Transport = awsECSRecorder()
	return ecs.New(cfg), context.Background()
}

func awsECSRecorder() *recorder.Recorder {
	r, err := recorder.New("_fixtures/cloudformation")
	if err != nil {
		panic(err)
	}
	defer r.Stop() // Make sure recorder is stopped once done with it

	return r
}
