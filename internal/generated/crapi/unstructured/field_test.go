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

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/crapi/unstructured"
)

func TestCreateAndAccessField(t *testing.T) {
	sampleMap, err := unstructured.ToUnstructured(sample)
	require.NoError(t, err)
	want := "value"
	require.NoError(t, unstructured.RecursiveCreateField(sampleMap, "value", "SubObj", "field"))
	got, err := unstructured.GetField[string](sampleMap, "SubObj", "field")
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestGetOrCreateField(t *testing.T) {
	for _, tc := range []struct {
		title        string
		obj          map[string]any
		path         []string
		defaultValue any
		want         any
		wantErr      string
	}{
		{
			title: "create missing leaf",
			obj: map[string]any{
				"deep": map[string]any{},
			},
			path:         []string{"deep", "field"},
			defaultValue: "some string",
			want:         "some string",
		},
		{
			title:        "fail to create branch",
			obj:          map[string]any{},
			path:         []string{"deep", "field"},
			defaultValue: "some string",
			wantErr:      "field \"deep\": not found",
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
			wantErr:      "field \"deep\": not an object",
		},
		{
			title: "find array",
			obj: map[string]any{
				"deep": []map[string]any{
					{
						"key-a":       "A",
						"other-field": "-",
					},
					{
						"key-b":       "B",
						"other-field": "value",
					},
				},
			},
			path: []string{"deep"},
			want: []map[string]any{
				{
					"key-a":       "A",
					"other-field": "-",
				},
				{
					"key-b":       "B",
					"other-field": "value",
				},
			},
		},
		{
			title: "find in array",
			obj: map[string]any{
				"deep": []any{
					map[string]any{
						"key-a":       "A",
						"other-field": "-",
					},
					map[string]any{
						"key-b":       "B",
						"other-field": "value",
					},
				},
			},
			path: []string{"deep", "[]", "key-b"},
			want: "B",
		},
		{
			title: "created in array",
			obj: map[string]any{
				"deep": []any{
					map[string]any{
						"key-a":       "A",
						"other-field": "-",
					},
				},
			},
			path: []string{"deep", "[]", "key-b"},
			defaultValue: map[string]any{
				"key-b":       "B",
				"other-field": "value",
			},
			want: map[string]any{
				"key-b":       "B",
				"other-field": "value",
			},
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

func TestRecursiveCreateField(t *testing.T) {
	for _, tc := range []struct {
		title   string
		obj     map[string]any
		path    []string
		value   any
		want    map[string]any
		wantErr string
	}{
		{
			title: "creates full branch",
			obj:   map[string]any{},
			path:  []string{"deep", "field"},
			value: "some string",
			want: map[string]any{
				"deep": map[string]any{
					"field": "some string",
				},
			},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			err := unstructured.RecursiveCreateField(tc.obj, tc.value, tc.path...)
			if tc.wantErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, tc.obj)
			} else {
				assert.ErrorContains(t, err, tc.wantErr)
			}
		})
	}
}

func TestGetFieldObject(t *testing.T) {
	for _, tc := range []struct {
		title   string
		obj     map[string]any
		path    []string
		want    map[string]any
		wantErr string
	}{
		{
			title: "grabs holder of field",
			obj: map[string]any{
				"deep": map[string]any{
					"field": map[string]any{
						"subfield": "something",
					},
				},
			},
			path: []string{"deep", "field", "subfield"},
			want: map[string]any{
				"subfield": "something",
			},
		},
		{
			title: "grabs holder of field in array",
			obj: map[string]any{
				"deep": map[string]any{
					"fields": []any{
						map[string]any{
							"notit": "something else",
						},
						map[string]any{
							"matchingKey": "value",
							"otherField":  "otherValue",
						},
					},
				},
			},
			path: []string{"deep", "fields", "[]", "matchingKey"},
			want: map[string]any{
				"matchingKey": "value",
				"otherField":  "otherValue",
			},
		},
		{
			title: "fails to grabs holder if field not in",
			obj: map[string]any{
				"deep": map[string]any{
					"field": map[string]any{
						"subfield": "something",
					},
				},
			},
			path:    []string{"deep", "field", "notSubfield"},
			wantErr: "not found",
		},
		{
			title: "fails to grabs holder if field not in array",
			obj: map[string]any{
				"deep": map[string]any{
					"fields": []any{
						map[string]any{
							"notit": "something else",
						},
						map[string]any{
							"nonMatchingKey": "value",
							"otherField":     "otherValue",
						},
					},
				},
			},
			path:    []string{"deep", "fields", "[]", "matchingKey"},
			wantErr: "not found",
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			got, err := unstructured.GetFieldObject(tc.obj, tc.path...)
			if tc.wantErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			} else {
				assert.ErrorContains(t, err, tc.wantErr)
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
