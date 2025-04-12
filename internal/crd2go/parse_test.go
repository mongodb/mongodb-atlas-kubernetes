package crd2go_test

import (
	"embed"
	"testing"

	"github.com/josvazg/crd2go/internal/crd2go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//go:embed samples/*
var samples embed.FS

func TestParseCRD(t *testing.T) {
	f, err := samples.Open("samples/group.crd.yaml")
	require.NoError(t, err)
	defer f.Close()

	crd, err := crd2go.ParseCRD(f)
	require.NoError(t, err)
	want := expectedCRD()
	require.NotNil(t, crd)
	assert.Equal(t, want, crd)
}

func expectedCRD() *apiextensions.CustomResourceDefinition {
	return &apiextensions.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apiextensions.k8s.io/v1",
			Kind:       "CustomResourceDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "groups.atlas.generated.mongodb.com",
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{
			Group: "atlas.generated.mongodb.com",
			Scope: "Namespaced",
			Names: apiextensions.CustomResourceDefinitionNames{
				Kind:     "Group",
				ListKind: "GroupList",
			},
			Versions: []apiextensions.CustomResourceDefinitionVersion{
				{
					Name:    "v1",
					Served:  true,
					Storage: true,
					Schema: &apiextensions.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
							Type: "object",
							Properties: map[string]apiextensions.JSONSchemaProps{
								"spec": {
									Type: "object",
									Properties: map[string]apiextensions.JSONSchemaProps{
										"v20231115": {
											Type: "object",
											Properties: map[string]apiextensions.JSONSchemaProps{
												"parameters": {
													Type: "object",
													Properties: map[string]apiextensions.JSONSchemaProps{
														"projectOwnerId": {
															Type:        "string",
															Description: `Unique 24-hexadecimal digit string that identifies the MongoDB Cloud user to whom to grant the Project Owner role on the specified project. If you set this parameter, it overrides the default value of the oldest Organization Owner.`,
															MaxLength:   ptr(int64(24)),
															MinLength:   ptr(int64(24)),
														},
													},
												},
												"entry": {
													Type: "object",
													Properties: map[string]apiextensions.JSONSchemaProps{
														"name": {
															Description: `Human-readable label that identifies the project included in the MongoDB Cloud organization.`,
															MaxLength:   ptr(int64(64)),
															MinLength:   ptr(int64(1)),
															Type:        "string",
														},
														"orgId": {
															Description: `Unique 24-hexadecimal digit string that identifies the MongoDB Cloud organization to which the project belongs.`,
															MaxLength:   ptr(int64(24)),
															MinLength:   ptr(int64(24)),
															Type:        "string",
														},
														"regionUsageRestrictions": {
															Description: `Applies to Atlas for Government only.
	In Commercial Atlas, this field will be rejected in requests and missing in responses.
	This field sets restrictions on available regions in the project.
	| Value                             | Available Regions |
	|-----------------------------------|------------|
	| ` + "`COMMERCIAL_FEDRAMP_REGIONS_ONLY` | Only allows deployments in FedRAMP Moderate regions.|" + `
	| ` + "`GOV_REGIONS_ONLY`                | Only allows deployments in GovCloud regions.|",
															Type: "string",
														},
														"tags": {
															Description: "List that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the project.",
															Type:        "array",
															Items: &apiextensions.JSONSchemaPropsOrArray{
																Schema: &apiextensions.JSONSchemaProps{
																	Description: `Key-value pair that tags and categorizes a MongoDB Cloud organization, project, or cluster. For example, ` + "`environment : production`.",
																	Title:       "Resource Tag",
																	Type:        "object",
																	Properties: map[string]apiextensions.JSONSchemaProps{
																		"key": {
																			Description: `Constant that defines the set of the tag. For example, ` + "`environment` in the `environment : production` tag.",
																			MaxLength:   ptr(int64(255)),
																			MinLength:   ptr(int64(1)),
																			Type:        "string",
																		},
																		"value": {
																			Description: `Variable that belongs to the set of the tag. For example, ` + "`production` in the `environment : production` tag.",
																			MaxLength:   ptr(int64(255)),
																			MinLength:   ptr(int64(1)),
																			Type:        "string",
																		},
																	},
																	Required: []string{"key", "value"},
																},
															},
														},
														"withDefaultAlertsSettings": {
															Description: `Flag that indicates whether to create the project with default alert settings.`,
															Type:        "boolean",
														},
													},
													Required: []string{"name", "orgId"},
												},
											},
										},
									},
								},
								"status": {
									Type: "object",
									Properties: map[string]apiextensions.JSONSchemaProps{
										"conditions": {
											Description: "Represents the latest available observations of a resource's current state.",
											Type:        "array",
											Items: &apiextensions.JSONSchemaPropsOrArray{
												Schema: &apiextensions.JSONSchemaProps{
													Properties: map[string]apiextensions.JSONSchemaProps{
														"lastTransitionTime": {
															Description: "Last time the condition transitioned from one status to another.",
															Format:      "date-time",
															Type:        "string",
														},
														"message": {
															Description: "A human readable message indicating details about the transition.",
															Type:        "string",
														},
														"observedGeneration": {
															Description: "observedGeneration represents the .metadata.generation that the condition was set based upon.",
															Type:        "integer",
														},
														"reason": {
															Description: "The reason for the condition's last transition.",
															Type:        "string",
														},
														"status": {
															Description: "Status of the condition, one of True, False, Unknown.",
															Type:        "string",
														},
														"type": {
															Description: "Type of condition.",
															Type:        "string",
														},
													},
													Required: []string{"type", "status"},
													Type:     "object",
												},
											},
										},
										"v20231115": {
											Properties: map[string]apiextensions.JSONSchemaProps{
												"clusterCount": {
													Description: "Quantity of MongoDB Cloud clusters deployed in this project.",
													Format:      "int64",
													Type:        "integer",
												},
												"created": {
													Description: "Date and time when MongoDB Cloud created this project. This parameter expresses its value in the ISO 8601 timestamp format in UTC.",
													Format:      "date-time",
													Type:        "string",
												},
												"id": {
													Description: "Unique 24-hexadecimal digit string that identifies the MongoDB Cloud project.",
													MaxLength:   ptr(int64(24)),
													MinLength:   ptr(int64(24)),
													Type:        "string",
												},
											},
											Required: []string{"clusterCount", "created"},
											Type:     "object",
										},
									},
								},
							},
						},
					},
					Subresources: &apiextensions.CustomResourceSubresources{
						Status: &apiextensions.CustomResourceSubresourceStatus{},
					},
				},
			},
		},
		Status: apiextensions.CustomResourceDefinitionStatus{},
	}
}

func ptr[T any](t T) *T {
	return &t
}
