package e2e_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	. "github.com/onsi/gomega/gbytes"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/mongocli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

// TODO add checks
type testDataProvider struct {
	description string
	confPath    string
	u           userInputs
}

var _ = Describe("[cluster-ns] Configuration namespaced. Deploy cluster", func() {

	var data testDataProvider

	var _ = AfterEach(func() {
		if CurrentGinkgoTestDescription().Failed {
			GinkgoWriter.Write([]byte("Resources wasn't clean\n"))
			utils.SaveToFile(
				"output/operator-logs.txt",
				kube.GetManagerLogs(data.u.namespace),
			)
			SaveK8sResources(
				[]string{"deploy", "atlasclusters", "atlasdatabaseusers", "atlasprojects"},
				data.u.namespace,
			)
		} else {
			GinkgoWriter.Write([]byte("Operator namespace: " + data.u.namespace))
			Eventually(kube.DeleteNamespace(data.u.namespace)).Should(Say("deleted"), "Cant delete namespace after testing")
		}
	})

	newData := func(description, path string, users []utils.DBUser) (string, testDataProvider) {
		var data testDataProvider
		data.description = description
		data.confPath = path
		Expect(users).ShouldNot(BeNil())
		data.u = NewUserInputs("my-atlas-key", users)
		return data.description, data
	}

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test testDataProvider) {
			data = test
			mainCycle(test.confPath, test.u)
		},
		Entry(newData("Trial - Simplest configuration with no backup and one Admin User",
			"data/atlascluster_basic.yaml",
			append(
				[]utils.DBUser{},
				*utils.NewDBUser("user1").
					WithSecretRef("dbuser-secret-u1").
					AddRole("readWriteAnyDatabase", "admin", ""),
			),
		)),
		Entry(newData("Almost Production - Backup and 2 users, one Admin and one read-only", "data/atlascluster_backup.yaml",
			append(
				[]utils.DBUser{},
				*utils.NewDBUser("admin").
					WithSecretRef("dbuser-admin-secret-u1").
					AddRole("atlasAdmin", "admin", ""),
				*utils.NewDBUser("user2").
					WithSecretRef("dbuser-secret-u2").
					AddRole("read", "testDB", ""),
			),
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

func mainCycle(clusterConfigurationFile string, u userInputs) {
	By("Prepare namespaces and project configuration", func() {
		kube.CreateNamespace(u.namespace)
		By("Create project spec", func() {
			project := utils.NewProject().
				ProjectName(u.projectName).
				SecretRef(u.keyName).
				CompleteK8sConfig(u.k8sProjectName)
			utils.SaveToFile(FilePathTo(u.projectName), project)
		})
		By("Create cluster spec", func() {
			u.clusters = append(u.clusters, utils.LoadUserClusterConfig(clusterConfigurationFile))
			u.clusters[0].Spec.Project.Name = u.k8sProjectName
			u.clusters[0].Spec.ProviderSettings.InstanceSizeName = "M10"
			utils.SaveToFile(
				u.clusters[0].ClusterFileName(),
				utils.JSONToYAMLConvert(u.clusters[0]),
			)
		})
		By("Create dbuser spec", func() {
			Expect(u.users).ShouldNot(BeNil())
			for _, user := range u.users {
				user.SaveConfigurationTo(u.namespace)
				kube.CreateUserSecret(user.Spec.PasswordSecret.Name, u.namespace)
			}
		})
	})

	By("Create namespaced Operator\n", func() {
		utils.CreateCopyKustomizeNamespace(u.namespace)
		kube.Apply("-k", "data/"+u.namespace)
		Eventually(
			kube.GetPodStatus(u.namespace),
			"5m", "3s",
		).Should(Equal("Running"), "The operator should successfully run")
	})

	By("Create users resources", func() {
		kube.CreateApiKeySecret(u.keyName, u.namespace)
		kube.Apply(FilePathTo(u.projectName), "-n", u.namespace)
		kube.Apply(u.clusters[0].ClusterFileName(), "-n", u.namespace)
		kube.Apply("data"+"/"+u.namespace+"/user/", "-n", u.namespace)
	})

	By("Wait project creation", func() {
		waitProject(u, "1")
		u.projectID = kube.GetProjectResource(u.namespace, u.k8sFullProjectName).Status.ID
	})

	By("Wait cluster creation", func() {
		waitCluster(u, "1")
	})

	By("check cluster Attribute", func() {
		cluster := mongocli.GetClustersInfo(u.projectID, u.clusters[0].Spec.Name)
		compareClustersSpec(u.clusters[0].Spec, cluster)
	})

	By("check database users Attibutes", func() {
		Eventually(checkIfUsersExist(u), "2m", "10s").Should(BeTrue())
		for _, user := range u.users {
			atlasUser := mongocli.GetUser(user.Spec.Username, u.projectID)
			// Required fields
			Expect(atlasUser).To(MatchFields(IgnoreExtras, Fields{
				"Username":     Equal(user.Spec.Username),
				"GroupID":      Equal(u.projectID),
				"DatabaseName": Or(Equal(user.Spec.DatabaseName), Equal("admin")),
			}), "Users attributes should be the same as requested by the user")

			for i, role := range atlasUser.Roles {
				Expect(role).To(MatchFields(IgnoreMissing, Fields{
					"RoleName":       Equal(user.Spec.Roles[i].RoleName),
					"DatabaseName":   Equal(user.Spec.Roles[i].DatabaseName),
					"CollectionName": Equal(user.Spec.Roles[i].CollectionName),
				}))
			}
		}
	})

	By("User can delete User", func() {
		// since it is could be several users, we should
		// - delete one user from the list,
		// - delete k8s resource
		// - check Atlas doesn't have the initial user and have the rest
		By("Delete k8s resources")
		kube.Delete(u.users[0].GetFilePath(u.namespace), "-n", u.namespace)
		Eventually(checkIfUserExist(u.users[0].Spec.Username, u.projectID)).Should(BeFalse())

		// the rest users should be still there
		u.users = u.users[1:]
		Eventually(checkIfUsersExist(u), "2m", "10s").Should(BeTrue())
		for _, user := range u.users {
			atlasUser := mongocli.GetUser(user.Spec.Username, u.projectID)
			// Requared fields
			Expect(atlasUser).To(MatchFields(IgnoreExtras, Fields{
				"Username":     Equal(user.Spec.Username),
				"GroupID":      Equal(u.projectID),
				"DatabaseName": Or(Equal(user.Spec.DatabaseName), Equal("admin")),
			}), "Users attributes should be the same as requested by the user")

			for i, role := range atlasUser.Roles {
				Expect(role).To(MatchFields(IgnoreMissing, Fields{
					"RoleName":       Equal(user.Spec.Roles[i].RoleName),
					"DatabaseName":   Equal(user.Spec.Roles[i].DatabaseName),
					"CollectionName": Equal(user.Spec.Roles[i].CollectionName),
				}))
			}
		}
	})

	By("Update cluster\n", func() {
		u.clusters[0].Spec.ProviderSettings.InstanceSizeName = "M20"
		utils.SaveToFile(
			u.clusters[0].ClusterFileName(),
			utils.JSONToYAMLConvert(u.clusters[0]),
		)
		kube.Apply(u.clusters[0].ClusterFileName(), "-n", u.namespace)
	})

	By("Wait cluster updating", func() {
		waitCluster(u, "2")
	})

	By("Check attributes", func() {
		uCluster := mongocli.GetClustersInfo(u.projectID, u.clusters[0].Spec.Name)
		compareClustersSpec(u.clusters[0].Spec, uCluster)
	})

	By("Delete cluster", func() {
		kube.Delete(u.clusters[0].ClusterFileName(), "-n", u.namespace)
		Eventually(
			checkIfClusterExist(u),
			"10m", "1m",
		).Should(BeFalse(), "Cluster should be deleted from Atlas")
	})

	By("Delete project", func() {
		kube.Delete(u.projectPath, "-n", u.namespace)
		Eventually(
			func() bool {
				return mongocli.IsProjectInfoExist(u.projectID)
			},
			"5m", "20s",
		).Should(BeFalse(), "Project should be deleted from Atlas")
	})
}
