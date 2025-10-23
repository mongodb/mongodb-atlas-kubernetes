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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	admin2025 "go.mongodb.org/atlas-sdk/v20250312006/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/autogen/translate/unstructured"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

type testStruct struct {
	ID        int
	YesNo     bool
	Text      string
	Data      float32
	Timestamp time.Time
	SubStruct []subStruct
}

type subStruct struct {
	SubID   int
	YesNo   bool
	Text    string
	Data    float32
	TheTime time.Time
}

var sample = testStruct{
	ID:        12,
	YesNo:     true,
	Text:      "some text",
	Data:      2.4,
	Timestamp: time.Now().UTC().Round(time.Nanosecond),
	SubStruct: []subStruct{
		{
			SubID:   15,
			YesNo:   true,
			Text:    "more text",
			Data:    0.67,
			TheTime: time.Now().Add(-24 * 7 * time.Hour).UTC().Round(time.Nanosecond),
		},
		{
			SubID:   14,
			YesNo:   false,
			Text:    "and even more text",
			Data:    1.67,
			TheTime: time.Now().Add(-24 * 15 * time.Hour).UTC().Round(time.Nanosecond),
		},
	},
}

func TestParamsFill(t *testing.T) {
	unstructuredSample := map[string]any{
		"groupid": "62b6e34b3d91647abb20e7b8",
		"groupalertsconfig": map[string]any{
			"enabled":       true,
			"eventtypename": "some-event",
		},
	}
	result := admin2025.CreateAlertConfigurationApiParams{}
	require.NoError(t, unstructured.FromUnstructured(&result, unstructuredSample))
	assert.Equal(t, admin2025.CreateAlertConfigurationApiParams{
		GroupId: "62b6e34b3d91647abb20e7b8",
		GroupAlertsConfig: &admin2025.GroupAlertsConfig{
			Enabled:       pointer.MakePtr(true),
			EventTypeName: pointer.MakePtr("some-event"),
		},
	}, result)
}

func TestToAndFromAndCopyUnstructured(t *testing.T) {
	sampleMap, err := unstructured.ToUnstructured(sample)
	require.NoError(t, err)
	clone := map[string]any{}
	unstructured.CopyFields(clone, sampleMap)
	result := testStruct{}
	require.NoError(t, unstructured.FromUnstructured(&result, sampleMap))
	assert.Equal(t, sample, result)
}

func TestCreateAndAccessField(t *testing.T) {
	sampleMap, err := unstructured.ToUnstructured(sample)
	require.NoError(t, err)
	want := "value"
	require.NoError(t, unstructured.CreateField(sampleMap, "value", "SubObj", "field"))
	got, err := unstructured.GetField[string](sampleMap, "SubObj", "field")
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestAsPath(t *testing.T) {
	for _, tc := range []struct {
		title string
		input string
		want  []string
	}{
		{
			title: "nothing",
			input: "",
			want:  []string{""},
		},
		{
			title: "single",
			input: "dir",
			want:  []string{"dir"},
		},
		{
			title: "double",
			input: "dir0.dir1",
			want:  []string{"dir0", "dir1"},
		},
		{
			title: "rooted double",
			input: ".dir0.dir1",
			want:  []string{"dir0", "dir1"},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			assert.Equal(t, tc.want, unstructured.AsPath(tc.input))
		})
	}
}

func TestBase(t *testing.T) {
	for _, tc := range []struct {
		title string
		input []string
		want  string
	}{
		{
			title: "none",
			input: []string{},
			want:  "",
		},
		{
			title: "single",
			input: []string{"dir"},
			want:  "dir",
		},
		{
			title: "double",
			input: []string{"dir0", "dir1"},
			want:  "dir1",
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			assert.Equal(t, tc.want, unstructured.Base(tc.input))
		})
	}
}

func TestSkipKeysAndFieldsOf(t *testing.T) {
	for _, tc := range []struct {
		title string
		skips []string
		want  []string
	}{
		{
			title: "no skips",
			skips: []string{},
			want:  []string{"ID", "YesNo", "Text", "Data", "Timestamp", "SubStruct"},
		},
		{
			title: "no matchs",
			skips: []string{"bleh", "ah"},
			want:  []string{"ID", "YesNo", "Text", "Data", "Timestamp", "SubStruct"},
		},
		{
			title: "last missing",
			skips: []string{"SubStruct"},
			want:  []string{"ID", "YesNo", "Text", "Data", "Timestamp"},
		},
		{
			title: "all gone missing",
			skips: []string{"ID", "YesNo", "Text", "Data", "Timestamp", "SubStruct"},
			want:  []string{},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			sampleMap, err := unstructured.ToUnstructured(sample)
			require.NoError(t, err)
			sort.Strings(tc.want)
			filtered := unstructured.SkipKeys(sampleMap, tc.skips...)
			fields := unstructured.FieldsOf(filtered)
			sort.Strings(fields)
			assert.Equal(t, tc.want, fields)
		})
	}
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
