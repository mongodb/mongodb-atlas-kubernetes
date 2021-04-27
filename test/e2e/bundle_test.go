package e2e_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	actions "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/mongocli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

var _ = Describe("[bundle-test] User can deploy operator from bundles", func() {
	var data model.TestDataProvider
	var imageURL string

	var _ = BeforeEach(func() {
		imageURL = os.Getenv("BUNDLE_IMAGE")
		Expect(imageURL).ShouldNot(BeEmpty(), "SetUP BUNDLE_IMAGE")
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
			} else {
				Eventually(kube.DeleteNamespace(data.Resources.Namespace)).Should(Say("deleted"), "Cant delete namespace after testing")
			}
		})
	})

	It("User can install", func() {
		Eventually(cli.Execute("operator-sdk", "olm", "install"), "5m").Should(gexec.Exit(0))
		Eventually(cli.Execute("operator-sdk", "run", "bundle", imageURL), "5m").Should(gexec.Exit(0))

		By("User creates configuration for a new Project and Cluster", func() {
			data = model.NewTestDataProvider(
				[]string{"data/atlascluster_basic.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("reader").
						WithSecretRef("dbuser-secret-u1").
						AddRole("readWrite", "Ships", ""),
				},
				30005,
				[]func(*model.TestDataProvider){},
			)
			Expect(len(data.Resources.Users)).Should(Equal(1))
			actions.PrepareUsersConfigurations(&data)
		})

		By("Apply configuration", func() {
			actions.DeployUserResourcesAction(&data)
		})

		By("Delete cluster", func() {
			kube.Delete(data.Resources.Clusters[0].ClusterFileName(data.Resources), "-n", data.Resources.Namespace)
			Eventually(
				actions.CheckIfClusterExist(data.Resources),
				"10m", "1m",
			).Should(BeFalse(), "Cluster should be deleted from Atlas")
		})

		By("Delete project", func() {
			kube.Delete(data.Resources.ProjectPath, "-n", data.Resources.Namespace)
			Eventually(
				func() bool {
					return mongocli.IsProjectInfoExist(data.Resources.ProjectID)
				},
				"5m", "20s",
			).Should(BeFalse(), "Project should be deleted from Atlas")
		})
	})
})
