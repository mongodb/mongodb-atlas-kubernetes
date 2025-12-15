// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package refs_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi/crds"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi/refs"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi/testdata"
	samplesv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi/testdata/samples/v1"
)

const (
	version = "v1"

	sdkVersion = "v20250312"
)

func TestMappings(t *testing.T) {
	groupAlertsConfigSecretsPath := []string{
		"spec",
		"v20250312",
		"entry",
		"notifications",
		"[]",
	}
	for _, tc := range []struct {
		title string
		kind  string
		want  []*refs.Mapping
	}{
		{
			title: "Map NetworkPeeringConnection",
			kind:  "NetworkPeeringConnection",
			want: []*refs.Mapping{
				groupRefmapping([]string{"spec", sdkVersion}),
			},
		},
		{
			title: "Map GroupAlertsConfig",
			kind:  "GroupAlertsConfig",
			want: []*refs.Mapping{
				secretRefMapping(groupAlertsConfigSecretsPath,
					"apiTokenSecretRef", ".apiToken"),
				secretRefMapping(groupAlertsConfigSecretsPath,
					"datadogApiKeySecretRef", ".datadogApiKey"),
				secretRefMapping(groupAlertsConfigSecretsPath,
					"microsoftTeamsWebhookUrlSecretRef", ".microsoftTeamsWebhookUrl"),
				secretRefMapping(groupAlertsConfigSecretsPath,
					"notificationTokenSecretRef", ".notificationToken"),
				secretRefMapping(groupAlertsConfigSecretsPath,
					"opsGenieApiKeySecretRef", ".opsGenieApiKey"),
				secretRefMapping(groupAlertsConfigSecretsPath,
					"serviceKeySecretRef", ".serviceKey"),
				secretRefMapping(groupAlertsConfigSecretsPath,
					"victorOpsApiKeySecretRef", ".victorOpsApiKey"),
				secretRefMapping(groupAlertsConfigSecretsPath,
					"victorOpsRoutingKeySecretRef", ".victorOpsRoutingKey"),
				secretRefMapping(groupAlertsConfigSecretsPath,
					"webhookSecretSecretRef", ".webhookSecret"),
				secretRefMapping(groupAlertsConfigSecretsPath,
					"webhookUrlSecretRef", ".webhookUrl"),
				groupRefmapping([]string{"spec", sdkVersion}),
			},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			scheme := testScheme(t)
			crdsYML := bytes.NewBuffer(testdata.SampleCRDs)
			crd, err := extractCRD(tc.kind, bufio.NewScanner(crdsYML))
			require.NoError(t, err)
			tr, err := crapi.NewTranslator(scheme, crd, version, sdkVersion)
			require.NoError(t, err)
			got, err := tr.Mappings()
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func extractCRD(kind string, scanner *bufio.Scanner) (*apiextensionsv1.CustomResourceDefinition, error) {
	for {
		crd, err := crds.Parse(scanner)
		if err != nil {
			return nil, fmt.Errorf("failed to extract CRD schema for kind %q: %w", kind, err)
		}
		if crd.Spec.Names.Kind == kind {
			return crd, nil
		}
	}
}

func secretRefMapping(dir []string, name string, prop string) *refs.Mapping {
	return refs.NewMapping(
		append(dir, name),
		name,
		refs.KubeMapping{
			NameSelector: ".name",
			PropertySelectors: []string{
				"$.data.#",
			},
			Type: refs.KubeType{
				Kind:     "Secret",
				Resource: "secrets",
				Version:  "v1",
			},
		},
		refs.OpenAPIMapping{
			Property: prop,
			Type:     "string",
		})
}

func groupRefmapping(dir []string) *refs.Mapping {
	return refs.NewMapping(
		append(dir, "groupRef"),
		"groupRef",
		refs.KubeMapping{
			NameSelector: ".name",
			Properties: []string{
				"$.status.v20250312.id",
			},
			Type: refs.KubeType{
				Kind:     "Group",
				Group:    "atlas.generated.mongodb.com",
				Resource: "groups",
				Version:  "v1",
			},
		},
		refs.OpenAPIMapping{
			Property: "$.groupId",
		})
}

var (
	secretRefStructMapping = map[string]any{
		"properties": map[string]any{
			"spec": map[string]any{
				"properties": map[string]any{
					"credentials": map[string]any{
						"properties": map[string]any{
							"apiKeyRef": map[string]any{
								"x-kubernetes-mapping": map[string]any{
									"nameSelector": ".name",
									"propertySelectors": []string{
										"$.data.#",
									},
									"type": map[string]any{
										"kind":     "Secret",
										"resource": "secrets",
										"version":  "v1",
									},
								},
								"x-openapi-mapping": map[string]any{
									"property": ".apiKey",
									"type":     "string",
								},
							},
						},
					},
				},
			},
		},
	}

	secretRefArrayMapping = map[string]any{
		"properties": map[string]any{
			"spec": map[string]any{
				"properties": map[string]any{
					"credentials": map[string]any{
						"items": map[string]any{
							"properties": map[string]any{
								"apiKeyRef": map[string]any{
									"x-kubernetes-mapping": map[string]any{
										"nameSelector": ".name",
										"propertySelectors": []string{
											"$.data.#",
										},
										"type": map[string]any{
											"kind":     "Secret",
											"resource": "secrets",
											"version":  "v1",
										},
									},
									"x-openapi-mapping": map[string]any{
										"property": ".apiKey",
										"type":     "string",
									},
								},
								"properties": map[string]any{
									"passwordRef": map[string]any{
										"x-kubernetes-mapping": map[string]any{
											"nameSelector": ".name",
											"propertySelectors": []string{
												"$.data.#",
											},
											"type": map[string]any{
												"kind":     "Secret",
												"resource": "secrets",
												"version":  "v1",
											},
										},
										"x-openapi-mapping": map[string]any{
											"property": ".password",
											"type":     "string",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
)

func TestExpandAll(t *testing.T) {
	mainObj := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "main", Namespace: "default"}}

	testCases := []struct {
		title     string
		obj       map[string]any
		deps      []client.Object
		mappings  map[string]any
		want      map[string]any
		wantAdded []client.Object
		wantErr   string
	}{
		{
			title: "expands a nested struct field reference successfully",
			obj: map[string]any{
				"spec": map[string]any{
					"credentials": map[string]any{"apiKey": "the-real-key"},
				},
			},
			mappings: secretRefStructMapping,
			want: map[string]any{
				"spec": map[string]any{
					"credentials": map[string]any{
						"apiKey": "the-real-key",
						"apiKeyRef": map[string]any{
							"key":  "apiKey",
							"name": "main-5885dc7ff969f8f85bd",
						}},
				},
			},
			wantAdded: []client.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "main-5885dc7ff969f8f85bd",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"apiKey": []byte("the-real-key"),
					},
				},
			},
		},

		{
			title: "expands an array nested field reference successfully",
			obj: map[string]any{
				"spec": map[string]any{
					"credentials": []any{
						map[string]any{"apiKey": "the-real-key"},
					},
				},
			},
			mappings: secretRefArrayMapping,
			want: map[string]any{
				"spec": map[string]any{
					"credentials": []any{
						map[string]any{
							"apiKey": "the-real-key",
							"apiKeyRef": map[string]any{
								"key":  "apiKey",
								"name": "main-f8ccbdf648bfc5f758d",
							},
						},
					},
				},
			},
			wantAdded: []client.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "main-f8ccbdf648bfc5f758d",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"apiKey": []byte("the-real-key"),
					},
				},
			},
		},

		{
			title: "no-op if secret does not match",
			obj: map[string]any{
				"spec": map[string]any{
					"credentials": map[string]any{"otherApiKey": "the-real-key"},
				},
			},
			mappings: secretRefStructMapping,
			want: map[string]any{
				"spec": map[string]any{
					"credentials": map[string]any{"otherApiKey": "the-real-key"},
				},
			},
		},

		{
			title: "no-op if mapping does not exist",
			obj: map[string]any{
				"spec": map[string]any{
					"nonCredentials": map[string]any{"apiKey": "the-real-key"},
				},
			},
			mappings: secretRefStructMapping,
			want: map[string]any{
				"spec": map[string]any{
					"nonCredentials": map[string]any{"apiKey": "the-real-key"},
				},
			},
		},

		{
			title:    "no-op when obj is empty",
			obj:      map[string]any{},
			mappings: secretRefStructMapping,
			want:     map[string]any{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			schema, err := mappingSchemaFrom(tc.mappings)
			require.NoError(t, err)
			refMappings, err := refs.FindMappings(schema, []string{})
			require.NoError(t, err)
			require.NotEmpty(t, refMappings)

			added, err := refs.ExpandAll(testScheme(t), refMappings, mainObj, tc.deps, tc.obj)

			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.want, tc.obj, "The object was not mutated as expected")
			assert.Equal(t, tc.wantAdded, added, "Added list is not as expected")
		})
	}
}

func TestCollapseReferences(t *testing.T) {
	mainObj := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "main", Namespace: "default"}}

	testCases := []struct {
		title    string
		obj      map[string]any
		deps     []client.Object
		mappings map[string]any
		want     map[string]any
		wantErr  string
	}{
		{
			title: "collapses a secret reference in struct",
			obj: map[string]any{
				"spec": map[string]any{
					"credentials": map[string]any{
						"apiKeyRef": map[string]any{
							"key":  "apiKey",
							"name": "main-89467cc86dbc9c79fdb",
						}},
				},
			},
			deps: []client.Object{
				&corev1.Secret{
					TypeMeta: metav1.TypeMeta{
						Kind:       "Secret",
						APIVersion: "v1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "main-89467cc86dbc9c79fdb",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"apiKey": []byte("the-real-key"),
					},
				},
			},
			mappings: secretRefStructMapping,
			want: map[string]any{
				"spec": map[string]any{
					"credentials": map[string]any{
						"apiKey": "the-real-key",
						"apiKeyRef": map[string]any{
							"key":  "apiKey",
							"name": "main-89467cc86dbc9c79fdb",
						}},
				},
			},
		},

		{
			title: "collapses a secret reference inside an array",
			obj: map[string]any{
				"spec": map[string]any{
					"credentials": []any{
						map[string]any{
							"apiKeyRef": map[string]any{
								"key":  "apiKey",
								"name": "main-89467cc86dbc9c79fdb",
							},
						},
					},
				},
			},
			deps: []client.Object{
				&corev1.Secret{
					TypeMeta: metav1.TypeMeta{
						Kind:       "Secret",
						APIVersion: "v1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "main-89467cc86dbc9c79fdb",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"apiKey": []byte("the-real-key"),
					},
				},
			},
			mappings: secretRefArrayMapping,
			want: map[string]any{
				"spec": map[string]any{
					"credentials": []any{
						map[string]any{
							"apiKey": "the-real-key",
							"apiKeyRef": map[string]any{
								"key":  "apiKey",
								"name": "main-89467cc86dbc9c79fdb",
							},
						},
					},
				},
			},
		},

		{
			title: "missing secret dep",
			obj: map[string]any{
				"spec": map[string]any{
					"credentials": map[string]any{
						"apiKeyRef": map[string]any{
							"key":  "apiKey",
							"name": "main-89467cc86dbc9c79fdb",
						}},
				},
			},
			mappings: secretRefStructMapping,
			want: map[string]any{
				"spec": map[string]any{
					"credentials": map[string]any{
						"apiKeyRef": map[string]any{
							"key":  "apiKey",
							"name": "main-89467cc86dbc9c79fdb",
						}},
				},
			},
			wantErr: "failed to find Kubernetes resource",
		},

		{
			title: "missing secret dep",
			obj: map[string]any{
				"spec": map[string]any{
					"credentials": map[string]any{
						"apiKeyRef": map[string]any{
							"key":  "apiKey",
							"name": "main-89467cc86dbc9c79fdb",
						}},
				},
			},
			mappings: secretRefStructMapping,
			want: map[string]any{
				"spec": map[string]any{
					"credentials": map[string]any{
						"apiKeyRef": map[string]any{
							"key":  "apiKey",
							"name": "main-89467cc86dbc9c79fdb",
						}},
				},
			},
			wantErr: "failed to find Kubernetes resource",
		},

		{
			title: "no-op if secret does not match",
			obj: map[string]any{
				"spec": map[string]any{
					"credentials": map[string]any{"otherApiKeyRef": "someref"},
				},
			},
			mappings: secretRefStructMapping,
			want: map[string]any{
				"spec": map[string]any{
					"credentials": map[string]any{"otherApiKeyRef": "someref"},
				},
			},
		},

		{
			title: "no-op if mapping does not exist",
			obj: map[string]any{
				"spec": map[string]any{
					"nonCredentials": map[string]any{"apiKeyRef": map[string]any{
						"key":  "apiKey",
						"name": "main-89467cc86dbc9c79fdb",
					}},
				},
			},
			mappings: secretRefStructMapping,
			want: map[string]any{
				"spec": map[string]any{
					"nonCredentials": map[string]any{"apiKeyRef": map[string]any{
						"key":  "apiKey",
						"name": "main-89467cc86dbc9c79fdb",
					}},
				},
			},
		},

		{
			title:    "no-op when obj is empty",
			obj:      map[string]any{},
			mappings: secretRefStructMapping,
			want:     map[string]any{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			schema, err := mappingSchemaFrom(tc.mappings)
			require.NoError(t, err)
			refMappings, err := refs.FindMappings(schema, []string{})
			require.NoError(t, err)

			err = refs.CollapseAll(testScheme(t), refMappings, mainObj, tc.deps, tc.obj)

			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.want, tc.obj, "The object was not mutated as expected")
		})
	}
}

func mappingSchemaFrom(obj map[string]any) (*openapi3.Schema, error) {
	var mappingSchema openapi3.Schema
	js, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal test mappings to JSON: %w", err)
	}
	if err := json.Unmarshal(js, &mappingSchema); err != nil {
		return nil, fmt.Errorf("failed to unmarshal test JSON into schema: %w", err)
	}
	return &mappingSchema, nil
}

func testScheme(t *testing.T) *runtime.Scheme {
	t.Helper()

	scheme := runtime.NewScheme()
	require.NoError(t, k8sscheme.AddToScheme(scheme))
	require.NoError(t, samplesv1.AddToScheme(scheme))
	return scheme
}
