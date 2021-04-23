package e2e_test

import (
	"os"

	. "github.com/onsi/ginkgo"
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
	var userSpec model.UserInputs

	var _ = BeforeEach(func() {
		imageURL := os.Getenv("IMAGE_URL")
		Expect(imageURL).ShouldNot(BeEmpty(), "SetUP IMAGE_URL")
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
					userSpec.Namespace,
				)
			}
		})
	})
	It("[helm-ns] User can deploy operator namespaces by using HELM", func() {
		By("User creates configuration for a new Project and Cluster", func() {
			userSpec = model.NewUserInputs(
				"only-key",
				[]model.DBUser{
					*model.NewDBUser("reader").
						WithSecretRef("dbuser-secret-u1").
						AddRole("readWrite", "Ships", "").
						WithAuthDatabase("admin"),
				},
			)
			userSpec.Clusters = append(userSpec.Clusters, model.LoadUserClusterConfig("data/atlascluster_basic_helm.yaml"))
			userSpec.Clusters[0].Spec.Project.Name = userSpec.Project.GetK8sMetaName()
			userSpec.Clusters[0].ObjectMeta.Name = "cluster-from-helm" // helm template has equal ObjectMeta.Name and Spec.Name
			userSpec.Clusters[0].Spec.Name = "cluster-from-helm"
		})
		By("User use helm for deploying operator", func() {
			helm.AddMongoDBRepo()
			helm.InstallKubernetesOperatorNS(userSpec)
		})
		deployCluster(userSpec, 30006)
	})

	It("[helm-wide] User can deploy operator namespaces by using HELM", func() {
		By("User creates configuration for a new Project and Cluster", func() {
			userSpec = model.NewUserInputs(
				"only-key",
				[]model.DBUser{
					*model.NewDBUser("reader2").
						WithSecretRef("dbuser-secret-u2").
						AddRole("readWrite", "Ships", "").
						WithAuthDatabase("admin"),
				},
			)
			userSpec.Clusters = append(userSpec.Clusters, model.LoadUserClusterConfig("data/atlascluster_basic_helm.yaml"))
			userSpec.Clusters[0].Spec.Project.Name = userSpec.Project.GetK8sMetaName()
			userSpec.Clusters[0].ObjectMeta.Name = "cluster-from-helm-wide" // helm template has equal ObjectMeta.Name and Spec.Name
			userSpec.Clusters[0].Spec.Name = "cluster-from-helm-wide"
		})
		By("User use helm for deploying operator", func() {
			helm.AddMongoDBRepo()
			helm.InstallKubernetesOperatorWide(userSpec)
		})
		deployCluster(userSpec, 30007)
	})

})

func deployCluster(userSpec model.UserInputs, appPort int) {
	By("User deploy cluster by helm", func() {
		kube.CreateApiKeySecret(userSpec.KeyName, userSpec.Namespace)
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

	By("Delete HELM release", func() {
		helm.UninstallKubernetesOperator(userSpec)
	})
}
