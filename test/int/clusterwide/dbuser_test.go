// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/secretservice"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/httputil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/conditions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/resources"
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
		createdProject       *akov2.AtlasProject
		createdDeploymentAWS *akov2.AtlasDeployment
		createdDBUser        *akov2.AtlasDatabaseUser
		secondDBUser         *akov2.AtlasDatabaseUser
	)

	BeforeEach(func() {
		namespace = corev1.Namespace{ObjectMeta: metav1.ObjectMeta{GenerateName: "test"}}
		Expect(k8sClient.Create(context.Background(), &namespace)).ToNot(HaveOccurred())

		createdDBUser = &akov2.AtlasDatabaseUser{}

		connectionSecret = buildConnectionSecret("my-atlas-key")
		Expect(k8sClient.Create(context.Background(), &connectionSecret)).To(Succeed())

		By("Creating the project", func() {
			// adding whitespace to the name to check normalization for connection secrets names
			createdProject = akov2.DefaultProject(namespace.Name, connectionSecret.Name).
				WithAtlasName(namespace.Name + " some").
				WithIPAccessList(project.NewIPAccessList().WithCIDR("0.0.0.0/0"))
			if DevMode {
				// While developing tests we need to reuse the same project
				createdProject.Spec.Name = "dev-test atlas-project"
			}
			Expect(k8sClient.Create(context.Background(), createdProject)).To(Succeed())
			Eventually(func() bool {
				return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType))
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
			list := akov2.AtlasDeploymentList{}
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

			createdDeploymentAWS = akov2.DefaultAWSDeployment(deploymentNS.Name, createdProject.Name).Lightweight()
			// The project namespace is different from the deployment one - need to specify explicitly
			createdDeploymentAWS.Spec.ProjectRef.Namespace = namespace.Name

			Expect(k8sClient.Create(context.Background(), createdDeploymentAWS)).ToNot(HaveOccurred())

			Eventually(func(g Gomega) bool {
				return resources.CheckCondition(k8sClient, createdDeploymentAWS, api.TrueCondition(api.ReadyType), validateDeploymentCreatingFunc(g))
			}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())

			createdDBUser = akov2.DefaultDBUser(userNS.Name, "test-db-user", createdProject.Name).WithPasswordSecret(UserPasswordSecret)
			createdDBUser.Spec.ProjectRef.Namespace = namespace.Name
			Expect(k8sClient.Create(context.Background(), createdDBUser)).To(Succeed())
			Eventually(func() bool {
				return resources.CheckCondition(k8sClient, createdDBUser, api.TrueCondition(api.ReadyType))
			}).WithTimeout(DBUserUpdateTimeout).WithPolling(interval).Should(BeTrue())

			By("Removing the deployment", func() {
				Expect(k8sClient.Delete(context.Background(), createdDeploymentAWS)).To(Succeed())
				Eventually(checkAtlasDeploymentRemoved(createdProject.ID(), createdDeploymentAWS.GetDeploymentName()), 600, interval).Should(BeTrue())
			})
		})
	})
})

func buildConnectionSecret(name string) corev1.Secret {
	return corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace.Name,
			Labels: map[string]string{
				secretservice.TypeLabelKey: secretservice.CredLabelVal,
			},
		},
		StringData: map[string]string{"orgId": orgID, "publicApiKey": publicKey, "privateApiKey": privateKey},
	}
}

func buildPasswordSecret(namespace, name, password string) corev1.Secret {
	return corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				secretservice.TypeLabelKey: secretservice.CredLabelVal,
			},
		},
		StringData: map[string]string{"password": password},
	}
}

func checkAtlasDatabaseUserRemoved(projectID string, user akov2.AtlasDatabaseUser) func() bool {
	return func() bool {
		_, r, err := atlasClient.DatabaseUsersApi.
			GetDatabaseUser(context.Background(), user.Spec.DatabaseName, projectID, user.Spec.Username).
			Execute()
		if err != nil {
			if httputil.StatusCode(r) == http.StatusNotFound {
				return true
			}
		}

		return false
	}
}

func checkAtlasDeploymentRemoved(projectID string, deploymentName string) func() bool {
	return func() bool {
		_, r, err := atlasClient.ClustersApi.
			GetCluster(context.Background(), projectID, deploymentName).
			Execute()
		if err != nil {
			if httputil.StatusCode(r) == http.StatusNotFound {
				return true
			}
		}

		return false
	}
}

func checkAtlasProjectRemoved(projectID string) func() bool {
	return func() bool {
		_, r, err := atlasClient.ProjectsApi.GetGroup(context.Background(), projectID).Execute()
		if err != nil {
			if httputil.StatusCode(r) == http.StatusNotFound {
				return true
			}
		}
		return false
	}
}

func validateDeploymentCreatingFunc(g Gomega) func(a api.AtlasCustomResource) {
	startedCreation := false
	return func(a api.AtlasCustomResource) {
		c := a.(*akov2.AtlasDeployment)
		if c.Status.StateName != "" {
			startedCreation = true
		}
		// When the create request has been made to Atlas - we expect the following status
		if startedCreation {
			g.Expect(c.Status.StateName).To(Equal("CREATING"), fmt.Sprintf("Current conditions: %+v", c.Status.Conditions))
			expectedConditionsMatchers := conditions.MatchConditions(
				api.FalseCondition(api.DeploymentReadyType).WithReason(string(workflow.DeploymentCreating)).WithMessageRegexp("deployment is provisioning"),
				api.FalseCondition(api.ReadyType),
				api.TrueCondition(api.ValidationSucceeded),
			)
			g.Expect(c.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
		} else {
			// Otherwise there could have been some exception in Atlas on creation - let's check the conditions
			condition, ok := conditions.FindConditionByType(c.Status.Conditions, api.DeploymentReadyType)
			g.Expect(ok).To(BeFalse(), fmt.Sprintf("Unexpected condition: %v", condition))
		}
	}
}
