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

package hooks

import (
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/internal/crd"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/internal/gotype"
)

func UnstructuredHookFn(_ *gotype.TypeDict, _ []crd.OpenAPI2GoHook, crdType *crd.CRDType) (*gotype.GoType, error) {
	if crdType.Schema.Type != crd.OpenAPIObject || !isUnstructured(crdType.Schema) {
		return nil, fmt.Errorf(
			"%s is not unstructured (has %d properties and x-preserve-unknown-fields is %v): %w",
			crdType.Schema.Type,
			len(crdType.Schema.Properties),
			crdType.Schema.XPreserveUnknownFields,
			crd.ErrNotProcessed,
		)
	}

	return gotype.JSONType, nil
}

func isUnstructured(schema *apiextensionsv1.JSONSchemaProps) bool {
	return len(schema.Properties) == 0 && schema.XPreserveUnknownFields != nil && *schema.XPreserveUnknownFields
}
