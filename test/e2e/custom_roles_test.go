package e2e_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

var _ = Describe("UserLogin", Label("custom-roles"), func() {
	var testData *model.TestDataProvider

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
		func(test *model.TestDataProvider, customRoles []akov2.CustomRole) {
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
			[]akov2.CustomRole{
				{
					Name: "ShardingAdmin",
					InheritedRoles: []akov2.Role{
						{
							Name:     "enableSharding",
							Database: "admin",
						},
						{
							Name:     "backup",
							Database: "admin",
						},
					},
					Actions: []akov2.Action{
						{
							Name: "LIST_SESSIONS",
							Resources: []akov2.Resource{
								{
									Cluster: pointer.MakePtr(true),
								},
							},
						},
						{
							Name: "KILL_ANY_SESSION",
							Resources: []akov2.Resource{
								{
									Cluster: pointer.MakePtr(true),
								},
							},
						},
					},
				},
				{
					Name: "test",
					InheritedRoles: []akov2.Role{
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

func projectCustomRolesFlow(userData *model.TestDataProvider, customRoles []akov2.CustomRole) {
	By("Add Custom Roles to the project", func() {
		userData.Project.Spec.CustomRoles = customRoles
		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		actions.WaitForConditionsToBecomeTrue(userData, api.ProjectCustomRolesReadyType, api.ReadyType)
	})

	By("Update Custom Role from the project", func() {
		crActions := userData.Project.Spec.CustomRoles[0].Actions
		crActions = append(crActions, akov2.Action{
			Name: "USE_UUID",
			Resources: []akov2.Resource{
				{
					Cluster: pointer.MakePtr(true),
				},
			},
		})
		userData.Project.Spec.CustomRoles[0].Actions = crActions

		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		actions.WaitForConditionsToBecomeTrue(userData, api.ProjectCustomRolesReadyType, api.ReadyType)
	})

	By("Remove one Custom Roles from the project", func() {
		Eventually(func(g Gomega) {
			g.Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
				Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())
			cr := userData.Project.Spec.CustomRoles
			userData.Project.Spec.CustomRoles = cr[:1]
			g.Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		}).Should(Succeed())
		actions.WaitForConditionsToBecomeTrue(userData, api.ProjectCustomRolesReadyType, api.ReadyType)
	})

	By("Remove all Custom Roles from the project", func() {
		Eventually(func(g Gomega) {
			g.Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
				Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())
			userData.Project.Spec.CustomRoles = nil
			g.Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		}).Should(Succeed())
		actions.CheckProjectConditionsNotSet(userData, api.ProjectCustomRolesReadyType)
	})
}
