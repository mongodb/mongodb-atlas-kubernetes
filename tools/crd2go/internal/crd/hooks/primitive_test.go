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

func TestPrimitiveHookFn(t *testing.T) {
	tests := map[string]struct {
		crdType      *crd.CRDType
		expectedType *gotype.GoType
		expectedErr  error
	}{
		"not a primitive type": {
			crdType: &crd.CRDType{
				Name: "NotAPrimitive",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type: crd.OpenAPIArray,
				},
			},
			expectedErr: fmt.Errorf("array is not a primitive type: %w", crd.ErrNotProcessed),
		},
		"string type": {
			crdType: &crd.CRDType{
				Name: "StringType",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type: crd.OpenAPIString,
				},
			},
			expectedType: gotype.NewPrimitive(gotype.StringKind, gotype.StringKind),
		},
		"integer type": {
			crdType: &crd.CRDType{
				Name: "IntegerType",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type: crd.OpenAPIInteger,
				},
			},
			expectedType: gotype.NewPrimitive(gotype.IntKind, gotype.IntKind),
		},
		"number type": {
			crdType: &crd.CRDType{
				Name: "NumberType",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type: crd.OpenAPINumber,
				},
			},
			expectedType: gotype.NewPrimitive(gotype.FloatKind, gotype.FloatKind),
		},
		"boolean type": {
			crdType: &crd.CRDType{
				Name: "BooleanType",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type: crd.OpenAPIBoolean,
				},
			},
			expectedType: gotype.NewPrimitive(gotype.BoolKind, gotype.BoolKind),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := PrimitiveHookFn(nil, nil, tt.crdType)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedType, got)
		})
	}
}
