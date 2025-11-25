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

package crapi

import (
	"encoding/json"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/santhosh-tekuri/jsonschema/v5"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/yaml"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/crapi/crds"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/crapi/refs"
)

// translator implements Translator to translate from a given CRD to and from
// a given SDK version using the same upstream OpenAPI schema
type translator struct {
	majorVersion     string
	validationSchema *jsonschema.Schema
	mappingSchema    *openapi3.SchemaRef
}

// Annotation returns the annotation value from the pinned translated schema
// CRD version
// FindAllMappings traverses the mapping schema and returns all found extensions.
func (tr *translator) Mappings() ([]*refs.Mapping, error) {
	noMappings := []*refs.Mapping{}
	if tr.mappingSchema == nil || tr.mappingSchema.Value == nil {
		return noMappings, nil
	}
	spec := propertyValueOrNil(tr.mappingSchema.Value, "spec")
	if spec == nil {
		return noMappings, nil
	}
	version := propertyValueOrNil(spec, tr.majorVersion)
	if version == nil {
		return noMappings, nil
	}
	return refs.FindMappings(version, []string{"spec", tr.majorVersion})
}

// MajorVersion returns the CRD pinned version
func (tr *translator) MajorVersion() string {
	return tr.majorVersion
}

// Validate would return any errors of the given unstructured object against the
// pinned schema version being translated, or nil if the object is compliant
func (tr *translator) Validate(unstructuredObj map[string]any) error {
	// This correctly uses the validator from crds.CompileCRDSchema
	if err := tr.validationSchema.Validate(unstructuredObj); err != nil {
		return fmt.Errorf("object validation failed against CRD schema: %w", err)
	}
	return nil
}

// NewTranslator creates a translator for a particular CRD version. It is also
// locked into a particular API majorVersion.
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
	specVersion := crds.SelectVersion(&crd.Spec, crdVersion)
	if err := crds.AssertMajorVersion(specVersion, crd.Spec.Names.Kind, majorVersion); err != nil {
		return nil, fmt.Errorf("failed to assert major version %s in CRD: %w", majorVersion, err)
	}
	validationSchema, err := crds.CompileCRDSchema(specVersion.Schema.OpenAPIV3Schema)
	if err != nil {
		return nil, fmt.Errorf("failed to compile schema for validation: %w", err)
	}
	var mappingSchema openapi3.Schema
	mappingString, ok := crd.Annotations["api-mappings"]
	if ok && mappingString != "" {
		jsonBytes, err := yaml.YAMLToJSON([]byte(mappingString))
		if err != nil {
			return nil, fmt.Errorf("failed to convert 'api-mappings' YAML to JSON: %w", err)
		}
		if err := json.Unmarshal(jsonBytes, &mappingSchema); err != nil {
			return nil, fmt.Errorf("failed to unmarshal 'api-mappings' JSON into schema: %w", err)
		}
	}

	return &translator{
		majorVersion:     majorVersion,
		validationSchema: validationSchema,
		mappingSchema:    &openapi3.SchemaRef{Value: &mappingSchema},
	}, nil
}
