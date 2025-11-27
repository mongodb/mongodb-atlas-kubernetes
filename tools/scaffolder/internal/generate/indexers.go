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
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dave/jennifer/jen"
	"gopkg.in/yaml.v3"
)

type ReferenceField struct {
	FieldName         string
	FieldPath         string
	ReferencedKind    string
	ReferencedGroup   string
	ReferencedVersion string
	IndexerType       string
}

type IndexerInfo struct {
	ResourceName    string
	IndexerName     string
	TargetKind      string
	ConstantName    string
	FunctionName    string
	ReferenceFields []ReferenceField // All fields that reference this kind
}

func ParseReferenceFields(resultPath, crdKind string) ([]ReferenceField, error) {
	file, err := os.Open(resultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open result file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for {
		refs, err := parseNextCRDReferences(scanner, crdKind)
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}

		// We found the target CRD, return the references (even if empty)
		return refs, nil
	}

	return nil, fmt.Errorf("CRD kind '%s' not found in result file", crdKind)
}

func parseNextCRDReferences(scanner *bufio.Scanner, targetKind string) ([]ReferenceField, error) {
	var buffer bytes.Buffer

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "---" {
			if len(strings.TrimSpace(buffer.String())) > 0 {
				refs, err := extractReferencesFromCRD(buffer.Bytes(), targetKind)
				if err != nil {
					buffer.Reset()
					continue
				}
				return refs, nil
			}
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		buffer.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	if buffer.Len() > 0 {
		refs, err := extractReferencesFromCRD(buffer.Bytes(), targetKind)
		if err != nil {
			// Ignore "not target CRD" errors, just like in the loop above
			if err.Error() != "not target CRD" {
				return nil, err
			}
		} else {
			return refs, nil
		}
	}

	return nil, io.EOF
}

func extractReferencesFromCRD(content []byte, targetKind string) ([]ReferenceField, error) {
	var crd CRDDocument
	if err := yaml.Unmarshal(content, &crd); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	if crd.Kind != "CustomResourceDefinition" || crd.Spec.Names.Kind != targetKind {
		return nil, fmt.Errorf("not target CRD")
	}

	var references []ReferenceField

	if apiMappingsStr, exists := crd.Metadata.Annotations["api-mappings"]; exists {
		refs, err := parseAPIMapping(apiMappingsStr)
		if err == nil {
			references = append(references, refs...)
		}
	}

	// Return empty slice for CRDs with no references (this is valid)
	return references, nil
}

type VersionProperties struct {
	Properties map[string]PropertyMapping `yaml:"properties"`
}

type PropertyMapping struct {
	KubernetesMapping *KubernetesMapping `yaml:"x-kubernetes-mapping,omitempty"`
	OpenAPIMapping    *OpenAPIMapping    `yaml:"x-openapi-mapping,omitempty"`
}

type KubernetesMapping struct {
	NameSelector string            `yaml:"nameSelector"`
	Properties   []string          `yaml:"properties"`
	Type         KubernetesRefType `yaml:"type"`
}

type KubernetesRefType struct {
	Group    string `yaml:"group"`
	Kind     string `yaml:"kind"`
	Resource string `yaml:"resource"`
	Version  string `yaml:"version"`
}

type OpenAPIMapping struct {
	Property string `yaml:"property"`
	Type     string `yaml:"type,omitempty"`
}

func parseAPIMapping(apiMappingsStr string) ([]ReferenceField, error) {
	// Parse the api-mappings as a generic map to handle any nesting level
	var mapping map[string]any
	if err := yaml.Unmarshal([]byte(apiMappingsStr), &mapping); err != nil {
		return nil, fmt.Errorf("failed to parse api-mappings: %w", err)
	}

	var references []ReferenceField

	// Recursively search for x-kubernetes-mapping
	findReferences(mapping, "", &references)

	return references, nil
}

func processKubernetesMapping(v map[string]any, path string, references *[]ReferenceField) {
	kubeMapping, exists := v["x-kubernetes-mapping"]
	if !exists {
		return
	}

	mapping, ok := kubeMapping.(map[string]any)
	if !ok {
		return
	}

	typeInfo, exists := mapping["type"]
	if !exists {
		return
	}

	typeMap, ok := typeInfo.(map[string]any)
	if !ok {
		return
	}

	kind, _ := typeMap["kind"].(string)
	group, _ := typeMap["group"].(string)
	version, _ := typeMap["version"].(string)

	// Extract field name from path
	pathParts := strings.Split(path, ".")
	fieldName := ""
	if len(pathParts) > 0 {
		fieldName = pathParts[len(pathParts)-1]
	}

	if kind != "" && fieldName != "" {
		ref := ReferenceField{
			FieldName:         fieldName,
			FieldPath:         path,
			ReferencedKind:    kind,
			ReferencedGroup:   group,
			ReferencedVersion: version,
		}
		ref.IndexerType = determineIndexerType(ref.ReferencedKind, ref.ReferencedGroup)
		*references = append(*references, ref)
	}
}

func findReferences(data any, path string, references *[]ReferenceField) {
	switch v := data.(type) {
	case map[string]any:
		// Check if this is a kubernetes mapping and process it
		processKubernetesMapping(v, path, references)

		for key, value := range v {
			newPath := path
			if newPath != "" {
				newPath += "."
			}
			newPath += key
			findReferences(value, newPath, references)
		}
	case []any:
		for i, item := range v {
			newPath := fmt.Sprintf("%s[%d]", path, i)
			findReferences(item, newPath, references)
		}
	}
}

func determineIndexerType(referencedKind, referencedGroup string) string {
	if referencedKind == "Group" && referencedGroup == "atlas.generated.mongodb.com" {
		return "project"
	}

	if referencedKind == "Secret" && referencedGroup == "" {
		return "credentials"
	}

	return "resource"
}

func GenerateIndexers(resultPath, crdKind, indexerOutDir string) error {
	references, err := ParseReferenceFields(resultPath, crdKind)
	if err != nil {
		return fmt.Errorf("failed to parse reference fields: %w", err)
	}

	if len(references) == 0 {
		fmt.Printf("No reference fields found for CRD %s, skipping indexer generation\n", crdKind)
		return nil
	}

	// Group references by target kind (e.g., all Secret refs together, all Group refs together)
	// Skip array-based references for now as they require iteration logic
	refsByKind := make(map[string][]ReferenceField)
	for _, ref := range references {
		// Skip references that are arrays for now
		if strings.Contains(ref.FieldPath, ".items.") {
			fmt.Printf("Skipping array-based reference %s in %s (array indexing not yet supported)\n", ref.FieldName, crdKind)
			continue
		}
		refsByKind[ref.ReferencedKind] = append(refsByKind[ref.ReferencedKind], ref)
	}

	// Generate one indexer per target kind
	for kind, refs := range refsByKind {
		indexerInfo := createIndexerInfoForKind(crdKind, kind, refs)
		if err := generateIndexerFile(crdKind, indexerInfo, indexerOutDir); err != nil {
			return fmt.Errorf("failed to generate indexer for kind %s: %w", kind, err)
		}
	}

	fmt.Printf("Generated indexers for CRD %s: %v\n", crdKind, refsByKind)
	return nil
}

func createIndexerInfoForKind(crdKind, targetKind string, refs []ReferenceField) IndexerInfo {
	resourceName := strings.ToLower(crdKind)

	// Build index name from all field paths
	indexParts := make([]string, 0, len(refs))
	for _, ref := range refs {
		indexParts = append(indexParts, ref.FieldName)
	}
	indexName := fmt.Sprintf("%s.%s", resourceName, strings.Join(indexParts, ","))

	return IndexerInfo{
		ResourceName:    crdKind,
		IndexerName:     indexName,
		TargetKind:      targetKind,
		ConstantName:    fmt.Sprintf("%sBy%sIndex", crdKind, targetKind),
		FunctionName:    fmt.Sprintf("New%sBy%sIndexer", crdKind, targetKind),
		ReferenceFields: refs,
	}
}

func generateIndexerFile(crdKind string, indexer IndexerInfo, indexerOutDir string) error {
	// Set default directory if not provided
	if indexerOutDir == "" {
		indexerOutDir = filepath.Join("..", "mongodb-atlas-kubernetes", "internal", "indexer")
	}

	// Ensure output directory exists
	if err := os.MkdirAll(indexerOutDir, 0755); err != nil {
		return fmt.Errorf("failed to create indexer directory: %w", err)
	}

	filename := fmt.Sprintf("%sby%s.go", strings.ToLower(crdKind), strings.ToLower(indexer.TargetKind))
	filePath := filepath.Join(indexerOutDir, filename)

	f := jen.NewFile("indexer")
	AddLicenseHeader(f)
	f.Comment("nolint:dupl")

	f.Const().Id(indexer.ConstantName).Op("=").Lit(indexer.IndexerName)

	// Add struct type with logger field
	structName := fmt.Sprintf("%sBy%sIndexer", crdKind, indexer.TargetKind)
	f.Type().Id(structName).Struct(
		jen.Id("logger").Op("*").Qual("go.uber.org/zap", "SugaredLogger"),
	)

	// Constructor
	f.Func().Id(indexer.FunctionName).Params(
		jen.Id("logger").Op("*").Qual("go.uber.org/zap", "Logger"),
	).Op("*").Id(structName).Block(
		jen.Return(jen.Op("&").Id(structName).Values(jen.Dict{
			jen.Id("logger"): jen.Id("logger").Dot("Named").Call(jen.Id(indexer.ConstantName)).Dot("Sugar").Call(),
		})),
	)

	// Object method
	f.Func().Params(jen.Op("*").Id(structName)).Id("Object").Params().Qual("sigs.k8s.io/controller-runtime/pkg/client", "Object").Block(
		jen.Return(jen.Op("&").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1", crdKind).Values()),
	)

	// Name method
	f.Func().Params(jen.Op("*").Id(structName)).Id("Name").Params().String().Block(
		jen.Return(jen.Id(indexer.ConstantName)),
	)

	// Keys method with logic for all reference fields
	generateKeysMethod(f, structName, crdKind, indexer)

	// Always generate helper Requests function for all reference types
	generateRequestsFunction(f, crdKind, indexer.TargetKind)

	if err := f.Save(filePath); err != nil {
		return fmt.Errorf("failed to save file %s: %w", filePath, err)
	}

	fmt.Printf("Generated indexer: %s\n", filePath)
	return nil
}

func generateKeysMethod(f *jen.File, structName, crdKind string, indexer IndexerInfo) {
	f.Comment("Keys extracts the index key(s) from the given object")

	// Build the block statements
	blockStatements := []jen.Code{
		// Type assertion
		jen.List(jen.Id("resource"), jen.Id("ok")).Op(":=").Id("object").Assert(jen.Op("*").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1", crdKind)),
		jen.If(jen.Op("!").Id("ok")).Block(
			jen.Id("i").Dot("logger").Dot("Errorf").Call(
				jen.Lit(fmt.Sprintf("expected *v1.%s but got %%T", crdKind)),
				jen.Id("object"),
			),
			jen.Return(jen.Nil()),
		),
		jen.Var().Id("keys").Index().String(),
	}

	// Add field extraction logic
	blockStatements = append(blockStatements, generateFieldExtractionCode(indexer.ReferenceFields)...)

	// Add return statement
	blockStatements = append(blockStatements, jen.Return(jen.Id("keys")))

	f.Func().Params(jen.Id("i").Op("*").Id(structName)).Id("Keys").Params(
		jen.Id("object").Qual("sigs.k8s.io/controller-runtime/pkg/client", "Object"),
	).Index().String().Block(blockStatements...)
}

func generateFieldExtractionCode(fields []ReferenceField) []jen.Code {
	code := make([]jen.Code, 0)

	for _, field := range fields {
		// Build the field path from the FieldPath
		// FieldPath looks like: "properties.spec.properties.<version>.properties.groupRef"
		// We need to convert this to: resource.Spec.<version>.GroupRef
		fieldAccessPath := buildFieldAccessPath(field.FieldPath)

		// Generate: if resource.Spec.<version>.GroupRef != nil && resource.Spec.<version>.GroupRef.Name != "" {
		//   keys = append(keys, types.NamespacedName{Name: resource.Spec.<version>.GroupRef.Name, Namespace: resource.Namespace}.String())
		// }
		code = append(code,
			jen.If(
				jen.Op("").Add(
					jen.Id(fieldAccessPath).Op("!=").Nil(),
				).Op("&&").Add(
					jen.Id(fieldAccessPath).Dot("Name").Op("!=").Lit(""),
				),
			).Block(
				jen.Id("keys").Op("=").Append(
					jen.Id("keys"),
					jen.Qual("k8s.io/apimachinery/pkg/types", "NamespacedName").Values(jen.Dict{
						jen.Id("Name"):      jen.Id(fieldAccessPath).Dot("Name"),
						jen.Id("Namespace"): jen.Id("resource").Dot("Namespace"),
					}).Dot("String").Call(),
				),
			),
		)
	}

	return code
}

func buildFieldAccessPath(fieldPath string) string {
	parts := strings.Split(fieldPath, ".")
	accessPath := []string{"resource"}

	for i := 0; i < len(parts); i++ {
		part := parts[i]

		// Skip "properties" and "items" keywords. Array based indexers are not supported for now
		if part == "properties" || part == "items" {
			continue
		}

		// Capitalize the first letter
		accessPath = append(accessPath, capitalizeFirst(part))
	}

	return strings.Join(accessPath, ".")
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func generateRequestsFunction(f *jen.File, crdKind string, targetKind string) {
	listTypeName := fmt.Sprintf("%sList", crdKind)
	// Make function name unique per targetKind to avoid duplicates when multiple indexers exist for same CRD
	requestsFuncName := fmt.Sprintf("%sRequestsFrom%s", crdKind, targetKind)
	f.Func().Id(requestsFuncName).Params(
		jen.Id("list").Op("*").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1", listTypeName),
	).Index().Qual("sigs.k8s.io/controller-runtime/pkg/reconcile", "Request").Block(
		jen.Id("requests").Op(":=").Make(jen.Index().Qual("sigs.k8s.io/controller-runtime/pkg/reconcile", "Request"), jen.Lit(0), jen.Len(jen.Id("list").Dot("Items"))),
		jen.For(jen.List(jen.Id("_"), jen.Id("item")).Op(":=").Range().Id("list").Dot("Items")).Block(
			jen.Id("requests").Op("=").Append(jen.Id("requests"), jen.Qual("sigs.k8s.io/controller-runtime/pkg/reconcile", "Request").Values(jen.Dict{
				jen.Id("NamespacedName"): jen.Qual("k8s.io/apimachinery/pkg/types", "NamespacedName").Values(jen.Dict{
					jen.Id("Name"):      jen.Id("item").Dot("Name"),
					jen.Id("Namespace"): jen.Id("item").Dot("Namespace"),
				}),
			})),
		),
		jen.Return(jen.Id("requests")),
	)
}

// TODO: UpdateIndexerRegistry needs to be reimplemented to work with the new kind-based indexer approach
// For now, indexers need to be manually registered in the indexer registry
