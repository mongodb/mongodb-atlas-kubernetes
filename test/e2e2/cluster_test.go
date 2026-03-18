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

package e2e2_test

import (
	"fmt"
	"time"

	"github.com/crd2go/crd2go/k8s"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas-sdk/v20250312016/admin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/controller-runtime/pkg/client"

	generatedv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/generated/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/operator"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/resources"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/samples"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/testparams"
)

const (
	ClusterCRDName = "clusters.atlas.generated.mongodb.com"

	clusterCreateTimeout   = 20 * time.Minute
	clusterUpdateTimeout   = 20 * time.Minute
	clusterConvertTimeout  = 25 * time.Minute
	clusterDeleteTimeout   = 15 * time.Minute
	clusterPollingInterval = 20 * time.Second
)

var _ = Describe("Cluster Generated v1", Ordered, func() {
	var ako operator.Operator
	var kubeClient client.Client
	var atlasClient *admin.APIClient

	var orgID string

	var testNamespace *corev1.Namespace
	var cluster *generatedv1.Cluster
	var group *generatedv1.Group
	var ctx = suiteCtx

	groupName := fmt.Sprintf("group-%s", rand.String(6))
	clusterName := fmt.Sprintf("cluster-%s", rand.String(6))

	_ = BeforeAll(func() {
		By("Should start the operator", func() {
			ako = runTestAKO(DefaultGlobalCredentials, control.MustEnvVar("OPERATOR_NAMESPACE"), false)
			ako.Start(ctx, GinkgoT())

			DeferCleanup(func() {
				if ako != nil {
					ako.Stop(GinkgoT())
				}
			})
		})

		By("Should create a kubernetes client", func() {
			testClient, err := kube.NewTestClient()
			Expect(err).To(Succeed())
			kubeClient = testClient
			Expect(kube.AssertCRDNames(ctx, kubeClient, GroupCRDName, ClusterCRDName)).To(Succeed())
		})

		By("Should create an Atlas Client", func() {
			atlasClient, orgID = newTestAtlasClient()
		})

		By("Should create namespace and copy credentials", func() {
			Expect(ako.Running()).To(BeTrue(), "Operator must be running")

			testNamespace = &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("cluster-ns-%s", rand.String(6)),
			}}
			Expect(kubeClient.Create(ctx, testNamespace)).To(Succeed())
			Expect(resources.CopyCredentialsToNamespace(
				ctx,
				kubeClient,
				DefaultGlobalCredentials,
				control.MustEnvVar("OPERATOR_NAMESPACE"),
				testNamespace.Name, GinkGoFieldOwner),
			).To(Succeed())
		})

		By("Should create a Group", func() {
			testParams := testparams.New(orgID, control.MustEnvVar("OPERATOR_NAMESPACE"), DefaultGlobalCredentials).
				WithGroupName(groupName).
				WithNamespace(testNamespace.Name)

			objs := samples.MustLoadSampleObjects("atlas_generated_v1_group.yaml")
			Expect(len(objs)).To(Equal(1))
			group = objs[0].(*generatedv1.Group)
			applyTestParamsToGroup(group, testParams)
			Expect(kubeClient.Create(ctx, group)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(resources.CheckResourceReady(ctx, kubeClient, group)).To(Succeed())
			}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

			Expect(group.Status.V20250312).NotTo(BeNil())
			Expect(group.Status.V20250312.Id).NotTo(BeNil())
		})
	})

	_ = AfterAll(func() {
		if kubeClient == nil {
			return
		}

		By("Should delete cluster", func() {
			if cluster != nil {
				_ = kubeClient.Delete(ctx, cluster)
				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceDeleted(ctx, kubeClient, cluster)).To(Succeed())
				}).WithContext(ctx).WithTimeout(clusterDeleteTimeout).WithPolling(clusterPollingInterval).Should(Succeed())
			}
		})

		By("Should delete group", func() {
			if group != nil {
				_ = kubeClient.Delete(ctx, group)
				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceDeleted(ctx, kubeClient, group)).To(Succeed())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			}
		})

		By("Should delete namespace", func() {
			Expect(
				kubeClient.Delete(ctx, testNamespace),
			).To(Succeed())
			Eventually(func(g Gomega) bool {
				return kubeClient.Get(ctx, client.ObjectKeyFromObject(testNamespace), testNamespace) == nil
			}).WithContext(ctx).WithTimeout(time.Minute).WithPolling(time.Second).To(BeFalse())
		})
	})

	Context("Share cluster with resource policy", Label("focus-cluster-shared"), func() {
		It("Should create shared cluster (M0) with ResourcePolicy Keep annotation", Label("focus-cluster-create-shared"), func() {
			cluster = newSharedCluster(clusterName, testNamespace.Name, groupName)
			cluster.SetAnnotations(map[string]string{
				customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep,
			})
			Expect(kubeClient.Create(ctx, cluster)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(resources.CheckResourceReady(ctx, kubeClient, cluster)).To(Succeed())
			}).WithContext(ctx).WithTimeout(clusterCreateTimeout).WithPolling(clusterPollingInterval).Should(Succeed())

			Expect(cluster.Status.V20250312).NotTo(BeNil())
			Expect(cluster.Status.V20250312.Id).NotTo(BeNil())
			Expect(*cluster.Status.V20250312.Id).NotTo(BeEmpty())
		})

		It("Should delete cluster from Kubernetes but should not delete from Atlas", func() {
			Expect(kubeClient.Delete(ctx, cluster)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(resources.CheckResourceDeleted(ctx, kubeClient, cluster)).To(Succeed())
			}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

			atlasCluster, _, err := atlasClient.ClustersApi.GetCluster(ctx, *group.Status.V20250312.Id, clusterName).Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(atlasCluster).NotTo(BeNil())
		})

		It("Should re-import cluster using external-id annotation", func() {
			cluster = newSharedCluster(clusterName, testNamespace.Name, groupName)
			cluster.SetAnnotations(map[string]string{
				"mongodb.com/external-id": clusterName,
			})
			Expect(kubeClient.Create(ctx, cluster)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(resources.CheckResourceReady(ctx, kubeClient, cluster)).To(Succeed())
			}).WithContext(ctx).WithTimeout(clusterCreateTimeout).WithPolling(clusterPollingInterval).Should(Succeed())

			Expect(cluster.Status.V20250312).NotTo(BeNil())
			Expect(cluster.Status.V20250312.Id).NotTo(BeNil())
			Expect(*cluster.Status.V20250312.Id).NotTo(BeEmpty())
		})

		It("Should delete cluster and delete from Atlas", func() {
			Expect(kubeClient.Delete(ctx, cluster)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(resources.CheckResourceDeleted(ctx, kubeClient, cluster)).To(Succeed())
			}).WithContext(ctx).WithTimeout(clusterDeleteTimeout).WithPolling(clusterPollingInterval).Should(Succeed())

			atlasCluster, _, err := atlasClient.ClustersApi.GetCluster(ctx, *group.Status.V20250312.Id, clusterName).Execute()
			Expect(err).To(HaveOccurred())
			Expect(atlasCluster).To(BeNil())
		})
	})

	Context("Cluster complex lifecycle", Label("focus-cluster-lifecycle"), func() {
		It("Should create a REPLICASET cluster", func() {
			cluster = newReplicaSetCluster(clusterName, testNamespace.GetName(), *group.Status.V20250312.Id, DefaultGlobalCredentials)
			Expect(kubeClient.Create(ctx, cluster)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(resources.CheckResourceReady(ctx, kubeClient, cluster)).To(Succeed())
			}).WithContext(ctx).WithTimeout(clusterCreateTimeout).WithPolling(clusterPollingInterval).Should(Succeed())

			Expect(cluster.Status.V20250312).NotTo(BeNil())
			Expect(cluster.Status.V20250312.Id).NotTo(BeNil())
			Expect(*cluster.Status.V20250312.Id).NotTo(BeEmpty())
		})

		It("Should update cluster with autoscaling and extra config", func() {
			Expect(kubeClient.Get(ctx, client.ObjectKeyFromObject(cluster), cluster)).To(Succeed())

			entry := cluster.Spec.V20250312.Entry
			entry.BackupEnabled = pointer.MakePtr(true)
			entry.PitEnabled = pointer.MakePtr(true)
			entry.MongoDBMajorVersion = pointer.MakePtr("8.0")
			entry.RedactClientLogData = pointer.MakePtr(true)
			entry.Tags = &[]generatedv1.Tags{
				{Key: "environment", Value: "test"},
				{Key: "managed-by", Value: "ako"},
			}

			specs := *entry.ReplicationSpecs
			rc := &(*specs[0].RegionConfigs)[0]
			rc.ElectableSpecs.DiskSizeGB = pointer.MakePtr(20.0)
			rc.AutoScaling = &generatedv1.AnalyticsAutoScaling{
				Compute: &generatedv1.Compute{
					Enabled:          pointer.MakePtr(true),
					MaxInstanceSize:  pointer.MakePtr("M30"),
					MinInstanceSize:  pointer.MakePtr("M10"),
					ScaleDownEnabled: pointer.MakePtr(true),
				},
				DiskGB: &generatedv1.DiskGB{
					Enabled: pointer.MakePtr(true),
				},
			}

			Expect(kubeClient.Update(ctx, cluster)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(resources.CheckResourceUpdated(ctx, kubeClient, cluster)).To(Succeed())
			}).WithContext(ctx).WithTimeout(clusterUpdateTimeout).WithPolling(clusterPollingInterval).Should(Succeed())
		})

		It("Should upgrade cluster from REPLICASET to SHARDED", func() {
			Expect(kubeClient.Get(ctx, client.ObjectKeyFromObject(cluster), cluster)).To(Succeed())

			entry := cluster.Spec.V20250312.Entry
			entry.ClusterType = pointer.MakePtr("SHARDED")

			existingSpecs := *entry.ReplicationSpecs
			secondShard := generatedv1.ReplicationSpecs{
				RegionConfigs: &[]generatedv1.RegionConfigs{
					{
						ProviderName: pointer.MakePtr("AWS"),
						RegionName:   pointer.MakePtr("US_EAST_1"),
						Priority:     pointer.MakePtr(7),
						ElectableSpecs: &generatedv1.ElectableSpecs{
							InstanceSize: pointer.MakePtr("M10"),
							NodeCount:    pointer.MakePtr(3),
						},
						AutoScaling: &generatedv1.AnalyticsAutoScaling{
							Compute: &generatedv1.Compute{
								Enabled:          pointer.MakePtr(true),
								MaxInstanceSize:  pointer.MakePtr("M30"),
								MinInstanceSize:  pointer.MakePtr("M10"),
								ScaleDownEnabled: pointer.MakePtr(true),
							},
							DiskGB: &generatedv1.DiskGB{
								Enabled: pointer.MakePtr(true),
							},
						},
					},
				},
			}
			updatedSpecs := append(existingSpecs, secondShard)
			entry.ReplicationSpecs = &updatedSpecs

			Expect(kubeClient.Update(ctx, cluster)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(resources.CheckResourceUpdated(ctx, kubeClient, cluster)).To(Succeed())
			}).WithContext(ctx).WithTimeout(clusterConvertTimeout).WithPolling(clusterPollingInterval).Should(Succeed())
		})

		It("Should update cluster with Independent Sharding Scaling", func() {
			Expect(kubeClient.Get(ctx, client.ObjectKeyFromObject(cluster), cluster)).To(Succeed())

			entry := cluster.Spec.V20250312.Entry
			entry.ReplicaSetScalingStrategy = pointer.MakePtr("NODE_TYPE")

			specs := *entry.ReplicationSpecs
			for i := range specs {
				rcs := *specs[i].RegionConfigs
				for j := range rcs {
					rcs[j].AutoScaling = &generatedv1.AnalyticsAutoScaling{
						Compute: &generatedv1.Compute{
							Enabled: pointer.MakePtr(false),
						},
						DiskGB: &generatedv1.DiskGB{
							Enabled: pointer.MakePtr(false),
						},
					}
				}
			}

			shard1RC := &(*(*entry.ReplicationSpecs)[0].RegionConfigs)[0]
			shard1RC.ElectableSpecs.InstanceSize = pointer.MakePtr("M10")
			shard1RC.ElectableSpecs.DiskSizeGB = pointer.MakePtr(30.0)

			shard2RC := &(*(*entry.ReplicationSpecs)[1].RegionConfigs)[0]
			shard2RC.ElectableSpecs.InstanceSize = pointer.MakePtr("M20")
			shard2RC.ElectableSpecs.DiskSizeGB = pointer.MakePtr(40.0)

			Expect(kubeClient.Update(ctx, cluster)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(resources.CheckResourceUpdated(ctx, kubeClient, cluster)).To(Succeed())
			}).WithContext(ctx).WithTimeout(clusterUpdateTimeout).WithPolling(clusterPollingInterval).Should(Succeed())
		})

		It("Should upgrade cluster from SHARDED to GEOSHARDED", func() {
			Expect(kubeClient.Get(ctx, client.ObjectKeyFromObject(cluster), cluster)).To(Succeed())

			entry := cluster.Spec.V20250312.Entry
			entry.ClusterType = pointer.MakePtr("GEOSHARDED")
			entry.GlobalClusterSelfManagedSharding = pointer.MakePtr(false)

			specs := *entry.ReplicationSpecs

			specs[0].ZoneName = pointer.MakePtr("Zone 1")
			(*specs[0].RegionConfigs)[0].RegionName = pointer.MakePtr("US_EAST_1")
			(*specs[0].RegionConfigs)[0].ElectableSpecs.InstanceSize = pointer.MakePtr("M10")

			specs[1].ZoneName = pointer.MakePtr("Zone 2")
			(*specs[1].RegionConfigs)[0].RegionName = pointer.MakePtr("EU_WEST_1")
			(*specs[1].RegionConfigs)[0].ElectableSpecs.InstanceSize = pointer.MakePtr("M10")

			Expect(kubeClient.Update(ctx, cluster)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(resources.CheckResourceUpdated(ctx, kubeClient, cluster)).To(Succeed())
			}).WithContext(ctx).WithTimeout(clusterConvertTimeout).WithPolling(clusterPollingInterval).Should(Succeed())
		})
	})
})

func newReplicaSetCluster(name, namespace, groupId, secretName string) *generatedv1.Cluster {
	return &generatedv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: generatedv1.ClusterSpec{
			ConnectionSecretRef: &k8s.LocalReference{
				Name: secretName,
			},
			V20250312: &generatedv1.ClusterSpecV20250312{
				GroupId: pointer.MakePtr(groupId),
				Entry: &generatedv1.V20250312Entry{
					Name:        pointer.MakePtr(name),
					ClusterType: pointer.MakePtr("REPLICASET"),
					ReplicationSpecs: &[]generatedv1.ReplicationSpecs{
						{
							RegionConfigs: &[]generatedv1.RegionConfigs{
								{
									ProviderName: pointer.MakePtr("AWS"),
									RegionName:   pointer.MakePtr("US_EAST_1"),
									Priority:     pointer.MakePtr(7),
									ElectableSpecs: &generatedv1.ElectableSpecs{
										InstanceSize: pointer.MakePtr("M10"),
										NodeCount:    pointer.MakePtr(3),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func newSharedCluster(name, namespace, groupRefName string) *generatedv1.Cluster {
	return &generatedv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: generatedv1.ClusterSpec{
			V20250312: &generatedv1.ClusterSpecV20250312{
				GroupRef: &k8s.LocalReference{
					Name: groupRefName,
				},
				Entry: &generatedv1.V20250312Entry{
					Name:        pointer.MakePtr(name),
					ClusterType: pointer.MakePtr("REPLICASET"),
					ReplicationSpecs: &[]generatedv1.ReplicationSpecs{
						{
							RegionConfigs: &[]generatedv1.RegionConfigs{
								{
									ProviderName:        pointer.MakePtr("TENANT"),
									BackingProviderName: pointer.MakePtr("AWS"),
									RegionName:          pointer.MakePtr("US_EAST_1"),
									Priority:            pointer.MakePtr(7),
									ElectableSpecs: &generatedv1.ElectableSpecs{
										InstanceSize: pointer.MakePtr("M0"),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
