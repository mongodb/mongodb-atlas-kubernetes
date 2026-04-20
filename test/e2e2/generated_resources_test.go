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
	"time"

	k8s "github.com/crd2go/crd2go/k8s"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/controller-runtime/pkg/client"

	generatedv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/generated/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/operator"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/resources"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/samples"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/testparams"
)

var _ = Describe("Generated Resources Integration", Ordered, Label("generated-resources"), func() {
	var ctx context.Context
	var kubeClient client.Client
	var ako operator.Operator
	var testNamespace *corev1.Namespace
	var orgID string

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
			"clusters.atlas.generated.mongodb.com",
			"flexclusters.atlas.generated.mongodb.com",
			"databaseusers.atlas.generated.mongodb.com",
			"ipaccesslistentries.atlas.generated.mongodb.com",
		)).To(Succeed())

		_, orgID = newTestAtlasClient()
	})

	_ = BeforeEach(func() {
		testNamespace = &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("gen-res-ns-%s", rand.String(6)),
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
		Expect(kubeClient.Delete(ctx, testNamespace)).To(Succeed())
		Eventually(func(g Gomega) bool {
			return kubeClient.Get(ctx, client.ObjectKeyFromObject(testNamespace), testNamespace) == nil
		}).WithContext(ctx).WithTimeout(time.Minute).WithPolling(time.Second).To(BeFalse())
	})

	Describe("All generated resource types in a single project", func() {
		It("Should create connection secrets for DatabaseUser for both Cluster and FlexCluster",
			Label("focus-generated-resources-all"),
			func() {
				groupName := fmt.Sprintf("test-group-%s", rand.String(6))
				groupParams := testparams.New(orgID, control.MustEnvVar("OPERATOR_NAMESPACE"), DefaultGlobalCredentials).
					WithGroupName(groupName).
					WithNamespace(testNamespace.Name)

				clusterName := fmt.Sprintf("cluster-%s", rand.String(6))
				flexClusterName := fmt.Sprintf("flex-%s", rand.String(6))
				username := fmt.Sprintf("testuser-%s", rand.String(6))
				passwordSecretName := fmt.Sprintf("dbuser-pass-%s", rand.String(6))
				ipalEntryName := fmt.Sprintf("ipal-%s", rand.String(6))

				const ipalCIDR = "203.0.113.128/25"

				var (
					testGroup       *generatedv1.Group
					testCluster     *generatedv1.Cluster
					testFlexCluster *generatedv1.FlexCluster
					testDBUser      *generatedv1.DatabaseUser
					testIAL         *generatedv1.IPAccessListEntry
				)

				By("Create Group", func() {
					objs := samples.MustLoadSampleObjects("atlas_generated_v1_group.yaml")
					Expect(len(objs)).To(Equal(1))
					testGroup = objs[0].(*generatedv1.Group)
					applyTestParamsToGroup(testGroup, groupParams)
					Expect(kubeClient.Create(ctx, testGroup)).To(Succeed())
				})

				By("Create Cluster (M0) referencing Group", func() {
					testCluster = newSharedCluster(clusterName, testNamespace.Name, testGroup.GetName())
					Expect(kubeClient.Create(ctx, testCluster)).To(Succeed())
				})

				By("Create FlexCluster referencing Group", func() {
					testFlexCluster = newSharedFlexCluster(flexClusterName, testNamespace.Name, testGroup.GetName())
					Expect(kubeClient.Create(ctx, testFlexCluster)).To(Succeed())
				})

				By("Create IPAccessListEntry referencing Group", func() {
					testIAL = &generatedv1.IPAccessListEntry{
						ObjectMeta: metav1.ObjectMeta{
							Name:      ipalEntryName,
							Namespace: testNamespace.Name,
						},
						Spec: generatedv1.IPAccessListEntrySpec{
							V20250312: &generatedv1.IPAccessListEntrySpecV20250312{
								GroupRef: &k8s.LocalReference{Name: testGroup.GetName()},
								Entry: &generatedv1.IPAccessListEntrySpecV20250312Entry{
									CidrBlock: pointer.MakePtr(ipalCIDR),
									Comment:   pointer.MakePtr("generated resources integration test"),
								},
							},
						},
					}
					Expect(kubeClient.Create(ctx, testIAL)).To(Succeed())
				})

				By("Create DatabaseUser referencing Group", func() {
					Expect(kubeClient.Create(ctx, newPasswordSecret(testNamespace.Name, passwordSecretName))).To(Succeed())

					objs := samples.MustLoadSampleObjects("atlas_generated_v1_databaseuser_with_groupref.yaml")
					Expect(len(objs)).To(Equal(1))
					testDBUser = objs[0].(*generatedv1.DatabaseUser)
					applyTestParamsToDBUser(testDBUser, testNamespace.Name, username, testGroup.GetName(), passwordSecretName)
					Expect(kubeClient.Create(ctx, testDBUser)).To(Succeed())
				})

				By("Wait for Group to be Ready", func() {
					Eventually(func(g Gomega) {
						g.Expect(resources.CheckResourceReady(ctx, kubeClient, testGroup)).To(Succeed())
					}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
					Expect(testGroup.Status.V20250312).NotTo(BeNil())
					Expect(testGroup.Status.V20250312.Id).NotTo(BeNil())
				})

				By("Wait for Cluster to be Ready", func() {
					Eventually(func(g Gomega) {
						g.Expect(resources.CheckResourceReady(ctx, kubeClient, testCluster)).To(Succeed())
					}).WithContext(ctx).WithTimeout(clusterCreateTimeout).WithPolling(clusterPollingInterval).Should(Succeed())
					Expect(testCluster.Status.V20250312).NotTo(BeNil())
					Expect(testCluster.Status.V20250312.ConnectionStrings).NotTo(BeNil())
				})

				By("Wait for FlexCluster to be Ready", func() {
					Eventually(func(g Gomega) {
						g.Expect(resources.CheckResourceReady(ctx, kubeClient, testFlexCluster)).To(Succeed())
					}).WithContext(ctx).WithTimeout(clusterCreateTimeout).WithPolling(clusterPollingInterval).Should(Succeed())
					Expect(testFlexCluster.Status.V20250312).NotTo(BeNil())
					Expect(testFlexCluster.Status.V20250312.ConnectionStrings).NotTo(BeNil())
				})

				By("Wait for IPAccessListEntry to be Ready", func() {
					Eventually(func(g Gomega) {
						g.Expect(resources.CheckResourceReady(ctx, kubeClient, testIAL)).To(Succeed())
					}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
					Expect(testIAL.Status.V20250312).NotTo(BeNil())
					Expect(testIAL.Status.V20250312.CidrBlock).To(Equal(pointer.MakePtr(ipalCIDR)))
				})

				By("Wait for DatabaseUser to be Ready", func() {
					Eventually(func(g Gomega) {
						g.Expect(resources.CheckResourceReady(ctx, kubeClient, testDBUser)).To(Succeed())
					}).WithContext(ctx).WithTimeout(10 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
					Expect(testDBUser.Status.V20250312).NotTo(BeNil())
					Expect(testDBUser.Status.V20250312.Username).To(Equal(username))
				})

				verifyConnectionSecret := func(targetName, clusterKind string) {
					secretName := connectionsecret.K8sConnectionSecretName(testGroup.GetName(), targetName, username, clusterKind)
					GinkgoWriter.Printf("Expected %s connection secret name: %q\n", clusterKind, secretName)
					connSecret := &corev1.Secret{}
					Eventually(func(g Gomega) {
						g.Expect(kubeClient.Get(ctx, client.ObjectKey{
							Namespace: testNamespace.Name,
							Name:      secretName,
						}, connSecret)).To(Succeed())
					}).WithContext(ctx).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

					Expect(connSecret.Data).To(HaveKey("username"))
					Expect(connSecret.Data).To(HaveKey("password"))
					Expect(connSecret.Data).To(HaveKey("connectionStringStandard"))
					Expect(connSecret.Data).To(HaveKey("connectionStringStandardSrv"))
					Expect(string(connSecret.Data["username"])).To(Equal(username))
					Expect(string(connSecret.Data["connectionStringStandard"])).To(ContainSubstring(username))
					Expect(string(connSecret.Data["connectionStringStandardSrv"])).To(ContainSubstring(username))
					Expect(connSecret.Labels).To(HaveKeyWithValue(connectionsecret.TargetLabelKey, targetName))
					Expect(connSecret.Labels).To(HaveKeyWithValue(connectionsecret.DatabaseUserLabelKey, username))
				}

				By("Verify connection secret is created for Cluster", func() {
					verifyConnectionSecret(clusterName, "cluster")
				})

				By("Verify connection secret is created for FlexCluster", func() {
					verifyConnectionSecret(flexClusterName, "flexcluster")
				})

				By("Delete DatabaseUser", func() {
					Expect(kubeClient.Delete(ctx, testDBUser)).To(Succeed())
					Eventually(func(g Gomega) {
						err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testDBUser), testDBUser)
						g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
					}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
				})

				By("Delete IPAccessListEntry", func() {
					Expect(kubeClient.Delete(ctx, testIAL)).To(Succeed())
					Eventually(func(g Gomega) {
						err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testIAL), testIAL)
						g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
					}).WithContext(ctx).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
				})

				By("Delete Cluster", func() {
					Expect(kubeClient.Delete(ctx, testCluster)).To(Succeed())
					Eventually(func(g Gomega) {
						g.Expect(resources.CheckResourceDeleted(ctx, kubeClient, testCluster)).To(Succeed())
					}).WithContext(ctx).WithTimeout(clusterDeleteTimeout).WithPolling(clusterPollingInterval).Should(Succeed())
				})

				By("Delete FlexCluster", func() {
					Expect(kubeClient.Delete(ctx, testFlexCluster)).To(Succeed())
					Eventually(func(g Gomega) {
						g.Expect(resources.CheckResourceDeleted(ctx, kubeClient, testFlexCluster)).To(Succeed())
					}).WithContext(ctx).WithTimeout(clusterDeleteTimeout).WithPolling(clusterPollingInterval).Should(Succeed())
				})

				By("Delete Group", func() {
					Expect(kubeClient.Delete(ctx, testGroup)).To(Succeed())
					Eventually(func(g Gomega) {
						err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testGroup), testGroup)
						g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
					}).WithContext(ctx).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
				})
			},
		)
	})
})
