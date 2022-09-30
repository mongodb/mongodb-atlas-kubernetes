package e2e_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

var _ = Describe("Annotations base test.", Label("deployment-annotations-ns"), func() {
	var testData *model.TestDataProvider

	BeforeEach(func() {
		Eventually(kubecli.GetVersionOutput()).Should(Say(K8sVersion))
	})
	AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			GinkgoWriter.Write([]byte("Test has been failed. Trying to save logs...\n"))
			utils.SaveToFile(
				fmt.Sprintf("output/%s/operatorDecribe.txt", testData.Resources.Namespace),
				[]byte(kubecli.DescribeOperatorPod(testData.Resources.Namespace)),
			)
			utils.SaveToFile(
				fmt.Sprintf("output/%s/operator-logs.txt", testData.Resources.Namespace),
				kubecli.GetManagerLogs(testData.Resources.Namespace),
			)
			actions.SaveTestAppLogs(testData.Resources)
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
			actions.SaveK8sResources(
				[]string{"deploy", "atlasdeployments"},
				testData.Resources.Namespace,
			)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
			actions.DeleteTestDataDeployments(testData)
			actions.DeleteTestDataProject(testData)
		}
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test *model.TestDataProvider) {
			testData = test
			mainCycle(test)
		},
		Entry("Simple configuration with keep resource policy annotation on deployment", Label("ns-crd"),
			model.DataProvider(
				"operator-ns-crd",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
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
