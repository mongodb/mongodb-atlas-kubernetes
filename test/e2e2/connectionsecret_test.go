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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/operator"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/resources"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/samples"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/testparams"
)

var _ = Describe("ConnectionSecret", Ordered, Label("connectionsecret"), func() {
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
			"databaseusers.atlas.generated.mongodb.com",
			"clusters.atlas.generated.mongodb.com",
			"ipaccesslistentries.atlas.generated.mongodb.com",
		)).To(Succeed())

		_, orgID = newTestAtlasClient()
	})

	_ = BeforeEach(func() {
		testNamespace = &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("connsecret-ns-%s", rand.String(6)),
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

	createPrerequisites := func(clusterName, username string) (*generatedv1.Group, *generatedv1.Cluster, *generatedv1.DatabaseUser) {
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

		var testCluster *generatedv1.Cluster
		By("Create Cluster (M0)", func() {
			testCluster = newSharedCluster(clusterName, testNamespace.Name, testGroup.GetName())
			Expect(kubeClient.Create(ctx, testCluster)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(resources.CheckResourceReady(ctx, kubeClient, testCluster)).To(Succeed())
			}).WithContext(ctx).WithTimeout(clusterCreateTimeout).WithPolling(clusterPollingInterval).Should(Succeed())
			Expect(testCluster.Status.V20250312).NotTo(BeNil())
			Expect(testCluster.Status.V20250312.ConnectionStrings).NotTo(BeNil())
		})

		passwordSecretName := fmt.Sprintf("dbuser-pass-%s", rand.String(6))
		var testDBUser *generatedv1.DatabaseUser
		By("Create DatabaseUser", func() {
			Expect(kubeClient.Create(ctx, newPasswordSecret(testNamespace.Name, passwordSecretName))).To(Succeed())

			objs := samples.MustLoadSampleObjects("atlas_generated_v1_databaseuser_with_groupref.yaml")
			Expect(len(objs)).To(Equal(1))
			testDBUser = objs[0].(*generatedv1.DatabaseUser)
			applyTestParamsToDBUser(testDBUser, testNamespace.Name, username, testGroup.GetName(), passwordSecretName)
			Expect(kubeClient.Create(ctx, testDBUser)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(resources.CheckResourceReady(ctx, kubeClient, testDBUser)).To(Succeed())
			}).WithContext(ctx).WithTimeout(10 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
		})

		return testGroup, testCluster, testDBUser
	}

	connectionSecretName := func(group *generatedv1.Group, clusterName, username string) string {
		return connectionsecret.K8sConnectionSecretName(group.GetName(), clusterName, username, "cluster")
	}

	waitForSecret := func(secretName string) *corev1.Secret {
		connSecret := &corev1.Secret{}
		Eventually(func(g Gomega) {
			g.Expect(kubeClient.Get(ctx, client.ObjectKey{
				Namespace: testNamespace.Name,
				Name:      secretName,
			}, connSecret)).To(Succeed())
		}).WithContext(ctx).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
		return connSecret
	}

	Describe("Connection secret lifecycle", func() {
		It("Should create connection secret with valid keys when Cluster and DatabaseUser are ready", Label("focus-connectionsecret-create"), func() {
			clusterName := fmt.Sprintf("cluster-%s", rand.String(6))
			username := fmt.Sprintf("testuser-%s", rand.String(6))

			testGroup, testCluster, testDBUser := createPrerequisites(clusterName, username)

			By("Verify connection secret is created with expected keys and credentials embedded", func() {
				secretName := connectionSecretName(testGroup, clusterName, username)
				connSecret := waitForSecret(secretName)

				Expect(connSecret.Data).To(HaveKey("username"))
				Expect(connSecret.Data).To(HaveKey("password"))
				Expect(connSecret.Data).To(HaveKey("connectionStringStandard"))
				Expect(connSecret.Data).To(HaveKey("connectionStringStandardSrv"))
				Expect(string(connSecret.Data["username"])).To(Equal(username))
				Expect(string(connSecret.Data["connectionStringStandard"])).To(ContainSubstring(username))
				Expect(string(connSecret.Data["connectionStringStandardSrv"])).To(ContainSubstring(username))
			})

			By("Verify connection secret has correct labels", func() {
				secretName := connectionSecretName(testGroup, clusterName, username)
				connSecret := waitForSecret(secretName)

				Expect(connSecret.Labels).To(HaveKey(connectionsecret.TypeLabelKey))
				Expect(connSecret.Labels).To(HaveKey(connectionsecret.ProjectLabelKey))
				Expect(connSecret.Labels).To(HaveKey(connectionsecret.TargetLabelKey))
				Expect(connSecret.Labels).To(HaveKey(connectionsecret.DatabaseUserLabelKey))
				Expect(connSecret.Labels[connectionsecret.TargetLabelKey]).To(Equal(clusterName))
				Expect(connSecret.Labels[connectionsecret.DatabaseUserLabelKey]).To(Equal(username))
			})

			By("Delete DatabaseUser", func() {
				Expect(kubeClient.Delete(ctx, testDBUser)).To(Succeed())
				Eventually(func(g Gomega) {
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testDBUser), testDBUser)
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})

			By("Delete Cluster", func() {
				Expect(kubeClient.Delete(ctx, testCluster)).To(Succeed())
				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceDeleted(ctx, kubeClient, testCluster)).To(Succeed())
				}).WithContext(ctx).WithTimeout(clusterDeleteTimeout).WithPolling(clusterPollingInterval).Should(Succeed())
			})

			By("Delete Group", func() {
				Expect(kubeClient.Delete(ctx, testGroup)).To(Succeed())
				Eventually(func(g Gomega) {
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testGroup), testGroup)
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})
		})

		It("Should delete connection secret when DatabaseUser is deleted", Label("focus-connectionsecret-delete-on-user-removal"), func() {
			clusterName := fmt.Sprintf("cluster-%s", rand.String(6))
			username := fmt.Sprintf("testuser-%s", rand.String(6))

			testGroup, testCluster, testDBUser := createPrerequisites(clusterName, username)

			secretName := connectionSecretName(testGroup, clusterName, username)
			By("Wait for connection secret to be created", func() {
				waitForSecret(secretName)
			})

			By("Delete DatabaseUser: connection secret should be garbage collected", func() {
				Expect(kubeClient.Delete(ctx, testDBUser)).To(Succeed())
				Eventually(func(g Gomega) {
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testDBUser), testDBUser)
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

				Eventually(func(g Gomega) {
					connSecret := &corev1.Secret{}
					err := kubeClient.Get(ctx, client.ObjectKey{
						Namespace: testNamespace.Name,
						Name:      secretName,
					}, connSecret)
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})

			By("Delete Cluster", func() {
				Expect(kubeClient.Delete(ctx, testCluster)).To(Succeed())
				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceDeleted(ctx, kubeClient, testCluster)).To(Succeed())
				}).WithContext(ctx).WithTimeout(clusterDeleteTimeout).WithPolling(clusterPollingInterval).Should(Succeed())
			})

			By("Delete Group", func() {
				Expect(kubeClient.Delete(ctx, testGroup)).To(Succeed())
				Eventually(func(g Gomega) {
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testGroup), testGroup)
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})
		})

		It("Should delete connection secret when Cluster is deleted", Label("focus-connectionsecret-delete-on-cluster-removal"), func() {
			clusterName := fmt.Sprintf("cluster-%s", rand.String(6))
			username := fmt.Sprintf("testuser-%s", rand.String(6))

			testGroup, testCluster, testDBUser := createPrerequisites(clusterName, username)

			secretName := connectionSecretName(testGroup, clusterName, username)
			By("Wait for connection secret to be created", func() {
				waitForSecret(secretName)
			})

			By("Delete Cluster: stale connection secret should be removed", func() {
				Expect(kubeClient.Delete(ctx, testCluster)).To(Succeed())
				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceDeleted(ctx, kubeClient, testCluster)).To(Succeed())
				}).WithContext(ctx).WithTimeout(clusterDeleteTimeout).WithPolling(clusterPollingInterval).Should(Succeed())

				Eventually(func(g Gomega) {
					connSecret := &corev1.Secret{}
					err := kubeClient.Get(ctx, client.ObjectKey{
						Namespace: testNamespace.Name,
						Name:      secretName,
					}, connSecret)
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})

			By("Delete DatabaseUser", func() {
				Expect(kubeClient.Delete(ctx, testDBUser)).To(Succeed())
				Eventually(func(g Gomega) {
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testDBUser), testDBUser)
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})

			By("Delete Group", func() {
				Expect(kubeClient.Delete(ctx, testGroup)).To(Succeed())
				Eventually(func(g Gomega) {
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testGroup), testGroup)
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})
		})

		It("Should only create connection secret for Cluster within DatabaseUser scopes", Label("focus-connectionsecret-scope-filter"), func() {
			clusterName1 := fmt.Sprintf("cluster-%s", rand.String(6))
			clusterName2 := fmt.Sprintf("cluster-%s", rand.String(6))
			username := fmt.Sprintf("testuser-%s", rand.String(6))

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

			var testCluster1, testCluster2 *generatedv1.Cluster
			By("Create two Clusters (M0)", func() {
				testCluster1 = newSharedCluster(clusterName1, testNamespace.Name, testGroup.GetName())
				Expect(kubeClient.Create(ctx, testCluster1)).To(Succeed())

				testCluster2 = newDedicatedCluster(clusterName2, testNamespace.Name, testGroup.GetName())
				Expect(kubeClient.Create(ctx, testCluster2)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, testCluster1)).To(Succeed())
				}).WithContext(ctx).WithTimeout(clusterCreateTimeout).WithPolling(clusterPollingInterval).Should(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, testCluster2)).To(Succeed())
				}).WithContext(ctx).WithTimeout(clusterCreateTimeout).WithPolling(clusterPollingInterval).Should(Succeed())
			})

			passwordSecretName := fmt.Sprintf("dbuser-pass-%s", rand.String(6))
			var testDBUser *generatedv1.DatabaseUser
			By("Create DatabaseUser scoped to clusterName1 only", func() {
				Expect(kubeClient.Create(ctx, newPasswordSecret(testNamespace.Name, passwordSecretName))).To(Succeed())

				// This user is for Cluster1 only, so the secret for Cluster2 must not be created
				testDBUser = newDBUserWithScopes(testNamespace.Name, username, testGroup.GetName(), passwordSecretName, []generatedv1.Scopes{
					{Name: clusterName1, Type: "CLUSTER"},
				})
				Expect(kubeClient.Create(ctx, testDBUser)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, testDBUser)).To(Succeed())
				}).WithContext(ctx).WithTimeout(10 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})

			By("Verify connection secret is created only for clusterName1, not for clusterName2", func() {
				// Secret for the in-scope cluster should exist
				secretName1 := connectionSecretName(testGroup, clusterName1, username)
				waitForSecret(secretName1)

				// Secret for the out-of-scope cluster should never appear
				secretName2 := connectionSecretName(testGroup, clusterName2, username)
				Consistently(func(g Gomega) {
					connSecret := &corev1.Secret{}
					err := kubeClient.Get(ctx, client.ObjectKey{
						Namespace: testNamespace.Name,
						Name:      secretName2,
					}, connSecret)
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(30 * time.Second).WithPolling(5 * time.Second).Should(Succeed())
			})

			By("Delete DatabaseUser", func() {
				Expect(kubeClient.Delete(ctx, testDBUser)).To(Succeed())
				Eventually(func(g Gomega) {
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testDBUser), testDBUser)
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})

			By("Delete both Clusters", func() {
				Expect(kubeClient.Delete(ctx, testCluster1)).To(Succeed())
				Expect(kubeClient.Delete(ctx, testCluster2)).To(Succeed())
				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceDeleted(ctx, kubeClient, testCluster1)).To(Succeed())
				}).WithContext(ctx).WithTimeout(clusterDeleteTimeout).WithPolling(clusterPollingInterval).Should(Succeed())
				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceDeleted(ctx, kubeClient, testCluster2)).To(Succeed())
				}).WithContext(ctx).WithTimeout(clusterDeleteTimeout).WithPolling(clusterPollingInterval).Should(Succeed())
			})

			By("Delete Group", func() {
				Expect(kubeClient.Delete(ctx, testGroup)).To(Succeed())
				Eventually(func(g Gomega) {
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testGroup), testGroup)
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})
		})

		It("Should update connection secret when DatabaseUser password is rotated", Label("focus-connectionsecret-password-rotation"), func() {
			clusterName := fmt.Sprintf("cluster-%s", rand.String(6))
			username := fmt.Sprintf("testuser-%s", rand.String(6))

			testGroup, testCluster, testDBUser := createPrerequisites(clusterName, username)

			secretName := connectionSecretName(testGroup, clusterName, username)
			var initialResourceVersion string
			By("Wait for initial connection secret to be created", func() {
				connSecret := waitForSecret(secretName)
				initialResourceVersion = connSecret.ResourceVersion
				Expect(string(connSecret.Data["password"])).To(Equal("Passw0rd!"))
			})

			By("Rotate DatabaseUser password", func() {
				Expect(kubeClient.Get(ctx, client.ObjectKeyFromObject(testDBUser), testDBUser)).To(Succeed())
				passwordSecretName := testDBUser.Spec.V20250312.Entry.PasswordSecretRef.Name

				passwordSecret := &corev1.Secret{}
				Expect(kubeClient.Get(ctx, client.ObjectKey{
					Namespace: testNamespace.Name,
					Name:      passwordSecretName,
				}, passwordSecret)).To(Succeed())
				passwordSecret.StringData = map[string]string{"password": "NewPassw0rd!"}
				Expect(kubeClient.Update(ctx, passwordSecret)).To(Succeed())
			})

			By("Verify connection secret is updated with new password", func() {
				Eventually(func(g Gomega) {
					connSecret := &corev1.Secret{}
					g.Expect(kubeClient.Get(ctx, client.ObjectKey{
						Namespace: testNamespace.Name,
						Name:      secretName,
					}, connSecret)).To(Succeed())
					g.Expect(connSecret.ResourceVersion).NotTo(Equal(initialResourceVersion))
					g.Expect(string(connSecret.Data["password"])).To(Equal("NewPassw0rd!"))
				}).WithContext(ctx).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})

			By("Delete DatabaseUser", func() {
				Expect(kubeClient.Delete(ctx, testDBUser)).To(Succeed())
				Eventually(func(g Gomega) {
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testDBUser), testDBUser)
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})

			By("Delete Cluster", func() {
				Expect(kubeClient.Delete(ctx, testCluster)).To(Succeed())
				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceDeleted(ctx, kubeClient, testCluster)).To(Succeed())
				}).WithContext(ctx).WithTimeout(clusterDeleteTimeout).WithPolling(clusterPollingInterval).Should(Succeed())
			})

			By("Delete Group", func() {
				Expect(kubeClient.Delete(ctx, testGroup)).To(Succeed())
				Eventually(func(g Gomega) {
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testGroup), testGroup)
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})
		})

		It("Should create separate connection secrets for each Cluster in the same project", Label("focus-connectionsecret-multi-cluster"), func() {
			clusterName1 := fmt.Sprintf("cluster-%s", rand.String(6))
			clusterName2 := fmt.Sprintf("cluster-%s", rand.String(6))
			username := fmt.Sprintf("testuser-%s", rand.String(6))

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

			var testCluster1, testCluster2 *generatedv1.Cluster
			By("Create two Clusters (M0)", func() {
				testCluster1 = newSharedCluster(clusterName1, testNamespace.Name, testGroup.GetName())
				Expect(kubeClient.Create(ctx, testCluster1)).To(Succeed())

				testCluster2 = newDedicatedCluster(clusterName2, testNamespace.Name, testGroup.GetName())
				Expect(kubeClient.Create(ctx, testCluster2)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, testCluster1)).To(Succeed())
				}).WithContext(ctx).WithTimeout(clusterCreateTimeout).WithPolling(clusterPollingInterval).Should(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, testCluster2)).To(Succeed())
				}).WithContext(ctx).WithTimeout(clusterCreateTimeout).WithPolling(clusterPollingInterval).Should(Succeed())
			})

			passwordSecretName := fmt.Sprintf("dbuser-pass-%s", rand.String(6))
			var testDBUser *generatedv1.DatabaseUser
			By("Create DatabaseUser without scopes (applies to all clusters)", func() {
				Expect(kubeClient.Create(ctx, newPasswordSecret(testNamespace.Name, passwordSecretName))).To(Succeed())

				objs := samples.MustLoadSampleObjects("atlas_generated_v1_databaseuser_with_groupref.yaml")
				Expect(len(objs)).To(Equal(1))
				testDBUser = objs[0].(*generatedv1.DatabaseUser)
				applyTestParamsToDBUser(testDBUser, testNamespace.Name, username, testGroup.GetName(), passwordSecretName)
				Expect(kubeClient.Create(ctx, testDBUser)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, testDBUser)).To(Succeed())
				}).WithContext(ctx).WithTimeout(10 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})

			By("Verify connection secrets are created for both clusters", func() {
				secretName1 := connectionSecretName(testGroup, clusterName1, username)
				secretName2 := connectionSecretName(testGroup, clusterName2, username)

				waitForSecret(secretName1)
				waitForSecret(secretName2)

				connSecret1 := &corev1.Secret{}
				Expect(kubeClient.Get(ctx, client.ObjectKey{Namespace: testNamespace.Name, Name: secretName1}, connSecret1)).To(Succeed())
				connSecret2 := &corev1.Secret{}
				Expect(kubeClient.Get(ctx, client.ObjectKey{Namespace: testNamespace.Name, Name: secretName2}, connSecret2)).To(Succeed())

				Expect(connSecret1.Labels[connectionsecret.TargetLabelKey]).To(Equal(clusterName1))
				Expect(connSecret2.Labels[connectionsecret.TargetLabelKey]).To(Equal(clusterName2))
			})

			By("Delete DatabaseUser", func() {
				Expect(kubeClient.Delete(ctx, testDBUser)).To(Succeed())
				Eventually(func(g Gomega) {
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testDBUser), testDBUser)
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})

			By("Delete both Clusters", func() {
				Expect(kubeClient.Delete(ctx, testCluster1)).To(Succeed())
				Expect(kubeClient.Delete(ctx, testCluster2)).To(Succeed())
				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceDeleted(ctx, kubeClient, testCluster1)).To(Succeed())
				}).WithContext(ctx).WithTimeout(clusterDeleteTimeout).WithPolling(clusterPollingInterval).Should(Succeed())
				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceDeleted(ctx, kubeClient, testCluster2)).To(Succeed())
				}).WithContext(ctx).WithTimeout(clusterDeleteTimeout).WithPolling(clusterPollingInterval).Should(Succeed())
			})

			By("Delete Group", func() {
				Expect(kubeClient.Delete(ctx, testGroup)).To(Succeed())
				Eventually(func(g Gomega) {
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testGroup), testGroup)
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})
		})
	})

	Describe("FlexCluster connection secret lifecycle", func() {
		It("Should create connection secret when FlexCluster and DatabaseUser are ready", Label("focus-connectionsecret-flexcluster"), func() {
			clusterName := fmt.Sprintf("flexy-%s", rand.String(6))
			username := fmt.Sprintf("testuser-%s", rand.String(6))

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
			})

			var testFlexCluster *generatedv1.FlexCluster
			By("Create FlexCluster", func() {
				testFlexCluster = newSharedFlexCluster(clusterName, testNamespace.Name, testGroup.GetName())
				Expect(kubeClient.Create(ctx, testFlexCluster)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, testFlexCluster)).To(Succeed())
				}).WithContext(ctx).WithTimeout(clusterCreateTimeout).WithPolling(clusterPollingInterval).Should(Succeed())
				Expect(testFlexCluster.Status.V20250312).NotTo(BeNil())
				Expect(testFlexCluster.Status.V20250312.ConnectionStrings).NotTo(BeNil())
			})

			passwordSecretName := fmt.Sprintf("dbuser-pass-%s", rand.String(6))
			var testDBUser *generatedv1.DatabaseUser
			By("Create DatabaseUser", func() {
				Expect(kubeClient.Create(ctx, newPasswordSecret(testNamespace.Name, passwordSecretName))).To(Succeed())

				objs := samples.MustLoadSampleObjects("atlas_generated_v1_databaseuser_with_groupref.yaml")
				Expect(len(objs)).To(Equal(1))
				testDBUser = objs[0].(*generatedv1.DatabaseUser)
				applyTestParamsToDBUser(testDBUser, testNamespace.Name, username, testGroup.GetName(), passwordSecretName)
				Expect(kubeClient.Create(ctx, testDBUser)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, testDBUser)).To(Succeed())
				}).WithContext(ctx).WithTimeout(10 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})

			By("Verify connection secret is created for FlexCluster", func() {
				secretName := connectionsecret.K8sConnectionSecretName(testGroup.GetName(), clusterName, username, "flexcluster")
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
			})

			By("Verify connection secret is deleted when DatabaseUser is deleted", func() {
				secretName := connectionsecret.K8sConnectionSecretName(testGroup.GetName(), clusterName, username, "flexcluster")
				Expect(kubeClient.Delete(ctx, testDBUser)).To(Succeed())

				Eventually(func(g Gomega) {
					err := kubeClient.Get(ctx, client.ObjectKey{
						Namespace: testNamespace.Name,
						Name:      secretName,
					}, &corev1.Secret{})
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
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
		})
	})
})

func newDBUserWithScopes(namespace, username, groupRefName, passwordSecretName string, scopes []generatedv1.Scopes) *generatedv1.DatabaseUser {
	dbUser := &generatedv1.DatabaseUser{
		ObjectMeta: metav1.ObjectMeta{
			Name:      username,
			Namespace: namespace,
		},
		Spec: generatedv1.DatabaseUserSpec{
			V20250312: &generatedv1.DatabaseUserSpecV20250312{
				GroupRef: &k8s.LocalReference{Name: groupRefName},
				Entry: &generatedv1.DatabaseUserSpecV20250312Entry{
					Username:     username,
					DatabaseName: "admin",
					PasswordSecretRef: &generatedv1.PasswordSecretRef{
						Name: passwordSecretName,
					},
					Roles: []generatedv1.Roles{
						{RoleName: "readAnyDatabase", DatabaseName: "admin"},
					},
				},
			},
		},
	}
	if len(scopes) > 0 {
		dbUser.Spec.V20250312.Entry.Scopes = &scopes
	}
	return dbUser
}

// The reason for this is that Atlas only allows only one M0 cluster per project
func newDedicatedCluster(name, namespace, groupRefName string) *generatedv1.Cluster {
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
					Name:        new(name),
					ClusterType: new("REPLICASET"),
					ReplicationSpecs: &[]generatedv1.ReplicationSpecs{
						{
							RegionConfigs: &[]generatedv1.RegionConfigs{
								{
									ProviderName: new("AWS"),
									RegionName:   new("US_EAST_1"),
									Priority:     new(7),
									ElectableSpecs: &generatedv1.ElectableSpecs{
										InstanceSize: new("M10"),
										NodeCount:    new(3),
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

func newSharedFlexCluster(name, namespace, groupRefName string) *generatedv1.FlexCluster {
	return &generatedv1.FlexCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: generatedv1.FlexClusterSpec{
			V20250312: &generatedv1.FlexClusterSpecV20250312{
				GroupRef: &k8s.LocalReference{
					Name: groupRefName,
				},
				Entry: &generatedv1.FlexClusterSpecV20250312Entry{
					Name: name,
					ProviderSettings: generatedv1.ProviderSettings{
						BackingProviderName: "AWS",
						RegionName:          "US_EAST_1",
					},
				},
			},
		},
	}
}
