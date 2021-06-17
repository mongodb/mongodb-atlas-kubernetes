package e2e_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

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
		Eventually(kube.GetVersionOutput()).Should(Say(K8sVersion))
	})

	var _ = AfterEach(func() {
		By("Atfer each.", func() {
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
	It("[helm-ns] User can deploy operator namespaces by using HELM", func() {
		By("User creates configuration for a new Project and Cluster", func() {
			data = model.NewTestDataProvider(
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
				[]func(*model.TestDataProvider){},
			)
		})
		By("User use helm for deploying operator", func() {
			helm.AddMongoDBRepo()
			helm.InstallKubernetesOperatorNS(data.Resources)
		})
		deployCluster(data.Resources, data.PortGroup)
	})

	It("[helm-wide] User can deploy operator namespaces by using HELM", func() {
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
			data.Resources.Clusters[0].Spec.Project.Name = data.Resources.Project.GetK8sMetaName()
			// TODO helm template has equal ObjectMeta.Name and Spec.Name
			data.Resources.Clusters[0].ObjectMeta.Name = "cluster-from-helm-wide"
			data.Resources.Clusters[0].Spec.Name = "cluster-from-helm-wide"
		})
		By("User use helm for deploying operator", func() {
			helm.AddMongoDBRepo()
			helm.InstallKubernetesOperatorWide(data.Resources)
		})
		deployCluster(data.Resources, data.PortGroup)
	})

})

func deployCluster(userSpec model.UserInputs, appPort int) {
	By("User deploy cluster by helm", func() {
		helm.InstallCluster(userSpec)
	})
	By("Wait creation until is done", func() {
		actions.WaitProject(userSpec, "1")
		userSpec.ProjectID = kube.GetProjectResource(userSpec.Namespace, userSpec.K8sFullProjectName).Status.ID
		actions.WaitCluster(userSpec, "1")
	})

	By("Check attributes", func() {
		uCluster := mongocli.GetClustersInfo(userSpec.ProjectID, userSpec.Clusters[0].Spec.Name)
		actions.CompareClustersSpec(userSpec.Clusters[0].Spec, uCluster)
	})

	By("check database users Attibutes", func() {
		Eventually(actions.CheckIfUsersExist(userSpec), "2m", "10s").Should(BeTrue())
		actions.CheckUsersAttributes(userSpec)
	})

	By("Deploy application for user", func() {
		// kube apply application
		// send data
		// retrieve data
		actions.CheckUsersCanUseApplication(appPort, userSpec)
	})

	By("Check project, cluster does not exist", func() {
		helm.Uninstall(userSpec.Clusters[0].Spec.Name, userSpec.Namespace)
		Eventually(
			func() bool {
				return mongocli.IsProjectInfoExist(userSpec.ProjectID)
			},
			"5m", "20s",
		).Should(BeFalse(), "Project and cluster should be deleted from Atlas")
	})

	By("Delete HELM releases", func() {
		helm.UninstallKubernetesOperator(userSpec)
	})
}
