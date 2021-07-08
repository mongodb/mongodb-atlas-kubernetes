package e2e_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	actions "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	kube "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kube"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

var _ = Describe("[cluster-ns] Configuration namespaced. Deploy cluster", func() {
	var data model.TestDataProvider // TODO check it

	_ = BeforeEach(func() {
		Eventually(kube.GetVersionOutput()).Should(Say(K8sVersion))
	})
	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + data.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentGinkgoTestDescription().Failed {
			GinkgoWriter.Write([]byte("Test has been failed. Trying to save logs...\n"))
			utils.SaveToFile(
				fmt.Sprintf("output/%s/operatorDecribe.txt", data.Resources.Namespace),
				[]byte(kube.DescribeOperatorPod(data.Resources.Namespace)),
			)
			utils.SaveToFile(
				fmt.Sprintf("output/%s/operator-logs.txt", data.Resources.Namespace),
				kube.GetManagerLogs(data.Resources.Namespace),
			)
			actions.SaveTestAppLogs(data.Resources)
			actions.SaveK8sResources(
				[]string{"deploy", "atlasclusters", "atlasdatabaseusers", "atlasprojects"},
				data.Resources.Namespace,
			)
		} else {
			actions.AfterEachFinalCleanup([]model.TestDataProvider{data})
		}
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test model.TestDataProvider) {
			data = test
			mainCycle(test)
		},
		Entry("Trial - Simplest configuration with no backup and one Admin User",
			model.NewTestDataProvider(
				"operator-ns-trial",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlascluster_basic.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("user1").
						WithSecretRef("dbuser-secret-u1").
						AddBuildInAdminRole(),
				},
				30000,
				[]func(*model.TestDataProvider){
					actions.DeleteFirstUser,
				},
			),
		),
		Entry("Almost Production - Backup and 2 DB users: one Admin and one read-only",
			model.NewTestDataProvider(
				"operator-ns-prodlike",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlascluster_backup.yaml"},
				[]string{"data/atlascluster_backup_update.yaml"},
				[]model.DBUser{
					*model.NewDBUser("admin").
						WithSecretRef("dbuser-admin-secret-u1").
						AddBuildInAdminRole(),
					*model.NewDBUser("user2").
						WithSecretRef("dbuser-secret-u2").
						AddCustomRole(model.RoleCustomReadWrite, "Ships", ""),
				},
				30001,
				[]func(*model.TestDataProvider){
					actions.UpdateClusterFromUpdateConfig,
					actions.SuspendCluster,
					actions.ReactivateCluster,
					actions.DeleteFirstUser,
				},
			),
		),
		Entry("Multiregion, Backup and 2 DBUsers",
			model.NewTestDataProvider(
				"operator-ns-multiregion",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlascluster_multiregion.yaml"},
				[]string{"data/atlascluster_multiregion_update.yaml"},
				[]model.DBUser{
					*model.NewDBUser("user1").
						WithSecretRef("dbuser-secret-u1").
						AddBuildInAdminRole(),
					*model.NewDBUser("user2").
						WithSecretRef("dbuser-secret-u2").
						AddBuildInAdminRole(),
				},
				30003,
				[]func(*model.TestDataProvider){
					actions.SuspendCluster,
					actions.ReactivateCluster,
					actions.DeleteFirstUser,
				},
			),
		),
		Entry("Product Owner - Simplest configuration with ProjectOwner and update cluster to have backup",
			model.NewTestDataProvider(
				"operator-ns-product-owner",
				model.NewEmptyAtlasKeyType().WithRoles([]model.AtlasRoles{model.GroupOwner}).WithWhiteList([]string{"0.0.0.1/1", "128.0.0.0/1"}),
				[]string{"data/atlascluster_backup.yaml"},
				[]string{"data/atlascluster_backup_update_remove_backup.yaml"},
				[]model.DBUser{
					*model.NewDBUser("user1").
						WithSecretRef("dbuser-secret-u1").
						AddBuildInAdminRole(),
				},
				30010,
				[]func(*model.TestDataProvider){
					actions.UpdateClusterFromUpdateConfig,
				},
			),
		),
		Entry("Trial - Global connection",
			model.NewTestDataProvider(
				"operator-ns-trial-global",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess().CreateAsGlobalLevelKey(),
				[]string{"data/atlascluster_basic.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("user1").
						WithSecretRef("dbuser-secret-u1").
						AddBuildInAdminRole(),
				},
				30011,
				[]func(*model.TestDataProvider){
					actions.DeleteFirstUser,
				},
			),
		),
	)
})

func mainCycle(data model.TestDataProvider) {
	actions.PrepareUsersConfigurations(&data)
	actions.DeployNamespacedOperatorKuber(&data)

	By("Deploy User Resouces", func() {
		actions.DeployUserResourcesAction(&data)
		Expect(data.Resources.ProjectID).ShouldNot(BeEmpty())
	})

	By("Additional check for the current data set", func() {
		for _, check := range data.Actions {
			check(&data)
		}
	})
	By("Delete User Resources", func() {
		actions.DeleteUserResources(&data)
	})
}
