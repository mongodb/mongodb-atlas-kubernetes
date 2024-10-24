package maintenancewindow

import (
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"
)

func TestRoundtrip_MaintenanceWindow(t *testing.T) {
	// Atlas has no 'Defer' field so we don't receive this back via fromAtlas
	f := fuzz.New().SkipFieldsWithPattern(regexp.MustCompile("Defer"))

	for i := 0; i < 100; i++ {
		fuzzed := &MaintenanceWindow{}
		f.Fuzz(fuzzed)
		fuzzed = NewMaintenanceWindow(fuzzed.MaintenanceWindow)

		toAtlasResult := toAtlas(fuzzed)
		fromAtlasResult := fromAtlas(toAtlasResult)

		equals := fuzzed.EqualTo(fromAtlasResult)
		if !equals {
			t.Log(cmp.Diff(fuzzed, fromAtlasResult))
		}
		require.True(t, equals)
	}
}
