package int

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas/mongodbatlas"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/testutil"
)

const UserPasswordSecret = "user-password-secret"
const DevMode = true

var _ = FDescribe("AtlasDatabaseUser", func() {
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
			createdProject = mdbv1.DefaultProject(namespace.Name, connectionSecret.Name)
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
			return
		}
		if createdProject != nil && createdProject.Status.ID != "" {
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
		It("Should Succeed", func() {
			createdDBUser = mdbv1.DefaultDBUser(namespace.Name, "test-db-user", createdProject.Name).WithPasswordSecret(UserPasswordSecret)

			By(fmt.Sprintf("Creating the Database User %s", kube.ObjectKeyFromObject(createdDBUser)), func() {
				Expect(k8sClient.Create(context.Background(), createdDBUser)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType)),
					20, interval).Should(BeTrue())

				checkUserInAtlas()
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
		StringData: map[string]string{"password": "Passw0rd!"},
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
