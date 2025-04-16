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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/resources"
)

// nolint:dupl
var _ = Describe("AtlasBackupSchedule Deletion Protected",
	Ordered,
	Label("AtlasDeployment", "AtlasBackupSchedule", "deletion-protection", "deletion-protection-backup"), func() {
		var testNamespace *corev1.Namespace
		var stopManager context.CancelFunc
		var connectionSecret corev1.Secret
		var testProject *akov2.AtlasProject
		var testDeployment *akov2.AtlasDeployment

		BeforeAll(func() {
			By("Starting the operator with protection ON", func() {
				testNamespace, stopManager = prepareControllers(true)
				Expect(testNamespace).ToNot(BeNil())
				Expect(stopManager).ToNot(BeNil())
			})

			By("Creating project connection secret", func() {
				connectionSecret = buildConnectionSecret(fmt.Sprintf("%s-atlas-key", testNamespace.Name))
				Expect(k8sClient.Create(context.Background(), &connectionSecret)).To(Succeed())
			})

			By("Creating a project with deletion annotation", func() {
				testProject = akov2.DefaultProject(testNamespace.Name, connectionSecret.Name).WithIPAccessList(project.NewIPAccessList().WithCIDR("0.0.0.0/0"))
				customresource.SetAnnotation( // this test project must be deleted
					testProject,
					customresource.ResourcePolicyAnnotation,
					customresource.ResourcePolicyDelete,
				)
				Expect(k8sClient.Create(context.Background(), testProject, &client.CreateOptions{})).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testProject, api.TrueCondition(api.ReadyType))
				}).WithTimeout(3 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})
		})

		AfterAll(func() {
			By("Deleting deployment from Atlas", func() {
				if testDeployment == nil {
					return
				}
				Expect(deleteAtlasDeployment(testProject.Status.ID, testDeployment.Spec.DeploymentSpec.Name)).To(Succeed())
			})
			By("Deleting deployment from Kubernetes", func() {
				Expect(k8sClient.Delete(context.Background(), testDeployment)).To(Succeed())
			})
			By("Deleting project from k8s and atlas", func() {
				Expect(k8sClient.Delete(context.Background(), testProject, &client.DeleteOptions{})).To(Succeed())
				Eventually(
					checkAtlasProjectRemoved(testProject.Status.ID),
				).WithTimeout(20 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})

			By("Deleting project connection secret", func() {
				Expect(k8sClient.Delete(context.Background(), &connectionSecret)).To(Succeed())
			})

			By("Stopping the operator", func() {
				stopManager()
				err := k8sClient.Delete(context.Background(), testNamespace)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		It("Should process BackupSchedule with deletion protection ON", func() {
			var bsPolicy *akov2.AtlasBackupPolicy
			var bsSchedule *akov2.AtlasBackupSchedule
			By("Creating AtlasBackupPolicy resource", func() {
				bsPolicy = &akov2.AtlasBackupPolicy{
					TypeMeta: metav1.TypeMeta{
						Kind:       "atlas.mongodb.com/v1",
						APIVersion: "AtlasBackupPolicy",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:        "bs-policy",
						Namespace:   testNamespace.Name,
						Labels:      map[string]string{},
						Annotations: map[string]string{},
					},
					Spec: akov2.AtlasBackupPolicySpec{
						Items: []akov2.AtlasBackupPolicyItem{
							{
								FrequencyType:     "daily",
								FrequencyInterval: 1,
								RetentionUnit:     "days",
								RetentionValue:    20,
							},
						},
					},
				}
				Expect(k8sClient.Create(context.Background(), bsPolicy)).To(Succeed())
			})

			By("Creating AtlasBackupSchedule resource", func() {
				bsSchedule = &akov2.AtlasBackupSchedule{
					TypeMeta: metav1.TypeMeta{
						Kind:       "atlas.mongodb.com/v1",
						APIVersion: "AtlasBackupSchedule",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bs-schedule",
						Namespace: testNamespace.Name,
					},
					Spec: akov2.AtlasBackupScheduleSpec{
						AutoExportEnabled: false,
						PolicyRef: common.ResourceRefNamespaced{
							Name:      bsPolicy.Name,
							Namespace: bsPolicy.Namespace,
						},
						ReferenceHourOfDay:                12,
						ReferenceMinuteOfHour:             20,
						RestoreWindowDays:                 2,
						UpdateSnapshots:                   false,
						UseOrgAndGroupNamesInExportPrefix: false,
						CopySettings:                      []akov2.CopySetting{},
					},
				}
				Expect(k8sClient.Create(context.Background(), bsSchedule)).To(Succeed())
			})

			By("Creating a deployment with backups enabled (default)", func() {
				testDeployment = akov2.DefaultAWSDeployment(testNamespace.Name, testProject.Name)
				testDeployment.Spec.DeploymentSpec.BackupEnabled = pointer.MakePtr(true)
				Expect(k8sClient.Create(context.Background(), testDeployment)).To(Succeed())
			})

			By("Deployment should be Ready", func() {
				Eventually(func(g Gomega) bool {
					return resources.CheckCondition(k8sClient, testDeployment, api.TrueCondition(api.ReadyType), validateDeploymentUpdatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})

			By("Add custom BackupSchedule for the Deployment", func() {
				Eventually(func(g Gomega) {
					deployment := &akov2.AtlasDeployment{}
					g.Expect(
						k8sClient.Get(context.Background(),
							kube.ObjectKeyFromObject(testDeployment),
							deployment),
					).To(Succeed())

					deployment.Spec.BackupScheduleRef = common.ResourceRefNamespaced{
						Name:      bsSchedule.Name,
						Namespace: bsSchedule.Namespace,
					}

					g.Expect(k8sClient.Update(context.Background(), deployment)).To(Succeed())
				}).WithTimeout(2 * time.Minute).WithPolling(20 * time.Second).Should(Succeed())
			})

			By("Deployment should be Ready", func() {
				Eventually(func(g Gomega) bool {
					return resources.CheckCondition(
						k8sClient,
						testDeployment,
						api.TrueCondition(api.DeploymentReadyType),
						validateDeploymentUpdatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})
		})
	})
