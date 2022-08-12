package e2e_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/cloud"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/kube"
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
	var data model.TestDataProvider

	_ = BeforeEach(func() {
		Eventually(kubecli.GetVersionOutput()).Should(Say(K8sVersion))
		checkUpAWSEnviroment()
		checkUpAzureEnviroment()
		checkNSetUpGCPEnviroment()
	})

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Network Peering Test\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + data.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			By("Save logs to output directory ", func() {
				GinkgoWriter.Write([]byte("Test has been failed. Trying to save logs...\n"))
				utils.SaveToFile(
					fmt.Sprintf("output/%s/operatorDecribe.txt", data.Resources.Namespace),
					[]byte(kubecli.DescribeOperatorPod(data.Resources.Namespace)),
				)
				utils.SaveToFile(
					fmt.Sprintf("output/%s/operator-logs.txt", data.Resources.Namespace),
					kubecli.GetManagerLogs(data.Resources.Namespace),
				)
				actions.SaveTestAppLogs(data.Resources)
				actions.SaveK8sResources(
					[]string{"deploy", "atlasprojects"},
					data.Resources.Namespace,
				)
			})
		}
		By("Clean Cloud", func() {
			DeleteAllNetworkPeering(&data)
		})
		By("Delete Resources, Project with NetworkPeering", func() {
			actions.DeleteUserResourcesProject(&data)
			actions.DeleteGlobalKeyIfExist(data)
		})
	})

	DescribeTable("NetworkPeering",
		func(test model.TestDataProvider, networkPeers []v1.NetworkPeer) {
			data = test
			networkPeerFlow(&data, networkPeers)
		},
		Entry("Test[networkpeering-aws-1]: User has project which was updated with AWS PrivateEndpoint",
			Label("network-peering-aws-1"),
			model.NewTestDataProvider(
				"networkpeering-aws-1",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlasdeployment_backup.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("user1").
						WithSecretRef("dbuser-secret-u1").
						AddBuildInAdminRole(),
				},
				40000,
				[]func(*model.TestDataProvider){},
			),
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
			model.NewTestDataProvider(
				"networkpeering-aws-2",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlasdeployment_backup.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("user1").
						WithSecretRef("dbuser-secret-u1").
						AddBuildInAdminRole(),
				},
				40000,
				[]func(*model.TestDataProvider){},
			),
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
			model.NewTestDataProvider(
				"networkpeering-aws-3",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlasdeployment_backup.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("user1").
						WithSecretRef("dbuser-secret-u1").
						AddBuildInAdminRole(),
				},
				40000,
				[]func(*model.TestDataProvider){},
			),
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
			model.NewTestDataProvider(
				"networkpeering-gcp-1",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlasdeployment_backup.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("user1").
						WithSecretRef("dbuser-secret-u1").
						AddBuildInAdminRole(),
				},
				40000,
				[]func(*model.TestDataProvider){},
			),
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
			model.NewTestDataProvider(
				"networkpeering-azure-1",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlasdeployment_backup.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("user1").
						WithSecretRef("dbuser-secret-u1").
						AddBuildInAdminRole(),
				},
				40000,
				[]func(*model.TestDataProvider){},
			),
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

func DeleteAllNetworkPeering(data *model.TestDataProvider) {
	errList := make([]error, 0)
	project, err := kube.GetProjectResource(data)
	Expect(err).ToNot(HaveOccurred())
	errors := networkpeer.DeletePeerVPC(project.Status.NetworkPeers)
	errList = append(errList, errors...)
	Expect(errList).To(BeEmpty())
}

func networkPeerFlow(userData *model.TestDataProvider, peers []v1.NetworkPeer) {
	By("Deploy Project with requested configuration", func() {
		actions.PrepareUsersConfigurations(userData)
		deploy.NamespacedOperator(userData)
		actions.DeployProjectAndWait(userData, "1")
	})

	By("Prepare network peers cloud infrastructure", func() {
		err := networkpeer.PreparePeerVPC(peers, userData.Resources.Namespace)
		Expect(err).ToNot(HaveOccurred())
	})

	By("Create network peers", func() {
		for _, peer := range peers {
			userData.Resources.Project.WithNetworkPeer(peer)
		}
		actions.PrepareUsersConfigurations(userData)
		actions.DeployProject(userData, "2")
	})

	By("Establish network peers connection", func() {
		Eventually(func() bool {
			return EnsurePeersReadyToConnect(*userData, len(peers))
		}).WithTimeout(5*time.Minute).WithPolling(20*time.Second).Should(BeTrue(), "Network Peering should be ready to establish connection")
		project, err := kube.GetProjectResource(userData)
		Expect(err).ShouldNot(HaveOccurred())
		err = networkpeer.EstablishPeerConnections(project.Status.NetworkPeers)
		Expect(err).ShouldNot(HaveOccurred())
		Eventually(kube.GetProjectNetworkPeerStatus(userData), "1m", "20s").Should(Equal("True"), "NetworkPeerStatus should be True")
	})

	By("Check network peers connection status state", func() {
		Eventually(kube.GetReadyProjectStatus(userData)).Should(Equal("True"), "Condition status 'Ready' is not 'True'")
		project, err := kube.GetProjectResource(userData)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(project.Status.NetworkPeers).Should(HaveLen(len(peers)))
	})
}

func EnsurePeersReadyToConnect(userData model.TestDataProvider, lenOfSpec int) bool {
	project, err := kube.GetProjectResource(&userData)
	if err != nil {
		return false
	}
	if len(project.Status.NetworkPeers) != lenOfSpec {
		return false
	}
	for _, networkPeering := range project.Status.NetworkPeers {
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
