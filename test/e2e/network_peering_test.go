package e2e_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/data"

	"k8s.io/apimachinery/pkg/types"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/cloud"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/networkpeer"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

const (
	statusPendingAcceptance = "PENDING_ACCEPTANCE"
	statusWaitingUser       = "WAITING_FOR_USER"
)

type networkPeerWithUpdate struct {
	InitialSpecForProject v1.NetworkPeer
	InitialSpecForCloud   *v1.NetworkPeer
}

func (peer *networkPeerWithUpdate) IsInitialMatch() bool {
	if peer.InitialSpecForCloud == nil {
		return true
	}
	return peer.InitialSpecForProject == *peer.InitialSpecForCloud
}

func newRandomVPCName(base string) string {
	randomSuffix := uuid.New().String()[0:6]
	return fmt.Sprintf("%s-%s", base, randomSuffix)
}

func TestName(t *testing.T) {
	fmt.Println(newRandomVPCName("test"))
}

var _ = Describe("NetworkPeering", Label("networkpeering"), func() {
	var testData *model.TestDataProvider

	_ = BeforeEach(func() {
		Eventually(kubecli.GetVersionOutput()).Should(Say(K8sVersion))
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
			By("Save logs to output directory ", func() {
				GinkgoWriter.Write([]byte("Test has been failed. Trying to save logs...\n"))
				utils.SaveToFile(
					fmt.Sprintf("output/%s/operatorDecribe.txt", testData.Resources.Namespace),
					[]byte(kubecli.DescribeOperatorPod(testData.Resources.Namespace)),
				)
				utils.SaveToFile(
					fmt.Sprintf("output/%s/operator-logs.txt", testData.Resources.Namespace),
					kubecli.GetManagerLogs(testData.Resources.Namespace),
				)
				actions.SaveTestAppLogs(testData.Resources)
				actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)
				actions.SaveK8sResources(
					[]string{"deploy"},
					testData.Resources.Namespace,
				)
			})
		}
		By("Clean Cloud", func() {
			DeleteAllNetworkPeering(testData)
		})
		By("Delete Resources, Project with NetworkPeering", func() {
			actions.DeleteTestDataProject(testData)
			actions.DeleteGlobalKeyIfExist(*testData)
		})
	})

	DescribeTable("NetworkPeering",
		func(test *model.TestDataProvider, networkPeers []networkPeerWithUpdate) {
			testData = test
			actions.ProjectCreationFlow(test)
			networkPeerFlowWithUpdate(test, networkPeers)
		},
		Entry("Test[networkpeering-aws-1]: User has project which was updated with AWS Network Peering",
			Label("network-peering-aws-1"),
			model.DataProvider(
				"networkpeering-aws-1",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]networkPeerWithUpdate{
				{
					InitialSpecForProject: v1.NetworkPeer{
						ProviderName:        "AWS",
						AccepterRegionName:  "us-east-2",
						RouteTableCIDRBlock: "10.0.0.0/24",
						AtlasCIDRBlock:      "10.8.0.0/22",
					},
				},
			},
		),
		Entry("Test[networkpeering-aws-2]: User has project which was updated with AWS Network Peering",
			Label("network-peering-aws-2"),

			model.DataProvider(
				"networkpeering-aws-2",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),

			[]networkPeerWithUpdate{
				{
					InitialSpecForProject: v1.NetworkPeer{
						ProviderName:        "AWS",
						AccepterRegionName:  config.AWSRegionEU,
						ContainerRegion:     config.AWSRegionUS,
						RouteTableCIDRBlock: "10.0.0.0/24",
						AtlasCIDRBlock:      "10.8.0.0/22",
					},
				},
			},
		),

		Entry("Test[networkpeering-aws-3]: User has project which was updated with AWS Network Peering",
			Label("network-peering-aws-3"),

			model.DataProvider(
				"networkpeering-aws-3",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),

			[]networkPeerWithUpdate{
				{
					InitialSpecForProject: v1.NetworkPeer{
						ProviderName:        "AWS",
						AccepterRegionName:  config.AWSRegionEU,
						ContainerRegion:     config.AWSRegionUS,
						RouteTableCIDRBlock: "192.168.0.0/16",
						AtlasCIDRBlock:      "10.8.0.0/22",
					},
				},

				{
					InitialSpecForProject: v1.NetworkPeer{
						ProviderName:        "AWS",
						AccepterRegionName:  "us-east-2",
						RouteTableCIDRBlock: "10.0.0.0/24",
						AtlasCIDRBlock:      "10.8.0.0/22",
					},
				},
			},
		),
		Entry("Test[networkpeering-gcp-1]: User has project which was updated with GCP Network Peering",
			Label("network-peering-gcp-1"),

			model.DataProvider(
				"networkpeering-gcp-1",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),

			[]networkPeerWithUpdate{
				{
					InitialSpecForProject: v1.NetworkPeer{
						ProviderName:        "GCP",
						AccepterRegionName:  config.GCPRegion,
						RouteTableCIDRBlock: "192.168.0.0/16",
						AtlasCIDRBlock:      "10.8.0.0/18",
						NetworkName:         newRandomVPCName("network-peering-gcp-1-vpc"),
						GCPProjectID:        cloud.GoogleProjectID,
					},
				},
			},
		),

		Entry("Test[networkpeering-azure-1]: User has project which was updated with Azure Network Peering",
			Label("network-peering-azure-1"),

			model.DataProvider(
				"networkpeering-azure-1",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),

			[]networkPeerWithUpdate{
				{
					InitialSpecForProject: v1.NetworkPeer{
						ProviderName:        "AZURE",
						AccepterRegionName:  "US_EAST_2",
						AtlasCIDRBlock:      "192.168.248.0/21",
						VNetName:            newRandomVPCName("test-vnet"),
						AzureSubscriptionID: os.Getenv(networkpeer.SubscriptionID),
						ResourceGroupName:   networkpeer.AzureResourceGroupName,
						AzureDirectoryID:    os.Getenv(networkpeer.DirectoryID),
					},
				},
			},
		),

		Entry("Test[networkpeering-aws-4]: User has project which was updated with AWS Network Peering",
			Label("network-peering-aws-4"),
			model.DataProvider(
				"networkpeering-aws-4",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]networkPeerWithUpdate{
				{
					InitialSpecForProject: v1.NetworkPeer{
						ProviderName:        "AWS",
						AccepterRegionName:  config.AWSRegionUS,
						RouteTableCIDRBlock: "10.0.0.0/24",
						AtlasCIDRBlock:      "10.8.0.0/22",
					},
					InitialSpecForCloud: &v1.NetworkPeer{
						ProviderName:        "AWS",
						AccepterRegionName:  config.AWSRegionEU,
						RouteTableCIDRBlock: "10.0.0.0/24",
						AtlasCIDRBlock:      "10.8.0.0/22",
					},
				},
			},
		),
	)

})

func DeleteAllNetworkPeering(testData *model.TestDataProvider) {
	errList := make([]error, 0)
	Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{Name: testData.Project.Name, Namespace: testData.Project.Namespace}, testData.Project)).ToNot(HaveOccurred())
	errors := networkpeer.DeletePeerVPC(testData.Project.Status.NetworkPeers)
	errList = append(errList, errors...)
	Expect(errList).To(BeEmpty())
}

func networkPeerFlowWithUpdate(userData *model.TestDataProvider, peers []networkPeerWithUpdate) {
	shouldUpdate := networkPeerFlow(userData, peers, true)
	if shouldUpdate {
		networkPeerFlow(userData, peers, false)
	}
}

func networkPeerFlow(userData *model.TestDataProvider, peers []networkPeerWithUpdate, createPeersInCloud bool) bool {
	isCorrect := true
	var projectPeers []v1.NetworkPeer
	var cloudPeers []v1.NetworkPeer

	By("Preparing network peer specs", func() {
		for _, peer := range peers {
			projectPeers = append(projectPeers, peer.InitialSpecForProject)
			if peer.IsInitialMatch() {
				cloudPeers = append(cloudPeers, peer.InitialSpecForProject)
			} else {
				isCorrect = false
				cloudPeers = append(cloudPeers, *peer.InitialSpecForCloud)
			}
		}
		Expect(len(projectPeers)).To(Equal(len(cloudPeers)), "Number of peers in project and cloud should be equal")
	})

	if createPeersInCloud {
		By("Prepare network peers cloud infrastructure", func() {
			err := networkpeer.PreparePeerVPC(cloudPeers)
			Expect(err).ToNot(HaveOccurred())
			for i, cPeer := range cloudPeers {
				projectPeers[i].AWSAccountID = cPeer.AWSAccountID
				projectPeers[i].VpcID = cPeer.VpcID
			}
		})
	}

	setNetworkPeeringSpec(userData, projectPeers)

	if isCorrect {
		acceptNetworkPeeringConnections(userData, projectPeers)
	} else {
		ensureNetworkPeeringNotReady(userData)
	}

	// Filling network peering specs with cloud data
	if !isCorrect {
		for i, peer := range peers {
			if !peer.IsInitialMatch() {
				peers[i].InitialSpecForProject = cloudPeers[i]
				peers[i].InitialSpecForCloud = nil
			}
		}
	}
	return !isCorrect
}

func setNetworkPeeringSpec(userData *model.TestDataProvider, peers []v1.NetworkPeer) {
	By("Create network peers", func() {
		Eventually(func() error {
			return updateNetworkPeers(userData, peers)
		}, 2*time.Minute, 10*time.Second).Should(Succeed())
	})
}

func updateNetworkPeers(userData *model.TestDataProvider, peers []v1.NetworkPeer) error {
	err := userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
		Namespace: userData.Project.Namespace}, userData.Project)
	if err != nil {
		return err
	}
	userData.Project.Spec.NetworkPeers = peers
	err = userData.K8SClient.Update(userData.Context, userData.Project)
	return err
}

func ensureNetworkPeeringNotReady(userData *model.TestDataProvider) {
	By("Check network peers connection false status state", func() {
		Eventually(func(g Gomega) bool {
			g.Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name, Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())
			for _, condition := range userData.Project.Status.Conditions {
				if condition.Type == status.NetworkPeerReadyType {
					if condition.Status == corev1.ConditionTrue {
						return false
					}
				}
			}
			return true
		}).Should(BeTrue(), "Network Peering should not be ready")
	})
}

func acceptNetworkPeeringConnections(userData *model.TestDataProvider, peers []v1.NetworkPeer) {
	By("Establish network peers connection", func() {
		Eventually(func(g Gomega) bool {
			return ensurePeersReadyToConnect(g, userData, len(peers))
		}).WithTimeout(7*time.Minute).WithPolling(20*time.Second).Should(BeTrue(), "Network Peering should be ready to establish connection")
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name, Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())
		Expect(networkpeer.EstablishPeerConnections(userData.Project.Status.NetworkPeers)).Should(Succeed())
	})

	By("Check network peers connection true status state", func() {
		actions.WaitForConditionsToBecomeTrue(userData, status.NetworkPeerReadyType, status.ReadyType)
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name, Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())
		Expect(userData.Project.Status.NetworkPeers).Should(HaveLen(len(peers)))
	})
}

func ensurePeersReadyToConnect(g Gomega, userData *model.TestDataProvider, lenOfSpec int) bool {
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
