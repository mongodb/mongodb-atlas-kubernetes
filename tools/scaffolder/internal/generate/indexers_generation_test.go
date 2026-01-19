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

func TestGenerateIndexers_Integration(t *testing.T) {
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

	outputDir := filepath.Join(tmpDir, "indexers")

	t.Run("GenerateIndexerFiles", func(t *testing.T) {
		err := GenerateIndexers(testFile, "Cluster", outputDir)
		require.NoError(t, err)

		indexerFile := filepath.Join(outputDir, "clusterbygroup.go")
		assert.FileExists(t, indexerFile)

		content, err := os.ReadFile(indexerFile)
		require.NoError(t, err)

		contentStr := string(content)
		assert.Contains(t, contentStr, "package indexer")
		assert.Contains(t, contentStr, "type ClusterByGroupIndexer struct")
		assert.Contains(t, contentStr, "const ClusterByGroupIndex")
		assert.Contains(t, contentStr, "func NewClusterByGroupIndexer")
		assert.Contains(t, contentStr, "func (*ClusterByGroupIndexer) Object()")
		assert.Contains(t, contentStr, "func (*ClusterByGroupIndexer) Name()")
		assert.Contains(t, contentStr, "func (i *ClusterByGroupIndexer) Keys(")
		assert.Contains(t, contentStr, "func NewClusterByGroupMapFunc")
		assert.Contains(t, contentStr, `"k8s.io/apimachinery/pkg/types"`)
		assert.Contains(t, contentStr, `"sigs.k8s.io/controller-runtime/pkg/reconcile"`)
	})

	t.Run("GenerateSingleLevelArrayIndexers", func(t *testing.T) {
		arrayYAML := `
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
`
		arrayFile := filepath.Join(tmpDir, "array.yaml")
		err := os.WriteFile(arrayFile, []byte(arrayYAML), 0644)
		require.NoError(t, err)

		arrayOutputDir := filepath.Join(tmpDir, "array-indexers")
		err = GenerateIndexers(arrayFile, "AlertConfig", arrayOutputDir)
		require.NoError(t, err)

		// Single-level arrays should now generate indexers
		files, err := os.ReadDir(arrayOutputDir)
		require.NoError(t, err)
		assert.NotEmpty(t, files, "Indexer files should be generated for single-level arrays")

		// Verify the generated indexer has loop code
		indexerFile := filepath.Join(arrayOutputDir, "alertconfigbysecret.go")
		content, err := os.ReadFile(indexerFile)
		require.NoError(t, err)
		assert.Contains(t, string(content), "for _, ", "Should have for-loop for array iteration")
	})
}

func TestGenerateIndexers_NoReferences(t *testing.T) {
	testYAML := `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: teams.atlas.generated.mongodb.com
spec:
  group: atlas.generated.mongodb.com
  names:
    kind: Team
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        properties:
          spec:
            properties:
              v20250312:
                properties:
                  name:
                    type: string
`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(testFile, []byte(testYAML), 0644)
	require.NoError(t, err)

	outputDir := filepath.Join(tmpDir, "indexers")

	err = GenerateIndexers(testFile, "Team", outputDir)
	require.NoError(t, err)

	files, err := os.ReadDir(outputDir)
	if err == nil {
		assert.Empty(t, files)
	}
}

func TestGenerateRequestsFunction_UniqueNames(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "indexers")
	err := os.MkdirAll(outputDir, 0755)
	require.NoError(t, err)

	testYAML := `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: integrations.atlas.generated.mongodb.com
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
                apiKeyRef:
                  x-kubernetes-mapping:
                    type:
                      kind: Secret
spec:
  group: atlas.generated.mongodb.com
  names:
    kind: Integration
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
                  apiKeyRef:
                    type: object
`

	testFile := filepath.Join(tmpDir, "test.yaml")
	err = os.WriteFile(testFile, []byte(testYAML), 0644)
	require.NoError(t, err)

	err = GenerateIndexers(testFile, "Integration", outputDir)
	require.NoError(t, err)

	groupFile, err := os.ReadFile(filepath.Join(outputDir, "integrationbygroup.go"))
	require.NoError(t, err)

	secretFile, err := os.ReadFile(filepath.Join(outputDir, "integrationbysecret.go"))
	require.NoError(t, err)

	assert.Contains(t, string(groupFile), "NewIntegrationByGroupMapFunc")
	assert.Contains(t, string(secretFile), "NewIntegrationBySecretMapFunc")
	assert.NotContains(t, string(groupFile), "func IntegrationRequests(")
	assert.NotContains(t, string(secretFile), "func IntegrationRequests(")
}
