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
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/autogen/translate/crds"
)

// translator implements Translator to translate from a given CRD to and from
// a given SDK version using the same upstream OpenAPI schema
type translator struct {
	majorVersion string
	jsonSchema   *jsonschema.Schema
	annotations  map[string]string
}

// Annotation returns the annotation value from the pinned translated schema
// CRD version
func (tr *translator) Annotation(annotation string) string {
	return tr.annotations[annotation]
}

// MajorVersion returns the CRD pinned version
func (tr *translator) MajorVersion() string {
	return tr.majorVersion
}

// Validate woudl return any errors of teh given unstructured object against the
// pinned schema version being translated, or nil if the object is compliant
func (tr *translator) Validate(unstructuredObj map[string]any) error {
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
	specVersion := crds.SelectVersion(&crd.Spec, crdVersion)
	if err := crds.AssertMajorVersion(specVersion, crd.Spec.Names.Kind, majorVersion); err != nil {
		return nil, fmt.Errorf("failed to assert major version %s in CRD: %w", majorVersion, err)
	}
	schema, err := crds.CompileCRDSchema(specVersion.Schema.OpenAPIV3Schema)
	if err != nil {
		return nil, fmt.Errorf("failed to compile schema: %w", err)
	}
	return &translator{
		majorVersion: majorVersion,
		jsonSchema:   schema,
		annotations:  crd.Annotations,
	}, nil
}
