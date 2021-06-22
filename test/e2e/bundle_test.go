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
			} else {
				actions.AfterEachFinalCleanup([]model.TestDataProvider{data})
			}
		})
	})

	It("User can install", func() {
		Eventually(cli.Execute("operator-sdk", "olm", "install"), "3m").Should(gexec.Exit(0))
		Eventually(cli.Execute("operator-sdk", "run", "bundle", imageURL), "5m").Should(gexec.Exit(0))

		By("User creates configuration for a new Project and Cluster", func() {
			data = model.NewTestDataProvider(
				"bundle-wide",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlascluster_basic.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("reader").
						WithSecretRef("dbuser-secret-u1").
						AddCustomRole(model.RoleCustomReadWrite, "Ships", ""),
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

		By("Delete user resources(project/cluster)", func() {
			actions.DeleteUserResources(&data)
		})
	})
})
