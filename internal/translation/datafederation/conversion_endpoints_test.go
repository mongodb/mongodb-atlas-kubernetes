package datafederation

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"
)

func TestRoundtrip_DataFederationPE(t *testing.T) {
	f := fuzz.New()

	for i := 0; i < 100; i++ {
		fuzzed := PrivateEndpoint{}
		f.Fuzz(&fuzzed)
		// ignore non-Atlas fields
		fuzzed.ProjectID = ""

		toAtlasResult := endpointToAtlas(&fuzzed)
		fromAtlasResult := endpointFromAtlas("", toAtlasResult)

		equals := reflect.DeepEqual(fuzzed, fromAtlasResult)
		if !equals {
			t.Log(cmp.Diff(fuzzed, fromAtlasResult))
		}
		require.True(t, equals)
	}
}
