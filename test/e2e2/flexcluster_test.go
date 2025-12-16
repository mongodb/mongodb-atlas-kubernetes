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
	"context"
	"os"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	nextapiv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/version"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e2/flexsamples"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/operator"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/yml"
)

const (
	FlexClusterCRDName = "flexclusters.atlas.generated.mongodb.com"
	GroupCRDName       = "groups.atlas.generated.mongodb.com"
)

var _ = Describe("FlexCluster CRUD", Ordered, Label("flexcluster-ctlr"), func() {
	var ctx context.Context
	var kubeClient client.Client
	var ako operator.Operator
	var testNamespace *corev1.Namespace
	var testGroup *nextapiv1.Group
	var groupID string
	var orgID string

	_ = BeforeAll(func() {
		if !version.IsExperimental() {
			Skip("FlexCluster is an experimental CRD and controller. Skipping test as experimental features are not enabled.")
		}

		orgID = os.Getenv("MCLI_ORG_ID")
		Expect(orgID).NotTo(BeEmpty(), "MCLI_ORG_ID environment variable must be set")

		deletionProtectionOff := false
		ako = runTestAKO(DefaultGlobalCredentials, control.MustEnvVar("OPERATOR_NAMESPACE"), deletionProtectionOff)
		ako.Start(GinkgoT())

		ctx = context.Background()
		testClient, err := kube.NewTestClient()
		Expect(err).To(Succeed())
		kubeClient = testClient
		Expect(kube.AssertCRDs(ctx, kubeClient,
			&apiextensionsv1.CustomResourceDefinition{
				ObjectMeta: v1.ObjectMeta{Name: FlexClusterCRDName},
			},
			&apiextensionsv1.CustomResourceDefinition{
				ObjectMeta: v1.ObjectMeta{Name: GroupCRDName},
			},
		)).To(Succeed())

		By("Create test Group", func() {
			operatorNamespace := control.MustEnvVar("OPERATOR_NAMESPACE")
			groupName := utils.RandomName("flexcluster-test-group")
			// Replace placeholders in the Group YAML template
			groupYAML := strings.ReplaceAll(string(flexsamples.TestGroup), "__GROUP_NAME__", groupName)
			groupYAML = strings.ReplaceAll(groupYAML, "__OPERATOR_NAMESPACE__", operatorNamespace)
			groupYAML = strings.ReplaceAll(groupYAML, "__CREDENTIALS_SECRET_NAME__", DefaultGlobalCredentials)
			groupYAML = strings.ReplaceAll(groupYAML, "__ORG_ID__", orgID)
			objs := yml.MustParseObjects(strings.NewReader(groupYAML))
			Expect(len(objs)).To(Equal(1))
			testGroup = objs[0].(*nextapiv1.Group)
			Expect(kubeClient.Create(ctx, testGroup)).To(Succeed())
		})

		By("Wait for Group to be Ready and get its ID", func() {
			Eventually(func(g Gomega) bool {
				g.Expect(
					kubeClient.Get(ctx, client.ObjectKeyFromObject(testGroup), testGroup),
				).To(Succeed())
				if condition := meta.FindStatusCondition(testGroup.GetConditions(), "Ready"); condition != nil {
					if condition.Status == metav1.ConditionTrue {
						if testGroup.Status.V20250312 != nil && testGroup.Status.V20250312.Id != nil {
							groupID = *testGroup.Status.V20250312.Id
							return true
						}
					}
				}
				return false
			}).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).To(BeTrue())
			Expect(groupID).NotTo(BeEmpty())
		})
	})

	_ = AfterAll(func() {
		if kubeClient != nil && testGroup != nil {
			By("Clean up test Group", func() {
				Expect(kubeClient.Delete(ctx, testGroup)).To(Succeed())
				Eventually(func(g Gomega) error {
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testGroup), testGroup)
					return err
				}).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).NotTo(Succeed())
			})
		}
		if ako != nil {
			ako.Stop(GinkgoT())
		}
	})

	_ = BeforeEach(func() {
		testNamespace = &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{
			Name: utils.RandomName("flexcluster-ctlr-ns"),
		}}
		Expect(kubeClient.Create(ctx, testNamespace)).To(Succeed())
		Expect(ako.Running()).To(BeTrue(), "Operator must be running")
	})

	_ = AfterEach(func() {
		if kubeClient == nil {
			return
		}
		Expect(
			kubeClient.Delete(ctx, testNamespace),
		).To(Succeed())
		Eventually(func(g Gomega) bool {
			return kubeClient.Get(ctx, client.ObjectKeyFromObject(testNamespace), testNamespace) == nil
		}).WithTimeout(time.Minute).WithPolling(time.Second).To(BeFalse())
	})

	DescribeTable("FlexCluster CRUD lifecycle",
		func(createYAML, updateYAML []byte, clusterName string) {
			By("Copy credentials secret to test namespace", func() {
				globalCredsKey := client.ObjectKey{
					Name:      DefaultGlobalCredentials,
					Namespace: control.MustEnvVar("OPERATOR_NAMESPACE"),
				}
				credentialsSecret, err := copySecretToNamespace(ctx, kubeClient, globalCredsKey, testNamespace.Name)
				Expect(err).NotTo(HaveOccurred())
				Expect(
					kubeClient.Patch(ctx, credentialsSecret, client.Apply, client.ForceOwnership, GinkGoFieldOwner),
				).To(Succeed())
			})

			By("Create resources from YAML", func() {
				// Replace placeholders with actual values
				createYAMLStr := strings.ReplaceAll(string(createYAML), "__GROUP_ID__", groupID)
				createYAMLStr = strings.ReplaceAll(createYAMLStr, "__ORG_ID__", orgID)
				objs := yml.MustParseObjects(strings.NewReader(createYAMLStr))
				for _, obj := range objs {
					objToApply := kube.WithRenamedNamespace(obj, testNamespace.Name)
					Expect(
						kubeClient.Patch(ctx, objToApply, client.Apply, client.ForceOwnership, GinkGoFieldOwner),
					).To(Succeed())
				}
			})

			By("Wait for Group to be Ready (if using groupRef)", func() {
				createYAMLStr := strings.ReplaceAll(string(createYAML), "__GROUP_ID__", groupID)
				createYAMLStr = strings.ReplaceAll(createYAMLStr, "__ORG_ID__", orgID)
				objs := yml.MustParseObjects(strings.NewReader(createYAMLStr))
				for _, obj := range objs {
					if group, ok := obj.(*nextapiv1.Group); ok {
						groupInKube := nextapiv1.Group{
							ObjectMeta: metav1.ObjectMeta{Name: group.Name, Namespace: testNamespace.Name},
						}
						Eventually(func(g Gomega) bool {
							g.Expect(
								kubeClient.Get(ctx, client.ObjectKeyFromObject(&groupInKube), &groupInKube),
							).To(Succeed())
							if condition := meta.FindStatusCondition(groupInKube.GetConditions(), "Ready"); condition != nil {
								return condition.Status == metav1.ConditionTrue
							}
							return false
						}).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).To(BeTrue())
					}
				}
			})

			cluster := nextapiv1.FlexCluster{
				ObjectMeta: metav1.ObjectMeta{Name: clusterName, Namespace: testNamespace.Name},
			}

			By("Wait for FlexCluster to be Ready", func() {
				Eventually(func(g Gomega) bool {
					g.Expect(
						kubeClient.Get(ctx, client.ObjectKeyFromObject(&cluster), &cluster),
					).To(Succeed())
					if condition := meta.FindStatusCondition(cluster.GetConditions(), "Ready"); condition != nil {
						return condition.Status == metav1.ConditionTrue
					}
					return false
				}).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).To(BeTrue())
			})

			By("Verify cluster was created", func() {
				Expect(cluster.Status.V20250312).NotTo(BeNil())
				Expect(cluster.Status.V20250312.Id).NotTo(BeNil())
				Expect(*cluster.Status.V20250312.Id).NotTo(BeEmpty())
			})

			By("Update FlexCluster", func() {
				if len(updateYAML) > 0 {
					// Replace placeholders with actual values
					updateYAMLStr := strings.ReplaceAll(string(updateYAML), "__GROUP_ID__", groupID)
					updateYAMLStr = strings.ReplaceAll(updateYAMLStr, "__ORG_ID__", orgID)
					updateObjs := yml.MustParseObjects(strings.NewReader(updateYAMLStr))
					for _, obj := range updateObjs {
						objToPatch := kube.WithRenamedNamespace(obj, testNamespace.Name)
						Expect(
							kubeClient.Patch(ctx, objToPatch, client.Apply, client.ForceOwnership, GinkGoFieldOwner),
						).To(Succeed())
					}
				}
			})

			By("Wait for FlexCluster to be Ready & updated", func() {
				if len(updateYAML) > 0 {
					Eventually(func(g Gomega) bool {
						g.Expect(
							kubeClient.Get(ctx, client.ObjectKeyFromObject(&cluster), &cluster),
						).To(Succeed())
						ready := false
						if condition := meta.FindStatusCondition(cluster.GetConditions(), "Ready"); condition != nil {
							ready = (condition.Status == metav1.ConditionTrue)
						}
						if ready {
							if condition := meta.FindStatusCondition(cluster.GetConditions(), "State"); condition != nil {
								return state.ResourceState(condition.Reason) == state.StateUpdated
							}
						}
						return false
					}).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).To(BeTrue())
				}
			})

			By("Delete FlexCluster", func() {
				Expect(kubeClient.Delete(ctx, &cluster)).To(Succeed())
			})

			By("Wait for FlexCluster to be deleted", func() {
				Eventually(func(g Gomega) error {
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(&cluster), &cluster)
					return err
				}).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).NotTo(Succeed())
			})
		},
		Entry("With direct groupId",
			flexsamples.WithGroupIdCreate,
			flexsamples.WithGroupIdUpdate,
			"flexy",
		),
		Entry("With groupRef",
			flexsamples.WithGroupRefCreate,
			flexsamples.WithGroupRefUpdate,
			"flexy",
		),
	)
})
