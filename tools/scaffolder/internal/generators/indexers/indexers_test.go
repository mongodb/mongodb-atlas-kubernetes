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

// Package generate contains tests for indexer utility functions.
// For parsing tests, see indexers_parsing_test.go
// For generation/integration tests, see indexers_generation_test.go
// For array-related tests, see indexers_array_test.go
// For nil check tests, see indexers_nil_check_test.go
// For dependent reference tests, see indexers_dependents_test.go
package indexers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildFieldAccessPath(t *testing.T) {
	tests := []struct {
		name      string
		fieldPath string
		expected  string
	}{
		{
			name:      "Simple field path",
			fieldPath: "properties.spec.properties.v20250312.properties.groupRef",
			expected:  "resource.Spec.V20250312.GroupRef",
		},
		{
			name:      "Field path with array items",
			fieldPath: "properties.spec.properties.v20250312.properties.notifications.items.properties.secretRef",
			expected:  "resource.Spec.V20250312.Notifications.SecretRef",
		},
		{
			name:      "Nested field path",
			fieldPath: "properties.spec.properties.entry.properties.apiKeyRef",
			expected:  "resource.Spec.Entry.ApiKeyRef",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildFieldAccessPath(tt.fieldPath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCapitalizeFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"groupRef", "GroupRef"},
		{"v20250312", "V20250312"},
		{"spec", "Spec"},
		{"", ""},
		{"a", "A"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := CapitalizeFirst(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCreateIndexerInfoForKind(t *testing.T) {
	refs := []ReferenceField{
		{
			FieldName:      "groupRef",
			FieldPath:      "properties.spec.properties.v20250312.properties.groupRef",
			ReferencedKind: "Group",
		},
	}

	indexer := createIndexerInfoForKind("Cluster", "Group", refs)

	assert.Equal(t, "Group", indexer.TargetKind)
	assert.Equal(t, "ClusterByGroupIndex", indexer.ConstantName)
	assert.Equal(t, "Cluster", indexer.ResourceName)
	assert.Equal(t, "cluster.groupRef", indexer.IndexerName)
	assert.Equal(t, refs, indexer.ReferenceFields)
}
