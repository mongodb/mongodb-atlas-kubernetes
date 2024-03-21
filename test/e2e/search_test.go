package e2e_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

var _ = Describe("Search Nodes", Label("atlas-search"), func() {
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

	It("Creates, upgrades, and deletes search nodes", func() {
		By("Setting up project", func() {
			actions.ProjectCreationFlow(testData.WithProject(data.DefaultProject()).WithInitialDeployments(data.CreateAdvancedDeployment("search-simple-s20")))
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
					if condition.Type == status.DeploymentReadyType {
						return condition.Status == v1.ConditionTrue
					}
				}
				return false
			}).WithTimeout(20 * time.Minute).Should(BeTrue())
		})
		By("Upgrading the deployment with different search nodes", func() {
			testData.InitialDeployments[0].Spec.DeploymentSpec.SearchNodes[0].InstanceSize = "S30_HIGHCPU_NVME"
			Expect(testData.K8SClient.Update(testData.Context, testData.InitialDeployments[0])).To(Succeed())

			Eventually(func(g Gomega) bool {
				g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{
					Name:      testData.InitialDeployments[0].Name,
					Namespace: testData.InitialDeployments[0].Namespace,
				}, testData.InitialDeployments[0])).To(Succeed())
				for _, condition := range testData.InitialDeployments[0].Status.Conditions {
					if condition.Type == status.DeploymentReadyType {
						return condition.Status == v1.ConditionTrue
					}
				}
				return false
			}).WithTimeout(20 * time.Minute).Should(BeTrue())

		})
		By("Removing the search nodes from the deployment", func() {
			testData.InitialDeployments[0].Spec.DeploymentSpec.SearchNodes = nil
			Expect(testData.K8SClient.Update(testData.Context, testData.InitialDeployments[0])).To(Succeed())

			Eventually(func(g Gomega) bool {
				g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{
					Name:      testData.InitialDeployments[0].Name,
					Namespace: testData.InitialDeployments[0].Namespace,
				}, testData.InitialDeployments[0])).To(Succeed())
				for _, condition := range testData.InitialDeployments[0].Status.Conditions {
					if condition.Type == status.DeploymentReadyType {
						return condition.Status == v1.ConditionTrue
					}
				}
				return false
			}).WithTimeout(20 * time.Minute).Should(BeTrue())

		})
	})

})
