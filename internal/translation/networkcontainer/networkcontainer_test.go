package networkcontainer_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"

	"github.com/stretchr/testify/assert"

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
		container         *networkcontainer.NetworkContainer
		api               admin.NetworkPeeringApi
		expectedContainer *networkcontainer.NetworkContainer
		expectedError     error
	}{
		{
			title: "successful api create for AWS returns success",
			container: &networkcontainer.NetworkContainer{
				Provider:                    string(provider.ProviderAWS),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api: testCreateNetworkPeeringAPI(
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
				ID:                          testContainerID,
				Provider:                    string(provider.ProviderAWS),
				AtlasNetworkContainerConfig: testContainerConfig(),
				AWSStatus:                   &networkcontainer.AWSContainerStatus{VpcID: testVpcID},
			},
			expectedError: nil,
		},

		{
			title: "successful api create for AWS returns success without VPC ID",
			container: &networkcontainer.NetworkContainer{
				Provider:                    string(provider.ProviderAWS),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api: testCreateNetworkPeeringAPI(
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
				ID:                          testContainerID,
				Provider:                    string(provider.ProviderAWS),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			expectedError: nil,
		},

		{
			title: "successful api create for Azure returns success",
			container: &networkcontainer.NetworkContainer{
				Provider:                    string(provider.ProviderAzure),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api: testCreateNetworkPeeringAPI(&admin.CloudProviderContainer{
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
				ID:                          testContainerID,
				Provider:                    string(provider.ProviderAzure),
				AtlasNetworkContainerConfig: testContainerConfig(),
				AzureStatus: &networkcontainer.AzureContainerStatus{
					AzureSubscriptionID: testAzureSubcriptionID,
					VnetName:            testVnet,
				},
			},
			expectedError: nil,
		},

		{
			title: "successful api create for Azure without status updates returns success",
			container: &networkcontainer.NetworkContainer{
				Provider:                    string(provider.ProviderAzure),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api: testCreateNetworkPeeringAPI(&admin.CloudProviderContainer{
				Id:             pointer.MakePtr(testContainerID),
				ProviderName:   pointer.MakePtr(string(provider.ProviderAzure)),
				Provisioned:    pointer.MakePtr(false),
				AtlasCidrBlock: pointer.MakePtr(testContainerConfig().CIDRBlock),
				Region:         pointer.MakePtr(testContainerConfig().Region),
			},
				nil,
			),
			expectedContainer: &networkcontainer.NetworkContainer{
				ID:                          testContainerID,
				Provider:                    string(provider.ProviderAzure),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			expectedError: nil,
		},

		{
			title: "successful api create for GCP returns success",
			container: &networkcontainer.NetworkContainer{
				Provider:                    string(provider.ProviderGCP),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api: testCreateNetworkPeeringAPI(&admin.CloudProviderContainer{
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
				ID:                          testContainerID,
				Provider:                    string(provider.ProviderGCP),
				AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{CIDRBlock: "1.1.1.1/2"},
				GCPStatus: &networkcontainer.GoogleContainerStatus{
					GCPProjectID: testGCPProjectID,
					NetworkName:  testNetworkName,
				},
			},
			expectedError: nil,
		},

		{
			title: "successful api create for GCP without status returns success",
			container: &networkcontainer.NetworkContainer{
				Provider:                    string(provider.ProviderGCP),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api: testCreateNetworkPeeringAPI(&admin.CloudProviderContainer{
				Id:             pointer.MakePtr(testContainerID),
				ProviderName:   pointer.MakePtr(string(provider.ProviderGCP)),
				Provisioned:    pointer.MakePtr(false),
				AtlasCidrBlock: pointer.MakePtr(testContainerConfig().CIDRBlock),
			},
				nil,
			),
			expectedContainer: &networkcontainer.NetworkContainer{
				ID:                          testContainerID,
				Provider:                    string(provider.ProviderGCP),
				AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{CIDRBlock: "1.1.1.1/2"},
			},
			expectedError: nil,
		},

		{
			title: "failed api create returns failure",
			container: &networkcontainer.NetworkContainer{
				Provider:                    "bad-provider",
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api:               testCreateNetworkPeeringAPI(nil, ErrFakeFailure),
			expectedContainer: nil,
			expectedError:     ErrFakeFailure,
		},
	} {
		ctx := context.Background()
		t.Run(tc.title, func(t *testing.T) {
			s := networkcontainer.NewNetworkContainerService(tc.api)
			container, err := s.Create(ctx, testProjectID, tc.container)
			assert.Equal(t, tc.expectedContainer, container)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestNetworkContainerGet(t *testing.T) {
	for _, tc := range []struct {
		title             string
		container         *networkcontainer.NetworkContainer
		api               admin.NetworkPeeringApi
		expectedContainer *networkcontainer.NetworkContainer
		expectedError     error
	}{
		{
			title: "successful api get returns success",
			container: &networkcontainer.NetworkContainer{
				Provider:                    string(provider.ProviderAWS),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api: testGetNetworkPeeringAPI(
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
				ID:                          testContainerID,
				Provider:                    string(provider.ProviderAWS),
				AtlasNetworkContainerConfig: testContainerConfig(),
				AWSStatus:                   &networkcontainer.AWSContainerStatus{VpcID: testVpcID},
			},
			expectedError: nil,
		},

		{
			title: "not found api get returns wrapped not found error",
			container: &networkcontainer.NetworkContainer{
				Provider:                    string(provider.ProviderAWS),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api:               testGetNetworkPeeringAPI(nil, testAPIError("CLOUD_PROVIDER_CONTAINER_NOT_FOUND")),
			expectedContainer: nil,
			expectedError:     networkcontainer.ErrNotFound,
		},

		{
			title: "other api get failure returns wrapped error",
			container: &networkcontainer.NetworkContainer{
				Provider:                    string(provider.ProviderAWS),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api:               testGetNetworkPeeringAPI(nil, ErrFakeFailure),
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

func TestNetworkContainerUpdate(t *testing.T) {
	for _, tc := range []struct {
		title             string
		container         *networkcontainer.NetworkContainer
		api               admin.NetworkPeeringApi
		expectedContainer *networkcontainer.NetworkContainer
		expectedError     error
	}{
		{
			title: "successful api update returns success",
			container: &networkcontainer.NetworkContainer{
				Provider:                    string(provider.ProviderAWS),
				ID:                          testContainerID,
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api: testUpdateNetworkPeeringAPI(
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
				ID:                          testContainerID,
				Provider:                    string(provider.ProviderAWS),
				AtlasNetworkContainerConfig: testContainerConfig(),
				AWSStatus:                   &networkcontainer.AWSContainerStatus{VpcID: testVpcID},
			},
			expectedError: nil,
		},

		{
			title: "api update failure returns wrapped error",
			container: &networkcontainer.NetworkContainer{
				Provider:                    string(provider.ProviderAWS),
				ID:                          testContainerID,
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api:               testUpdateNetworkPeeringAPI(nil, ErrFakeFailure),
			expectedContainer: nil,
			expectedError:     ErrFakeFailure,
		},
	} {
		ctx := context.Background()
		t.Run(tc.title, func(t *testing.T) {
			s := networkcontainer.NewNetworkContainerService(tc.api)
			container, err := s.Update(ctx, testProjectID, tc.container)
			assert.Equal(t, tc.expectedContainer, container)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestNetworkContainerDelete(t *testing.T) {
	for _, tc := range []struct {
		title         string
		container     *networkcontainer.NetworkContainer
		api           admin.NetworkPeeringApi
		expectedError error
	}{
		{
			title: "successful api delete returns success",
			container: &networkcontainer.NetworkContainer{
				Provider:                    string(provider.ProviderAWS),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api:           testDeleteNetworkPeeringAPI(nil),
			expectedError: nil,
		},

		{
			title: "not found api delete failure returns wrapped not found error",
			container: &networkcontainer.NetworkContainer{
				Provider:                    string(provider.ProviderAWS),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api:           testDeleteNetworkPeeringAPI(testAPIError("CLOUD_PROVIDER_CONTAINER_NOT_FOUND")),
			expectedError: networkcontainer.ErrNotFound,
		},

		{
			title: "container in api delete failure returns wrapped container in use",
			container: &networkcontainer.NetworkContainer{
				Provider:                    string(provider.ProviderAWS),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api:           testDeleteNetworkPeeringAPI(testAPIError("CONTAINERS_IN_USE")),
			expectedError: networkcontainer.ErrContainerInUse,
		},

		{
			title: "other api get failure returns wrapped error",
			container: &networkcontainer.NetworkContainer{
				Provider:                    string(provider.ProviderAWS),
				AtlasNetworkContainerConfig: testContainerConfig(),
			},
			api:           testDeleteNetworkPeeringAPI(ErrFakeFailure),
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

func testCreateNetworkPeeringAPI(apiContainer *admin.CloudProviderContainer, err error) admin.NetworkPeeringApi {
	var apiMock mockadmin.NetworkPeeringApi

	apiMock.EXPECT().CreatePeeringContainer(
		mock.Anything, testProjectID, mock.Anything,
	).Return(admin.CreatePeeringContainerApiRequest{ApiService: &apiMock})

	apiMock.EXPECT().CreatePeeringContainerExecute(
		mock.AnythingOfType("admin.CreatePeeringContainerApiRequest"),
	).Return(apiContainer, nil, err)
	return &apiMock
}

func testGetNetworkPeeringAPI(apiContainer *admin.CloudProviderContainer, err error) admin.NetworkPeeringApi {
	var apiMock mockadmin.NetworkPeeringApi

	apiMock.EXPECT().GetPeeringContainer(
		mock.Anything, testProjectID, mock.Anything,
	).Return(admin.GetPeeringContainerApiRequest{ApiService: &apiMock})

	apiMock.EXPECT().GetPeeringContainerExecute(
		mock.AnythingOfType("admin.GetPeeringContainerApiRequest"),
	).Return(apiContainer, nil, err)
	return &apiMock
}

func testAPIError(code string) error {
	err := &admin.GenericOpenAPIError{}
	err.SetModel(admin.ApiError{
		ErrorCode: pointer.MakePtr(code),
	})
	return err
}

func testUpdateNetworkPeeringAPI(apiContainer *admin.CloudProviderContainer, err error) admin.NetworkPeeringApi {
	var apiMock mockadmin.NetworkPeeringApi

	apiMock.EXPECT().UpdatePeeringContainer(
		mock.Anything, testProjectID, testContainerID, mock.Anything,
	).Return(admin.UpdatePeeringContainerApiRequest{ApiService: &apiMock})

	apiMock.EXPECT().UpdatePeeringContainerExecute(
		mock.AnythingOfType("admin.UpdatePeeringContainerApiRequest"),
	).Return(apiContainer, nil, err)
	return &apiMock
}

func testDeleteNetworkPeeringAPI(err error) admin.NetworkPeeringApi {
	var apiMock mockadmin.NetworkPeeringApi

	apiMock.EXPECT().DeletePeeringContainer(
		mock.Anything, testProjectID, testContainerID,
	).Return(admin.DeletePeeringContainerApiRequest{ApiService: &apiMock})

	apiMock.EXPECT().DeletePeeringContainerExecute(
		mock.AnythingOfType("admin.DeletePeeringContainerApiRequest"),
	).Return(nil, nil, err)
	return &apiMock
}
