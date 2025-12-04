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
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas-sdk/v20250312009/admin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/httputil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/timeutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/version"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/access"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/conditions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/events"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/maintenance"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/resources"
	akoretry "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/retry"
)

const (
	ProjectCreationTimeout = 5 * time.Minute
)

var _ = Describe("AtlasProject", Label("int", "AtlasProject"), func() {
	const interval = time.Second * 2

	var (
		connectionSecret corev1.Secret
		createdProject   *akov2.AtlasProject
	)

	BeforeEach(func() {
		prepareControllers(false)

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

	checkIPAccessListInAtlas := func() {
		list, _, err := atlasClient.ProjectIPAccessListApi.
			ListAccessListEntries(context.Background(), createdProject.ID()).
			Execute()
		Expect(err).NotTo(HaveOccurred())

		Expect(list.GetTotalCount()).To(Equal(len(createdProject.Spec.ProjectIPAccessList)))
		Expect(list.GetResults()[0]).To(access.MatchIPAccessList(createdProject.Spec.ProjectIPAccessList[0]))
	}

	checkExpiredAccessLists := func(lists []project.IPAccessList) {
		currentStatusIPs := createdProject.Status.ExpiredIPAccessList
		if currentStatusIPs == nil {
			currentStatusIPs = []project.IPAccessList{}
		}
		Expect(currentStatusIPs).To(Equal(lists))
	}

	checkMaintenanceWindowInAtlas := func() {
		window, _, err := atlasClient.MaintenanceWindowsApi.
			GetMaintenanceWindow(context.Background(), createdProject.ID()).
			Execute()
		Expect(err).NotTo(HaveOccurred())
		Expect(window).To(maintenance.MatchMaintenanceWindow(createdProject.Spec.MaintenanceWindow))
	}

	checkAtlasProjectIsReady := func() {
		projectReadyConditions := conditions.MatchConditions(
			api.TrueCondition(api.ProjectReadyType),
			api.TrueCondition(api.ReadyType),
			api.TrueCondition(api.ValidationSucceeded),
		)
		Expect(createdProject.Status.ID).NotTo(BeNil())
		Expect(createdProject.Status.Conditions).To(ContainElements((projectReadyConditions)))
		Expect(createdProject.Status.ObservedGeneration).To(Equal(createdProject.Generation))
	}

	Describe("Creating the project", func() {
		It("Should Succeed", func() {
			expectedProject := akov2.DefaultProject(namespace.Name, connectionSecret.Name)
			createdProject.ObjectMeta = expectedProject.ObjectMeta
			Expect(k8sClient.Create(context.Background(), expectedProject)).ToNot(HaveOccurred())

			Eventually(func() bool {
				return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			checkAtlasProjectIsReady()

			// Atlas
			atlasProject, _, err := atlasClient.ProjectsApi.
				GetGroup(context.Background(), createdProject.Status.ID).
				Execute()
			Expect(err).ToNot(HaveOccurred())

			Expect(atlasProject.Name).To(Equal(expectedProject.Spec.Name))

			events.EventExists(k8sClient, createdProject, "Normal", "Ready", "")
		})
		It("Should Succeed with previous version of the operator", func() {
			version.Version = "1.0.0"
			expectedProject := akov2.DefaultProject(namespace.Name, connectionSecret.Name).WithLabels(map[string]string{
				customresource.ResourceVersion: "0.0.1",
			})
			createdProject.ObjectMeta = expectedProject.ObjectMeta
			Expect(k8sClient.Create(context.Background(), expectedProject)).ToNot(HaveOccurred())

			Eventually(func() bool {
				return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			expectedConditionsMatchers := conditions.MatchConditions(
				api.TrueCondition(api.ProjectReadyType),
				api.TrueCondition(api.ReadyType),
				api.TrueCondition(api.ValidationSucceeded),
				api.TrueCondition(api.ResourceVersionStatus),
			)
			Expect(createdProject.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))

			events.EventExists(k8sClient, createdProject, "Normal", "Ready", "")
		})
		It("Should Succeed with current version of the operator", func() {
			version.Version = "1.0.0"
			expectedProject := akov2.DefaultProject(namespace.Name, connectionSecret.Name).WithLabels(map[string]string{
				customresource.ResourceVersion: version.Version,
			})
			createdProject.ObjectMeta = expectedProject.ObjectMeta
			Expect(k8sClient.Create(context.Background(), expectedProject)).ToNot(HaveOccurred())

			Eventually(func() bool {
				return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			expectedConditionsMatchers := conditions.MatchConditions(
				api.TrueCondition(api.ProjectReadyType),
				api.TrueCondition(api.ReadyType),
				api.TrueCondition(api.ValidationSucceeded),
				api.TrueCondition(api.ResourceVersionStatus),
			)
			Expect(createdProject.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))

			events.EventExists(k8sClient, createdProject, "Normal", "Ready", "")
		})
		It("Should Fail with newer version of the operator", func() {
			version.Version = "1.0.0"
			expectedProject := akov2.DefaultProject(namespace.Name, connectionSecret.Name).WithLabels(map[string]string{
				customresource.ResourceVersion: "2.3.0",
			})
			createdProject.ObjectMeta = expectedProject.ObjectMeta
			Expect(k8sClient.Create(context.Background(), expectedProject)).ToNot(HaveOccurred())

			expectedCondition := api.FalseCondition(api.ResourceVersionStatus).WithReason(string(workflow.AtlasResourceVersionMismatch))
			Eventually(func() bool {
				return resources.CheckCondition(k8sClient, createdProject, expectedCondition)
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			Eventually(func(g Gomega) bool {
				expectedConditionsMatchers := conditions.MatchConditions(
					api.FalseCondition(api.ReadyType),
					api.FalseCondition(api.ResourceVersionStatus),
				)
				return g.Expect(createdProject.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())
		})
		It("Should Succeed with newer version of the operator and the override label", func() {
			version.Version = "1.0.0"
			expectedProject := akov2.DefaultProject(namespace.Name, connectionSecret.Name).WithLabels(map[string]string{
				customresource.ResourceVersion: "2.3.0",
			}).WithAnnotations(map[string]string{
				customresource.ResourceVersionOverride: customresource.ResourceVersionAllow,
			})
			createdProject.ObjectMeta = expectedProject.ObjectMeta
			Expect(k8sClient.Create(context.Background(), expectedProject)).ToNot(HaveOccurred())

			Eventually(func() bool {
				return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			expectedConditionsMatchers := conditions.MatchConditions(
				api.TrueCondition(api.ProjectReadyType),
				api.TrueCondition(api.ReadyType),
				api.TrueCondition(api.ValidationSucceeded),
				api.TrueCondition(api.ResourceVersionStatus),
			)
			Expect(createdProject.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))

			events.EventExists(k8sClient, createdProject, "Normal", "Ready", "")
		})
		It("Should fail if Secret is wrong", func() {
			expectedProject := akov2.DefaultProject(namespace.Name, "non-existent-secret")
			createdProject.ObjectMeta = expectedProject.ObjectMeta
			Expect(k8sClient.Create(context.Background(), expectedProject)).ToNot(HaveOccurred())

			expectedCondition := api.FalseCondition(api.ProjectReadyType).WithReason(string(workflow.AtlasAPIAccessNotConfigured))
			Eventually(func() bool {
				return resources.CheckCondition(k8sClient, createdProject, expectedCondition)
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			expectedConditionsMatchers := conditions.MatchConditions(
				api.FalseCondition(api.ProjectReadyType),
				api.FalseCondition(api.ReadyType),
				api.TrueCondition(api.ValidationSucceeded),
				api.TrueCondition(api.ResourceVersionStatus),
			)
			Expect(createdProject.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
			Expect(createdProject.ID()).To(BeEmpty())
			Expect(createdProject.Status.ObservedGeneration).To(Equal(createdProject.Generation))
			events.EventExists(k8sClient, createdProject, "Warning", string(workflow.AtlasAPIAccessNotConfigured), "Secret .* not found")

			// Atlas
			_, _, err := atlasClient.ProjectsApi.
				GetGroupByName(context.Background(), expectedProject.Spec.Name).
				Execute()

			// "NOT_IN_GROUP" is what is returned if the project is not found
			Expect(admin.IsErrorCode(err, atlas.NotInGroup)).To(BeTrue())
		})
	})

	Describe("Deleting the project (not cleaning Atlas)", func() {
		It("Should Succeed", func() {
			By(`Creating the project with retention policy "keep" first`, func() {
				createdProject = akov2.DefaultProject(namespace.Name, connectionSecret.Name)
				createdProject.ObjectMeta.Annotations = map[string]string{customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep}
				Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())
			})
			By("Deleting the project", func() {
				Expect(k8sClient.Delete(context.Background(), createdProject)).To(Succeed())
				time.Sleep(10 * time.Second)
				Expect(checkAtlasProjectRemoved(createdProject.Status.ID)()).Should(BeFalse())
			})
			By("Manually deleting the project from Atlas", func() {
				_, err := atlasClient.ProjectsApi.DeleteGroup(context.Background(), createdProject.ID()).Execute()
				Expect(err).ToNot(HaveOccurred())
				createdProject = nil
			})
		})
	})

	Describe("Deleting the project twice", func() {
		It("Should Succeed", func() {
			By(`Creating the project`, func() {
				createdProject = akov2.DefaultProject(namespace.Name, connectionSecret.Name)
				Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())
			})
			By("Deleting the project", func() {
				Expect(k8sClient.Delete(context.Background(), createdProject)).To(Succeed())
				Eventually(checkAtlasProjectRemoved(createdProject.Status.ID)).Should(BeTrue())
				time.Sleep(1 * time.Minute)
				Expect(checkAtlasProjectRemoved(createdProject.Status.ID)()).Should(BeTrue())
				createdProject = nil
			})
		})
	})

	Describe("Deleting the project several times", func() {
		// Should show that deleted project wasn't created again (depends on Atlas)
		It("Should Succeed", func() {
			const totalProject = 10
			var wg sync.WaitGroup
			wg.Add(totalProject)
			createdProjects := make([]*akov2.AtlasProject, totalProject)
			projectPrefix := fmt.Sprintf("project-%s", namespace.Name)

			By("Creating global key", func() {
				globalConnectionSecret := buildConnectionSecret("atlas-operator-api-key")
				Expect(k8sClient.Create(context.Background(), &globalConnectionSecret)).To(Succeed())
			})

			for i := 0; i < totalProject; i++ {
				go func(i int) {
					defer GinkgoRecover()
					defer wg.Done()
					projectName := fmt.Sprintf("%s-%v", projectPrefix, i)

					By(fmt.Sprintf("Creating several projects: %s", projectName))
					createdProjects[i] = akov2.DefaultProject(namespace.Name, "").WithAtlasName(projectName).WithName(projectName)
					Expect(k8sClient.Create(context.Background(), createdProjects[i])).ShouldNot(HaveOccurred())
					GinkgoWriter.Write([]byte(fmt.Sprintf("%+v", createdProjects[i])))

					Eventually(func() bool {
						return resources.CheckCondition(k8sClient, createdProjects[i], api.TrueCondition(api.ReadyType))
					}).WithTimeout(5 * time.Minute).WithPolling(interval).Should(BeTrue())

					By(fmt.Sprintf("Deleting the project: %s", projectName))
					Expect(k8sClient.Delete(context.Background(), createdProjects[i])).Should(Succeed())
					GinkgoWriter.Write([]byte(fmt.Sprintf("%+v\n", createdProjects[i])))
					GinkgoWriter.Write([]byte(fmt.Sprintf("%v=======================NAME: %s\n", i, projectName)))
					GinkgoWriter.Write([]byte(fmt.Sprintf("%v=========================ID: %s\n", i, createdProjects[i].Status.ID)))
					Eventually(checkAtlasProjectRemoved(createdProjects[i].Status.ID), 2*time.Minute, 5*time.Second).Should(BeTrue())

					By(fmt.Sprintf("Check if project wasn't created again: %s", projectName))
					time.Sleep(1 * time.Minute)
					GinkgoWriter.Write([]byte(fmt.Sprintf("%+v\n", createdProjects[i])))
					GinkgoWriter.Write([]byte(fmt.Sprintf("%v=======================NAME: %s\n", i, projectName)))
					GinkgoWriter.Write([]byte(fmt.Sprintf("%v=========================ID: %s\n", i, createdProjects[i].Status.ID)))
					Expect(checkAtlasProjectRemoved(createdProjects[i].Status.ID)()).Should(BeTrue())
				}(i)
			}
			wg.Wait()
			createdProject = nil
		})
	})

	Describe("Updating the project", func() {
		It("Should Succeed", func() {
			By("Creating the project first")

			expectedProject := akov2.DefaultProject(namespace.Name, connectionSecret.Name)
			createdProject.ObjectMeta = expectedProject.ObjectMeta
			Expect(k8sClient.Create(context.Background(), expectedProject)).ToNot(HaveOccurred())

			Eventually(func() bool {
				return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			// Updating (the existing project is expected to be read from Atlas)
			By("Updating the project")

			var err error
			createdProject, err = akoretry.RetryUpdateOnConflict(context.Background(), k8sClient, client.ObjectKeyFromObject(createdProject), func(p *akov2.AtlasProject) {
				p.Spec.ProjectIPAccessList = []project.IPAccessList{{CIDRBlock: "0.0.0.0/0"}}
				p.Spec.MaintenanceWindow = project.MaintenanceWindow{
					DayOfWeek: 4,
					HourOfDay: 11,
					AutoDefer: true,
					StartASAP: false,
					Defer:     false,
				}
			})
			Expect(err).To(BeNil())

			Eventually(func() bool {
				return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			Expect(resources.ReadAtlasResource(context.Background(), k8sClient, createdProject)).To(BeTrue())
			Expect(createdProject.Status.Conditions).To(ContainElement(conditions.MatchCondition(api.TrueCondition(api.ProjectReadyType))))

			// Atlas
			atlasProject, _, err := atlasClient.ProjectsApi.GetGroup(context.Background(), createdProject.ID()).Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(atlasProject.Name).To(Equal(expectedProject.Spec.Name))
		})
	})

	Describe("Two projects watching the Connection Secret", func() {
		var secondProject *akov2.AtlasProject
		It("Should Succeed", func() {
			By("Creating two projects first", func() {
				createdProject = akov2.DefaultProject(namespace.Name, connectionSecret.Name).WithName("first-project")
				Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

				secondProject = akov2.DefaultProject(namespace.Name, connectionSecret.Name).WithName("second-project").WithAtlasName("second-project")
				Expect(k8sClient.Create(context.Background(), secondProject)).ToNot(HaveOccurred())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, secondProject, api.TrueCondition(api.ReadyType))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())
			})
			By("Breaking the Connection Secret", func() {
				connectionSecret = buildConnectionSecret("my-atlas-key")
				_, err := akoretry.RetryUpdateOnConflict(context.Background(), k8sClient, client.ObjectKeyFromObject(&connectionSecret), func(s *corev1.Secret) {
					s.StringData = buildConnectionSecret("my-atlas-key").StringData
					s.StringData["publicApiKey"] = "non-existing"
				})
				Expect(err).To(BeNil())

				createdProject, err = akoretry.RetryUpdateOnConflict(context.Background(), k8sClient, client.ObjectKeyFromObject(createdProject), func(p *akov2.AtlasProject) {
					p.Spec.AlertConfigurationSyncEnabled = true
				})
				Expect(err).To(BeNil())

				secondProject, err = akoretry.RetryUpdateOnConflict(context.Background(), k8sClient, client.ObjectKeyFromObject(secondProject), func(p *akov2.AtlasProject) {
					p.Spec.AlertConfigurationSyncEnabled = true
				})
				Expect(err).To(BeNil())

				// Both projects are expected to get to Failed state right away
				expectedCondition := api.FalseCondition(api.ProjectReadyType).WithReason(string(workflow.ProjectNotCreatedInAtlas))
				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, createdProject, expectedCondition)
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())
				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, secondProject, expectedCondition)
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())
			})
			By("Fixing the Connection Secret", func() {
				connectionSecret = buildConnectionSecret("my-atlas-key")
				_, err := akoretry.RetryUpdateOnConflict(context.Background(), k8sClient, client.ObjectKeyFromObject(&connectionSecret), func(s *corev1.Secret) {
					s.StringData = buildConnectionSecret("my-atlas-key").StringData
				})
				Expect(err).To(BeNil())

				// Both projects are expected to recover
				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())
				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, secondProject, api.TrueCondition(api.ReadyType))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())
			})
			By("Removing (second) Atlas Project "+secondProject.Status.ID, func() {
				if secondProject != nil && secondProject.Status.ID != "" {
					Expect(k8sClient.Delete(context.Background(), secondProject)).To(Succeed())
					Eventually(checkAtlasProjectRemoved(secondProject.Status.ID), 20, interval).Should(BeTrue())
				}
			})
		})
	})

	Describe("Creating the project IP access list", func() {
		It("Should Succeed (single)", func() {
			createdProject = akov2.DefaultProject(namespace.Name, connectionSecret.Name).
				WithIPAccessList(project.NewIPAccessList().WithComment("bla").WithIP("192.0.2.15"))
			Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

			Eventually(func(g Gomega) bool {
				return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType), validateNoErrorsIPAccessListDuringCreate(g))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			checkAtlasProjectIsReady()
			checkExpiredAccessLists([]project.IPAccessList{})
			checkIPAccessListInAtlas()
		})
		It("Should Succeed (multiple)", func() {
			tenHoursLater := time.Now().Add(time.Hour * 10).Format("2006-01-02T15:04:05-0700")
			createdProject = akov2.DefaultProject(namespace.Name, connectionSecret.Name).
				WithIPAccessList(project.NewIPAccessList().WithComment("bla").WithCIDR("203.0.113.0/24").WithDeleteAfterDate(tenHoursLater)).
				WithIPAccessList(project.NewIPAccessList().WithComment("foo").WithIP("192.0.2.20"))

			Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

			Eventually(func(g Gomega) bool {
				return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType), validateNoErrorsIPAccessListDuringCreate(g))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			checkAtlasProjectIsReady()
			checkExpiredAccessLists([]project.IPAccessList{})
			checkIPAccessListInAtlas()
		})
		It("Should Succeed (1 expired)", func() {
			tenHoursBefore := time.Now().Add(time.Hour * -10)
			expiredList := project.IPAccessList{Comment: "bla", CIDRBlock: "203.0.113.0/24", DeleteAfterDate: timeutil.FormatISO8601(tenHoursBefore)}
			activeList := project.IPAccessList{Comment: "foo", IPAddress: "192.0.2.20"}

			createdProject = akov2.DefaultProject(namespace.Name, connectionSecret.Name).WithIPAccessList(expiredList).WithIPAccessList(activeList)

			Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

			Eventually(func(g Gomega) bool {
				return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType), validateNoErrorsIPAccessListDuringCreate(g))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			checkAtlasProjectIsReady()
			checkExpiredAccessLists([]project.IPAccessList{expiredList})

			// Atlas
			list, _, err := atlasClient.ProjectIPAccessListApi.
				ListAccessListEntries(context.Background(), createdProject.ID()).
				Execute()
			Expect(err).NotTo(HaveOccurred())

			Expect(list.GetTotalCount()).To(Equal(1))
			Expect(list.GetResults()[0]).To(access.MatchIPAccessList(activeList))
		})
		It("Should Fail (AWS security group not supported without VPC)", func() {
			createdProject = akov2.DefaultProject(namespace.Name, connectionSecret.Name).
				WithIPAccessList(project.NewIPAccessList().WithAWSGroup("sg-0026348ec11780bd1"))

			Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

			Eventually(func() bool {
				return resources.CheckCondition(k8sClient, createdProject, api.FalseCondition(api.IPAccessListReadyType))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			ipAccessFailedCondition := api.FalseCondition(api.IPAccessListReadyType).
				WithReason(string(workflow.ProjectIPNotCreatedInAtlas)).
				WithMessageRegexp(".*CANNOT_USE_AWS_SECURITY_GROUP_WITHOUT_VPC_PEERING_CONNECTION.*")

			expectedConditionsMatchers := conditions.MatchConditions(
				api.TrueCondition(api.ProjectReadyType),
				api.TrueCondition(api.ValidationSucceeded),
				ipAccessFailedCondition,
				api.FalseCondition(api.ReadyType),
				api.TrueCondition(api.ResourceVersionStatus),
			)
			Expect(createdProject.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
			checkExpiredAccessLists([]project.IPAccessList{})
		})
	})

	Describe("Updating the project IP access list", func() {
		It("Should Succeed (single)", func() {
			By("Creating the project first", func() {
				tenMinutesLater := time.Now().Add(time.Minute * 10).Format("2006-01-02T15:04:05-0700")
				createdProject = akov2.DefaultProject(namespace.Name, connectionSecret.Name).
					WithIPAccessList(project.NewIPAccessList().WithComment("bla").WithIP("192.0.2.15").WithDeleteAfterDate(tenMinutesLater))

				Expect(k8sClient.Create(context.Background(), createdProject)).To(Succeed())

				Eventually(func(g Gomega) bool {
					return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType), validateNoErrorsIPAccessListDuringCreate(g))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())
			})
			By("Updating the IP Access List comment and delete date", func() {
				// Just a note: Atlas doesn't allow to make the "permanent" entity "temporary". But it works the other way
				var err error
				createdProject, err = akoretry.RetryUpdateOnConflict(context.Background(), k8sClient, client.ObjectKeyFromObject(createdProject), func(p *akov2.AtlasProject) {
					p.Spec.ProjectIPAccessList[0].Comment = "new comment"
					p.Spec.ProjectIPAccessList[0].DeleteAfterDate = ""
				})
				Expect(err).To(BeNil())

				Eventually(func(g Gomega) bool {
					return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType), validateNoErrorsIPAccessListDuringUpdate(g))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

				checkAtlasProjectIsReady()
				checkExpiredAccessLists([]project.IPAccessList{})
				checkIPAccessListInAtlas()
			})
		})

		It("Should Succeed (multiple)", func() {
			By("Creating the project first", func() {
				thirtyHoursLater := time.Now().Add(time.Hour * 30).Format("2006-01-02T15:04:05-0700")
				createdProject = akov2.DefaultProject(namespace.Name, connectionSecret.Name).
					WithIPAccessList(project.NewIPAccessList().WithComment("bla").WithCIDR("203.0.113.0/24").WithDeleteAfterDate(thirtyHoursLater)).
					WithIPAccessList(project.NewIPAccessList().WithComment("foo").WithIP("192.0.2.20"))

				Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

				Eventually(func(g Gomega) bool {
					return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType), validateNoErrorsIPAccessListDuringCreate(g))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())
			})
			By("Updating the IP Access List IPAddress", func() {
				twoDaysLater := time.Now().Add(time.Hour * 48).Format("2006-01-02T15:04:05Z")
				var err error
				createdProject, err = akoretry.RetryUpdateOnConflict(context.Background(), k8sClient, client.ObjectKeyFromObject(createdProject), func(p *akov2.AtlasProject) {
					p.Spec.ProjectIPAccessList[0].DeleteAfterDate = twoDaysLater
					// Update of the IP address will result in delete for the old IP address first and then the new
					// IP address will be created
					p.Spec.ProjectIPAccessList[1].IPAddress = "168.32.54.0"
				})
				Expect(err).To(BeNil())

				Eventually(func(g Gomega) bool {
					return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType), validateNoErrorsIPAccessListDuringUpdate(g))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

				checkAtlasProjectIsReady()
				checkExpiredAccessLists([]project.IPAccessList{})
				checkIPAccessListInAtlas()
			})
		})
	})

	// Here we do not test defer and startASAP requests because we cannot check the maintenance status from the API
	Describe("Updating the project Maintenance Window", func() {
		It("Should Succeed (single)", func() {
			By("Creating the project first", func() {
				createdProject = akov2.DefaultProject(namespace.Name, connectionSecret.Name).
					WithMaintenanceWindow(project.MaintenanceWindow{
						DayOfWeek: 2,
						HourOfDay: 2,
						AutoDefer: false,
					})

				Expect(k8sClient.Create(context.Background(), createdProject)).To(Succeed())

				Eventually(func(g Gomega) bool {
					return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType), validateNoErrorsMaintenanceWindowDuringCreate(g))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			})
			By("Updating the project maintenance window hour and enabling auto-defer", func() {
				var err error
				createdProject, err = akoretry.RetryUpdateOnConflict(context.Background(), k8sClient, client.ObjectKeyFromObject(createdProject), func(p *akov2.AtlasProject) {
					p.Spec.MaintenanceWindow.HourOfDay = 3
					p.Spec.MaintenanceWindow.AutoDefer = true
				})
				Expect(err).To(BeNil())

				Eventually(func(g Gomega) bool {
					return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType), validateNoErrorsMaintenanceWindowDuringUpdate(g))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

				// TODO: Refactor check functions to use Eventually assertion
				time.Sleep(10 * time.Second)
				checkAtlasProjectIsReady()
				checkMaintenanceWindowInAtlas()
			})
			By("Toggling auto-defer to false", func() {
				var err error
				createdProject, err = akoretry.RetryUpdateOnConflict(context.Background(), k8sClient, client.ObjectKeyFromObject(createdProject), func(p *akov2.AtlasProject) {
					p.Spec.MaintenanceWindow.AutoDefer = false
				})
				Expect(err).To(BeNil())

				Eventually(func(g Gomega) bool {
					return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType), validateNoErrorsMaintenanceWindowDuringUpdate(g))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

				checkAtlasProjectIsReady()
				checkMaintenanceWindowInAtlas()
			})
		})
	})

	Describe("Using the global Connection Secret", func() {
		It("Should Succeed", func() {
			globalConnectionSecret := buildConnectionSecret("atlas-operator-api-key")
			Expect(k8sClient.Create(context.Background(), &globalConnectionSecret)).To(Succeed())

			// We don't specify the connection Secret per project - the global one must be used
			createdProject = akov2.DefaultProject(namespace.Name, "")

			Expect(k8sClient.Create(context.Background(), createdProject)).To(Succeed())

			Eventually(func() bool {
				return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			expectedConditionsMatchers := conditions.MatchConditions(
				api.TrueCondition(api.ProjectReadyType),
				api.TrueCondition(api.ReadyType),
				api.TrueCondition(api.ValidationSucceeded),
			)
			Expect(createdProject.Status.Conditions).To(ContainElements(expectedConditionsMatchers))
			Expect(createdProject.Status.ObservedGeneration).To(Equal(createdProject.Generation))
		})
		It("Should Fail if the global Secret doesn't exist", func() {
			By("Creating without a global Secret", func() {
				createdProject = akov2.DefaultProject(namespace.Name, "").WithName("project-no-secret")

				Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, createdProject, api.FalseCondition(api.ReadyType))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

				expectedConditionsMatchers := conditions.MatchConditions(
					api.FalseCondition(api.ProjectReadyType).
						WithReason(string(workflow.AtlasAPIAccessNotConfigured)),
					api.FalseCondition(api.ReadyType),
					api.TrueCondition(api.ValidationSucceeded),
					api.TrueCondition(api.ResourceVersionStatus),
				)
				Expect(createdProject.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
				Expect(createdProject.ID()).To(BeEmpty())
				Expect(createdProject.Status.ObservedGeneration).To(Equal(createdProject.Generation))
			})
			By("Creating a global Secret - should get fixed", func() {
				globalConnectionSecret := buildConnectionSecret("atlas-operator-api-key")
				Expect(k8sClient.Create(context.Background(), &globalConnectionSecret)).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())
			})
		})
	})
})

func buildConnectionSecret(name string) corev1.Secret {
	return corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace.Name,
			Labels: map[string]string{
				"atlas.mongodb.com/type": "credentials",
			},
		},
		StringData: secretData(),
	}
}

// checkAtlasProjectRemoved returns true if the Atlas Project is removed from Atlas.
func checkAtlasProjectRemoved(projectID string) func() bool {
	return func() bool {
		_, r, err := atlasClient.ProjectsApi.GetGroup(context.Background(), projectID).Execute()
		if err != nil {
			statusCode := httputil.StatusCode(r)
			if statusCode == http.StatusNotFound || statusCode == http.StatusUnauthorized {
				return true
			}
		}
		return false
	}
}

// validateNoErrorsIPAccessListDuringCreate performs check that no problems happen to IP Access list during the creation.
// This allows the test to fail fast instead by timeout if there are any troubles.
func validateNoErrorsIPAccessListDuringCreate(g Gomega) func(a api.AtlasCustomResource) {
	return func(a api.AtlasCustomResource) {
		c := a.(*akov2.AtlasProject)

		if condition, ok := conditions.FindConditionByType(c.Status.Conditions, api.IPAccessListReadyType); ok {
			g.Expect(condition.Status).To(Equal(api.TrueCondition(api.IPAccessListReadyType).Status), fmt.Sprintf("Unexpected condition: %v", condition))
		}
	}
}

// validateNoErrorsIPAccessListDuringUpdate performs check that no problems happen to IP Access list during the update.
func validateNoErrorsIPAccessListDuringUpdate(g Gomega) func(a api.AtlasCustomResource) {
	return func(a api.AtlasCustomResource) {
		c := a.(*akov2.AtlasProject)
		condition, ok := conditions.FindConditionByType(c.Status.Conditions, api.IPAccessListReadyType)
		g.Expect(ok).To(BeTrue())
		g.Expect(condition.Reason).To(BeEmpty())
	}
}

// validateNoErrorsMaintenanceWindowDuringCreate performs check that no problems happen to Maintenance Window during the creation.
// This allows the test to fail fast instead by timeout if there are any troubles.
func validateNoErrorsMaintenanceWindowDuringCreate(g Gomega) func(a api.AtlasCustomResource) {
	return func(a api.AtlasCustomResource) {
		c := a.(*akov2.AtlasProject)

		if condition, ok := conditions.FindConditionByType(c.Status.Conditions, api.MaintenanceWindowReadyType); ok {
			g.Expect(condition.Status).To(Equal(api.TrueCondition(api.MaintenanceWindowReadyType).Status), fmt.Sprintf("Unexpected condition: %v", condition))
		}
	}
}

// validateNoErrorsMaintenanceWindowDuringUpdate performs check that no problems happen to Maintenance Window during the update.
func validateNoErrorsMaintenanceWindowDuringUpdate(g Gomega) func(a api.AtlasCustomResource) {
	return func(a api.AtlasCustomResource) {
		c := a.(*akov2.AtlasProject)
		condition, ok := conditions.FindConditionByType(c.Status.Conditions, api.MaintenanceWindowReadyType)
		g.Expect(ok).To(BeTrue())
		g.Expect(condition.Reason).To(BeEmpty())
	}
}

func deleteK8sObject(obj client.Object) func() bool {
	return func() bool {
		nn := kube.ObjectKeyFromObject(obj)
		GinkgoWriter.Printf("Deleting %s/%s\n", nn.Namespace, nn.Name)
		err := k8sClient.Get(context.Background(), nn, obj)
		if err == nil {
			err = k8sClient.Delete(context.Background(), obj)
		}
		if err != nil {
			GinkgoWriter.Printf("Attempt to delete %s/%s failed: %v\n", nn.Namespace, nn.Name, err)
			return false
		}
		GinkgoWriter.Printf("Deleted %s/%s\n", nn.Namespace, nn.Name)
		return true
	}
}
