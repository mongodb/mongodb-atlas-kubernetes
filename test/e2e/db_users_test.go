package e2e

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

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/api/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/resources"
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

	It("Operator run on global namespace", func() {
		By("Setting up test data", func() {
			project := data.DefaultProject()
			project.Namespace = config.DefaultOperatorNS

			deployment := data.CreateBasicDeployment("dbusers-operator-global")
			deployment.Namespace = config.DefaultOperatorNS

			testData = model.DataProvider(
				"dbusers-operator-global",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				30008,
				[]func(*model.TestDataProvider){},
			).WithProject(project).
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
			Expect(k8s.CreateNamespace(testData.Context, testData.K8SClient, config.DefaultOperatorNS)).NotTo(HaveOccurred())
			k8s.CreateDefaultSecret(testData.Context, testData.K8SClient, config.DefaultOperatorGlobalKey, config.DefaultOperatorNS)
			Expect(k8s.CreateNamespace(testData.Context, testData.K8SClient, secondNamespace)).NotTo(HaveOccurred())
			k8s.CreateDefaultSecret(testData.Context, testData.K8SClient, localSecretName, secondNamespace)

			mgr, err := k8s.BuildManager(&k8s.Config{
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
				err := mgr.Start(ctx)
				Expect(err).NotTo(HaveOccurred())
				return ctx
			}(testData.Context)
		})
		By("Creating project and database users resources", func() {
			deploy.CreateProject(testData)
			deploy.CreateUsers(testData)

			Eventually(func(g Gomega) bool {
				for i := range testData.Users {
					dbUser := testData.Users[i]

					g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(dbUser), dbUser)).To(Succeed())
					g.Expect(resources.CheckCondition(testData.K8SClient, dbUser, api.TrueCondition(api.ReadyType))).To(BeTrue())
				}

				return true
			}).WithTimeout(5 * time.Minute).WithPolling(10 * time.Second).Should(BeTrue())

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

var _ = Describe("Operator fails if local credentials is mentioned but unavailable", Label("users", "users-no-creds"), func() {
	var testData *model.TestDataProvider
	namespace := utils.RandomName("namespace")

	_ = AfterEach(func() {
		actions.DeleteTestDataUsers(testData)
		actions.DeleteTestDataProject(testData)
		actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		Expect(k8s.DeleteNamespace(testData.Context, testData.K8SClient, namespace)).Should(Succeed())
	})

	It("Operator run on global namespace to test bogus local credential", func() {
		By("Setting up test data", func() {
			project := data.DefaultProject()
			project.Namespace = namespace

			testData = model.DataProvider(
				"dbusers-operator-global",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				30008,
				[]func(*model.TestDataProvider){},
			).WithProject(project).
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

			mgr, err := k8s.BuildManager(&k8s.Config{
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
				err := mgr.Start(ctx)
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
