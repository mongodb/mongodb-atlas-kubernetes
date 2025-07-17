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

package e2e_test

import (
	"context"
	"os/exec"
	"syscall"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	akov2common "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
)

var _ = Describe("Flex to Dedicated Upgrade", Ordered, Label("flex-to-dedicated"), func() {
	var kubeClient client.Client
	var testNamespace *corev1.Namespace
	var resourcePrefix string
	var akoCmd *exec.Cmd

	_ = BeforeAll(func(ctx context.Context) {
		cmd, err := k8s.RunManagerBinary(false)
		Expect(err).To(BeNil())
		akoCmd = cmd

		client, err := k8s.CreateNewClient()
		Expect(err).To(BeNil())
		kubeClient = client
	})

	_ = AfterAll(func(ctx context.Context) {
		if akoCmd != nil {
			if akoCmd.Process != nil {
				Expect(syscall.Kill(akoCmd.Process.Pid, syscall.SIGTERM)).To(Succeed())
			}

			Expect(akoCmd.Wait()).To(Succeed())
		}
	})

	_ = BeforeEach(func(ctx context.Context) {
		Expect(akoCmd.ProcessState).To(BeNil())

		resourcePrefix = utils.RandomName("flex-to-dedicated")

		testNamespace = &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{
			Name: resourcePrefix + "-ns",
		}}
		Expect(kubeClient.Create(ctx, testNamespace)).To(Succeed())
	})

	_ = AfterEach(func(ctx context.Context) {
		if kubeClient == nil {
			return
		}
		Expect(kubeClient.Delete(ctx, testNamespace)).To(Succeed())

		Eventually(func(g Gomega) {
			g.Expect(kubeClient.Get(ctx, client.ObjectKeyFromObject(testNamespace), testNamespace)).To(Succeed())
		}).WithTimeout(time.Minute).WithPolling(time.Second).To(Succeed())
	})

	It("Should upgrade a Flex cluster to a Dedicated cluster", func(ctx context.Context) {
		By("Create a secret with Atlas credentials", func() {
			k8s.CreateDefaultSecret(ctx, kubeClient, resourcePrefix+"-secret", testNamespace.Name)
		})

		By("Create Atlas Project", func() {
			project := akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      resourcePrefix + "-project",
					Namespace: testNamespace.Name,
				},
				Spec: akov2.AtlasProjectSpec{
					Name: resourcePrefix + "-project",
					ConnectionSecret: &akov2common.ResourceRefNamespaced{
						Name: resourcePrefix + "-secret",
					},
				},
			}

			Expect(kubeClient.Create(ctx, &project)).To(Succeed())

			Eventually(func(g Gomega) {
				condition, err := k8s.GetProjectStatusCondition(ctx, kubeClient, api.ReadyType, testNamespace.Name, resourcePrefix+"-project")
				g.Expect(err).To(BeNil())
				g.Expect(condition).To(Equal("True"))
			}).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
		})

		By("Create a Flex cluster", func() {
			flexCluster := akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      resourcePrefix + "-cluster",
					Namespace: testNamespace.Name,
				},
				Spec: akov2.AtlasDeploymentSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &akov2common.ResourceRefNamespaced{
							Name: resourcePrefix + "-project",
						},
					},
					FlexSpec: &akov2.FlexSpec{
						Name: resourcePrefix + "-cluster",
						ProviderSettings: &akov2.FlexProviderSettings{
							BackingProviderName: "AWS",
							RegionName:          "US_EAST_2",
						},
					},
				},
			}

			Expect(kubeClient.Create(ctx, &flexCluster)).To(Succeed())
			Eventually(func(g Gomega) {
				condition, err := k8s.GetDeploymentStatusCondition(ctx, kubeClient, api.ReadyType, testNamespace.Name, resourcePrefix+"-cluster")
				g.Expect(err).To(BeNil())
				g.Expect(condition).To(Equal("True"))
			}).WithTimeout(10 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
		})

		By("Upgrade Flex cluster to Dedicated cluster", func() {
			var deployment akov2.AtlasDeployment
			Expect(kubeClient.Get(ctx, client.ObjectKey{Namespace: testNamespace.Name, Name: resourcePrefix + "-cluster"}, &deployment)).To(Succeed())

			deployment.Spec.UpgradeToDedicated = true
			deployment.Spec.FlexSpec = nil
			deployment.Spec.DeploymentSpec = &akov2.AdvancedDeploymentSpec{
				Name:        resourcePrefix + "-cluster",
				ClusterType: "REPLICASET",
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{
								ProviderName: "AWS",
								RegionName:   "EU_CENTRAL_1",
								Priority:     pointer.MakePtr(7),
								ElectableSpecs: &akov2.Specs{
									InstanceSize: "M30",
									NodeCount:    pointer.MakePtr(3),
								},
							},
						},
					},
				},
			}

			Expect(kubeClient.Update(ctx, &deployment)).To(Succeed())
			Eventually(func(g Gomega) {
				updatedDeployment := akov2.AtlasDeployment{}
				g.Expect(kubeClient.Get(ctx, client.ObjectKey{Namespace: testNamespace.Name, Name: resourcePrefix + "-cluster"}, &updatedDeployment)).To(Succeed())
				for _, c := range updatedDeployment.GetStatus().GetConditions() {
					switch c.Type {
					case "Ready":
						g.Expect(c.Status).To(Equal(corev1.ConditionTrue))
						g.Expect(c.Message).To(Equal("Cluster is already dedicated. It’s recommended to remove or set the upgrade flag to false"))
					default:
						g.Expect(c.Status).To(Equal(corev1.ConditionTrue))
					}
				}
			}).WithTimeout(30 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
		})

		By("Delete resources", func() {
			var deployment akov2.AtlasDeployment
			Expect(kubeClient.Get(ctx, client.ObjectKey{Namespace: testNamespace.Name, Name: resourcePrefix + "-cluster"}, &deployment)).To(Succeed())
			Expect(kubeClient.Delete(ctx, &deployment)).To(Succeed())

			var project akov2.AtlasProject
			Expect(kubeClient.Get(ctx, client.ObjectKey{Namespace: testNamespace.Name, Name: resourcePrefix + "-project"}, &project)).To(Succeed())
			Expect(kubeClient.Delete(ctx, &project)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(kubeClient.Get(ctx, client.ObjectKey{Namespace: testNamespace.Name, Name: resourcePrefix + "-project"}, &project)).ToNot(Succeed())
			}).WithTimeout(30 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
		})
	})

	DescribeTable(
		"Should handle invalid upgrade scenarios",
		func(ctx context.Context, project *akov2.AtlasProject, deployment *akov2.AtlasDeployment, updateDeployment *akov2.AtlasDeployment, errorMessage string) {
			By("Create a secret with Atlas credentials", func() {
				k8s.CreateDefaultSecret(ctx, kubeClient, resourcePrefix+"-secret", testNamespace.Name)
			})

			By("Create Atlas Project", func() {
				project.WithName(resourcePrefix + "-project")
				project.WithAtlasName(resourcePrefix + "-project")
				project.WithConnectionSecret(resourcePrefix + "-secret")
				project.Namespace = testNamespace.Name

				Expect(kubeClient.Create(ctx, project)).To(Succeed())

				Eventually(func(g Gomega) {
					condition, err := k8s.GetProjectStatusCondition(ctx, kubeClient, api.ReadyType, testNamespace.Name, resourcePrefix+"-project")
					g.Expect(err).To(BeNil())
					g.Expect(condition).To(Equal("True"))
				}).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})

			By("Create a cluster", func() {
				deployment.WithName(resourcePrefix + "-cluster")
				deployment.WithAtlasName(resourcePrefix + "-cluster")
				deployment.WithProjectName(resourcePrefix + "-project")
				deployment.Namespace = testNamespace.Name

				Expect(kubeClient.Create(ctx, deployment)).To(Succeed())
				Eventually(func(g Gomega) {
					condition, err := k8s.GetDeploymentStatusCondition(ctx, kubeClient, api.ReadyType, testNamespace.Name, resourcePrefix+"-cluster")
					g.Expect(err).To(BeNil())
					g.Expect(condition).To(Equal("True"))
				}).WithTimeout(10 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})

			By("Upgrade cluster to Dedicated cluster", func() {
				Expect(kubeClient.Get(ctx, client.ObjectKey{Namespace: testNamespace.Name, Name: resourcePrefix + "-cluster"}, deployment)).To(Succeed())

				updateDeployment.ObjectMeta = deployment.ObjectMeta
				updateDeployment.WithAtlasName(resourcePrefix + "-cluster")
				updateDeployment.WithProjectName(resourcePrefix + "-project")
				updateDeployment.Namespace = testNamespace.Name
				updateDeployment.Spec.UpgradeToDedicated = true
				updateDeployment.Spec.DeploymentSpec.ClusterType = "SHARDED"

				Expect(kubeClient.Update(ctx, updateDeployment)).To(Succeed())
				Eventually(func(g Gomega) {
					updatedDeployment := akov2.AtlasDeployment{}
					g.Expect(kubeClient.Get(ctx, client.ObjectKey{Namespace: testNamespace.Name, Name: resourcePrefix + "-cluster"}, &updatedDeployment)).To(Succeed())
					for _, c := range updatedDeployment.GetStatus().GetConditions() {
						if c.Type == "DeploymentReady" {
							g.Expect(c.Status).To(Equal(corev1.ConditionFalse))
							g.Expect(c.Reason).To(Equal("DedicatedMigrationFailed"))
							g.Expect(c.Message).To(ContainSubstring(errorMessage))
						}
					}
				}).WithTimeout(30 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
			})

			By("Delete resources", func() {
				Expect(kubeClient.Get(ctx, client.ObjectKey{Namespace: testNamespace.Name, Name: resourcePrefix + "-cluster"}, deployment)).To(Succeed())
				Expect(kubeClient.Delete(ctx, deployment)).To(Succeed())

				Expect(kubeClient.Get(ctx, client.ObjectKey{Namespace: testNamespace.Name, Name: resourcePrefix + "-project"}, project)).To(Succeed())
				Expect(kubeClient.Delete(ctx, project)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(kubeClient.Get(ctx, client.ObjectKey{Namespace: testNamespace.Name, Name: resourcePrefix + "-project"}, project)).ToNot(Succeed())
				}).WithTimeout(30 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
			})
		},
		Entry(
			"Cannot upgrade a shared cluster to dedicated",
			utils.ParseObjectFromYAMLFile("../helper/e2e/data/atlasproject.yaml", &akov2.AtlasProject{}),
			utils.ParseObjectFromYAMLFile("../helper/e2e/data/atlasdeployment_basic_free.yaml", &akov2.AtlasDeployment{}),
			utils.ParseObjectFromYAMLFile("../helper/e2e/data/atlasdeployment_standard.yaml", &akov2.AtlasDeployment{}),
			"failed to upgrade cluster: upgrade from shared to dedicated is not supported",
		),
		Entry(
			"Cannot upgrade a flex cluster to dedicated with wrong spec",
			utils.ParseObjectFromYAMLFile("../helper/e2e/data/atlasproject.yaml", &akov2.AtlasProject{}),
			utils.ParseObjectFromYAMLFile("../helper/e2e/data/atlasdeployment_flex.yaml", &akov2.AtlasDeployment{}),
			utils.ParseObjectFromYAMLFile("../helper/e2e/data/atlasdeployment_standard.yaml", &akov2.AtlasDeployment{}),
			"Cannot upgrade a shared-tier cluster to a sharded cluster. Please upgrade to a dedicated replica set before converting to a sharded cluster",
		),
	)
})
