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
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

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

type containerAndPeering struct {
	container *akov2.AtlasNetworkContainer
	peering   *akov2.AtlasNetworkPeering
}

var _ = Describe("NetworkPeeringController", Label("networkpeering-controller"), FlakeAttempts(3), func() {
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
			actions.DeleteTestDataNetworkContainers(testData)
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	DescribeTable("NetworkPeeringController",
		func(ctx SpecContext, test func(context.Context) *model.TestDataProvider, pairs []containerAndPeering) {
			testData = test(ctx)
			actions.ProjectCreationFlow(testData)
			networkPeerControllerFlow(ctx, testData, pairs)
		},
		Entry("Test[networkpeering-aws-1]: AWS Network Peering CR within a region and without existent Atlas Container",
			Label("focus-network-peering-cr-aws-1"),
			func(ctx SpecContext) *model.TestDataProvider {
				return model.DataProvider(
					ctx,
					"network-peering-cr-aws-1",
					model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
					40000,
					[]func(*model.TestDataProvider){},
				).WithProject(data.DefaultProject())
			},
			[]containerAndPeering{
				{
					container: &akov2.AtlasNetworkContainer{
						Spec: akov2.AtlasNetworkContainerSpec{
							Provider: string(provider.ProviderAWS),
							AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
								Region:    "US_EAST_1",
								CIDRBlock: "10.8.0.0/22",
							},
						},
					},
					peering: &akov2.AtlasNetworkPeering{
						Spec: akov2.AtlasNetworkPeeringSpec{
							AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
								Provider: string(provider.ProviderAWS),
								AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
									AccepterRegionName:  "us-east-1", // AccepterRegionName uses AWS region names
									AWSAccountID:        os.Getenv("AWS_ACCOUNT_ID"),
									RouteTableCIDRBlock: "10.0.0.0/24",
								},
							},
						},
					},
				},
			},
		),
		Entry("Test[networkpeering-aws-2]: AWS Network Peering CR between different regions and without existent Atlas Container",
			Label("focus-network-peering-cr-aws-2"),
			func(ctx SpecContext) *model.TestDataProvider {
				return model.DataProvider(ctx, "network-peering-cr-aws-2", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			[]containerAndPeering{
				{
					container: &akov2.AtlasNetworkContainer{
						Spec: akov2.AtlasNetworkContainerSpec{
							Provider: string(provider.ProviderAWS),
							AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
								Region:    "US_EAST_1",
								CIDRBlock: "10.8.0.0/22",
							},
						},
					},
					peering: &akov2.AtlasNetworkPeering{
						Spec: akov2.AtlasNetworkPeeringSpec{
							AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
								Provider: string(provider.ProviderAWS),
								AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
									AccepterRegionName:  "eu-west-2",
									AWSAccountID:        os.Getenv("AWS_ACCOUNT_ID"),
									RouteTableCIDRBlock: "10.0.0.0/24",
								},
							},
						},
					},
				},
			},
		),
		Entry("Test[networkpeering-aws-3]: AWS Network Peering CRs between different regions and without container region specified",
			Label("focus-network-peering-cr-aws-3"),
			func(ctx SpecContext) *model.TestDataProvider {
				return model.DataProvider(ctx, "network-peering-cr-aws-3", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			[]containerAndPeering{
				{
					container: &akov2.AtlasNetworkContainer{
						Spec: akov2.AtlasNetworkContainerSpec{
							Provider: string(provider.ProviderAWS),
							AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
								Region:    "US_EAST_1",
								CIDRBlock: "10.64.0.0/22",
							},
						},
					},
					peering: &akov2.AtlasNetworkPeering{
						Spec: akov2.AtlasNetworkPeeringSpec{
							AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
								Provider: string(provider.ProviderAWS),
								AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
									AccepterRegionName:  "eu-west-1",
									AWSAccountID:        os.Getenv("AWS_ACCOUNT_ID"),
									RouteTableCIDRBlock: "192.168.0.0/16",
								},
							},
						},
					},
				},
				{
					container: &akov2.AtlasNetworkContainer{
						Spec: akov2.AtlasNetworkContainerSpec{
							Provider: string(provider.ProviderAWS),
							AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
								CIDRBlock: "10.128.0.0/22",
								Region:    "US_WEST_1",
							},
						},
					},
					peering: &akov2.AtlasNetworkPeering{
						Spec: akov2.AtlasNetworkPeeringSpec{
							AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
								Provider: string(provider.ProviderAWS),
								AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
									AccepterRegionName:  "us-east-1",
									AWSAccountID:        os.Getenv("AWS_ACCOUNT_ID"),
									RouteTableCIDRBlock: "10.0.0.0/24",
								},
							},
						},
					},
				},
			},
		),
		Entry("Test[networkpeering-gcp-1]: GCP Network Peering CR",
			Label("focus-network-peering-cr-gcp-1"),
			func(ctx SpecContext) *model.TestDataProvider {
				return model.DataProvider(ctx, "network-peering-cr-gcp-1", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			[]containerAndPeering{
				{
					container: &akov2.AtlasNetworkContainer{
						Spec: akov2.AtlasNetworkContainerSpec{
							Provider: string(provider.ProviderGCP),
							AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
								CIDRBlock: "10.8.0.0/18",
							},
						},
					},
					peering: &akov2.AtlasNetworkPeering{
						Spec: akov2.AtlasNetworkPeeringSpec{
							AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
								Provider: string(provider.ProviderGCP),
								GCPConfiguration: &akov2.GCPNetworkPeeringConfiguration{
									GCPProjectID: cloud.GoogleProjectID,
									NetworkName:  newRandomName(GCPVPCName),
								},
							},
						},
					},
				},
			},
		),
		Entry("Test[networkpeering-azure-1]: Azure Network Peering CR",
			Label("focus-network-peering-cr-azure-1"),
			func(ctx SpecContext) *model.TestDataProvider {
				return model.DataProvider(ctx, "network-peering-cr-azure-1", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			[]containerAndPeering{
				{
					container: &akov2.AtlasNetworkContainer{
						Spec: akov2.AtlasNetworkContainerSpec{
							Provider: string(provider.ProviderAzure),
							AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
								CIDRBlock: "192.168.248.0/21",
								Region:    "US_EAST_2",
							},
						},
					},
					peering: &akov2.AtlasNetworkPeering{
						Spec: akov2.AtlasNetworkPeeringSpec{
							AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
								Provider: string(provider.ProviderAzure),
								AzureConfiguration: &akov2.AzureNetworkPeeringConfiguration{
									AzureDirectoryID:    os.Getenv(DirectoryID),
									AzureSubscriptionID: os.Getenv(SubscriptionID),
									ResourceGroupName:   cloud.ResourceGroupName,
									VNetName:            newRandomName(AzureVPCName),
								},
							},
						},
					},
				},
			},
		),
	)
})

func networkPeerControllerFlow(ctx context.Context, userData *model.TestDataProvider, pairs []containerAndPeering) {
	providerActions := make([]cloud.Provider, len(pairs))

	By("Prepare network peers cloud infrastructure from CRs", func() {
		for i, pair := range pairs {
			peer := pair.peering
			providerAction, err := prepareProviderAction(ctx)
			Expect(err).To(BeNil())
			providerActions[i] = providerAction

			providerName := provider.ProviderName(peer.Spec.Provider)
			switch providerName {
			case provider.ProviderAWS:
				peer.Spec.AWSConfiguration.AWSAccountID = providerActions[i].GetAWSAccountID(ctx)
				cfg := &cloud.AWSConfig{
					Region:        peer.Spec.AWSConfiguration.AccepterRegionName,
					VPC:           newRandomName("ao-vpc-peering-e2e"),
					CIDR:          peer.Spec.AWSConfiguration.RouteTableCIDRBlock,
					Subnets:       map[string]string{"ao-peering-e2e-subnet": peer.Spec.AWSConfiguration.RouteTableCIDRBlock},
					EnableCleanup: true,
				}
				peer.Spec.AWSConfiguration.VpcID = providerActions[i].SetupNetwork(ctx, providerName, cloud.WithAWSConfig(cfg))
			case provider.ProviderGCP:
				cfg := &cloud.GCPConfig{
					VPC:           peer.Spec.GCPConfiguration.NetworkName,
					EnableCleanup: true,
				}
				providerActions[i].SetupNetwork(ctx, providerName, cloud.WithGCPConfig(cfg))
			case provider.ProviderAzure:
				cfg := &cloud.AzureConfig{
					VPC:           peer.Spec.AzureConfiguration.VNetName,
					EnableCleanup: true,
				}
				providerActions[i].SetupNetwork(ctx, providerName, cloud.WithAzureConfig(cfg))
			}
		}
	})

	By("Create network containers from CRs and update their IDs", func() {
		for i, pair := range pairs {
			container := pair.container
			container.Spec.ProjectRef = &common.ResourceRefNamespaced{
				Name:      userData.Project.Name,
				Namespace: userData.Project.Namespace,
			}
			container.Name = fmt.Sprintf("container-%s-item-%d", userData.Prefix, i)
			container.Namespace = userData.Project.Namespace
			Expect(userData.K8SClient.Create(userData.Context, container)).Should(Succeed())
		}
		for _, pair := range pairs {
			key := client.ObjectKeyFromObject(pair.container)
			Eventually(func(g Gomega) bool {
				Expect(userData.K8SClient.Get(userData.Context, key, pair.container)).Should(Succeed())
				return pair.container.Status.ID != ""
			}).WithTimeout(3*time.Minute).WithPolling(20*time.Second).Should(
				BeTrue(),
				"Network Containers CRs should be created with an Atlas ID set in the status",
			)
		}
	})

	By("Create network peer from CRs", func() {
		for i, pair := range pairs {
			peer := pair.peering
			peer.Spec.ProjectRef = &common.ResourceRefNamespaced{
				Name:      userData.Project.Name,
				Namespace: userData.Project.Namespace,
			}
			peer.Name = fmt.Sprintf("%s-item-%d", userData.Prefix, i)
			peer.Namespace = userData.Project.Namespace
			peer.Spec.ContainerRef.ID = pair.container.Status.ID
			Expect(userData.K8SClient.Create(userData.Context, peer)).Should(Succeed())
		}
	})

	By("Establish network peers connection with CRs", func() {
		Eventually(func(g Gomega) bool {
			return EnsurePeersReadyToConnect(g, userData, pairs)
		}).WithTimeout(15*time.Minute).WithPolling(20*time.Second).Should(
			BeTrue(),
			"Network Peering CRs should be ready to establish connection",
		)

		for ix, pair := range pairs {
			peer := pair.peering
			providerName := provider.ProviderName(peer.Spec.Provider)
			switch providerName {
			case provider.ProviderAWS:
				providerActions[ix].SetupNetworkPeering(ctx, providerName, peer.Status.AWSStatus.ConnectionID, "")
			case provider.ProviderGCP:
				providerActions[ix].SetupNetworkPeering(ctx, providerName, pair.peering.Status.GCPStatus.GCPProjectID, pair.peering.Status.GCPStatus.NetworkName)
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

	By("Check network containers & peers CRs to be Ready", func() {
		for _, pair := range pairs {
			containerKey := client.ObjectKeyFromObject(pair.container)
			Expect(userData.K8SClient.Get(userData.Context, containerKey, pair.container)).Should(Succeed())
			Expect(networkContainerReady(pair.container)).Should(BeTrue())

			key := client.ObjectKeyFromObject(pair.peering)
			Expect(userData.K8SClient.Get(userData.Context, key, pair.peering)).Should(Succeed())
			Expect(networkPeeringReady(pair.peering)).Should(BeTrue())
		}
	})
}

func EnsurePeersReadyToConnect(g Gomega, userData *model.TestDataProvider, pairs []containerAndPeering) bool {
	for _, pair := range pairs {
		key := client.ObjectKeyFromObject(pair.peering)
		g.Expect(userData.K8SClient.Get(userData.Context, key, pair.peering)).Should(Succeed())
		if pair.peering.Spec.Provider == string(provider.ProviderAzure) {
			continue
		}
		statusMsg := pair.peering.Status.Status
		if statusMsg != statusPendingAcceptance && statusMsg != statusWaitingUser {
			return false
		}
		if pair.peering.Spec.Provider == string(provider.ProviderAWS) &&
			pair.peering.Status.AWSStatus == nil {
			return false
		}
	}
	By("Network containers & peers are ready to connect", func() {})
	return true
}

func networkPeeringReady(peer *akov2.AtlasNetworkPeering) bool {
	for _, condition := range peer.Status.Conditions {
		if condition.Type == api.ReadyType && condition.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}
