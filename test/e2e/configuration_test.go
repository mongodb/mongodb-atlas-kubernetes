package e2e_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onsi/gomega/gbytes"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/mongocli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

var _ = Describe("Deploy simple cluster", func() {

	userSpec := NewUserInputs("my-atlas-key")

	var _ = AfterEach(func() {
		GinkgoWriter.Write([]byte(userSpec.projectID))
		Eventually(kube.DeleteNamespace(userSpec.namespace)).Should(Say("deleted"))
		// mongocli.DeleteCluster(ID, "cluster45") // TODO struct
	})

	It("Release sample all-in-one.yaml should work", func() {
		By("Prepare namespaces and project configuration", func() {
			kube.CreateNamespace(userSpec.namespace)
			project := utils.NewProject().
				ProjectName(userSpec.projectName).
				SecretRef(userSpec.keyName).
				CompleteK8sConfig(userSpec.k8sProjectName)
			utils.SaveToFile(FilePathTo(userSpec.projectName), project)

			userSpec.clusters = append(userSpec.clusters, utils.LoadUserClusterConfig(ClusterSample))
			userSpec.clusters[0].Spec.Project.Name = userSpec.k8sProjectName
			userSpec.clusters[0].Spec.ProviderSettings.InstanceSizeName = "M10"
			userSpec.clusters[0].ObjectMeta.Name = "init-cluster"
			utils.SaveToFile(
				userSpec.clusters[0].ClusterFileName(),
				utils.JSONToYAMLConvert(userSpec.clusters[0]),
			)
		})

		By("Apply All-in-one operator wide configuration\n", func() {
			kube.Apply(ConfigAll)
			Eventually(
				kube.GetPodStatus(defaultOperatorNS),
				"5m", "3s",
			).Should(Equal("Running"))
		})

		By("Create users resources", func() {
			kube.CreateKeySecret(userSpec.keyName, userSpec.namespace)
			kube.Apply(FilePathTo(userSpec.projectName), "-n", userSpec.namespace)
			kube.Apply(userSpec.clusters[0].ClusterFileName(), "-n", userSpec.namespace)
		})

		By("Wait project creation", func() {
			waitProject(userSpec, "1")
			userSpec.projectID = kube.GetProjectResource(userSpec.namespace, userSpec.k8sFullProjectName).Status.ID
		})

		By("Wait cluster creation", func() {
			waitCluster(userSpec, "1")
		})

		By("check cluster Attribute", func() { // TODO ...
			cluster := mongocli.GetClustersInfo(userSpec.projectID, userSpec.clusters[0].Spec.Name)
			Expect(
				cluster.ProviderSettings.InstanceSizeName,
			).Should(Equal(userSpec.clusters[0].Spec.ProviderSettings.InstanceSizeName))
			Expect(
				cluster.ProviderSettings.ProviderName,
			).Should(Equal(string(userSpec.clusters[0].Spec.ProviderSettings.ProviderName)))
			Expect(
				cluster.ProviderSettings.RegionName,
			).Should(Equal(userSpec.clusters[0].Spec.ProviderSettings.RegionName))
		})

		By("Update cluster\n", func() {
			userSpec.clusters[0].Spec.ProviderSettings.InstanceSizeName = "M20"
			utils.SaveToFile(
				userSpec.clusters[0].ClusterFileName(),
				utils.JSONToYAMLConvert(userSpec.clusters[0]),
			)
			kube.Apply(userSpec.clusters[0].ClusterFileName(), "-n", userSpec.namespace)
		})

		By("Wait cluster updating", func() {
			waitCluster(userSpec, "2")
		})

		By("Check attributes", func() {
			uCluster := mongocli.GetClustersInfo(userSpec.projectID, userSpec.clusters[0].Spec.Name)
			Expect(
				uCluster.ProviderSettings.InstanceSizeName,
			).Should(Equal(
				userSpec.clusters[0].Spec.ProviderSettings.InstanceSizeName,
			))
		})

		By("Delete cluster", func() {
			kube.Delete(userSpec.clusters[0].ClusterFileName(), "-n", userSpec.namespace)
			Eventually(
				checkIfClusterExist(userSpec),
				"10m", "1m",
			).Should(BeFalse())
		})

		By("Delete project", func() {
			kube.Delete(userSpec.projectPath, "-n", userSpec.namespace)
			Eventually(
				checkIfProjectExist(userSpec),
				"5m", "20s",
			).Should(BeFalse())
		})
	})
})
