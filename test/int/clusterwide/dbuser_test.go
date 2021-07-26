package int

import (
	"context"
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/testutil"
)

const (
	DevMode            = false
	UserPasswordSecret = "user-password-secret"
	DBUserPassword     = "Passw0rd!"
	// M2 clusters take longer time to apply changes
	DBUserUpdateTimeout    = 170
	ProjectCreationTimeout = 40
)

var _ = Describe("ClusterWide", func() {
	const interval = time.Second * 1

	var (
		connectionSecret  corev1.Secret
		createdProject    *mdbv1.AtlasProject
		createdClusterAWS *mdbv1.AtlasCluster
		createdDBUser     *mdbv1.AtlasDatabaseUser
		secondDBUser      *mdbv1.AtlasDatabaseUser
	)

	BeforeEach(func() {
		namespace = corev1.Namespace{ObjectMeta: metav1.ObjectMeta{GenerateName: "test"}}
		Expect(k8sClient.Create(context.Background(), &namespace)).ToNot(HaveOccurred())

		createdDBUser = &mdbv1.AtlasDatabaseUser{}

		connectionSecret = buildConnectionSecret("my-atlas-key")
		Expect(k8sClient.Create(context.Background(), &connectionSecret)).To(Succeed())

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
	})

	Describe("Create user and cluster in different namespaces", func() {
		It("Should Succeed", func() {
			clusterNS := corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace.Name + "-other-cluster"}}
			Expect(k8sClient.Create(context.Background(), &clusterNS)).ToNot(HaveOccurred())

			userNS := corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace.Name + "-other-user"}}
			Expect(k8sClient.Create(context.Background(), &userNS)).ToNot(HaveOccurred())

			By(fmt.Sprintf("Creating password Secret %s", UserPasswordSecret))
			passwordSecret := buildPasswordSecret(userNS.Name, UserPasswordSecret, DBUserPassword)
			Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())

			createdClusterAWS = mdbv1.DefaultAWSCluster(clusterNS.Name, createdProject.Name)
			// The project namespace is different from the cluster one - need to specify explicitly
			createdClusterAWS.Spec.Project.Namespace = namespace.Name

			Expect(k8sClient.Create(context.Background(), createdClusterAWS)).ToNot(HaveOccurred())

			Eventually(testutil.WaitFor(k8sClient, createdClusterAWS, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
				30*time.Minute, interval).Should(BeTrue())

			createdDBUser = mdbv1.DefaultDBUser(userNS.Name, "test-db-user", createdProject.Name).WithPasswordSecret(UserPasswordSecret)
			createdDBUser.Spec.Project.Namespace = namespace.Name
			Expect(k8sClient.Create(context.Background(), createdDBUser)).To(Succeed())
			Eventually(testutil.WaitFor(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType)),
				DBUserUpdateTimeout, interval, validateDatabaseUserUpdatingFunc()).Should(BeTrue())

			By("Removing the cluster", func() {
				Expect(k8sClient.Delete(context.Background(), createdClusterAWS)).To(Succeed())
				Eventually(checkAtlasClusterRemoved(createdProject.ID(), createdClusterAWS.Spec.Name), 600, interval).Should(BeTrue())
			})
		})
	})
})

func buildConnectionSecret(name string) corev1.Secret {
	return corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace.Name,
		},
		StringData: map[string]string{"orgId": connection.OrgID, "publicApiKey": connection.PublicKey, "privateApiKey": connection.PrivateKey},
	}
}

func buildPasswordSecret(namespace, name, password string) corev1.Secret {
	return corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		StringData: map[string]string{"password": password},
	}
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

func checkAtlasClusterRemoved(projectID string, clusterName string) func() bool {
	return func() bool {
		_, r, err := atlasClient.Clusters.Get(context.Background(), projectID, clusterName)
		if err != nil {
			if r != nil && r.StatusCode == http.StatusNotFound {
				return true
			}
		}

		return false
	}
}

func checkAtlasProjectRemoved(projectID string) func() bool {
	return func() bool {
		_, r, err := atlasClient.Projects.GetOneProject(context.Background(), projectID)
		if err != nil {
			if r != nil && r.StatusCode == http.StatusNotFound {
				return true
			}
		}
		return false
	}
}

func validateClusterCreatingFunc() func(a mdbv1.AtlasCustomResource) {
	startedCreation := false
	return func(a mdbv1.AtlasCustomResource) {
		c := a.(*mdbv1.AtlasCluster)
		if c.Status.StateName != "" {
			startedCreation = true
		}
		// When the create request has been made to Atlas - we expect the following status
		if startedCreation {
			Expect(c.Status.StateName).To(Equal("CREATING"), fmt.Sprintf("Current conditions: %+v", c.Status.Conditions))
			expectedConditionsMatchers := testutil.MatchConditions(
				status.FalseCondition(status.ClusterReadyType).WithReason(string(workflow.ClusterCreating)).WithMessageRegexp("cluster is provisioning"),
				status.FalseCondition(status.ReadyType),
			)
			Expect(c.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
		} else {
			// Otherwise there could have been some exception in Atlas on creation - let's check the conditions
			condition, ok := testutil.FindConditionByType(c.Status.Conditions, status.ClusterReadyType)
			Expect(ok).To(BeFalse(), fmt.Sprintf("Unexpected condition: %v", condition))
		}
	}
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
