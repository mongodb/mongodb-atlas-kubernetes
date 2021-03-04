package e2e_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/mongocli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

// TODO add checks
type testDataProvider struct {
	description string
	confPath    string
	userSpec    userInputs
}

var _ = Describe("[cluster-ns] Configuration namespaced. Deploy cluster", func() {

	var data testDataProvider

	var _ = AfterEach(func() {
		if CurrentGinkgoTestDescription().Failed {
			GinkgoWriter.Write([]byte("Resources wasn't clean\n"))
			utils.SaveToFile(
				"output/operator-logs.txt",
				kube.GetManagerLogs(data.userSpec.namespace),
			)
			SaveK8sResources(
				[]string{"deploy", "atlasclusters", "atlasdatabaseusers", "atlasprojects"},
				data.userSpec.namespace,
			)
		} else {
			Eventually(kube.DeleteNamespace(data.userSpec.namespace)).Should(Say("deleted"))
		}
	})

	newData := func(description, path string) (string, testDataProvider) {
		var data testDataProvider
		data.description = description
		data.confPath = path
		data.userSpec = NewUserInputs("my-atlas-key")
		return data.description, data
	}

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test testDataProvider) {
			data = test
			mainCycle(test.confPath, test.userSpec)
		},
		Entry(newData("Simple Test with default cluster", "data/atlascluster_basic.yaml")),
		Entry(newData("Simple Backup Configuration and 2 users", "data/atlascluster_backup.yaml")),
		// Entry(newData("Multiregion, Backup and 2 users", "data/atlascluster_multiregion.yaml")), // TODO CLOUDP-83419
	)
})

func mainCycle(clusterConfigurationFile string, userSpec userInputs) {
	By("Prepare namespaces and project configuration", func() {
		kube.CreateNamespace(userSpec.namespace)
		project := utils.NewProject().
			ProjectName(userSpec.projectName).
			SecretRef(userSpec.keyName).
			CompleteK8sConfig(userSpec.k8sProjectName)
		utils.SaveToFile(FilePathTo(userSpec.projectName), project)

		userSpec.clusters = append(userSpec.clusters, utils.LoadUserClusterConfig(clusterConfigurationFile))
		userSpec.clusters[0].Spec.Project.Name = userSpec.k8sProjectName
		userSpec.clusters[0].Spec.ProviderSettings.InstanceSizeName = "M10"
		utils.SaveToFile(
			userSpec.clusters[0].ClusterFileName(),
			utils.JSONToYAMLConvert(userSpec.clusters[0]),
		)
	})

	By("Create namespaced Operator\n", func() {
		utils.CreateCopyKustomizeNamespace(userSpec.namespace)
		kube.Apply("-k", "data/"+userSpec.namespace)
		Eventually(
			kube.GetPodStatus(userSpec.namespace),
			"5m", "3s",
		).Should(Equal("Running"), "The operator should successfully run")
	})

	By("Create users resources", func() {
		kube.CreateApiKeySecret(userSpec.keyName, userSpec.namespace)
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

	By("check cluster Attribute", func() {
		cluster := mongocli.GetClustersInfo(userSpec.projectID, userSpec.clusters[0].Spec.Name)
		compareClustersSpec(userSpec.clusters[0].Spec, cluster)
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
		compareClustersSpec(userSpec.clusters[0].Spec, uCluster)
	})

	By("Delete cluster", func() {
		kube.Delete(userSpec.clusters[0].ClusterFileName(), "-n", userSpec.namespace)
		Eventually(
			checkIfClusterExist(userSpec),
			"10m", "1m",
		).Should(BeFalse(), "Cluster should be deleted from Atlas")
	})

	By("Delete project", func() {
		kube.Delete(userSpec.projectPath, "-n", userSpec.namespace)
		Eventually(
			func() bool {
				return mongocli.IsProjectInfoExist(userSpec.projectID)
			},
			"5m", "20s",
		).Should(BeFalse(), "Project should be deleted from Atlas")
	})
}
