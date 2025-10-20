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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/translate/refs"
)

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

func TestExpandReferences(t *testing.T) {
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
							"name": "main-89467cc86dbc9c79fdb",
						}},
				},
			},
			wantAdded: []client.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "main-89467cc86dbc9c79fdb",
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
								"name": "main-89467cc86dbc9c79fdb",
							},
						},
					},
				},
			},
			wantAdded: []client.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "main-89467cc86dbc9c79fdb",
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
			title:    "fails when obj is empty",
			obj:      map[string]any{},
			mappings: secretRefStructMapping,
			want:     map[string]any{},
			wantErr:  "path [spec] not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			h := refs.NewHandler(mainObj, tc.deps)
			err := h.ExpandReferences(tc.obj, tc.mappings, "spec")

			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.want, tc.obj, "The object was not mutated as expected")
			assert.Equal(t, tc.wantAdded, h.Added(), "Added list is not as expected")
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
			title:    "fails when obj is empty",
			obj:      map[string]any{},
			mappings: secretRefStructMapping,
			want:     map[string]any{},
			wantErr:  "path [spec] not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			handler := refs.NewHandler(mainObj, tc.deps)
			err := handler.CollapseReferences(tc.obj, tc.mappings, "spec")

			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.want, tc.obj, "The object was not mutated as expected")
			assert.Equal(t, handler.Added(), []client.Object(nil), "No additions were expected")
		})
	}
}
