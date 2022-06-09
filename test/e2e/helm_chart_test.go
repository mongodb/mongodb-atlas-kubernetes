package e2e_test

import (
	"fmt"
	"os"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/atlas"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	actions "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	kube "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/kube"
	helm "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/helm"

	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/mongocli"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

var _ = Describe("HELM charts", func() {
	var data model.TestDataProvider

	_ = BeforeEach(func() {
		imageURL := os.Getenv("IMAGE_URL")
		Expect(imageURL).ShouldNot(BeEmpty(), "SetUP IMAGE_URL")
		Eventually(kubecli.GetVersionOutput()).Should(Say(K8sVersion))
	})

	_ = AfterEach(func() {
		By("Atfer each.", func() {
			GinkgoWriter.Write([]byte("\n"))
			GinkgoWriter.Write([]byte("===============================================\n"))
			GinkgoWriter.Write([]byte("Operator namespace: " + data.Resources.Namespace + "\n"))
			GinkgoWriter.Write([]byte("===============================================\n"))
			if CurrentSpecReport().Failed() {
				GinkgoWriter.Write([]byte("Resources wasn't clean"))
				utils.SaveToFile(
					fmt.Sprintf("output/%s/operator-logs-default.txt", data.Resources.Namespace),
					kubecli.GetManagerLogs(config.DefaultOperatorNS),
				)
				utils.SaveToFile(
					fmt.Sprintf("output/%s/operator-logs.txt", data.Resources.Namespace),
					kubecli.GetManagerLogs(data.Resources.Namespace),
				)
				actions.SaveK8sResourcesTo(
					[]string{"deploy"},
					"default",
					data.Resources.Namespace,
				)
				actions.SaveK8sResources(
					[]string{"atlasdeployments", "atlasdatabaseusers", "atlasprojects"},
					data.Resources.Namespace,
				)
				actions.SaveTestAppLogs(data.Resources)
				actions.DeleteGlobalKeyIfExist(&data)
			}
		})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration", Label("helm-ns"),
		func(test model.TestDataProvider, clusterType string) { // clusterType - probably will be moved later ()
			data = test
			GinkgoWriter.Println(data.Resources.KeyName)
			switch clusterType {
			case "advanced":
				data.Resources.Clusters[0].Spec.AdvancedDeploymentSpec.Name = data.Resources.KeyName
			case "serverless":
				data.Resources.Clusters[0].Spec.ServerlessSpec.Name = data.Resources.KeyName
			default:
				data.Resources.Clusters[0].Spec.DeploymentSpec.Name = data.Resources.KeyName
			}
			data.Resources.Clusters[0].ObjectMeta.Name = data.Resources.KeyName

			By("Install CRD", func() {
				helm.InstallCRD(data.Resources)
			})
			By("User use helm for deploying namespaces operator", func() {
				helm.InstallOperatorNamespacedSubmodule(data.Resources)
			})
			By("User deploy cluster by helm", func() {
				helm.InstallClusterSubmodule(data.Resources)
			})
			waitClusterWithChecks(&data)
			By("Additional check for the current data set", func() {
				for _, check := range data.Actions {
					check(&data)
				}
			})
			deleteClusterAndOperator(&data)
		},
		Entry("Several actions with helm update", Label("helm-ns-flow"),
			model.NewTestDataProvider(
				"helm-ns",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlascluster_basic_helm.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("reader").
						WithSecretRef("dbuser-secret-u1").
						AddCustomRole(model.RoleCustomReadWrite, "Ships", "").
						WithAuthDatabase("admin"),
				},
				30006,
				[]func(*model.TestDataProvider){
					actions.HelmDefaultUpgradeResouces,
					actions.HelmUpgradeUsersRoleAddAdminUser,
					actions.HelmUpgradeDeleteFirstUser,
				},
			),
			"default",
		),
		Entry("Advanced cluster by helm chart",
			model.NewTestDataProvider(
				"helm-advanced",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlascluster_advanced_helm.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("reader2").
						WithSecretRef("dbuser-secret-u2").
						AddCustomRole(model.RoleCustomReadWrite, "Ships", "").
						WithAuthDatabase("admin"),
				},
				30014,
				[]func(*model.TestDataProvider){},
			),
			"advanced",
		),
		Entry("Advanced multiregion cluster by helm chart",
			model.NewTestDataProvider(
				"helm-advanced-multiregion",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlascluster_advanced_multi_region_helm.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("reader2").
						WithSecretRef("dbuser-secret-u2").
						AddCustomRole(model.RoleCustomReadWrite, "Ships", "").
						WithAuthDatabase("admin"),
				},
				30015,
				[]func(*model.TestDataProvider){},
			),
			"advanced",
		),
		Entry("Serverless cluster by helm chart",
			model.NewTestDataProvider(
				"helm-serverless",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlascluster_serverless.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("reader2").
						WithSecretRef("dbuser-secret-u2").
						AddCustomRole(model.RoleCustomReadWrite, "Ships", "").
						WithAuthDatabase("admin"),
				},
				30016,
				[]func(*model.TestDataProvider){},
			),
			"serverless",
		),
	)

	Describe("HELM charts.", Label("helm-wide"), func() {
		It("User can deploy operator namespaces by using HELM", func() {
			By("User creates configuration for a new Project and Cluster", func() {
				data = model.NewTestDataProvider(
					"helm-wide",
					model.AProject{},
					model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
					[]string{"data/atlascluster_basic_helm.yaml"},
					[]string{},
					[]model.DBUser{
						*model.NewDBUser("reader2").
							WithSecretRef("dbuser-secret-u2").
							AddCustomRole(model.RoleCustomReadWrite, "Ships", "").
							WithAuthDatabase("admin"),
					},
					30007,
					[]func(*model.TestDataProvider){},
				)
				// helm template has equal ObjectMeta.Name and Spec.Name
				data.Resources.Clusters[0].ObjectMeta.Name = "cluster-from-helm-wide"
				data.Resources.Clusters[0].Spec.DeploymentSpec.Name = "cluster-from-helm-wide"
			})
			By("User use helm for deploying operator", func() {
				helm.InstallOperatorWideSubmodule(data.Resources)
			})
			By("User deploy cluster by helm", func() {
				helm.InstallClusterSubmodule(data.Resources)
			})
			waitClusterWithChecks(&data)
			deleteClusterAndOperator(&data)
		})
	})

	Describe("HELM charts.", Label("helm-update"), func() {
		It("User deploy operator and later deploy new version of the Atlas operator", func() {
			By("User creates configuration for a new Project, Cluster, DBUser", func() {
				data = model.NewTestDataProvider(
					"helm-upgrade",
					model.AProject{},
					model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
					[]string{"data/atlascluster_basic_helm.yaml"},
					[]string{},
					[]model.DBUser{
						*model.NewDBUser("admin").
							WithSecretRef("dbuser-secret-u2").
							AddBuildInAdminRole().
							WithAuthDatabase("admin"),
					},
					30010,
					[]func(*model.TestDataProvider){},
				)
				// helm template has equal ObjectMeta.Name and Spec.Name
				data.Resources.Clusters[0].ObjectMeta.Name = "cluster-from-helm-upgrade"
				data.Resources.Clusters[0].Spec.DeploymentSpec.Name = "cluster-from-helm-upgrade"
			})
			By("User use helm for last released version of operator and deploy his resouces", func() {
				helm.AddMongoDBRepo()
				helm.InstallOperatorNamespacedFromLatestRelease(data.Resources)
				helm.InstallClusterRelease(data.Resources)
				waitClusterWithChecks(&data)
			})
			By("User update new released operator", func() {
				backup := true
				data.Resources.Clusters[0].Spec.DeploymentSpec.ProviderBackupEnabled = &backup
				actions.HelmUpgradeChartVersions(&data)
				actions.CheckUsersCanUseOldApp(&data)
			})
			By("Delete Resources", func() {
				deleteClusterAndOperator(&data)
			})
		})
	})
})

func waitClusterWithChecks(data *model.TestDataProvider) {
	By("Wait creation until is done", func() {
		actions.WaitProject(data, "1")
		resource, err := kube.GetProjectResource(data)
		Expect(err).Should(BeNil())
		data.Resources.ProjectID = resource.Status.ID
		actions.WaitCluster(data.Resources, "1")
	})

	By("Check attributes", func() {
		cluster := data.Resources.Clusters[0]
		switch {
		case cluster.Spec.AdvancedDeploymentSpec != nil:
			atlasClient, err := atlas.AClient()
			Expect(err).To(BeNil())
			advancedCluster, err := atlasClient.GetAdvancedDeployment(data.Resources.ProjectID, cluster.Spec.AdvancedDeploymentSpec.Name)
			Expect(err).To(BeNil())
			actions.CompareAdvancedDeploymentsSpec(cluster.Spec, *advancedCluster)
		case cluster.Spec.ServerlessSpec != nil:
			atlasClient, err := atlas.AClient()
			Expect(err).To(BeNil())
			serverlessInstance, err := atlasClient.GetServerlessInstance(data.Resources.ProjectID, cluster.Spec.ServerlessSpec.Name)
			Expect(err).To(BeNil())
			actions.CompareServerlessSpec(cluster.Spec, *serverlessInstance)
		default:
			uCluster := mongocli.GetClustersInfo(data.Resources.ProjectID, data.Resources.Clusters[0].Spec.DeploymentSpec.Name)
			actions.CompareClustersSpec(cluster.Spec, uCluster)
		}
	})

	By("check database users Attributes", func() {
		Eventually(actions.CheckIfUsersExist(data.Resources), "2m", "10s").Should(BeTrue())
		actions.CheckUsersAttributes(data.Resources)
	})

	if !data.SkipAppConnectivityCheck {
		By("Deploy application for user", func() {
			actions.CheckUsersCanUseApp(data)
		})
	}
}

func deleteClusterAndOperator(data *model.TestDataProvider) {
	By("Check project, cluster does not exist", func() {
		helm.Uninstall(data.Resources.Clusters[0].Spec.GetClusterName(), data.Resources.Namespace)
		Eventually(
			func() bool {
				return mongocli.IsProjectInfoExist(data.Resources.ProjectID)
			},
			"7m", "20s",
		).Should(BeFalse(), "Project and cluster should be deleted from Atlas")
	})

	By("Delete HELM releases", func() {
		helm.UninstallKubernetesOperator(data.Resources)
	})
}
