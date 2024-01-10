package e2e

import (
	"context"
	"time"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	dbuserController "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasdatabaseuser"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/resources"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Operator to run db-user with the OIDC feature flags", Ordered, Label("users-oidc"), func() {
	var testData *model.TestDataProvider

	_ = AfterEach(func() {
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
			Expect(actions.SaveDeploymentsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
			Expect(actions.SaveUsersToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		actions.DeleteTestDataUsers(testData)
		actions.DeleteTestDataProject(testData)
		actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
	})

	It("Operator run on global namespace with the OIDC feature disabled", func() {
		By("Setting up test data with the DB User OIDC disabled", func() {
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
					),
				)
			testData.Resources.Namespace = config.DefaultOperatorNS
		})

		By("Running operator watching global namespace with OIDC disabled", func() {
			Eventually(k8s.CreateNamespace(testData.Context, testData.K8SClient, config.DefaultOperatorNS)).WithTimeout(5 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
			k8s.CreateDefaultSecret(testData.Context, testData.K8SClient, config.DefaultOperatorGlobalKey, config.DefaultOperatorNS)

			mgr, err := k8s.BuildManager(&k8s.Config{
				GlobalAPISecret: client.ObjectKey{
					Namespace: config.DefaultOperatorNS,
					Name:      config.DefaultOperatorGlobalKey,
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
					g.Expect(resources.CheckCondition(testData.K8SClient, dbUser, status.TrueCondition(status.ReadyType))).To(BeTrue())
				}

				return true
			}).WithTimeout(5 * time.Minute).WithPolling(10 * time.Second).Should(BeTrue())
		})
		By("Try to enabled the OIDC feature for the user", func() {
			currentUser := &mdbv1.AtlasDatabaseUser{}
			Expect(testData.K8SClient.Get(context.Background(),
				types.NamespacedName{
					Name:      testData.Users[0].Name,
					Namespace: testData.Users[0].Namespace,
				}, currentUser)).NotTo(HaveOccurred())

			currentUser.Spec.OIDCAuthType = "IDP_GROUP"
			Expect(testData.K8SClient.Update(context.Background(), currentUser)).NotTo(HaveOccurred())
		})

		By("Verify if user is ready. It shouldn't be", func() {
			currentUser := &mdbv1.AtlasDatabaseUser{}
			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(context.Background(),
					types.NamespacedName{
						Name:      testData.Users[0].Name,
						Namespace: testData.Users[0].Namespace,
					}, currentUser)).NotTo(HaveOccurred())
				for _, condition := range currentUser.Status.Conditions {
					if condition.Type == status.ReadyType {
						g.Expect(condition.Status).Should(Equal(corev1.ConditionFalse))
						g.Expect(condition.Message).To(ContainSubstring(dbuserController.ErrOIDCNotEnabled.Error()))
					}
				}
			}).WithTimeout(1 * time.Minute).WithPolling(20 * time.Second).Should(Succeed())
		})
	})

	// TODO: Enable this test as soon as API for configuring OpenID providers becomes available
	// It("Operator run on global namespace with the OIDC feature enabled", func() {
	// 	By("Setting up test data with the DB User OIDC enabled", func() {
	// 		project := data.DefaultProject()
	// 		project.Namespace = config.DefaultOperatorNS

	// 		deployment := data.CreateBasicDeployment("dbusers-operator-global")
	// 		deployment.Namespace = config.DefaultOperatorNS

	// 		testData = model.DataProvider(
	// 			"dbusers-operator-global",
	// 			model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
	// 			30008,
	// 			[]func(*model.TestDataProvider){},
	// 		).WithProject(project).
	// 			WithInitialDeployments(deployment).
	// 			WithUsers(
	// 				data.BasicUser(
	// 					"reader1",
	// 					"reader1",
	// 					data.WithReadWriteRole(),
	// 					data.WithNamespace(config.DefaultOperatorNS),
	// 					data.WithOIDCEnabled(),
	// 				),
	// 			)
	// 		testData.Resources.Namespace = config.DefaultOperatorNS
	// 	})

	// 	By("Running operator watching global namespace with OIDC enabled", func() {
	// 		Eventually(k8s.CreateNamespace(testData.Context, testData.K8SClient, config.DefaultOperatorNS)).WithTimeout(5 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
	// 		k8s.CreateDefaultSecret(testData.Context, testData.K8SClient, config.DefaultOperatorGlobalKey, config.DefaultOperatorNS)
	// 		k8s.CreateNamespace(testData.Context, testData.K8SClient, secondNamespace)

	// 		mgr, err := k8s.BuildManager(&k8s.Config{
	// 			GlobalAPISecret: client.ObjectKey{
	// 				Namespace: config.DefaultOperatorNS,
	// 				Name:      config.DefaultOperatorGlobalKey,
	// 			},
	// 			FeatureFlags: featureflags.NewFeatureFlags(func() []string { return []string{featureflags.FeatureOIDC} }),
	// 		})
	// 		Expect(err).NotTo(HaveOccurred())

	// 		go func(ctx context.Context) context.Context {
	// 			err := mgr.Start(ctx)
	// 			Expect(err).NotTo(HaveOccurred())
	// 			return ctx
	// 		}(testData.Context)
	// 	})
	// 	By("Creating project and database users resources", func() {
	// 		deploy.CreateProject(testData)
	// 		deploy.CreateUsers(testData)

	// 		Eventually(func(g Gomega) bool {
	// 			for i := range testData.Users {
	// 				dbUser := testData.Users[i]

	// 				g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(dbUser), dbUser)).To(Succeed())
	// 				g.Expect(resources.CheckCondition(testData.K8SClient, dbUser, status.TrueCondition(status.ReadyType))).To(BeTrue())
	// 			}

	// 			return true
	// 		}).WithTimeout(5 * time.Minute).WithPolling(10 * time.Second).Should(BeTrue())
	// 	})
	// 	By("Try to enabled the OIDC feature for the user", func() {
	// 		currentUser := &mdbv1.AtlasDatabaseUser{}
	// 		Expect(testData.K8SClient.Get(context.Background(),
	// 			types.NamespacedName{
	// 				Name:      testData.Users[0].Name,
	// 				Namespace: testData.Users[0].Namespace,
	// 			}, currentUser)).NotTo(HaveOccurred())

	// 		currentUser.Spec.OIDCAuthType = "IDP_GROUP"
	// 		Expect(testData.K8SClient.Update(context.Background(), currentUser)).NotTo(HaveOccurred())
	// 	})

	// 	By("Verify if user is ready", func() {
	// 		currentUser := &mdbv1.AtlasDatabaseUser{}
	// 		Eventually(func(g Gomega) {
	// 			g.Expect(testData.K8SClient.Get(context.Background(),
	// 				types.NamespacedName{
	// 					Name:      testData.Users[0].Name,
	// 					Namespace: testData.Users[0].Namespace,
	// 				}, currentUser)).NotTo(HaveOccurred())
	// 			for _, condition := range currentUser.Status.Conditions {
	// 				if condition.Type == status.ReadyType {
	// 					g.Expect(condition.Status).Should(Equal(corev1.ConditionTrue))
	// 				}
	// 			}
	// 		}).WithTimeout(1 * time.Minute).WithPolling(20 * time.Second).Should(Succeed())
	// 	})
	// })
})
