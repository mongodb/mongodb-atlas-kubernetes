// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2e_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/cli"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

var _ = Describe("User can deploy operator from bundles", Label("focus-quatantine-bundle-test"), func() {
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

	It("User can install operator with OLM", func(ctx SpecContext) {
		By("User creates configuration for a new Project and Deployment", func() {
			testData = model.DataProvider(ctx, "bundle-wide", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 30005, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject()).
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
