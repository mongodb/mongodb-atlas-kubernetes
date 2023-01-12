package e2e_test

import (
	"k8s.io/apimachinery/pkg/types"

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

var _ = Describe("UserLogin", Label("custom-roles"), func() {
	var testData *model.TestDataProvider

	_ = BeforeEach(func() {
		Eventually(kubecli.GetVersionOutput()).Should(Say(K8sVersion))
	})

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Custom Roles Test\n"))
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
		func(test *model.TestDataProvider, customRoles []v1.CustomRole) {
			testData = test
			actions.ProjectCreationFlow(test)
			projectCustomRolesFlow(test, customRoles)
		},
		Entry("Test[custom-roles-1]: User has project to which custom roles where added",
			model.DataProvider(
				"custom-roles-1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]v1.CustomRole{
				{
					Name: "ShardingAdmin",
					InheritedRoles: []v1.Role{
						{
							Name:     "enableSharding",
							Database: "admin",
						},
						{
							Name:     "backup",
							Database: "admin",
						},
					},
					Actions: []v1.Action{
						{
							Name: "LIST_SESSIONS",
							Resources: []v1.Resource{
								{
									Cluster: toptr.MakePtr(true),
								},
							},
						},
						{
							Name: "KILL_ANY_SESSION",
							Resources: []v1.Resource{
								{
									Cluster: toptr.MakePtr(true),
								},
							},
						},
					},
				},
				{
					Name: "test",
					InheritedRoles: []v1.Role{
						{
							Name:     "readWrite",
							Database: "test",
						},
						{
							Name:     "dbAdmin",
							Database: "test",
						},
					},
				},
			},
		),
	)
})

func projectCustomRolesFlow(userData *model.TestDataProvider, customRoles []v1.CustomRole) {
	By("Add Custom Roles to the project", func() {
		userData.Project.Spec.CustomRoles = customRoles
		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		actions.WaitForConditionsToBecomeTrue(userData, status.ProjectCustomRolesReadyType, status.ReadyType)
	})

	By("Update Custom Role from the project", func() {
		crActions := userData.Project.Spec.CustomRoles[0].Actions
		crActions = append(crActions, v1.Action{
			Name: "USE_UUID",
			Resources: []v1.Resource{
				{
					Cluster: toptr.MakePtr(true),
				},
			},
		})
		userData.Project.Spec.CustomRoles[0].Actions = crActions

		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		actions.WaitForConditionsToBecomeTrue(userData, status.ProjectCustomRolesReadyType, status.ReadyType)
	})

	By("Remove one Custom Roles from the project", func() {
		Eventually(func(g Gomega) {
			g.Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
				Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())
			cr := userData.Project.Spec.CustomRoles
			userData.Project.Spec.CustomRoles = cr[:1]
			g.Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		}).Should(Succeed())
		actions.WaitForConditionsToBecomeTrue(userData, status.ProjectCustomRolesReadyType, status.ReadyType)
	})

	By("Remove all Custom Roles from the project", func() {
		Eventually(func(g Gomega) {
			g.Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
				Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())
			userData.Project.Spec.CustomRoles = nil
			g.Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		}).Should(Succeed())
		actions.CheckProjectConditionsNotSet(userData, status.ProjectCustomRolesReadyType)
	})
}
