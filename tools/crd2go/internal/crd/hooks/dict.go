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

	"github.com/josvazg/crd2go/internal/crd"
	"github.com/josvazg/crd2go/internal/gotype"
)

func DictHookFn(td *gotype.TypeDict, hooks []crd.OpenAPI2GoHook, crdType *crd.CRDType) (*gotype.GoType, error) {
	if !isDict(crdType.Schema) {
		return nil, fmt.Errorf("%s is not a dictionary (additionalProperties is %v): %w",
			crdType.Schema.Type, crdType.Schema.AdditionalProperties, crd.ErrNotProcessed)
	}
	return fromOpenAPIDict(td, hooks, crdType)
}

// fromOpenAPIDict converts an OpenAPI dictionary to a GoType map
func fromOpenAPIDict(td *gotype.TypeDict, hooks []crd.OpenAPI2GoHook, crdType *crd.CRDType) (*gotype.GoType, error) {
	elemType := gotype.JSONType
	if crdType.Schema.AdditionalProperties.Schema != nil {
		var err error
		elemType, err = crd.FromOpenAPIType(td, hooks, &crd.CRDType{
			Name:    crdType.Name,
			Parents: crdType.Parents,
			Schema:  crdType.Schema.AdditionalProperties.Schema,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to check map value type: %w", err)
		}
	}
	return &gotype.GoType{Name: gotype.MapKind, Kind: gotype.MapKind, Element: elemType}, nil
}

func isDict(schema *apiextensionsv1.JSONSchemaProps) bool {
	return schema.AdditionalProperties != nil
}
