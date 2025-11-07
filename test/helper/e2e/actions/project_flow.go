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

package actions

import (
	"context"
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/deploy"
	helper "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/api/aws"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

func ProjectCreationFlow(userData *model.TestDataProvider) {
	By("Prepare operator configurations", func() {
		r := PrepareOperatorConfigurations(userData)
		ctx := context.Background()
		go func(ctx context.Context) context.Context {
			err := r.Start(ctx)
			Expect(err).NotTo(HaveOccurred())
			return ctx
		}(ctx)
		deploy.CreateProject(userData)
		userData.ManagerContext = ctx
	})
}

func PrepareOperatorConfigurations(userData *model.TestDataProvider) manager.Runnable {
	CreateNamespaceAndSecrets(userData)
	c, err := k8s.BuildCluster(&k8s.Config{
		WatchedNamespaces: map[string]bool{
			userData.Resources.Namespace: true,
		},
		GlobalAPISecret: client.ObjectKey{
			Namespace: userData.Resources.Namespace,
			Name:      config.DefaultOperatorGlobalKey,
		},
		ObjectDeletionProtection:    userData.ObjectDeletionProtection,
		SubObjectDeletionProtection: userData.SubObjectDeletionProtection,
		FeatureFlags:                featureflags.NewFeatureFlags(os.Environ),
	})
	Expect(err).NotTo(HaveOccurred())
	return c
}

func CreateNamespaceAndSecrets(userData *model.TestDataProvider) {
	By(fmt.Sprintf("Create namespace %s", userData.Resources.Namespace))
	Expect(k8s.CreateNamespace(userData.Context, userData.K8SClient, userData.Resources.Namespace)).Should(Succeed())
	k8s.CreateDefaultSecret(userData.Context, userData.K8SClient, config.DefaultOperatorGlobalKey, userData.Resources.Namespace)
	if !userData.Resources.AtlasKeyAccessType.GlobalLevelKey {
		CreateConnectionAtlasKey(userData)
	}
}

func CreateProjectWithCloudProviderAccess(ctx context.Context, testData *model.TestDataProvider, atlasIAMRoleName string) {
	ProjectCreationFlow(testData)

	By("Configure cloud provider access", func() {
		testData.Project.Spec.CloudProviderIntegrations = []akov2.CloudProviderIntegration{
			{
				ProviderName: "AWS",
			},
		}
		Expect(testData.K8SClient.Update(testData.Context, testData.Project)).To(Succeed())

		Eventually(func(g Gomega) bool {
			g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{
				Name:      testData.Project.Name,
				Namespace: testData.Project.Namespace,
			}, testData.Project)).To(Succeed())

			g.Expect(testData.Project.Status.CloudProviderIntegrations).ShouldNot(BeEmpty())

			return testData.Project.Status.CloudProviderIntegrations[0].Status == status.CloudProviderIntegrationStatusCreated
		}).WithTimeout(5 * time.Minute).WithPolling(10 * time.Second).Should(BeTrue())

		roleArn, err := testData.AWSResourcesGenerator.CreateIAMRole(ctx, atlasIAMRoleName, func() helper.IAMPolicy {
			cloudProviderAccess := testData.Project.Status.CloudProviderIntegrations[0]
			return helper.CloudProviderAccessPolicy(cloudProviderAccess.AtlasAWSAccountArn, cloudProviderAccess.AtlasAssumedRoleExternalID)
		})

		Expect(err).Should(BeNil())
		Expect(roleArn).ShouldNot(BeEmpty())

		DeferCleanup(func(ctx SpecContext) {
			Expect(testData.AWSResourcesGenerator.DeleteIAMRole(ctx, atlasIAMRoleName)).To(Succeed())
		})

		testData.Project.Spec.CloudProviderIntegrations[0].IamAssumedRoleArn = roleArn
		Expect(testData.K8SClient.Update(testData.Context, testData.Project)).To(Succeed())

		WaitForConditionsToBecomeTrue(testData, api.CloudProviderIntegrationReadyType, api.ReadyType)
	})
}
