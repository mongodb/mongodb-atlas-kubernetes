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
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
	akoretry "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/retry"
)

var _ = Describe("UserLogin", Label("global-deployment"), func() {
	var testData *model.TestDataProvider

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
			Expect(actions.SaveDeploymentsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Delete Resources", func() {
			actions.DeleteTestDataDeployments(testData)
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(ctx SpecContext, test func(context.Context) *model.TestDataProvider, mapping []akov2.CustomZoneMapping, ns []akov2.ManagedNamespace) {
			testData = test(ctx)
			actions.ProjectCreationFlow(testData)
			globalClusterFlow(testData, mapping, ns)
		},
		Entry("Test[gc-advanced-deployment]: Advanced", Label("focus-gc-advanced-deployment"),
			func(ctx context.Context) *model.TestDataProvider {
				return model.DataProvider(ctx, "gc-advanced-deployment", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject()).WithInitialDeployments(data.CreateAdvancedGeoshardedDeployment("gc-advanced-deployment"))
			},
			[]akov2.CustomZoneMapping{
				{
					Zone:     "Zone 1",
					Location: "AO",
				},
				{
					Zone:     "Zone 2",
					Location: "CA",
				},
			},
			[]akov2.ManagedNamespace{
				{
					Collection:             "somecollection",
					Db:                     "somedb",
					CustomShardKey:         "somekey",
					PresplitHashedZones:    new(true),
					IsCustomShardKeyHashed: new(true),
					NumInitialChunks:       4,
				},
			},
		),
	)
})

// Regression test for AUTO_SCALINGS_MUST_MATCH (HTTP 400).
// Atlas requires all autoScaling (and analyticsAutoScaling) objects to match across all
// regionConfigs in a cluster. When the operator adds a new region while autoscaling is
// enabled on existing regions, it must send consistent autoScaling values in the PATCH
// body — otherwise Atlas rejects the request.
var _ = Describe("AtlasDeployment autoscaling region add", Label("deployment-autoscaling-region-add"), func() {
	var testData *model.TestDataProvider

	AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
			Expect(actions.SaveDeploymentsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Delete Resources", func() {
			actions.DeleteTestDataDeployments(testData)
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	BeforeEach(func(ctx SpecContext) {
		testData = model.DataProvider(
			ctx,
			"autoscaling-region-add",
			model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
			40000,
			[]func(*model.TestDataProvider){},
		).WithProject(data.DefaultProject())
		actions.ProjectCreationFlow(testData)
	})

	It("Should add a readOnly region to a cluster with autoscaling without AUTO_SCALINGS_MUST_MATCH",
		Label("focus-autoscaling-region-add"),
		func(ctx context.Context) {
			name := "autoscaling-region-add"
			d := data.CreateDeploymentWithAutoscaling(name)
			d.Namespace = testData.Resources.Namespace
			d.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].ElectableSpecs.NodeCount = new(3)

			By("Creating the initial deployment with one region and autoscaling enabled", func() {
				Expect(testData.K8SClient.Create(testData.Context, d)).To(Succeed())
				Eventually(func(g Gomega) bool {
					g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{
						Name:      d.Name,
						Namespace: d.Namespace,
					}, d)).To(Succeed())
					return d.Status.StateName == status.StateIDLE
				}).WithTimeout(30 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})

			By("Adding a second readOnly region with the same autoscaling settings", func() {
				_, err := akoretry.RetryUpdateOnConflict(
					ctx,
					testData.K8SClient,
					client.ObjectKeyFromObject(d),
					func(dep *akov2.AtlasDeployment) {
						data.UpdateDeploymentAddReadOnlyRegion(dep)
					},
				)
				Expect(err).To(BeNil())
			})

			By("Waiting for the cluster to leave IDLE (operator picked up the change)", func() {
				Eventually(func(g Gomega) bool {
					g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{
						Name:      d.Name,
						Namespace: d.Namespace,
					}, d)).To(Succeed())
					return d.Status.StateName != status.StateIDLE
				}).WithTimeout(5 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})

			By("Waiting for the cluster to reach IDLE state after adding the second region", func() {
				Eventually(func(g Gomega) bool {
					g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{
						Name:      d.Name,
						Namespace: d.Namespace,
					}, d)).To(Succeed())
					for _, condition := range d.Status.Conditions {
						if condition.Type == api.DeploymentReadyType {
							g.Expect(condition.Status).To(Equal(corev1.ConditionTrue),
								"DeploymentReady condition is False: %s", condition.Message)
							return condition.Status == corev1.ConditionTrue
						}
					}
					return false
				}).WithTimeout(30 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})
		},
	)
})

func globalClusterFlow(userData *model.TestDataProvider, mapping []akov2.CustomZoneMapping, managedNamespace []akov2.ManagedNamespace) {
	By("Apply deployment", func() {
		Expect(userData.InitialDeployments).ShouldNot(BeEmpty())
		userData.InitialDeployments[0].Namespace = userData.Resources.Namespace
		Expect(userData.K8SClient.Create(userData.Context, userData.InitialDeployments[0])).To(Succeed())
		Eventually(func(g Gomega) bool {
			g.Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{
				Name:      userData.InitialDeployments[0].Name,
				Namespace: userData.InitialDeployments[0].Namespace,
			}, userData.InitialDeployments[0])).To(Succeed())
			return userData.InitialDeployments[0].Status.StateName == status.StateIDLE
		}).WithTimeout(30 * time.Minute).Should(BeTrue())
	})

	By("Applying global cluster config to Deployment", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{
			Name:      userData.InitialDeployments[0].Name,
			Namespace: userData.InitialDeployments[0].Namespace,
		}, userData.InitialDeployments[0])).To(Succeed())
		if userData.InitialDeployments[0].Spec.DeploymentSpec != nil {
			userData.InitialDeployments[0].Spec.DeploymentSpec.ManagedNamespaces = managedNamespace
			userData.InitialDeployments[0].Spec.DeploymentSpec.CustomZoneMapping = mapping
		}

		Expect(userData.K8SClient.Update(userData.Context, userData.InitialDeployments[0])).To(Succeed())
	})

	By("Wait and check global zone mapping status", func() {
		Eventually(func(g Gomega) bool {
			g.Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{
				Name:      userData.InitialDeployments[0].Name,
				Namespace: userData.InitialDeployments[0].Namespace,
			}, userData.InitialDeployments[0])).To(Succeed())
			for _, condition := range userData.InitialDeployments[0].Status.Conditions {
				if condition.Type == api.CustomZoneMappingReadyType {
					return condition.Status == corev1.ConditionTrue
				}
			}
			return false
		}).WithTimeout(10 * time.Minute).Should(BeTrue())
	})

	By("Wait and check global managed namespaces status", func() {
		Eventually(func(g Gomega) bool {
			g.Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{
				Name:      userData.InitialDeployments[0].Name,
				Namespace: userData.InitialDeployments[0].Namespace,
			}, userData.InitialDeployments[0])).To(Succeed())
			for _, condition := range userData.InitialDeployments[0].Status.Conditions {
				if condition.Type == api.ManagedNamespacesReadyType {
					return condition.Status == corev1.ConditionTrue
				}
			}
			return false
		}).WithTimeout(10 * time.Minute).Should(BeTrue())
	})

	By("Delete global cluster config and wait idle state of cluster", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{
			Name:      userData.InitialDeployments[0].Name,
			Namespace: userData.InitialDeployments[0].Namespace,
		}, userData.InitialDeployments[0])).To(Succeed())
		if userData.InitialDeployments[0].Spec.DeploymentSpec != nil {
			userData.InitialDeployments[0].Spec.DeploymentSpec.ManagedNamespaces = nil
			userData.InitialDeployments[0].Spec.DeploymentSpec.CustomZoneMapping = nil
		}
		Expect(userData.K8SClient.Update(userData.Context, userData.InitialDeployments[0])).To(Succeed())
		Eventually(func(g Gomega) bool {
			g.Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{
				Name:      userData.InitialDeployments[0].Name,
				Namespace: userData.InitialDeployments[0].Namespace,
			}, userData.InitialDeployments[0])).To(Succeed())
			for _, condition := range userData.InitialDeployments[0].Status.Conditions {
				if condition.Type == api.DeploymentReadyType {
					return condition.Status == corev1.ConditionTrue
				}
			}
			return false
		}).WithTimeout(30 * time.Minute).Should(BeTrue())
	})
}
