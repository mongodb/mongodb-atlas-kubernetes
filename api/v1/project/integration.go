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

package project

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
)

// Integration for the project between Atlas and a third party service.
// Deprecated: Migrate to the AtlasThirdPartyIntegration custom resource in accordance with the migration guide
// at https://www.mongodb.com/docs/atlas/operator/current/migrate-parameter-to-resource/#std-label-ak8so-migrate-ptr
type Integration struct {
	// Third Party Integration type such as Slack, New Relic, etc.
	// Each integration type requires a distinct set of configuration fields.
	// For example, if you set type to DATADOG, you must configure only datadog subfields.
	// +kubebuilder:validation:Enum=PAGER_DUTY;SLACK;DATADOG;NEW_RELIC;OPS_GENIE;VICTOR_OPS;FLOWDOCK;WEBHOOK;MICROSOFT_TEAMS;PROMETHEUS
	// +optional
	Type string `json:"type,omitempty"`
	// Reference to a Kubernetes Secret containing your Unique 40-hexadecimal digit string that identifies your New Relic license.
	// +optional
	LicenseKeyRef common.ResourceRefNamespaced `json:"licenseKeyRef,omitempty"`
	// Unique 40-hexadecimal digit string that identifies your New Relic account.
	// +optional
	AccountID string `json:"accountId,omitempty"`
	// Reference to a Kubernetes Secret containing the insert key associated with your New Relic account.
	// +optional
	WriteTokenRef common.ResourceRefNamespaced `json:"writeTokenRef,omitempty"`
	// Reference to a Kubernetes Secret containing the query key associated with your New Relic account.
	// +optional
	ReadTokenRef common.ResourceRefNamespaced `json:"readTokenRef,omitempty"`
	// Reference to a Kubernetes Secret containing your API Key for Datadog, OpsGenie or Victor Ops.
	// +optional
	APIKeyRef common.ResourceRefNamespaced `json:"apiKeyRef,omitempty"`
	// Region code indicating which regional API Atlas uses to access PagerDuty, Datadog, or OpsGenie.
	// +optional
	Region string `json:"region,omitempty"`
	// Reference to a Kubernetes Secret containing the service key associated with your PagerDuty account.
	// +optional
	ServiceKeyRef common.ResourceRefNamespaced `json:"serviceKeyRef,omitempty"`
	// Reference to a Kubernetes Secret containing the Key that allows Atlas to access your Slack account.
	// +optional
	APITokenRef common.ResourceRefNamespaced `json:"apiTokenRef,omitempty"`
	// Human-readable label that identifies your Slack team.
	// +optional
	TeamName string `json:"teamName,omitempty"`
	// Name of the Slack channel to which Atlas sends alert notifications.
	// +optional
	ChannelName string `json:"channelName,omitempty"`
	// Reference to a Kubernetes Secret containing the Routing key associated with your Splunk On-Call account.
	// Used for Victor Ops.
	// +optional
	RoutingKeyRef common.ResourceRefNamespaced `json:"routingKeyRef,omitempty"`
	// +optional
	FlowName string `json:"flowName,omitempty"`
	// +optional
	OrgName string `json:"orgName,omitempty"`
	// Endpoint web address to which Atlas sends notifications.
	// Used for Webhooks.
	// +optional
	URL string `json:"url,omitempty"`
	// Reference to a Kubernetes Secret containing the secret for your Webhook.
	// +optional
	SecretRef common.ResourceRefNamespaced `json:"secretRef,omitempty"`
	// +optional
	Name string `json:"name,omitempty"`
	// Endpoint web address of the Microsoft Teams webhook to which Atlas sends notifications.
	// +optional
	MicrosoftTeamsWebhookURL string `json:"microsoftTeamsWebhookUrl,omitempty"`
	// Human-readable label that identifies your Prometheus incoming webhook.
	// +optional
	UserName string `json:"username,omitempty"`
	// Reference to a Kubernetes Secret containing the password to allow Atlas to access your Prometheus account.
	// +optional
	PasswordRef common.ResourceRefNamespaced `json:"passwordRef,omitempty"`
	// Desired method to discover the Prometheus service.
	// +optional
	ServiceDiscovery string `json:"serviceDiscovery,omitempty"`
	// +optional
	Scheme string `json:"scheme,omitempty"`
	//
	// +optional
	Enabled bool `json:"enabled,omitempty"`
}

func (i Integration) Identifier() interface{} {
	return i.Type
}
