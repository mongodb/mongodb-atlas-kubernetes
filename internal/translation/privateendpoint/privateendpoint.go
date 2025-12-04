// Copyright 2024 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package privateendpoint

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312009/admin"
	"golang.org/x/exp/slices"
)

const (
	ErrorServiceNotFound = "PRIVATE_ENDPOINT_SERVICE_NOT_FOUND"
)

type PrivateEndpointService interface {
	ListPrivateEndpoints(ctx context.Context, projectID, provider string) ([]EndpointService, error)
	GetPrivateEndpoint(ctx context.Context, projectID, provider, ID string) (EndpointService, error)
	CreatePrivateEndpointService(ctx context.Context, projectID string, peService EndpointService) (EndpointService, error)
	DeleteEndpointService(ctx context.Context, projectID, provider, ID string) error
	CreatePrivateEndpointInterface(ctx context.Context, projectID, provider, serviceID, gcpProjectID string, peInterface EndpointInterface) (EndpointInterface, error)
	DeleteEndpointInterface(ctx context.Context, projectID, provider, serviceID, ID string) error
}

type PrivateEndpoint struct {
	api admin.PrivateEndpointServicesApi
}

func (pe *PrivateEndpoint) ListPrivateEndpoints(ctx context.Context, projectID, provider string) ([]EndpointService, error) {
	services, _, err := pe.api.ListPrivateEndpointService(ctx, projectID, provider).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the list of private endpoints: %w", err)
	}

	peServices := make([]EndpointService, 0, len(services))
	for _, service := range services {
		interfaceIDs := getInterfacesIDs(&service)
		peInterfaces := make([]EndpointInterface, 0, len(interfaceIDs))

		for _, interfaceID := range interfaceIDs {
			peInterface, err := pe.getEndpointInterfaces(ctx, projectID, provider, service.GetId(), interfaceID)
			if err != nil {
				return nil, err
			}

			if peInterface != nil {
				peInterfaces = append(peInterfaces, peInterface)
			}
		}

		slices.SortFunc(peInterfaces, func(a, b EndpointInterface) int {
			return strings.Compare(a.InterfaceID(), b.InterfaceID())
		})

		peServices = append(
			peServices,
			serviceFromAtlas(&service, peInterfaces),
		)
	}

	return peServices, nil
}

func (pe *PrivateEndpoint) GetPrivateEndpoint(ctx context.Context, projectID, provider, ID string) (EndpointService, error) {
	service, _, err := pe.api.GetPrivateEndpointService(ctx, projectID, provider, ID).
		Execute()
	if admin.IsErrorCode(err, ErrorServiceNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the private endpoint: %w", err)
	}

	interfaceIDs := getInterfacesIDs(service)
	peInterfaces := make([]EndpointInterface, 0, len(interfaceIDs))

	for _, interfaceID := range interfaceIDs {
		peInterface, err := pe.getEndpointInterfaces(ctx, projectID, provider, service.GetId(), interfaceID)
		if err != nil {
			return nil, err
		}

		if peInterface != nil {
			peInterfaces = append(peInterfaces, peInterface)
		}
	}

	slices.SortFunc(peInterfaces, func(a, b EndpointInterface) int {
		return strings.Compare(a.InterfaceID(), b.InterfaceID())
	})

	return serviceFromAtlas(service, peInterfaces), nil
}

func (pe *PrivateEndpoint) CreatePrivateEndpointService(ctx context.Context, projectID string, peService EndpointService) (EndpointService, error) {
	service, _, err := pe.api.CreatePrivateEndpointService(ctx, projectID, serviceCreateToAtlas(peService)).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to create the private endpoint service: %w", err)
	}

	return serviceFromAtlas(service, []EndpointInterface{}), nil
}

func (pe *PrivateEndpoint) DeleteEndpointService(ctx context.Context, projectID, provider, ID string) error {
	_, err := pe.api.DeletePrivateEndpointService(ctx, projectID, provider, ID).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete the private endpoint service: %w", err)
	}

	return nil
}

func (pe *PrivateEndpoint) CreatePrivateEndpointInterface(ctx context.Context, projectID, provider, serviceID, gcpProjectID string, peInterface EndpointInterface) (EndpointInterface, error) {
	i, _, err := pe.api.CreatePrivateEndpoint(ctx, projectID, provider, serviceID, interfaceCreateToAtlas(peInterface, gcpProjectID)).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to create the private endpoint interface: %w", err)
	}

	return interfaceFromAtlas(i), nil
}

func (pe *PrivateEndpoint) DeleteEndpointInterface(ctx context.Context, projectID, provider, serviceID, ID string) error {
	_, err := pe.api.DeletePrivateEndpoint(ctx, projectID, provider, ID, serviceID).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete the private endpoint interface: %w", err)
	}

	return nil
}

func (pe *PrivateEndpoint) getEndpointInterfaces(ctx context.Context, projectID, provider, serviceID, ID string) (EndpointInterface, error) {
	i, _, err := pe.api.GetPrivateEndpoint(ctx, projectID, provider, ID, serviceID).Execute()
	if admin.IsErrorCode(err, ErrorServiceNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the private endpoint interface: %w", err)
	}

	return interfaceFromAtlas(i), nil
}

func getInterfacesIDs(peService *admin.EndpointService) []string {
	switch peService.GetCloudProvider() {
	case ProviderAWS:
		return peService.GetInterfaceEndpoints()
	case ProviderAzure:
		return peService.GetPrivateEndpoints()
	case ProviderGCP:
		return peService.GetEndpointGroupNames()
	}

	return nil
}

func NewPrivateEndpointAPI(api admin.PrivateEndpointServicesApi) PrivateEndpointService {
	return &PrivateEndpoint{
		api: api,
	}
}
