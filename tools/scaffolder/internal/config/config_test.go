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

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCRDConfig(t *testing.T) {
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
              x-atlas-sdk-version: go.mongodb.org/atlas-sdk/v20250312008/admin
spec:
  group: atlas.generated.mongodb.com
  names:
    kind: Cluster
    plural: clusters
  versions:
  - name: v1
`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(testFile, []byte(testYAML), 0644)
	require.NoError(t, err)

	t.Run("ParseValidCRD", func(t *testing.T) {
		config, err := ParseCRDConfig(testFile, "Cluster")
		require.NoError(t, err)
		assert.Equal(t, "Cluster", config.ResourceName)
		assert.Len(t, config.Mappings, 1)
		assert.Equal(t, "v20250312", config.Mappings[0].Version)
		assert.Equal(t, "go.mongodb.org/atlas-sdk/v20250312008/admin", config.Mappings[0].OpenAPIConfig.Package)
	})

	t.Run("ParseNonExistentCRD", func(t *testing.T) {
		_, err := ParseCRDConfig(testFile, "NonExistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("ParseInvalidFile", func(t *testing.T) {
		_, err := ParseCRDConfig("/nonexistent/file.yaml", "Cluster")
		assert.Error(t, err)
	})
}

func TestListCRDs(t *testing.T) {
	testYAML := `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: clusters.atlas.generated.mongodb.com
spec:
  group: atlas.generated.mongodb.com
  names:
    kind: Cluster
    plural: clusters
    categories: [atlas]
  versions:
  - name: v1
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: groups.atlas.generated.mongodb.com
spec:
  group: atlas.generated.mongodb.com
  names:
    kind: Group
    plural: groups
  versions:
  - name: v1
`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(testFile, []byte(testYAML), 0644)
	require.NoError(t, err)

	crds, err := ListCRDs(testFile)
	require.NoError(t, err)
	assert.Len(t, crds, 2)

	assert.Equal(t, "Cluster", crds[0].Kind)
	assert.Equal(t, "atlas.generated.mongodb.com", crds[0].Group)
	assert.Equal(t, "v1", crds[0].Version)
	assert.Contains(t, crds[0].Categories, "atlas")

	assert.Equal(t, "Group", crds[1].Kind)
}
