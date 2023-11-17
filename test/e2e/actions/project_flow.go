package actions

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/types"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/model"
)

func ProjectCreationFlow(userData *model.TestDataProvider) {
	By("Prepare operator configurations", func() {
		mgr := PrepareOperatorConfigurations(userData)
		ctx := context.Background()
		go func(ctx context.Context) context.Context {
			err := mgr.Start(ctx)
			Expect(err).NotTo(HaveOccurred())
			return ctx
		}(ctx)
		deploy.CreateProject(userData)
		userData.ManagerContext = ctx
	})
}

func PrepareOperatorConfigurations(userData *model.TestDataProvider) manager.Manager {
	CreateNamespaceAndSecrets(userData)
	mgr, err := k8s.BuildManager(&k8s.Config{
		Namespace: userData.Resources.Namespace,
		WatchedNamespaces: map[string]bool{
			userData.Resources.Namespace: true,
		},
		GlobalAPISecret: client.ObjectKey{
			Namespace: userData.Resources.Namespace,
			Name:      config.DefaultOperatorGlobalKey,
		},
		ObjectDeletionProtection:    userData.ObjectDeletionProtection,
		SubObjectDeletionProtection: userData.SubObjectDeletionProtection,
	})
	Expect(err).NotTo(HaveOccurred())
	return mgr
}

func CreateNamespaceAndSecrets(userData *model.TestDataProvider) {
	By(fmt.Sprintf("Create namespace %s", userData.Resources.Namespace))
	Expect(k8s.CreateNamespace(userData.Context, userData.K8SClient, userData.Resources.Namespace)).Should(Succeed())
	k8s.CreateDefaultSecret(userData.Context, userData.K8SClient, config.DefaultOperatorGlobalKey, userData.Resources.Namespace)
	if !userData.Resources.AtlasKeyAccessType.GlobalLevelKey {
		CreateConnectionAtlasKey(userData)
	}
}

func CreateProjectWithCloudProviderAccess(testData *model.TestDataProvider, atlasIAMRoleName string) {
	ProjectCreationFlow(testData)

	By("Configure cloud provider access", func() {
		testData.Project.Spec.CloudProviderAccessRoles = []mdbv1.CloudProviderAccessRole{
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

			g.Expect(testData.Project.Status.CloudProviderAccessRoles).ShouldNot(BeEmpty())

			return testData.Project.Status.CloudProviderAccessRoles[0].Status == status.CloudProviderAccessStatusCreated
		}).WithTimeout(5 * time.Minute).WithPolling(10 * time.Second).Should(BeTrue())

		roleArn, err := testData.AWSResourcesGenerator.CreateIAMRole(atlasIAMRoleName, func() helper.IAMPolicy {
			cloudProviderAccess := testData.Project.Status.CloudProviderAccessRoles[0]
			return helper.CloudProviderAccessPolicy(cloudProviderAccess.AtlasAWSAccountArn, cloudProviderAccess.AtlasAssumedRoleExternalID)
		})

		Expect(err).Should(BeNil())
		Expect(roleArn).ShouldNot(BeEmpty())

		testData.AWSResourcesGenerator.Cleanup(func() {
			Expect(testData.AWSResourcesGenerator.DeleteIAMRole(atlasIAMRoleName)).To(Succeed())
		})

		testData.Project.Spec.CloudProviderAccessRoles[0].IamAssumedRoleArn = roleArn
		Expect(testData.K8SClient.Update(testData.Context, testData.Project)).To(Succeed())

		WaitForConditionsToBecomeTrue(testData, status.CloudProviderAccessReadyType, status.ReadyType)
	})
}
