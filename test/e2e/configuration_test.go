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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

var _ = Describe("Configuration namespaced. Deploy deployment", Label("deployment-ns"), func() {
	var testData *model.TestDataProvider

	AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
			Expect(actions.SaveDeploymentsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
			Expect(actions.SaveUsersToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
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
		Entry("Trial - Simplest configuration with no backup and one Admin User", Label("focus-ns-trial"),
			func(ctx context.Context) *model.TestDataProvider {
				return model.DataProvider(
					ctx,
					"operator-ns-trial",
					model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
					30000,
					[]func(*model.TestDataProvider){
						actions.DeleteFirstUser,
					},
				).WithProject(data.DefaultProject()).
					WithInitialDeployments(data.CreateBasicDeployment("basic-deployment")).
					WithUsers(data.BasicUser("user1", "user1", data.WithSecretRef("dbuser-secret-u1"), data.WithAdminRole()))
			},
		),
		Entry("Almost Production - Backup and 2 DB users: one Admin and one read-only", Label("focus-ns-backup2db", "focus-long-run"),
			func(ctx context.Context) *model.TestDataProvider {
				return model.DataProvider(
					ctx,
					"operator-ns-prodlike",
					model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
					30001,
					[]func(*model.TestDataProvider){
						actions.UpdateSpecOfSelectedDeployment(data.NewDeploymentWithBackupSpec(), 0),
						actions.SuspendDeployment,
						actions.ReactivateDeployment,
						actions.DeleteFirstUser,
					},
				).WithProject(data.DefaultProject()).
					WithInitialDeployments(data.CreateDeploymentWithBackup("backup-deployment")).
					WithUsers(
						data.BasicUser("admin", "user1", data.WithSecretRef("dbuser-secret-u1"), data.WithAdminRole()),
						data.BasicUser("user2", "user2", data.WithSecretRef("dbuser-secret-u2"), data.WithCustomRole(string(model.RoleCustomReadWrite), "Ships", "readWrite")),
					)
			},
		),
		Entry("Multiregion AWS, Backup and 2 DBUsers", Label("focus-ns-multiregion-aws-2"),
			func(ctx context.Context) *model.TestDataProvider {
				return model.DataProvider(
					ctx,
					"operator-ns-multiregion-aws",
					model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
					30003,
					[]func(*model.TestDataProvider){
						actions.SuspendDeployment,
						actions.ReactivateDeployment,
						actions.DeleteFirstUser,
					},
				).WithProject(data.DefaultProject()).
					WithInitialDeployments(data.CreateDeploymentWithMultiregionAWS("multiregion-aws-deployment")).
					WithUsers(data.BasicUser("user1", "user1", data.WithSecretRef("dbuser-secret-u1"), data.WithAdminRole()),
						data.BasicUser("user2", "user2", data.WithSecretRef("dbuser-secret-u2"), data.WithAdminRole()))
			},
		),
		Entry("Multiregion Azure, Backup and 1 DBUser", Label("focus-ns-multiregion-azure-1"),
			func(ctx context.Context) *model.TestDataProvider {
				return model.DataProvider(
					ctx,
					"operator-multiregion-azure",
					model.NewEmptyAtlasKeyType().UseDefaultFullAccess().CreateAsGlobalLevelKey(),
					30012,
					[]func(*model.TestDataProvider){
						actions.DeleteFirstUser,
					},
				).WithProject(data.DefaultProject()).
					WithInitialDeployments(data.CreateDeploymentWithMultiregionAzure("multiregion-azure-deployment")).
					WithUsers(data.BasicUser("user1", "user1", data.WithSecretRef("dbuser-secret-azure"), data.WithAdminRole()))
			},
		),
		Entry("Multiregion GCP, Backup and 1 DBUser", Label("focus-ns-multiregion-gcp-1"),
			func(ctx context.Context) *model.TestDataProvider {
				return model.DataProvider(
					ctx,
					"operator-multiregion-gcp",
					model.NewEmptyAtlasKeyType().UseDefaultFullAccess().CreateAsGlobalLevelKey(),
					30013,
					[]func(*model.TestDataProvider){
						actions.DeleteFirstUser,
					},
				).WithProject(data.DefaultProject()).
					WithInitialDeployments(data.CreateDeploymentWithMultiregionGCP("multiregion-gcp-deployment")).
					WithUsers(data.BasicUser("user1", "user1", data.WithSecretRef("dbuser-secret-gcp"), data.WithAdminRole()))
			},
		),
		Entry("Product Owner - Simplest configuration with ProjectOwner and update deployment to have backup", Label("focus-ns-owner", "long-run"),
			func(ctx context.Context) *model.TestDataProvider {
				return model.DataProvider(
					ctx,
					"operator-ns-product-owner",
					model.NewEmptyAtlasKeyType().WithRoles([]model.AtlasRoles{model.GroupOwner}).WithWhiteList([]string{"0.0.0.1/1", "128.0.0.0/1"}),
					30010,
					[]func(*model.TestDataProvider){
						actions.UpdateSpecOfSelectedDeployment(data.NewDeploymentWithBackupSpec(), 0),
					},
				).WithProject(data.DefaultProject()).
					WithInitialDeployments(data.CreateDeploymentWithBackup("backup-deployment-owner")).
					WithUsers(
						data.BasicUser("user1", "user1", data.WithSecretRef("dbuser-secret-u1"), data.WithAdminRole()),
					)
			},
		),
		Entry("Trial - Global connection", Label("focus-ns-global-key"),
			func(ctx context.Context) *model.TestDataProvider {
				return model.DataProvider(
					ctx,
					"operator-ns-trial-global",
					model.NewEmptyAtlasKeyType().UseDefaultFullAccess().CreateAsGlobalLevelKey(),
					30011,
					[]func(*model.TestDataProvider){
						actions.DeleteFirstUser,
					},
				).WithProject(data.DefaultProject()).
					WithInitialDeployments(data.CreateBasicDeployment("trial")).
					WithUsers(
						data.BasicUser("user1", "user1", data.WithSecretRef("dbuser-secret-u1"), data.WithAdminRole()),
					)
			},
		),
		Entry("Free - Users can use M0, default key", Label("focus-ns-m0"),
			func(ctx context.Context) *model.TestDataProvider {
				return model.DataProvider(
					ctx,
					"operator-ns-free",
					model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
					30016,
					[]func(*model.TestDataProvider){
						actions.DeleteFirstUser,
					},
				).WithProject(data.DefaultProject()).
					WithInitialDeployments(data.CreateFreeAdvancedDeployment("basic-free-deployment")).
					WithUsers(data.BasicUser("user", "user1", data.WithSecretRef("dbuser-secret"), data.WithAdminRole()))
			},
		),
		Entry("Free - Users can use M0, global", Label("focus-ns-global-key-m0"),
			func(ctx context.Context) *model.TestDataProvider {
				return model.DataProvider(
					ctx,
					"operator-ns-free",
					model.NewEmptyAtlasKeyType().UseDefaultFullAccess().CreateAsGlobalLevelKey(),
					30017,
					[]func(*model.TestDataProvider){
						actions.DeleteFirstUser,
					},
				).WithProject(data.DefaultProject()).
					WithInitialDeployments(data.CreateFreeAdvancedDeployment("basic-free-deployment")).
					WithUsers(data.BasicUser("user", "user1", data.WithSecretRef("dbuser-secret"), data.WithAdminRole()))
			},
		),
	)
})

func mainCycle(testData *model.TestDataProvider) {
	r := actions.PrepareOperatorConfigurations(testData)
	ctx := context.Background()
	go func(ctx context.Context) {
		err := r.Start(ctx)
		Expect(err).NotTo(HaveOccurred())
	}(ctx)

	By("Deploy User Resources", func() {
		deploy.CreateProject(testData)
		deploy.CreateInitialDeployments(testData)
		deploy.CreateUsers(testData)
	})

	By("Additional check for the current data set", func() {
		for _, check := range testData.Actions {
			check(testData)
		}
	})
}
