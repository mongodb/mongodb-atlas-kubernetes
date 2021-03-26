package e2e_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	// . "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	// . "github.com/onsi/gomega/gstruct"

	// v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	// "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/mongocli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

var _ = Describe("[bundle-test] User can", func() {
	var userSpec userInputs
	var imageURL string

	var _ = BeforeEach(func() {
		imageURL = os.Getenv("BUNDLE_IMAGE")
		// imageURL := "docker.io/leori/test-bundles:CLOUDP-84848-bundle-test-252453e" // TODO get it from env . Registry is nessary
	})
	var _ = AfterEach(func() {
		// By("Delete clusters", func() {
		// 	if CurrentGinkgoTestDescription().Failed {
		// 		GinkgoWriter.Write([]byte("userSpec wasn't clean"))
		// 		utils.SaveToFile(
		// 			"output/operator-logs.txt",
		// 			kube.GetManagerLogs(defaultOperatorNS),
		// 		)
		// 		SaveK8suserSpec(
		// 			[]string{"deploy"},
		// 			defaultOperatorNS, //TODO bundles
		// 		)
		// 		SaveK8suserSpec(
		// 			[]string{"atlasclusters", "atlasdatabaseusers", "atlasprojects"},
		// 			userSpec.namespace,
		// 		)
		// 		SaveK8suserSpec(
		// 			[]string{"atlasclusters", "atlasdatabaseusers", "atlasprojects"},
		// 			userSpec.namespace,
		// 		)
		// 	} else {
		// 		// kube.Delete(userSpec.clusters[0].ClusterFileName(), "-n", userSpec.namespace)
		// 		// do not wait it
		// 	}
		// })
	})

	It("User can install", func() {
		// task:
		// operator-sdk olm install
		// # replace the image with the correct one
		// make bundle VERSION=0.4.0 IMG=mongodb/mongodb-atlas-kubernetes-operator:0.4.0
		// docker build -f bundle.Dockerfile -t antonlisovenko/test-bundle:v0.4.0 .
		// docker push antonlisovenko/test-bundle:v0.4.0
		// operator-sdk run bundle docker.io/antonlisovenko/test-bundle:v0.4.0

		// Eventually(cli.Execute("operator-sdk", "olm", "install"), "5m").Should(gexec.Exit(0))

		Eventually(cli.Execute("operator-sdk", "run", "bundle", imageURL), "5m").Should(gexec.Exit(0))

		By("User creates configuration for a new Project and Cluster", func() {
			// user := *utils.NewDBUser("admin").
			// 	WithSecretRef("dbuser-secret-u1").
			// 	AddRole("read", "Ships", "")
			userSpec = NewUserInputs(
				"only-key",
				[]utils.DBUser{
					*utils.NewDBUser("admin").
						WithSecretRef("dbuser-secret-u1").
						AddRole("read", "Ships", ""),
				},
			)

			utils.SaveToFile(
				FilePathTo(userSpec.projectName),
				utils.NewProject().
					ProjectName(userSpec.projectName).
					SecretRef(userSpec.keyName).
					CompleteK8sConfig(userSpec.k8sProjectName),
		 	)
			userSpec.clusters = append(userSpec.clusters, utils.LoadUserClusterConfig(ClusterSample))
			userSpec.clusters[0].Spec.Project.Name = userSpec.k8sProjectName
			userSpec.clusters[0].ObjectMeta.Name = "from-bundle"
			utils.SaveToFile(
				userSpec.clusters[0].ClusterFileName(),
				utils.JSONToYAMLConvert(userSpec.clusters[0]),
			)
		})

		By("Apply configuration", func() {
			kube.CreateNamespace(userSpec.namespace)
			kube.CreateApiKeySecret(userSpec.keyName, userSpec.namespace)
			kube.Apply(FilePathTo(userSpec.projectName), "-n", userSpec.namespace)
			kube.Apply(userSpec.clusters[0].ClusterFileName(), "-n", userSpec.namespace)
		})

		By("Wait project creation", func() {
			waitProject(userSpec, "1")
			userSpec.projectID = kube.GetProjectResource(userSpec.namespace, userSpec.k8sFullProjectName).Status.ID
		})
		//wait cluster + check if user exist
		By("Wait cluster creation", func() {
			waitCluster(userSpec, "1")
		})

		By("check cluster Attribute", func() {
			cluster := mongocli.GetClustersInfo(userSpec.projectID, userSpec.clusters[0].Spec.Name)
			compareClustersSpec(userSpec.clusters[0].Spec, cluster)
		})

		By("check database users Attibutes", func() {
			Eventually(checkIfUsersExist(userSpec), "2m", "10s").Should(BeTrue())
			checkUsersAttributes(userSpec)
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
	})
})
