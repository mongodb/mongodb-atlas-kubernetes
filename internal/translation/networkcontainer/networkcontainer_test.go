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

package networkcontainer_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20250312010/admin"
	"go.mongodb.org/atlas-sdk/v20250312010/mockadmin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkcontainer"
)

const (
	testProjectID = "fake-test-project-id"

	testContainerID = "fake-container-id"

	testVpcID = "fake-vpc-id"

	testAzureSubcriptionID = "fake-azure-subcription-id"

	testVnet = "fake-vnet"

	testGCPProjectID = "fake-test-project"

	testNetworkName = "fake-test-network"
)

var (
	ErrFakeFailure = errors.New("fake-failure")
)

func TestNetworkContainerCreate(t *testing.T) {
	for _, tc := range []struct {
		title             string
		cfg               *networkcontainer.NetworkContainerConfig
		api               admin.NetworkPeeringApi
		expectedContainer *networkcontainer.NetworkContainer
		expectedError     error
	}{
		{
			title: "successful api create for AWS returns success",
			cfg: &networkcontainer.NetworkContainerConfig{
				Provider:                    string(provider.ProviderAWS),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api: testCreateNetworkContainerAPI(
				&admin.CloudProviderContainer{
					Id:             pointer.MakePtr(testContainerID),
					ProviderName:   pointer.MakePtr(string(provider.ProviderAWS)),
					Provisioned:    pointer.MakePtr(false),
					AtlasCidrBlock: pointer.MakePtr(testContainerConfig().CIDRBlock),
					RegionName:     pointer.MakePtr(testContainerConfig().Region),
					VpcId:          pointer.MakePtr(testVpcID),
				},
				nil,
			),
			expectedContainer: &networkcontainer.NetworkContainer{
				NetworkContainerConfig: networkcontainer.NetworkContainerConfig{
					Provider:                    string(provider.ProviderAWS),
					AtlasNetworkContainerConfig: testContainerConfig(),
				},
				ID:        testContainerID,
				AWSStatus: &networkcontainer.AWSContainerStatus{VpcID: testVpcID},
			},
			expectedError: nil,
		},

		{
			title: "successful api create for AWS returns success without VPC ID",
			cfg: &networkcontainer.NetworkContainerConfig{
				Provider:                    string(provider.ProviderAWS),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api: testCreateNetworkContainerAPI(
				&admin.CloudProviderContainer{
					Id:             pointer.MakePtr(testContainerID),
					ProviderName:   pointer.MakePtr(string(provider.ProviderAWS)),
					Provisioned:    pointer.MakePtr(false),
					AtlasCidrBlock: pointer.MakePtr(testContainerConfig().CIDRBlock),
					RegionName:     pointer.MakePtr(testContainerConfig().Region),
				},
				nil,
			),
			expectedContainer: &networkcontainer.NetworkContainer{
				NetworkContainerConfig: networkcontainer.NetworkContainerConfig{
					Provider:                    string(provider.ProviderAWS),
					AtlasNetworkContainerConfig: testContainerConfig(),
				},
				ID: testContainerID,
			},
			expectedError: nil,
		},

		{
			title: "successful api create for Azure returns success",
			cfg: &networkcontainer.NetworkContainerConfig{
				Provider:                    string(provider.ProviderAzure),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api: testCreateNetworkContainerAPI(&admin.CloudProviderContainer{
				Id:                  pointer.MakePtr(testContainerID),
				ProviderName:        pointer.MakePtr(string(provider.ProviderAzure)),
				Provisioned:         pointer.MakePtr(false),
				AtlasCidrBlock:      pointer.MakePtr(testContainerConfig().CIDRBlock),
				Region:              pointer.MakePtr(testContainerConfig().Region),
				AzureSubscriptionId: pointer.MakePtr(string(testAzureSubcriptionID)),
				VnetName:            pointer.MakePtr(testVnet),
			},
				nil,
			),
			expectedContainer: &networkcontainer.NetworkContainer{
				NetworkContainerConfig: networkcontainer.NetworkContainerConfig{
					Provider:                    string(provider.ProviderAzure),
					AtlasNetworkContainerConfig: testContainerConfig(),
				},
				ID: testContainerID,
				AzureStatus: &networkcontainer.AzureContainerStatus{
					AzureSubscriptionID: testAzureSubcriptionID,
					VnetName:            testVnet,
				},
			},
			expectedError: nil,
		},

		{
			title: "successful api create for Azure without status updates returns success",
			cfg: &networkcontainer.NetworkContainerConfig{
				Provider:                    string(provider.ProviderAzure),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api: testCreateNetworkContainerAPI(&admin.CloudProviderContainer{
				Id:             pointer.MakePtr(testContainerID),
				ProviderName:   pointer.MakePtr(string(provider.ProviderAzure)),
				Provisioned:    pointer.MakePtr(false),
				AtlasCidrBlock: pointer.MakePtr(testContainerConfig().CIDRBlock),
				Region:         pointer.MakePtr(testContainerConfig().Region),
			},
				nil,
			),
			expectedContainer: &networkcontainer.NetworkContainer{
				NetworkContainerConfig: networkcontainer.NetworkContainerConfig{
					Provider:                    string(provider.ProviderAzure),
					AtlasNetworkContainerConfig: testContainerConfig(),
				},
				ID: testContainerID,
			},
			expectedError: nil,
		},

		{
			title: "successful api create for GCP returns success",
			cfg: &networkcontainer.NetworkContainerConfig{
				Provider:                    string(provider.ProviderGCP),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api: testCreateNetworkContainerAPI(&admin.CloudProviderContainer{
				Id:             pointer.MakePtr(testContainerID),
				ProviderName:   pointer.MakePtr(string(provider.ProviderGCP)),
				Provisioned:    pointer.MakePtr(false),
				AtlasCidrBlock: pointer.MakePtr(testContainerConfig().CIDRBlock),
				GcpProjectId:   pointer.MakePtr(testGCPProjectID),
				NetworkName:    pointer.MakePtr(testNetworkName),
			},
				nil,
			),
			expectedContainer: &networkcontainer.NetworkContainer{
				NetworkContainerConfig: networkcontainer.NetworkContainerConfig{
					Provider:                    string(provider.ProviderGCP),
					AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{CIDRBlock: "1.1.1.1/2"},
				},
				ID: testContainerID,
				GCPStatus: &networkcontainer.GCPContainerStatus{
					GCPProjectID: testGCPProjectID,
					NetworkName:  testNetworkName,
				},
			},
			expectedError: nil,
		},

		{
			title: "successful api create for GCP without status returns success",
			cfg: &networkcontainer.NetworkContainerConfig{
				Provider:                    string(provider.ProviderGCP),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api: testCreateNetworkContainerAPI(&admin.CloudProviderContainer{
				Id:             pointer.MakePtr(testContainerID),
				ProviderName:   pointer.MakePtr(string(provider.ProviderGCP)),
				Provisioned:    pointer.MakePtr(false),
				AtlasCidrBlock: pointer.MakePtr(testContainerConfig().CIDRBlock),
			},
				nil,
			),
			expectedContainer: &networkcontainer.NetworkContainer{
				NetworkContainerConfig: networkcontainer.NetworkContainerConfig{
					Provider:                    string(provider.ProviderGCP),
					AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{CIDRBlock: "1.1.1.1/2"},
				},
				ID: testContainerID,
			},
			expectedError: nil,
		},

		{
			title: "failed api create returns failure",
			cfg: &networkcontainer.NetworkContainerConfig{
				Provider:                    "bad-provider",
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api:               testCreateNetworkContainerAPI(nil, ErrFakeFailure),
			expectedContainer: nil,
			expectedError:     ErrFakeFailure,
		},
	} {
		ctx := context.Background()
		t.Run(tc.title, func(t *testing.T) {
			s := networkcontainer.NewNetworkContainerService(tc.api)
			container, err := s.Create(ctx, testProjectID, tc.cfg)
			assert.Equal(t, tc.expectedContainer, container)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestNetworkContainerGet(t *testing.T) {
	for _, tc := range []struct {
		title             string
		api               admin.NetworkPeeringApi
		expectedContainer *networkcontainer.NetworkContainer
		expectedError     error
	}{
		{
			title: "successful api get returns success",
			api: testGetNetworkContainerAPI(
				&admin.CloudProviderContainer{
					Id:             pointer.MakePtr(testContainerID),
					ProviderName:   pointer.MakePtr(string(provider.ProviderAWS)),
					Provisioned:    pointer.MakePtr(false),
					AtlasCidrBlock: pointer.MakePtr(testContainerConfig().CIDRBlock),
					RegionName:     pointer.MakePtr(testContainerConfig().Region),
					VpcId:          pointer.MakePtr(testVpcID),
				},
				nil,
			),
			expectedContainer: &networkcontainer.NetworkContainer{
				NetworkContainerConfig: networkcontainer.NetworkContainerConfig{
					Provider:                    string(provider.ProviderAWS),
					AtlasNetworkContainerConfig: testContainerConfig(),
				},
				ID:        testContainerID,
				AWSStatus: &networkcontainer.AWSContainerStatus{VpcID: testVpcID},
			},
			expectedError: nil,
		},

		{
			title:             "not found api get returns wrapped not found error",
			api:               testGetNetworkContainerAPI(nil, testAPIError("CLOUD_PROVIDER_CONTAINER_NOT_FOUND")),
			expectedContainer: nil,
			expectedError:     networkcontainer.ErrNotFound,
		},

		{
			title:             "other api get failure returns wrapped error",
			api:               testGetNetworkContainerAPI(nil, ErrFakeFailure),
			expectedContainer: nil,
			expectedError:     ErrFakeFailure,
		},
	} {
		ctx := context.Background()
		t.Run(tc.title, func(t *testing.T) {
			s := networkcontainer.NewNetworkContainerService(tc.api)
			container, err := s.Get(ctx, testProjectID, testContainerID)
			assert.Equal(t, tc.expectedContainer, container)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestNetworkContainerFind(t *testing.T) {
	for _, tc := range []struct {
		title             string
		cfg               *networkcontainer.NetworkContainerConfig
		api               admin.NetworkPeeringApi
		expectedContainer *networkcontainer.NetworkContainer
		expectedError     error
	}{
		{
			title: "successful find returns success",
			cfg: &networkcontainer.NetworkContainerConfig{
				Provider:                    string(provider.ProviderAWS),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api: testFindNetworkContainerAPI(
				[]admin.CloudProviderContainer{
					{
						Id:             pointer.MakePtr(testContainerID),
						ProviderName:   pointer.MakePtr(string(provider.ProviderAWS)),
						Provisioned:    pointer.MakePtr(false),
						AtlasCidrBlock: pointer.MakePtr(testContainerConfig().CIDRBlock),
						RegionName:     pointer.MakePtr(testContainerConfig().Region),
						VpcId:          pointer.MakePtr(testVpcID),
					},
				},
				nil,
			),
			expectedContainer: &networkcontainer.NetworkContainer{
				NetworkContainerConfig: networkcontainer.NetworkContainerConfig{
					Provider:                    string(provider.ProviderAWS),
					AtlasNetworkContainerConfig: testContainerConfig(),
				},
				ID:        testContainerID,
				AWSStatus: &networkcontainer.AWSContainerStatus{VpcID: testVpcID},
			},
			expectedError: nil,
		},

		{
			title: "find fails other error",
			cfg: &networkcontainer.NetworkContainerConfig{
				Provider:                    string(provider.ProviderAWS),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api: testFindNetworkContainerAPI(
				nil,
				ErrFakeFailure,
			),
			expectedContainer: nil,
			expectedError:     ErrFakeFailure,
		},

		{
			title: "find fails not found",
			cfg: &networkcontainer.NetworkContainerConfig{
				Provider:                    string(provider.ProviderAWS),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api: testFindNetworkContainerAPI(
				[]admin.CloudProviderContainer{},
				nil,
			),
			expectedContainer: nil,
			expectedError:     networkcontainer.ErrNotFound,
		},

		{
			title: "successful find on GCP",
			cfg: &networkcontainer.NetworkContainerConfig{
				Provider: string(provider.ProviderGCP),
				AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
					CIDRBlock: "18.18.192.0/18",
				},
			},
			api: testFindNetworkContainerAPI(
				[]admin.CloudProviderContainer{
					{
						Id:             pointer.MakePtr(testContainerID),
						ProviderName:   pointer.MakePtr(string(provider.ProviderGCP)),
						Provisioned:    pointer.MakePtr(false),
						AtlasCidrBlock: pointer.MakePtr("18.18.192.0/18"),
						GcpProjectId:   pointer.MakePtr(testGCPProjectID),
						NetworkName:    pointer.MakePtr(testNetworkName),
					},
				},
				nil,
			),
			expectedContainer: &networkcontainer.NetworkContainer{
				NetworkContainerConfig: networkcontainer.NetworkContainerConfig{
					Provider: string(provider.ProviderGCP),
					AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
						CIDRBlock: "18.18.192.0/18",
					},
				},
				ID: testContainerID,
				GCPStatus: &networkcontainer.GCPContainerStatus{
					GCPProjectID: testGCPProjectID,
					NetworkName:  testNetworkName,
				},
			},
			expectedError: nil,
		},

		{
			title: "successful find on Azure",
			cfg: &networkcontainer.NetworkContainerConfig{
				Provider: string(provider.ProviderAzure),
				AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
					CIDRBlock: "11.11.0.0/16",
					Region:    "US_EAST_2",
				},
			},
			api: testFindNetworkContainerAPI(
				[]admin.CloudProviderContainer{
					{
						Id:             pointer.MakePtr(testContainerID),
						ProviderName:   pointer.MakePtr(string(provider.ProviderAzure)),
						Provisioned:    pointer.MakePtr(false),
						AtlasCidrBlock: pointer.MakePtr("11.11.0.0/16"),
						Region:         pointer.MakePtr("US_EAST_2"),
					},
				},
				nil,
			),
			expectedContainer: &networkcontainer.NetworkContainer{
				NetworkContainerConfig: networkcontainer.NetworkContainerConfig{
					Provider: string(provider.ProviderAzure),
					AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
						CIDRBlock: "11.11.0.0/16",
						Region:    "US_EAST_2",
					},
				},
				ID: testContainerID,
			},
			expectedError: nil,
		},

		{
			title: "not found on Azure",
			cfg: &networkcontainer.NetworkContainerConfig{
				Provider: string(provider.ProviderAzure),
				AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
					CIDRBlock: "11.11.0.0/16",
					Region:    "US_CENTRAL_5",
				},
			},
			api: testFindNetworkContainerAPI(
				[]admin.CloudProviderContainer{
					{
						Id:             pointer.MakePtr(testContainerID),
						ProviderName:   pointer.MakePtr(string(provider.ProviderAzure)),
						Provisioned:    pointer.MakePtr(false),
						AtlasCidrBlock: pointer.MakePtr("11.11.0.0/16"),
						Region:         pointer.MakePtr("US_EAST_2"),
					},
				},
				nil,
			),
			expectedContainer: nil,
			expectedError:     networkcontainer.ErrNotFound,
		},
	} {
		ctx := context.Background()
		t.Run(tc.title, func(t *testing.T) {
			s := networkcontainer.NewNetworkContainerService(tc.api)
			container, err := s.Find(ctx, testProjectID, tc.cfg)
			assert.Equal(t, tc.expectedContainer, container)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestNetworkContainerUpdate(t *testing.T) {
	for _, tc := range []struct {
		title             string
		cfg               *networkcontainer.NetworkContainerConfig
		api               admin.NetworkPeeringApi
		expectedContainer *networkcontainer.NetworkContainer
		expectedError     error
	}{
		{
			title: "successful api update returns success",
			cfg: &networkcontainer.NetworkContainerConfig{
				Provider:                    string(provider.ProviderAWS),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api: testUpdateNetworkContainerAPI(
				&admin.CloudProviderContainer{
					Id:             pointer.MakePtr(testContainerID),
					ProviderName:   pointer.MakePtr(string(provider.ProviderAWS)),
					Provisioned:    pointer.MakePtr(false),
					AtlasCidrBlock: pointer.MakePtr(testContainerConfig().CIDRBlock),
					RegionName:     pointer.MakePtr(testContainerConfig().Region),
					VpcId:          pointer.MakePtr(testVpcID),
				},
				nil,
			),
			expectedContainer: &networkcontainer.NetworkContainer{
				NetworkContainerConfig: networkcontainer.NetworkContainerConfig{
					Provider:                    string(provider.ProviderAWS),
					AtlasNetworkContainerConfig: testContainerConfig(),
				},
				ID:        testContainerID,
				AWSStatus: &networkcontainer.AWSContainerStatus{VpcID: testVpcID},
			},
			expectedError: nil,
		},

		{
			title: "api update failure returns wrapped error",
			cfg: &networkcontainer.NetworkContainerConfig{
				Provider:                    string(provider.ProviderAWS),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api:               testUpdateNetworkContainerAPI(nil, ErrFakeFailure),
			expectedContainer: nil,
			expectedError:     ErrFakeFailure,
		},
	} {
		ctx := context.Background()
		t.Run(tc.title, func(t *testing.T) {
			s := networkcontainer.NewNetworkContainerService(tc.api)
			container, err := s.Update(ctx, testProjectID, testContainerID, tc.cfg)
			assert.Equal(t, tc.expectedContainer, container)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestNetworkContainerDelete(t *testing.T) {
	for _, tc := range []struct {
		title         string
		api           admin.NetworkPeeringApi
		expectedError error
	}{
		{
			title:         "successful api delete returns success",
			api:           testDeleteNetworkContainerAPI(nil),
			expectedError: nil,
		},

		{
			title:         "not found api delete failure returns wrapped not found error",
			api:           testDeleteNetworkContainerAPI(testAPIError("CLOUD_PROVIDER_CONTAINER_NOT_FOUND")),
			expectedError: networkcontainer.ErrNotFound,
		},

		{
			title:         "container in api delete failure returns wrapped container in use",
			api:           testDeleteNetworkContainerAPI(testAPIError("CONTAINERS_IN_USE")),
			expectedError: networkcontainer.ErrContainerInUse,
		},

		{
			title:         "other api get failure returns wrapped error",
			api:           testDeleteNetworkContainerAPI(ErrFakeFailure),
			expectedError: ErrFakeFailure,
		},
	} {
		ctx := context.Background()
		t.Run(tc.title, func(t *testing.T) {
			s := networkcontainer.NewNetworkContainerService(tc.api)
			err := s.Delete(ctx, testProjectID, testContainerID)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func testContainerConfig() akov2.AtlasNetworkContainerConfig {
	return akov2.AtlasNetworkContainerConfig{
		Region:    "sample-region",
		CIDRBlock: "1.1.1.1/2",
	}
}

func testCreateNetworkContainerAPI(apiContainer *admin.CloudProviderContainer, err error) admin.NetworkPeeringApi {
	var apiMock mockadmin.NetworkPeeringApi

	apiMock.EXPECT().CreatePeeringContainer(
		mock.Anything, testProjectID, mock.Anything,
	).Return(admin.CreatePeeringContainerApiRequest{ApiService: &apiMock})

	apiMock.EXPECT().CreatePeeringContainerExecute(
		mock.AnythingOfType("admin.CreatePeeringContainerApiRequest"),
	).Return(apiContainer, nil, err)
	return &apiMock
}

func testGetNetworkContainerAPI(apiContainer *admin.CloudProviderContainer, err error) admin.NetworkPeeringApi {
	var apiMock mockadmin.NetworkPeeringApi

	apiMock.EXPECT().GetPeeringContainer(
		mock.Anything, testProjectID, mock.Anything,
	).Return(admin.GetPeeringContainerApiRequest{ApiService: &apiMock})

	apiMock.EXPECT().GetPeeringContainerExecute(
		mock.AnythingOfType("admin.GetPeeringContainerApiRequest"),
	).Return(apiContainer, nil, err)
	return &apiMock
}

func testFindNetworkContainerAPI(apiContainers []admin.CloudProviderContainer, err error) admin.NetworkPeeringApi {
	var apiMock mockadmin.NetworkPeeringApi

	apiMock.EXPECT().ListPeeringContainerByCloudProvider(mock.Anything, testProjectID).Return(
		admin.ListPeeringContainerByCloudProviderApiRequest{ApiService: &apiMock},
	)

	results := admin.PaginatedCloudProviderContainer{
		Results: &apiContainers,
	}
	apiMock.EXPECT().ListPeeringContainerByCloudProviderExecute(
		mock.AnythingOfType("admin.ListPeeringContainerByCloudProviderApiRequest"),
	).Return(&results, nil, err)
	return &apiMock
}

func testAPIError(code string) error {
	err := &admin.GenericOpenAPIError{}
	err.SetModel(admin.ApiError{
		ErrorCode: code,
	})
	return err
}

func testUpdateNetworkContainerAPI(apiContainer *admin.CloudProviderContainer, err error) admin.NetworkPeeringApi {
	var apiMock mockadmin.NetworkPeeringApi

	apiMock.EXPECT().UpdatePeeringContainer(
		mock.Anything, testProjectID, testContainerID, mock.Anything,
	).Return(admin.UpdatePeeringContainerApiRequest{ApiService: &apiMock})

	apiMock.EXPECT().UpdatePeeringContainerExecute(
		mock.AnythingOfType("admin.UpdatePeeringContainerApiRequest"),
	).Return(apiContainer, nil, err)
	return &apiMock
}

func testDeleteNetworkContainerAPI(err error) admin.NetworkPeeringApi {
	var apiMock mockadmin.NetworkPeeringApi

	apiMock.EXPECT().DeletePeeringContainer(
		mock.Anything, testProjectID, testContainerID,
	).Return(admin.DeletePeeringContainerApiRequest{ApiService: &apiMock})

	apiMock.EXPECT().DeletePeeringContainerExecute(
		mock.AnythingOfType("admin.DeletePeeringContainerApiRequest"),
	).Return(nil, err)
	return &apiMock
}
