package int

import (
	"context"
	"fmt"
	"net/http"
	"time"

	corev1 "k8s.io/api/core/v1"

	"go.mongodb.org/atlas/mongodbatlas"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/testutil"
)

var _ = Describe("Atlas Database User", Label("int", "AtlasDatabaseUser", "protection-enabled"), func() {
	var testNamespace *corev1.Namespace
	var stopManager context.CancelFunc
	var projectName string
	projectNamePrefix := "database-user-protected"
	dbUserName1 := "db-user1"
	dbUserName2 := "db-user2"
	dbUserName3 := "db-user3"
	testProject := &mdbv1.AtlasProject{}
	testDeployment := &mdbv1.AtlasDeployment{}
	testDBUser1 := &mdbv1.AtlasDatabaseUser{}
	testDBUser2 := &mdbv1.AtlasDatabaseUser{}
	testDBUser3 := &mdbv1.AtlasDatabaseUser{}

	BeforeEach(func() {
		testNamespace, stopManager = prepareControllers(true)
		projectName = fmt.Sprintf("%s-%s", projectNamePrefix, testNamespace.Name)

		By("Creating a project", func() {
			connSecret := buildConnectionSecret("my-atlas-key")
			Expect(k8sClient.Create(context.TODO(), &connSecret)).To(Succeed())

			testProject = mdbv1.NewProject(testNamespace.Name, projectName, projectName).
				WithConnectionSecret(connSecret.Name).
				WithIPAccessList(project.NewIPAccessList().WithCIDR("0.0.0.0/0"))
			Expect(k8sClient.Create(context.TODO(), testProject, &client.CreateOptions{})).To(Succeed())

			Eventually(func() bool {
				return testutil.CheckCondition(k8sClient, testProject, status.TrueCondition(status.ReadyType))
			}).WithTimeout(15 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
		})

		By("Creating a deployment", func() {
			testDeployment = mdbv1.DefaultAWSDeployment(testNamespace.Name, projectName).Lightweight()
			customresource.SetAnnotation( // this test deployment must be deleted
				testDeployment,
				customresource.ResourcePolicyAnnotation,
				customresource.ResourcePolicyDelete,
			)
			Expect(k8sClient.Create(context.TODO(), testDeployment)).To(Succeed())

			Eventually(func() bool {
				return testutil.CheckCondition(k8sClient, testDeployment, status.TrueCondition(status.ReadyType))
			}).WithTimeout(20 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
		})

		By("Creating database user", func() {
			dbUser := &mongodbatlas.DatabaseUser{
				Username:     dbUserName3,
				Password:     "mypass",
				DatabaseName: "admin",
				Roles: []mongodbatlas.Role{
					{
						RoleName:     "readAnyDatabase",
						DatabaseName: "admin",
					},
				},
				Scopes: []mongodbatlas.Scope{},
			}
			_, _, err := atlasClient.DatabaseUsers.Create(context.TODO(), testProject.ID(), dbUser)
			Expect(err).To(BeNil())
		})
	})

	Describe("Operator is running with deletion protection enabled", func() {
		It("Adds database users and protect them to be deleted when operator doesn't own resource", func() {
			By("First without setting atlas-resource-policy annotation", func() {
				passwordSecret := buildPasswordSecret(testNamespace.Name, UserPasswordSecret, DBUserPassword)
				Expect(k8sClient.Create(context.TODO(), &passwordSecret)).To(Succeed())

				testDBUser1 = mdbv1.NewDBUser(testNamespace.Name, dbUserName1, dbUserName1, projectName).
					WithPasswordSecret(UserPasswordSecret).
					WithRole("readWriteAnyDatabase", "admin", "")
				Expect(k8sClient.Create(context.TODO(), testDBUser1)).To(Succeed())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, testDBUser1, status.TrueCondition(status.ReadyType))
				}).WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())

				validateSecret(k8sClient, *testProject, *testDeployment, *testDBUser1)

				Expect(tryConnect(testProject.ID(), *testDeployment, *testDBUser1)).Should(Succeed())
			})

			// nolint:dupl
			By("Second setting atlas-resource-policy annotation to delete", func() {
				passwordSecret := buildPasswordSecret(testNamespace.Name, UserPasswordSecret2, DBUserPassword2)
				Expect(k8sClient.Create(context.TODO(), &passwordSecret)).To(Succeed())

				testDBUser2 = mdbv1.NewDBUser(testNamespace.Name, dbUserName2, dbUserName2, projectName).
					WithPasswordSecret(UserPasswordSecret2).
					WithRole("readWriteAnyDatabase", "admin", "")
				testDBUser2.SetAnnotations(map[string]string{customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyDelete})
				Expect(k8sClient.Create(context.TODO(), testDBUser2)).To(Succeed())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, testDBUser2, status.TrueCondition(status.ReadyType))
				}).WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())

				validateSecret(k8sClient, *testProject, *testDeployment, *testDBUser2)

				Expect(tryConnect(testProject.ID(), *testDeployment, *testDBUser2)).Should(Succeed())
			})

			By("Third previously added in Atlas", func() {
				passwordSecret := buildPasswordSecret(testNamespace.Name, "third-pass-secret", "mypass")
				Expect(k8sClient.Create(context.TODO(), &passwordSecret)).To(Succeed())

				testDBUser3 = mdbv1.NewDBUser(testNamespace.Name, dbUserName3, dbUserName3, projectName).
					WithPasswordSecret("third-pass-secret").
					WithRole("readWriteAnyDatabase", "admin", "")
				Expect(k8sClient.Create(context.TODO(), testDBUser3)).To(Succeed())

				Eventually(func(g Gomega) bool {
					expectedConditions := testutil.MatchConditions(
						status.TrueCondition(status.ValidationSucceeded),
						status.FalseCondition(status.ReadyType),
						status.FalseCondition(status.DatabaseUserReadyType).
							WithReason(string(workflow.AtlasDeletionProtection)).
							WithMessageRegexp("unable to reconcile database user: it already exists in Atlas, it was not previously managed by the operator, and the deletion protection is enabled."),
					)

					g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(testDBUser3), testDBUser3, &client.GetOptions{}))
					g.Expect(testDBUser3.Status.Conditions).To(ContainElements(expectedConditions))

					return true
				}).WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})

			By("Deleting AtlasDatabaseUser custom resource", func() {
				By("Keeping database user 1 in Atlas", func() {
					Expect(k8sClient.Delete(context.TODO(), testDBUser1)).To(Succeed())

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
					Expect(k8sClient.Delete(context.TODO(), testDBUser2)).To(Succeed())

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
					Expect(k8sClient.Delete(context.TODO(), testDBUser3)).To(Succeed())

					Eventually(checkAtlasDatabaseUserRemoved(testProject.ID(), *testDBUser3)).
						WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeFalse())
				})
			})
		})

		It("Adds database users and manage them when operator take ownership of existing resources", func() {
			By("First without setting atlas-resource-policy annotation", func() {
				passwordSecret := buildPasswordSecret(testNamespace.Name, UserPasswordSecret, DBUserPassword)
				Expect(k8sClient.Create(context.TODO(), &passwordSecret)).To(Succeed())

				testDBUser1 = mdbv1.NewDBUser(testNamespace.Name, dbUserName1, dbUserName1, projectName).
					WithPasswordSecret(UserPasswordSecret).
					WithRole("readWriteAnyDatabase", "admin", "")
				Expect(k8sClient.Create(context.TODO(), testDBUser1)).To(Succeed())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, testDBUser1, status.TrueCondition(status.ReadyType))
				}).WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())

				validateSecret(k8sClient, *testProject, *testDeployment, *testDBUser1)

				Expect(tryConnect(testProject.ID(), *testDeployment, *testDBUser1)).Should(Succeed())
			})

			// nolint:dupl
			By("Second setting atlas-resource-policy annotation to delete", func() {
				passwordSecret := buildPasswordSecret(testNamespace.Name, UserPasswordSecret2, DBUserPassword2)
				Expect(k8sClient.Create(context.TODO(), &passwordSecret)).To(Succeed())

				testDBUser2 = mdbv1.NewDBUser(testNamespace.Name, dbUserName2, dbUserName2, projectName).
					WithPasswordSecret(UserPasswordSecret2).
					WithRole("readWriteAnyDatabase", "admin", "")
				testDBUser2.SetAnnotations(map[string]string{customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyDelete})
				Expect(k8sClient.Create(context.TODO(), testDBUser2)).To(Succeed())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, testDBUser2, status.TrueCondition(status.ReadyType))
				}).WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())

				validateSecret(k8sClient, *testProject, *testDeployment, *testDBUser2)

				Expect(tryConnect(testProject.ID(), *testDeployment, *testDBUser2)).Should(Succeed())
			})

			By("Third previously added in Atlas", func() {
				passwordSecret := buildPasswordSecret(testNamespace.Name, "third-pass-secret", "mypass")
				Expect(k8sClient.Create(context.TODO(), &passwordSecret)).To(Succeed())

				testDBUser3 = mdbv1.NewDBUser(testNamespace.Name, dbUserName3, dbUserName3, projectName).
					WithPasswordSecret("third-pass-secret").
					WithRole("readAnyDatabase", "admin", "")
				Expect(k8sClient.Create(context.TODO(), testDBUser3)).To(Succeed())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, testDBUser3, status.TrueCondition(status.ReadyType))
				}).WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())

				validateSecret(k8sClient, *testProject, *testDeployment, *testDBUser3)

				Expect(tryConnect(testProject.ID(), *testDeployment, *testDBUser3)).Should(Succeed())
			})

			By("Deleting AtlasDatabaseUser custom resource", func() {
				By("Keeping database user 1 in Atlas", func() {
					Expect(k8sClient.Delete(context.TODO(), testDBUser1)).To(Succeed())

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
					Expect(k8sClient.Delete(context.TODO(), testDBUser2)).To(Succeed())

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
					Expect(k8sClient.Delete(context.TODO(), testDBUser3)).To(Succeed())

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
			Expect(k8sClient.Delete(context.TODO(), testDeployment)).To(Succeed())

			Eventually(func() bool {
				_, r, err := atlasClient.AdvancedClusters.Get(context.TODO(), testProject.ID(), deploymentName)
				if err != nil {
					if r != nil && r.StatusCode == http.StatusNotFound {
						return true
					}
				}

				return false
			}).WithTimeout(20 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
		})

		By("Deleting project", func() {
			projectID := testProject.ID()
			Expect(k8sClient.Delete(context.TODO(), testProject)).To(Succeed())

			Eventually(func() bool {
				_, r, err := atlasClient.Projects.GetOneProject(context.TODO(), projectID)
				if err != nil {
					if r != nil && r.StatusCode == http.StatusNotFound {
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
