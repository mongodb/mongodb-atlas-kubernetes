package e2e_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/kube"
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
			SaveDump(testData)
		}
		By("Delete Resources", func() {
			actions.DeleteTestDataProject(testData)
			actions.DeleteGlobalKeyIfExist(*testData)
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
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
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
		userData.Resources.Project.WithAuditing(auditing)
		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
	})

	By("Check project status with auditing", func() {
		Eventually(func(g Gomega) string {
			condition, err := kube.GetProjectStatusCondition(userData, status.ReadyType)
			g.Expect(err).ShouldNot(HaveOccurred())
			return condition
		}).Should(Equal("True"), "Condition status 'Ready' is not 'True'")
	})

	By("Remove Auditing from the project", func() {
		userData.Resources.Project.WithAuditing(nil)
		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
	})

	By("Check project status with auditing removed", func() {
		Eventually(func(g Gomega) string {
			condition, err := kube.GetProjectStatusCondition(userData, status.ReadyType)
			g.Expect(err).ShouldNot(HaveOccurred())
			return condition
		}).Should(Equal("True"), "Condition status 'Ready' is not 'True'")
	})
}

func exampleFilter() string {
	return "{\"atype\" : \"authenticate\", \"param\" : {\"user\" : \"auditReadOnly\", \"db\" : \"admin\", \"mechanism\" : \"SCRAM-SHA-1\"} }"
}
