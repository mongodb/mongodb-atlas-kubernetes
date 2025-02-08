package e2e_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/conditions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

var _ = Describe("Serverless", Label("serverless"), func() {
	var testData *model.TestDataProvider

	_ = BeforeEach(OncePerOrdered, func() {
		checkUpAWSEnvironment()
		checkUpAzureEnvironment()
		checkNSetUpGCPEnvironment()
	})

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

	DescribeTable("providers",
		func(test *model.TestDataProvider) {
			testData = test
			actions.ProjectCreationFlow(test)
			serverlessFlow(test)
		},
		Entry("Test[spe-aws-1]: Serverless deployment on AWS", Label("spe-aws-1"),
			model.DataProvider(
				"spe-aws-1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()).WithInitialDeployments(data.CreateServerlessDeployment("spe-test-1", "AWS", "US_EAST_1")),
		),
		Entry("Test[spe-azure-1]: Serverless deployment on Azure", Label("spe-azure-1"),
			model.DataProvider(
				"spe-azure-1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()).WithInitialDeployments(data.CreateServerlessDeployment("spe-test-2", "AZURE", "US_EAST_2")),
		),
	)
})

func serverlessFlow(userData *model.TestDataProvider) {
	By("Apply deployment", func() {
		Expect(userData.InitialDeployments).ShouldNot(BeEmpty())
		userData.InitialDeployments[0].Namespace = userData.Resources.Namespace
		Expect(userData.K8SClient.Create(userData.Context, userData.InitialDeployments[0])).To(Succeed())

		Eventually(func(g Gomega) {
			deployment := userData.InitialDeployments[0]
			g.Expect(userData.K8SClient.Get(userData.Context, client.ObjectKeyFromObject(deployment), deployment)).To(Succeed())
			g.Expect(deployment.Status.Conditions).To(ContainElement(conditions.MatchCondition(api.TrueCondition(api.DeploymentReadyType))))
		}).WithTimeout(time.Minute * 15).WithPolling(time.Second * 15).Should(Succeed())
	})
}
