package e2e_test

import (
	"fmt"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	. "github.com/onsi/gomega/gbytes"

	appclient "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/appclient"
	helm "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/helm"
	kube "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kube"
	mongocli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/mongocli"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

// TODO add checks
type testDataProvider struct {
	description string
	confPath    string
	resources   model.UserInputs
}

var _ = Describe("[cluster-ns] Configuration namespaced. Deploy cluster", func() {

	var data testDataProvider

	var _ = AfterEach(func() {
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + data.resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentGinkgoTestDescription().Failed {
			GinkgoWriter.Write([]byte("Resources wasn't clean\n"))
			utils.SaveToFile(
				"output/operator-logs.txt",
				kube.GetManagerLogs(data.resources.Namespace),
			)
			SaveK8sResources(
				[]string{"deploy", "atlasclusters", "atlasdatabaseusers", "atlasprojects"},
				data.resources.Namespace,
			)
		} else {
			Eventually(kube.DeleteNamespace(data.resources.Namespace)).Should(Say("deleted"), "Cant delete namespace after testing")
		}

	})

	// TODO remove portGroup (nodePort for the app)
	newData := func(description, path string, users []model.DBUser, portGroup int) (string, testDataProvider, int) {
		var data testDataProvider
		data.description = description
		data.confPath = path
		Expect(users).ShouldNot(BeNil())
		data.resources = model.NewUserInputs("my-atlas-key", users)
		return data.description, data, portGroup
	}

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test testDataProvider, portGroup int) {
			data = test
			mainCycle(test.confPath, test.resources, portGroup)
		},
		Entry(newData("Trial - Simplest configuration with no backup and one Admin User",
			"data/atlascluster_basic.yaml",
			[]model.DBUser{
				*model.NewDBUser("user1").
					WithSecretRef("dbuser-secret-u1").
					AddRole("readWriteAnyDatabase", "admin", ""),
			},
			30000,
		)),
		Entry(newData("Almost Production - Backup and 2 users, one Admin and one read-only",
			"data/atlascluster_backup.yaml",
			[]model.DBUser{
				*model.NewDBUser("admin").
					WithSecretRef("dbuser-admin-secret-u1").
					AddRole("atlasAdmin", "admin", ""),
				*model.NewDBUser("user2").
					WithSecretRef("dbuser-secret-u2").
					AddRole("readWrite", "Ships", ""),
			},
			30002,
		)),
		// Entry(newData("Multiregion, Backup and 2 users", "data/atlascluster_multiregion.yaml",
		// 	append(
		// 		[]utils.DBUser{},
		// 		*utils.NewDBUser("user1").
		// 			WithSecretRef("dbuser-secret-u1").
		// 			AddRole("atlasAdmin", "admin", ""),
		// 		*utils.NewDBUser("user2").
		// 			WithSecretRef("dbuser-secret-u2").
		// 			AddRole("atlasAdmin", "admin", ""),
		// 	),
		// )), // TODO CLOUDP-83419
	)
})

func mainCycle(clusterConfigurationFile string, resources model.UserInputs, portGroup int) {
	By("Prepare namespaces and project configuration", func() {
		kube.CreateNamespace(resources.Namespace)
		By("Create project spec", func() {
			project := model.NewProject().
				ProjectName(resources.ProjectName).
				SecretRef(resources.KeyName).
				WithIpAccess("0.0.0.0/0", "everyone").
				CompleteK8sConfig(resources.K8sProjectName)
			GinkgoWriter.Write([]byte(resources.ProjectPath + "\n"))
			utils.SaveToFile(resources.ProjectPath, project)
		})
		By("Create cluster spec", func() {
			resources.Clusters = append(resources.Clusters, model.LoadUserClusterConfig(clusterConfigurationFile))
			resources.Clusters[0].Spec.Project.Name = resources.K8sProjectName
			resources.Clusters[0].Spec.ProviderSettings.InstanceSizeName = "M10"
			utils.SaveToFile(
				resources.Clusters[0].ClusterFileName(resources),
				utils.JSONToYAMLConvert(resources.Clusters[0]),
			)
		})
		By("Create dbuser spec", func() {
			Expect(resources.Users).ShouldNot(BeNil())
			for _, user := range resources.Users {
				user.SaveConfigurationTo(resources.ProjectPath)
				kube.CreateUserSecret(user.Spec.PasswordSecret.Name, resources.Namespace)
			}
		})
	})

	By("Create namespaced Operator\n", func() {
		CopyKustomizeNamespaceOperator(resources)
		// CreateCopyKustomizeNamespace(resources.namespace)
		kube.Apply("-k", resources.GetOperatorFolder())
		Eventually(
			kube.GetPodStatus(resources.Namespace),
			"5m", "3s",
		).Should(Equal("Running"), "The operator should successfully run")
	})

	By("Create users resources", func() {
		kube.CreateApiKeySecret(resources.KeyName, resources.Namespace)
		kube.Apply(resources.ProjectPath, "-n", resources.Namespace)
		kube.Apply(resources.Clusters[0].ClusterFileName(resources), "-n", resources.Namespace)
		kube.Apply(resources.GetResourceFolder()+"/user/", "-n", resources.Namespace)
	})

	By("Wait project creation", func() {
		waitProject(resources, "1")
		resources.ProjectID = kube.GetProjectResource(resources.Namespace, resources.K8sFullProjectName).Status.ID
	})

	By("Wait cluster creation", func() {
		waitCluster(resources, "1")
	})

	By("check cluster Attribute", func() {
		cluster := mongocli.GetClustersInfo(resources.ProjectID, resources.Clusters[0].Spec.Name)
		compareClustersSpec(resources.Clusters[0].Spec, cluster)
	})

	By("check database users Attibutes", func() {
		Eventually(checkIfUsersExist(resources), "2m", "10s").Should(BeTrue())
		checkUsersAttributes(resources)
	})

	By("Deploy application for user", func() {
		// 	// kube apply application
		// 	// send ddata
		// 	// retrieve data
		for i, user := range resources.Users { // TODO in parallel(?)
			// data
			port := strconv.Itoa(i + portGroup)
			key := port
			data := fmt.Sprintf("{\"key\":\"%s\",\"shipmodel\":\"heavy\",\"hp\":150}", key)

			helm.InstallTestApplication(resources, user, port)
			waitTestApplication(resources.Namespace, "app=test-app-"+user.Spec.Username)

			app := appclient.NewTestAppClient(port)
			Expect(app.Get("")).Should(Equal("It is working"))
			Expect(app.Post(data)).ShouldNot(HaveOccurred())
			Expect(app.Get("/mongo/" + key)).Should(Equal(data))
		}
	})

	By("Update cluster\n", func() {
		resources.Clusters[0].Spec.ProviderSettings.InstanceSizeName = "M20"
		utils.SaveToFile(
			resources.Clusters[0].ClusterFileName(resources),
			utils.JSONToYAMLConvert(resources.Clusters[0]),
		)
		kube.Apply(resources.Clusters[0].ClusterFileName(resources), "-n", resources.Namespace)
	})

	By("Wait cluster updating", func() {
		waitCluster(resources, "2")
	})

	By("Check attributes", func() {
		uCluster := mongocli.GetClustersInfo(resources.ProjectID, resources.Clusters[0].Spec.Name)
		compareClustersSpec(resources.Clusters[0].Spec, uCluster)
	})

	By("Check user data still in the cluster", func() {
		for i := range resources.Users { // TODO in parallel(?)
			port := strconv.Itoa(i + portGroup)
			key := port
			data := fmt.Sprintf("{\"key\":\"%s\",\"shipmodel\":\"heavy\",\"hp\":150}", key)
			app := appclient.NewTestAppClient(port)
			Expect(app.Get("/mongo/" + key)).Should(Equal(data))
		}
	})

	By("User can delete Database User", func() {
		// since it is could be several users, we should
		// - delete k8s resource
		// - delete one user from the list,
		// - check Atlas doesn't have the initial user and have the rest
		By("Delete k8s resources")
		Eventually(kube.Delete(resources.GetResourceFolder()+"/user/user-"+resources.Users[0].ObjectMeta.Name+".yaml", "-n", resources.Namespace)).Should(Say("deleted"))
		Eventually(checkIfUserExist(resources.Users[0].Spec.Username, resources.ProjectID)).Should(BeFalse())

		// the rest users should be still there
		resources.Users = resources.Users[1:]
		Eventually(checkIfUsersExist(resources), "2m", "10s").Should(BeTrue())
		checkUsersAttributes(resources)
	})

	By("Delete cluster", func() {
		kube.Delete(resources.Clusters[0].ClusterFileName(resources), "-n", resources.Namespace)
		Eventually(
			checkIfClusterExist(resources),
			"10m", "1m",
		).Should(BeFalse(), "Cluster should be deleted from Atlas")
	})

	By("Delete project", func() {
		kube.Delete(resources.ProjectPath, "-n", resources.Namespace)
		Eventually(
			func() bool {
				return mongocli.IsProjectInfoExist(resources.ProjectID)
			},
			"5m", "20s",
		).Should(BeFalse(), "Project should be deleted from Atlas")
	})
}
