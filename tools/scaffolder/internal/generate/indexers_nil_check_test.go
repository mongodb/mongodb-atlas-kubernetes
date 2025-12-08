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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildDotChain(t *testing.T) {
	tests := []struct {
		name     string
		segments []string
		expected string
	}{
		{
			name:     "single segment",
			segments: []string{"resource"},
			expected: "resource",
		},
		{
			name:     "two segments",
			segments: []string{"resource", "Spec"},
			expected: "resource.Spec",
		},
		{
			name:     "multiple segments",
			segments: []string{"cluster", "Spec", "V20250312", "GroupRef"},
			expected: "cluster.Spec.V20250312.GroupRef",
		},
		{
			name:     "empty segments",
			segments: []string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := buildDotChain(tt.segments)
			
			// For empty segments, jen returns empty/null statement
			if len(tt.segments) == 0 {
				// Empty segments return a null statement, just verify it doesn't panic
				assert.NotNil(t, stmt)
				return
			}
			
			result := stmt.GoString()
			
			// Check that all segments appear in the result
			for _, seg := range tt.segments {
				assert.Contains(t, result, seg)
			}
		})
	}
}

func TestBuildNilCheckConditions_NoRequiredSegments(t *testing.T) {
	// Test fallback behavior when requiredSegments is empty
	fieldPath := "resource.Spec.V20250312.GroupRef"
	requiredSegments := []bool{}
	
	stmt := buildNilCheckConditions(fieldPath, requiredSegments)
	result := stmt.GoString()
	
	// Should fall back to checking the last field
	assert.Contains(t, result, "GroupRef")
	assert.Contains(t, result, "!=")
	assert.Contains(t, result, "nil")
}

func TestBuildNilCheckConditions_LengthMismatch(t *testing.T) {
	// Test fallback behavior when requiredSegments length doesn't match
	fieldPath := "resource.Spec.V20250312.GroupRef"
	requiredSegments := []bool{true} // Wrong length (should be 3: Spec, V20250312, GroupRef)
	
	stmt := buildNilCheckConditions(fieldPath, requiredSegments)
	result := stmt.GoString()
	
	// Should fall back to checking the last field
	assert.Contains(t, result, "GroupRef")
	assert.Contains(t, result, "!=")
	assert.Contains(t, result, "nil")
}

func TestBuildNilCheckConditions_AllRequired(t *testing.T) {
	// When all segments are required (non-pointer structs), should only check last field
	fieldPath := "resource.Spec.V20250312.GroupRef"
	requiredSegments := []bool{false, true, true} // Spec is never required by convention, V20250312 and GroupRef are required
	
	stmt := buildNilCheckConditions(fieldPath, requiredSegments)
	result := stmt.GoString()
	
	// Since V20250312 and GroupRef are required, and Spec is skipped by special case,
	// should only check the last field
	assert.Contains(t, result, "GroupRef")
	assert.Contains(t, result, "!=")
	assert.Contains(t, result, "nil")
}

func TestBuildNilCheckConditions_MixedRequired(t *testing.T) {
	// When some segments are optional (pointers), should check each optional segment
	fieldPath := "resource.Spec.V20250312.OptionalSection.GroupRef"
	requiredSegments := []bool{false, true, false, true} // Spec (convention), V20250312 required, OptionalSection optional, GroupRef required
	
	stmt := buildNilCheckConditions(fieldPath, requiredSegments)
	result := stmt.GoString()
	
	// V20250312 is required so skipped, OptionalSection is optional so should be checked
	// GroupRef is required but is the final field so it should be in the check
	assert.Contains(t, result, "OptionalSection")
	assert.Contains(t, result, "!=")
	assert.Contains(t, result, "nil")
	
	// Since only OptionalSection is optional (not counting final field), may or may not have &&
	// depending on whether the final required field is also checked
}

func TestBuildNilCheckConditions_AllOptional(t *testing.T) {
	// When all segments are optional, should check each one (except Spec by convention)
	fieldPath := "resource.Spec.V20250312.GroupRef"
	requiredSegments := []bool{false, false, false} // All optional
	
	stmt := buildNilCheckConditions(fieldPath, requiredSegments)
	result := stmt.GoString()
	
	// Should check V20250312 and GroupRef (not Spec due to special case)
	assert.Contains(t, result, "V20250312")
	assert.Contains(t, result, "GroupRef")
	assert.Contains(t, result, "!=")
	assert.Contains(t, result, "nil")
	
	// Should have && operator
	assert.Contains(t, result, "&&")
}

func TestBuildNilCheckConditions_SpecAlwaysSkipped(t *testing.T) {
	// Spec should always be skipped regardless of required status
	fieldPath := "resource.Spec.GroupRef"
	requiredSegments := []bool{false, false} // Both optional, but Spec should still be skipped
	
	stmt := buildNilCheckConditions(fieldPath, requiredSegments)
	result := stmt.GoString()
	
	// Should only check GroupRef, not Spec
	assert.Contains(t, result, "GroupRef")
	assert.Contains(t, result, "!=")
	assert.Contains(t, result, "nil")
	
	// Should not have Spec in the nil check (no &&)
	parts := strings.Split(result, "&&")
	// If only one part, it means no && operator, which is correct
	if len(parts) > 1 {
		// If there is an &&, make sure it's not checking Spec
		for _, part := range parts {
			if strings.Contains(part, "Spec") && strings.Contains(part, "!=") {
				t.Errorf("Spec should not be checked for nil, but found in: %s", part)
			}
		}
	}
}

func TestBuildNilCheckConditions_DeepNesting(t *testing.T) {
	// Test with deeply nested optional fields
	fieldPath := "resource.Spec.V20250312.Config.Advanced.Settings.GroupRef"
	requiredSegments := []bool{false, true, false, false, false, false} // Only V20250312 is required
	
	stmt := buildNilCheckConditions(fieldPath, requiredSegments)
	result := stmt.GoString()
	
	// Should check all optional nested fields
	assert.Contains(t, result, "V20250312")
	assert.Contains(t, result, "Config")
	assert.Contains(t, result, "Advanced")
	assert.Contains(t, result, "Settings")
	assert.Contains(t, result, "GroupRef")
	assert.Contains(t, result, "!=")
	assert.Contains(t, result, "nil")
	
	// Should have multiple && operators
	andCount := strings.Count(result, "&&")
	assert.Greater(t, andCount, 2, "Should have multiple nil checks with && operators")
}

func TestBuildNilCheckConditions_SingleOptionalField(t *testing.T) {
	// Test with only one optional field after Spec
	fieldPath := "resource.Spec.GroupRef"
	requiredSegments := []bool{false, false} // Both optional
	
	stmt := buildNilCheckConditions(fieldPath, requiredSegments)
	result := stmt.GoString()
	
	// Should only check GroupRef (Spec skipped by convention)
	assert.Contains(t, result, "GroupRef")
	assert.Contains(t, result, "!=")
	assert.Contains(t, result, "nil")
	
	// Should not contain && since only one check
	assert.NotContains(t, result, "&&")
}

func TestGenerateFieldExtractionCodeWithNilChecks(t *testing.T) {
	// Integration test: verify that generateFieldExtractionCode uses buildNilCheckConditions correctly
	fields := []ReferenceField{
		{
			FieldName:         "groupRef",
			FieldPath:         "properties.spec.properties.v20250312.properties.groupRef",
			ReferencedKind:    "Group",
			RequiredSegments:  []bool{false, true, false}, // Spec (convention), v20250312 required, groupRef optional
		},
	}
	
	code := generateFieldExtractionCode(fields)
	assert.Len(t, code, 1)
	
	// Can't directly inspect jen.Code, but we can verify it was created
	// The actual behavior will be tested in integration test
	assert.NotNil(t, code[0])
}

func TestGenerateFieldExtractionCode_MultipleReferences(t *testing.T) {
	// Test with multiple references with different required patterns
	fields := []ReferenceField{
		{
			FieldName:         "groupRef",
			FieldPath:         "properties.spec.properties.v20250312.properties.groupRef",
			ReferencedKind:    "Group",
			RequiredSegments:  []bool{false, true, false},
		},
		{
			FieldName:         "secretRef",
			FieldPath:         "properties.spec.properties.v20250312.properties.secretRef",
			ReferencedKind:    "Secret",
			RequiredSegments:  []bool{false, true, true}, // secretRef is required
		},
	}
	
	code := generateFieldExtractionCode(fields)
	assert.Len(t, code, 2)
	
	// Can't directly inspect jen.Code, but we can verify both were created
	assert.NotNil(t, code[0])
	assert.NotNil(t, code[1])
}

func TestGenerateIndexerWithNilChecks_Integration(t *testing.T) {
	// Full integration test: generate an indexer and verify it has proper nil checks
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
                config:
                  properties:
                    groupRef:
                      x-kubernetes-mapping:
                        type:
                          kind: Group
                          group: atlas.generated.mongodb.com
                          version: v1
spec:
  group: atlas.generated.mongodb.com
  names:
    kind: Cluster
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            required:
            - v20250312
            properties:
              v20250312:
                type: object
                properties:
                  config:
                    type: object
                    properties:
                      groupRef:
                        type: object
                        properties:
                          name:
                            type: string
`
	
	tmpDir := t.TempDir()
	testFile := tmpDir + "/test.yaml"
	err := writeFile(testFile, testYAML)
	assert.NoError(t, err)
	
	outputDir := tmpDir + "/indexers"
	err = GenerateIndexers(testFile, "Cluster", outputDir)
	assert.NoError(t, err)
	
	// Read generated file
	content, err := readFile(outputDir + "/clusterbygroup.go")
	assert.NoError(t, err)
	
	contentStr := string(content)
	
	// Verify proper nil checks are generated
	// Should check V20250312, Config, and GroupRef
	assert.Contains(t, contentStr, "V20250312")
	assert.Contains(t, contentStr, "Config")
	assert.Contains(t, contentStr, "GroupRef")
	assert.Contains(t, contentStr, "!= nil")
	
	// Should have multiple nil checks
	nilCheckCount := strings.Count(contentStr, "!= nil")
	assert.GreaterOrEqual(t, nilCheckCount, 2, "Should have multiple nil checks for nested optional fields")
}

// Helper functions for tests
func writeFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func readFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

