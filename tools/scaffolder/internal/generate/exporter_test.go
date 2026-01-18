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
)

func assertContains(t *testing.T, content, needle string) {
	t.Helper()

	if !strings.Contains(content, needle) {
		t.Fatalf("expected to find %q in content:\n%s", needle, content)
	}
}

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write file %s: %v", path, err)
	}
	return path
}

func TestGenerateResourceExporter_UsesGetBlock(t *testing.T) {
	dir := t.TempDir()
	crdPath := writeTempFile(t, dir, "crd.yaml", minimalCRD("Group", false))

	req := &exporterRequest{
		crdPath:            crdPath,
		kind:               "Group",
		resourceName:       "Group",
		resourceImportPath: "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1",
		destination:        dir,
		mapping: MappingWithConfig{
			Version: "v1",
			OpenAPIConfig: OpenAPIConfig{
				Package: "go.mongodb.org/atlas-sdk/v20250312011/admin",
			},
		},
	}

	if err := generateResourceExporter(req); err != nil {
		t.Fatalf("generateResourceExporter: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(dir, "group_exporter.go"))
	if err != nil {
		t.Fatalf("read generated file: %v", err)
	}

	output := string(content)
	assertContains(t, output, "type GroupExporter struct")
	assertContains(t, output, "RESORCE_GET_METHOD")
	if strings.Contains(output, "AllPages(") {
		t.Fatalf("expected get block output, found list block")
	}
}

func TestGenerateResourceExporter_UsesListBlock(t *testing.T) {
	dir := t.TempDir()
	crdPath := writeTempFile(t, dir, "crd.yaml", minimalCRD("Cluster", true))

	req := &exporterRequest{
		crdPath:            crdPath,
		kind:               "Cluster",
		resourceName:       "Cluster",
		resourceImportPath: "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1",
		destination:        dir,
		mapping: MappingWithConfig{
			Version: "v1",
			OpenAPIConfig: OpenAPIConfig{
				Package: "go.mongodb.org/atlas-sdk/v20250312011/admin",
			},
		},
	}

	if err := generateResourceExporter(req); err != nil {
		t.Fatalf("generateResourceExporter: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(dir, "cluster_exporter.go"))
	if err != nil {
		t.Fatalf("read generated file: %v", err)
	}

	output := string(content)
	assertContains(t, output, "AllPages(")
	assertContains(t, output, "List_RESOURCE(")
	assertContains(t, output, "e.identifiers[0]")
}

func TestGenerateResourceExporter_ParseError(t *testing.T) {
	dir := t.TempDir()
	req := &exporterRequest{
		crdPath:            filepath.Join(dir, "missing.yaml"),
		kind:               "Foo",
		resourceName:       "Foo",
		resourceImportPath: "github.com/acme/resource",
		destination:        dir,
		mapping: MappingWithConfig{
			Version: "v1",
			OpenAPIConfig: OpenAPIConfig{
				Package: "example.com/sdk",
			},
		},
	}

	if err := generateResourceExporter(req); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func minimalCRD(kind string, withRefs bool) string {
	if withRefs {
		return `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: ` + strings.ToLower(kind) + `s.example.com
  annotations:
    api-mappings: |
      properties:
        spec:
          properties:
            refField:
              x-kubernetes-mapping:
                type:
                  kind: Group
                  group: atlas.generated.mongodb.com
                  version: v1
spec:
  group: atlas.generated.mongodb.com
  names:
    kind: ` + kind + `
    plural: ` + kind + `s
  scope: Namespaced
  versions:
  - name: v1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
`
	}

	return `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: ` + kind + `s.example.com
spec:
  group: atlas.generated.mongodb.com
  names:
    kind: ` + kind + `
    plural: ` + kind + `s
  scope: Namespaced
  versions:
  - name: v1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
`
}
