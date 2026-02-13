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

package atlasexporters

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/scaffolder/internal/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/scaffolder/internal/generators/indexers"
)

func assertContains(t *testing.T, content, needle string) {
	t.Helper()

	if !strings.Contains(content, needle) {
		t.Fatalf("expected to find %q in content:\n%s", needle, content)
	}
}

func TestGenerateResourceExporter_UsesGetBlock(t *testing.T) {
	dir := t.TempDir()

	mapping := config.MappingWithConfig{
		Version: "v1",
		OpenAPIConfig: config.OpenAPIConfig{
			Package: "go.mongodb.org/atlas-sdk/v20250312011/admin",
		},
	}

	// No reference fields means it should use the get block
	var referenceFields []indexers.ReferenceField

	err := GenerateExporter(
		"Group",
		"Group",
		"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1",
		dir,
		mapping,
		referenceFields,
	)
	if err != nil {
		t.Fatalf("GenerateExporter: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(dir, "group_exporter.go"))
	if err != nil {
		t.Fatalf("read generated file: %v", err)
	}

	output := string(content)
	assertContains(t, output, "type GroupExporter struct")
	assertContains(t, output, "Export(ctx context.Context, referencedObjects []client.Object)")
	assertContains(t, output, "GetGroup")
	assertContains(t, output, "FromAPI(resource, atlasResource, referencedObjects...)")
	assertContains(t, output, "resource.GetObjectKind().SetGroupVersionKind")
	assertContains(t, output, `GroupVersion.WithKind("Group")`)
	assertContains(t, output, "append([]client.Object{resource}, resources...)")
	if strings.Contains(output, "AllPages(") {
		t.Fatalf("expected get block output, found list block")
	}
}

func TestGenerateResourceExporter_UsesListBlock(t *testing.T) {
	dir := t.TempDir()

	mapping := config.MappingWithConfig{
		Version: "v1",
		OpenAPIConfig: config.OpenAPIConfig{
			Package: "go.mongodb.org/atlas-sdk/v20250312011/admin",
		},
	}

	// With reference fields, it should use the list block
	referenceFields := []indexers.ReferenceField{
		{
			FieldName:      "groupRef",
			ReferencedKind: "Group",
		},
	}

	err := GenerateExporter(
		"Cluster",
		"Cluster",
		"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1",
		dir,
		mapping,
		referenceFields,
	)
	if err != nil {
		t.Fatalf("GenerateExporter: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(dir, "cluster_exporter.go"))
	if err != nil {
		t.Fatalf("read generated file: %v", err)
	}

	output := string(content)
	assertContains(t, output, "Export(ctx context.Context, referencedObjects []client.Object)")
	assertContains(t, output, "var atlasResources []any")
	assertContains(t, output, "for pageNum := 1; ; pageNum++")
	assertContains(t, output, "resp.GetResults()")
	assertContains(t, output, "resp.GetTotalCount()")
	assertContains(t, output, "ListClusters(")
	assertContains(t, output, "e.identifiers[0]")
	assertContains(t, output, "FromAPI(resource, atlasResource, referencedObjects...)")
	assertContains(t, output, "resource.GetObjectKind().SetGroupVersionKind")
	assertContains(t, output, `GroupVersion.WithKind("Cluster")`)
	assertContains(t, output, "resources = append(resources, resource)")
	assertContains(t, output, "resources = append(resources, translatedResources...)")
}

func TestGenerateResourceExporter_InvalidDestination(t *testing.T) {
	// Use a path that cannot be created
	invalidPath := "/nonexistent/path/that/cannot/be/created/\x00invalid"

	mapping := config.MappingWithConfig{
		Version: "v1",
		OpenAPIConfig: config.OpenAPIConfig{
			Package: "example.com/sdk",
		},
	}

	err := GenerateExporter(
		"Foo",
		"Foo",
		"github.com/acme/resource",
		invalidPath,
		mapping,
		nil,
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
