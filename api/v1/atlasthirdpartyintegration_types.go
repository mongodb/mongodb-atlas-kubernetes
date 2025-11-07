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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
)

func init() {
	SchemeBuilder.Register(&AtlasThirdPartyIntegration{}, &AtlasThirdPartyIntegrationList{})
}

// AtlasThirdPartyIntegration is the Schema for the atlas 3rd party integrations API.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:subresource:status
// +groupName:=atlas.mongodb.com
// +kubebuilder:resource:categories=atlas,shortName=atpi
type AtlasThirdPartyIntegration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasThirdPartyIntegrationSpec          `json:"spec,omitempty"`
	Status status.AtlasThirdPartyIntegrationStatus `json:"status,omitempty"`
}

func (np *AtlasThirdPartyIntegration) Credentials() *api.LocalObjectReference {
	return np.Spec.ConnectionSecret
}

func (atpi *AtlasThirdPartyIntegration) GetConditions() []metav1.Condition {
	return atpi.Status.Conditions
}

// +k8s:deepcopy-gen=true

// AtlasThirdPartyIntegrationSpec contains the expected configuration for an integration
// +kubebuilder:validation:XValidation:rule="(has(self.externalProjectRef) && !has(self.projectRef)) || (!has(self.externalProjectRef) && has(self.projectRef))",message="must define only one project reference through externalProjectRef or projectRef"
// +kubebuilder:validation:XValidation:rule="(has(self.externalProjectRef) && has(self.connectionSecret)) || !has(self.externalProjectRef)",message="must define a local connection secret when referencing an external project"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type.size() != 0",message="must define a type of integration"
// +kubebuilder:validation:XValidation:rule="!has(self.datadog) || (self.type == 'DATADOG' && has(self.datadog))",message="only DATADOG type may set datadog fields"
// +kubebuilder:validation:XValidation:rule="!has(self.microsoftTeams) || (self.type == 'MICROSOFT_TEAMS' && has(self.microsoftTeams))",message="only MICROSOFT_TEAMS type may set microsoftTeams fields"
// +kubebuilder:validation:XValidation:rule="!has(self.newRelic) || (self.type == 'NEW_RELIC' && has(self.newRelic))",message="only NEW_RELIC type may set newRelic fields"
// +kubebuilder:validation:XValidation:rule="!has(self.opsGenie) || (self.type == 'OPS_GENIE' && has(self.opsGenie))",message="only OPS_GENIE type may set opsGenie fields"
// +kubebuilder:validation:XValidation:rule="!has(self.prometheus) || (self.type == 'PROMETHEUS' && has(self.prometheus))",message="only PROMETHEUS type may set prometheus fields"
// +kubebuilder:validation:XValidation:rule="!has(self.pagerDuty) || (self.type == 'PAGER_DUTY' && has(self.pagerDuty))",message="only PAGER_DUTY type may set pagerDuty fields"
// +kubebuilder:validation:XValidation:rule="!has(self.slack) || (self.type == 'SLACK' && has(self.slack))",message="only SLACK type may set slack fields"
// +kubebuilder:validation:XValidation:rule="!has(self.victorOps) || (self.type == 'VICTOR_OPS' && has(self.victorOps))",message="only VICTOR_OPS type may set victorOps fields"
// +kubebuilder:validation:XValidation:rule="!has(self.webhook) || (self.type == 'WEBHOOK' && has(self.webhook))",message="only WEBHOOK type may set webhook fields"
type AtlasThirdPartyIntegrationSpec struct {
	ProjectDualReference `json:",inline"`

	// Type of the integration.
	// +kubebuilder:validation:Enum:=DATADOG;MICROSOFT_TEAMS;NEW_RELIC;OPS_GENIE;PAGER_DUTY;PROMETHEUS;SLACK;VICTOR_OPS;WEBHOOK
	// +kubebuilder:validation:Required
	Type string `json:"type"`

	// Datadog contains the config fields for Datadog's Integration.
	// +kubebuilder:validation:Optional
	Datadog *DatadogIntegration `json:"datadog,omitempty"`

	// MicrosoftTeams contains the config fields for Microsoft Teams's Integration.
	// +kubebuilder:validation:Optional
	MicrosoftTeams *MicrosoftTeamsIntegration `json:"microsoftTeams,omitempty"`

	// NewRelic contains the config fields for New Relic's Integration.
	// +kubebuilder:validation:Optional
	NewRelic *NewRelicIntegration `json:"newRelic,omitempty"`

	// OpsGenie contains the config fields for Ops Genie's Integration.
	// +kubebuilder:validation:Optional
	OpsGenie *OpsGenieIntegration `json:"opsGenie,omitempty"`

	// PagerDuty contains the config fields for PagerDuty's Integration.
	// +kubebuilder:validation:Optional
	PagerDuty *PagerDutyIntegration `json:"pagerDuty,omitempty"`

	// Prometheus contains the config fields for Prometheus's Integration.
	// +kubebuilder:validation:Optional
	Prometheus *PrometheusIntegration `json:"prometheus,omitempty"`

	// Slack contains the config fields for Slack's Integration.
	// +kubebuilder:validation:Optional
	Slack *SlackIntegration `json:"slack,omitempty"`

	// VictorOps contains the config fields for VictorOps's Integration.
	// +kubebuilder:validation:Optional
	VictorOps *VictorOpsIntegration `json:"victorOps,omitempty"`

	// Webhook contains the config fields for Webhook's Integration.
	// +kubebuilder:validation:Optional
	Webhook *WebhookIntegration `json:"webhook,omitempty"`
}

type DatadogIntegration struct {
	// APIKeySecretRef holds the name of a secret containing the Datadog API key.
	// +kubebuilder:validation:Required
	APIKeySecretRef api.LocalObjectReference `json:"apiKeySecretRef"`

	// Region is the Datadog region
	// +kubebuilder:validation:Required
	Region string `json:"region"`

	// SendCollectionLatencyMetrics toggles sending collection latency metrics.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=enabled;disabled
	// +kubebuilder:default:=disabled
	SendCollectionLatencyMetrics *string `json:"sendCollectionLatencyMetrics"`

	// SendDatabaseMetrics toggles sending database metrics,
	// including database and collection names
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=enabled;disabled
	// +kubebuilder:default:=disabled
	SendDatabaseMetrics *string `json:"sendDatabaseMetrics"`
}

type MicrosoftTeamsIntegration struct {
	// URLSecretRef holds the name of a secret containing the Microsoft Teams secret URL.
	// +kubebuilder:validation:Required
	URLSecretRef api.LocalObjectReference `json:"urlSecretRef"`
}

type NewRelicIntegration struct {
	// CredentialsSecretRef holds the name of a secret containing new relic's credentials:
	// account id, license key, read and write tokens.
	// +kubebuilder:validation:Required
	CredentialsSecretRef api.LocalObjectReference `json:"credentialsSecretRef"`
}

type OpsGenieIntegration struct {
	// APIKeySecretRef holds the name of a secret containing Ops Genie's API key.
	// +kubebuilder:validation:Required
	APIKeySecretRef api.LocalObjectReference `json:"apiKeySecretRef"`

	// Region is the Ops Genie region.
	// +kubebuilder:validation:Required
	Region string `json:"region"`
}

type PagerDutyIntegration struct {
	// ServiceKeySecretRef holds the name of a secret containing Pager Duty service key.
	// +kubebuilder:validation:Required
	ServiceKeySecretRef api.LocalObjectReference `json:"serviceKeySecretRef"`

	// Region is the Pager Duty region.
	// +kubebuilder:validation:Required
	Region string `json:"region"`
}

type PrometheusIntegration struct {
	// Enabled is true when Prometheus integration is enabled.
	Enabled *string `json:"enabled"`

	// PrometheusCredentialsSecretRef holds the name of a secret containing the Prometheus.
	// username & password
	// +kubebuilder:validation:Required
	PrometheusCredentialsSecretRef api.LocalObjectReference `json:"prometheusCredentialsSecretRef"`

	// ServiceDiscovery to be used by Prometheus.
	// +kubebuilder:validation:Enum:=file;http
	// +kubebuilder:validation:Required
	ServiceDiscovery string `json:"serviceDiscovery"`
}

type SlackIntegration struct {
	// APITokenSecretRef holds the name of a secret containing the Slack API token.
	// +kubebuilder:validation:Required
	APITokenSecretRef api.LocalObjectReference `json:"apiTokenSecretRef"`

	// ChannelName to be used by Prometheus.
	// +kubebuilder:validation:Required
	ChannelName string `json:"channelName"`

	// TeamName flags whether Prometheus integration is enabled.
	// +kubebuilder:validation:Required
	TeamName string `json:"teamName"`
}

type VictorOpsIntegration struct {
	// RoutingKey holds VictorOps routing key.
	// +kubebuilder:validation:Required
	RoutingKey string `json:"routingKey"`

	// APIKeySecretRef is the name of a secret containing Victor Ops API key.
	// +kubebuilder:validation:Required
	APIKeySecretRef api.LocalObjectReference `json:"apiKeySecretRef"`
}

type WebhookIntegration struct {
	// URLSecretRef holds the name of a secret containing Webhook URL and secret.
	// +kubebuilder:validation:Required
	URLSecretRef api.LocalObjectReference `json:"urlSecretRef"`
}

// +kubebuilder:object:root=true

// AtlasThirdPartyIntegrationList contains a list of Atlas Integrations.
type AtlasThirdPartyIntegrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasThirdPartyIntegration `json:"items"`
}

func (i *AtlasThirdPartyIntegration) ProjectDualRef() *ProjectDualReference {
	return &i.Spec.ProjectDualReference
}
