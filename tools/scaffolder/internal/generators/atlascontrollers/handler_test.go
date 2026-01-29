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

package atlascontrollers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	exporterDir := filepath.Join(tmpDir, "exporters")

	typesPath := "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	indexerImportPath := "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/indexers"
	err = runAllGenerators(testFile, "Resource", controllerDir, indexerDir, exporterDir, typesPath, typesPath, indexerImportPath, true)
	require.NoError(t, err)

	handlerFile := filepath.Join(controllerDir, "resource", "handler.go")
	content, err := os.ReadFile(handlerFile)
	require.NoError(t, err)

	contentStr := string(content)

	assert.Contains(t, contentStr, "func (h *Handler) getHandlerForResource")
	assert.Contains(t, contentStr, "func (h *Handler) HandleInitial")
	assert.Contains(t, contentStr, "func (h *Handler) HandleCreating")
	assert.Contains(t, contentStr, "func (h *Handler) HandleDeletionRequested")
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
	exporterDir := filepath.Join(tmpDir, "exporters")

	typesPath := "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	indexerImportPath := "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/indexers"
	err = runAllGenerators(testFile, "Group", controllerDir, indexerDir, exporterDir, typesPath, typesPath, indexerImportPath, true)
	require.NoError(t, err)

	// Test handler.go contains certain package-level functions
	handlerFile := filepath.Join(controllerDir, "group", "handler.go")
	content, err := os.ReadFile(handlerFile)
	require.NoError(t, err)
	contentStr := string(content)

	// Verify package-level getHandlerForResource method
	assert.Contains(t, contentStr, "func (h *Handler) getHandlerForResource(")
	assert.Contains(t, contentStr, "ctx context.Context")
	assert.Contains(t, contentStr, "h.getSDKClientSet")
	assert.Contains(t, contentStr, "h.translators")

	// Verify getSDKClientSet method
	assert.Contains(t, contentStr, "func (h *Handler) getSDKClientSet(")
	assert.Contains(t, contentStr, "GetConnectionConfig")
	assert.Contains(t, contentStr, "SdkClientSet")
	assert.Contains(t, contentStr, "ConnectionSecretRef")
}
