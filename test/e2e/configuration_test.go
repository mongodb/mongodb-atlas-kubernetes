package e2e_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/deploy"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

var _ = Describe("Configuration namespaced. Deploy deployment", Label("deployment-ns"), func() {
	var data model.TestDataProvider

	BeforeEach(func() {
		Eventually(kubecli.GetVersionOutput()).Should(Say(K8sVersion))
	})
	AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + data.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			GinkgoWriter.Write([]byte("Test has been failed. Trying to save logs...\n"))
			utils.SaveToFile(
				fmt.Sprintf("output/%s/operatorDecribe.txt", data.Resources.Namespace),
				[]byte(kubecli.DescribeOperatorPod(data.Resources.Namespace)),
			)
			utils.SaveToFile(
				fmt.Sprintf("output/%s/operator-logs.txt", data.Resources.Namespace),
				kubecli.GetManagerLogs(data.Resources.Namespace),
			)
			actions.SaveTestAppLogs(data.Resources)
			actions.SaveProjectsToFile(data.Context, data.K8SClient, data.Resources.Namespace)
			actions.SaveK8sResources(
				[]string{"deploy", "atlasdeployments", "atlasdatabaseusers"},
				data.Resources.Namespace,
			)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{data})
		}
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test model.TestDataProvider) {
			data = test
			mainCycle(test)
		},
		Entry("Trial - Simplest configuration with no backup and one Admin User", Label("ns-trial"),
			model.NewTestDataProvider(
				"operator-ns-trial",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlasdeployment_basic.yaml"},
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
		Entry("Almost Production - Backup and 2 DB users: one Admin and one read-only", Label("ns-backup2db", "long-run"),
			model.NewTestDataProvider(
				"operator-ns-prodlike",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlasdeployment_backup.yaml"},
				[]string{"data/atlasdeployment_backup_update.yaml"},
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
					actions.UpdateDeploymentFromUpdateConfig,
					actions.SuspendDeployment,
					actions.ReactivateDeployment,
					actions.DeleteFirstUser,
				},
			),
		),
		Entry("Multiregion AWS, Backup and 2 DBUsers", Label("ns-multiregion-aws-2"),
			model.NewTestDataProvider(
				"operator-ns-multiregion-aws",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlasdeployment_multiregion_aws.yaml"},
				[]string{"data/atlasdeployment_multiregion_aws_update.yaml"},
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
					actions.SuspendDeployment,
					actions.ReactivateDeployment,
					actions.DeleteFirstUser,
				},
			),
		),
		Entry("Multiregion Azure, Backup and 1 DBUser", Label("ns-multiregion-azure-1"),
			model.NewTestDataProvider(
				"operator-multiregion-azure",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess().CreateAsGlobalLevelKey(),
				[]string{"data/atlasdeployment_multiregion_azure.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("user1").
						WithSecretRef("dbuser-secret-u1").
						AddBuildInAdminRole(),
				},
				30012,
				[]func(*model.TestDataProvider){
					actions.DeleteFirstUser,
				},
			),
		),
		Entry("Multiregion GCP, Backup and 1 DBUser", Label("ns-multiregion-gcp-1"),
			model.NewTestDataProvider(
				"operator-multiregion-gcp",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess().CreateAsGlobalLevelKey(),
				[]string{"data/atlasdeployment_multiregion_gcp.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("user1").
						WithSecretRef("dbuser-secret-u1").
						AddBuildInAdminRole(),
				},
				30013,
				[]func(*model.TestDataProvider){
					actions.DeleteFirstUser,
				},
			),
		),
		Entry("Product Owner - Simplest configuration with ProjectOwner and update deployment to have backup", Label("ns-owner", "long-run"),
			model.NewTestDataProvider(
				"operator-ns-product-owner",
				model.AProject{},
				model.NewEmptyAtlasKeyType().WithRoles([]model.AtlasRoles{model.GroupOwner}).WithWhiteList([]string{"0.0.0.1/1", "128.0.0.0/1"}),
				[]string{"data/atlasdeployment_backup.yaml"},
				[]string{"data/atlasdeployment_backup_update_remove_backup.yaml"},
				[]model.DBUser{
					*model.NewDBUser("user1").
						WithSecretRef("dbuser-secret-u1").
						AddBuildInAdminRole(),
				},
				30010,
				[]func(*model.TestDataProvider){
					actions.UpdateDeploymentFromUpdateConfig,
				},
			),
		),
		Entry("Trial - Global connection", Label("ns-global-key"),
			model.NewTestDataProvider(
				"operator-ns-trial-global",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess().CreateAsGlobalLevelKey(),
				[]string{"data/atlasdeployment_basic.yaml"},
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
		Entry("Free - Users can use M0, default key", Label("ns-m0"),
			model.NewTestDataProvider(
				"operator-ns-free",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlasdeployment_basic_free.yaml"},
				[]string{""},
				[]model.DBUser{
					*model.NewDBUser("user").
						WithSecretRef("dbuser-secret").
						AddBuildInAdminRole(),
				},
				30016,
				[]func(*model.TestDataProvider){
					actions.DeleteFirstUser,
				},
			),
		),
		Entry("Free - Users can use M0, global", Label("ns-global-key-m0"),
			model.NewTestDataProvider(
				"operator-ns-free",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess().CreateAsGlobalLevelKey(),
				[]string{"data/atlasdeployment_basic_free.yaml"},
				[]string{""},
				[]model.DBUser{
					*model.NewDBUser("user").
						WithSecretRef("dbuser-secret").
						AddBuildInAdminRole(),
				},
				30017,
				[]func(*model.TestDataProvider){
					actions.DeleteFirstUser,
				},
			),
		),
	)
})

func mainCycle(data model.TestDataProvider) {
	actions.PrepareUsersConfigurations(&data)
	deploy.NamespacedOperator(&data)

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
