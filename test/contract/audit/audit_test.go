package audit

import (
	"context"
	_ "embed"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/audit"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/contract"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/launcher"
)

//go:embed test.yml
var testYml string

const (
	testVersion = "2.1.0"
)

func TestMain(m *testing.M) {
	if !control.Enabled("AKO_CONTRACT_TEST") {
		log.Printf("Skipping contract test as AKO_CONTRACT_TEST is unset")
		return
	}

	l := launcher.NewFromEnv(testVersion)
	if err := l.Launch(
		testYml,
		launcher.WaitReady("atlasprojects/my-project", time.Minute)); err != nil {
		log.Fatalf("Failed to launch test bed: %v", err)
	}

	if !control.Enabled("SKIP_CLEANUP") { // allow to reuse Atlas resources for local tests
		defer l.Cleanup()
	}
	os.Exit(m.Run())
}

func TestDefaultAuditingGet(t *testing.T) {
	testProjectID := mustReadProjectID("atlasprojects/my-project2")
	ctx := context.Background()
	as := audit.NewAuditLog(contract.MustVersionedClient(t, ctx).AuditingApi)

	result, err := as.Get(ctx, testProjectID)
	require.NoError(t, err)
	assert.Equal(t, audit.NewAuditConfig(nil), result)
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
	testProjectID := mustReadProjectID("atlasprojects/my-project")
	ctx := context.Background()
	as := audit.NewAuditLog(contract.MustVersionedClient(t, ctx).AuditingApi)

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			err := as.Update(ctx, testProjectID, tc.auditing)
			require.NoError(t, err)

			result, err := as.Get(ctx, testProjectID)
			require.NoError(t, err)
			assert.Equal(t, tc.auditing, result)
		})
	}
}

func mustReadProjectID(namespacedName string) string {
	l := launcher.NewFromEnv(testVersion)
	output, err := l.Kubectl("get", namespacedName, "-o=jsonpath={.status.id}")
	if err != nil {
		log.Fatalf("Failed to get test project id: %v", err)
	}
	return output
}
