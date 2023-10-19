package e2e_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

var _ = Describe("Free tier", Label("free-tier"), func() {
	var testData *model.TestDataProvider

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Free tier test\n"))
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

	DescribeTable("Operator should support exported CR for free tier deployments",
		func(test *model.TestDataProvider) {
			testData = test
			actions.ProjectCreationFlow(test)
			freeTierDeploymentFlow(test)
		},
		Entry("Test free tier deployment",
			model.DataProvider(
				"free-tier",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()).WithInitialDeployments(data.CreateFreeAdvancedDeployment("free-tier")),
		),
		Entry("Test free tier advanced deployment",
			model.DataProvider(
				"free-tier-advanced",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()).WithInitialDeployments(data.CreateFreeAdvancedDeployment("free-tier")),
		),
	)
})

func freeTierDeploymentFlow(userData *model.TestDataProvider) {
	By("Create free cluster in Atlas", func() {
		aClient := atlas.GetClientOrFail()
		Expect(userData.InitialDeployments).Should(HaveLen(1))
		name := userData.InitialDeployments[0].GetDeploymentName()
		_, _, err := aClient.Client.Clusters.Create(userData.Context, userData.Project.ID(), &mongodbatlas.Cluster{
			Name: name,
			ProviderSettings: &mongodbatlas.ProviderSettings{
				ProviderName:        string(provider.ProviderTenant),
				RegionName:          "US_EAST_1",
				InstanceSizeName:    data.InstanceSizeM0,
				BackingProviderName: string(provider.ProviderAWS),
			},
		})
		Expect(err).ShouldNot(HaveOccurred())
	})

	By("Apply deployment CR", func() {
		deploy.CreateInitialDeployments(userData)
	})
}
