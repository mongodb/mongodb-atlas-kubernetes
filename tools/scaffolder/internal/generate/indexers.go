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
	"regexp"
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
	ResourceName string
	IndexerName  string
	IndexerType  string
	ConstantName string
	FunctionName string
	BaseType     string
	RefFieldName string
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

		if refs != nil {
			return refs, nil
		}
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
			return nil, err
		}
		return refs, nil
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

	indexersByType := make(map[string][]IndexerInfo)

	for _, ref := range references {
		info := createIndexerInfo(crdKind, ref)
		indexersByType[ref.IndexerType] = append(indexersByType[ref.IndexerType], info)
	}

	for indexerType, indexers := range indexersByType {
		if err := generateIndexerFiles(crdKind, indexerType, indexers, indexerOutDir); err != nil {
			return fmt.Errorf("failed to generate %s indexers: %w", indexerType, err)
		}
	}

	fmt.Printf("Generated indexers for CRD %s: %v\n", crdKind, indexersByType)
	return nil
}

func createIndexerInfo(crdKind string, ref ReferenceField) IndexerInfo {
	resourceName := strings.ToLower(crdKind)

	switch ref.IndexerType {
	case "project":
		return IndexerInfo{
			ResourceName: crdKind,
			IndexerName:  fmt.Sprintf("atlas%s.spec.projectRef", resourceName),
			IndexerType:  "project",
			ConstantName: fmt.Sprintf("Atlas%sByProjectIndex", crdKind),
			FunctionName: fmt.Sprintf("NewAtlas%sByProjectIndexer", crdKind),
			BaseType:     "AtlasReferrerByProjectIndexerBase",
			RefFieldName: ref.FieldName,
		}
	case "credentials":
		return IndexerInfo{
			ResourceName: crdKind,
			IndexerName:  fmt.Sprintf("atlas%s.credentials", resourceName),
			IndexerType:  "credentials",
			ConstantName: fmt.Sprintf("Atlas%sCredentialsIndex", crdKind),
			FunctionName: fmt.Sprintf("NewAtlas%sByCredentialIndexer", crdKind),
			BaseType:     "LocalCredentialIndexer",
			RefFieldName: ref.FieldName,
		}
	default:
		return IndexerInfo{
			ResourceName: crdKind,
			IndexerName:  fmt.Sprintf("atlas%s.%s", resourceName, strings.ToLower(ref.FieldName)),
			IndexerType:  "resource",
			ConstantName: fmt.Sprintf("Atlas%sBy%sIndex", crdKind, ref.ReferencedKind),
			FunctionName: fmt.Sprintf("NewAtlas%sBy%sIndexer", crdKind, ref.ReferencedKind),
			BaseType:     "ResourceIndexer",
			RefFieldName: ref.FieldName,
		}
	}
}

func generateIndexerFiles(crdKind, indexerType string, indexers []IndexerInfo, indexerOutDir string) error {
	// Set default directory if not provided
	if indexerOutDir == "" {
		indexerOutDir = filepath.Join("..", "mongodb-atlas-kubernetes", "internal", "indexer")
	}

	switch indexerType {
	case "project":
		return generateProjectIndexer(crdKind, indexers[0], indexerOutDir)
	case "credentials":
		return generateCredentialsIndexer(crdKind, indexers[0], indexerOutDir)
	default:
		for _, indexer := range indexers {
			if err := generateResourceIndexer(crdKind, indexer, indexerOutDir); err != nil {
				return err
			}
		}
		return nil
	}
}

func generateProjectIndexer(crdKind string, indexer IndexerInfo, indexerOutDir string) error {
	// Ensure output directory exists
	if err := os.MkdirAll(indexerOutDir, 0755); err != nil {
		return fmt.Errorf("failed to create indexer directory: %w", err)
	}

	filename := fmt.Sprintf("atlas%sprojects.go", strings.ToLower(crdKind))
	filePath := filepath.Join(indexerOutDir, filename)

	f := jen.NewFile("indexer")

	AddLicenseHeader(f)
	f.Comment("nolint:dupl")

	f.ImportAlias("github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1", "akov2")

	f.Const().Id(indexer.ConstantName).Op("=").Lit(indexer.IndexerName)

	// Add struct type
	structName := fmt.Sprintf("Atlas%sByProjectIndexer", crdKind)
	f.Type().Id(structName).Struct(
		jen.Id("AtlasReferrerByProjectIndexerBase"),
	)

	f.Func().Id(indexer.FunctionName).Params(
		jen.Id("logger").Op("*").Qual("go.uber.org/zap", "Logger"),
	).Op("*").Id(structName).Block(
		jen.Return(jen.Op("&").Id(structName).Values(jen.Dict{
			jen.Id("AtlasReferrerByProjectIndexerBase"): jen.Op("*").Id("NewAtlasReferrerByProjectIndexer").Call(
				jen.Id("logger"),
				jen.Id(indexer.ConstantName),
			),
		})),
	)

	f.Func().Params(jen.Op("*").Id(structName)).Id("Object").Params().Qual("sigs.k8s.io/controller-runtime/pkg/client", "Object").Block(
		jen.Return(jen.Op("&").Id("akov2").Dot("Atlas" + crdKind).Values()),
	)

	if err := f.Save(filePath); err != nil {
		return fmt.Errorf("failed to save file %s: %w", filePath, err)
	}

	fmt.Printf("Generated project indexer: %s\n", filePath)
	return nil
}

func generateCredentialsIndexer(crdKind string, indexer IndexerInfo, indexerOutDir string) error {
	if err := os.MkdirAll(indexerOutDir, 0755); err != nil {
		return fmt.Errorf("failed to create indexer directory: %w", err)
	}

	filename := fmt.Sprintf("atlas%scredentials.go", strings.ToLower(crdKind))
	filePath := filepath.Join(indexerOutDir, filename)

	f := jen.NewFile("indexer")

	AddLicenseHeader(f)
	f.ImportAlias("github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1", "akov2")

	f.Const().Id(indexer.ConstantName).Op("=").Lit(indexer.IndexerName)

	f.Func().Id(indexer.FunctionName).Params(
		jen.Id("logger").Op("*").Qual("go.uber.org/zap", "Logger"),
	).Op("*").Id("LocalCredentialIndexer").Block(
		jen.Return(jen.Id("NewLocalCredentialsIndexer").Call(
			jen.Id(indexer.ConstantName),
			jen.Op("&").Id("akov2").Dot("Atlas"+crdKind).Values(),
			jen.Id("logger"),
		)),
	)

	listTypeName := fmt.Sprintf("Atlas%sList", crdKind)
	requestsFuncName := fmt.Sprintf("%sRequests", crdKind)
	f.Func().Id(requestsFuncName).Params(
		jen.Id("list").Op("*").Id("akov2").Dot(listTypeName),
	).Index().Qual("sigs.k8s.io/controller-runtime/pkg/reconcile", "Request").Block(
		jen.Id("requests").Op(":=").Make(jen.Index().Qual("sigs.k8s.io/controller-runtime/pkg/reconcile", "Request"), jen.Lit(0), jen.Len(jen.Id("list").Dot("Items"))),
		jen.For(jen.List(jen.Id("_"), jen.Id("item")).Op(":=").Range().Id("list").Dot("Items")).Block(
			jen.Id("requests").Op("=").Append(jen.Id("requests"), jen.Id("toRequest").Call(jen.Op("&").Id("item"))),
		),
		jen.Return(jen.Id("requests")),
	)

	if err := f.Save(filePath); err != nil {
		return fmt.Errorf("failed to save file %s: %w", filePath, err)
	}

	fmt.Printf("Generated credentials indexer: %s\n", filePath)
	return nil
}

func generateResourceIndexer(crdKind string, indexer IndexerInfo, indexerOutDir string) error {
	// For now, resource indexers follow the same pattern to project indexers. Can be extend later
	return generateProjectIndexer(crdKind, indexer, indexerOutDir)
}

// TODO: this is kinda hacky, consider regenerating the entire file, or have a separate one for new indexers
func UpdateIndexerRegistry(crdKind, indexerOutDir string) error {
	const idxArrayLen = 4
	registryFile := filepath.Join(indexerOutDir, "indexer.go")

	content, err := os.ReadFile(registryFile)
	if err != nil {
		return fmt.Errorf("failed to read indexer.go: %w", err)
	}

	projectIndexerCall := fmt.Sprintf("NewAtlas%sByProjectIndexer(logger),", crdKind)
	credentialsIndexerCall := fmt.Sprintf("NewAtlas%sByCredentialIndexer(logger),", crdKind)

	contentStr := string(content)

	// Check if indexers are already registered
	if strings.Contains(contentStr, projectIndexerCall) && strings.Contains(contentStr, credentialsIndexerCall) {
		fmt.Printf("Indexers for %s are already registered\n", crdKind)
		return nil
	}

	// Find the "IndexesArray"
	re := regexp.MustCompile(`(indexers = append\(indexers,\s*\n)([\s\S]*?)(\s*\)\s*if version\.IsExperimental\(\))`)
	matches := re.FindStringSubmatch(contentStr)

	if len(matches) != idxArrayLen {
		return fmt.Errorf("could not find indexers append block in indexer.go")
	}

	newIndexers := matches[2]

	if !strings.Contains(newIndexers, projectIndexerCall) {
		newIndexers += fmt.Sprintf("\t\t%s\n", projectIndexerCall)
	}

	if !strings.Contains(newIndexers, credentialsIndexerCall) {
		newIndexers += fmt.Sprintf("\t\t%s\n", credentialsIndexerCall)
	}

	newContent := strings.Replace(contentStr, matches[0], matches[1]+newIndexers+matches[3], 1)

	if err := os.WriteFile(registryFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write indexer.go: %w", err)
	}

	fmt.Printf("Updated indexer registry for %s\n", crdKind)
	return nil
}
