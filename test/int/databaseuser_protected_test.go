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

package int

import (
	"context"
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas-sdk/v20250312014/admin"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/httputil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/conditions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/resources"
)

var _ = Describe("Atlas Database User", Label("int", "AtlasDatabaseUser", "focus-protection-enabled"), func() {
	var testNamespace *corev1.Namespace
	var stopManager context.CancelFunc
	var projectName string
	projectNamePrefix := "database-user-protected"
	dbUserName1 := "db-user1"
	dbUserName2 := "db-user2"
	dbUserName3 := "db-user3"
	testProject := &akov2.AtlasProject{}
	testDeployment := &akov2.AtlasDeployment{}
	testDBUser1 := &akov2.AtlasDatabaseUser{}
	testDBUser2 := &akov2.AtlasDatabaseUser{}
	testDBUser3 := &akov2.AtlasDatabaseUser{}

	BeforeEach(func() {
		testNamespace, stopManager = prepareControllers(true)
		projectName = fmt.Sprintf("%s-%s", projectNamePrefix, testNamespace.Name)

		By("Creating a project", func() {
			connSecret := buildConnectionSecret("my-atlas-key")
			Expect(k8sClient.Create(context.Background(), &connSecret)).To(Succeed())

			testProject = akov2.NewProject(testNamespace.Name, projectName, projectName).
				WithConnectionSecret(connSecret.Name).
				WithIPAccessList(project.NewIPAccessList().WithCIDR("0.0.0.0/0"))
			Expect(k8sClient.Create(context.Background(), testProject, &client.CreateOptions{})).To(Succeed())

			Eventually(func() bool {
				return resources.CheckCondition(k8sClient, testProject, api.TrueCondition(api.ReadyType))
			}).WithTimeout(15 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
		})

		By("Creating a deployment", func() {
			testDeployment = akov2.NewDefaultAWSFlexInstance(testNamespace.Name, projectName)
			customresource.SetAnnotation( // this test deployment must be deleted
				testDeployment,
				customresource.ResourcePolicyAnnotation,
				customresource.ResourcePolicyDelete,
			)
			Expect(k8sClient.Create(context.Background(), testDeployment)).To(Succeed())

			Eventually(func() bool {
				return resources.CheckCondition(k8sClient, testDeployment, api.TrueCondition(api.ReadyType))
			}).WithTimeout(20 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
		})

		By("Creating database user", func() {
			dbUser := admin.NewCloudDatabaseUser("admin", testProject.ID(), dbUserName3)
			dbUser.SetPassword("mypass")
			dbUser.SetRoles(
				[]admin.DatabaseUserRole{
					{
						RoleName:     "readAnyDatabase",
						DatabaseName: "admin",
					},
				},
			)
			_, _, err := atlasClient.DatabaseUsersApi.CreateDatabaseUser(context.Background(), testProject.ID(), dbUser).Execute()
			Expect(err).To(BeNil())
		})
	})

	Describe("Operator is running with deletion protection enabled", func() {
		It("Adds database users and protect them to be deleted when operator doesn't own resource", Label("focus-unowned-protected"), func() {
			By("First without setting atlas-resource-policy annotation", func() {
				passwordSecret := buildPasswordSecret(testNamespace.Name, UserPasswordSecret, DBUserPassword)
				Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())

				testDBUser1 = akov2.NewDBUser(testNamespace.Name, dbUserName1, dbUserName1, projectName).
					WithPasswordSecret(UserPasswordSecret).
					WithRole("readWriteAnyDatabase", "admin", "")
				Expect(k8sClient.Create(context.Background(), testDBUser1)).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testDBUser1, api.TrueCondition(api.ReadyType))
				}).WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())

				validateSecret(k8sClient, *testProject, *testDeployment, *testDBUser1)

				Expect(tryConnect(testProject.ID(), *testDeployment, *testDBUser1)).Should(Succeed())
			})

			// nolint:dupl
			By("Second setting atlas-resource-policy annotation to delete", func() {
				passwordSecret := buildPasswordSecret(testNamespace.Name, UserPasswordSecret2, DBUserPassword2)
				Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())

				testDBUser2 = akov2.NewDBUser(testNamespace.Name, dbUserName2, dbUserName2, projectName).
					WithPasswordSecret(UserPasswordSecret2).
					WithRole("readWriteAnyDatabase", "admin", "")
				testDBUser2.SetAnnotations(map[string]string{customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyDelete})
				Expect(k8sClient.Create(context.Background(), testDBUser2)).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testDBUser2, api.TrueCondition(api.ReadyType))
				}).WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())

				validateSecret(k8sClient, *testProject, *testDeployment, *testDBUser2)

				Expect(tryConnect(testProject.ID(), *testDeployment, *testDBUser2)).Should(Succeed())
			})

			By("Third previously added in Atlas", func() {
				passwordSecret := buildPasswordSecret(testNamespace.Name, "third-pass-secret", "mypass")
				Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())

				testDBUser3 = akov2.NewDBUser(testNamespace.Name, dbUserName3, dbUserName3, projectName).
					WithPasswordSecret("third-pass-secret").
					WithRole("readWriteAnyDatabase", "admin", "")
				Expect(k8sClient.Create(context.Background(), testDBUser3)).To(Succeed())

				Eventually(func(g Gomega) bool {
					expectedConditions := conditions.MatchConditions(
						api.TrueCondition(api.ReadyType),
						api.TrueCondition(api.ResourceVersionStatus),
						api.TrueCondition(api.ValidationSucceeded),
						api.TrueCondition(api.DatabaseUserReadyType),
					)

					g.Expect(k8sClient.Get(context.Background(), client.ObjectKeyFromObject(testDBUser3), testDBUser3, &client.GetOptions{}))
					g.Expect(testDBUser3.Status.Conditions).To(ContainElements(expectedConditions))

					return true
				}).WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})

			By("Deleting AtlasDatabaseUser custom resource", func() {
				By("Keeping database user 1 in Atlas", func() {
					Expect(k8sClient.Delete(context.Background(), testDBUser1)).To(Succeed())

					secretName := fmt.Sprintf(
						"%s-%s-%s",
						kube.NormalizeIdentifier(projectName),
						kube.NormalizeIdentifier(testDeployment.GetDeploymentName()),
						kube.NormalizeIdentifier(dbUserName1),
					)
					Eventually(checkSecretsDontExist(testProject.ID(), []string{secretName})).
						WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())

					Eventually(checkAtlasDatabaseUserRemoved(testProject.ID(), *testDBUser1)).
						WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeFalse())
				})

				By("Deleting database user 2 in Atlas", func() {
					Expect(k8sClient.Delete(context.Background(), testDBUser2)).To(Succeed())

					secretName := fmt.Sprintf(
						"%s-%s-%s",
						kube.NormalizeIdentifier(projectName),
						kube.NormalizeIdentifier(testDeployment.GetDeploymentName()),
						kube.NormalizeIdentifier(dbUserName2),
					)
					Eventually(checkSecretsDontExist(testProject.ID(), []string{secretName})).
						WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())

					Eventually(checkAtlasDatabaseUserRemoved(testProject.ID(), *testDBUser2)).
						WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
				})

				By("Keeping database user 3 in Atlas", func() {
					Expect(k8sClient.Delete(context.Background(), testDBUser3)).To(Succeed())

					Eventually(checkAtlasDatabaseUserRemoved(testProject.ID(), *testDBUser3)).
						WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeFalse())
				})
			})
		})

		It("Adds database users and manage them when operator take ownership of existing resources", Label("focus-owning-protected"), func() {
			By("First without setting atlas-resource-policy annotation", func() {
				passwordSecret := buildPasswordSecret(testNamespace.Name, UserPasswordSecret, DBUserPassword)
				Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())

				testDBUser1 = akov2.NewDBUser(testNamespace.Name, dbUserName1, dbUserName1, projectName).
					WithPasswordSecret(UserPasswordSecret).
					WithRole("readWriteAnyDatabase", "admin", "")
				Expect(k8sClient.Create(context.Background(), testDBUser1)).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testDBUser1, api.TrueCondition(api.ReadyType))
				}).WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())

				validateSecret(k8sClient, *testProject, *testDeployment, *testDBUser1)

				Expect(tryConnect(testProject.ID(), *testDeployment, *testDBUser1)).Should(Succeed())
			})

			// nolint:dupl
			By("Second setting atlas-resource-policy annotation to delete", func() {
				passwordSecret := buildPasswordSecret(testNamespace.Name, UserPasswordSecret2, DBUserPassword2)
				Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())

				testDBUser2 = akov2.NewDBUser(testNamespace.Name, dbUserName2, dbUserName2, projectName).
					WithPasswordSecret(UserPasswordSecret2).
					WithRole("readWriteAnyDatabase", "admin", "")
				testDBUser2.SetAnnotations(map[string]string{customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyDelete})
				Expect(k8sClient.Create(context.Background(), testDBUser2)).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testDBUser2, api.TrueCondition(api.ReadyType))
				}).WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())

				validateSecret(k8sClient, *testProject, *testDeployment, *testDBUser2)

				Expect(tryConnect(testProject.ID(), *testDeployment, *testDBUser2)).Should(Succeed())
			})

			By("Third previously added in Atlas", func() {
				passwordSecret := buildPasswordSecret(testNamespace.Name, "third-pass-secret", "mypass")
				Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())

				testDBUser3 = akov2.NewDBUser(testNamespace.Name, dbUserName3, dbUserName3, projectName).
					WithPasswordSecret("third-pass-secret").
					WithRole("readAnyDatabase", "admin", "")
				Expect(k8sClient.Create(context.Background(), testDBUser3)).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testDBUser3, api.TrueCondition(api.ReadyType))
				}).WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())

				validateSecret(k8sClient, *testProject, *testDeployment, *testDBUser3)

				Expect(tryConnect(testProject.ID(), *testDeployment, *testDBUser3)).Should(Succeed())
			})

			By("Deleting AtlasDatabaseUser custom resource", func() {
				By("Keeping database user 1 in Atlas", func() {
					Expect(k8sClient.Delete(context.Background(), testDBUser1)).To(Succeed())

					secretName := fmt.Sprintf(
						"%s-%s-%s",
						kube.NormalizeIdentifier(projectName),
						kube.NormalizeIdentifier(testDeployment.GetDeploymentName()),
						kube.NormalizeIdentifier(dbUserName1),
					)
					Eventually(checkSecretsDontExist(testProject.ID(), []string{secretName})).
						WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())

					Eventually(checkAtlasDatabaseUserRemoved(testProject.ID(), *testDBUser1)).
						WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeFalse())
				})

				By("Deleting database user 2 in Atlas", func() {
					Expect(k8sClient.Delete(context.Background(), testDBUser2)).To(Succeed())

					secretName := fmt.Sprintf(
						"%s-%s-%s",
						kube.NormalizeIdentifier(projectName),
						kube.NormalizeIdentifier(testDeployment.GetDeploymentName()),
						kube.NormalizeIdentifier(dbUserName2),
					)
					Eventually(checkSecretsDontExist(testProject.ID(), []string{secretName})).
						WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())

					Eventually(checkAtlasDatabaseUserRemoved(testProject.ID(), *testDBUser2)).
						WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
				})

				By("Keeping database user 3 in Atlas", func() {
					Expect(k8sClient.Delete(context.Background(), testDBUser3)).To(Succeed())

					secretName := fmt.Sprintf(
						"%s-%s-%s",
						kube.NormalizeIdentifier(projectName),
						kube.NormalizeIdentifier(testDeployment.GetDeploymentName()),
						kube.NormalizeIdentifier(dbUserName3),
					)
					Eventually(checkSecretsDontExist(testProject.ID(), []string{secretName})).
						WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())

					Eventually(checkAtlasDatabaseUserRemoved(testProject.ID(), *testDBUser3)).
						WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeFalse())
				})
			})
		})
	})

	// nolint:dupl
	AfterEach(func() {
		By("Deleting deployment", func() {
			deploymentName := testDeployment.GetDeploymentName()
			Expect(k8sClient.Delete(context.Background(), testDeployment)).To(Succeed())

			Eventually(func() bool {
				_, r, err := atlasClient.FlexClustersApi.
					GetFlexCluster(context.Background(), testProject.ID(), deploymentName).
					Execute()
				if err != nil {
					if httputil.StatusCode(r) == http.StatusNotFound {
						return true
					}
				}

				return false
			}).WithTimeout(20 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
		})

		By("Deleting project", func() {
			projectID := testProject.ID()
			Expect(k8sClient.Delete(context.Background(), testProject)).To(Succeed())

			_, err := atlasClient.ProjectsApi.DeleteGroup(context.Background(), projectID).Execute()
			Expect(err).To(BeNil())

			Eventually(func() bool {
				_, r, err := atlasClient.ProjectsApi.GetGroup(context.Background(), projectID).Execute()
				if err != nil {
					if httputil.StatusCode(r) == http.StatusNotFound {
						return true
					}
				}

				return false
			}).WithTimeout(15 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
		})

		By("Stopping the operator", func() {
			stopManager()

			By("Removing the namespace " + testNamespace.Name)
			err := k8sClient.Delete(context.Background(), testNamespace)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
