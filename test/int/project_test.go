package int

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/version"

	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"go.mongodb.org/atlas/mongodbatlas"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/testutil"
)

const (
	ProjectCreationTimeout = 40 * time.Second
)

var _ = Describe("AtlasProject", Label("int", "AtlasProject"), func() {
	const interval = time.Second * 2

	var (
		connectionSecret corev1.Secret
		createdProject   *mdbv1.AtlasProject
	)

	BeforeEach(func() {
		prepareControllers()

		createdProject = &mdbv1.AtlasProject{}

		connectionSecret = buildConnectionSecret("my-atlas-key")
		By(fmt.Sprintf("Creating the Secret %s", kube.ObjectKeyFromObject(&connectionSecret)))
		Expect(k8sClient.Create(context.Background(), &connectionSecret)).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		if createdProject != nil && createdProject.Status.ID != "" {
			By("Removing Atlas Project " + createdProject.Status.ID)
			Expect(k8sClient.Delete(context.Background(), createdProject)).To(Succeed())
			Eventually(checkAtlasProjectRemoved(createdProject.Status.ID), 20, interval).Should(BeTrue())
		}
		removeControllersAndNamespace()
	})

	checkIPAccessListInAtlas := func() {
		list, _, err := atlasClient.ProjectIPAccessList.List(context.Background(), createdProject.ID(), &mongodbatlas.ListOptions{})
		Expect(err).NotTo(HaveOccurred())

		Expect(list.Results).To(HaveLen(len(createdProject.Spec.ProjectIPAccessList)))
		Expect(list.Results[0]).To(testutil.MatchIPAccessList(createdProject.Spec.ProjectIPAccessList[0]))
	}

	checkExpiredAccessLists := func(lists []project.IPAccessList) {
		currentStatusIPs := createdProject.Status.ExpiredIPAccessList
		if currentStatusIPs == nil {
			currentStatusIPs = []project.IPAccessList{}
		}
		Expect(currentStatusIPs).To(Equal(lists))
	}

	checkMaintenanceWindowInAtlas := func() {
		window, _, err := atlasClient.MaintenanceWindows.Get(context.Background(), createdProject.ID())
		Expect(err).NotTo(HaveOccurred())
		Expect(window).To(testutil.MatchMaintenanceWindow(createdProject.Spec.MaintenanceWindow))
	}

	checkAtlasProjectIsReady := func() {
		projectReadyConditions := testutil.MatchConditions(
			status.TrueCondition(status.ProjectReadyType),
			status.TrueCondition(status.ReadyType),
			status.TrueCondition(status.ValidationSucceeded),
		)
		Expect(createdProject.Status.ID).NotTo(BeNil())
		Expect(createdProject.Status.Conditions).To(ContainElements((projectReadyConditions)))
		Expect(createdProject.Status.ObservedGeneration).To(Equal(createdProject.Generation))
	}

	Describe("Creating the project", func() {
		It("Should Succeed", func() {
			expectedProject := mdbv1.DefaultProject(namespace.Name, connectionSecret.Name)
			createdProject.ObjectMeta = expectedProject.ObjectMeta
			Expect(k8sClient.Create(context.Background(), expectedProject)).ToNot(HaveOccurred())

			Eventually(func() bool {
				return testutil.CheckCondition(k8sClient, createdProject, status.TrueCondition(status.ReadyType))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			checkAtlasProjectIsReady()

			// Atlas
			atlasProject, _, err := atlasClient.Projects.GetOneProject(context.Background(), createdProject.Status.ID)
			Expect(err).ToNot(HaveOccurred())

			Expect(atlasProject.Name).To(Equal(expectedProject.Spec.Name))

			testutil.EventExists(k8sClient, createdProject, "Normal", "Ready", "")
		})
		It("Should Succeed with previous version of the operator", func() {
			version.Version = "1.0.0"
			expectedProject := mdbv1.DefaultProject(namespace.Name, connectionSecret.Name).WithLabels(map[string]string{
				customresource.ResourceVersion: "0.0.1",
			})
			createdProject.ObjectMeta = expectedProject.ObjectMeta
			Expect(k8sClient.Create(context.Background(), expectedProject)).ToNot(HaveOccurred())

			Eventually(func() bool {
				return testutil.CheckCondition(k8sClient, createdProject, status.TrueCondition(status.ReadyType))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			expectedConditionsMatchers := testutil.MatchConditions(
				status.TrueCondition(status.ProjectReadyType),
				status.TrueCondition(status.ReadyType),
				status.TrueCondition(status.ValidationSucceeded),
				status.TrueCondition(status.ResourceVersionStatus),
			)
			Expect(createdProject.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))

			testutil.EventExists(k8sClient, createdProject, "Normal", "Ready", "")
		})
		It("Should Succeed with current version of the operator", func() {
			version.Version = "1.0.0"
			expectedProject := mdbv1.DefaultProject(namespace.Name, connectionSecret.Name).WithLabels(map[string]string{
				customresource.ResourceVersion: version.Version,
			})
			createdProject.ObjectMeta = expectedProject.ObjectMeta
			Expect(k8sClient.Create(context.Background(), expectedProject)).ToNot(HaveOccurred())

			Eventually(func() bool {
				return testutil.CheckCondition(k8sClient, createdProject, status.TrueCondition(status.ReadyType))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			expectedConditionsMatchers := testutil.MatchConditions(
				status.TrueCondition(status.ProjectReadyType),
				status.TrueCondition(status.ReadyType),
				status.TrueCondition(status.ValidationSucceeded),
				status.TrueCondition(status.ResourceVersionStatus),
			)
			Expect(createdProject.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))

			testutil.EventExists(k8sClient, createdProject, "Normal", "Ready", "")
		})
		It("Should Fail with newer version of the operator", func() {
			version.Version = "1.0.0"
			expectedProject := mdbv1.DefaultProject(namespace.Name, connectionSecret.Name).WithLabels(map[string]string{
				customresource.ResourceVersion: "2.3.0",
			})
			createdProject.ObjectMeta = expectedProject.ObjectMeta
			Expect(k8sClient.Create(context.Background(), expectedProject)).ToNot(HaveOccurred())

			expectedCondition := status.FalseCondition(status.ResourceVersionStatus).WithReason(string(workflow.AtlasResourceVersionMismatch))
			Eventually(func() bool {
				return testutil.CheckCondition(k8sClient, createdProject, expectedCondition)
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			Eventually(func(g Gomega) bool {
				expectedConditionsMatchers := testutil.MatchConditions(
					status.FalseCondition(status.ReadyType),
					status.FalseCondition(status.ResourceVersionStatus),
				)
				return g.Expect(createdProject.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())
		})
		It("Should Succeed with newer version of the operator and the override label", func() {
			version.Version = "1.0.0"
			expectedProject := mdbv1.DefaultProject(namespace.Name, connectionSecret.Name).WithLabels(map[string]string{
				customresource.ResourceVersion: "2.3.0",
			}).WithAnnotations(map[string]string{
				customresource.ResourceVersionOverride: customresource.ResourceVersionAllow,
			})
			createdProject.ObjectMeta = expectedProject.ObjectMeta
			Expect(k8sClient.Create(context.Background(), expectedProject)).ToNot(HaveOccurred())

			Eventually(func() bool {
				return testutil.CheckCondition(k8sClient, createdProject, status.TrueCondition(status.ReadyType))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			expectedConditionsMatchers := testutil.MatchConditions(
				status.TrueCondition(status.ProjectReadyType),
				status.TrueCondition(status.ReadyType),
				status.TrueCondition(status.ValidationSucceeded),
				status.TrueCondition(status.ResourceVersionStatus),
			)
			Expect(createdProject.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))

			testutil.EventExists(k8sClient, createdProject, "Normal", "Ready", "")
		})
		It("Should fail if Secret is wrong", func() {
			expectedProject := mdbv1.DefaultProject(namespace.Name, "non-existent-secret")
			createdProject.ObjectMeta = expectedProject.ObjectMeta
			Expect(k8sClient.Create(context.Background(), expectedProject)).ToNot(HaveOccurred())

			expectedCondition := status.FalseCondition(status.ProjectReadyType).WithReason(string(workflow.AtlasCredentialsNotProvided))
			Eventually(func() bool {
				return testutil.CheckCondition(k8sClient, createdProject, expectedCondition)
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			expectedConditionsMatchers := testutil.MatchConditions(
				status.FalseCondition(status.ProjectReadyType),
				status.FalseCondition(status.ReadyType),
				status.TrueCondition(status.ValidationSucceeded),
				status.TrueCondition(status.ResourceVersionStatus),
			)
			Expect(createdProject.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
			Expect(createdProject.ID()).To(BeEmpty())
			Expect(createdProject.Status.ObservedGeneration).To(Equal(createdProject.Generation))
			testutil.EventExists(k8sClient, createdProject, "Warning", string(workflow.AtlasCredentialsNotProvided), "Secret .* not found")

			// Atlas
			_, _, err := atlasClient.Projects.GetOneProjectByName(context.Background(), expectedProject.Spec.Name)

			// "NOT_IN_GROUP" is what is returned if the project is not found
			var apiError *mongodbatlas.ErrorResponse
			Expect(errors.As(err, &apiError)).To(BeTrue(), "Error occurred: "+err.Error())
			Expect(apiError.ErrorCode).To(Equal(atlas.NotInGroup))
		})
	})

	Describe("Deleting the project (not cleaning Atlas)", func() {
		It("Should Succeed", func() {
			By(`Creating the project with retention policy "keep" first`, func() {
				createdProject = mdbv1.DefaultProject(namespace.Name, connectionSecret.Name)
				createdProject.ObjectMeta.Annotations = map[string]string{customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep}
				Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdProject, status.TrueCondition(status.ReadyType))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())
			})
			By("Deleting the project", func() {
				Expect(k8sClient.Delete(context.Background(), createdProject)).To(Succeed())
				time.Sleep(10 * time.Second)
				Expect(checkAtlasProjectRemoved(createdProject.Status.ID)()).Should(BeFalse())
			})
			By("Manually deleting the project from Atlas", func() {
				_, _ = atlasClient.Projects.Delete(context.Background(), createdProject.ID())
				createdProject = nil
			})
		})
	})

	Describe("Deleting the project twice", func() {
		It("Should Succeed", func() {
			By(`Creating the project`, func() {
				createdProject = mdbv1.DefaultProject(namespace.Name, connectionSecret.Name)
				Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdProject, status.TrueCondition(status.ReadyType))
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
			createdProjects := make([]*mdbv1.AtlasProject, totalProject)
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
					createdProjects[i] = mdbv1.DefaultProject(namespace.Name, "").WithAtlasName(projectName).WithName(projectName)
					Expect(k8sClient.Create(context.Background(), createdProjects[i])).ShouldNot(HaveOccurred())
					GinkgoWriter.Write([]byte(fmt.Sprintf("%+v", createdProjects[i])))

					Eventually(func() bool {
						return testutil.CheckCondition(k8sClient, createdProjects[i], status.TrueCondition(status.ReadyType))
					}).WithTimeout(5 * time.Minute).WithPolling(interval).Should(BeTrue())

					By(fmt.Sprintf("Deleting the project: %s", projectName))
					Expect(k8sClient.Delete(context.Background(), createdProjects[i])).Should(Succeed())
					GinkgoWriter.Write([]byte(fmt.Sprintf("%+v\n", createdProjects[i])))
					GinkgoWriter.Write([]byte(fmt.Sprintf("%v=======================NAME: %s\n", i, projectName)))
					GinkgoWriter.Write([]byte(fmt.Sprintf("%v=========================ID: %s\n", i, createdProjects[i].Status.ID)))
					Eventually(checkAtlasProjectRemoved(createdProjects[i].Status.ID), 1*time.Minute, 5*time.Second).Should(BeTrue())

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

			expectedProject := mdbv1.DefaultProject(namespace.Name, connectionSecret.Name)
			createdProject.ObjectMeta = expectedProject.ObjectMeta
			Expect(k8sClient.Create(context.Background(), expectedProject)).ToNot(HaveOccurred())

			Eventually(func() bool {
				return testutil.CheckCondition(k8sClient, createdProject, status.TrueCondition(status.ReadyType))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			// Updating (the existing project is expected to be read from Atlas)
			By("Updating the project")

			createdProject.Spec.ProjectIPAccessList = []project.IPAccessList{{CIDRBlock: "0.0.0.0/0"}}
			createdProject.Spec.MaintenanceWindow = project.MaintenanceWindow{
				DayOfWeek: 4,
				HourOfDay: 11,
				AutoDefer: true,
				StartASAP: false,
				Defer:     false,
			}
			Expect(k8sClient.Update(context.Background(), createdProject)).To(Succeed())

			Eventually(func() bool {
				return testutil.CheckCondition(k8sClient, createdProject, status.TrueCondition(status.ReadyType))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			Expect(testutil.ReadAtlasResource(k8sClient, createdProject)).To(BeTrue())
			Expect(createdProject.Status.Conditions).To(ContainElement(testutil.MatchCondition(status.TrueCondition(status.ProjectReadyType))))

			// Atlas
			atlasProject, _, err := atlasClient.Projects.GetOneProject(context.Background(), createdProject.ID())
			Expect(err).ToNot(HaveOccurred())
			Expect(atlasProject.Name).To(Equal(expectedProject.Spec.Name))
		})
	})

	Describe("Two projects watching the Connection Secret", func() {
		var secondProject *mdbv1.AtlasProject
		AfterEach(func() {
			if secondProject != nil && secondProject.Status.ID != "" {
				By("Removing (second) Atlas Project " + secondProject.Status.ID)
				Expect(k8sClient.Delete(context.Background(), secondProject)).To(Succeed())
				Eventually(checkAtlasProjectRemoved(secondProject.Status.ID), 20, interval).Should(BeTrue())
			}
		})
		It("Should Succeed", func() {
			By("Creating two projects first", func() {
				createdProject = mdbv1.DefaultProject(namespace.Name, connectionSecret.Name)
				Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

				secondProject = mdbv1.DefaultProject(namespace.Name, connectionSecret.Name).WithName("second-project").WithAtlasName("second Project")
				Expect(k8sClient.Create(context.Background(), secondProject)).ToNot(HaveOccurred())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdProject, status.TrueCondition(status.ReadyType))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, secondProject, status.TrueCondition(status.ReadyType))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())
			})
			By("Breaking the Connection Secret", func() {
				connectionSecret = buildConnectionSecret("my-atlas-key")
				connectionSecret.StringData["publicApiKey"] = "non-existing"
				Expect(k8sClient.Update(context.Background(), &connectionSecret)).To(Succeed())

				// Both projects are expected to get to Failed state right away
				expectedCondition := status.FalseCondition(status.ProjectReadyType).WithReason(string(workflow.ProjectNotCreatedInAtlas))
				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdProject, expectedCondition)
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())
				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, secondProject, expectedCondition)
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())
			})
			By("Fixing the Connection Secret", func() {
				connectionSecret = buildConnectionSecret("my-atlas-key")
				Expect(k8sClient.Update(context.Background(), &connectionSecret)).To(Succeed())

				// Both projects are expected to recover
				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdProject, status.TrueCondition(status.ReadyType))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())
				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, secondProject, status.TrueCondition(status.ReadyType))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())
			})
		})
	})

	Describe("Creating the project IP access list", func() {
		It("Should Succeed (single)", func() {
			createdProject = mdbv1.DefaultProject(namespace.Name, connectionSecret.Name).
				WithIPAccessList(project.NewIPAccessList().WithComment("bla").WithIP("192.0.2.15"))
			Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

			Eventually(func(g Gomega) bool {
				return testutil.CheckCondition(k8sClient, createdProject, status.TrueCondition(status.ReadyType), validateNoErrorsIPAccessListDuringCreate(g))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			checkAtlasProjectIsReady()
			checkExpiredAccessLists([]project.IPAccessList{})
			checkIPAccessListInAtlas()
		})
		It("Should Succeed (multiple)", func() {
			tenHoursLater := time.Now().Add(time.Hour * 10).Format("2006-01-02T15:04:05-0700")
			createdProject = mdbv1.DefaultProject(namespace.Name, connectionSecret.Name).
				WithIPAccessList(project.NewIPAccessList().WithComment("bla").WithCIDR("203.0.113.0/24").WithDeleteAfterDate(tenHoursLater)).
				WithIPAccessList(project.NewIPAccessList().WithComment("foo").WithIP("192.0.2.20"))

			Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

			Eventually(func(g Gomega) bool {
				return testutil.CheckCondition(k8sClient, createdProject, status.TrueCondition(status.ReadyType), validateNoErrorsIPAccessListDuringCreate(g))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			checkAtlasProjectIsReady()
			checkExpiredAccessLists([]project.IPAccessList{})
			checkIPAccessListInAtlas()
		})
		It("Should Succeed (1 expired)", func() {
			tenHoursBefore := time.Now().Add(time.Hour * -10).Format("2006-01-02T15:04:05-0700")
			expiredList := project.IPAccessList{Comment: "bla", CIDRBlock: "203.0.113.0/24", DeleteAfterDate: tenHoursBefore}
			activeList := project.IPAccessList{Comment: "foo", IPAddress: "192.0.2.20"}

			createdProject = mdbv1.DefaultProject(namespace.Name, connectionSecret.Name).WithIPAccessList(expiredList).WithIPAccessList(activeList)

			Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

			Eventually(func(g Gomega) bool {
				return testutil.CheckCondition(k8sClient, createdProject, status.TrueCondition(status.ReadyType), validateNoErrorsIPAccessListDuringCreate(g))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			checkAtlasProjectIsReady()
			checkExpiredAccessLists([]project.IPAccessList{expiredList})

			// Atlas
			list, _, err := atlasClient.ProjectIPAccessList.List(context.Background(), createdProject.ID(), &mongodbatlas.ListOptions{})
			Expect(err).NotTo(HaveOccurred())

			Expect(list.Results).To(HaveLen(1))
			Expect(list.Results[0]).To(testutil.MatchIPAccessList(activeList))
		})
		It("Should Fail (AWS security group not supported without VPC)", func() {
			createdProject = mdbv1.DefaultProject(namespace.Name, connectionSecret.Name).
				WithIPAccessList(project.NewIPAccessList().WithAWSGroup("sg-0026348ec11780bd1"))

			Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

			Eventually(func() bool {
				return testutil.CheckCondition(k8sClient, createdProject, status.FalseCondition(status.IPAccessListReadyType))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			ipAccessFailedCondition := status.FalseCondition(status.IPAccessListReadyType).
				WithReason(string(workflow.ProjectIPNotCreatedInAtlas)).
				WithMessageRegexp(".*CANNOT_USE_AWS_SECURITY_GROUP_WITHOUT_VPC_PEERING_CONNECTION.*")

			expectedConditionsMatchers := testutil.MatchConditions(
				status.TrueCondition(status.ProjectReadyType),
				status.TrueCondition(status.ValidationSucceeded),
				ipAccessFailedCondition,
				status.FalseCondition(status.ReadyType),
				status.TrueCondition(status.ResourceVersionStatus),
			)
			Expect(createdProject.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
			checkExpiredAccessLists([]project.IPAccessList{})
		})
	})

	Describe("Updating the project IP access list", func() {
		It("Should Succeed (single)", func() {
			By("Creating the project first", func() {
				tenMinutesLater := time.Now().Add(time.Minute * 10).Format("2006-01-02T15:04:05-0700")
				createdProject = mdbv1.DefaultProject(namespace.Name, connectionSecret.Name).
					WithIPAccessList(project.NewIPAccessList().WithComment("bla").WithIP("192.0.2.15").WithDeleteAfterDate(tenMinutesLater))

				Expect(k8sClient.Create(context.Background(), createdProject)).To(Succeed())

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdProject, status.TrueCondition(status.ReadyType), validateNoErrorsIPAccessListDuringCreate(g))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())
			})
			By("Updating the IP Access List comment and delete date", func() {
				// Just a note: Atlas doesn't allow to make the "permanent" entity "temporary". But it works the other way
				createdProject.Spec.ProjectIPAccessList[0].Comment = "new comment"
				createdProject.Spec.ProjectIPAccessList[0].DeleteAfterDate = ""
				Expect(k8sClient.Update(context.Background(), createdProject)).To(Succeed())

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdProject, status.TrueCondition(status.ReadyType), validateNoErrorsIPAccessListDuringUpdate(g))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

				checkAtlasProjectIsReady()
				checkExpiredAccessLists([]project.IPAccessList{})
				checkIPAccessListInAtlas()
			})
		})

		It("Should Succeed (multiple)", func() {
			By("Creating the project first", func() {
				thirtyHoursLater := time.Now().Add(time.Hour * 30).Format("2006-01-02T15:04:05-0700")
				createdProject = mdbv1.DefaultProject(namespace.Name, connectionSecret.Name).
					WithIPAccessList(project.NewIPAccessList().WithComment("bla").WithCIDR("203.0.113.0/24").WithDeleteAfterDate(thirtyHoursLater)).
					WithIPAccessList(project.NewIPAccessList().WithComment("foo").WithIP("192.0.2.20"))

				Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdProject, status.TrueCondition(status.ReadyType), validateNoErrorsIPAccessListDuringCreate(g))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())
			})
			By("Updating the IP Access List IPAddress", func() {
				twoDaysLater := time.Now().Add(time.Hour * 48).Format("2006-01-02T15:04:05Z")
				createdProject.Spec.ProjectIPAccessList[0].DeleteAfterDate = twoDaysLater
				// Update of the IP address will result in delete for the old IP address first and then the new
				// IP address will be created
				createdProject.Spec.ProjectIPAccessList[1].IPAddress = "168.32.54.0"
				Expect(k8sClient.Update(context.Background(), createdProject)).To(Succeed())

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdProject, status.TrueCondition(status.ReadyType), validateNoErrorsIPAccessListDuringUpdate(g))
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
				createdProject = mdbv1.DefaultProject(namespace.Name, connectionSecret.Name).
					WithMaintenanceWindow(project.NewMaintenanceWindow().WithDay(2).WithHour(2).WithAutoDefer(false))

				Expect(k8sClient.Create(context.Background(), createdProject)).To(Succeed())

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdProject, status.TrueCondition(status.ReadyType), validateNoErrorsMaintenanceWindowDuringCreate(g))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			})
			By("Updating the project maintenance window hour and enabling auto-defer", func() {
				createdProject.Spec.MaintenanceWindow.HourOfDay = 3
				createdProject.Spec.MaintenanceWindow.AutoDefer = true
				Expect(k8sClient.Update(context.Background(), createdProject)).To(Succeed())

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdProject, status.TrueCondition(status.ReadyType), validateNoErrorsMaintenanceWindowDuringUpdate(g))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

				checkAtlasProjectIsReady()
				checkMaintenanceWindowInAtlas()
			})
			By("Toggling auto-defer to false", func() {
				createdProject.Spec.MaintenanceWindow.AutoDefer = false
				Expect(k8sClient.Update(context.Background(), createdProject)).To(Succeed())

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdProject, status.TrueCondition(status.ReadyType), validateNoErrorsMaintenanceWindowDuringUpdate(g))
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
			createdProject = mdbv1.DefaultProject(namespace.Name, "")

			Expect(k8sClient.Create(context.Background(), createdProject)).To(Succeed())

			Eventually(func() bool {
				return testutil.CheckCondition(k8sClient, createdProject, status.TrueCondition(status.ReadyType))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

			expectedConditionsMatchers := testutil.MatchConditions(
				status.TrueCondition(status.ProjectReadyType),
				status.TrueCondition(status.ReadyType),
				status.TrueCondition(status.ValidationSucceeded),
			)
			Expect(createdProject.Status.Conditions).To(ContainElements(expectedConditionsMatchers))
			Expect(createdProject.Status.ObservedGeneration).To(Equal(createdProject.Generation))
		})
		It("Should Fail if the global Secret doesn't exist", func() {
			By("Creating without a global Secret", func() {
				createdProject = mdbv1.DefaultProject(namespace.Name, "")

				Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdProject, status.FalseCondition(status.ReadyType))
				}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())

				expectedConditionsMatchers := testutil.MatchConditions(
					status.FalseCondition(status.ProjectReadyType).WithReason(string(workflow.AtlasCredentialsNotProvided)),
					status.FalseCondition(status.ReadyType),
					status.TrueCondition(status.ValidationSucceeded),
					status.TrueCondition(status.ResourceVersionStatus),
				)
				Expect(createdProject.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
				Expect(createdProject.ID()).To(BeEmpty())
				Expect(createdProject.Status.ObservedGeneration).To(Equal(createdProject.Generation))
			})
			By("Creating a global Secret - should get fixed", func() {
				globalConnectionSecret := buildConnectionSecret("atlas-operator-api-key")
				Expect(k8sClient.Create(context.Background(), &globalConnectionSecret)).To(Succeed())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdProject, status.TrueCondition(status.ReadyType))
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
		StringData: map[string]string{"orgId": connection.OrgID, "publicApiKey": connection.PublicKey, "privateApiKey": connection.PrivateKey},
	}
}

// checkAtlasProjectRemoved returns true if the Atlas Project is removed from Atlas.
func checkAtlasProjectRemoved(projectID string) func() bool {
	return func() bool {
		_, r, err := atlasClient.Projects.GetOneProject(context.Background(), projectID)
		if err != nil {
			if r != nil && r.StatusCode == http.StatusNotFound {
				return true
			}
		}
		return false
	}
}

// validateNoErrorsIPAccessListDuringCreate performs check that no problems happen to IP Access list during the creation.
// This allows the test to fail fast instead by timeout if there are any troubles.
func validateNoErrorsIPAccessListDuringCreate(g Gomega) func(a mdbv1.AtlasCustomResource) {
	return func(a mdbv1.AtlasCustomResource) {
		c := a.(*mdbv1.AtlasProject)

		if condition, ok := testutil.FindConditionByType(c.Status.Conditions, status.IPAccessListReadyType); ok {
			g.Expect(condition.Status).To(Equal(status.TrueCondition(status.IPAccessListReadyType).Status), fmt.Sprintf("Unexpected condition: %v", condition))
		}
	}
}

// validateNoErrorsIPAccessListDuringUpdate performs check that no problems happen to IP Access list during the update.
func validateNoErrorsIPAccessListDuringUpdate(g Gomega) func(a mdbv1.AtlasCustomResource) {
	return func(a mdbv1.AtlasCustomResource) {
		c := a.(*mdbv1.AtlasProject)
		condition, ok := testutil.FindConditionByType(c.Status.Conditions, status.IPAccessListReadyType)
		g.Expect(ok).To(BeTrue())
		g.Expect(condition.Reason).To(BeEmpty())
	}
}

// validateNoErrorsMaintenanceWindowDuringCreate performs check that no problems happen to Maintenance Window during the creation.
// This allows the test to fail fast instead by timeout if there are any troubles.
func validateNoErrorsMaintenanceWindowDuringCreate(g Gomega) func(a mdbv1.AtlasCustomResource) {
	return func(a mdbv1.AtlasCustomResource) {
		c := a.(*mdbv1.AtlasProject)

		if condition, ok := testutil.FindConditionByType(c.Status.Conditions, status.MaintenanceWindowReadyType); ok {
			g.Expect(condition.Status).To(Equal(status.TrueCondition(status.MaintenanceWindowReadyType).Status), fmt.Sprintf("Unexpected condition: %v", condition))
		}
	}
}

// validateNoErrorsMaintenanceWindowDuringUpdate performs check that no problems happen to Maintenance Window during the update.
func validateNoErrorsMaintenanceWindowDuringUpdate(g Gomega) func(a mdbv1.AtlasCustomResource) {
	return func(a mdbv1.AtlasCustomResource) {
		c := a.(*mdbv1.AtlasProject)
		condition, ok := testutil.FindConditionByType(c.Status.Conditions, status.MaintenanceWindowReadyType)
		g.Expect(ok).To(BeTrue())
		g.Expect(condition.Reason).To(BeEmpty())
	}
}
