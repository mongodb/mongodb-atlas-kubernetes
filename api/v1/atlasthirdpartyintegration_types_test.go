// Copyright 2025 MongoDB.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/cel"
)

func TestIntegrationCELChecks(t *testing.T) {
	for _, tc := range []struct {
		title          string
		obj            *AtlasThirdPartyIntegration
		expectedErrors []string
	}{
		{
			title:          "fails with no type",
			obj:            &AtlasThirdPartyIntegration{},
			expectedErrors: []string{"spec: Invalid value: \"object\": must define a type of integration"},
		},
		{
			title: "Datadog works",
			obj: &AtlasThirdPartyIntegration{
				Spec: AtlasThirdPartyIntegrationSpec{
					Type: "DATADOG",
					Datadog: &DatadogIntegration{
						APIKeySecret: api.LocalObjectReference{
							Name: "api-key-secretname",
						},
						Region: "US",
					},
				},
			},
		},
		{
			title: "Microsoft Teams works",
			obj: &AtlasThirdPartyIntegration{
				Spec: AtlasThirdPartyIntegrationSpec{
					Type: "MICROSOFT_TEAMS",
					MicrosoftTeams: &MicrosoftTeamsIntegration{
						URLSecret: api.LocalObjectReference{
							Name: "url-secretname",
						},
					},
				},
			},
		},
		{
			title: "New Relic works",
			obj: &AtlasThirdPartyIntegration{
				Spec: AtlasThirdPartyIntegrationSpec{
					Type: "NEW_RELIC",
					NewRelic: &NewRelicIntegration{
						CredentialsSecret: api.LocalObjectReference{
							Name: "credentials-secretname",
						},
					},
				},
			},
		},
		{
			title: "Ops Genie works",
			obj: &AtlasThirdPartyIntegration{
				Spec: AtlasThirdPartyIntegrationSpec{
					Type: "OPS_GENIE",
					OpsGenie: &OpsGenieIntegration{
						APIKeySecret: api.LocalObjectReference{
							Name: "api-key-secretname",
						},
						Region: "US",
					},
				},
			},
		},
		{
			title: "Pager Duty works",
			obj: &AtlasThirdPartyIntegration{
				Spec: AtlasThirdPartyIntegrationSpec{
					Type: "PAGER_DUTY",
					PagerDuty: &PagerDutyIntegration{
						ServiceKeySecret: api.LocalObjectReference{
							Name: "service-key-secretname",
						},
						Region: "US",
					},
				},
			},
		},
		{
			title: "Prometheus duty works",
			obj: &AtlasThirdPartyIntegration{
				Spec: AtlasThirdPartyIntegrationSpec{
					Type: "PROMETHEUS",
					Prometheus: &PrometheusIntegration{
						PrometheusCredentials: api.LocalObjectReference{
							Name: "prometheus-credentials",
						},
						ServiceDiscovery: "http",
					},
				},
			},
		},
		{
			title: "Slack works",
			obj: &AtlasThirdPartyIntegration{
				Spec: AtlasThirdPartyIntegrationSpec{
					Type: "SLACK",
					Slack: &SlackIntegration{
						APITokenSecret: api.LocalObjectReference{
							Name: "api-tooken-secretname",
						},
						ChannelName: "channel",
						TeamName:    "team",
					},
				},
			},
		},
		{
			title: "Victor ops works",
			obj: &AtlasThirdPartyIntegration{
				Spec: AtlasThirdPartyIntegrationSpec{
					Type: "VICTOR_OPS",
					VictorOps: &VictorOpsIntegration{
						RoutingKey: api.LocalObjectReference{
							Name: "routing-key",
						},
						APIKeySecret: "keys-secetname",
					},
				},
			},
		},
		{
			title: "Webhook works",
			obj: &AtlasThirdPartyIntegration{
				Spec: AtlasThirdPartyIntegrationSpec{
					Type: "WEBHOOK",
					Webhook: &WebhookIntegration{
						URLSecret: api.LocalObjectReference{
							Name: "url-secretname",
						},
					},
				},
			},
		},
		{
			title: "Prometheus on Pager Duty type fails",
			obj: &AtlasThirdPartyIntegration{
				Spec: AtlasThirdPartyIntegrationSpec{
					Type:      "PAGER_DUTY",
					PagerDuty: &PagerDutyIntegration{},
					Prometheus: &PrometheusIntegration{
						PrometheusCredentials: api.LocalObjectReference{
							Name: "prometheus-credentials",
						},
						ServiceDiscovery: "http",
					},
				},
			},
			expectedErrors: []string{"spec: Invalid value: \"object\": only PROMETHEUS type may set prometheus fields"},
		},
		{
			title: "Datadog on Webhook type fails",
			obj: &AtlasThirdPartyIntegration{
				Spec: AtlasThirdPartyIntegrationSpec{
					Type: "WEBHOOK",
					Datadog: &DatadogIntegration{
						APIKeySecret: api.LocalObjectReference{
							Name: "api-key-secretname",
						},
						Region: "US",
					},
					Webhook: &WebhookIntegration{
						URLSecret: api.LocalObjectReference{
							Name: "url-secretname",
						},
					},
				},
			},
			expectedErrors: []string{"spec: Invalid value: \"object\": only DATADOG type may set datadog fields"},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			// inject a project to avoid other CEL validations being hit
			tc.obj.Spec.ProjectRef = &common.ResourceRefNamespaced{Name: "some-project"}
			unstructuredObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&tc.obj)
			require.NoError(t, err)

			crdPath := "../../config/crd/bases/atlas.mongodb.com_atlasthirdpartyintegrations.yaml"
			validator, err := cel.VersionValidatorFromFile(t, crdPath, "v1")
			assert.NoError(t, err)
			errs := validator(unstructuredObject, nil)

			require.Equal(t, tc.expectedErrors, cel.ErrorListAsStrings(errs))
		})
	}
}
