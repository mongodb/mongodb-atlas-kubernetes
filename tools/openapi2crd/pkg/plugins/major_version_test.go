package plugins

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

func TestMajorVersionName(t *testing.T) {
	p := &MajorVersion{}
	assert.Equal(t, "major_version", p.Name())
}

func TestMajorVersionProcess(t *testing.T) {
	tests := map[string]struct {
		request             *MappingProcessorRequest
		expectedVersionSpec apiextensions.JSONSchemaProps
		expectedErr         error
	}{
		"add major version schema to the CRD": {
			request:             groupMappingRequest(t, groupBaseCRD(t), majorVersionInitialExtensionsSchema(t), nil),
			expectedVersionSpec: majorVersionSchema(t),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := &MajorVersion{}
			err := p.Process(tt.request)
			assert.Equal(t, tt.expectedErr, err)
			versionSpec := tt.request.CRD.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[tt.request.MappingConfig.MajorVersion]
			assert.Equal(t, tt.expectedVersionSpec, versionSpec)
		})
	}
}

func majorVersionSchema(t *testing.T) apiextensions.JSONSchemaProps {
	t.Helper()

	return apiextensions.JSONSchemaProps{
		Type:        "object",
		Description: "The spec of the group resource for version v20250312.",
		Properties:  map[string]apiextensions.JSONSchemaProps{},
	}
}

func majorVersionInitialExtensionsSchema(t *testing.T) *openapi3.Schema {
	t.Helper()

	return &openapi3.Schema{
		Type: &openapi3.Types{"object"},
		Properties: map[string]*openapi3.SchemaRef{
			"spec": {
				Value: &openapi3.Schema{
					Type: &openapi3.Types{"object"},
					Properties: map[string]*openapi3.SchemaRef{
						"v20250312": {
							Value: &openapi3.Schema{
								Type: &openapi3.Types{"object"},
								Properties: map[string]*openapi3.SchemaRef{
									"entry": {Value: &openapi3.Schema{}},
								},
							},
						},
					},
				},
			},
		},
	}
}
