package project

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
)

type Integration struct {
	// Third Party Integration type such as Slack, New Relic, etc
	// +kubebuilder:validation:Enum=PAGER_DUTY;SLACK;DATADOG;NEW_RELIC;OPS_GENIE;VICTOR_OPS;FLOWDOCK;WEBHOOK;MICROSOFT_TEAMS;PROMETHEUS
	// +optional
	Type string `json:"type,omitempty"`
	// Reference to a secret containing a unique 40-hexadecimal digit string that idenfies your New Relic license.
	// +optional
	LicenseKeyRef common.ResourceRefNamespaced `json:"licenseKeyRef,omitempty"`
	// Unique 40-hexadecimal digit string identifying your New Relic account.
	// +optional
	AccountID string `json:"accountId,omitempty"`
	// Reference to a secret containing an insert key associated with your New Relic account.
	// +optional
	WriteTokenRef common.ResourceRefNamespaced `json:"writeTokenRef,omitempty"`
	// Reference to a secret containing a query key associated with your New Relic account.
	// +optional
	ReadTokenRef common.ResourceRefNamespaced `json:"readTokenRef,omitempty"`
	// Reference to a secret containing the API key for your Datadog, OpsGenie, or Victor Ops account.
	// +optional
	APIKeyRef common.ResourceRefNamespaced `json:"apiKeyRef,omitempty"`
	// Region code indicating which region to use for Datadog, OpsGenie, or PagerDuty.
	// +optional
	Region string `json:"region,omitempty"`
	// Reference to a secret containing the service key for your PagerDuty account.
	// +optional
	ServiceKeyRef common.ResourceRefNamespaced `json:"serviceKeyRef,omitempty"`
	// Reference to a secret containing the API token for your Slack account.
	// +optional
	APITokenRef common.ResourceRefNamespaced `json:"apiTokenRef,omitempty"`
	// Human-readable label that identies your Slack team.
	// +optional
	TeamName string `json:"teamName,omitempty"`
	// Name of the Slack channel to which notifications are sent.
	// +optional
	ChannelName string `json:"channelName,omitempty"`
	// Reference to a secret containing your routing key for Splunk On-Call, for use with Victor Ops.
	// +optional
	RoutingKeyRef common.ResourceRefNamespaced `json:"routingKeyRef,omitempty"`
	// DEPRECATED: No longer available in Atlas API.
	// +optional
	FlowName string `json:"flowName,omitempty"`
	// DEPRECATED: No longer available in Atlas API.
	// +optional
	OrgName string `json:"orgName,omitempty"`
	// Endpoint web address to which notifications are sent.
	// +optional
	URL string `json:"url,omitempty"`
	// Reference to a secret containing the secret that secures your webhook.
	// +optional
	SecretRef common.ResourceRefNamespaced `json:"secretRef,omitempty"`
	// DEPRECATED: No longer available in Atlas API.
	// +optional
	Name string `json:"name,omitempty"`
	// Endpoint web address of the Microsoft Teams webhook to which notifications are sent.
	// +optional
	MicrosoftTeamsWebhookURL string `json:"microsoftTeamsWebhookUrl,omitempty"`
	// Human-readable label that identifies your Prometheus incoming webhook.
	// +optional
	UserName string `json:"username,omitempty"`
	// Reference to a secret containing the password for your Prometheus account.
	// +optional
	PasswordRef common.ResourceRefNamespaced `json:"passwordRef,omitempty"`
	// Desired method to discover the Prometheus service.
	// +optional
	ServiceDiscovery string `json:"serviceDiscovery,omitempty"`
	// DEPRECATED: No longer available in Atlas API.
	// +optional
	Scheme string `json:"scheme,omitempty"`
	// Flag indicating whether the Prometheus integration is activated.
	// +optional
	Enabled bool `json:"enabled,omitempty"`
}
