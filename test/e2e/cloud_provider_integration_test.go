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
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/cloudaccess"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
)

const awsRoleNameBase = "atlas-operator-test-aws-role"

var _ = Describe("UserLogin", Label("cloud-access-role"), func() {
	var testData *model.TestDataProvider

	_ = BeforeEach(func() {
		checkUpAWSEnvironment()
	})

	_ = AfterEach(func(ctx SpecContext) {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Clean Roles", func() {
			DeleteAllRoles(ctx, testData)
		})
		By("Delete Resources, Project with Cloud provider access roles", func() {
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(ctx SpecContext, test func(context.Context) *model.TestDataProvider, roles []cloudaccess.Role) {
			testData = test(ctx)
			actions.ProjectCreationFlow(testData)
			cloudAccessRolesFlow(ctx, testData, roles)
		},
		Entry("Test[cloud-access-role-aws-1]: User has project which was updated with AWS custom role", Label("focus-cloud-access-role-aws-1"),
			func(ctx context.Context) *model.TestDataProvider {
				return model.DataProvider(ctx, "cloud-access-role-aws-1", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			[]cloudaccess.Role{
				{
					Name: utils.RandomName(awsRoleNameBase),
					AccessRole: akov2.CloudProviderIntegration{
						ProviderName: "AWS",
						// IamAssumedRoleArn will be filled after role creation
					},
				},
				{
					Name: utils.RandomName(awsRoleNameBase),
					AccessRole: akov2.CloudProviderIntegration{
						ProviderName: "AWS",
						// IamAssumedRoleArn will be filled after role creation
					},
				},
			},
		),
	)
})

func DeleteAllRoles(ctx context.Context, testData *model.TestDataProvider) {
	Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{Name: testData.Project.Name, Namespace: testData.Project.Namespace}, testData.Project)).Should(Succeed())
	errorList := cloudaccess.DeleteCloudProviderIntegrations(ctx, testData.Project.Spec.CloudProviderIntegrations)
	Expect(len(errorList)).Should(Equal(0), errorList)
}

func cloudAccessRolesFlow(ctx context.Context, userData *model.TestDataProvider, roles []cloudaccess.Role) {
	By("Create AWS role", func() {
		err := cloudaccess.CreateRoles(ctx, roles)
		Expect(err).ShouldNot(HaveOccurred())
	})

	By("Create project with cloud access role", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{
			Name:      userData.Project.Name,
			Namespace: userData.Project.Namespace,
		}, userData.Project)).Should(Succeed())
		for _, role := range roles {
			userData.Project.Spec.CloudProviderIntegrations = append(userData.Project.Spec.CloudProviderIntegrations, role.AccessRole)
		}
		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
	})

	By("Establish connection between Atlas and cloud roles", func() {
		Eventually(func(g Gomega) {
			EnsureAllRolesCreated(g, *userData, len(roles))
		}).WithTimeout(5*time.Minute).WithPolling(20*time.Second).Should(Succeed(), "Cloud access roles are not created")

		project := &akov2.AtlasProject{}
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name, Namespace: userData.Project.Namespace}, project)).Should(Succeed())

		err := cloudaccess.AddAtlasStatementToRole(ctx, roles, project.Status.CloudProviderIntegrations)
		Expect(err).ShouldNot(HaveOccurred())

		actions.WaitForConditionsToBecomeTrue(userData, api.CloudProviderIntegrationReadyType, api.ReadyType)
	})

	By("Check cloud access roles status state", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name, Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())
		Expect(userData.Project.Status.CloudProviderIntegrations).Should(HaveLen(len(roles)))
	})
}

func EnsureAllRolesCreated(g Gomega, testData model.TestDataProvider, rolesLen int) {
	project := &akov2.AtlasProject{}
	g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{Name: testData.Project.Name, Namespace: testData.Project.Namespace}, project)).Should(Succeed())
	g.Expect(project.Status.CloudProviderIntegrations).Should(HaveLen(rolesLen))

	for _, role := range project.Status.CloudProviderIntegrations {
		g.Expect(role.Status).Should(BeElementOf([2]string{status.CloudProviderIntegrationStatusCreated, status.CloudProviderIntegrationStatusFailedToAuthorize}))
	}
}
