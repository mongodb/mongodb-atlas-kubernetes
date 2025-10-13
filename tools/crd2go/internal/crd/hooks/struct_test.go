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

func TestStructHookFn(t *testing.T) {
	typeToRename := gotype.GoType{
		Name:    "Array",
		Kind:    gotype.ArrayKind,
		Element: &gotype.GoType{Name: "Struct", Kind: gotype.StructKind},
	}
	fieldToRename := gotype.GoField{
		Name:     "Array",
		Key:      "Array",
		GoType:   &typeToRename,
		Required: false,
	}

	tests := map[string]struct {
		hooks         []crd.OpenAPI2GoHook
		crdType       *crd.CRDType
		expectedType  *gotype.GoType
		expectedError error
	}{
		"not an object type": {
			crdType: &crd.CRDType{
				Name: "NotAnObject",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type: crd.OpenAPIString,
				},
			},
			expectedError: fmt.Errorf("string is not an object: %w", crd.ErrNotProcessed),
		},
		"object type with unspported properties": {
			crdType: &crd.CRDType{
				Name: "ObjectWithUnsupportedProperties",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type:       crd.OpenAPIObject,
					Properties: map[string]apiextensionsv1.JSONSchemaProps{"prop": {Type: crd.OpenAPIString}},
				},
			},
			expectedError: fmt.Errorf("failed to parse prop type: %w", fmt.Errorf("unsupported Open API type \"prop\"")),
		},
		"failed to rename field": {
			hooks: []crd.OpenAPI2GoHook{
				hookMock(t, &typeToRename),
			},
			crdType: &crd.CRDType{
				Name: "Object",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type: crd.OpenAPIObject,
					Properties: map[string]apiextensionsv1.JSONSchemaProps{
						"Array": {Type: crd.OpenAPIArray},
					},
				},
				Parents: []string{"Parent"},
			},
			expectedError: fmt.Errorf(
				"failed to rename field %v: %w",
				&fieldToRename,
				fmt.Errorf(
					"failed to rename field type: %w",
					fmt.Errorf("failed to find a free type name for type %v", &typeToRename),
				),
			),
		},
		"object type with no properties": {
			crdType: &crd.CRDType{
				Name: "EmptyObject",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type: crd.OpenAPIObject,
				},
			},
			expectedType: &gotype.GoType{Name: "EmptyObject", Kind: gotype.StructKind, Fields: []*gotype.GoField{}},
		},
		"object type with properties": {
			hooks: []crd.OpenAPI2GoHook{
				hookMock(t, &gotype.GoType{Name: "string", Kind: gotype.StringKind}),
			},
			crdType: &crd.CRDType{
				Name: "ObjectWithProperties",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type: crd.OpenAPIObject,
					Properties: map[string]apiextensionsv1.JSONSchemaProps{
						"prop": {Type: crd.OpenAPIString},
					},
				},
			},
			expectedType: &gotype.GoType{
				Name: "ObjectWithProperties",
				Kind: gotype.StructKind,
				Fields: []*gotype.GoField{
					{Name: "Prop", Key: "prop", GoType: &gotype.GoType{Name: "string", Kind: gotype.StringKind}},
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			td := gotype.NewTypeDict(map[string]string{"Struct": "Data"}, gotype.KnownTypes()...)
			td.Add(&gotype.GoType{Name: "Data", Kind: gotype.StringKind})
			td.Add(&gotype.GoType{Name: "ObjectData", Kind: gotype.StringKind})
			td.Add(&gotype.GoType{Name: "ParentObjectData", Kind: gotype.StringKind})
			got, err := StructHookFn(td, tt.hooks, tt.crdType)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedType, got)
		})
	}
}
