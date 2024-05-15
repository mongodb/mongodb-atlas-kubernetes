package auditing

import (
	"context"
	_ "embed"
	"log"
	"testing"
	"time"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translayer/auditing"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1alpha1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/launcher"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/contract"
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
	l.Launch(
		testYml,
		launcher.WaitReady("atlasprojects/my-project", 30*time.Second))
	if !control.Enabled("SKIP_CLEANUP") { // allow to reuse Atlas resources for local tests
		defer l.Cleanup()
	}
	m.Run()
}

func TestDefaultAuditingGet(t *testing.T) {
	testProjectID := mustReadProjectID()
	ctx := context.Background()
	as := auditing.NewFromAuditingAPI(contract.MustVersionedClient(t, ctx).AuditingApi)

	result, err := as.Get(ctx, testProjectID)

	require.NoError(t, err)
	result.ConfigurationType = "" // Do not expect the returned  cfg type to match
	if result.AuditFilter != nil && string(result.AuditFilter.Raw) == "{}" {
		// Support re-runs, as we cannot get the filter back to unset
		result.AuditFilter = nil
	}
	assert.Equal(t, defaultAtlasAuditing(), result)
}

func defaultAtlasAuditing() *v1alpha1.AtlasAuditingSpec {
	return &v1alpha1.AtlasAuditingSpec{
		Enabled:                   false,
		AuditAuthorizationSuccess: false,
		AuditFilter:               nil,
	}
}

func TestSyncs(t *testing.T) {
	testCases := []struct {
		title    string
		auditing *v1alpha1.AtlasAuditingSpec
	}{
		{
			title: "Just enabled",
			auditing: &v1alpha1.AtlasAuditingSpec{
				Enabled:                   true,
				AuditAuthorizationSuccess: false,
				AuditFilter:               jsonType("{}"), // must sent empty JSON to overwrite previous state
			},
		},
		{
			title: "Auth success logs as well",
			auditing: &v1alpha1.AtlasAuditingSpec{
				Enabled:                   true,
				AuditAuthorizationSuccess: true,
				AuditFilter:               jsonType("{}"),
			},
		},
		{
			title: "With a filter",
			auditing: &v1alpha1.AtlasAuditingSpec{
				Enabled:                   true,
				AuditAuthorizationSuccess: false,
				AuditFilter:               jsonType(`{"atype":"authenticate"}`),
			},
		},
		{
			title: "With a filter and success logs",
			auditing: &v1alpha1.AtlasAuditingSpec{
				Enabled:                   true,
				AuditAuthorizationSuccess: true,
				AuditFilter:               jsonType(`{"atype":"authenticate"}`),
			},
		},
		{
			title: "All set but disabled",
			auditing: &v1alpha1.AtlasAuditingSpec{
				Enabled:                   false,
				AuditAuthorizationSuccess: true,
				AuditFilter:               jsonType(`{"atype":"authenticate"}`),
			},
		},
		{
			title: "Default (disabled) case",
			auditing: &v1alpha1.AtlasAuditingSpec{
				Enabled:                   false,
				AuditAuthorizationSuccess: false,
				AuditFilter:               jsonType("{}"),
			},
		},
	}
	testProjectID := mustReadProjectID()
	ctx := context.Background()
	as := auditing.NewFromAuditingAPI(contract.MustVersionedClient(t, ctx).AuditingApi)

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			err := as.Set(ctx, testProjectID, tc.auditing)
			require.NoError(t, err)

			result, err := as.Get(ctx, testProjectID)
			require.NoError(t, err)
			result.ConfigurationType = "" // Do not expect the returned  cfg type to match
			assert.Equal(t, tc.auditing, result)
		})
	}
}

func mustReadProjectID() string {
	l := launcher.NewFromEnv(testVersion)
	output, err := l.Kubectl("get", "atlasprojects/my-project", "-o=jsonpath={.status.id}")
	if err != nil {
		log.Fatalf("Failed to get test project id: %v", err)
	}
	return output
}

func jsonType(js string) *apiextensionsv1.JSON {
	return &apiextensionsv1.JSON{Raw: []byte(js)}
}
