package plugins

import (
	"errors"
	"testing"
	configv1alpha1 "tools/openapi2crd/pkg/apis/config/v1alpha1"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
)

func TestReferenceMetadataName(t *testing.T) {
	p := &ReferencesMetadata{}
	assert.Equal(t, "reference_metadata", p.Name())
}

func TestReferenceMetadataProcess(t *testing.T) {
	tests := map[string]struct {
		request            *ExtensionProcessorRequest
		expectedExtensions map[string]interface{}
		expectedErr        error
	}{
		"do nothing when no references": {
			request: &ExtensionProcessorRequest{
				MappingConfig: &configv1alpha1.CRDMapping{
					MajorVersion: "v20250312",
				},
				ExtensionsSchema: &openapi3.Schema{
					Properties: map[string]*openapi3.SchemaRef{
						"spec": {
							Value: &openapi3.Schema{
								Type: &openapi3.Types{"object"},
								Properties: map[string]*openapi3.SchemaRef{
									"v20250312": {
										Value: &openapi3.Schema{
											Type:       &openapi3.Types{"object"},
											Properties: map[string]*openapi3.SchemaRef{},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedExtensions: nil,
			expectedErr:        nil,
		},
		"add reference metadata": {
			request: &ExtensionProcessorRequest{
				MappingConfig: &configv1alpha1.CRDMapping{
					MajorVersion: "v20250312",
					ParametersMapping: configv1alpha1.PropertyMapping{
						References: []configv1alpha1.Reference{
							{
								Name:     "myRef",
								Property: "spec.myRef",
								Target: configv1alpha1.Target{
									Type: configv1alpha1.Type{
										Kind:     "MyKind",
										Group:    "mygroup.example.com",
										Version:  "v1",
										Resource: "myresources",
									},
									Properties: []string{"status.id"},
								},
							},
						},
					},
				},
				ExtensionsSchema: &openapi3.Schema{
					Properties: map[string]*openapi3.SchemaRef{
						"spec": {
							Value: &openapi3.Schema{
								Type: &openapi3.Types{"object"},
								Properties: map[string]*openapi3.SchemaRef{
									"v20250312": {
										Value: &openapi3.Schema{
											Type:       &openapi3.Types{"object"},
											Properties: map[string]*openapi3.SchemaRef{},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedExtensions: map[string]interface{}{
				"x-kubernetes-mapping": map[string]interface{}{
					"type": map[string]interface{}{
						"kind":     "MyKind",
						"group":    "mygroup.example.com",
						"version":  "v1",
						"resource": "myresources",
					},
					"nameSelector": ".name",
					"properties": []string{
						"status.id",
					},
				},
				"x-openapi-mapping": map[string]interface{}{
					"property": "spec.myRef",
				},
			},
			expectedErr: nil,
		},
		"error when reference target has no properties": {
			request: &ExtensionProcessorRequest{
				MappingConfig: &configv1alpha1.CRDMapping{
					MajorVersion: "v20250312",
					ParametersMapping: configv1alpha1.PropertyMapping{
						References: []configv1alpha1.Reference{
							{
								Name:     "myRef",
								Property: "spec.myRef",
								Target: configv1alpha1.Target{
									Type: configv1alpha1.Type{
										Kind:     "MyKind",
										Group:    "mygroup.example.com",
										Version:  "v1",
										Resource: "myresources",
									},
									Properties: []string{},
								},
							},
						},
					},
				},
				ExtensionsSchema: &openapi3.Schema{
					Properties: map[string]*openapi3.SchemaRef{
						"spec": {
							Value: &openapi3.Schema{
								Type: &openapi3.Types{"object"},
								Properties: map[string]*openapi3.SchemaRef{
									"v20250312": {
										Value: &openapi3.Schema{
											Type:       &openapi3.Types{"object"},
											Properties: map[string]*openapi3.SchemaRef{},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedExtensions: nil,
			expectedErr:        errors.New("reference target must have at least one property defined"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := &ReferencesMetadata{}
			err := p.Process(tt.request)
			assert.Equal(t, tt.expectedErr, err)
			if tt.expectedExtensions != nil {
				assert.Equal(t, tt.expectedExtensions, tt.request.ExtensionsSchema.Properties["spec"].Value.Properties[tt.request.MappingConfig.MajorVersion].Value.Properties["myRef"].Value.Extensions)
			} else {
				assert.Empty(t, tt.request.ExtensionsSchema.Properties["spec"].Value.Properties[tt.request.MappingConfig.MajorVersion].Value.Properties)
			}
		})
	}
}
