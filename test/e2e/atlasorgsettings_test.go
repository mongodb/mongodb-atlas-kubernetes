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
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
)

var _ = Describe("AtlasOrgSettings", Label("atlas-org-settings"), func() {
	var testData *model.TestDataProvider

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("AtlasOrgSettings Test\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
			Expect(actions.SaveAtlasOrgSettingsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Delete Resources", func() {
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	It("Should create AtlasOrgSettings with strict configuration and verify its status", func(ctx SpecContext) {
		testData = model.DataProvider(ctx, "atlas-org-settings", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())

		actions.ProjectCreationFlow(testData)

		orgSettings := akov2.AtlasOrgSettings{
			ObjectMeta: metav1.ObjectMeta{
				Name: utils.RandomName("org-settings-strict"),
			},
			Spec: akov2.AtlasOrgSettingsSpec{
				OrgID:                                  "",
				ApiAccessListRequired:                  pointer.MakePtr(true),
				GenAIFeaturesEnabled:                   pointer.MakePtr(false),
				MultiFactorAuthRequired:                pointer.MakePtr(true),
				RestrictEmployeeAccess:                 pointer.MakePtr(true),
				SecurityContact:                        pointer.MakePtr("security@example.com"),
				MaxServiceAccountSecretValidityInHours: pointer.MakePtr(24),
			},
		}

		By("Set the organization ID from environment", func() {
			orgSettings.Spec.OrgID = os.Getenv("MCLI_ORG_ID")
			Expect(orgSettings.Spec.OrgID).ShouldNot(BeEmpty(), "MCLI_ORG_ID environment variable must be set")
		})

		By("Create AtlasOrgSettings", func() {
			orgSettings.Namespace = testData.Resources.Namespace
			Expect(testData.K8SClient.Create(testData.Context, &orgSettings)).Should(Succeed())
		})

		By("Wait for AtlasOrgSettings to be ready", func() {
			Eventually(func(g Gomega) bool {
				return atlasOrgSettingsIsReady(g, testData, orgSettings.Name)
			}).WithTimeout(10 * time.Minute).WithPolling(20 * time.Second).Should(BeTrue())
		})

		By("Check AtlasOrgSettings status conditions", func() {
			Eventually(func(g Gomega) bool {
				currentOrgSettings := &akov2.AtlasOrgSettings{}
				err := testData.K8SClient.Get(
					testData.Context,
					types.NamespacedName{Name: orgSettings.Name, Namespace: testData.Resources.Namespace},
					currentOrgSettings,
				)
				g.Expect(err).ShouldNot(HaveOccurred())

				for _, condition := range currentOrgSettings.Status.Conditions {
					if condition.Type == string(api.ReadyType) && condition.Status == metav1.ConditionTrue {
						return true
					}
				}
				return false
			}).WithTimeout(10 * time.Minute).WithPolling(20 * time.Second).Should(BeTrue())
		})

		By("Verify AtlasOrgSettings spec fields are preserved", func() {
			currentOrgSettings := &akov2.AtlasOrgSettings{}
			Expect(testData.K8SClient.Get(
				testData.Context,
				types.NamespacedName{Name: orgSettings.Name, Namespace: testData.Resources.Namespace},
				currentOrgSettings,
			)).Should(Succeed())

			Expect(currentOrgSettings.Spec.OrgID).Should(Equal(orgSettings.Spec.OrgID))
			Expect(currentOrgSettings.Spec.ConnectionSecretRef).Should(Equal(orgSettings.Spec.ConnectionSecretRef))

			if orgSettings.Spec.ApiAccessListRequired != nil {
				Expect(currentOrgSettings.Spec.ApiAccessListRequired).Should(Equal(orgSettings.Spec.ApiAccessListRequired))
			}
			if orgSettings.Spec.GenAIFeaturesEnabled != nil {
				Expect(currentOrgSettings.Spec.GenAIFeaturesEnabled).Should(Equal(orgSettings.Spec.GenAIFeaturesEnabled))
			}
			if orgSettings.Spec.SecurityContact != nil {
				Expect(currentOrgSettings.Spec.SecurityContact).Should(Equal(orgSettings.Spec.SecurityContact))
			}
		})

		By("Delete AtlasOrgSettings", func() {
			currentOrgSettings := &akov2.AtlasOrgSettings{}
			Expect(testData.K8SClient.Get(
				testData.Context,
				types.NamespacedName{Name: orgSettings.Name, Namespace: testData.Resources.Namespace},
				currentOrgSettings,
			)).Should(Succeed())

			Expect(testData.K8SClient.Delete(testData.Context, currentOrgSettings)).Should(Succeed())
		})

		By("Wait for AtlasOrgSettings to be deleted", func() {
			Eventually(func(g Gomega) bool {
				currentOrgSettings := &akov2.AtlasOrgSettings{}
				err := testData.K8SClient.Get(
					testData.Context,
					types.NamespacedName{Name: orgSettings.Name, Namespace: testData.Resources.Namespace},
					currentOrgSettings,
				)
				return err != nil // Should return not found error
			}).WithTimeout(5 * time.Minute).WithPolling(10 * time.Second).Should(BeTrue())
		})
	})
})

func atlasOrgSettingsIsReady(g Gomega, userData *model.TestDataProvider, orgSettingsName string) bool {
	currentOrgSettings := &akov2.AtlasOrgSettings{}
	err := userData.K8SClient.Get(
		userData.Context,
		types.NamespacedName{Name: orgSettingsName, Namespace: userData.Resources.Namespace},
		currentOrgSettings,
	)
	g.Expect(err).ShouldNot(HaveOccurred())

	if currentOrgSettings.ResourceVersion == "" {
		return false
	}

	if len(currentOrgSettings.Status.Conditions) == 0 {
		return false
	}

	for _, condition := range currentOrgSettings.Status.Conditions {
		if condition.Type == string(api.ReadyType) {
			return condition.Status == metav1.ConditionTrue
		}
	}

	return false
}
