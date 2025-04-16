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

package cmp

import (
	"reflect"
	"testing"

	v1apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestEndlessRecursion(t *testing.T) {
	for _, tc := range []struct {
		name string
		data any
		want any
	}{
		{
			name: "pointer",
			data: &struct {
				Slice []string
			}{
				Slice: []string{"C", "B", "A", "sort"},
			},
			want: &struct {
				Slice []string
			}{
				Slice: []string{"A", "B", "C", "sort"},
			},
		},
		{
			name: "nested JSON",
			data: struct {
				NestedJSON v1apiextensions.JSON
			}{
				NestedJSON: v1apiextensions.JSON{Raw: []byte("CBA")},
			},
			want: struct {
				NestedJSON v1apiextensions.JSON
			}{
				NestedJSON: v1apiextensions.JSON{Raw: []byte("CBA")},
			},
		},
		{
			name: "ignore byte slices",
			data: struct {
				ByteSlice []byte
			}{
				ByteSlice: []byte("CBA"),
			},
			want: struct {
				ByteSlice []byte
			}{
				ByteSlice: []byte("CBA"),
			},
		},
		{
			name: "ignore zero length slices",
			data: struct {
				EmptySlice []string
			}{},
			want: struct {
				EmptySlice []string
			}{},
		},
		{
			name: "ignore unexported fields",
			data: struct {
				ExportedFields  []string
				unExportedField []string
			}{
				ExportedFields:  []string{"C", "B", "A", "sort"},
				unExportedField: []string{"C", "B", "A", `don't sort'`},
			},
			want: struct {
				ExportedFields  []string
				unExportedField []string
			}{
				ExportedFields:  []string{"A", "B", "C", "sort"},
				unExportedField: []string{"C", "B", "A", `don't sort'`},
			},
		},
		{
			name: "nested slices",
			data: struct {
				Nested []struct {
					Slice []string
				}
			}{
				Nested: []struct {
					Slice []string
				}{
					{Slice: []string{"Z", "Y", "X"}},
					{Slice: []string{"C", "B", "A"}},
				},
			},
			want: struct {
				Nested []struct {
					Slice []string
				}
			}{
				Nested: []struct {
					Slice []string
				}{
					{Slice: []string{"A", "B", "C"}},
					{Slice: []string{"X", "Y", "Z"}},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := Normalize(tc.data)
			if err != nil {
				t.Error(err)
				return
			}
			if !reflect.DeepEqual(tc.data, tc.want) {
				t.Errorf("want normalized value %+v, got %v", tc.want, tc.data)
				return
			}
		})
	}
}
