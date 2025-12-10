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

package int

import (
	"context"
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas-sdk/v20250312010/admin"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/httputil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/resources"
)

var _ = Describe("AtlasProject", Label("int", "AtlasProject", "protection-enabled"), func() {
	var testNamespace *corev1.Namespace
	var stopManager context.CancelFunc
	var connectionSecret corev1.Secret

	BeforeEach(func() {
		By("Starting the operator", func() {
			testNamespace, stopManager = prepareControllers(true)
			Expect(testNamespace).ToNot(BeNil())
			Expect(stopManager).ToNot(BeNil())
		})

		By("Creating project connection secret", func() {
			connectionSecret = buildConnectionSecret(fmt.Sprintf("%s-atlas-key", testNamespace.Name))
			Expect(k8sClient.Create(context.Background(), &connectionSecret)).To(Succeed())
		})
	})

	Describe("Operator is running with deletion protection enabled", func() {
		It("Creates a project and protect it to be deleted", func() {
			testProject := &akov2.AtlasProject{}
			projectName := fmt.Sprintf("new-project-%s", testNamespace.Name)

			By("Creating a project in the cluster", func() {
				testProject = akov2.NewProject(testNamespace.Name, projectName, projectName).
					WithConnectionSecret(connectionSecret.Name)
				Expect(k8sClient.Create(context.Background(), testProject, &client.CreateOptions{})).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testProject, api.TrueCondition(api.ReadyType))
				}).WithTimeout(15 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})

			// nolint:dupl
			By("Deleting project in cluster doesn't delete from Atlas", func() {
				projectID := testProject.ID()
				Expect(k8sClient.Delete(context.Background(), testProject, &client.DeleteOptions{})).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(k8sClient.Get(context.Background(), client.ObjectKeyFromObject(testProject), testProject, &client.GetOptions{})).ToNot(Succeed())

					atlasProject, _, err := atlasClient.ProjectsApi.GetGroupByName(context.Background(), projectName).Execute()
					g.Expect(err).To(BeNil())
					g.Expect(atlasProject).ToNot(BeNil())
				}).WithTimeout(5 * time.Minute).WithPolling(PollingInterval).Should(Succeed())

				_, err := atlasClient.ProjectsApi.DeleteGroup(context.Background(), projectID).Execute()
				Expect(err).To(BeNil())
			})
		})

		It("Adds an existing Atlas project and protect it to be deleted", func() {
			testProject := &akov2.AtlasProject{}
			projectName := fmt.Sprintf("existing-project-%s", testNamespace.Name)

			By("Creating a project in Atlas", func() {
				atlasProject := admin.Group{
					OrgId:                     orgID,
					Name:                      projectName,
					WithDefaultAlertsSettings: pointer.MakePtr(true),
				}
				_, _, err := atlasClient.ProjectsApi.CreateGroup(context.Background(), &atlasProject).Execute()
				Expect(err).To(BeNil())
			})

			By("Creating a project in the cluster", func() {
				testProject = akov2.NewProject(testNamespace.Name, projectName, projectName).
					WithConnectionSecret(connectionSecret.Name)
				Expect(k8sClient.Create(context.Background(), testProject, &client.CreateOptions{})).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testProject, api.TrueCondition(api.ReadyType))
				}).WithTimeout(15 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})

			// nolint:dupl
			By("Deleting project in cluster doesn't delete from Atlas", func() {
				projectID := testProject.ID()
				Expect(k8sClient.Delete(context.Background(), testProject, &client.DeleteOptions{})).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(k8sClient.Get(context.Background(), client.ObjectKeyFromObject(testProject), testProject, &client.GetOptions{})).ToNot(Succeed())

					atlasProject, _, err := atlasClient.ProjectsApi.GetGroupByName(context.Background(), projectName).Execute()
					g.Expect(err).To(BeNil())
					g.Expect(atlasProject).ToNot(BeNil())
				}).WithTimeout(5 * time.Minute).WithPolling(PollingInterval).Should(Succeed())

				_, err := atlasClient.ProjectsApi.DeleteGroup(context.Background(), projectID).Execute()
				Expect(err).To(BeNil())
			})
		})

		It("Creates a project and annotate it to be deleted", func() {
			testProject := &akov2.AtlasProject{}
			projectName := fmt.Sprintf("new-project-%s", testNamespace.Name)

			By("Creating a project in the cluster", func() {
				testProject = akov2.NewProject(testNamespace.Name, projectName, projectName).
					WithAnnotations(map[string]string{customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyDelete}).
					WithConnectionSecret(connectionSecret.Name)
				Expect(k8sClient.Create(context.Background(), testProject, &client.CreateOptions{})).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testProject, api.TrueCondition(api.ReadyType))
				}).WithTimeout(15 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})

			// nolint:dupl
			By("Deleting project in cluster should delete it from Atlas", func() {
				projectID := testProject.ID()
				Expect(k8sClient.Delete(context.Background(), testProject, &client.DeleteOptions{})).To(Succeed())

				Eventually(func(g Gomega) {
					_, r, err := atlasClient.ProjectsApi.GetGroup(context.Background(), projectID).Execute()
					g.Expect(err).ToNot(BeNil())
					g.Expect(httputil.StatusCode(r)).To(Equal(http.StatusNotFound))
				}).WithTimeout(5 * time.Minute).WithPolling(PollingInterval).Should(Succeed())
			})
		})
	})

	AfterEach(func() {
		By("Deleting project connection secret", func() {
			Expect(k8sClient.Delete(context.Background(), &connectionSecret)).To(Succeed())
		})

		By("Stopping the operator", func() {
			stopManager()
			err := k8sClient.Delete(context.Background(), testNamespace)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
