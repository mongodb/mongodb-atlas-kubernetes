package datafederation

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
)

func TestRoundtrip_DataFederationPE(t *testing.T) {
	f := fuzz.New()

	for i := 0; i < 100; i++ {
		var atlas admin.PrivateNetworkEndpointIdEntry
		f.Fuzz(&atlas)

		fromAtlasResult := endpointFromAtlas(atlas, "")
		toAtlasResult := endpointToAtlas(fromAtlasResult)
		roundtripAtlasResult := endpointFromAtlas(*toAtlasResult, "")

		equals := reflect.DeepEqual(fromAtlasResult, roundtripAtlasResult)
		if !equals {
			t.Log(cmp.Diff(fromAtlasResult, roundtripAtlasResult))
		}
		require.True(t, equals)
	}
}
