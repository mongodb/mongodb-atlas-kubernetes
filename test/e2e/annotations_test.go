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
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

var _ = Describe("Annotations base test.", Label("deployment-annotations-ns"), func() {
	var testData *model.TestDataProvider

	AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
			Expect(actions.SaveDeploymentsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		actions.DeleteTestDataDeployments(testData)
		actions.DeleteTestDataProject(testData)
		actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(ctx SpecContext, test func(context.Context) *model.TestDataProvider) {
			testData = test(ctx)
			mainCycle(testData)
		},
		// TODO: fix test for deletion protection on, as it would fail to re-take the cluster after deletion
		Entry("Simple configuration with keep resource policy annotation on deployment", Label("focus-ns-crd"),
			func(ctx context.Context) *model.TestDataProvider {
				return model.DataProvider(ctx, "operator-ns-crd", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 30000, []func(*model.TestDataProvider){
					actions.DeleteDeploymentCRWithKeepAnnotation,
					actions.RedeployDeployment,
					actions.RemoveKeepAnnotation,
				}).WithInitialDeployments(data.CreateDeploymentWithKeepPolicy("atlascluster-annotation")).
					WithProject(data.DefaultProject())
			},
		),
	)
})
