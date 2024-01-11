package atlasproject

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func TestCanNetworkPeeringReconcile(t *testing.T) {
	t.Run("should return true when subResourceDeletionProtection is disabled", func(t *testing.T) {
		workflowCtx := &workflow.Context{
			Client:  &mongodbatlas.Client{},
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
			Client:  &mongodbatlas.Client{},
			Context: context.Background(),
		}
		result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)
		require.EqualError(t, err, "invalid character 'w' looking for beginning of object key string")
		require.False(t, result)
	})

	t.Run("should return error when unable to fetch container data from Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Containers: &atlas.ContainerClientMock{
				ListFunc: func(projectID string) ([]mongodbatlas.Container, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)

		require.EqualError(t, err, "failed to retrieve data")
		require.False(t, result)
	})

	t.Run("should return error when unable to fetch peers data from Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Containers: &atlas.ContainerClientMock{
				ListFunc: func(projectID string) ([]mongodbatlas.Container, *mongodbatlas.Response, error) {
					return []mongodbatlas.Container{}, nil, nil
				},
			},
			Peers: &atlas.NetworkPeeringClientMock{
				ListFunc: func(projectID string) ([]mongodbatlas.Peer, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)

		require.EqualError(t, err, "failed to retrieve data")
		require.False(t, result)
	})

	t.Run("should return true when there are no container items in Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Containers: &atlas.ContainerClientMock{
				ListFunc: func(projectID string) ([]mongodbatlas.Container, *mongodbatlas.Response, error) {
					return []mongodbatlas.Container{}, nil, nil
				},
			},
			Peers: &atlas.NetworkPeeringClientMock{
				ListFunc: func(projectID string) ([]mongodbatlas.Peer, *mongodbatlas.Response, error) {
					return []mongodbatlas.Peer{}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{\"networkPeers\":[{\"providerName\":\"AWS\",\"accepterRegionName\":\"eu-west-1\",\"atlasCidrBlock\":\"192.168.0.0/24\"}]}"})
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return true when there are no difference between current Atlas and previous applied configuration", func(t *testing.T) {
		t.Run("should return true for AWS configuration", func(t *testing.T) {
			atlasClient := mongodbatlas.Client{
				Containers: &atlas.ContainerClientMock{
					ListFunc: func(projectID string) ([]mongodbatlas.Container, *mongodbatlas.Response, error) {
						return []mongodbatlas.Container{
							{
								ProviderName:   "AWS",
								RegionName:     "EU_WEST_1",
								AtlasCIDRBlock: "192.168.0.0/24",
							},
						}, nil, nil
					},
				},
				Peers: &atlas.NetworkPeeringClientMock{
					ListFunc: func(projectID string) ([]mongodbatlas.Peer, *mongodbatlas.Response, error) {
						return []mongodbatlas.Peer{
							{
								ProviderName:        "AWS",
								AccepterRegionName:  "eu-west-1",
								RouteTableCIDRBlock: "10.8.0.0/22",
								AWSAccountID:        "123456",
								VpcID:               "654321",
							},
						}, nil, nil
					},
				},
			}
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
				Client:  &atlasClient,
				Context: context.Background(),
			}
			result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)

			require.NoError(t, err)
			require.True(t, result)
		})

		t.Run("should return true for GCP configuration", func(t *testing.T) {
			atlasClient := mongodbatlas.Client{
				Containers: &atlas.ContainerClientMock{
					ListFunc: func(projectID string) ([]mongodbatlas.Container, *mongodbatlas.Response, error) {
						return []mongodbatlas.Container{
							{
								ProviderName:   "GCP",
								AtlasCIDRBlock: "192.168.0.0/24",
							},
						}, nil, nil
					},
				},
				Peers: &atlas.NetworkPeeringClientMock{
					ListFunc: func(projectID string) ([]mongodbatlas.Peer, *mongodbatlas.Response, error) {
						return []mongodbatlas.Peer{
							{
								ProviderName:       "GCP",
								AccepterRegionName: "europe-west-1",
								GCPProjectID:       "my-project",
								NetworkName:        "my-network",
							},
						}, nil, nil
					},
				},
			}
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
				Client:  &atlasClient,
				Context: context.Background(),
			}
			result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)

			require.NoError(t, err)
			require.True(t, result)
		})

		t.Run("should return true for Azure configuration", func(t *testing.T) {
			atlasClient := mongodbatlas.Client{
				Containers: &atlas.ContainerClientMock{
					ListFunc: func(projectID string) ([]mongodbatlas.Container, *mongodbatlas.Response, error) {
						return []mongodbatlas.Container{
							{
								ProviderName:   "AZURE",
								Region:         "GERMANY_CENTRAL",
								AtlasCIDRBlock: "192.168.0.0/24",
							},
						}, nil, nil
					},
				},
				Peers: &atlas.NetworkPeeringClientMock{
					ListFunc: func(projectID string) ([]mongodbatlas.Peer, *mongodbatlas.Response, error) {
						return []mongodbatlas.Peer{
							{
								ProviderName:        "AZURE",
								AccepterRegionName:  "GERMANY_CENTRAL",
								AzureSubscriptionID: "123",
								AzureDirectoryID:    "456",
								ResourceGroupName:   "my-rg",
								VNetName:            "my-vnet",
							},
						}, nil, nil
					},
				},
			}
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
				Client:  &atlasClient,
				Context: context.Background(),
			}
			result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)

			require.NoError(t, err)
			require.True(t, result)
		})
	})

	t.Run("should return false when unable to reconcile due to containers config mismatch", func(t *testing.T) {
		t.Run("should return false for AWS configuration", func(t *testing.T) {
			atlasClient := mongodbatlas.Client{
				Containers: &atlas.ContainerClientMock{
					ListFunc: func(projectID string) ([]mongodbatlas.Container, *mongodbatlas.Response, error) {
						return []mongodbatlas.Container{
							{
								ProviderName:   "AWS",
								RegionName:     "EU_WEST_1",
								AtlasCIDRBlock: "192.168.1.0/24",
							},
						}, nil, nil
					},
				},
				Peers: &atlas.NetworkPeeringClientMock{
					ListFunc: func(projectID string) ([]mongodbatlas.Peer, *mongodbatlas.Response, error) {
						return []mongodbatlas.Peer{
							{
								ProviderName:        "AWS",
								AccepterRegionName:  "eu-west-1",
								RouteTableCIDRBlock: "10.8.0.0/22",
								AWSAccountID:        "123456",
								VpcID:               "654321",
							},
						}, nil, nil
					},
				},
			}
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
				Client:  &atlasClient,
				Context: context.Background(),
			}
			result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)

			require.NoError(t, err)
			require.False(t, result)
		})

		t.Run("should return false for GCP configuration", func(t *testing.T) {
			atlasClient := mongodbatlas.Client{
				Containers: &atlas.ContainerClientMock{
					ListFunc: func(projectID string) ([]mongodbatlas.Container, *mongodbatlas.Response, error) {
						return []mongodbatlas.Container{
							{
								ProviderName:   "GCP",
								AtlasCIDRBlock: "192.168.1.0/24",
							},
						}, nil, nil
					},
				},
				Peers: &atlas.NetworkPeeringClientMock{
					ListFunc: func(projectID string) ([]mongodbatlas.Peer, *mongodbatlas.Response, error) {
						return []mongodbatlas.Peer{
							{
								ProviderName:       "GCP",
								AccepterRegionName: "europe-west-1",
								GCPProjectID:       "my-project",
								NetworkName:        "my-network",
							},
						}, nil, nil
					},
				},
			}
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
				Client:  &atlasClient,
				Context: context.Background(),
			}
			result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)

			require.NoError(t, err)
			require.False(t, result)
		})

		t.Run("should return false for Azure configuration", func(t *testing.T) {
			atlasClient := mongodbatlas.Client{
				Containers: &atlas.ContainerClientMock{
					ListFunc: func(projectID string) ([]mongodbatlas.Container, *mongodbatlas.Response, error) {
						return []mongodbatlas.Container{
							{
								ProviderName:   "AZURE",
								Region:         "GERMANY_CENTRAL",
								AtlasCIDRBlock: "192.168.1.0/24",
							},
						}, nil, nil
					},
				},
				Peers: &atlas.NetworkPeeringClientMock{
					ListFunc: func(projectID string) ([]mongodbatlas.Peer, *mongodbatlas.Response, error) {
						return []mongodbatlas.Peer{
							{
								ProviderName:        "AZURE",
								AccepterRegionName:  "GERMANY_CENTRAL",
								AzureSubscriptionID: "123",
								AzureDirectoryID:    "456",
								ResourceGroupName:   "my-rg",
								VNetName:            "my-vnet",
							},
						}, nil, nil
					},
				},
			}
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
				Client:  &atlasClient,
				Context: context.Background(),
			}
			result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)

			require.NoError(t, err)
			require.False(t, result)
		})
	})

	t.Run("should return false when unable to reconcile due to peering config mismatch", func(t *testing.T) {
		t.Run("should return false for AWS configuration", func(t *testing.T) {
			atlasClient := mongodbatlas.Client{
				Containers: &atlas.ContainerClientMock{
					ListFunc: func(projectID string) ([]mongodbatlas.Container, *mongodbatlas.Response, error) {
						return []mongodbatlas.Container{
							{
								ProviderName:   "AWS",
								RegionName:     "EU_WEST_1",
								AtlasCIDRBlock: "192.168.0.0/24",
							},
						}, nil, nil
					},
				},
				Peers: &atlas.NetworkPeeringClientMock{
					ListFunc: func(projectID string) ([]mongodbatlas.Peer, *mongodbatlas.Response, error) {
						return []mongodbatlas.Peer{
							{
								ProviderName:        "AWS",
								AccepterRegionName:  "eu-west-1",
								RouteTableCIDRBlock: "10.9.0.0/22",
								AWSAccountID:        "123456",
								VpcID:               "654321",
							},
						}, nil, nil
					},
				},
			}
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
				Client:  &atlasClient,
				Context: context.Background(),
			}
			result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)

			require.NoError(t, err)
			require.False(t, result)
		})

		t.Run("should return false for GCP configuration", func(t *testing.T) {
			atlasClient := mongodbatlas.Client{
				Containers: &atlas.ContainerClientMock{
					ListFunc: func(projectID string) ([]mongodbatlas.Container, *mongodbatlas.Response, error) {
						return []mongodbatlas.Container{
							{
								ProviderName:   "GCP",
								AtlasCIDRBlock: "192.168.0.0/24",
							},
						}, nil, nil
					},
				},
				Peers: &atlas.NetworkPeeringClientMock{
					ListFunc: func(projectID string) ([]mongodbatlas.Peer, *mongodbatlas.Response, error) {
						return []mongodbatlas.Peer{
							{
								ProviderName:       "GCP",
								AccepterRegionName: "europe-west-1",
								GCPProjectID:       "my-project2",
								NetworkName:        "my-network",
							},
						}, nil, nil
					},
				},
			}
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
				Client:  &atlasClient,
				Context: context.Background(),
			}
			result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)

			require.NoError(t, err)
			require.False(t, result)
		})

		t.Run("should return false for Azure configuration", func(t *testing.T) {
			atlasClient := mongodbatlas.Client{
				Containers: &atlas.ContainerClientMock{
					ListFunc: func(projectID string) ([]mongodbatlas.Container, *mongodbatlas.Response, error) {
						return []mongodbatlas.Container{
							{
								ProviderName:   "AZURE",
								Region:         "GERMANY_CENTRAL",
								AtlasCIDRBlock: "192.168.0.0/24",
							},
						}, nil, nil
					},
				},
				Peers: &atlas.NetworkPeeringClientMock{
					ListFunc: func(projectID string) ([]mongodbatlas.Peer, *mongodbatlas.Response, error) {
						return []mongodbatlas.Peer{
							{
								ProviderName:        "AZURE",
								AccepterRegionName:  "GERMANY_CENTRAL",
								AzureSubscriptionID: "123",
								AzureDirectoryID:    "456",
								ResourceGroupName:   "my-rg2",
								VNetName:            "my-vnet",
							},
						}, nil, nil
					},
				},
			}
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
				Client:  &atlasClient,
				Context: context.Background(),
			}
			result, err := canNetworkPeeringReconcile(workflowCtx, true, akoProject)

			require.NoError(t, err)
			require.False(t, result)
		})
	})
}

func TestEnsureNetworkPeers(t *testing.T) {
	t.Run("should failed to reconcile when unable to decide resource ownership", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Containers: &atlas.ContainerClientMock{
				ListFunc: func(projectID string) ([]mongodbatlas.Container, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		result := ensureNetworkPeers(workflowCtx, akoProject, true)

		require.Equal(t, workflow.Terminate(workflow.Internal, "unable to resolve ownership for deletion protection: failed to retrieve data"), result)
	})

	t.Run("should failed to reconcile when unable to synchronize with Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Containers: &atlas.ContainerClientMock{
				ListFunc: func(projectID string) ([]mongodbatlas.Container, *mongodbatlas.Response, error) {
					return []mongodbatlas.Container{
						{
							ProviderName:   "AWS",
							RegionName:     "EU_WEST_1",
							AtlasCIDRBlock: "192.168.0.0/24",
						},
					}, nil, nil
				},
			},
			Peers: &atlas.NetworkPeeringClientMock{
				ListFunc: func(projectID string) ([]mongodbatlas.Peer, *mongodbatlas.Response, error) {
					return []mongodbatlas.Peer{
						{
							ProviderName:        "AWS",
							AccepterRegionName:  "eu-west-1",
							RouteTableCIDRBlock: "10.9.0.0/22",
							AWSAccountID:        "123456",
							VpcID:               "654321",
						},
					}, nil, nil
				},
			},
		}
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
			Client:  &atlasClient,
			Context: context.Background(),
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
