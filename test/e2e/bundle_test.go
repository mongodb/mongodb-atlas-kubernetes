package e2e_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/mongocli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

var _ = Describe("[bundle-test] User can", func() {
	var userSpec model.UserInputs
	var imageURL string

	var _ = BeforeEach(func() {
		imageURL = os.Getenv("BUNDLE_IMAGE")
		Expect(imageURL).ShouldNot(BeNil())
	})
	var _ = AfterEach(func() {
		By("Atfer each.", func() {
			if CurrentGinkgoTestDescription().Failed {
				GinkgoWriter.Write([]byte("Resources wasn't clean"))
				utils.SaveToFile(
					"output/operator-logs.txt",
					kube.GetManagerLogs(config.DefaultOperatorNS),
				)
				SaveK8sResources(
					[]string{"deploy"},
					"default",
				)
				SaveK8sResources(
					[]string{"atlasclusters", "atlasdatabaseusers", "atlasprojects"},
					userSpec.Namespace,
				)
			} else {
				Eventually(kube.DeleteNamespace(userSpec.Namespace)).Should(Say("deleted"), "Cant delete namespace after testing")
			}
		})
	})

	It("User can install", func() {
		Eventually(cli.Execute("operator-sdk", "olm", "install"), "5m").Should(gexec.Exit(0))
		Eventually(cli.Execute("operator-sdk", "run", "bundle", imageURL), "5m").Should(gexec.Exit(0))

		By("User creates configuration for a new Project and Cluster", func() {
			userSpec = model.NewUserInputs(
				"only-key",
				[]model.DBUser{
					*model.NewDBUser("reader").
						WithSecretRef("dbuser-secret-u1").
						AddRole("read", "Ships", ""),
				},
			)

			utils.SaveToFile(
				userSpec.ProjectPath,
				model.NewProject().
					ProjectName(userSpec.ProjectName).
					SecretRef(userSpec.KeyName).
					CompleteK8sConfig(userSpec.K8sProjectName),
		 	)
			userSpec.Clusters = append(userSpec.Clusters, model.LoadUserClusterConfig(config.ClusterSample))
			userSpec.Clusters[0].Spec.Project.Name = userSpec.K8sProjectName
			userSpec.Clusters[0].ObjectMeta.Name = "cluster-from-bundle"
			utils.SaveToFile(
				userSpec.Clusters[0].ClusterFileName(userSpec),
				utils.JSONToYAMLConvert(userSpec.Clusters[0]),
			)
		})

		By("Apply configuration", func() {
			kube.CreateNamespace(userSpec.Namespace)
			kube.CreateApiKeySecret(userSpec.KeyName, userSpec.Namespace)
			kube.Apply(userSpec.GetResourceFolder(), "-n", userSpec.Namespace)
			for _, user := range userSpec.Users {
				user.SaveConfigurationTo(userSpec.ProjectPath)
				kube.CreateUserSecret(user.Spec.PasswordSecret.Name, userSpec.Namespace)
			}
			kube.Apply(userSpec.GetResourceFolder()+"/user/", "-n", userSpec.Namespace)
		})

		By("Wait project creation", func() {
			waitProject(userSpec, "1")
			userSpec.ProjectID = kube.GetProjectResource(userSpec.Namespace, userSpec.K8sFullProjectName).Status.ID
		})

		By("Wait cluster creation", func() {
			waitCluster(userSpec, "1")
		})

		By("Check attributes", func() {
			uCluster := mongocli.GetClustersInfo(userSpec.ProjectID, userSpec.Clusters[0].Spec.Name)
			compareClustersSpec(userSpec.Clusters[0].Spec, uCluster)
		})

		By("check database users Attibutes", func() {
			Eventually(checkIfUsersExist(userSpec), "2m", "10s").Should(BeTrue())
			checkUsersAttributes(userSpec)
		})

		By("Delete cluster", func() {
			kube.Delete(userSpec.Clusters[0].ClusterFileName(userSpec), "-n", userSpec.Namespace)
			Eventually(
				checkIfClusterExist(userSpec),
				"10m", "1m",
			).Should(BeFalse(), "Cluster should be deleted from Atlas")
		})

		By("Delete project", func() {
			kube.Delete(userSpec.ProjectPath, "-n", userSpec.Namespace)
			Eventually(
				func() bool {
					return mongocli.IsProjectInfoExist(userSpec.ProjectID)
				},
				"5m", "20s",
			).Should(BeFalse(), "Project should be deleted from Atlas")
		})
	})
})
