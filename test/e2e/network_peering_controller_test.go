package e2e_test

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/cloud"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
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

type netPeerTestCase struct {
	ProviderName provider.ProviderName
	AWS          *akov2.AWSNetworkPeeringConfiguration
	Azure        *akov2.AzureNetworkPeeringConfiguration
	GCP          *akov2.GCPNetworkPeeringConfiguration
	Container    akov2.AtlasProviderContainerConfig
}

var _ = Describe("NetworkPeeringController", Label("networkpeering-controller"), func() {
	var testData *model.TestDataProvider

	_ = BeforeEach(OncePerOrdered, func() {
		checkUpAWSEnvironment()
		checkUpAzureEnvironment()
		checkNSetUpGCPEnvironment()
	})

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Network Peering Controller Test\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Delete Resources, Project with NetworkPeering", func() {
			actions.DeleteTestDataNetworkPeerings(testData)
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	DescribeTable("NetworkPeeringController",
		func(test *model.TestDataProvider, networkPeers []*akov2.AtlasNetworkPeering) {
			testData = test
			actions.ProjectCreationFlow(test)
			networkPeerControllerFlow(test, networkPeers)
		},
		Entry("Test[networkpeering-aws-1]: AWS Network Peering CR within a region and without existent Atlas Container",
			Label("network-peering-cr-aws-1"),
			model.DataProvider(
				"networkpeering-cr-aws-1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]*akov2.AtlasNetworkPeering{
				newPeeringTestCase(netPeerTestCase{
					ProviderName: provider.ProviderAWS,
					Container: akov2.AtlasProviderContainerConfig{
						ContainerRegion: "US_EAST_1",
						AtlasCIDRBlock:  "10.8.0.0/22",
					},
					AWS: &akov2.AWSNetworkPeeringConfiguration{
						AccepterRegionName:  "us-east-1", // AccepterRegionName uses AWS region names
						AWSAccountID:        os.Getenv("AWS_ACCOUNT_ID"),
						RouteTableCIDRBlock: "10.0.0.0/24",
					},
				}),
			},
		),
		Entry("Test[networkpeering-aws-2]: AWS Network Peering CR between different regions and without existent Atlas Container",
			Label("network-peering-cr-aws-2"),
			model.DataProvider(
				"networkpeering-aws-2",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]*akov2.AtlasNetworkPeering{
				newPeeringTestCase(netPeerTestCase{
					ProviderName: provider.ProviderAWS,
					Container: akov2.AtlasProviderContainerConfig{
						ContainerRegion: "US_EAST_1",
						AtlasCIDRBlock:  "10.8.0.0/22",
					},
					AWS: &akov2.AWSNetworkPeeringConfiguration{
						AccepterRegionName:  "eu-west-2",
						AWSAccountID:        os.Getenv("AWS_ACCOUNT_ID"),
						RouteTableCIDRBlock: "10.0.0.0/24",
					},
				}),
			},
		),
		Entry("Test[networkpeering-aws-3]: AWS Network Peering CRs between different regions and without container region specified",
			Label("network-peering-cr-aws-3"),
			model.DataProvider(
				"networkpeering-aws-3",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]*akov2.AtlasNetworkPeering{
				newPeeringTestCase(netPeerTestCase{
					ProviderName: provider.ProviderAWS,
					Container: akov2.AtlasProviderContainerConfig{
						ContainerRegion: "US_EAST_1",
						AtlasCIDRBlock:  "10.8.0.0/22",
					},
					AWS: &akov2.AWSNetworkPeeringConfiguration{
						AccepterRegionName:  "eu-west-1",
						AWSAccountID:        os.Getenv("AWS_ACCOUNT_ID"),
						RouteTableCIDRBlock: "192.168.0.0/16",
					},
				}),
				newPeeringTestCase(netPeerTestCase{
					ProviderName: provider.ProviderAWS,
					Container: akov2.AtlasProviderContainerConfig{
						AtlasCIDRBlock: "10.8.0.0/22",
					},
					AWS: &akov2.AWSNetworkPeeringConfiguration{
						AccepterRegionName:  "us-east-1",
						AWSAccountID:        os.Getenv("AWS_ACCOUNT_ID"),
						RouteTableCIDRBlock: "10.0.0.0/24",
					},
				}),
			},
		),
		Entry("Test[networkpeering-gcp-1]: GCP Network Peering CR",
			Label("network-peering-cr-gcp-1"),
			model.DataProvider(
				"networkpeering-gcp-1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]*akov2.AtlasNetworkPeering{
				newPeeringTestCase(netPeerTestCase{
					ProviderName: provider.ProviderGCP,
					Container: akov2.AtlasProviderContainerConfig{
						AtlasCIDRBlock: "10.8.0.0/18",
					},
					GCP: &akov2.GCPNetworkPeeringConfiguration{
						GCPProjectID: cloud.GoogleProjectID,
						NetworkName:  newRandomName(GCPVPCName),
					},
				}),
			},
		),
		Entry("Test[networkpeering-azure-1]: User has project which was updated with Azure PrivateEndpoint",
			Label("network-peering-cr-azure-1"),
			model.DataProvider(
				"networkpeering-azure-1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]*akov2.AtlasNetworkPeering{
				newPeeringTestCase(netPeerTestCase{
					ProviderName: provider.ProviderAzure,
					Container: akov2.AtlasProviderContainerConfig{
						AtlasCIDRBlock:  "192.168.248.0/21",
						ContainerRegion: "US_EAST_2",
					},
					Azure: &akov2.AzureNetworkPeeringConfiguration{
						AzureDirectoryID:    os.Getenv(DirectoryID),
						AzureSubscriptionID: os.Getenv(SubscriptionID),
						ResourceGroupName:   cloud.ResourceGroupName,
						VNetName:            newRandomName(AzureVPCName),
					},
				}),
			},
		),
	)
})

func networkPeerControllerFlow(userData *model.TestDataProvider, peers []*akov2.AtlasNetworkPeering) {
	providerActions := make([]cloud.Provider, len(peers))

	By("Prepare network peers cloud infrastructure from CRs", func() {
		for ix, peer := range peers {
			providerAction, err := prepareProviderAction()
			Expect(err).To(BeNil())
			providerActions[ix] = providerAction

			providerName := provider.ProviderName(peer.Spec.Provider)
			switch providerName {
			case provider.ProviderAWS:
				peers[ix].Spec.AWSConfiguration.AWSAccountID = providerActions[ix].GetAWSAccountID()
				cfg := &cloud.AWSConfig{
					Region:        peer.Spec.AWSConfiguration.AccepterRegionName,
					VPC:           newRandomName("ao-vpc-peering-e2e"),
					CIDR:          peer.Spec.AWSConfiguration.RouteTableCIDRBlock,
					Subnets:       map[string]string{"ao-peering-e2e-subnet": peer.Spec.AWSConfiguration.RouteTableCIDRBlock},
					EnableCleanup: true,
				}
				peers[ix].Spec.AWSConfiguration.VpcID = providerActions[ix].SetupNetwork(providerName, cloud.WithAWSConfig(cfg))
			case provider.ProviderGCP:
				cfg := &cloud.GCPConfig{
					// Region:        peer.Spec.GCPConfiguration.AccepterRegionName, GCP does not use regions
					VPC:           peer.Spec.GCPConfiguration.NetworkName,
					EnableCleanup: true,
				}
				providerActions[ix].SetupNetwork(providerName, cloud.WithGCPConfig(cfg))
			case provider.ProviderAzure:
				cfg := &cloud.AzureConfig{
					VPC:           peer.Spec.AzureConfiguration.VNetName,
					EnableCleanup: true,
				}
				providerActions[ix].SetupNetwork(providerName, cloud.WithAzureConfig(cfg))
			}
		}
	})

	By("Create network peer from CRs", func() {
		for i, peer := range peers {
			peer.Spec.ProjectRef.Name = userData.Project.Name
			peer.Spec.ProjectRef.Namespace = userData.Project.Namespace
			peer.Name = fmt.Sprintf("%s-item-%d", userData.Prefix, i)
			peer.Namespace = userData.Project.Namespace
			Expect(userData.K8SClient.Create(userData.Context, peer)).Should(Succeed())
		}
	})

	By("Establish network peers connection with CRs", func() {
		Eventually(func(g Gomega) bool {
			return EnsurePeersReadyToConnect(g, userData, peers)
		}).WithTimeout(15*time.Minute).WithPolling(20*time.Second).Should(
			BeTrue(),
			"Network Peering CRs should be ready to establish connection",
		)

		for ix, peer := range peers {
			providerName := provider.ProviderName(peer.Spec.Provider)
			switch providerName {
			case provider.ProviderAWS:
				providerActions[ix].SetupNetworkPeering(
					providerName,
					peer.Status.AWSStatus.ConnectionID,
					"",
				)
			case provider.ProviderGCP:
				providerActions[ix].SetupNetworkPeering(
					providerName,
					peer.Status.GoogleStatus.GCPProjectID,
					peer.Status.GoogleStatus.NetworkName,
				)
			}
			key := types.NamespacedName{Name: peer.Name, Namespace: peer.Namespace}
			Eventually(func(g Gomega) bool {
				Expect(userData.K8SClient.Get(userData.Context, key, peer)).Should(Succeed())
				return peer.Status.Status == "AVAILABLE"
			}).WithTimeout(15*time.Minute).WithPolling(5*time.Second).Should(
				BeTrue(),
				"Network Peering CRs should become available",
			)
		}
	})

	By("Check network peers CRs to be Ready", func() {
		for _, peer := range peers {
			key := types.NamespacedName{Name: peer.Name, Namespace: peer.Namespace}
			Expect(userData.K8SClient.Get(userData.Context, key, peer)).Should(Succeed())
			Expect(networkPeeringReady(peer)).Should(BeTrue())
		}
	})
}

func EnsurePeersReadyToConnect(g Gomega, userData *model.TestDataProvider, peers []*akov2.AtlasNetworkPeering) bool {
	for _, peer := range peers {
		key := types.NamespacedName{Name: peer.Name, Namespace: peer.Namespace}
		g.Expect(userData.K8SClient.Get(userData.Context, key, peer)).Should(Succeed())
		if peer.Spec.Provider == string(provider.ProviderAzure) {
			continue
		}
		statusMsg := peer.Status.Status
		if statusMsg != statusPendingAcceptance && statusMsg != statusWaitingUser {
			return false
		}
	}
	By("Network peers are ready to connect", func() {})
	return true
}

func newPeeringTestCase(tc netPeerTestCase) *akov2.AtlasNetworkPeering {
	np := &akov2.AtlasNetworkPeering{
		Spec: akov2.AtlasNetworkPeeringSpec{
			ProjectDualReference: akov2.ProjectDualReference{
				ProjectRef: &common.ResourceRefNamespaced{
					Name: data.ProjectName,
				},
			},
			AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
				Provider: string(tc.ProviderName),
			},
			AtlasProviderContainerConfig: tc.Container,
		},
	}
	switch tc.ProviderName {
	case provider.ProviderAWS:
		np.Spec.AWSConfiguration = tc.AWS
	case provider.ProviderAzure:
		np.Spec.AzureConfiguration = tc.Azure
	case provider.ProviderGCP:
		np.Spec.GCPConfiguration = tc.GCP
	}
	return np
}

func networkPeeringReady(peer *akov2.AtlasNetworkPeering) bool {
	for _, condition := range peer.Status.Conditions {
		GinkgoWriter.Printf("TODO: REMOVE LOG peer %s condition type=%s status=%s",
			peer.Status.ID, condition.Type, condition.Status)
		if condition.Type == api.ReadyType && condition.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}
