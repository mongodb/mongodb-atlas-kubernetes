package e2e_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/deploy"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

var _ = Describe("Deploy cluster", Label("cluster-ns-ct"), func() {
	var data model.TestDataProvider

	BeforeEach(func() {
		Eventually(kubecli.GetVersionOutput()).Should(Say(K8sVersion))
	})
	AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + data.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
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
				[]string{"deploy", "atlasclusters", "atlasprojects"},
				data.Resources.Namespace,
			)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{data})
		}
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test model.TestDataProvider) {
			data = test
			advancedMainCycle(test)
		},
		Entry("Trial - Simplest configuration with no backup and no user", Label("ns-trial-ct"),
			model.NewTestDataProvider(
				"operator-ns-trial",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlascluster_basic_with_keep_resource_policy.yaml"},
				[]string{},
				[]model.DBUser{},
				30000,
				[]func(*model.TestDataProvider){
					actions.DeleteCRDs,
				},
			),
		),
	)
})

func advancedMainCycle(data model.TestDataProvider) {
	actions.PrepareUsersConfigurations(&data)
	deploy.NamespacedOperator(&data)

	By("Deploy Advanced Resources", func() {
		actions.DeployAdvancedResourcesAction(&data)
		Expect(data.Resources.ProjectID).ShouldNot(BeEmpty())
	})

	By("Additional check for the current data set", func() {
		for _, check := range data.Actions {
			check(&data)
		}
	})
	By("Delete Advanced Resources", func() {
		actions.DeleteAdvancedResources(&data)
	})
}
