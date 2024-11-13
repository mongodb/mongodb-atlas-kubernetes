package datafederation

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
)

type DataFederationPrivateEndpointService interface {
	List(ctx context.Context, projectID string) ([]PrivateEndpoint, error)
	Create(context.Context, *PrivateEndpoint) error
	Delete(context.Context, *PrivateEndpoint) error
}

type PrivateEndpoints struct {
	api admin.DataFederationApi
}

func NewPrivateEndpointService(api admin.DataFederationApi) *PrivateEndpoints {
	return &PrivateEndpoints{api: api}
}

func (d *PrivateEndpoints) List(ctx context.Context, projectID string) ([]PrivateEndpoint, error) {
	paginatedResponse, _, err := d.api.ListDataFederationPrivateEndpoints(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to list data federation private endpoints from Atlas: %w", err)
	}

	return endpointsFromAtlas(projectID, paginatedResponse.GetResults())
}

func (d *PrivateEndpoints) Create(ctx context.Context, aep *PrivateEndpoint) error {
	_, _, err := d.api.CreateDataFederationPrivateEndpoint(ctx, aep.ProjectID, endpointToAtlas(aep)).Execute()
	if err != nil {
		return fmt.Errorf("failed to create data federation private endpoint: %w", err)
	}

	return nil
}

func (d *PrivateEndpoints) Delete(ctx context.Context, aep *PrivateEndpoint) error {
	_, _, err := d.api.DeleteDataFederationPrivateEndpoint(ctx, aep.ProjectID, aep.EndpointID).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete data federation private endpoint: %w", err)
	}

	return nil
}
