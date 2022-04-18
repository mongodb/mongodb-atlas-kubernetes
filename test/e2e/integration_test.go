package e2e_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/deploy"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

var _ =  Describe("Configuration namespaced. Deploy cluster", Label("integration-ns"), func() {
	var data model.TestDataProvider
	var key string

	BeforeEach(func() {
		Eventually(kubecli.GetVersionOutput()).Should(Say(K8sVersion))
		key = os.Getenv("DATADOG_KEY")
		Expect(key).ShouldNot(BeEmpty())
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
				[]string{"deploy", "atlasclusters", "atlasdatabaseusers", "atlasprojects"},
				data.Resources.Namespace,
			)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{data})
		}
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test model.TestDataProvider) {
			data = test
			integrationCycle(test, key)
		},
		Entry("Users can use integration section", Label("project-integration"),
			model.NewTestDataProvider(
				"operator-ns-project-integration-cr",
				*model.NewProject("project-integration-cr").
					ProjectName("project-integration").
					WithIntegration(
						*model.NewPIntegration("DATADOG").WithAPIKeyRef("test-int", ""),
					),
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{},
				[]string{},
				[]model.DBUser{},
				30018,
				[]func(*model.TestDataProvider){
				},
			),
		),
	)
})

func integrationCycle(data model.TestDataProvider, key string) {
	actions.PrepareUsersConfigurations(&data)
	deploy.NamespacedOperator(&data)

	By("Create Secret for integration", func() {
		for _, i := range data.Resources.Project.Spec.Integrations {
			kubecli.CreateUserSecret(key, i.APIKeyRef.Name, i.APIKeyRef.Namespace)
		}
	})

	By("Deploy User Resouces", func() {
		actions.DeployProjectAndWait(&data, "1")
		Expect(data.Resources.ProjectID).ShouldNot(BeEmpty())
	})

	By("Check statuses", func() {
		// TODO to kube
		status := kubecli.GetStatusCondition(string(status.IntegrationReadyType), data.Resources.Namespace, data.Resources.GetAtlasProjectFullKubeName())
		Expect(status).Should(Equal("True"))
	})

	By("Additional check for the current data set", func() {
		for _, check := range data.Actions {
			check(&data)
		}
	})
	By("Delete User Resources", func() {
		actions.DeleteUserResourcesProject(&data)
	})
}
