package e2e_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/cli"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

var _ = Describe("User can deploy operator from bundles", func() {
	var testData *model.TestDataProvider
	var imageURL string

	_ = BeforeEach(func() {
		imageURL = os.Getenv("BUNDLE_IMAGE")
		Expect(imageURL).ShouldNot(BeEmpty(), "SetUP BUNDLE_IMAGE")
	})
	_ = AfterEach(func() {
		By("After each.", func() {
			if CurrentSpecReport().Failed() {
				Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
				Expect(actions.SaveDeploymentsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
			}
			actions.DeleteTestDataDeployments(testData)
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	It("User can install operator with OLM", Label("bundle-test"), func() {
		By("User creates configuration for a new Project and Deployment", func() {
			testData = model.DataProvider(
				"bundle-wide",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				30005,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()).
				WithInitialDeployments(data.CreateBasicDeployment("basic-deployment")).
				WithUsers(data.BasicUser("user1", "user1",
					data.WithSecretRef("dbuser-secret-u1"),
					data.WithCustomRole(string(model.RoleCustomReadWrite), "Ships", "")))
			actions.PrepareUsersConfigurations(testData)
		})

		By("OLM install", func() {
			Eventually(cli.Execute("operator-sdk", "olm", "install"), "3m").Should(gexec.Exit(0))
			Eventually(cli.Execute("operator-sdk", "run", "bundle", imageURL, fmt.Sprintf("--namespace=%s", testData.Resources.Namespace), "--verbose", "--timeout", "30m"), "30m").Should(gexec.Exit(0)) // timeout of operator-sdk is bigger then our default
		})

		By("Apply configuration", func() {
			By(fmt.Sprintf("Create namespace %s", testData.Resources.Namespace))
			Expect(k8s.CreateNamespace(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
			k8s.CreateDefaultSecret(testData.Context, testData.K8SClient, config.DefaultOperatorGlobalKey, testData.Resources.Namespace)
			if !testData.Resources.AtlasKeyAccessType.GlobalLevelKey {
				actions.CreateConnectionAtlasKey(testData)
			}
			deploy.CreateProject(testData)
			By(fmt.Sprintf("project namespace %v", testData.Project.Namespace))
			actions.WaitForConditionsToBecomeTrue(testData, api.ReadyType)
			deploy.CreateInitialDeployments(testData)
			deploy.CreateUsers(testData)
		})
	})
})
