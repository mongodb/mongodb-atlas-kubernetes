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

func TestFromConfig_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

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
`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(testFile, []byte(testYAML), 0644)
	require.NoError(t, err)

	controllerDir := filepath.Join(tmpDir, "controllers")
	indexerDir := filepath.Join(tmpDir, "indexers")

	err = FromConfig(testFile, "Cluster", controllerDir, indexerDir, "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1", true)
	require.NoError(t, err)

	clusterControllerDir := filepath.Join(controllerDir, "cluster")
	assert.DirExists(t, clusterControllerDir)

	controllerFile := filepath.Join(clusterControllerDir, "cluster_controller.go")
	assert.FileExists(t, controllerFile)

	handlerFile := filepath.Join(clusterControllerDir, "handler.go")
	assert.FileExists(t, handlerFile)

	versionHandlerFile := filepath.Join(clusterControllerDir, "handler_v20250312.go")
	assert.FileExists(t, versionHandlerFile)

	indexerFile := filepath.Join(indexerDir, "clusterbygroup.go")
	assert.FileExists(t, indexerFile)
}

func TestGenerateSetupWithManager_Watches(t *testing.T) {
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
`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(testFile, []byte(testYAML), 0644)
	require.NoError(t, err)

	controllerDir := filepath.Join(tmpDir, "controllers")
	indexerDir := filepath.Join(tmpDir, "indexers")

	err = FromConfig(testFile, "Cluster", controllerDir, indexerDir, "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1", true)
	require.NoError(t, err)

	handlerFile := filepath.Join(controllerDir, "cluster", "handler.go")
	content, err := os.ReadFile(handlerFile)
	require.NoError(t, err)

	contentStr := string(content)

	assert.Contains(t, contentStr, "func (h *Handler) SetupWithManager")

	assert.Contains(t, contentStr, "Watches")
	assert.Contains(t, contentStr, "&akov2generated.Group{}", "Should contain Group reference")

	assert.Contains(t, contentStr, "NewClusterByGroupMapFunc")

	assert.Contains(t, contentStr, "ResourceVersionChangedPredicate")
}

func TestGenerateMapperFunctions_MultipleReferences(t *testing.T) {
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
              x-atlas-sdk-version: go.mongodb.org/atlas-sdk/v20250312008/admin
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

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(testFile, []byte(testYAML), 0644)
	require.NoError(t, err)

	controllerDir := filepath.Join(tmpDir, "controllers")
	indexerDir := filepath.Join(tmpDir, "indexers")

	err = FromConfig(testFile, "Integration", controllerDir, indexerDir, "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1", true)
	require.NoError(t, err)

	handlerFile := filepath.Join(controllerDir, "integration", "handler.go")
	content, err := os.ReadFile(handlerFile)
	require.NoError(t, err)

	contentStr := string(content)

	// Check for Group and Secret references (may be aliased due to import conflicts)
	assert.Contains(t, contentStr, ".Group{}")
	assert.Contains(t, contentStr, ".Secret{}")

	assert.Contains(t, contentStr, "NewIntegrationBySecretMapFunc")
}

func TestGetWatchedTypeInstance(t *testing.T) {
	tests := []struct {
		kind     string
		expected string
	}{
		{"Secret", "&corev1.Secret{}"},
		{"Group", "&v1.Group{}"},
		{"CustomResource", "&v1.CustomResource{}"},
	}

	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			stmt := getWatchedTypeInstance(tt.kind)
			result := stmt.GoString()
			assert.Contains(t, result, tt.kind)
		})
	}
}

func TestGenerateController_NoReferences(t *testing.T) {
	testYAML := `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: teams.atlas.generated.mongodb.com
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

	controllerDir := filepath.Join(tmpDir, "controllers")
	indexerDir := filepath.Join(tmpDir, "indexers")

	err = FromConfig(testFile, "Team", controllerDir, indexerDir, "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1", true)
	require.NoError(t, err)

	handlerFile := filepath.Join(controllerDir, "team", "handler.go")
	content, err := os.ReadFile(handlerFile)
	require.NoError(t, err)

	contentStr := string(content)

	assert.Contains(t, contentStr, "func (h *Handler) SetupWithManager")

	// Should not have Watches() calls because there are no refs
	assert.NotContains(t, contentStr, ".Watches(")
}

func TestGeneratedControllerStructure(t *testing.T) {
	testYAML := `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: resources.atlas.generated.mongodb.com
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
    kind: Resource
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        properties:
          spec:
            properties:
              v20250312:
                type: object
`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(testFile, []byte(testYAML), 0644)
	require.NoError(t, err)

	controllerDir := filepath.Join(tmpDir, "controllers")
	indexerDir := filepath.Join(tmpDir, "indexers")

	err = FromConfig(testFile, "Resource", controllerDir, indexerDir, "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1", true)
	require.NoError(t, err)

	controllerFile := filepath.Join(controllerDir, "resource", "resource_controller.go")
	content, err := os.ReadFile(controllerFile)
	require.NoError(t, err)

	contentStr := string(content)

	assert.Contains(t, contentStr, "package resource")
	assert.Contains(t, contentStr, "type Handler struct")
	assert.Contains(t, contentStr, "+kubebuilder:rbac:groups=atlas.generated.mongodb.com,resources=resources")
	assert.Contains(t, contentStr, "func NewResourceReconciler")
}
