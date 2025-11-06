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
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkcontainer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/api/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

const (
	createMeID = "create-me"
)

var _ = Describe("NetworkContainerController", Label("networkcontainer-controller"), func() {
	var testData *model.TestDataProvider

	_ = BeforeEach(OncePerOrdered, func() {
		checkUpAWSEnvironment()
		checkUpAzureEnvironment()
		checkNSetUpGCPEnvironment()
	})

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Network Container Controller Test\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Delete Resources, Project with NetworkContainer", func() {
			actions.DeleteTestDataNetworkContainers(testData)
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	DescribeTable("NetworkContainerController",
		func(ctx SpecContext, test func(context.Context) *model.TestDataProvider, useProjectID bool, networkPeers []*akov2.AtlasNetworkContainer) {
			testData = test(ctx)
			actions.ProjectCreationFlow(testData)
			networkContainerControllerFlow(testData, useProjectID, networkPeers)
		},
		Entry("Test[networkpeering-aws-1]: New AWS Network Container is created successfully",
			Label("focus-network-container-cr-aws-1"),
			func(ctx SpecContext) *model.TestDataProvider {
				return model.DataProvider(ctx, "networkcontainer-cr-aws-1", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			false,
			[]*akov2.AtlasNetworkContainer{
				{
					Spec: akov2.AtlasNetworkContainerSpec{
						ProjectDualReference: akov2.ProjectDualReference{
							ProjectRef: &common.ResourceRefNamespaced{Name: data.ProjectName},
						},
						Provider: string(provider.ProviderAWS),
						AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
							Region:    "US_EAST_1",
							CIDRBlock: "10.128.0.0/21",
						},
					},
				},
			},
		),
		Entry("Test[networkpeering-azure-2]: New Azure Network Container is created successfully",
			Label("focus-network-container-cr-azure-2"),
			func(ctx SpecContext) *model.TestDataProvider {
				return model.DataProvider(ctx, "networkcontainer-cr-azure-2", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			false,
			[]*akov2.AtlasNetworkContainer{
				{
					Spec: akov2.AtlasNetworkContainerSpec{
						ProjectDualReference: akov2.ProjectDualReference{
							ProjectRef: &common.ResourceRefNamespaced{Name: data.ProjectName},
						},
						Provider: string(provider.ProviderAzure),
						AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
							Region:    "US_EAST_2",
							CIDRBlock: "10.128.0.0/21",
						},
					},
				},
			},
		),
		Entry("Test[networkpeering-gcp-3]: New GCP Network Container is created successfully",
			Label("focus-network-container-cr-gcp-3"),
			func(ctx SpecContext) *model.TestDataProvider {
				return model.DataProvider(ctx, "networkcontainer-cr-gcp-3", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			false,
			[]*akov2.AtlasNetworkContainer{
				{
					Spec: akov2.AtlasNetworkContainerSpec{
						ProjectDualReference: akov2.ProjectDualReference{
							ProjectRef: &common.ResourceRefNamespaced{Name: data.ProjectName},
						},
						Provider: string(provider.ProviderGCP),
						AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
							CIDRBlock: "10.128.0.0/18",
						},
					},
				},
			},
		),
		Entry("Test[networkpeering-all-5]: Existing Network Containers from all providers with direct ids are taken over successfully",
			Label("focus-network-container-cr-all-5"),
			func(ctx SpecContext) *model.TestDataProvider {
				return model.DataProvider(ctx, "networkcontainer-cr-all-5", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			true,
			[]*akov2.AtlasNetworkContainer{
				{
					Spec: akov2.AtlasNetworkContainerSpec{
						ProjectDualReference: akov2.ProjectDualReference{
							ProjectRef: &common.ResourceRefNamespaced{Name: data.ProjectName},
						},
						Provider: string(provider.ProviderAWS),
						AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
							ID:        createMeID,
							Region:    "US_EAST_1",
							CIDRBlock: "10.128.0.0/21",
						},
					},
				},
				{
					Spec: akov2.AtlasNetworkContainerSpec{
						ProjectDualReference: akov2.ProjectDualReference{
							ProjectRef: &common.ResourceRefNamespaced{Name: data.ProjectName},
						},
						Provider: string(provider.ProviderAzure),
						AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
							ID:        createMeID,
							Region:    "US_EAST_2",
							CIDRBlock: "10.128.0.0/21",
						},
					},
				},
				{
					Spec: akov2.AtlasNetworkContainerSpec{
						ProjectDualReference: akov2.ProjectDualReference{
							ProjectRef: &common.ResourceRefNamespaced{Name: data.ProjectName},
						},
						Provider: string(provider.ProviderGCP),
						AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
							ID:        createMeID,
							CIDRBlock: "10.128.0.0/18",
						},
					},
				},
			},
		),
	)
})

func networkContainerControllerFlow(userData *model.TestDataProvider, useProjectID bool, containers []*akov2.AtlasNetworkContainer) {
	By("Create network containers from CRs", func() {
		atlasClient, err := atlas.AClient()
		Expect(err).To(Succeed())
		projectID := ""
		if useProjectID {
			createdProject := akov2.AtlasProject{}
			key := types.NamespacedName{Name: userData.Project.Name, Namespace: userData.Project.Namespace}
			Expect(userData.K8SClient.Get(userData.Context, key, &createdProject)).Should(Succeed())
			projectID = createdProject.Status.ID
		}
		for i, container := range containers {
			if useProjectID {
				container.Spec.ExternalProjectRef = &akov2.ExternalProjectReference{
					ID: projectID,
				}
				container.Spec.ProjectRef = nil
				container.Spec.ConnectionSecret = &api.LocalObjectReference{
					Name: config.DefaultOperatorGlobalKey,
				}
			} else {
				container.Spec.ProjectRef = &common.ResourceRefNamespaced{
					Name:      userData.Project.Name,
					Namespace: userData.Project.Namespace,
				}
				container.Spec.ExternalProjectRef = nil
			}
			container.Name = fmt.Sprintf("%s-item-%d", userData.Prefix, i)
			container.Namespace = userData.Project.Namespace
			if container.Spec.ID == createMeID {
				id, err := createTestContainer(userData.Context, atlasClient, userData.Project.Status.ID, &container.Spec)
				Expect(err).To(Succeed())
				container.Spec.ID = id
			}
			Expect(userData.K8SClient.Create(userData.Context, container)).Should(Succeed())
		}
	})

	By("Check network container CRs to be Ready", func() {
		for _, container := range containers {
			key := types.NamespacedName{Name: container.Name, Namespace: container.Namespace}
			Eventually(func(g Gomega) bool {
				Expect(userData.K8SClient.Get(userData.Context, key, container)).Should(Succeed())
				return networkContainerReady(container)
			}).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(BeTrue())
		}
	})
}

func networkContainerReady(container *akov2.AtlasNetworkContainer) bool {
	for _, condition := range container.Status.Conditions {
		if condition.Type == api.ReadyType && condition.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}

func createTestContainer(ctx context.Context, atlasClient atlas.Atlas, projectID string, container *akov2.AtlasNetworkContainerSpec) (string, error) {
	service := networkcontainer.NewNetworkContainerService(atlasClient.Client.NetworkPeeringApi)
	cfg := networkcontainer.NewNetworkContainerConfig(container.Provider, &container.AtlasNetworkContainerConfig)
	createdContainer, err := service.Create(ctx, projectID, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to pre-providiong test container for %s config %v: %w",
			container.Provider, container.AtlasNetworkContainerConfig, err)
	}
	return createdContainer.ID, nil
}
