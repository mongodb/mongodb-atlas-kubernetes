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
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dave/jennifer/jen"
	"gopkg.in/yaml.v3"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	clientsetscheme "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/scaffolder/internal/boilerplate"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/scaffolder/internal/config"
)

// ReferenceField represents a field that references another Kubernetes resource.
type ReferenceField struct {
	FieldName         string
	FieldPath         string
	ReferencedKind    string
	ReferencedGroup   string
	ReferencedVersion string
	RequiredSegments  []bool
	ArrayBoundaries   []ArrayBoundary // all array boundaries (empty if not array-based)
}

// ArrayBoundary represents a single array level in a nested array path.
type ArrayBoundary struct {
	ArrayPath string // path to the array container
	ItemPath  string // path within array item to the next array or final field
}

// IsArrayBased returns true if the reference is inside an array.
func (r ReferenceField) IsArrayBased() bool {
	return len(r.ArrayBoundaries) > 0
}

// ArrayPath returns the path to the first array container (for backwards compatibility).
func (r ReferenceField) ArrayPath() string {
	if len(r.ArrayBoundaries) == 0 {
		return ""
	}
	return r.ArrayBoundaries[0].ArrayPath
}

// ItemPath returns the path within the last array item (for backwards compatibility).
func (r ReferenceField) ItemPath() string {
	if len(r.ArrayBoundaries) == 0 {
		return ""
	}
	return r.ArrayBoundaries[len(r.ArrayBoundaries)-1].ItemPath
}

// IndexerInfo contains information needed to generate an indexer.
type IndexerInfo struct {
	ResourceName    string
	IndexerName     string
	TargetKind      string
	ConstantName    string
	FunctionName    string
	ReferenceFields []ReferenceField // All fields that reference this kind
}

// DependentInfo represents a resource that depends on another resource (has a reference to it).
type DependentInfo struct {
	// DependentKind is the kind of the resource that has the reference (e.g., "FlexCluster")
	DependentKind string
	// TargetKind is the kind being referenced (e.g., "Group")
	TargetKind string
	// IndexerConstantName is the constant name for the indexer (e.g., "FlexClusterByGroupIndex")
	IndexerConstantName string
	// MapFuncName is the name of the MapFunc function (e.g., "NewFlexClusterByGroupMapFunc")
	MapFuncName string
}

// GenerateIndexers generates indexer files for a CRD kind.
func GenerateIndexers(resultPath, crdKind, indexerOutDir, indexerTypesPath string) error {
	references, err := ParseReferenceFields(resultPath, crdKind)
	if err != nil {
		return fmt.Errorf("failed to parse reference fields: %w", err)
	}

	if len(references) == 0 {
		fmt.Printf("No reference fields found for CRD %s, skipping indexer generation\n", crdKind)
		return nil
	}

	// Group references by target kind (e.g., all Secret refs together, all Group refs together)
	refsByKind := make(map[string][]ReferenceField)
	for _, ref := range references {
		refsByKind[ref.ReferencedKind] = append(refsByKind[ref.ReferencedKind], ref)
	}

	// Generate one indexer per target kind
	for kind, refs := range refsByKind {
		indexerInfo := createIndexerInfoForKind(crdKind, kind, refs)
		if err := generateIndexerFile(crdKind, indexerInfo, indexerOutDir, indexerTypesPath); err != nil {
			return fmt.Errorf("failed to generate indexer for kind %s: %w", kind, err)
		}
	}

	fmt.Printf("Generated indexers for CRD %s: %v\n", crdKind, refsByKind)
	return nil
}

// ParseReferenceFields parses reference fields from a CRD YAML file.
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

// ParseDependentReferences finds all resources that have references to the given targetKind.
func ParseDependentReferences(resultPath, targetKind string) ([]DependentInfo, error) {
	crds, err := config.ListCRDs(resultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to list CRDs: %w", err)
	}

	var dependents []DependentInfo

	for _, crd := range crds {
		if crd.Kind == targetKind {
			continue
		}

		refs, err := ParseReferenceFields(resultPath, crd.Kind)
		if err != nil {
			return nil, err
		}

		for _, ref := range refs {
			if ref.ReferencedKind == targetKind {
				dependent := DependentInfo{
					DependentKind:       crd.Kind,
					TargetKind:          targetKind,
					IndexerConstantName: fmt.Sprintf("%sBy%sIndex", crd.Kind, targetKind),
					MapFuncName:         fmt.Sprintf("New%sBy%sMapFunc", crd.Kind, targetKind),
				}
				dependents = append(dependents, dependent)
				break
			}
		}
	}

	return dependents, nil
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
	crd, err := decode(content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode YAML: %w", err)
	}

	if crd.Spec.Names.Kind != targetKind {
		return nil, fmt.Errorf("not target CRD")
	}

	var references []ReferenceField

	if apiMappingsStr, exists := crd.GetAnnotations()["api-mappings"]; exists {
		refs, err := parseAPIMapping(apiMappingsStr, crd.Spec.Versions[0].Schema.OpenAPIV3Schema)
		if err == nil {
			references = append(references, refs...)
		}
	}

	return references, nil
}

func decode(content []byte) (*apiextensionsv1.CustomResourceDefinition, error) {
	sch := runtime.NewScheme()
	_ = clientsetscheme.AddToScheme(sch)
	_ = apiextensions.AddToScheme(sch)
	_ = apiextensionsv1.AddToScheme(sch)
	_ = apiextensionsv1.RegisterConversions(sch)
	_ = apiextensionsv1beta1.AddToScheme(sch)
	_ = apiextensionsv1beta1.RegisterConversions(sch)

	decodeFunc := serializer.NewCodecFactory(sch).UniversalDeserializer().Decode

	obj, _, err := decodeFunc(content, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decode YAML: %w", err)
	}

	kind := obj.GetObjectKind().GroupVersionKind().Kind
	if kind != "CustomResourceDefinition" {
		return nil, fmt.Errorf("unexpected kind %q: %w", kind, err)
	}

	crd := &apiextensionsv1.CustomResourceDefinition{}
	err = sch.Convert(obj, crd, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to convert CRD object: %w", err)
	}
	return crd, nil
}

func parseAPIMapping(apiMappingsStr string, schema *apiextensionsv1.JSONSchemaProps) ([]ReferenceField, error) {
	var mapping map[string]any
	if err := yaml.Unmarshal([]byte(apiMappingsStr), &mapping); err != nil {
		return nil, fmt.Errorf("failed to parse api-mappings: %w", err)
	}

	var references []ReferenceField
	findReferences(mapping, "", schema, nil, &references)

	return references, nil
}

func processKubernetesMapping(v map[string]any, path string, requiredSegments []bool, references *[]ReferenceField) {
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

	pathParts := strings.Split(path, ".")
	fieldName := ""
	if len(pathParts) > 0 {
		fieldName = pathParts[len(pathParts)-1]
	}

	if kind != "" && fieldName != "" {
		reqCopy := make([]bool, len(requiredSegments))
		copy(reqCopy, requiredSegments)

		arrayBoundaries := parseAllArrayBoundaries(path)

		ref := ReferenceField{
			FieldName:         fieldName,
			FieldPath:         path,
			ReferencedKind:    kind,
			ReferencedGroup:   group,
			ReferencedVersion: version,
			RequiredSegments:  reqCopy,
			ArrayBoundaries:   arrayBoundaries,
		}
		*references = append(*references, ref)
	}
}

func findReferences(data any, path string, schema *apiextensionsv1.JSONSchemaProps, requiredSegments []bool, references *[]ReferenceField) {
	switch v := data.(type) {
	case map[string]any:
		processKubernetesMapping(v, path, requiredSegments, references)

		for key, value := range v {
			newPath := path
			if newPath != "" {
				newPath += "."
			}
			newPath += key

			required, childSchema := getSchemaForPathSegment(schema, key)

			newRequiredSegments := requiredSegments
			if key != "properties" && key != "items" {
				if key == "spec" && len(requiredSegments) == 0 {
					required = false
				}
				newRequiredSegments = append(requiredSegments, required)
			}

			findReferences(value, newPath, childSchema, newRequiredSegments, references)
		}
	case []any:
		for i, item := range v {
			newPath := fmt.Sprintf("%s[%d]", path, i)

			var childSchema *apiextensionsv1.JSONSchemaProps
			if schema != nil {
				_, childSchema = getSchemaForPathSegment(schema, "items")
			}
			findReferences(item, newPath, childSchema, requiredSegments, references)
		}
	}
}

func getSchemaForPathSegment(schema *apiextensionsv1.JSONSchemaProps, key string) (bool, *apiextensionsv1.JSONSchemaProps) {
	if schema == nil {
		return false, nil
	}

	if key == "properties" {
		return false, schema
	}

	required := false
	for _, req := range schema.Required {
		if req == key {
			required = true
			break
		}
	}

	if key == "items" {
		if schema.Items != nil && schema.Items.Schema != nil {
			return false, schema.Items.Schema
		}
		return false, nil
	}

	if schema.Properties != nil {
		if p, ok := schema.Properties[key]; ok {
			return required, &p
		}
	}

	return false, nil
}

func createIndexerInfoForKind(crdKind, targetKind string, refs []ReferenceField) IndexerInfo {
	resourceName := strings.ToLower(crdKind)

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

func generateIndexerFile(crdKind string, indexer IndexerInfo, indexerOutDir, indexerTypesPath string) error {
	if indexerOutDir == "" {
		indexerOutDir = filepath.Join("..", "mongodb-atlas-kubernetes", "internal", "indexer")
	}

	if err := os.MkdirAll(indexerOutDir, 0755); err != nil {
		return fmt.Errorf("failed to create indexer directory: %w", err)
	}

	filename := fmt.Sprintf("%sby%s.go", strings.ToLower(crdKind), strings.ToLower(indexer.TargetKind))
	filePath := filepath.Join(indexerOutDir, filename)

	f := jen.NewFile("indexer")
	boilerplate.AddLicenseHeader(f)
	f.Comment("nolint:dupl")

	f.Const().Id(indexer.ConstantName).Op("=").Lit(indexer.IndexerName)

	structName := fmt.Sprintf("%sBy%sIndexer", crdKind, indexer.TargetKind)
	f.Type().Id(structName).Struct(
		jen.Id("logger").Op("*").Qual("go.uber.org/zap", "SugaredLogger"),
	)

	f.Func().Id(indexer.FunctionName).Params(
		jen.Id("logger").Op("*").Qual("go.uber.org/zap", "Logger"),
	).Op("*").Id(structName).Block(
		jen.Return(jen.Op("&").Id(structName).Values(jen.Dict{
			jen.Id("logger"): jen.Id("logger").Dot("Named").Call(jen.Id(indexer.ConstantName)).Dot("Sugar").Call(),
		})),
	)

	f.Func().Params(jen.Op("*").Id(structName)).Id("Object").Params().Qual("sigs.k8s.io/controller-runtime/pkg/client", "Object").Block(
		jen.Return(jen.Op("&").Qual(indexerTypesPath, crdKind).Values()),
	)

	f.Func().Params(jen.Op("*").Id(structName)).Id("Name").Params().String().Block(
		jen.Return(jen.Id(indexer.ConstantName)),
	)

	generateKeysMethod(f, structName, crdKind, indexer, indexerTypesPath)

	f.Line()

	generateMapFunc(f, crdKind, indexer, indexerTypesPath)

	if err := f.Save(filePath); err != nil {
		return fmt.Errorf("failed to save file %s: %w", filePath, err)
	}

	fmt.Printf("Generated indexer: %s\n", filePath)
	return nil
}

func generateKeysMethod(f *jen.File, structName, crdKind string, indexer IndexerInfo, indexerTypesPath string) {
	f.Comment("Keys extracts the index key(s) from the given object")

	blockStatements := []jen.Code{
		jen.List(jen.Id("resource"), jen.Id("ok")).Op(":=").Id("object").Assert(jen.Op("*").Qual(indexerTypesPath, crdKind)),
		jen.If(jen.Op("!").Id("ok")).Block(
			jen.Id("i").Dot("logger").Dot("Errorf").Call(
				jen.Lit(fmt.Sprintf("expected *v1.%s but got %%T", crdKind)),
				jen.Id("object"),
			),
			jen.Return(jen.Nil()),
		),
		jen.Var().Id("keys").Index().String(),
	}

	blockStatements = append(blockStatements, generateFieldExtractionCode(indexer.ReferenceFields)...)
	blockStatements = append(blockStatements, jen.Return(jen.Id("keys")))

	f.Func().Params(jen.Id("i").Op("*").Id(structName)).Id("Keys").Params(
		jen.Id("object").Qual("sigs.k8s.io/controller-runtime/pkg/client", "Object"),
	).Index().String().Block(blockStatements...)
}

func generateFieldExtractionCode(fields []ReferenceField) []jen.Code {
	code := make([]jen.Code, 0)

	for _, field := range fields {
		if field.IsArrayBased() {
			code = append(code, generateArrayFieldExtractionCode(field))
		} else {
			fieldAccessPath := BuildFieldAccessPath(field.FieldPath)
			nilCheckCondition := BuildNilCheckConditions(fieldAccessPath, field.RequiredSegments)
			condition := nilCheckCondition.Op("&&").Add(
				jen.Id(fieldAccessPath).Dot("Name").Op("!=").Lit(""),
			)

			code = append(code,
				jen.If(condition).Block(
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
	}

	return code
}

func generateArrayFieldExtractionCode(field ReferenceField) jen.Code {
	boundaries := field.ArrayBoundaries

	lastBoundary := boundaries[len(boundaries)-1]
	lastItemPath := BuildFieldAccessPath(lastBoundary.ItemPath)

	lastArrayParts := strings.Split(lastBoundary.ArrayPath, ".")
	lastArrayFieldName := lastArrayParts[len(lastArrayParts)-1]
	lastLoopVar := GenerateLoopVariableName(lastArrayFieldName)

	finalFieldPath := strings.Replace(lastItemPath, "resource", lastLoopVar, 1)

	innermostBlock := jen.If(
		jen.Id(finalFieldPath).Op("!=").Nil().Op("&&").Add(
			jen.Id(finalFieldPath).Dot("Name").Op("!=").Lit(""),
		),
	).Block(
		jen.Id("keys").Op("=").Append(
			jen.Id("keys"),
			jen.Qual("k8s.io/apimachinery/pkg/types", "NamespacedName").Values(jen.Dict{
				jen.Id("Name"):      jen.Id(finalFieldPath).Dot("Name"),
				jen.Id("Namespace"): jen.Id("resource").Dot("Namespace"),
			}).Dot("String").Call(),
		),
	)

	currentBlock := innermostBlock

	for i := len(boundaries) - 1; i >= 1; i-- {
		boundary := boundaries[i]
		prevBoundary := boundaries[i-1]

		prevArrayParts := strings.Split(prevBoundary.ArrayPath, ".")
		prevArrayFieldName := prevArrayParts[len(prevArrayParts)-1]
		prevLoopVar := GenerateLoopVariableName(prevArrayFieldName)

		arrayParts := strings.Split(boundary.ArrayPath, ".")
		arrayFieldName := arrayParts[len(arrayParts)-1]
		loopVar := GenerateLoopVariableName(arrayFieldName)

		arrayAccessFromPrev := BuildFieldAccessPath(prevBoundary.ItemPath)
		arrayAccessInLoop := strings.Replace(arrayAccessFromPrev, "resource", prevLoopVar, 1)

		currentBlock = jen.If(jen.Id(arrayAccessInLoop).Op("!=").Nil()).Block(
			jen.For(
				jen.List(jen.Id("_"), jen.Id(loopVar)).Op(":=").Range().Op("*").Id(arrayAccessInLoop),
			).Block(currentBlock),
		)
	}

	firstBoundary := boundaries[0]
	firstArrayAccessPath := BuildFieldAccessPath(firstBoundary.ArrayPath)
	firstArrayParts := strings.Split(firstBoundary.ArrayPath, ".")
	firstArrayFieldName := firstArrayParts[len(firstArrayParts)-1]
	firstLoopVar := GenerateLoopVariableName(firstArrayFieldName)

	return jen.If(jen.Id(firstArrayAccessPath).Op("!=").Nil()).Block(
		jen.For(
			jen.List(jen.Id("_"), jen.Id(firstLoopVar)).Op(":=").Range().Op("*").Id(firstArrayAccessPath),
		).Block(currentBlock),
	)
}

// BuildFieldAccessPath converts a field path to a Go field access path.
func BuildFieldAccessPath(fieldPath string) string {
	parts := strings.Split(fieldPath, ".")
	accessPath := []string{"resource"}

	for i := 0; i < len(parts); i++ {
		part := parts[i]

		if part == "properties" {
			continue
		}

		if part == "items" {
			if i+1 < len(parts) && parts[i+1] == "properties" {
				continue
			}
		}

		accessPath = append(accessPath, CapitalizeFirst(part))
	}

	return strings.Join(accessPath, ".")
}

// CapitalizeFirst capitalizes the first letter of a string.
func CapitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func parseAllArrayBoundaries(fieldPath string) []ArrayBoundary {
	var boundaries []ArrayBoundary

	const delimiter = ".items.properties."
	parts := strings.Split(fieldPath, delimiter)

	if len(parts) <= 1 {
		return nil
	}

	for i := 0; i < len(parts)-1; i++ {
		arrayPath := parts[i]

		if i > 0 {
			arrayPath = "properties." + extractArrayFieldFromItemPath(parts[i])
		}

		itemPath := "properties." + parts[i+1]

		boundaries = append(boundaries, ArrayBoundary{
			ArrayPath: arrayPath,
			ItemPath:  itemPath,
		})
	}

	return boundaries
}

func extractArrayFieldFromItemPath(path string) string {
	path = strings.TrimPrefix(path, "properties.")

	lastPropertiesIdx := strings.LastIndex(path, ".properties.")
	if lastPropertiesIdx >= 0 {
		return path[lastPropertiesIdx+len(".properties."):]
	}

	return path
}

// GenerateLoopVariableName generates a loop variable name from an array field name.
func GenerateLoopVariableName(arrayFieldName string) string {
	if arrayFieldName == "" {
		return "item"
	}

	name := strings.ToLower(arrayFieldName)

	if strings.HasSuffix(name, "ies") {
		return name[:len(name)-3] + "y"
	}
	if strings.HasSuffix(name, "ses") || strings.HasSuffix(name, "ches") || strings.HasSuffix(name, "xes") {
		return name[:len(name)-2]
	}
	if strings.HasSuffix(name, "s") {
		return name[:len(name)-1]
	}

	return name + "Item"
}

// BuildNilCheckConditions creates a compound nil check condition for a field access path.
func BuildNilCheckConditions(fieldAccessPath string, requiredSegments []bool) *jen.Statement {
	segments := strings.Split(fieldAccessPath, ".")

	if len(requiredSegments) == 0 {
		return BuildDotChain(segments).Op("!=").Nil()
	}

	if len(requiredSegments) != len(segments)-1 {
		return BuildDotChain(segments).Op("!=").Nil()
	}

	var conditions *jen.Statement

	for i := 1; i < len(segments); i++ {
		requiredIndex := i - 1

		if requiredSegments[requiredIndex] {
			continue
		}

		if segments[i] == "Spec" {
			continue
		}

		pathSegments := segments[:i+1]
		nilCheck := BuildDotChain(pathSegments).Op("!=").Nil()

		if conditions == nil {
			conditions = nilCheck
		} else {
			conditions = conditions.Op("&&").Add(nilCheck)
		}
	}

	if conditions == nil {
		return BuildDotChain(segments).Op("!=").Nil()
	}

	return conditions
}

// BuildDotChain creates a jen.Statement for a dot-separated path.
func BuildDotChain(segments []string) *jen.Statement {
	if len(segments) == 0 {
		return jen.Null()
	}

	stmt := jen.Id(segments[0])
	for i := 1; i < len(segments); i++ {
		stmt = stmt.Dot(segments[i])
	}

	return stmt
}

func generateMapFunc(f *jen.File, crdKind string, indexer IndexerInfo, indexerTypesPath string) {
	f.Func().
		Id(fmt.Sprintf("New%sBy%sMapFunc", crdKind, indexer.TargetKind)).
		Params(
			jen.Id("kubeClient").Qual("sigs.k8s.io/controller-runtime/pkg/client", "Client"),
		).
		Qual("sigs.k8s.io/controller-runtime/pkg/handler", "MapFunc").
		Block(
			jen.Return(
				jen.Func().
					Params(
						jen.Id("ctx").Qual("context", "Context"),
						jen.Id("obj").Qual("sigs.k8s.io/controller-runtime/pkg/client", "Object"),
					).
					Index().
					Qual("sigs.k8s.io/controller-runtime/pkg/reconcile", "Request").
					Block(
						jen.Id("logger").Op(":=").
							Qual("sigs.k8s.io/controller-runtime/pkg/log", "FromContext").
							Call(jen.Id("ctx")),
						jen.Line(),
						jen.Id("listOpts").Op(":=").
							Op("&").
							Qual("sigs.k8s.io/controller-runtime/pkg/client", "ListOptions").
							Values(
								jen.Dict{
									jen.Id("FieldSelector"): jen.
										Qual("k8s.io/apimachinery/pkg/fields", "OneTermEqualSelector").
										Call(
											jen.Id(indexer.ConstantName),
											jen.
												Qual("k8s.io/apimachinery/pkg/types", "NamespacedName").
												Values(
													jen.Dict{
														jen.Id("Name"):      jen.Id("obj").Dot("GetName").Call(),
														jen.Id("Namespace"): jen.Id("obj").Dot("GetNamespace").Call(),
													},
												).
												Dot("String").
												Call(),
										),
								},
							),
						jen.Line(),
						jen.Id("list").Op(":=").
							Op("&").
							Qual(indexerTypesPath, fmt.Sprintf("%sList", crdKind)).
							Values(),
						jen.Id("err").Op(":=").
							Id("kubeClient").Dot("List").
							Call(
								jen.Id("ctx"),
								jen.Id("list"),
								jen.Id("listOpts"),
							),
						jen.If(jen.Id("err").Op("!=").Nil()).Block(
							jen.Id("logger").Dot("Error").
								Call(
									jen.Id("err"),
									jen.Lit(fmt.Sprintf("failed to list %v objects", crdKind)),
								),
							jen.Return(jen.Nil()),
						),
						jen.Line(),
						jen.Id("requests").Op(":=").
							Make(
								jen.Index().Qual("sigs.k8s.io/controller-runtime/pkg/reconcile", "Request"),
								jen.Lit(0),
								jen.Len(jen.Id("list").Dot("Items")),
							),
						jen.For(
							jen.List(
								jen.Id("_"),
								jen.Id("item"),
							).
								Op(":=").
								Range().
								Id("list").Dot("Items"),
						).Block(
							jen.Id("requests").Op("=").
								Append(
									jen.Id("requests"),
									jen.Qual("sigs.k8s.io/controller-runtime/pkg/reconcile", "Request").
										Values(
											jen.Dict{
												jen.Id("NamespacedName"): jen.
													Qual("k8s.io/apimachinery/pkg/types", "NamespacedName").
													Values(
														jen.Dict{
															jen.Id("Name"):      jen.Id("item").Dot("Name"),
															jen.Id("Namespace"): jen.Id("item").Dot("Namespace"),
														},
													),
											},
										),
								),
						),
						jen.Line(),
						jen.Return(jen.Id("requests")),
					),
			),
		)
}
