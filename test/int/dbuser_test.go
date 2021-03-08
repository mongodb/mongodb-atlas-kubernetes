package int

import (
	"context"
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

	"go.mongodb.org/mongo-driver/mongo"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/testutil"
)

const DevMode = false

const UserPasswordSecret = "user-password-secret"
const DBUserPassword = "Passw0rd!"
const UserPasswordSecret2 = "second-user-password-secret"
const DBUserPassword2 = "H@lla#!"

var _ = Describe("AtlasDatabaseUser", func() {
	const interval = time.Second * 1

	var (
		connectionSecret  corev1.Secret
		createdProject    *mdbv1.AtlasProject
		createdClusterAWS *mdbv1.AtlasCluster
		createdClusterGCP *mdbv1.AtlasCluster
		createdDBUser     *mdbv1.AtlasDatabaseUser
	)

	BeforeEach(func() {
		prepareControllers()
		createdDBUser = &mdbv1.AtlasDatabaseUser{}

		connectionSecret = buildConnectionSecret("my-atlas-key")
		Expect(k8sClient.Create(context.Background(), &connectionSecret)).To(Succeed())

		passwordSecret := buildPasswordSecret(UserPasswordSecret, DBUserPassword)
		Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())

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
				20, interval).Should(BeTrue())
		})

		By("Creating cluster", func() {
			createdClusterAWS = mdbv1.DefaultAWSCluster(namespace.Name, createdProject.Name)
			Expect(k8sClient.Create(context.Background(), createdClusterAWS)).ToNot(HaveOccurred())

			createdClusterGCP = mdbv1.DefaultGCPCluster(namespace.Name, createdProject.Name)
			Expect(k8sClient.Create(context.Background(), createdClusterGCP)).ToNot(HaveOccurred())

			Eventually(testutil.WaitFor(k8sClient, createdClusterAWS, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
				1800, interval).Should(BeTrue())

			Eventually(testutil.WaitFor(k8sClient, createdClusterGCP, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
				500, interval).Should(BeTrue())
		})
	})

	AfterEach(func() {
		if DevMode {
			// No tearDown in dev mode - projects and both clusters will stay in Atlas so it's easier to develop
			// tests. Just rerun the test and the project + clusters in Atlas will be reused.
			// We only need to wipe data in the databases.
			dbClient, err := mongoClient(createdProject.ID(), *createdClusterAWS, *createdDBUser)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbClient.Database("test").Collection("operatortest").Drop(context.Background())).To(Succeed())

			dbClient, err = mongoClient(createdProject.ID(), *createdClusterGCP, *createdDBUser)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbClient.Database("test").Collection("operatortest").Drop(context.Background())).To(Succeed())

			return
		}
		if createdProject != nil && createdProject.ID() != "" {
			if createdClusterGCP != nil {
				By("Removing Atlas Cluster " + createdClusterGCP.Name)
				Expect(k8sClient.Delete(context.Background(), createdClusterGCP)).To(Succeed())
				Eventually(checkAtlasClusterRemoved(createdProject.Status.ID, createdClusterGCP.Name), 600, interval).Should(BeTrue())
			}
			if createdClusterAWS != nil {
				By("Removing Atlas Cluster " + createdClusterAWS.Name)
				Expect(k8sClient.Delete(context.Background(), createdClusterAWS)).To(Succeed())
				Eventually(checkAtlasClusterRemoved(createdProject.Status.ID, createdClusterAWS.Name), 600, interval).Should(BeTrue())
			}

			By("Removing Atlas Project " + createdProject.Status.ID)
			Eventually(removeAtlasProject(createdProject.Status.ID), 600, interval).Should(BeTrue())
		}
		removeControllersAndNamespace()
	})

	checkUserInAtlas := func(user mdbv1.AtlasDatabaseUser) {
		By("Verifying Database User state in Atlas", func() {
			atlasDBUser, _, err := atlasClient.DatabaseUsers.Get(context.Background(), user.Spec.DatabaseName, createdProject.ID(), user.Spec.Username)
			Expect(err).ToNot(HaveOccurred())
			operatorDBUser, err := user.ToAtlas(k8sClient)
			Expect(err).ToNot(HaveOccurred())

			Expect(*atlasDBUser).To(Equal(normalize(*operatorDBUser, createdProject.ID())))
		})
	}

	connSecretname := func(suffix string) string {
		return kube.NormalizeIdentifier(createdProject.Spec.Name) + suffix
	}

	Describe("Create/Update the db user", func() {
		It("Should be created successfully", func() {
			createdDBUser = mdbv1.DefaultDBUser(namespace.Name, "test-db-user", createdProject.Name).WithPasswordSecret(UserPasswordSecret)

			By(fmt.Sprintf("Creating the Database User %s", kube.ObjectKeyFromObject(createdDBUser)), func() {
				Expect(k8sClient.Create(context.Background(), createdDBUser)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType)),
					20, interval).Should(BeTrue())

				checkUserInAtlas(*createdDBUser)

				// TODO CLOUDP-83026 and CLOUDP-83098 remove Eventually in favor of Expect
				Eventually(tryConnect(createdProject.ID(), *createdClusterGCP, *createdDBUser), 90, interval).Should(Succeed())
				Eventually(tryConnect(createdProject.ID(), *createdClusterAWS, *createdDBUser), 90, interval).Should(Succeed())
				By("Checking connection Secrets", func() {
					validateSecret(k8sClient, *createdProject, *createdClusterGCP, *createdDBUser)
					validateSecret(k8sClient, *createdProject, *createdClusterAWS, *createdDBUser)
					checkNumberOfConnectionSecrets(k8sClient, *createdProject, 2)

					expectedSecretsInStatus := map[string]string{
						"test-cluster-aws": "dev-test-atlas-project-test-cluster-aws-test-db-user",
						"test-cluster-gcp": "dev-test-atlas-project-test-cluster-gcp-test-db-user",
					}
					Expect(createdDBUser.Status.ConnectionSecrets).To(Equal(expectedSecretsInStatus))
				})
				By("Checking connectivity to Clusters", func() {
					// The user created lacks read/write roles
					err := tryWrite(createdProject.ID(), *createdClusterGCP, *createdDBUser, "test", "operatortest")()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(MatchRegexp("not authorized"))

					err = tryWrite(createdProject.ID(), *createdClusterAWS, *createdDBUser, "test", "operatortest")()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(MatchRegexp("not authorized"))
				})
			})
			By("Update database user - give readWrite permissions", func() {
				// Adding the role allowing read/write
				createdDBUser = createdDBUser.WithRole("readWriteAnyDatabase", "admin", "")

				Expect(k8sClient.Update(context.Background(), createdDBUser)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType)),
					20, interval).Should(BeTrue())

				checkUserInAtlas(*createdDBUser)

				By("Checking connection Secrets", func() {
					validateSecret(k8sClient, *createdProject, *createdClusterGCP, *createdDBUser)
					validateSecret(k8sClient, *createdProject, *createdClusterAWS, *createdDBUser)
					checkNumberOfConnectionSecrets(k8sClient, *createdProject, 2)

					expectedSecretsInStatus := map[string]string{
						"test-cluster-aws": connSecretname("-test-cluster-aws-test-db-user"),
						"test-cluster-gcp": connSecretname("-test-cluster-gcp-test-db-user"),
					}
					Expect(createdDBUser.Status.ConnectionSecrets).To(Equal(expectedSecretsInStatus))
				})

				By("Checking write permissions for Clusters", func() {
					// TODO CLOUDP-83026 remove Eventually in favor of Expect
					Eventually(tryWrite(createdProject.ID(), *createdClusterGCP, *createdDBUser, "test", "operatortest"), 60, interval).Should(Succeed())
					Eventually(tryWrite(createdProject.ID(), *createdClusterAWS, *createdDBUser, "test", "operatortest"), 60, interval).Should(Succeed())
				})
			})
			By("Adding another user for GCP cluster only", func() {
				passwordSecret := buildPasswordSecret(UserPasswordSecret2, DBUserPassword2)
				Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())
				secondDBUser := mdbv1.DefaultDBUser(namespace.Name, "second-db-user", createdProject.Name).
					WithPasswordSecret(UserPasswordSecret2).
					WithRole("readWrite", "someDB", "thisIsTheOnlyAllowedCollection").
					WithScope(mdbv1.ClusterScopeType, createdClusterGCP.Spec.Name)

				Expect(k8sClient.Create(context.Background(), secondDBUser)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, secondDBUser, status.TrueCondition(status.ReadyType)),
					20, interval).Should(BeTrue())

				checkUserInAtlas(*secondDBUser)
				By("Checking connection Secrets", func() {
					validateSecret(k8sClient, *createdProject, *createdClusterGCP, *createdDBUser)
					validateSecret(k8sClient, *createdProject, *createdClusterAWS, *createdDBUser)
					validateSecret(k8sClient, *createdProject, *createdClusterGCP, *secondDBUser)
					checkNumberOfConnectionSecrets(k8sClient, *createdProject, 3)
					expectedSecretsInStatus := map[string]string{"test-cluster-gcp": connSecretname("-test-cluster-gcp-second-db-user")}
					Expect(secondDBUser.Status.ConnectionSecrets).To(Equal(expectedSecretsInStatus))
				})

				By("Checking write permissions for Clusters", func() {
					// We still can write by the first user
					Expect(tryWrite(createdProject.ID(), *createdClusterGCP, *createdDBUser, "test", "testCollection")()).Should(Succeed())
					Expect(tryWrite(createdProject.ID(), *createdClusterAWS, *createdDBUser, "test", "testCollection")()).Should(Succeed())

					// The second user can eventually write to one collection only
					Eventually(tryConnect(createdProject.ID(), *createdClusterGCP, *secondDBUser), 90, interval).Should(Succeed())
					Eventually(tryWrite(createdProject.ID(), *createdClusterGCP, *secondDBUser, "someDB", "thisIsTheOnlyAllowedCollection"), 60, interval).Should(Succeed())

					err := tryWrite(createdProject.ID(), *createdClusterGCP, *secondDBUser, "test", "someNotAllowedCollection")()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(MatchRegexp("not authorized"))
				})
				By("Removing Second user", func() {
					Expect(k8sClient.Delete(context.Background(), secondDBUser)).To(Succeed())
					Eventually(checkAtlasDatabaseUserRemoved(createdProject.Status.ID, *secondDBUser), 50, interval).Should(BeTrue())

					secretNames := []string{connSecretname("-test-cluster-gcp-second-db-user")}
					Eventually(checkSecretsDontExist(namespace.Name, secretNames), 50, interval).Should(BeTrue())
				})
			})
			By("Removing First user", func() {
				Expect(k8sClient.Delete(context.Background(), createdDBUser)).To(Succeed())
				Eventually(checkAtlasDatabaseUserRemoved(createdProject.Status.ID, *createdDBUser), 50, interval).Should(BeTrue())

				secretNames := []string{connSecretname("-test-cluster-aws-test-db-user"), connSecretname("-test-cluster-gcp-test-db-user")}
				Eventually(checkSecretsDontExist(namespace.Name, secretNames), 50, interval).Should(BeTrue())

				checkNumberOfConnectionSecrets(k8sClient, *createdProject, 0)
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
	user.GroupID = projectID
	user.Password = ""
	return user
}
func tryConnect(projectID string, cluster mdbv1.AtlasCluster, user mdbv1.AtlasDatabaseUser) func() error {
	return func() error {
		_, err := mongoClient(projectID, cluster, user)
		return err
	}
}
func mongoClient(projectID string, cluster mdbv1.AtlasCluster, user mdbv1.AtlasDatabaseUser) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c, _, err := atlasClient.Clusters.Get(context.Background(), projectID, cluster.Spec.Name)
	Expect(err).NotTo(HaveOccurred())

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

func tryWrite(projectID string, cluster mdbv1.AtlasCluster, user mdbv1.AtlasDatabaseUser, dbName, collectionName string) func() error {
	return func() error {
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
}

func validateSecret(k8sClient client.Client, project mdbv1.AtlasProject, cluster mdbv1.AtlasCluster, user mdbv1.AtlasDatabaseUser) {
	secret := corev1.Secret{}
	username := user.Spec.Username
	secretName := fmt.Sprintf("%s-%s-%s", kube.NormalizeIdentifier(project.Spec.Name), kube.NormalizeIdentifier(cluster.Spec.Name), kube.NormalizeIdentifier(username))
	Expect(k8sClient.Get(context.Background(), kube.ObjectKey(project.Namespace, secretName), &secret)).To(Succeed())

	password, err := user.ReadPassword(k8sClient)
	Expect(err).NotTo(HaveOccurred())

	expectedData := map[string][]byte{
		"connectionString.standard":    []byte(buildConnectionURL(cluster.Status.ConnectionStrings.Standard, username, password)),
		"connectionString.standardSrv": []byte(buildConnectionURL(cluster.Status.ConnectionStrings.StandardSrv, username, password)),
		"username":                     []byte(username),
		"password":                     []byte(password),
	}
	expectedLabels := map[string]string{
		"atlas.mongodb.com/project-id":   project.ID(),
		"atlas.mongodb.com/cluster-name": cluster.Spec.Name,
	}
	Expect(secret.Data).To(Equal(expectedData))
	Expect(secret.Labels).To(Equal(expectedLabels))
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
