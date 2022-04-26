package project

import (
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Integration struct {
	// Third Party Integration type such as Slack, New Relic, etc
	// +kubebuilder:validation:Enum=PAGER_DUTY;SLACK;DATADOG;NEW_RELIC;OPS_GENIE;VICTOR_OPS;FLOWDOCK;WEBHOOK;MICROSOFT_TEAMS
	// +optional
	Type string `json:"type,omitempty"`
	// +optional
	LicenseKeyRef common.ResourceRefNamespaced `json:"licenseKeyRef,omitempty"`
	// +optional
	AccountID string `json:"accountId,omitempty"`
	// +optional
	WriteTokenRef common.ResourceRefNamespaced `json:"writeTokenRef,omitempty"`
	// +optional
	ReadTokenRef common.ResourceRefNamespaced `json:"readTokenRef,omitempty"`
	// +optional
	APIKeyRef common.ResourceRefNamespaced `json:"apiKeyRef,omitempty"`
	// +optional
	Region string `json:"region,omitempty"`
	// +optional
	ServiceKeyRef common.ResourceRefNamespaced `json:"serviceKeyRef,omitempty"`
	// +optional
	APITokenRef common.ResourceRefNamespaced `json:"apiTokenRef,omitempty"`
	// +optional
	TeamName string `json:"teamName,omitempty"`
	// +optional
	ChannelName string `json:"channelName,omitempty"`
	// +optional
	RoutingKeyRef common.ResourceRefNamespaced `json:"routingKeyRef,omitempty"`
	// +optional
	FlowName string `json:"flowName,omitempty"`
	// +optional
	OrgName string `json:"orgName,omitempty"`
	// +optional
	URL string `json:"url,omitempty"`
	// +optional
	SecretRef common.ResourceRefNamespaced `json:"secret,omitempty"`
}

func (i Integration) ToAtlas(defaultNS string, c client.Client) *mongodbatlas.ThirdPartyIntegration {
	result := mongodbatlas.ThirdPartyIntegration{}
	result.Type = i.Type
	result.LicenseKey, _ = i.LicenseKeyRef.ReadPassword(c, defaultNS)
	result.AccountID = i.AccountID
	result.WriteToken, _ = i.WriteTokenRef.ReadPassword(c, defaultNS)
	result.ReadToken, _ = i.ReadTokenRef.ReadPassword(c, defaultNS)
	result.APIKey, _ = i.APIKeyRef.ReadPassword(c, defaultNS)
	result.Region = i.Region
	result.ServiceKey, _ = i.ServiceKeyRef.ReadPassword(c, defaultNS)
	result.APIToken, _ = i.APITokenRef.ReadPassword(c, defaultNS)
	result.TeamName = i.TeamName
	result.ChannelName = i.ChannelName
	result.RoutingKey, _ = i.RoutingKeyRef.ReadPassword(c, defaultNS)
	result.FlowName = i.FlowName
	result.OrgName = i.OrgName
	result.URL = i.URL
	result.Secret, _ = i.SecretRef.ReadPassword(c, defaultNS)
	return &result
}

func (i Integration) Identifier() interface{} {
	return i.Type
}
