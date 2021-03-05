package int

import (
	"context"
	"fmt"
	"net/url"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"go.mongodb.org/mongo-driver/mongo"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/testutil"
)

const (
	UserPasswordSecret = "user-password-secret"
	DevMode            = false
	DBUserPassword     = "Passw0rd!"
)

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

		passwordSecret := buildPasswordSecret(UserPasswordSecret)
		Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())

		By("Creating the project", func() {
			createdProject = mdbv1.DefaultProject(namespace.Name, connectionSecret.Name).
				WithIPAccessList(project.NewIPAccessList().WithIP("0.0.0.0/0"))
			if DevMode {
				// While developing tests we need to reuse the same project
				createdProject.Spec.Name = "dev-test-atlas-project"
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
			Expect(k8sClient.Delete(context.Background(), createdProject)).To(Succeed())
			Eventually(checkAtlasProjectRemoved(createdProject.Status.ID), 20, interval).Should(BeTrue())
		}
		removeControllersAndNamespace()
	})

	checkUserInAtlas := func() {
		By("Verifying Database User state in Atlas", func() {
			atlasDBUser, _, err := atlasClient.DatabaseUsers.Get(context.Background(), createdDBUser.Spec.DatabaseName, createdProject.ID(), createdDBUser.Spec.Username)
			Expect(err).ToNot(HaveOccurred())
			operatorDBUser, err := createdDBUser.ToAtlas(k8sClient)
			Expect(err).ToNot(HaveOccurred())

			Expect(*atlasDBUser).To(Equal(normalize(*operatorDBUser, createdProject.ID())))
		})
	}

	Describe("Create/Update the db user", func() {
		It("Should be created successfully", func() {
			createdDBUser = mdbv1.DefaultDBUser(namespace.Name, "test-db-user", createdProject.Name).WithPasswordSecret(UserPasswordSecret)

			By(fmt.Sprintf("Creating the Database User %s", kube.ObjectKeyFromObject(createdDBUser)), func() {
				Expect(k8sClient.Create(context.Background(), createdDBUser)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType)),
					20, interval).Should(BeTrue())

				checkUserInAtlas()

				By("Checking connectivity to Clusters", func() {
					// TODO CLOUDP-83026 and CLOUDP-83098 remove Eventually in favor of Expect
					Eventually(tryConnect(createdProject.ID(), *createdClusterGCP, *createdDBUser), 90, interval).Should(Succeed())
					Eventually(tryConnect(createdProject.ID(), *createdClusterAWS, *createdDBUser), 90, interval).Should(Succeed())

					// The user created lacks read/write roles
					err := tryWrite(createdProject.ID(), *createdClusterGCP, *createdDBUser)()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(MatchRegexp("not authorized"))

					err = tryWrite(createdProject.ID(), *createdClusterAWS, *createdDBUser)()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(MatchRegexp("not authorized"))
				})
			})
			By("Should get readWrite permissions", func() {
				// Adding the role allowing read/write
				createdDBUser = createdDBUser.WithRole("readWriteAnyDatabase", "admin", "")

				Expect(k8sClient.Update(context.Background(), createdDBUser)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType)),
					20, interval).Should(BeTrue())

				checkUserInAtlas()

				By("Checking write permissions for Clusters", func() {
					// TODO CLOUDP-83026 remove Eventually in favor of Expect
					Eventually(tryWrite(createdProject.ID(), *createdClusterGCP, *createdDBUser), 60, interval).Should(Succeed())
					Eventually(tryWrite(createdProject.ID(), *createdClusterAWS, *createdDBUser), 60, interval).Should(Succeed())
				})
			})
		})
	})
})

func buildPasswordSecret(name string) corev1.Secret {
	return corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace.Name,
		},
		StringData: map[string]string{"password": DBUserPassword},
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
	cs.User = url.UserPassword(user.Spec.Username, DBUserPassword)

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

func tryWrite(projectID string, cluster mdbv1.AtlasCluster, user mdbv1.AtlasDatabaseUser) func() error {
	return func() error {
		dbClient, err := mongoClient(projectID, cluster, user)
		Expect(err).NotTo(HaveOccurred())
		defer func() {
			if err = dbClient.Disconnect(context.Background()); err != nil {
				panic(err)
			}
		}()

		collection := dbClient.Database("test").Collection("operatortest")

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
		fmt.Println("Inserted a single document")
		return nil
	}
}
