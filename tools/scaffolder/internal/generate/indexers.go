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
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	clientsetscheme "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

type ReferenceField struct {
	FieldName         string
	FieldPath         string
	ReferencedKind    string
	ReferencedGroup   string
	ReferencedVersion string
	RequiredSegments  []bool
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
	crd, err := Decode(content)
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

	// Return empty slice for CRDs with no references (this is valid)
	return references, nil
}

func Decode(content []byte) (*apiextensionsv1.CustomResourceDefinition, error) {
	sch := runtime.NewScheme()
	_ = clientsetscheme.AddToScheme(sch)
	_ = apiextensions.AddToScheme(sch)
	_ = apiextensionsv1.AddToScheme(sch)
	_ = apiextensionsv1.RegisterConversions(sch)
	_ = apiextensionsv1beta1.AddToScheme(sch)
	_ = apiextensionsv1beta1.RegisterConversions(sch)

	decode := serializer.NewCodecFactory(sch).UniversalDeserializer().Decode

	obj, _, err := decode(content, nil, nil)
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
	// Parse the api-mappings as a generic map to handle any nesting level
	var mapping map[string]any
	if err := yaml.Unmarshal([]byte(apiMappingsStr), &mapping); err != nil {
		return nil, fmt.Errorf("failed to parse api-mappings: %w", err)
	}

	var references []ReferenceField

	// Recursively search for x-kubernetes-mapping
	// Start with empty requiredSegments slice
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

	// Extract field name from path
	pathParts := strings.Split(path, ".")
	fieldName := ""
	if len(pathParts) > 0 {
		fieldName = pathParts[len(pathParts)-1]
	}

	if kind != "" && fieldName != "" {
		// Make a copy of requiredSegments to avoid sharing the underlying array
		reqCopy := make([]bool, len(requiredSegments))
		copy(reqCopy, requiredSegments)

		ref := ReferenceField{
			FieldName:         fieldName,
			FieldPath:         path,
			ReferencedKind:    kind,
			ReferencedGroup:   group,
			ReferencedVersion: version,
			RequiredSegments:  reqCopy,
		}
		*references = append(*references, ref)
	}
}

func findReferences(data any, path string, schema *apiextensionsv1.JSONSchemaProps, requiredSegments []bool, references *[]ReferenceField) {
	switch v := data.(type) {
	case map[string]any:
		// Check if this is a kubernetes mapping and process it
		processKubernetesMapping(v, path, requiredSegments, references)

		for key, value := range v {
			newPath := path
			if newPath != "" {
				newPath += "."
			}
			newPath += key

			// Resolve the child schema for this key and pass it into the recursive call
			required, childSchema := getSchemaForPathSegment(schema, key)

			// Track required status for meaningful segments (not "properties" or "items")
			newRequiredSegments := requiredSegments
			if key != "properties" && key != "items" {
				// Special case: the top-level "spec" property is never required per Kubernetes convention
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

			// For arrays pass the items schema if available
			var childSchema *apiextensionsv1.JSONSchemaProps
			if schema != nil {
				_, childSchema = getSchemaForPathSegment(schema, "items")
			}
			findReferences(item, newPath, childSchema, requiredSegments, references)
		}
	}
}

// getSchemaForPathSegment returns the nested JSONSchemaProps for a given mapping key.
// - "properties" returns the current schema (so inner property names will be looked up in schema.Properties)
// - "items" returns the schema for array items (handles single Schema or first element of JSONSchemas)
// - normal property names return schema.Properties[name] (strips trailing "[index]" if present)
func getSchemaForPathSegment(schema *apiextensionsv1.JSONSchemaProps, key string) (bool, *apiextensionsv1.JSONSchemaProps) {
	if schema == nil {
		return false, nil
	}

	if key == "properties" {
		// The next level will contain actual property names; keep the current schema so those names can be looked up in schema.Properties
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
		// Return the items schema for arrays
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

	f.Line()

	// Always generate helper Requests function for all reference types
	generateMapFunc(f, crdKind, indexer)

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

		// Nil check if conditions
		nilCheckCondition := buildNilCheckConditions(fieldAccessPath, field.RequiredSegments)

		// Add check that the Name field is not empty
		condition := nilCheckCondition.Op("&&").Add(
			jen.Id(fieldAccessPath).Dot("Name").Op("!=").Lit(""),
		)

		// Generate: if <nil checks> && resource.Spec.<version>.GroupRef.Name != "" {
		//   keys = append(keys, types.NamespacedName{...}.String())
		// }
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

// buildNilCheckConditions creates a compound nil check condition for a field access path
// based on which segments are required (non-pointer) vs optional (pointer).
// Examples:
//   - fieldAccessPath: "cluster.Spec.V20250312.GroupRef"
//   - requiredSegments: [false, true, false] (for Spec, V20250312, GroupRef - excludes variable name)
func buildNilCheckConditions(fieldAccessPath string, requiredSegments []bool) *jen.Statement {
	segments := strings.Split(fieldAccessPath, ".")

	if len(requiredSegments) == 0 {
		return buildDotChain(segments).Op("!=").Nil()
	}

	if len(requiredSegments) != len(segments)-1 {
		return buildDotChain(segments).Op("!=").Nil()
	}

	var conditions *jen.Statement

	for i := 1; i < len(segments); i++ {
		requiredIndex := i - 1

		// Skip if segment is required, means it's a non-pointer struct
		if requiredSegments[requiredIndex] {
			continue
		}

		// Special case: "Spec" is always a non-pointer struct in Kubernetes resources
		if segments[i] == "Spec" {
			continue
		}

		// Build the path up to this segment
		pathSegments := segments[:i+1]
		nilCheck := buildDotChain(pathSegments).Op("!=").Nil()

		if conditions == nil {
			conditions = nilCheck
		} else {
			conditions = conditions.Op("&&").Add(nilCheck)
		}
	}

	// If no conditions were added (all segments required), check only the last field
	if conditions == nil {
		return buildDotChain(segments).Op("!=").Nil()
	}

	return conditions
}

// buildDotChain creates a jen.Statement for a dot-separated path
// For example: ["cluster", "Spec", "V20250312"] -> jen.Id("cluster").Dot("Spec").Dot("V20250312")
func buildDotChain(segments []string) *jen.Statement {
	if len(segments) == 0 {
		return jen.Null()
	}

	stmt := jen.Id(segments[0])
	for i := 1; i < len(segments); i++ {
		stmt = stmt.Dot(segments[i])
	}

	return stmt
}

func generateMapFunc(f *jen.File, crdKind string, indexer IndexerInfo) {
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
							Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1", fmt.Sprintf("%sList", crdKind)).
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

// TODO: UpdateIndexerRegistry needs to be reimplemented to work with the new kind-based indexer approach
// For now, indexers need to be manually registered in the indexer registry
