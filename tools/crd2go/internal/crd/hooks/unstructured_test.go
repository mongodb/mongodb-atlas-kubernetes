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
	"testing"

	"github.com/stretchr/testify/assert"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/internal/crd"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/internal/gotype"
)

func TestUnstructuredHookFn(t *testing.T) {
	trueVal := true
	falseVal := false

	tests := map[string]struct {
		crdType      *crd.CRDType
		expectedType *gotype.GoType
		expectedErr  error
	}{
		"not an unstructured type": {
			crdType: &crd.CRDType{
				Name: "NotAnUnstructured",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type: crd.OpenAPIString,
				},
			},
			expectedErr: fmt.Errorf("string is not unstructured (has 0 properties and x-preserve-unknown-fields is <nil>): %w", crd.ErrNotProcessed),
		},
		"object type but not unstructured (has properties)": {
			crdType: &crd.CRDType{
				Name: "ObjectButNotUnstructured",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type: crd.OpenAPIObject,
					Properties: map[string]apiextensionsv1.JSONSchemaProps{
						"prop1": {Type: crd.OpenAPIString},
					},
					XPreserveUnknownFields: &trueVal,
				},
			},
			expectedErr: fmt.Errorf("object is not unstructured (has 1 properties and x-preserve-unknown-fields is %v): %w", &trueVal, crd.ErrNotProcessed),
		},
		"object type but not unstructured (x-preserve-unknown-fields is false)": {
			crdType: &crd.CRDType{
				Name: "ObjectButNotUnstructured",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type:                   crd.OpenAPIObject,
					XPreserveUnknownFields: &falseVal,
				},
			},
			expectedErr: fmt.Errorf("object is not unstructured (has 0 properties and x-preserve-unknown-fields is %v): %w", &falseVal, crd.ErrNotProcessed),
		},
		"object type but not unstructured (x-preserve-unknown-fields is nil)": {
			crdType: &crd.CRDType{
				Name: "ObjectButNotUnstructured",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type: crd.OpenAPIObject,
				},
			},
			expectedErr: fmt.Errorf("object is not unstructured (has 0 properties and x-preserve-unknown-fields is %v): %w", nil, crd.ErrNotProcessed),
		},
		"valid unstructured type": {
			crdType: &crd.CRDType{
				Name: "ValidUnstructured",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type:                   crd.OpenAPIObject,
					XPreserveUnknownFields: &trueVal,
				},
			},
			expectedType: gotype.JSONType,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := UnstructuredHookFn(nil, nil, tt.crdType)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedType, got)
		})
	}
}
