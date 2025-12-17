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
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// yamlPlaceholders holds all placeholder values for YAML template replacement.
type yamlPlaceholders struct {
	GroupID               string
	OrgID                 string
	GroupName             string
	OperatorNamespace     string
	CredentialsSecretName string
}

var _ = Describe("FlexCluster CRUD", Ordered, Label("flexcluster-ctlr"), func() {
	var ctx context.Context
	var kubeClient client.Client
	var ako operator.Operator
	var testNamespace *corev1.Namespace
	var sharedGroupNamespace *corev1.Namespace
	var testGroup *nextapiv1.Group
	var groupID string
	var orgID string
	var sharedPlaceholders yamlPlaceholders

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
			copyCredentialsToNamespace(ctx, kubeClient, sharedGroupNamespace.Name)
		})

		By("Create test Group", func() {
			groupName := utils.RandomName("flexcluster-test-group")
			// Set up shared placeholders for Group YAML template
			sharedPlaceholders = yamlPlaceholders{
				GroupName:             groupName,
				OperatorNamespace:     sharedGroupNamespace.Name,
				CredentialsSecretName: DefaultGlobalCredentials,
				OrgID:                 orgID,
			}
			// Replace placeholders in the Group YAML template
			groupYAML := replaceYAMLPlaceholders(string(flexsamples.TestGroup), sharedPlaceholders)
			objs := yml.MustParseObjects(strings.NewReader(groupYAML))
			Expect(len(objs)).To(Equal(1))
			testGroup = objs[0].(*nextapiv1.Group)
			Expect(kubeClient.Create(ctx, testGroup)).To(Succeed())
		})

		By("Wait for Group to be Ready and get its ID", func() {
			waitForResourceReady(ctx, kubeClient, testGroup)
			Expect(testGroup.Status.V20250312).NotTo(BeNil())
			Expect(testGroup.Status.V20250312.Id).NotTo(BeNil())
			groupID = *testGroup.Status.V20250312.Id
			Expect(groupID).NotTo(BeEmpty())
			// Update shared placeholders with groupID now that it's available
			sharedPlaceholders.GroupID = groupID
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

			// Set up placeholders for this test case (reuse shared values, override groupName)
			testPlaceholders := sharedPlaceholders
			testPlaceholders.GroupName = groupName

			// Track created objects for cleanup
			var createdObjects []client.Object

			By("Copy credentials secret to test namespace", func() {
				copyCredentialsToNamespace(ctx, kubeClient, testNamespace.Name)
			})

			By("Create resources from YAML", func() {
				objs := applyYAMLToNamespace(ctx, kubeClient, createYAML, testPlaceholders, testNamespace.Name)
				createdObjects = append(createdObjects, objs...)
			})

			By("Wait for Group to be Ready (if using groupRef)", func() {
				createYAMLStr := replaceYAMLPlaceholders(string(createYAML), testPlaceholders)
				objs := yml.MustParseObjects(strings.NewReader(createYAMLStr))
				for _, obj := range objs {
					if group, ok := obj.(*nextapiv1.Group); ok {
						waitForResourceReady(ctx, kubeClient, &nextapiv1.Group{
							ObjectMeta: metav1.ObjectMeta{Name: group.Name, Namespace: testNamespace.Name},
						})
					}
				}
			})

			cluster := nextapiv1.FlexCluster{
				ObjectMeta: metav1.ObjectMeta{Name: clusterName, Namespace: testNamespace.Name},
			}

			By("Wait for FlexCluster to be Ready", func() {
				waitForResourceReady(ctx, kubeClient, &cluster)
			})

			By("Verify cluster was created", func() {
				Expect(cluster.Status.V20250312).NotTo(BeNil())
				Expect(cluster.Status.V20250312.Id).NotTo(BeNil())
				Expect(*cluster.Status.V20250312.Id).NotTo(BeEmpty())
			})

			By("Update FlexCluster", func() {
				if len(updateYAML) > 0 {
					applyYAMLToNamespace(ctx, kubeClient, updateYAML, testPlaceholders, testNamespace.Name)
				}
			})

			By("Wait for FlexCluster to be Ready & updated", func() {
				if len(updateYAML) > 0 {
					waitForResourceUpdated(ctx, kubeClient, &cluster)
				}
			})

			By("Delete all created resources", func() {
				for _, obj := range createdObjects {
					_ = kubeClient.Delete(ctx, obj)
				}
			})

			By("Wait for all resources to be deleted", func() {
				for _, obj := range createdObjects {
					waitForResourceDeleted(ctx, kubeClient, obj)
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

// replaceYAMLPlaceholders replaces placeholders in YAML templates with actual values from the struct.
func replaceYAMLPlaceholders(yaml string, p yamlPlaceholders) string {
	result := yaml
	result = strings.ReplaceAll(result, "__GROUP_ID__", p.GroupID)
	result = strings.ReplaceAll(result, "__ORG_ID__", p.OrgID)
	result = strings.ReplaceAll(result, "__GROUP_NAME__", p.GroupName)
	result = strings.ReplaceAll(result, "__OPERATOR_NAMESPACE__", p.OperatorNamespace)
	result = strings.ReplaceAll(result, "__CREDENTIALS_SECRET_NAME__", p.CredentialsSecretName)
	return result
}

// copyCredentialsToNamespace copies the default global credentials secret to the specified namespace.
func copyCredentialsToNamespace(ctx context.Context, kubeClient client.Client, namespace string) {
	globalCredsKey := client.ObjectKey{
		Name:      DefaultGlobalCredentials,
		Namespace: control.MustEnvVar("OPERATOR_NAMESPACE"),
	}
	credentialsSecret, err := copySecretToNamespace(ctx, kubeClient, globalCredsKey, namespace)
	Expect(err).NotTo(HaveOccurred())
	Expect(
		kubeClient.Patch(ctx, credentialsSecret, client.Apply, client.ForceOwnership, GinkGoFieldOwner),
	).To(Succeed())
}

// applyYAMLToNamespace applies YAML objects to a namespace after replacing placeholders.
// Returns the list of applied objects.
func applyYAMLToNamespace(ctx context.Context, kubeClient client.Client, yaml []byte, placeholders yamlPlaceholders, namespace string) []client.Object {
	yamlStr := replaceYAMLPlaceholders(string(yaml), placeholders)
	objs := yml.MustParseObjects(strings.NewReader(yamlStr))
	for _, obj := range objs {
		obj.SetNamespace(namespace)
		Expect(
			kubeClient.Patch(ctx, obj, client.Apply, client.ForceOwnership, GinkGoFieldOwner),
		).To(Succeed())
	}
	return objs
}

// waitForResourceReady waits for a resource to have Ready condition set to True.
func waitForResourceReady(ctx context.Context, kubeClient client.Client, obj kube.ObjectWithStatus) {
	Eventually(func(g Gomega) bool {
		g.Expect(
			kubeClient.Get(ctx, client.ObjectKeyFromObject(obj), obj),
		).To(Succeed())
		if condition := meta.FindStatusCondition(obj.GetConditions(), "Ready"); condition != nil {
			return condition.Status == metav1.ConditionTrue
		}
		return false
	}).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).To(BeTrue())
}

// waitForResourceUpdated waits for a resource to be Ready and in Updated state.
func waitForResourceUpdated(ctx context.Context, kubeClient client.Client, obj kube.ObjectWithStatus) {
	Eventually(func(g Gomega) bool {
		g.Expect(
			kubeClient.Get(ctx, client.ObjectKeyFromObject(obj), obj),
		).To(Succeed())
		ready := false
		if condition := meta.FindStatusCondition(obj.GetConditions(), "Ready"); condition != nil {
			ready = (condition.Status == metav1.ConditionTrue)
		}
		if ready {
			if condition := meta.FindStatusCondition(obj.GetConditions(), "State"); condition != nil {
				return state.ResourceState(condition.Reason) == state.StateUpdated
			}
		}
		return false
	}).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).To(BeTrue())
}

// waitForResourceDeleted waits for a resource to be deleted from the cluster.
func waitForResourceDeleted(ctx context.Context, kubeClient client.Client, obj client.Object) {
	Eventually(func(g Gomega) error {
		return kubeClient.Get(ctx, client.ObjectKeyFromObject(obj), obj)
	}).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).ShouldNot(Succeed())
}
