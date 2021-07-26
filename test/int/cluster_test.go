package int

import (
	"context"
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlascluster"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/testutil"
)

const (
	// Set this to true if you are debugging cluster creation.
	// This may not help much if there was the update though...
	ClusterDevMode       = false
	ClusterUpdateTimeout = 40 * time.Minute
)

var _ = Describe("AtlasCluster", func() {
	const (
		interval      = PollingInterval
		intervalShort = time.Second * 2
	)

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

		createdProject = mdbv1.DefaultProject(namespace.Name, connectionSecret.Name).WithIPAccessList(project.NewIPAccessList().WithIP("0.0.0.0/0"))
		if ClusterDevMode {
			// While developing tests we need to reuse the same project
			createdProject.Spec.Name = "dev-test atlas-project"
		}
		By("Creating the project " + createdProject.Name)
		Expect(k8sClient.Create(context.Background(), createdProject)).To(Succeed())
		Eventually(testutil.WaitFor(k8sClient, createdProject, status.TrueCondition(status.ReadyType)),
			ProjectCreationTimeout, intervalShort).Should(BeTrue())
	})

	AfterEach(func() {
		if ClusterDevMode {
			return
		}
		if createdProject != nil && createdProject.Status.ID != "" {
			if createdCluster != nil {
				By("Removing Atlas Cluster " + createdCluster.Name)
				Expect(k8sClient.Delete(context.Background(), createdCluster)).To(Succeed())
				Eventually(checkAtlasClusterRemoved(createdProject.Status.ID, createdCluster.Spec.Name), 600, interval).Should(BeTrue())
			}

			By("Removing Atlas Project " + createdProject.Status.ID)
			Expect(k8sClient.Delete(context.Background(), createdProject)).To(Succeed())
			Eventually(checkAtlasProjectRemoved(createdProject.Status.ID), 60, interval).Should(BeTrue())
		}
		removeControllersAndNamespace()
	})

	doCommonStatusChecks := func() {
		By("Checking observed Cluster state", func() {
			atlasCluster, _, err := atlasClient.Clusters.Get(context.Background(), createdProject.Status.ID, createdCluster.Spec.Name)
			Expect(err).ToNot(HaveOccurred())

			Expect(createdCluster.Status.ConnectionStrings).NotTo(BeNil())
			Expect(createdCluster.Status.ConnectionStrings.Standard).To(Equal(atlasCluster.ConnectionStrings.Standard))
			Expect(createdCluster.Status.ConnectionStrings.StandardSrv).To(Equal(atlasCluster.ConnectionStrings.StandardSrv))
			Expect(createdCluster.Status.MongoDBVersion).To(Equal(atlasCluster.MongoDBVersion))
			Expect(createdCluster.Status.MongoURIUpdated).To(Equal(atlasCluster.MongoURIUpdated))
			Expect(createdCluster.Status.StateName).To(Equal("IDLE"))
			Expect(createdCluster.Status.Conditions).To(HaveLen(2))
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

			mergedCluster, err := atlascluster.MergedCluster(*atlasCluster, createdCluster.Spec)
			Expect(err).ToNot(HaveOccurred())

			Expect(atlascluster.ClustersEqual(zap.S(), *atlasCluster, mergedCluster)).To(BeTrue())

			for _, check := range additionalChecks {
				check(atlasCluster)
			}
		})
	}

	performUpdate := func(timeout interface{}) {
		Expect(k8sClient.Update(context.Background(), createdCluster)).To(Succeed())

		Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validateClusterUpdatingFunc()),
			timeout, interval).Should(BeTrue())

		lastGeneration++
	}

	Describe("Create cluster & change ReplicationSpecs", func() {
		It("Should Succeed", func() {
			createdCluster = mdbv1.DefaultAWSCluster(namespace.Name, createdProject.Name)

			// Atlas will add some defaults in case the Atlas Operator doesn't set them
			replicationSpecsCheck := func(cluster *mongodbatlas.Cluster) {
				Expect(cluster.ReplicationSpecs).To(HaveLen(1))
				Expect(cluster.ReplicationSpecs[0].ID).NotTo(BeNil())
				Expect(cluster.ReplicationSpecs[0].ZoneName).To(Equal("Zone 1"))
				Expect(cluster.ReplicationSpecs[0].RegionsConfig).To(HaveLen(1))
				Expect(cluster.ReplicationSpecs[0].RegionsConfig[createdCluster.Spec.ProviderSettings.RegionName]).NotTo(BeNil())
			}

			By(fmt.Sprintf("Creating the Cluster %s", kube.ObjectKeyFromObject(createdCluster)), func() {
				Expect(k8sClient.Create(context.Background(), createdCluster)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
					30*time.Minute, interval).Should(BeTrue())

				doCommonStatusChecks()

				singleNumShard := func(cluster *mongodbatlas.Cluster) {
					Expect(cluster.ReplicationSpecs[0].NumShards).To(Equal(int64ptr(1)))
				}
				checkAtlasState(replicationSpecsCheck, singleNumShard)
			})

			By("Updating ReplicationSpecs", func() {
				createdCluster.Spec.ReplicationSpecs = append(createdCluster.Spec.ReplicationSpecs, mdbv1.ReplicationSpec{
					NumShards: int64ptr(2),
				})
				createdCluster.Spec.ClusterType = "SHARDED"

				performUpdate(40 * time.Minute)
				doCommonStatusChecks()

				twoNumShard := func(cluster *mongodbatlas.Cluster) {
					Expect(cluster.ReplicationSpecs[0].NumShards).To(Equal(int64ptr(2)))
				}
				// ReplicationSpecs has the same defaults but the number of shards has changed
				checkAtlasState(replicationSpecsCheck, twoNumShard)
			})
		})
	})

	Describe("Create cluster & increase DiskSizeGB", func() {
		It("Should Succeed", func() {
			expectedCluster := mdbv1.DefaultAWSCluster(namespace.Name, createdProject.Name)

			By(fmt.Sprintf("Creating the Cluster %s", kube.ObjectKeyFromObject(expectedCluster)), func() {
				createdCluster.ObjectMeta = expectedCluster.ObjectMeta
				Expect(k8sClient.Create(context.Background(), expectedCluster)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
					1800, interval).Should(BeTrue())

				doCommonStatusChecks()
				checkAtlasState()
			})

			By("Increasing InstanceSize", func() {
				createdCluster.Spec.ProviderSettings.InstanceSizeName = "M30"
				performUpdate(40 * time.Minute)
				doCommonStatusChecks()
				checkAtlasState()
			})
		})
	})

	Describe("Create cluster & change it to GEOSHARDED", func() {
		It("Should Succeed", func() {
			expectedCluster := mdbv1.DefaultAWSCluster(namespace.Name, createdProject.Name)

			By(fmt.Sprintf("Creating the Cluster %s", kube.ObjectKeyFromObject(expectedCluster)), func() {
				createdCluster.ObjectMeta = expectedCluster.ObjectMeta
				Expect(k8sClient.Create(context.Background(), expectedCluster)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
					ClusterUpdateTimeout, interval).Should(BeTrue())

				doCommonStatusChecks()
				checkAtlasState()
			})

			By("Change cluster to GEOSHARDED", func() {
				createdCluster.Spec.ClusterType = "GEOSHARDED"
				createdCluster.Spec.ReplicationSpecs = []mdbv1.ReplicationSpec{
					{
						NumShards: int64ptr(1),
						ZoneName:  "Zone 1",
						RegionsConfig: map[string]mdbv1.RegionsConfig{
							"US_EAST_1": {
								AnalyticsNodes: int64ptr(1),
								ElectableNodes: int64ptr(2),
								Priority:       int64ptr(7),
								ReadOnlyNodes:  int64ptr(0),
							},
							"US_WEST_1": {
								AnalyticsNodes: int64ptr(0),
								ElectableNodes: int64ptr(1),
								Priority:       int64ptr(6),
								ReadOnlyNodes:  int64ptr(0),
							},
						},
					},
				}
				performUpdate(80 * time.Minute)
				doCommonStatusChecks()
				checkAtlasState()
			})
		})
	})

	Describe("Create/Update the cluster (more complex scenario)", func() {
		It("Should be created", func() {
			createdCluster = mdbv1.DefaultAWSCluster(namespace.Name, createdProject.Name)
			createdCluster.Spec.ClusterType = mdbv1.TypeReplicaSet
			createdCluster.Spec.AutoScaling = &mdbv1.AutoScalingSpec{
				Compute: &mdbv1.ComputeSpec{
					Enabled:          boolptr(true),
					ScaleDownEnabled: boolptr(true),
				},
			}
			createdCluster.Spec.ProviderSettings.AutoScaling = &mdbv1.AutoScalingSpec{
				Compute: &mdbv1.ComputeSpec{
					MaxInstanceSize: "M20",
					MinInstanceSize: "M10",
				},
			}
			createdCluster.Spec.ProviderSettings.InstanceSizeName = "M10"
			createdCluster.Spec.Labels = []mdbv1.LabelSpec{{Key: "createdBy", Value: "Atlas Operator"}}
			createdCluster.Spec.ReplicationSpecs = []mdbv1.ReplicationSpec{{
				NumShards: int64ptr(1),
				ZoneName:  "Zone 1",
				// One interesting thing: if the regionsConfig is not empty - Atlas nullifies the 'providerSettings.regionName' field
				RegionsConfig: map[string]mdbv1.RegionsConfig{
					"US_EAST_1": {AnalyticsNodes: int64ptr(0), ElectableNodes: int64ptr(1), Priority: int64ptr(6), ReadOnlyNodes: int64ptr(0)},
					"US_WEST_2": {AnalyticsNodes: int64ptr(0), ElectableNodes: int64ptr(2), Priority: int64ptr(7), ReadOnlyNodes: int64ptr(0)},
				},
			}}

			replicationSpecsCheckFunc := func(c *mongodbatlas.Cluster) {
				cluster, err := createdCluster.Spec.Cluster()
				Expect(err).NotTo(HaveOccurred())
				expectedReplicationSpecs := cluster.ReplicationSpecs

				// The ID field is added by Atlas - we don't have it in our specs
				Expect(c.ReplicationSpecs[0].ID).NotTo(BeNil())
				c.ReplicationSpecs[0].ID = ""
				// Apart from 'ID' all other fields are equal to the ones sent by the Operator
				Expect(c.ReplicationSpecs).To(Equal(expectedReplicationSpecs))
			}

			By("Creating the Cluster", func() {
				Expect(k8sClient.Create(context.Background(), createdCluster)).To(Succeed())

				Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
					ClusterUpdateTimeout, interval).Should(BeTrue())

				doCommonStatusChecks()

				checkAtlasState(replicationSpecsCheckFunc)
			})

			By("Updating the cluster (multiple operations)", func() {
				delete(createdCluster.Spec.ReplicationSpecs[0].RegionsConfig, "US_WEST_2")
				createdCluster.Spec.ReplicationSpecs[0].RegionsConfig["US_WEST_1"] = mdbv1.RegionsConfig{AnalyticsNodes: int64ptr(0), ElectableNodes: int64ptr(2), Priority: int64ptr(6), ReadOnlyNodes: int64ptr(0)}
				config := createdCluster.Spec.ReplicationSpecs[0].RegionsConfig["US_EAST_1"]
				// Note, that Atlas has strict requirements to priorities - they must start with 7 and be in descending order over the regions
				config.Priority = int64ptr(7)
				createdCluster.Spec.ReplicationSpecs[0].RegionsConfig["US_EAST_1"] = config

				createdCluster.Spec.ProviderSettings.AutoScaling.Compute.MaxInstanceSize = "M30"

				performUpdate(ClusterUpdateTimeout)

				Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validateClusterUpdatingFunc()),
					ClusterUpdateTimeout, interval).Should(BeTrue())

				doCommonStatusChecks()

				checkAtlasState(replicationSpecsCheckFunc)
			})
		})
	})

	Describe("Create/Update the cluster (GCP)", func() {
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
				testutil.EventExists(k8sClient, createdCluster, "Warning", string(workflow.Internal), "name is invalid because must be set")

				lastGeneration++

			})

			By("Fixing the cluster", func() {
				createdCluster.Spec.Name = "fixed-cluster"

				Expect(k8sClient.Update(context.Background(), createdCluster)).To(Succeed())

				Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType)),
					20*time.Minute, interval).Should(BeTrue())

				doCommonStatusChecks()
				checkAtlasState()
			})
		})

		It("Should Succeed", func() {
			createdCluster = mdbv1.DefaultAWSCluster(namespace.Name, createdProject.Name)

			By(fmt.Sprintf("Creating the Cluster %s", kube.ObjectKeyFromObject(createdCluster)), func() {
				Expect(k8sClient.Create(context.Background(), createdCluster)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
					ClusterUpdateTimeout, interval).Should(BeTrue())

				doCommonStatusChecks()
				checkAtlasState()
			})

			By("Updating the Cluster labels", func() {
				createdCluster.Spec.Labels = []mdbv1.LabelSpec{{Key: "int-test", Value: "true"}}
				performUpdate(20 * time.Minute)
				doCommonStatusChecks()
				checkAtlasState()
			})

			By("Updating the Cluster backups settings", func() {
				createdCluster.Spec.ProviderBackupEnabled = boolptr(true)
				performUpdate(20 * time.Minute)
				doCommonStatusChecks()
				checkAtlasState(func(c *mongodbatlas.Cluster) {
					Expect(c.ProviderBackupEnabled).To(Equal(createdCluster.Spec.ProviderBackupEnabled))
				})
			})

			By("Decreasing the Cluster disk size", func() {
				createdCluster.Spec.DiskSizeGB = intptr(10)
				performUpdate(20 * time.Minute)
				doCommonStatusChecks()
				checkAtlasState(func(c *mongodbatlas.Cluster) {
					Expect(*c.DiskSizeGB).To(BeEquivalentTo(*createdCluster.Spec.DiskSizeGB))

					// check whether https://github.com/mongodb/go-client-mongodb-atlas/issues/140 is fixed
					Expect(c.DiskSizeGB).To(BeAssignableToTypeOf(float64ptr(0)), "DiskSizeGB is no longer a *float64, please check the spec!")
				})
			})

			By("Pausing the cluster", func() {
				createdCluster.Spec.Paused = boolptr(true)
				performUpdate(20 * time.Minute)
				doCommonStatusChecks()
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
				performUpdate(20 * time.Minute)
				doCommonStatusChecks()
				checkAtlasState(func(c *mongodbatlas.Cluster) {
					Expect(c.Paused).To(Equal(createdCluster.Spec.Paused))
				})
			})

			By("Checking that modifications were applied after unpausing", func() {
				doCommonStatusChecks()
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
					performUpdate(20 * time.Minute)
					doCommonStatusChecks()
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
					performUpdate(20 * time.Minute)
					doCommonStatusChecks()
					checkAtlasState()
				})
			})
		})
	})

	Describe("Create DBUser before cluster & check secrets", func() {
		It("Should Succeed", func() {
			By(fmt.Sprintf("Creating password Secret %s", UserPasswordSecret), func() {
				passwordSecret := buildPasswordSecret(UserPasswordSecret, DBUserPassword)
				Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())
			})

			createdDBUser := mdbv1.DefaultDBUser(namespace.Name, "test-db-user", createdProject.Name).WithPasswordSecret(UserPasswordSecret)
			By(fmt.Sprintf("Creating the Database User %s", kube.ObjectKeyFromObject(createdDBUser)), func() {
				Expect(k8sClient.Create(context.Background(), createdDBUser)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType))).Should(BeTrue())

				checkUserInAtlas(createdProject.ID(), *createdDBUser)
			})

			createdDBUserFakeScope := mdbv1.DefaultDBUser(namespace.Name, "test-db-user-fake-scope", createdProject.Name).
				WithPasswordSecret(UserPasswordSecret).
				WithScope(mdbv1.ClusterScopeType, "fake-cluster")
			By(fmt.Sprintf("Creating the Database User %s", kube.ObjectKeyFromObject(createdDBUserFakeScope)), func() {
				Expect(k8sClient.Create(context.Background(), createdDBUserFakeScope)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdDBUserFakeScope, status.FalseCondition(status.DatabaseUserReadyType).WithReason(string(workflow.DatabaseUserInvalidSpec))),
					20, intervalShort).Should(BeTrue())
			})
			checkNumberOfConnectionSecrets(k8sClient, *createdProject, 0)

			createdCluster = mdbv1.DefaultAWSCluster(namespace.Name, createdProject.Name)
			By(fmt.Sprintf("Creating the Cluster %s", kube.ObjectKeyFromObject(createdCluster)), func() {
				Expect(k8sClient.Create(context.Background(), createdCluster)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
					ClusterUpdateTimeout, interval).Should(BeTrue())

				doCommonStatusChecks()
				checkAtlasState()
			})

			By("Checking connection Secrets", func() {
				Expect(tryConnect(createdProject.ID(), *createdCluster, *createdDBUser)).To(Succeed())
				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 1)
				validateSecret(k8sClient, *createdProject, *createdCluster, *createdDBUser)
			})
		})
	})

	Describe("Create cluster, user, delete cluster and check secrets are removed", func() {
		It("Should Succeed", func() {
			createdCluster = mdbv1.DefaultAWSCluster(namespace.Name, createdProject.Name)
			By(fmt.Sprintf("Creating the Cluster %s", kube.ObjectKeyFromObject(createdCluster)), func() {
				Expect(k8sClient.Create(context.Background(), createdCluster)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
					ClusterUpdateTimeout, interval).Should(BeTrue())

				doCommonStatusChecks()
				checkAtlasState()
			})

			passwordSecret := buildPasswordSecret(UserPasswordSecret, DBUserPassword)
			Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())

			createdDBUser := mdbv1.DefaultDBUser(namespace.Name, "test-db-user", createdProject.Name).WithPasswordSecret(UserPasswordSecret)
			By(fmt.Sprintf("Creating the Database User %s", kube.ObjectKeyFromObject(createdDBUser)), func() {
				Expect(k8sClient.Create(context.Background(), createdDBUser)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType)),
					100, interval).Should(BeTrue())
			})

			By("Removing Atlas Cluster "+createdCluster.Name, func() {
				Expect(k8sClient.Delete(context.Background(), createdCluster)).To(Succeed())
				Eventually(checkAtlasClusterRemoved(createdProject.Status.ID, createdCluster.Spec.Name), 600, interval).Should(BeTrue())
			})

			By("Checking that Secrets got removed", func() {
				secretNames := []string{kube.NormalizeIdentifier(fmt.Sprintf("%s-%s-%s", createdProject.Spec.Name, createdCluster.Spec.Name, createdDBUser.Spec.Username))}
				Eventually(checkSecretsDontExist(namespace.Name, secretNames), 50, interval).Should(BeTrue())
				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 0)
			})

			// prevent cleanup from failing due to cluster already deleted
			createdCluster = nil
		})
	})

	Describe("Deleting the cluster (not cleaning Atlas)", func() {
		It("Should Succeed", func() {
			By(`Creating the cluster with retention policy "keep" first`, func() {
				createdCluster = mdbv1.DefaultAWSCluster(namespace.Name, createdProject.Name)
				createdCluster.ObjectMeta.Annotations = map[string]string{customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep}
				Expect(k8sClient.Create(context.Background(), createdCluster)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
					30*time.Minute, interval).Should(BeTrue())
			})
			By("Deleting the cluster - stays in Atlas", func() {
				Expect(k8sClient.Delete(context.Background(), createdCluster)).To(Succeed())
				time.Sleep(5 * time.Minute)
				Expect(checkAtlasClusterRemoved(createdProject.Status.ID, createdCluster.Spec.Name)()).Should(BeFalse())

				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 0)
			})
			By("Deleting the cluster in Atlas manually", func() {
				// We need to remove the cluster in Atlas manually to let project get removed
				_, err := atlasClient.Clusters.Delete(context.Background(), createdProject.ID(), createdCluster.Spec.Name)
				Expect(err).NotTo(HaveOccurred())
				Eventually(checkAtlasClusterRemoved(createdProject.Status.ID, createdCluster.Spec.Name), 600, interval).Should(BeTrue())
				createdCluster = nil
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
				return true
			}
		}

		return false
	}
}

func int64ptr(i int64) *int64 {
	return &i
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
