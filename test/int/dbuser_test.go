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
	"go.mongodb.org/atlas/mongodbatlas"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/testutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/timeutil"
)

const (
	DevMode             = false
	UserPasswordSecret  = "user-password-secret"
	DBUserPassword      = "Passw0rd!"
	UserPasswordSecret2 = "second-user-password-secret"
	DBUserPassword2     = "H@lla#!"
	// M2 deployments take longer time to apply changes
	DBUserUpdateTimeout = time.Minute * 4
)

var _ = Describe("AtlasDatabaseUser", Label("int", "AtlasDatabaseUser"), func() {
	const (
		interval      = PollingInterval
		intervalShort = time.Second * 2
	)

	var (
		connectionSecret       corev1.Secret
		createdProject         *mdbv1.AtlasProject
		createdDeploymentAWS   *mdbv1.AtlasDeployment
		createdDeploymentGCP   *mdbv1.AtlasDeployment
		createdDeploymentAzure *mdbv1.AtlasDeployment
		createdDBUser          *mdbv1.AtlasDatabaseUser
		secondDBUser           *mdbv1.AtlasDatabaseUser
	)

	BeforeEach(func() {
		prepareControllers()
		createdDBUser = &mdbv1.AtlasDatabaseUser{}

		connectionSecret = buildConnectionSecret("my-atlas-key")
		Expect(k8sClient.Create(context.Background(), &connectionSecret)).To(Succeed())

		By(fmt.Sprintf("Creating password Secret %s", UserPasswordSecret))
		passwordSecret := buildPasswordSecret(UserPasswordSecret, DBUserPassword)
		Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())

		By(fmt.Sprintf("Creating password Secret %s", UserPasswordSecret2))
		passwordSecret2 := buildPasswordSecret(UserPasswordSecret2, DBUserPassword2)
		Expect(k8sClient.Create(context.Background(), &passwordSecret2)).To(Succeed())

		By("Creating the project", func() {
			// adding whitespace to the name to check normalization for connection secrets names
			createdProject = mdbv1.DefaultProject(namespace.Name, connectionSecret.Name).
				WithAtlasName(namespace.Name + " some").
				WithIPAccessList(project.NewIPAccessList().WithCIDR("0.0.0.0/0"))
			if DevMode {
				// While developing tests we need to reuse the same project
				createdProject.Spec.Name = "dev-test atlas-project"
			}

			Expect(k8sClient.Create(context.Background(), createdProject)).To(Succeed())
			Eventually(func() bool {
				return testutil.CheckCondition(k8sClient, createdProject, status.TrueCondition(status.ReadyType))
			}).WithTimeout(ProjectCreationTimeout).WithPolling(interval).Should(BeTrue())
		})
	})

	AfterEach(func() {
		if DevMode {
			// No tearDown in dev mode - projects and both deployments will stay in Atlas so it's easier to develop
			// tests. Just rerun the test and the project + deployments in Atlas will be reused.
			// We only need to wipe data in the databases.
			if createdDeploymentAWS != nil {
				dbClient, err := mongoClient(createdProject.ID(), *createdDeploymentAWS, *createdDBUser)
				if err == nil {
					_ = dbClient.Database("test").Collection("operatortest").Drop(context.Background())
				}
			}
			if createdDeploymentGCP != nil {
				dbClient, err := mongoClient(createdProject.ID(), *createdDeploymentGCP, *createdDBUser)
				if err == nil {
					_ = dbClient.Database("test").Collection("operatortest").Drop(context.Background())
				}
			}
			if createdDeploymentAzure != nil {
				dbClient, err := mongoClient(createdProject.ID(), *createdDeploymentAzure, *createdDBUser)
				if err == nil {
					_ = dbClient.Database("test").Collection("operatortest").Drop(context.Background())
				}
			}

			Expect(k8sClient.Delete(context.Background(), createdDBUser)).To(Succeed())
			Eventually(checkAtlasDatabaseUserRemoved(createdProject.ID(), *createdDBUser), 20, interval).Should(BeTrue())
			if secondDBUser != nil {
				Expect(k8sClient.Delete(context.Background(), secondDBUser)).To(Succeed())
				Eventually(checkAtlasDatabaseUserRemoved(createdProject.ID(), *secondDBUser), 20, interval).Should(BeTrue())
			}
			return
		}

		if createdProject != nil && createdProject.ID() != "" {
			list := mdbv1.AtlasDeploymentList{}
			Expect(k8sClient.List(context.Background(), &list, client.InNamespace(namespace.Name))).To(Succeed())

			for i := range list.Items {
				By("Removing Atlas Deployment " + list.Items[i].Name)
				Expect(k8sClient.Delete(context.Background(), &list.Items[i])).To(Succeed())
			}
			for i := range list.Items {
				Eventually(checkAtlasDeploymentRemoved(createdProject.ID(), list.Items[i].GetDeploymentName()), 600, interval).Should(BeTrue())
			}

			By("Removing Atlas Project " + createdProject.Status.ID)
			Expect(k8sClient.Delete(context.Background(), createdProject)).To(Succeed())
			Eventually(checkAtlasProjectRemoved(createdProject.Status.ID), 60, interval).Should(BeTrue())
		}
		removeControllersAndNamespace()
	})

	connSecretname := func(suffix string) string {
		return kube.NormalizeIdentifier(createdProject.Spec.Name) + suffix
	}

	byCreatingDefaultAWSandAzureDeployments := func() {
		By("Creating deployments", func() {
			createdDeploymentAWS = mdbv1.DefaultAWSDeployment(namespace.Name, createdProject.Name).Lightweight()
			Expect(k8sClient.Create(context.Background(), createdDeploymentAWS)).ToNot(HaveOccurred())

			createdDeploymentAzure = mdbv1.DefaultAzureDeployment(namespace.Name, createdProject.Name).Lightweight()
			Expect(k8sClient.Create(context.Background(), createdDeploymentAzure)).ToNot(HaveOccurred())

			Eventually(func(g Gomega) bool {
				return testutil.CheckCondition(k8sClient, createdDeploymentAWS, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc(g))
			}).WithTimeout(DeploymentUpdateTimeout).WithPolling(interval).Should(BeTrue())
			Eventually(func(g Gomega) bool {
				return testutil.CheckCondition(k8sClient, createdDeploymentAzure, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc(g))
			}).WithTimeout(DeploymentUpdateTimeout).WithPolling(interval).Should(BeTrue())

		})
	}

	Describe("Create/Update two users, two deployments", func() {
		It("They should be created successfully", func() {
			byCreatingDefaultAWSandAzureDeployments()
			createdDBUser = mdbv1.DefaultDBUser(namespace.Name, "test-db-user", createdProject.Name).WithPasswordSecret(UserPasswordSecret)

			By(fmt.Sprintf("Creating the Database User %s", kube.ObjectKeyFromObject(createdDBUser)), func() {
				Expect(k8sClient.Create(context.Background(), createdDBUser)).ToNot(HaveOccurred())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType))
				}).WithTimeout(DBUserUpdateTimeout).WithPolling(interval).Should(BeTrue())

				checkUserInAtlas(createdProject.ID(), *createdDBUser)

				Expect(tryConnect(createdProject.ID(), *createdDeploymentAzure, *createdDBUser)).Should(Succeed())
				Expect(tryConnect(createdProject.ID(), *createdDeploymentAWS, *createdDBUser)).Should(Succeed())
				By("Checking connection Secrets", func() {
					validateSecret(k8sClient, *createdProject, *createdDeploymentAzure, *createdDBUser)
					validateSecret(k8sClient, *createdProject, *createdDeploymentAWS, *createdDBUser)
					checkNumberOfConnectionSecrets(k8sClient, *createdProject, 2)
				})
				By("Checking connectivity to Deployments", func() {
					// The user created lacks read/write roles
					err := tryWrite(createdProject.ID(), *createdDeploymentAzure, *createdDBUser, "test", "operatortest")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(MatchRegexp("user is not allowed"))

					err = tryWrite(createdProject.ID(), *createdDeploymentAWS, *createdDBUser, "test", "operatortest")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(MatchRegexp("user is not allowed"))
				})
			})
			By("Update database user - give readWrite permissions", func() {
				// Adding the role allowing read/write
				createdDBUser = createdDBUser.WithRole("readWriteAnyDatabase", "admin", "")

				Expect(k8sClient.Update(context.Background(), createdDBUser)).ToNot(HaveOccurred())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType))
				}).WithTimeout(DBUserUpdateTimeout).WithPolling(interval).Should(BeTrue())

				checkUserInAtlas(createdProject.ID(), *createdDBUser)

				By("Checking connection Secrets", func() {
					validateSecret(k8sClient, *createdProject, *createdDeploymentAzure, *createdDBUser)
					validateSecret(k8sClient, *createdProject, *createdDeploymentAWS, *createdDBUser)
					checkNumberOfConnectionSecrets(k8sClient, *createdProject, 2)
				})

				By("Checking write permissions for Deployments", func() {
					Expect(tryWrite(createdProject.ID(), *createdDeploymentAzure, *createdDBUser, "test", "operatortest")).Should(Succeed())
					Expect(tryWrite(createdProject.ID(), *createdDeploymentAWS, *createdDBUser, "test", "operatortest")).Should(Succeed())
				})
			})
			By("Adding second user for Azure deployment only (fails, wrong scope)", func() {
				secondDBUser = mdbv1.DefaultDBUser(namespace.Name, "second-db-user", createdProject.Name).
					WithPasswordSecret(UserPasswordSecret2).
					WithRole("readWrite", "someDB", "thisIsTheOnlyAllowedCollection").
					// Deployment doesn't exist
					WithScope(mdbv1.DeploymentScopeType, createdDeploymentAzure.GetDeploymentName()+"-foo")

				Expect(k8sClient.Create(context.Background(), secondDBUser)).ToNot(HaveOccurred())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, secondDBUser,
						status.
							FalseCondition(status.DatabaseUserReadyType).
							WithReason(string(workflow.DatabaseUserInvalidSpec)).
							WithMessageRegexp("such deployment doesn't exist in Atlas"))
				}).WithTimeout(DBUserUpdateTimeout).WithPolling(intervalShort).Should(BeTrue())
			})
			By("Fixing second user", func() {
				secondDBUser = secondDBUser.ClearScopes().WithScope(mdbv1.DeploymentScopeType, createdDeploymentAzure.Spec.DeploymentSpec.Name)

				Expect(k8sClient.Update(context.Background(), secondDBUser)).ToNot(HaveOccurred())

				// First we need to wait for "such deployment doesn't exist in Atlas" error to be gone
				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, secondDBUser,
						status.FalseCondition(status.DatabaseUserReadyType).
							WithReason(string(workflow.DatabaseUserDeploymentAppliedChanges)))
				}).WithTimeout(DBUserUpdateTimeout).WithPolling(intervalShort).Should(BeTrue())

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, secondDBUser,
						status.TrueCondition(status.ReadyType), validateDatabaseUserUpdatingFunc(g))
				}).WithTimeout(DBUserUpdateTimeout).WithPolling(interval).Should(BeTrue())

				checkUserInAtlas(createdProject.ID(), *secondDBUser)

				By("Checking connection Secrets", func() {
					validateSecret(k8sClient, *createdProject, *createdDeploymentAzure, *createdDBUser)
					validateSecret(k8sClient, *createdProject, *createdDeploymentAWS, *createdDBUser)
					validateSecret(k8sClient, *createdProject, *createdDeploymentAzure, *secondDBUser)
					checkNumberOfConnectionSecrets(k8sClient, *createdProject, 3)
				})

				By("Checking write permissions for Deployments", func() {
					// We still can write by the first user
					Expect(tryWrite(createdProject.ID(), *createdDeploymentAzure, *createdDBUser, "test", "testCollection")).Should(Succeed())
					Expect(tryWrite(createdProject.ID(), *createdDeploymentAWS, *createdDBUser, "test", "testCollection")).Should(Succeed())

					// The second user can eventually write to one collection only
					Expect(tryConnect(createdProject.ID(), *createdDeploymentAzure, *secondDBUser)).Should(Succeed())
					Expect(tryWrite(createdProject.ID(), *createdDeploymentAzure, *secondDBUser, "someDB", "thisIsTheOnlyAllowedCollection")).Should(Succeed())

					err := tryWrite(createdProject.ID(), *createdDeploymentAzure, *secondDBUser, "test", "someNotAllowedCollection")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(MatchRegexp("user is not allowed"))
				})
				By("Removing Second user", func() {
					Expect(k8sClient.Delete(context.Background(), secondDBUser)).To(Succeed())
					Eventually(checkAtlasDatabaseUserRemoved(createdProject.Status.ID, *secondDBUser), 50, interval).Should(BeTrue())

					secretNames := []string{connSecretname("-test-deployment-azure-second-db-user")}
					Eventually(checkSecretsDontExist(namespace.Name, secretNames), 50, interval).Should(BeTrue())
				})
			})
			By("Removing First user", func() {
				Expect(k8sClient.Delete(context.Background(), createdDBUser)).To(Succeed())
				Eventually(checkAtlasDatabaseUserRemoved(createdProject.Status.ID, *createdDBUser), 50, interval).Should(BeTrue())

				secretNames := []string{connSecretname("-test-deployment-aws-test-db-user"), connSecretname("-test-deployment-azure-test-db-user")}
				Eventually(checkSecretsDontExist(namespace.Name, secretNames), 50, interval).Should(BeTrue())

				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 0)
			})
		})
	})

	// Note, that this test doesn't work with "DevMode=true" as requires the deployment to get created
	Describe("Check the reverse order of deployment-user creation (user - first, then - the deployment)", func() {
		It("Should succeed", func() {
			// Here we create a database user first - then the deployment
			By("Creating database user", func() {
				createdDBUser = mdbv1.DefaultDBUser(namespace.Name, "test-db-user", createdProject.Name).WithPasswordSecret(UserPasswordSecret)

				Expect(k8sClient.Create(context.Background(), createdDBUser)).To(Succeed())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType))
				}).WithTimeout(DBUserUpdateTimeout).WithPolling(interval).Should(BeTrue())

				checkUserInAtlas(createdProject.ID(), *createdDBUser)
				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 0)
			})
			By("Creating deployment", func() {
				createdDeploymentAWS = mdbv1.DefaultAWSDeployment(namespace.Name, createdProject.Name).Lightweight()
				Expect(k8sClient.Create(context.Background(), createdDeploymentAWS)).ToNot(HaveOccurred())

				// We don't wait for the full deployment creation - only when it has started the process
				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdDeploymentAWS, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc(g))
				}).WithTimeout(DeploymentUpdateTimeout).WithPolling(interval).Should(BeTrue())
			})
			By("Updating the database user while the deployment is being created", func() {
				createdDBUser = createdDBUser.WithRole("read", "test", "somecollection")
				Expect(k8sClient.Update(context.Background(), createdDBUser)).To(Succeed())

				// DatabaseUser will wait for the deployment to get created.
				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType))
				}).WithTimeout(DBUserUpdateTimeout).WithPolling(interval).Should(BeTrue())

				expectedConditionsMatchers := testutil.MatchConditions(
					status.TrueCondition(status.DatabaseUserReadyType),
					status.TrueCondition(status.ReadyType),
					status.TrueCondition(status.ValidationSucceeded),
					status.TrueCondition(status.ResourceVersionStatus),
				)
				Expect(createdDBUser.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))

				checkUserInAtlas(createdProject.ID(), *createdDBUser)
				Expect(tryConnect(createdProject.ID(), *createdDeploymentAWS, *createdDBUser)).Should(Succeed())
				By("Checking connection Secrets", func() {
					validateSecret(k8sClient, *createdProject, *createdDeploymentAWS, *createdDBUser)
					checkNumberOfConnectionSecrets(k8sClient, *createdProject, 1)
				})
			})
		})
	})
	Describe("Check the password Secret is watched", func() {
		It("Should succeed", func() {
			By("Creating deployments", func() {
				createdDeploymentAWS = mdbv1.DefaultAWSDeployment(namespace.Name, createdProject.Name).Lightweight()
				Expect(k8sClient.Create(context.Background(), createdDeploymentAWS)).ToNot(HaveOccurred())

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdDeploymentAWS, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc(g))
				}).WithTimeout(DeploymentUpdateTimeout).WithPolling(interval).Should(BeTrue())

			})
			createdDBUser = mdbv1.DefaultDBUser(namespace.Name, "test-db-user", createdProject.Name).WithPasswordSecret(UserPasswordSecret)
			var connSecretInitial corev1.Secret
			var pwdSecret corev1.Secret

			By(fmt.Sprintf("Creating the Database User %s", kube.ObjectKeyFromObject(createdDBUser)), func() {
				Expect(k8sClient.Create(context.Background(), createdDBUser)).ToNot(HaveOccurred())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType))
				}).WithTimeout(DBUserUpdateTimeout).WithPolling(interval).Should(BeTrue())
				testutil.EventExists(k8sClient, createdDBUser, "Normal", "Ready", "")

				Expect(tryConnect(createdProject.ID(), *createdDeploymentAWS, *createdDBUser)).Should(Succeed())

				connSecretInitial = validateSecret(k8sClient, *createdProject, *createdDeploymentAWS, *createdDBUser)
				Expect(k8sClient.Get(context.Background(), kube.ObjectKey(namespace.Name, UserPasswordSecret), &pwdSecret)).To(Succeed())
				Expect(createdDBUser.Status.PasswordVersion).To(Equal(pwdSecret.ResourceVersion))
			})

			By("Breaking the password secret", func() {
				passwordSecret := buildPasswordSecret(UserPasswordSecret, "")
				Expect(k8sClient.Update(context.Background(), &passwordSecret)).To(Succeed())

				expectedCondition := status.FalseCondition(status.DatabaseUserReadyType).WithReason(string(workflow.Internal)).WithMessageRegexp("the 'password' field is empty")
				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdDBUser, expectedCondition)
				}).WithTimeout(DBUserUpdateTimeout).WithPolling(interval).Should(BeTrue())
				testutil.EventExists(k8sClient, createdDBUser, "Warning", string(workflow.Internal), "the 'password' field is empty")
			})
			By("Fixing the password secret", func() {
				passwordSecret := buildPasswordSecret(UserPasswordSecret, "someNewPassw00rd")
				Expect(k8sClient.Update(context.Background(), &passwordSecret)).To(Succeed())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType))
				}).WithTimeout(DBUserUpdateTimeout).WithPolling(interval).Should(BeTrue())

				// We need to make sure that the new connection secret is different from the initial one
				connSecretUpdated := validateSecret(k8sClient, *createdProject, *createdDeploymentAWS, *createdDBUser)
				Expect(string(connSecretInitial.Data["password"])).To(Equal(DBUserPassword))
				Expect(string(connSecretUpdated.Data["password"])).To(Equal("someNewPassw00rd"))

				var updatedPwdSecret corev1.Secret
				Expect(k8sClient.Get(context.Background(), kube.ObjectKey(namespace.Name, UserPasswordSecret), &updatedPwdSecret)).To(Succeed())
				Expect(updatedPwdSecret.ResourceVersion).NotTo(Equal(pwdSecret.ResourceVersion))
				Expect(createdDBUser.Status.PasswordVersion).To(Equal(updatedPwdSecret.ResourceVersion))

				Expect(tryConnect(createdProject.ID(), *createdDeploymentAWS, *createdDBUser)).Should(Succeed())
			})
		})
	})
	Describe("Change database users (make sure all stale secrets are removed)", func() {
		It("Should succeed", func() {
			byCreatingDefaultAWSandAzureDeployments()
			createdDBUser = mdbv1.DefaultDBUser(namespace.Name, "test-db-user", createdProject.Name).WithPasswordSecret(UserPasswordSecret)

			By(fmt.Sprintf("Creating the Database User %s (no scopes)", kube.ObjectKeyFromObject(createdDBUser)), func() {
				Expect(k8sClient.Create(context.Background(), createdDBUser)).To(Succeed())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType))
				}).WithTimeout(DBUserUpdateTimeout).WithPolling(interval).Should(BeTrue())

				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 2)

				s1 := validateSecret(k8sClient, *createdProject, *createdDeploymentAWS, *createdDBUser)
				s2 := validateSecret(k8sClient, *createdProject, *createdDeploymentAzure, *createdDBUser)

				testutil.EventExists(k8sClient, createdDBUser, "Normal", connectionsecret.ConnectionSecretsEnsuredEvent,
					fmt.Sprintf("Connection Secrets were created/updated: (%s|%s|, ){3}", s1.Name, s2.Name))
			})
			By("Changing the db user name - two stale secret are expected to be removed, two added instead", func() {
				oldName := createdDBUser.Spec.Username
				createdDBUser = createdDBUser.WithAtlasUserName("new-user")
				Expect(k8sClient.Update(context.Background(), createdDBUser)).To(Succeed())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType))
				}).WithTimeout(DBUserUpdateTimeout).WithPolling(interval).Should(BeTrue())

				checkUserInAtlas(createdProject.ID(), *createdDBUser)
				// Old user has been removed
				_, _, err := atlasClient.DatabaseUsers.Get(context.Background(), createdDBUser.Spec.DatabaseName, createdProject.ID(), oldName)
				Expect(err).To(HaveOccurred())

				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 2)
				secret := validateSecret(k8sClient, *createdProject, *createdDeploymentAzure, *createdDBUser)
				Expect(secret.Name).To(Equal(connSecretname("-test-deployment-azure-new-user")))
				secret = validateSecret(k8sClient, *createdProject, *createdDeploymentAWS, *createdDBUser)
				Expect(secret.Name).To(Equal(connSecretname("-test-deployment-aws-new-user")))

				Expect(tryConnect(createdProject.ID(), *createdDeploymentAzure, *createdDBUser)).Should(Succeed())
				Expect(tryConnect(createdProject.ID(), *createdDeploymentAWS, *createdDBUser)).Should(Succeed())
			})
			By("Changing the scopes - one stale secret is expected to be removed", func() {
				createdDBUser = createdDBUser.ClearScopes().WithScope(mdbv1.DeploymentScopeType, createdDeploymentAzure.Spec.DeploymentSpec.Name)
				Expect(k8sClient.Update(context.Background(), createdDBUser)).To(Succeed())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType))
				}).WithTimeout(DBUserUpdateTimeout).WithPolling(interval).Should(BeTrue())

				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 1)
				validateSecret(k8sClient, *createdProject, *createdDeploymentAzure, *createdDBUser)

				Expect(tryConnect(createdProject.ID(), *createdDeploymentAzure, *createdDBUser)).Should(Succeed())
				Expect(tryConnect(createdProject.ID(), *createdDeploymentAWS, *createdDBUser)).ShouldNot(Succeed())
			})
		})
	})
	Describe("Check the user expiration", func() {
		It("Should succeed", func() {
			By("Creating a AWS deployment", func() {
				createdDeploymentAWS = mdbv1.DefaultAWSDeployment(namespace.Name, createdProject.Name).Lightweight()
				Expect(k8sClient.Create(context.Background(), createdDeploymentAWS)).To(Succeed())

				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, createdDeploymentAWS, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc(g))
				}).WithTimeout(DeploymentUpdateTimeout).WithPolling(intervalShort).Should(BeTrue())
			})

			By("Creating the expired Database User - no user created in Atlas", func() {
				before := time.Now().UTC().Add(time.Minute * -10).Format("2006-01-02T15:04:05")
				createdDBUser = mdbv1.DefaultDBUser(namespace.Name, "test-db-user", createdProject.Name).
					WithPasswordSecret(UserPasswordSecret).
					WithDeleteAfterDate(before)

				Expect(k8sClient.Create(context.Background(), createdDBUser)).To(Succeed())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdDBUser,
						status.FalseCondition(status.DatabaseUserReadyType).WithReason(string(workflow.DatabaseUserExpired)))
				}).WithTimeout(DBUserUpdateTimeout).WithPolling(intervalShort).Should(BeTrue())

				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 0)

				// no user in Atlas
				_, _, err := atlasClient.DatabaseUsers.Get(context.Background(), createdDBUser.Spec.DatabaseName, createdProject.ID(), createdDBUser.Spec.Username)
				Expect(err).To(HaveOccurred())
			})
			By("Fixing the Database User - setting the expiration to future", func() {
				after := time.Now().UTC().Add(time.Hour * 10).Format("2006-01-02T15:04:05")
				createdDBUser = createdDBUser.WithDeleteAfterDate(after)

				Expect(k8sClient.Update(context.Background(), createdDBUser)).To(Succeed())
				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType))
				}).WithTimeout(DBUserUpdateTimeout).WithPolling(interval).Should(BeTrue())

				checkUserInAtlas(createdProject.ID(), *createdDBUser)
				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 1)
				Expect(tryConnect(createdProject.ID(), *createdDeploymentAWS, *createdDBUser)).Should(Succeed())
			})
			By("Extending the expiration", func() {
				after := time.Now().UTC().Add(time.Hour * 30).Format("2006-01-02T15:04:05")
				createdDBUser = createdDBUser.WithDeleteAfterDate(after)

				Expect(k8sClient.Update(context.Background(), createdDBUser)).To(Succeed())
				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType))
				}).WithTimeout(DBUserUpdateTimeout).WithPolling(interval).Should(BeTrue())

				checkUserInAtlas(createdProject.ID(), *createdDBUser)
			})
			By("Emulating expiration of the User - connection secret must be removed", func() {
				before := time.Now().UTC().Add(time.Minute * -5).Format("2006-01-02T15:04:05")
				createdDBUser = createdDBUser.WithDeleteAfterDate(before)

				Expect(k8sClient.Update(context.Background(), createdDBUser)).To(Succeed())
				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, createdDBUser, status.FalseCondition(status.DatabaseUserReadyType).WithReason(string(workflow.DatabaseUserExpired)))
				}).WithTimeout(DBUserUpdateTimeout).WithPolling(intervalShort).Should(BeTrue())

				expectedConditionsMatchers := testutil.MatchConditions(
					status.FalseCondition(status.DatabaseUserReadyType),
					status.FalseCondition(status.ReadyType),
					status.TrueCondition(status.ValidationSucceeded),
					status.TrueCondition(status.ResourceVersionStatus),
				)
				Expect(createdDBUser.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))

				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 0)
			})
		})
		Describe("Deleting the db user (not cleaning Atlas)", func() {
			It("Should Succeed", func() {
				By(`Creating the db user with retention policy "keep" first`, func() {
					createdDeploymentAWS = mdbv1.DefaultAWSDeployment(namespace.Name, createdProject.Name).Lightweight()
					Expect(k8sClient.Create(context.Background(), createdDeploymentAWS)).ToNot(HaveOccurred())

					Eventually(func(g Gomega) bool {
						return testutil.CheckCondition(k8sClient, createdDeploymentAWS, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc(g))
					}).WithTimeout(DeploymentUpdateTimeout).WithPolling(interval).Should(BeTrue())

					createdDBUser = mdbv1.DefaultDBUser(namespace.Name, "test-db-user", createdProject.Name).WithPasswordSecret(UserPasswordSecret)
					createdDBUser.ObjectMeta.Annotations = map[string]string{customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep}
					Expect(k8sClient.Create(context.Background(), createdDBUser)).To(Succeed())
					Eventually(func() bool {
						return testutil.CheckCondition(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType))
					}).WithTimeout(DBUserUpdateTimeout).WithPolling(interval).Should(BeTrue())
				})
				By("Deleting the db user - stays in Atlas", func() {
					Expect(k8sClient.Delete(context.Background(), createdDBUser)).To(Succeed())

					time.Sleep(1 * time.Minute)
					Expect(checkAtlasDatabaseUserRemoved(createdProject.Status.ID, *createdDBUser)()).Should(BeFalse())

					checkNumberOfConnectionSecrets(k8sClient, *createdProject, 0)
				})
			})
		})

		Describe("Setting the user skip annotation should skip reconciliations.", func() {
			It("Should Succeed", func() {

				By(`Creating the user with reconciliation policy "skip" first`, func() {
					createdDeploymentAWS = mdbv1.DefaultAWSDeployment(namespace.Name, createdProject.Name).Lightweight()
					Expect(k8sClient.Create(context.Background(), createdDeploymentAWS)).ToNot(HaveOccurred())
					Eventually(func(g Gomega) bool {
						return testutil.CheckCondition(k8sClient, createdDeploymentAWS, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc(g))
					}).WithTimeout(DeploymentUpdateTimeout).WithPolling(interval).Should(BeTrue())

					createdDBUser = mdbv1.DefaultDBUser(namespace.Name, "test-db-user", createdProject.Name).WithPasswordSecret(UserPasswordSecret)

					Expect(k8sClient.Create(context.Background(), createdDBUser)).To(Succeed())
					Eventually(func() bool {
						return testutil.CheckCondition(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType))
					}).WithTimeout(DBUserUpdateTimeout).WithPolling(interval).Should(BeTrue())

					createdDBUser.ObjectMeta.Annotations = map[string]string{customresource.ReconciliationPolicyAnnotation: customresource.ReconciliationPolicySkip}
					createdDBUser.Spec.Roles = append(createdDBUser.Spec.Roles, mdbv1.RoleSpec{
						RoleName:       "new-role",
						DatabaseName:   "new-database",
						CollectionName: "new-collection",
					})

					// add the annotation to skip reconciliation and a new role. This new role should not be seen in
					// atlas.
					Expect(k8sClient.Update(context.Background(), createdDBUser)).To(Succeed())

					ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
					defer cancel()

					containsDatabaseUser := func(dbUser *mongodbatlas.DatabaseUser) bool {
						for _, role := range dbUser.Roles {
							if role.RoleName == "new-role" && role.DatabaseName == "new-database" && role.CollectionName == "new-collection" {
								return true
							}
						}
						return false
					}

					Eventually(testutil.WaitForAtlasDatabaseUserStateToNotBeReached(ctx, atlasClient, "admin", createdProject.Name, createdDeploymentAWS.GetDeploymentName(), containsDatabaseUser))
				})
			})
		})
	})
})

func buildPasswordSecret(name, password string) corev1.Secret {
	return corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace.Name,
			Labels: map[string]string{
				connectionsecret.TypeLabelKey: connectionsecret.CredLabelVal,
			},
		},
		StringData: map[string]string{"password": password},
	}
}

// normalize brings the operator 'user' to the user returned by Atlas that allows to perform comparison for equality
func normalize(user mongodbatlas.DatabaseUser, projectID string) mongodbatlas.DatabaseUser {
	if user.Scopes == nil {
		user.Scopes = []mongodbatlas.Scope{}
	}
	if user.Labels == nil {
		user.Labels = []mongodbatlas.Label{}
	}
	if user.LDAPAuthType == "" {
		user.LDAPAuthType = "NONE"
	}
	if user.AWSIAMType == "" {
		user.AWSIAMType = "NONE"
	}
	if user.X509Type == "" {
		user.X509Type = "NONE"
	}
	if user.DeleteAfterDate != "" {
		user.DeleteAfterDate = timeutil.FormatISO8601(timeutil.MustParseISO8601(user.DeleteAfterDate))
	}
	user.GroupID = projectID
	user.Password = ""
	return user
}

func tryConnect(projectID string, deployment mdbv1.AtlasDeployment, user mdbv1.AtlasDatabaseUser) error {
	_, err := mongoClient(projectID, deployment, user)
	return err
}

func mongoClient(projectID string, deployment mdbv1.AtlasDeployment, user mdbv1.AtlasDatabaseUser) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c, _, err := atlasClient.Clusters.Get(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name)
	Expect(err).NotTo(HaveOccurred())

	if c.ConnectionStrings == nil {
		return nil, errors.New("Connection strings are not provided!")
	}

	cs, err := url.Parse(c.ConnectionStrings.StandardSrv)
	Expect(err).NotTo(HaveOccurred())

	password, err := user.ReadPassword(k8sClient)
	Expect(err).NotTo(HaveOccurred())
	cs.User = url.UserPassword(user.Spec.Username, password)

	dbClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cs.String()))
	if err != nil {
		return nil, err
	}
	err = dbClient.Ping(context.TODO(), nil)

	return dbClient, err
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func tryWrite(projectID string, deployment mdbv1.AtlasDeployment, user mdbv1.AtlasDatabaseUser, dbName, collectionName string) error {
	dbClient, err := mongoClient(projectID, deployment, user)
	Expect(err).NotTo(HaveOccurred())
	defer func() {
		if err = dbClient.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	collection := dbClient.Database(dbName).Collection(collectionName)

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
	fmt.Fprintf(GinkgoWriter, "User %s (deployment %s) has inserted a single document to %s/%s\n", user.Spec.Username, deployment.GetDeploymentName(), dbName, collectionName)
	return nil
}

func validateSecret(k8sClient client.Client, project mdbv1.AtlasProject, deployment mdbv1.AtlasDeployment, user mdbv1.AtlasDatabaseUser) corev1.Secret {
	secret := corev1.Secret{}
	username := user.Spec.Username
	secretName := fmt.Sprintf("%s-%s-%s", kube.NormalizeIdentifier(project.Spec.Name), kube.NormalizeIdentifier(deployment.GetDeploymentName()), kube.NormalizeIdentifier(username))
	Expect(k8sClient.Get(context.Background(), kube.ObjectKey(project.Namespace, secretName), &secret)).To(Succeed())
	GinkgoWriter.Write([]byte(fmt.Sprintf("!! Secret: %v (%v)\n", kube.ObjectKey(project.Namespace, secretName), secret.Namespace+"/"+secret.Name)))

	password, err := user.ReadPassword(k8sClient)
	Expect(err).NotTo(HaveOccurred())

	c, _, err := atlasClient.Clusters.Get(context.Background(), project.ID(), deployment.Spec.DeploymentSpec.Name)
	Expect(err).NotTo(HaveOccurred())

	expectedData := map[string][]byte{
		"connectionStringStandard":    []byte(buildConnectionURL(c.ConnectionStrings.Standard, username, password)),
		"connectionStringStandardSrv": []byte(buildConnectionURL(c.ConnectionStrings.StandardSrv, username, password)),
		"connectionStringPrivate":     []byte(buildConnectionURL(c.ConnectionStrings.Private, username, password)),
		"connectionStringPrivateSrv":  []byte(buildConnectionURL(c.ConnectionStrings.PrivateSrv, username, password)),
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
	GinkgoWriter.Write([]byte(fmt.Sprintf("!! Secret 2: %v \n", secret.Namespace+"/"+secret.Name)))
	return secret
}

func checkNumberOfConnectionSecrets(k8sClient client.Client, project mdbv1.AtlasProject, length int) {
	secretList := corev1.SecretList{}
	Expect(k8sClient.List(context.Background(), &secretList, client.InNamespace(namespace.Name))).To(Succeed())

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

func checkAtlasDatabaseUserRemoved(projectID string, user mdbv1.AtlasDatabaseUser) func() bool {
	return func() bool {
		_, r, err := atlasClient.DatabaseUsers.Get(context.Background(), user.Spec.DatabaseName, projectID, user.Spec.Username)
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

func checkUserInAtlas(projectID string, user mdbv1.AtlasDatabaseUser) {
	By("Verifying Database User state in Atlas", func() {
		atlasDBUser, _, err := atlasClient.DatabaseUsers.Get(context.Background(), user.Spec.DatabaseName, projectID, user.Spec.Username)
		Expect(err).ToNot(HaveOccurred())
		operatorDBUser, err := user.ToAtlas(k8sClient)
		Expect(err).ToNot(HaveOccurred())

		Expect(*atlasDBUser).To(Equal(normalize(*operatorDBUser, projectID)))
	})
}

func validateDatabaseUserUpdatingFunc(g Gomega) func(a mdbv1.AtlasCustomResource) {
	return func(a mdbv1.AtlasCustomResource) {
		d := a.(*mdbv1.AtlasDatabaseUser)
		expectedConditionsMatchers := testutil.MatchConditions(
			status.FalseCondition(status.DatabaseUserReadyType).WithReason(string(workflow.DatabaseUserDeploymentAppliedChanges)),
			status.FalseCondition(status.ReadyType),
			status.TrueCondition(status.ValidationSucceeded),
			status.TrueCondition(status.ResourceVersionStatus),
		)
		g.Expect(d.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
	}
}

// nolint
func validateDatabaseUserWaitingForCluster() func(a mdbv1.AtlasCustomResource) {
	return func(a mdbv1.AtlasCustomResource) {
		d := a.(*mdbv1.AtlasDatabaseUser)
		// this is the first status that db user gets after update
		userChangesApplied := testutil.MatchConditions(
			status.FalseCondition(status.DatabaseUserReadyType).WithReason(string(workflow.DatabaseUserDeploymentAppliedChanges)),
			status.FalseCondition(status.ReadyType),
			status.TrueCondition(status.ValidationSucceeded),
		)
		// this is the status the db user gets to when tries to create connection secrets and sees that the deployment
		// is not ready
		waitingForDeployment := testutil.MatchConditions(
			status.FalseCondition(status.DatabaseUserReadyType).
				WithReason(string(workflow.DatabaseUserConnectionSecretsNotCreated)).
				WithMessageRegexp("Waiting for deployments to get created/updated"),
			status.FalseCondition(status.ReadyType),
			status.TrueCondition(status.ResourceVersionStatus),
		)
		Expect(d.Status.Conditions).To(Or(ConsistOf(waitingForDeployment), ConsistOf(userChangesApplied)))
	}
}
