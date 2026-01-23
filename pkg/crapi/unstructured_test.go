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
//

package crapi_test

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	admin2025 "go.mongodb.org/atlas-sdk/v20250312012/admin"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/crapi"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/crapi/testdata"
)

//nolint:dupl
func TestToAPIUnstructured(t *testing.T) {
	for _, tc := range []struct {
		crd    string
		input  client.Object
		deps   []client.Object
		target admin2025.GroupAlertsConfig
		want   admin2025.GroupAlertsConfig
	}{
		{
			crd: "GroupAlertsConfig",
			input: &unstructured.Unstructured{
				Object: map[string]any{
					"apiVersion": "atlas.generated.mongodb.com/v1",
					"kind":       "GroupAlertsConfig",
					"metadata": map[string]any{
						"name":      "my-group-alerts-config",
						"namespace": "ns",
					},
					"spec": map[string]any{
						"v20250312": map[string]any{
							"entry": map[string]any{
								"enabled":       true,
								"eventTypeName": "some-event",
								"matchers": []any{
									map[string]any{
										"fieldName": "field1",
										"operator":  "op1",
										"value":     "value1",
									},
									map[string]any{
										"fieldName": "field2",
										"operator":  "op2",
										"value":     "value2",
									},
								},
								"metricThreshold": map[string]any{
									"metricName": "metric",
									"mode":       "mode",
									"operator":   "operator",
									"threshold":  1.0,
									"units":      "unit",
								},
								"notifications": []any{
									map[string]any{
										"datadogApiKeySecretRef": map[string]any{
											"name": "alert-secrets-0",
											"key":  "apiKey",
										},
										"datadogRegion": "US",
									},
									map[string]any{
										"webhookSecretSecretRef": map[string]any{
											"name": "alert-secrets-0",
											"key":  "webhookSecret",
										},
										"webhookUrlSecretRef": map[string]any{
											"name": "alert-secrets-1",
											"key":  "webhookUrl",
										},
									},
								},
								"severityOverride": "severe",
								"threshold": map[string]any{
									"metricName": "metric",
									"mode":       "mode-t",
									"operator":   "op-t",
									"threshold":  2.0,
									"units":      "unit-t",
								},
							},
							"groupRef": map[string]any{
								"name": "my-project",
							},
						},
					},
				},
			},
			deps: []client.Object{
				&unstructured.Unstructured{
					Object: map[string]any{
						"apiVersion": "atlas.generated.mongodb.com/v1",
						"kind":       "Group",
						"metadata": map[string]any{
							"name":      "my-project",
							"namespace": "ns",
						},
						"spec": map[string]any{
							"v20250312": map[string]any{
								"entry": map[string]any{
									"name":  "some-project",
									"orgId": "621454123423x125235142",
								},
							},
						},
						"status": map[string]any{
							"v20250312": map[string]any{
								"id": "62b6e34b3d91647abb20e7b8",
							},
						},
					},
				},
				&unstructured.Unstructured{
					Object: map[string]any{
						"apiVersion": "v1",
						"kind":       "Secret",
						"metadata": map[string]any{
							"name":      "alert-secrets-0",
							"namespace": "ns",
						},
						"data": map[string]any{
							"apiKey":        "c2FtcGxlLWFwaS1rZXk=",
							"webhookSecret": "c2FtcGxlLXdlYmhvb2stc2VjcmV0",
						},
					},
				},
				&unstructured.Unstructured{
					Object: map[string]any{
						"apiVersion": "v1",
						"kind":       "Secret",
						"metadata": map[string]any{
							"name":      "alert-secrets-1",
							"namespace": "ns",
						},
						"data": map[string]any{
							"webhookUrl": "c2FtcGxlLXdlYmhvb2stdXJs",
						},
					},
				},
			},
			want: admin2025.GroupAlertsConfig{
				Enabled:       pointer.MakePtr(true),
				EventTypeName: pointer.MakePtr("some-event"),
				GroupId:       pointer.MakePtr("62b6e34b3d91647abb20e7b8"),
				Matchers: &[]admin2025.StreamsMatcher{
					{
						FieldName: "field1",
						Operator:  "op1",
						Value:     "value1",
					},
					{
						FieldName: "field2",
						Operator:  "op2",
						Value:     "value2",
					},
				},
				Notifications: &[]admin2025.AlertsNotificationRootForGroup{
					{
						DatadogApiKey: pointer.MakePtr("sample-api-key"),
						DatadogRegion: pointer.MakePtr("US"),
					},
					{
						WebhookSecret: pointer.MakePtr("sample-webhook-secret"),
						WebhookUrl:    pointer.MakePtr("sample-webhook-url"),
					},
				},
				SeverityOverride: pointer.MakePtr("severe"),
				MetricThreshold: &admin2025.FlexClusterMetricThreshold{
					MetricName: "metric",
					Mode:       pointer.MakePtr("mode"),
					Operator:   pointer.MakePtr("operator"),
					Threshold:  pointer.MakePtr(1.0),
					Units:      pointer.MakePtr("unit"),
				},
				Threshold: &admin2025.StreamProcessorMetricThreshold{
					MetricName: pointer.MakePtr("metric"),
					Mode:       pointer.MakePtr("mode-t"),
					Operator:   pointer.MakePtr("op-t"),
					Threshold:  pointer.MakePtr(2.0),
					Units:      pointer.MakePtr("unit-t"),
				},
			},
			target: admin2025.GroupAlertsConfig{},
		},
	} {
		scheme := testScheme(t)
		crdsYML := bytes.NewBuffer(testdata.SampleCRDs)
		crd, err := extractCRD(tc.crd, bufio.NewScanner(crdsYML))
		require.NoError(t, err)
		tr, err := crapi.NewTranslator(scheme, crd, version, sdkVersion)
		require.NoError(t, err)
		require.NoError(t, tr.ToAPI(&tc.target, tc.input, tc.deps...))
		assert.Equal(t, tc.want, tc.target)
	}
}

//nolint:dupl
func TestFromAPIUnstructured(t *testing.T) {
	for _, tc := range []struct {
		crd      string
		input    any
		target   client.Object
		deps     []client.Object
		want     client.Object
		wantDeps []client.Object
	}{
		{
			crd: "GroupAlertsConfig",
			input: admin2025.GroupAlertsConfig{
				Enabled:       pointer.MakePtr(true),
				EventTypeName: pointer.MakePtr("some-event"),
				GroupId:       pointer.MakePtr("62b6e34b3d91647abb20e7b8"),
				Matchers: &[]admin2025.StreamsMatcher{
					{
						FieldName: "field1",
						Operator:  "op1",
						Value:     "value1",
					},
					{
						FieldName: "field2",
						Operator:  "op2",
						Value:     "value2",
					},
				},
				Notifications: &[]admin2025.AlertsNotificationRootForGroup{
					{
						DatadogApiKey: pointer.MakePtr("sample-api-key"),
						DatadogRegion: pointer.MakePtr("US"),
					},
					{
						WebhookSecret: pointer.MakePtr("sample-webhook-secret"),
						WebhookUrl:    pointer.MakePtr("sample-webhook-url"),
					},
				},
				SeverityOverride: pointer.MakePtr("severe"),
				MetricThreshold: &admin2025.FlexClusterMetricThreshold{
					MetricName: "metric",
					Mode:       pointer.MakePtr("mode"),
					Operator:   pointer.MakePtr("operator"),
					Threshold:  pointer.MakePtr(1.0),
					Units:      pointer.MakePtr("unit"),
				},
				Threshold: &admin2025.StreamProcessorMetricThreshold{
					MetricName: pointer.MakePtr("metric"),
					Mode:       pointer.MakePtr("mode-t"),
					Operator:   pointer.MakePtr("op-t"),
					Threshold:  pointer.MakePtr(2.0),
					Units:      pointer.MakePtr("unit-t"),
				},
			},
			target: &unstructured.Unstructured{
				Object: map[string]any{
					"apiVersion": "atlas.generated.mongodb.com/v1",
					"kind":       "GroupAlertsConfig",
					"metadata": map[string]any{
						"name":      "my-group-alerts-config",
						"namespace": "ns",
					},
				},
			},
			deps: []client.Object{
				&unstructured.Unstructured{
					Object: map[string]any{
						"apiVersion": "atlas.generated.mongodb.com/v1",
						"kind":       "Group",
						"metadata": map[string]any{
							"name":      "my-project",
							"namespace": "ns",
						},
						"spec": map[string]any{
							"v20250312": map[string]any{
								"entry": map[string]any{
									"name":  "some-project",
									"orgId": "621454123423x125235142",
								},
							},
						},
						"status": map[string]any{
							"v20250312": map[string]any{
								"id": "62b6e34b3d91647abb20e7b8",
							},
						},
					},
				},
			},
			want: &unstructured.Unstructured{
				Object: map[string]any{
					"apiVersion": "atlas.generated.mongodb.com/v1",
					"kind":       "GroupAlertsConfig",
					"metadata": map[string]any{
						"name":      "my-group-alerts-config",
						"namespace": "ns",
					},
					"spec": map[string]any{
						"v20250312": map[string]any{
							"enabled":       true,                       // extra field
							"eventTypeName": "some-event",               // extra field
							"groupId":       "62b6e34b3d91647abb20e7b8", // extra field
							"matchers": []any{ // extra field structure
								map[string]any{
									"fieldName": "field1",
									"operator":  "op1",
									"value":     "value1",
								},
								map[string]any{
									"fieldName": "field2",
									"operator":  "op2",
									"value":     "value2",
								},
							},
							"metricThreshold": map[string]any{ // extra field structure
								"metricName": "metric",
								"mode":       "mode",
								"operator":   "operator",
								"threshold":  int64(1),
								"units":      "unit",
							},
							"notifications": []any{ // extra field structure
								map[string]any{
									"datadogApiKey": "sample-api-key", // extra field
									"datadogApiKeySecretRef": map[string]any{
										"name": "my-group-alerts-config-f4f4b5f9c849fc4cbdc",
										"key":  "datadogApiKey",
									},
									"datadogRegion": "US",
								},
								map[string]any{
									"webhookSecret": "sample-webhook-secret", // extra field
									"webhookSecretSecretRef": map[string]any{
										"name": "my-group-alerts-config-54854758bf4fbd4fb559",
										"key":  "webhookSecret",
									},
									"webhookUrl": "sample-webhook-url", // extra field
									"webhookUrlSecretRef": map[string]any{
										"name": "my-group-alerts-config-d9f674cf64c78546d88",
										"key":  "webhookUrl",
									},
								},
							},
							"severityOverride": "severe", // extra field
							"threshold": map[string]any{ // extra field structure
								"metricName": "metric",
								"mode":       "mode-t",
								"operator":   "op-t",
								"threshold":  int64(2),
								"units":      "unit-t",
							},
							"entry": map[string]any{
								"groupId":       "62b6e34b3d91647abb20e7b8", // extra field
								"enabled":       true,
								"eventTypeName": "some-event",
								"matchers": []any{
									map[string]any{
										"fieldName": "field1",
										"operator":  "op1",
										"value":     "value1",
									},
									map[string]any{
										"fieldName": "field2",
										"operator":  "op2",
										"value":     "value2",
									},
								},
								"metricThreshold": map[string]any{
									"metricName": "metric",
									"mode":       "mode",
									"operator":   "operator",
									"threshold":  int64(1),
									"units":      "unit",
								},
								"notifications": []any{
									map[string]any{
										"datadogApiKey": "sample-api-key", // extra field
										"datadogApiKeySecretRef": map[string]any{
											"name": "my-group-alerts-config-f4f4b5f9c849fc4cbdc",
											"key":  "datadogApiKey",
										},
										"datadogRegion": "US",
									},
									map[string]any{
										"webhookSecret": "sample-webhook-secret", // extra field
										"webhookSecretSecretRef": map[string]any{
											"name": "my-group-alerts-config-54854758bf4fbd4fb559",
											"key":  "webhookSecret",
										},
										"webhookUrl": "sample-webhook-url", // extra field
										"webhookUrlSecretRef": map[string]any{
											"name": "my-group-alerts-config-d9f674cf64c78546d88",
											"key":  "webhookUrl",
										},
									},
								},
								"severityOverride": "severe",
								"threshold": map[string]any{
									"metricName": "metric",
									"mode":       "mode-t",
									"operator":   "op-t",
									"threshold":  int64(2),
									"units":      "unit-t",
								},
							},
						},
					},
					"status": map[string]any{
						"v20250312": map[string]any{
							"groupId":       "62b6e34b3d91647abb20e7b8",
							"enabled":       true,         // extra field
							"eventTypeName": "some-event", // extra field
							"matchers": []any{ // extra field structure
								map[string]any{
									"fieldName": "field1",
									"operator":  "op1",
									"value":     "value1",
								},
								map[string]any{
									"fieldName": "field2",
									"operator":  "op2",
									"value":     "value2",
								},
							},
							"metricThreshold": map[string]any{ // extra field structure
								"metricName": "metric",
								"mode":       "mode",
								"operator":   "operator",
								"threshold":  int64(1),
								"units":      "unit",
							},
							"notifications": []any{ // extra field structure
								map[string]any{
									"datadogApiKey": "sample-api-key", // extra field
									"datadogApiKeySecretRef": map[string]any{
										"name": "my-group-alerts-config-f4f4b5f9c849fc4cbdc",
										"key":  "datadogApiKey",
									},
									"datadogRegion": "US",
								},
								map[string]any{
									"webhookSecret": "sample-webhook-secret", // extra field
									"webhookSecretSecretRef": map[string]any{
										"name": "my-group-alerts-config-54854758bf4fbd4fb559",
										"key":  "webhookSecret",
									},
									"webhookUrl": "sample-webhook-url", // extra field
									"webhookUrlSecretRef": map[string]any{
										"name": "my-group-alerts-config-d9f674cf64c78546d88",
										"key":  "webhookUrl",
									},
								},
							},
							"severityOverride": "severe", // extra field
							"threshold": map[string]any{ // extra field structure
								"metricName": "metric",
								"mode":       "mode-t",
								"operator":   "op-t",
								"threshold":  int64(2),
								"units":      "unit-t",
							},
						},
					},
				},
			},
			wantDeps: []client.Object{
				&corev1.Secret{
					ObjectMeta: v1.ObjectMeta{
						Name:      "my-group-alerts-config-f4f4b5f9c849fc4cbdc",
						Namespace: "ns",
					},
					Data: map[string][]byte{
						"datadogApiKey": ([]byte)("sample-api-key"),
					},
				},
				&corev1.Secret{
					ObjectMeta: v1.ObjectMeta{
						Name:      "my-group-alerts-config-54854758bf4fbd4fb559",
						Namespace: "ns",
					},
					Data: map[string][]byte{
						"webhookSecret": ([]byte)("sample-webhook-secret"),
					},
				},
				&corev1.Secret{
					ObjectMeta: v1.ObjectMeta{
						Name:      "my-group-alerts-config-d9f674cf64c78546d88",
						Namespace: "ns",
					},
					Data: map[string][]byte{
						"webhookUrl": ([]byte)("sample-webhook-url"),
					},
				},
			},
		},
	} {
		scheme := testScheme(t)
		crdsYML := bytes.NewBuffer(testdata.SampleCRDs)
		crd, err := extractCRD(tc.crd, bufio.NewScanner(crdsYML))
		require.NoError(t, err)
		tr, err := crapi.NewTranslator(scheme, crd, version, sdkVersion)
		require.NoError(t, err)
		results, err := tr.FromAPI(tc.target, tc.input, tc.deps...)
		require.NoError(t, err)
		assert.Equal(t, tc.want, tc.target)
		assert.Equal(t, tc.wantDeps, results)
	}
}
