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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/v1/status"
)

func init() {
	SchemeBuilder.Register(&AtlasThirdPartyIntegration{}, &AtlasThirdPartyIntegrationList{})
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AtlasThirdPartyIntegration is the Schema for the atlas 3rd party inegrations API.
type AtlasThirdPartyIntegration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasThirdPartyIntegrationSpec          `json:"spec,omitempty"`
	Status status.AtlasThirdPartyIntegrationStatus `json:"status,omitempty"`
}

func (in *AtlasThirdPartyIntegration) GetStatus() api.Status {
	return in.Status
}

func (b *AtlasThirdPartyIntegration) UpdateStatus(conditions []api.Condition, _ ...api.Option) {
	b.Status.Conditions = conditions
	b.Status.ObservedGeneration = b.ObjectMeta.Generation
}

func (np *AtlasThirdPartyIntegration) Credentials() *api.LocalObjectReference {
	return np.Spec.ConnectionSecret
}

// +k8s:deepcopy-gen=true

// AtlasThirdPartyIntegrationSpec contains the expected configuration for an integration
type AtlasThirdPartyIntegrationSpec struct {
	akov2.ProjectDualReference `json:",inline"`

	// Type of the integration
	Type string `json:"type"`

	// Datadog contains the config fields for Datadog's Integration
	Datadog *DatadogIntegration `json:"datadog,omitempty"`

	// MicrosoftTeams contains the config fields for Microsoft Teams's Integration
	MicrosoftTeams *MicrosoftTeamsIntegration `json:"microsoftTeams,omitempty"`

	// NewRelic contains the config fields for New Relic's Integration
	NewRelic *NewRelicIntegration `json:"newRelic,omitempty"`

	// OpsGenie contains the config fields for Ops Genie's Integration
	OpsGenie *OpsGenieIntegration `json:"opsGenie,omitempty"`

	// PagerDuty contains the config fields for PagerDuty's Integration
	PagerDuty *PagerDutyIntegration `json:"pagerDuty,omitempty"`

	// Prometheus contains the config fields for Prometheus's Integration
	Prometheus *PrometheusIntegration `json:"prometheus,omitempty"`

	// Slack contains the config fields for Slack's Integration
	Slack *SlackIntegration `json:"slack,omitempty"`

	// VictorOps contains the config fields for VictorOps's Integration
	VictorOps *VictorOpsIntegration `json:"victorOps,omitempty"`

	// Webhook contains the config fields for Webhook's Integration
	Webhook *WebhookIntegration `json:"webhook,omitempty"`
}

// +k8s:deepcopy-gen=true

type DatadogIntegration struct {
	// APIKeySecret holds the name of a secret containing the datadog api key
	APIKeySecret api.LocalObjectReference `json:"apiKeySecret"`

	// Region is the Datadog region
	Region string `json:"region"`

	// SendCollectionLatencyMetrics toggles sending collection latency metrics
	SendCollectionLatencyMetrics *string `json:"sendCollectionLatencyMetrics"`

	// SendDatabaseMetrics toggles sending database metrics,
	// including database and collection names
	SendDatabaseMetrics *string `json:"sendDatabaseMetrics"`
}

// +k8s:deepcopy-gen=true

type MicrosoftTeamsIntegration struct {
	// URLSecret holds the name of a secret containing the microsoft teams secret URL
	URLSecret api.LocalObjectReference `json:"urlSecret"`
}

// +k8s:deepcopy-gen=true

type NewRelicIntegration struct {
	// CredentialsSecret holds the name of a secret containing new relic's credentials:
	// account id, license key, read and write tokens
	CredentialsSecret api.LocalObjectReference `json:"credentialsSecret"`
}

// +k8s:deepcopy-gen=true

type OpsGenieIntegration struct {
	// APIKeySecret holds the name of a secret containing Ops Genie's API key
	APIKeySecret api.LocalObjectReference `json:"apiKeySecret"`

	// Region is the Ops Genie region
	Region string `json:"region"`
}

// +k8s:deepcopy-gen=true

type PagerDutyIntegration struct {
	// ServiceKeySecret holds the name of a secret containing Pager Duty service key
	ServiceKeySecret api.LocalObjectReference `json:"serviceKeySecret"`

	// Region is the Pager Duty region
	Region string `json:"region"`
}

// +k8s:deepcopy-gen=true

type PrometheusIntegration struct {
	// PrometheusCredentials holds the name of a secret containing the Prometheus
	// username & password
	PrometheusCredentials api.LocalObjectReference `json:"prometheusCredentials"`

	// ServiceDiscovery to be used by Prometheus
	ServiceDiscovery string `json:"serviceDiscovery"`
}

// +k8s:deepcopy-gen=true

type SlackIntegration struct {
	// APITokenSecret holds the name of a secret containing the Slack API token
	APITokenSecret api.LocalObjectReference `json:"apiTokenSecret"`

	// ChannelName to be used by Prometheus
	ChannelName string `json:"channelName"`

	// TeamName flags whether or not Prometheus integration is enabled
	TeamName string `json:"teamName"`
}

// +k8s:deepcopy-gen=true

type VictorOpsIntegration struct {
	// RoutingKey holds VictorOps routing key
	RoutingKey api.LocalObjectReference `json:"routingKey"`

	// APIKeySecret is the name of a secret containing Victor Ops API key
	APIKeySecret string `json:"apiKeySecret"`
}

// +k8s:deepcopy-gen=true

type WebhookIntegration struct {
	// URLSecret holds the name of a secret containing Webhook URL and secret
	URLSecret api.LocalObjectReference `json:"urlSecret"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AtlasThirdPartyIntegrationList contains a list of Atlas Integrations.
type AtlasThirdPartyIntegrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasThirdPartyIntegration `json:"items"`
}

func (i *AtlasThirdPartyIntegration) ProjectDualRef() *akov2.ProjectDualReference {
	return &i.Spec.ProjectDualReference
}
