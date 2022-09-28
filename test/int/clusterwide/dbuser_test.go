package int

import (
	"context"
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
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
	// M2 Deployments take longer time to apply changes
	DBUserUpdateTimeout    = 170 * time.Second
	ProjectCreationTimeout = 40 * time.Second
)

var _ = Describe("clusterwide", Label("int", "clusterwide"), func() {
	const interval = time.Second * 1

	var (
		connectionSecret     corev1.Secret
		createdProject       *mdbv1.AtlasProject
		createdDeploymentAWS *mdbv1.AtlasDeployment
		createdDBUser        *mdbv1.AtlasDatabaseUser
		secondDBUser         *mdbv1.AtlasDatabaseUser
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
				Eventually(checkAtlasDeploymentRemoved(createdProject.ID(), list.Items[i].Spec.DeploymentSpec.Name), 600, interval).Should(BeTrue())
			}

			By("Removing Atlas Project " + createdProject.Status.ID)
			Expect(k8sClient.Delete(context.Background(), createdProject)).To(Succeed())
			Eventually(checkAtlasProjectRemoved(createdProject.Status.ID), 60, interval).Should(BeTrue())
		}
	})

	Describe("Create user and deployment in different namespaces", func() {
		It("Should Succeed", func() {
			deploymentNS := corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace.Name + "-other-deployment"}}
			Expect(k8sClient.Create(context.Background(), &deploymentNS)).ToNot(HaveOccurred())

			userNS := corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace.Name + "-other-user"}}
			Expect(k8sClient.Create(context.Background(), &userNS)).ToNot(HaveOccurred())

			By(fmt.Sprintf("Creating password Secret %s", UserPasswordSecret))
			passwordSecret := buildPasswordSecret(userNS.Name, UserPasswordSecret, DBUserPassword)
			Expect(k8sClient.Create(context.Background(), &passwordSecret)).To(Succeed())

			createdDeploymentAWS = mdbv1.DefaultAWSDeployment(deploymentNS.Name, createdProject.Name).Lightweight()
			// The project namespace is different from the deployment one - need to specify explicitly
			createdDeploymentAWS.Spec.Project.Namespace = namespace.Name

			Expect(k8sClient.Create(context.Background(), createdDeploymentAWS)).ToNot(HaveOccurred())

			Eventually(func(g Gomega) bool {
				return testutil.CheckCondition(k8sClient, createdDeploymentAWS, status.TrueCondition(status.ReadyType), validateDeploymentCreatingFunc(g))
			}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

			createdDBUser = mdbv1.DefaultDBUser(userNS.Name, "test-db-user", createdProject.Name).WithPasswordSecret(UserPasswordSecret)
			createdDBUser.Spec.Project.Namespace = namespace.Name
			Expect(k8sClient.Create(context.Background(), createdDBUser)).To(Succeed())
			Eventually(func() bool {
				return testutil.CheckCondition(k8sClient, createdDBUser, status.TrueCondition(status.ReadyType))
			}).WithTimeout(DBUserUpdateTimeout).WithPolling(interval).Should(BeTrue())

			By("Removing the deployment", func() {
				Expect(k8sClient.Delete(context.Background(), createdDeploymentAWS)).To(Succeed())
				Eventually(checkAtlasDeploymentRemoved(createdProject.ID(), createdDeploymentAWS.Spec.DeploymentSpec.Name), 600, interval).Should(BeTrue())
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

func checkAtlasDeploymentRemoved(projectID string, deploymentName string) func() bool {
	return func() bool {
		_, r, err := atlasClient.Clusters.Get(context.Background(), projectID, deploymentName)
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

func validateDeploymentCreatingFunc(g Gomega) func(a mdbv1.AtlasCustomResource) {
	startedCreation := false
	return func(a mdbv1.AtlasCustomResource) {
		c := a.(*mdbv1.AtlasDeployment)
		if c.Status.StateName != "" {
			startedCreation = true
		}
		// When the create request has been made to Atlas - we expect the following status
		if startedCreation {
			g.Expect(c.Status.StateName).To(Equal("CREATING"), fmt.Sprintf("Current conditions: %+v", c.Status.Conditions))
			expectedConditionsMatchers := testutil.MatchConditions(
				status.FalseCondition(status.DeploymentReadyType).WithReason(string(workflow.DeploymentCreating)).WithMessageRegexp("deployment is provisioning"),
				status.FalseCondition(status.ReadyType),
				status.TrueCondition(status.ValidationSucceeded),
			)
			g.Expect(c.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
		} else {
			// Otherwise there could have been some exception in Atlas on creation - let's check the conditions
			condition, ok := testutil.FindConditionByType(c.Status.Conditions, status.DeploymentReadyType)
			g.Expect(ok).To(BeFalse(), fmt.Sprintf("Unexpected condition: %v", condition))
		}
	}
}
