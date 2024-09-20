package datafederation

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	fuzz "github.com/google/gofuzz"

	akocmp "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/cmp"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
)

func TestRoundtrip_DataFederationPE(t *testing.T) {
	f := fuzz.New().NilChance(0.0).NumElements(1, 10)
	f.Funcs(
		NonEmptyString,
		EnsureEmptySliceOf[string](f),
	)

	for i := 0; i < 100; i++ {
		var atlas admin.PrivateNetworkEndpointIdEntry
		f.Fuzz(&atlas)

		fromAtlasResult := endpointFromAtlas(atlas, "")
		toAtlasResult := endpointToAtlas(fromAtlasResult)

		require.NoError(t, akocmp.Normalize(&atlas))
		require.NoError(t, akocmp.Normalize(toAtlasResult))

		equals := reflect.DeepEqual(&atlas, toAtlasResult)
		if !equals {
			t.Log(cmp.Diff(&atlas, toAtlasResult))
		}
		require.True(t, equals)
	}
}
