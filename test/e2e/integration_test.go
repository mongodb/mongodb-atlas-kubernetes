package e2e_test

import (
	"fmt"
	"os"
	"time"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/kube"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/k8s"

	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/data"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/atlas"

	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

var _ = Describe("Configuration namespaced. Deploy deployment", Label("integration-ns"), func() {
	var testData *model.TestDataProvider
	var key string

	BeforeEach(func() {
		Eventually(kubecli.GetVersionOutput()).Should(Say(K8sVersion))
		key = os.Getenv("DATADOG_KEY")
		Expect(key).ShouldNot(BeEmpty())
	})

	AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			SaveDump(testData)
		}
		actions.DeleteTestDataProject(testData)
		actions.DeleteGlobalKeyIfExist(*testData)
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test *model.TestDataProvider) {
			testData = test
			actions.ProjectCreationFlow(test)
			integrationCycle(test, key)
		},
		Entry("Users can use integration section", Label("project-integration"),
			model.DataProvider(
				"operator-integration-cr",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				30018,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
		),
	)
})

func integrationCycle(data *model.TestDataProvider, key string) {
	t := "DATADOG"

	By("Deploy User Resouces", func() {
		projectStatus := GetProjectIntegrationStatus(data)
		Expect(projectStatus).Should(BeEmpty())
	})

	By("Add integration", func() {
		Expect(data.K8SClient.Get(data.Context, types.NamespacedName{Name: data.Project.Name,
			Namespace: data.Resources.Namespace}, data.Project)).Should(Succeed())
		newIntegration := project.Integration{
			Type: t,
			APIKeyRef: common.ResourceRefNamespaced{
				Name:      "test-int",
				Namespace: data.Resources.Namespace,
			},
			Region: "EU",
		}
		data.Project.Spec.Integrations = append(data.Project.Spec.Integrations, newIntegration)
		By("Create Secret for integration", func() {
			for _, i := range data.Project.Spec.Integrations {
				Expect(k8s.CreateUserSecret(data.Context, data.K8SClient, key, i.APIKeyRef.Name, i.APIKeyRef.Namespace)).Should(Succeed())
			}
		})
		Expect(data.K8SClient.Update(data.Context, data.Project)).Should(Succeed())
		Eventually(func() string {
			return GetProjectIntegrationStatus(data)
		}).WithTimeout(5 * time.Minute).WithPolling(20 * time.Second).Should(Equal("True"))
		Eventually(func(g Gomega) string {
			condition, err := kube.GetProjectStatusCondition(data, status.ReadyType)
			g.Expect(err).ShouldNot(HaveOccurred())
			return condition
		}).WithTimeout(5 * time.Minute).WithPolling(20 * time.Second).Should(Equal("True"))
	})
	atlasClient := atlas.GetClientOrFail()
	By("Check statuses", func() {
		var projectStatus string
		projectStatus, err := k8s.GetProjectStatusCondition(data.Context, data.K8SClient, status.IntegrationReadyType, data.Resources.Namespace, data.Project.GetName())
		Expect(err).ShouldNot(HaveOccurred())
		Expect(projectStatus).Should(Equal("True"))

		Expect(err).ShouldNot(HaveOccurred())

		dog, err := atlasClient.GetIntegrationbyType(data.Project.ID(), t)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(dog.APIKey).Should(Equal(key))
	})

	By("Delete integration", func() {
		Expect(data.K8SClient.Get(data.Context, types.NamespacedName{Name: data.Project.Name,
			Namespace: data.Resources.Namespace}, data.Project)).Should(Succeed())
		data.Project.Spec.Integrations = []project.Integration{}
		Expect(data.K8SClient.Update(data.Context, data.Project)).Should(Succeed())
		Eventually(func() string {
			return GetProjectIntegrationStatus(data)
		}).WithTimeout(5 * time.Minute).WithPolling(20 * time.Second).Should(BeEmpty())
		Eventually(func(g Gomega) string {
			condition, err := kube.GetProjectStatusCondition(data, status.ReadyType)
			g.Expect(err).ShouldNot(HaveOccurred())
			return condition
		}).WithTimeout(5 * time.Minute).WithPolling(20 * time.Second).Should(Equal("True"))
	})

	By("Delete integration check", func() {
		integration, err := atlasClient.GetIntegrationbyType(data.Project.ID(), t)
		By(fmt.Sprintf("Integration %v", integration))
		Expect(err).Should(HaveOccurred())

		// TODO uncomment with
		// status := kubecli.GetStatusCondition(string(status.IntegrationReadyType), data.Resources.Namespace, data.Resources.GetAtlasProjectFullKubeName())
		// Expect(status).Should(BeEmpty())
	})
}

func GetProjectIntegrationStatus(testData *model.TestDataProvider) string {
	Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{Name: testData.Project.Name, Namespace: testData.Project.Namespace}, testData.Project)).Should(Succeed())
	for _, condition := range testData.Project.Status.Conditions {
		if condition.Type == status.IntegrationReadyType {
			return string(condition.Status)
		}
	}
	return ""
}
