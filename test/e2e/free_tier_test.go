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
	"go.mongodb.org/atlas-sdk/v20250312006/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/api/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

var _ = Describe("Free tier", Label("free-tier"), func() {
	var testData *model.TestDataProvider

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Free tier test\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
			Expect(actions.SaveDeploymentsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Delete Resources", func() {
			actions.DeleteTestDataDeployments(testData)
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	DescribeTable("Operator should support exported CR for free tier deployments",
		func(ctx SpecContext, test func(context.Context) *model.TestDataProvider) {
			testData = test(ctx)
			actions.ProjectCreationFlow(testData)
			freeTierDeploymentFlow(testData)
		},
		Entry("Test free tier deployment",
			func(ctx SpecContext) *model.TestDataProvider {
				return model.DataProvider(ctx, "free-tier", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject()).WithInitialDeployments(data.CreateFreeAdvancedDeployment("free-tier"))
			},
		),
		Entry("Test free tier advanced deployment",
			func(ctx SpecContext) *model.TestDataProvider {
				return model.DataProvider(ctx, "free-tier-advanced", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject()).WithInitialDeployments(data.CreateFreeAdvancedDeployment("free-tier"))
			},
		),
	)
})

func freeTierDeploymentFlow(userData *model.TestDataProvider) {
	By("Create free cluster in Atlas", func() {
		aClient := atlas.GetClientOrFail()
		Expect(userData.InitialDeployments).Should(HaveLen(1))
		name := userData.InitialDeployments[0].GetDeploymentName()
		_, _, err := aClient.Client.ClustersApi.
			CreateCluster(
				userData.Context,
				userData.Project.ID(),
				&admin.ClusterDescription20240805{
					Name:        &name,
					ClusterType: pointer.MakePtr("REPLICASET"),
					ReplicationSpecs: &[]admin.ReplicationSpec20240805{
						{
							ZoneName: pointer.MakePtr("Zone 1"),
							RegionConfigs: &[]admin.CloudRegionConfig20240805{
								{
									ProviderName:        pointer.MakePtr("TENANT"),
									BackingProviderName: pointer.MakePtr("AWS"),
									Priority:            pointer.MakePtr(7),
									RegionName:          pointer.MakePtr("US_EAST_1"),
									ElectableSpecs: &admin.HardwareSpec20240805{
										InstanceSize: pointer.MakePtr(data.InstanceSizeM0),
										NodeCount:    pointer.MakePtr(3),
									},
								},
							},
						},
					},
				},
			).Execute()
		Expect(err).ShouldNot(HaveOccurred())
	})

	By("Apply deployment CR", func() {
		deploy.CreateInitialDeployments(userData)
	})
}
