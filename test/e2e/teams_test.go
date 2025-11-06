// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2e_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/api/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
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
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(ctx SpecContext, test func(context.Context) *model.TestDataProvider, teams []akov2.Team) {
			testData = test(ctx)
			actions.ProjectCreationFlow(testData)
			actions.AddTeamResourcesWithNUsers(testData, teams, 1)
			projectTeamsFlow(testData, teams)
		},
		Entry("Test[teams-1]: User has project to which a team was added",
			func(ctx context.Context) *model.TestDataProvider {
				return model.DataProvider(ctx, "teams-1", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			[]akov2.Team{
				{
					TeamRef: common.ResourceRefNamespaced{
						Name: utils.RandomName("my-team-1"),
					},
					Roles: []akov2.TeamRole{
						akov2.TeamRoleOwner,
					},
				},
				{
					TeamRef: common.ResourceRefNamespaced{
						Name: utils.RandomName("my-team-2"),
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
			return ensureTeamsStatus(g, *userData, teams, teamWasAssigned)
		}).WithTimeout(5*time.Minute).WithPolling(10*time.Second).Should(BeTrue(), "Teams were not assigned")

		actions.WaitForConditionsToBecomeTrue(userData, api.ProjectTeamsReadyType, api.ReadyType)
	})

	By("De-assign one team from the project", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
			Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())

		assignedTeams := userData.Project.Spec.Teams
		userData.Project.Spec.Teams = assignedTeams[:1]

		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		Eventually(func(g Gomega) bool {
			return ensureTeamsStatus(g, *userData, teams[1:], teamWasDeAssigned)
		}).WithTimeout(5*time.Minute).WithPolling(10*time.Second).Should(BeTrue(), "Team were not removed")

		actions.WaitForConditionsToBecomeTrue(userData, api.ProjectTeamsReadyType, api.ReadyType)
	})

	By("Update team role in the project", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
			Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())
		userData.Project.Spec.Teams[0].Roles = []akov2.TeamRole{akov2.TeamRoleReadOnly}

		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		Eventually(func(g Gomega) bool {
			return ensureTeamsStatus(g, *userData, userData.Project.Spec.Teams, teamWasAssigned)
		}).WithTimeout(5*time.Minute).WithPolling(10*time.Second).Should(BeTrue(), "Teams were not assigned")

		actions.WaitForConditionsToBecomeTrue(userData, api.ProjectTeamsReadyType, api.ReadyType)
	})

	By("De-assign all teams from the project", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
			Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())

		userData.Project.Spec.Teams = nil

		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		Eventually(func(g Gomega) bool {
			return ensureTeamsStatus(g, *userData, teams, teamWasDeAssigned)
		}).WithTimeout(5*time.Minute).WithPolling(10*time.Second).Should(BeTrue(), "Teams were not de-assigned")

		actions.CheckProjectConditionsNotSet(userData, api.ProjectTeamsReadyType)
	})

	By("Cleanup Atlas Teams", func() {
		aClient := atlas.GetClientOrFail()

		for _, AssociatedTeam := range teams {
			team := &akov2.AtlasTeam{}
			Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: AssociatedTeam.TeamRef.Name, Namespace: userData.Resources.Namespace}, team)).Should(Succeed())
			_, err := aClient.Client.TeamsApi.DeleteTeam(userData.Context, aClient.OrgID, team.Status.ID).Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(userData.K8SClient.Delete(userData.Context, team)).To(Succeed())
		}
	})
}

func ensureTeamsStatus(g Gomega, testData model.TestDataProvider, teams []akov2.Team, check func(team *akov2.AtlasTeam, project *akov2.AtlasProject) bool) bool {
	for _, team := range teams {
		resource := &akov2.AtlasTeam{}
		g.Expect(testData.K8SClient.Get(testData.Context,
			types.NamespacedName{Name: team.TeamRef.Name, Namespace: testData.Resources.Namespace}, resource)).Should(Succeed())

		if !check(resource, testData.Project) {
			return false
		}
	}
	return true
}

func teamWasAssigned(team *akov2.AtlasTeam, project *akov2.AtlasProject) bool {
	if team.Status.ID == "" {
		return false
	}

	for _, p := range team.Status.Projects {
		if p.ID == project.ID() {
			return true
		}
	}

	return len(team.Finalizers) > 0
}

func teamWasDeAssigned(team *akov2.AtlasTeam, project *akov2.AtlasProject) bool {
	if team.Status.ID == "" {
		return false
	}

	for _, p := range team.Status.Projects {
		if p.ID == project.ID() {
			return false
		}
	}

	return len(team.Finalizers) == 0
}
