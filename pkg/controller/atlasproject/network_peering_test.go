package atlasproject

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115004/admin"

	"github.com/stretchr/testify/require"

	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func TestCanNetworkPeeringReconcile(t *testing.T) {
	t.Run("should return true when subResourceDeletionProtection is disabled", func(t *testing.T) {
		workflowCtx := &workflow.Context{
			Context: context.Background(),
		}
		result, err := canNetworkPeeringReconcile(workflowCtx, false, &mdbv1.AtlasProject{})
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return error when unable to deserialize last applied configuration", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{wrong}"})
		workflowCtx := &workflow.Context{
			Context: context.Background(),
		}
		result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)
		require.EqualError(t, err, "invalid character 'w' looking for beginning of object key string")
		require.False(t, result)
	})

	t.Run("should return error when unable to fetch container data from Atlas", func(t *testing.T) {
		m := atlasmock.NewNetworkPeeringApiMock(t)
		m.EXPECT().ListPeeringContainers(mock.Anything, mock.Anything).Return(admin.ListPeeringContainersApiRequest{ApiService: m})
		m.EXPECT().ListPeeringContainersExecute(mock.Anything).Return(nil, nil, errors.New("failed to retrieve data"))
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{NetworkPeeringApi: m},
			Context:   context.Background(),
		}
		result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)

		require.EqualError(t, err, "failed to retrieve data")
		require.False(t, result)
	})

	t.Run("should return error when unable to fetch peers data from Atlas", func(t *testing.T) {
		m := atlasmock.NewNetworkPeeringApiMock(t)
		m.EXPECT().ListPeeringContainers(mock.Anything, mock.Anything).Return(admin.ListPeeringContainersApiRequest{ApiService: m})
		m.EXPECT().ListPeeringContainersExecute(mock.Anything).Return(&admin.PaginatedCloudProviderContainer{}, nil, nil)
		m.EXPECT().ListPeeringConnections(mock.Anything, mock.Anything).Return(admin.ListPeeringConnectionsApiRequest{ApiService: m})
		m.EXPECT().ListPeeringConnectionsExecute(mock.Anything).Return(nil, nil, errors.New("failed to retrieve data"))
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{NetworkPeeringApi: m},
			Context:   context.Background(),
		}
		result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)

		require.EqualError(t, err, "failed to retrieve data")
		require.False(t, result)
	})

	t.Run("should return true when there are no container items in Atlas", func(t *testing.T) {
		m := atlasmock.NewNetworkPeeringApiMock(t)
		m.EXPECT().ListPeeringContainers(mock.Anything, mock.Anything).Return(admin.ListPeeringContainersApiRequest{ApiService: m})
		m.EXPECT().ListPeeringContainersExecute(mock.Anything).Return(&admin.PaginatedCloudProviderContainer{}, nil, nil)
		m.EXPECT().ListPeeringConnections(mock.Anything, mock.Anything).Return(admin.ListPeeringConnectionsApiRequest{ApiService: m})
		m.EXPECT().ListPeeringConnectionsExecute(mock.Anything).Return(&admin.PaginatedContainerPeer{}, nil, nil)
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{\"networkPeers\":[{\"providerName\":\"AWS\",\"accepterRegionName\":\"eu-west-1\",\"atlasCidrBlock\":\"192.168.0.0/24\"}]}"})
		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{NetworkPeeringApi: m},
			Context:   context.Background(),
		}
		result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return true when there are no difference between current Atlas and previous applied configuration", func(t *testing.T) {
		t.Run("should return true for AWS configuration", func(t *testing.T) {
			m := atlasmock.NewNetworkPeeringApiMock(t)
			m.EXPECT().ListPeeringContainers(mock.Anything, mock.Anything).Return(admin.ListPeeringContainersApiRequest{ApiService: m})
			m.EXPECT().ListPeeringContainersExecute(mock.Anything).Return(&admin.PaginatedCloudProviderContainer{
				Results: &[]admin.CloudProviderContainer{
					{
						ProviderName:   admin.PtrString("AWS"),
						RegionName:     admin.PtrString("EU_WEST_1"),
						AtlasCidrBlock: admin.PtrString("192.168.0.0/24"),
					},
				},
			}, nil, nil)
			m.EXPECT().ListPeeringConnections(mock.Anything, mock.Anything).Return(admin.ListPeeringConnectionsApiRequest{ApiService: m})
			m.EXPECT().ListPeeringConnectionsExecute(mock.Anything).Return(&admin.PaginatedContainerPeer{
				Results: &[]admin.BaseNetworkPeeringConnectionSettings{
					{
						ProviderName:        admin.PtrString("AWS"),
						AccepterRegionName:  admin.PtrString("eu-west-1"),
						RouteTableCidrBlock: admin.PtrString("10.8.0.0/22"),
						AwsAccountId:        admin.PtrString("123456"),
						VpcId:               admin.PtrString("654321"),
					},
				},
			}, nil, nil)
			akoProject := &mdbv1.AtlasProject{
				Spec: mdbv1.AtlasProjectSpec{
					NetworkPeers: []mdbv1.NetworkPeer{
						{
							ProviderName:        provider.ProviderAWS,
							AccepterRegionName:  "eu-west-1",
							AtlasCIDRBlock:      "192.168.0.0/24",
							RouteTableCIDRBlock: "10.8.0.0/22",
							AWSAccountID:        "123456",
							VpcID:               "654321",
						},
					},
				},
			}
			akoProject.WithAnnotations(
				map[string]string{
					customresource.AnnotationLastAppliedConfiguration: "{\"networkPeers\":[{\"accepterRegionName\":\"eu-west-1\",\"awsAccountId\":\"123456\",\"providerName\":\"AWS\",\"routeTableCidrBlock\":\"10.8.0.0/22\",\"vpcId\":\"654321\",\"atlasCidrBlock\":\"192.168.0.0/24\"}]}",
				},
			)
			workflowCtx := &workflow.Context{
				SdkClient: &admin.APIClient{NetworkPeeringApi: m},
				Context:   context.Background(),
			}
			result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)

			require.NoError(t, err)
			require.True(t, result)
		})

		t.Run("should return true for GCP configuration", func(t *testing.T) {
			m := atlasmock.NewNetworkPeeringApiMock(t)
			m.EXPECT().ListPeeringContainers(mock.Anything, mock.Anything).Return(admin.ListPeeringContainersApiRequest{ApiService: m})
			m.EXPECT().ListPeeringContainersExecute(mock.Anything).Return(&admin.PaginatedCloudProviderContainer{
				Results: &[]admin.CloudProviderContainer{
					{
						ProviderName:   admin.PtrString("GCP"),
						AtlasCidrBlock: admin.PtrString("192.168.0.0/24"),
					},
				},
			}, nil, nil)
			m.EXPECT().ListPeeringConnections(mock.Anything, mock.Anything).Return(admin.ListPeeringConnectionsApiRequest{ApiService: m})
			m.EXPECT().ListPeeringConnectionsExecute(mock.Anything).Return(&admin.PaginatedContainerPeer{
				Results: &[]admin.BaseNetworkPeeringConnectionSettings{
					{
						ProviderName:       admin.PtrString("GCP"),
						AccepterRegionName: admin.PtrString("europe-west-1"),
						GcpProjectId:       admin.PtrString("my-project"),
						NetworkName:        admin.PtrString("my-network"),
					},
				},
			}, nil, nil)
			akoProject := &mdbv1.AtlasProject{
				Spec: mdbv1.AtlasProjectSpec{
					NetworkPeers: []mdbv1.NetworkPeer{
						{
							ProviderName:       provider.ProviderGCP,
							AccepterRegionName: "europe-west-1",
							AtlasCIDRBlock:     "192.168.0.0/24",
							GCPProjectID:       "my-project",
							NetworkName:        "my-network",
						},
					},
				},
			}
			akoProject.WithAnnotations(
				map[string]string{
					customresource.AnnotationLastAppliedConfiguration: "{\"networkPeers\":[{\"accepterRegionName\":\"europe-west-1\",\"providerName\":\"GCP\",\"gcpProjectId\":\"my-project\",\"networkName\":\"my-network\",\"atlasCidrBlock\":\"192.168.0.0/24\"}]}",
				},
			)
			workflowCtx := &workflow.Context{
				SdkClient: &admin.APIClient{NetworkPeeringApi: m},
				Context:   context.Background(),
			}
			result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)

			require.NoError(t, err)
			require.True(t, result)
		})

		t.Run("should return true for Azure configuration", func(t *testing.T) {
			m := atlasmock.NewNetworkPeeringApiMock(t)
			m.EXPECT().ListPeeringContainers(mock.Anything, mock.Anything).Return(admin.ListPeeringContainersApiRequest{ApiService: m})
			m.EXPECT().ListPeeringContainersExecute(mock.Anything).Return(&admin.PaginatedCloudProviderContainer{
				Results: &[]admin.CloudProviderContainer{
					{
						ProviderName:   admin.PtrString("AZURE"),
						Region:         admin.PtrString("GERMANY_CENTRAL"),
						AtlasCidrBlock: admin.PtrString("192.168.0.0/24"),
					},
				},
			}, nil, nil)
			m.EXPECT().ListPeeringConnections(mock.Anything, mock.Anything).Return(admin.ListPeeringConnectionsApiRequest{ApiService: m})
			m.EXPECT().ListPeeringConnectionsExecute(mock.Anything).Return(&admin.PaginatedContainerPeer{
				Results: &[]admin.BaseNetworkPeeringConnectionSettings{
					{
						ProviderName:        admin.PtrString("AZURE"),
						AccepterRegionName:  admin.PtrString("GERMANY_CENTRAL"),
						AzureSubscriptionId: admin.PtrString("123"),
						AzureDirectoryId:    admin.PtrString("456"),
						ResourceGroupName:   admin.PtrString("my-rg"),
						VnetName:            admin.PtrString("my-vnet"),
					},
				},
			}, nil, nil)
			akoProject := &mdbv1.AtlasProject{
				Spec: mdbv1.AtlasProjectSpec{
					NetworkPeers: []mdbv1.NetworkPeer{
						{
							ProviderName:        provider.ProviderAzure,
							AccepterRegionName:  "GERMANY_CENTRAL",
							AtlasCIDRBlock:      "192.168.0.0/24",
							AzureSubscriptionID: "123",
							AzureDirectoryID:    "456",
							ResourceGroupName:   "my-rg",
							VNetName:            "my-vnet",
						},
					},
				},
			}
			akoProject.WithAnnotations(
				map[string]string{
					customresource.AnnotationLastAppliedConfiguration: "{\"networkPeers\":[{\"accepterRegionName\":\"GERMANY_CENTRAL\",\"providerName\":\"AZURE\",\"atlasCidrBlock\":\"192.168.0.0/24\",\"azureSubscriptionId\":\"123\",\"azureDirectoryId\":\"456\",\"resourceGroupName\":\"my-rg\",\"vnetName\":\"my-vnet\"}]}",
				},
			)
			workflowCtx := &workflow.Context{
				SdkClient: &admin.APIClient{NetworkPeeringApi: m},
				Context:   context.Background(),
			}
			result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)

			require.NoError(t, err)
			require.True(t, result)
		})
	})

	t.Run("should return false when unable to reconcile due to containers config mismatch", func(t *testing.T) {
		t.Run("should return false for AWS configuration", func(t *testing.T) {
			m := atlasmock.NewNetworkPeeringApiMock(t)
			m.EXPECT().ListPeeringContainers(mock.Anything, mock.Anything).Return(admin.ListPeeringContainersApiRequest{ApiService: m})
			m.EXPECT().ListPeeringContainersExecute(mock.Anything).Return(&admin.PaginatedCloudProviderContainer{
				Results: &[]admin.CloudProviderContainer{
					{
						ProviderName:   admin.PtrString("AWS"),
						RegionName:     admin.PtrString("EU_WEST_1"),
						AtlasCidrBlock: admin.PtrString("192.168.1.0/24"),
					},
				},
			}, nil, nil)
			akoProject := &mdbv1.AtlasProject{
				Spec: mdbv1.AtlasProjectSpec{
					NetworkPeers: []mdbv1.NetworkPeer{
						{
							ProviderName:        provider.ProviderAWS,
							AccepterRegionName:  "eu-west-1",
							AtlasCIDRBlock:      "192.168.0.0/24",
							RouteTableCIDRBlock: "10.8.0.0/22",
							AWSAccountID:        "123456",
							VpcID:               "654321",
						},
					},
				},
			}
			akoProject.WithAnnotations(
				map[string]string{
					customresource.AnnotationLastAppliedConfiguration: "{\"networkPeers\":[{\"accepterRegionName\":\"eu-west-1\",\"awsAccountId\":\"123456\",\"providerName\":\"AWS\",\"routeTableCidrBlock\":\"10.8.0.0/22\",\"vpcId\":\"654321\",\"atlasCidrBlock\":\"192.168.0.0/24\"}]}",
				},
			)
			workflowCtx := &workflow.Context{
				SdkClient: &admin.APIClient{NetworkPeeringApi: m},
				Context:   context.Background(),
			}
			result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)

			require.NoError(t, err)
			require.False(t, result)
		})

		t.Run("should return false for GCP configuration", func(t *testing.T) {
			m := atlasmock.NewNetworkPeeringApiMock(t)
			m.EXPECT().ListPeeringContainers(mock.Anything, mock.Anything).Return(admin.ListPeeringContainersApiRequest{ApiService: m})
			m.EXPECT().ListPeeringContainersExecute(mock.Anything).Return(&admin.PaginatedCloudProviderContainer{
				Results: &[]admin.CloudProviderContainer{
					{
						ProviderName:   admin.PtrString("GCP"),
						AtlasCidrBlock: admin.PtrString("192.168.1.0/24"),
					},
				},
			}, nil, nil)
			akoProject := &mdbv1.AtlasProject{
				Spec: mdbv1.AtlasProjectSpec{
					NetworkPeers: []mdbv1.NetworkPeer{
						{
							ProviderName:       provider.ProviderGCP,
							AccepterRegionName: "europe-west-1",
							AtlasCIDRBlock:     "192.168.0.0/24",
							GCPProjectID:       "my-project",
							NetworkName:        "my-network",
						},
					},
				},
			}
			akoProject.WithAnnotations(
				map[string]string{
					customresource.AnnotationLastAppliedConfiguration: "{\"networkPeers\":[{\"accepterRegionName\":\"europe-west-1\",\"providerName\":\"GCP\",\"gcpProjectId\":\"my-project\",\"networkName\":\"my-network\",\"atlasCidrBlock\":\"192.168.0.0/24\"}]}",
				},
			)
			workflowCtx := &workflow.Context{
				SdkClient: &admin.APIClient{NetworkPeeringApi: m},
				Context:   context.Background(),
			}
			result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)

			require.NoError(t, err)
			require.False(t, result)
		})

		t.Run("should return false for Azure configuration", func(t *testing.T) {
			m := atlasmock.NewNetworkPeeringApiMock(t)
			m.EXPECT().ListPeeringContainers(mock.Anything, mock.Anything).Return(admin.ListPeeringContainersApiRequest{ApiService: m})
			m.EXPECT().ListPeeringContainersExecute(mock.Anything).Return(&admin.PaginatedCloudProviderContainer{
				Results: &[]admin.CloudProviderContainer{
					{
						ProviderName:   admin.PtrString("AZURE"),
						Region:         admin.PtrString("GERMANY_CENTRAL"),
						AtlasCidrBlock: admin.PtrString("192.168.1.0/24"),
					},
				},
			}, nil, nil)
			akoProject := &mdbv1.AtlasProject{
				Spec: mdbv1.AtlasProjectSpec{
					NetworkPeers: []mdbv1.NetworkPeer{
						{
							ProviderName:        provider.ProviderAzure,
							AccepterRegionName:  "GERMANY_CENTRAL",
							AtlasCIDRBlock:      "192.168.0.0/24",
							AzureSubscriptionID: "123",
							AzureDirectoryID:    "456",
							ResourceGroupName:   "my-rg",
							VNetName:            "my-vnet",
						},
					},
				},
			}
			akoProject.WithAnnotations(
				map[string]string{
					customresource.AnnotationLastAppliedConfiguration: "{\"networkPeers\":[{\"accepterRegionName\":\"GERMANY_CENTRAL\",\"providerName\":\"AZURE\",\"atlasCidrBlock\":\"192.168.0.0/24\",\"azureSubscriptionId\":\"123\",\"azureDirectoryId\":\"456\",\"resourceGroupName\":\"my-rg\",\"vnetName\":\"my-vnet\"}]}",
				},
			)
			workflowCtx := &workflow.Context{
				SdkClient: &admin.APIClient{NetworkPeeringApi: m},
				Context:   context.Background(),
			}
			result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)

			require.NoError(t, err)
			require.False(t, result)
		})
	})

	t.Run("should return false when unable to reconcile due to peering config mismatch", func(t *testing.T) {
		t.Run("should return false for AWS configuration", func(t *testing.T) {
			m := atlasmock.NewNetworkPeeringApiMock(t)
			m.EXPECT().ListPeeringContainers(mock.Anything, mock.Anything).Return(admin.ListPeeringContainersApiRequest{ApiService: m})
			m.EXPECT().ListPeeringContainersExecute(mock.Anything).Return(&admin.PaginatedCloudProviderContainer{
				Results: &[]admin.CloudProviderContainer{
					{
						ProviderName:   admin.PtrString("AWS"),
						Region:         admin.PtrString("EU_WEST_1"),
						AtlasCidrBlock: admin.PtrString("192.168.0.0/24"),
					},
				},
			}, nil, nil)
			akoProject := &mdbv1.AtlasProject{
				Spec: mdbv1.AtlasProjectSpec{
					NetworkPeers: []mdbv1.NetworkPeer{
						{
							ProviderName:        provider.ProviderAWS,
							AccepterRegionName:  "eu-west-1",
							AtlasCIDRBlock:      "192.168.0.0/24",
							RouteTableCIDRBlock: "10.8.0.0/22",
							AWSAccountID:        "123456",
							VpcID:               "654321",
						},
					},
				},
			}
			akoProject.WithAnnotations(
				map[string]string{
					customresource.AnnotationLastAppliedConfiguration: "{\"networkPeers\":[{\"accepterRegionName\":\"eu-west-1\",\"awsAccountId\":\"123456\",\"providerName\":\"AWS\",\"routeTableCidrBlock\":\"10.8.0.0/22\",\"vpcId\":\"654321\",\"atlasCidrBlock\":\"192.168.0.0/24\"}]}",
				},
			)
			workflowCtx := &workflow.Context{
				SdkClient: &admin.APIClient{NetworkPeeringApi: m},
				Context:   context.Background(),
			}
			result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)

			require.NoError(t, err)
			require.False(t, result)
		})

		t.Run("should return false for GCP configuration", func(t *testing.T) {
			m := atlasmock.NewNetworkPeeringApiMock(t)
			m.EXPECT().ListPeeringContainers(mock.Anything, mock.Anything).Return(admin.ListPeeringContainersApiRequest{ApiService: m})
			m.EXPECT().ListPeeringContainersExecute(mock.Anything).Return(&admin.PaginatedCloudProviderContainer{
				Results: &[]admin.CloudProviderContainer{
					{
						ProviderName:   admin.PtrString("GCP"),
						AtlasCidrBlock: admin.PtrString("192.168.0.0/24"),
					},
				},
			}, nil, nil)
			m.EXPECT().ListPeeringConnections(mock.Anything, mock.Anything).Return(admin.ListPeeringConnectionsApiRequest{ApiService: m})
			m.EXPECT().ListPeeringConnectionsExecute(mock.Anything).Return(&admin.PaginatedContainerPeer{
				Results: &[]admin.BaseNetworkPeeringConnectionSettings{
					{
						ProviderName:       admin.PtrString("GCP"),
						AccepterRegionName: admin.PtrString("europe-west-1"),
						GcpProjectId:       admin.PtrString("my-project2"),
						NetworkName:        admin.PtrString("my-network"),
					},
				},
			}, nil, nil)
			akoProject := &mdbv1.AtlasProject{
				Spec: mdbv1.AtlasProjectSpec{
					NetworkPeers: []mdbv1.NetworkPeer{
						{
							ProviderName:       provider.ProviderGCP,
							AccepterRegionName: "europe-west-1",
							GCPProjectID:       "my-project",
							NetworkName:        "my-network",
						},
					},
				},
			}
			akoProject.WithAnnotations(
				map[string]string{
					customresource.AnnotationLastAppliedConfiguration: "{\"networkPeers\":[{\"accepterRegionName\":\"europe-west-1\",\"providerName\":\"GCP\",\"gcpProjectId\":\"my-project\",\"networkName\":\"my-network\",\"atlasCidrBlock\":\"192.168.0.0/24\"}]}",
				},
			)
			workflowCtx := &workflow.Context{
				SdkClient: &admin.APIClient{NetworkPeeringApi: m},
				Context:   context.Background(),
			}
			result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)

			require.NoError(t, err)
			require.False(t, result)
		})

		t.Run("should return false for Azure configuration", func(t *testing.T) {
			m := atlasmock.NewNetworkPeeringApiMock(t)
			m.EXPECT().ListPeeringContainers(mock.Anything, mock.Anything).Return(admin.ListPeeringContainersApiRequest{ApiService: m})
			m.EXPECT().ListPeeringContainersExecute(mock.Anything).Return(&admin.PaginatedCloudProviderContainer{
				Results: &[]admin.CloudProviderContainer{
					{
						ProviderName:   admin.PtrString("AZURE"),
						Region:         admin.PtrString("GERMANY_CENTRAL"),
						AtlasCidrBlock: admin.PtrString("192.168.0.0/24"),
					},
				},
			}, nil, nil)
			m.EXPECT().ListPeeringConnections(mock.Anything, mock.Anything).Return(admin.ListPeeringConnectionsApiRequest{ApiService: m})
			m.EXPECT().ListPeeringConnectionsExecute(mock.Anything).Return(&admin.PaginatedContainerPeer{
				Results: &[]admin.BaseNetworkPeeringConnectionSettings{
					{
						ProviderName:        admin.PtrString("GCP"),
						AccepterRegionName:  admin.PtrString("GERMANY_CENTRAL"),
						AzureSubscriptionId: admin.PtrString("123"),
						AzureDirectoryId:    admin.PtrString("456"),
						ResourceGroupName:   admin.PtrString("my-rg2"),
						VnetName:            admin.PtrString("my-vnet"),
					},
				},
			}, nil, nil)
			akoProject := &mdbv1.AtlasProject{
				Spec: mdbv1.AtlasProjectSpec{
					NetworkPeers: []mdbv1.NetworkPeer{
						{
							ProviderName:        provider.ProviderAzure,
							AccepterRegionName:  "GERMANY_CENTRAL",
							AtlasCIDRBlock:      "192.168.0.0/24",
							AzureSubscriptionID: "123",
							AzureDirectoryID:    "456",
							ResourceGroupName:   "my-rg",
							VNetName:            "my-vnet",
						},
					},
				},
			}
			akoProject.WithAnnotations(
				map[string]string{
					customresource.AnnotationLastAppliedConfiguration: "{\"networkPeers\":[{\"accepterRegionName\":\"GERMANY_CENTRAL\",\"providerName\":\"AZURE\",\"atlasCidrBlock\":\"192.168.0.0/24\",\"azureSubscriptionId\":\"123\",\"azureDirectoryId\":\"456\",\"resourceGroupName\":\"my-rg\",\"vnetName\":\"my-vnet\"}]}",
				},
			)
			workflowCtx := &workflow.Context{
				SdkClient: &admin.APIClient{NetworkPeeringApi: m},
				Context:   context.Background(),
			}
			result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)

			require.NoError(t, err)
			require.False(t, result)
		})
	})
}

func TestEnsureNetworkPeers(t *testing.T) {
	t.Run("should failed to reconcile when unable to decide resource ownership", func(t *testing.T) {
		m := atlasmock.NewNetworkPeeringApiMock(t)
		m.EXPECT().ListPeeringContainers(mock.Anything, mock.Anything).Return(admin.ListPeeringContainersApiRequest{ApiService: m})
		m.EXPECT().ListPeeringContainersExecute(mock.Anything).Return(nil, nil, errors.New("failed to retrieve data"))
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{NetworkPeeringApi: m},
			Context:   context.Background(),
		}
		result := ensureNetworkPeers(workflowCtx, akoProject, true)

		require.Equal(t, workflow.Terminate(workflow.Internal, "unable to resolve ownership for deletion protection: failed to retrieve data"), result)
	})

	t.Run("should failed to reconcile when unable to synchronize with Atlas", func(t *testing.T) {
		m := atlasmock.NewNetworkPeeringApiMock(t)
		m.EXPECT().ListPeeringContainers(mock.Anything, mock.Anything).Return(admin.ListPeeringContainersApiRequest{ApiService: m})
		m.EXPECT().ListPeeringContainersExecute(mock.Anything).Return(&admin.PaginatedCloudProviderContainer{
			Results: &[]admin.CloudProviderContainer{
				{
					ProviderName:   admin.PtrString("AWS"),
					Region:         admin.PtrString("EU_WEST_1"),
					AtlasCidrBlock: admin.PtrString("192.168.0.0/24"),
				},
			},
		}, nil, nil)
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				NetworkPeers: []mdbv1.NetworkPeer{
					{
						ProviderName:        provider.ProviderAWS,
						AccepterRegionName:  "eu-west-1",
						AtlasCIDRBlock:      "192.168.0.0/24",
						RouteTableCIDRBlock: "10.8.0.0/22",
						AWSAccountID:        "123456",
						VpcID:               "654321",
					},
				},
			},
		}
		akoProject.WithAnnotations(
			map[string]string{
				customresource.AnnotationLastAppliedConfiguration: "{\"networkPeers\":[{\"accepterRegionName\":\"eu-west-1\",\"awsAccountId\":\"123456\",\"providerName\":\"AWS\",\"routeTableCidrBlock\":\"10.8.0.0/22\",\"vpcId\":\"654321\",\"atlasCidrBlock\":\"192.168.0.0/24\"}]}",
			},
		)
		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{NetworkPeeringApi: m},
			Context:   context.Background(),
		}
		result := ensureNetworkPeers(workflowCtx, akoProject, true)

		require.Equal(
			t,
			workflow.Terminate(
				workflow.AtlasDeletionProtection,
				"unable to reconcile Network Peering due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
			),
			result,
		)
	})
}
