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
