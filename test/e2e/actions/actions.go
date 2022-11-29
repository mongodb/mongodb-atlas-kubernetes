// `actions` additional functions which accept testDataProvider struct and could be used as additional acctions in the tests since they all typical

package actions

import (
	"fmt"
	"time"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/kube"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"

	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/atlas"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

func UpdateSpecOfSelectedDeployment(spec v1.AtlasDeploymentSpec, indexOfDeployment int) func(data *model.TestDataProvider) {
	return func(data *model.TestDataProvider) {
		if len(data.InitialDeployments) < indexOfDeployment+1 {
			Fail("Index is out of range")
		}
		By(fmt.Sprintf("Update Deployment %s", data.InitialDeployments[indexOfDeployment].GetName()), func() {
			Expect(data.K8SClient.Get(data.Context, types.NamespacedName{Name: data.InitialDeployments[indexOfDeployment].ObjectMeta.GetName(),
				Namespace: data.Resources.Namespace}, data.InitialDeployments[indexOfDeployment])).To(Succeed())
			data.InitialDeployments[indexOfDeployment].Spec = spec
			Expect(data.K8SClient.Update(data.Context, data.InitialDeployments[indexOfDeployment])).To(Succeed())
			Eventually(kube.DeploymentReadyCondition(data)).WithTimeout(15*time.Minute).WithPolling(20*time.Second).Should(Equal("True"),
				fmt.Sprintf("Deployment is not ready. Status: %v", kube.GetDeploymentStatus(data)))
		})
	}
}

func changeFirstDeploymentPauseSpec(data *model.TestDataProvider, paused bool) {
	By(fmt.Sprintf("Setting pause to %v", paused), func() {
		Expect(data.K8SClient.Get(data.Context,
			types.NamespacedName{Name: data.InitialDeployments[0].GetName(),
				Namespace: data.Resources.Namespace},
			data.InitialDeployments[0])).Should(Succeed())
		updateSpec := data.InitialDeployments[0].Spec
		updateSpec.DeploymentSpec.Paused = &paused
		data.InitialDeployments[0].Spec = updateSpec
		Expect(data.K8SClient.Update(data.Context, data.InitialDeployments[0])).Should(Succeed())
		Eventually(kube.DeploymentReadyCondition(data)).WithTimeout(30*time.Minute).WithPolling(20*time.Second).Should(Equal("True"),
			fmt.Sprintf("Deployment is not ready. Status: %v", kube.GetDeploymentStatus(data)))
	})
	By("Check additional Deployment field `paused`\n", func() {
		aClient := atlas.GetClientOrFail()
		Eventually(func(g Gomega) {
			uDeployment := aClient.GetDeployment(data.Project.ID(), data.InitialDeployments[0].AtlasName())
			g.Expect(*uDeployment.Paused).Should(Equal(paused))
		}).WithTimeout(5 * time.Minute).WithPolling(20 * time.Second).Should(Succeed())
	})
}

func SuspendDeployment(data *model.TestDataProvider) {
	changeFirstDeploymentPauseSpec(data, true)
}

func ReactivateDeployment(data *model.TestDataProvider) {
	changeFirstDeploymentPauseSpec(data, false)
}

func DeleteFirstUser(data *model.TestDataProvider) {
	By("User can delete Database User", func() {
		// data.Resources.ProjectID = kube.GetProjectResource(data.Resources.Namespace, data.Resources.GetAtlasProjectFullKubeName()).Status.ID
		// since it is could be several users, we should
		// - delete k8s resource
		// - delete one user from the list,
		// - check Atlas doesn't have the initial user and have the rest
		By("Delete k8s resources")
		if len(data.Users) == 0 {
			Fail("No users to delete")
		}
		Expect(data.K8SClient.Get(data.Context, types.NamespacedName{Name: data.Users[0].Name, Namespace: data.Users[0].Namespace}, data.Users[0])).Should(Succeed())
		Expect(data.K8SClient.Delete(data.Context, data.Users[0])).Should(Succeed())
		Eventually(func(g Gomega) {
			aClient := atlas.GetClientOrFail()
			user, err := aClient.GetDBUser(data.Users[0].Spec.DatabaseName, data.Users[0].Spec.Username, data.Project.ID())
			g.Expect(err).To(BeNil())
			g.Expect(user).To(BeNil())
		}).WithTimeout(5 * time.Minute).WithPolling(20 * time.Second).Should(Succeed())

		// the rest users should be still there
		data.Users = data.Users[1:]
	})
}

func AddTeamResourcesWithNUsers(data *model.TestDataProvider, teams []v1.Team, n int) {
	By("Setup Teams", func() {
		aClient := atlas.GetClientOrFail()
		users, err := aClient.GetOrgUsers(data.Project.ID())
		Expect(err).ShouldNot(HaveOccurred())
		Expect(len(users) < n).ShouldNot(BeTrue())

		for _, team := range teams {
			By(fmt.Sprintf("Add Team \"%s\" resource to k8s", team.TeamRef.Name), func() {
				usernames := make([]v1.TeamUser, 0, n)
				for i := 0; i < n; i++ {
					usernames = append(usernames, v1.TeamUser(users[i].Username))
				}

				resource := model.NewTeam(team.TeamRef.Name, data.Resources.Namespace)
				resource.Spec.Usernames = usernames

				Expect(data.K8SClient.Create(data.Context, resource)).Should(Succeed())
				data.Teams = append(data.Teams, resource)
			})
		}
	})
}
