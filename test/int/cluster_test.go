package int

import (
	"context"
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
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/testutil"
)

var _ = Describe("AtlasCluster", func() {
	const interval = time.Second * 1

	var (
		connectionSecret corev1.Secret
		createdProject   *mdbv1.AtlasProject
		createdCluster   *mdbv1.AtlasCluster
		lastGeneration   int64
	)

	BeforeEach(func() {
		prepareControllers()

		createdCluster = &mdbv1.AtlasCluster{}

		lastGeneration = 0

		connectionSecret = corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-atlas-key",
				Namespace: namespace.Name,
			},
			StringData: map[string]string{"orgId": connection.OrgID, "publicApiKey": connection.PublicKey, "privateApiKey": connection.PrivateKey},
		}
		By(fmt.Sprintf("Creating the Secret %s", kube.ObjectKeyFromObject(&connectionSecret)))
		Expect(k8sClient.Create(context.Background(), &connectionSecret)).To(Succeed())

		createdProject = mdbv1.DefaultProject(namespace.Name, connectionSecret.Name)
		By("Creating the project " + createdProject.Name)
		Expect(k8sClient.Create(context.Background(), createdProject)).To(Succeed())
		Eventually(testutil.WaitFor(k8sClient, createdProject, status.TrueCondition(status.ReadyType)),
			20, interval).Should(BeTrue())
	})

	AfterEach(func() {
		if createdProject != nil && createdProject.Status.ID != "" {
			if createdCluster != nil {
				By("Removing Atlas Cluster " + createdCluster.Name)
				Expect(k8sClient.Delete(context.Background(), createdCluster)).To(Succeed())
				Eventually(checkAtlasClusterRemoved(createdProject.Status.ID, createdCluster.Name), 600, interval).Should(BeTrue())
			}

			By("Removing Atlas Project " + createdProject.Status.ID)
			Expect(k8sClient.Delete(context.Background(), createdProject)).To(Succeed())
			Eventually(checkAtlasProjectRemoved(createdProject.Status.ID), 600, interval).Should(BeTrue())
		}
		removeControllersAndNamespace()
	})

	doCommonChecks := func() {
		By("Checking observed Cluster state", func() {
			Expect(createdCluster.Status.ConnectionStrings).NotTo(BeNil())
			Expect(createdCluster.Status.ConnectionStrings.Standard).NotTo(BeNil())
			Expect(createdCluster.Status.ConnectionStrings.StandardSrv).NotTo(BeNil())
			Expect(createdCluster.Status.MongoDBVersion).NotTo(BeNil())
			Expect(createdCluster.Status.MongoURIUpdated).NotTo(BeNil())
			Expect(createdCluster.Status.StateName).To(Equal("IDLE"))
			Expect(createdCluster.Status.Conditions).To(ConsistOf(testutil.MatchConditions(
				status.TrueCondition(status.ClusterReadyType),
				status.TrueCondition(status.ReadyType),
			)))
			Expect(createdCluster.Status.ObservedGeneration).To(Equal(createdCluster.Generation))
			Expect(createdCluster.Status.ObservedGeneration).To(Equal(lastGeneration + 1))
		})
	}

	checkAtlasState := func(additionalChecks ...func(c *mongodbatlas.Cluster)) {
		By("Verifying Cluster state in Atlas", func() {
			atlasCluster, _, err := atlasClient.Clusters.Get(context.Background(), createdProject.Status.ID, createdCluster.Spec.Name)
			Expect(err).ToNot(HaveOccurred())

			createdAtlasCluster, err := createdCluster.Spec.Cluster()
			Expect(err).ToNot(HaveOccurred())

			Expect(atlasCluster.Name).To(Equal(createdAtlasCluster.Name))
			Expect(atlasCluster.Labels).To(ConsistOf(createdAtlasCluster.Labels))
			Expect(atlasCluster.ProviderSettings.InstanceSizeName).To(Equal(createdAtlasCluster.ProviderSettings.InstanceSizeName))
			Expect(atlasCluster.ProviderSettings.ProviderName).To(Equal(createdAtlasCluster.ProviderSettings.ProviderName))
			Expect(atlasCluster.ProviderSettings.RegionName).To(Equal(createdAtlasCluster.ProviderSettings.RegionName))

			for _, check := range additionalChecks {
				check(atlasCluster)
			}
		})
	}

	performUpdate := func() {
		Expect(k8sClient.Update(context.Background(), createdCluster)).To(Succeed())

		Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validateClusterUpdatingFunc()),
			1200, interval).Should(BeTrue())

		lastGeneration++
	}

	Describe("Create/Update the cluster", func() {
		It("Should fail, then be fixed", func() {
			createdCluster = mdbv1.DefaultGCPCluster(namespace.Name, createdProject.Name).WithAtlasName("")

			By(fmt.Sprintf("Creating the Cluster %s with invalid parameters", kube.ObjectKeyFromObject(createdCluster)), func() {
				Expect(k8sClient.Create(context.Background(), createdCluster)).ToNot(HaveOccurred())

				Eventually(
					testutil.WaitFor(
						k8sClient,
						createdCluster,
						status.
							FalseCondition(status.ClusterReadyType).
							WithReason(string(workflow.Internal)). // Internal due to reconciliation failing on the initial GET request
							WithMessageRegexp("name is invalid because must be set"),
					),
					60,
					interval,
				).Should(BeTrue())

				lastGeneration++
			})

			By("Fixing the cluster", func() {
				createdCluster.Spec.Name = "fixed-cluster"

				Expect(k8sClient.Update(context.Background(), createdCluster)).To(Succeed())

				Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType)),
					1200, interval).Should(BeTrue())

				doCommonChecks()
				checkAtlasState()
			})
		})

		It("Should Succeed", func() {
			createdCluster = mdbv1.DefaultGCPCluster(namespace.Name, createdProject.Name)

			By(fmt.Sprintf("Creating the Cluster %s", kube.ObjectKeyFromObject(createdCluster)), func() {
				Expect(k8sClient.Create(context.Background(), createdCluster)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
					1800, interval).Should(BeTrue())

				doCommonChecks()
				checkAtlasState()
			})

			By("Updating the Cluster labels", func() {
				createdCluster.Spec.Labels = []mdbv1.LabelSpec{{Key: "int-test", Value: "true"}}
				performUpdate()
				doCommonChecks()
				checkAtlasState()
			})

			By("Updating the Cluster backups settings", func() {
				createdCluster.Spec.ProviderBackupEnabled = boolptr(true)
				performUpdate()
				doCommonChecks()
				checkAtlasState(func(c *mongodbatlas.Cluster) {
					Expect(c.ProviderBackupEnabled).To(Equal(createdCluster.Spec.ProviderBackupEnabled))
				})
			})

			By("Decreasing the Cluster disk size", func() {
				createdCluster.Spec.DiskSizeGB = intptr(10)
				performUpdate()
				doCommonChecks()
				checkAtlasState(func(c *mongodbatlas.Cluster) {
					Expect(*c.DiskSizeGB).To(BeEquivalentTo(*createdCluster.Spec.DiskSizeGB))

					// check whether https://github.com/mongodb/go-client-mongodb-atlas/issues/140 is fixed
					Expect(c.DiskSizeGB).To(BeAssignableToTypeOf(float64ptr(0)), "DiskSizeGB is no longer a *float64, please check the spec!")
				})
			})

			By("Pausing the cluster", func() {
				createdCluster.Spec.Paused = boolptr(true)
				performUpdate()
				doCommonChecks()
				checkAtlasState(func(c *mongodbatlas.Cluster) {
					Expect(c.Paused).To(Equal(createdCluster.Spec.Paused))
				})
			})

			By("Updating the Cluster configuration while paused (should fail)", func() {
				createdCluster.Spec.ProviderBackupEnabled = boolptr(false)

				Expect(k8sClient.Update(context.Background(), createdCluster)).To(Succeed())
				Eventually(
					testutil.WaitFor(
						k8sClient,
						createdCluster,
						status.
							FalseCondition(status.ClusterReadyType).
							WithReason(string(workflow.ClusterNotUpdatedInAtlas)).
							WithMessageRegexp("CANNOT_UPDATE_PAUSED_CLUSTER"),
					),
					60,
					interval,
				).Should(BeTrue())

				lastGeneration++
			})

			By("Unpausing the cluster", func() {
				createdCluster.Spec.Paused = boolptr(false)
				performUpdate()
				doCommonChecks()
				checkAtlasState(func(c *mongodbatlas.Cluster) {
					Expect(c.Paused).To(Equal(createdCluster.Spec.Paused))
				})
			})

			By("Checking that modifications were applied after unpausing", func() {
				doCommonChecks()
				checkAtlasState(func(c *mongodbatlas.Cluster) {
					Expect(c.ProviderBackupEnabled).To(Equal(createdCluster.Spec.ProviderBackupEnabled))
				})
			})

			By("Setting AutoScaling.Compute.Enabled to false (should fail)", func() {
				createdCluster.Spec.ProviderSettings.AutoScaling = &mdbv1.AutoScalingSpec{
					Compute: &mdbv1.ComputeSpec{
						Enabled: boolptr(false),
					},
				}

				Expect(k8sClient.Update(context.Background(), createdCluster)).To(Succeed())
				Eventually(
					testutil.WaitFor(
						k8sClient,
						createdCluster,
						status.
							FalseCondition(status.ClusterReadyType).
							WithReason(string(workflow.ClusterNotUpdatedInAtlas)).
							WithMessageRegexp("INVALID_ATTRIBUTE"),
					),
					60,
					interval,
				).Should(BeTrue())

				lastGeneration++

				By("Fixing the Cluster", func() {
					createdCluster.Spec.ProviderSettings.AutoScaling = nil
					performUpdate()
					doCommonChecks()
					checkAtlasState()
				})
			})

			By("Setting incorrect instance size (should fail)", func() {
				oldSizeName := createdCluster.Spec.ProviderSettings.InstanceSizeName
				createdCluster.Spec.ProviderSettings.InstanceSizeName = "M42"

				Expect(k8sClient.Update(context.Background(), createdCluster)).To(Succeed())
				Eventually(
					testutil.WaitFor(
						k8sClient,
						createdCluster,
						status.
							FalseCondition(status.ClusterReadyType).
							WithReason(string(workflow.ClusterNotUpdatedInAtlas)).
							WithMessageRegexp("INVALID_ENUM_VALUE"),
					),
					60,
					interval,
				).Should(BeTrue())

				lastGeneration++

				By("Fixing the Cluster", func() {
					createdCluster.Spec.ProviderSettings.InstanceSizeName = oldSizeName
					performUpdate()
					doCommonChecks()
					checkAtlasState()
				})
			})
		})
	})
})

func validateClusterCreatingFunc() func(a mdbv1.AtlasCustomResource) {
	startedCreation := false
	return func(a mdbv1.AtlasCustomResource) {
		c := a.(*mdbv1.AtlasCluster)
		if c.Status.StateName != "" {
			startedCreation = true
		}
		// When the create request has been made to Atlas - we expect the following status
		if startedCreation {
			Expect(c.Status.StateName).To(Equal("CREATING"), fmt.Sprintf("Current conditions: %+v", c.Status.Conditions))
			expectedConditionsMatchers := testutil.MatchConditions(
				status.FalseCondition(status.ClusterReadyType).WithReason(string(workflow.ClusterCreating)).WithMessageRegexp("cluster is provisioning"),
				status.FalseCondition(status.ReadyType),
			)
			Expect(c.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
		} else {
			// Otherwise there could have been some exception in Atlas on creation - let's check the conditions
			condition, ok := testutil.FindConditionByType(c.Status.Conditions, status.ClusterReadyType)
			Expect(ok).To(BeFalse(), fmt.Sprintf("Unexpected condition: %v", condition))
		}
	}
}

func validateClusterUpdatingFunc() func(a mdbv1.AtlasCustomResource) {
	isIdle := true
	return func(a mdbv1.AtlasCustomResource) {
		c := a.(*mdbv1.AtlasCluster)
		// It's ok if the first invocations see IDLE
		if c.Status.StateName != "IDLE" {
			isIdle = false
		}
		// When the create request has been made to Atlas - we expect the following status
		if !isIdle {
			Expect(c.Status.StateName).To(Or(Equal("UPDATING"), Equal("REPAIRING")), fmt.Sprintf("Current conditions: %+v", c.Status.Conditions))
			expectedConditionsMatchers := testutil.MatchConditions(
				status.FalseCondition(status.ClusterReadyType).WithReason(string(workflow.ClusterUpdating)).WithMessageRegexp("cluster is updating"),
				status.FalseCondition(status.ReadyType),
			)
			Expect(c.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
		}
	}
}

// checkAtlasClusterRemoved returns true if the Atlas Cluster is removed from Atlas. Note the behavior: the cluster
// is removed from Atlas as soon as the DELETE API call has been made. This is different from the case when the
// cluster is terminated from UI (in this case GET request succeeds while the cluster is being terminated)
func checkAtlasClusterRemoved(projectID string, clusterName string) func() bool {
	return func() bool {
		_, r, err := atlasClient.Clusters.Get(context.Background(), projectID, clusterName)
		if err != nil {
			if r != nil && r.StatusCode == http.StatusNotFound {
				fmt.Println("cluster removed!", time.Now(), projectID, clusterName)
				return true
			}
		}

		fmt.Println("cluster exists", time.Now())
		return false
	}
}

func intptr(i int) *int {
	return &i
}

func float64ptr(f float64) *float64 {
	return &f
}

func boolptr(b bool) *bool {
	return &b
}
