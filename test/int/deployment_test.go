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
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas-sdk/v20250312012/admin"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/compat"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasdeployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/secretservice"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/httputil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/deployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/conditions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/resources"
	akoretry "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/retry"
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

const (
	interval = PollingInterval
)

var _ = Describe("AtlasDeployment", Label("int", "AtlasDeployment", "deployment-non-backups"), func() {
	var (
		deploymentService deployment.AtlasDeploymentsService
		connectionSecret  *corev1.Secret
		createdProject    *akov2.AtlasProject
		createdDeployment *akov2.AtlasDeployment
		manualDeletion    bool
	)

	BeforeEach(func() {
		prepareControllers(false)

		deploymentService = deployment.NewAtlasDeployments(atlasClient.ClustersApi, atlasClient.GlobalClustersApi, atlasClient.FlexClustersApi, false)
		createdDeployment = &akov2.AtlasDeployment{}

		manualDeletion = false

		connectionSecret = createConnectionSecret()
		createdProject = createProject(connectionSecret)
	})

	AfterEach(func() {
		if DeploymentDevMode {
			return
		}
		if manualDeletion && createdProject != nil {
			By("Deleting the deployment in Atlas manually", func() {
				// We need to remove the deployment in Atlas to let project get removed
				_, err := atlasClient.ClustersApi.
					DeleteCluster(context.Background(), createdProject.ID(), createdDeployment.GetDeploymentName()).
					Execute()
				Expect(err).NotTo(HaveOccurred())
				Eventually(checkAtlasDeploymentRemoved(createdProject.Status.ID, createdDeployment.GetDeploymentName()), 600, interval).Should(BeTrue())
				createdDeployment = nil
			})
		}
		if createdProject != nil && createdProject.Status.ID != "" {
			if createdDeployment != nil {
				deleteDeploymentFromKubernetes(createdProject, createdDeployment)
			}

			deleteProjectFromKubernetes(createdProject)
		}
		removeControllersAndNamespace()
	})

	doDeploymentStatusChecks := func() {
		By("Checking observed Deployment state", func() {
			doDeploymentStatusChecksFor(createdProject, createdDeployment)
		})
	}

	checkAtlasState := func(additionalChecks ...func(c *admin.ClusterDescription20240805)) {
		By("Verifying Deployment state in Atlas", func() {
			atlasDeploymentAsAtlas, _, err := atlasClient.ClustersApi.
				GetCluster(context.Background(), createdProject.Status.ID, createdDeployment.GetDeploymentName()).
				Execute()
			Expect(err).ToNot(HaveOccurred())

			for _, check := range additionalChecks {
				check(atlasDeploymentAsAtlas)
			}
		})
	}

	checkAdvancedAtlasState := func() {
		By("Verifying Advanced Deployment state in Atlas", func() {
			deploymentInAtlas, err := deploymentService.GetDeployment(context.Background(), createdProject.ID(), createdDeployment)
			Expect(err).ToNot(HaveOccurred())

			deploymentInAKO := deployment.NewDeployment(createdProject.ID(), createdDeployment)
			_, hasChanges := deployment.ComputeChanges(deploymentInAKO.(*deployment.Cluster), deploymentInAtlas.(*deployment.Cluster))
			Expect(hasChanges).ShouldNot(BeTrue())
		})
	}

	checkAdvancedDeploymentOptions := func(ctx context.Context, projectID string, atlasDeployment *akov2.AtlasDeployment) {
		By("Checking that Atlas Advanced Options are equal to the Spec Options", func() {
			deploymentInAKO := deployment.NewDeployment(projectID, atlasDeployment).(*deployment.Cluster)
			deploymentInAtlas, err := deploymentService.GetDeployment(ctx, projectID, atlasDeployment)
			Expect(err).ToNot(HaveOccurred())

			cluster := deploymentInAtlas.(*deployment.Cluster)
			err = deploymentService.ClusterWithProcessArgs(ctx, cluster)
			Expect(err).ToNot(HaveOccurred())

			Expect(cluster.ProcessArgs).To(Equal(deploymentInAKO.ProcessArgs))
		})
	}

	performCreate := func(deployment *akov2.AtlasDeployment, timeout time.Duration) {
		Expect(k8sClient.Create(context.Background(), deployment)).To(Succeed())

		Eventually(func(g Gomega) bool {
			return resources.CheckCondition(k8sClient, createdDeployment, api.TrueCondition(api.ReadyType), validateDeploymentCreatingFunc(g))
		}).WithTimeout(timeout).WithPolling(interval).Should(BeTrue())
	}

	Describe("Deployment with Termination Protection should remain in Atlas after the CR is deleted", Label("focus-dedicated-termination-protection", "slow"), func() {
		It("Should succeed", func() {
			createdDeployment = akov2.DefaultAWSDeployment(namespace.Name, createdProject.Name)
			deploymentName := createdDeployment.GetDeploymentName()

			By(fmt.Sprintf("Creating the Deployment %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				createdDeployment.Spec.DeploymentSpec.TerminationProtectionEnabled = true

				performCreate(createdDeployment, 30*time.Minute)

				doDeploymentStatusChecks()
				checkAtlasState()
			})

			By("Removing deployment", func() {
				Expect(k8sClient.Delete(context.Background(), createdDeployment)).To(Succeed())
			})

			By("Verifying the deployment is still in Atlas", func() {
				Eventually(func(g Gomega) {
					ctx, cancelF := context.WithTimeout(context.Background(), 20*time.Second)
					defer cancelF()
					aCluster, _, err := atlasClient.ClustersApi.GetCluster(ctx, createdProject.ID(),
						deploymentName).Execute()
					g.Expect(err).NotTo(HaveOccurred())
					Expect(aCluster.GetName()).Should(BeEquivalentTo(deploymentName))
				}).WithTimeout(30 * time.Second).WithPolling(5 * time.Second)
			})

			By("Disabling Termination protection", func() {
				ctx, cancelF := context.WithTimeout(context.Background(), 20*time.Second)
				defer cancelF()
				aCluster, _, err := atlasClient.ClustersApi.GetCluster(ctx, createdProject.ID(),
					deploymentName).Execute()
				Expect(err).NotTo(HaveOccurred())
				aCluster.TerminationProtectionEnabled = pointer.MakePtr(false)
				aCluster.ConnectionStrings = nil
				_, _, err = atlasClient.ClustersApi.UpdateCluster(ctx, createdProject.ID(), deploymentName, aCluster).Execute()
				Expect(err).NotTo(HaveOccurred())
			})

			By("Waiting for Termination protection to be disabled", func() {
				Eventually(func(g Gomega) {
					ctx, cancelF := context.WithTimeout(context.Background(), 20*time.Second)
					defer cancelF()
					aCluster, _, err := atlasClient.ClustersApi.GetCluster(ctx, createdProject.ID(),
						deploymentName).Execute()
					g.Expect(err).NotTo(HaveOccurred())
					g.Expect(aCluster.TerminationProtectionEnabled).NotTo(BeNil())
					g.Expect(*aCluster.TerminationProtectionEnabled).To(BeFalse())
				}).WithTimeout(2 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
			})

			By("Manually deleting the cluster", func() {
				ctx, cancelF := context.WithTimeout(context.Background(), 20*time.Second)
				defer cancelF()
				_, err := atlasClient.ClustersApi.DeleteCluster(ctx, createdProject.ID(),
					deploymentName).Execute()
				Expect(err).NotTo(HaveOccurred())
				createdDeployment = nil
			})

			By("Waiting for Deployment termination", func() {
				Eventually(func(g Gomega) {
					ctx, cancelF := context.WithTimeout(context.Background(), 20*time.Second)
					defer cancelF()
					_, resp, _ := atlasClient.ClustersApi.GetCluster(ctx, createdProject.ID(),
						deploymentName).Execute()
					g.Expect(httputil.StatusCode(resp)).To(Equal(http.StatusNotFound))
				}).WithTimeout(10 * time.Minute).WithPolling(20 * time.Second).Should(Succeed())
			})
		})
	})

	Describe("Deployment CR should exist if it is tried to delete and the token is not valid", func() {
		It("Should Succeed", func(ctx context.Context) {
			expectedDeployment := akov2.DefaultAWSDeployment(namespace.Name, createdProject.Name)

			By(fmt.Sprintf("Creating the Deployment %s", kube.ObjectKeyFromObject(expectedDeployment)), func() {
				createdDeployment.ObjectMeta = expectedDeployment.ObjectMeta

				performCreate(expectedDeployment, 30*time.Minute)

				createdDeployment.Spec.DeploymentSpec = expectedDeployment.Spec.DeploymentSpec

				doDeploymentStatusChecks()
				checkAtlasState()
			})

			By("Filling token secret with invalid data", func() {
				_, err := akoretry.RetryUpdateOnConflict(ctx, k8sClient, client.ObjectKeyFromObject(connectionSecret), func(secret *corev1.Secret) {
					secret.StringData = map[string]string{
						OrgID: "fake", PrivateAPIKey: "fake", PublicAPIKey: "fake",
					}
				})
				Expect(err).To(BeNil())
			})

			By("Deleting the Deployment", func() {
				Expect(k8sClient.Delete(context.Background(), createdDeployment)).To(Succeed())
			})

			By("Checking that the Deployment still exists", func() {
				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, createdDeployment, api.FalseCondition(api.DeploymentReadyType).
						WithMessageRegexp(strconv.Itoa(http.StatusUnauthorized)))
				}).WithTimeout(5 * time.Minute).WithPolling(20 * time.Second).Should(BeTrue())
			})

			By("Fix the token secret", func() {
				_, err := akoretry.RetryUpdateOnConflict(ctx, k8sClient, types.NamespacedName{Namespace: namespace.Name, Name: ConnectionSecretName}, func(secret *corev1.Secret) {
					secret.StringData = secretData()
				})
				Expect(err).To(BeNil())
			})

			By("Checking that the Deployment is deleted", func() {
				Eventually(checkAtlasDeploymentRemoved(createdProject.Status.ID, createdDeployment.GetDeploymentName())).
					WithTimeout(600 * time.Second).WithPolling(interval).Should(BeTrue())
			})

			// it's needed to skip deployment deletion in AfterEach
			createdDeployment = nil
		})
	})

	Describe("Create deployment & increase InstanceSize", func() {
		It("Should Succeed", func(ctx context.Context) {
			expectedDeployment := akov2.DefaultAWSDeployment(namespace.Name, createdProject.Name)

			By(fmt.Sprintf("Creating the Deployment %s", kube.ObjectKeyFromObject(expectedDeployment)), func() {
				createdDeployment.ObjectMeta = expectedDeployment.ObjectMeta

				performCreate(expectedDeployment, 30*time.Minute)

				createdDeployment.Spec.DeploymentSpec = expectedDeployment.Spec.DeploymentSpec

				doDeploymentStatusChecks()
				checkAtlasState()
			})

			By("Increasing InstanceSize", func() {
				createdDeployment = performUpdate(ctx, 40*time.Minute, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
					deployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].ElectableSpecs = &akov2.Specs{
						InstanceSize: "M30",
						NodeCount:    pointer.MakePtr(3),
					}
					deployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].ReadOnlySpecs = &akov2.Specs{
						InstanceSize: "M30",
						NodeCount:    pointer.MakePtr(0),
					}
					deployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].AnalyticsSpecs = &akov2.Specs{
						InstanceSize: "M30",
						NodeCount:    pointer.MakePtr(0),
					}
				})
				doDeploymentStatusChecks()
				checkAtlasState()
			})
		})
	})

	Describe("Create deployment & change it to GEOSHARDED", Label("focus-int", "focus-geosharded", "focus-slow"), func() {
		It("Should Succeed", func(ctx context.Context) {
			expectedDeployment := akov2.DefaultAWSDeployment(namespace.Name, createdProject.Name)

			By(fmt.Sprintf("Creating the Deployment %s", kube.ObjectKeyFromObject(expectedDeployment)), func() {
				createdDeployment.ObjectMeta = expectedDeployment.ObjectMeta
				performCreate(expectedDeployment, 30*time.Minute)

				doDeploymentStatusChecks()
				checkAtlasState()
			})

			By("Change deployment to GEOSHARDED", func() {
				createdDeployment = performUpdate(ctx, 90*time.Minute, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
					deployment.Spec.DeploymentSpec.ClusterType = "GEOSHARDED"
					deployment.Spec.DeploymentSpec.ReplicationSpecs = []*akov2.AdvancedReplicationSpec{
						{
							NumShards: 1,
							ZoneName:  "Zone 1",
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									AutoScaling: &akov2.AdvancedAutoScalingSpec{
										DiskGB: &akov2.DiskGB{
											Enabled: pointer.MakePtr(false),
										},
										Compute: &akov2.ComputeSpec{
											Enabled:          pointer.MakePtr(false),
											ScaleDownEnabled: pointer.MakePtr(false),
										},
									},
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M10",
										NodeCount:    pointer.MakePtr(2),
									},
									AnalyticsSpecs: &akov2.Specs{
										InstanceSize: "M10",
										NodeCount:    pointer.MakePtr(1),
									},
									Priority:     pointer.MakePtr(7),
									ProviderName: "AWS",
									RegionName:   "US_EAST_1",
								},
								{
									AutoScaling: &akov2.AdvancedAutoScalingSpec{
										DiskGB: &akov2.DiskGB{
											Enabled: pointer.MakePtr(false),
										},
										Compute: &akov2.ComputeSpec{
											Enabled:          pointer.MakePtr(false),
											ScaleDownEnabled: pointer.MakePtr(false),
										},
									},
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M10",
										NodeCount:    pointer.MakePtr(1),
									},
									Priority:     pointer.MakePtr(6),
									ProviderName: "AWS",
									RegionName:   "US_WEST_1",
								},
							},
						},
					}
				})
				doDeploymentStatusChecks()
				checkAtlasState()
			})
		})
	})

	Describe("Create/Update the deployment (more complex scenario)", Label("focus-int", "focus-create-update-complex-deployment", "slow"), func() {
		It("Should be created", func(ctx context.Context) {
			createdDeployment = akov2.DefaultAWSDeployment(namespace.Name, createdProject.Name)
			createdDeployment.Spec.DeploymentSpec.ClusterType = string(akov2.TypeReplicaSet)
			createdDeployment.Spec.DeploymentSpec.Labels = []common.LabelSpec{{Key: "createdBy", Value: "Atlas Operator"}}
			createdDeployment.Spec.DeploymentSpec.ReplicationSpecs[0] = &akov2.AdvancedReplicationSpec{
				NumShards: 1,
				ZoneName:  "Zone 1",
				RegionConfigs: []*akov2.AdvancedRegionConfig{
					{
						AutoScaling: &akov2.AdvancedAutoScalingSpec{
							DiskGB: &akov2.DiskGB{
								Enabled: pointer.MakePtr(true),
							},
							Compute: &akov2.ComputeSpec{
								Enabled:          pointer.MakePtr(true),
								ScaleDownEnabled: pointer.MakePtr(true),
								MinInstanceSize:  "M10",
								MaxInstanceSize:  "M20",
							},
						},
						ElectableSpecs: &akov2.Specs{
							InstanceSize: "M10",
							NodeCount:    pointer.MakePtr(2),
						},
						Priority:     pointer.MakePtr(7),
						ProviderName: "AWS",
						RegionName:   "US_EAST_1",
					},
					{
						AutoScaling: &akov2.AdvancedAutoScalingSpec{
							DiskGB: &akov2.DiskGB{
								Enabled: pointer.MakePtr(true),
							},
							Compute: &akov2.ComputeSpec{
								Enabled:          pointer.MakePtr(true),
								ScaleDownEnabled: pointer.MakePtr(true),
								MinInstanceSize:  "M10",
								MaxInstanceSize:  "M20",
							},
						},
						ElectableSpecs: &akov2.Specs{
							InstanceSize: "M10",
							NodeCount:    pointer.MakePtr(1),
						},
						Priority:     pointer.MakePtr(6),
						ProviderName: "AWS",
						RegionName:   "US_WEST_2",
					},
				},
			}

			createdDeployment.Spec.DeploymentSpec.DiskSizeGB = pointer.MakePtr(10)

			replicationSpecsCheckFunc := func(c *admin.ClusterDescription20240805) {
				mergedDeployment, _, err := mergedAdvancedDeployment(*c, *createdDeployment.Spec.DeploymentSpec)
				Expect(err).ToNot(HaveOccurred())

				expectedReplicationSpecs := mergedDeployment.ReplicationSpecs

				// The ID field is added by Atlas - we don't have it in our specs
				Expect(c.GetReplicationSpecs()[0].GetId()).NotTo(BeEmpty())

				// Apart from 'ID' all other fields are equal to the ones sent by the Operator
				Expect(len(c.GetReplicationSpecs())).To(Equal(expectedReplicationSpecs[0].NumShards))
				Expect(c.GetReplicationSpecs()[0].GetZoneName()).To(Equal(expectedReplicationSpecs[0].ZoneName))

				less := func(a, b *admin.CloudRegionConfig20240805) bool { return a.GetRegionName() < b.GetRegionName() }
				Expect(cmp.Diff(c.GetReplicationSpecs()[0].RegionConfigs, expectedReplicationSpecs[0].RegionConfigs, cmpopts.SortSlices(less)))
			}

			By("Creating the Deployment", func() {
				performCreate(createdDeployment, 30*time.Minute)
				doDeploymentStatusChecks()
				checkAtlasState(replicationSpecsCheckFunc)
			})

			By("Updating the deployment (multiple operations)", func() {
				var legacySpec *akov2.AdvancedDeploymentSpec
				createdDeployment = performUpdate(ctx, DeploymentUpdateTimeout, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
					deployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs = []*akov2.AdvancedRegionConfig{
						{
							AutoScaling: &akov2.AdvancedAutoScalingSpec{
								DiskGB: &akov2.DiskGB{
									Enabled: pointer.MakePtr(true),
								},
								Compute: &akov2.ComputeSpec{
									Enabled:          pointer.MakePtr(true),
									ScaleDownEnabled: pointer.MakePtr(true),
									MinInstanceSize:  "M10",
									MaxInstanceSize:  "M30",
								},
							},
							ElectableSpecs: &akov2.Specs{
								InstanceSize: "M10",
								NodeCount:    pointer.MakePtr(2),
							},
							Priority:     pointer.MakePtr(7),
							ProviderName: "AWS",
							RegionName:   "US_EAST_1",
						},
						{
							ElectableSpecs: &akov2.Specs{
								InstanceSize: "M10",
								NodeCount:    pointer.MakePtr(1),
							},
							Priority:     pointer.MakePtr(6),
							ProviderName: "AWS",
							RegionName:   "US_WEST_1",
							AutoScaling: &akov2.AdvancedAutoScalingSpec{
								DiskGB: &akov2.DiskGB{
									Enabled: pointer.MakePtr(true),
								},
								Compute: &akov2.ComputeSpec{
									Enabled:          pointer.MakePtr(true),
									ScaleDownEnabled: pointer.MakePtr(true),
									MinInstanceSize:  "M10",
									MaxInstanceSize:  "M30",
								},
							},
						},
					}

					legacySpec = deployment.Spec.DeploymentSpec
				})

				Eventually(func(g Gomega) bool {
					return resources.CheckCondition(k8sClient, createdDeployment, api.TrueCondition(api.ReadyType), validateDeploymentUpdatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

				doDeploymentStatusChecks()

				checkAtlasState(replicationSpecsCheckFunc)

				createdDeployment.Spec.DeploymentSpec = legacySpec
			})

			By("Disable deployment and disk AutoScaling", func() {
				createdDeployment = performUpdate(ctx, DeploymentUpdateTimeout, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
					deployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].AutoScaling = &akov2.AdvancedAutoScalingSpec{
						DiskGB: &akov2.DiskGB{
							Enabled: pointer.MakePtr(false),
						},
						Compute: &akov2.ComputeSpec{
							Enabled: pointer.MakePtr(false),
						},
					}
					deployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[1].AutoScaling = &akov2.AdvancedAutoScalingSpec{
						DiskGB: &akov2.DiskGB{
							Enabled: pointer.MakePtr(false),
						},
						Compute: &akov2.ComputeSpec{
							Enabled: pointer.MakePtr(false),
						},
					}
				})

				Eventually(func(g Gomega) bool {
					return resources.CheckCondition(k8sClient, createdDeployment, api.TrueCondition(api.ReadyType), validateDeploymentUpdatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())
				doDeploymentStatusChecks()

				checkAtlasState(func(c *admin.ClusterDescription20240805) {
					d := deployment.NewDeployment(createdProject.ID(), createdDeployment)

					autoScalingInput := c.GetReplicationSpecs()[0].GetRegionConfigs()[0].GetAutoScaling()
					autoScalingSpec := d.GetCustomResource().Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].AutoScaling
					Expect(autoScalingInput.Compute.Enabled).To(Equal(autoScalingSpec.Compute.Enabled))
					Expect(autoScalingInput.Compute.GetMaxInstanceSize()).To(Equal(autoScalingSpec.Compute.MaxInstanceSize))
					Expect(autoScalingInput.Compute.GetMinInstanceSize()).To(Equal(autoScalingSpec.Compute.MinInstanceSize))
				})
			})
		})
	})

	Describe("Create/Update the cluster", func() {
		It("Should fail, then be fixed (GCP)", func() {
			createdDeployment = akov2.DefaultGCPDeployment(namespace.Name, createdProject.Name).WithAtlasName("----")

			By(fmt.Sprintf("Trying to create the Deployment %s with invalid parameters", kube.ObjectKeyFromObject(createdDeployment)), func() {
				err := k8sClient.Create(context.Background(), createdDeployment)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(MatchRegexp("is invalid: spec.deploymentSpec.name"))
			})

			By("Creating the fixed deployment", func() {
				createdDeployment.Spec.DeploymentSpec.Name = "fixed-deployment"
				performCreate(createdDeployment, 30*time.Minute)

				doDeploymentStatusChecks()
				checkAtlasState()
			})
		})

		It("Should succeed (AWS) with enabled autoscaling for Disk size", func(ctx context.Context) {
			createdDeployment = akov2.DefaultAWSDeployment(namespace.Name, createdProject.Name)

			createdDeployment.Spec.DeploymentSpec.DiskSizeGB = pointer.MakePtr(20)
			createdDeployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].AutoScaling = &akov2.AdvancedAutoScalingSpec{
				DiskGB: &akov2.DiskGB{
					Enabled: pointer.MakePtr(true),
				},
			}

			By(fmt.Sprintf("Creating the Deployment %s with autoscaling", kube.ObjectKeyFromObject(createdDeployment)), func() {
				performCreate(createdDeployment, 30*time.Minute)

				doDeploymentStatusChecks()
				checkAtlasState()
			})

			By("Decreasing the Deployment disk size should not take effect", func() {
				// prevDiskSize := *createdDeployment.Spec.DeploymentSpec.DiskSizeGB
				createdDeployment = performUpdate(ctx, 30*time.Minute, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
					deployment.Spec.DeploymentSpec.DiskSizeGB = pointer.MakePtr(14)
				})

				doDeploymentStatusChecks()
				checkAtlasState(func(c *admin.ClusterDescription20240805) {
					// Expect(*c.DiskSizeGB).To(BeEquivalentTo(prevDiskSize)) // todo: find out if this should still work for advanced clusters

					Expect(c.GetReplicationSpecs()[0].GetRegionConfigs()[0].ElectableSpecs.DiskSizeGB).To(BeAssignableToTypeOf(pointer.MakePtr[float64](0)), "DiskSizeGB is no longer a *float64, please check the spec!")
				})
			})
		})

		It("Should Succeed (AWS)", func(ctx context.Context) {
			createdDeployment = akov2.DefaultAWSDeployment(namespace.Name, createdProject.Name)
			createdDeployment.Spec.DeploymentSpec.DiskSizeGB = pointer.MakePtr(20)
			createdDeployment = createdDeployment.WithAutoscalingDisabled()

			By(fmt.Sprintf("Creating the Deployment %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				performCreate(createdDeployment, 30*time.Minute)

				doDeploymentStatusChecks()
				checkAtlasState()
			})

			By("Updating the Deployment labels", func() {
				createdDeployment = performUpdate(ctx, 30*time.Minute, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
					deployment.Spec.DeploymentSpec.Labels = []common.LabelSpec{{Key: "int-test", Value: "true"}}
				})
				doDeploymentStatusChecks()
				checkAtlasState()
			})

			By("Updating the Deployment tags", func() {
				createdDeployment = performUpdate(ctx, 30*time.Minute, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
					deployment.Spec.DeploymentSpec.Tags = []*akov2.TagSpec{{Key: "test 1", Value: "value 1"}, {Key: "test-2", Value: "value-2"}}
				})
				doDeploymentStatusChecks()
				checkAtlasState(func(c *admin.ClusterDescription20240805) {
					for i, tag := range createdDeployment.Spec.DeploymentSpec.Tags {
						Expect(c.GetTags()[i].GetKey() == tag.Key).To(BeTrue())
						Expect(c.GetTags()[i].GetValue() == tag.Value).To(BeTrue())
					}
				})
			})

			By("Updating the Deployment tags with a duplicate key and removing all tags", func() {
				_, err := akoretry.RetryUpdateOnConflict(ctx, k8sClient, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
					deployment.Spec.DeploymentSpec.Tags = []*akov2.TagSpec{{Key: "test-1", Value: "value-1"}, {Key: "test-1", Value: "value-2"}}
				})
				Expect(err).To(BeNil())
				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, createdDeployment, api.FalseCondition(api.ValidationSucceeded))
				}).WithTimeout(DeploymentUpdateTimeout).Should(BeTrue())
				createdDeployment = performUpdate(ctx, 30*time.Minute, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
					// Removing tags for next tests
					deployment.Spec.DeploymentSpec.Tags = []*akov2.TagSpec{}
				})
			})

			By("Updating the Deployment backups settings", func() {
				createdDeployment = performUpdate(ctx, 30*time.Minute, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
					// createdDeployment.Spec.DeploymentSpec.ProviderBackupEnabled = pointer.MakePtr(true)
					deployment.Spec.DeploymentSpec.BackupEnabled = pointer.MakePtr(true)
				})
				doDeploymentStatusChecks()
				checkAtlasState(func(c *admin.ClusterDescription20240805) {
					Expect(c.BackupEnabled).To(Equal(createdDeployment.Spec.DeploymentSpec.BackupEnabled))
				})
			})

			By("Decreasing the Deployment disk size", func() {
				createdDeployment = performUpdate(ctx, 30*time.Minute, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
					deployment.Spec.DeploymentSpec.DiskSizeGB = pointer.MakePtr(15)
				})
				doDeploymentStatusChecks()
				checkAtlasState(func(c *admin.ClusterDescription20240805) {
					Expect(int(c.GetReplicationSpecs()[0].GetRegionConfigs()[0].ElectableSpecs.GetDiskSizeGB())).To(BeEquivalentTo(*createdDeployment.Spec.DeploymentSpec.DiskSizeGB))

					Expect(c.GetReplicationSpecs()[0].GetRegionConfigs()[0].ElectableSpecs.DiskSizeGB).To(BeAssignableToTypeOf(pointer.MakePtr[float64](0)), "DiskSizeGB is no longer a *float64, please check the spec!")
				})
			})

			By("Pausing the deployment", func() {
				createdDeployment = performUpdate(ctx, 30*time.Minute, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
					deployment.Spec.DeploymentSpec.Paused = pointer.MakePtr(true)
				})
				doDeploymentStatusChecks()
				checkAtlasState(func(c *admin.ClusterDescription20240805) {
					Expect(c.Paused).To(Equal(createdDeployment.Spec.DeploymentSpec.Paused))
				})
			})

			By("Updating the Deployment configuration while paused (should fail)", func() {
				_, err := akoretry.RetryUpdateOnConflict(ctx, k8sClient, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
					deployment.Spec.DeploymentSpec.BackupEnabled = pointer.MakePtr(false)
				})
				Expect(err).To(BeNil())

				Eventually(func() bool {
					return resources.CheckCondition(
						k8sClient,
						createdDeployment,
						api.FalseCondition(api.DeploymentReadyType).
							WithReason(string(workflow.DeploymentNotUpdatedInAtlas)).
							WithMessageRegexp("CANNOT_UPDATE_PAUSED_CLUSTER"),
					)
				}).
					WithTimeout(DeploymentUpdateTimeout).
					WithPolling(interval).
					Should(BeTrue())
			})

			By("Unpausing the deployment", func() {
				createdDeployment = performUpdate(ctx, 30*time.Minute, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
					deployment.Spec.DeploymentSpec.Paused = pointer.MakePtr(false)
				})
				doDeploymentStatusChecks()
				checkAtlasState(func(c *admin.ClusterDescription20240805) {
					Expect(c.Paused).To(Equal(createdDeployment.Spec.DeploymentSpec.Paused))
				})
			})

			By("Checking that modifications were applied after unpausing", func() {
				doDeploymentStatusChecks()
				checkAtlasState(func(c *admin.ClusterDescription20240805) {
					Expect(c.BackupEnabled).To(Equal(createdDeployment.Spec.DeploymentSpec.BackupEnabled))
				})
			})

			By("Setting incorrect instance size (should fail)", func() {
				var (
					oldSizeName string
					err         error
				)
				createdDeployment, err = akoretry.RetryUpdateOnConflict(ctx, k8sClient, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
					oldSizeName = deployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].ElectableSpecs.InstanceSize
					deployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].ElectableSpecs = &akov2.Specs{
						InstanceSize: "M42",
						NodeCount:    pointer.MakePtr(3),
					}
					deployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].ReadOnlySpecs = &akov2.Specs{
						InstanceSize: "M42",
						NodeCount:    pointer.MakePtr(0),
					}
					deployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].AnalyticsSpecs = &akov2.Specs{
						InstanceSize: "M42",
						NodeCount:    pointer.MakePtr(0),
					}
				})
				Expect(err).To(BeNil())

				Eventually(func() bool {
					return resources.CheckCondition(
						k8sClient,
						createdDeployment,
						api.FalseCondition(api.DeploymentReadyType).
							WithReason(string(workflow.DeploymentNotUpdatedInAtlas)).
							WithMessageRegexp(".*INVALID_ENUM_VALUE.*"),
					)
				}).WithTimeout(DeploymentUpdateTimeout).
					WithPolling(interval).
					Should(BeTrue())

				By("Fixing the Deployment", func() {
					createdDeployment = performUpdate(ctx, 30*time.Minute, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
						deployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].ElectableSpecs.InstanceSize = oldSizeName
						deployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].AnalyticsSpecs.InstanceSize = oldSizeName
						deployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].ReadOnlySpecs.InstanceSize = oldSizeName
					})
					doDeploymentStatusChecks()
					checkAtlasState()
				})
			})
		})
	})

	Describe("Create DBUser before deployment & check secrets", func() {
		It("Should Succeed", func() {
			By(fmt.Sprintf("Creating password Secret %s", UserPasswordSecret), func() {
				passwordSecret := buildPasswordSecret(namespace.Name, UserPasswordSecret, DBUserPassword)
				Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())
			})

			createdDBUser := akov2.DefaultDBUser(namespace.Name, "test-db-user", createdProject.Name).WithPasswordSecret(UserPasswordSecret)
			By(fmt.Sprintf("Creating the Database User %s", kube.ObjectKeyFromObject(createdDBUser)), func() {
				Expect(k8sClient.Create(context.Background(), createdDBUser)).ToNot(HaveOccurred())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, createdDBUser, api.TrueCondition(api.ReadyType))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

				checkUserInAtlas(createdProject.ID(), *createdDBUser)
			})

			createdDBUserFakeScope := akov2.DefaultDBUser(namespace.Name, "test-db-user-fake-scope", createdProject.Name).
				WithPasswordSecret(UserPasswordSecret).
				WithScope(akov2.DeploymentScopeType, "fake-deployment")
			By(fmt.Sprintf("Creating the Database User %s", kube.ObjectKeyFromObject(createdDBUserFakeScope)), func() {
				Expect(k8sClient.Create(context.Background(), createdDBUserFakeScope)).ToNot(HaveOccurred())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, createdDBUserFakeScope, api.FalseCondition(api.DatabaseUserReadyType).WithReason(string(workflow.DatabaseUserInvalidSpec)))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())
			})
			checkNumberOfConnectionSecrets(k8sClient, *createdProject, namespace.Name, 0)

			createdDeployment = akov2.DefaultAWSDeployment(namespace.Name, createdProject.Name).Lightweight()
			By(fmt.Sprintf("Creating the Deployment %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				performCreate(createdDeployment, 30*time.Minute)

				doDeploymentStatusChecks()
				checkAtlasState()
			})

			By("Checking connection Secrets", func() {
				Expect(tryConnect(createdProject.ID(), *createdDeployment, *createdDBUser)).To(Succeed())
				checkNumberOfConnectionSecrets(k8sClient, *createdProject, namespace.Name, 1)
				validateSecret(k8sClient, *createdProject, *createdDeployment, *createdDBUser)
			})
		})
	})

	Describe("Create deployment, user, delete deployment and check secrets are removed", func() {
		It("Should Succeed", func() {
			createdDeployment = akov2.DefaultAWSDeployment(namespace.Name, createdProject.Name).Lightweight()
			By(fmt.Sprintf("Creating the Deployment %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				Expect(k8sClient.Create(context.Background(), createdDeployment)).ToNot(HaveOccurred())

				Eventually(func(g Gomega) bool {
					return resources.CheckCondition(k8sClient, createdDeployment, api.TrueCondition(api.ReadyType), validateDeploymentCreatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

				doDeploymentStatusChecks()
				checkAtlasState()
			})

			passwordSecret := buildPasswordSecret(namespace.Name, UserPasswordSecret, DBUserPassword)
			Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())

			createdDBUser := akov2.DefaultDBUser(namespace.Name, "test-db-user", createdProject.Name).WithPasswordSecret(UserPasswordSecret)
			By(fmt.Sprintf("Creating the Database User %s", kube.ObjectKeyFromObject(createdDBUser)), func() {
				Expect(k8sClient.Create(context.Background(), createdDBUser)).ToNot(HaveOccurred())
				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, createdDBUser, api.TrueCondition(api.ReadyType))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())
			})

			By("Removing Atlas Deployment "+createdDeployment.Name, func() {
				Expect(k8sClient.Delete(context.Background(), createdDeployment)).To(Succeed())
				Eventually(checkAtlasDeploymentRemoved(createdProject.Status.ID, createdDeployment.GetDeploymentName()), 600, interval).Should(BeTrue())
			})

			By("Checking that Secrets got removed", func() {
				secretNames := []string{kube.NormalizeIdentifier(fmt.Sprintf("%s-%s-%s", createdProject.Spec.Name, createdDeployment.GetDeploymentName(), createdDBUser.Spec.Username))}
				createdDeployment = nil // prevent cleanup from failing due to deployment already deleted
				Eventually(checkSecretsDontExist(namespace.Name, secretNames), 50, interval).Should(BeTrue())
				checkNumberOfConnectionSecrets(k8sClient, *createdProject, namespace.Name, 0)
			})
		})
	})

	Describe("Deleting the deployment (not cleaning Atlas)", func() {
		It("Should Succeed", func() {
			By(`Creating the deployment with retention policy "keep" first`, func() {
				createdDeployment = akov2.DefaultAWSDeployment(namespace.Name, createdProject.Name).Lightweight()
				createdDeployment.ObjectMeta.Annotations = map[string]string{customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep}
				manualDeletion = true // We need to remove the deployment in Atlas manually to let project get removed
				performCreate(createdDeployment, 30*time.Minute)
			})
			By("Deleting the deployment - stays in Atlas", func() {
				Expect(k8sClient.Delete(context.Background(), createdDeployment)).To(Succeed())
				Eventually(func() {
					Expect(checkAtlasDeploymentRemoved(createdProject.Status.ID, createdDeployment.GetDeploymentName())()).Should(BeFalse())
					checkNumberOfConnectionSecrets(k8sClient, *createdProject, namespace.Name, 0)
				}).WithTimeout(5 * time.Minute).WithPolling(10 * time.Second)
			})
		})
	})

	Describe("Setting the deployment skip annotation should skip reconciliations.", func() {
		It("Should Succeed", func(ctx context.Context) {
			By(`Creating the deployment with reconciliation policy "skip" first`, func() {
				createdDeployment = akov2.DefaultAWSDeployment(namespace.Name, createdProject.Name).Lightweight()
				performCreate(createdDeployment, 30*time.Minute)

				var err error
				createdDeployment, err = akoretry.RetryUpdateOnConflict(ctx, k8sClient, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
					deployment.ObjectMeta.Annotations = map[string]string{customresource.ReconciliationPolicyAnnotation: customresource.ReconciliationPolicySkip}
					deployment.Spec.DeploymentSpec.Labels = append(createdDeployment.Spec.DeploymentSpec.Labels, common.LabelSpec{
						Key:   "some-key",
						Value: "some-value",
					})
				})
				Expect(err).To(BeNil())

				containsLabel := func(ac *admin.ClusterDescription20240805) bool {
					for _, label := range ac.GetLabels() {
						if label.GetKey() == "some-key" && label.GetValue() == "some-value" {
							return true
						}
					}
					return false
				}

				timeoutCtx, cancel := context.WithTimeout(ctx, time.Minute*2)
				defer cancel()
				Eventually(atlas.WaitForAtlasDeploymentStateToNotBeReached(timeoutCtx, atlasClient, createdProject.Name, createdDeployment.GetDeploymentName(), containsLabel))
			})
		})
	})

	Describe("Create the advanced deployment & change the InstanceSize", func() {
		It("Should Succeed", func(ctx context.Context) {
			createdDeployment = akov2.DefaultAwsAdvancedDeployment(namespace.Name, createdProject.Name)

			By(fmt.Sprintf("Creating the Advanced Deployment %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				performCreate(createdDeployment, 30*time.Minute)

				doDeploymentStatusChecks()
				checkAdvancedAtlasState()
			})

			By(fmt.Sprintf("Updating the InstanceSize of Advanced Deployment %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				var err error
				createdDeployment, err = akoretry.RetryUpdateOnConflict(ctx, k8sClient, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
					deployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].ElectableSpecs = &akov2.Specs{
						InstanceSize: "M20",
						NodeCount:    pointer.MakePtr(3),
					}
					deployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].ReadOnlySpecs = &akov2.Specs{
						InstanceSize: "M20",
						NodeCount:    pointer.MakePtr(0),
					}
					deployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].AnalyticsSpecs = &akov2.Specs{
						InstanceSize: "M20",
						NodeCount:    pointer.MakePtr(0),
					}
				})
				Expect(err).To(BeNil())

				Eventually(func(g Gomega) bool {
					return resources.CheckCondition(k8sClient, createdDeployment, api.TrueCondition(api.ReadyType), validateDeploymentUpdatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

				doDeploymentStatusChecks()
				checkAdvancedAtlasState()
			})

			By(fmt.Sprintf("Enable AutoScaling for the Advanced Deployment %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				var err error
				createdDeployment, err = akoretry.RetryUpdateOnConflict(ctx, k8sClient, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
					regionConfig := deployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0]
					regionConfig.ElectableSpecs.InstanceSize = "M10"
					regionConfig.ReadOnlySpecs.InstanceSize = "M10"
					regionConfig.AnalyticsSpecs.InstanceSize = "M10"
					regionConfig.AutoScaling = &akov2.AdvancedAutoScalingSpec{
						Compute: &akov2.ComputeSpec{
							Enabled:          pointer.MakePtr(true),
							MaxInstanceSize:  "M30",
							MinInstanceSize:  "M10",
							ScaleDownEnabled: pointer.MakePtr(true),
						},
						DiskGB: &akov2.DiskGB{
							Enabled: pointer.MakePtr(true),
						},
					}
				})
				Expect(err).To(BeNil())

				Eventually(func(g Gomega) bool {
					return resources.CheckCondition(k8sClient, createdDeployment, api.TrueCondition(api.ReadyType), validateDeploymentUpdatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

				doDeploymentStatusChecks()
				checkAdvancedAtlasState()
			})

			By(fmt.Sprintf("Update Instance Size Margins with AutoScaling for Deployment %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				var err error
				createdDeployment, err = akoretry.RetryUpdateOnConflict(ctx, k8sClient, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
					regionConfig := deployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0]
					regionConfig.AutoScaling.Compute.MinInstanceSize = "M20"
					regionConfig.ElectableSpecs.InstanceSize = "M20"
					regionConfig.ReadOnlySpecs.InstanceSize = "M20"
					regionConfig.AnalyticsSpecs.InstanceSize = "M20"
				})
				Expect(err).To(BeNil())

				Eventually(func(g Gomega) bool {
					return resources.CheckCondition(k8sClient, createdDeployment, api.TrueCondition(api.ReadyType), validateDeploymentUpdatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

				doDeploymentStatusChecks()
				checkAdvancedAtlasState()
			})
		})
	})

	Describe("Create the advanced deployment with enabled autoscaling", func() {
		It("Should Succeed", func(ctx context.Context) {
			createdDeployment = akov2.DefaultAwsAdvancedDeployment(namespace.Name, createdProject.Name)

			createdDeployment.Spec.DeploymentSpec.ReplicationSpecs = []*akov2.AdvancedReplicationSpec{
				{
					NumShards: 1,
					ZoneName:  "US_EAST_1",
					RegionConfigs: []*akov2.AdvancedRegionConfig{
						{
							AnalyticsSpecs: &akov2.Specs{
								DiskIOPS:      nil,
								EbsVolumeType: "",
								InstanceSize:  "M10",
								NodeCount:     pointer.MakePtr(1),
							},
							ElectableSpecs: &akov2.Specs{
								DiskIOPS:      nil,
								EbsVolumeType: "",
								InstanceSize:  "M10",
								NodeCount:     pointer.MakePtr(3),
							},
							ReadOnlySpecs: &akov2.Specs{
								DiskIOPS:      nil,
								EbsVolumeType: "",
								InstanceSize:  "M10",
								NodeCount:     pointer.MakePtr(1),
							},
							AutoScaling: &akov2.AdvancedAutoScalingSpec{
								DiskGB: &akov2.DiskGB{
									Enabled: pointer.MakePtr(true),
								},
								Compute: &akov2.ComputeSpec{
									Enabled:          pointer.MakePtr(true),
									ScaleDownEnabled: pointer.MakePtr(true),
									MinInstanceSize:  "M10",
									MaxInstanceSize:  "M40",
								},
							},
							BackingProviderName: "AWS",
							Priority:            pointer.MakePtr(7),
							ProviderName:        "AWS",
							RegionName:          "US_EAST_1",
						},
					},
				},
			}

			By(fmt.Sprintf("Creating the Advanced Deployment %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				performCreate(createdDeployment, 30*time.Minute)

				doDeploymentStatusChecks()
				checkAdvancedAtlasState()
			})

			By(fmt.Sprintf("Update autoscaling configuration with wrong values it should fail %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				previousDeployment := akov2.AtlasDeployment{}
				err := compat.JSONCopy(&previousDeployment, createdDeployment)
				Expect(err).NotTo(HaveOccurred())

				createdDeployment, err = akoretry.RetryUpdateOnConflict(ctx, k8sClient, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
					deployment.Spec.DeploymentSpec.ReplicationSpecs[0].
						RegionConfigs[0].
						AutoScaling.
						Compute.
						MinInstanceSize = "S"
				})
				Expect(err).To(BeNil())

				Eventually(func() bool {
					return resources.CheckCondition(
						k8sClient,
						createdDeployment,
						api.FalseCondition(api.ValidationSucceeded).
							WithReason(string(workflow.Internal)).
							WithMessageRegexp("instance size is invalid"),
					)
				}).WithTimeout(DeploymentUpdateTimeout).
					WithPolling(interval).
					Should(BeTrue())

				By(fmt.Sprintf("Update autoscaling configuration should update InstanceSize and DiskSizeGB of Advanced deployment %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
					previousDeployment := akov2.AtlasDeployment{}
					err := compat.JSONCopy(&previousDeployment, createdDeployment)
					Expect(err).NotTo(HaveOccurred())

					createdDeployment, err = akoretry.RetryUpdateOnConflict(ctx, k8sClient, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
						deployment.Spec.DeploymentSpec.ReplicationSpecs[0].
							RegionConfigs[0].
							ElectableSpecs.InstanceSize = "M20"
						deployment.Spec.DeploymentSpec.ReplicationSpecs[0].
							RegionConfigs[0].
							ReadOnlySpecs.InstanceSize = "M20"
						deployment.Spec.DeploymentSpec.ReplicationSpecs[0].
							RegionConfigs[0].
							AnalyticsSpecs.InstanceSize = "M20"
						deployment.Spec.DeploymentSpec.ReplicationSpecs[0].
							RegionConfigs[0].
							AutoScaling.
							Compute.
							MinInstanceSize = "M20"
					})
					Expect(err).To(BeNil())

					Eventually(func(g Gomega) bool {
						GinkgoWriter.Println("ProjectID", createdProject.ID(), "DeploymentName", createdDeployment.GetDeploymentName())
						current, _, err := atlasClient.ClustersApi.
							GetCluster(context.Background(), createdProject.ID(), createdDeployment.GetDeploymentName()).
							Execute()
						g.Expect(err).NotTo(HaveOccurred())
						g.Expect(current).NotTo(BeNil())

						replicas := current.GetReplicationSpecs()
						g.Expect(replicas[0].GetRegionConfigs()[0].AnalyticsSpecs.GetInstanceSize()).To(Equal("M20"))
						g.Expect(replicas[0].GetRegionConfigs()[0].ElectableSpecs.GetInstanceSize()).To(Equal("M20"))
						g.Expect(replicas[0].GetRegionConfigs()[0].ReadOnlySpecs.GetInstanceSize()).To(Equal("M20"))
						return true
					}).WithTimeout(2 * time.Minute).WithPolling(interval).Should(BeTrue())
				})
			})
		})
	})

	Describe("Set advanced deployment options", func() {
		It("Should Succeed", func(ctx context.Context) {
			createdDeployment = akov2.DefaultAWSDeployment(namespace.Name, createdProject.Name)
			createdDeployment.Spec.ProcessArgs = &akov2.ProcessArgs{
				JavascriptEnabled:  pointer.MakePtr(true),
				DefaultReadConcern: "available",
			}

			By(fmt.Sprintf("Creating the Deployment with Advanced Options %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				performCreate(createdDeployment, 30*time.Minute)

				doDeploymentStatusChecks()
				checkAdvancedDeploymentOptions(ctx, createdProject.ID(), createdDeployment)
			})

			By("Updating Advanced Deployment Options", func() {
				createdDeployment = performUpdate(ctx, 40*time.Minute, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
					deployment.Spec.ProcessArgs.JavascriptEnabled = pointer.MakePtr(false)
				})
				doDeploymentStatusChecks()
				checkAdvancedDeploymentOptions(ctx, createdProject.ID(), createdDeployment)
			})
		})
	})
})

var _ = Describe("AtlasDeployment", Ordered, Label("int", "AtlasDeployment", "deployment-backups"), func() {
	var (
		connectionSecret  *corev1.Secret
		createdProject    *akov2.AtlasProject
		createdDeployment *akov2.AtlasDeployment

		backupPolicyDefault   *akov2.AtlasBackupPolicy
		backupScheduleDefault *akov2.AtlasBackupSchedule
	)

	BeforeAll(func() {
		prepareControllers(false)
		connectionSecret = createConnectionSecret()
		createdProject = createProject(connectionSecret)
	})

	AfterAll(func() {
		deleteProjectFromKubernetes(createdProject)
		removeControllersAndNamespace()
	})

	Describe("Create default deployment with backups enabled", Label("focus-basic-backups"), func() {
		BeforeEach(func() {
			backupPolicyDefault = &akov2.AtlasBackupPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "policy-1",
					Namespace: namespace.Name,
				},
				Spec: akov2.AtlasBackupPolicySpec{
					Items: []akov2.AtlasBackupPolicyItem{
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

			backupScheduleDefault = &akov2.AtlasBackupSchedule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "schedule-1",
					Namespace: namespace.Name,
				},
				Spec: akov2.AtlasBackupScheduleSpec{
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

			createdDeployment = akov2.DefaultAWSDeployment(namespace.Name, createdProject.Name).WithBackupScheduleRef(common.ResourceRefNamespaced{
				Name:      backupScheduleDefault.Name,
				Namespace: backupScheduleDefault.Namespace,
			})
		})

		AfterEach(func() {
			deleteDeploymentFromKubernetes(createdProject, createdDeployment)
			deleteBackupDefsFromKubernetes(backupScheduleDefault, backupPolicyDefault)
		})

		It("Should succeed", func() {
			By(fmt.Sprintf("Creating deployment with backups enabled: %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				Expect(k8sClient.Create(context.Background(), createdDeployment)).NotTo(HaveOccurred())

				// Do not use Gomega function here like func(g Gomega) as it seems to hang when tests run in parallel
				Eventually(
					func() error {
						deployment, _, err := atlasClient.ClustersApi.
							GetCluster(context.Background(), createdProject.ID(), createdDeployment.GetDeploymentName()).
							Execute()
						if err != nil {
							return err
						}
						if deployment.GetStateName() != "IDLE" {
							return errors.New("deployment is not IDLE yet")
						}
						time.Sleep(10 * time.Second)
						return nil
					}).WithTimeout(40 * time.Minute).WithPolling(15 * time.Second).Should(Not(HaveOccurred()))

				Eventually(func() error {
					actualPolicy, _, err := atlasClient.CloudBackupsApi.
						GetBackupSchedule(context.Background(), createdProject.ID(), createdDeployment.GetDeploymentName()).
						Execute()
					if err != nil {
						return err
					}
					if len(actualPolicy.GetPolicies()[0].GetPolicyItems()) == 0 {
						return errors.New("policies == 0")
					}
					ap := actualPolicy.GetPolicies()[0].GetPolicyItems()[0]
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

	Describe("Create deployment with backups enabled and snapshot distribution", Label("focus-snapshot-distribution"), func() {
		var secondDeployment *akov2.AtlasDeployment
		bScheduleName := "schedule-1"

		AfterEach(func() {
			deleteDeploymentFromKubernetes(createdProject, createdDeployment)
			deleteDeploymentFromKubernetes(createdProject, secondDeployment)
			deleteBackupDefsFromKubernetes(backupScheduleDefault, backupPolicyDefault)
		})

		It("Should succeed", func(ctx context.Context) {
			By("Creating deployment with backups enabled", func() {
				createdDeployment = akov2.DefaultAwsAdvancedDeployment(namespace.Name, createdProject.Name)
				createdDeployment.Spec.DeploymentSpec.BackupEnabled = pointer.MakePtr(true)
				Expect(k8sClient.Create(context.Background(), createdDeployment)).NotTo(HaveOccurred())

				Eventually(func(g Gomega) {
					deployment, _, err := atlasClient.ClustersApi.
						GetCluster(context.Background(), createdProject.ID(), createdDeployment.Spec.DeploymentSpec.Name).
						Execute()
					g.Expect(err).Should(BeNil())
					g.Expect(deployment.GetStateName()).Should(Equal("IDLE"))
					g.Expect(deployment.GetBackupEnabled()).Should(BeTrue())
					g.Expect(len(deployment.GetReplicationSpecs())).ShouldNot(Equal(0))
				}).WithTimeout(40 * time.Minute).WithPolling(15 * time.Second).Should(Succeed())
			})

			By("Adding BackupSchedule with Snapshot distribution", func() {
				backupPolicyDefault = &akov2.AtlasBackupPolicy{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "policy-1",
						Namespace: namespace.Name,
					},
					Spec: akov2.AtlasBackupPolicySpec{
						Items: []akov2.AtlasBackupPolicyItem{
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
				backupScheduleDefault = &akov2.AtlasBackupSchedule{
					ObjectMeta: metav1.ObjectMeta{
						Name:      bScheduleName,
						Namespace: namespace.Name,
					},
					Spec: akov2.AtlasBackupScheduleSpec{
						AutoExportEnabled: false,
						PolicyRef: common.ResourceRefNamespaced{
							Name:      backupPolicyDefault.Name,
							Namespace: backupPolicyDefault.Namespace,
						},
						ReferenceHourOfDay:    12,
						ReferenceMinuteOfHour: 10,
						RestoreWindowDays:     5,
						UpdateSnapshots:       false,
						CopySettings: []akov2.CopySetting{
							{
								CloudProvider:    pointer.MakePtr("AWS"),
								RegionName:       pointer.MakePtr("US_WEST_1"),
								ShouldCopyOplogs: pointer.MakePtr(false),
								Frequencies:      []string{"MONTHLY"},
							},
						},
					},
				}
				Expect(k8sClient.Create(context.Background(), backupPolicyDefault)).Should(Succeed())
				Expect(k8sClient.Create(context.Background(), backupScheduleDefault)).Should(Succeed())

				Expect(k8sClient.Get(context.Background(), client.ObjectKeyFromObject(createdDeployment), createdDeployment)).Should(Succeed())
				var err error
				createdDeployment, err = akoretry.RetryUpdateOnConflict(ctx, k8sClient, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
					deployment.Spec.BackupScheduleRef = common.ResourceRefNamespaced{
						Name:      bScheduleName,
						Namespace: namespace.Name,
					}
				})
				Expect(err).To(BeNil())
			})

			By("Deployment is ready with backup and snapshot distribution configured", func() {
				Eventually(func(g Gomega) {
					validateDeploymentWithSnapshotDistribution(
						g,
						createdProject.ID(),
						createdDeployment.GetDeploymentName(),
						[]admin.DiskBackupCopySetting20240805{
							{
								CloudProvider:    pointer.MakePtr("AWS"),
								RegionName:       pointer.MakePtr("US_WEST_1"),
								ShouldCopyOplogs: pointer.MakePtr(false),
								Frequencies:      &[]string{"MONTHLY"},
							},
						},
					)
				}).WithTimeout(10 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
			})

			By("Creating a second deployment with backups enabled using same snapshot distribution configuration", func() {
				secondDeployment = akov2.DefaultAwsAdvancedDeployment(namespace.Name, createdProject.Name)
				secondDeployment.WithName("deployment-advanced-k8s-2")
				secondDeployment.Spec.DeploymentSpec.Name = "deployment-advanced-2"
				secondDeployment.Spec.DeploymentSpec.BackupEnabled = pointer.MakePtr(true)
				secondDeployment.Spec.BackupScheduleRef = common.ResourceRefNamespaced{
					Name:      bScheduleName,
					Namespace: namespace.Name,
				}
				Expect(k8sClient.Create(context.Background(), secondDeployment)).Should(Succeed())

				Eventually(func(g Gomega) {
					deployment, _, err := atlasClient.ClustersApi.
						GetCluster(context.Background(), createdProject.ID(), secondDeployment.Spec.DeploymentSpec.Name).
						Execute()
					g.Expect(err).Should(BeNil())
					g.Expect(deployment.GetStateName()).Should(Equal("IDLE"))
					g.Expect(deployment.GetBackupEnabled()).Should(BeTrue())
					g.Expect(len(deployment.GetReplicationSpecs())).ShouldNot(Equal(0))
				}).WithTimeout(40 * time.Minute).WithPolling(15 * time.Second).Should(Succeed())
			})

			By("The second Deployment is ready with backup and snapshot distribution configured", func() {
				Eventually(func(g Gomega) {
					validateDeploymentWithSnapshotDistribution(
						g,
						createdProject.ID(),
						secondDeployment.GetDeploymentName(),
						[]admin.DiskBackupCopySetting20240805{
							{
								CloudProvider:    pointer.MakePtr("AWS"),
								RegionName:       pointer.MakePtr("US_WEST_1"),
								ShouldCopyOplogs: pointer.MakePtr(false),
								Frequencies:      &[]string{"MONTHLY"},
							},
						},
					)
				}).WithTimeout(10 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
			})
		})
	})
})

var _ = Describe("AtlasDeploymentSharding", Label("int", "AtlasDeploymentSharding", "deployment-non-backups"), func() {
	var (
		connectionSecret  *corev1.Secret
		createdProject    *akov2.AtlasProject
		createdDeployment *akov2.AtlasDeployment
		manualDeletion    bool
	)

	BeforeEach(func() {
		prepareControllers(false)

		deployment.NewAtlasDeployments(atlasClient.ClustersApi, atlasClient.GlobalClustersApi, atlasClient.FlexClustersApi, false)
		createdDeployment = &akov2.AtlasDeployment{}

		manualDeletion = false

		connectionSecret = createConnectionSecret()
		createdProject = createProject(connectionSecret)
	})

	AfterEach(func() {
		if DeploymentDevMode {
			return
		}
		if manualDeletion && createdProject != nil {
			By("Deleting the deployment in Atlas manually", func() {
				// We need to remove the deployment in Atlas to let project get removed
				_, err := atlasClient.ClustersApi.
					DeleteCluster(context.Background(), createdProject.ID(), createdDeployment.GetDeploymentName()).
					Execute()
				Expect(err).NotTo(HaveOccurred())
				Eventually(checkAtlasDeploymentRemoved(createdProject.Status.ID, createdDeployment.GetDeploymentName()), 600, interval).Should(BeTrue())
				createdDeployment = nil
			})
		}
		if createdProject != nil && createdProject.Status.ID != "" {
			if createdDeployment != nil {
				deleteDeploymentFromKubernetes(createdProject, createdDeployment)
			}

			deleteProjectFromKubernetes(createdProject)
		}
		removeControllersAndNamespace()
	})

	doDeploymentStatusChecks := func() {
		By("Checking observed Deployment state", func() {
			doDeploymentStatusChecksFor(createdProject, createdDeployment)
		})
	}

	checkAtlasState := func(additionalChecks ...func(c *admin.ClusterDescription20240805)) {
		By("Verifying Deployment state in Atlas", func() {
			atlasDeploymentAsAtlas, _, err := atlasClient.ClustersApi.
				GetCluster(context.Background(), createdProject.Status.ID, createdDeployment.GetDeploymentName()).
				Execute()
			Expect(err).ToNot(HaveOccurred())

			for _, check := range additionalChecks {
				check(atlasDeploymentAsAtlas)
			}
		})
	}

	performCreate := func(deployment *akov2.AtlasDeployment, timeout time.Duration) {
		Expect(k8sClient.Create(context.Background(), deployment)).To(Succeed())

		Eventually(func(g Gomega) bool {
			return resources.CheckCondition(k8sClient, createdDeployment, api.TrueCondition(api.ReadyType), validateDeploymentCreatingFunc(g))
		}).WithTimeout(timeout).WithPolling(interval).Should(BeTrue())
	}

	Describe("Create deployment & change ReplicationSpecs", func() {
		It("Should Succeed", func(ctx context.Context) {
			createdDeployment = akov2.DefaultAWSDeployment(namespace.Name, createdProject.Name).
				WithInstanceSize("M30")

			// Atlas will add some defaults in case the Atlas Operator doesn't set them
			replicationSpecsCheck := func(deployment *admin.ClusterDescription20240805) {
				Expect(deployment.GetReplicationSpecs()[0].GetId()).NotTo(BeEmpty())
				Expect(deployment.GetReplicationSpecs()[0].GetZoneName()).To(Equal("Zone 1"))
				Expect(deployment.GetReplicationSpecs()[0].GetRegionConfigs()).To(HaveLen(1))
				Expect(deployment.GetReplicationSpecs()[0].GetRegionConfigs()[0]).NotTo(BeNil())
			}

			By(fmt.Sprintf("Creating the Deployment %s", kube.ObjectKeyFromObject(createdDeployment)), func() {
				performCreate(createdDeployment, 30*time.Minute)

				doDeploymentStatusChecks()

				singleNumShard := func(deployment *admin.ClusterDescription20240805) {
					Expect(len(deployment.GetReplicationSpecs())).To(Equal(1))
				}
				checkAtlasState(replicationSpecsCheck, singleNumShard)
			})

			By("Upgrade to sharded", func() {
				createdDeployment = performUpdate(ctx, 40*time.Minute, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
					deployment.Spec.DeploymentSpec.ClusterType = "SHARDED"
				})
				doDeploymentStatusChecks()

				singleNumShard := func(deployment *admin.ClusterDescription20240805) {
					Expect(len(deployment.GetReplicationSpecs())).To(Equal(1))
				}
				// ReplicationSpecs has the same defaults but the number of shards has changed
				checkAtlasState(replicationSpecsCheck, singleNumShard)
			})

			By("Increase number of shards", func() {
				numShards := 2
				createdDeployment = performUpdate(ctx, 40*time.Minute, client.ObjectKeyFromObject(createdDeployment), func(deployment *akov2.AtlasDeployment) {
					deployment.Spec.DeploymentSpec.ReplicationSpecs[0].NumShards = numShards
				})
				doDeploymentStatusChecks()

				twoNumShard := func(deployment *admin.ClusterDescription20240805) {
					Expect(len(deployment.GetReplicationSpecs())).To(Equal(numShards))
				}
				// ReplicationSpecs has the same defaults but the number of shards has changed
				checkAtlasState(replicationSpecsCheck, twoNumShard)
			})
		})
	})
})

func doDeploymentStatusChecksFor(createdProject *akov2.AtlasProject, createdDeployment *akov2.AtlasDeployment) {
	deploymentName := createdDeployment.GetDeploymentName()
	Expect(deploymentName).ToNot(BeEmpty())

	atlasDeployment, _, err := atlasClient.ClustersApi.
		GetCluster(context.Background(), createdProject.Status.ID, deploymentName).
		Execute()
	Expect(err).ToNot(HaveOccurred())

	Expect(createdDeployment.Status.ConnectionStrings).NotTo(BeNil())
	Expect(createdDeployment.Status.ConnectionStrings.Standard).To(Equal(atlasDeployment.ConnectionStrings.GetStandard()))
	Expect(createdDeployment.Status.ConnectionStrings.StandardSrv).To(Equal(atlasDeployment.ConnectionStrings.GetStandardSrv()))
	Expect(createdDeployment.Status.MongoDBVersion).To(Equal(atlasDeployment.GetMongoDBVersion()))
	Expect(createdDeployment.Status.StateName).To(Equal("IDLE"))
	Expect(createdDeployment.Status.Conditions).To(HaveLen(4))
	Expect(createdDeployment.Status.Conditions).To(ConsistOf(conditions.MatchConditions(
		api.TrueCondition(api.DeploymentReadyType),
		api.TrueCondition(api.ReadyType),
		api.TrueCondition(api.ValidationSucceeded),
		api.TrueCondition(api.ResourceVersionStatus),
	)))
	Expect(createdDeployment.Status.ObservedGeneration).To(Equal(createdDeployment.Generation))
}

func validateDeploymentCreatingFunc(g Gomega) func(a api.AtlasCustomResource) {
	startedCreation := false
	return func(a api.AtlasCustomResource) {
		c := a.(*akov2.AtlasDeployment)
		if c.Status.StateName != "" {
			startedCreation = true
		}
		// When the create request has been made to Atlas - we expect the following status
		if startedCreation {
			g.Expect(c.Status.StateName).To(Or(Equal("CREATING"), Equal("IDLE")), fmt.Sprintf("Current conditions: %+v", c.Status.Conditions))
			expectedConditionsMatchers := conditions.MatchConditions(
				api.FalseCondition(api.DeploymentReadyType).WithReason(string(workflow.DeploymentCreating)).WithMessageRegexp("deployment is provisioning"),
				api.FalseCondition(api.ReadyType),
				api.TrueCondition(api.ValidationSucceeded),
				api.TrueCondition(api.ResourceVersionStatus),
			)
			g.Expect(c.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
		} else {
			// Otherwise there could have been some exception in Atlas on creation - let's check the conditions
			condition, ok := conditions.FindConditionByType(c.Status.Conditions, api.DeploymentReadyType)
			g.Expect(ok).To(BeFalse(), fmt.Sprintf("Unexpected condition: %v", condition))
		}
	}
}

func validateDeploymentUpdatingFunc(g Gomega) func(a api.AtlasCustomResource) {
	isIdle := true
	return func(a api.AtlasCustomResource) {
		c := a.(*akov2.AtlasDeployment)
		// It's ok if the first invocations see IDLE
		if c.Status.StateName != "IDLE" {
			isIdle = false
		}
		// When the create request has been made to Atlas - we expect the following status
		if !isIdle {
			g.Expect(c.Status.StateName).To(Or(Equal("UPDATING"), Equal("REPAIRING")), fmt.Sprintf("Current conditions: %+v", c.Status.Conditions))
			expectedConditionsMatchers := conditions.MatchConditions(
				api.FalseCondition(api.DeploymentReadyType).WithReason(string(workflow.DeploymentUpdating)).WithMessageRegexp("deployment is updating"),
				api.FalseCondition(api.ReadyType),
				api.TrueCondition(api.ValidationSucceeded),
				api.TrueCondition(api.ResourceVersionStatus),
			)
			g.Expect(c.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
		}
	}
}

func validateDeploymentWithSnapshotDistribution(g Gomega, projectID, deploymentName string, copySettings []admin.DiskBackupCopySetting20240805) {
	atlasCluster, _, err := atlasClient.ClustersApi.GetCluster(context.Background(), projectID, deploymentName).Execute()
	g.Expect(err).Should(BeNil())
	g.Expect(atlasCluster.GetStateName()).Should(Equal("IDLE"))
	g.Expect(atlasCluster.GetBackupEnabled()).Should(BeTrue())

	for i := range copySettings {
		copySettings[i].SetZoneId(atlasCluster.GetReplicationSpecs()[0].GetZoneId())
	}

	atlasBSchedule, _, err := atlasClient.CloudBackupsApi.
		GetBackupSchedule(context.Background(), projectID, deploymentName).
		Execute()
	g.Expect(err).Should(BeNil())
	g.Expect(len(atlasBSchedule.GetCopySettings())).ShouldNot(Equal(0))
	g.Expect(atlasBSchedule.GetCopySettings()).Should(Equal(copySettings))
}

// checkAtlasDeploymentRemoved returns true if the Atlas Deployment is removed from Atlas. Note the behavior: the deployment
// is removed from Atlas as soon as the DELETE API call has been made. This is different from the case when the
// deployment is terminated from UI (in this case GET request succeeds while the deployment is being terminated)
func checkAtlasDeploymentRemoved(projectID string, deploymentName string) func() bool {
	return func() bool {
		_, r, err := atlasClient.ClustersApi.GetCluster(context.Background(), projectID, deploymentName).Execute()
		if err != nil {
			if httputil.StatusCode(r) == http.StatusNotFound {
				return true
			}
		}

		return false
	}
}

func checkAtlasFlexInstanceRemoved(projectID string, deploymentName string) func() bool {
	return func() bool {
		_, r, err := atlasClient.FlexClustersApi.
			GetFlexCluster(context.Background(), projectID, deploymentName).
			Execute()
		if err != nil {
			if httputil.StatusCode(r) == http.StatusNotFound {
				return true
			}
		}

		return false
	}
}

func deleteAtlasDeployment(projectID string, deploymentName string) error {
	_, err := atlasClient.ClustersApi.DeleteCluster(context.Background(), projectID, deploymentName).Execute()
	return err
}

func deleteFlexInstance(projectID string, deploymentName string) error {
	_, err := atlasClient.FlexClustersApi.
		DeleteFlexCluster(context.Background(), projectID, deploymentName).
		Execute()
	return err
}

func createConnectionSecret() *corev1.Secret {
	connectionSecret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ConnectionSecretName,
			Namespace: namespace.Name,
			Labels: map[string]string{
				secretservice.TypeLabelKey: secretservice.CredLabelVal,
			},
		},
		StringData: secretData(),
	}
	By(fmt.Sprintf("Creating the Secret %s", kube.ObjectKeyFromObject(&connectionSecret)), func() {
		Expect(k8sClient.Create(context.Background(), &connectionSecret)).To(Succeed())
	})
	return &connectionSecret
}

func createProject(connectionSecret *corev1.Secret) *akov2.AtlasProject {
	createdProject := akov2.DefaultProject(namespace.Name, connectionSecret.Name).WithIPAccessList(
		project.NewIPAccessList().WithCIDR("0.0.0.0/0"),
	)
	By("Creating the project "+createdProject.Name, func() {
		if DeploymentDevMode {
			// While developing tests we need to reuse the same project
			createdProject.Spec.Name = "dev-test atlas-project"
		}
		Expect(k8sClient.Create(context.Background(), createdProject)).To(Succeed())
		Eventually(func() bool {
			return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType))
		}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())
	})
	return createdProject
}

func deleteBackupDefsFromKubernetes(schedule *akov2.AtlasBackupSchedule, policy *akov2.AtlasBackupPolicy) {
	By("Deleting the schedule and policy in Kubernetes (should have no finalizers by now)", func() {
		Expect(k8sClient.Delete(context.Background(), schedule)).NotTo(HaveOccurred())
		Expect(k8sClient.Delete(context.Background(), policy)).NotTo(HaveOccurred())

		policyRef := kube.ObjectKey(policy.Namespace, policy.Name)
		Eventually(func() bool {
			p := &akov2.AtlasBackupPolicy{}
			return k8serrors.IsNotFound(k8sClient.Get(context.Background(), policyRef, p))
		}).WithTimeout(30 * time.Second).WithPolling(PollingInterval).Should(BeTrue())

		scheduleRef := kube.ObjectKey(schedule.Namespace, schedule.Name)
		Eventually(func() bool {
			s := &akov2.AtlasBackupSchedule{}
			return k8serrors.IsNotFound(k8sClient.Get(context.Background(), scheduleRef, s))
		}).WithTimeout(30 * time.Second).WithPolling(PollingInterval).Should(BeTrue())
	})
}

func deleteDeploymentFromKubernetes(project *akov2.AtlasProject, deployment *akov2.AtlasDeployment) {
	By(fmt.Sprintf("Removing Atlas Deployment %q", deployment.Name), func() {
		Expect(k8sClient.Delete(context.Background(), deployment)).To(Succeed())
		deploymentName := deployment.GetDeploymentName()
		if customresource.IsResourcePolicyKeep(deployment) || customresource.ReconciliationShouldBeSkipped(deployment) {
			By("Removing Atlas Deployment " + deployment.Name + " from Atlas manually")
			Expect(deleteAtlasDeployment(project.Status.ID, deploymentName)).To(Succeed())
		}
		Eventually(checkAtlasDeploymentRemoved(project.Status.ID, deploymentName), 600, interval).Should(BeTrue())
	})
}

func deleteProjectFromKubernetes(project *akov2.AtlasProject) {
	By(fmt.Sprintf("Removing Atlas Project %s", project.Status.ID), func() {
		Expect(k8sClient.Delete(context.Background(), project)).To(Succeed())
		Eventually(checkAtlasProjectRemoved(project.Status.ID), 240, interval).Should(BeTrue())
	})
}

// mergedAdvancedDeployment is clone of atlasdeployment.MergedAdvancedDeployment
func mergedAdvancedDeployment(
	atlasDeploymentAsAtlas admin.ClusterDescription20240805,
	specDeployment akov2.AdvancedDeploymentSpec,
) (mergedDeployment akov2.AdvancedDeploymentSpec, atlasDeployment akov2.AdvancedDeploymentSpec, err error) {
	if atlasDeploymentAsAtlas.ReplicationSpecs != nil {
		for _, replicationSpec := range atlasDeploymentAsAtlas.GetReplicationSpecs() {
			for _, regionConfig := range replicationSpec.GetRegionConfigs() {
				if regionConfig.ElectableSpecs != nil &&
					regionConfig.ElectableSpecs.GetInstanceSize() == atlasdeployment.FreeTier {
				}
			}
		}
	}

	var value float64

	if specs := atlasDeploymentAsAtlas.GetReplicationSpecs(); len(specs) > 0 {
		if configs := specs[0].GetRegionConfigs(); len(configs) > 0 {
			if e, ok := configs[0].GetElectableSpecsOk(); ok {
				value = e.GetDiskSizeGB()
			} else if r, ok := configs[0].GetReadOnlySpecsOk(); ok {
				value = r.GetDiskSizeGB()
			} else if a, ok := configs[0].GetAnalyticsSpecsOk(); ok {
				value = a.GetDiskSizeGB()
			}
		}
	}
	if value >= 1 {
		atlasDeployment.DiskSizeGB = pointer.MakePtr(int(value))
	}

	if err = compat.JSONCopy(&atlasDeployment, atlasDeploymentAsAtlas); err != nil {
		return mergedDeployment, atlasDeployment, err
	}

	for _, region := range specDeployment.ReplicationSpecs[0].RegionConfigs {
		if region == nil {
			return
		}

		var notNilSpecs akov2.Specs
		if region.ElectableSpecs != nil {
			notNilSpecs = *region.ElectableSpecs
		} else if region.ReadOnlySpecs != nil {
			notNilSpecs = *region.ReadOnlySpecs
		} else if region.AnalyticsSpecs != nil {
			notNilSpecs = *region.AnalyticsSpecs
		}

		if region.ElectableSpecs == nil {
			region.ElectableSpecs = &notNilSpecs
			region.ElectableSpecs.NodeCount = pointer.MakePtr(0)
		}

		if region.ReadOnlySpecs == nil {
			region.ReadOnlySpecs = &notNilSpecs
			region.ReadOnlySpecs.NodeCount = pointer.MakePtr(0)
		}

		if region.AnalyticsSpecs == nil {
			region.AnalyticsSpecs = &notNilSpecs
			region.AnalyticsSpecs.NodeCount = pointer.MakePtr(0)
		}
	}

	mergedDeployment = akov2.AdvancedDeploymentSpec{}

	if err = compat.JSONCopy(&mergedDeployment, atlasDeployment); err != nil {
		return
	}

	if err = compat.JSONCopy(&mergedDeployment, specDeployment); err != nil {
		return
	}

	for i, replicationSpec := range atlasDeployment.ReplicationSpecs {
		for k, v := range replicationSpec.RegionConfigs {
			// the response does not return backing provider names in some situations.
			// if this is the case, we want to strip these fields so they do not cause a bad comparison.
			if v.BackingProviderName == "" && k < len(mergedDeployment.ReplicationSpecs[i].RegionConfigs) {
				mergedDeployment.ReplicationSpecs[i].RegionConfigs[k].BackingProviderName = ""
			}
		}
	}

	atlasDeployment.MongoDBVersion = ""
	mergedDeployment.MongoDBVersion = ""

	return
}

func performUpdate[T any](ctx context.Context, timeout time.Duration, key client.ObjectKey, mutator func(*T)) *T {
	obj, err := akoretry.RetryUpdateOnConflict(ctx, k8sClient, key, mutator)
	Expect(err).To(BeNil())

	clientObj := any(obj).(api.AtlasCustomResource)
	Eventually(func(g Gomega) bool {
		return resources.CheckCondition(k8sClient, clientObj, api.TrueCondition(api.ReadyType), validateDeploymentUpdatingFunc(g))
	}).WithTimeout(timeout).WithPolling(interval).Should(BeTrue())

	return obj
}
