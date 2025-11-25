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

package crds

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// AssertMajorVersion checks the given SDK majorVersion exists for the given
// Kubernetes kind and CRD version
func AssertMajorVersion(specVersion *apiextensionsv1.CustomResourceDefinitionVersion, kind string, majorVersion string) error {
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

// CompileCRDSchema compiles the given JSON properties as a type schema.
// Enables validating whether or not some JSON or unstructured data comforms
// with such type schema.
func CompileCRDSchema(openAPISchema *apiextensionsv1.JSONSchemaProps) (*jsonschema.Schema, error) {
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

// SelectVersion extracts the CRD version definition matching the given version.
func SelectVersion(spec *apiextensionsv1.CustomResourceDefinitionSpec, version string) *apiextensionsv1.CustomResourceDefinitionVersion {
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

// getOpenAPIProperties extracts the schema properties of a given CRD version.
// The kind is passed as a reference of the object this version is from.
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

// getSpecPropertiesFor takes the properties value of a given field of a kind's
// properties set
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
