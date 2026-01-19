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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseReferenceFields(t *testing.T) {
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
        properties:
          spec:
            properties:
              v20250312:
                properties:
                  groupRef:
                    type: object
                    properties:
                      name:
                        type: string
`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(testFile, []byte(testYAML), 0644)
	require.NoError(t, err)

	t.Run("ParseGroupReference", func(t *testing.T) {
		refs, err := ParseReferenceFields(testFile, "Cluster")
		require.NoError(t, err)
		require.Len(t, refs, 1)

		ref := refs[0]
		assert.Equal(t, "groupRef", ref.FieldName)
		assert.Equal(t, "Group", ref.ReferencedKind)
		assert.Contains(t, ref.FieldPath, "groupRef")
	})

	t.Run("ParseNonExistentCRD", func(t *testing.T) {
		refs, err := ParseReferenceFields(testFile, "NonExistent")
		assert.Error(t, err)
		assert.Nil(t, refs)
	})
}

func TestParseReferenceFields_ArrayReferences(t *testing.T) {
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
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        properties:
          spec:
            properties:
              v20250312:
                properties:
                  notifications:
                    type: array
                    items:
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

	refs, err := ParseReferenceFields(testFile, "AlertConfig")
	require.NoError(t, err)
	require.Len(t, refs, 1)

	// Array references are now supported (single-level)
	assert.Contains(t, refs[0].FieldPath, ".items.")
	assert.True(t, refs[0].IsArrayBased(), "Should be marked as array-based")
	assert.Equal(t, "secretRef", refs[0].FieldName)
}

func TestParseReferenceFields_RequiredSegments(t *testing.T) {
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
        required:
        - spec
        properties:
          spec:
            type: object
            required:
            - v20250312
            properties:
              v20250312:
                type: object
                required:
                - groupRef
                properties:
                  groupRef:
                    type: object
                    properties:
                      name:
                        type: string
`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(testFile, []byte(testYAML), 0644)
	require.NoError(t, err)

	refs, err := ParseReferenceFields(testFile, "Cluster")
	require.NoError(t, err)
	require.Len(t, refs, 1)

	ref := refs[0]
	assert.Equal(t, "groupRef", ref.FieldName)
	assert.Equal(t, "Group", ref.ReferencedKind)
	assert.Equal(t, "properties.spec.properties.v20250312.properties.groupRef", ref.FieldPath)
	// spec is never required (Kubernetes convention), v20250312 is required in spec, groupRef is required in v20250312
	assert.Equal(t, []bool{false, true, true}, ref.RequiredSegments)
}

func TestParseReferenceFields_MixedRequiredSegments(t *testing.T) {
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
                optionalSection:
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
        required:
        - spec
        properties:
          spec:
            type: object
            required:
            - v20250312
            properties:
              v20250312:
                type: object
                properties:
                  optionalSection:
                    type: object
                    required:
                    - groupRef
                    properties:
                      groupRef:
                        type: object
                        properties:
                          name:
                            type: string
`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(testFile, []byte(testYAML), 0644)
	require.NoError(t, err)

	refs, err := ParseReferenceFields(testFile, "Cluster")
	require.NoError(t, err)
	require.Len(t, refs, 1)

	ref := refs[0]
	assert.Equal(t, "groupRef", ref.FieldName)
	// spec is never required (Kubernetes convention), v20250312 is required, optionalSection is NOT required, groupRef IS required
	assert.Equal(t, []bool{false, true, false, true}, ref.RequiredSegments)
}
