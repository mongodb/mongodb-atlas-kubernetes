package integrations

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
)

type AtlasIntegrationsService interface {
	Get(ctx context.Context, projectID string, integrationType string) (*Integration, error)
	List(ctx context.Context, projectID string) ([]Integration, error)
	Create(ctx context.Context, projectID string, integrationType string, integration Integration, secrets map[string]string) error
	Update(ctx context.Context, projectID string, integrationType string, integration Integration, secrets map[string]string) error
	Delete(ctx context.Context, projectID string, integrationType string) error
}

type AtlasIntegrations struct {
	integrationsAPI admin.ThirdPartyIntegrationsApi
}

func NewAtlasIntegrationsAPIService(api admin.ThirdPartyIntegrationsApi) *AtlasIntegrations {
	return &AtlasIntegrations{integrationsAPI: api}
}

func (i *AtlasIntegrations) Get(ctx context.Context, projectID string, integrationType string) (*Integration, error) {
	integration, _, err := i.integrationsAPI.GetThirdPartyIntegration(ctx, projectID, integrationType).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get integration from Atlas: %w", err)
	}
	return fromAtlas(integration)
}

func (i *AtlasIntegrations) List(ctx context.Context, projectID string) ([]Integration, error) {
	paginatedIntegrations, _, err := i.integrationsAPI.ListThirdPartyIntegrations(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to list integrations from Atlas: %w", err)
	}
	list := make([]Integration, 0, len(paginatedIntegrations.GetResults()))
	for _, integration := range paginatedIntegrations.GetResults() {
		int, err := fromAtlas(&integration)
		if err != nil {
			return nil, err
		}
		list = append(list, *int)
	}
	return list, nil
}

func (i *AtlasIntegrations) Create(ctx context.Context, projectID string, integrationType string, integration Integration, secrets map[string]string) error {
	atlasIntegration := toAtlas(integration, secrets)

	_, _, err := i.integrationsAPI.CreateThirdPartyIntegration(ctx, integrationType, projectID, atlasIntegration).Execute()
	if err != nil {
		return fmt.Errorf("failed to create integration in Atlas: %w", err)
	}
	return nil
}

func (i *AtlasIntegrations) Update(ctx context.Context, projectID string, integrationType string, integration Integration, secrets map[string]string) error {
	atlasIntegration := toAtlas(integration, secrets)

	_, _, err := i.integrationsAPI.UpdateThirdPartyIntegration(ctx, integrationType, projectID, atlasIntegration).Execute()
	if err != nil {
		return fmt.Errorf("failed to update integration in Atlas: %w", err)
	}
	return nil
}

func (i *AtlasIntegrations) Delete(ctx context.Context, projectID string, integrationType string) error {
	_, _, err := i.integrationsAPI.DeleteThirdPartyIntegration(ctx, integrationType, projectID).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete integration in Atlas: %w", err)
	}
	return nil
}
