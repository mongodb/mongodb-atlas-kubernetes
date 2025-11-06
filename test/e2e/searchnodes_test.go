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
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/api/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

var _ = Describe("Search Nodes", Label("atlas-search-nodes"), func() {
	var testData *model.TestDataProvider

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
			Expect(actions.SaveDeploymentsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Delete Resources", func() {
			actions.DeleteTestDataDeployments(testData)
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	It("Creates, upgrades, and deletes search nodes", func(ctx SpecContext) {
		testData = model.DataProvider(ctx, "atlas-search-nodes", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject()).WithInitialDeployments(data.CreateAdvancedDeployment("search-nodes-test"))
		atlasClient = atlas.GetClientOrFail()
		By("Setting up project", func() {
			actions.ProjectCreationFlow(testData)
		})
		By("Creating a deployment with search nodes", func() {
			search := []akov2.SearchNode{
				{
					InstanceSize: "S20_HIGHCPU_NVME",
					NodeCount:    2,
				},
			}
			testData.InitialDeployments[0].Spec.DeploymentSpec.SearchNodes = search
			testData.InitialDeployments[0].Namespace = testData.Resources.Namespace
			Expect(testData.K8SClient.Create(testData.Context, testData.InitialDeployments[0])).To(Succeed())

			Eventually(func(g Gomega) bool {
				g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{
					Name:      testData.InitialDeployments[0].Name,
					Namespace: testData.InitialDeployments[0].Namespace,
				}, testData.InitialDeployments[0])).To(Succeed())
				for _, condition := range testData.InitialDeployments[0].Status.Conditions {
					if condition.Type == api.DeploymentReadyType {
						return condition.Status == v1.ConditionTrue
					}
				}
				return false
			}).WithTimeout(60 * time.Minute).Should(BeTrue())

			Eventually(func(g Gomega) {
				atlasSearchNodes, _, err := atlasClient.Client.AtlasSearchApi.GetAtlasSearchDeployment(testData.Context, testData.Project.ID(), testData.InitialDeployments[0].Name).Execute()
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(atlasSearchNodes.GetSpecs()[0].InstanceSize).Should(Equal("S20_HIGHCPU_NVME"))
				g.Expect(atlasSearchNodes.GetSpecs()[0].NodeCount).Should(Equal(2))
			}).WithPolling(10 * time.Second).WithTimeout(5 * time.Minute).Should(Succeed())
		})
		By("Upgrading the deployment with different search nodes", func() {
			testData.InitialDeployments[0].Spec.DeploymentSpec.SearchNodes[0].InstanceSize = "S30_HIGHCPU_NVME"
			Expect(testData.K8SClient.Update(testData.Context, testData.InitialDeployments[0])).To(Succeed())

			Eventually(func(g Gomega) {
				atlasSearchNodes, _, err := atlasClient.Client.AtlasSearchApi.GetAtlasSearchDeployment(testData.Context, testData.Project.ID(), testData.InitialDeployments[0].Name).Execute()
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(atlasSearchNodes.GetSpecs()[0].InstanceSize).Should(Equal("S30_HIGHCPU_NVME"))
				g.Expect(atlasSearchNodes.GetSpecs()[0].NodeCount).Should(Equal(2))
			}).WithPolling(10 * time.Second).WithTimeout(5 * time.Minute).Should(Succeed())

			Eventually(func(g Gomega) bool {
				g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{
					Name:      testData.InitialDeployments[0].Name,
					Namespace: testData.InitialDeployments[0].Namespace,
				}, testData.InitialDeployments[0])).To(Succeed())
				for _, condition := range testData.InitialDeployments[0].Status.Conditions {
					if condition.Type == api.DeploymentReadyType {
						return condition.Status == v1.ConditionTrue
					}
				}
				return false
			}).WithTimeout(60 * time.Minute).Should(BeTrue())

		})
		By("Removing the search nodes from the deployment", func() {
			testData.InitialDeployments[0].Spec.DeploymentSpec.SearchNodes = nil
			Expect(testData.K8SClient.Update(testData.Context, testData.InitialDeployments[0])).To(Succeed())

			Eventually(func(g Gomega) {
				response, httpResponse, _ := atlasClient.Client.AtlasSearchApi.GetAtlasSearchDeployment(testData.Context, testData.Project.ID(), testData.InitialDeployments[0].Name).Execute()
				g.Expect(httpResponse).NotTo(BeNil())
				g.Expect(response).NotTo(BeNil())
				g.Expect(len(response.GetSpecs())).To(Equal(0))
			}).WithPolling(10 * time.Second).WithTimeout(15 * time.Minute).Should(Succeed())

			Eventually(func(g Gomega) bool {
				g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{
					Name:      testData.InitialDeployments[0].Name,
					Namespace: testData.InitialDeployments[0].Namespace,
				}, testData.InitialDeployments[0])).To(Succeed())
				for _, condition := range testData.InitialDeployments[0].Status.Conditions {
					if condition.Type == api.DeploymentReadyType {
						return condition.Status == v1.ConditionTrue
					}
				}
				return false
			}).WithTimeout(20 * time.Minute).Should(BeTrue())

		})
	})
})
