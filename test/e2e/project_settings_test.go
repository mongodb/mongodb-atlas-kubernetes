package e2e_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

var _ = Describe("UserLogin", Label("project-settings"), func() {
	var testData *model.TestDataProvider

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Project Settings Test\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Delete Resources", func() {
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test *model.TestDataProvider, settings akov2.ProjectSettings) {
			testData = test
			actions.ProjectCreationFlow(test)
			projectSettingsFlow(test, &settings)
		},
		Entry("Test[project-settings]: User has project to which Project Settings was added", Label("project-settings"),
			model.DataProvider(
				"project-settings",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			akov2.ProjectSettings{
				IsCollectDatabaseSpecificsStatisticsEnabled: pointer.MakePtr(false),
				IsDataExplorerEnabled:                       pointer.MakePtr(false),
				IsPerformanceAdvisorEnabled:                 pointer.MakePtr(false),
				IsRealtimePerformancePanelEnabled:           pointer.MakePtr(false),
				IsSchemaAdvisorEnabled:                      pointer.MakePtr(false),
			},
		),
	)
})

func projectSettingsFlow(userData *model.TestDataProvider, settings *akov2.ProjectSettings) {
	By("Add Project Settings to the project", func() {
		userData.Project.Spec.Settings = settings
		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		actions.WaitForConditionsToBecomeTrue(userData, status.ProjectSettingsReadyType, status.ReadyType)
	})

	By("Remove Project Settings from the project", func() {
		userData.Project.Spec.Settings = nil
		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		actions.CheckProjectConditionsNotSet(userData, status.ProjectSettingsReadyType)
	})
}
