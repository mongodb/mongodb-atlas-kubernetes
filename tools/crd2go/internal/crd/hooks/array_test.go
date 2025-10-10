package hooks

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/internal/crd"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/internal/gotype"
)

func TestArrayHookFn(t *testing.T) {
	tests := map[string]struct {
		hooks        []crd.OpenAPI2GoHook
		crdType      *crd.CRDType
		expectedType *gotype.GoType
		expectedErr  error
	}{
		"not an array": {
			hooks: []crd.OpenAPI2GoHook{},
			crdType: &crd.CRDType{
				Name: "NotAnArray",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type: crd.OpenAPIString,
				},
			},
			expectedErr: fmt.Errorf("string is not an array: %w", crd.ErrNotProcessed),
		},
		"array with no items": {
			hooks: []crd.OpenAPI2GoHook{},
			crdType: &crd.CRDType{
				Name: "ArrayWithNoItems",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type: crd.OpenAPIArray,
				},
			},
			expectedErr: fmt.Errorf("array ArrayWithNoItems has no items"),
		},
		"array with no items schema": {
			hooks: []crd.OpenAPI2GoHook{},
			crdType: &crd.CRDType{
				Name: "ArrayWithNoItemsSchema",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type:  crd.OpenAPIArray,
					Items: &apiextensionsv1.JSONSchemaPropsOrArray{},
				},
			},
			expectedErr: fmt.Errorf("array ArrayWithNoItemsSchema has no items schema"),
		},
		"failed to parse element type": {
			hooks: []crd.OpenAPI2GoHook{},
			crdType: &crd.CRDType{
				Name: "ArrayWithInvalidElementType",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type: crd.OpenAPIArray,
					Items: &apiextensionsv1.JSONSchemaPropsOrArray{
						Schema: &apiextensionsv1.JSONSchemaProps{
							Type: crd.OpenAPIString,
						},
					},
				},
			},
			expectedErr: fmt.Errorf(
				"failed to parse array %s element type: %w",
				"ArrayWithInvalidElementType",
				fmt.Errorf("unsupported Open API type %q", "ArrayWithInvalidElementType"),
			),
		},
		"failed to rename element type": {
			hooks: []crd.OpenAPI2GoHook{
				hookMock(t, &gotype.GoType{Name: "Object", Kind: gotype.StructKind}),
			},
			crdType: &crd.CRDType{
				Name: "Array",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type: crd.OpenAPIArray,
					Items: &apiextensionsv1.JSONSchemaPropsOrArray{
						Schema: &apiextensionsv1.JSONSchemaProps{
							Type: crd.OpenAPIObject,
						},
					},
				},
				Parents: []string{"Parent"},
			},
			expectedErr: fmt.Errorf(
				"failed to rename element type under %s: %w",
				"Array",
				fmt.Errorf("failed to find a free type name for type %v", &gotype.GoType{Name: "Data", Kind: gotype.StructKind}),
			),
		},
		"array of strings": {
			hooks: []crd.OpenAPI2GoHook{
				hookMock(t, &gotype.GoType{Kind: gotype.StringKind}),
			},
			crdType: &crd.CRDType{
				Name: "ArrayOfStrings",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type: crd.OpenAPIArray,
					Items: &apiextensionsv1.JSONSchemaPropsOrArray{
						Schema: &apiextensionsv1.JSONSchemaProps{
							Type: crd.OpenAPIString,
						},
					},
				},
			},
			expectedType: gotype.NewArray(&gotype.GoType{Kind: gotype.StringKind}),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			td := gotype.NewTypeDict(map[string]string{"Object": "Data"}, gotype.KnownTypes()...)
			td.Add(&gotype.GoType{Name: "Data", Kind: gotype.ArrayKind})
			td.Add(&gotype.GoType{Name: "ParentData", Kind: gotype.ArrayKind})

			got, err := ArrayHookFn(td, tt.hooks, tt.crdType)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedType, got)
		})
	}
}

func hookMock(t *testing.T, goType *gotype.GoType) crd.OpenAPI2GoHook {
	t.Helper()

	return func(_ *gotype.TypeDict, _ []crd.OpenAPI2GoHook, _ *crd.CRDType) (*gotype.GoType, error) {
		return goType, nil
	}
}
