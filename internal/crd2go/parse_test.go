package crd2go_test

import (
	"embed"
	"testing"

	"github.com/josvazg/crd2go/internal/crd2go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed samples/*
var samples embed.FS

func TestParseCRD(t *testing.T) {
	f, err := samples.Open("samples/group.crd.yaml")
	require.NoError(t, err)
	defer f.Close()

	crd, err := crd2go.ParseCRD(f)
	require.NoError(t, err)
	want := &crd2go.CRD{
		APIVersion: "apiextensions.k8s.io/v1",
		Kind:       "CustomResourceDefinition",
		Metadata: crd2go.Metadata{
			Name: "groups.atlas.generated.mongodb.com",
		},
		Spec: crd2go.Spec{
			Group: "atlas.generated.mongodb.com",
			Scope: "Namespaced",
			Names: crd2go.SpecNames{
				Kind:     "Group",
				ListKind: "GroupList",
			},
			Versions: []crd2go.Version{
				{
					Name:    "v1",
					Served:  true,
					Storage: true,
					Schema: crd2go.VersionSchema{
						OpenAPIV3Schema: crd2go.OpenAPISchema{
							Type: "object",
							Properties: map[string]crd2go.OpenAPISchema{
								"spec": {
									Type: "object",
									Properties: map[string]crd2go.OpenAPISchema{
										"v20231115": {
											Type: "object",
											Properties: map[string]crd2go.OpenAPISchema{
												"parameters": {
													Type: "object",
													Properties: map[string]crd2go.OpenAPISchema{
														"projectOwnerId": {
															Type:        "string",
															Description: ptr(`Unique 24-hexadecimal digit string that identifies the MongoDB Cloud user to whom to grant the Project Owner role on the specified project. If you set this parameter, it overrides the default value of the oldest Organization Owner.`),
															MaxLength:   ptr(24),
															MinLength:   ptr(24),
														},
													},
												},
												"entry": {
													Type: "object",
													Properties: map[string]crd2go.OpenAPISchema{
														"name": {
															Description: ptr(`Human-readable label that identifies the project included in the MongoDB Cloud organization.`),
															MaxLength:   ptr(64),
															MinLength:   ptr(1),
															Type:        "string",
														},
														"orgId": {
															Description: ptr(`Unique 24-hexadecimal digit string that identifies the MongoDB Cloud organization to which the project belongs.`),
															MaxLength:   ptr(24),
															MinLength:   ptr(24),
															Type:        "string",
														},
														"regionUsageRestrictions": {
															Description: ptr(`Applies to Atlas for Government only.

In Commercial Atlas, this field will be rejected in requests and missing in responses.

This field sets restrictions on available regions in the project.

| Value                             | Available Regions |
|-----------------------------------|------------|
| ` + "`COMMERCIAL_FEDRAMP_REGIONS_ONLY` | Only allows deployments in FedRAMP Moderate regions.|" + `
| ` + "`GOV_REGIONS_ONLY`                | Only allows deployments in GovCloud regions.|"),

															Type: "string",
														},
														"tags": {
															Description: ptr("List that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the project."),
															Type:        "array",
															Items: &crd2go.OpenAPISchema{
																Description: ptr(`Key-value pair that tags and categorizes a MongoDB Cloud organization, project, or cluster. For example, ` + "`environment : production`."),
																Title:       ptr("Resource Tag"),
																Type:        "object",
																Properties: map[string]crd2go.OpenAPISchema{
																	"key": {
																		Description: ptr(`Constant that defines the set of the tag. For example, ` + "`environment` in the `environment : production` tag."),
																		MaxLength:   ptr(255),
																		MinLength:   ptr(1),
																		Type:        "string",
																	},
																	"value": {
																		Description: ptr(`Variable that belongs to the set of the tag. For example, ` + "`production` in the `environment : production` tag."),
																		MaxLength:   ptr(255),
																		MinLength:   ptr(1),
																		Type:        "string",
																	},
																},
																Required: []string{"key", "value"},
															},
														},
														"withDefaultAlertsSettings": {
															Description: ptr(`Flag that indicates whether to create the project with default alert settings.`),
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
									Properties: map[string]crd2go.OpenAPISchema{
										"conditions": {
											Description: ptr("Represents the latest available observations of a resource's current state."),
											Type:        "array",
											Items: &crd2go.OpenAPISchema{
												Properties: map[string]crd2go.OpenAPISchema{
													"lastTransitionTime": {
														Description: ptr("Last time the condition transitioned from one status to another."),
														Format:      ptr("date-time"),
														Type:        "string",
													},
													"message": {
														Description: ptr("A human readable message indicating details about the transition."),
														Type:        "string",
													},
													"observedGeneration": {
														Description: ptr("observedGeneration represents the .metadata.generation that the condition was set based upon."),
														Type:        "integer",
													},
													"reason": {
														Description: ptr("The reason for the condition's last transition."),
														Type:        "string",
													},
													"status": {
														Description: ptr("Status of the condition, one of True, False, Unknown."),
														Type:        "string",
													},
													"type": {
														Description: ptr("Type of condition."),
														Type:        "string",
													},
												},
												Required: []string{"type", "status"},
												Type:     "object",
											},
										},
										"v20231115": {
											Properties: map[string]crd2go.OpenAPISchema{
												"clusterCount": {
													Description: ptr("Quantity of MongoDB Cloud clusters deployed in this project."),
													Format:      ptr("int64"),
													Type:        "integer",
												},
												"created": {
													Description: ptr("Date and time when MongoDB Cloud created this project. This parameter expresses its value in the ISO 8601 timestamp format in UTC."),
													Format:      ptr("date-time"),
													Type:        "string",
												},
												"id": {
													Description: ptr("Unique 24-hexadecimal digit string that identifies the MongoDB Cloud project."),
													MaxLength:   ptr(24),
													MinLength:   ptr(24),
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
					Subresources: map[string]any{
						"status": map[string]any{},
					},
				},
			},
		},
	}
	require.NotNil(t, crd)
	assert.Equal(t, want, crd)
}

func ptr[T any](t T) *T {
	return &t
}
