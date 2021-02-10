package int

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas/mongodbatlas"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/testutil"
)

var _ = Describe("AtlasProject", func() {
	const interval = time.Second * 1

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

	checkExpiredAccessLists := func(lists []mdbv1.ProjectIPAccessList) {
		expiredCopy := make([]status.ProjectIPAccessList, len(lists))
		for i, list := range lists {
			expiredCopy[i] = status.ProjectIPAccessList(list)
		}
		currentStatusIPs := createdProject.Status.ExpiredIPAccessList
		if currentStatusIPs == nil {
			currentStatusIPs = []status.ProjectIPAccessList{}
		}
		Expect(currentStatusIPs).To(Equal(expiredCopy))
	}

	checkAtlasProjectIsReady := func() {
		projectReadyConditions := testutil.MatchConditions(
			status.TrueCondition(status.ProjectReadyType),
			status.TrueCondition(status.IPAccessListReadyType),
			status.TrueCondition(status.ReadyType),
		)
		Expect(createdProject.Status.ID).NotTo(BeNil())
		Expect(createdProject.Status.Conditions).To(ConsistOf(projectReadyConditions))
		Expect(createdProject.Status.ObservedGeneration).To(Equal(createdProject.Generation))
	}

	Describe("Creating the project", func() {
		It("Should Succeed", func() {
			expectedProject := testAtlasProject(namespace.Name, "test-project", namespace.Name, connectionSecret.Name)
			createdProject.ObjectMeta = expectedProject.ObjectMeta
			Expect(k8sClient.Create(context.Background(), expectedProject)).ToNot(HaveOccurred())

			Eventually(testutil.WaitFor(k8sClient, createdProject, status.TrueCondition(status.ReadyType)),
				20, interval).Should(BeTrue())

			checkAtlasProjectIsReady()

			// Atlas
			atlasProject, _, err := atlasClient.Projects.GetOneProject(context.Background(), createdProject.Status.ID)
			Expect(err).ToNot(HaveOccurred())

			Expect(atlasProject.Name).To(Equal(expectedProject.Spec.Name))
		})
		It("Should fail if Secret is wrong", func() {
			expectedProject := testAtlasProject(namespace.Name, "test-project", namespace.Name, "non-existent-secret")
			createdProject.ObjectMeta = expectedProject.ObjectMeta
			Expect(k8sClient.Create(context.Background(), expectedProject)).ToNot(HaveOccurred())

			expectedCondition := status.FalseCondition(status.ProjectReadyType).WithReason(string(workflow.AtlasCredentialsNotProvided))
			Eventually(testutil.WaitFor(k8sClient, createdProject, expectedCondition),
				20, interval).Should(BeTrue())

			expectedConditionsMatchers := testutil.MatchConditions(
				status.FalseCondition(status.ProjectReadyType),
				status.FalseCondition(status.ReadyType),
			)
			Expect(createdProject.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
			Expect(createdProject.ID()).To(BeEmpty())
			Expect(createdProject.Status.ObservedGeneration).To(Equal(createdProject.Generation))

			// Atlas
			_, _, err := atlasClient.Projects.GetOneProjectByName(context.Background(), expectedProject.Spec.Name)

			// "NOT_IN_GROUP" is what is returned if the project is not found
			var apiError *mongodbatlas.ErrorResponse
			Expect(errors.As(err, &apiError)).To(BeTrue(), "Error occurred: "+err.Error())
			Expect(apiError.ErrorCode).To(Equal(atlas.NotInGroup))
		})
	})

	Describe("Updating the project", func() {
		It("Should Succeed", func() {
			By("Creating the project first")

			expectedProject := testAtlasProject(namespace.Name, "test-project", namespace.Name, connectionSecret.Name)
			createdProject.ObjectMeta = expectedProject.ObjectMeta
			Expect(k8sClient.Create(context.Background(), expectedProject)).ToNot(HaveOccurred())

			Eventually(testutil.WaitFor(k8sClient, createdProject, status.TrueCondition(status.ReadyType)),
				20, interval).Should(BeTrue())

			// Updating (the existing project is expected to be read from Atlas)
			By("Updating the project")

			createdProject.Spec.ProjectIPAccessList = []mdbv1.ProjectIPAccessList{{CIDRBlock: "0.0.0.0/0"}}
			Expect(k8sClient.Update(context.Background(), createdProject)).To(Succeed())

			Eventually(testutil.WaitFor(k8sClient, createdProject, status.TrueCondition(status.ReadyType)),
				20, interval).Should(BeTrue())

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
				createdProject = testAtlasProject(namespace.Name, "test-project", namespace.Name, connectionSecret.Name)
				Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

				secondProject = testAtlasProject(namespace.Name, "second-project", "second Project", connectionSecret.Name)
				Expect(k8sClient.Create(context.Background(), secondProject)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdProject, status.TrueCondition(status.ReadyType)),
					20, interval).Should(BeTrue())

				Eventually(testutil.WaitFor(k8sClient, secondProject, status.TrueCondition(status.ReadyType)),
					20, interval).Should(BeTrue())
			})
			By("Breaking the Connection Secret", func() {
				connectionSecret.StringData["publicApiKey"] = "non-existing"
				Expect(k8sClient.Update(context.Background(), &connectionSecret)).To(Succeed())

				// Both projects are expected to get to Failed state right away
				expectedCondition := status.FalseCondition(status.ProjectReadyType).WithReason(string(workflow.ProjectNotCreatedInAtlas))
				Eventually(testutil.WaitFor(k8sClient, createdProject, expectedCondition),
					20, interval).Should(BeTrue())
				Eventually(testutil.WaitFor(k8sClient, secondProject, expectedCondition),
					20, interval).Should(BeTrue())
			})
			By("Fixing the Connection Secret", func() {
				connectionSecret.StringData["publicApiKey"] = connection.PublicKey
				Expect(k8sClient.Update(context.Background(), &connectionSecret)).To(Succeed())

				// Both projects are expected to recover
				Eventually(testutil.WaitFor(k8sClient, createdProject, status.TrueCondition(status.ReadyType)),
					20, interval).Should(BeTrue())
				Eventually(testutil.WaitFor(k8sClient, secondProject, status.TrueCondition(status.ReadyType)),
					20, interval).Should(BeTrue())
			})
		})
	})

	Describe("Creating the project IP access list", func() {
		It("Should Succeed (single)", func() {
			createdProject = testAtlasProject(namespace.Name, "test-project", namespace.Name, connectionSecret.Name)
			createdProject.Spec.ProjectIPAccessList = []mdbv1.ProjectIPAccessList{{Comment: "bla", IPAddress: "192.0.2.15"}}
			Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

			Eventually(testutil.WaitFor(k8sClient, createdProject, status.TrueCondition(status.ReadyType), validateNoErrorsIPAccessListDuringCreate),
				20, interval).Should(BeTrue())

			checkAtlasProjectIsReady()
			checkExpiredAccessLists([]mdbv1.ProjectIPAccessList{})
			checkIPAccessListInAtlas()
		})
		It("Should Succeed (multiple)", func() {
			createdProject = testAtlasProject(namespace.Name, "test-project", namespace.Name, connectionSecret.Name)
			tenHoursLater := time.Now().Add(time.Hour * 10).Format("2006-01-02T15:04:05-0700")

			createdProject.Spec.ProjectIPAccessList = []mdbv1.ProjectIPAccessList{
				{Comment: "bla", CIDRBlock: "203.0.113.0/24", DeleteAfterDate: tenHoursLater},
				{Comment: "foo", IPAddress: "192.0.2.20"},
			}
			Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

			Eventually(testutil.WaitFor(k8sClient, createdProject, status.TrueCondition(status.ReadyType), validateNoErrorsIPAccessListDuringCreate),
				20, interval).Should(BeTrue())

			checkAtlasProjectIsReady()
			checkExpiredAccessLists([]mdbv1.ProjectIPAccessList{})
			checkIPAccessListInAtlas()
		})
		It("Should Succeed (1 expired)", func() {
			createdProject = testAtlasProject(namespace.Name, "test-project", namespace.Name, connectionSecret.Name)
			tenHoursBefore := time.Now().Add(time.Hour * -10).Format("2006-01-02T15:04:05-0700")

			expiredList := mdbv1.ProjectIPAccessList{Comment: "bla", CIDRBlock: "203.0.113.0/24", DeleteAfterDate: tenHoursBefore}
			activeList := mdbv1.ProjectIPAccessList{Comment: "foo", IPAddress: "192.0.2.20"}
			createdProject.Spec.ProjectIPAccessList = []mdbv1.ProjectIPAccessList{expiredList, activeList}

			Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

			Eventually(testutil.WaitFor(k8sClient, createdProject, status.TrueCondition(status.ReadyType), validateNoErrorsIPAccessListDuringCreate),
				20, interval).Should(BeTrue())

			checkAtlasProjectIsReady()
			checkExpiredAccessLists([]mdbv1.ProjectIPAccessList{expiredList})

			// Atlas
			list, _, err := atlasClient.ProjectIPAccessList.List(context.Background(), createdProject.ID(), &mongodbatlas.ListOptions{})
			Expect(err).NotTo(HaveOccurred())

			Expect(list.Results).To(HaveLen(1))
			Expect(list.Results[0]).To(testutil.MatchIPAccessList(activeList))
		})
		It("Should Fail (AWS security group not supported without VPC)", func() {
			createdProject = testAtlasProject(namespace.Name, "test-project", namespace.Name, connectionSecret.Name)
			createdProject.Spec.ProjectIPAccessList = []mdbv1.ProjectIPAccessList{{AwsSecurityGroup: "sg-0026348ec11780bd1"}}
			Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

			Eventually(testutil.WaitFor(k8sClient, createdProject, status.FalseCondition(status.IPAccessListReadyType)),
				20, interval).Should(BeTrue())

			ipAccessFailedCondition := status.FalseCondition(status.IPAccessListReadyType).
				WithReason(string(workflow.ProjectIPNotCreatedInAtlas)).
				WithMessageRegexp(".*CANNOT_USE_AWS_SECURITY_GROUP_WITHOUT_VPC_PEERING_CONNECTION.*")

			expectedConditionsMatchers := testutil.MatchConditions(
				status.TrueCondition(status.ProjectReadyType),
				ipAccessFailedCondition,
				status.FalseCondition(status.ReadyType),
			)
			Expect(createdProject.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
			checkExpiredAccessLists([]mdbv1.ProjectIPAccessList{})
		})
	})

	Describe("Updating the project IP access list", func() {
		It("Should Succeed (single)", func() {
			By("Creating the project first", func() {
				tenMinutesLater := time.Now().Add(time.Minute * 10).Format("2006-01-02T15:04:05-0700")
				createdProject = testAtlasProject(namespace.Name, "test-project", namespace.Name, connectionSecret.Name)
				createdProject.Spec.ProjectIPAccessList = []mdbv1.ProjectIPAccessList{{Comment: "bla", IPAddress: "192.0.2.15", DeleteAfterDate: tenMinutesLater}}
				Expect(k8sClient.Create(context.Background(), createdProject)).To(Succeed())

				Eventually(testutil.WaitFor(k8sClient, createdProject, status.TrueCondition(status.ReadyType), validateNoErrorsIPAccessListDuringCreate),
					20, interval).Should(BeTrue())
			})
			By("Updating the IP Access List comment and delete date", func() {
				// Just a note: Atlas doesn't allow to make the "permanent" entity "temporary". But it works the other way
				createdProject.Spec.ProjectIPAccessList[0].Comment = "new comment"
				createdProject.Spec.ProjectIPAccessList[0].DeleteAfterDate = ""
				Expect(k8sClient.Update(context.Background(), createdProject)).To(Succeed())

				Eventually(testutil.WaitFor(k8sClient, createdProject, status.TrueCondition(status.ReadyType), validateNoErrorsIPAccessListDuringUpdate),
					20, interval).Should(BeTrue())

				checkAtlasProjectIsReady()
				checkExpiredAccessLists([]mdbv1.ProjectIPAccessList{})
				checkIPAccessListInAtlas()
			})
		})

		It("Should Succeed (multiple)", func() {
			By("Creating the project first", func() {
				createdProject = testAtlasProject(namespace.Name, "test-project", namespace.Name, connectionSecret.Name)
				thirtyHoursLater := time.Now().Add(time.Hour * 30).Format("2006-01-02T15:04:05-0700")

				createdProject.Spec.ProjectIPAccessList = []mdbv1.ProjectIPAccessList{
					{Comment: "bla", CIDRBlock: "203.0.113.0/24", DeleteAfterDate: thirtyHoursLater},
					{Comment: "foo", IPAddress: "192.0.2.20"},
				}
				Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdProject, status.TrueCondition(status.ReadyType), validateNoErrorsIPAccessListDuringCreate),
					20, interval).Should(BeTrue())
			})
			By("Updating the IP Access List IPAddress", func() {
				twoDaysLater := time.Now().Add(time.Hour * 48).Format("2006-01-02T15:04:05Z")
				createdProject.Spec.ProjectIPAccessList[0].DeleteAfterDate = twoDaysLater
				// Update of the IP address will result in delete for the old IP address first and then the new
				// IP address will be created
				createdProject.Spec.ProjectIPAccessList[1].IPAddress = "168.32.54.0"
				Expect(k8sClient.Update(context.Background(), createdProject)).To(Succeed())

				Eventually(testutil.WaitFor(k8sClient, createdProject, status.TrueCondition(status.ReadyType), validateNoErrorsIPAccessListDuringUpdate),
					20, interval).Should(BeTrue())

				checkAtlasProjectIsReady()
				checkExpiredAccessLists([]mdbv1.ProjectIPAccessList{})
				checkIPAccessListInAtlas()
			})
		})
	})

	Describe("Using the global Connection Secret", func() {
		It("Should Succeed", func() {
			globalConnectionSecret := buildConnectionSecret("atlas-operator-api-key")
			Expect(k8sClient.Create(context.Background(), &globalConnectionSecret)).To(Succeed())

			// We don't specify the connection Secret per project - the global one must be used
			createdProject = testAtlasProject(namespace.Name, "test-project", namespace.Name, "")

			Expect(k8sClient.Create(context.Background(), createdProject)).To(Succeed())

			Eventually(testutil.WaitFor(k8sClient, createdProject, status.TrueCondition(status.ReadyType)),
				20, interval).Should(BeTrue())

			expectedConditionsMatchers := testutil.MatchConditions(
				status.TrueCondition(status.ProjectReadyType),
				status.TrueCondition(status.IPAccessListReadyType),
				status.TrueCondition(status.ReadyType),
			)
			Expect(createdProject.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
			Expect(createdProject.Status.ObservedGeneration).To(Equal(createdProject.Generation))
		})
		It("Should Fail if the global Secret doesn't exist", func() {
			By("Creating without a global Secret", func() {
				createdProject = testAtlasProject(namespace.Name, "test-project", namespace.Name, "")

				Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdProject, status.FalseCondition(status.ReadyType)),
					20, interval).Should(BeTrue())

				expectedConditionsMatchers := testutil.MatchConditions(
					status.FalseCondition(status.ProjectReadyType).WithReason(string(workflow.AtlasCredentialsNotProvided)),
					status.FalseCondition(status.ReadyType),
				)
				Expect(createdProject.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
				Expect(createdProject.ID()).To(BeEmpty())
				Expect(createdProject.Status.ObservedGeneration).To(Equal(createdProject.Generation))
			})
			By("Creating a global Secret - should get fixed", func() {
				globalConnectionSecret := buildConnectionSecret("atlas-operator-api-key")
				Expect(k8sClient.Create(context.Background(), &globalConnectionSecret)).To(Succeed())

				Eventually(testutil.WaitFor(k8sClient, createdProject, status.TrueCondition(status.ReadyType)),
					20, interval).Should(BeTrue())
			})

		})
	})

})

func buildConnectionSecret(name string) corev1.Secret {
	return corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace.Name,
		},
		StringData: map[string]string{"orgId": connection.OrgID, "publicApiKey": connection.PublicKey, "privateApiKey": connection.PrivateKey},
	}
}

// TODO builders
func testAtlasProject(namespace, name, atlasName, connectionSecretName string) *mdbv1.AtlasProject {
	project := mdbv1.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: mdbv1.AtlasProjectSpec{
			Name: atlasName,
		},
	}
	if connectionSecretName != "" {
		project.Spec.ConnectionSecret = &mdbv1.ResourceRef{Name: connectionSecretName}
	}
	return &project
}

func removeAtlasProject(projectID string) func() bool {
	return func() bool {
		_, err := atlasClient.Projects.Delete(context.Background(), projectID)
		if err != nil {
			var apiError *mongodbatlas.ErrorResponse
			Expect(errors.As(err, &apiError)).To(BeTrue())
			Expect(apiError.ErrorCode).To(Equal(atlas.CannotCloseGroupActiveAtlasCluster))
			return false
		}
		return true
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

// validateNoErrorsIPAccessListDuringCreate performs check that no problems happen to IP Access list during the create.
// This allows the test to fail fast instead by timeout if there are any troubles.
func validateNoErrorsIPAccessListDuringCreate(a mdbv1.AtlasCustomResource) {
	c := a.(*mdbv1.AtlasProject)
	condition, ok := testutil.FindConditionByType(c.Status.Conditions, status.IPAccessListReadyType)
	Expect(ok).To(BeFalse(), fmt.Sprintf("Unexpected condition: %v", condition))
}

// validateNoErrorsIPAccessListDuringUpdate performs check that no problems happen to IP Access list during the update.
func validateNoErrorsIPAccessListDuringUpdate(a mdbv1.AtlasCustomResource) {
	c := a.(*mdbv1.AtlasProject)
	condition, ok := testutil.FindConditionByType(c.Status.Conditions, status.IPAccessListReadyType)
	Expect(ok).To(BeTrue())
	Expect(condition.Reason).To(BeEmpty())
}
