// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package specmatch_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/specmatch"
)

type container struct {
	ProviderName   *string   `json:"providerName,omitempty"`
	AtlasCidrBlock *string   `json:"atlasCidrBlock,omitempty"`
	Regions        *[]string `json:"regions,omitempty"`
	Id             *string   `json:"id,omitempty"`
}

func TestAtlasMatchesSpec(t *testing.T) {
	for _, tc := range []struct {
		title   string
		spec    container
		atlas   container
		opts    []specmatch.Option
		wantErr string
	}{
		{
			title: "fields Atlas adds are not compared",
			spec:  container{ProviderName: new("AWS")},
			atlas: container{ProviderName: new("AWS"), Id: new("abc")},
		},
		{
			title:   "changed value is reported with its path",
			spec:    container{AtlasCidrBlock: new("10.0.0.0/21")},
			atlas:   container{AtlasCidrBlock: new("10.1.0.0/21")},
			wantErr: "$.atlasCidrBlock: want 10.0.0.0/21, got 10.1.0.0/21",
		},
		{
			title:   "field dropped on the way to Atlas is reported",
			spec:    container{ProviderName: new("AWS"), AtlasCidrBlock: new("10.0.0.0/21")},
			atlas:   container{ProviderName: new("AWS")},
			wantErr: "$.atlasCidrBlock: want 10.0.0.0/21, missing in Atlas",
		},
		{
			title:   "list order matters by default",
			spec:    container{Regions: &[]string{"A", "B"}},
			atlas:   container{Regions: &[]string{"B", "A"}},
			wantErr: "$.regions[]",
		},
		{
			title: "unordered lists match in any order",
			spec:  container{Regions: &[]string{"A", "B"}},
			atlas: container{Regions: &[]string{"B", "A"}},
			opts:  []specmatch.Option{specmatch.UnorderedLists("$.regions")},
		},
		{
			title: "ignored fields are skipped",
			spec:  container{AtlasCidrBlock: new("10.0.0.0/21")},
			atlas: container{},
			opts:  []specmatch.Option{specmatch.IgnoreFields("$.atlasCidrBlock")},
		},
		{
			title: "normalized values are compared after normalization",
			spec:  container{ProviderName: new("aws")},
			atlas: container{ProviderName: new("AWS")},
			opts: []specmatch.Option{specmatch.Normalize("$.providerName", func(v any) any {
				return strings.ToUpper(v.(string))
			})},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			err := specmatch.AtlasMatchesSpec(tc.spec, tc.atlas, tc.opts...)
			if tc.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
