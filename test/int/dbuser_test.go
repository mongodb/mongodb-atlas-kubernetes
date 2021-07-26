package int

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasdatabaseuser"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"

	"go.mongodb.org/mongo-driver/mongo"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/connectionsecret"
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
	// M2 clusters take longer time to apply changes
	DBUserUpdateTimeout = time.Minute * 4
)

var _ = Describe("AtlasDatabaseUser", func() {
	const (
		interval      = PollingInterval
		intervalShort = time.Second * 2
	)

	var (
		connectionSecret    corev1.Secret
		createdProject      *mdbv1.AtlasProject
		createdClusterAWS   *mdbv1.AtlasCluster
		createdClusterGCP   *mdbv1.AtlasCluster
		createdClusterAzure *mdbv1.AtlasCluster
		createdDBUser       *mdbv1.AtlasDatabaseUser
		secondDBUser        *mdbv1.AtlasDatabaseUser
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
				WithIPAccessList(project.NewIPAccessList().WithIP("0.0.0.0/0"))
			if DevMode {
				// While developing tests we need to reuse the same project
				createdProject.Spec.Name = "dev-test atlas-project"
			}

			Expect(k8sClient.Create(context.Background(), createdProject)).To(Succeed())
			Eventually(testutil.WaitFor(k8sClient, createdProject, status.TrueCondition(status.ReadyType)),
				ProjectCreationTimeout, interval).Should(BeTrue())
		})
	})

	AfterEach(func() {
		if DevMode {
			// No tearDown in dev mode - projects and both clusters will stay in Atlas so it's easier to develop
			// tests. Just rerun the test and the project + clusters in Atlas will be reused.
			// We only need to wipe data in the databases.
			if createdClusterAWS != nil {
				dbClient, err := mongoClient(createdProject.ID(), *createdClusterAWS, *createdDBUser)
				if err == nil {
					_ = dbClient.Database("test").Collection("operatortest").Drop(context.Background())
				}
			}
			if createdClusterGCP != nil {
				dbClient, err := mongoClient(createdProject.ID(), *createdClusterGCP, *createdDBUser)
				if err == nil {
					_ = dbClient.Database("test").Collection("operatortest").Drop(context.Background())
				}
			}
			if createdClusterAzure != nil {
				dbClient, err := mongoClient(createdProject.ID(), *createdClusterAzure, *createdDBUser)
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
			list := mdbv1.AtlasClusterList{}
			Expect(k8sClient.List(context.Background(), &list, client.InNamespace(namespace.Name))).To(Succeed())

			for i := range list.Items {
				By("Removing Atlas Cluster " + list.Items[i].Name)
				Expect(k8sClient.Delete(context.Background(), &list.Items[i])).To(Succeed())
			}
			for i := range list.Items {
				Eventually(checkAtlasClusterRemoved(createdProject.ID(), list.Items[i].Spec.Name), 600, interval).Should(BeTrue())
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

	Describe("Create/Update two users, two clusters", func() {
		It("They should be created successfully", func() {
			By("Creating clusters", func() {
				createdClusterAWS = mdbv1.DefaultAWSCluster(namespace.Name, createdProject.Name)
				Expect(k8sClient.Create(context.Background(), createdClusterAWS)).ToNot(HaveOccurred())

				createdClusterAzure = mdbv1.DefaultAzureCluster(namespace.Name, createdProject.Name)
				Expect(k8sClient.Create(context.Background(), createdClusterAzure)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdClusterAWS, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
					ClusterUpdateTimeout, interval).Should(BeTrue())

				Eventually(testutil.WaitFor(k8sClient, createdClusterAzure, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
					500, interval).Should(BeTrue())
			})
			createdDBUser = mdbv1.DefaultDBUser(namespace.Name, "test-db-user", createdProject.Name).WithPasswordSecret(UserPasswordSecret)

			By(fmt.Sprintf("Creating the Database User %s", kube.ObjectKeyFromObject(createdDBUser)), func() {
				Expect(k8sClient.Create(context.Background(), createdDBUser)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType)),
					DBUserUpdateTimeout, interval, validateDatabaseUserUpdatingFunc()).Should(BeTrue())

				checkUserInAtlas(createdProject.ID(), *createdDBUser)

				Expect(tryConnect(createdProject.ID(), *createdClusterAzure, *createdDBUser)).Should(Succeed())
				Expect(tryConnect(createdProject.ID(), *createdClusterAWS, *createdDBUser)).Should(Succeed())
				By("Checking connection Secrets", func() {
					validateSecret(k8sClient, *createdProject, *createdClusterAzure, *createdDBUser)
					validateSecret(k8sClient, *createdProject, *createdClusterAWS, *createdDBUser)
					checkNumberOfConnectionSecrets(k8sClient, *createdProject, 2)
				})
				By("Checking connectivity to Clusters", func() {
					// The user created lacks read/write roles
					err := tryWrite(createdProject.ID(), *createdClusterAzure, *createdDBUser, "test", "operatortest")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(MatchRegexp("not authorized on test to execute command")) // TODO check this is the same "user is not allowed"

					err = tryWrite(createdProject.ID(), *createdClusterAWS, *createdDBUser, "test", "operatortest")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(MatchRegexp("not authorized on test to execute command")) // TODO check this is the same "user is not allowed"
				})
			})
			By("Update database user - give readWrite permissions", func() {
				// Adding the role allowing read/write
				createdDBUser = createdDBUser.WithRole("readWriteAnyDatabase", "admin", "")

				Expect(k8sClient.Update(context.Background(), createdDBUser)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType)),
					DBUserUpdateTimeout, interval).Should(BeTrue())

				checkUserInAtlas(createdProject.ID(), *createdDBUser)

				By("Checking connection Secrets", func() {
					validateSecret(k8sClient, *createdProject, *createdClusterAzure, *createdDBUser)
					validateSecret(k8sClient, *createdProject, *createdClusterAWS, *createdDBUser)
					checkNumberOfConnectionSecrets(k8sClient, *createdProject, 2)
				})

				By("Checking write permissions for Clusters", func() {
					Expect(tryWrite(createdProject.ID(), *createdClusterAzure, *createdDBUser, "test", "operatortest")).Should(Succeed())
					Expect(tryWrite(createdProject.ID(), *createdClusterAWS, *createdDBUser, "test", "operatortest")).Should(Succeed())
				})
			})
			By("Adding second user for Azure cluster only (fails, wrong scope)", func() {
				secondDBUser = mdbv1.DefaultDBUser(namespace.Name, "second-db-user", createdProject.Name).
					WithPasswordSecret(UserPasswordSecret2).
					WithRole("readWrite", "someDB", "thisIsTheOnlyAllowedCollection").
					// Cluster doesn't exist
					WithScope(mdbv1.ClusterScopeType, createdClusterAzure.Spec.Name+"-foo")

				Expect(k8sClient.Create(context.Background(), secondDBUser)).ToNot(HaveOccurred())

				Eventually(
					testutil.WaitFor(
						k8sClient,
						secondDBUser,
						status.
							FalseCondition(status.DatabaseUserReadyType).
							WithReason(string(workflow.DatabaseUserInvalidSpec)).
							WithMessageRegexp("such cluster doesn't exist in Atlas"),
					),
					20,
					intervalShort,
				).Should(BeTrue())
			})
			By("Fixing second user", func() {
				secondDBUser = secondDBUser.ClearScopes().WithScope(mdbv1.ClusterScopeType, createdClusterAzure.Spec.Name)

				Expect(k8sClient.Update(context.Background(), secondDBUser)).ToNot(HaveOccurred())

				// First we need to wait for "such cluster doesn't exist in Atlas" error to be gone
				Eventually(testutil.WaitFor(k8sClient, secondDBUser, status.FalseCondition(status.DatabaseUserReadyType).WithReason(string(workflow.DatabaseUserClustersAppliedChanges))),
					20, intervalShort).Should(BeTrue())

				Eventually(testutil.WaitFor(k8sClient, secondDBUser, status.TrueCondition(status.ReadyType), validateDatabaseUserUpdatingFunc()),
					DBUserUpdateTimeout, interval).Should(BeTrue())

				checkUserInAtlas(createdProject.ID(), *secondDBUser)

				By("Checking connection Secrets", func() {
					validateSecret(k8sClient, *createdProject, *createdClusterAzure, *createdDBUser)
					validateSecret(k8sClient, *createdProject, *createdClusterAWS, *createdDBUser)
					validateSecret(k8sClient, *createdProject, *createdClusterAzure, *secondDBUser)
					checkNumberOfConnectionSecrets(k8sClient, *createdProject, 3)
				})

				By("Checking write permissions for Clusters", func() {
					// We still can write by the first user
					Expect(tryWrite(createdProject.ID(), *createdClusterAzure, *createdDBUser, "test", "testCollection")).Should(Succeed())
					Expect(tryWrite(createdProject.ID(), *createdClusterAWS, *createdDBUser, "test", "testCollection")).Should(Succeed())

					// The second user can eventually write to one collection only
					Expect(tryConnect(createdProject.ID(), *createdClusterAzure, *secondDBUser)).Should(Succeed())
					Expect(tryWrite(createdProject.ID(), *createdClusterAzure, *secondDBUser, "someDB", "thisIsTheOnlyAllowedCollection")).Should(Succeed())

					err := tryWrite(createdProject.ID(), *createdClusterAzure, *secondDBUser, "test", "someNotAllowedCollection")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(MatchRegexp("not authorized on test to execute command")) // TODO check this is the same "user is not allowed"
				})
				By("Removing Second user", func() {
					Expect(k8sClient.Delete(context.Background(), secondDBUser)).To(Succeed())
					Eventually(checkAtlasDatabaseUserRemoved(createdProject.Status.ID, *secondDBUser), 50, interval).Should(BeTrue())

					secretNames := []string{connSecretname("-test-cluster-azure-second-db-user")}
					Eventually(checkSecretsDontExist(namespace.Name, secretNames), 50, interval).Should(BeTrue())
				})
			})
			By("Removing First user", func() {
				Expect(k8sClient.Delete(context.Background(), createdDBUser)).To(Succeed())
				Eventually(checkAtlasDatabaseUserRemoved(createdProject.Status.ID, *createdDBUser), 50, interval).Should(BeTrue())

				secretNames := []string{connSecretname("-test-cluster-aws-test-db-user"), connSecretname("-test-cluster-azure-test-db-user")}
				Eventually(checkSecretsDontExist(namespace.Name, secretNames), 50, interval).Should(BeTrue())

				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 0)
			})
		})
	})

	// Note, that this test doesn't work with "DevMode=true" as requires the cluster to get created
	Describe("Check the reverse order of cluster-user creation (user - first, then - the cluster)", func() {
		It("Should succeed", func() {
			// Here we create a database user first - then the cluster
			By("Creating database user", func() {
				createdDBUser = mdbv1.DefaultDBUser(namespace.Name, "test-db-user", createdProject.Name).WithPasswordSecret(UserPasswordSecret)

				Expect(k8sClient.Create(context.Background(), createdDBUser)).To(Succeed())

				Eventually(testutil.WaitFor(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType)),
					DBUserUpdateTimeout, interval).Should(BeTrue())

				checkUserInAtlas(createdProject.ID(), *createdDBUser)
				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 0)
			})
			By("Creating cluster", func() {
				createdClusterAWS = mdbv1.DefaultAWSCluster(namespace.Name, createdProject.Name)
				Expect(k8sClient.Create(context.Background(), createdClusterAWS)).ToNot(HaveOccurred())

				// We don't wait for the full cluster creation - only when it has started the process
				Eventually(testutil.WaitFor(k8sClient, createdClusterAWS, status.FalseCondition(status.ClusterReadyType).WithReason(string(workflow.ClusterCreating))),
					20, intervalShort).Should(BeTrue())
			})
			By("Updating the database user while the cluster is being created", func() {
				createdDBUser = createdDBUser.WithRole("read", "test", "somecollection")
				Expect(k8sClient.Update(context.Background(), createdDBUser)).To(Succeed())

				// DatabaseUser will wait for the cluster to get created.
				Eventually(testutil.WaitFor(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType)),
					ClusterUpdateTimeout, interval).Should(BeTrue())

				expectedConditionsMatchers := testutil.MatchConditions(
					status.TrueCondition(status.DatabaseUserReadyType),
					status.TrueCondition(status.ReadyType),
				)
				Expect(createdDBUser.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))

				checkUserInAtlas(createdProject.ID(), *createdDBUser)
				Expect(tryConnect(createdProject.ID(), *createdClusterAWS, *createdDBUser)).Should(Succeed())
				By("Checking connection Secrets", func() {
					validateSecret(k8sClient, *createdProject, *createdClusterAWS, *createdDBUser)
					checkNumberOfConnectionSecrets(k8sClient, *createdProject, 1)
				})
			})
		})
	})
	Describe("Check the password Secret is watched", func() {
		It("Should succeed", func() {
			By("Creating clusters", func() {
				createdClusterAWS = mdbv1.DefaultAWSCluster(namespace.Name, createdProject.Name)
				Expect(k8sClient.Create(context.Background(), createdClusterAWS)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdClusterAWS, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
					ClusterUpdateTimeout, interval).Should(BeTrue())
			})
			createdDBUser = mdbv1.DefaultDBUser(namespace.Name, "test-db-user", createdProject.Name).WithPasswordSecret(UserPasswordSecret)
			var connSecretInitial corev1.Secret
			var pwdSecret corev1.Secret

			By(fmt.Sprintf("Creating the Database User %s", kube.ObjectKeyFromObject(createdDBUser)), func() {
				Expect(k8sClient.Create(context.Background(), createdDBUser)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType)),
					DBUserUpdateTimeout, interval, validateDatabaseUserUpdatingFunc()).Should(BeTrue())
				testutil.EventExists(k8sClient, createdDBUser, "Normal", "Ready", "")

				Expect(tryConnect(createdProject.ID(), *createdClusterAWS, *createdDBUser)).Should(Succeed())

				connSecretInitial = validateSecret(k8sClient, *createdProject, *createdClusterAWS, *createdDBUser)
				Expect(k8sClient.Get(context.Background(), kube.ObjectKey(namespace.Name, UserPasswordSecret), &pwdSecret)).To(Succeed())
				Expect(createdDBUser.Status.PasswordVersion).To(Equal(pwdSecret.ResourceVersion))
			})

			By("Breaking the password secret", func() {
				passwordSecret := buildPasswordSecret(UserPasswordSecret, "")
				Expect(k8sClient.Update(context.Background(), &passwordSecret)).To(Succeed())

				expectedCondition := status.FalseCondition(status.DatabaseUserReadyType).WithReason(string(workflow.Internal)).WithMessageRegexp("the 'password' field is empty")
				Eventually(testutil.WaitFor(k8sClient, createdDBUser, expectedCondition), 20, interval).Should(BeTrue())
				testutil.EventExists(k8sClient, createdDBUser, "Warning", string(workflow.Internal), "the 'password' field is empty")
			})
			By("Fixing the password secret", func() {
				passwordSecret := buildPasswordSecret(UserPasswordSecret, "someNewPassw00rd")
				Expect(k8sClient.Update(context.Background(), &passwordSecret)).To(Succeed())

				Eventually(testutil.WaitFor(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType)),
					DBUserUpdateTimeout, interval, validateDatabaseUserUpdatingFunc()).Should(BeTrue())

				// We need to make sure that the new connection secret is different from the initial one
				connSecretUpdated := validateSecret(k8sClient, *createdProject, *createdClusterAWS, *createdDBUser)
				Expect(string(connSecretInitial.Data["password"])).To(Equal(DBUserPassword))
				Expect(string(connSecretUpdated.Data["password"])).To(Equal("someNewPassw00rd"))

				var updatedPwdSecret corev1.Secret
				Expect(k8sClient.Get(context.Background(), kube.ObjectKey(namespace.Name, UserPasswordSecret), &updatedPwdSecret)).To(Succeed())
				Expect(updatedPwdSecret.ResourceVersion).NotTo(Equal(pwdSecret.ResourceVersion))
				Expect(createdDBUser.Status.PasswordVersion).To(Equal(updatedPwdSecret.ResourceVersion))

				Expect(tryConnect(createdProject.ID(), *createdClusterAWS, *createdDBUser)).Should(Succeed())
			})
		})
	})
	Describe("Change database users (make sure all stale secrets are removed)", func() {
		It("Should succeed", func() {
			By("Creating AWS and Azure clusters", func() {
				createdClusterAWS = mdbv1.DefaultAWSCluster(namespace.Name, createdProject.Name)
				Expect(k8sClient.Create(context.Background(), createdClusterAWS)).ToNot(HaveOccurred())

				createdClusterAzure = mdbv1.DefaultAzureCluster(namespace.Name, createdProject.Name)
				Expect(k8sClient.Create(context.Background(), createdClusterAzure)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdClusterAWS, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
					ClusterUpdateTimeout, interval).Should(BeTrue())

				Eventually(testutil.WaitFor(k8sClient, createdClusterAzure, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
					ClusterUpdateTimeout, interval).Should(BeTrue())
			})
			createdDBUser = mdbv1.DefaultDBUser(namespace.Name, "test-db-user", createdProject.Name).WithPasswordSecret(UserPasswordSecret)

			By(fmt.Sprintf("Creating the Database User %s (no scopes)", kube.ObjectKeyFromObject(createdDBUser)), func() {
				Expect(k8sClient.Create(context.Background(), createdDBUser)).To(Succeed())

				Eventually(testutil.WaitFor(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType)),
					DBUserUpdateTimeout, interval, validateDatabaseUserUpdatingFunc()).Should(BeTrue())

				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 2)

				s1 := validateSecret(k8sClient, *createdProject, *createdClusterAWS, *createdDBUser)
				s2 := validateSecret(k8sClient, *createdProject, *createdClusterAzure, *createdDBUser)

				testutil.EventExists(k8sClient, createdDBUser, "Normal", atlasdatabaseuser.ConnectionSecretsEnsuredEvent,
					fmt.Sprintf("Connection Secrets were created/updated: %s, %s", s1.Name, s2.Name))

			})
			By("Changing the db user name - two stale secret are expected to be removed, two added instead", func() {
				oldName := createdDBUser.Spec.Username
				createdDBUser = createdDBUser.WithAtlasUserName("new-user")
				Expect(k8sClient.Update(context.Background(), createdDBUser)).To(Succeed())

				Eventually(testutil.WaitFor(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType)),
					DBUserUpdateTimeout, interval, validateDatabaseUserUpdatingFunc()).Should(BeTrue())

				checkUserInAtlas(createdProject.ID(), *createdDBUser)
				// Old user has been removed
				_, _, err := atlasClient.DatabaseUsers.Get(context.Background(), createdDBUser.Spec.DatabaseName, createdProject.ID(), oldName)
				Expect(err).To(HaveOccurred())

				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 2)
				secret := validateSecret(k8sClient, *createdProject, *createdClusterAzure, *createdDBUser)
				Expect(secret.Name).To(Equal(connSecretname("-test-cluster-azure-new-user")))
				secret = validateSecret(k8sClient, *createdProject, *createdClusterAWS, *createdDBUser)
				Expect(secret.Name).To(Equal(connSecretname("-test-cluster-aws-new-user")))

				Expect(tryConnect(createdProject.ID(), *createdClusterAzure, *createdDBUser)).Should(Succeed())
				Expect(tryConnect(createdProject.ID(), *createdClusterAWS, *createdDBUser)).Should(Succeed())
			})
			By("Changing the scopes - one stale secret is expected to be removed", func() {
				createdDBUser = createdDBUser.ClearScopes().WithScope(mdbv1.ClusterScopeType, createdClusterAzure.Spec.Name)
				Expect(k8sClient.Update(context.Background(), createdDBUser)).To(Succeed())

				Eventually(testutil.WaitFor(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType)),
					DBUserUpdateTimeout, interval, validateDatabaseUserUpdatingFunc()).Should(BeTrue())

				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 1)
				validateSecret(k8sClient, *createdProject, *createdClusterAzure, *createdDBUser)

				Expect(tryConnect(createdProject.ID(), *createdClusterAzure, *createdDBUser)).Should(Succeed())
				Expect(tryConnect(createdProject.ID(), *createdClusterAWS, *createdDBUser)).ShouldNot(Succeed())
			})
		})
	})
	Describe("Check the user expiration", func() {
		It("Should succeed", func() {
			By("Creating a AWS cluster", func() {
				createdClusterAWS = mdbv1.DefaultAWSCluster(namespace.Name, createdProject.Name)
				Expect(k8sClient.Create(context.Background(), createdClusterAWS)).To(Succeed())

				Eventually(testutil.WaitFor(k8sClient, createdClusterAWS, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
					ClusterUpdateTimeout, interval).Should(BeTrue())
			})

			By("Creating the expired Database User - no user created in Atlas", func() {
				before := time.Now().UTC().Add(time.Minute * -10).Format("2006-01-02T15:04:05")
				createdDBUser = mdbv1.DefaultDBUser(namespace.Name, "test-db-user", createdProject.Name).
					WithPasswordSecret(UserPasswordSecret).
					WithDeleteAfterDate(before)

				Expect(k8sClient.Create(context.Background(), createdDBUser)).To(Succeed())

				Eventually(testutil.WaitFor(k8sClient, createdDBUser, status.FalseCondition(status.DatabaseUserReadyType).WithReason(string(workflow.DatabaseUserExpired))),
					15, intervalShort).Should(BeTrue())

				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 0)

				// no user in Atlas
				_, _, err := atlasClient.DatabaseUsers.Get(context.Background(), createdDBUser.Spec.DatabaseName, createdProject.ID(), createdDBUser.Spec.Username)
				Expect(err).To(HaveOccurred())
			})
			By("Fixing the Database User - setting the expiration to future", func() {
				after := time.Now().UTC().Add(time.Hour * 10).Format("2006-01-02T15:04:05")
				createdDBUser = createdDBUser.WithDeleteAfterDate(after)

				Expect(k8sClient.Update(context.Background(), createdDBUser)).To(Succeed())
				Eventually(testutil.WaitFor(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType)),
					DBUserUpdateTimeout, interval, validateDatabaseUserUpdatingFunc()).Should(BeTrue())

				checkUserInAtlas(createdProject.ID(), *createdDBUser)
				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 1)
				Expect(tryConnect(createdProject.ID(), *createdClusterAWS, *createdDBUser)).Should(Succeed())
			})
			By("Extending the expiration", func() {
				after := time.Now().UTC().Add(time.Hour * 30).Format("2006-01-02T15:04:05")
				createdDBUser = createdDBUser.WithDeleteAfterDate(after)

				Expect(k8sClient.Update(context.Background(), createdDBUser)).To(Succeed())
				Eventually(testutil.WaitFor(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType)),
					DBUserUpdateTimeout, interval, validateDatabaseUserUpdatingFunc()).Should(BeTrue())

				checkUserInAtlas(createdProject.ID(), *createdDBUser)
			})
			By("Emulating expiration of the User - connection secret must be removed", func() {
				before := time.Now().UTC().Add(time.Minute * -5).Format("2006-01-02T15:04:05")
				createdDBUser = createdDBUser.WithDeleteAfterDate(before)

				Expect(k8sClient.Update(context.Background(), createdDBUser)).To(Succeed())
				Eventually(testutil.WaitFor(k8sClient, createdDBUser, status.FalseCondition(status.DatabaseUserReadyType).WithReason(string(workflow.DatabaseUserExpired))),
					20, intervalShort).Should(BeTrue())

				expectedConditionsMatchers := testutil.MatchConditions(
					status.FalseCondition(status.DatabaseUserReadyType),
					status.FalseCondition(status.ReadyType),
				)
				Expect(createdDBUser.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))

				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 0)
			})
		})
		Describe("Deleting the db user (not cleaning Atlas)", func() {
			It("Should Succeed", func() {
				By(`Creating the db user with retention policy "keep" first`, func() {
					createdClusterAWS = mdbv1.DefaultAWSCluster(namespace.Name, createdProject.Name)
					Expect(k8sClient.Create(context.Background(), createdClusterAWS)).ToNot(HaveOccurred())

					Eventually(testutil.WaitFor(k8sClient, createdClusterAWS, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
						30*time.Minute, interval).Should(BeTrue())

					createdDBUser = mdbv1.DefaultDBUser(namespace.Name, "test-db-user", createdProject.Name).WithPasswordSecret(UserPasswordSecret)
					createdDBUser.ObjectMeta.Annotations = map[string]string{customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep}
					Expect(k8sClient.Create(context.Background(), createdDBUser)).To(Succeed())
					Eventually(testutil.WaitFor(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType)),
						DBUserUpdateTimeout, interval, validateDatabaseUserUpdatingFunc()).Should(BeTrue())
				})
				By("Deleting the db user - stays in Atlas", func() {
					Expect(k8sClient.Delete(context.Background(), createdDBUser)).To(Succeed())

					time.Sleep(1 * time.Minute)
					Expect(checkAtlasDatabaseUserRemoved(createdProject.Status.ID, *createdDBUser)()).Should(BeFalse())

					checkNumberOfConnectionSecrets(k8sClient, *createdProject, 0)
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

func tryConnect(projectID string, cluster mdbv1.AtlasCluster, user mdbv1.AtlasDatabaseUser) error {
	_, err := mongoClient(projectID, cluster, user)
	return err
}

func mongoClient(projectID string, cluster mdbv1.AtlasCluster, user mdbv1.AtlasDatabaseUser) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c, _, err := atlasClient.Clusters.Get(context.Background(), projectID, cluster.Spec.Name)
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

func tryWrite(projectID string, cluster mdbv1.AtlasCluster, user mdbv1.AtlasDatabaseUser, dbName, collectionName string) error {
	dbClient, err := mongoClient(projectID, cluster, user)
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
	fmt.Fprintf(GinkgoWriter, "User %s (cluster %s) has inserted a single document to %s/%s\n", user.Spec.Username, cluster.Spec.Name, dbName, collectionName)
	return nil
}

func validateSecret(k8sClient client.Client, project mdbv1.AtlasProject, cluster mdbv1.AtlasCluster, user mdbv1.AtlasDatabaseUser) corev1.Secret {
	secret := corev1.Secret{}
	username := user.Spec.Username
	secretName := fmt.Sprintf("%s-%s-%s", kube.NormalizeIdentifier(project.Spec.Name), kube.NormalizeIdentifier(cluster.Spec.Name), kube.NormalizeIdentifier(username))
	Expect(k8sClient.Get(context.Background(), kube.ObjectKey(project.Namespace, secretName), &secret)).To(Succeed())
	fmt.Printf("!! Secret: %v (%v)\n", kube.ObjectKey(project.Namespace, secretName), secret.Namespace+"/"+secret.Name)

	password, err := user.ReadPassword(k8sClient)
	Expect(err).NotTo(HaveOccurred())

	c, _, err := atlasClient.Clusters.Get(context.Background(), project.ID(), cluster.Spec.Name)
	Expect(err).NotTo(HaveOccurred())

	expectedData := map[string][]byte{
		"connectionStringStandard":    []byte(buildConnectionURL(c.ConnectionStrings.Standard, username, password)),
		"connectionStringStandardSrv": []byte(buildConnectionURL(c.ConnectionStrings.StandardSrv, username, password)),
		"username":                    []byte(username),
		"password":                    []byte(password),
	}
	expectedLabels := map[string]string{
		"atlas.mongodb.com/project-id":   project.ID(),
		"atlas.mongodb.com/cluster-name": cluster.Spec.Name,
	}
	Expect(secret.Data).To(Equal(expectedData))
	Expect(secret.Labels).To(Equal(expectedLabels))
	fmt.Printf("!! Secret 2: %v \n", secret.Namespace+"/"+secret.Name)
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

func validateDatabaseUserUpdatingFunc() func(a mdbv1.AtlasCustomResource) {
	return func(a mdbv1.AtlasCustomResource) {
		d := a.(*mdbv1.AtlasDatabaseUser)
		expectedConditionsMatchers := testutil.MatchConditions(
			status.FalseCondition(status.DatabaseUserReadyType).WithReason(string(workflow.DatabaseUserClustersAppliedChanges)),
			status.FalseCondition(status.ReadyType),
		)
		Expect(d.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
	}
}

//nolint
func validateDatabaseUserWaitingForCluster() func(a mdbv1.AtlasCustomResource) {
	return func(a mdbv1.AtlasCustomResource) {
		d := a.(*mdbv1.AtlasDatabaseUser)
		// this is the first status that db user gets after update
		userChangesApplied := testutil.MatchConditions(
			status.FalseCondition(status.DatabaseUserReadyType).WithReason(string(workflow.DatabaseUserClustersAppliedChanges)),
			status.FalseCondition(status.ReadyType),
		)
		// this is the status the db user gets to when tries to create connection secrets and sees that the cluster
		// is not ready
		waitingForCluster := testutil.MatchConditions(
			status.FalseCondition(status.DatabaseUserReadyType).
				WithReason(string(workflow.DatabaseUserConnectionSecretsNotCreated)).
				WithMessageRegexp("Waiting for clusters to get created/updated"),
			status.FalseCondition(status.ReadyType),
		)
		Expect(d.Status.Conditions).To(Or(ConsistOf(waitingForCluster), ConsistOf(userChangesApplied)))
	}
}
