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
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/conditions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/events"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/resources"
)

var _ = Describe("Sync Period test", Label("int", "sync-period"), func() {
	const interval = time.Second * 2
	const syncInterval = 40 * time.Second

	var (
		connectionSecret        corev1.Secret
		createdProject          *akov2.AtlasProject
		previousResourceVersion string
	)

	BeforeEach(func() {
		prepareControllersWithSyncPeriod(false, syncInterval)

		createdProject = &akov2.AtlasProject{}

		connectionSecret = buildConnectionSecret("my-atlas-key")
		By(fmt.Sprintf("Creating the Secret %s", kube.ObjectKeyFromObject(&connectionSecret)))
		Expect(k8sClient.Create(context.Background(), &connectionSecret)).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		if createdProject != nil && createdProject.Status.ID != "" {
			By("Removing Atlas Project " + createdProject.Status.ID)
			Eventually(deleteK8sObject(createdProject), 20, interval).Should(BeTrue())
			Eventually(checkAtlasProjectRemoved(createdProject.Status.ID), 20, interval).Should(BeTrue())
		}
		removeControllersAndNamespace()
	})
	It("Should reconcile after defined SyncPeriod", func() {
		By("Should Succeed with creating the project", func() {
			expectedProject := akov2.DefaultProject(namespace.Name, connectionSecret.Name)
			createdProject.ObjectMeta = expectedProject.ObjectMeta
			Expect(k8sClient.Create(context.Background(), expectedProject)).ToNot(HaveOccurred())

			Eventually(func() bool {
				return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			projectReadyConditions := conditions.MatchConditions(
				api.TrueCondition(api.ProjectReadyType),
				api.TrueCondition(api.ReadyType),
				api.TrueCondition(api.ValidationSucceeded),
			)
			Expect(createdProject.Status.ID).NotTo(BeNil())
			Expect(createdProject.Status.Conditions).To(ContainElements((projectReadyConditions)))
			Expect(createdProject.Status.ObservedGeneration).To(Equal(createdProject.Generation))

			atlasProject, _, err := atlasClient.ProjectsApi.
				GetGroup(context.Background(), createdProject.Status.ID).
				Execute()
			Expect(err).ToNot(HaveOccurred())

			Expect(atlasProject.Name).To(Equal(expectedProject.Spec.Name))

			events.EventExists(k8sClient, createdProject, "Normal", "Ready", "")

			Eventually(func(g Gomega) bool {
				if !resources.ReadAtlasResource(context.Background(), k8sClient, createdProject) {
					return false
				}
				previousResourceVersion = createdProject.ResourceVersion
				return true
			}).WithTimeout(10 * time.Second).WithPolling(2 * time.Second).Should(BeTrue())
		})

		By(fmt.Sprintf("Should wait for at least %f seconds", (syncInterval*2).Seconds()), func() {
			time.Sleep(syncInterval * 2)
		})

		By("Project resource version should be different", func() {
			var currentResourceVersion string
			Eventually(func(g Gomega) bool {
				if !resources.ReadAtlasResource(context.Background(), k8sClient, createdProject) {
					return false
				}
				currentResourceVersion = createdProject.ResourceVersion
				return true
			}).WithTimeout(10 * time.Second).WithPolling(2 * time.Second).Should(BeTrue())

			Expect(currentResourceVersion).ToNot(BeEquivalentTo(previousResourceVersion))
		})
	})
})
