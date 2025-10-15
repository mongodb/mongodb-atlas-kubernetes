// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package translate

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// crdTranslator implements Translator to translate from a given CRD to and from
// a given SDK version using the same upstream OpenAPI schema
type crdTranslator struct {
	majorVersion string
	jsonSchema   *jsonschema.Schema
	annotations  map[string]string
}

func (tr *crdTranslator) Annotation(annotation string) string {
	return tr.annotations[annotation]
}

func (tr *crdTranslator) MajorVersion() string {
	return tr.majorVersion
}

func (tr *crdTranslator) Validate(unstructuredObj map[string]any) error {
	if err := tr.jsonSchema.Validate(unstructuredObj); err != nil {
		return fmt.Errorf("object validation failed against CRD schema: %w", err)
	}
	return nil
}

// NewTranslator creates a translator for a particular CRD and major version pairs,
// and with a particular set of known Kubernetes object dependencies.
//
// Given the following example resource:
//
//	apiVersion: atlas.generated.mongodb.com/v1
//	kind: SearchIndex
//	metadata:
//	  name: search-index
//	spec:
//	  v20250312:
//
// In the above case crdVersion is "v1" and majorVersion is "v20250312".
func NewTranslator(crd *apiextensionsv1.CustomResourceDefinition, crdVersion string, majorVersion string) (Translator, error) {
	specVersion := selectVersion(&crd.Spec, crdVersion)
	if err := assertMajorVersion(specVersion, crd.Spec.Names.Kind, majorVersion); err != nil {
		return nil, fmt.Errorf("failed to assert major version %s in CRD: %w", majorVersion, err)
	}
	schema, err := compileCRDSchema(specVersion.Schema.OpenAPIV3Schema)
	if err != nil {
		return nil, fmt.Errorf("failed to compile schema: %w", err)
	}
	return &crdTranslator{
		majorVersion: majorVersion,
		jsonSchema:   schema,
		annotations:  crd.Annotations,
	}, nil
}

func assertMajorVersion(specVersion *apiextensionsv1.CustomResourceDefinitionVersion, kind string, majorVersion string) error {
	props, err := getOpenAPIProperties(kind, specVersion)
	if err != nil {
		return fmt.Errorf("failed to enumerate CRD schema properties: %w", err)
	}
	specProps, err := getSpecPropertiesFor(kind, props, "spec")
	if err != nil {
		return fmt.Errorf("failed to enumerate CRD spec properties: %w", err)
	}
	_, ok := specProps[majorVersion]
	if !ok {
		return fmt.Errorf("failed to match the CRD spec version %q in schema", majorVersion)
	}
	return nil
}

func compileCRDSchema(openAPISchema *apiextensionsv1.JSONSchemaProps) (*jsonschema.Schema, error) {
	schemaBytes, err := json.Marshal(openAPISchema)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CRD schema to JSON: %w", err)
	}
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", bytes.NewReader(schemaBytes)); err != nil {
		return nil, fmt.Errorf("failed to add schema resource: %w", err)
	}
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return nil, fmt.Errorf("failed to compile schema: %w", err)
	}
	return schema, nil
}

// selectVersion returns the version from the CRD spec that matches the given version string
func selectVersion(spec *apiextensionsv1.CustomResourceDefinitionSpec, version string) *apiextensionsv1.CustomResourceDefinitionVersion {
	if len(spec.Versions) == 0 {
		return nil
	}
	if version == "" {
		return &spec.Versions[0]
	}
	for _, v := range spec.Versions {
		if v.Name == version {
			return &v
		}
	}
	return nil
}

func getOpenAPIProperties(kind string, version *apiextensionsv1.CustomResourceDefinitionVersion) (map[string]apiextensionsv1.JSONSchemaProps, error) {
	if version == nil {
		return nil, fmt.Errorf("missing version (nil) from %v spec", kind)
	}
	if version.Schema == nil {
		return nil, fmt.Errorf("missing version schema from %v spec", kind)
	}
	if version.Schema.OpenAPIV3Schema == nil {
		return nil, fmt.Errorf("missing version OpenAPI Schema from %v spec", kind)
	}
	if version.Schema.OpenAPIV3Schema.Properties == nil {
		return nil, fmt.Errorf("missing version OpenAPI Properties from %v spec", kind)
	}
	return version.Schema.OpenAPIV3Schema.Properties, nil
}

func getSpecPropertiesFor(kind string, props map[string]apiextensionsv1.JSONSchemaProps, field string) (map[string]apiextensionsv1.JSONSchemaProps, error) {
	prop, ok := props[field]
	if !ok {
		return nil, fmt.Errorf("kind %q spec is missing field %q on", kind, field)
	}
	if prop.Type != "object" {
		return nil, fmt.Errorf("kind %q field %q expected to be object but is %v", kind, field, prop.Type)
	}
	return prop.Properties, nil
}
