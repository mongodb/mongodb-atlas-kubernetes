// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2e_test

import (
	"context"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/cloud"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

var _ = Describe("NetworkPeering", Label("networkpeering"), FlakeAttempts(2), func() {
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
		func(ctx SpecContext, test func(context.Context) *model.TestDataProvider, networkPeers []akov2.NetworkPeer) {
			testData = test(ctx)
			actions.ProjectCreationFlow(testData)
			networkPeerFlow(ctx, testData, networkPeers)
		},
		Entry("Test[networkpeering-aws-1]: User has project which was updated with AWS PrivateEndpoint",
			Label("focus-network-peering-aws-1"),
			func(ctx SpecContext) *model.TestDataProvider {
				return model.DataProvider(ctx, "networkpeering-aws-1", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			[]akov2.NetworkPeer{
				{
					ProviderName: provider.ProviderAWS,
					// Container config
					ContainerRegion: config.AWSRegionUS,
					AtlasCIDRBlock:  "10.8.0.0/22",
					// Peering config
					AccepterRegionName:  config.AWSRegionUS,
					RouteTableCIDRBlock: "10.0.0.0/24",
				},
			},
		),
		Entry("Test[networkpeering-aws-2]: User has project which was updated with AWS PrivateEndpoint",
			Label("focus-network-peering-aws-2"),
			func(ctx SpecContext) *model.TestDataProvider {
				return model.DataProvider(ctx, "networkpeering-aws-2", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			[]akov2.NetworkPeer{
				{
					ProviderName: provider.ProviderAWS,
					// Container config
					ContainerRegion: config.AWSRegionUS,
					AtlasCIDRBlock:  "10.8.0.0/22",
					// Peering config
					AccepterRegionName:  config.AWSRegionEU,
					RouteTableCIDRBlock: "10.0.0.0/24",
				},
			},
		),
		Entry("Test[networkpeering-aws-3]: User has project which was updated with AWS PrivateEndpoint",
			Label("focus-network-peering-aws-3"),
			func(ctx SpecContext) *model.TestDataProvider {
				return model.DataProvider(ctx, "networkpeering-aws-3", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			[]akov2.NetworkPeer{
				{
					ProviderName: provider.ProviderAWS,
					// Container config
					ContainerRegion: config.AWSRegionUS,
					AtlasCIDRBlock:  "10.8.0.0/22",
					// Peering config
					AccepterRegionName:  config.AWSRegionEU,
					RouteTableCIDRBlock: "192.168.0.0/16",
				},
				{
					ProviderName: provider.ProviderAWS,
					// Container config
					// Missing ContainerRegion would match AccepterRegionName
					AtlasCIDRBlock: "10.8.0.0/22",
					// Peering config
					AccepterRegionName:  config.AWSRegionUS,
					RouteTableCIDRBlock: "10.0.0.0/24",
				},
			},
		),
		Entry("Test[networkpeering-gcp-1]: User has project which was updated with GCP PrivateEndpoint",
			Label("focus-network-peering-gcp-1"),
			func(ctx SpecContext) *model.TestDataProvider {
				return model.DataProvider(ctx, "networkpeering-gcp-1", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			[]akov2.NetworkPeer{
				{
					ProviderName: provider.ProviderGCP,
					// Container config (no region setting for GCP)
					AtlasCIDRBlock: "10.8.0.0/18",
					// Peering config
					GCPProjectID: cloud.GoogleProjectID,
					NetworkName:  newRandomName(GCPVPCName),
				},
			},
		),
		Entry("Test[networkpeering-azure-1]: User has project which was updated with Azure PrivateEndpoint",
			Label("focus-network-peering-azure-1"),
			func(ctx SpecContext) *model.TestDataProvider {
				return model.DataProvider(ctx, "networkpeering-azure-1", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			[]akov2.NetworkPeer{
				{
					ProviderName: provider.ProviderAzure,
					// Container config
					ContainerRegion: "US_EAST_2",
					AtlasCIDRBlock:  "192.168.248.0/21",
					// Peering config
					AzureSubscriptionID: os.Getenv(SubscriptionID),
					AzureDirectoryID:    os.Getenv(DirectoryID),
					ResourceGroupName:   cloud.ResourceGroupName,
					VNetName:            newRandomName(AzureVPCName),
				},
			},
		),
	)
})

func networkPeerFlow(ctx context.Context, userData *model.TestDataProvider, peers []akov2.NetworkPeer) {
	providerActions := make([]cloud.Provider, len(peers))

	By("Prepare network peers cloud infrastructure", func() {
		for ix, peer := range peers {
			providerAction, err := prepareProviderAction(ctx)
			Expect(err).To(BeNil())
			providerActions[ix] = providerAction

			switch peer.ProviderName {
			case provider.ProviderAWS:
				peers[ix].AWSAccountID = providerActions[ix].GetAWSAccountID(ctx)
				cfg := &cloud.AWSConfig{
					Region:        peer.AccepterRegionName,
					VPC:           newRandomName("ao-vpc-peering-e2e"),
					CIDR:          peer.RouteTableCIDRBlock,
					Subnets:       map[string]string{"ao-peering-e2e-subnet": peer.RouteTableCIDRBlock},
					EnableCleanup: true,
				}
				peers[ix].VpcID = providerActions[ix].SetupNetwork(ctx, peer.ProviderName, cloud.WithAWSConfig(cfg))
			case provider.ProviderGCP:
				cfg := &cloud.GCPConfig{
					Region:        peer.AccepterRegionName,
					VPC:           peer.NetworkName,
					EnableCleanup: true,
				}
				providerActions[ix].SetupNetwork(ctx, peer.ProviderName, cloud.WithGCPConfig(cfg))
			case provider.ProviderAzure:
				cfg := &cloud.AzureConfig{
					VPC:           peer.VNetName,
					EnableCleanup: true,
				}
				providerActions[ix].SetupNetwork(ctx, peer.ProviderName, cloud.WithAzureConfig(cfg))
			}
		}
	})

	By("Create network peers", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{
			Name:      userData.Project.Name,
			Namespace: userData.Project.Namespace,
		}, userData.Project)).Should(Succeed())
		userData.Project.Spec.NetworkPeers = append(userData.Project.Spec.NetworkPeers, peers...)
		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
	})

	By("Establish network peers connection", func() {
		Eventually(func(g Gomega) bool {
			return EnsureProjectPeersReadyToConnect(g, userData, len(peers))
		}).WithTimeout(15*time.Minute).WithPolling(20*time.Second).Should(BeTrue(), "Network Peering should be ready to establish connection")
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name, Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())

		for ix, peer := range userData.Project.Status.NetworkPeers {
			switch peer.ProviderName {
			case provider.ProviderAWS:
				providerActions[ix].SetupNetworkPeering(ctx, peer.ProviderName, peer.ConnectionID, "")
			case provider.ProviderGCP:
				providerActions[ix].SetupNetworkPeering(ctx, peer.ProviderName, peer.AtlasGCPProjectID, peer.AtlasNetworkName)
			}
		}
		actions.WaitForConditionsToBecomeTrue(userData, api.NetworkPeerReadyType, api.ReadyType)
	})

	By("Check network peers connection status state", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name, Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())
		Expect(userData.Project.Status.NetworkPeers).Should(HaveLen(len(peers)))
	})
}

func EnsureProjectPeersReadyToConnect(g Gomega, userData *model.TestDataProvider, lenOfSpec int) bool {
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
