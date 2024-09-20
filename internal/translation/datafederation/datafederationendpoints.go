package datafederation

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
)

type DatafederationPrivateEndpointService interface {
	List(ctx context.Context, projectID string) ([]*DatafederationPrivateEndpointEntry, error)
	Create(context.Context, *DatafederationPrivateEndpointEntry) error
	Delete(context.Context, *DatafederationPrivateEndpointEntry) error
}

type DatafederationPrivateEndpoints struct {
	api admin.DataFederationApi
}

func NewDatafederationPrivateEndpointService(ctx context.Context, provider atlas.Provider, secretRef *types.NamespacedName, log *zap.SugaredLogger) (*DatafederationPrivateEndpoints, error) {
	client, err := translation.NewVersionedClient(ctx, provider, secretRef, log)
	if err != nil {
		return nil, fmt.Errorf("failed to create versioned client: %w", err)
	}
	return &DatafederationPrivateEndpoints{client.DataFederationApi}, nil
}

func (i *DatafederationPrivateEndpoints) List(ctx context.Context, projectID string) ([]*DatafederationPrivateEndpointEntry, error) {
	paginatedResponse, _, err := i.api.ListDataFederationPrivateEndpoints(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to list data federation private endpoints from Atlas: %w", err)
	}

	return endpointsFromAtlas(paginatedResponse.GetResults(), projectID)
}

func (service *DatafederationPrivateEndpoints) Create(ctx context.Context, aep *DatafederationPrivateEndpointEntry) error {
	ep := endpointToAtlas(aep)
	_, _, err := service.api.CreateDataFederationPrivateEndpoint(ctx, aep.ProjectID, ep).Execute()
	if err != nil {
		return fmt.Errorf("failed to create data federation private endpoint: %w", err)
	}
	return nil
}

func (service *DatafederationPrivateEndpoints) Delete(ctx context.Context, aep *DatafederationPrivateEndpointEntry) error {
	_, _, err := service.api.DeleteDataFederationPrivateEndpoint(ctx, aep.ProjectID, aep.EndpointID).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete data federation private endpoint: %w", err)
	}
	return nil
}
