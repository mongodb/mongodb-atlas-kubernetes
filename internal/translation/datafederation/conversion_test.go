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

package datafederation

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"
)

func TestRoundtrip_DataFederation(t *testing.T) {
	f := fuzz.New()

	for range 100 {
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
