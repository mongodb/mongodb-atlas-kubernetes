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
	resources   userInputs
}

var _ = Describe("[cluster-ns] Configuration namespaced. Deploy cluster", func() {

	var data testDataProvider

	var _ = AfterEach(func() {
		if CurrentGinkgoTestDescription().Failed {
			GinkgoWriter.Write([]byte("Resources wasn't clean\n"))
			utils.SaveToFile(
				"output/operator-logs.txt",
				kube.GetManagerLogs(data.resources.namespace),
			)
			SaveK8sResources(
				[]string{"deploy", "atlasclusters", "atlasdatabaseusers", "atlasprojects"},
				data.resources.namespace,
			)
		} else {
			GinkgoWriter.Write([]byte("Operator namespace: " + data.resources.namespace))
			Eventually(kube.DeleteNamespace(data.resources.namespace)).Should(Say("deleted"), "Cant delete namespace after testing")
		}
	})

	newData := func(description, path string, users []utils.DBUser) (string, testDataProvider) {
		var data testDataProvider
		data.description = description
		data.confPath = path
		Expect(users).ShouldNot(BeNil())
		data.resources = NewUserInputs("my-atlas-key", users)
		return data.description, data
	}

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test testDataProvider) {
			data = test
			mainCycle(test.confPath, test.resources)
		},
		Entry(newData("Trial - Simplest configuration with no backup and one Admin User",
			"data/atlascluster_basic.yaml",
			[]utils.DBUser{
				*utils.NewDBUser("user1").
					WithSecretRef("dbuser-secret-u1").
					AddRole("readWriteAnyDatabase", "admin", ""),
			},
		)),
		Entry(newData("Almost Production - Backup and 2 users, one Admin and one read-only",
			"data/atlascluster_backup.yaml",
			[]utils.DBUser{
				*utils.NewDBUser("admin").
					WithSecretRef("dbuser-admin-secret-u1").
					AddRole("atlasAdmin", "admin", ""),
				*utils.NewDBUser("user2").
					WithSecretRef("dbuser-secret-u2").
					AddRole("read", "testDB", ""),
			},
		)),
		// Entry(newData("Multiregion, Backup and 2 users", "data/atlascluster_multiregion.yaml",
		// 	append(
		// 		[]utils.DBUser{},
		// 		*utils.NewDBUser("user1").
		// 			WithSecretRef("dbuser-secret-u1").
		// 			AddRole("readWriteAnyDatabase", "admin", ""),
		// 		*utils.NewDBUser("user2").
		// 			WithSecretRef("dbuser-secret-u2").
		// 			AddRole("readWriteAnyDatabase", "admin", ""),
		// 	),
		// )), // TODO CLOUDP-83419
	)
})

func mainCycle(clusterConfigurationFile string, resources userInputs) {
	By("Prepare namespaces and project configuration", func() {
		kube.CreateNamespace(resources.namespace)
		By("Create project spec", func() {
			project := utils.NewProject().
				ProjectName(resources.projectName).
				SecretRef(resources.keyName).
				CompleteK8sConfig(resources.k8sProjectName)
			utils.SaveToFile(FilePathTo(resources.projectName), project)
		})
		By("Create cluster spec", func() {
			resources.clusters = append(resources.clusters, utils.LoadUserClusterConfig(clusterConfigurationFile))
			resources.clusters[0].Spec.Project.Name = resources.k8sProjectName
			resources.clusters[0].Spec.ProviderSettings.InstanceSizeName = "M10"
			utils.SaveToFile(
				resources.clusters[0].ClusterFileName(),
				utils.JSONToYAMLConvert(resources.clusters[0]),
			)
		})
		By("Create dbuser spec", func() {
			Expect(resources.users).ShouldNot(BeNil())
			for _, user := range resources.users {
				user.SaveConfigurationTo(resources.namespace)
				kube.CreateUserSecret(user.Spec.PasswordSecret.Name, resources.namespace)
			}
		})
	})

	By("Create namespaced Operator\n", func() {
		utils.CreateCopyKustomizeNamespace(resources.namespace)
		kube.Apply("-k", "data/"+resources.namespace)
		Eventually(
			kube.GetPodStatus(resources.namespace),
			"5m", "3s",
		).Should(Equal("Running"), "The operator should successfully run")
	})

	By("Create users resources", func() {
		kube.CreateApiKeySecret(resources.keyName, resources.namespace)
		kube.Apply(FilePathTo(resources.projectName), "-n", resources.namespace)
		kube.Apply(resources.clusters[0].ClusterFileName(), "-n", resources.namespace)
		kube.Apply("data"+"/"+resources.namespace+"/user/", "-n", resources.namespace)
	})

	By("Wait project creation", func() {
		waitProject(resources, "1")
		resources.projectID = kube.GetProjectResource(resources.namespace, resources.k8sFullProjectName).Status.ID
	})

	By("Wait cluster creation", func() {
		waitCluster(resources, "1")
	})

	By("check cluster Attribute", func() {
		cluster := mongocli.GetClustersInfo(resources.projectID, resources.clusters[0].Spec.Name)
		compareClustersSpec(resources.clusters[0].Spec, cluster)
	})

	By("check database users Attibutes", func() {
		Eventually(checkIfUsersExist(resources), "2m", "10s").Should(BeTrue())
		checkUsersAttributes(resources)
	})

	By("User can delete Database User", func() {
		// since it is could be several users, we should
		// - delete one user from the list,
		// - delete k8s resource
		// - check Atlas doesn't have the initial user and have the rest
		By("Delete k8s resources")
		kube.Delete(resources.users[0].GetFilePath(resources.namespace), "-n", resources.namespace)
		Eventually(checkIfUserExist(resources.users[0].Spec.Username, resources.projectID)).Should(BeFalse())

		// the rest users should be still there
		resources.users = resources.users[1:]
		Eventually(checkIfUsersExist(resources), "2m", "10s").Should(BeTrue())
		checkUsersAttributes(resources)
	})

	By("Update cluster\n", func() {
		resources.clusters[0].Spec.ProviderSettings.InstanceSizeName = "M20"
		utils.SaveToFile(
			resources.clusters[0].ClusterFileName(),
			utils.JSONToYAMLConvert(resources.clusters[0]),
		)
		kube.Apply(resources.clusters[0].ClusterFileName(), "-n", resources.namespace)
	})

	By("Wait cluster updating", func() {
		waitCluster(resources, "2")
	})

	By("Check attributes", func() {
		uCluster := mongocli.GetClustersInfo(resources.projectID, resources.clusters[0].Spec.Name)
		compareClustersSpec(resources.clusters[0].Spec, uCluster)
	})

	By("Delete cluster", func() {
		kube.Delete(resources.clusters[0].ClusterFileName(), "-n", resources.namespace)
		Eventually(
			checkIfClusterExist(resources),
			"10m", "1m",
		).Should(BeFalse(), "Cluster should be deleted from Atlas")
	})

	By("Delete project", func() {
		kube.Delete(resources.projectPath, "-n", resources.namespace)
		Eventually(
			func() bool {
				return mongocli.IsProjectInfoExist(resources.projectID)
			},
			"5m", "20s",
		).Should(BeFalse(), "Project should be deleted from Atlas")
	})
}
