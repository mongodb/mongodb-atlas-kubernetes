package audit

import (
	"context"
	_ "embed"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/audit"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/contract"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
)

func TestDefaultAuditingGet(t *testing.T) {
	ctx := context.Background()
	contract.RunGoContractTest(ctx, t, "get default auditing", func(ch contract.ContractHelper) {
		projectName := utils.RandomName("default-auditing-project")
		require.NoError(t, ch.AddResources(ctx, 5*time.Minute, contract.DefaultAtlasProject(projectName)))
		testProjectID, err := ch.ProjectID(ctx, projectName)
		require.NoError(t, err)
		as := audit.NewAuditLog(ch.AtlasClient().AuditingApi)

		result, err := as.Get(ctx, testProjectID)
		require.NoError(t, err)
		assert.Equal(t, audit.NewAuditConfig(nil), result)
	})
}

func TestSyncs(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		title    string
		auditing *audit.AuditConfig
	}{
		{
			title: "Just enabled",
			auditing: audit.NewAuditConfig(
				&akov2.Auditing{
					Enabled: true,
				},
			),
		},
		{
			title: "Auth success logs as well",
			auditing: audit.NewAuditConfig(
				&akov2.Auditing{
					Enabled:                   true,
					AuditAuthorizationSuccess: true,
				},
			),
		},
		{
			title: "With a filter",
			auditing: audit.NewAuditConfig(
				&akov2.Auditing{
					Enabled:     true,
					AuditFilter: `{"atype":"authenticate"}`,
				},
			),
		},
		{
			title: "With a filter and success logs",
			auditing: audit.NewAuditConfig(
				&akov2.Auditing{
					Enabled:                   true,
					AuditAuthorizationSuccess: true,
					AuditFilter:               `{"atype":"authenticate"}`,
				},
			),
		},
		{
			title: "All set but disabled",
			auditing: audit.NewAuditConfig(
				&akov2.Auditing{
					AuditAuthorizationSuccess: true,
					AuditFilter:               `{"atype":"authenticate"}`,
				},
			),
		},
		{
			title: "Default (disabled) case",
			auditing: audit.NewAuditConfig(
				&akov2.Auditing{},
			),
		},
	}
	contract.RunGoContractTest(ctx, t, "test syncs", func(ch contract.ContractHelper) {
		projectName := utils.RandomName("audit-syncs-project")
		require.NoError(t, ch.AddResources(ctx, 5 * time.Minute, contract.DefaultAtlasProject(projectName)))
		testProjectID, err := ch.ProjectID(ctx, projectName)
		require.NoError(t, err)
		as := audit.NewAuditLog(ch.AtlasClient().AuditingApi)

		for _, tc := range testCases {
			t.Run(tc.title, func(t *testing.T) {
				err := as.Update(ctx, testProjectID, tc.auditing)
				require.NoError(t, err)

				result, err := as.Get(ctx, testProjectID)
				require.NoError(t, err)
				assert.Equal(t, tc.auditing, result)
			})
		}
	})
}
