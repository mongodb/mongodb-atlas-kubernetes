package networkpeering_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkpeering"
)

const (
	testProjectID = "fake-test-project-id"

	testAWSSubscriptionID = "fake-subscription-id"

	testVpcID = "fake-vpc-id"

	testPeerID = "fake-peering-id"

	testContainerID = "fake-container-id"

	testAzureDirectoryID = "fake-azure-directorty-id"

	testAzureSubcriptionID = "fake-azure-subcription-id"

	testAzureResourceGroup = "fake-azure-resource-group"

	testVnet = "fake-vnet"

	testGCPProjectID = "fake-test-project"

	testNetworkName = "fake-test-network"
)

var (
	ErrFakeFailure = errors.New("fake-failure")
)

func TestNetworkPeeringCreate(t *testing.T) {
	for _, tc := range []struct {
		title         string
		cfg           *akov2.AtlasNetworkPeeringConfig
		api           admin.NetworkPeeringApi
		expectedPeer  *networkpeering.NetworkPeer
		expectedError error
	}{
		{
			title: "successful api create for AWS returns success",
			cfg: &akov2.AtlasNetworkPeeringConfig{
				Provider: string(provider.ProviderAWS),
				AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
					AccepterRegionName:  "US_EAST_1",
					AWSAccountID:        testAWSSubscriptionID,
					RouteTableCIDRBlock: "10.0.0.0/18",
					VpcID:               testVpcID,
				},
			},
			api: testCreateNetworkPeeringAPI(
				&admin.BaseNetworkPeeringConnectionSettings{
					ContainerId:         testContainerID,
					Id:                  pointer.MakePtr(testPeerID),
					ProviderName:        pointer.MakePtr(string(provider.ProviderAWS)),
					AccepterRegionName:  pointer.MakePtr("US_EAST_1"),
					AwsAccountId:        pointer.MakePtr(testAWSSubscriptionID),
					RouteTableCidrBlock: pointer.MakePtr("10.0.0.0/18"),
					VpcId:               pointer.MakePtr(testVpcID),
				},
				nil,
			),
			expectedPeer: &networkpeering.NetworkPeer{
				AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
					ID:       testPeerID,
					Provider: string(provider.ProviderAWS),
					AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
						AccepterRegionName:  "US_EAST_1",
						AWSAccountID:        testAWSSubscriptionID,
						RouteTableCIDRBlock: "10.0.0.0/18",
						VpcID:               testVpcID,
					},
				},
				ContainerID: testContainerID,
			},
			expectedError: nil,
		},

		{
			title: "API failure gets passed through",
			cfg: &akov2.AtlasNetworkPeeringConfig{
				Provider: string(provider.ProviderAWS),
				AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
					AccepterRegionName:  "US_EAST_1",
					AWSAccountID:        testAWSSubscriptionID,
					RouteTableCIDRBlock: "10.0.0.0/18",
					VpcID:               testVpcID,
				},
			},
			api: testCreateNetworkPeeringAPI(
				nil,
				ErrFakeFailure,
			),
			expectedPeer:  nil,
			expectedError: ErrFakeFailure,
		},

		{
			title: "failure to parse config returns before calling API",
			cfg: &akov2.AtlasNetworkPeeringConfig{
				Provider: "invalid provider",
				AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
					AccepterRegionName:  "US_EAST_1",
					AWSAccountID:        testAWSSubscriptionID,
					RouteTableCIDRBlock: "10.0.0.0/18",
					VpcID:               testVpcID,
				},
			},
			expectedPeer:  nil,
			expectedError: networkpeering.ErrUnsupportedProvider,
		},

		{
			title: "failure to parse API reply",
			cfg: &akov2.AtlasNetworkPeeringConfig{
				Provider: string(provider.ProviderAWS),
				AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
					AccepterRegionName:  "US_EAST_1",
					AWSAccountID:        testAWSSubscriptionID,
					RouteTableCIDRBlock: "10.0.0.0/18",
					VpcID:               testVpcID,
				},
			},
			api: testCreateNetworkPeeringAPI(
				&admin.BaseNetworkPeeringConnectionSettings{
					ContainerId:         testContainerID,
					Id:                  pointer.MakePtr(testPeerID),
					ProviderName:        pointer.MakePtr("oops also invalid provider"),
					AccepterRegionName:  pointer.MakePtr("US_EAST_1"),
					AwsAccountId:        pointer.MakePtr(testAWSSubscriptionID),
					RouteTableCidrBlock: pointer.MakePtr("10.0.0.0/18"),
					VpcId:               pointer.MakePtr(testVpcID),
				},
				nil,
			),
			expectedPeer:  nil,
			expectedError: networkpeering.ErrUnsupportedProvider,
		},
	} {
		ctx := context.Background()
		t.Run(tc.title, func(t *testing.T) {
			s := networkpeering.NewNetworkPeeringService(tc.api)
			container, err := s.Create(ctx, testProjectID, testContainerID, tc.cfg)
			assert.Equal(t, tc.expectedPeer, container)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestNetworkPeeringGet(t *testing.T) {
	for _, tc := range []struct {
		title         string
		api           admin.NetworkPeeringApi
		expectedPeer  *networkpeering.NetworkPeer
		expectedError error
	}{
		{
			title: "successful api get for Azure returns success",
			api: testGetNetworkPeeringAPI(
				&admin.BaseNetworkPeeringConnectionSettings{
					ContainerId:         testContainerID,
					Id:                  pointer.MakePtr(testPeerID),
					ProviderName:        pointer.MakePtr(string(provider.ProviderAzure)),
					AzureDirectoryId:    pointer.MakePtr(testAzureDirectoryID),
					AzureSubscriptionId: pointer.MakePtr(testAzureSubcriptionID),
					ResourceGroupName:   pointer.MakePtr(testAzureResourceGroup),
					VnetName:            pointer.MakePtr(testVnet),
				},
				nil,
			),
			expectedPeer: &networkpeering.NetworkPeer{
				AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
					ID:       testPeerID,
					Provider: string(provider.ProviderAzure),
					AzureConfiguration: &akov2.AzureNetworkPeeringConfiguration{
						AzureDirectoryID:    testAzureDirectoryID,
						AzureSubscriptionID: testAzureSubcriptionID,
						ResourceGroupName:   testAzureResourceGroup,
						VNetName:            testVnet,
					},
				},
				ContainerID: testContainerID,
			},
			expectedError: nil,
		},

		{
			title: "API not found is detected",
			api: testGetNetworkPeeringAPI(
				nil,
				testAPIError("PEER_NOT_FOUND"),
			),
			expectedPeer:  nil,
			expectedError: networkpeering.ErrNotFound,
		},

		{
			title: "generic API failure passes though",
			api: testGetNetworkPeeringAPI(
				nil,
				ErrFakeFailure,
			),
			expectedPeer:  nil,
			expectedError: ErrFakeFailure,
		},

		{
			title: "failure to parse API reply",
			api: testGetNetworkPeeringAPI(
				&admin.BaseNetworkPeeringConnectionSettings{
					ContainerId:         testContainerID,
					Id:                  pointer.MakePtr(testPeerID),
					ProviderName:        pointer.MakePtr("invalid provider"),
					AzureDirectoryId:    pointer.MakePtr(testAzureDirectoryID),
					AzureSubscriptionId: pointer.MakePtr(testAzureSubcriptionID),
					ResourceGroupName:   pointer.MakePtr(testAzureResourceGroup),
					VnetName:            pointer.MakePtr(testVnet),
				},
				nil,
			),
			expectedPeer:  nil,
			expectedError: networkpeering.ErrUnsupportedProvider,
		},
	} {
		ctx := context.Background()
		t.Run(tc.title, func(t *testing.T) {
			s := networkpeering.NewNetworkPeeringService(tc.api)
			container, err := s.Get(ctx, testProjectID, testPeerID)
			assert.Equal(t, tc.expectedPeer, container)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestNetworkPeeringUpdate(t *testing.T) {
	for _, tc := range []struct {
		title         string
		cfg           *akov2.AtlasNetworkPeeringConfig
		api           admin.NetworkPeeringApi
		expectedPeer  *networkpeering.NetworkPeer
		expectedError error
	}{
		{
			title: "successful api update for GCP returns success",
			cfg: &akov2.AtlasNetworkPeeringConfig{
				Provider: string(provider.ProviderGCP),
				GCPConfiguration: &akov2.GCPNetworkPeeringConfiguration{
					GCPProjectID: testGCPProjectID,
					NetworkName:  testNetworkName,
				},
			},
			api: testUpdateNetworkPeeringAPI(
				&admin.BaseNetworkPeeringConnectionSettings{
					ContainerId:  testContainerID,
					Id:           pointer.MakePtr(testPeerID),
					ProviderName: pointer.MakePtr(string(provider.ProviderGCP)),
					GcpProjectId: pointer.MakePtr(testGCPProjectID),
					NetworkName:  pointer.MakePtr(testNetworkName),
				},
				nil,
			),
			expectedPeer: &networkpeering.NetworkPeer{
				AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
					ID:       testPeerID,
					Provider: string(provider.ProviderGCP),
					GCPConfiguration: &akov2.GCPNetworkPeeringConfiguration{
						GCPProjectID: testGCPProjectID,
						NetworkName:  testNetworkName,
					},
				},
				ContainerID: testContainerID,
			},
			expectedError: nil,
		},

		{
			title: "API failure gets passed through",
			cfg: &akov2.AtlasNetworkPeeringConfig{
				Provider: string(provider.ProviderGCP),
				GCPConfiguration: &akov2.GCPNetworkPeeringConfiguration{
					GCPProjectID: testGCPProjectID,
					NetworkName:  testNetworkName,
				},
			},
			api: testUpdateNetworkPeeringAPI(
				nil,
				ErrFakeFailure,
			),
			expectedPeer:  nil,
			expectedError: ErrFakeFailure,
		},

		{
			title: "failure to parse config returns before calling API",
			cfg: &akov2.AtlasNetworkPeeringConfig{
				Provider: "invalid provider",
				GCPConfiguration: &akov2.GCPNetworkPeeringConfiguration{
					GCPProjectID: testGCPProjectID,
					NetworkName:  testNetworkName,
				},
			},
			expectedPeer:  nil,
			expectedError: networkpeering.ErrUnsupportedProvider,
		},

		{
			title: "failure to parse API reply",
			cfg: &akov2.AtlasNetworkPeeringConfig{
				Provider: string(provider.ProviderGCP),
				GCPConfiguration: &akov2.GCPNetworkPeeringConfiguration{
					GCPProjectID: testGCPProjectID,
					NetworkName:  testNetworkName,
				},
			},
			api: testUpdateNetworkPeeringAPI(
				&admin.BaseNetworkPeeringConnectionSettings{
					ContainerId:  testContainerID,
					Id:           pointer.MakePtr(testPeerID),
					ProviderName: pointer.MakePtr("oops also invalid provider"),
					GcpProjectId: pointer.MakePtr(testGCPProjectID),
					NetworkName:  pointer.MakePtr(testNetworkName),
				},
				nil,
			),
			expectedPeer:  nil,
			expectedError: networkpeering.ErrUnsupportedProvider,
		},
	} {
		ctx := context.Background()
		t.Run(tc.title, func(t *testing.T) {
			s := networkpeering.NewNetworkPeeringService(tc.api)
			container, err := s.Update(ctx, testProjectID, testPeerID, testContainerID, tc.cfg)
			assert.Equal(t, tc.expectedPeer, container)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestNetworkPeeringDelete(t *testing.T) {
	for _, tc := range []struct {
		title         string
		api           admin.NetworkPeeringApi
		expectedError error
	}{
		{
			title:         "successful api delete returns success",
			api:           testDeleteNetworkPeeringAPI(nil),
			expectedError: nil,
		},

		{
			title:         "API not found is detected",
			api:           testDeleteNetworkPeeringAPI(testAPIError("PEER_NOT_FOUND")),
			expectedError: networkpeering.ErrNotFound,
		},

		{
			title:         "API already deleting also gets not found",
			api:           testDeleteNetworkPeeringAPI(testAPIError("PEER_ALREADY_REQUESTED_DELETION")),
			expectedError: networkpeering.ErrNotFound,
		},

		{
			title:         "generic API failure passes though",
			api:           testDeleteNetworkPeeringAPI(ErrFakeFailure),
			expectedError: ErrFakeFailure,
		},
	} {
		ctx := context.Background()
		t.Run(tc.title, func(t *testing.T) {
			s := networkpeering.NewNetworkPeeringService(tc.api)
			err := s.Delete(ctx, testProjectID, testPeerID)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func testCreateNetworkPeeringAPI(apiPeering *admin.BaseNetworkPeeringConnectionSettings, err error) admin.NetworkPeeringApi {
	var apiMock mockadmin.NetworkPeeringApi

	apiMock.EXPECT().CreatePeeringConnection(
		mock.Anything, testProjectID, mock.Anything,
	).Return(admin.CreatePeeringConnectionApiRequest{ApiService: &apiMock})

	apiMock.EXPECT().CreatePeeringConnectionExecute(
		mock.AnythingOfType("admin.CreatePeeringConnectionApiRequest"),
	).Return(apiPeering, nil, err)
	return &apiMock
}

func testGetNetworkPeeringAPI(apiPeering *admin.BaseNetworkPeeringConnectionSettings, err error) admin.NetworkPeeringApi {
	var apiMock mockadmin.NetworkPeeringApi

	apiMock.EXPECT().GetPeeringConnection(
		mock.Anything, testProjectID, testPeerID,
	).Return(admin.GetPeeringConnectionApiRequest{ApiService: &apiMock})

	apiMock.EXPECT().GetPeeringConnectionExecute(
		mock.AnythingOfType("admin.GetPeeringConnectionApiRequest"),
	).Return(apiPeering, nil, err)
	return &apiMock
}

func testUpdateNetworkPeeringAPI(apiPeering *admin.BaseNetworkPeeringConnectionSettings, err error) admin.NetworkPeeringApi {
	var apiMock mockadmin.NetworkPeeringApi

	apiMock.EXPECT().UpdatePeeringConnection(
		mock.Anything, testProjectID, testPeerID, mock.Anything,
	).Return(admin.UpdatePeeringConnectionApiRequest{ApiService: &apiMock})

	apiMock.EXPECT().UpdatePeeringConnectionExecute(
		mock.AnythingOfType("admin.UpdatePeeringConnectionApiRequest"),
	).Return(apiPeering, nil, err)
	return &apiMock
}

func testDeleteNetworkPeeringAPI(err error) admin.NetworkPeeringApi {
	var apiMock mockadmin.NetworkPeeringApi

	apiMock.EXPECT().DeletePeeringConnection(
		mock.Anything, testProjectID, testPeerID,
	).Return(admin.DeletePeeringConnectionApiRequest{ApiService: &apiMock})

	apiMock.EXPECT().DeletePeeringConnectionExecute(
		mock.AnythingOfType("admin.DeletePeeringConnectionApiRequest"),
	).Return(nil, nil, err)
	return &apiMock
}

func testAPIError(code string) error {
	err := &admin.GenericOpenAPIError{}
	err.SetModel(admin.ApiError{
		ErrorCode: pointer.MakePtr(code),
	})
	return err
}
