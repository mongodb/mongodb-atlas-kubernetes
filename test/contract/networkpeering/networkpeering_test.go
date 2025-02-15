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

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkcontainer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkpeering"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/cloud/aws"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/cloud/azure"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/cloud/google"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/contract"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
)

const (
	testVPCName = "ako-test-network-peering-vpc"
)

func TestPeerServiceCRUD(t *testing.T) {
	ctx := context.Background()
	contract.RunGoContractTest(ctx, t, "test peer CRUD", func(ch contract.ContractHelper) {
		projectName := utils.RandomName("peer-connection-crud-project")
		require.NoError(t, ch.AddResources(ctx, 5*time.Minute, contract.DefaultAtlasProject(projectName)))
		testProjectID, err := ch.ProjectID(ctx, projectName)
		require.NoError(t, err)
		ncs := networkcontainer.NewNetworkContainerService(ch.AtlasClient().NetworkPeeringApi)
		nps := networkpeering.NewNetworkPeeringService(ch.AtlasClient().NetworkPeeringApi)
		createdPeer := &networkpeering.NetworkPeer{}

		for _, tc := range []struct {
			provider          string
			preparedCloudTest func(func(containerID string, cfg *akov2.AtlasNetworkPeeringConfig))
		}{
			{
				provider: string(provider.ProviderAWS),
				preparedCloudTest: func(performTest func(string, *akov2.AtlasNetworkPeeringConfig)) {
					testContainer := testAWSPeeringContainer("10.10.0.0/21")
					awsRegionName := aws.RegionCode(testContainer.Region)
					vpcCIDR := "10.11.0.0/21"
					awsVPCid, err := aws.CreateVPC(utils.RandomName(testVPCName), vpcCIDR, awsRegionName)
					require.NoError(t, err)
					newContainer, err := ncs.Create(ctx, testProjectID, testContainer)
					require.NoError(t, err)
					assert.NotEmpty(t, newContainer.ID)
					defer func() {
						require.NoError(t, aws.DeleteVPC(awsVPCid, awsRegionName))
					}()
					performTest(newContainer.ID, testAWSPeerConnection(t, vpcCIDR, awsVPCid))
				},
			},
			{
				provider: string(provider.ProviderAzure),
				preparedCloudTest: func(performTest func(string, *akov2.AtlasNetworkPeeringConfig)) {
					testContainer := testAzurePeeringContainer("10.20.0.0/21")
					azureRegionName := azure.RegionCode(testContainer.Region)
					vpcCIDR := "10.21.0.0/21"
					azureVPC, err := azure.CreateVPC(ctx, utils.RandomName(testVPCName), vpcCIDR, azureRegionName)
					require.NoError(t, err)
					newContainer, err := ncs.Create(ctx, testProjectID, testContainer)
					require.NoError(t, err)
					assert.NotEmpty(t, newContainer.ID)
					defer func() {
						require.NoError(t, azure.DeleteVPC(ctx, azureVPC))
					}()
					performTest(newContainer.ID, testAzurePeerConnection(t, azureVPC))
				},
			},
			{
				provider: string(provider.ProviderGCP),
				preparedCloudTest: func(performTest func(string, *akov2.AtlasNetworkPeeringConfig)) {
					testContainer := testGooglePeeringContainer("10.30.0.0/18")
					vpcName := utils.RandomName(testVPCName)
					require.NoError(t, google.CreateVPC(ctx, vpcName))
					newContainer, err := ncs.Create(ctx, testProjectID, testContainer)
					require.NoError(t, err)
					assert.NotEmpty(t, newContainer.ID)
					defer func() {
						require.NoError(t, google.DeleteVPC(ctx, vpcName))
					}()
					performTest(newContainer.ID, testGooglePeerConnection(t, vpcName))
				},
			},
		} {
			tc.preparedCloudTest(func(containerID string, cfg *akov2.AtlasNetworkPeeringConfig) {
				t.Run(fmt.Sprintf("create %s peer connection", tc.provider), func(t *testing.T) {
					newPeer, err := nps.Create(ctx, testProjectID, containerID, cfg)
					require.NoError(t, err)
					assert.NotEmpty(t, newPeer)
					createdPeer = newPeer
				})

				t.Run(fmt.Sprintf("get %s peer connection", tc.provider), func(t *testing.T) {
					peer, err := nps.Get(ctx, testProjectID, createdPeer.ID)
					require.NoError(t, err)
					assert.Equal(t, createdPeer, peer)
				})

				t.Run(fmt.Sprintf("delete %s peer connection", tc.provider), func(t *testing.T) {
					assert.NoError(t, nps.Delete(ctx, testProjectID, createdPeer.ID))
				})
			})
		}
	})
}

func testAWSPeeringContainer(cidr string) *networkcontainer.NetworkContainerConfig {
	return &networkcontainer.NetworkContainerConfig{
		Provider: string(provider.ProviderAWS),
		AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
			Region:    "US_EAST_1",
			CIDRBlock: cidr,
		},
	}
}

func testAzurePeeringContainer(cidr string) *networkcontainer.NetworkContainerConfig {
	return &networkcontainer.NetworkContainerConfig{
		Provider: string(provider.ProviderAzure),
		AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
			Region:    "US_EAST_2",
			CIDRBlock: cidr,
		},
	}
}

func testGooglePeeringContainer(cidr string) *networkcontainer.NetworkContainerConfig {
	return &networkcontainer.NetworkContainerConfig{
		Provider: string(provider.ProviderGCP),
		AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
			CIDRBlock: cidr,
		},
	}
}

func testAWSPeerConnection(t *testing.T, vpcCIDR, vpcID string) *akov2.AtlasNetworkPeeringConfig {
	return &akov2.AtlasNetworkPeeringConfig{
		Provider: string(provider.ProviderAWS),
		AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
			AWSAccountID:        mustHaveEnvVar(t, "AWS_ACCOUNT_ID"),
			AccepterRegionName:  "us-east-1",
			RouteTableCIDRBlock: vpcCIDR,
			VpcID:               vpcID,
		},
	}
}

func testAzurePeerConnection(t *testing.T, vpcName string) *akov2.AtlasNetworkPeeringConfig {
	return &akov2.AtlasNetworkPeeringConfig{
		Provider: string(provider.ProviderAzure),
		AzureConfiguration: &akov2.AzureNetworkPeeringConfiguration{
			AzureDirectoryID:    mustHaveEnvVar(t, "AZURE_TENANT_ID"),
			AzureSubscriptionID: mustHaveEnvVar(t, "AZURE_SUBSCRIPTION_ID"),
			ResourceGroupName:   azure.TestResourceGroupName(),
			VNetName:            vpcName,
		},
	}
}

func testGooglePeerConnection(t *testing.T, vpcName string) *akov2.AtlasNetworkPeeringConfig {
	return &akov2.AtlasNetworkPeeringConfig{
		Provider: string(provider.ProviderGCP),
		GCPConfiguration: &akov2.GCPNetworkPeeringConfiguration{
			GCPProjectID: mustHaveEnvVar(t, "GOOGLE_PROJECT_ID"),
			NetworkName:  vpcName,
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
