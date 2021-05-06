package e2e_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	actions "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	helm "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/helm"
	kube "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/mongocli"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

var _ = Describe("HELM charts", func() {
	var data model.TestDataProvider

	var _ = BeforeEach(func() {
		imageURL := os.Getenv("IMAGE_URL")
		Expect(imageURL).ShouldNot(BeEmpty(), "SetUP IMAGE_URL")
	})

	var _ = AfterEach(func() {
		By("Atfer each.", func() {
			GinkgoWriter.Write([]byte("\n"))
			GinkgoWriter.Write([]byte("===============================================\n"))
			GinkgoWriter.Write([]byte("Operator namespace: " + data.Resources.Namespace + "\n"))
			GinkgoWriter.Write([]byte("===============================================\n"))
			if CurrentGinkgoTestDescription().Failed {
				GinkgoWriter.Write([]byte("Resources wasn't clean"))
				utils.SaveToFile(
					"output/operator-logs.txt",
					kube.GetManagerLogs(config.DefaultOperatorNS),
				)
				actions.SaveK8sResources(
					[]string{"deploy"},
					"default",
				)
				actions.SaveK8sResources(
					[]string{"atlasclusters", "atlasdatabaseusers", "atlasprojects"},
					data.Resources.Namespace,
				)
				actions.SaveTestAppLogs(data.Resources)
			}
		})
	})

	DescribeTable("[helm-ns] Namespaced operators working only with its own namespace with different configuration",
		func(test model.TestDataProvider) {
			data = test
			By("User use helm for deploying namespaces operator", func() {
				helm.AddMongoDBRepo()
				helm.InstallKubernetesOperatorNS(data.Resources)
			})

			deployCluster(&data)
			By("Additional check for the current data set", func() {
				for _, check := range data.Actions {
					check(&data)
				}
			})
			deleteClusterAndOperator(&data)
		},
		Entry("Several actions with helm update",
			model.NewTestDataProvider(
				"helm-ns",
				[]string{"data/atlascluster_basic_helm.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("reader").
						WithSecretRef("dbuser-secret-u1").
						AddRole("readWrite", "Ships", "").
						WithAuthDatabase("admin"),
				},
				30006,
				[]func(*model.TestDataProvider){
					actions.HelmDefaultUpgradeResouces,
					actions.HelmUpgradeUsersRoleAddAdminUser,
					actions.HelmUpgradeDeleteFirstUser,
				},
			),
		),
	)

	Describe("[helm-wide] HELM charts.", func() {
		It("User can deploy operator namespaces by using HELM", func() {
			By("User creates configuration for a new Project and Cluster", func() {
				data = model.NewTestDataProvider(
					"helm-wide",
					[]string{"data/atlascluster_basic_helm.yaml"},
					[]string{},
					[]model.DBUser{
						*model.NewDBUser("reader2").
							WithSecretRef("dbuser-secret-u2").
							AddRole("readWrite", "Ships", "").
							WithAuthDatabase("admin"),
					},
					30007,
					[]func(*model.TestDataProvider){},
				)
				// helm template has equal ObjectMeta.Name and Spec.Name
				data.Resources.Clusters[0].ObjectMeta.Name = "cluster-from-helm-wide"
				data.Resources.Clusters[0].Spec.Name = "cluster-from-helm-wide"
			})
			By("User use helm for deploying operator", func() {
				helm.AddMongoDBRepo()
				helm.InstallKubernetesOperatorWide(data.Resources)
			})
			deployCluster(&data)
			deleteClusterAndOperator(&data)
		})
	})
})

func deployCluster(data *model.TestDataProvider) {
	By("User deploy cluster by helm", func() {
		helm.InstallCluster(data.Resources)
	})
	By("Wait creation until is done", func() {
		actions.WaitProject(data.Resources, "1")
		data.Resources.ProjectID = kube.GetProjectResource(data.Resources.Namespace, data.Resources.K8sFullProjectName).Status.ID
		actions.WaitCluster(data.Resources, "1")
	})

	By("Check attributes", func() {
		uCluster := mongocli.GetClustersInfo(data.Resources.ProjectID, data.Resources.Clusters[0].Spec.Name)
		actions.CompareClustersSpec(data.Resources.Clusters[0].Spec, uCluster)
	})

	By("check database users Attibutes", func() {
		Eventually(actions.CheckIfUsersExist(data.Resources), "2m", "10s").Should(BeTrue())
		actions.CheckUsersAttributes(data.Resources)
	})

	By("Deploy application for user", func() {
		actions.CheckUsersCanUseApplication(data.PortGroup, data.Resources)
	})
}

func deleteClusterAndOperator(data *model.TestDataProvider) {
	By("Check project, cluster does not exist", func() {
		helm.Uninstall(data.Resources.Clusters[0].Spec.Name, data.Resources.Namespace)
		Eventually(
			func() bool {
				return mongocli.IsProjectInfoExist(data.Resources.ProjectID)
			},
			"5m", "20s",
		).Should(BeFalse(), "Project and cluster should be deleted from Atlas")
	})

	By("Delete HELM releases", func() {
		helm.UninstallKubernetesOperator(data.Resources)
	})
}
