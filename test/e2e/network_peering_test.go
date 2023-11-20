package e2e_test

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/helper/e2e/actions/cloud"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/helper/e2e/model"
)

const (
	statusPendingAcceptance = "PENDING_ACCEPTANCE"
	statusWaitingUser       = "WAITING_FOR_USER"
	SubscriptionID          = "AZURE_SUBSCRIPTION_ID"
	DirectoryID             = "AZURE_TENANT_ID"
	GCPVPCName              = "network-peering-gcp-1-vpc"
	AzureVPCName            = "test-vnet"
)

func newRandomName(base string) string {
	randomSuffix := uuid.New().String()[0:6]
	return fmt.Sprintf("%s-%s", base, randomSuffix)
}

var _ = Describe("NetworkPeering", Label("networkpeering"), func() {
	var testData *model.TestDataProvider

	_ = BeforeEach(OncePerOrdered, func() {
		checkUpAWSEnvironment()
		checkUpAzureEnvironment()
		checkNSetUpGCPEnvironment()
	})

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Network Peering Test\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Delete Resources, Project with NetworkPeering", func() {
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	DescribeTable("NetworkPeering",
		func(test *model.TestDataProvider, networkPeers []v1.NetworkPeer) {
			testData = test
			actions.ProjectCreationFlow(test)
			networkPeerFlow(test, networkPeers)
		},
		Entry("Test[networkpeering-aws-1]: User has project which was updated with AWS PrivateEndpoint",
			Label("network-peering-aws-1"),
			model.DataProvider(
				"networkpeering-aws-1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]v1.NetworkPeer{
				{
					ProviderName:        provider.ProviderAWS,
					AccepterRegionName:  config.AWSRegionUS,
					ContainerRegion:     config.AWSRegionUS,
					RouteTableCIDRBlock: "10.0.0.0/24",
					AtlasCIDRBlock:      "10.8.0.0/22",
				},
			},
		),
		Entry("Test[networkpeering-aws-2]: User has project which was updated with AWS PrivateEndpoint",
			Label("network-peering-aws-2"),
			model.DataProvider(
				"networkpeering-aws-2",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]v1.NetworkPeer{
				{
					ProviderName:        provider.ProviderAWS,
					AccepterRegionName:  config.AWSRegionEU,
					ContainerRegion:     config.AWSRegionUS,
					RouteTableCIDRBlock: "10.0.0.0/24",
					AtlasCIDRBlock:      "10.8.0.0/22",
				},
			},
		),
		Entry("Test[networkpeering-aws-3]: User has project which was updated with AWS PrivateEndpoint",
			Label("network-peering-aws-3"),
			model.DataProvider(
				"networkpeering-aws-3",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]v1.NetworkPeer{
				{
					ProviderName:        provider.ProviderAWS,
					AccepterRegionName:  config.AWSRegionEU,
					ContainerRegion:     config.AWSRegionUS,
					RouteTableCIDRBlock: "192.168.0.0/16",
					AtlasCIDRBlock:      "10.8.0.0/22",
				},
				{
					ProviderName:        provider.ProviderAWS,
					AccepterRegionName:  config.AWSRegionUS,
					RouteTableCIDRBlock: "10.0.0.0/24",
					AtlasCIDRBlock:      "10.8.0.0/22",
				},
			},
		),
		Entry("Test[networkpeering-gcp-1]: User has project which was updated with GCP PrivateEndpoint",
			Label("network-peering-gcp-1"),
			model.DataProvider(
				"networkpeering-gcp-1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]v1.NetworkPeer{
				{
					ProviderName:        provider.ProviderGCP,
					AccepterRegionName:  config.GCPRegion,
					RouteTableCIDRBlock: "192.168.0.0/16",
					AtlasCIDRBlock:      "10.8.0.0/18",
					NetworkName:         newRandomName(GCPVPCName),
					GCPProjectID:        cloud.GoogleProjectID,
				},
			},
		),
		Entry("Test[networkpeering-azure-1]: User has project which was updated with Azure PrivateEndpoint",
			Label("network-peering-azure-1"),
			model.DataProvider(
				"networkpeering-azure-1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]v1.NetworkPeer{
				{
					ProviderName:        provider.ProviderAzure,
					AccepterRegionName:  "US_EAST_2",
					AtlasCIDRBlock:      "192.168.248.0/21",
					VNetName:            newRandomName(AzureVPCName),
					AzureSubscriptionID: os.Getenv(SubscriptionID),
					ResourceGroupName:   cloud.ResourceGroupName,
					AzureDirectoryID:    os.Getenv(DirectoryID),
				},
			},
		),
	)

})

func networkPeerFlow(userData *model.TestDataProvider, peers []v1.NetworkPeer) {
	providerActions := make([]cloud.Provider, len(peers))

	By("Prepare network peers cloud infrastructure", func() {
		for ix, peer := range peers {
			providerAction, err := prepareProviderAction()
			Expect(err).To(BeNil())
			providerActions[ix] = providerAction

			switch peer.ProviderName {
			case provider.ProviderAWS:
				peers[ix].AWSAccountID = providerActions[ix].GetAWSAccountID()
				cfg := &cloud.AWSConfig{
					Region:        peer.AccepterRegionName,
					VPC:           newRandomName("ao-vpc-peering-e2e"),
					CIDR:          peer.RouteTableCIDRBlock,
					Subnets:       map[string]string{"ao-peering-e2e-subnet": peer.RouteTableCIDRBlock},
					EnableCleanup: true,
				}
				peers[ix].VpcID = providerActions[ix].SetupNetwork(peer.ProviderName, cloud.WithAWSConfig(cfg))
			case provider.ProviderGCP:
				cfg := &cloud.GCPConfig{
					Region:        peer.AccepterRegionName,
					VPC:           peer.NetworkName,
					EnableCleanup: true,
				}
				providerActions[ix].SetupNetwork(peer.ProviderName, cloud.WithGCPConfig(cfg))
			case provider.ProviderAzure:
				cfg := &cloud.AzureConfig{
					VPC:           peer.VNetName,
					EnableCleanup: true,
				}
				providerActions[ix].SetupNetwork(peer.ProviderName, cloud.WithAzureConfig(cfg))
			}
		}
	})

	By("Create network peers", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
			Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())
		userData.Project.Spec.NetworkPeers = append(userData.Project.Spec.NetworkPeers, peers...)
		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
	})

	By("Establish network peers connection", func() {
		Eventually(func(g Gomega) bool {
			return EnsurePeersReadyToConnect(g, userData, len(peers))
		}).WithTimeout(15*time.Minute).WithPolling(20*time.Second).Should(BeTrue(), "Network Peering should be ready to establish connection")
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name, Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())

		for ix, peer := range userData.Project.Status.NetworkPeers {
			switch peer.ProviderName {
			case provider.ProviderAWS:
				providerActions[ix].SetupNetworkPeering(peer.ProviderName, peer.ConnectionID, "")
			case provider.ProviderGCP:
				providerActions[ix].SetupNetworkPeering(peer.ProviderName, peer.AtlasGCPProjectID, peer.AtlasNetworkName)
			}
		}
		actions.WaitForConditionsToBecomeTrue(userData, status.NetworkPeerReadyType, status.ReadyType)
	})

	By("Check network peers connection status state", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name, Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())
		Expect(userData.Project.Status.NetworkPeers).Should(HaveLen(len(peers)))
	})
}

func EnsurePeersReadyToConnect(g Gomega, userData *model.TestDataProvider, lenOfSpec int) bool {
	g.Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name, Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())
	if len(userData.Project.Status.NetworkPeers) != lenOfSpec {
		return false
	}
	for _, networkPeering := range userData.Project.Status.NetworkPeers {
		if networkPeering.ProviderName == provider.ProviderAzure {
			continue
		}
		if networkPeering.GetStatus() != statusPendingAcceptance && networkPeering.GetStatus() != statusWaitingUser {
			return false
		}
	}
	By("Network peers are ready to connect", func() {})
	return true
}
