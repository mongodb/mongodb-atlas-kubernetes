package e2e_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

var _ = Describe("UserLogin", Label("auditing"), func() {
	var testData *model.TestDataProvider

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Auditing Test\n"))
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
		func(test *model.TestDataProvider, auditing akov2.Auditing) {
			testData = test
			actions.ProjectCreationFlow(test)
			auditingFlow(test, &auditing)
		},
		Entry("Test[auditing]: User has project to which Auditing was added", Label("auditing"),
			model.DataProvider(
				"auditing",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			akov2.Auditing{
				AuditAuthorizationSuccess: false,
				AuditFilter:               exampleFilter(),
				Enabled:                   true,
			},
		),
	)
})

func auditingFlow(userData *model.TestDataProvider, auditing *akov2.Auditing) {
	By("Add auditing to the project", func() {
		userData.Project.Spec.Auditing = auditing
		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		actions.WaitForConditionsToBecomeTrue(userData, status.AuditingReadyType, status.ReadyType)
	})

	By("Remove Auditing from the project", func() {
		userData.Project.Spec.Auditing = nil
		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		actions.CheckProjectConditionsNotSet(userData, status.AuditingReadyType)
	})
}

func exampleFilter() string {
	return `{"atype" : "authenticate", "param" : {"user" : "auditReadOnly", "db" : "admin", "mechanism" : "SCRAM-SHA-1"} }`
}
