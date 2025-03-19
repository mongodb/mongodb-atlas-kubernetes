package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"

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
						APIKeySecret: "api-key-secretname",
						Region:       "US",
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
						URLSecret: "url-secretname",
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
						CredentialsSecret: "credentials-secretname",
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
						APIKeySecret: "api-key-secretname",
						Region:       "US",
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
						ServiceKeySecret: "service-key-secretname",
						Region:           "US",
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
						PrometheusCredentials: "prometheus-credentials",
						ServiceDiscovery:      "http",
						Enabled:               false,
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
						APITokenSecret: "api-tooken-secretname",
						ChannelName:    "channel",
						TeamName:       "team",
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
						KeysSecret: "keys-secetname",
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
						URLSecret: "url-secretname",
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
						PrometheusCredentials: "prometheus-credentials",
						ServiceDiscovery:      "http",
						Enabled:               false,
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
						APIKeySecret: "api-key-secretname",
						Region:       "US",
					},
					Webhook: &WebhookIntegration{
						URLSecret: "url-secretname",
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

			crdPath := "../../config/crd/bases/atlas.mongodb.com_atlasintegrations.yaml"
			validator, err := cel.VersionValidatorFromFile(t, crdPath, "v1")
			assert.NoError(t, err)
			errs := validator(unstructuredObject, nil)

			require.Equal(t, tc.expectedErrors, cel.ErrorListAsStrings(errs))
		})
	}
}
