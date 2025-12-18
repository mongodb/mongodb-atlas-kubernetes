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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	nextapiv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/version"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e2/flexsamples"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/operator"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/resources"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/testparams"
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
	var sharedGroupNamespace *corev1.Namespace
	var testGroup *nextapiv1.Group
	var groupID string
	var orgID string
	var sharedTestParams *testparams.TestParams

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
		Expect(kube.AssertCRDNames(ctx, kubeClient, FlexClusterCRDName, GroupCRDName)).To(Succeed())

		By("Create namespace and credentials for shared test Group", func() {
			sharedGroupNamespace = &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{
				Name: utils.RandomName("flex-shared-grp-ns"),
			}}
			Expect(kubeClient.Create(ctx, sharedGroupNamespace)).To(Succeed())
			Expect(resources.CopyCredentialsToNamespace(ctx, kubeClient, DefaultGlobalCredentials, control.MustEnvVar("OPERATOR_NAMESPACE"), sharedGroupNamespace.Name, GinkGoFieldOwner)).To(Succeed())
		})

		By("Create test Group", func() {
			groupName := utils.RandomName("flexcluster-test-group")
			// Set up shared test params for Group YAML template
			sharedTestParams = testparams.New(orgID, sharedGroupNamespace.Name, DefaultGlobalCredentials).
				WithGroupName(groupName)
			// Replace placeholders in the Group YAML template
			groupYAML := sharedTestParams.ReplaceYAML(string(flexsamples.TestGroup))
			objs := yml.MustParseObjects(strings.NewReader(groupYAML))
			Expect(len(objs)).To(Equal(1))
			testGroup = objs[0].(*nextapiv1.Group)
			Expect(kubeClient.Create(ctx, testGroup)).To(Succeed())
		})

		By("Wait for Group to be Ready and get its ID", func() {
			Eventually(func(g Gomega) {
				g.Expect(resources.CheckResourceReady(ctx, kubeClient, testGroup)).To(Succeed())
			}).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			Expect(testGroup.Status.V20250312).NotTo(BeNil())
			Expect(testGroup.Status.V20250312.Id).NotTo(BeNil())
			groupID = *testGroup.Status.V20250312.Id
			Expect(groupID).NotTo(BeEmpty())
			// Update shared test params with groupID now that it's available
			sharedTestParams = sharedTestParams.WithGroupID(groupID)
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
		if kubeClient != nil && sharedGroupNamespace != nil {
			By("Clean up shared group namespace", func() {
				Expect(kubeClient.Delete(ctx, sharedGroupNamespace)).To(Succeed())
				Eventually(func(g Gomega) bool {
					return kubeClient.Get(ctx, client.ObjectKeyFromObject(sharedGroupNamespace), sharedGroupNamespace) == nil
				}).WithTimeout(time.Minute).WithPolling(time.Second).To(BeFalse())
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
			// Generate randomized group name for this test run (cluster names are unique per group)
			groupName := utils.RandomName("flex-grp")

			// Set up test params for this test case (reuse shared values, override groupName)
			testParams := sharedTestParams.WithGroupName(groupName)

			// Track created objects for cleanup
			var createdObjects []client.Object

			By("Copy credentials secret to test namespace", func() {
				Expect(resources.CopyCredentialsToNamespace(ctx, kubeClient, DefaultGlobalCredentials, control.MustEnvVar("OPERATOR_NAMESPACE"), testNamespace.Name, GinkGoFieldOwner)).To(Succeed())
			})

			By("Create resources from YAML", func() {
				objs, err := resources.ApplyYAMLToNamespace(ctx, kubeClient, createYAML, testParams, testNamespace.Name, GinkGoFieldOwner)
				Expect(err).NotTo(HaveOccurred())
				createdObjects = append(createdObjects, objs...)
			})

			By("Wait for Group to be Ready (if using groupRef)", func() {
				createYAMLStr := testParams.ReplaceYAML(string(createYAML))
				objs := yml.MustParseObjects(strings.NewReader(createYAMLStr))
				for _, obj := range objs {
					if group, ok := obj.(*nextapiv1.Group); ok {
						groupObj := &nextapiv1.Group{
							ObjectMeta: metav1.ObjectMeta{Name: group.Name, Namespace: testNamespace.Name},
						}
						Eventually(func(g Gomega) {
							g.Expect(resources.CheckResourceReady(ctx, kubeClient, groupObj)).To(Succeed())
						}).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
					}
				}
			})

			cluster := nextapiv1.FlexCluster{
				ObjectMeta: metav1.ObjectMeta{Name: clusterName, Namespace: testNamespace.Name},
			}

			By("Wait for FlexCluster to be Ready", func() {
				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, &cluster)).To(Succeed())
				}).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})

			By("Verify cluster was created", func() {
				Expect(cluster.Status.V20250312).NotTo(BeNil())
				Expect(cluster.Status.V20250312.Id).NotTo(BeNil())
				Expect(*cluster.Status.V20250312.Id).NotTo(BeEmpty())
			})

			By("Update FlexCluster", func() {
				if len(updateYAML) > 0 {
					_, err := resources.ApplyYAMLToNamespace(ctx, kubeClient, updateYAML, testParams, testNamespace.Name, GinkGoFieldOwner)
					Expect(err).NotTo(HaveOccurred())
				}
			})

			By("Wait for FlexCluster to be Ready & updated", func() {
				if len(updateYAML) > 0 {
					Eventually(func(g Gomega) {
						g.Expect(resources.CheckResourceUpdated(ctx, kubeClient, &cluster)).To(Succeed())
					}).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
				}
			})

			By("Delete all created resources", func() {
				for _, obj := range createdObjects {
					_ = kubeClient.Delete(ctx, obj)
				}
			})

			By("Wait for all resources to be deleted", func() {
				for _, obj := range createdObjects {
					Eventually(func(g Gomega) {
						g.Expect(resources.CheckResourceDeleted(ctx, kubeClient, obj)).To(Succeed())
					}).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
				}
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
