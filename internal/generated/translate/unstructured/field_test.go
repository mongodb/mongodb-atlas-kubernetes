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
//

package unstructured_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/translate/unstructured"
)

func TestCreateAndAccessField(t *testing.T) {
	sampleMap, err := unstructured.ToUnstructured(sample)
	require.NoError(t, err)
	want := "value"
	require.NoError(t, unstructured.CreateField(sampleMap, "value", "SubObj", "field"))
	got, err := unstructured.GetField[string](sampleMap, "SubObj", "field")
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestAccessOrCreateField(t *testing.T) {
	for _, tc := range []struct {
		title        string
		obj          map[string]any
		path         []string
		defaultValue any
		want         any
		wantErr      string
	}{
		{
			title:        "create missing",
			obj:          map[string]any{},
			path:         []string{"deep", "field"},
			defaultValue: "some string",
			want:         "some string",
		},
		{
			title: "read existing",
			obj: map[string]any{
				"deep": map[string]any{
					"field": "other string",
				},
			},
			path:         []string{"deep", "field"},
			defaultValue: "some string",
			want:         "other string",
		},
		{
			title:        "nil object",
			obj:          nil,
			path:         []string{"deep", "field"},
			defaultValue: "some string",
			wantErr:      "nil object",
		},
		{
			title: "wrong path type",
			obj: map[string]any{
				"deep": "an string",
			},
			path:         []string{"deep", "field"},
			defaultValue: "some string",
			wantErr:      "intermediate path [deep] exists but is of type string",
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			got, err := unstructured.GetOrCreateField(tc.obj, tc.defaultValue, tc.path...)
			if tc.wantErr != "" {
				require.Nil(t, got)
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestFieldsOf(t *testing.T) {
	sampleMap, err := unstructured.ToUnstructured(sample)
	require.NoError(t, err)
	fields := unstructured.FieldsOf(sampleMap)
	want := []string{"SubStruct", "ID", "YesNo", "Text", "Data", "Timestamp"}
	sort.Strings(want)
	sort.Strings(fields)
	assert.Equal(t, want, fields)
}
