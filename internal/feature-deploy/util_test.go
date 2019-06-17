package feature_deploy

import (
	"gotest.tools/assert"
	"testing"
)

func TestStackNameFromBranchName(t *testing.T) {
	assert.Equal(t, stackNameFromBranchName("WebFeatureDeploy-", "fe!@atu-re$tes%!@ting"), "WebFeatureDeploy-featu-retesting")
	assert.Equal(t, stackNameFromBranchName("WebFeatureDeploy-", "feature testing"), "WebFeatureDeploy-feature-testing")
	assert.Equal(t, stackNameFromBranchName("WebFeatureDeploy-", "fe!@atu-re$tes%!@ting"), "WebFeatureDeploy-featu-retesting")
	assert.Equal(t, stackNameFromBranchName("WebFeatureDeploy-", "feature/tEsting  "), "WebFeatureDeploy-featuretesting")
	assert.Equal(t, stackNameFromBranchName("WebFeatureDeploy-", "   feature/tesTing/taKe2 "), "WebFeatureDeploy-featuretestingtake2")
}
