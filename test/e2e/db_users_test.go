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
	"fmt"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/api/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
)

const (
	localSecretName = "local-secret"
)

var _ = Describe("Operator watch all namespace should create connection secrets for database users in any namespace", Label("users", "users-ns"), func() {
	var testData *model.TestDataProvider
	secondNamespace := "second-namespace"

	_ = AfterEach(func() {
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
			Expect(actions.SaveDeploymentsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
			Expect(actions.SaveUsersToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		actions.DeleteTestDataUsers(testData)
		actions.DeleteTestDataProject(testData)
		actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		Expect(k8s.DeleteNamespace(testData.Context, testData.K8SClient, secondNamespace)).Should(Succeed())
	})

	It("Operator run on global namespace", func(ctx SpecContext) {
		By("Setting up test data", func() {
			project := data.DefaultProject()
			project.Namespace = config.DefaultOperatorNS

			deployment := data.CreateBasicDeployment("dbusers-operator-global")
			deployment.Namespace = config.DefaultOperatorNS

			testData = model.DataProvider(ctx, "dbusers-operator-global", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 30008, []func(*model.TestDataProvider){}).WithProject(project).
				WithInitialDeployments(deployment).
				WithUsers(
					data.BasicUser(
						"reader1",
						"reader1",
						data.WithSecretRef("dbuser-secret-u1"),
						data.WithReadWriteRole(),
						data.WithNamespace(config.DefaultOperatorNS),
						data.WithLabels([]common.LabelSpec{
							{Key: "type", Value: "e2e-test"},
							{Key: "context", Value: "cloud"},
						}),
					),
					data.BasicUser(
						"reader2",
						"reader2",
						data.WithProject(project),
						data.WithSecretRef("dbuser-secret-u2"),
						data.WithReadWriteRole(),
						data.WithNamespace(secondNamespace),
						// user 2 access Atlas a local secret
						data.WithCredentials(localSecretName),
					),
				)
			testData.Resources.Namespace = config.DefaultOperatorNS
		})
		By("Running operator watching global namespace", func() {
			Expect(k8s.CreateNamespace(testData.Context, testData.K8SClient, config.DefaultOperatorNS)).To(Succeed())
			k8s.CreateDefaultSecret(testData.Context, testData.K8SClient, config.DefaultOperatorGlobalKey, config.DefaultOperatorNS)
			k8s.CreateDefaultSecret(testData.Context, testData.K8SClient, localSecretName, config.DefaultOperatorNS)
			Expect(k8s.CreateNamespace(testData.Context, testData.K8SClient, secondNamespace)).To(Succeed())
			k8s.CreateDefaultSecret(testData.Context, testData.K8SClient, localSecretName, secondNamespace)

			c, err := k8s.BuildCluster(&k8s.Config{
				GlobalAPISecret: client.ObjectKey{
					Namespace: config.DefaultOperatorNS,
					Name:      config.DefaultOperatorGlobalKey,
				},
				WatchedNamespaces: map[string]bool{
					config.DefaultOperatorNS: true,
					secondNamespace:          true,
				},
				FeatureFlags: featureflags.NewFeatureFlags(func() []string { return []string{} }),
			})
			Expect(err).NotTo(HaveOccurred())

			go func(ctx context.Context) context.Context {
				err := c.Start(ctx)
				Expect(err).NotTo(HaveOccurred())
				return ctx
			}(testData.Context)
		})
		By("Creating the project", func() {
			deploy.CreateProject(testData)
		})
		By("Failing when an user has both project and atlas references are set", func() {
			testData.Users[0].Spec.ExternalProjectRef = &akov2.ExternalProjectReference{
				ID: testData.Project.ID(),
			}

			Expect(testData.K8SClient.Create(testData.Context, testData.Users[0])).ToNot(Succeed())
		})
		By("Creating a linked and a standalone users", func() {
			data.WithExternalProjectRef(testData.Project.ID(), localSecretName)(testData.Users[0])
			deploy.CreateUsers(testData)

			Expect(countConnectionSecrets(testData.K8SClient, testData.Project.Spec.Name)).To(Equal(0))
		})
		By("Create deployment and connection secrets for all related users", func() {
			Expect(testData.K8SClient.Create(testData.Context, testData.InitialDeployments[0])).To(Succeed())

			Eventually(func(g Gomega) bool {
				g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{
					Name:      testData.InitialDeployments[0].Name,
					Namespace: testData.InitialDeployments[0].Namespace,
				}, testData.InitialDeployments[0])).To(Succeed())

				return testData.InitialDeployments[0].Status.StateName == status.StateIDLE
			}).WithTimeout(30 * time.Minute).Should(BeTrue())

			Eventually(func(g Gomega) bool {
				return g.Expect(countConnectionSecrets(testData.K8SClient, testData.Project.Spec.Name)).To(Equal(2))
			}).WithTimeout(5 * time.Minute).WithPolling(10 * time.Second).Should(BeTrue())
		})
		By("Delete deployment and connection secrets for all related users", func() {
			Expect(testData.K8SClient.Delete(testData.Context, testData.InitialDeployments[0])).To(Succeed())

			projectID := testData.Project.Status.ID
			deploymentName := testData.InitialDeployments[0].AtlasName()
			Eventually(func(g Gomega) bool {
				aClient := atlas.GetClientOrFail()
				return aClient.IsDeploymentExist(projectID, deploymentName)
			},
			).WithTimeout(30 * time.Minute).WithPolling(10 * time.Second).Should(BeFalse())

			Eventually(func(g Gomega) bool {
				return g.Expect(countConnectionSecrets(testData.K8SClient, testData.Project.Spec.Name)).To(Equal(0))
			}).WithTimeout(5 * time.Minute).WithPolling(10 * time.Second).Should(BeTrue())
		})
	})
})

var _ = Describe("Operator fails if local credentials is mentioned but unavailable", Label("focus-users", "focus-users-no-creds"), func() {
	var testData *model.TestDataProvider
	namespace := utils.RandomName("namespace")

	_ = AfterEach(func() {
		actions.DeleteTestDataUsers(testData)
		actions.DeleteTestDataProject(testData)
		actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		Expect(k8s.DeleteNamespace(testData.Context, testData.K8SClient, namespace)).Should(Succeed())
	})

	It("Operator run on global namespace to test bogus local credential", func(ctx SpecContext) {
		By("Setting up test data", func() {
			project := data.DefaultProject()
			project.Namespace = namespace

			testData = model.DataProvider(ctx, "dbusers-operator-global", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 30008, []func(*model.TestDataProvider){}).WithProject(project).
				WithUsers(
					data.BasicUser(
						"reader1",
						"reader1",
						data.WithSecretRef("dbuser-secret-u1"),
						data.WithReadWriteRole(),
						data.WithNamespace(namespace),
						data.WithLabels([]common.LabelSpec{
							{Key: "type", Value: "e2e-test"},
							{Key: "context", Value: "cloud"},
						}),
						// user 2 access Atlas a local secret
						data.WithCredentials(localSecretName),
					),
				)
			testData.Resources.Namespace = namespace
		})
		By("Running operator watching global namespace missing credentials", func() {
			Expect(k8s.CreateNamespace(testData.Context, testData.K8SClient, namespace)).NotTo(HaveOccurred())
			k8s.CreateDefaultSecret(testData.Context, testData.K8SClient, config.DefaultOperatorGlobalKey, namespace)

			c, err := k8s.BuildCluster(&k8s.Config{
				GlobalAPISecret: client.ObjectKey{
					Namespace: namespace,
					Name:      config.DefaultOperatorGlobalKey,
				},
				WatchedNamespaces: map[string]bool{
					namespace: true,
				},
				FeatureFlags: featureflags.NewFeatureFlags(func() []string { return []string{} }),
			})
			Expect(err).NotTo(HaveOccurred())

			go func(ctx context.Context) context.Context {
				err := c.Start(ctx)
				Expect(err).NotTo(HaveOccurred())
				return ctx
			}(testData.Context)
		})
		By("Creating project", func() {
			deploy.CreateProject(testData)
		})
		By("Creating user with missing credentials", func() {
			user := testData.Users[0]
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
				g.Expect(user.Status.Conditions).ShouldNot(BeEmpty())
				for _, condition := range user.Status.Conditions {
					if condition.Type == api.ReadyType {
						g.Expect(condition.Status).ShouldNot(Equal(corev1.ConditionTrue), "User should NOT be ready")
					}
					if condition.Type == api.DatabaseUserReadyType {
						g.Expect(condition.Message).Should(ContainSubstring(`Secret "local-secret" not found`))
					}
				}
			}).WithTimeout(2*time.Minute).WithPolling(20*time.Second).Should(Succeed(), "User did not fail as expected")
		})
	})
})

func countConnectionSecrets(k8sClient client.Client, projectName string) int {
	secretList := corev1.SecretList{}
	Expect(k8sClient.List(context.Background(), &secretList)).To(Succeed())

	names := make([]string, 0)
	for _, item := range secretList.Items {
		if strings.HasPrefix(item.Name, kube.NormalizeIdentifier(projectName)) {
			names = append(names, item.Name)
		}
	}

	return len(names)
}
