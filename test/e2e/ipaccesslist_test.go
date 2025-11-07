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
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/conditions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

var _ = Describe("Migrate ip access list from sub-resources to separate custom resources", Label("ip-access-list"), func() {
	var testData *model.TestDataProvider
	var ial *akov2.AtlasIPAccessList

	_ = AfterEach(func() {
		GinkgoWriter.Println()
		GinkgoWriter.Println("===============================================")
		GinkgoWriter.Println("Operator namespace: " + testData.Resources.Namespace)
		GinkgoWriter.Println("===============================================")
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Delete Project and cluster resources", func() {
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	It("Should migrate a ip access list configured in a project as sub-resource to a separate custom resource", func(ctx SpecContext) {
		By("Setting up project", func() {
			p := data.DefaultProject()
			p.Spec.ProjectIPAccessList = nil
			testData = model.DataProvider(ctx, "migrate-ip-access-list", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(p)

			actions.ProjectCreationFlow(testData)
		})

		By("Configuring ip-access-list as a sub-resource", func() {
			testData.Project.Spec.ProjectIPAccessList = []project.IPAccessList{
				{
					CIDRBlock: "10.1.1.0/24",
					Comment:   "Company Network",
				},
				{
					IPAddress: "192.168.0.100",
					Comment:   "Private IP",
				},
			}

			Expect(testData.K8SClient.Update(testData.Context, testData.Project)).To(Succeed())
			Eventually(func(g Gomega) {
				expectedConditions := conditions.MatchConditions(
					api.TrueCondition(api.IPAccessListReadyType),
					api.TrueCondition(api.ReadyType),
				)

				g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(5 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
		})
		//nolint:dupl
		By("Stopping reconciling project and its sub-resources", func() {
			Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Annotations[customresource.ReconciliationPolicyAnnotation] = customresource.ReconciliationPolicySkip
			testData.Project.Spec.ProjectIPAccessList = nil

			Expect(testData.K8SClient.Update(testData.Context, testData.Project)).To(Succeed())
			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElement(conditions.MatchCondition(api.TrueCondition(api.ReadyType))))
			}).WithTimeout(5 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
		})

		By("Migrate ip access list as separate custom resource", func() {
			ial = &akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ial-" + testData.Resources.TestID,
					Namespace: testData.Resources.Namespace,
				},
				Spec: akov2.AtlasIPAccessListSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      testData.Project.Name,
							Namespace: testData.Project.Namespace,
						},
					},
					Entries: []akov2.IPAccessEntry{
						{
							CIDRBlock: "10.1.1.0/24",
							Comment:   "Company Network",
						},
						{
							IPAddress: "192.168.0.100",
							Comment:   "Private IP",
						},
					},
				},
			}

			Expect(testData.K8SClient.Create(testData.Context, ial)).To(Succeed())
			Eventually(func(g Gomega) {
				expectedConditions := conditions.MatchConditions(
					api.TrueCondition(api.IPAccessListReady),
					api.TrueCondition(api.ReadyType),
				)
				g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(ial), ial)).To(Succeed())
				g.Expect(ial.Status.Conditions).To(ContainElements(expectedConditions))
				for _, eStatus := range ial.Status.Entries {
					g.Expect(eStatus.Status).To(Equal("ACTIVE"))
				}
			}).WithTimeout(5 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
		})
		//nolint:dupl
		By("Restating project reconciliation", func() {
			Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			delete(testData.Project.Annotations, customresource.ReconciliationPolicyAnnotation)

			Expect(testData.K8SClient.Update(testData.Context, testData.Project)).To(Succeed())
			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElement(conditions.MatchCondition(api.TrueCondition(api.ReadyType))))
			}).WithTimeout(5 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
		})

		By("Updating project doesn't affect ip access list", func() {
			Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Spec.Settings = &akov2.ProjectSettings{
				IsSchemaAdvisorEnabled: pointer.MakePtr(true),
			}

			Expect(testData.K8SClient.Update(testData.Context, testData.Project)).To(Succeed())
			Eventually(func(g Gomega) {
				notExpectedConditions := conditions.MatchConditions(
					api.TrueCondition(api.IPAccessListReadyType),
				)

				g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).ToNot(ContainElements(notExpectedConditions))
				g.Expect(testData.Project.Status.Conditions).To(ContainElement(conditions.MatchCondition(api.TrueCondition(api.ReadyType))))
			}).WithTimeout(5 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
		})

		By("IP Access List are still ready", func() {
			Eventually(func(g Gomega) { //nolint:dupl
				expectedConditions := conditions.MatchConditions(
					api.TrueCondition(api.IPAccessListReady),
					api.TrueCondition(api.ReadyType),
				)

				g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(ial), ial)).To(Succeed())
				g.Expect(ial.Status.Conditions).To(ContainElements(expectedConditions))
				for _, eStatus := range ial.Status.Entries {
					g.Expect(eStatus.Status).To(Equal("ACTIVE"))
				}
			}).WithTimeout(5 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
		})

		By("Removing ip access list", func() {
			Expect(testData.K8SClient.Delete(testData.Context, ial)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(ial), ial)).ShouldNot(Succeed())
			}).WithTimeout(5 * time.Minute).WithPolling(20 * time.Second).Should(Succeed())
		})
	})
})
