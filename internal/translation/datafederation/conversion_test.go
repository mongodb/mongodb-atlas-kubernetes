package datafederation

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"
)

func TestRoundtrip_DataFederation(t *testing.T) {
	f := fuzz.New()

	for i := 0; i < 100; i++ {
		fuzzed := &DataFederation{}
		f.Fuzz(fuzzed)
		fuzzed, err := NewDataFederation(fuzzed.DataFederationSpec, fuzzed.ProjectID, fuzzed.Hostnames)
		require.NoError(t, err)

		toAtlasResult := toAtlas(fuzzed)
		fromAtlasResult, err := fromAtlas(toAtlasResult)
		require.NoError(t, err)

		equals := fuzzed.SpecEqualsTo(fromAtlasResult)
		if !equals {
			t.Log(cmp.Diff(fuzzed, fromAtlasResult))
		}
		require.True(t, equals)
	}
}
