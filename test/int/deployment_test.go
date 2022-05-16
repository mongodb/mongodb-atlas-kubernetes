package int

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/connectionsecret"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasdeployment"
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

var _ = Describe("AtlasDeployment", Label("int", "AtlasDeployment"), func() {
	const (
		interval      = PollingInterval
		intervalShort = time.Second * 2
	)

	var (
		connectionSecret corev1.Secret
		createdProject   *mdbv1.AtlasProject
		createdCluster   *mdbv1.AtlasDeployment
		lastGeneration   int64
		manualDeletion   bool
	)

	BeforeEach(func() {
		prepareControllers()

		createdCluster = &mdbv1.AtlasDeployment{}

		lastGeneration = 0
		manualDeletion = false

		connectionSecret = corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-atlas-key",
				Namespace: namespace.Name,
				Labels: map[string]string{
					connectionsecret.TypeLabelKey: connectionsecret.CredLabelVal,
				},
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
		if manualDeletion && createdProject != nil {
			By("Deleting the cluster in Atlas manually", func() {
				// We need to remove the cluster in Atlas manually to let project get removed
				_, err := atlasClient.Clusters.Delete(context.Background(), createdProject.ID(), createdCluster.Spec.DeploymentSpec.Name)
				Expect(err).NotTo(HaveOccurred())
				Eventually(checkAtlasDeploymentRemoved(createdProject.Status.ID, createdCluster.Spec.DeploymentSpec.Name), 600, interval).Should(BeTrue())
				createdCluster = nil
			})
		}
		if createdProject != nil && createdProject.Status.ID != "" {
			if createdCluster != nil {
				By("Removing Atlas Cluster " + createdCluster.Name)
				Expect(k8sClient.Delete(context.Background(), createdCluster)).To(Succeed())

				Eventually(checkAtlasDeploymentRemoved(createdProject.Status.ID, createdCluster.GetClusterName()), 600, interval).Should(BeTrue())
			}

			By("Removing Atlas Project " + createdProject.Status.ID)
			Expect(k8sClient.Delete(context.Background(), createdProject)).To(Succeed())
			Eventually(checkAtlasProjectRemoved(createdProject.Status.ID), 60, interval).Should(BeTrue())
		}
		removeControllersAndNamespace()
	})

	doRegularClusterStatusChecks := func() {
		By("Checking observed Cluster state", func() {
			atlasCluster, _, err := atlasClient.Clusters.Get(context.Background(), createdProject.Status.ID, createdCluster.Spec.DeploymentSpec.Name)
			Expect(err).ToNot(HaveOccurred())

			Expect(createdCluster.Status.ConnectionStrings).NotTo(BeNil())
			Expect(createdCluster.Status.ConnectionStrings.Standard).To(Equal(atlasCluster.ConnectionStrings.Standard))
			Expect(createdCluster.Status.ConnectionStrings.StandardSrv).To(Equal(atlasCluster.ConnectionStrings.StandardSrv))
			Expect(createdCluster.Status.MongoDBVersion).To(Equal(atlasCluster.MongoDBVersion))
			Expect(createdCluster.Status.MongoURIUpdated).To(Equal(atlasCluster.MongoURIUpdated))
			Expect(createdCluster.Status.StateName).To(Equal("IDLE"))
			Expect(createdCluster.Status.Conditions).To(HaveLen(3))
			Expect(createdCluster.Status.Conditions).To(ConsistOf(testutil.MatchConditions(
				status.TrueCondition(status.ClusterReadyType),
				status.TrueCondition(status.ReadyType),
				status.TrueCondition(status.ValidationSucceeded),
			)))
			Expect(createdCluster.Status.ObservedGeneration).To(Equal(createdCluster.Generation))
			Expect(createdCluster.Status.ObservedGeneration).To(Equal(lastGeneration + 1))
		})
	}

	doAdvancedDeploymentStatusChecks := func() {
		By("Checking observed Advanced Cluster state", func() {
			atlasCluster, _, err := atlasClient.AdvancedClusters.Get(context.Background(), createdProject.Status.ID, createdCluster.Spec.AdvancedDeploymentSpec.Name)
			Expect(err).ToNot(HaveOccurred())

			Expect(createdCluster.Status.ConnectionStrings).NotTo(BeNil())
			Expect(createdCluster.Status.ConnectionStrings.Standard).To(Equal(atlasCluster.ConnectionStrings.Standard))
			Expect(createdCluster.Status.ConnectionStrings.StandardSrv).To(Equal(atlasCluster.ConnectionStrings.StandardSrv))
			Expect(createdCluster.Status.MongoDBVersion).To(Equal(atlasCluster.MongoDBVersion))
			Expect(createdCluster.Status.StateName).To(Equal("IDLE"))
			Expect(createdCluster.Status.Conditions).To(HaveLen(3))
			Expect(createdCluster.Status.Conditions).To(ConsistOf(testutil.MatchConditions(
				status.TrueCondition(status.ClusterReadyType),
				status.TrueCondition(status.ReadyType),
				status.TrueCondition(status.ValidationSucceeded),
			)))
			Expect(createdCluster.Status.ObservedGeneration).To(Equal(createdCluster.Generation))
			Expect(createdCluster.Status.ObservedGeneration).To(Equal(lastGeneration + 1))
		})
	}

	doServerlessClusterStatusChecks := func() {
		By("Checking observed Serverless state", func() {
			atlasCluster, _, err := atlasClient.ServerlessInstances.Get(context.Background(), createdProject.Status.ID, createdCluster.Spec.ServerlessSpec.Name)
			Expect(err).ToNot(HaveOccurred())

			Expect(createdCluster.Status.ConnectionStrings).NotTo(BeNil())
			Expect(createdCluster.Status.ConnectionStrings.Standard).To(Equal(atlasCluster.ConnectionStrings.Standard))
			Expect(createdCluster.Status.ConnectionStrings.StandardSrv).To(Equal(atlasCluster.ConnectionStrings.StandardSrv))
			Expect(createdCluster.Status.MongoDBVersion).To(Not(BeEmpty()))
			Expect(createdCluster.Status.StateName).To(Equal("IDLE"))
			Expect(createdCluster.Status.Conditions).To(HaveLen(3))
			Expect(createdCluster.Status.Conditions).To(ConsistOf(testutil.MatchConditions(
				status.TrueCondition(status.ClusterReadyType),
				status.TrueCondition(status.ReadyType),
				status.TrueCondition(status.ValidationSucceeded),
			)))
			Expect(createdCluster.Status.ObservedGeneration).To(Equal(createdCluster.Generation))
			Expect(createdCluster.Status.ObservedGeneration).To(Equal(lastGeneration + 1))
		})
	}

	checkAtlasState := func(additionalChecks ...func(c *mongodbatlas.Cluster)) {
		By("Verifying Cluster state in Atlas", func() {
			atlasCluster, _, err := atlasClient.Clusters.Get(context.Background(), createdProject.Status.ID, createdCluster.Spec.DeploymentSpec.Name)
			Expect(err).ToNot(HaveOccurred())

			mergedCluster, err := atlasdeployment.MergedCluster(*atlasCluster, createdCluster.Spec)
			Expect(err).ToNot(HaveOccurred())

			Expect(atlasdeployment.ClustersEqual(zap.S(), *atlasCluster, mergedCluster)).To(BeTrue())

			for _, check := range additionalChecks {
				check(atlasCluster)
			}
		})
	}

	checkAdvancedAtlasState := func(additionalChecks ...func(c *mongodbatlas.AdvancedCluster)) {
		By("Verifying Cluster state in Atlas", func() {
			atlasCluster, _, err := atlasClient.AdvancedClusters.Get(context.Background(), createdProject.Status.ID, createdCluster.Spec.AdvancedDeploymentSpec.Name)
			Expect(err).ToNot(HaveOccurred())

			mergedCluster, err := atlasdeployment.MergedAdvancedDeployment(*atlasCluster, createdCluster.Spec)
			Expect(err).ToNot(HaveOccurred())

			Expect(atlasdeployment.AdvancedDeploymentsEqual(zap.S(), *atlasCluster, mergedCluster)).To(BeTrue())

			for _, check := range additionalChecks {
				check(atlasCluster)
			}
		})
	}

	checkAdvancedDeploymentOptions := func(specOptions *mdbv1.ProcessArgs) {
		By("Checking that Atlas Advanced Options are equal to the Spec Options", func() {
			atlasOptions, _, err := atlasClient.Clusters.GetProcessArgs(context.Background(), createdProject.Status.ID, createdCluster.Spec.DeploymentSpec.Name)
			Expect(err).ToNot(HaveOccurred())

			Expect(specOptions.IsEqual(atlasOptions)).To(BeTrue())
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
				Expect(cluster.ReplicationSpecs[0].RegionsConfig[createdCluster.Spec.DeploymentSpec.ProviderSettings.RegionName]).NotTo(BeNil())
			}

			By(fmt.Sprintf("Creating the Cluster %s", kube.ObjectKeyFromObject(createdCluster)), func() {
				Expect(k8sClient.Create(context.Background(), createdCluster)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
					30*time.Minute, interval).Should(BeTrue())

				doRegularClusterStatusChecks()

				singleNumShard := func(cluster *mongodbatlas.Cluster) {
					Expect(cluster.ReplicationSpecs[0].NumShards).To(Equal(int64ptr(1)))
				}
				checkAtlasState(replicationSpecsCheck, singleNumShard)
			})

			By("Updating ReplicationSpecs", func() {
				createdCluster.Spec.DeploymentSpec.ReplicationSpecs = append(createdCluster.Spec.DeploymentSpec.ReplicationSpecs, mdbv1.ReplicationSpec{
					NumShards: int64ptr(2),
				})
				createdCluster.Spec.DeploymentSpec.ClusterType = "SHARDED"

				performUpdate(40 * time.Minute)
				doRegularClusterStatusChecks()

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

				doRegularClusterStatusChecks()
				checkAtlasState()
			})

			By("Increasing InstanceSize", func() {
				createdCluster.Spec.DeploymentSpec.ProviderSettings.InstanceSizeName = "M30"
				performUpdate(40 * time.Minute)
				doRegularClusterStatusChecks()
				checkAtlasState()
			})
		})
	})

	Describe("Create cluster & change it to GEOSHARDED", Label("int", "geosharded", "slow"), func() {
		It("Should Succeed", func() {
			expectedCluster := mdbv1.DefaultAWSCluster(namespace.Name, createdProject.Name)

			By(fmt.Sprintf("Creating the Cluster %s", kube.ObjectKeyFromObject(expectedCluster)), func() {
				createdCluster.ObjectMeta = expectedCluster.ObjectMeta
				Expect(k8sClient.Create(context.Background(), expectedCluster)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
					ClusterUpdateTimeout, interval).Should(BeTrue())

				doRegularClusterStatusChecks()
				checkAtlasState()
			})

			By("Change cluster to GEOSHARDED", func() {
				createdCluster.Spec.DeploymentSpec.ClusterType = "GEOSHARDED"
				createdCluster.Spec.DeploymentSpec.ReplicationSpecs = []mdbv1.ReplicationSpec{
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
				doRegularClusterStatusChecks()
				checkAtlasState()
			})
		})
	})

	Describe("Create/Update the cluster (more complex scenario)", func() {
		It("Should be created", func() {
			createdCluster = mdbv1.DefaultAWSCluster(namespace.Name, createdProject.Name)
			createdCluster.Spec.DeploymentSpec.ClusterType = mdbv1.TypeReplicaSet
			createdCluster.Spec.DeploymentSpec.AutoScaling = &mdbv1.AutoScalingSpec{
				Compute: &mdbv1.ComputeSpec{
					Enabled:          boolptr(true),
					ScaleDownEnabled: boolptr(true),
				},
				DiskGBEnabled: boolptr(false),
			}
			createdCluster.Spec.DeploymentSpec.ProviderSettings.AutoScaling = &mdbv1.AutoScalingSpec{
				Compute: &mdbv1.ComputeSpec{
					MaxInstanceSize: "M20",
					MinInstanceSize: "M10",
				},
			}
			createdCluster.Spec.DeploymentSpec.ProviderSettings.InstanceSizeName = "M10"
			createdCluster.Spec.DeploymentSpec.Labels = []common.LabelSpec{{Key: "createdBy", Value: "Atlas Operator"}}
			createdCluster.Spec.DeploymentSpec.ReplicationSpecs = []mdbv1.ReplicationSpec{{
				NumShards: int64ptr(1),
				ZoneName:  "Zone 1",
				// One interesting thing: if the regionsConfig is not empty - Atlas nullifies the 'providerSettings.regionName' field
				RegionsConfig: map[string]mdbv1.RegionsConfig{
					"US_EAST_1": {AnalyticsNodes: int64ptr(0), ElectableNodes: int64ptr(1), Priority: int64ptr(6), ReadOnlyNodes: int64ptr(0)},
					"US_WEST_2": {AnalyticsNodes: int64ptr(0), ElectableNodes: int64ptr(2), Priority: int64ptr(7), ReadOnlyNodes: int64ptr(0)},
				},
			}}
			createdCluster.Spec.DeploymentSpec.DiskSizeGB = intptr(10)

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

				doRegularClusterStatusChecks()

				checkAtlasState(replicationSpecsCheckFunc)
			})

			By("Updating the cluster (multiple operations)", func() {
				delete(createdCluster.Spec.DeploymentSpec.ReplicationSpecs[0].RegionsConfig, "US_WEST_2")
				createdCluster.Spec.DeploymentSpec.ReplicationSpecs[0].RegionsConfig["US_WEST_1"] = mdbv1.RegionsConfig{AnalyticsNodes: int64ptr(0), ElectableNodes: int64ptr(2), Priority: int64ptr(6), ReadOnlyNodes: int64ptr(0)}
				config := createdCluster.Spec.DeploymentSpec.ReplicationSpecs[0].RegionsConfig["US_EAST_1"]
				// Note, that Atlas has strict requirements to priorities - they must start with 7 and be in descending order over the regions
				config.Priority = int64ptr(7)
				createdCluster.Spec.DeploymentSpec.ReplicationSpecs[0].RegionsConfig["US_EAST_1"] = config

				createdCluster.Spec.DeploymentSpec.ProviderSettings.AutoScaling.Compute.MaxInstanceSize = "M30"

				performUpdate(ClusterUpdateTimeout)

				Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validateClusterUpdatingFunc()),
					ClusterUpdateTimeout, interval).Should(BeTrue())

				doRegularClusterStatusChecks()

				checkAtlasState(replicationSpecsCheckFunc)
			})

			By("Disable cluster and disk AutoScaling", func() {
				createdCluster.Spec.DeploymentSpec.AutoScaling = &mdbv1.AutoScalingSpec{
					Compute: &mdbv1.ComputeSpec{
						Enabled:          boolptr(false),
						ScaleDownEnabled: boolptr(false),
					},
					DiskGBEnabled: boolptr(false),
				}

				performUpdate(ClusterUpdateTimeout)

				Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validateClusterUpdatingFunc()),
					ClusterUpdateTimeout, interval).Should(BeTrue())

				doRegularClusterStatusChecks()

				checkAtlasState(func(c *mongodbatlas.Cluster) {
					cluster, err := createdCluster.Spec.Cluster()
					Expect(err).NotTo(HaveOccurred())

					Expect(c.AutoScaling.Compute).To(Equal(cluster.AutoScaling.Compute))
				})
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
				createdCluster.Spec.DeploymentSpec.Name = "fixed-cluster"

				Expect(k8sClient.Update(context.Background(), createdCluster)).To(Succeed())

				Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType)),
					20*time.Minute, interval).Should(BeTrue())

				doRegularClusterStatusChecks()
				checkAtlasState()
			})
		})

		It("Should Succeed", func() {
			createdCluster = mdbv1.DefaultAWSCluster(namespace.Name, createdProject.Name)

			By(fmt.Sprintf("Creating the Cluster %s", kube.ObjectKeyFromObject(createdCluster)), func() {
				Expect(k8sClient.Create(context.Background(), createdCluster)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
					ClusterUpdateTimeout, interval).Should(BeTrue())

				doRegularClusterStatusChecks()
				checkAtlasState()
			})

			By("Updating the Cluster labels", func() {
				createdCluster.Spec.DeploymentSpec.Labels = []common.LabelSpec{{Key: "int-test", Value: "true"}}
				performUpdate(20 * time.Minute)
				doRegularClusterStatusChecks()
				checkAtlasState()
			})

			By("Updating the Cluster backups settings", func() {
				createdCluster.Spec.DeploymentSpec.ProviderBackupEnabled = boolptr(true)
				performUpdate(20 * time.Minute)
				doRegularClusterStatusChecks()
				checkAtlasState(func(c *mongodbatlas.Cluster) {
					Expect(c.ProviderBackupEnabled).To(Equal(createdCluster.Spec.DeploymentSpec.ProviderBackupEnabled))
				})
			})

			By("Decreasing the Cluster disk size", func() {
				createdCluster.Spec.DeploymentSpec.DiskSizeGB = intptr(10)
				performUpdate(20 * time.Minute)
				doRegularClusterStatusChecks()
				checkAtlasState(func(c *mongodbatlas.Cluster) {
					Expect(*c.DiskSizeGB).To(BeEquivalentTo(*createdCluster.Spec.DeploymentSpec.DiskSizeGB))

					// check whether https://github.com/mongodb/go-client-mongodb-atlas/issues/140 is fixed
					Expect(c.DiskSizeGB).To(BeAssignableToTypeOf(float64ptr(0)), "DiskSizeGB is no longer a *float64, please check the spec!")
				})
			})

			By("Pausing the cluster", func() {
				createdCluster.Spec.DeploymentSpec.Paused = boolptr(true)
				performUpdate(20 * time.Minute)
				doRegularClusterStatusChecks()
				checkAtlasState(func(c *mongodbatlas.Cluster) {
					Expect(c.Paused).To(Equal(createdCluster.Spec.DeploymentSpec.Paused))
				})
			})

			By("Updating the Cluster configuration while paused (should fail)", func() {
				createdCluster.Spec.DeploymentSpec.ProviderBackupEnabled = boolptr(false)

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
				createdCluster.Spec.DeploymentSpec.Paused = boolptr(false)
				performUpdate(20 * time.Minute)
				doRegularClusterStatusChecks()
				checkAtlasState(func(c *mongodbatlas.Cluster) {
					Expect(c.Paused).To(Equal(createdCluster.Spec.DeploymentSpec.Paused))
				})
			})

			By("Checking that modifications were applied after unpausing", func() {
				doRegularClusterStatusChecks()
				checkAtlasState(func(c *mongodbatlas.Cluster) {
					Expect(c.ProviderBackupEnabled).To(Equal(createdCluster.Spec.DeploymentSpec.ProviderBackupEnabled))
				})
			})

			By("Setting incorrect instance size (should fail)", func() {
				oldSizeName := createdCluster.Spec.DeploymentSpec.ProviderSettings.InstanceSizeName
				createdCluster.Spec.DeploymentSpec.ProviderSettings.InstanceSizeName = "M42"

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
					createdCluster.Spec.DeploymentSpec.ProviderSettings.InstanceSizeName = oldSizeName
					performUpdate(20 * time.Minute)
					doRegularClusterStatusChecks()
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

				doRegularClusterStatusChecks()
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

				doRegularClusterStatusChecks()
				checkAtlasState()
			})

			passwordSecret := buildPasswordSecret(UserPasswordSecret, DBUserPassword)
			Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())

			createdDBUser := mdbv1.DefaultDBUser(namespace.Name, "test-db-user", createdProject.Name).WithPasswordSecret(UserPasswordSecret)
			By(fmt.Sprintf("Creating the Database User %s", kube.ObjectKeyFromObject(createdDBUser)), func() {
				Expect(k8sClient.Create(context.Background(), createdDBUser)).ToNot(HaveOccurred())
				Eventually(testutil.WaitFor(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType)),
					DBUserUpdateTimeout, interval).Should(BeTrue())
			})

			By("Removing Atlas Cluster "+createdCluster.Name, func() {
				Expect(k8sClient.Delete(context.Background(), createdCluster)).To(Succeed())
				Eventually(checkAtlasDeploymentRemoved(createdProject.Status.ID, createdCluster.Spec.DeploymentSpec.Name), 600, interval).Should(BeTrue())
			})

			By("Checking that Secrets got removed", func() {
				secretNames := []string{kube.NormalizeIdentifier(fmt.Sprintf("%s-%s-%s", createdProject.Spec.Name, createdCluster.Spec.DeploymentSpec.Name, createdDBUser.Spec.Username))}
				createdCluster = nil // prevent cleanup from failing due to cluster already deleted
				Eventually(checkSecretsDontExist(namespace.Name, secretNames), 50, interval).Should(BeTrue())
				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 0)
			})
		})
	})

	Describe("Deleting the cluster (not cleaning Atlas)", func() {
		It("Should Succeed", func() {
			By(`Creating the cluster with retention policy "keep" first`, func() {
				createdCluster = mdbv1.DefaultAWSCluster(namespace.Name, createdProject.Name).Lightweight()
				createdCluster.ObjectMeta.Annotations = map[string]string{customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep}
				manualDeletion = true // We need to remove the cluster in Atlas manually to let project get removed
				Expect(k8sClient.Create(context.Background(), createdCluster)).ToNot(HaveOccurred())

				Eventually(
					func(g Gomega) {
						success := testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validateClusterCreatingFuncGContext(g))()
						g.Expect(success).To(BeTrue())
					}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(Succeed())
			})
			By("Deleting the cluster - stays in Atlas", func() {
				Expect(k8sClient.Delete(context.Background(), createdCluster)).To(Succeed())
				time.Sleep(5 * time.Minute)
				Expect(checkAtlasDeploymentRemoved(createdProject.Status.ID, createdCluster.Spec.DeploymentSpec.Name)()).Should(BeFalse())
				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 0)
			})
		})
	})

	Describe("Setting the cluster skip annotation should skip reconciliations.", func() {
		It("Should Succeed", func() {

			By(`Creating the cluster with reconciliation policy "skip" first`, func() {
				createdCluster = mdbv1.DefaultAWSCluster(namespace.Name, createdProject.Name).Lightweight()
				Expect(k8sClient.Create(context.Background(), createdCluster)).ToNot(HaveOccurred())

				Eventually(
					func(g Gomega) {
						success := testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validateClusterCreatingFuncGContext(g))()
						g.Expect(success).To(BeTrue())
					}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(Succeed())

				createdCluster.ObjectMeta.Annotations = map[string]string{customresource.ReconciliationPolicyAnnotation: customresource.ReconciliationPolicySkip}
				createdCluster.Spec.DeploymentSpec.Labels = append(createdCluster.Spec.DeploymentSpec.Labels, common.LabelSpec{
					Key:   "some-key",
					Value: "some-value",
				})

				ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
				defer cancel()

				containsLabel := func(ac *mongodbatlas.Cluster) bool {
					for _, label := range ac.Labels {
						if label.Key == "some-key" && label.Value == "some-value" {
							return true
						}
					}
					return false
				}

				Expect(k8sClient.Update(context.Background(), createdCluster)).ToNot(HaveOccurred())
				Eventually(testutil.WaitForAtlasDeploymentStateToNotBeReached(ctx, atlasClient, createdProject.Name, createdCluster.GetClusterName(), containsLabel))
			})
		})
	})

	Describe("Create advanced cluster", func() {
		It("Should Succeed", func() {
			createdCluster = mdbv1.DefaultAwsAdvancedDeployment(namespace.Name, createdProject.Name)

			By(fmt.Sprintf("Creating the Advanced Cluster %s", kube.ObjectKeyFromObject(createdCluster)), func() {
				Expect(k8sClient.Create(context.Background(), createdCluster)).ToNot(HaveOccurred())

				Eventually(
					func(g Gomega) {
						success := testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validateClusterCreatingFuncGContext(g))()
						g.Expect(success).To(BeTrue())
					}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(Succeed())

				doAdvancedDeploymentStatusChecks()
				checkAdvancedAtlasState()
			})
		})
	})

	Describe("Set advanced cluster options", func() {
		It("Should Succeed", func() {
			createdCluster = mdbv1.DefaultAWSCluster(namespace.Name, createdProject.Name)
			createdCluster.Spec.ProcessArgs = &mdbv1.ProcessArgs{
				JavascriptEnabled:  boolptr(true),
				DefaultReadConcern: "available",
			}

			By(fmt.Sprintf("Creating the Cluster with Advanced Options %s", kube.ObjectKeyFromObject(createdCluster)), func() {
				Expect(k8sClient.Create(context.Background(), createdCluster)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
					30*time.Minute, interval).Should(BeTrue())

				doRegularClusterStatusChecks()
				checkAdvancedDeploymentOptions(createdCluster.Spec.ProcessArgs)
			})

			By("Updating Advanced Cluster Options", func() {
				createdCluster.Spec.ProcessArgs.JavascriptEnabled = boolptr(false)
				performUpdate(40 * time.Minute)
				doRegularClusterStatusChecks()
				checkAdvancedDeploymentOptions(createdCluster.Spec.ProcessArgs)
			})
		})
	})

	Describe("Create serverless instance", func() {
		It("Should Succeed", func() {
			createdCluster = mdbv1.NewDefaultAWSServerlessInstance(namespace.Name, createdProject.Name)

			By(fmt.Sprintf("Creating the Serverless Instance %s", kube.ObjectKeyFromObject(createdCluster)), func() {
				Expect(k8sClient.Create(context.Background(), createdCluster)).ToNot(HaveOccurred())

				Eventually(
					func(g Gomega) {
						success := testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validateClusterCreatingFuncGContext(g))()
						g.Expect(success).To(BeTrue())
					}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(Succeed())

				doServerlessClusterStatusChecks()
			})
		})
	})

	Describe("Create default cluster with backups enabled", func() {
		It("Should succeed", func() {
			backupPolicyDefault := &mdbv1.AtlasBackupPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "policy-1",
					Namespace: namespace.Name,
				},
				Spec: mdbv1.AtlasBackupPolicySpec{
					Items: []mdbv1.AtlasBackupPolicyItem{
						{
							FrequencyType:     "weekly",
							FrequencyInterval: 1,
							RetentionUnit:     "days",
							RetentionValue:    7,
						},
					},
				},
				Status: mdbv1.AtlasBackupPolicyStatus{},
			}

			backupScheduleDefault := &mdbv1.AtlasBackupSchedule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "schedule-1",
					Namespace: namespace.Name,
				},
				Spec: mdbv1.AtlasBackupScheduleSpec{
					AutoExportEnabled: false,
					PolicyRef: common.ResourceRefNamespaced{
						Name:      backupPolicyDefault.Name,
						Namespace: backupPolicyDefault.Namespace,
					},
					ReferenceHourOfDay:    12,
					ReferenceMinuteOfHour: 10,
					RestoreWindowDays:     5,
					UpdateSnapshots:       false,
					Export:                mdbv1.AtlasBackupExportSpec{FrequencyType: "MONTHLY"},
				},
			}

			Expect(k8sClient.Create(context.Background(), backupPolicyDefault)).NotTo(HaveOccurred())
			Expect(k8sClient.Create(context.Background(), backupScheduleDefault)).NotTo(HaveOccurred())

			createdCluster = mdbv1.DefaultAWSCluster(namespace.Name, createdProject.Name).WithBackupScheduleRef(common.ResourceRefNamespaced{
				Name:      backupScheduleDefault.Name,
				Namespace: backupScheduleDefault.Namespace,
			})

			By(fmt.Sprintf("Creating cluster with backups enabled: %s", kube.ObjectKeyFromObject(createdCluster)), func() {
				Expect(k8sClient.Create(context.Background(), createdCluster)).NotTo(HaveOccurred())

				// Do not use Gomega function here like func(g Gomega) as it seems to hang when tests run in parallel
				Eventually(
					func() error {
						cluster, _, err := atlasClient.Clusters.Get(context.Background(), createdProject.ID(), createdCluster.Spec.DeploymentSpec.Name)
						if err != nil {
							return err
						}
						if cluster.StateName != "IDLE" {
							return errors.New("cluster is not IDLE yet")
						}
						time.Sleep(10 * time.Second)
						return nil
					}).WithTimeout(30 * time.Minute).WithPolling(5 * time.Second).Should(Not(HaveOccurred()))

				Eventually(func() error {
					actualPolicy, _, err := atlasClient.CloudProviderSnapshotBackupPolicies.Get(context.Background(), createdProject.ID(), createdCluster.Spec.DeploymentSpec.Name)
					if err != nil {
						return err
					}
					if len(actualPolicy.Policies[0].PolicyItems) == 0 {
						return errors.New("policies == 0")
					}
					ap := actualPolicy.Policies[0].PolicyItems[0]
					cp := backupPolicyDefault.Spec.Items[0]
					if ap.FrequencyType != cp.FrequencyType {
						return fmt.Errorf("incorrect frequency type. got: %v. expected: %v", ap.FrequencyType, cp.FrequencyType)
					}
					if ap.FrequencyInterval != cp.FrequencyInterval {
						return fmt.Errorf("incorrect frequency interval. got: %v. expected: %v", ap.FrequencyInterval, cp.FrequencyInterval)
					}
					if ap.RetentionValue != cp.RetentionValue {
						return fmt.Errorf("incorrect retention value. got: %v. expected: %v", ap.RetentionValue, cp.RetentionValue)
					}
					if ap.RetentionUnit != cp.RetentionUnit {
						return fmt.Errorf("incorrect retention unit. got: %v. expected: %v", ap.RetentionUnit, cp.RetentionUnit)
					}
					return nil
				}).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Not(HaveOccurred()))

			})
		})
	})
})

func validateClusterCreatingFunc() func(a mdbv1.AtlasCustomResource) {
	startedCreation := false
	return func(a mdbv1.AtlasCustomResource) {
		c := a.(*mdbv1.AtlasDeployment)
		if c.Status.StateName != "" {
			startedCreation = true
		}
		// When the create request has been made to Atlas - we expect the following status
		if startedCreation {
			Expect(c.Status.StateName).To(Or(Equal("CREATING"), Equal("IDLE")), fmt.Sprintf("Current conditions: %+v", c.Status.Conditions))
			expectedConditionsMatchers := testutil.MatchConditions(
				status.FalseCondition(status.ClusterReadyType).WithReason(string(workflow.ClusterCreating)).WithMessageRegexp("cluster is provisioning"),
				status.FalseCondition(status.ReadyType),
				status.TrueCondition(status.ValidationSucceeded),
			)
			Expect(c.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
		} else {
			// Otherwise there could have been some exception in Atlas on creation - let's check the conditions
			condition, ok := testutil.FindConditionByType(c.Status.Conditions, status.ClusterReadyType)
			Expect(ok).To(BeFalse(), fmt.Sprintf("Unexpected condition: %v", condition))
		}
	}
}

func validateClusterCreatingFuncGContext(g Gomega) func(a mdbv1.AtlasCustomResource) {
	startedCreation := false
	return func(a mdbv1.AtlasCustomResource) {
		c := a.(*mdbv1.AtlasDeployment)
		if c.Status.StateName != "" {
			startedCreation = true
		}
		// When the create request has been made to Atlas - we expect the following status
		if startedCreation {
			g.Expect(c.Status.StateName).To(Or(Equal("CREATING"), Equal("IDLE")), fmt.Sprintf("Current conditions: %+v", c.Status.Conditions))
			expectedConditionsMatchers := testutil.MatchConditions(
				status.FalseCondition(status.ClusterReadyType).WithReason(string(workflow.ClusterCreating)).WithMessageRegexp("cluster is provisioning"),
				status.FalseCondition(status.ReadyType),
				status.TrueCondition(status.ValidationSucceeded),
			)
			g.Expect(c.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
		} else {
			// Otherwise there could have been some exception in Atlas on creation - let's check the conditions
			condition, ok := testutil.FindConditionByType(c.Status.Conditions, status.ClusterReadyType)
			g.Expect(ok).To(BeFalse(), fmt.Sprintf("Unexpected condition: %v", condition))
		}
	}
}

func validateClusterUpdatingFunc() func(a mdbv1.AtlasCustomResource) {
	isIdle := true
	return func(a mdbv1.AtlasCustomResource) {
		c := a.(*mdbv1.AtlasDeployment)
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
				status.TrueCondition(status.ValidationSucceeded),
			)
			Expect(c.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
		}
	}
}

// checkAtlasDeploymentRemoved returns true if the Atlas Cluster is removed from Atlas. Note the behavior: the cluster
// is removed from Atlas as soon as the DELETE API call has been made. This is different from the case when the
// cluster is terminated from UI (in this case GET request succeeds while the cluster is being terminated)
func checkAtlasDeploymentRemoved(projectID string, clusterName string) func() bool {
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
