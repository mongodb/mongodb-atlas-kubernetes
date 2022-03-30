package project

import (
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/compat"
)

type IntegrationType string

const (
	PagerDuty       IntegrationType = "PAGER_DUTY"
	Slack           IntegrationType = "SLACK"
	Datadog         IntegrationType = "DATADOG"
	NewRelic        IntegrationType = "NEW_RELIC"
	Opsgenie        IntegrationType = "OPS_GENIE"
	VictorOps       IntegrationType = "VICTOR_OPS"
	Flowdock        IntegrationType = "FLOWDOCK"
	WebhookSettings IntegrationType = "WEBHOOK"
	MicrosoftTeams  IntegrationType = "MICROSOFT_TEAMS"
)

type Intergation struct {
	// Third Party Integration type such as Slack, New Relic, etc
	// +kubebuilder:validation:UniqueItems
	// +kubebuilder:validation:Enum=PAGER_DUTY;SLACK;DATADOG;NEW_RELIC;OPS_GENIE;VICTOR_OPS;FLOWDOCK;WEBHOOK;MICROSOFT_TEAMS
	// +optional
	Type IntegrationType `json:"type,omitempty"`
	// +optional
	LicenseKey string `json:"licenseKey,omitempty"`
	// +optional
	AccountID string `json:"accountId,omitempty"`
	// +optional
	WriteToken string `json:"writeToken,omitempty"`
	// +optional
	ReadToken string `json:"readToken,omitempty"`
	// +optional
	APIKey string `json:"apiKey,omitempty"`
	// +optional
	Region string `json:"region,omitempty"`
	// +optional
	ServiceKey string `json:"serviceKey,omitempty"`
	// +optional
	APIToken string `json:"apiToken,omitempty"`
	// +optional
	TeamName string `json:"teamName,omitempty"`
	// +optional
	ChannelName string `json:"channelName,omitempty"`
	// +optional
	RoutingKey string `json:"routingKey,omitempty"`
	// +optional
	FlowName string `json:"flowName,omitempty"`
	// +optional
	OrgName string `json:"orgName,omitempty"`
	// +optional
	URL string `json:"url,omitempty"`
	// +optional
	Secret string `json:"secret,omitempty"`
}

func (i Intergation) ToAtlas() (*mongodbatlas.ThirdPartyIntegration, error) {
	result := &mongodbatlas.ThirdPartyIntegration{}
	err := compat.JSONCopy(result, i)
	return result, err
}

// TODO identifier?
func (i Intergation) Identifier() interface{} {
	return i.Type
}
