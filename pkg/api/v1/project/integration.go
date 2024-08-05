package project

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Integration struct {
	// Third Party Integration type such as Slack, New Relic, etc
	// +kubebuilder:validation:Enum=PAGER_DUTY;SLACK;DATADOG;NEW_RELIC;OPS_GENIE;VICTOR_OPS;FLOWDOCK;WEBHOOK;MICROSOFT_TEAMS;PROMETHEUS
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
	SecretRef common.ResourceRefNamespaced `json:"secretRef,omitempty"`
	// +optional
	Name string `json:"name,omitempty"`
	// +optional
	MicrosoftTeamsWebhookURL string `json:"microsoftTeamsWebhookUrl,omitempty"`
	// +optional
	UserName string `json:"username,omitempty"`
	// +optional
	PasswordRef common.ResourceRefNamespaced `json:"passwordRef,omitempty"`
	// +optional
	ServiceDiscovery string `json:"serviceDiscovery,omitempty"`
	// +optional
	Scheme string `json:"scheme,omitempty"`
	// +optional
	Enabled bool `json:"enabled,omitempty"`
}

func (i Integration) ToAtlas(ctx context.Context, c client.Client, defaultNS string) (result *admin.ThirdPartyIntegration, err error) {
	result = &admin.ThirdPartyIntegration{
		Type:                     pointer.MakePtr(i.Type),
		AccountId:                pointer.MakePtr(i.AccountID),
		Region:                   pointer.MakePtr(i.Region),
		TeamName:                 pointer.MakePtr(i.TeamName),
		ChannelName:              pointer.MakePtr(i.ChannelName),
		FlowName:                 pointer.MakePtr(i.FlowName),
		OrgName:                  pointer.MakePtr(i.OrgName),
		Url:                      pointer.MakePtr(i.URL),
		Name:                     pointer.MakePtr(i.Name),
		MicrosoftTeamsWebhookUrl: pointer.MakePtr(i.MicrosoftTeamsWebhookURL),
		Username:                 pointer.MakePtr(i.UserName),
		ServiceDiscovery:         pointer.MakePtr(i.ServiceDiscovery),
		Scheme:                   pointer.MakePtr(i.Scheme),
		Enabled:                  pointer.MakePtr(i.Enabled),
	}

	readPassword := func(passwordField common.ResourceRefNamespaced, target *string, errors *[]error) {
		if passwordField.Name == "" {
			return
		}

		*target, err = passwordField.ReadPassword(ctx, c, defaultNS)
		storeError(err, errors)
	}

	errorList := make([]error, 0)
	readPassword(i.LicenseKeyRef, result.LicenseKey, &errorList)
	readPassword(i.WriteTokenRef, result.WriteToken, &errorList)
	readPassword(i.ReadTokenRef, result.ReadToken, &errorList)
	readPassword(i.APIKeyRef, result.ApiKey, &errorList)
	readPassword(i.ServiceKeyRef, result.ServiceKey, &errorList)
	readPassword(i.APITokenRef, result.ApiToken, &errorList)
	readPassword(i.RoutingKeyRef, result.RoutingKey, &errorList)
	readPassword(i.SecretRef, result.Secret, &errorList)
	readPassword(i.PasswordRef, result.Password, &errorList)

	if len(errorList) != 0 {
		firstError := (errorList)[0]
		return nil, firstError
	}
	return result, nil
}

func (i Integration) Identifier() interface{} {
	return i.Type
}

func storeError(err error, errors *[]error) {
	if err != nil {
		*errors = append(*errors, err)
	}
}
