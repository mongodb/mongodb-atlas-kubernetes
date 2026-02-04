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

package main

import (
	"bytes"
	"encoding/json"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterLabelsContain(t *testing.T) {
	tests := []struct {
		name           string
		labels         []string
		substr         string
		expectedResult []string
	}{
		{
			name:           "Single match",
			labels:         []string{"atlas-gov", "atlas", "cloud"},
			substr:         "gov",
			expectedResult: []string{"atlas-gov"},
		},
		{
			name:           "Multiple matches",
			labels:         []string{"atlas-gov", "atlas-gov-cloud", "cloud"},
			substr:         "gov",
			expectedResult: []string{"atlas-gov", "atlas-gov-cloud"},
		},
		{
			name:           "No matches",
			labels:         []string{"atlas", "cloud"},
			substr:         "gov",
			expectedResult: []string{},
		},
		{
			name:           "Empty labels",
			labels:         []string{},
			substr:         "gov",
			expectedResult: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterLabelsContain(tt.labels, tt.substr)
			if len(result) != len(tt.expectedResult) {
				t.Errorf("Test %s failed: expected %v, got %v", tt.name, tt.expectedResult, result)
			}
			for _, label := range tt.expectedResult {
				found := slices.Contains(result, label)
				if !found {
					t.Errorf("Test %s failed: expected %v to be in the result", tt.name, label)
				}
			}
		})
	}
}

func TestFilterLabelsDoNotContain(t *testing.T) {
	tests := []struct {
		name           string
		labels         []string
		substr         string
		expectedResult []string
	}{
		{
			name:           "Single exclusion",
			labels:         []string{"atlas-gov", "atlas", "cloud"},
			substr:         "gov",
			expectedResult: []string{"atlas", "cloud"},
		},
		{
			name:           "Multiple exclusions",
			labels:         []string{"atlas-gov", "atlas-gov-cloud", "cloud"},
			substr:         "gov",
			expectedResult: []string{"cloud"},
		},
		{
			name:           "No exclusions",
			labels:         []string{"atlas", "cloud"},
			substr:         "gov",
			expectedResult: []string{"atlas", "cloud"},
		},
		{
			name:           "Empty labels",
			labels:         []string{},
			substr:         "gov",
			expectedResult: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterLabelsDoNotContain(tt.labels, tt.substr)
			if len(result) != len(tt.expectedResult) {
				t.Errorf("Test %s failed: expected %v, got %v", tt.name, tt.expectedResult, result)
			}
			for _, label := range tt.expectedResult {
				found := slices.Contains(result, label)
				if !found {
					t.Errorf("Test %s failed: expected %v to be in the result", tt.name, label)
				}
			}
		})
	}
}

func TestMatchWildcards(t *testing.T) {
	tests := []struct {
		name           string
		prLabels       []string
		testLabels     []string
		testType       string
		expectedResult []string
	}{
		{
			name:           "Exact Match Test",
			prLabels:       []string{"test/int/users"},
			testLabels:     []string{"users", "books", "orders"},
			testType:       "int",
			expectedResult: []string{"users"},
		},
		{
			name:           "Wildcard match Test",
			prLabels:       []string{"test/int/user*"},
			testLabels:     []string{"users", "user-profiles", "books", "orders"},
			testType:       "int",
			expectedResult: []string{"users", "user-profiles"},
		},
		{
			name:           "Wildcard match all tests",
			prLabels:       []string{"test/e2e/*"},
			testLabels:     []string{"login", "signup", "users", "books", "orders"},
			testType:       "e2e",
			expectedResult: []string{"login", "signup", "users", "books", "orders"},
		},
		{
			name:           "Wildcards with prefix",
			prLabels:       []string{"test/e2e/login*"},
			testLabels:     []string{"login", "login-new", "signup", "users", "books", "orders"},
			testType:       "e2e",
			expectedResult: []string{"login", "login-new"},
		},
		{
			name:           "No tests should match",
			prLabels:       []string{"test/int/unknown"},
			testLabels:     []string{"users", "books", "orders"},
			testType:       "int",
			expectedResult: []string{},
		},
		{
			name:           "Wildcard on integration and e2e tests but for e2e",
			prLabels:       []string{"test/int/user*", "test/e2e/login*"},
			testLabels:     []string{"users", "user-profiles", "login", "signup"},
			testType:       "e2e",
			expectedResult: []string{"login"},
		},
		{
			name:           "Wildcard on integration and e2e tests but for integration",
			prLabels:       []string{"test/int/user*", "test/e2e/login*"},
			testLabels:     []string{"users", "user-profiles", "login", "signup"},
			testType:       "int",
			expectedResult: []string{"users", "user-profiles"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MatchWildcards(tt.prLabels, tt.testLabels, tt.testType)

			gotMap := make(map[string]struct{})
			for _, label := range result {
				gotMap[label] = struct{}{}
			}

			expectedMap := make(map[string]struct{})
			for _, label := range tt.expectedResult {
				expectedMap[label] = struct{}{}
			}

			if len(gotMap) != len(expectedMap) {
				t.Errorf("Test %s failed: expected %v, got %v", tt.name, tt.expectedResult, result)
			}

			for label := range expectedMap {
				if _, found := gotMap[label]; !found {
					t.Errorf("Test %s failed: expected %v to be in the result.\r\n\tExpected: %s\r\n\tGot: %s",
						tt.name, label, jsonDump(expectedMap), jsonDump(gotMap))
				}
			}
		})
	}
}

func TestComputeTestLabel(t *testing.T) {
	outputJSON := true
	for _, tc := range []struct {
		name   string
		inputs labelSet
		want   string
	}{
		{
			name: "empty input renders nothing",
			inputs: labelSet{
				prLabels: "[]",
			},
			want: `{"e2e":[],"e2e2":[],"e2e_gov":[],"e2e_helm":[],"int":[]}` + "\n",
		},
		{
			name: "e2e2 explicit name is targeted",
			inputs: labelSet{
				prLabels:   `["test/e2e2/some-test"]`,
				e2e2Labels: `["some-test"]`,
			},
			want: `{"e2e":[],"e2e2":["some-test"],"e2e_gov":[],"e2e_helm":[],"int":[]}` + "\n",
		},
		{
			name: "e2e2 wildcard ",
			inputs: labelSet{
				prLabels:   `["test/e2e2/some*"]`,
				e2e2Labels: `["some-other-test"]`,
			},
			want: `{"e2e":[],"e2e2":["some-other-test"],"e2e_gov":[],"e2e_helm":[],"int":[]}` + "\n",
		},
		{
			name: "helm tests are separated from e2e tests - exact match",
			inputs: labelSet{
				prLabels:  `["test/e2e/helm-ns"]`,
				e2eLabels: `["deployment", "helm-ns", "helm-wide", "network-peering"]`,
			},
			want: `{"e2e":[],"e2e2":[],"e2e_gov":[],"e2e_helm":["helm-ns"],"int":[]}` + "\n",
		},
		{
			name: "helm tests are separated from e2e tests - wildcard all",
			inputs: labelSet{
				prLabels:  `["test/e2e/*"]`,
				e2eLabels: `["deployment", "helm-ns", "helm-wide", "network-peering"]`,
			},
			// Order can be different
			want: `{"e2e":["deployment","network-peering"],"e2e2":[],"e2e_gov":[],"e2e_helm":["helm-ns","helm-wide"],"int":[]}` + "\n",
		},
		{
			name: "helm wildcard pattern matches only helm tests",
			inputs: labelSet{
				prLabels:  `["test/e2e/helm-*"]`,
				e2eLabels: `["deployment", "helm-ns", "helm-wide", "helm-update", "network-peering"]`,
			},
			// Order can be different
			want: `{"e2e":[],"e2e2":[],"e2e_gov":[],"e2e_helm":["helm-ns","helm-update","helm-wide"],"int":[]}` + "\n",
		},
		{
			name: "non-helm e2e test does not match helm",
			inputs: labelSet{
				prLabels:  `["test/e2e/deployment"]`,
				e2eLabels: `["deployment", "helm-ns", "network-peering"]`,
			},
			want: `{"e2e":["deployment"],"e2e2":[],"e2e_gov":[],"e2e_helm":[],"int":[]}` + "\n",
		},
		{
			name: "helm tests extracted before gov tests",
			inputs: labelSet{
				prLabels:  `["test/e2e/*"]`,
				e2eLabels: `["deployment", "helm-ns", "atlas-gov-peering", "network-peering"]`,
			},
			want: `{"e2e":["deployment","network-peering"],"e2e2":[],"e2e_gov":["atlas-gov-peering"],"e2e_helm":["helm-ns"],"int":[]}` + "\n",
		},
		{
			name: "multiple label patterns with helm",
			inputs: labelSet{
				prLabels:  `["test/e2e/deployment", "test/e2e/helm-ns"]`,
				e2eLabels: `["deployment", "helm-ns", "helm-wide", "network-peering"]`,
			},
			want: `{"e2e":["deployment"],"e2e2":[],"e2e_gov":[],"e2e_helm":["helm-ns"],"int":[]}` + "\n",
		},
		{
			name: "skip prefixes applied to helm tests",
			inputs: labelSet{
				prLabels:     `["test/e2e/*"]`,
				e2eLabels:    `["deployment", "helm-ns", "focus-helm-test", "network-peering"]`,
				skipPrefixes: `["focus"]`,
			},
			want: `{"e2e":["deployment","network-peering"],"e2e2":[],"e2e_gov":[],"e2e_helm":["helm-ns"],"int":[]}` + "\n",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.NewBufferString("")
			err := computeTestLabels(buf, outputJSON, &tc.inputs)
			require.NoError(t, err)

			var want, got map[string]any
			require.NoError(t, json.Unmarshal([]byte(tc.want), &want))
			require.NoError(t, json.Unmarshal(buf.Bytes(), &got))

			assert.Equal(t, len(want), len(got), "Number of keys should match")
			for key, wantVal := range want {
				gotVal, exists := got[key]
				assert.True(t, exists, "Key %s should exist in result", key)

				wantArray, wantIsArray := wantVal.([]any)
				gotArray, gotIsArray := gotVal.([]any)
				if wantIsArray && gotIsArray {
					assert.ElementsMatch(t, wantArray, gotArray, "Arrays for key %s should match (ignoring order)", key)
				} else {
					assert.Equal(t, wantVal, gotVal, "Value for key %s should match", key)
				}
			}
		})
	}
}

func TestSkipLabelsByPrefix(t *testing.T) {
	for _, tc := range []struct {
		name         string
		labels       []string
		skipPrefixes []string
		want         []string
	}{
		{
			name:         "no filtering",
			labels:       []string{"a", "b", "c"},
			skipPrefixes: []string{},
			want:         []string{"a", "b", "c"},
		},
		{
			name:         "several prefixes matched",
			labels:       []string{"something", "anotherthing", "andanother"},
			skipPrefixes: []string{"something", "another"},
			want:         []string{"andanother"},
		},
		{
			name:         "no matches",
			labels:       []string{"something", "anotherthing", "andanother"},
			skipPrefixes: []string{"blah"},
			want:         []string{"something", "anotherthing", "andanother"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := SkipLabelsByPrefix(tc.labels, tc.skipPrefixes)
			assert.Equal(t, tc.want, got)
		})
	}
}
