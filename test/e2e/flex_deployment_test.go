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
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas-sdk/v20250312014/admin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/httputil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
	akoretry "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/retry"
)

var _ = Describe("Flex", Label("flex"), func() {
	var testData *model.TestDataProvider
	var flexDeployment *akov2.AtlasDeployment
	var apiClient *admin.APIClient

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Flex test\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
			Expect(actions.SaveDeploymentsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Delete Resources", func() {
			actions.DeleteTestDataDeployments(testData)
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	BeforeEach(func(ctx SpecContext) {
		By("Setting up cloud environment", func() {
			testData = model.DataProvider(ctx, "atlas-flex", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 30005, []func(*model.TestDataProvider){}).
				WithProject(data.DefaultProject())
			actions.ProjectCreationFlow(testData)
		})
		By("Creating the API Client", func() {
			var err error
			domain := os.Getenv("MCLI_OPS_MANAGER_URL")
			pubKey := os.Getenv("MCLI_PUBLIC_API_KEY")
			prvKey := os.Getenv("MCLI_PRIVATE_API_KEY")
			apiClient, err = admin.NewClient(
				admin.UseBaseURL(domain),
				admin.UseDigestAuth(pubKey, prvKey),
			)
			Expect(err).To(BeNil())

		})
	})

	It("Creates a Flex Cluster", func(ctx context.Context) {
		name := fmt.Sprintf("%s-flex", testData.Project.Name)

		By("Creating an AtlasDeployment CR with a Flex cluster defined", func() {
			flexDeployment = &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: testData.Resources.Namespace,
				},
				Spec: akov2.AtlasDeploymentSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name: testData.Project.Name,
						},
					},
					FlexSpec: &akov2.FlexSpec{
						Name: name,
						ProviderSettings: &akov2.FlexProviderSettings{
							BackingProviderName: "AWS",
							RegionName:          "US_EAST_1",
						},
					},
				},
			}
			Expect(testData.K8SClient.Create(testData.Context, flexDeployment)).Should(Succeed())
		})
		By("Checking Deployment status conditions in Kube", func() {
			Eventually(func(g Gomega) bool {
				err := testData.K8SClient.Get(testData.Context, types.NamespacedName{Name: name, Namespace: testData.Resources.Namespace}, flexDeployment)
				Expect(err).To(BeNil())

				conditionTypes := []api.ConditionType{api.DeploymentReadyType, api.ReadyType}

				for _, conditionType := range conditionTypes {
					found := false
					for _, con := range flexDeployment.Status.Conditions {
						if con.Type == conditionType && con.Status == corev1.ConditionTrue {
							found = true
							break
						}
					}
					if !found {
						return false
					}
				}
				return true
			}).WithTimeout(15 * time.Minute).WithPolling(20 * time.Second).Should(BeTrue())
		})
		By("Checking the Flex cluster is ready in Atlas", func() {

			Eventually(func(g Gomega) {
				flex, _, err := apiClient.FlexClustersApi.GetFlexCluster(testData.Context, testData.Project.ID(), name).Execute()
				g.Expect(err).To(BeNil())
				g.Expect(flex.GetStateName()).To(Equal("IDLE"))
			}).WithTimeout(1 * time.Minute).WithPolling(PollingInterval).Should(Succeed())

		})
		By("Adding tags to the Flex cluster", func() {
			tags := []*akov2.TagSpec{
				{Key: "test-key", Value: "test-value"},
				{Key: "another-test-key", Value: "test-value-2"},
			}
			_, err := akoretry.RetryUpdateOnConflict(
				ctx,
				testData.K8SClient,
				client.ObjectKeyFromObject(flexDeployment),
				func(flex *akov2.AtlasDeployment) {
					flex.Spec.FlexSpec.Tags = tags
				})
			Expect(err).To(BeNil())
		})
		By("Deleting the AtlasDeployment CR", func() {
			Expect(testData.K8SClient.Delete(ctx, flexDeployment)).Should(Succeed())
		})
		By("Checking the Flex cluster is deleted in Atlas", func() {
			Eventually(func(g Gomega) {
				_, resp, err := apiClient.FlexClustersApi.GetFlexCluster(testData.Context, testData.Project.ID(), name).Execute()
				g.Expect(err).ToNot(BeNil())
				g.Expect(httputil.StatusCode(resp)).Should(Equal(404))
			}).WithTimeout(5 * time.Minute).WithPolling(PollingInterval).Should(Succeed())
		})
	})
})
