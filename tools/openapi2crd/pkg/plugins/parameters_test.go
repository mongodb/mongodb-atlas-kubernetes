package plugins

import (
	"errors"
	"fmt"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"k8s.io/utils/ptr"

	configv1alpha1 "tools/openapi2crd/pkg/apis/config/v1alpha1"
	"tools/openapi2crd/pkg/converter"
)

func TestParameterName(t *testing.T) {
	p := &Parameters{}
	assert.Equal(t, "parameters", p.Name())
}

func TestParameterProcess(t *testing.T) {
	tests := map[string]struct {
		request             *MappingProcessorRequest
		expectedVersionSpec apiextensions.JSONSchemaProps
		expectedErr         error
	}{
		"add parameter schema to the CRD": {
			request: groupMappingRequest(t, groupBaseCRDWithMajorVersion(t), entryInitialExtensionsSchema(t), parameterConverterMock(t)),
			expectedVersionSpec: apiextensions.JSONSchemaProps{
				Description: "The spec of the group resource for version v20250312.",
				Type:        "object",
				Properties: map[string]apiextensions.JSONSchemaProps{
					"projectOwnerId": {
						Description: "Unique 24-hexadecimal digit string that identifies the MongoDB Cloud user to whom to grant the Project Owner role on the specified project. If you set this parameter, it overrides the default value of the oldest Organization Owner.",
						Type:        "object",
						XValidations: apiextensions.ValidationRules{
							{
								Rule:    "self == oldSelf",
								Message: "projectOwnerId cannot be modified after creation",
							},
						},
					},
				},
				Required: []string{"projectOwnerId"},
			},
		},
		"missing path in OpenAPI spec": {
			request: &MappingProcessorRequest{
				CRD: groupBaseCRDWithMajorVersion(t),
				MappingConfig: &configv1alpha1.CRDMapping{
					MajorVersion: "v20250312",
					ParametersMapping: configv1alpha1.PropertyMapping{
						Path: configv1alpha1.PropertyPath{
							Name: "/api/atlas/v2/nonexistent",
							Verb: "post",
						},
					},
				},
				OpenAPISpec: &openapi3.T{
					Paths: openapi3.NewPaths(),
				},
			},
			expectedErr: fmt.Errorf("OpenAPI path %v does not exist", configv1alpha1.PropertyMapping{
				Path: configv1alpha1.PropertyPath{
					Name: "/api/atlas/v2/nonexistent",
					Verb: "post",
				},
			}),
			expectedVersionSpec: apiextensions.JSONSchemaProps{
				Description: "The spec of the group resource for version v20250312.",
				Type:        "object",
				Properties:  map[string]apiextensions.JSONSchemaProps{},
			},
		},
		"unsupported operation": {
			request: &MappingProcessorRequest{
				CRD: groupBaseCRDWithMajorVersion(t),
				MappingConfig: &configv1alpha1.CRDMapping{
					MajorVersion: "v20250312",
					ParametersMapping: configv1alpha1.PropertyMapping{
						Path: configv1alpha1.PropertyPath{
							Name: "/api/atlas/v2/groups",
							Verb: "delete",
						},
					},
				},
				OpenAPISpec: &openapi3.T{
					Paths: openapi3.NewPaths(
						openapi3.WithPath(
							"/api/atlas/v2/groups",
							&openapi3.PathItem{},
						),
					),
				},
			},
			expectedErr: errors.New("verb \"delete\" unsupported"),
			expectedVersionSpec: apiextensions.JSONSchemaProps{
				Description: "The spec of the group resource for version v20250312.",
				Type:        "object",
				Properties:  map[string]apiextensions.JSONSchemaProps{},
			},
		},
		"skipped parameters are not added to the CRD": {
			request: groupMappingRequestWithSkippedParameters(t, groupBaseCRDWithMajorVersion(t), entryInitialExtensionsSchema(t), parameterConverterMock(t)),
			expectedVersionSpec: apiextensions.JSONSchemaProps{
				Description: "The spec of the group resource for version v20250312.",
				Type:        "object",
				Properties: map[string]apiextensions.JSONSchemaProps{
					"projectOwnerId": {
						Description: "Unique 24-hexadecimal digit string that identifies the MongoDB Cloud user to whom to grant the Project Owner role on the specified project. If you set this parameter, it overrides the default value of the oldest Organization Owner.",
						Type:        "object",
						XValidations: apiextensions.ValidationRules{
							{
								Rule:    "self == oldSelf",
								Message: "projectOwnerId cannot be modified after creation",
							},
						},
					},
				},
				Required: []string{"projectOwnerId"},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := &Parameters{}
			err := p.Process(tt.request)
			assert.Equal(t, tt.expectedErr, err)
			versionSpec := tt.request.CRD.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[tt.request.MappingConfig.MajorVersion]
			assert.Equal(t, tt.expectedVersionSpec, versionSpec)
		})
	}
}

func groupMappingRequest(
	t *testing.T,
	crd *apiextensions.CustomResourceDefinition,
	extensionsSchema *openapi3.Schema,
	converterFunc converterFuncMock,
) *MappingProcessorRequest {
	t.Helper()

	return &MappingProcessorRequest{
		CRD: crd,
		MappingConfig: &configv1alpha1.CRDMapping{
			MajorVersion: "v20250312",
			OpenAPIRef: configv1alpha1.LocalObjectReference{
				Name: "v20250312",
			},
			ParametersMapping: configv1alpha1.PropertyMapping{
				Path: configv1alpha1.PropertyPath{
					Name: "/api/atlas/v2/groups",
					Verb: "post",
				},
			},
			EntryMapping: configv1alpha1.PropertyMapping{
				Schema: "Group",
				Filters: configv1alpha1.Filters{
					ReadWriteOnly: true,
				},
			},
			StatusMapping: configv1alpha1.PropertyMapping{
				Schema: "Group",
				Filters: configv1alpha1.Filters{
					ReadOnly:       true,
					SkipProperties: []string{"$.links"},
				},
			},
		},
		OpenAPISpec: &openapi3.T{
			Paths: openapi3.NewPaths(
				openapi3.WithPath(
					"/api/atlas/v2/groups",
					&openapi3.PathItem{
						Post: &openapi3.Operation{
							Parameters: openapi3.Parameters{
								&openapi3.ParameterRef{
									Value: &openapi3.Parameter{
										Name:        "projectOwnerId",
										In:          "query",
										Description: "Unique 24-hexadecimal digit string that identifies the MongoDB Cloud user to whom to grant the Project Owner role on the specified project. If you set this parameter, it overrides the default value of the oldest Organization Owner.",
										Schema: &openapi3.SchemaRef{
											Value: &openapi3.Schema{
												Type:    &openapi3.Types{"boolean"},
												Default: false,
											},
										},
									},
								},
							},
							RequestBody: &openapi3.RequestBodyRef{
								Value: &openapi3.RequestBody{
									Description: "Request body to create a new project.",
									Content: openapi3.Content{
										"application/vnd.atlas.2025-03-12+json": &openapi3.MediaType{
											Schema: &openapi3.SchemaRef{
												Ref: "#/components/schemas/Group",
											},
										},
									},
									Required: true,
								},
							},
							Responses: openapi3.NewResponses(
								openapi3.WithStatus(200, &openapi3.ResponseRef{
									Value: &openapi3.Response{
										Description: ptr.To("OK"),
										Content: openapi3.Content{
											"application/vnd.atlas.2025-03-12+json": &openapi3.MediaType{
												Schema: &openapi3.SchemaRef{
													Ref: "#/components/schemas/Group",
												},
											},
										},
									},
								}),
							),
						},
					},
				),
			),
			Components: &openapi3.Components{
				Schemas: map[string]*openapi3.SchemaRef{
					"Group": {
						Value: &openapi3.Schema{
							Type: &openapi3.Types{"object"},
							Properties: map[string]*openapi3.SchemaRef{
								"id": {
									Value: &openapi3.Schema{
										Type:        &openapi3.Types{"string"},
										Description: "Unique 24-hexadecimal digit string that identifies this project.",
									},
								},
								"name": {
									Value: &openapi3.Schema{
										Type:        &openapi3.Types{"string"},
										Description: "Human-readable label that identifies the project.",
									},
								},
							},
						},
					},
				},
			},
		},
		ExtensionsSchema: extensionsSchema,
		Converter: &dummyConverter{
			ConverterFunc: func(input converter.PropertyConvertInput) *apiextensions.JSONSchemaProps {
				if converterFunc == nil {
					return &apiextensions.JSONSchemaProps{}
				}

				return converterFunc(input)
			},
		},
	}
}

func groupMappingRequestWithSkippedParameters(
	t *testing.T,
	crd *apiextensions.CustomResourceDefinition,
	extensionsSchema *openapi3.Schema,
	converterFunc converterFuncMock,
) *MappingProcessorRequest {
	t.Helper()

	req := groupMappingRequest(t, crd, extensionsSchema, converterFunc)
	req.OpenAPISpec.Paths.Find("/api/atlas/v2/groups").Post.Parameters = append(
		req.OpenAPISpec.Paths.Find("/api/atlas/v2/groups").Post.Parameters,
		&openapi3.ParameterRef{
			Value: &openapi3.Parameter{
				Name:        "includeCount",
				In:          "query",
				Description: "A parameter that should be skipped.",
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{"boolean"},
					},
				},
			},
		},
		&openapi3.ParameterRef{
			Value: &openapi3.Parameter{
				Name:        "itemsPerPage",
				In:          "query",
				Description: "A parameter that should be skipped.",
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{"boolean"},
					},
				},
			},
		},
		&openapi3.ParameterRef{
			Value: &openapi3.Parameter{
				Name:        "pageNum",
				In:          "query",
				Description: "A parameter that should be skipped.",
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{"boolean"},
					},
				},
			},
		},
		&openapi3.ParameterRef{
			Value: &openapi3.Parameter{
				Name:        "envelope",
				In:          "query",
				Description: "A parameter that should be skipped.",
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{"boolean"},
					},
				},
			},
		},
		&openapi3.ParameterRef{
			Value: &openapi3.Parameter{
				Name:        "pretty",
				In:          "query",
				Description: "A parameter that should be skipped.",
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{"boolean"},
					},
				},
			},
		},
	)

	return req
}

func parameterConverterMock(t *testing.T) converterFuncMock {
	t.Helper()

	return func(input converter.PropertyConvertInput) *apiextensions.JSONSchemaProps {
		return &apiextensions.JSONSchemaProps{
			Description: "Unique 24-hexadecimal digit string that identifies the MongoDB Cloud user to whom to grant the Project Owner role on the specified project. If you set this parameter, it overrides the default value of the oldest Organization Owner.",
			Type:        "object",
			XValidations: apiextensions.ValidationRules{
				{
					Rule:    "self == oldSelf",
					Message: "projectOwnerId cannot be modified after creation",
				},
			},
		}
	}
}

type converterFuncMock func(input converter.PropertyConvertInput) *apiextensions.JSONSchemaProps

type dummyConverter struct {
	ConverterFunc func(input converter.PropertyConvertInput) *apiextensions.JSONSchemaProps
}

func (d *dummyConverter) Convert(input converter.PropertyConvertInput) *apiextensions.JSONSchemaProps {
	return d.ConverterFunc(input)
}
