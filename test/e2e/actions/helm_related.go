package actions

import (
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	helm "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/helm"
	kube "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kube"
	mongocli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/mongocli"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

// helm update should change at least 1 field: databaseusers, project, cluster
func HelmDefaultUpgradeResouces(data *model.TestDataProvider) {
	By("User use HELM upgrade command for changing atlas resources\n", func() {
		data.Resources.Project.Spec.ProjectIPAccessList[0].Comment = "updated"
		enabled := true
		data.Resources.Clusters[0].Spec.ProviderBackupEnabled = &enabled
		data.Resources.Users[0].DeleteAllRoles()
		data.Resources.Users[0].AddBuildInAdminRole()
		data.Resources.Users[0].Spec.Project.Name = data.Resources.K8sFullProjectName
		generation, _ := strconv.Atoi(kube.GetGeneration(data.Resources.Namespace, data.Resources.Clusters[0].GetClusterNameResource()))
		helm.UpgradeAtlasClusterChart(data.Resources)

		By("Wait project creation", func() {
			WaitCluster(data.Resources, strconv.Itoa(generation+1))
			ExpectWithOffset(1, data.Resources.ProjectID).ShouldNot(BeEmpty())
		})
		updatedCluster := mongocli.GetClustersInfo(data.Resources.ProjectID, data.Resources.Clusters[0].Spec.Name)
		CompareClustersSpec(data.Resources.Clusters[0].Spec, updatedCluster)
		user := mongocli.GetUser(data.Resources.Users[0].Spec.Username, data.Resources.ProjectID)
		ExpectWithOffset(1, user.Roles[0].RoleName).Should(Equal(model.RoleBuildInAdmin))
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
		helm.UpgradeAtlasClusterChart(data.Resources)
		CheckUsersAttributes(data.Resources)
	})
}

// helm update: delete user
func HelmUpgradeDeleteFirstUser(data *model.TestDataProvider) {
	By("User delete database user from the Atlas\n", func() {
		data.Resources.Users = data.Resources.Users[1:]
		helm.UpgradeAtlasClusterChart(data.Resources)
		CheckUsersAttributes(data.Resources)
	})
}

// HelmUpgradeChartVersions upgrade chart version of crd, operator, and
func HelmUpgradeChartVersions(data *model.TestDataProvider) {
	By("User update helm chart (used main-branch)", func() {
		generation, _ := strconv.Atoi(kube.GetGeneration(data.Resources.Namespace, data.Resources.Clusters[0].GetClusterNameResource()))
		helm.UpgradeOperatorChart(data.Resources)
		helm.UpgradeAtlasClusterChartDev(data.Resources)

		By("Wait updating")
		WaitCluster(data.Resources, strconv.Itoa(generation+1))

		updatedCluster := mongocli.GetClustersInfo(data.Resources.ProjectID, data.Resources.Clusters[0].Spec.Name)
		CompareClustersSpec(data.Resources.Clusters[0].Spec, updatedCluster)
		CheckUsersAttributes(data.Resources)
	})
}
