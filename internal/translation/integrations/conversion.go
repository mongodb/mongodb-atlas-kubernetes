package integrations

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/cmp"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
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

func toAtlas(in Integration, ctx context.Context, c client.Client, defaultNS string) (result *admin.ThirdPartyIntegration, err error) {
	result = &admin.ThirdPartyIntegration{
		Type:                     pointer.MakePtr(in.Type),
		AccountId:                pointer.MakePtr(in.AccountID),
		Region:                   pointer.MakePtr(in.Region),
		TeamName:                 pointer.MakePtr(in.TeamName),
		ChannelName:              pointer.MakePtr(in.ChannelName),
		Url:                      pointer.MakePtr(in.URL),
		MicrosoftTeamsWebhookUrl: pointer.MakePtr(in.MicrosoftTeamsWebhookURL),
		Username:                 pointer.MakePtr(in.UserName),
		ServiceDiscovery:         pointer.MakePtr(in.ServiceDiscovery),
		LicenseKey:               pointer.MakePtr(""),
		Enabled:                  pointer.MakePtr(in.Enabled),
	}

	readPassword := func(passwordField common.ResourceRefNamespaced, setFunc func(string), errors *[]error) {
		if passwordField.Name == "" {
			return
		}

		target, err := passwordField.ReadPassword(ctx, c, defaultNS)
		if err != nil {
			*errors = append(*errors, err)
		}
		setFunc(target)
	}

	errorList := make([]error, 0)
	readPassword(in.LicenseKeyRef, result.SetLicenseKey, &errorList)
	readPassword(in.WriteTokenRef, result.SetWriteToken, &errorList)
	readPassword(in.ReadTokenRef, result.SetReadToken, &errorList)
	readPassword(in.APIKeyRef, result.SetApiKey, &errorList)
	readPassword(in.ServiceKeyRef, result.SetServiceKey, &errorList)
	readPassword(in.APITokenRef, result.SetApiToken, &errorList)
	readPassword(in.RoutingKeyRef, result.SetRoutingKey, &errorList)
	readPassword(in.SecretRef, result.SetSecret, &errorList)
	readPassword(in.PasswordRef, result.SetPassword, &errorList)

	if len(errorList) != 0 {
		firstError := (errorList)[0]
		return nil, firstError
	}
	return result, nil
}
