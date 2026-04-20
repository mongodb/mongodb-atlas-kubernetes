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
	"fmt"
	"net/http"
	"time"

	k8s "github.com/crd2go/crd2go/k8s"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas-sdk/v20250312018/admin"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/controller-runtime/pkg/client"

	generatedv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/generated/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/httputil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/operator"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/resources"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/samples"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/testparams"
)

const testCIDR = "203.0.113.0/24"

var _ = Describe("IPAccessList CRUD", Ordered, Label("ipaccesslist"), func() {
	var ctx context.Context
	var kubeClient client.Client
	var ako operator.Operator
	var testNamespace *corev1.Namespace
	var orgID string
	var atlasClient *admin.APIClient

	_ = BeforeAll(func() {
		ctx = context.Background()

		deletionProtectionOff := false
		ako = runTestAKO(DefaultGlobalCredentials, control.MustEnvVar("OPERATOR_NAMESPACE"), deletionProtectionOff)
		ako.Start(ctx, GinkgoT())

		DeferCleanup(func() {
			if ako != nil {
				ako.Stop(GinkgoT())
			}
		})

		testClient, err := kube.NewTestClient()
		Expect(err).To(Succeed())
		kubeClient = testClient
		Expect(kube.AssertCRDNames(ctx, kubeClient,
			"groups.atlas.generated.mongodb.com",
			"ipaccesslistentries.atlas.generated.mongodb.com",
		)).To(Succeed())

		atlasClient, orgID = newTestAtlasClient()
	})

	_ = BeforeEach(func() {
		testNamespace = &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("ipal-ns-%s", rand.String(6)),
		}}
		Expect(kubeClient.Create(ctx, testNamespace)).To(Succeed())
		Expect(ako.Running()).To(BeTrue(), "Operator must be running")
		Expect(resources.CopyCredentialsToNamespace(
			ctx,
			kubeClient,
			DefaultGlobalCredentials,
			control.MustEnvVar("OPERATOR_NAMESPACE"),
			testNamespace.Name, GinkGoFieldOwner),
		).To(Succeed())
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
		}).WithContext(ctx).WithTimeout(time.Minute).WithPolling(time.Second).To(BeFalse())
	})

	Describe("IPAccessList CRUD lifecycle", func() {
		It("Should create, update comment, and delete IPAccessList", Label("focus-ipaccesslist-crud"), func() {
			groupName := fmt.Sprintf("test-group-%s", rand.String(6))
			groupParams := testparams.New(orgID, control.MustEnvVar("OPERATOR_NAMESPACE"), DefaultGlobalCredentials).
				WithGroupName(groupName).
				WithNamespace(testNamespace.Name)

			var testGroup *generatedv1.Group
			By("Create prerequisite Group", func() {
				objs := samples.MustLoadSampleObjects("atlas_generated_v1_group.yaml")
				Expect(len(objs)).To(Equal(1))
				testGroup = objs[0].(*generatedv1.Group)
				applyTestParamsToGroup(testGroup, groupParams)
				Expect(kubeClient.Create(ctx, testGroup)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, testGroup)).To(Succeed())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
				Expect(testGroup.Status.V20250312).NotTo(BeNil())
				Expect(testGroup.Status.V20250312.Id).NotTo(BeNil())
			})

			entryName := fmt.Sprintf("ipal-%s", rand.String(6))
			var testIAL *generatedv1.IPAccessListEntry
			By("Create IPAccessList with cidrBlock via groupRef", func() {
				objs := samples.MustLoadSampleObjects("atlas_generated_v1_ipaccesslistentry_with_groupref.yaml")
				Expect(len(objs)).To(Equal(1))
				testIAL = objs[0].(*generatedv1.IPAccessListEntry)
				applyTestParamsToIPAccessList(testIAL, testNamespace.Name, entryName, testGroup.GetName())
				Expect(kubeClient.Create(ctx, testIAL)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, testIAL)).To(Succeed())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

				Expect(testIAL.Status.V20250312).NotTo(BeNil())
				Expect(testIAL.Status.V20250312.GroupId).NotTo(BeNil())
				Expect(meta.IsStatusConditionTrue(testIAL.GetConditions(), "Ready")).To(BeTrue())

				groupID := *testGroup.Status.V20250312.Id
				atlasEntry, _, err := atlasClient.ProjectIPAccessListApi.GetAccessListEntry(ctx, groupID, testCIDR).Execute()
				Expect(err).ToNot(HaveOccurred())
				Expect(atlasEntry.GetCidrBlock()).To(Equal(testCIDR))
				Expect(atlasEntry.GetComment()).To(Equal("e2e test entry"))
			})

			By("Update IPAccessList comment", func() {
				Expect(kubeClient.Get(ctx, client.ObjectKeyFromObject(testIAL), testIAL)).To(Succeed())
				testIAL.Spec.V20250312.Entry.Comment = pointer.MakePtr("updated comment")
				Expect(kubeClient.Update(ctx, testIAL)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceUpdated(ctx, kubeClient, testIAL)).To(Succeed())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

				groupID := *testGroup.Status.V20250312.Id
				atlasEntry, _, err := atlasClient.ProjectIPAccessListApi.GetAccessListEntry(ctx, groupID, testCIDR).Execute()
				Expect(err).ToNot(HaveOccurred())
				Expect(atlasEntry.GetComment()).To(Equal("updated comment"))
			})

			By("Delete IPAccessList from cluster - should delete from Atlas", func() {
				groupID := *testIAL.Status.V20250312.GroupId
				Expect(kubeClient.Delete(ctx, testIAL)).To(Succeed())

				Eventually(func(g Gomega) {
					ial := &generatedv1.IPAccessListEntry{}
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testIAL), ial)
					g.Expect(err).NotTo(Succeed())
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

				Eventually(func(g Gomega) {
					_, r, err := atlasClient.ProjectIPAccessListApi.GetAccessListEntry(ctx, groupID, testCIDR).Execute()
					g.Expect(err).ToNot(BeNil())
					g.Expect(httputil.StatusCode(r)).To(Equal(http.StatusNotFound))
				}).WithContext(ctx).WithTimeout(time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})

			By("Delete prerequisite Group", func() {
				Expect(kubeClient.Delete(ctx, testGroup)).To(Succeed())
				Eventually(func(g Gomega) {
					group := &generatedv1.Group{}
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testGroup), group)
					g.Expect(err).NotTo(Succeed())
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})
		})

		It("Should NOT delete from Atlas when ResourcePolicyKeep annotation is set", Label("focus-ipaccesslist-kept"), func() {
			groupName := fmt.Sprintf("test-group-%s", rand.String(6))
			groupParams := testparams.New(orgID, control.MustEnvVar("OPERATOR_NAMESPACE"), DefaultGlobalCredentials).
				WithGroupName(groupName).
				WithNamespace(testNamespace.Name)

			var testGroup *generatedv1.Group
			By("Create prerequisite Group", func() {
				objs := samples.MustLoadSampleObjects("atlas_generated_v1_group.yaml")
				Expect(len(objs)).To(Equal(1))
				testGroup = objs[0].(*generatedv1.Group)
				applyTestParamsToGroup(testGroup, groupParams)
				Expect(kubeClient.Create(ctx, testGroup)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, testGroup)).To(Succeed())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
				Expect(testGroup.Status.V20250312).NotTo(BeNil())
				Expect(testGroup.Status.V20250312.Id).NotTo(BeNil())
			})

			entryName := fmt.Sprintf("ipal-%s", rand.String(6))
			var testIAL *generatedv1.IPAccessListEntry
			By("Create IPAccessList with keep annotation", func() {
				objs := samples.MustLoadSampleObjects("atlas_generated_v1_ipaccesslistentry_with_groupref.yaml")
				Expect(len(objs)).To(Equal(1))
				testIAL = objs[0].(*generatedv1.IPAccessListEntry)
				applyTestParamsToIPAccessList(testIAL, testNamespace.Name, entryName, testGroup.GetName())
				testIAL.SetAnnotations(map[string]string{
					customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep,
				})
				Expect(kubeClient.Create(ctx, testIAL)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, testIAL)).To(Succeed())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
				Expect(testIAL.Status.V20250312).NotTo(BeNil())
			})

			By("Delete IPAccessList from cluster - should NOT delete from Atlas", func() {
				groupID := *testIAL.Status.V20250312.GroupId
				Expect(kubeClient.Delete(ctx, testIAL)).To(Succeed())

				Eventually(func(g Gomega) {
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testIAL), testIAL)
					g.Expect(err).ToNot(Succeed())
				}).WithContext(ctx).WithTimeout(time.Minute).WithPolling(time.Second).Should(Succeed())

				Eventually(func(g Gomega) {
					atlasEntry, _, err := atlasClient.ProjectIPAccessListApi.GetAccessListEntry(ctx, groupID, testCIDR).Execute()
					g.Expect(err).ToNot(HaveOccurred())
					g.Expect(atlasEntry).ToNot(BeNil())
				}).WithContext(ctx).WithTimeout(time.Minute).WithPolling(time.Second).Should(Succeed())

				groupID2 := *testGroup.Status.V20250312.Id
				_, err := atlasClient.ProjectIPAccessListApi.DeleteAccessListEntry(ctx, groupID2, testCIDR).Execute()
				Expect(err).ToNot(HaveOccurred())
			})

			By("Delete prerequisite Group", func() {
				Expect(kubeClient.Delete(ctx, testGroup)).To(Succeed())
				Eventually(func(g Gomega) {
					group := &generatedv1.Group{}
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testGroup), group)
					g.Expect(err).NotTo(Succeed())
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})
		})

		It("Should mark entry as Expired when deleteAfterDate passes", Label("focus-ipaccesslist-expiry"), func() {
			groupName := fmt.Sprintf("test-group-%s", rand.String(6))
			groupParams := testparams.New(orgID, control.MustEnvVar("OPERATOR_NAMESPACE"), DefaultGlobalCredentials).
				WithGroupName(groupName).
				WithNamespace(testNamespace.Name)

			var testGroup *generatedv1.Group
			By("Create prerequisite Group", func() {
				objs := samples.MustLoadSampleObjects("atlas_generated_v1_group.yaml")
				Expect(len(objs)).To(Equal(1))
				testGroup = objs[0].(*generatedv1.Group)
				applyTestParamsToGroup(testGroup, groupParams)
				Expect(kubeClient.Create(ctx, testGroup)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, testGroup)).To(Succeed())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
				Expect(testGroup.Status.V20250312).NotTo(BeNil())
				Expect(testGroup.Status.V20250312.Id).NotTo(BeNil())
			})

			entryName := fmt.Sprintf("ipal-%s", rand.String(6))
			var testIAL *generatedv1.IPAccessListEntry

			// 1 min ahead of current time
			deleteAfterDate := time.Now().UTC().Add(time.Minute).Format(time.RFC3339)

			By("Create IPAccessList with deleteAfterDate 1 minute from now (UTC)", func() {
				objs := samples.MustLoadSampleObjects("atlas_generated_v1_ipaccesslistentry_with_groupref.yaml")
				Expect(len(objs)).To(Equal(1))
				testIAL = objs[0].(*generatedv1.IPAccessListEntry)
				applyTestParamsToIPAccessList(testIAL, testNamespace.Name, entryName, testGroup.GetName())
				testIAL.Spec.V20250312.Entry.DeleteAfterDate = pointer.MakePtr(deleteAfterDate)
				Expect(kubeClient.Create(ctx, testIAL)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, testIAL)).To(Succeed())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})

			By("Wait for Atlas to expire the entry and operator to mark it Expired", func() {
				Eventually(func(g Gomega) {
					g.Expect(kubeClient.Get(ctx, client.ObjectKeyFromObject(testIAL), testIAL)).To(Succeed())
					stateCondition := meta.FindStatusCondition(testIAL.GetConditions(), "State")
					g.Expect(stateCondition).NotTo(BeNil())
					g.Expect(stateCondition.Message).To(Equal("Expired"))
				}).WithContext(ctx).WithTimeout(10 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
			})

			By("Verify entry still exists in Kubernetes (not deleted)", func() {
				Expect(kubeClient.Get(ctx, client.ObjectKeyFromObject(testIAL), testIAL)).To(Succeed())
				Expect(testIAL.GetDeletionTimestamp()).To(BeNil())
			})

			By("Verify entry no longer exists in Atlas", func() {
				groupID := *testIAL.Status.V20250312.GroupId
				_, r, err := atlasClient.ProjectIPAccessListApi.GetAccessListEntry(ctx, groupID, testCIDR).Execute()
				Expect(err).To(HaveOccurred())
				Expect(httputil.StatusCode(r)).To(Equal(http.StatusNotFound))
			})

			By("Delete expired entry from Kubernetes", func() {
				Expect(kubeClient.Delete(ctx, testIAL)).To(Succeed())
				// The operator treats ATLAS_NETWORK_PERMISSION_ENTRY_NOT_FOUND as success so
				// the finalizer is removed and the resource disappears without manual intervention.
				Eventually(func(g Gomega) {
					ial := &generatedv1.IPAccessListEntry{}
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testIAL), ial)
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})

			By("Delete prerequisite Group", func() {
				Expect(kubeClient.Delete(ctx, testGroup)).To(Succeed())
				Eventually(func(g Gomega) {
					group := &generatedv1.Group{}
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testGroup), group)
					g.Expect(err).NotTo(Succeed())
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})
		})

		It("Should fail if groupRef points to non-existent Group", Label("focus-ipaccesslist-fail"), func() {
			entryName := fmt.Sprintf("ipal-%s", rand.String(6))

			var testIAL *generatedv1.IPAccessListEntry
			By("Create IPAccessList with non-existent groupRef", func() {
				objs := samples.MustLoadSampleObjects("atlas_generated_v1_ipaccesslistentry_with_groupref.yaml")
				Expect(len(objs)).To(Equal(1))
				testIAL = objs[0].(*generatedv1.IPAccessListEntry)
				applyTestParamsToIPAccessList(testIAL, testNamespace.Name, entryName, "non-existent-group")
				Expect(kubeClient.Create(ctx, testIAL)).To(Succeed())
			})

			By("Wait for IPAccessList to report an error condition", func() {
				ial := &generatedv1.IPAccessListEntry{
					ObjectMeta: metav1.ObjectMeta{Name: entryName, Namespace: testNamespace.Name},
				}
				Eventually(func(g Gomega) {
					g.Expect(kubeClient.Get(ctx, client.ObjectKeyFromObject(ial), ial)).To(Succeed())
					g.Expect(ial.GetConditions()).NotTo(BeEmpty())
				}).WithContext(ctx).WithTimeout(30 * time.Second).WithPolling(time.Second).Should(Succeed())
				testIAL = ial
				Expect(meta.IsStatusConditionTrue(testIAL.GetConditions(), "Ready")).To(BeFalse())
				readyCondition := meta.FindStatusCondition(testIAL.GetConditions(), "Ready")
				Expect(readyCondition.Reason).To(Equal("Error"))
			})

			By("Force delete IPAccessList by removing finalizers", func() {
				Expect(kubeClient.Get(ctx, client.ObjectKeyFromObject(testIAL), testIAL)).To(Succeed())
				testIAL.SetFinalizers([]string{})
				Expect(kubeClient.Update(ctx, testIAL)).To(Succeed())
				Expect(kubeClient.Delete(ctx, testIAL)).To(Succeed())
				Eventually(func(g Gomega) {
					ial := &generatedv1.IPAccessListEntry{}
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testIAL), ial)
					g.Expect(err).NotTo(Succeed())
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})
		})
	})
})

var _ = Describe("IPAccessList with Deletion Protection", Ordered, Label("ipaccesslist"), func() {
	var ctx context.Context
	var kubeClient client.Client
	var ako operator.Operator
	var testNamespace *corev1.Namespace
	var orgID string
	var atlasClient *admin.APIClient

	_ = BeforeAll(func() {
		ctx = context.Background()

		deletionProtectionOn := true
		ako = runTestAKO(DefaultGlobalCredentials, control.MustEnvVar("OPERATOR_NAMESPACE"), deletionProtectionOn)
		ako.Start(ctx, GinkgoT())

		DeferCleanup(func() {
			if ako != nil {
				ako.Stop(GinkgoT())
			}
		})

		testClient, err := kube.NewTestClient()
		Expect(err).To(Succeed())
		kubeClient = testClient
		Expect(kube.AssertCRDNames(ctx, kubeClient,
			"groups.atlas.generated.mongodb.com",
			"ipaccesslistentries.atlas.generated.mongodb.com",
		)).To(Succeed())

		atlasClient, orgID = newTestAtlasClient()
	})

	_ = BeforeEach(func() {
		testNamespace = &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("ipal-protect-ns-%s", rand.String(6)),
		}}
		Expect(kubeClient.Create(ctx, testNamespace)).To(Succeed())
		Expect(ako.Running()).To(BeTrue(), "Operator must be running")
		Expect(resources.CopyCredentialsToNamespace(
			ctx,
			kubeClient,
			DefaultGlobalCredentials,
			control.MustEnvVar("OPERATOR_NAMESPACE"),
			testNamespace.Name, GinkGoFieldOwner),
		).To(Succeed())
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
		}).WithContext(ctx).WithTimeout(time.Minute).WithPolling(time.Second).To(BeFalse())
	})

	Describe("Deleting the IPAccessList", Label("focus-ipaccesslist-deletion-protected"), func() {
		It("Should NOT delete from Atlas when deletion protection is enabled", func() {
			groupName := fmt.Sprintf("test-group-%s", rand.String(6))
			groupParams := testparams.New(orgID, control.MustEnvVar("OPERATOR_NAMESPACE"), DefaultGlobalCredentials).
				WithGroupName(groupName).
				WithNamespace(testNamespace.Name)

			var testGroup *generatedv1.Group
			By("Create prerequisite Group", func() {
				objs := samples.MustLoadSampleObjects("atlas_generated_v1_group.yaml")
				Expect(len(objs)).To(Equal(1))
				testGroup = objs[0].(*generatedv1.Group)
				applyTestParamsToGroup(testGroup, groupParams)
				Expect(kubeClient.Create(ctx, testGroup)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, testGroup)).To(Succeed())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
				Expect(testGroup.Status.V20250312).NotTo(BeNil())
				Expect(testGroup.Status.V20250312.Id).NotTo(BeNil())
			})

			entryName := fmt.Sprintf("ipal-%s", rand.String(6))
			var testIAL *generatedv1.IPAccessListEntry
			By("Create IPAccessList", func() {
				objs := samples.MustLoadSampleObjects("atlas_generated_v1_ipaccesslistentry_with_groupref.yaml")
				Expect(len(objs)).To(Equal(1))
				testIAL = objs[0].(*generatedv1.IPAccessListEntry)
				applyTestParamsToIPAccessList(testIAL, testNamespace.Name, entryName, testGroup.GetName())
				Expect(kubeClient.Create(ctx, testIAL)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, testIAL)).To(Succeed())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

				Expect(testIAL.Status.V20250312).NotTo(BeNil())
				Expect(meta.IsStatusConditionTrue(testIAL.GetConditions(), "Ready")).To(BeTrue())
			})

			By("Delete IPAccessList from cluster - should NOT delete from Atlas", func() {
				groupID := *testIAL.Status.V20250312.GroupId
				Expect(kubeClient.Delete(ctx, testIAL)).To(Succeed())

				Eventually(func(g Gomega) {
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testIAL), testIAL)
					g.Expect(err).To(HaveOccurred())
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(time.Minute).WithPolling(time.Second).Should(Succeed())

				Eventually(func(g Gomega) {
					atlasEntry, _, err := atlasClient.ProjectIPAccessListApi.GetAccessListEntry(ctx, groupID, testCIDR).Execute()
					g.Expect(err).ToNot(HaveOccurred())
					g.Expect(atlasEntry).NotTo(BeNil())
					g.Expect(atlasEntry.GetCidrBlock()).To(Equal(testCIDR))
				}).WithContext(ctx).WithTimeout(30 * time.Second).WithPolling(5 * time.Second).Should(Succeed())
			})

			By("Clean up Atlas entry manually", func() {
				groupID := *testGroup.Status.V20250312.Id
				_, err := atlasClient.ProjectIPAccessListApi.DeleteAccessListEntry(ctx, groupID, testCIDR).Execute()
				Expect(err).ToNot(HaveOccurred())

				Eventually(func(g Gomega) {
					_, r, err := atlasClient.ProjectIPAccessListApi.GetAccessListEntry(ctx, groupID, testCIDR).Execute()
					g.Expect(err).ToNot(BeNil())
					g.Expect(httputil.StatusCode(r)).To(Equal(http.StatusNotFound))
				}).WithContext(ctx).WithTimeout(time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})

			By("Delete prerequisite Group", func() {
				Expect(kubeClient.Delete(ctx, testGroup)).To(Succeed())
				Eventually(func(g Gomega) {
					group := &generatedv1.Group{}
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testGroup), group)
					g.Expect(err).NotTo(Succeed())
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})
		})
	})
})

func applyTestParamsToIPAccessList(ial *generatedv1.IPAccessListEntry, namespace, name, groupRefName string) {
	ial.SetNamespace(namespace)
	ial.SetName(name)

	if ial.Spec.V20250312 == nil {
		ial.Spec.V20250312 = &generatedv1.IPAccessListEntrySpecV20250312{}
	}
	ial.Spec.V20250312.GroupRef = &k8s.LocalReference{Name: groupRefName}
	ial.Spec.V20250312.GroupId = nil
	ial.Spec.V20250312.Entry = &generatedv1.IPAccessListEntrySpecV20250312Entry{
		CidrBlock: pointer.MakePtr(testCIDR),
		Comment:   pointer.MakePtr("e2e test entry"),
	}
}
