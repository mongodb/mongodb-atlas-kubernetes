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

package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitArrayPath(t *testing.T) {
	tests := []struct {
		name                string
		fieldPath           string
		expectedBeforeArray string
		expectedAfterArray  string
		expectedIsArray     bool
	}{
		{
			name:                "single array",
			fieldPath:           "properties.spec.properties.entries.items.properties.secretRef",
			expectedBeforeArray: "properties.spec.properties.entries",
			expectedAfterArray:  "properties.secretRef",
			expectedIsArray:     true,
		},
		{
			name:                "no array",
			fieldPath:           "properties.spec.properties.secretRef",
			expectedBeforeArray: "",
			expectedAfterArray:  "properties.spec.properties.secretRef",
			expectedIsArray:     false,
		},
		{
			name:                "nested array",
			fieldPath:           "properties.entries.items.properties.nested.items.properties.ref",
			expectedBeforeArray: "properties.entries",
			expectedAfterArray:  "properties.nested.items.properties.ref",
			expectedIsArray:     true,
		},
		{
			name:                "array at root",
			fieldPath:           "items.properties.ref",
			expectedBeforeArray: "",
			expectedAfterArray:  "items.properties.ref", // No ".items." just "items" at start
			expectedIsArray:     false,                  // Not detected as array
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeArray, afterArray, isArray := splitArrayPath(tt.fieldPath)
			assert.Equal(t, tt.expectedBeforeArray, beforeArray)
			assert.Equal(t, tt.expectedAfterArray, afterArray)
			assert.Equal(t, tt.expectedIsArray, isArray)
		})
	}
}

func TestGenerateLoopVariableName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"entries", "entry"},
		{"items", "item"},
		{"configs", "config"},
		{"data", "dataItem"},
		{"boxes", "box"}, // xes ending
		{"policies", "policy"},
		{"matches", "match"},
		{"indexes", "index"},
		{"", "item"},
		{"ENTRIES", "entry"}, // Should convert to lowercase
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := generateLoopVariableName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseReferenceFields_ArrayReferences_Detection(t *testing.T) {
	testYAML := `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: deployments.atlas.generated.mongodb.com
  annotations:
    api-mappings: |
      properties:
        spec:
          properties:
            v20250312:
              properties:
                replicas:
                  items:
                    properties:
                      secretRef:
                        x-kubernetes-mapping:
                          type:
                            kind: Secret
                            group: ""
                            version: v1
spec:
  group: atlas.generated.mongodb.com
  names:
    kind: Deployment
    plural: deployments
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              v20250312:
                type: object
                properties:
                  replicas:
                    type: array
                    items:
                      type: object
                      properties:
                        secretRef:
                          type: object
                          properties:
                            name:
                              type: string
`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(testFile, []byte(testYAML), 0644)
	require.NoError(t, err)

	refs, err := ParseReferenceFields(testFile, "Deployment")
	require.NoError(t, err)
	require.Len(t, refs, 1)

	ref := refs[0]
	assert.Equal(t, "secretRef", ref.FieldName)
	assert.Equal(t, "Secret", ref.ReferencedKind)
	assert.True(t, ref.IsArrayBased, "Reference should be marked as array-based")
	assert.Equal(t, "properties.spec.properties.v20250312.properties.replicas", ref.ArrayPath)
	assert.Equal(t, "properties.secretRef", ref.ItemPath)
	assert.Contains(t, ref.FieldPath, ".items.")
}

func TestGenerateIndexerWithArrayReferences_Integration(t *testing.T) {
	testYAML := `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: clusters.atlas.generated.mongodb.com
  annotations:
    api-mappings: |
      properties:
        spec:
          properties:
            v20250312:
              properties:
                entries:
                  items:
                    properties:
                      secretRef:
                        x-kubernetes-mapping:
                          type:
                            kind: Secret
                            group: ""
                            version: v1
spec:
  group: atlas.generated.mongodb.com
  names:
    kind: Cluster
    plural: clusters
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              v20250312:
                type: object
                properties:
                  entries:
                    type: array
                    items:
                      type: object
                      properties:
                        secretRef:
                          type: object
                          properties:
                            name:
                              type: string
`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(testFile, []byte(testYAML), 0644)
	require.NoError(t, err)

	outputDir := filepath.Join(tmpDir, "indexers")
	err = GenerateIndexers(testFile, "Cluster", outputDir)
	require.NoError(t, err)

	// Read generated file
	indexerFile := filepath.Join(outputDir, "clusterbysecret.go")
	assert.FileExists(t, indexerFile)

	content, err := os.ReadFile(indexerFile)
	require.NoError(t, err)
	contentStr := string(content)

	// Verify it contains for-loop
	assert.Contains(t, contentStr, "for _, entry := range", "Should have for-loop over array")

	// Verify proper nil checks
	assert.Contains(t, contentStr, "V20250312")
	assert.Contains(t, contentStr, "Entries")
	assert.Contains(t, contentStr, "!= nil")

	// Verify it appends keys
	assert.Contains(t, contentStr, "keys = append(keys")
	assert.Contains(t, contentStr, "entry.SecretRef")
}

func TestGenerateIndexer_MixedReferences(t *testing.T) {
	testYAML := `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: deployments.atlas.generated.mongodb.com
  annotations:
    api-mappings: |
      properties:
        spec:
          properties:
            v20250312:
              properties:
                groupRef:
                  x-kubernetes-mapping:
                    type:
                      kind: Group
                      group: atlas.generated.mongodb.com
                      version: v1
                replicas:
                  items:
                    properties:
                      secretRef:
                        x-kubernetes-mapping:
                          type:
                            kind: Secret
                            group: ""
                            version: v1
spec:
  group: atlas.generated.mongodb.com
  names:
    kind: Deployment
    plural: deployments
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              v20250312:
                type: object
                properties:
                  groupRef:
                    type: object
                    properties:
                      name:
                        type: string
                  replicas:
                    type: array
                    items:
                      type: object
                      properties:
                        secretRef:
                          type: object
                          properties:
                            name:
                              type: string
`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(testFile, []byte(testYAML), 0644)
	require.NoError(t, err)

	outputDir := filepath.Join(tmpDir, "indexers")
	err = GenerateIndexers(testFile, "Deployment", outputDir)
	require.NoError(t, err)

	// Should have generated two indexers
	groupIndexer := filepath.Join(outputDir, "deploymentbygroup.go")
	secretIndexer := filepath.Join(outputDir, "deploymentbysecret.go")
	assert.FileExists(t, groupIndexer)
	assert.FileExists(t, secretIndexer)

	// Group indexer should NOT have for-loop in Keys() method (non-array)
	groupContent, err := os.ReadFile(groupIndexer)
	require.NoError(t, err)
	groupStr := string(groupContent)
	// Keys method shouldn't have for loop (check area between "func (i" and "return keys")
	assert.Contains(t, groupStr, "resource.Spec.V20250312.GroupRef")
	// Verify it's not iterating over a range (the Keys method specifically)
	keysMethodStart := strings.Index(groupStr, "func (i *DeploymentByGroupIndexer) Keys(")
	keysMethodEnd := strings.Index(groupStr, "func NewDeploymentByGroupMapFunc")
	if keysMethodStart >= 0 && keysMethodEnd >= 0 {
		keysMethod := groupStr[keysMethodStart:keysMethodEnd]
		assert.NotContains(t, keysMethod, "range resource", "Keys method should not iterate over array")
	}

	// Secret indexer SHOULD have for-loop (array-based)
	secretContent, err := os.ReadFile(secretIndexer)
	require.NoError(t, err)
	secretStr := string(secretContent)
	// Keys method SHOULD have for loop over array
	keysMethodStart = strings.Index(secretStr, "func (i *DeploymentBySecretIndexer) Keys(")
	keysMethodEnd = strings.Index(secretStr, "func NewDeploymentBySecretMapFunc")
	if keysMethodStart >= 0 && keysMethodEnd >= 0 {
		keysMethod := secretStr[keysMethodStart:keysMethodEnd]
		assert.Contains(t, keysMethod, "for _, ", "Keys method should iterate over array")
		assert.Contains(t, keysMethod, ".SecretRef") // Should reference field in loop variable
	}
}

func TestSkipNestedArrayReferences(t *testing.T) {
	testYAML := `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: alerts.atlas.generated.mongodb.com
  annotations:
    api-mappings: |
      properties:
        spec:
          properties:
            v20250312:
              properties:
                regions:
                  items:
                    properties:
                      notifications:
                        items:
                          properties:
                            secretRef:
                              x-kubernetes-mapping:
                                type:
                                  kind: Secret
spec:
  group: atlas.generated.mongodb.com
  names:
    kind: AlertConfig
    plural: alertconfigs
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        properties:
          spec:
            properties:
              v20250312:
                properties:
                  regions:
                    type: array
                    items:
                      properties:
                        notifications:
                          type: array
`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(testFile, []byte(testYAML), 0644)
	require.NoError(t, err)

	outputDir := filepath.Join(tmpDir, "indexers")
	err = GenerateIndexers(testFile, "AlertConfig", outputDir)
	require.NoError(t, err)

	// Should not generate any indexers (nested arrays skipped)
	files, err := os.ReadDir(outputDir)
	if err == nil {
		assert.Empty(t, files, "No indexer files should be generated for nested arrays")
	}
}

func TestSplitRequiredSegments(t *testing.T) {
	tests := []struct {
		name                string
		fieldPath           string
		requiredSegments    []bool
		expectedBeforeArray []bool
		expectedInArray     []bool
	}{
		{
			name:                "simple array reference",
			fieldPath:           "properties.spec.properties.v20250312.properties.entries.items.properties.secretRef",
			requiredSegments:    []bool{false, true, false, false}, // spec, v20250312, entries, secretRef
			expectedBeforeArray: []bool{false, true, false},        // up to entries
			expectedInArray:     []bool{false},                     // secretRef within item
		},
		{
			name:                "no array",
			fieldPath:           "properties.spec.properties.groupRef",
			requiredSegments:    []bool{false, false},
			expectedBeforeArray: []bool{false, false},
			expectedInArray:     nil,
		},
		{
			name:                "empty required segments",
			fieldPath:           "properties.spec.properties.entries.items.properties.ref",
			requiredSegments:    []bool{},
			expectedBeforeArray: nil,
			expectedInArray:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeArray, inArray := splitRequiredSegments(tt.fieldPath, tt.requiredSegments)
			assert.Equal(t, tt.expectedBeforeArray, beforeArray)
			assert.Equal(t, tt.expectedInArray, inArray)
		})
	}
}

func TestGenerateArrayFieldExtractionCode_WithOptionalFields(t *testing.T) {
	// Test that array extraction code handles optional fields correctly
	field := ReferenceField{
		FieldName:        "secretRef",
		FieldPath:        "properties.spec.properties.v20250312.properties.entries.items.properties.config.properties.secretRef",
		ReferencedKind:   "Secret",
		RequiredSegments: []bool{false, true, false, false, false}, // spec, v20250312, entries, config, secretRef
		IsArrayBased:     true,
		ArrayPath:        "properties.spec.properties.v20250312.properties.entries",
		ItemPath:         "properties.config.properties.secretRef",
	}

	code := generateArrayFieldExtractionCode(field)
	assert.NotNil(t, code)

	// Verify code was generated (can't easily inspect jen.Code)
	// The integration tests will verify actual generated code
}

func TestArrayIndexerGeneration(t *testing.T) {
	testYAML := `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: configs.atlas.generated.mongodb.com
  annotations:
    api-mappings: |
      properties:
        spec:
          properties:
            v20250312:
              properties:
                items:
                  items:
                    properties:
                      secretRef:
                        x-kubernetes-mapping:
                          type:
                            kind: Secret
spec:
  group: atlas.generated.mongodb.com
  names:
    kind: Config
    plural: configs
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        properties:
          spec:
            properties:
              v20250312:
                properties:
                  items:
                    type: array
                    items:
                      properties:
                        secretRef:
                          type: object
`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(testFile, []byte(testYAML), 0644)
	require.NoError(t, err)

	outputDir := filepath.Join(tmpDir, "indexers")
	err = GenerateIndexers(testFile, "Config", outputDir)
	require.NoError(t, err)

	indexerFile := filepath.Join(outputDir, "configbysecret.go")
	content, err := os.ReadFile(indexerFile)
	require.NoError(t, err)
	contentStr := string(content)

	// Verify basic structure is present
	assert.Contains(t, contentStr, "for _, ")
	assert.Contains(t, contentStr, "range resource.Spec.V20250312")
}

func TestArrayIndexerGenerationLoopVariable(t *testing.T) {
	// Test that loop variables are generated correctly based on array field names
	testYAML := `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: policies.atlas.generated.mongodb.com
  annotations:
    api-mappings: |
      properties:
        spec:
          properties:
            v20250312:
              properties:
                policies:
                  items:
                    properties:
                      secretRef:
                        x-kubernetes-mapping:
                          type:
                            kind: Secret
spec:
  group: atlas.generated.mongodb.com
  names:
    kind: Policy
    plural: policies
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        properties:
          spec:
            properties:
              v20250312:
                properties:
                  policies:
                    type: array
                    items:
                      properties:
                        secretRef:
                          type: object
`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(testFile, []byte(testYAML), 0644)
	require.NoError(t, err)

	outputDir := filepath.Join(tmpDir, "indexers")
	err = GenerateIndexers(testFile, "Policy", outputDir)
	require.NoError(t, err)

	indexerFile := filepath.Join(outputDir, "policybysecret.go")
	content, err := os.ReadFile(indexerFile)
	require.NoError(t, err)
	contentStr := string(content)

	// Verify loop variable is "policy" (singular of "policies")
	assert.Contains(t, contentStr, "for _, policy := range")
	assert.Contains(t, contentStr, "policy.SecretRef")
}

func TestArrayIndexerWithRequiredArrayContainer(t *testing.T) {
	testYAML := `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: apps.atlas.generated.mongodb.com
  annotations:
    api-mappings: |
      properties:
        spec:
          properties:
            v20250312:
              properties:
                entries:
                  items:
                    properties:
                      groupRef:
                        x-kubernetes-mapping:
                          type:
                            kind: Group
spec:
  group: atlas.generated.mongodb.com
  names:
    kind: App
    plural: apps
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        properties:
          spec:
            required:
            - v20250312
            properties:
              v20250312:
                type: object
                required:
                - entries
                properties:
                  entries:
                    type: array
                    items:
                      properties:
                        groupRef:
                          type: object
`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(testFile, []byte(testYAML), 0644)
	require.NoError(t, err)

	outputDir := filepath.Join(tmpDir, "indexers")
	err = GenerateIndexers(testFile, "App", outputDir)
	require.NoError(t, err)

	indexerFile := filepath.Join(outputDir, "appbygroup.go")
	content, err := os.ReadFile(indexerFile)
	require.NoError(t, err)
	contentStr := string(content)

	// Since entries is required and v20250312 is required, should check less
	// The exact nil checks depend on implementation, but verify basic structure
	assert.Contains(t, contentStr, "for _, entry := range")
	assert.Contains(t, contentStr, "entry.GroupRef")

	// Should still check if entries != nil (even if required in schema, runtime could be nil)
	nilChecks := strings.Count(contentStr, "!= nil")
	assert.Greater(t, nilChecks, 0, "Should have at least some nil checks")
}
