package networkpeering

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115004/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkpeering"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/cloud/aws"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/cloud/azure"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/cloud/google"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/contract"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
)

const (
	testVPCName = "ako-test-network-peering-vpc"
)

func TestPeerContainerServiceCRUD(t *testing.T) {
	ctx := context.Background()
	contract.RunGoContractTest(ctx, t, "test container CRUD", func(ch contract.ContractHelper) {
		projectName := "peer-container-crud-project"
		require.NoError(t, ch.AddResources(ctx, time.Minute, contract.DefaultAtlasProject(projectName)))
		testProjectID, err := ch.ProjectID(ctx, projectName)
		require.NoError(t, err)
		nps := networkpeering.NewNetworkPeeringService(ch.AtlasClient().NetworkPeeringApi)
		cs := nps.(networkpeering.PeeringContainerService)
		for _, tc := range []struct {
			provider  string
			container *networkpeering.ProviderContainer
		}{
			{
				provider:  "AWS",
				container: testAWSPeeringContainer("10.1.0.0/21"),
			},
			{
				provider:  "Azure",
				container: testAzurePeeringContainer("10.2.0.0/21"),
			},
			{
				provider:  "Google",
				container: testGooglePeeringContainer("10.3.0.0/18"), // .../21 is not allowed in GCP
			},
		} {
			createdContainer := &networkpeering.ProviderContainer{}
			t.Run(fmt.Sprintf("create %s container", tc.provider), func(t *testing.T) {
				newContainer, err := cs.CreateContainer(ctx, testProjectID, tc.container)
				require.NoError(t, err)
				assert.NotEmpty(t, newContainer.ID)
				createdContainer = newContainer
			})

			t.Run(fmt.Sprintf("list %s containers", tc.provider), func(t *testing.T) {
				containers, err := cs.ListContainers(ctx, testProjectID, tc.container.Provider)
				require.NoError(t, err)
				assert.NotEmpty(t, containers)
				assert.GreaterOrEqual(t, len(containers), 1)
			})

			t.Run(fmt.Sprintf("get %s container", tc.provider), func(t *testing.T) {
				container, err := cs.GetContainer(ctx, testProjectID, createdContainer.ID)
				require.NoError(t, err)
				assert.NotEmpty(t, container)
				assert.Equal(t, createdContainer.ID, container.ID)
				assert.Equal(t, tc.container.ContainerRegion, container.ContainerRegion)
				assert.Equal(t, tc.container.AtlasCIDRBlock, container.AtlasCIDRBlock)
			})

			t.Run(fmt.Sprintf("delete %s container", tc.provider), func(t *testing.T) {
				time.Sleep(time.Second) // Atlas may reject removal if it happened within a second of creation
				assert.NoErrorf(t, ignoreRemoved(cs.DeleteContainer(ctx, testProjectID, createdContainer.ID)),
					"failed cleanup for provider %s Atlas project ID %s and container id %s",
					tc.provider, testProjectID, createdContainer.ID)
			})
		}
	})
}

func TestPeerServiceCRUD(t *testing.T) {
	ctx := context.Background()
	contract.RunGoContractTest(ctx, t, "test container CRUD", func(ch contract.ContractHelper) {
		projectName := "peer-connection-crud-project"
		require.NoError(t, ch.AddResources(ctx, time.Minute, contract.DefaultAtlasProject(projectName)))
		testProjectID, err := ch.ProjectID(ctx, projectName)
		require.NoError(t, err)
		nps := networkpeering.NewNetworkPeeringService(ch.AtlasClient().NetworkPeeringApi)
		ps := nps.(networkpeering.PeerConnectionsService)
		createdPeer := &networkpeering.NetworkPeer{}

		for _, tc := range []struct {
			provider          string
			preparedCloudTest func(func(peerRequest *networkpeering.NetworkPeer))
		}{
			{
				provider: "AWS",
				preparedCloudTest: func(performTest func(*networkpeering.NetworkPeer)) {
					testContainer := testAWSPeeringContainer("10.10.0.0/21")
					awsRegionName := aws.RegionCode(testContainer.ContainerRegion)
					vpcCIDR := "10.11.0.0/21"
					awsVPCid, err := aws.CreateVPC(utils.RandomName(testVPCName), vpcCIDR, awsRegionName)
					require.NoError(t, err)
					newContainer, err := nps.CreateContainer(ctx, testProjectID, testContainer)
					require.NoError(t, err)
					assert.NotEmpty(t, newContainer.ID)
					defer func() {
						require.NoError(t, aws.DeleteVPC(awsVPCid, awsRegionName))
					}()
					performTest(testAWSPeerConnection(t, newContainer.ID, vpcCIDR, awsVPCid))
				},
			},
			{
				provider: "AZURE",
				preparedCloudTest: func(performTest func(*networkpeering.NetworkPeer)) {
					testContainer := testAzurePeeringContainer("10.20.0.0/21")
					azureRegionName := azure.RegionCode(testContainer.ContainerRegion)
					vpcCIDR := "10.21.0.0/21"
					azureVPC, err := azure.CreateVPC(ctx, utils.RandomName(testVPCName), vpcCIDR, azureRegionName)
					require.NoError(t, err)
					newContainer, err := nps.CreateContainer(ctx, testProjectID, testContainer)
					require.NoError(t, err)
					assert.NotEmpty(t, newContainer.ID)
					defer func() {
						require.NoError(t, azure.DeleteVPC(ctx, azureVPC))
					}()
					performTest(testAzurePeerConnection(t, newContainer.ID, azureVPC))
				},
			},
			{
				provider: "GOOGLE",
				preparedCloudTest: func(performTest func(*networkpeering.NetworkPeer)) {
					testContainer := testGooglePeeringContainer("10.30.0.0/18")
					vpcName := utils.RandomName(testVPCName)
					require.NoError(t, google.CreateVPC(ctx, vpcName))
					newContainer, err := nps.CreateContainer(ctx, testProjectID, testContainer)
					require.NoError(t, err)
					assert.NotEmpty(t, newContainer.ID)
					defer func() {
						require.NoError(t, google.DeleteVPC(ctx, vpcName))
					}()
					performTest(testGooglePeerConnection(t, newContainer.ID, vpcName))
				},
			},
		} {
			tc.preparedCloudTest(func(peerRequest *networkpeering.NetworkPeer) {
				t.Run(fmt.Sprintf("create %s peer connection", tc.provider), func(t *testing.T) {
					newPeer, err := ps.CreatePeer(ctx, testProjectID, peerRequest)
					require.NoError(t, err)
					assert.NotEmpty(t, newPeer)
					createdPeer = newPeer
				})

				t.Run(fmt.Sprintf("list %s peer connections", tc.provider), func(t *testing.T) {
					containers, err := ps.ListPeers(ctx, testProjectID)
					require.NoError(t, err)
					assert.NotEmpty(t, containers)
					assert.GreaterOrEqual(t, len(containers), 1)
				})

				t.Run(fmt.Sprintf("delete %s peer connection", tc.provider), func(t *testing.T) {
					assert.NoError(t, ps.DeletePeer(ctx, testProjectID, createdPeer.ID))
				})
			})
		}
	})
}

func testAWSPeeringContainer(cidr string) *networkpeering.ProviderContainer {
	return &networkpeering.ProviderContainer{
		Provider: "AWS",
		AtlasProviderContainerConfig: akov2.AtlasProviderContainerConfig{
			ContainerRegion: "US_EAST_1",
			AtlasCIDRBlock:  cidr,
		},
	}
}

func testAzurePeeringContainer(cidr string) *networkpeering.ProviderContainer {
	return &networkpeering.ProviderContainer{
		Provider: "AZURE",
		AtlasProviderContainerConfig: akov2.AtlasProviderContainerConfig{
			ContainerRegion: "US_EAST_2",
			AtlasCIDRBlock:  cidr,
		},
	}
}

func testGooglePeeringContainer(cidr string) *networkpeering.ProviderContainer {
	return &networkpeering.ProviderContainer{
		Provider: "GCP",
		AtlasProviderContainerConfig: akov2.AtlasProviderContainerConfig{
			AtlasCIDRBlock: cidr,
		},
	}
}

func testAWSPeerConnection(t *testing.T, containerID string, vpcCIDR, vpcID string) *networkpeering.NetworkPeer {
	return &networkpeering.NetworkPeer{
		AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
			Provider:    "AWS",
			ContainerID: containerID,
			AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
				AWSAccountID:        mustHaveEnvVar(t, "AWS_ACCOUNT_ID"),
				AccepterRegionName:  "us-east-1",
				RouteTableCIDRBlock: vpcCIDR,
				VpcID:               vpcID,
			},
		},
	}
}

func testAzurePeerConnection(t *testing.T, containerID string, vpcName string) *networkpeering.NetworkPeer {
	return &networkpeering.NetworkPeer{
		AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
			Provider:    "AZURE",
			ContainerID: containerID,
			AzureConfiguration: &akov2.AzureNetworkPeeringConfiguration{
				AzureDirectoryID:    mustHaveEnvVar(t, "AZURE_TENANT_ID"),
				AzureSubscriptionID: mustHaveEnvVar(t, "AZURE_SUBSCRIPTION_ID"),
				ResourceGroupName:   azure.TestResourceGroupName(),
				VNetName:            vpcName,
			},
		},
	}
}

func testGooglePeerConnection(t *testing.T, containerID string, vpcName string) *networkpeering.NetworkPeer {
	return &networkpeering.NetworkPeer{
		AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
			Provider:    "GCP",
			ContainerID: containerID,
			GCPConfiguration: &akov2.GCPNetworkPeeringConfiguration{
				GCPProjectID: mustHaveEnvVar(t, "GOOGLE_PROJECT_ID"),
				NetworkName:  vpcName,
			},
		},
	}
}

func mustHaveEnvVar(t *testing.T, name string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		t.Fatalf("Unexpected unset env var %q", name)
	}
	return value
}

func ignoreRemoved(err error) error {
	if admin.IsErrorCode(err, "CLOUD_PROVIDER_CONTAINER_NOT_FOUND") {
		return nil
	}
	return err
}
