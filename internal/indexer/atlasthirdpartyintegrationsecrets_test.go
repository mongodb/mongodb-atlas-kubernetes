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

package indexer

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

func TestAtlasThirdPartyIntgerationBySecretsIndexer(t *testing.T) {
	for _, tc := range []struct {
		name     string
		object   client.Object
		wantKeys []string
	}{
		{
			name:   "should return nil on wrong type",
			object: &akov2.AtlasProject{},
		},
		{
			name:   "should return nil when there are no references",
			object: &akov2.AtlasThirdPartyIntegration{},
		},
		{
			name: "should return nil when there is an empty reference",
			object: &akov2.AtlasThirdPartyIntegration{
				Spec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "DATADOG",
					Datadog: &akov2.DatadogIntegration{
						APIKeySecretRef: api.LocalObjectReference{},
					},
				},
			},
		},
		{
			name: "should return the datadog secret name",
			object: &akov2.AtlasThirdPartyIntegration{
				ObjectMeta: v1.ObjectMeta{Namespace: "ns"},
				Spec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "DATADOG",
					Datadog: &akov2.DatadogIntegration{
						APIKeySecretRef: api.LocalObjectReference{
							Name: "datadogSecret",
						},
					},
				},
			},
			wantKeys: []string{"ns/datadogSecret"},
		},
		{
			name: "should return the microsoft teams secret name",
			object: &akov2.AtlasThirdPartyIntegration{
				ObjectMeta: v1.ObjectMeta{Namespace: "ns"},
				Spec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "MICROSOFT_TEAMS",
					MicrosoftTeams: &akov2.MicrosoftTeamsIntegration{
						URLSecretRef: api.LocalObjectReference{
							Name: "microsoftTeamsSecret",
						},
					},
				},
			},
			wantKeys: []string{"ns/microsoftTeamsSecret"},
		},
		{
			name: "should return the new relic secret name",
			object: &akov2.AtlasThirdPartyIntegration{
				ObjectMeta: v1.ObjectMeta{Namespace: "ns"},
				Spec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "NEW_RELIC",
					NewRelic: &akov2.NewRelicIntegration{
						CredentialsSecretRef: api.LocalObjectReference{
							Name: "newRelicSecret",
						},
					},
				},
			},
			wantKeys: []string{"ns/newRelicSecret"},
		},
		{
			name: "should return the ops genie secret name",
			object: &akov2.AtlasThirdPartyIntegration{
				ObjectMeta: v1.ObjectMeta{Namespace: "ns"},
				Spec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "OPS_GENIE",
					OpsGenie: &akov2.OpsGenieIntegration{
						APIKeySecretRef: api.LocalObjectReference{
							Name: "opsGenieSecret",
						},
					},
				},
			},
			wantKeys: []string{"ns/opsGenieSecret"},
		},
		{
			name: "should return the pager duty secret name",
			object: &akov2.AtlasThirdPartyIntegration{
				ObjectMeta: v1.ObjectMeta{Namespace: "ns"},
				Spec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "PAGER_DUTY",
					PagerDuty: &akov2.PagerDutyIntegration{
						ServiceKeySecretRef: api.LocalObjectReference{
							Name: "pagerDutySecret",
						},
					},
				},
			},
			wantKeys: []string{"ns/pagerDutySecret"},
		},
		{
			name: "should return the prometheus secret name",
			object: &akov2.AtlasThirdPartyIntegration{
				ObjectMeta: v1.ObjectMeta{Namespace: "ns"},
				Spec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "PROMETHEUS",
					Prometheus: &akov2.PrometheusIntegration{
						PrometheusCredentialsSecretRef: api.LocalObjectReference{
							Name: "prometheusSecret",
						},
					},
				},
			},
			wantKeys: []string{"ns/prometheusSecret"},
		},
		{
			name: "should return the slack secret name",
			object: &akov2.AtlasThirdPartyIntegration{
				ObjectMeta: v1.ObjectMeta{Namespace: "ns"},
				Spec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "SLACK",
					Slack: &akov2.SlackIntegration{
						APITokenSecretRef: api.LocalObjectReference{
							Name: "slackSecret",
						},
					},
				},
			},
			wantKeys: []string{"ns/slackSecret"},
		},
		{
			name: "should return the victor ops secret name",
			object: &akov2.AtlasThirdPartyIntegration{
				ObjectMeta: v1.ObjectMeta{Namespace: "ns"},
				Spec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "VICTOR_OPS",
					VictorOps: &akov2.VictorOpsIntegration{
						APIKeySecretRef: api.LocalObjectReference{
							Name: "victorOpsSecret",
						},
					},
				},
			},
			wantKeys: []string{"ns/victorOpsSecret"},
		},
		{
			name: "should return the webhook api key name",
			object: &akov2.AtlasThirdPartyIntegration{
				ObjectMeta: v1.ObjectMeta{Namespace: "ns"},
				Spec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "WEBHOOK",
					Webhook: &akov2.WebhookIntegration{
						URLSecretRef: api.LocalObjectReference{
							Name: "webhookSecret",
						},
					},
				},
			},
			wantKeys: []string{"ns/webhookSecret"},
		},
		{
			name: "wrong type returns nothing",
			object: &akov2.AtlasThirdPartyIntegration{
				Spec: akov2.AtlasThirdPartyIntegrationSpec{
					Type: "DATADOG",
					Webhook: &akov2.WebhookIntegration{
						URLSecretRef: api.LocalObjectReference{
							Name: "webhookSecret",
						},
					},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			indexer := NewAtlasThirdPartyIntegrationBySecretsIndexer(zaptest.NewLogger(t))
			keys := indexer.Keys(tc.object)
			sort.Strings(keys)
			assert.Equal(t, tc.wantKeys, keys)
		})
	}
}
