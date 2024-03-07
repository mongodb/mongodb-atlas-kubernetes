package e2e_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

var _ = Describe("Teams", Label("teams"), func() {
	var testData *model.TestDataProvider

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Teams Test\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
			Expect(actions.SaveTeamsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Delete Resources", func() {
			actions.DeleteTestDataTeams(testData)
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test *model.TestDataProvider, teams []akov2.Team) {
			testData = test
			actions.ProjectCreationFlow(test)
			actions.AddTeamResourcesWithNUsers(test, teams, 1)
			projectTeamsFlow(test, teams)
		},
		Entry("Test[teams-1]: User has project to which a team was added",
			model.DataProvider(
				"teams-1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]akov2.Team{
				{
					TeamRef: common.ResourceRefNamespaced{
						Name: "my-team-1",
					},
					Roles: []akov2.TeamRole{
						akov2.TeamRoleOwner,
					},
				},
				{
					TeamRef: common.ResourceRefNamespaced{
						Name: "my-team-2",
					},
					Roles: []akov2.TeamRole{
						akov2.TeamRoleOwner,
					},
				},
			},
		),
	)
})

func projectTeamsFlow(userData *model.TestDataProvider, teams []akov2.Team) {
	By("Add Teams to project", func() {
		Expect(userData.K8SClient.Get(userData.Context, client.ObjectKeyFromObject(userData.Project), userData.Project)).Should(Succeed())
		userData.Project.Spec.Teams = teams
		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		Eventually(func(g Gomega) bool {
			return ensureTeamsStatus(g, *userData, teams, teamWasCreated)
		}).WithTimeout(10*time.Minute).WithPolling(20*time.Second).Should(BeTrue(), "Teams were not created")

		actions.WaitForConditionsToBecomeTrue(userData, status.ProjectTeamsReadyType, status.ReadyType)
	})

	By("Remove one team from the project", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
			Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())

		assignedTeams := userData.Project.Spec.Teams
		userData.Project.Spec.Teams = assignedTeams[:1]

		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		Eventually(func(g Gomega) bool {
			return ensureTeamsStatus(g, *userData, teams[1:], teamWasRemoved)
		}).WithTimeout(10*time.Minute).WithPolling(20*time.Second).Should(BeTrue(), "Team were not removed")

		actions.WaitForConditionsToBecomeTrue(userData, status.ProjectTeamsReadyType, status.ReadyType)
	})

	By("Update team role in the project", func() {
		userData.Project.Spec.Teams[0].Roles = []akov2.TeamRole{akov2.TeamRoleReadOnly}

		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		Eventually(func(g Gomega) bool {
			return ensureTeamsStatus(g, *userData, userData.Project.Spec.Teams, teamWasCreated)
		}).WithTimeout(10*time.Minute).WithPolling(20*time.Second).Should(BeTrue(), "Teams were not created")

		actions.WaitForConditionsToBecomeTrue(userData, status.ProjectTeamsReadyType, status.ReadyType)
	})

	By("Remove all teams from the project", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
			Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())

		userData.Project.Spec.Teams = nil

		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		Eventually(func(g Gomega) bool {
			return ensureTeamsStatus(g, *userData, teams, teamWasRemoved)
		}).WithTimeout(10*time.Minute).WithPolling(20*time.Second).Should(BeTrue(), "Team were not removed")

		actions.CheckProjectConditionsNotSet(userData, status.ProjectTeamsReadyType)
	})
}

func ensureTeamsStatus(g Gomega, testData model.TestDataProvider, teams []akov2.Team, check func(res *akov2.AtlasTeam) bool) bool {
	for _, team := range teams {
		resource := &akov2.AtlasTeam{}
		g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{Name: team.TeamRef.Name, Namespace: testData.Resources.Namespace}, resource)).Should(Succeed())

		if !check(resource) {
			return false
		}
	}

	return true
}

func teamWasCreated(team *akov2.AtlasTeam) bool {
	return team.Status.ID != ""
}

func teamWasRemoved(team *akov2.AtlasTeam) bool {
	return team.Status.ID == ""
}
