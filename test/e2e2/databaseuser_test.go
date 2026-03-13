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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/secretservice"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/httputil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/operator"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/resources"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/samples"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/testparams"
)

var _ = Describe("DatabaseUser CRUD", Ordered, Label("databaseuser"), func() {
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
			"databaseusers.atlas.generated.mongodb.com",
			"clusters.atlas.generated.mongodb.com",
		)).To(Succeed())

		atlasClient, orgID = newTestAtlasClient()
	})

	_ = BeforeEach(func() {
		testNamespace = &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("dbuser-ns-%s", rand.String(6)),
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

	Describe("DatabaseUser CRUD lifecycle", func() {
		It("Should create, update roles, and delete DatabaseUser", Label("focus-databaseuser-crud"), func() {
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

			username := fmt.Sprintf("testuser-%s", rand.String(6))
			passwordSecretName := fmt.Sprintf("dbuser-pass-%s", rand.String(6))

			var testDBUser *generatedv1.DatabaseUser
			By("Create DatabaseUser with groupRef", func() {
				Expect(kubeClient.Create(ctx, newPasswordSecret(testNamespace.Name, passwordSecretName))).To(Succeed())

				objs := samples.MustLoadSampleObjects("atlas_generated_v1_databaseuser_with_groupref.yaml")
				Expect(len(objs)).To(Equal(1))
				testDBUser = objs[0].(*generatedv1.DatabaseUser)
				applyTestParamsToDBUser(testDBUser, testNamespace.Name, username, testGroup.GetName(), passwordSecretName)
				Expect(kubeClient.Create(ctx, testDBUser)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, testDBUser)).To(Succeed())
				}).WithContext(ctx).WithTimeout(10 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

				Expect(testDBUser.Status.V20250312).NotTo(BeNil())
				Expect(testDBUser.Status.V20250312.GroupId).NotTo(BeEmpty())
				Expect(testDBUser.Status.V20250312.DatabaseName).To(Equal("admin"))
				Expect(testDBUser.Status.V20250312.Username).To(Equal(username))
				Expect(meta.IsStatusConditionTrue(testDBUser.GetConditions(), "Ready")).To(BeTrue())

				groupID := *testGroup.Status.V20250312.Id
				atlasUser, _, err := atlasClient.DatabaseUsersApi.GetDatabaseUser(ctx, groupID, "admin", username).Execute()
				Expect(err).ToNot(HaveOccurred())
				Expect(atlasUser.GetUsername()).To(Equal(username))
			})

			By("Update DatabaseUser roles", func() {
				Expect(kubeClient.Get(ctx, client.ObjectKeyFromObject(testDBUser), testDBUser)).To(Succeed())
				testDBUser.Spec.V20250312.Entry.Roles = &[]generatedv1.Roles{
					{RoleName: "readWriteAnyDatabase", DatabaseName: "admin"},
				}
				Expect(kubeClient.Update(ctx, testDBUser)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceUpdated(ctx, kubeClient, testDBUser)).To(Succeed())
				}).WithContext(ctx).WithTimeout(10 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

				groupID := *testGroup.Status.V20250312.Id
				atlasUser, _, err := atlasClient.DatabaseUsersApi.GetDatabaseUser(ctx, groupID, "admin", username).Execute()
				Expect(err).ToNot(HaveOccurred())
				atlasRoles := atlasUser.GetRoles()
				Expect(len(atlasRoles)).To(Equal(1))
				Expect(atlasRoles[0].GetRoleName()).To(Equal("readWriteAnyDatabase"))
			})

			By("Delete DatabaseUser from cluster - should delete from Atlas", func() {
				groupID := testDBUser.Status.V20250312.GroupId
				databaseName := testDBUser.Status.V20250312.DatabaseName
				Expect(kubeClient.Delete(ctx, testDBUser)).To(Succeed())

				Eventually(func(g Gomega) {
					dbUser := &generatedv1.DatabaseUser{}
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testDBUser), dbUser)
					g.Expect(err).NotTo(Succeed())
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

				Eventually(func(g Gomega) {
					_, r, err := atlasClient.DatabaseUsersApi.GetDatabaseUser(ctx, groupID, databaseName, username).Execute()
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

		It("Should create connection secret when Cluster and DatabaseUser are ready", Label("focus-databaseuser-connection-secret"), func() {
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

			clusterName := fmt.Sprintf("cluster-%s", rand.String(6))
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

			username := fmt.Sprintf("testuser-%s", rand.String(6))
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

			By("Verify connection secret is created with correct keys", func() {
				projectName := testGroup.GetName()
				secretName := kube.NormalizeIdentifier(
					fmt.Sprintf("%s-%s-%s", projectName, clusterName, username),
				)
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
			})

			By("Delete DatabaseUser", func() {
				Expect(kubeClient.Delete(ctx, testDBUser)).To(Succeed())
				Eventually(func(g Gomega) {
					dbUser := &generatedv1.DatabaseUser{}
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testDBUser), dbUser)
					g.Expect(err).NotTo(Succeed())
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})

			By("Delete Cluster", func() {
				Expect(kubeClient.Delete(ctx, testCluster)).To(Succeed())
				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceDeleted(ctx, kubeClient, testCluster)).To(Succeed())
				}).WithContext(ctx).WithTimeout(clusterDeleteTimeout).WithPolling(clusterPollingInterval).Should(Succeed())
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

		It("Should fail if password Secret is missing", Label("focus-databaseuser-fail-secret"), func() {
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

			username := fmt.Sprintf("testuser-%s", rand.String(6))

			var testDBUser *generatedv1.DatabaseUser
			By("Create DatabaseUser with non-existent password secret", func() {
				objs := samples.MustLoadSampleObjects("atlas_generated_v1_databaseuser_with_groupref.yaml")
				Expect(len(objs)).To(Equal(1))
				testDBUser = objs[0].(*generatedv1.DatabaseUser)
				applyTestParamsToDBUser(testDBUser, testNamespace.Name, username, testGroup.GetName(), "non-existent-secret")
				Expect(kubeClient.Create(ctx, testDBUser)).To(Succeed())
			})

			By("Wait for DatabaseUser to fail", func() {
				dbUser := &generatedv1.DatabaseUser{
					ObjectMeta: metav1.ObjectMeta{Name: username, Namespace: testNamespace.Name},
				}
				Eventually(func(g Gomega) {
					g.Expect(kubeClient.Get(ctx, client.ObjectKeyFromObject(dbUser), dbUser)).To(Succeed())
					g.Expect(dbUser.GetConditions()).NotTo(BeEmpty())
				}).WithContext(ctx).WithTimeout(30 * time.Second).WithPolling(time.Second).Should(Succeed())
				testDBUser = dbUser
				Expect(meta.IsStatusConditionTrue(testDBUser.GetConditions(), "Ready")).To(BeFalse())
				readyCondition := meta.FindStatusCondition(testDBUser.GetConditions(), "Ready")
				Expect(readyCondition.Reason).To(Equal("Error"))
			})

			By("Force delete DatabaseUser by removing finalizers", func() {
				Expect(kubeClient.Get(ctx, client.ObjectKeyFromObject(testDBUser), testDBUser)).To(Succeed())
				testDBUser.SetFinalizers([]string{})
				Expect(kubeClient.Update(ctx, testDBUser)).To(Succeed())
				Expect(kubeClient.Delete(ctx, testDBUser)).To(Succeed())
				Eventually(func(g Gomega) {
					dbUser := &generatedv1.DatabaseUser{}
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testDBUser), dbUser)
					g.Expect(err).NotTo(Succeed())
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
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

		It("Should NOT delete from Atlas when ResourcePolicyKeep annotation is set", Label("focus-databaseuser-kept"), func() {
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

			username := fmt.Sprintf("testuser-%s", rand.String(6))
			passwordSecretName := fmt.Sprintf("dbuser-pass-%s", rand.String(6))

			var testDBUser *generatedv1.DatabaseUser
			By("Create DatabaseUser with keep annotation", func() {
				Expect(kubeClient.Create(ctx, newPasswordSecret(testNamespace.Name, passwordSecretName))).To(Succeed())

				objs := samples.MustLoadSampleObjects("atlas_generated_v1_databaseuser_with_groupref.yaml")
				Expect(len(objs)).To(Equal(1))
				testDBUser = objs[0].(*generatedv1.DatabaseUser)
				applyTestParamsToDBUser(testDBUser, testNamespace.Name, username, testGroup.GetName(), passwordSecretName)
				testDBUser.SetAnnotations(map[string]string{
					customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep,
				})
				Expect(kubeClient.Create(ctx, testDBUser)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, testDBUser)).To(Succeed())
				}).WithContext(ctx).WithTimeout(10 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
				Expect(testDBUser.Status.V20250312).NotTo(BeNil())
			})

			By("Delete DatabaseUser from cluster - should NOT delete from Atlas", func() {
				groupID := testDBUser.Status.V20250312.GroupId
				databaseName := testDBUser.Status.V20250312.DatabaseName
				Expect(kubeClient.Delete(ctx, testDBUser)).To(Succeed())

				Eventually(func(g Gomega) {
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testDBUser), testDBUser)
					g.Expect(err).ToNot(Succeed())
				}).WithContext(ctx).WithTimeout(time.Minute).WithPolling(time.Second).Should(Succeed())

				Eventually(func(g Gomega) {
					atlasUser, _, err := atlasClient.DatabaseUsersApi.GetDatabaseUser(ctx, groupID, databaseName, username).Execute()
					g.Expect(err).ToNot(HaveOccurred())
					g.Expect(atlasUser).ToNot(BeNil())
				}).WithContext(ctx).WithTimeout(time.Minute).WithPolling(time.Second).Should(Succeed())

				_, err := atlasClient.DatabaseUsersApi.DeleteDatabaseUser(ctx, groupID, databaseName, username).Execute()
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
	})

	Describe("Importing existing DatabaseUser from Atlas", func() {
		It("Should reconcile existing Atlas user", Label("focus-databaseuser-existing"), func() {
			groupName := fmt.Sprintf("test-group-%s", rand.String(6))
			groupParams := testparams.New(orgID, control.MustEnvVar("OPERATOR_NAMESPACE"), DefaultGlobalCredentials).
				WithGroupName(groupName).
				WithNamespace(testNamespace.Name)

			var testGroup *generatedv1.Group
			var atlasGroupID string
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
				atlasGroupID = *testGroup.Status.V20250312.Id
			})

			username := fmt.Sprintf("testuser-%s", rand.String(6))
			passwordSecretName := fmt.Sprintf("dbuser-pass-%s", rand.String(6))

			By("Create user in Atlas directly", func() {
				atlasDBUser := &admin.CloudDatabaseUser{
					GroupId:      atlasGroupID,
					DatabaseName: "admin",
					Username:     username,
					Password:     pointer.MakePtr("Passw0rd!"),
					Roles: &[]admin.DatabaseUserRole{
						{RoleName: "readAnyDatabase", DatabaseName: "admin"},
					},
				}
				_, _, err := atlasClient.DatabaseUsersApi.CreateDatabaseUser(ctx, atlasGroupID, atlasDBUser).Execute()
				Expect(err).ToNot(HaveOccurred())
			})

			var testDBUser *generatedv1.DatabaseUser
			By("Create DatabaseUser K8s resource with import annotation", func() {
				Expect(kubeClient.Create(ctx, newPasswordSecret(testNamespace.Name, passwordSecretName))).To(Succeed())

				objs := samples.MustLoadSampleObjects("atlas_generated_v1_databaseuser_with_groupref.yaml")
				Expect(len(objs)).To(Equal(1))
				testDBUser = objs[0].(*generatedv1.DatabaseUser)
				applyTestParamsToDBUser(testDBUser, testNamespace.Name, username, testGroup.GetName(), passwordSecretName)
				testDBUser.SetAnnotations(map[string]string{
					"mongodb.com/external-id": fmt.Sprintf("%s:admin:%s", atlasGroupID, username),
				})
				Expect(kubeClient.Create(ctx, testDBUser)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, testDBUser)).To(Succeed())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

				Expect(testDBUser.Status.V20250312).NotTo(BeNil())
				Expect(testDBUser.Status.V20250312.GroupId).To(Equal(atlasGroupID))
				Expect(testDBUser.Status.V20250312.Username).To(Equal(username))

				atlasUser, _, err := atlasClient.DatabaseUsersApi.GetDatabaseUser(ctx, atlasGroupID, "admin", username).Execute()
				Expect(err).ToNot(HaveOccurred())
				Expect(atlasUser.GetUsername()).To(Equal(username))

				Expect(meta.IsStatusConditionTrue(testDBUser.GetConditions(), "Ready")).To(BeTrue())
				Expect(meta.IsStatusConditionTrue(testDBUser.GetConditions(), "State")).To(BeTrue())
			})

			By("Delete DatabaseUser from cluster", func() {
				Expect(kubeClient.Delete(ctx, testDBUser)).To(Succeed())

				Eventually(func(g Gomega) {
					dbUser := &generatedv1.DatabaseUser{}
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testDBUser), dbUser)
					g.Expect(err).NotTo(Succeed())
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
	})
})

var _ = Describe("DatabaseUser with Deletion Protection", Ordered, Label("databaseuser"), func() {
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
			"databaseusers.atlas.generated.mongodb.com",
		)).To(Succeed())

		atlasClient, orgID = newTestAtlasClient()
	})

	_ = BeforeEach(func() {
		testNamespace = &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("dbuser-protect-ns-%s", rand.String(6)),
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

	Describe("Deleting the DatabaseUser", Label("focus-databaseuser-deletion-protected"), func() {
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

			username := fmt.Sprintf("testuser-%s", rand.String(6))
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

				Expect(testDBUser.Status.V20250312).NotTo(BeNil())
				Expect(meta.IsStatusConditionTrue(testDBUser.GetConditions(), "Ready")).To(BeTrue())
			})

			By("Delete DatabaseUser from cluster - should NOT delete from Atlas", func() {
				groupID := testDBUser.Status.V20250312.GroupId
				databaseName := testDBUser.Status.V20250312.DatabaseName
				Expect(kubeClient.Delete(ctx, testDBUser)).To(Succeed())

				Eventually(func(g Gomega) {
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(testDBUser), testDBUser)
					g.Expect(err).To(HaveOccurred())
					g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
				}).WithContext(ctx).WithTimeout(time.Minute).WithPolling(time.Second).Should(Succeed())

				Eventually(func(g Gomega) {
					atlasUser, _, err := atlasClient.DatabaseUsersApi.GetDatabaseUser(ctx, groupID, databaseName, username).Execute()
					g.Expect(err).ToNot(HaveOccurred())
					g.Expect(atlasUser).NotTo(BeNil())
					g.Expect(atlasUser.GetUsername()).To(Equal(username))
				}).WithContext(ctx).WithTimeout(30 * time.Second).WithPolling(5 * time.Second).Should(Succeed())
			})

			By("Clean up Atlas user manually", func() {
				groupID := testDBUser.Status.V20250312.GroupId
				databaseName := testDBUser.Status.V20250312.DatabaseName
				_, err := atlasClient.DatabaseUsersApi.DeleteDatabaseUser(ctx, groupID, databaseName, username).Execute()
				Expect(err).ToNot(HaveOccurred())

				Eventually(func(g Gomega) {
					_, r, err := atlasClient.DatabaseUsersApi.GetDatabaseUser(ctx, groupID, databaseName, username).Execute()
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

// applyTestParamsToDBUser configures a DatabaseUser object for a test run.
func applyTestParamsToDBUser(dbUser *generatedv1.DatabaseUser, namespace, username, groupRefName, passwordSecretName string) {
	dbUser.SetNamespace(namespace)
	dbUser.SetName(username)

	if dbUser.Spec.V20250312 == nil {
		dbUser.Spec.V20250312 = &generatedv1.DatabaseUserSpecV20250312{}
	}
	dbUser.Spec.V20250312.GroupRef = &k8s.LocalReference{Name: groupRefName}
	dbUser.Spec.V20250312.GroupId = nil

	if dbUser.Spec.V20250312.Entry == nil {
		dbUser.Spec.V20250312.Entry = &generatedv1.DatabaseUserV20250312Entry{}
	}
	dbUser.Spec.V20250312.Entry.Username = username
	dbUser.Spec.V20250312.Entry.DatabaseName = "admin"
	dbUser.Spec.V20250312.Entry.PasswordSecretRef = &generatedv1.PasswordSecretRef{Name: passwordSecretName}
}

// newPasswordSecret creates a Kubernetes Secret containing a MongoDB user password.
func newPasswordSecret(namespace, name string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				secretservice.TypeLabelKey: secretservice.CredLabelVal,
			},
		},
		StringData: map[string]string{
			"password": "Passw0rd!",
		},
	}
}
