package integrations

import (
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/cmp"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/project"
)

type Integration struct {
	project.Integration
}

func (i Integration) Identifier() interface{} {
	return i.Type
}

func NewIntegration(spec *project.Integration) (*Integration, error) {
	if spec == nil {
		return nil, nil
	}
	if err := cmp.Normalize(spec); err != nil {
		return nil, fmt.Errorf("failed to normalize integration: %w", err)
	}
	return &Integration{Integration: *spec}, nil
}

func NewIntegrations(specs []project.Integration) ([]Integration, error) {
	list := make([]Integration, 0, len(specs))
	for _, spec := range specs {
		integration, err := NewIntegration(&spec)
		if err != nil {
			return nil, err
		}
		list = append(list, *integration)
	}
	return list, nil
}

func fromAtlas(in *admin.ThirdPartyIntegration) (*Integration, error) {
	return NewIntegration(
		&project.Integration{
			Type:                     in.GetType(),
			AccountID:                in.GetAccountId(),
			Region:                   in.GetRegion(),
			TeamName:                 in.GetTeamName(),
			ChannelName:              in.GetChannelName(),
			URL:                      in.GetUrl(),
			MicrosoftTeamsWebhookURL: in.GetMicrosoftTeamsWebhookUrl(),
			UserName:                 in.GetUsername(),
			ServiceDiscovery:         in.GetServiceDiscovery(),
			Enabled:                  in.GetEnabled(),
		},
	)
}

func toAtlas(in Integration, secrets map[string]string) *admin.ThirdPartyIntegration {
	result := &admin.ThirdPartyIntegration{
		Type:                     pointer.MakePtr(in.Type),
		AccountId:                pointer.MakePtr(in.AccountID),
		Region:                   pointer.MakePtr(in.Region),
		TeamName:                 pointer.MakePtr(in.TeamName),
		ChannelName:              pointer.MakePtr(in.ChannelName),
		Url:                      pointer.MakePtr(in.URL),
		MicrosoftTeamsWebhookUrl: pointer.MakePtr(in.MicrosoftTeamsWebhookURL),
		Username:                 pointer.MakePtr(in.UserName),
		ServiceDiscovery:         pointer.MakePtr(in.ServiceDiscovery),
		Enabled:                  pointer.MakePtr(in.Enabled),
	}

	if licenseKey, ok := secrets["licenseKey"]; ok {
		result.SetLicenseKey(licenseKey)
	}
	if writeToken, ok := secrets["writeToken"]; ok {
		result.SetWriteToken(writeToken)
	}
	if readToken, ok := secrets["readToken"]; ok {
		result.SetReadToken(readToken)
	}
	if apiKey, ok := secrets["apiKey"]; ok {
		result.SetApiKey(apiKey)
	}
	if serviceKey, ok := secrets["serviceKey"]; ok {
		result.SetServiceKey(serviceKey)
	}
	if apiToken, ok := secrets["apiToken"]; ok {
		result.SetApiToken(apiToken)
	}
	if routingKey, ok := secrets["routingKey"]; ok {
		result.SetRoutingKey(routingKey)
	}
	if secret, ok := secrets["secret"]; ok {
		result.SetSecret(secret)
	}
	if password, ok := secrets["password"]; ok {
		result.SetPassword(password)
	}

	return result
}
