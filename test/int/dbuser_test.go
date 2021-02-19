package int

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/testutil"
)

const UserPasswordSecret = "user-password-secret"

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
			createdProject = mdbv1.DefaultProject(namespace.Name, connectionSecret.Name)

			Expect(k8sClient.Create(context.Background(), createdProject)).To(Succeed())
			Eventually(testutil.WaitFor(k8sClient, createdProject, status.TrueCondition(status.ReadyType)),
				20, interval).Should(BeTrue())
		})

		By("Creating cluster", func() {
			createdClusterAWS = mdbv1.DefaultAWSCluster(namespace.Name, createdProject.Name)
			Expect(k8sClient.Create(context.Background(), createdClusterAWS)).ToNot(HaveOccurred())

			createdClusterGCP = mdbv1.DefaultGCPCluster(namespace.Name, createdProject.Name)
			Expect(k8sClient.Create(context.Background(), createdClusterAWS)).ToNot(HaveOccurred())

			Eventually(testutil.WaitFor(k8sClient, createdClusterAWS, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
				1800, interval).Should(BeTrue())

			Eventually(testutil.WaitFor(k8sClient, createdClusterGCP, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
				500, interval).Should(BeTrue())
		})
	})

	AfterEach(func() {
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

	Describe("Create/Update the db user", func() {
		It("Should Succeed", func() {
			createdDBUser = mdbv1.DefaultDBUser(namespace.Name, "test-db-user", createdProject.Name).WithPasswordSecret(UserPasswordSecret)

			By(fmt.Sprintf("Creating the Database User %s", kube.ObjectKeyFromObject(createdDBUser)), func() {
				Expect(k8sClient.Create(context.Background(), createdDBUser)).ToNot(HaveOccurred())

				Eventually(testutil.WaitFor(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType)),
					20, interval).Should(BeTrue())
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
