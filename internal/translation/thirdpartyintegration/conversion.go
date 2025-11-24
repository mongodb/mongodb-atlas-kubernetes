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

package thirdpartyintegration

import (
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312009/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

var (
	// ErrUnsupportedIntegrationType when the integration type is not supported
	ErrUnsupportedIntegrationType = errors.New("unsupported integration type")
)

func NewFromSpec(crd *akov2.AtlasThirdPartyIntegration, secrets map[string][]byte) (*ThirdPartyIntegration, error) {
	tpi := ThirdPartyIntegration{
		AtlasThirdPartyIntegrationSpec: crd.Spec,
		ID:                             crd.Status.ID,
	}
	switch tpi.Type {
	case "DATADOG":
		tpi.DatadogSecrets = &DatadogSecrets{
			APIKey: string(secrets["apiKey"]),
		}
	case "MICROSOFT_TEAMS":
		tpi.MicrosoftTeamsSecrets = &MicrosoftTeamsSecrets{
			WebhookUrl: string(secrets["webhookURL"]),
		}
	case "NEW_RELIC":
		tpi.NewRelicSecrets = &NewRelicSecrets{
			AccountID:  string(secrets["accountId"]),
			LicenseKey: string(secrets["licenseKey"]),
			ReadToken:  string(secrets["readToken"]),
			WriteToken: string(secrets["writeToken"]),
		}
	case "OPS_GENIE":
		tpi.OpsGenieSecrets = &OpsGenieSecrets{
			APIKey: string(secrets["apiKey"]),
		}
	case "PAGER_DUTY":
		tpi.PagerDutySecrets = &PagerDutySecrets{
			ServiceKey: string(secrets["serviceKey"]),
		}
	case "PROMETHEUS":
		tpi.PrometheusSecrets = &PrometheusSecrets{
			Username: string(secrets["username"]),
			Password: string(secrets["password"]),
		}
	case "SLACK":
		tpi.SlackSecrets = &SlackSecrets{
			APIToken: string(secrets["apiToken"]),
		}
	case "VICTOR_OPS":
		tpi.VictorOpsSecrets = &VictorOpsSecrets{
			APIKey: string(secrets["apiKey"]),
		}
	case "WEBHOOK":
		tpi.WebhookSecrets = &WebhookSecrets{
			URL:    string(secrets["url"]),
			Secret: string(secrets["secret"]),
		}
	default:
		return nil, fmt.Errorf("%w %v", ErrUnsupportedIntegrationType, tpi.Type)
	}
	return &tpi, nil
}

type ThirdPartyIntegration struct {
	akov2.AtlasThirdPartyIntegrationSpec
	ID                    string
	DatadogSecrets        *DatadogSecrets
	MicrosoftTeamsSecrets *MicrosoftTeamsSecrets
	NewRelicSecrets       *NewRelicSecrets
	OpsGenieSecrets       *OpsGenieSecrets
	PagerDutySecrets      *PagerDutySecrets
	PrometheusSecrets     *PrometheusSecrets
	SlackSecrets          *SlackSecrets
	VictorOpsSecrets      *VictorOpsSecrets
	WebhookSecrets        *WebhookSecrets
}

// Comparable returns a copy of ThirdPartyIntegration without secrets,
// so that it is comparable
func (tpi *ThirdPartyIntegration) Comparable() *ThirdPartyIntegration {
	comparable := &ThirdPartyIntegration{
		AtlasThirdPartyIntegrationSpec: *tpi.AtlasThirdPartyIntegrationSpec.DeepCopy(),
		ID:                             tpi.ID,
	}
	comparable.AtlasThirdPartyIntegrationSpec.ConnectionSecret = nil
	comparable.AtlasThirdPartyIntegrationSpec.ExternalProjectRef = nil
	comparable.AtlasThirdPartyIntegrationSpec.ProjectRef = nil
	if comparable.AtlasThirdPartyIntegrationSpec.Datadog != nil {
		comparable.AtlasThirdPartyIntegrationSpec.Datadog.APIKeySecretRef.Name = ""
	}
	if comparable.AtlasThirdPartyIntegrationSpec.MicrosoftTeams != nil {
		comparable.AtlasThirdPartyIntegrationSpec.MicrosoftTeams.URLSecretRef.Name = ""
	}
	if comparable.AtlasThirdPartyIntegrationSpec.NewRelic != nil {
		comparable.AtlasThirdPartyIntegrationSpec.NewRelic.CredentialsSecretRef.Name = ""
	}
	if comparable.AtlasThirdPartyIntegrationSpec.OpsGenie != nil {
		comparable.AtlasThirdPartyIntegrationSpec.OpsGenie.APIKeySecretRef.Name = ""
	}
	if comparable.AtlasThirdPartyIntegrationSpec.PagerDuty != nil {
		comparable.AtlasThirdPartyIntegrationSpec.PagerDuty.ServiceKeySecretRef.Name = ""
	}
	if comparable.AtlasThirdPartyIntegrationSpec.Prometheus != nil {
		comparable.AtlasThirdPartyIntegrationSpec.Prometheus.PrometheusCredentialsSecretRef.Name = ""
	}
	if comparable.AtlasThirdPartyIntegrationSpec.Slack != nil {
		comparable.AtlasThirdPartyIntegrationSpec.Slack.APITokenSecretRef.Name = ""
	}
	if comparable.AtlasThirdPartyIntegrationSpec.VictorOps != nil {
		comparable.AtlasThirdPartyIntegrationSpec.VictorOps.APIKeySecretRef.Name = ""
	}
	if comparable.AtlasThirdPartyIntegrationSpec.Webhook != nil {
		comparable.AtlasThirdPartyIntegrationSpec.Webhook.URLSecretRef.Name = ""
	}
	// unset ID for comparison, as the ID changes on each update
	comparable.ID = ""
	return comparable
}

type DatadogSecrets struct {
	APIKey string
}

type MicrosoftTeamsSecrets struct {
	WebhookUrl string
}

type NewRelicSecrets struct {
	AccountID  string
	LicenseKey string
	ReadToken  string
	WriteToken string
}

type OpsGenieSecrets struct {
	APIKey string
}

type PagerDutySecrets struct {
	ServiceKey string
}

type PrometheusSecrets struct {
	Username string
	Password string
}

type SlackSecrets struct {
	APIToken string
}

type VictorOpsSecrets struct {
	APIKey string
}

type WebhookSecrets struct {
	URL    string
	Secret string
}

func toAtlas(tpi *ThirdPartyIntegration) (*admin.ThirdPartyIntegration, error) {
	ai := &admin.ThirdPartyIntegration{
		Id:   &tpi.ID,
		Type: &tpi.Type,
	}
	switch tpi.Type {
	case "DATADOG":
		if tpi.Datadog == nil || tpi.DatadogSecrets == nil {
			return nil, errors.New("missing Datadog settings")
		}
		ai.ApiKey = &tpi.DatadogSecrets.APIKey
		ai.Region = &tpi.Datadog.Region
		ai.SendCollectionLatencyMetrics = pointer.MakePtr(isEnabled(tpi.Datadog.SendCollectionLatencyMetrics))
		ai.SendDatabaseMetrics = pointer.MakePtr(isEnabled(tpi.Datadog.SendDatabaseMetrics))
	case "MICROSOFT_TEAMS":
		if tpi.MicrosoftTeams == nil || tpi.MicrosoftTeamsSecrets == nil {
			return nil, errors.New("missing Microsoft teams settings")
		}
		ai.MicrosoftTeamsWebhookUrl = &tpi.MicrosoftTeamsSecrets.WebhookUrl
	case "NEW_RELIC":
		if tpi.NewRelic == nil || tpi.NewRelicSecrets == nil {
			return nil, errors.New("missing New Relic settings")
		}
		ai.AccountId = &tpi.NewRelicSecrets.AccountID
		ai.LicenseKey = &tpi.NewRelicSecrets.LicenseKey
		ai.ReadToken = &tpi.NewRelicSecrets.ReadToken
		ai.WriteToken = &tpi.NewRelicSecrets.WriteToken
	case "OPS_GENIE":
		if tpi.OpsGenie == nil || tpi.OpsGenieSecrets == nil {
			return nil, errors.New("missing OpsGenie settings")
		}
		ai.ApiKey = &tpi.OpsGenieSecrets.APIKey
		ai.Region = &tpi.OpsGenie.Region
	case "PAGER_DUTY":
		if tpi.PagerDuty == nil || tpi.PagerDutySecrets == nil {
			return nil, errors.New("missing Pager Duty settings")
		}
		ai.ServiceKey = &tpi.PagerDutySecrets.ServiceKey
		ai.Region = &tpi.PagerDuty.Region
	case "PROMETHEUS":
		if tpi.Prometheus == nil || tpi.PrometheusSecrets == nil {
			return nil, errors.New("missing Prometheus settings")
		}
		ai.Enabled = pointer.MakePtr(isEnabled(tpi.Prometheus.Enabled))
		ai.Username = &tpi.PrometheusSecrets.Username
		ai.Password = &tpi.PrometheusSecrets.Password
		ai.ServiceDiscovery = &tpi.Prometheus.ServiceDiscovery
	case "SLACK":
		if tpi.Slack == nil || tpi.SlackSecrets == nil {
			return nil, errors.New("missing Slack settings")
		}
		ai.ApiToken = &tpi.SlackSecrets.APIToken
		ai.ChannelName = &tpi.Slack.ChannelName
		ai.TeamName = &tpi.Slack.TeamName
	case "VICTOR_OPS":
		if tpi.VictorOps == nil || tpi.VictorOpsSecrets == nil {
			return nil, errors.New("missing Victor Ops settings")
		}
		ai.ApiKey = &tpi.VictorOpsSecrets.APIKey
		ai.RoutingKey = &tpi.VictorOps.RoutingKey
	case "WEBHOOK":
		if tpi.Webhook == nil || tpi.WebhookSecrets == nil {
			return nil, errors.New("missing Webhook settings")
		}
		ai.Url = &tpi.WebhookSecrets.URL
		ai.Secret = &tpi.WebhookSecrets.Secret
	default:
		return nil, fmt.Errorf("%w %v", ErrUnsupportedIntegrationType, tpi.Type)
	}
	return ai, nil
}

func fromAtlas(ai *admin.ThirdPartyIntegration) (*ThirdPartyIntegration, error) {
	tpi := &ThirdPartyIntegration{
		AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
			Type: ai.GetType(),
		},
		ID: ai.GetId(),
	}
	switch ai.GetType() {
	case "DATADOG":
		tpi.DatadogSecrets = &DatadogSecrets{
			APIKey: ai.GetApiKey(),
		}
		tpi.Datadog = &akov2.DatadogIntegration{
			Region:                       ai.GetRegion(),
			SendCollectionLatencyMetrics: encodeEnabled(ai.GetSendCollectionLatencyMetrics()),
			SendDatabaseMetrics:          encodeEnabled(ai.GetSendDatabaseMetrics()),
		}
	case "MICROSOFT_TEAMS":
		tpi.MicrosoftTeamsSecrets = &MicrosoftTeamsSecrets{
			WebhookUrl: ai.GetMicrosoftTeamsWebhookUrl(),
		}
		tpi.MicrosoftTeams = &akov2.MicrosoftTeamsIntegration{}
	case "NEW_RELIC":
		tpi.NewRelicSecrets = &NewRelicSecrets{
			AccountID:  ai.GetAccountId(),
			LicenseKey: ai.GetLicenseKey(),
			ReadToken:  ai.GetReadToken(),
			WriteToken: ai.GetWriteToken(),
		}
		tpi.NewRelic = &akov2.NewRelicIntegration{}
	case "OPS_GENIE":
		tpi.OpsGenieSecrets = &OpsGenieSecrets{
			APIKey: ai.GetApiKey(),
		}
		tpi.OpsGenie = &akov2.OpsGenieIntegration{
			Region: ai.GetRegion(),
		}
	case "PAGER_DUTY":
		tpi.PagerDutySecrets = &PagerDutySecrets{
			ServiceKey: ai.GetServiceKey(),
		}
		tpi.PagerDuty = &akov2.PagerDutyIntegration{
			Region: ai.GetRegion(),
		}
	case "PROMETHEUS":
		tpi.Prometheus = &akov2.PrometheusIntegration{
			Enabled:          encodeEnabled(ai.GetEnabled()),
			ServiceDiscovery: ai.GetServiceDiscovery(),
		}
		tpi.PrometheusSecrets = &PrometheusSecrets{
			Username: ai.GetUsername(),
			Password: ai.GetPassword(),
		}
	case "SLACK":
		tpi.Slack = &akov2.SlackIntegration{
			ChannelName: ai.GetChannelName(),
			TeamName:    ai.GetTeamName(),
		}
		tpi.SlackSecrets = &SlackSecrets{
			APIToken: ai.GetApiToken(),
		}
	case "VICTOR_OPS":
		tpi.VictorOps = &akov2.VictorOpsIntegration{
			RoutingKey: ai.GetRoutingKey(),
		}
		tpi.VictorOpsSecrets = &VictorOpsSecrets{
			APIKey: ai.GetApiKey(),
		}
	case "WEBHOOK":
		tpi.Webhook = &akov2.WebhookIntegration{}
		tpi.WebhookSecrets = &WebhookSecrets{
			URL:    ai.GetUrl(),
			Secret: ai.GetSecret(),
		}
	default:
		return nil, fmt.Errorf("%w %v", ErrUnsupportedIntegrationType, tpi.Type)
	}
	return tpi, nil
}

func assertType(typeName string) error {
	switch typeName {
	case "DATADOG", "MICROSOFT_TEAMS", "NEW_RELIC",
		"OPS_GENIE", "PAGER_DUTY", "PROMETHEUS", "SLACK",
		"VICTOR_OPS", "WEBHOOK":
		return nil
	default:
		return fmt.Errorf("%w %v", ErrUnsupportedIntegrationType, typeName)
	}
}

func isEnabled(field *string) bool {
	if field == nil {
		return false
	}
	return strings.EqualFold(*field, "enabled")
}

func encodeEnabled(on bool) *string {
	if on {
		return pointer.MakePtr("enabled")
	}
	return pointer.MakePtr("disabled")
}
