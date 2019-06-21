package feature_deploy

import (
	"github.com/Pocket/ops-cli/internal/util"
)

func stackNameFromBranchName(prefix string, branchName string) string {
	return prefix + util.DomainSafeString(branchName)
}
