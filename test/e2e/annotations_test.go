package e2e_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

var _ = Describe("Annotations base test.", Label("deployment-annotations-ns"), func() {
	var testData *model.TestDataProvider

	AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
			Expect(actions.SaveDeploymentsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		actions.DeleteTestDataDeployments(testData)
		actions.DeleteTestDataProject(testData)
		actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test *model.TestDataProvider) {
			testData = test
			mainCycle(test)
		},
		// TODO: fix test for deletion protection on, as it would fail to re-take the cluster after deletion
		Entry("Simple configuration with keep resource policy annotation on deployment", Label("ns-crd"),
			model.DataProvider(
				"operator-ns-crd",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				30000,
				[]func(*model.TestDataProvider){
					actions.DeleteDeploymentCRWithKeepAnnotation,
					actions.RedeployDeployment,
					actions.RemoveKeepAnnotation,
				},
			).WithInitialDeployments(data.CreateDeploymentWithKeepPolicy("atlascluster-annotation")).
				WithProject(data.DefaultProject()),
		),
	)
})
