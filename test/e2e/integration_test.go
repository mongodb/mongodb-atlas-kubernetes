package e2e_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"

	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

var _ = Describe("Configuration namespaced. Deploy deployment", Label("integration-ns"), func() {
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
			actions.SaveProjectsToFile(data.Context, data.K8SClient, data.Resources.Namespace)
			actions.SaveK8sResources(
				[]string{"deploy"},
				data.Resources.Namespace,
			)
			actions.DeleteUserResourcesProject(&data)
			actions.DeleteGlobalKeyIfExist(data)
		}
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test model.TestDataProvider) {
			data = test
			integrationCycle(test, key)
		},
		Entry("Users can use integration section", Label("project-integration"),
			model.NewTestDataProvider(
				"operator-integration-cr",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{},
				[]string{},
				[]model.DBUser{},
				30018,
				[]func(*model.TestDataProvider){},
			),
		),
	)
})

func integrationCycle(data model.TestDataProvider, key string) {
	actions.PrepareUsersConfigurations(&data)
	deploy.NamespacedOperator(&data)
	t := "DATADOG"

	By("Deploy User Resouces", func() {
		actions.DeployProjectAndWait(&data, 1)
		Expect(data.Resources.ProjectID).ShouldNot(BeEmpty())
		projectStatus, _ := kubecli.GetProjectStatusCondition(data.Context, data.K8SClient, status.IntegrationReadyType, data.Resources.Namespace, data.Resources.Project.ObjectMeta.GetName())
		Expect(projectStatus).Should(BeEmpty())
	})

	By("Add integration", func() {
		newIntegration := model.NewPIntegration(t).WithAPIKeyRef("test-int", data.Resources.Namespace).WithRegion("EU")
		data.Resources.Project = data.Resources.Project.WithIntegration(*newIntegration)
		By("Create Secret for integration", func() {
			for _, i := range data.Resources.Project.Spec.Integrations {
				kubecli.CreateUserSecret(key, i.APIKeyRef.Name, i.APIKeyRef.Namespace)
			}
		})
		actions.PrepareUsersConfigurations(&data)
		actions.DeployProjectAndWait(&data, 2)
	})
	atlasClient, err := atlas.AClient()
	By("Check statuses", func() {
		var projectStatus string
		projectStatus, err = kubecli.GetProjectStatusCondition(data.Context, data.K8SClient, status.IntegrationReadyType, data.Resources.Namespace, data.Resources.Project.ObjectMeta.GetName())
		Expect(err).ShouldNot(HaveOccurred())
		Expect(projectStatus).Should(Equal("True"))

		Expect(err).ShouldNot(HaveOccurred())

		dog, err := atlasClient.GetIntegrationbyType(data.Resources.ProjectID, t)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(dog.APIKey).Should(Equal(key))
	})

	By("Delete integration", func() {
		data.Resources.Project.Spec.Integrations = []project.Integration{}
		actions.PrepareUsersConfigurations(&data)
		actions.DeployProjectAndWait(&data, 3)
	})

	By("Delete integration check", func() {
		_, err := atlasClient.GetIntegrationbyType(data.Resources.ProjectID, t)
		Expect(err).Should(HaveOccurred())

		// TODO uncomment with
		// status := kubecli.GetStatusCondition(string(status.IntegrationReadyType), data.Resources.Namespace, data.Resources.GetAtlasProjectFullKubeName())
		// Expect(status).Should(BeEmpty())
	})

	By("Delete User Resources", func() {
		actions.DeleteUserResourcesProject(&data)
	})
}
