package int

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/conditions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/events"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/resources"
)

const (
	databaseUserTimeout = 10 * time.Minute
	UserPasswordSecret  = "user-password-secret"
	DBUserPassword      = "Passw0rd!"
	UserPasswordSecret2 = "second-user-password-secret"
	DBUserPassword2     = "H@lla#!"
)

var _ = Describe("Atlas Database User", Label("int", "AtlasDatabaseUser", "protection-disabled"), func() {
	var testNamespace *corev1.Namespace
	var stopManager context.CancelFunc
	var projectName string
	projectNamePrefix := "database-user-unprotected"
	dbUserName1 := "db-user1"
	dbUserName2 := "db-user2"
	dbUserName3 := "db-user3"
	testProject := &akov2.AtlasProject{}
	testDeployment := &akov2.AtlasDeployment{}
	testDBUser1 := &akov2.AtlasDatabaseUser{}
	testDBUser2 := &akov2.AtlasDatabaseUser{}
	testDBUser3 := &akov2.AtlasDatabaseUser{}

	BeforeEach(func() {
		testNamespace, stopManager = prepareControllers(false)
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
			testDeployment = akov2.DefaultAWSDeployment(testNamespace.Name, projectName).Lightweight()
			Expect(k8sClient.Create(context.Background(), testDeployment)).To(Succeed())

			Eventually(func() bool {
				return resources.CheckCondition(k8sClient, testDeployment, api.TrueCondition(api.ReadyType))
			}).WithTimeout(20 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
		})
	})

	Describe("Operator is running with deletion protection disabled", func() {
		It("Adds database users and allow them to be deleted", func() {
			By("Creating a database user previously on Atlas", func() {
				dbUser := admin.NewCloudDatabaseUser("admin", testProject.ID(), dbUserName3)
				dbUser.SetPassword("mypass")
				dbUser.SetRoles(
					[]admin.DatabaseUserRole{
						{
							RoleName:     "readWriteAnyDatabase",
							DatabaseName: "admin",
						},
					},
				)
				_, _, err := atlasClient.DatabaseUsersApi.CreateDatabaseUser(context.Background(), testProject.ID(), dbUser).Execute()
				Expect(err).To(BeNil())
			})

			By("First without setting atlas-resource-policy annotation", func() {
				passwordSecret := buildPasswordSecret(testNamespace.Name, UserPasswordSecret, DBUserPassword)
				Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())

				testDBUser1 = akov2.NewDBUser(testNamespace.Name, dbUserName1, dbUserName1, projectName).
					WithPasswordSecret(UserPasswordSecret).
					WithRole("readWriteAnyDatabase", "admin", "")
				Expect(k8sClient.Create(context.Background(), testDBUser1)).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testDBUser1, api.TrueCondition(api.ReadyType))
				}).WithTimeout(databaseUserTimeout).WithPolling(PollingInterval).Should(BeTrue())

				validateSecret(k8sClient, *testProject, *testDeployment, *testDBUser1)

				Expect(tryConnect(testProject.ID(), *testDeployment, *testDBUser1)).Should(Succeed())
			})

			By("Second setting atlas-resource-policy annotation to keep", func() {
				passwordSecret := buildPasswordSecret(testNamespace.Name, UserPasswordSecret2, DBUserPassword2)
				Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())

				testDBUser2 = akov2.NewDBUser(testNamespace.Name, dbUserName2, dbUserName2, projectName).
					WithPasswordSecret(UserPasswordSecret2).
					WithRole("readWriteAnyDatabase", "admin", "")
				testDBUser2.SetAnnotations(map[string]string{customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep})
				Expect(k8sClient.Create(context.Background(), testDBUser2)).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testDBUser2, api.TrueCondition(api.ReadyType))
				}).WithTimeout(databaseUserTimeout).WithPolling(PollingInterval).Should(BeTrue())

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

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testDBUser3, api.TrueCondition(api.ReadyType))
				}).WithTimeout(databaseUserTimeout).WithPolling(PollingInterval).Should(BeTrue())

				validateSecret(k8sClient, *testProject, *testDeployment, *testDBUser3)

				Expect(tryConnect(testProject.ID(), *testDeployment, *testDBUser3)).Should(Succeed())
			})

			By("Deleting AtlasDatabaseUser custom resource", func() {
				By("Deleting database user 1 in Atlas", func() {
					deleteSecret(testDBUser1)
					Expect(k8sClient.Delete(context.Background(), testDBUser1)).To(Succeed())

					secretName := fmt.Sprintf(
						"%s-%s-%s",
						kube.NormalizeIdentifier(projectName),
						kube.NormalizeIdentifier(testDeployment.GetDeploymentName()),
						kube.NormalizeIdentifier(dbUserName1),
					)
					Eventually(checkSecretsDontExist(testProject.ID(), []string{secretName})).
						WithTimeout(databaseUserTimeout).WithPolling(PollingInterval).Should(BeTrue())

					Eventually(checkAtlasDatabaseUserRemoved(testProject.ID(), *testDBUser1)).
						WithTimeout(databaseUserTimeout).WithPolling(PollingInterval).Should(BeTrue())
				})

				By("Keeping database user 2 in Atlas", func() {
					deleteSecret(testDBUser2)
					Expect(k8sClient.Delete(context.Background(), testDBUser2)).To(Succeed())

					secretName := fmt.Sprintf(
						"%s-%s-%s",
						kube.NormalizeIdentifier(projectName),
						kube.NormalizeIdentifier(testDeployment.GetDeploymentName()),
						kube.NormalizeIdentifier(dbUserName2),
					)
					Eventually(checkSecretsDontExist(testProject.ID(), []string{secretName})).
						WithTimeout(databaseUserTimeout).WithPolling(PollingInterval).Should(BeTrue())

					Eventually(checkAtlasDatabaseUserRemoved(testProject.ID(), *testDBUser2)).
						WithTimeout(databaseUserTimeout).WithPolling(PollingInterval).Should(BeFalse())

					_, _, err := atlasClient.DatabaseUsersApi.
						DeleteDatabaseUser(context.Background(), testProject.ID(), "admin", dbUserName2).
						Execute()
					Expect(err).To(BeNil())
				})

				By("Deleting database user 3 in Atlas", func() {
					deleteSecret(testDBUser3)
					Expect(k8sClient.Delete(context.Background(), testDBUser3)).To(Succeed())

					secretName := fmt.Sprintf(
						"%s-%s-%s",
						kube.NormalizeIdentifier(projectName),
						kube.NormalizeIdentifier(testDeployment.GetDeploymentName()),
						kube.NormalizeIdentifier(dbUserName3),
					)
					Eventually(checkSecretsDontExist(testProject.ID(), []string{secretName})).
						WithTimeout(databaseUserTimeout).WithPolling(PollingInterval).Should(BeTrue())

					Eventually(checkAtlasDatabaseUserRemoved(testProject.ID(), *testDBUser3)).
						WithTimeout(databaseUserTimeout).WithPolling(PollingInterval).Should(BeTrue())
				})
			})
		})

		It("Adds an user and manage roles", func() {
			By("Creating an user with clusterMonitor role", func() {
				passwordSecret := buildPasswordSecret(testNamespace.Name, UserPasswordSecret, DBUserPassword)
				Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())

				testDBUser1 = akov2.NewDBUser(testNamespace.Name, dbUserName1, dbUserName1, projectName).
					WithPasswordSecret(UserPasswordSecret).
					WithRole("clusterMonitor", "admin", "")
				Expect(k8sClient.Create(context.Background(), testDBUser1)).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testDBUser1, api.TrueCondition(api.ReadyType))
				}).WithTimeout(databaseUserTimeout).WithPolling(PollingInterval).Should(BeTrue())
			})

			By("Validating credentials and cluster access", func() {
				validateSecret(k8sClient, *testProject, *testDeployment, *testDBUser1)

				Expect(tryConnect(testProject.ID(), *testDeployment, *testDBUser1)).Should(Succeed())

				err := tryWrite(testProject.ID(), *testDeployment, *testDBUser1, "test", "operatortest")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(MatchRegexp("user is not allowed"))
			})

			By("Giving user readWrite permissions", func() {
				// Adding the role allowing read/write
				testDBUser1 = testDBUser1.WithRole("readWriteAnyDatabase", "admin", "")

				Expect(k8sClient.Update(context.Background(), testDBUser1)).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testDBUser1, api.TrueCondition(api.ReadyType))
				}).WithTimeout(databaseUserTimeout).WithPolling(PollingInterval).Should(BeTrue())
			})

			By("Validating user has permission to write", func() {
				validateSecret(k8sClient, *testProject, *testDeployment, *testDBUser1)

				Expect(tryConnect(testProject.ID(), *testDeployment, *testDBUser1)).Should(Succeed())

				Expect(tryWrite(testProject.ID(), *testDeployment, *testDBUser1, "test", "operatortest")).To(Succeed())
			})

			By("Deleting database user", func() {
				deleteSecret(testDBUser1)
				Expect(k8sClient.Delete(context.Background(), testDBUser1)).To(Succeed())

				secretName := fmt.Sprintf(
					"%s-%s-%s",
					kube.NormalizeIdentifier(projectName),
					kube.NormalizeIdentifier(testDeployment.GetDeploymentName()),
					kube.NormalizeIdentifier(dbUserName1),
				)
				Eventually(checkSecretsDontExist(testProject.ID(), []string{secretName})).
					WithTimeout(databaseUserTimeout).WithPolling(PollingInterval).Should(BeTrue())

				Eventually(checkAtlasDatabaseUserRemoved(testProject.ID(), *testDBUser1)).
					WithTimeout(databaseUserTimeout).WithPolling(PollingInterval).Should(BeTrue())
			})
		})

		It("Adds connection secret when new deployment is created", func() {
			secondDeployment := &akov2.AtlasDeployment{}

			By("Creating a database user", func() {
				passwordSecret := buildPasswordSecret(testNamespace.Name, UserPasswordSecret, DBUserPassword)
				Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())

				testDBUser1 = akov2.NewDBUser(testNamespace.Name, dbUserName1, dbUserName1, projectName).
					WithPasswordSecret(UserPasswordSecret).
					WithRole("readWriteAnyDatabase", "admin", "")
				Expect(k8sClient.Create(context.Background(), testDBUser1)).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testDBUser1, api.TrueCondition(api.ReadyType))
				}).WithTimeout(databaseUserTimeout).WithPolling(PollingInterval).Should(BeTrue())

				validateSecret(k8sClient, *testProject, *testDeployment, *testDBUser1)

				Expect(tryConnect(testProject.ID(), *testDeployment, *testDBUser1)).Should(Succeed())
			})

			By("Creating a  second deployment", func() {
				secondDeployment = akov2.DefaultAzureDeployment(testNamespace.Name, projectName).Lightweight()
				Expect(k8sClient.Create(context.Background(), secondDeployment)).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, secondDeployment, api.TrueCondition(api.ReadyType))
				}).WithTimeout(20 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})

			By("Validating connection secrets were created", func() {
				validateSecret(k8sClient, *testProject, *testDeployment, *testDBUser1)
				validateSecret(k8sClient, *testProject, *secondDeployment, *testDBUser1)

				Expect(tryConnect(testProject.ID(), *testDeployment, *testDBUser1)).Should(Succeed())
				Expect(tryConnect(testProject.ID(), *secondDeployment, *testDBUser1)).Should(Succeed())
			})

			By("Deleting the second deployment", func() {
				deploymentName := secondDeployment.GetDeploymentName()
				Expect(k8sClient.Delete(context.Background(), secondDeployment)).To(Succeed())

				Eventually(func() bool {
					_, r, err := atlasClient.ClustersApi.
						GetCluster(context.Background(), testProject.ID(), deploymentName).
						Execute()
					if err != nil {
						if r != nil && r.StatusCode == http.StatusNotFound {
							return true
						}
					}

					return false
				}).WithTimeout(20 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})
		})

		It("Watches password secret", func() {
			By("Creating a database user", func() {
				passwordSecret := buildPasswordSecret(testNamespace.Name, UserPasswordSecret, DBUserPassword)
				Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())

				testDBUser1 = akov2.NewDBUser(testNamespace.Name, dbUserName1, dbUserName1, projectName).
					WithPasswordSecret(UserPasswordSecret).
					WithRole("readWriteAnyDatabase", "admin", "")
				Expect(k8sClient.Create(context.Background(), testDBUser1)).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testDBUser1, api.TrueCondition(api.ReadyType))
				}).WithTimeout(databaseUserTimeout).WithPolling(PollingInterval).Should(BeTrue())

				validateSecret(k8sClient, *testProject, *testDeployment, *testDBUser1)

				Expect(tryConnect(testProject.ID(), *testDeployment, *testDBUser1)).Should(Succeed())
			})

			By("Breaking the password secret", func() {
				passwordSecret := buildPasswordSecret(testNamespace.Name, UserPasswordSecret, "")
				Expect(k8sClient.Update(context.Background(), &passwordSecret)).To(Succeed())

				expectedCondition := api.FalseCondition(api.DatabaseUserReadyType).WithReason(string(workflow.Internal)).WithMessageRegexp("the 'password' field is empty")
				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testDBUser1, expectedCondition)
				}).WithTimeout(databaseUserTimeout).WithPolling(PollingInterval).Should(BeTrue())

				events.EventExists(k8sClient, testDBUser1, "Warning", string(workflow.Internal), "the 'password' field is empty")
			})

			By("Fixing the password secret", func() {
				passwordSecret := buildPasswordSecret(testNamespace.Name, UserPasswordSecret, "someNewPassw00rd")
				Expect(k8sClient.Update(context.Background(), &passwordSecret)).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testDBUser1, api.TrueCondition(api.ReadyType))
				}).WithTimeout(databaseUserTimeout).WithPolling(PollingInterval).Should(BeTrue())

				// We need to make sure that the new connection secret is different from the initial one
				connSecretUpdated := validateSecret(k8sClient, *testProject, *testDeployment, *testDBUser1)
				Expect(string(connSecretUpdated.Data["password"])).To(Equal("someNewPassw00rd"))

				var updatedPwdSecret corev1.Secret
				Expect(k8sClient.Get(context.Background(), kube.ObjectKey(testNamespace.Name, UserPasswordSecret), &updatedPwdSecret)).To(Succeed())
				Expect(testDBUser1.Status.PasswordVersion).To(Equal(updatedPwdSecret.ResourceVersion))

				Expect(tryConnect(testProject.ID(), *testDeployment, *testDBUser1)).Should(Succeed())
			})
		})

		It("Remove stale secrets", func() {
			secondTestDeployment := &akov2.AtlasDeployment{}

			By("Creating a second deployment", func() {
				secondTestDeployment = akov2.DefaultAzureDeployment(testNamespace.Name, projectName).Lightweight()
				Expect(k8sClient.Create(context.Background(), secondTestDeployment)).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, secondTestDeployment, api.TrueCondition(api.ReadyType))
				}).WithTimeout(20 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})

			By("Creating a database user", func() {
				passwordSecret := buildPasswordSecret(testNamespace.Name, UserPasswordSecret, DBUserPassword)
				Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())

				testDBUser1 = akov2.NewDBUser(testNamespace.Name, dbUserName1, dbUserName1, projectName).
					WithPasswordSecret(UserPasswordSecret).
					WithRole("readWriteAnyDatabase", "admin", "")
				Expect(k8sClient.Create(context.Background(), testDBUser1)).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testDBUser1, api.TrueCondition(api.ReadyType))
				}).WithTimeout(databaseUserTimeout).WithPolling(PollingInterval).Should(BeTrue())

				validateSecret(k8sClient, *testProject, *testDeployment, *testDBUser1)
				validateSecret(k8sClient, *testProject, *secondTestDeployment, *testDBUser1)

				Expect(tryConnect(testProject.ID(), *testDeployment, *testDBUser1)).Should(Succeed())
				Expect(tryConnect(testProject.ID(), *secondTestDeployment, *testDBUser1)).Should(Succeed())
			})

			By("Renaming username, new user is added and stale secrets are removed", func() {
				oldName := testDBUser1.Spec.Username
				testDBUser1 = testDBUser1.WithAtlasUserName("new-user")
				Expect(k8sClient.Update(context.Background(), testDBUser1)).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testDBUser1, api.TrueCondition(api.ReadyType))
				}).WithTimeout(databaseUserTimeout).WithPolling(PollingInterval).Should(BeTrue())

				_, _, err := atlasClient.DatabaseUsersApi.
					GetDatabaseUser(context.Background(), testProject.ID(), testDBUser1.Spec.DatabaseName, oldName).
					Execute()
				Expect(err).To(HaveOccurred())

				checkNumberOfConnectionSecrets(k8sClient, *testProject, testNamespace.Name, 2)
				secret := validateSecret(k8sClient, *testProject, *testDeployment, *testDBUser1)
				Expect(secret.Name).To(Equal(fmt.Sprintf("%s-test-deployment-aws-new-user", kube.NormalizeIdentifier(testProject.Spec.Name))))
				secret = validateSecret(k8sClient, *testProject, *secondTestDeployment, *testDBUser1)
				Expect(secret.Name).To(Equal(fmt.Sprintf("%s-test-deployment-azure-new-user", kube.NormalizeIdentifier(testProject.Spec.Name))))

				Expect(tryConnect(testProject.ID(), *testDeployment, *testDBUser1)).Should(Succeed())
				Expect(tryConnect(testProject.ID(), *secondTestDeployment, *testDBUser1)).Should(Succeed())
			})

			By("Scoping user to one cluster, a stale secret is removed", func() {
				testDBUser1 = testDBUser1.ClearScopes().WithScope(akov2.DeploymentScopeType, testDeployment.GetDeploymentName())
				Expect(k8sClient.Update(context.Background(), testDBUser1)).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testDBUser1, api.TrueCondition(api.ReadyType))
				}).WithTimeout(databaseUserTimeout).WithPolling(PollingInterval).Should(BeTrue())

				checkNumberOfConnectionSecrets(k8sClient, *testProject, testNamespace.Name, 1)
				validateSecret(k8sClient, *testProject, *testDeployment, *testDBUser1)

				Expect(tryConnect(testProject.ID(), *testDeployment, *testDBUser1)).Should(Succeed())
				Expect(tryConnect(testProject.ID(), *secondTestDeployment, *testDBUser1)).ShouldNot(Succeed())
			})

			By("Deleting second deployment", func() {
				deploymentName := secondTestDeployment.GetDeploymentName()
				Expect(k8sClient.Delete(context.Background(), secondTestDeployment)).To(Succeed())

				Eventually(func() bool {
					_, r, err := atlasClient.ClustersApi.
						GetCluster(context.Background(), testProject.ID(), deploymentName).
						Execute()
					if err != nil {
						if r != nil && r.StatusCode == http.StatusNotFound {
							return true
						}
					}

					return false
				}).WithTimeout(20 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})
		})

		It("Validates user date expiration", func() {
			By("Creating expired user", func() {
				passwordSecret := buildPasswordSecret(testNamespace.Name, UserPasswordSecret, DBUserPassword)
				Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())

				before := time.Now().UTC().Add(time.Minute * -10).Format("2006-01-02T15:04:05")

				testDBUser1 = akov2.NewDBUser(testNamespace.Name, dbUserName1, dbUserName1, projectName).
					WithPasswordSecret(UserPasswordSecret).
					WithRole("readWriteAnyDatabase", "admin", "").
					WithDeleteAfterDate(before)
				Expect(k8sClient.Create(context.Background(), testDBUser1)).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testDBUser1, api.FalseCondition(api.DatabaseUserReadyType).WithReason(string(workflow.DatabaseUserExpired)))
				}).WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())

				checkNumberOfConnectionSecrets(k8sClient, *testProject, testNamespace.Name, 0)

				_, _, err := atlasClient.DatabaseUsersApi.
					GetDatabaseUser(context.Background(), testProject.ID(), testDBUser1.Spec.DatabaseName, testDBUser1.Spec.Username).
					Execute()
				Expect(err).To(HaveOccurred())
			})

			By("Fixing the user date expiration", func() {
				after := time.Now().UTC().Add(time.Hour * 10).Format("2006-01-02T15:04:05")
				testDBUser1 = testDBUser1.WithDeleteAfterDate(after)

				Expect(k8sClient.Update(context.Background(), testDBUser1)).To(Succeed())
				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testDBUser1, api.TrueCondition(api.ReadyType))
				}).WithTimeout(databaseUserTimeout).WithPolling(PollingInterval).Should(BeTrue())

				checkNumberOfConnectionSecrets(k8sClient, *testProject, testNamespace.Name, 1)
				Expect(tryConnect(testProject.ID(), *testDeployment, *testDBUser1)).Should(Succeed())
			})

			By("Expiring the User", func() {
				before := time.Now().UTC().Add(time.Minute * -5).Format("2006-01-02T15:04:05")
				testDBUser1 = testDBUser1.WithDeleteAfterDate(before)

				Expect(k8sClient.Update(context.Background(), testDBUser1)).To(Succeed())
				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testDBUser1, api.FalseCondition(api.DatabaseUserReadyType).WithReason(string(workflow.DatabaseUserExpired)))
				}).WithTimeout(databaseUserTimeout).WithPolling(PollingInterval).Should(BeTrue())

				expectedConditionsMatchers := conditions.MatchConditions(
					api.FalseCondition(api.DatabaseUserReadyType),
					api.FalseCondition(api.ReadyType),
					api.TrueCondition(api.ValidationSucceeded),
					api.TrueCondition(api.ResourceVersionStatus),
				)
				Expect(testDBUser1.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))

				checkNumberOfConnectionSecrets(k8sClient, *testProject, testNamespace.Name, 0)
			})
		})

		It("Skips reconciliations.", func() {
			By("Creating a database user", func() {
				passwordSecret := buildPasswordSecret(testNamespace.Name, UserPasswordSecret, DBUserPassword)
				Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())

				testDBUser1 = akov2.NewDBUser(testNamespace.Name, dbUserName1, dbUserName1, projectName).
					WithPasswordSecret(UserPasswordSecret).
					WithRole("readWriteAnyDatabase", "admin", "")
				Expect(k8sClient.Create(context.Background(), testDBUser1)).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testDBUser1, api.TrueCondition(api.ReadyType))
				}).WithTimeout(databaseUserTimeout).WithPolling(PollingInterval).Should(BeTrue())

				validateSecret(k8sClient, *testProject, *testDeployment, *testDBUser1)

				Expect(tryConnect(testProject.ID(), *testDeployment, *testDBUser1)).Should(Succeed())
			})

			By("Skipping reconciliation", func() {
				testDBUser1.ObjectMeta.Annotations = map[string]string{customresource.ReconciliationPolicyAnnotation: customresource.ReconciliationPolicySkip}
				testDBUser1.Spec.Roles = append(testDBUser1.Spec.Roles, akov2.RoleSpec{
					RoleName:       "new-role",
					DatabaseName:   "new-database",
					CollectionName: "new-collection",
				})

				Expect(k8sClient.Update(context.Background(), testDBUser1)).To(Succeed())

				ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
				defer cancel()
				containsDatabaseUser := func(dbUser *admin.CloudDatabaseUser) bool {
					for _, role := range dbUser.GetRoles() {
						if role.RoleName == "new-role" && role.DatabaseName == "new-database" && role.GetCollectionName() == "new-collection" {
							return true
						}
					}
					return false
				}

				Eventually(atlas.WaitForAtlasDatabaseUserStateToNotBeReached(ctx, atlasClient, "admin", testProject.Name, testDeployment.GetDeploymentName(), containsDatabaseUser))
			})
		})
	})

	// nolint:dupl
	AfterEach(func() {
		By("Deleting the deployment", func() {
			deploymentName := testDeployment.GetDeploymentName()
			Expect(k8sClient.Delete(context.Background(), testDeployment)).To(Succeed())

			Eventually(func() bool {
				_, r, err := atlasClient.ClustersApi.
					GetCluster(context.Background(), testProject.ID(), deploymentName).
					Execute()
				if err != nil {
					if r != nil && r.StatusCode == http.StatusNotFound {
						return true
					}
				}

				return false
			}).WithTimeout(20 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
		})

		By("Deleting the project", func() {
			projectID := testProject.ID()
			Expect(k8sClient.Delete(context.Background(), testProject)).To(Succeed())

			Eventually(func() bool {
				_, r, err := atlasClient.ProjectsApi.GetProject(context.Background(), projectID).Execute()
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

func buildPasswordSecret(namespace, name, password string) corev1.Secret {
	return corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				connectionsecret.TypeLabelKey: connectionsecret.CredLabelVal,
			},
		},
		StringData: map[string]string{"password": password},
	}
}

func validateSecret(k8sClient client.Client, project akov2.AtlasProject, deployment akov2.AtlasDeployment, user akov2.AtlasDatabaseUser) corev1.Secret {
	secret := corev1.Secret{}
	username := user.Spec.Username
	secretName := fmt.Sprintf("%s-%s-%s", kube.NormalizeIdentifier(project.Spec.Name), kube.NormalizeIdentifier(deployment.GetDeploymentName()), kube.NormalizeIdentifier(username))
	Expect(k8sClient.Get(context.Background(), kube.ObjectKey(project.Namespace, secretName), &secret)).To(Succeed())

	password, err := user.ReadPassword(context.Background(), k8sClient)
	Expect(err).NotTo(HaveOccurred())

	c, _, err := atlasClient.ClustersApi.
		GetCluster(context.Background(), project.ID(), deployment.GetDeploymentName()).
		Execute()
	Expect(err).NotTo(HaveOccurred())

	connectionStrings := c.GetConnectionStrings()

	expectedData := map[string][]byte{
		"connectionStringStandard":    []byte(buildConnectionURL(connectionStrings.GetStandard(), username, password)),
		"connectionStringStandardSrv": []byte(buildConnectionURL(connectionStrings.GetStandardSrv(), username, password)),
		"connectionStringPrivate":     []byte(buildConnectionURL(connectionStrings.GetPrivate(), username, password)),
		"connectionStringPrivateSrv":  []byte(buildConnectionURL(connectionStrings.GetPrivateSrv(), username, password)),
		"username":                    []byte(username),
		"password":                    []byte(password),
	}
	expectedLabels := map[string]string{
		"atlas.mongodb.com/project-id":   project.ID(),
		"atlas.mongodb.com/cluster-name": deployment.GetDeploymentName(),
		connectionsecret.TypeLabelKey:    connectionsecret.CredLabelVal,
	}
	Expect(secret.Data).To(Equal(expectedData))
	Expect(secret.Labels).To(Equal(expectedLabels))

	return secret
}

func checkNumberOfConnectionSecrets(k8sClient client.Client, project akov2.AtlasProject, namespace string, length int) {
	secretList := corev1.SecretList{}
	Expect(k8sClient.List(context.Background(), &secretList, client.InNamespace(namespace))).To(Succeed())

	names := make([]string, 0)
	for _, item := range secretList.Items {
		if strings.HasPrefix(item.Name, kube.NormalizeIdentifier(project.Spec.Name)) {
			names = append(names, item.Name)
		}
	}
	Expect(names).To(HaveLen(length), fmt.Sprintf("Expected %d items, but found %d (%v)", length, len(names), names))
}

func buildConnectionURL(connURL, userName, password string) string {
	if connURL == "" {
		return ""
	}

	u, err := connectionsecret.AddCredentialsToConnectionURL(connURL, userName, password)
	Expect(err).NotTo(HaveOccurred())
	return u
}

func mongoClient(projectID string, deployment akov2.AtlasDeployment, user akov2.AtlasDatabaseUser) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c, _, err := atlasClient.ClustersApi.
		GetCluster(context.Background(), projectID, deployment.GetDeploymentName()).
		Execute()
	Expect(err).NotTo(HaveOccurred())

	if c.ConnectionStrings == nil {
		return nil, errors.New("connection strings are not provided")
	}

	cs, err := url.Parse(c.ConnectionStrings.GetStandardSrv())
	Expect(err).NotTo(HaveOccurred())

	password, err := user.ReadPassword(context.Background(), k8sClient)
	Expect(err).NotTo(HaveOccurred())
	cs.User = url.UserPassword(user.Spec.Username, password)

	dbClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cs.String()))
	if err != nil {
		return nil, err
	}
	err = dbClient.Ping(context.Background(), nil)

	return dbClient, err
}

func tryConnect(projectID string, deployment akov2.AtlasDeployment, user akov2.AtlasDatabaseUser) error {
	_, err := mongoClient(projectID, deployment, user)
	return err
}

func tryWrite(projectID string, deployment akov2.AtlasDeployment, user akov2.AtlasDatabaseUser, dbName, collectionName string) error {
	dbClient, err := mongoClient(projectID, deployment, user)
	Expect(err).NotTo(HaveOccurred())
	defer func() {
		if err = dbClient.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	collection := dbClient.Database(dbName).Collection(collectionName)

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	p := Person{
		Name: "Patrick",
		Age:  32,
	}

	_, err = collection.InsertOne(context.Background(), p)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "name", Value: "Patrick"}}

	var s Person

	err = collection.FindOne(context.Background(), filter).Decode(&s)
	Expect(err).NotTo(HaveOccurred())
	// Shouldn't return the error - by this step the roles should be propagated
	Expect(s).To(Equal(p))
	return nil
}

func checkAtlasDatabaseUserRemoved(projectID string, user akov2.AtlasDatabaseUser) func() bool {
	return func() bool {
		_, r, err := atlasClient.DatabaseUsersApi.
			GetDatabaseUser(context.Background(), projectID, user.Spec.DatabaseName, user.Spec.Username).
			Execute()
		if err != nil {
			if r != nil && r.StatusCode == http.StatusNotFound {
				return true
			}
		}

		return false
	}
}

func checkSecretsDontExist(namespace string, secretNames []string) func() bool {
	return func() bool {
		nonExisting := 0
		for _, name := range secretNames {
			s := corev1.Secret{}
			err := k8sClient.Get(context.Background(), kube.ObjectKey(namespace, name), &s)
			if err != nil && apiErrors.IsNotFound(err) {
				nonExisting++
			}
		}
		return nonExisting == len(secretNames)
	}
}

func checkUserInAtlas(projectID string, user akov2.AtlasDatabaseUser) {
	By("Verifying Database User state in Atlas", func() {
		atlasDBUser, _, err := atlasClient.DatabaseUsersApi.
			GetDatabaseUser(context.Background(), projectID, user.Spec.DatabaseName, user.Spec.Username).
			Execute()
		Expect(err).ToNot(HaveOccurred())
		atlasDBUser.Links = nil
		operatorDBUser, err := user.ToAtlasSDK(context.Background(), k8sClient)
		Expect(err).ToNot(HaveOccurred())

		Expect(*atlasDBUser).To(Equal(normalize(*operatorDBUser, projectID)))
	})
}

// normalize brings the operator 'user' to the user returned by Atlas that allows to perform comparison for equality
func normalize(user admin.CloudDatabaseUser, projectID string) admin.CloudDatabaseUser {
	if !user.HasScopes() {
		user.SetScopes([]admin.UserScope{})
	}
	if !user.HasLabels() {
		user.SetLabels([]admin.ComponentLabel{})
	}

	user.SetGroupId(projectID)
	user.Password = nil

	return user
}

func deleteSecret(user *akov2.AtlasDatabaseUser) {
	secret := &corev1.Secret{}
	Expect(
		k8sClient.Get(
			context.Background(),
			client.ObjectKey{Namespace: user.Namespace, Name: user.Spec.PasswordSecret.Name},
			secret,
			&client.GetOptions{},
		),
	).To(Succeed())

	Expect(k8sClient.Delete(context.Background(), secret)).To(Succeed())
}
