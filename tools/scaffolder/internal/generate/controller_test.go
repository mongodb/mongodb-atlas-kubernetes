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

	assert.Contains(t, contentStr, "func (h *ClusterHandler) SetupWithManager")

	assert.Contains(t, contentStr, "Watches")
	// Group conflicts with apiextensions, sometimes linted as v11
	assert.True(t,
		strings.Contains(contentStr, "&v1.Group{}") || strings.Contains(contentStr, "&v11.Group{}"),
		"Should contain Group reference")

	assert.Contains(t, contentStr, "func (h *ClusterHandler) clusterForGroupMapFunc()")
	assert.Contains(t, contentStr, "ProjectsIndexMapperFunc")

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

	assert.Contains(t, contentStr, "integrationForGroupMapFunc")
	assert.Contains(t, contentStr, "integrationForSecretMapFunc")

	// Check for Group and Secret references (may be aliased due to import conflicts)
	assert.Contains(t, contentStr, ".Group{}")
	assert.Contains(t, contentStr, ".Secret{}")

	assert.Contains(t, contentStr, "ProjectsIndexMapperFunc")
	assert.Contains(t, contentStr, "CredentialsIndexMapperFunc")
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

	assert.Contains(t, contentStr, "func (h *TeamHandler) SetupWithManager")

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
	assert.Contains(t, contentStr, "type ResourceHandler struct")
	assert.Contains(t, contentStr, "+kubebuilder:rbac:groups=atlas.generated.mongodb.com,resources=resources")
	assert.Contains(t, contentStr, "func NewResourceReconciler")
}

func TestGeneratedHandlerDelegation(t *testing.T) {
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
            v20250401:
              x-atlas-sdk-version: go.mongodb.org/atlas-sdk/v20250401001/admin
spec:
  group: atlas.generated.mongodb.com
  names:
    kind: Resource
    plural: resources
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        properties:
          spec:
            properties:
              v20250312:
                type: object
              v20250401:
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

	handlerFile := filepath.Join(controllerDir, "resource", "handler.go")
	content, err := os.ReadFile(handlerFile)
	require.NoError(t, err)

	contentStr := string(content)

	assert.Contains(t, contentStr, "func (h *ResourceHandler) getHandlerForResource")
	assert.Contains(t, contentStr, "func (h *ResourceHandler) HandleInitial")
	assert.Contains(t, contentStr, "func (h *ResourceHandler) HandleCreating")
	assert.Contains(t, contentStr, "func (h *ResourceHandler) HandleDeletionRequested")
	v1Handler := filepath.Join(controllerDir, "resource", "handler_v20250312.go")
	assert.FileExists(t, v1Handler)

	v2Handler := filepath.Join(controllerDir, "resource", "handler_v20250401.go")
	assert.FileExists(t, v2Handler)
}

func TestGeneratedHelperFunctions(t *testing.T) {
	testYAML := `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: groups.atlas.generated.mongodb.com
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
    kind: Group
    plural: groups
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

	err = FromConfig(testFile, "Group", controllerDir, indexerDir, "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1", true)
	require.NoError(t, err)

	// Test handler.go contains package-level getTranslationRequest function
	handlerFile := filepath.Join(controllerDir, "group", "handler.go")
	content, err := os.ReadFile(handlerFile)
	require.NoError(t, err)
	contentStr := string(content)

	// Verify package-level getTranslationRequest function
	assert.Contains(t, contentStr, "func getTranslationRequest(")
	assert.Contains(t, contentStr, "ctx context.Context")
	//assert.Contains(t, contentStr, "client client.Client")
	assert.Contains(t, contentStr, "crdName string")
	assert.Contains(t, contentStr, "storageVersion string")
	assert.Contains(t, contentStr, "targetVersion string")
	assert.Contains(t, contentStr, "NewTranslator")

	// Test versioned handler contains helper methods
	versionHandlerFile := filepath.Join(controllerDir, "group", "handler_v20250312.go")
	versionContent, err := os.ReadFile(versionHandlerFile)
	require.NoError(t, err)
	versionContentStr := string(versionContent)

	// Verify getSDKClientSet method
	assert.Contains(t, versionContentStr, "func (h *GroupHandlerv20250312) getSDKClientSet(")
	assert.Contains(t, versionContentStr, "GetConnectionConfig")
	assert.Contains(t, versionContentStr, "SdkClientSet")
	assert.Contains(t, versionContentStr, "ConnectionSecretRef")

	// Verify getTranslationRequest wrapper method
	assert.Contains(t, versionContentStr, "func (h *GroupHandlerv20250312) getTranslationRequest(")
	assert.Contains(t, versionContentStr, "return getTranslationRequest(")
	assert.Contains(t, versionContentStr, "groups.atlas.generated.mongodb.com")
	assert.Contains(t, versionContentStr, "\"v1\"")        // storage version
	assert.Contains(t, versionContentStr, "\"v20250312\"") // target version
}
