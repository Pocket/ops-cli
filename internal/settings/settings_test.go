package settings

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"gotest.tools/assert"
	"testing"
)

func TestCloudFormationSettings(t *testing.T) {
	pocketSettings := NewSettings("_data/parameters.json")

	settings := &Settings{
		Parameters: []cloudformation.Parameter{
			{
				ParameterKey:   getString("IamStack"),
				ParameterValue: getString("WebFeatureIAM"),
			},
			{
				ParameterKey:   getString("ALBStack"),
				ParameterValue: getString("WebFeatureShared"),
			},
			{
				ParameterKey:   getString("Env"),
				ParameterValue: getString("Feature"),
			},
			{
				ParameterKey:   getString("HostedZoneId"),
				ParameterValue: getString("1234"),
			},
			{
				ParameterKey:   getString("DomainBase"),
				ParameterValue: getString("web.test.com"),
			},
			{
				ParameterKey:   getString("BranchName"),
				ParameterValue: getString("feature-deploys"),
			},
			{
				ParameterKey:   getString("FormattedBranchName"),
				ParameterValue: getString("feature-deploys"),
			},
			{
				ParameterKey:   getString("GitSHA"),
				ParameterValue: getString("latest"),
			},
		},
		Tags: []cloudformation.Tag{
			{
				Key:   getString("GitSHA"),
				Value: getString("WebFeature"),
			},
			{
				Key:   getString("Env"),
				Value: getString("Feature"),
			},
			{
				Key:   getString("BranchName"),
				Value: getString("feature-deploys"),
			},
		},
		OnFailure: cloudformation.OnFailureDelete,
		Capabilities: []cloudformation.Capability{
			cloudformation.CapabilityCapabilityIam,
			cloudformation.CapabilityCapabilityNamedIam,
			cloudformation.CapabilityCapabilityAutoExpand,
		},
		TemplateBody:     getString("testbody"),
		TimeoutInMinutes: getInt64(5),
		StackName:        getString("stackname"),
		ECSCluster:       getString("WebFeatureShared"),
	}

	for k, value := range pocketSettings.Tags {
		assert.DeepEqual(t, value.Key, settings.Tags[k].Key)
		assert.DeepEqual(t, value.Value, settings.Tags[k].Value)
	}

	for k, value := range pocketSettings.Parameters {
		assert.DeepEqual(t, value.ParameterKey, settings.Parameters[k].ParameterKey)
		assert.DeepEqual(t, value.ParameterValue, settings.Parameters[k].ParameterValue)
	}

	assert.DeepEqual(t, pocketSettings.OnFailure, settings.OnFailure)
	assert.DeepEqual(t, pocketSettings.Capabilities, settings.Capabilities)
	assert.DeepEqual(t, pocketSettings.TemplateBody, settings.TemplateBody)
	assert.DeepEqual(t, pocketSettings.TimeoutInMinutes, settings.TimeoutInMinutes)
	assert.DeepEqual(t, pocketSettings.StackPrefix, settings.StackPrefix)

	assert.DeepEqual(t, pocketSettings.SlackCleanUpSettings.Icon, ":cleanup:")
	assert.DeepEqual(t, pocketSettings.SlackCleanUpSettings.Username, "Damage Control")
	assert.DeepEqual(t, pocketSettings.SlackCleanUpSettings.Channel, "#log-feature-deploys")

	assert.DeepEqual(t, pocketSettings.SlackDeploySettings.Icon, ":ship:")
	assert.DeepEqual(t, pocketSettings.SlackDeploySettings.Username, "Buster")
	assert.DeepEqual(t, pocketSettings.SlackDeploySettings.Channel, "#log-feature-deploys")

	assert.Assert(t, pocketSettings.StackName == nil)
	assert.Assert(t, pocketSettings.FormattedBranchName == nil)
	assert.Assert(t, pocketSettings.GitSHA == nil)
	assert.Assert(t, pocketSettings.BranchName == nil)
}

func TestCloudFormationSettingsOmitting(t *testing.T) {
	pocketSettings := NewSettings("_data/parameters_omitting.json")

	settings := &Settings{
		TemplateBody:     getString("file://test.yaml"),
		TimeoutInMinutes: getInt64(5),
		OnFailure:        cloudformation.OnFailureDelete,
	}
	assert.DeepEqual(t, pocketSettings.Parameters, []cloudformation.Parameter(nil))
	assert.DeepEqual(t, pocketSettings.Tags, []cloudformation.Tag(nil))
	assert.DeepEqual(t, pocketSettings.StackName, (*string)(nil))
	assert.DeepEqual(t, pocketSettings.ECSCluster, (*string)(nil))
	assert.DeepEqual(t, pocketSettings.Capabilities, []cloudformation.Capability(nil))
	assert.DeepEqual(t, pocketSettings.OnFailure, settings.OnFailure)
	assert.DeepEqual(t, pocketSettings.TemplateBody, settings.TemplateBody)
	assert.DeepEqual(t, pocketSettings.TimeoutInMinutes, settings.TimeoutInMinutes)
}

/**
 * Helper to get a referenced string
 */
func getString(theString string) *string {
	return &theString
}

/**
 * Helper to get a referenced int
 */
func getInt64(theNumber int64) *int64 {
	return &theNumber
}
