package util

import (
	"gotest.tools/assert"
	"testing"
)

func TestDomainSafeString(t *testing.T) {
	assert.Equal(t, DomainSafeString("feature_testing"), "feature-testing")
	assert.Equal(t, DomainSafeString("feature testing"), "feature-testing")
	assert.Equal(t, DomainSafeString("fe!@atu-re$tes%!@ting"), "featu-retesting")
	assert.Equal(t, DomainSafeString("feature/tEsting  "), "featuretesting")
	assert.Equal(t, DomainSafeString("   feature/tesTing/taKe2 "), "featuretestingtake2")
}
