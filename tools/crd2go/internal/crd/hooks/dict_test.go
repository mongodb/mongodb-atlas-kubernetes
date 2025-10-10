package hooks

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/internal/crd"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/internal/gotype"
)

func TestDictHookFn(t *testing.T) {
	tests := map[string]struct {
		hooks        []crd.OpenAPI2GoHook
		crdType      *crd.CRDType
		expectedType *gotype.GoType
		expectedErr  error
	}{
		"not a dictionary type": {
			crdType: &crd.CRDType{
				Name: "NotADict",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type: crd.OpenAPIString,
				},
			},
			expectedErr: fmt.Errorf("string is not a dictionary (additionalProperties is nil): %w", crd.ErrNotProcessed),
		},
		"dictionary with no specified value type": {
			crdType: &crd.CRDType{
				Name: "DictWithNoValueType",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type:                 crd.OpenAPIObject,
					AdditionalProperties: &apiextensionsv1.JSONSchemaPropsOrBool{Allows: true},
				},
			},
			expectedType: &gotype.GoType{Name: gotype.MapKind, Kind: gotype.MapKind, Element: gotype.JSONType},
		},
		"dictionary with unsupported value type": {
			hooks: []crd.OpenAPI2GoHook{},
			crdType: &crd.CRDType{
				Name: "DictWithValueType",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type: crd.OpenAPIObject,
					AdditionalProperties: &apiextensionsv1.JSONSchemaPropsOrBool{
						Schema: &apiextensionsv1.JSONSchemaProps{
							Type: crd.OpenAPIString,
						},
					},
				},
			},
			expectedErr: fmt.Errorf("failed to check map value type: %w", fmt.Errorf("unsupported Open API type \"DictWithValueType\"")),
		},
		"dictionary with specified value type": {
			hooks: []crd.OpenAPI2GoHook{
				hookMock(t, &gotype.GoType{Name: "string", Kind: gotype.StringKind}),
			},
			crdType: &crd.CRDType{
				Name: "DictWithValueType",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type: crd.OpenAPIObject,
					AdditionalProperties: &apiextensionsv1.JSONSchemaPropsOrBool{
						Schema: &apiextensionsv1.JSONSchemaProps{
							Type: crd.OpenAPIString,
						},
					},
				},
			},
			expectedType: &gotype.GoType{Name: gotype.MapKind, Kind: gotype.MapKind, Element: &gotype.GoType{Name: "string", Kind: gotype.StringKind}},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			td := gotype.NewTypeDict(nil, gotype.KnownTypes()...)
			got, err := DictHookFn(td, tt.hooks, tt.crdType)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedType, got)
		})
	}
}
