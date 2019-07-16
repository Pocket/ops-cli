package feature_deploy

import (
	"github.com/Pocket/ops-cli/internal/settings"
	"testing"
)

//Daniel - This just makes sure we don't panic for now.
//All functions should be updated to return errors we can test in a separate pr.
func TestClient_CleanUpBranch(t *testing.T) {
	client, r := testFeatureDeploysClient("_fixtures/cleanup_branch")
	defer r.Stop()

	templateFilePath := "_fixtures/template_cloudformation.yml"
	branchName := "master"
	createdSettings := settings.NewSettingsParams("_fixtures/template_parameters.json", &templateFilePath, nil, &branchName , &branchName)
	client.CleanUpBranch(createdSettings, nil, nil)
}
