package audit

import (
	"context"
	_ "embed"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/audit"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/contract"
)

func TestDefaultAuditingGet(t *testing.T) {
	contract.RunContractTest(t, "get default auditing", func(ct *contract.ContractTest) {
		projectName := "default-auditing-project"
		ct.AddResources(time.Minute, contract.DefaultAtlasProject(projectName))
		testProjectID := mustReadProjectID(t, ct.Ctx, ct.K8sClient, ct.Namespace(), projectName)
		as := audit.NewAuditLog(ct.AtlasClient.AuditingApi)

		result, err := as.Get(ct.Ctx, testProjectID)
		require.NoError(t, err)
		assert.Equal(t, audit.NewAuditConfig(nil), result)
	})
}

func TestSyncs(t *testing.T) {
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
	contract.RunContractTest(t, "test syncs", func(ct *contract.ContractTest) {
		projectName := "audit-syncs-project"
		ct.AddResources(time.Minute, contract.DefaultAtlasProject(projectName))
		testProjectID := mustReadProjectID(t, ct.Ctx, ct.K8sClient, ct.Namespace(), projectName)
		ctx := context.Background()
		as := audit.NewAuditLog(ct.AtlasClient.AuditingApi)

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

func mustReadProjectID(t *testing.T, ctx context.Context, k8sClient client.Client, ns, projectName string) string {
	t.Helper()
	project := akov2.AtlasProject{}
	key := types.NamespacedName{
		Namespace: ns,
		Name:      projectName,
	}
	require.NoError(t, k8sClient.Get(ctx, key, &project))
	return project.Status.ID
}
