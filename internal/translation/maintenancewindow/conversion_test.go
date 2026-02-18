// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

	for range 100 {
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
