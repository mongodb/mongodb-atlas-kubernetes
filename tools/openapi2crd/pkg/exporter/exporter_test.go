package exporter

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestStart(t *testing.T) {
	tests := map[string]struct {
		filename      string
		overwrite     bool
		expectedError error
	}{
		/*		"creates file if it does not exist": {
				filename:  "testdata/newfile.yaml",
				overwrite: false,
			},*/
		"fails if file exists and overwrite is false": {
			filename:      "testdata/groups.atlas.generated.mongodb.com.yaml",
			overwrite:     false,
			expectedError: fmt.Errorf("file %s already exists, use --force to overwrite", "testdata/groups.atlas.generated.mongodb.com.yaml"),
		},
		"overwrites file if it exists and overwrite is true": {
			filename:  "testdata/groups.atlas.generated.mongodb.com.yaml",
			overwrite: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			_, err := fs.Create(tt.filename)
			require.NoError(t, err)

			exporter, err := New(fs, tt.filename, tt.overwrite)
			require.NoError(t, err)

			err = exporter.Start()
			require.Equal(t, tt.expectedError, err)

			if err == nil {
				_, err = fs.Stat(tt.filename)
				require.NoError(t, err)
			}
		})
	}
}

func TestExport(t *testing.T) {
	tests := map[string]struct {
		filename      string
		overwrite     bool
		crd           *apiextensionsv1.CustomResourceDefinition
		yaml          string
		expectedError error
	}{
		"export to file": {
			filename:  "testdata/groups.atlas.generated.mongodb.com.yaml",
			overwrite: true,
			crd:       sampleCRD(),
			yaml:      sampleCRDYaml(false),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			_, err := fs.Create(tt.filename)
			require.NoError(t, err)

			exporter, err := New(fs, tt.filename, tt.overwrite)
			require.NoError(t, err)

			err = exporter.Start()
			require.NoError(t, err)

			err = exporter.Export(tt.crd)
			require.Equal(t, tt.expectedError, err)
			data, err := afero.ReadFile(fs, tt.filename)
			require.NoError(t, err)
			require.Equal(t, tt.yaml, string(data))

			err = exporter.Close()
			require.NoError(t, err)
		})
	}
}

func sampleCRD() *apiextensionsv1.CustomResourceDefinition {
	crd := apiextensionsv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "groups.atlas.generated.mongodb.com",
			Annotations: map[string]string{
				"api-mappings": "properties:\n  spec:\n    properties:\n      v20250312:\n        x-atlas-sdk-version: go.mongodb.org/atlas-sdk/v20250312005/admin\n",
			},
		},
		Spec: apiextensionsv1.CustomResourceDefinitionSpec{
			Group: "atlas.generated.mongodb.com",
			Scope: apiextensionsv1.NamespaceScoped,
			Names: apiextensionsv1.CustomResourceDefinitionNames{
				Plural:     "groups",
				Singular:   "group",
				ShortNames: []string{"ag"},
				Kind:       "Group",
				ListKind:   "GroupList",
				Categories: []string{"atlas"},
			},
			Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
				{
					Name:    "v1",
					Served:  true,
					Storage: true,
				},
			},
			Validation: &apiextensionsv1.CustomResourceValidation{
				OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
					Type:        "object",
					Description: "A group, managed by the MongoDB Kubernetes Atlas Operator.",
					Properties: map[string]apiextensionsv1.JSONSchemaProps{
						"spec": {
							Type:        "object",
							Description: "Specification of the group supporting the following versions:\n\n- v20250312\n\nAt most one versioned spec can be specified.",
							Properties: map[string]apiextensionsv1.JSONSchemaProps{
								"v20250312": {
									Type:        "object",
									Description: "The spec of the group resource for version v20250312.",
									Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"entry": {
											Type:        "object",
											Description: "The entry fields of the group resource spec. These fields can be set for creating and updating groups.",
											Properties: map[string]apiextensionsv1.JSONSchemaProps{
												"name": {
													Type:        "string",
													Description: "Human-readable label that identifies the project included in the MongoDB Cloud organization.",
												},
												"orgId": {
													Type:        "string",
													Description: "Unique 24-hexadecimal digit string that identifies the MongoDB Cloud organization to which the project belongs.",
												},
											},
											Required: []string{"name", "orgId"},
										},
									},
								},
							},
						},
						"status": {
							Type:        "object",
							Description: "Most recently observed read-only status of the group for the specified resource version.",
							Properties: map[string]apiextensionsv1.JSONSchemaProps{
								"v20250312": {
									Type:        "object",
									Description: "The last observed Atlas state of the group resource for version v20250312.",
									Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"id": {
											Type:        "string",
											Description: "Unique 24-hexadecimal digit string that identifies the MongoDB Cloud project.",
										},
										"clusterCount": {
											Type:        "integer",
											Description: "Quantity of MongoDB Cloud clusters deployed in this project.",
										},
									},
									Required: []string{"clusterCount"},
								},
							},
						},
					},
				},
			},
		},
		Status: apiextensionsv1.CustomResourceDefinitionStatus{
			StoredVersions: []string{"v1"},
		},
	}

	return &crd
}

func sampleCRDYaml(withHeader bool) string {
	crd := `# The file is generated by openapi2crd
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    api-mappings: |
      properties:
        spec:
          properties:
            v20250312:
              x-atlas-sdk-version: go.mongodb.org/atlas-sdk/v20250312005/admin
  creationTimestamp: null
  name: groups.atlas.generated.mongodb.com
spec:
  group: atlas.generated.mongodb.com
  names:
    categories:
    - atlas
    kind: Group
    listKind: GroupList
    plural: groups
    shortNames:
    - ag
    singular: group
  scope: Namespaced
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        description: A group, managed by the MongoDB Kubernetes Atlas Operator.
        properties:
          spec:
            description: |-
              Specification of the group supporting the following versions:

              - v20250312

              At most one versioned spec can be specified.
            properties:
              v20250312:
                description: The spec of the group resource for version v20250312.
                properties:
                  entry:
                    description: The entry fields of the group resource spec. These
                      fields can be set for creating and updating groups.
                    properties:
                      name:
                        description: Human-readable label that identifies the project
                          included in the MongoDB Cloud organization.
                        type: string
                      orgId:
                        description: Unique 24-hexadecimal digit string that identifies
                          the MongoDB Cloud organization to which the project belongs.
                        type: string
                    required:
                    - name
                    - orgId
                    type: object
                type: object
            type: object
          status:
            description: Most recently observed read-only status of the group for
              the specified resource version.
            properties:
              v20250312:
                description: The last observed Atlas state of the group resource for
                  version v20250312.
                properties:
                  clusterCount:
                    description: Quantity of MongoDB Cloud clusters deployed in this
                      project.
                    type: integer
                  id:
                    description: Unique 24-hexadecimal digit string that identifies
                      the MongoDB Cloud project.
                    type: string
                required:
                - clusterCount
                type: object
            type: object
        type: object
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: null
  storedVersions:
  - v1
`
	if withHeader {
		crd = "# The file is generated by openapi2crd\n" + crd
	}

	return crd
}
