package e2e

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

var _ = Describe("Operator to run db-user with the OIDC feature flags disabled", Ordered, Label("users-oidc"), func() {
	var testData *model.TestDataProvider

	_ = BeforeEach(func() {
		project := data.DefaultProject()

		deployment := data.CreateBasicDeployment("dbusers-operator-global")

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
					data.WithNamespace(project.Namespace),
					data.WithLabels([]common.LabelSpec{
						{Key: "type", Value: "e2e-test"},
						{Key: "context", Value: "cloud"},
					}),
				),
			)

		actions.CreateNamespaceAndSecrets(testData)
	})

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

	It("Operator run on global namespace with the OIDC feature enabled", func() {
		By("Running operator watching global namespace with OIDC enabled", func() {
			actions.ProjectCreationFlow(testData)
		})
		By("Creating database users resource", func() {
			deploy.CreateUsers(testData)
		})
		By("Creating a FedAuth resource and verify it is ready", func() {
			fedAuth := &akov2.AtlasFederatedAuth{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-fedauth", testData.Project.Name),
					Namespace: testData.Resources.Namespace,
				},
				Spec: akov2.AtlasFederatedAuthSpec{
					Enabled:                     true,
					DataAccessIdentityProviders: &[]string{"abc"},
					ConnectionSecretRef: common.ResourceRefNamespaced{
						Namespace: testData.Resources.Namespace,
						Name:      config.DefaultOperatorGlobalKey,
					},
				},
			}
			Expect(testData.K8SClient.Create(testData.Context, fedAuth)).Should(Succeed())
			Eventually(func(g Gomega) {
				currentFedAuth := &akov2.AtlasFederatedAuth{}
				g.Expect(testData.K8SClient.Get(context.Background(),
					types.NamespacedName{
						Name:      fedAuth.Name,
						Namespace: fedAuth.Namespace,
					}, currentFedAuth)).NotTo(HaveOccurred())
				for _, condition := range currentFedAuth.Status.Conditions {
					if condition.Type == api.ReadyType {
						g.Expect(condition.Status).Should(Equal(corev1.ConditionTrue))
					}
				}
			}).WithTimeout(1 * time.Minute).WithPolling(20 * time.Second).Should(Succeed())
		})
		By("Try to enabled the OIDC Group feature for the user", func() {
			currentUser := &akov2.AtlasDatabaseUser{}
			Expect(testData.K8SClient.Get(context.Background(),
				types.NamespacedName{
					Name:      testData.Users[0].Name,
					Namespace: testData.Users[0].Namespace,
				}, currentUser)).NotTo(HaveOccurred())

			currentUser.Spec.OIDCAuthType = "IDP_GROUP"
			currentUser.Spec.PasswordSecret = nil
			Expect(testData.K8SClient.Update(context.Background(), currentUser)).NotTo(HaveOccurred())
		})

		By("Verify if user is ready", func() {
			currentUser := &akov2.AtlasDatabaseUser{}
			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(context.Background(),
					types.NamespacedName{
						Name:      testData.Users[0].Name,
						Namespace: testData.Users[0].Namespace,
					}, currentUser)).NotTo(HaveOccurred())
				for _, condition := range currentUser.Status.Conditions {
					if condition.Type == api.ReadyType {
						g.Expect(condition.Status).Should(Equal(corev1.ConditionTrue))
					}
				}
			}).WithTimeout(1 * time.Minute).WithPolling(20 * time.Second).Should(Succeed())
		})
		By("Try to enabled the OIDC User feature for the user", func() {
			currentUser := &akov2.AtlasDatabaseUser{}
			Expect(testData.K8SClient.Get(context.Background(),
				types.NamespacedName{
					Name:      testData.Users[0].Name,
					Namespace: testData.Users[0].Namespace,
				}, currentUser)).NotTo(HaveOccurred())

			currentUser.Spec.OIDCAuthType = "USER"
			Expect(testData.K8SClient.Update(context.Background(), currentUser)).NotTo(HaveOccurred())
		})

		By("Verify if user is ready", func() {
			currentUser := &akov2.AtlasDatabaseUser{}
			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(context.Background(),
					types.NamespacedName{
						Name:      testData.Users[0].Name,
						Namespace: testData.Users[0].Namespace,
					}, currentUser)).NotTo(HaveOccurred())
				for _, condition := range currentUser.Status.Conditions {
					if condition.Type == api.ReadyType {
						g.Expect(condition.Status).Should(Equal(corev1.ConditionTrue))
					}
				}
			}).WithTimeout(1 * time.Minute).WithPolling(20 * time.Second).Should(Succeed())
		})
	})
})
