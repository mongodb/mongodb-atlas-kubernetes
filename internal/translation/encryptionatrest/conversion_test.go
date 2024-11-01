package encryptionatrest

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
)

func TestRoundtrip_EncryptionAtRest(t *testing.T) {
	f := fuzz.New()

	for range 100 {
		fuzzed := &EncryptionAtRest{}
		f.Fuzz(fuzzed)

		//ignore secret fields
		fuzzed.AWS.AccessKeyID = ""
		fuzzed.AWS.SecretAccessKey = ""
		fuzzed.AWS.CustomerMasterKeyID = ""
		fuzzed.AWS.RoleID = ""
		fuzzed.AWS.CloudProviderIntegrationRole = ""
		fuzzed.AWS.SecretRef = common.ResourceRefNamespaced{}

		fuzzed.Azure.SubscriptionID = ""
		fuzzed.Azure.KeyVaultName = ""
		fuzzed.Azure.KeyIdentifier = ""
		fuzzed.Azure.Secret = ""
		fuzzed.Azure.SecretRef = common.ResourceRefNamespaced{}

		fuzzed.GCP.ServiceAccountKey = ""
		fuzzed.GCP.KeyVersionResourceID = ""
		fuzzed.GCP.SecretRef = common.ResourceRefNamespaced{}

		//ignore read-only 'Valid' field
		fuzzed.AWS.Valid = nil
		fuzzed.Azure.Valid = nil
		fuzzed.GCP.Valid = nil

		toAtlasResult := toAtlas(fuzzed)
		fromAtlasResult := fromAtlas(toAtlasResult)

		equals := EqualSpecs(fuzzed, fromAtlasResult)
		if !equals {
			t.Log(cmp.Diff(fuzzed, fromAtlasResult))
		}
		require.True(t, equals)
	}
}
