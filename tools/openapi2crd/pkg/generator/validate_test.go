package generator

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateCRD(t *testing.T) {
	tests := map[string]struct {
		crd         *apiextensions.CustomResourceDefinition
		expectedErr error
	}{
		"valid CRD": {
			crd: &apiextensions.CustomResourceDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name: "examples.test.com",
				},
				Spec: apiextensions.CustomResourceDefinitionSpec{
					Group: "test.com",
					Names: apiextensions.CustomResourceDefinitionNames{
						Plural:   "examples",
						Singular: "example",
						Kind:     "Example",
						ListKind: "ExampleList",
					},
					Versions: []apiextensions.CustomResourceDefinitionVersion{
						{
							Name:    "v1",
							Served:  true,
							Storage: true,
						},
					},
					Validation: &apiextensions.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
							Type: "object",
						},
					},
					Scope: apiextensions.NamespaceScoped,
				},
				Status: apiextensions.CustomResourceDefinitionStatus{
					StoredVersions: []string{"v1"},
				},
			},
			expectedErr: nil,
		},
		"invalid CRD": {
			crd: &apiextensions.CustomResourceDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name: "examples.test.com",
				},
				Spec: apiextensions.CustomResourceDefinitionSpec{
					Group: "test.com",
					Names: apiextensions.CustomResourceDefinitionNames{
						Plural:   "wrongs",
						Singular: "wrong",
						Kind:     "Wrong",
						ListKind: "WrongList",
					},
					Versions: []apiextensions.CustomResourceDefinitionVersion{
						{
							Name:    "v1",
							Served:  true,
							Storage: true,
						},
					},
					Validation: &apiextensions.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
							Type: "object",
						},
					},
					Scope: apiextensions.NamespaceScoped,
				},
				Status: apiextensions.CustomResourceDefinitionStatus{
					StoredVersions: []string{"v1"},
				},
			},
			expectedErr: fmt.Errorf(
				"error validating CRD %v: %w",
				"examples.test.com",
				field.ErrorList{&field.Error{
					Type:     field.ErrorTypeInvalid,
					Field:    "metadata.name",
					BadValue: "examples.test.com",
					Detail:   "must be spec.names.plural+\".\"+spec.group"},
				}.ToAggregate(),
			),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := ValidateCRD(context.Background(), tt.crd)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
