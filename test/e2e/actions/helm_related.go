package actions

import (
	"fmt"
	"time"

	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/k8s"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/atlas"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	helm "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/helm"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

// helm update should change at least 1 field: databaseusers, project, deployment
func HelmDefaultUpgradeResources(data *model.TestDataProvider) {
	By("User use HELM upgrade command for changing atlas resources\n", func() {
		data.Resources.Project.Spec.ProjectIPAccessList[0].Comment = "updated"
		enabled := true
		data.Resources.Deployments[0].Spec.DeploymentSpec.ProviderBackupEnabled = &enabled
		data.Resources.Users[0].DeleteAllRoles()
		data.Resources.Users[0].AddBuildInAdminRole()
		data.Resources.Users[0].Spec.Project.Name = data.Resources.GetAtlasProjectFullKubeName()
		generation, err := kubecli.GetDeploymentObservedGeneration(data.Context, data.K8SClient, data.Resources.Namespace, data.Resources.Deployments[0].ObjectMeta.GetName())
		Expect(err).NotTo(HaveOccurred())
		helm.UpgradeAtlasDeploymentChartDev(data.Resources)

		By("Wait project creation", func() {
			WaitDeployment(data, generation+1)
			ExpectWithOffset(1, data.Resources.ProjectID).ShouldNot(BeEmpty())
		})
		aClient := atlas.GetClientOrFail()
		updatedDeployment := aClient.GetDeployment(data.Resources.ProjectID, data.Resources.Deployments[0].Spec.GetDeploymentName())
		CompareDeploymentsSpec(data.Resources.Deployments[0].Spec, updatedDeployment)
		Eventually(func() error {
			aClient := atlas.GetClientOrFail()
			user, err := aClient.GetDBUser("admin", data.Resources.Users[0].Spec.Username, data.Resources.ProjectID)
			if err != nil {
				return err
			}
			if user.Roles[0].RoleName != model.RoleBuildInAdmin {
				return fmt.Errorf("user role %s not equal to %s", user.Roles[0].RoleName, model.RoleBuildInAdmin)
			}
			return nil
		}).WithTimeout(7 * time.Minute).WithPolling(10 * time.Second).ShouldNot(HaveOccurred())
	})
}

// helm update: add user+change user role
func HelmUpgradeUsersRoleAddAdminUser(data *model.TestDataProvider) {
	By("User change role for all users and add new database user\n", func() {
		for i := range data.Resources.Users {
			data.Resources.Users[i].WithProjectRef(data.Resources.Project.GetK8sMetaName())
			data.Resources.Users[i].AddCustomRole(model.RoleCustomReadWrite, "Ships", "")
		}
		newUser := *model.NewDBUser("only-one-admin").
			WithAuthDatabase("admin").
			WithProjectRef(data.Resources.Project.GetK8sMetaName()).
			WithSecretRef("new-user-secret").
			AddBuildInAdminRole()
		data.Resources.Users = append(data.Resources.Users, newUser)
		helm.UpgradeAtlasDeploymentChartDev(data.Resources)
		CheckUsersAttributes(data)
	})
}

// helm update: delete user
func HelmUpgradeDeleteFirstUser(data *model.TestDataProvider) {
	By("User delete database user from the Atlas\n", func() {
		data.Resources.Users = data.Resources.Users[1:]
		helm.UpgradeAtlasDeploymentChartDev(data.Resources)
		CheckUsersAttributes(data)
	})
}

// HelmUpgradeChartVersions upgrade chart version of crd, operator, and
func HelmUpgradeChartVersions(data *model.TestDataProvider) {
	By("User update helm chart (used main-branch)", func() {
		generation, err := kubecli.GetDeploymentObservedGeneration(data.Context, data.K8SClient, data.Resources.Namespace, data.Resources.Deployments[0].ObjectMeta.GetName())
		Expect(err).NotTo(HaveOccurred())
		helm.UpgradeOperatorChart(data.Resources)
		helm.UpgradeAtlasDeploymentChartDev(data.Resources)

		By("Wait updating")
		WaitDeployment(data, generation+1)
		aClient := atlas.GetClientOrFail()
		updatedDeployment := aClient.GetDeployment(data.Resources.ProjectID, data.Resources.Deployments[0].Spec.GetDeploymentName())
		CompareDeploymentsSpec(data.Resources.Deployments[0].Spec, updatedDeployment)
		CheckUsersAttributes(data)
	})
}
