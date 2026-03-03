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
	"go.mongodb.org/atlas-sdk/v20250312014/admin"
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

var _ = Describe("Group CRUD", Ordered, Label("group-ctlr"), func() {
	var ctx context.Context
	var kubeClient client.Client
	var ako operator.Operator
	var testNamespace *corev1.Namespace
	var orgID string
	var atlasClient *admin.APIClient

	_ = BeforeAll(func() {
		deletionProtectionOff := false
		ako = runTestAKO(DefaultGlobalCredentials, control.MustEnvVar("OPERATOR_NAMESPACE"), deletionProtectionOff)
		ako.Start(GinkgoT())

		DeferCleanup(func() {
			if ako != nil {
				ako.Stop(GinkgoT())
			}
		})

		ctx = context.Background()
		testClient, err := kube.NewTestClient()
		Expect(err).To(Succeed())
		kubeClient = testClient
		Expect(kube.AssertCRDNames(ctx, kubeClient, "groups.atlas.generated.mongodb.com")).To(Succeed())

		atlasClient, orgID = newTestAtlasClient()
	})

	_ = BeforeEach(func() {
		testNamespace = &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("group-ctlr-ns-%s", rand.String(6)),
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

	Describe("Group CRUD lifecycle", func() {
		It("Should create, update, and delete Group", Label("focus-group-crud"), func() {
			groupName := fmt.Sprintf("test-group-crud-%s", rand.String(6))
			testParams := testparams.New(orgID, control.MustEnvVar("OPERATOR_NAMESPACE"), DefaultGlobalCredentials).
				WithGroupName(groupName).
				WithNamespace(testNamespace.Name)

			var testGroup *generatedv1.Group
			By("Create Group", func() {
				objs := samples.MustLoadSampleObjects("atlas_generated_v1_group.yaml")
				Expect(len(objs)).To(Equal(1))
				testGroup = objs[0].(*generatedv1.Group)
				applyTestParamsToGroup(testGroup, testParams)
				testGroup.SetAnnotations(map[string]string{
					customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyDelete,
				})
				Expect(kubeClient.Create(ctx, testGroup)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, testGroup)).To(Succeed())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

				// Verify status
				Expect(testGroup.Status.V20250312).NotTo(BeNil())
				Expect(testGroup.Status.V20250312.Id).NotTo(BeNil())
				Expect(meta.IsStatusConditionTrue(testGroup.GetConditions(), "Ready")).To(BeTrue())
				Expect(meta.IsStatusConditionTrue(testGroup.GetConditions(), "State")).To(BeTrue())

				// Verify in Atlas
				atlasGroup, _, err := atlasClient.ProjectsApi.GetGroup(ctx, *testGroup.Status.V20250312.Id).Execute()
				Expect(err).ToNot(HaveOccurred())
				Expect(atlasGroup.Name).To(Equal(groupName))
			})

			By("Update Group tags", func() {
				Expect(kubeClient.Get(ctx, client.ObjectKeyFromObject(testGroup), testGroup)).To(Succeed())
				testGroup.Spec.V20250312.Entry.Tags = &[]generatedv1.Tags{
					{Key: "environment", Value: "test"},
					{Key: "managed-by", Value: "ako"},
				}
				Expect(kubeClient.Update(ctx, testGroup)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceUpdated(ctx, kubeClient, testGroup)).To(Succeed())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

				// Verify in Atlas
				Expect(testGroup.Status.V20250312).NotTo(BeNil())
				Expect(testGroup.Status.V20250312.Id).NotTo(BeNil())
				atlasGroup, _, err := atlasClient.ProjectsApi.GetGroup(ctx, *testGroup.Status.V20250312.Id).Execute()
				Expect(err).ToNot(HaveOccurred())
				atlasTags := atlasGroup.GetTags()
				Expect(len(atlasTags)).To(Equal(2))
				tagMap := make(map[string]string)
				for _, tag := range atlasTags {
					tagMap[tag.GetKey()] = tag.GetValue()
				}
				Expect(tagMap["environment"]).To(Equal("test"))
				Expect(tagMap["managed-by"]).To(Equal("ako"))
			})

			By("Delete Group from cluster - should delete from Atlas", func() {
				groupID := testGroup.Status.V20250312.Id
				Expect(groupID).NotTo(BeNil())
				Expect(kubeClient.Delete(ctx, testGroup)).To(Succeed())

				Eventually(func(g Gomega) {
					group := &generatedv1.Group{}
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testGroup), group)
					g.Expect(err).NotTo(Succeed())
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

				Eventually(func(g Gomega) {
					_, r, err := atlasClient.ProjectsApi.GetGroup(ctx, *groupID).Execute()
					g.Expect(err).ToNot(BeNil())
					g.Expect(httputil.StatusCode(r)).To(Equal(http.StatusNotFound))
				}).WithContext(ctx).WithTimeout(time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})
		})

		It("Should Fail if Secret is wrong", Label("focus-group-fail-secret"), func() {
			groupName := fmt.Sprintf("test-group-fail-%s", rand.String(6))
			testParams := testparams.New(orgID, control.MustEnvVar("OPERATOR_NAMESPACE"), "non-existent-secret").
				WithGroupName(groupName).
				WithNamespace(testNamespace.Name)
			var testGroup *generatedv1.Group

			By("Load sample Group YAML with wrong secret", func() {
				objs := samples.MustLoadSampleObjects("atlas_generated_v1_group.yaml")
				Expect(len(objs)).To(Equal(1))
				testGroup := objs[0].(*generatedv1.Group)
				applyTestParamsToGroup(testGroup, testParams)
				Expect(kubeClient.Create(ctx, testGroup)).To(Succeed())
			})

			By("Wait for Group to fail", func() {
				testGroup = &generatedv1.Group{
					ObjectMeta: metav1.ObjectMeta{Name: groupName, Namespace: testNamespace.Name},
				}
				Eventually(func(g Gomega) {
					g.Expect(kubeClient.Get(ctx, client.ObjectKeyFromObject(testGroup), testGroup)).To(Succeed())
					g.Expect(testGroup.GetConditions()).NotTo(BeEmpty())
				}).WithContext(ctx).WithTimeout(time.Second * 5).WithPolling(time.Second).Should(Succeed())
				Expect(meta.IsStatusConditionTrue(testGroup.GetConditions(), "Ready")).To(BeFalse())
				Expect(meta.IsStatusConditionTrue(testGroup.GetConditions(), "State")).To(BeFalse())
				readyCondition := meta.FindStatusCondition(testGroup.GetConditions(), "Ready")
				Expect(readyCondition.Reason).To(Equal("Error"))
			})

			By("Enforce deletion of Group by removing finalizers before namespace deletion", func() {
				Expect(kubeClient.Get(ctx, client.ObjectKeyFromObject(testGroup), testGroup)).To(Succeed())
				testGroup.SetFinalizers([]string{})
				Expect(kubeClient.Update(ctx, testGroup)).To(Succeed())
				Expect(kubeClient.Delete(ctx, testGroup)).To(Succeed())
				Eventually(func(g Gomega) {
					group := &generatedv1.Group{}
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testGroup), group)
					g.Expect(err).NotTo(Succeed())
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})
		})

		It("Should NOT delete from Atlas when ResourcePolicyKeep annotation is set", Label("focus-group-kept"),
			func() {
				groupName := fmt.Sprintf("test-group-keep-%s", rand.String(6))
				testParams := testparams.New(orgID, control.MustEnvVar("OPERATOR_NAMESPACE"), DefaultGlobalCredentials).
					WithGroupName(groupName).
					WithNamespace(testNamespace.Name)

				var testGroup *generatedv1.Group
				By("Create Group with keep annotation", func() {
					objs := samples.MustLoadSampleObjects("atlas_generated_v1_group.yaml")
					Expect(len(objs)).To(Equal(1))
					testGroup = objs[0].(*generatedv1.Group)
					applyTestParamsToGroup(testGroup, testParams)
					testGroup.SetAnnotations(map[string]string{
						customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep,
					})
					Expect(kubeClient.Create(ctx, testGroup)).To(Succeed())

					Eventually(func(g Gomega) {
						g.Expect(resources.CheckResourceReady(ctx, kubeClient, testGroup)).To(Succeed())
					}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
				})

				By("Delete Group from cluster - should NOT delete from Atlas", func() {
					groupID := testGroup.Status.V20250312.Id
					Expect(groupID).NotTo(BeNil())
					Expect(kubeClient.Delete(ctx, testGroup)).To(Succeed())

					// Wait for K8s resource to be deleted
					Eventually(func(g Gomega) {
						err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testGroup), testGroup)
						g.Expect(err).ToNot(Succeed())
					}).WithContext(ctx).WithTimeout(time.Minute).WithPolling(time.Second).Should(Succeed())

					// Verify it still exists in Atlas
					time.Sleep(10 * time.Second)
					atlasGroup, _, err := atlasClient.ProjectsApi.GetGroup(ctx, *groupID).Execute()
					Expect(err).ToNot(HaveOccurred())
					Expect(atlasGroup).ToNot(BeNil())

					// Clean up manually
					_, err = atlasClient.ProjectsApi.DeleteGroup(ctx, *groupID).Execute()
					Expect(err).ToNot(HaveOccurred())
				})
			})
	})

	Describe("Using the global Connection Secret", func() {
		It("Should Succeed", Label("focus-group-global-creds"), func() {
			groupName := fmt.Sprintf("test-group-global-%s", rand.String(6))
			testParams := testparams.New(orgID, control.MustEnvVar("OPERATOR_NAMESPACE"), DefaultGlobalCredentials).
				WithGroupName(groupName).
				WithNamespace(testNamespace.Name)

			By("Create Group without explicit connectionSecretRef (uses global)", func() {
				objs := samples.MustLoadSampleObjects("atlas_generated_v1_group.yaml")
				Expect(len(objs)).To(Equal(1))
				testGroup := objs[0].(*generatedv1.Group)
				applyTestParamsToGroup(testGroup, testParams)
				// Remove connectionSecretRef to use global secret
				testGroup.Spec.ConnectionSecretRef = nil
				Expect(kubeClient.Create(ctx, testGroup)).To(Succeed())
			})

			By("Wait for Group to be Ready", func() {
				testGroup := &generatedv1.Group{
					ObjectMeta: metav1.ObjectMeta{Name: groupName, Namespace: testNamespace.Name},
				}
				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, testGroup)).To(Succeed())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

				Expect(meta.IsStatusConditionTrue(testGroup.GetConditions(), "Ready")).To(BeTrue())
				Expect(meta.IsStatusConditionTrue(testGroup.GetConditions(), "State")).To(BeTrue())
			})
		})
	})

	Describe("Importing existing Group in Atlas project", func() {
		It("Should reconcile existing project", Label("focus-group-existing"), func() {
			groupName := fmt.Sprintf("existing-group-%s", rand.String(6))
			testParams := testparams.New(orgID, control.MustEnvVar("OPERATOR_NAMESPACE"), DefaultGlobalCredentials).
				WithGroupName(groupName).
				WithNamespace(testNamespace.Name)
			var atlasGroupID string
			var testGroup *generatedv1.Group
			By("Create project in Atlas first", func() {
				atlasGroup := admin.Group{
					OrgId:                     orgID,
					Name:                      groupName,
					WithDefaultAlertsSettings: pointer.MakePtr(true),
				}
				createdGroup, _, err := atlasClient.ProjectsApi.CreateGroup(ctx, &atlasGroup).Execute()
				Expect(err).ToNot(HaveOccurred())
				atlasGroupID = createdGroup.GetId()
				Expect(atlasGroupID).NotTo(BeEmpty())
			})

			By("Create Group in cluster with import annotation - should reconcile existing Atlas project", func() {
				objs := samples.MustLoadSampleObjects("atlas_generated_v1_group.yaml")
				Expect(len(objs)).To(Equal(1))
				testGroup = objs[0].(*generatedv1.Group)
				applyTestParamsToGroup(testGroup, testParams)
				// Add import annotation: mongodb.com/external-id with the Atlas project ID
				// This tells the controller to import the existing resource instead of creating a new one
				testGroup.SetAnnotations(map[string]string{
					"mongodb.com/external-id": atlasGroupID,
				})
				Expect(kubeClient.Create(ctx, testGroup)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, testGroup)).To(Succeed())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

				// Verify it has the correct ID
				Expect(testGroup.Status.V20250312).NotTo(BeNil())
				Expect(testGroup.Status.V20250312.Id).NotTo(BeNil())

				// Verify in Atlas
				atlasGroup, _, err := atlasClient.ProjectsApi.GetGroup(ctx, *testGroup.Status.V20250312.Id).Execute()
				Expect(err).ToNot(HaveOccurred())
				Expect(atlasGroup.Name).To(Equal(groupName))

				Expect(meta.IsStatusConditionTrue(testGroup.GetConditions(), "Ready")).To(BeTrue())
				Expect(meta.IsStatusConditionTrue(testGroup.GetConditions(), "State")).To(BeTrue())
			})

			By("Delete Group from cluster before namespace gets deleted", func() {
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

var _ = Describe("Group with Deletion Protection", Ordered, Label("group-ctlr"), func() {
	var ctx context.Context
	var kubeClient client.Client
	var ako operator.Operator
	var testNamespace *corev1.Namespace
	var orgID string
	var atlasClient *admin.APIClient

	_ = BeforeAll(func() {
		deletionProtectionOn := true
		ako = runTestAKO(DefaultGlobalCredentials, control.MustEnvVar("OPERATOR_NAMESPACE"), deletionProtectionOn)
		ako.Start(GinkgoT())

		DeferCleanup(func() {
			if ako != nil {
				ako.Stop(GinkgoT())
			}
		})

		ctx = context.Background()
		testClient, err := kube.NewTestClient()
		Expect(err).To(Succeed())
		kubeClient = testClient
		Expect(kube.AssertCRDNames(ctx, kubeClient, "groups.atlas.generated.mongodb.com")).To(Succeed())

		atlasClient, orgID = newTestAtlasClient()
	})

	_ = BeforeEach(func() {
		testNamespace = &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("group-ctlr-protect-ns-%s", rand.String(6)),
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

	Describe("Deleting the Group", Label("focus-group-deletion-protected"), func() {
		It("Should NOT delete from Atlas when deletion protection is enabled", func() {
			groupName := fmt.Sprintf("test-group-protect-%s", rand.String(6))
			testParams := testparams.New(orgID, control.MustEnvVar("OPERATOR_NAMESPACE"), DefaultGlobalCredentials).
				WithGroupName(groupName).
				WithNamespace(testNamespace.Name)

			var testGroup *generatedv1.Group
			By("Create Group in cluster", func() {
				objs := samples.MustLoadSampleObjects("atlas_generated_v1_group.yaml")
				Expect(len(objs)).To(Equal(1))
				testGroup = objs[0].(*generatedv1.Group)
				applyTestParamsToGroup(testGroup, testParams)
				Expect(kubeClient.Create(ctx, testGroup)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, testGroup)).To(Succeed())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

				// Verify it was created in Atlas
				Expect(testGroup.Status.V20250312).NotTo(BeNil())
				Expect(testGroup.Status.V20250312.Id).NotTo(BeNil())
				Expect(meta.IsStatusConditionTrue(testGroup.GetConditions(), "Ready")).To(BeTrue())
				Expect(meta.IsStatusConditionTrue(testGroup.GetConditions(), "State")).To(BeTrue())
			})

			By("Delete Group from cluster - should NOT delete from Atlas", func() {
				groupID := testGroup.Status.V20250312.Id
				Expect(groupID).NotTo(BeNil())
				Expect(kubeClient.Delete(ctx, testGroup)).To(Succeed())

				// Verify Kubernetes resource is deleted
				Eventually(func(g Gomega) {
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testGroup), testGroup)
					g.Expect(err).To(HaveOccurred())
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(time.Minute).WithPolling(time.Second).Should(Succeed())

				// Verify Group still exists in Atlas (deletion protection prevented deletion)
				Eventually(func(g Gomega) {
					atlasGroup, _, err := atlasClient.ProjectsApi.GetGroup(ctx, *groupID).Execute()
					g.Expect(err).ToNot(HaveOccurred())
					g.Expect(atlasGroup).NotTo(BeNil())
					g.Expect(atlasGroup.GetId()).To(Equal(*groupID))
				}).WithContext(ctx).WithTimeout(30 * time.Second).WithPolling(5 * time.Second).Should(Succeed())
			})

			By("Clean up Atlas resource manually", func() {
				groupID := testGroup.Status.V20250312.Id
				Expect(groupID).NotTo(BeNil())
				_, err := atlasClient.ProjectsApi.DeleteGroup(ctx, *groupID).Execute()
				Expect(err).ToNot(HaveOccurred())

				// Verify it's deleted from Atlas
				Eventually(func(g Gomega) {
					_, r, err := atlasClient.ProjectsApi.GetGroup(ctx, *groupID).Execute()
					g.Expect(err).ToNot(BeNil())
					g.Expect(httputil.StatusCode(r)).To(Equal(http.StatusNotFound))
				}).WithContext(ctx).WithTimeout(time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})
		})
	})
})

// applyTestParamsToGroup applies test parameters to a Group object.
// This is similar to testparams.ApplyToGroup but works with generatedv1.Group instead of nextapiv1.Group.
func applyTestParamsToGroup(group *generatedv1.Group, params *testparams.TestParams) {
	group.SetNamespace(params.Namespace)
	group.SetName(params.GroupName)

	if group.Spec.ConnectionSecretRef == nil {
		group.Spec.ConnectionSecretRef = &k8s.LocalReference{}
	}
	group.Spec.ConnectionSecretRef.Name = params.CredentialsSecretName

	if group.Spec.V20250312 == nil {
		group.Spec.V20250312 = &generatedv1.V20250312{}
	}
	if group.Spec.V20250312.Entry == nil {
		group.Spec.V20250312.Entry = &generatedv1.Entry{}
	}
	group.Spec.V20250312.Entry.OrgId = params.OrgID
	group.Spec.V20250312.Entry.Name = params.GroupName
}
