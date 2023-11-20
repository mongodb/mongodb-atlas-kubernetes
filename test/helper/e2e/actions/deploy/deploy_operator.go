// different ways to deploy operator
package deploy

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/api/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/fixtest"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
)

func MultiNamespaceOperator(data *model.TestDataProvider, watchNamespace []string) {
	By("Deploy multinamespaced Operator \n", func() {
		watchNamespaceMap := make(map[string]bool, len(watchNamespace))
		for _, ns := range watchNamespace {
			watchNamespaceMap[ns] = true
		}
		mgr, err := k8s.BuildManager(&k8s.Config{
			Namespace: config.DefaultOperatorNS,
			GlobalAPISecret: client.ObjectKey{
				Namespace: config.DefaultOperatorNS,
				Name:      config.DefaultOperatorGlobalKey,
			},
			WatchedNamespaces: watchNamespaceMap,
		})
		Expect(err).Should(Succeed())
		ctx := context.Background()
		go func(ctx context.Context) {
			err = mgr.Start(ctx)
			Expect(err).Should(Succeed(), "Operator should be started")
		}(ctx)
		data.ManagerContext = ctx
	})
}

func CreateProject(testData *model.TestDataProvider) {
	if testData.Project.GetNamespace() == "" {
		testData.Project.Namespace = testData.Resources.Namespace
	}
	By(fmt.Sprintf("Deploy Project %s", testData.Project.GetName()), func() {
		aClient := atlas.GetClientOrFail()
		err := testData.K8SClient.Create(testData.Context, testData.Project)
		Expect(err).ShouldNot(HaveOccurred(), "Project %s was not created", testData.Project.GetName())
		Eventually(func(g Gomega) {
			// We reported Atlas creating duplicates of a project with the same name
			// See https://jira.mongodb.org/browse/CLOUDP-187749
			// this fix in our tests allows them to automatically fix this issue
			// and thus avoid a flaky failure when this duplicates happens
			g.Expect(fixtest.EnsureNoDuplicates(aClient.Client, fixtest.ZapLoggerFrom(GinkgoWriter), testData.Project.Spec.Name)).ToNot(HaveOccurred())

			condition, _ := k8s.GetProjectStatusCondition(
				testData.Context,
				testData.K8SClient,
				status.ReadyType,
				testData.Resources.Namespace,
				testData.Project.GetName(),
			)
			g.Expect(condition).To(Equal("True"))

			statuses := kube.GetProjectStatus(testData)
			g.Expect(statuses.ID).ToNot(BeEmpty())
		}).WithTimeout(10*time.Minute).WithPolling(10*time.Second).
			Should(Succeed(), "Project %s was not created", testData.Project.GetName())
	})
}

func CreateInitialDeployments(testData *model.TestDataProvider) {
	By("Deploy Initial Deployments", func() {
		for _, deployment := range testData.InitialDeployments {
			if deployment.Namespace == "" {
				deployment.Namespace = testData.Resources.Namespace
				deployment.Spec.Project.Namespace = testData.Resources.Namespace
			}
			err := testData.K8SClient.Create(testData.Context, deployment)
			Expect(err).ShouldNot(HaveOccurred(), fmt.Sprintf("Deployment was not created: %v", deployment))
			Eventually(kube.DeploymentReadyCondition(testData), time.Minute*60, time.Second*5).Should(Equal("True"), "Deployment was not created")
		}
	})
}

func CreateUsers(testData *model.TestDataProvider) {
	By("Deploy Users", func() {
		for _, user := range testData.Users {
			if user.Namespace == "" {
				user.Namespace = testData.Resources.Namespace
				user.Spec.Project.Namespace = testData.Resources.Namespace
			}
			if user.Spec.PasswordSecret != nil {
				secret := utils.UserSecretPassword()
				Expect(k8s.CreateUserSecret(testData.Context, testData.K8SClient, secret,
					user.Spec.PasswordSecret.Name, user.Namespace)).Should(Succeed(),
					"Create user secret failed")
			}
			err := testData.K8SClient.Create(testData.Context, user)
			Expect(err).ShouldNot(HaveOccurred(), fmt.Sprintf("User was not created: %v", user))
			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{Name: user.GetName(), Namespace: user.GetNamespace()}, user))
				for _, condition := range user.Status.Conditions {
					if condition.Type == status.ReadyType {
						g.Expect(condition.Status).Should(Equal(corev1.ConditionTrue), "User should be ready")
					}
				}
			}).WithTimeout(5*time.Minute).WithPolling(20*time.Second).Should(Succeed(), "User was not created")
		}
	})
}
