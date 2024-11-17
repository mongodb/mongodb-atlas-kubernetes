package integrations

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	fuzz "github.com/google/gofuzz"

	"github.com/stretchr/testify/require"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
)

func TestRoundTrip_Integrations(t *testing.T) {
	f := fuzz.New()

	for range 100 {
		fuzzed := &Integration{}
		f.Fuzz(fuzzed)
		fuzzed, err := NewIntegration(&fuzzed.Integration)
		require.NoError(t, err)

		// Don't fuzz secrets that we haven't created in the fake client
		fuzzed.LicenseKeyRef = common.ResourceRefNamespaced{}
		fuzzed.WriteTokenRef = common.ResourceRefNamespaced{}
		fuzzed.ReadTokenRef = common.ResourceRefNamespaced{}
		fuzzed.APIKeyRef = common.ResourceRefNamespaced{}
		fuzzed.ServiceKeyRef = common.ResourceRefNamespaced{}
		fuzzed.APITokenRef = common.ResourceRefNamespaced{}
		fuzzed.RoutingKeyRef = common.ResourceRefNamespaced{}
		fuzzed.SecretRef = common.ResourceRefNamespaced{}
		fuzzed.PasswordRef = common.ResourceRefNamespaced{}

		// Don't expect the 'dud' fields to be converted
		fuzzed.FlowName = ""
		fuzzed.OrgName = ""
		fuzzed.Name = ""
		fuzzed.Scheme = ""

		toAtlasResult := toAtlas(*fuzzed, nil)

		fromAtlasResult, err := fromAtlas(toAtlasResult)
		require.NoError(t, err)

		equals := cmp.Diff(fuzzed, fromAtlasResult) == ""
		require.True(t, equals)
	}
}
