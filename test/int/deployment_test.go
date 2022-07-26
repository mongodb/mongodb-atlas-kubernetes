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
	// Set this to true if you are debugging deployment creation.
	// This may not help much if there was the update though...
	DeploymentDevMode       = false
	DeploymentUpdateTimeout = 40 * time.Minute
)

var _ = Describe("AtlasDeployment", Label("int", "AtlasDeployment"), func() {
	const (
		interval      = PollingInterval
		intervalShort = time.Second * 2
	)

	var (
		connectionSecret  corev1.Secret
		createdProject    *mdbv1.AtlasProject
		createdDeployment *mdbv1.AtlasDeployment
		lastGeneration    int64
		manualDeletion    bool
	)

	BeforeEach(func() {
		prepareControllers()

		createdDeployment = &mdbv1.AtlasDeployment{}

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
		if DeploymentDevMode {
			// While developing tests we need to reuse the same project
			createdProject.Spec.Name = "dev-test atlas-project"
		}
		By("Creating the project " + createdProject.Name)
		Expect(k8sClient.Create(context.Background(), createdProject)).To(Succeed())
		Eventually(testutil.WaitFor(k8sClient, createdProject, status.TrueCondition(status.ReadyType)),
			ProjectCreationTimeout, intervalShort).Should(BeTrue())
	})

	AfterEach(func() {
		if DeploymentDevMode {
			return
		}
		if manualDeletion && createdProject != nil {
			By("Deleting the deployment in Atlas manually", func() {
				// We need to remove the deployment in Atlas manually to let project get removed
				_, err := atlasClient.Clusters.Delete(context.Background(), createdProject.ID(), createdDeployment.Spec.DeploymentSpec.Name)
				Expect(err).NotTo(HaveOccurred())
				Eventually(checkAtlasDeploymentRemoved(createdProject.Status.ID, createdDeployment.Spec.DeploymentSpec.Name), 600, interval).Should(BeTrue())
				createdDeployment = nil
			})
		}
		if createdProject != nil && createdProject.Status.ID != "" {
			if createdDeployment != nil {
				By("Removing Atlas Deployment " + createdDeployment.Name)
				Expect(k8sClient.Delete(context.Background(), createdDeployment)).To(Succeed())

				Eventually(checkAtlasDeploymentRemoved(createdProject.Status.ID, createdDeployment.GetDeploymentName()), 600, interval).Should(BeTrue())
			}

			By("Removing Atlas Project " + createdProject.Status.ID)
			Expect(k8sClient.Delete(context.Background(), createdProject)).To(Succeed())
			Eventually(checkAtlasProjectRemoved(createdProject.Status.ID), 60, interval).Should(BeTrue())
		}
		removeControllersAndNamespace()
	})

	doRegularDeploymentStatusChecks := func() {
		By("Checking observed Deployment state", func() {
			atlasDeployment, _, err := atlasClient.Clusters.Get(context.Background(), createdProject.Status.ID, createdDeployment.Spec.DeploymentSpec.Name)
			Expect(err).ToNot(HaveOccurred())

			Expect(createdDeployment.Status.ConnectionStrings).NotTo(BeNil())
			Expect(createdDeployment.Status.ConnectionStrings.Standard).To(Equal(atlasDeployment.ConnectionStrings.Standard))
			Expect(createdDeployment.Status.ConnectionStrings.StandardSrv).To(Equal(atlasDeployment.ConnectionStrings.StandardSrv))
			Expect(createdDeployment.Status.MongoDBVersion).To(Equal(atlasDeployment.MongoDBVersion))
			Expect(createdDeployment.Status.MongoURIUpdated).To(Equal(atlasDeployment.MongoURIUpdated))
			Expect(createdDeployment.Status.StateName).To(Equal("IDLE"))
			Expect(createdDeployment.Status.Conditions).To(HaveLen(3))
			Expect(createdDeployment.Status.Conditions).To(ConsistOf(testutil.MatchConditions(
				status.TrueCondition(status.DeploymentReadyType),
				status.TrueCondition(status.ReadyType),
				status.TrueCondition(status.ValidationSucceeded),
			)))
			Expect(createdDeployment.Status.ObservedGeneration).To(Equal(createdDeployment.Generation))
			Expect(createdDeployment.Status.ObservedGeneration).To(Equal(lastGeneration + 1))
		})
	}

	doAdvancedDeploymentStatusChecks := func() {
		By("Checking observed Advanced Deployment state", func() {
			atlasDeployment, _, err := atlasClient.AdvancedClusters.Get(context.Background(), createdProject.Status.ID, createdDeployment.Spec.AdvancedDeploymentSpec.Name)
			Expect(err).ToNot(HaveOccurred())

			Expect(createdDeployment.Status.ConnectionStrings).NotTo(BeNil())
			Expect(createdDeployment.Status.ConnectionStrings.Standard).To(Equal(atlasDeployment.ConnectionStrings.Standard))
			Expect(createdDeployment.Status.ConnectionStrings.StandardSrv).To(Equal(atlasDeployment.ConnectionStrings.StandardSrv))
			Expect(createdDeployment.Status.MongoDBVersion).To(Equal(atlasDeployment.MongoDBVersion))
			Expect(createdDeployment.Status.StateName).To(Equal("IDLE"))
			Expect(createdDeployment.Status.Conditions).To(HaveLen(3))
			Expect(createdDeployment.Status.Conditions).To(ConsistOf(testutil.MatchConditions(
				status.TrueCondition(status.DeploymentReadyType),
				status.TrueCondition(status.ReadyType),
				status.TrueCondition(status.ValidationSucceeded),
			)))
			Expect(createdDeployment.Status.ObservedGeneration).To(Equal(createdDeployment.Generation))
			Expect(createdDeployment.Status.ObservedGeneration).To(Equal(lastGeneration + 1))
		})
	}

	doServerlessDeploymentStatusChecks := func() {
		By("Checking observed Serverless state", func() {
			atlasDeployment, _, err := atlasClient.ServerlessInstances.Get(context.Background(), createdProject.Status.ID, createdDeployment.Spec.ServerlessSpec.Name)
			Expect(err).ToNot(HaveOccurred())

			Expect(createdDeployment.Status.ConnectionStrings).NotTo(BeNil())
			Expect(createdDeployment.Status.ConnectionStrings.Standard).To(Equal(atlasDeployment.ConnectionStrings.Standard))
			Expect(createdDeployment.Status.ConnectionStrings.StandardSrv).To(Equal(atlasDeployment.ConnectionStrings.StandardSrv))
			Expect(createdDeployment.Status.MongoDBVersion).To(Not(BeEmpty()))
			Expect(createdDeployment.Status.StateName).To(Equal("IDLE"))
			Expect(createdDeployment.Status.Conditions).To(HaveLen(3))
			Expect(createdDeployment.Status.Conditions).To(ConsistOf(testutil.MatchConditions(
				status.TrueCondition(status.DeploymentReadyType),
				status.TrueCondition(status.ReadyType),
				status.TrueCondition(status.ValidationSucceeded),
			)))
			Expect(createdDeployment.Status.ObservedGeneration).To(Equal(createdDeployment.Generation))
			Expect(createdDeployment.Status.ObservedGeneration).To(Equal(lastGeneration + 1))
		})
	}

	checkAtlasState := func(additionalChecks ...func(c *mongodbatlas.Cluster)) {
		By("Verifying Deployment state in Atlas", func() {
			atlasDeployment, _, err := atlasClient.Clusters.Get(context.Background(), createdProject.Status.ID, createdDeployment.Spec.DeploymentSpec.Name)
			Expect(err).ToNot(HaveOccurred())

			mergedDeployment, err := atlasdeployment.MergedDeployment(*atlasDeployment, createdDeployment.Spec)
			Expect(err).ToNot(HaveOccurred())

			Expect(atlasdeployment.DeploymentsEqual(zap.S(), *atlasDeployment, mergedDeployment)).To(BeTrue())

			for _, check := range additionalChecks {
				check(atlasDeployment)
			}
		})
	}

	checkAdvancedAtlasState := func(additionalChecks ...func(c *mongodbatlas.AdvancedCluster)) {
		By("Verifying Deployment state in Atlas", func() {
			atlasDeployment, _, err := atlasClient.AdvancedClusters.Get(context.Background(), createdProject.Status.ID, createdDeployment.Spec.AdvancedDeploymentSpec.Name)
			Expect(err).ToNot(HaveOccurred())

			mergedDeployment, err := atlasdeployment.MergedAdvancedDeployment(*atlasDeployment, createdDeployment.Spec)
			Expect(err).ToNot(HaveOccurred())

			Expect(atlasdeployment.AdvancedDeploymentsEqual(zap.S(), *atlasDeployment, mergedDeployment)).To(BeTrue())

			for _, check := range additionalChecks {
				check(atlasDeployment)
			}
		})
	}

	checkAdvancedDeploymentOptions := func(specOptions *mdbv1.ProcessArgs) {
		By("Checking that Atlas Advanced Options are equal to the Spec Options", func() {
			atlasOptions, _, err := atlasClient.Clusters.GetProcessArgs(context.Background(), createdProject.Status.ID, createdDeployment.Spec.DeploymentSpec.Name)
			Expect(err).ToNot(HaveOccurred())

			Expect(specOptions.IsEqual(atlasOptions)).To(BeTrue())
		})
	}

	performUpdate := func(timeout interface{}) {
		Expect(k8sClient.Update(context.Background(), createdDeployment)).To(Succeed())

		Eventually(testutil.WaitFor(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentUpdatingFunc()),
			timeout, interval).Should(BeTrue())

		lastGeneration++
	}

	Describe("Create deployment & change ReplicationSpecs", func() {
		It("Should Succeed", func() {
			createdDeployment = mdbv1.DefaultAWSDeployment(namespace.Name, createdProject.Name)

			// Atlas will add some defaults in case the Atlas Operator doesn't set them
			replicationSpecsCheck := func(deployment *mongodbatlas.Cluster) {
				Expect(deployment.ReplicationSpecs).To(HaveLen(1))
				Expect(deployment.ReplicationSpecs[0].ID).NotTo(BeNil())
				Expect(deployment.ReplicationSpecs[0].ZoneName).To(Equal("Zone 1"))
				Expect(deployment.ReplicationSpecs[0].RegionsConfig).To(HaveLen(1))
				Expect(deployment.ReplicationSpecs[0].RegionsConfig[createdDeployment.Spec.DeploymentSpec.ProviderSettings.RegionName]).NotTo(BeNil())
			}

			By(fmt.Sprintf("Creating the Deployment %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				Expect(k8sClient.Create(context.Background(), createdDeployment)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc()),
					30*time.Minute, interval).Should(BeTrue())

				doRegularDeploymentStatusChecks()

				singleNumShard := func(deployment *mongodbatlas.Cluster) {
					Expect(deployment.ReplicationSpecs[0].NumShards).To(Equal(int64ptr(1)))
				}
				checkAtlasState(replicationSpecsCheck, singleNumShard)
			})

			By("Updating ReplicationSpecs", func() {
				createdDeployment.Spec.DeploymentSpec.ReplicationSpecs = append(createdDeployment.Spec.DeploymentSpec.ReplicationSpecs, mdbv1.ReplicationSpec{
					NumShards: int64ptr(2),
				})
				createdDeployment.Spec.DeploymentSpec.ClusterType = "SHARDED"

				performUpdate(40 * time.Minute)
				doRegularDeploymentStatusChecks()

				twoNumShard := func(deployment *mongodbatlas.Cluster) {
					Expect(deployment.ReplicationSpecs[0].NumShards).To(Equal(int64ptr(2)))
				}
				// ReplicationSpecs has the same defaults but the number of shards has changed
				checkAtlasState(replicationSpecsCheck, twoNumShard)
			})
		})
	})

	Describe("Create deployment & increase DiskSizeGB", func() {
		It("Should Succeed", func() {
			expectedDeployment := mdbv1.DefaultAWSDeployment(namespace.Name, createdProject.Name)

			By(fmt.Sprintf("Creating the Deployment %s", kube.ObjectKeyFromObject(expectedDeployment)), func() {
				createdDeployment.ObjectMeta = expectedDeployment.ObjectMeta
				Expect(k8sClient.Create(context.Background(), expectedDeployment)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc()),
					1800, interval).Should(BeTrue())

				doRegularDeploymentStatusChecks()
				checkAtlasState()
			})

			By("Increasing InstanceSize", func() {
				createdDeployment.Spec.DeploymentSpec.ProviderSettings.InstanceSizeName = "M30"
				performUpdate(40 * time.Minute)
				doRegularDeploymentStatusChecks()
				checkAtlasState()
			})
		})
	})

	Describe("Create deployment & change it to GEOSHARDED", Label("int", "geosharded", "slow"), func() {
		It("Should Succeed", func() {
			expectedDeployment := mdbv1.DefaultAWSDeployment(namespace.Name, createdProject.Name)

			By(fmt.Sprintf("Creating the Deployment %s", kube.ObjectKeyFromObject(expectedDeployment)), func() {
				createdDeployment.ObjectMeta = expectedDeployment.ObjectMeta
				Expect(k8sClient.Create(context.Background(), expectedDeployment)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc()),
					DeploymentUpdateTimeout, interval).Should(BeTrue())

				doRegularDeploymentStatusChecks()
				checkAtlasState()
			})

			By("Change deployment to GEOSHARDED", func() {
				createdDeployment.Spec.DeploymentSpec.ClusterType = "GEOSHARDED"
				createdDeployment.Spec.DeploymentSpec.ReplicationSpecs = []mdbv1.ReplicationSpec{
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
				doRegularDeploymentStatusChecks()
				checkAtlasState()
			})
		})
	})

	Describe("Create/Update the deployment (more complex scenario)", func() {
		It("Should be created", func() {
			createdDeployment = mdbv1.DefaultAWSDeployment(namespace.Name, createdProject.Name)
			createdDeployment.Spec.DeploymentSpec.ClusterType = mdbv1.TypeReplicaSet
			createdDeployment.Spec.DeploymentSpec.AutoScaling = &mdbv1.AutoScalingSpec{
				Compute: &mdbv1.ComputeSpec{
					Enabled:          boolptr(true),
					ScaleDownEnabled: boolptr(true),
				},
				DiskGBEnabled: boolptr(false),
			}
			createdDeployment.Spec.DeploymentSpec.ProviderSettings.AutoScaling = &mdbv1.AutoScalingSpec{
				Compute: &mdbv1.ComputeSpec{
					MaxInstanceSize: "M20",
					MinInstanceSize: "M10",
				},
			}
			createdDeployment.Spec.DeploymentSpec.ProviderSettings.InstanceSizeName = "M10"
			createdDeployment.Spec.DeploymentSpec.Labels = []common.LabelSpec{{Key: "createdBy", Value: "Atlas Operator"}}
			createdDeployment.Spec.DeploymentSpec.ReplicationSpecs = []mdbv1.ReplicationSpec{{
				NumShards: int64ptr(1),
				ZoneName:  "Zone 1",
				// One interesting thing: if the regionsConfig is not empty - Atlas nullifies the 'providerSettings.regionName' field
				RegionsConfig: map[string]mdbv1.RegionsConfig{
					"US_EAST_1": {AnalyticsNodes: int64ptr(0), ElectableNodes: int64ptr(1), Priority: int64ptr(6), ReadOnlyNodes: int64ptr(0)},
					"US_WEST_2": {AnalyticsNodes: int64ptr(0), ElectableNodes: int64ptr(2), Priority: int64ptr(7), ReadOnlyNodes: int64ptr(0)},
				},
			}}
			createdDeployment.Spec.DeploymentSpec.DiskSizeGB = intptr(10)

			replicationSpecsCheckFunc := func(c *mongodbatlas.Cluster) {
				deployment, err := createdDeployment.Spec.Deployment()
				Expect(err).NotTo(HaveOccurred())
				expectedReplicationSpecs := deployment.ReplicationSpecs

				// The ID field is added by Atlas - we don't have it in our specs
				Expect(c.ReplicationSpecs[0].ID).NotTo(BeNil())
				c.ReplicationSpecs[0].ID = ""
				// Apart from 'ID' all other fields are equal to the ones sent by the Operator
				Expect(c.ReplicationSpecs).To(Equal(expectedReplicationSpecs))
			}

			By("Creating the Deployment", func() {
				Expect(k8sClient.Create(context.Background(), createdDeployment)).To(Succeed())

				Eventually(testutil.WaitFor(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc()),
					DeploymentUpdateTimeout, interval).Should(BeTrue())

				doRegularDeploymentStatusChecks()

				checkAtlasState(replicationSpecsCheckFunc)
			})

			By("Updating the deployment (multiple operations)", func() {
				delete(createdDeployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionsConfig, "US_WEST_2")
				createdDeployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionsConfig["US_WEST_1"] = mdbv1.RegionsConfig{AnalyticsNodes: int64ptr(0), ElectableNodes: int64ptr(2), Priority: int64ptr(6), ReadOnlyNodes: int64ptr(0)}
				config := createdDeployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionsConfig["US_EAST_1"]
				// Note, that Atlas has strict requirements to priorities - they must start with 7 and be in descending order over the regions
				config.Priority = int64ptr(7)
				createdDeployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionsConfig["US_EAST_1"] = config

				createdDeployment.Spec.DeploymentSpec.ProviderSettings.AutoScaling.Compute.MaxInstanceSize = "M30"

				performUpdate(DeploymentUpdateTimeout)

				Eventually(testutil.WaitFor(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentUpdatingFunc()),
					DeploymentUpdateTimeout, interval).Should(BeTrue())

				doRegularDeploymentStatusChecks()

				checkAtlasState(replicationSpecsCheckFunc)
			})

			By("Disable deployment and disk AutoScaling", func() {
				createdDeployment.Spec.DeploymentSpec.AutoScaling = &mdbv1.AutoScalingSpec{
					Compute: &mdbv1.ComputeSpec{
						Enabled:          boolptr(false),
						ScaleDownEnabled: boolptr(false),
					},
					DiskGBEnabled: boolptr(false),
				}

				performUpdate(DeploymentUpdateTimeout)

				Eventually(testutil.WaitFor(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentUpdatingFunc()),
					DeploymentUpdateTimeout, interval).Should(BeTrue())

				doRegularDeploymentStatusChecks()

				checkAtlasState(func(c *mongodbatlas.Cluster) {
					deployment, err := createdDeployment.Spec.Deployment()
					Expect(err).NotTo(HaveOccurred())

					Expect(c.AutoScaling.Compute).To(Equal(deployment.AutoScaling.Compute))
				})
			})
		})
	})

	Describe("Create/Update the cluster", func() {
		It("Should fail, then be fixed (GCP)", func() {
			createdDeployment = mdbv1.DefaultGCPDeployment(namespace.Name, createdProject.Name).WithAtlasName("")

			By(fmt.Sprintf("Trying to create the Deployment %s with invalid parameters", kube.ObjectKeyFromObject(createdDeployment)), func() {
				err := k8sClient.Create(context.Background(), createdDeployment)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(MatchRegexp("is invalid: spec.deploymentSpec.name"))
			})

			By("Creating the fixed deployment", func() {
				createdDeployment.Spec.DeploymentSpec.Name = "fixed-deployment"

				Expect(k8sClient.Create(context.Background(), createdDeployment)).To(Succeed())

				Eventually(testutil.WaitFor(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType)),
					20*time.Minute, interval).Should(BeTrue())

				doRegularDeploymentStatusChecks()
				checkAtlasState()
			})
		})

		It("Should Succeed (AWS)", func() {
			createdDeployment = mdbv1.DefaultAWSDeployment(namespace.Name, createdProject.Name)

			By(fmt.Sprintf("Creating the Deployment %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				Expect(k8sClient.Create(context.Background(), createdDeployment)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc()),
					DeploymentUpdateTimeout, interval).Should(BeTrue())

				doRegularDeploymentStatusChecks()
				checkAtlasState()
			})

			By("Updating the Deployment labels", func() {
				createdDeployment.Spec.DeploymentSpec.Labels = []common.LabelSpec{{Key: "int-test", Value: "true"}}
				performUpdate(20 * time.Minute)
				doRegularDeploymentStatusChecks()
				checkAtlasState()
			})

			By("Updating the Deployment backups settings", func() {
				createdDeployment.Spec.DeploymentSpec.ProviderBackupEnabled = boolptr(true)
				performUpdate(20 * time.Minute)
				doRegularDeploymentStatusChecks()
				checkAtlasState(func(c *mongodbatlas.Cluster) {
					Expect(c.ProviderBackupEnabled).To(Equal(createdDeployment.Spec.DeploymentSpec.ProviderBackupEnabled))
				})
			})

			By("Decreasing the Deployment disk size", func() {
				createdDeployment.Spec.DeploymentSpec.DiskSizeGB = intptr(10)
				performUpdate(20 * time.Minute)
				doRegularDeploymentStatusChecks()
				checkAtlasState(func(c *mongodbatlas.Cluster) {
					Expect(*c.DiskSizeGB).To(BeEquivalentTo(*createdDeployment.Spec.DeploymentSpec.DiskSizeGB))

					// check whether https://github.com/mongodb/go-client-mongodb-atlas/issues/140 is fixed
					Expect(c.DiskSizeGB).To(BeAssignableToTypeOf(float64ptr(0)), "DiskSizeGB is no longer a *float64, please check the spec!")
				})
			})

			By("Pausing the deployment", func() {
				createdDeployment.Spec.DeploymentSpec.Paused = boolptr(true)
				performUpdate(20 * time.Minute)
				doRegularDeploymentStatusChecks()
				checkAtlasState(func(c *mongodbatlas.Cluster) {
					Expect(c.Paused).To(Equal(createdDeployment.Spec.DeploymentSpec.Paused))
				})
			})

			By("Updating the Deployment configuration while paused (should fail)", func() {
				createdDeployment.Spec.DeploymentSpec.ProviderBackupEnabled = boolptr(false)

				Expect(k8sClient.Update(context.Background(), createdDeployment)).To(Succeed())
				Eventually(
					testutil.WaitFor(
						k8sClient,
						createdDeployment,
						status.
							FalseCondition(status.DeploymentReadyType).
							WithReason(string(workflow.DeploymentNotUpdatedInAtlas)).
							WithMessageRegexp("CANNOT_UPDATE_PAUSED_CLUSTER"),
					),
					60,
					interval,
				).Should(BeTrue())

				lastGeneration++
			})

			By("Unpausing the deployment", func() {
				createdDeployment.Spec.DeploymentSpec.Paused = boolptr(false)
				performUpdate(20 * time.Minute)
				doRegularDeploymentStatusChecks()
				checkAtlasState(func(c *mongodbatlas.Cluster) {
					Expect(c.Paused).To(Equal(createdDeployment.Spec.DeploymentSpec.Paused))
				})
			})

			By("Checking that modifications were applied after unpausing", func() {
				doRegularDeploymentStatusChecks()
				checkAtlasState(func(c *mongodbatlas.Cluster) {
					Expect(c.ProviderBackupEnabled).To(Equal(createdDeployment.Spec.DeploymentSpec.ProviderBackupEnabled))
				})
			})

			By("Setting incorrect instance size (should fail)", func() {
				oldSizeName := createdDeployment.Spec.DeploymentSpec.ProviderSettings.InstanceSizeName
				createdDeployment.Spec.DeploymentSpec.ProviderSettings.InstanceSizeName = "M42"

				Expect(k8sClient.Update(context.Background(), createdDeployment)).To(Succeed())
				Eventually(
					testutil.WaitFor(
						k8sClient,
						createdDeployment,
						status.
							FalseCondition(status.DeploymentReadyType).
							WithReason(string(workflow.DeploymentNotUpdatedInAtlas)).
							WithMessageRegexp("INVALID_ENUM_VALUE"),
					),
					60,
					interval,
				).Should(BeTrue())

				lastGeneration++

				By("Fixing the Deployment", func() {
					createdDeployment.Spec.DeploymentSpec.ProviderSettings.InstanceSizeName = oldSizeName
					performUpdate(20 * time.Minute)
					doRegularDeploymentStatusChecks()
					checkAtlasState()
				})
			})
		})
	})

	Describe("Create DBUser before deployment & check secrets", func() {
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
				WithScope(mdbv1.DeploymentScopeType, "fake-deployment")
			By(fmt.Sprintf("Creating the Database User %s", kube.ObjectKeyFromObject(createdDBUserFakeScope)), func() {
				Expect(k8sClient.Create(context.Background(), createdDBUserFakeScope)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdDBUserFakeScope, status.FalseCondition(status.DatabaseUserReadyType).WithReason(string(workflow.DatabaseUserInvalidSpec))),
					20, intervalShort).Should(BeTrue())
			})
			checkNumberOfConnectionSecrets(k8sClient, *createdProject, 0)

			createdDeployment = mdbv1.DefaultAWSDeployment(namespace.Name, createdProject.Name)
			By(fmt.Sprintf("Creating the Deployment %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				Expect(k8sClient.Create(context.Background(), createdDeployment)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc()),
					DeploymentUpdateTimeout, interval).Should(BeTrue())

				doRegularDeploymentStatusChecks()
				checkAtlasState()
			})

			By("Checking connection Secrets", func() {
				Expect(tryConnect(createdProject.ID(), *createdDeployment, *createdDBUser)).To(Succeed())
				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 1)
				validateSecret(k8sClient, *createdProject, *createdDeployment, *createdDBUser)
			})
		})
	})

	Describe("Create deployment, user, delete deployment and check secrets are removed", func() {
		It("Should Succeed", func() {
			createdDeployment = mdbv1.DefaultAWSDeployment(namespace.Name, createdProject.Name)
			By(fmt.Sprintf("Creating the Deployment %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				Expect(k8sClient.Create(context.Background(), createdDeployment)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc()),
					DeploymentUpdateTimeout, interval).Should(BeTrue())

				doRegularDeploymentStatusChecks()
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

			By("Removing Atlas Deployment "+createdDeployment.Name, func() {
				Expect(k8sClient.Delete(context.Background(), createdDeployment)).To(Succeed())
				Eventually(checkAtlasDeploymentRemoved(createdProject.Status.ID, createdDeployment.Spec.DeploymentSpec.Name), 600, interval).Should(BeTrue())
			})

			By("Checking that Secrets got removed", func() {
				secretNames := []string{kube.NormalizeIdentifier(fmt.Sprintf("%s-%s-%s", createdProject.Spec.Name, createdDeployment.Spec.DeploymentSpec.Name, createdDBUser.Spec.Username))}
				createdDeployment = nil // prevent cleanup from failing due to deployment already deleted
				Eventually(checkSecretsDontExist(namespace.Name, secretNames), 50, interval).Should(BeTrue())
				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 0)
			})
		})
	})

	Describe("Deleting the deployment (not cleaning Atlas)", func() {
		It("Should Succeed", func() {
			By(`Creating the deployment with retention policy "keep" first`, func() {
				createdDeployment = mdbv1.DefaultAWSDeployment(namespace.Name, createdProject.Name).Lightweight()
				createdDeployment.ObjectMeta.Annotations = map[string]string{customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep}
				manualDeletion = true // We need to remove the deployment in Atlas manually to let project get removed
				Expect(k8sClient.Create(context.Background(), createdDeployment)).ToNot(HaveOccurred())

				Eventually(
					func(g Gomega) {
						success := testutil.WaitFor(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFuncGContext(g))()
						g.Expect(success).To(BeTrue())
					}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(Succeed())
			})
			By("Deleting the deployment - stays in Atlas", func() {
				Expect(k8sClient.Delete(context.Background(), createdDeployment)).To(Succeed())
				time.Sleep(5 * time.Minute)
				Expect(checkAtlasDeploymentRemoved(createdProject.Status.ID, createdDeployment.Spec.DeploymentSpec.Name)()).Should(BeFalse())
				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 0)
			})
		})
	})

	Describe("Setting the deployment skip annotation should skip reconciliations.", func() {
		It("Should Succeed", func() {

			By(`Creating the deployment with reconciliation policy "skip" first`, func() {
				createdDeployment = mdbv1.DefaultAWSDeployment(namespace.Name, createdProject.Name).Lightweight()
				Expect(k8sClient.Create(context.Background(), createdDeployment)).ToNot(HaveOccurred())

				Eventually(
					func(g Gomega) {
						success := testutil.WaitFor(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFuncGContext(g))()
						g.Expect(success).To(BeTrue())
					}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(Succeed())

				createdDeployment.ObjectMeta.Annotations = map[string]string{customresource.ReconciliationPolicyAnnotation: customresource.ReconciliationPolicySkip}
				createdDeployment.Spec.DeploymentSpec.Labels = append(createdDeployment.Spec.DeploymentSpec.Labels, common.LabelSpec{
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

				Expect(k8sClient.Update(context.Background(), createdDeployment)).ToNot(HaveOccurred())
				Eventually(testutil.WaitForAtlasDeploymentStateToNotBeReached(ctx, atlasClient, createdProject.Name, createdDeployment.GetDeploymentName(), containsLabel))
			})
		})
	})

	Describe("Create the advanced deployment & change the InstanceSize", func() {
		It("Should Succeed", func() {
			createdDeployment = mdbv1.DefaultAwsAdvancedDeployment(namespace.Name, createdProject.Name)

			By(fmt.Sprintf("Creating the Advanced Deployment %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				Expect(k8sClient.Create(context.Background(), createdDeployment)).ToNot(HaveOccurred())

				Eventually(
					func(g Gomega) {
						success := testutil.WaitFor(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFuncGContext(g))()
						g.Expect(success).To(BeTrue())
					}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(Succeed())

				doAdvancedDeploymentStatusChecks()
				checkAdvancedAtlasState()
			})

			By(fmt.Sprintf("Updating the Advanced Deployment %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				createdDeployment.Spec.AdvancedDeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].ElectableSpecs.InstanceSize = "M10"
				Expect(k8sClient.Update(context.Background(), createdDeployment)).ToNot(HaveOccurred())

				Eventually(
					func(g Gomega) {
						success := testutil.WaitFor(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFuncGContext(g))()
						g.Expect(success).To(BeTrue())
					}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(Succeed())

				doAdvancedDeploymentStatusChecks()
				checkAdvancedAtlasState()
			})
		})
	})

	Describe("Set advanced deployment options", func() {
		It("Should Succeed", func() {
			createdDeployment = mdbv1.DefaultAWSDeployment(namespace.Name, createdProject.Name)
			createdDeployment.Spec.ProcessArgs = &mdbv1.ProcessArgs{
				JavascriptEnabled:  boolptr(true),
				DefaultReadConcern: "available",
			}

			By(fmt.Sprintf("Creating the Deployment with Advanced Options %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				Expect(k8sClient.Create(context.Background(), createdDeployment)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc()),
					30*time.Minute, interval).Should(BeTrue())

				doRegularDeploymentStatusChecks()
				checkAdvancedDeploymentOptions(createdDeployment.Spec.ProcessArgs)
			})

			By("Updating Advanced Deployment Options", func() {
				createdDeployment.Spec.ProcessArgs.JavascriptEnabled = boolptr(false)
				performUpdate(40 * time.Minute)
				doRegularDeploymentStatusChecks()
				checkAdvancedDeploymentOptions(createdDeployment.Spec.ProcessArgs)
			})
		})
	})

	Describe("Create serverless instance", func() {
		It("Should Succeed", func() {
			createdDeployment = mdbv1.NewDefaultAWSServerlessInstance(namespace.Name, createdProject.Name)

			By(fmt.Sprintf("Creating the Serverless Instance %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				Expect(k8sClient.Create(context.Background(), createdDeployment)).ToNot(HaveOccurred())

				Eventually(
					func(g Gomega) {
						success := testutil.WaitFor(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFuncGContext(g))()
						g.Expect(success).To(BeTrue())
					}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(Succeed())

				doServerlessDeploymentStatusChecks()
			})
		})
	})

	Describe("Create default deployment with backups enabled", func() {
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

			createdDeployment = mdbv1.DefaultAWSDeployment(namespace.Name, createdProject.Name).WithBackupScheduleRef(common.ResourceRefNamespaced{
				Name:      backupScheduleDefault.Name,
				Namespace: backupScheduleDefault.Namespace,
			})

			By(fmt.Sprintf("Creating deployment with backups enabled: %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				Expect(k8sClient.Create(context.Background(), createdDeployment)).NotTo(HaveOccurred())

				// Do not use Gomega function here like func(g Gomega) as it seems to hang when tests run in parallel
				Eventually(
					func() error {
						deployment, _, err := atlasClient.Clusters.Get(context.Background(), createdProject.ID(), createdDeployment.Spec.DeploymentSpec.Name)
						if err != nil {
							return err
						}
						if deployment.StateName != "IDLE" {
							return errors.New("deployment is not IDLE yet")
						}
						time.Sleep(10 * time.Second)
						return nil
					}).WithTimeout(30 * time.Minute).WithPolling(5 * time.Second).Should(Not(HaveOccurred()))

				Eventually(func() error {
					actualPolicy, _, err := atlasClient.CloudProviderSnapshotBackupPolicies.Get(context.Background(), createdProject.ID(), createdDeployment.Spec.DeploymentSpec.Name)
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

func validateDeploymentCreatingFunc() func(a mdbv1.AtlasCustomResource) {
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
				status.FalseCondition(status.DeploymentReadyType).WithReason(string(workflow.DeploymentCreating)).WithMessageRegexp("deployment is provisioning"),
				status.FalseCondition(status.ReadyType),
				status.TrueCondition(status.ValidationSucceeded),
			)
			Expect(c.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
		} else {
			// Otherwise there could have been some exception in Atlas on creation - let's check the conditions
			condition, ok := testutil.FindConditionByType(c.Status.Conditions, status.DeploymentReadyType)
			Expect(ok).To(BeFalse(), fmt.Sprintf("Unexpected condition: %v", condition))
		}
	}
}

func validateDeploymentCreatingFuncGContext(g Gomega) func(a mdbv1.AtlasCustomResource) {
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
				status.FalseCondition(status.DeploymentReadyType).WithReason(string(workflow.DeploymentCreating)).WithMessageRegexp("deployment is provisioning"),
				status.FalseCondition(status.ReadyType),
				status.TrueCondition(status.ValidationSucceeded),
			)
			g.Expect(c.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
		} else {
			// Otherwise there could have been some exception in Atlas on creation - let's check the conditions
			condition, ok := testutil.FindConditionByType(c.Status.Conditions, status.DeploymentReadyType)
			g.Expect(ok).To(BeFalse(), fmt.Sprintf("Unexpected condition: %v", condition))
		}
	}
}

func validateDeploymentUpdatingFunc() func(a mdbv1.AtlasCustomResource) {
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
				status.FalseCondition(status.DeploymentReadyType).WithReason(string(workflow.DeploymentUpdating)).WithMessageRegexp("deployment is updating"),
				status.FalseCondition(status.ReadyType),
				status.TrueCondition(status.ValidationSucceeded),
			)
			Expect(c.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
		}
	}
}

// checkAtlasDeploymentRemoved returns true if the Atlas Deployment is removed from Atlas. Note the behavior: the deployment
// is removed from Atlas as soon as the DELETE API call has been made. This is different from the case when the
// deployment is terminated from UI (in this case GET request succeeds while the deployment is being terminated)
func checkAtlasDeploymentRemoved(projectID string, deploymentName string) func() bool {
	return func() bool {
		_, r, err := atlasClient.Clusters.Get(context.Background(), projectID, deploymentName)
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
