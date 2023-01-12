package e2e_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

var _ = Describe("UserLogin", Label("auditing"), func() {
	var testData *model.TestDataProvider

	_ = BeforeEach(func() {
		Eventually(kubecli.GetVersionOutput()).Should(Say(K8sVersion))
	})

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
		func(test *model.TestDataProvider, auditing v1.Auditing) {
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
			v1.Auditing{
				AuditAuthorizationSuccess: toptr.MakePtr(false),
				AuditFilter:               exampleFilter(),
				Enabled:                   toptr.MakePtr(true),
			},
		),
	)
})

func auditingFlow(userData *model.TestDataProvider, auditing *v1.Auditing) {
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
