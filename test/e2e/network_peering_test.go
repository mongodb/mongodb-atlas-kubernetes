package e2e_test

import (
	"fmt"
	"log"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/kube"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/networkpeer"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/deploy"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

type networkPeer struct {
	provider        string
	peerRegion      string
	containerRegion string
	atlasCidr       string
	appCidr         string
	accountID       string
	vpc             string
}

var _ = Describe("NetworkPeering", Label("networkpeering"), func() {
	var data model.TestDataProvider

	_ = BeforeEach(func() {
		Eventually(kubecli.GetVersionOutput()).Should(Say(K8sVersion))
		checkUpAWSEnviroment()
		//checkUpAzureEnviroment()
		//checkNSetUpGCPEnviroment()
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

	DescribeTable("aueoueo",
		func(test model.TestDataProvider, networkPeers []networkPeer) {
			data = test
			networkPeerFlow(&data, networkPeers)
		},
		Entry("Test[networkpeering-aws-1]: User has project which was updated with AWS PrivateEndpoint", Label("network-peering-aws-1"),
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
			[]networkPeer{
				{
					provider:        "AWS",
					peerRegion:      config.AWSRegionUS,
					containerRegion: config.AWSRegionUS,
					appCidr:         "10.0.0.0/24",
					atlasCidr:       "10.8.0.0/22",
				},
			},
		),
		Entry("Test[networkpeering-aws-2]: User has project which was updated with AWS PrivateEndpoint", Label("network-peering-aws-1"),
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
			[]networkPeer{
				{
					provider:        "AWS",
					peerRegion:      config.AWSRegionEU,
					containerRegion: config.AWSRegionUS,
					appCidr:         "10.0.0.0/24",
					atlasCidr:       "10.8.0.0/22",
				},
			},
		),
		Entry("Test[networkpeering-aws-3]: User has project which was updated with AWS PrivateEndpoint", Label("network-peering-aws-1"),
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
			[]networkPeer{
				{
					provider:        "AWS",
					peerRegion:      config.AWSRegionEU,
					containerRegion: config.AWSRegionUS,
					appCidr:         "192.168.0.0/16",
					atlasCidr:       "10.8.0.0/22",
				},
				{
					provider:   "AWS",
					peerRegion: config.AWSRegionUS,
					appCidr:    "10.0.0.0/24",
					atlasCidr:  "10.8.0.0/22",
				},
			},
		),
	)

})

func DeleteAllNetworkPeering(data *model.TestDataProvider) {
	errList := make([]error, 0)
	project, err := kube.GetProjectResource(data)
	Expect(err).ToNot(HaveOccurred())

	for _, networkPeering := range project.Status.NetworkPeers {
		err := networkpeer.DeletePeerConnectionAndVPC(networkPeering.ConnectionID, networkPeering.Region)
		if err != nil {
			errList = append(errList, err)
		}
	}

	Expect(errList).To(BeEmpty())
}

func networkPeerFlow(userData *model.TestDataProvider, peers []networkPeer) {
	By("Deploy Project with requested configuration", func() {
		actions.PrepareUsersConfigurations(userData)
		deploy.NamespacedOperator(userData)
		actions.DeployProjectAndWait(userData, "1")
	})

	By("Prepare network peers cloud infrastructure", func() {
		for i, peer := range peers {
			awsNetworkPeer, err := networkpeer.NewAWSNetworkPeer(peer.peerRegion)
			Expect(err).ShouldNot(HaveOccurred())
			testID := fmt.Sprintf("%s-%d", userData.Resources.Namespace, i)
			accountID, vpcID, err := awsNetworkPeer.CreateVPC(peer.appCidr, testID) // TODO: get it from somewhere
			Expect(err).ShouldNot(HaveOccurred())
			peers[i].accountID = accountID
			peers[i].vpc = vpcID
		}
	})

	By("Create network peers", func() {
		for _, peer := range peers {
			log.Printf("Creating network peer for %s", peer.provider)
			networkPeerSpec := peerToPeerSpec(peer)
			userData.Resources.Project.WithNetworkPeer(networkPeerSpec)
		}
		log.Printf("peer spec %v", userData.Resources.Project.Spec.NetworkPeers)
		actions.PrepareUsersConfigurations(userData)
		actions.DeployProject(userData, "2")
	})

	By("Establish network peers connection", func() {
		Eventually(func() bool {
			project, err := kube.GetProjectResource(userData)
			if err != nil {
				return false
			}
			if len(project.Status.NetworkPeers) != len(peers) {
				return false
			}
			for _, networkPeering := range project.Status.NetworkPeers {
				if networkPeering.GetStatus() != "PENDING_ACCEPTANCE" {
					return false
				}
			}
			return true //TODO: move it to function. and use constant
		}, "3m").Should(BeTrue(), "Network Peering should be ready to establish connection")

		project, err := kube.GetProjectResource(userData)
		Expect(err).ShouldNot(HaveOccurred())
		for _, peerStatus := range project.Status.NetworkPeers {
			errEstablish := networkpeer.EstablishPeerConnection(peerStatus)
			Expect(errEstablish).ShouldNot(HaveOccurred())
		}
		Eventually(kube.GetProjectNetworkPeerStatus(userData), "1m", "20s").Should(Equal("True"), "NetworkPeerStatus should be True")
	})

	By("Check network peers connection status state", func() {
		Eventually(kube.GetReadyProjectStatus(userData)).Should(Equal("True"), "Condition status 'Ready' is not 'True'")
		project, err := kube.GetProjectResource(userData)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(project.Status.NetworkPeers).Should(HaveLen(len(peers)))
	})
}

func peerToPeerSpec(peer networkPeer) v1.NetworkPeer {
	return v1.NetworkPeer{
		AccepterRegionName:  peer.peerRegion,
		ContainerRegion:     peer.containerRegion,
		AWSAccountID:        peer.accountID,
		ProviderName:        provider.ProviderName(peer.provider),
		RouteTableCIDRBlock: peer.appCidr,
		VpcID:               peer.vpc,
		AtlasCIDRBlock:      peer.atlasCidr,
	}
}
