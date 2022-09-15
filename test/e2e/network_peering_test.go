package e2e_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/kube"

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
		func(test *model.TestDataProvider, networkPeers []v1.NetworkPeer) {
			testData = test
			actions.ProjectCreationFlow(test)
			networkPeerFlow(test, networkPeers)
		},
		Entry("Test[networkpeering-aws-1]: User has project which was updated with AWS PrivateEndpoint",
			Label("network-peering-aws-1"),
			model.DataProvider(
				"networkpeering-aws-1",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]v1.NetworkPeer{
				{
					ProviderName:        "AWS",
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
				"networkpeering-aws-1",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]v1.NetworkPeer{
				{
					ProviderName:        "AWS",
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
				"networkpeering-aws-1",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]v1.NetworkPeer{
				{
					ProviderName:        "AWS",
					AccepterRegionName:  config.AWSRegionEU,
					ContainerRegion:     config.AWSRegionUS,
					RouteTableCIDRBlock: "192.168.0.0/16",
					AtlasCIDRBlock:      "10.8.0.0/22",
				},
				{
					ProviderName:        "AWS",
					AccepterRegionName:  config.AWSRegionUS,
					RouteTableCIDRBlock: "10.0.0.0/24",
					AtlasCIDRBlock:      "10.8.0.0/22",
				},
			},
		),
		Entry("Test[networkpeering-gcp-1]: User has project which was updated with GCP PrivateEndpoint",
			Label("network-peering-gcp-1"),
			model.DataProvider(
				"networkpeering-aws-1",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]v1.NetworkPeer{
				{
					ProviderName:        "GCP",
					AccepterRegionName:  config.GCPRegion,
					RouteTableCIDRBlock: "192.168.0.0/16",
					AtlasCIDRBlock:      "10.8.0.0/18",
					NetworkName:         newRandomVPCName("network-peering-gcp-1-vpc"),
					GCPProjectID:        cloud.GoogleProjectID,
				},
			},
		),
		Entry("Test[networkpeering-azure-1]: User has project which was updated with Azure PrivateEndpoint",
			Label("network-peering-azure-1"),
			model.DataProvider(
				"networkpeering-aws-1",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]v1.NetworkPeer{
				{
					ProviderName:        "AZURE",
					AccepterRegionName:  "US_EAST_2",
					AtlasCIDRBlock:      "192.168.248.0/21",
					VNetName:            newRandomVPCName("test-vnet"),
					AzureSubscriptionID: os.Getenv(networkpeer.SubscriptionID),
					ResourceGroupName:   networkpeer.AzureResourceGroupName,
					AzureDirectoryID:    os.Getenv(networkpeer.DirectoryID),
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

func networkPeerFlow(userData *model.TestDataProvider, peers []v1.NetworkPeer) {
	By("Prepare network peers cloud infrastructure", func() {
		err := networkpeer.PreparePeerVPC(peers)
		Expect(err).ToNot(HaveOccurred())
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
		}).WithTimeout(5*time.Minute).WithPolling(20*time.Second).Should(BeTrue(), "Network Peering should be ready to establish connection")
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name, Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())
		Expect(networkpeer.EstablishPeerConnections(userData.Project.Status.NetworkPeers)).Should(Succeed())
		Eventually(func(g Gomega) string {
			condition, err := kube.GetProjectStatusCondition(userData, status.NetworkPeerReadyType)
			g.Expect(err).ShouldNot(HaveOccurred())
			return condition
		}).WithTimeout(3*time.Minute).WithPolling(20*time.Second).Should(Equal("True"), "Cloud Access Roles status should be True")
	})

	By("Check network peers connection status state", func() {
		Eventually(func(g Gomega) string {
			condition, err := kube.GetProjectStatusCondition(userData, status.ReadyType)
			g.Expect(err).ShouldNot(HaveOccurred())
			return condition
		}).Should(Equal("True"), "Condition status 'Ready' is not 'True'")
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
