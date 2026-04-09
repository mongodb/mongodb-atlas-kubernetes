// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package privateendpoint

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20250312018/admin"
	"go.mongodb.org/atlas-sdk/v20250312018/mockadmin"
)

func TestListPrivateEndpoints(t *testing.T) {
	tests := map[string]struct {
		provider                string
		mockListReturnFunc      func() ([]admin.EndpointService, *http.Response, error)
		mockInterfaceReturnFunc func() (*admin.PrivateLinkEndpoint, *http.Response, error)
		expectedPEs             []EndpointService
		expectedErr             error
	}{
		"failed to retrieve data": {
			provider: "AWS",
			mockListReturnFunc: func() ([]admin.EndpointService, *http.Response, error) {
				return nil, &http.Response{}, errors.New("atlas failed to list")
			},
			expectedErr: fmt.Errorf("failed to retrieve the list of private endpoints: %w", errors.New("atlas failed to list")),
		},
		"failed to retrieve existing interface for listed private endpoint service": {
			provider: "AWS",
			mockListReturnFunc: func() ([]admin.EndpointService, *http.Response, error) {
				return []admin.EndpointService{
					{
						CloudProvider:       "AWS",
						Id:                  new("pe-service-ID"),
						RegionName:          new("US_EAST_1"),
						Status:              new("AVAILABLE"),
						EndpointServiceName: new("aws/pe-service/name"),
						InterfaceEndpoints:  &[]string{"vpcpe-123456"},
					},
				}, &http.Response{}, nil
			},
			mockInterfaceReturnFunc: func() (*admin.PrivateLinkEndpoint, *http.Response, error) {
				return nil, &http.Response{}, errors.New("atlas failed to get")
			},
			expectedErr: fmt.Errorf("failed to retrieve the private endpoint interface: %w", errors.New("atlas failed to get")),
		},
		"list AWS private endpoints": {
			provider: "AWS",
			mockListReturnFunc: func() ([]admin.EndpointService, *http.Response, error) {
				return []admin.EndpointService{
					{
						CloudProvider:       "AWS",
						Id:                  new("pe-service-ID-1"),
						RegionName:          new("US_EAST_1"),
						Status:              new("AVAILABLE"),
						EndpointServiceName: new("aws/pe-service/name"),
						InterfaceEndpoints:  &[]string{"vpcpe-123456"},
					},
					{
						CloudProvider:       "AWS",
						Id:                  new("pe-service-ID-2"),
						RegionName:          new("US_EAST_2"),
						Status:              new("AVAILABLE"),
						EndpointServiceName: new("aws/pe-service/name"),
					},
				}, &http.Response{}, nil
			},
			mockInterfaceReturnFunc: func() (*admin.PrivateLinkEndpoint, *http.Response, error) {
				return &admin.PrivateLinkEndpoint{
					CloudProvider:       "AWS",
					ConnectionStatus:    new("AVAILABLE"),
					InterfaceEndpointId: new("vpcpe-123456"),
				}, &http.Response{}, nil
			},
			expectedPEs: []EndpointService{
				&AWSService{
					CommonEndpointService: CommonEndpointService{
						ID:            "pe-service-ID-1",
						CloudRegion:   "US_EAST_1",
						ServiceStatus: "AVAILABLE",
						Interfaces: EndpointInterfaces{
							&AWSInterface{
								CommonEndpointInterface{
									ID:              "vpcpe-123456",
									InterfaceStatus: "AVAILABLE",
								},
							},
						},
					},
					ServiceName: "aws/pe-service/name",
				},
				&AWSService{
					CommonEndpointService: CommonEndpointService{
						ID:            "pe-service-ID-2",
						CloudRegion:   "US_EAST_2",
						ServiceStatus: "AVAILABLE",
						Interfaces:    EndpointInterfaces{},
					},
					ServiceName: "aws/pe-service/name",
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			projectID := "project-ID"
			api := mockadmin.NewPrivateEndpointServicesApi(t)
			api.EXPECT().ListPrivateEndpointService(ctx, projectID, tt.provider).
				Return(admin.ListPrivateEndpointServiceApiRequest{ApiService: api})
			api.EXPECT().ListPrivateEndpointServiceExecute(mock.AnythingOfType("admin.ListPrivateEndpointServiceApiRequest")).
				Return(tt.mockListReturnFunc())

			if tt.mockInterfaceReturnFunc != nil {
				api.EXPECT().GetPrivateEndpoint(ctx, projectID, tt.provider, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(admin.GetPrivateEndpointApiRequest{ApiService: api})
				api.EXPECT().GetPrivateEndpointExecute(mock.AnythingOfType("admin.GetPrivateEndpointApiRequest")).
					Return(tt.mockInterfaceReturnFunc())
			}

			pe := &PrivateEndpoint{
				api: api,
			}

			items, err := pe.ListPrivateEndpoints(ctx, projectID, tt.provider)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedPEs, items)
		})
	}
}

func TestGetPrivateEndpoint(t *testing.T) {
	notFoundErr := &admin.GenericOpenAPIError{}
	notFoundErr.SetModel(admin.ApiError{ErrorCode: "PRIVATE_ENDPOINT_SERVICE_NOT_FOUND"})
	tests := map[string]struct {
		provider                string
		mockGetReturnFunc       func() (*admin.EndpointService, *http.Response, error)
		mockInterfaceReturnFunc func() (*admin.PrivateLinkEndpoint, *http.Response, error)
		expectedPE              EndpointService
		expectedErr             error
	}{
		"failed to retrieve data": {
			provider: "AWS",
			mockGetReturnFunc: func() (*admin.EndpointService, *http.Response, error) {
				return nil, &http.Response{}, errors.New("atlas failed to get")
			},
			expectedErr: fmt.Errorf("failed to retrieve the private endpoint: %w", errors.New("atlas failed to get")),
		},
		"service was not found": {
			provider: "AWS",
			mockGetReturnFunc: func() (*admin.EndpointService, *http.Response, error) {
				return nil, &http.Response{}, notFoundErr
			},
		},
		"failed to get interface for the service": {
			provider: "AZURE",
			mockGetReturnFunc: func() (*admin.EndpointService, *http.Response, error) {
				return &admin.EndpointService{
					CloudProvider:                "AZURE",
					Id:                           new("pe-service-ID"),
					RegionName:                   new("GERMANY_NORTH"),
					Status:                       new("AVAILABLE"),
					PrivateEndpoints:             &[]string{"long-azure-resource-ID"},
					PrivateLinkServiceName:       new("pls_name"),
					PrivateLinkServiceResourceId: new("long-azure-resource-ID"),
				}, &http.Response{}, nil
			},
			mockInterfaceReturnFunc: func() (*admin.PrivateLinkEndpoint, *http.Response, error) {
				return nil, &http.Response{}, errors.New("atlas failed to get")
			},
			expectedErr: fmt.Errorf("failed to retrieve the private endpoint interface: %w", errors.New("atlas failed to get")),
		},
		"get AZURE private endpoint": {
			provider: "AZURE",
			mockGetReturnFunc: func() (*admin.EndpointService, *http.Response, error) {
				return &admin.EndpointService{
					CloudProvider:                "AZURE",
					Id:                           new("pe-service-ID"),
					RegionName:                   new("GERMANY_NORTH"),
					Status:                       new("AVAILABLE"),
					PrivateEndpoints:             &[]string{"long-azure-resource-ID"},
					PrivateLinkServiceName:       new("pls_name"),
					PrivateLinkServiceResourceId: new("long-azure-resource-ID"),
				}, &http.Response{}, nil
			},
			mockInterfaceReturnFunc: func() (*admin.PrivateLinkEndpoint, *http.Response, error) {
				return &admin.PrivateLinkEndpoint{
					CloudProvider:                 "AZURE",
					PrivateEndpointConnectionName: new("atlas-resource-name"),
					PrivateEndpointIPAddress:      new("10.0.0.4"),
					PrivateEndpointResourceId:     new("long-azure-resource-ID"),
					Status:                        new("AVAILABLE"),
				}, &http.Response{}, nil
			},
			expectedPE: &AzureService{
				CommonEndpointService: CommonEndpointService{
					ID:            "pe-service-ID",
					CloudRegion:   "GERMANY_NORTH",
					ServiceStatus: "AVAILABLE",
					Interfaces: EndpointInterfaces{
						&AzureInterface{
							CommonEndpointInterface: CommonEndpointInterface{
								ID:              "long-azure-resource-ID",
								InterfaceStatus: "AVAILABLE",
							},
							IP:             "10.0.0.4",
							ConnectionName: "atlas-resource-name",
						},
					},
				},
				ServiceName: "pls_name",
				ResourceID:  "long-azure-resource-ID",
			},
		},
		"get GCP private endpoint": {
			provider: "GCP",
			mockGetReturnFunc: func() (*admin.EndpointService, *http.Response, error) {
				return &admin.EndpointService{
					CloudProvider:          "GCP",
					Id:                     new("pe-service-ID"),
					RegionName:             new("EUROPE_WEST_3"),
					Status:                 new("AVAILABLE"),
					EndpointGroupNames:     &[]string{"group-name"},
					ServiceAttachmentNames: &[]string{"service/attachment1", "service/attachment2", "service/attachment3"},
				}, &http.Response{}, nil
			},
			mockInterfaceReturnFunc: func() (*admin.PrivateLinkEndpoint, *http.Response, error) {
				return &admin.PrivateLinkEndpoint{
					CloudProvider:     "GCP",
					Status:            new("AVAILABLE"),
					EndpointGroupName: new("group-name"),
					Endpoints: &[]admin.GCPConsumerForwardingRule{
						{
							EndpointName: new("group-name-pe-1"),
							IpAddress:    new("10.0.0.1"),
							Status:       new("AVAILABLE"),
						},
						{
							EndpointName: new("group-name-pe-2"),
							IpAddress:    new("10.0.0.3"),
							Status:       new("AVAILABLE"),
						},
						{
							EndpointName: new("group-name-pe-3"),
							IpAddress:    new("10.0.0.3"),
							Status:       new("AVAILABLE"),
						},
					},
				}, &http.Response{}, nil
			},
			expectedPE: &GCPService{
				CommonEndpointService: CommonEndpointService{
					ID:            "pe-service-ID",
					CloudRegion:   "EUROPE_WEST_3",
					ServiceStatus: "AVAILABLE",
					Interfaces: EndpointInterfaces{
						&GCPInterface{
							CommonEndpointInterface: CommonEndpointInterface{
								ID:              "group-name",
								InterfaceStatus: "AVAILABLE",
							},
							Endpoints: []GCPInterfaceEndpoint{
								{
									Name:   "group-name-pe-1",
									IP:     "10.0.0.1",
									Status: "AVAILABLE",
								},
								{
									Name:   "group-name-pe-2",
									IP:     "10.0.0.3",
									Status: "AVAILABLE",
								},
								{
									Name:   "group-name-pe-3",
									IP:     "10.0.0.3",
									Status: "AVAILABLE",
								},
							},
						},
					},
				},
				AttachmentNames: []string{"service/attachment1", "service/attachment2", "service/attachment3"},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			projectID := "project-ID"
			api := mockadmin.NewPrivateEndpointServicesApi(t)
			api.EXPECT().GetPrivateEndpointService(ctx, projectID, tt.provider, "pe-service-ID").
				Return(admin.GetPrivateEndpointServiceApiRequest{ApiService: api})
			api.EXPECT().GetPrivateEndpointServiceExecute(mock.AnythingOfType("admin.GetPrivateEndpointServiceApiRequest")).
				Return(tt.mockGetReturnFunc())

			if tt.mockInterfaceReturnFunc != nil {
				api.EXPECT().GetPrivateEndpoint(ctx, projectID, tt.provider, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(admin.GetPrivateEndpointApiRequest{ApiService: api})
				api.EXPECT().GetPrivateEndpointExecute(mock.AnythingOfType("admin.GetPrivateEndpointApiRequest")).
					Return(tt.mockInterfaceReturnFunc())
			}

			pe := &PrivateEndpoint{
				api: api,
			}

			result, err := pe.GetPrivateEndpoint(ctx, projectID, tt.provider, "pe-service-ID")
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedPE, result)
		})
	}
}

func TestCreatePrivateEndpointService(t *testing.T) {
	tests := map[string]struct {
		service              EndpointService
		mockCreateReturnFunc func() (*admin.EndpointService, *http.Response, error)
		expectedPE           EndpointService
		expectedErr          error
	}{
		"failed to create the service": {
			service: &AWSService{
				CommonEndpointService: CommonEndpointService{
					CloudRegion: "US_EAST_1",
				},
			},
			mockCreateReturnFunc: func() (*admin.EndpointService, *http.Response, error) {
				return nil, &http.Response{}, errors.New("atlas failed to create")
			},
			expectedErr: fmt.Errorf("failed to create the private endpoint service: %w", errors.New("atlas failed to create")),
		},
		"create private endpoint service": {
			service: &AWSService{
				CommonEndpointService: CommonEndpointService{
					CloudRegion: "US_EAST_1",
				},
			},
			mockCreateReturnFunc: func() (*admin.EndpointService, *http.Response, error) {
				return &admin.EndpointService{
					CloudProvider: "AWS",
					Id:            new("pe-service-ID"),
					RegionName:    new("US_EAST_1"),
					Status:        new("INITIALIZING"),
				}, &http.Response{}, nil
			},
			expectedPE: &AWSService{
				CommonEndpointService: CommonEndpointService{
					ID:            "pe-service-ID",
					CloudRegion:   "US_EAST_1",
					ServiceStatus: "INITIALIZING",
					Interfaces:    EndpointInterfaces{},
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			projectID := "project-ID"
			api := mockadmin.NewPrivateEndpointServicesApi(t)
			api.EXPECT().CreatePrivateEndpointService(ctx, projectID, mock.AnythingOfType("*admin.CloudProviderEndpointServiceRequest")).
				Return(admin.CreatePrivateEndpointServiceApiRequest{ApiService: api})
			api.EXPECT().CreatePrivateEndpointServiceExecute(mock.AnythingOfType("admin.CreatePrivateEndpointServiceApiRequest")).
				Return(tt.mockCreateReturnFunc())

			pe := &PrivateEndpoint{
				api: api,
			}

			result, err := pe.CreatePrivateEndpointService(ctx, projectID, tt.service)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedPE, result)
		})
	}
}

func TestDeletePrivateEndpointService(t *testing.T) {
	tests := map[string]struct {
		mockDeleteReturnFunc func() (*http.Response, error)
		expectedErr          error
	}{
		"failed to delete the service": {
			mockDeleteReturnFunc: func() (*http.Response, error) {
				return &http.Response{}, errors.New("atlas failed to delete")
			},
			expectedErr: fmt.Errorf("failed to delete the private endpoint service: %w", errors.New("atlas failed to delete")),
		},
		"delete private endpoint service": {
			mockDeleteReturnFunc: func() (*http.Response, error) {
				return &http.Response{}, nil
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			projectID := "project-ID"
			api := mockadmin.NewPrivateEndpointServicesApi(t)
			api.EXPECT().DeletePrivateEndpointService(ctx, projectID, "AWS", "pe-service-ID").
				Return(admin.DeletePrivateEndpointServiceApiRequest{ApiService: api})
			api.EXPECT().DeletePrivateEndpointServiceExecute(mock.AnythingOfType("admin.DeletePrivateEndpointServiceApiRequest")).
				Return(tt.mockDeleteReturnFunc())

			pe := &PrivateEndpoint{
				api: api,
			}

			err := pe.DeleteEndpointService(ctx, projectID, "AWS", "pe-service-ID")
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestCreatePrivateEndpointInterface(t *testing.T) {
	tests := map[string]struct {
		provider             string
		gcpProjectID         string
		endpointInterface    EndpointInterface
		mockCreateReturnFunc func() (*admin.PrivateLinkEndpoint, *http.Response, error)
		expectedPE           EndpointInterface
		expectedErr          error
	}{
		"failed to create the endpoint interface": {
			provider: "AWS",
			endpointInterface: &AWSInterface{
				CommonEndpointInterface{
					ID: "vpcpe-123456",
				},
			},
			mockCreateReturnFunc: func() (*admin.PrivateLinkEndpoint, *http.Response, error) {
				return nil, &http.Response{}, errors.New("atlas failed to create")
			},
			expectedErr: fmt.Errorf("failed to create the private endpoint interface: %w", errors.New("atlas failed to create")),
		},
		"create AWS private endpoint": {
			provider: "AWS",
			endpointInterface: &AWSInterface{
				CommonEndpointInterface{
					ID: "vpcpe-123456",
				},
			},
			mockCreateReturnFunc: func() (*admin.PrivateLinkEndpoint, *http.Response, error) {
				return &admin.PrivateLinkEndpoint{
					CloudProvider:       "AWS",
					InterfaceEndpointId: new("vpcpe-123456"),
					ConnectionStatus:    new("INITIALIZING"),
				}, &http.Response{}, nil
			},
			expectedPE: &AWSInterface{
				CommonEndpointInterface{
					ID:              "vpcpe-123456",
					InterfaceStatus: "INITIALIZING",
				},
			},
		},
		"create AZURE private endpoint": {
			provider: "AZURE",
			endpointInterface: &AzureInterface{
				CommonEndpointInterface: CommonEndpointInterface{
					ID:              "long-azure-resource-ID",
					InterfaceStatus: "INITIALIZING",
				},
				IP: "10.0.0.2",
			},
			mockCreateReturnFunc: func() (*admin.PrivateLinkEndpoint, *http.Response, error) {
				return &admin.PrivateLinkEndpoint{
					CloudProvider:                 "AZURE",
					PrivateEndpointResourceId:     new("long-azure-resource-ID"),
					PrivateEndpointIPAddress:      new("10.0.0.2"),
					PrivateEndpointConnectionName: new("atlas-resource-name"),
					Status:                        new("INITIALIZING"),
				}, &http.Response{}, nil
			},
			expectedPE: &AzureInterface{
				CommonEndpointInterface: CommonEndpointInterface{
					ID:              "long-azure-resource-ID",
					InterfaceStatus: "INITIALIZING",
				},
				IP:             "10.0.0.2",
				ConnectionName: "atlas-resource-name",
			},
		},
		"create GCP private endpoint": {
			provider:     "GCP",
			gcpProjectID: "customer-project-ID",
			endpointInterface: &GCPInterface{
				CommonEndpointInterface: CommonEndpointInterface{
					ID: "group-name",
				},
				Endpoints: []GCPInterfaceEndpoint{
					{
						Name: "group-name-pe-1",
						IP:   "10.0.0.1",
					},
					{
						Name: "group-name-pe-2",
						IP:   "10.0.0.2",
					},
					{
						Name: "group-name-pe-3",
						IP:   "10.0.0.3",
					},
				},
			},
			mockCreateReturnFunc: func() (*admin.PrivateLinkEndpoint, *http.Response, error) {
				return &admin.PrivateLinkEndpoint{
					CloudProvider:     "GCP",
					EndpointGroupName: new("group-name"),
					Endpoints: &[]admin.GCPConsumerForwardingRule{
						{
							EndpointName: new("group-name-pe-1"),
							IpAddress:    new("10.0.0.1"),
							Status:       new("INITIALIZING"),
						},
						{
							EndpointName: new("group-name-pe-2"),
							IpAddress:    new("10.0.0.2"),
							Status:       new("INITIALIZING"),
						},
						{
							EndpointName: new("group-name-pe-3"),
							IpAddress:    new("10.0.0.3"),
							Status:       new("INITIALIZING"),
						},
					},
					Status: new("INITIALIZING"),
				}, &http.Response{}, nil
			},
			expectedPE: &GCPInterface{
				CommonEndpointInterface: CommonEndpointInterface{
					ID:              "group-name",
					InterfaceStatus: "INITIALIZING",
				},
				Endpoints: []GCPInterfaceEndpoint{
					{
						Name:   "group-name-pe-1",
						IP:     "10.0.0.1",
						Status: "INITIALIZING",
					},
					{
						Name:   "group-name-pe-2",
						IP:     "10.0.0.2",
						Status: "INITIALIZING",
					},
					{
						Name:   "group-name-pe-3",
						IP:     "10.0.0.3",
						Status: "INITIALIZING",
					},
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			projectID := "project-ID"
			serviceID := "pe-service-ID"
			api := mockadmin.NewPrivateEndpointServicesApi(t)
			api.EXPECT().CreatePrivateEndpoint(ctx, projectID, tt.provider, serviceID, mock.AnythingOfType("*admin.CreateEndpointRequest")).
				Return(admin.CreatePrivateEndpointApiRequest{ApiService: api})
			api.EXPECT().CreatePrivateEndpointExecute(mock.AnythingOfType("admin.CreatePrivateEndpointApiRequest")).
				Return(tt.mockCreateReturnFunc())

			pe := &PrivateEndpoint{
				api: api,
			}

			result, err := pe.CreatePrivateEndpointInterface(ctx, projectID, tt.provider, serviceID, tt.gcpProjectID, tt.endpointInterface)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedPE, result)
		})
	}
}

func TestDeletePrivateEndpointInterface(t *testing.T) {
	tests := map[string]struct {
		mockDeleteReturnFunc func() (*http.Response, error)
		expectedErr          error
	}{
		"failed to delete the interface": {
			mockDeleteReturnFunc: func() (*http.Response, error) {
				return &http.Response{}, errors.New("atlas failed to delete")
			},
			expectedErr: fmt.Errorf("failed to delete the private endpoint interface: %w", errors.New("atlas failed to delete")),
		},
		"delete private endpoint service": {
			mockDeleteReturnFunc: func() (*http.Response, error) {
				return &http.Response{}, nil
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			projectID := "project-ID"
			api := mockadmin.NewPrivateEndpointServicesApi(t)
			api.EXPECT().DeletePrivateEndpoint(ctx, projectID, "AWS", "endpoint-ID", "pe-service-ID").
				Return(admin.DeletePrivateEndpointApiRequest{ApiService: api})
			api.EXPECT().DeletePrivateEndpointExecute(mock.AnythingOfType("admin.DeletePrivateEndpointApiRequest")).
				Return(tt.mockDeleteReturnFunc())

			pe := &PrivateEndpoint{
				api: api,
			}

			err := pe.DeleteEndpointInterface(ctx, projectID, "AWS", "pe-service-ID", "endpoint-ID")
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
