package e2e_test

import (
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

const (
	datadogEnvKey         = "DATADOG_KEY"
	pagerDutyEnvKey       = "PAGER_DUTY_SERVICE_KEY"
	integrationSecretName = "integration-secret"
)

var _ = Describe("Project Third-Party Integration", Label("integration-ns"), func() {
	var testData *model.TestDataProvider

	BeforeEach(func() {
		Expect(os.Getenv(datadogEnvKey)).ShouldNot(BeEmpty())
		Expect(os.Getenv(pagerDutyEnvKey)).ShouldNot(BeEmpty())
	})

	AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		actions.DeleteTestDataProject(testData)
		actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
	})

	DescribeTable("Integration can be configured in a project",
		func(test *model.TestDataProvider, integration project.Integration, envKeyName string, setSecret configSecret) {
			testData = test
			actions.ProjectCreationFlow(test)
			integrationTest(test, integration, os.Getenv(envKeyName), setSecret)
		},

		Entry("Users can integrate DATADOG on region US1", Label("project-integration"),
			model.DataProvider(
				"datatog-us1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				30018,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			project.Integration{
				Type:   "DATADOG",
				Region: "US",
			},
			datadogEnvKey,
			func(integration *project.Integration, ref common.ResourceRefNamespaced) {
				integration.APIKeyRef = ref
			},
		),
		Entry("Users can integrate DATADOG on region US3", Label("project-integration"),
			model.DataProvider(
				"datatog-us3",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				30018,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			project.Integration{
				Type:   "DATADOG",
				Region: "US3",
			},
			datadogEnvKey,
			func(integration *project.Integration, ref common.ResourceRefNamespaced) {
				integration.APIKeyRef = ref
			},
		),
		Entry("Users can integrate DATADOG on region US5", Label("project-integration"),
			model.DataProvider(
				"datatog-us5",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				30018,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			project.Integration{
				Type:   "DATADOG",
				Region: "US5",
			},
			datadogEnvKey,
			func(integration *project.Integration, ref common.ResourceRefNamespaced) {
				integration.APIKeyRef = ref
			},
		),
		Entry("Users can integrate DATADOG on region EU1", Label("project-integration"),
			model.DataProvider(
				"datatog-eu1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				30018,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			project.Integration{
				Type:   "DATADOG",
				Region: "EU",
			},
			datadogEnvKey,
			func(integration *project.Integration, ref common.ResourceRefNamespaced) {
				integration.APIKeyRef = ref
			},
		),
		Entry("Users can integrate PagerDuty on region US", Label("project-integration"),
			model.DataProvider(
				"pager-duty-us",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				30018,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			project.Integration{
				Type:   "PAGER_DUTY",
				Region: "US",
			},
			pagerDutyEnvKey,
			func(integration *project.Integration, ref common.ResourceRefNamespaced) {
				integration.ServiceKeyRef = ref
			},
		),
	)
})

func integrationTest(data *model.TestDataProvider, integration project.Integration, key string, setSecret configSecret) {
	By("Create Secret for integration", func() {
		Expect(k8s.CreateUserSecret(data.Context, data.K8SClient, key, integrationSecretName, data.Resources.Namespace)).Should(Succeed())

		setSecret(&integration, common.ResourceRefNamespaced{Name: integrationSecretName, Namespace: data.Resources.Namespace})
	})

	By("Add integration", func() {
		Expect(data.K8SClient.Get(data.Context, types.NamespacedName{Name: data.Project.Name, Namespace: data.Resources.Namespace}, data.Project)).Should(Succeed())
		data.Project.Spec.Integrations = append(data.Project.Spec.Integrations, integration)

		Expect(data.K8SClient.Update(data.Context, data.Project)).Should(Succeed())
	})

	By("Integration is ready", func() {
		actions.WaitForConditionsToBecomeTrue(data, status.IntegrationReadyType, status.ReadyType)

		atlasIntegration, err := atlasClient.GetIntegrationByType(data.Project.ID(), integration.Type)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(strings.HasSuffix(key, strings.TrimLeft(atlasIntegration.GetApiKey(), "*"))).Should(BeTrue())
	})

	By("Delete integration", func() {
		Expect(data.K8SClient.Get(data.Context, types.NamespacedName{Name: data.Project.Name, Namespace: data.Resources.Namespace}, data.Project)).Should(Succeed())
		data.Project.Spec.Integrations = []project.Integration{}

		Expect(data.K8SClient.Update(data.Context, data.Project)).Should(Succeed())
	})

	By("Delete integration check", func() {
		actions.CheckProjectConditionsNotSet(data, status.IntegrationReadyType)

		atlasIntegration, err := atlasClient.GetIntegrationByType(data.Project.ID(), integration.Type)
		Expect(err).Should(HaveOccurred())
		Expect(atlasIntegration).To(BeNil())
	})
}

type configSecret func(integration *project.Integration, ref common.ResourceRefNamespaced)
