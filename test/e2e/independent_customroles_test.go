package e2e

import (
	"time"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/conditions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

var _ = Describe("Migrate one CustomRole from AtlasProject to AtlasCustomRole resource", Label("custom-roles", "independent-custom-roles"), func() {
	var testData *model.TestDataProvider
	var customRole = akov2.CustomRole{
		Name: "backupAdmin",
		InheritedRoles: []akov2.Role{
			{
				Name:     "backup",
				Database: "admin",
			},
		},
		Actions: []akov2.Action{
			{
				Name: "LIST_SESSIONS",
				Resources: []akov2.Resource{
					{
						Cluster: pointer.MakePtr(true),
					},
				},
			},
		},
	}

	_ = AfterEach(func() {
		GinkgoWriter.Println("")
		GinkgoWriter.Println("Independent Custom Roles Test")
		GinkgoWriter.Println("Operator namespace: " + testData.Resources.Namespace)
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Delete Resources", func() {
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	It("Should migrate one CustomRole from existing AtlasProject to dedicated AtlasCustomRole resource", func() {
		By("Creating AtlasProject", func() {
			testData = model.DataProvider(
				"project-custom-role",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject())

			actions.ProjectCreationFlow(testData)
		})

		By("Configuring one CustomRole", func() {
			Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).Should(Succeed())
			testData.Project.Spec.CustomRoles = []akov2.CustomRole{customRole}
			Expect(testData.K8SClient.Update(testData.Context, testData.Project)).To(Succeed())

			Eventually(func(g Gomega) {
				expectedConditions := conditions.MatchConditions(
					api.TrueCondition(api.ProjectCustomRolesReadyType),
					api.TrueCondition(api.ReadyType),
				)
				g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithPolling(10 * time.Second).WithTimeout(2 * time.Minute).Should(Succeed())
		})

		By("Disabling reconciliation for AtlasProject", func() {
			Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).Should(Succeed())
			testData.Project.Annotations[customresource.ReconciliationPolicyAnnotation] = customresource.ReconciliationPolicySkip
			Expect(testData.K8SClient.Update(testData.Context, testData.Project)).To(Succeed())

			Eventually(func(g Gomega) {
				Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).Should(Succeed())
				Expect(testData.Project.Annotations).To(HaveKeyWithValue(customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip))
			}).WithPolling(10 * time.Second).WithTimeout(2 * time.Minute).Should(Succeed())
		})

		By("Migrating the Role to the AtlasCustomRole resource", func() {
			akoCustomRole := &akov2.AtlasCustomRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-role",
					Namespace: testData.Project.GetNamespace(),
				},
				Spec: akov2.AtlasCustomRoleSpec{
					ProjectRef: &common.ResourceRefNamespaced{
						Name:      testData.Project.GetName(),
						Namespace: testData.Project.GetNamespace(),
					},
					Role: customRole,
					LocalCredentialHolder: api.LocalCredentialHolder{
						ConnectionSecret: &api.LocalObjectReference{Name: config.DefaultOperatorGlobalKey},
					},
				},
				Status: status.AtlasCustomRoleStatus{},
			}
			Expect(testData.K8SClient.Create(testData.Context, akoCustomRole)).To(Succeed())
			Eventually(func(g Gomega) {
				expectedConditions := conditions.MatchConditions(
					api.TrueCondition(api.ProjectCustomRolesReadyType),
					api.TrueCondition(api.ReadyType),
				)

				g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(akoCustomRole), akoCustomRole)).To(Succeed())
				g.Expect(akoCustomRole.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithPolling(10 * time.Second).WithTimeout(2 * time.Minute).Should(Succeed())
		})

		By("Removing custom roles from the AtlasProject", func() {
			Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).Should(Succeed())
			testData.Project.Spec.CustomRoles = []akov2.CustomRole{}
			Expect(testData.K8SClient.Update(testData.Context, testData.Project)).To(Succeed())
		})

		By("Enabled reconciliation for AtlasProject", func() {
			Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).Should(Succeed())
			testData.Project.Annotations = map[string]string{}
			Expect(testData.K8SClient.Update(testData.Context, testData.Project)).To(Succeed())

			Eventually(func(g Gomega) {
				expectedConditions := conditions.MatchConditions(
					api.TrueCondition(api.ReadyType),
				)
				g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).Should(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithPolling(10 * time.Second).WithTimeout(2 * time.Minute).Should(Succeed())
		})
	})
})
