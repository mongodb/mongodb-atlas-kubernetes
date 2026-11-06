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

package integrations

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

func TestFetchIntegrationSecrets(t *testing.T) {
	testScheme := runtime.NewScheme()
	assert.NoError(t, v1.AddToScheme(testScheme))
	assert.NoError(t, akov2.AddToScheme(testScheme))
	for _, tc := range []struct {
		name        string
		integration *akov2.AtlasThirdPartyIntegration
		secret      *v1.Secret
		want        map[string][]byte
		wantErr     string
	}{
		{
			name: "datadog secret",
			integration: &akov2.AtlasThirdPartyIntegration{
				Spec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "DATADOG",
					Datadog: &akov2.DatadogIntegration{
						APIKeySecretRef: api.LocalObjectReference{
							Name: "datadog-secret",
						},
						Region:                       "",
						SendCollectionLatencyMetrics: new(string),
						SendDatabaseMetrics:          new(string),
					},
				},
			},
			secret: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: "datadog-secret"},
				Data:       map[string][]byte{"secret": ([]byte)("datadog")},
			},
			want: map[string][]byte{"secret": ([]byte)("datadog")},
		},
		{
			name: "ms teams secret",
			integration: &akov2.AtlasThirdPartyIntegration{
				Spec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "MICROSOFT_TEAMS",
					MicrosoftTeams: &akov2.MicrosoftTeamsIntegration{
						URLSecretRef: api.LocalObjectReference{
							Name: "msteams-secret",
						},
					},
				},
			},
			secret: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: "msteams-secret"},
				Data:       map[string][]byte{"secret": ([]byte)("msteams")},
			},
			want: map[string][]byte{"secret": ([]byte)("msteams")},
		},
		{
			name: "new relic secret",
			integration: &akov2.AtlasThirdPartyIntegration{
				Spec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "NEW_RELIC",
					NewRelic: &akov2.NewRelicIntegration{
						CredentialsSecretRef: api.LocalObjectReference{
							Name: "new-relic-secret",
						},
					},
				},
			},
			secret: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: "new-relic-secret"},
				Data:       map[string][]byte{"secret": ([]byte)("new-relic")},
			},
			want: map[string][]byte{"secret": ([]byte)("new-relic")},
		},
		{
			name: "ops genie secret",
			integration: &akov2.AtlasThirdPartyIntegration{
				Spec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "OPS_GENIE",
					OpsGenie: &akov2.OpsGenieIntegration{
						APIKeySecretRef: api.LocalObjectReference{
							Name: "ops-genie-secret",
						},
						Region: "US",
					},
				},
			},
			secret: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: "ops-genie-secret"},
				Data:       map[string][]byte{"secret": ([]byte)("ops-genie")},
			},
			want: map[string][]byte{"secret": ([]byte)("ops-genie")},
		},
		{
			name: "pager duty secret",
			integration: &akov2.AtlasThirdPartyIntegration{
				Spec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "PAGER_DUTY",
					PagerDuty: &akov2.PagerDutyIntegration{
						ServiceKeySecretRef: api.LocalObjectReference{
							Name: "pager-duty-secret",
						},
						Region: "US",
					},
				},
			},
			secret: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: "pager-duty-secret"},
				Data:       map[string][]byte{"secret": ([]byte)("pager-duty")},
			},
			want: map[string][]byte{"secret": ([]byte)("pager-duty")},
		},
		{
			name: "prometheus secret",
			integration: &akov2.AtlasThirdPartyIntegration{
				Spec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "PROMETHEUS",
					Prometheus: &akov2.PrometheusIntegration{
						PrometheusCredentialsSecretRef: api.LocalObjectReference{
							Name: "prometheus-secret",
						},
						Enabled:          new(string),
						ServiceDiscovery: "",
					},
				},
			},
			secret: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: "prometheus-secret"},
				Data:       map[string][]byte{"secret": ([]byte)("prometheus")},
			},
			want: map[string][]byte{"secret": ([]byte)("prometheus")},
		},
		{
			name: "slack secret",
			integration: &akov2.AtlasThirdPartyIntegration{
				Spec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "SLACK",
					Slack: &akov2.SlackIntegration{
						APITokenSecretRef: api.LocalObjectReference{
							Name: "slack-secret",
						},
						ChannelName: "",
						TeamName:    "",
					},
				},
			},
			secret: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: "slack-secret"},
				Data:       map[string][]byte{"secret": ([]byte)("slack-secret")},
			},
			want: map[string][]byte{"secret": ([]byte)("slack-secret")},
		},
		{
			name: "victor ops secret",
			integration: &akov2.AtlasThirdPartyIntegration{
				Spec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "VICTOR_OPS",
					VictorOps: &akov2.VictorOpsIntegration{
						APIKeySecretRef: api.LocalObjectReference{
							Name: "victor-ops-secret",
						},
						RoutingKey: "",
					},
				},
			},
			secret: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: "victor-ops-secret"},
				Data:       map[string][]byte{"secret": ([]byte)("victor-ops")},
			},
			want: map[string][]byte{"secret": ([]byte)("victor-ops")},
		},
		{
			name: "webhook secret",
			integration: &akov2.AtlasThirdPartyIntegration{
				Spec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "WEBHOOK",
					Webhook: &akov2.WebhookIntegration{
						URLSecretRef: api.LocalObjectReference{
							Name: "webhook-secret",
						},
					},
				},
			},
			secret: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: "webhook-secret"},
				Data:       map[string][]byte{"secret": ([]byte)("webhook")},
			},
			want: map[string][]byte{"secret": ([]byte)("webhook")},
		},
		{
			name: "bad integration",
			integration: &akov2.AtlasThirdPartyIntegration{
				Spec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "MADE_UP_TYPE",
				},
			},
			secret:  &v1.Secret{},
			wantErr: "unsupported integration type MADE_UP_TYPE",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tc.integration, tc.secret).
				Build()
			ctx := context.Background()

			secretData, err := fetchIntegrationSecrets(ctx, k8sClient, tc.integration)

			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.wantErr)
			}
			assert.Equal(t, tc.want, secretData)
		})
	}
}
