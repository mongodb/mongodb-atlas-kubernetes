package int

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/compat"

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
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
)

const (
	// Set this to true if you are debugging deployment creation.
	// This may not help much if there was the update though...
	DeploymentDevMode       = false
	DeploymentUpdateTimeout = 40 * time.Minute
	ConnectionSecretName    = "my-atlas-key"
	PrivateAPIKey           = "privateApiKey"
	OrgID                   = "orgId"
	PublicAPIKey            = "publicApiKey"
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
				Name:      ConnectionSecretName,
				Namespace: namespace.Name,
				Labels: map[string]string{
					connectionsecret.TypeLabelKey: connectionsecret.CredLabelVal,
				},
			},
			StringData: map[string]string{OrgID: connection.OrgID, PublicAPIKey: connection.PublicKey, PrivateAPIKey: connection.PrivateKey},
		}
		By(fmt.Sprintf("Creating the Secret %s", kube.ObjectKeyFromObject(&connectionSecret)))
		Expect(k8sClient.Create(context.Background(), &connectionSecret)).To(Succeed())

		createdProject = mdbv1.DefaultProject(namespace.Name, connectionSecret.Name).WithIPAccessList(project.NewIPAccessList().WithCIDR("0.0.0.0/0"))
		if DeploymentDevMode {
			// While developing tests we need to reuse the same project
			createdProject.Spec.Name = "dev-test atlas-project"
		}
		By("Creating the project " + createdProject.Name)
		Expect(k8sClient.Create(context.Background(), createdProject)).To(Succeed())
		Eventually(func() bool {
			return testutil.CheckCondition(k8sClient, createdProject, status.TrueCondition(status.ReadyType))
		}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())
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
				deploymentName := createdDeployment.GetDeploymentName()
				if customresource.ResourceShouldBeLeftInAtlas(createdDeployment) || customresource.ReconciliationShouldBeSkipped(createdDeployment) {
					By("Removing Atlas Deployment " + createdDeployment.Name + " from Atlas manually")
					Expect(deleteAtlasDeployment(createdProject.Status.ID, deploymentName)).To(Succeed())
				}
				Eventually(checkAtlasDeploymentRemoved(createdProject.Status.ID, deploymentName), 600, interval).Should(BeTrue())
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
			Expect(createdDeployment.Status.Conditions).To(HaveLen(4))
			Expect(createdDeployment.Status.Conditions).To(ConsistOf(testutil.MatchConditions(
				status.TrueCondition(status.DeploymentReadyType),
				status.TrueCondition(status.ReadyType),
				status.TrueCondition(status.ValidationSucceeded),
				status.TrueCondition(status.ResourceVersionStatus),
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
			Expect(createdDeployment.Status.Conditions).To(HaveLen(4))
			Expect(createdDeployment.Status.Conditions).To(ConsistOf(testutil.MatchConditions(
				status.TrueCondition(status.DeploymentReadyType),
				status.TrueCondition(status.ReadyType),
				status.TrueCondition(status.ValidationSucceeded),
				status.TrueCondition(status.ResourceVersionStatus),
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
			Expect(createdDeployment.Status.Conditions).To(HaveLen(4))
			Expect(createdDeployment.Status.Conditions).To(ConsistOf(testutil.MatchConditions(
				status.TrueCondition(status.DeploymentReadyType),
				status.TrueCondition(status.ReadyType),
				status.TrueCondition(status.ValidationSucceeded),
				status.TrueCondition(status.ResourceVersionStatus),
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
			specDeployment := *createdDeployment.Spec.AdvancedDeploymentSpec
			atlasDeploymentAsAtlas, _, err := atlasClient.AdvancedClusters.Get(context.Background(), createdProject.Status.ID, createdDeployment.Spec.AdvancedDeploymentSpec.Name)
			Expect(err).ToNot(HaveOccurred())

			mergedDeployment, atlasDeployment, err := atlasdeployment.MergedAdvancedDeployment(*atlasDeploymentAsAtlas, specDeployment)
			Expect(err).ToNot(HaveOccurred())

			Expect(atlasdeployment.AdvancedDeploymentsEqual(zap.S(), mergedDeployment, atlasDeployment)).To(BeTrue())

			for _, check := range additionalChecks {
				check(atlasDeploymentAsAtlas)
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

	performUpdate := func(timeout time.Duration) {
		Expect(k8sClient.Update(context.Background(), createdDeployment)).To(Succeed())

		Eventually(func(g Gomega) bool {
			return testutil.CheckCondition(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentUpdatingFunc(g))
		}).WithTimeout(timeout).WithPolling(interval).Should(BeTrue())

		lastGeneration++
	}

	Describe("Deployment CR should exist if it is tried to delete and the token is not valid", func() {
		It("Should Succeed", func() {
			expectedDeployment := mdbv1.DefaultAWSDeployment(namespace.Name, createdProject.Name)

			By(fmt.Sprintf("Creating the Deployment %s", kube.ObjectKeyFromObject(expectedDeployment)), func() {
				createdDeployment.ObjectMeta = expectedDeployment.ObjectMeta
				Expect(k8sClient.Create(context.Background(), expectedDeployment)).ToNot(HaveOccurred())

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

				doRegularDeploymentStatusChecks()
				checkAtlasState()
			})

			By("Filling token secret with invalid data", func() {
				secret := &corev1.Secret{}
				Expect(k8sClient.Get(context.Background(), kube.ObjectKeyFromObject(&connectionSecret), secret)).To(Succeed())
				secret.StringData = map[string]string{
					OrgID: "fake", PrivateAPIKey: "fake", PublicAPIKey: "fake",
				}
				Expect(k8sClient.Update(context.Background(), secret)).To(Succeed())
			})

			By("Deleting the Deployment", func() {
				Expect(k8sClient.Delete(context.Background(), createdDeployment)).To(Succeed())
			})

			By("Checking that the Deployment still exists", func() {
				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdDeployment, status.FalseCondition(status.DeploymentReadyType).
						WithMessageRegexp(strconv.Itoa(http.StatusUnauthorized)))
				}).WithTimeout(5 * time.Minute).WithPolling(20 * time.Second).Should(BeTrue())
			})

			By("Fix the token secret", func() {
				secret := &corev1.Secret{}
				Expect(k8sClient.Get(context.Background(), types.NamespacedName{Namespace: namespace.Name, Name: ConnectionSecretName}, secret)).Should(Succeed())
				secret.StringData = map[string]string{
					OrgID: connection.OrgID, PublicAPIKey: connection.PublicKey, PrivateAPIKey: connection.PrivateKey,
				}
				Expect(k8sClient.Update(context.Background(), secret)).To(Succeed())
			})

			By("Checking that the Deployment is deleted", func() {
				Eventually(checkAtlasDeploymentRemoved(createdProject.Status.ID, createdDeployment.GetDeploymentName())).
					WithTimeout(600 * time.Second).WithPolling(interval).Should(BeTrue())
			})

			// it's needed to skip deployment deletion in AfterEach
			createdDeployment = nil
		})
	})

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

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

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

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

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

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

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

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

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

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentUpdatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

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

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentUpdatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())
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

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())
				doRegularDeploymentStatusChecks()
				checkAtlasState()
			})
		})

		It("Should Success (AWS) with enabled autoscaling", func() {
			createdDeployment = mdbv1.DefaultAWSDeployment(namespace.Name, createdProject.Name)
			createdDeployment.Spec.DeploymentSpec.DiskSizeGB = intptr(20)
			createdDeployment.Spec.DeploymentSpec.AutoScaling = &mdbv1.AutoScalingSpec{
				DiskGBEnabled: boolptr(true),
			}

			By(fmt.Sprintf("Creating the Deployment %s with autoscaling", kube.ObjectKeyFromObject(createdDeployment)), func() {
				Expect(k8sClient.Create(context.Background(), createdDeployment)).ToNot(HaveOccurred())

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

				doRegularDeploymentStatusChecks()
				checkAtlasState()
			})

			By("Decreasing the Deployment disk size should not take effect", func() {
				prevDiskSize := *createdDeployment.Spec.DeploymentSpec.DiskSizeGB
				createdDeployment.Spec.DeploymentSpec.DiskSizeGB = intptr(14)
				performUpdate(30 * time.Minute)
				doRegularDeploymentStatusChecks()
				checkAtlasState(func(c *mongodbatlas.Cluster) {
					Expect(*c.DiskSizeGB).To(BeEquivalentTo(prevDiskSize))

					// check whether https://github.com/mongodb/go-client-mongodb-atlas/issues/140 is fixed
					Expect(c.DiskSizeGB).To(BeAssignableToTypeOf(float64ptr(0)), "DiskSizeGB is no longer a *float64, please check the spec!")
				})
			})
		})

		It("Should Succeed (AWS)", func() {
			createdDeployment = mdbv1.DefaultAWSDeployment(namespace.Name, createdProject.Name)
			createdDeployment.Spec.DeploymentSpec.DiskSizeGB = intptr(20)
			createdDeployment = createdDeployment.WithAutoscalingDisabled()

			By(fmt.Sprintf("Creating the Deployment %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				Expect(k8sClient.Create(context.Background(), createdDeployment)).ToNot(HaveOccurred())

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

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
				createdDeployment.Spec.DeploymentSpec.DiskSizeGB = intptr(15)
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
				Eventually(func() bool {
					return testutil.CheckCondition(
						k8sClient,
						createdDeployment,
						status.
							FalseCondition(status.DeploymentReadyType).
							WithReason(string(workflow.DeploymentNotUpdatedInAtlas)).
							WithMessageRegexp("CANNOT_UPDATE_PAUSED_CLUSTER"),
					)
				}).
					WithTimeout(DeploymentUpdateTimeout).
					WithPolling(interval).
					Should(BeTrue())

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
				Eventually(func() bool {
					return testutil.CheckCondition(
						k8sClient,
						createdDeployment,
						status.
							FalseCondition(status.DeploymentReadyType).
							WithReason(string(workflow.DeploymentNotUpdatedInAtlas)).
							WithMessageRegexp(".*INVALID_ENUM_VALUE.*"),
					)
				}).WithTimeout(DeploymentUpdateTimeout).
					WithPolling(interval).
					Should(BeTrue())

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

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

				checkUserInAtlas(createdProject.ID(), *createdDBUser)
			})

			createdDBUserFakeScope := mdbv1.DefaultDBUser(namespace.Name, "test-db-user-fake-scope", createdProject.Name).
				WithPasswordSecret(UserPasswordSecret).
				WithScope(mdbv1.DeploymentScopeType, "fake-deployment")
			By(fmt.Sprintf("Creating the Database User %s", kube.ObjectKeyFromObject(createdDBUserFakeScope)), func() {
				Expect(k8sClient.Create(context.Background(), createdDBUserFakeScope)).ToNot(HaveOccurred())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdDBUserFakeScope, status.FalseCondition(status.DatabaseUserReadyType).WithReason(string(workflow.DatabaseUserInvalidSpec)))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())
			})
			checkNumberOfConnectionSecrets(k8sClient, *createdProject, 0)

			createdDeployment = mdbv1.DefaultAWSDeployment(namespace.Name, createdProject.Name)
			By(fmt.Sprintf("Creating the Deployment %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				Expect(k8sClient.Create(context.Background(), createdDeployment)).ToNot(HaveOccurred())

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

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

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

				doRegularDeploymentStatusChecks()
				checkAtlasState()
			})

			passwordSecret := buildPasswordSecret(UserPasswordSecret, DBUserPassword)
			Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())

			createdDBUser := mdbv1.DefaultDBUser(namespace.Name, "test-db-user", createdProject.Name).WithPasswordSecret(UserPasswordSecret)
			By(fmt.Sprintf("Creating the Database User %s", kube.ObjectKeyFromObject(createdDBUser)), func() {
				Expect(k8sClient.Create(context.Background(), createdDBUser)).ToNot(HaveOccurred())
				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())
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

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())
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

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

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

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

				doAdvancedDeploymentStatusChecks()
				checkAdvancedAtlasState()

				lastGeneration++
			})

			By(fmt.Sprintf("Updating the InstanceSize of Advanced Deployment %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				createdDeployment.Spec.AdvancedDeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].ElectableSpecs.InstanceSize = "M20"
				Expect(k8sClient.Update(context.Background(), createdDeployment)).ToNot(HaveOccurred())

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentUpdatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

				doAdvancedDeploymentStatusChecks()
				checkAdvancedAtlasState()

				lastGeneration++
			})

			By(fmt.Sprintf("Enable AutoScaling for the Advanced Deployment %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				regionConfig := createdDeployment.Spec.AdvancedDeploymentSpec.ReplicationSpecs[0].RegionConfigs[0]
				regionConfig.ElectableSpecs.InstanceSize = "M10"
				regionConfig.AutoScaling = &mdbv1.AdvancedAutoScalingSpec{
					Compute: &mdbv1.ComputeSpec{
						Enabled:          toptr.MakePtr(true),
						MaxInstanceSize:  "M30",
						MinInstanceSize:  "M10",
						ScaleDownEnabled: toptr.MakePtr(true),
					},
					DiskGB: &mdbv1.DiskGB{
						Enabled: toptr.MakePtr(true),
					},
				}
				Expect(k8sClient.Update(context.Background(), createdDeployment)).ToNot(HaveOccurred())

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentUpdatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

				doAdvancedDeploymentStatusChecks()
				checkAdvancedAtlasState()

				lastGeneration++
			})

			By(fmt.Sprintf("Update Instance Size Margins with AutoScaling for Deployemnt %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				regionConfig := createdDeployment.Spec.AdvancedDeploymentSpec.ReplicationSpecs[0].RegionConfigs[0]
				regionConfig.AutoScaling.Compute.MinInstanceSize = "M20"
				regionConfig.ElectableSpecs.InstanceSize = "M20"
				Expect(k8sClient.Update(context.Background(), createdDeployment)).ToNot(HaveOccurred())

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentUpdatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

				doAdvancedDeploymentStatusChecks()
				checkAdvancedAtlasState()
			})
		})
	})

	Describe("Create the advanced deployment with enabled autoscaling", func() {
		It("Should Succeed", func() {
			createdDeployment = mdbv1.DefaultAwsAdvancedDeployment(namespace.Name, createdProject.Name)

			createdDeployment.Spec.AdvancedDeploymentSpec.ReplicationSpecs = []*mdbv1.AdvancedReplicationSpec{
				{
					NumShards: 1,
					ZoneName:  "US_EAST_1",
					RegionConfigs: []*mdbv1.AdvancedRegionConfig{
						{
							AnalyticsSpecs: &mdbv1.Specs{
								DiskIOPS:      nil,
								EbsVolumeType: "",
								InstanceSize:  "M10",
								NodeCount:     intptr(1),
							},
							ElectableSpecs: &mdbv1.Specs{
								DiskIOPS:      nil,
								EbsVolumeType: "",
								InstanceSize:  "M10",
								NodeCount:     intptr(3),
							},
							ReadOnlySpecs: &mdbv1.Specs{
								DiskIOPS:      nil,
								EbsVolumeType: "",
								InstanceSize:  "M10",
								NodeCount:     intptr(1),
							},
							AutoScaling: &mdbv1.AdvancedAutoScalingSpec{
								DiskGB: &mdbv1.DiskGB{
									Enabled: toptr.MakePtr(true),
								},
								Compute: &mdbv1.ComputeSpec{
									Enabled:          toptr.MakePtr(true),
									ScaleDownEnabled: toptr.MakePtr(true),
									MinInstanceSize:  "M10",
									MaxInstanceSize:  "M40",
								},
							},
							BackingProviderName: "AWS",
							Priority:            intptr(7),
							ProviderName:        "AWS",
							RegionName:          "US_EAST_1",
						},
					},
				},
			}

			By(fmt.Sprintf("Creating the Advanced Deployment %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				Expect(k8sClient.Create(context.Background(), createdDeployment)).ToNot(HaveOccurred())

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

				doAdvancedDeploymentStatusChecks()
				checkAdvancedAtlasState()

				lastGeneration++
			})

			By(fmt.Sprintf("Update autoscaling configuration with wrong values it should fail %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				previousDeployment := mdbv1.AtlasDeployment{}
				err := compat.JSONCopy(&previousDeployment, createdDeployment)
				Expect(err).NotTo(HaveOccurred())

				createdDeployment.Spec.AdvancedDeploymentSpec.ReplicationSpecs[0].
					RegionConfigs[0].
					AutoScaling.
					Compute.
					MinInstanceSize = "S"
				Expect(k8sClient.Update(context.Background(), createdDeployment)).ToNot(HaveOccurred())

				Eventually(func() bool {
					return testutil.CheckCondition(
						k8sClient,
						createdDeployment,
						status.
							FalseCondition(status.DeploymentReadyType).
							WithReason(string(workflow.Internal)).
							WithMessageRegexp("instance size is invalid"),
					)
				}).WithTimeout(DeploymentUpdateTimeout).
					WithPolling(interval).
					Should(BeTrue())

				lastGeneration++

				By(fmt.Sprintf("Update autoscaling configuration should update InstanceSize and DiskSizeGB of Advanced deployment %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
					previousDeployment := mdbv1.AtlasDeployment{}
					err := compat.JSONCopy(&previousDeployment, createdDeployment)
					Expect(err).NotTo(HaveOccurred())

					createdDeployment.Spec.AdvancedDeploymentSpec.ReplicationSpecs[0].
						RegionConfigs[0].
						AutoScaling.
						Compute.
						MinInstanceSize = "M20"
					Expect(k8sClient.Update(context.Background(), createdDeployment)).ToNot(HaveOccurred())

					Eventually(func(g Gomega) bool {
						GinkgoWriter.Println("ProjectID", createdProject.ID(), "DeploymentName", createdDeployment.GetDeploymentName())
						current, _, err := atlasClient.AdvancedClusters.Get(context.Background(), createdProject.ID(), createdDeployment.GetDeploymentName())
						g.Expect(err).NotTo(HaveOccurred())
						g.Expect(current).NotTo(BeNil())

						g.Expect(current.ReplicationSpecs[0].RegionConfigs[0].AnalyticsSpecs.InstanceSize).To(Equal("M20"))
						g.Expect(current.ReplicationSpecs[0].RegionConfigs[0].ElectableSpecs.InstanceSize).To(Equal("M20"))
						g.Expect(current.ReplicationSpecs[0].RegionConfigs[0].ReadOnlySpecs.InstanceSize).To(Equal("M20"))
						return true
					}).WithTimeout(2 * time.Minute).WithPolling(interval).Should(BeTrue())
				})
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

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

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

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdDeployment, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

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
				Status: status.BackupPolicyStatus{},
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
					}).WithTimeout(40 * time.Minute).WithPolling(15 * time.Second).Should(Not(HaveOccurred()))

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

	Describe("Create deployment with backups enabled and snapshot distribution", func() {
		It("Should succeed", func() {
			By("Creating deployment with backups enabled", func() {
				createdDeployment = mdbv1.DefaultAwsAdvancedDeployment(namespace.Name, createdProject.Name)
				createdDeployment.Spec.AdvancedDeploymentSpec.BackupEnabled = toptr.MakePtr(true)
				Expect(k8sClient.Create(context.Background(), createdDeployment)).NotTo(HaveOccurred())

				Eventually(func(g Gomega) {
					deployment, _, err := atlasClient.AdvancedClusters.Get(context.Background(), createdProject.ID(), createdDeployment.Spec.AdvancedDeploymentSpec.Name)
					g.Expect(err).Should(BeNil())
					g.Expect(deployment.StateName).Should(Equal("IDLE"))
					g.Expect(*deployment.BackupEnabled).Should(BeTrue())
					g.Expect(len(deployment.ReplicationSpecs)).ShouldNot(Equal(0))
				}).WithTimeout(40 * time.Minute).WithPolling(15 * time.Second).Should(Not(HaveOccurred()))
			})

			By("Adding BackupSchedule with Snapshot distribution", func() {
				bScheduleName := "schedule-1"

				Expect(
					k8sClient.Get(context.Background(), types.NamespacedName{Namespace: namespace.Name, Name: "test-deployment-advanced-k8s"}, createdDeployment),
				).NotTo(HaveOccurred())

				replicaSetID := createdDeployment.Status.ReplicaSets[0].ID
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
					Status: status.BackupPolicyStatus{},
				}
				backupScheduleDefault := &mdbv1.AtlasBackupSchedule{
					ObjectMeta: metav1.ObjectMeta{
						Name:      bScheduleName,
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
						CopySettings: []mdbv1.CopySetting{
							{
								CloudProvider:     toptr.MakePtr("AWS"),
								RegionName:        toptr.MakePtr("US_WEST_1"),
								ReplicationSpecID: toptr.MakePtr(replicaSetID),
								ShouldCopyOplogs:  toptr.MakePtr(false),
								Frequencies:       []string{"MONTHLY"},
							},
						},
					},
				}
				Expect(k8sClient.Create(context.Background(), backupPolicyDefault)).NotTo(HaveOccurred())
				Expect(k8sClient.Create(context.Background(), backupScheduleDefault)).NotTo(HaveOccurred())

				createdDeployment.Spec.BackupScheduleRef = common.ResourceRefNamespaced{
					Name:      bScheduleName,
					Namespace: namespace.Name,
				}
				Expect(k8sClient.Update(context.Background(), createdDeployment)).NotTo(HaveOccurred())

				Eventually(func(g Gomega) {
					atlasCluster, _, err := atlasClient.AdvancedClusters.Get(context.Background(), createdProject.ID(), createdDeployment.Spec.AdvancedDeploymentSpec.Name)
					g.Expect(err).Should(BeNil())
					g.Expect(atlasCluster.StateName).Should(Equal("IDLE"))
					g.Expect(*atlasCluster.BackupEnabled).Should(BeTrue())

					atlasBSchedule, _, err := atlasClient.CloudProviderSnapshotBackupPolicies.Get(context.Background(), createdProject.ID(), createdDeployment.Spec.AdvancedDeploymentSpec.Name)
					g.Expect(err).Should(BeNil())
					g.Expect(len(atlasBSchedule.CopySettings)).ShouldNot(Equal(0))
					g.Expect(atlasBSchedule.CopySettings[0]).
						Should(Equal(
							mongodbatlas.CopySetting{
								CloudProvider:     toptr.MakePtr("AWS"),
								RegionName:        toptr.MakePtr("US_WEST_1"),
								ReplicationSpecID: toptr.MakePtr(replicaSetID),
								ShouldCopyOplogs:  toptr.MakePtr(false),
								Frequencies:       []string{"MONTHLY"},
							},
						))
				}).WithTimeout(15 * time.Minute).WithPolling(10 * time.Second).Should(Not(HaveOccurred()))
			})
		})
	})
})

func validateDeploymentCreatingFunc(g Gomega) func(a mdbv1.AtlasCustomResource) {
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
				status.TrueCondition(status.ResourceVersionStatus),
			)
			g.Expect(c.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
		} else {
			// Otherwise there could have been some exception in Atlas on creation - let's check the conditions
			condition, ok := testutil.FindConditionByType(c.Status.Conditions, status.DeploymentReadyType)
			g.Expect(ok).To(BeFalse(), fmt.Sprintf("Unexpected condition: %v", condition))
		}
	}
}

func validateDeploymentUpdatingFunc(g Gomega) func(a mdbv1.AtlasCustomResource) {
	isIdle := true
	return func(a mdbv1.AtlasCustomResource) {
		c := a.(*mdbv1.AtlasDeployment)
		// It's ok if the first invocations see IDLE
		if c.Status.StateName != "IDLE" {
			isIdle = false
		}
		// When the create request has been made to Atlas - we expect the following status
		if !isIdle {
			g.Expect(c.Status.StateName).To(Or(Equal("UPDATING"), Equal("REPAIRING")), fmt.Sprintf("Current conditions: %+v", c.Status.Conditions))
			expectedConditionsMatchers := testutil.MatchConditions(
				status.FalseCondition(status.DeploymentReadyType).WithReason(string(workflow.DeploymentUpdating)).WithMessageRegexp("deployment is updating"),
				status.FalseCondition(status.ReadyType),
				status.TrueCondition(status.ValidationSucceeded),
				status.TrueCondition(status.ResourceVersionStatus),
			)
			g.Expect(c.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
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

func deleteAtlasDeployment(projectID string, deploymentName string) error {
	_, err := atlasClient.Clusters.Delete(context.Background(), projectID, deploymentName)
	return err
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
