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

package indexers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDependentReferences(t *testing.T) {
	// Create a test YAML with multiple CRDs where some reference Group
	testYAML := `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: groups.atlas.generated.mongodb.com
spec:
  group: atlas.generated.mongodb.com
  names:
    kind: Group
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
---
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
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: databaseusers.atlas.generated.mongodb.com
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
    kind: DatabaseUser
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
---
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

	t.Run("FindDependentsOfGroup", func(t *testing.T) {
		dependents, err := ParseDependentReferences(testFile, "Group")
		require.NoError(t, err)

		// Should find Cluster and DatabaseUser as dependents of Group
		assert.Len(t, dependents, 2)

		// Check that the dependent info is correct
		dependentKinds := make(map[string]DependentInfo)
		for _, dep := range dependents {
			dependentKinds[dep.DependentKind] = dep
		}

		// Verify Cluster dependent
		clusterDep, ok := dependentKinds["Cluster"]
		assert.True(t, ok, "Cluster should be a dependent of Group")
		assert.Equal(t, "Group", clusterDep.TargetKind)
		assert.Equal(t, "ClusterByGroupIndex", clusterDep.IndexerConstantName)
		assert.Equal(t, "NewClusterByGroupMapFunc", clusterDep.MapFuncName)

		// Verify DatabaseUser dependent
		dbUserDep, ok := dependentKinds["DatabaseUser"]
		assert.True(t, ok, "DatabaseUser should be a dependent of Group")
		assert.Equal(t, "Group", dbUserDep.TargetKind)
		assert.Equal(t, "DatabaseUserByGroupIndex", dbUserDep.IndexerConstantName)
		assert.Equal(t, "NewDatabaseUserByGroupMapFunc", dbUserDep.MapFuncName)
	})

	t.Run("NoDependentsForTeam", func(t *testing.T) {
		dependents, err := ParseDependentReferences(testFile, "Team")
		require.NoError(t, err)

		// Team has no dependents - no one references it
		assert.Empty(t, dependents)
	})

	t.Run("DependentDoesNotIncludeSelf", func(t *testing.T) {
		dependents, err := ParseDependentReferences(testFile, "Cluster")
		require.NoError(t, err)

		// Cluster should not have itself as a dependent
		for _, dep := range dependents {
			assert.NotEqual(t, "Cluster", dep.DependentKind)
		}
	})
}

func TestParseDependentReferences_ArrayBasedReferences(t *testing.T) {
	// Test that ParseDependentReferences correctly finds dependents even when
	// the reference is inside an array (e.g., spec.entries[].groupRef)
	testYAML := `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: groups.atlas.generated.mongodb.com
spec:
  group: atlas.generated.mongodb.com
  names:
    kind: Group
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
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: alertconfigs.atlas.generated.mongodb.com
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
                      groupRef:
                        x-kubernetes-mapping:
                          type:
                            kind: Group
                            group: atlas.generated.mongodb.com
                            version: v1
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
                      type: object
                      properties:
                        groupRef:
                          type: object
                          properties:
                            name:
                              type: string
---
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
                regions:
                  items:
                    properties:
                      configs:
                        items:
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
                  regions:
                    type: array
                    items:
                      type: object
                      properties:
                        configs:
                          type: array
                          items:
                            type: object
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

	t.Run("FindDependentsWithSingleLevelArrayReference", func(t *testing.T) {
		dependents, err := ParseDependentReferences(testFile, "Group")
		require.NoError(t, err)

		// Should find AlertConfig and Integration as dependents of Group
		// even though the references are inside arrays
		dependentKinds := make(map[string]DependentInfo)
		for _, dep := range dependents {
			dependentKinds[dep.DependentKind] = dep
		}

		// Verify AlertConfig dependent (single-level array: notifications[].groupRef)
		alertDep, ok := dependentKinds["AlertConfig"]
		assert.True(t, ok, "AlertConfig should be a dependent of Group (array-based reference)")
		assert.Equal(t, "Group", alertDep.TargetKind)
		assert.Equal(t, "AlertConfigByGroupIndex", alertDep.IndexerConstantName)
		assert.Equal(t, "NewAlertConfigByGroupMapFunc", alertDep.MapFuncName)
	})

	t.Run("FindDependentsWithNestedArrayReference", func(t *testing.T) {
		dependents, err := ParseDependentReferences(testFile, "Group")
		require.NoError(t, err)

		dependentKinds := make(map[string]DependentInfo)
		for _, dep := range dependents {
			dependentKinds[dep.DependentKind] = dep
		}

		// Verify Integration dependent (nested array: regions[].configs[].groupRef)
		integrationDep, ok := dependentKinds["Integration"]
		assert.True(t, ok, "Integration should be a dependent of Group (nested array reference)")
		assert.Equal(t, "Group", integrationDep.TargetKind)
		assert.Equal(t, "IntegrationByGroupIndex", integrationDep.IndexerConstantName)
		assert.Equal(t, "NewIntegrationByGroupMapFunc", integrationDep.MapFuncName)
	})
}
