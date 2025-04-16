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
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/resources"
)

var _ = Describe("AtlasDeployment", Label("int", "AtlasDeployment", "independent-crd"), func() {
	var project *akov2.AtlasProject
	var deployment *akov2.AtlasDeployment

	BeforeEach(func() {
		By("Starting the operator", func() {
			prepareControllers(false)
		})

		By("Creating connection secret and project", func() {
			project = createProject(createConnectionSecret())
		})
	})

	Describe("Manage an independent deployment resource", func() {
		It("Successfully creates and manage a independent CRD", func(ctx context.Context) {
			deployment = akov2.DefaultAWSDeployment(namespace.Name, project.Name).
				Lightweight()

			By("Failing to create deployment with duplicated reference to the project", func() {
				deployment.Spec.ExternalProjectRef = &akov2.ExternalProjectReference{
					ID: project.ID(),
				}
				deployment.Spec.ConnectionSecret = &api.LocalObjectReference{
					Name: project.Spec.ConnectionSecret.Name,
				}

				Expect(k8sClient.Create(ctx, deployment)).ToNot(Succeed())
			})

			By("Creating a independent deployment resource", func() {
				deployment.Spec.ProjectRef = nil
				Expect(k8sClient.Create(ctx, deployment)).To(Succeed())
				Eventually(func(g Gomega) bool {
					return resources.CheckCondition(k8sClient, deployment, api.TrueCondition(api.ReadyType), validateDeploymentCreatingFunc(g))
				}).WithTimeout(10 * time.Minute).WithPolling(interval).Should(BeTrue())
			})

			By("Removing deployment", func() {
				Expect(k8sClient.Delete(ctx, deployment)).To(Succeed())
				Eventually(checkAtlasDeploymentRemoved(project.ID(), deployment.GetDeploymentName())).
					WithTimeout(5 * time.Minute).WithPolling(interval).Should(BeTrue())
			})
		})
	})

	Describe("Independent deployment created and 2 users are added, 1 linked and other standalone", func() {
		It("Should reconcile resource and create connection secrets", func(ctx context.Context) {
			var linkedDBUser *akov2.AtlasDatabaseUser
			var independentDBUser *akov2.AtlasDatabaseUser

			By("Creating an independent deployment", func() {
				deployment = akov2.DefaultAWSDeployment(namespace.Name, project.Name).
					WithExternaLProject(project.ID(), project.Spec.ConnectionSecret.Name).
					Lightweight()

				Expect(k8sClient.Create(ctx, deployment)).To(Succeed())
				Eventually(func(g Gomega) bool {
					return resources.CheckCondition(k8sClient, deployment, api.TrueCondition(api.ReadyType), validateDeploymentCreatingFunc(g))
				}).WithTimeout(10 * time.Minute).WithPolling(interval).Should(BeTrue())
			})

			By("Creating password secret for the users", func() {
				passwordSecret := buildPasswordSecret(namespace.Name, UserPasswordSecret, DBUserPassword)

				Expect(k8sClient.Create(ctx, &passwordSecret)).To(Succeed())
			})

			By("Creating linked database user", func() {
				linkedDBUser = akov2.DefaultDBUser(namespace.Name, "linked-db-user", project.Name).
					WithPasswordSecret(UserPasswordSecret)

				Expect(k8sClient.Create(ctx, linkedDBUser)).ToNot(HaveOccurred())
				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, linkedDBUser, api.TrueCondition(api.ReadyType))
				}).WithTimeout(5 * time.Minute).WithPolling(interval).Should(BeTrue())
			})

			By("Creating independent database user", func() {
				independentDBUser = akov2.DefaultDBUser(namespace.Name, "independent-db-user", project.Name).
					WithExternaLProject(project.ID(), project.Spec.ConnectionSecret.Name).
					WithPasswordSecret(UserPasswordSecret)

				Expect(k8sClient.Create(ctx, independentDBUser)).ToNot(HaveOccurred())
				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, independentDBUser, api.TrueCondition(api.ReadyType))
				}).WithTimeout(5 * time.Minute).WithPolling(interval).Should(BeTrue())
			})

			By("Users should have their connections secrets created", func() {
				checkNumberOfConnectionSecrets(k8sClient, *project, namespace.Name, 2)
				validateSecret(k8sClient, *project, *deployment, *linkedDBUser)
				validateSecret(k8sClient, *project, *deployment, *independentDBUser)
			})

			By("Removing secrets when deployment is deleted", func() {
				Expect(k8sClient.Delete(ctx, deployment)).To(Succeed())
				Eventually(func(g Gomega) {
					g.Expect(getNumOfConnectionSecrets(ctx, project)).To(Equal(0))
				}).WithTimeout(1 * time.Minute).WithPolling(interval).Should(Succeed())
			})

			By("Removing users", func() {
				Expect(k8sClient.Delete(context.Background(), linkedDBUser)).To(Succeed())
				Eventually(checkAtlasDatabaseUserRemoved(project.ID(), *linkedDBUser)).
					WithTimeout(2 * time.Minute).WithPolling(PollingInterval).Should(BeFalse())

				Expect(k8sClient.Delete(context.Background(), independentDBUser)).To(Succeed())
				Eventually(checkAtlasDatabaseUserRemoved(project.ID(), *independentDBUser)).
					WithTimeout(2 * time.Minute).WithPolling(PollingInterval).Should(BeFalse())
			})
		})
	})

	Describe("Add 2 users are, 1 linked and other standalone and then create an independent deployment", func() {
		It("Should reconcile resource and create connection secrets", func(ctx context.Context) {
			var linkedDBUser *akov2.AtlasDatabaseUser
			var independentDBUser *akov2.AtlasDatabaseUser

			By("Creating password secret for the users", func() {
				passwordSecret := buildPasswordSecret(namespace.Name, UserPasswordSecret, DBUserPassword)

				Expect(k8sClient.Create(ctx, &passwordSecret)).To(Succeed())
			})

			By("Creating linked database user", func() {
				linkedDBUser = akov2.DefaultDBUser(namespace.Name, "linked-db-user", project.Name).
					WithPasswordSecret(UserPasswordSecret)

				Expect(k8sClient.Create(ctx, linkedDBUser)).ToNot(HaveOccurred())
				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, linkedDBUser, api.TrueCondition(api.ReadyType))
				}).WithTimeout(5 * time.Minute).WithPolling(interval).Should(BeTrue())
			})

			By("Creating independent database user", func() {
				independentDBUser = akov2.DefaultDBUser(namespace.Name, "independent-db-user", project.Name).
					WithExternaLProject(project.ID(), project.Spec.ConnectionSecret.Name).
					WithPasswordSecret(UserPasswordSecret)

				Expect(k8sClient.Create(ctx, independentDBUser)).ToNot(HaveOccurred())
				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, independentDBUser, api.TrueCondition(api.ReadyType))
				}).WithTimeout(5 * time.Minute).WithPolling(interval).Should(BeTrue())
			})

			By("Users should not have connections secrets created", func() {
				checkNumberOfConnectionSecrets(k8sClient, *project, namespace.Name, 0)
			})

			By("Creating an independent deployment", func() {
				deployment = akov2.DefaultAWSDeployment(namespace.Name, project.Name).
					WithExternaLProject(project.ID(), project.Spec.ConnectionSecret.Name).
					Lightweight()

				Expect(k8sClient.Create(ctx, deployment)).To(Succeed())
				Eventually(func(g Gomega) bool {
					return resources.CheckCondition(k8sClient, deployment, api.TrueCondition(api.ReadyType), validateDeploymentCreatingFunc(g))
				}).WithTimeout(10 * time.Minute).WithPolling(interval).Should(BeTrue())
			})

			By("Users should have their connections secrets created", func() {
				checkNumberOfConnectionSecrets(k8sClient, *project, namespace.Name, 2)
				validateSecret(k8sClient, *project, *deployment, *linkedDBUser)
				validateSecret(k8sClient, *project, *deployment, *independentDBUser)
			})

			By("Removing secrets when deployment is deleted", func() {
				Expect(k8sClient.Delete(ctx, deployment)).To(Succeed())
				Eventually(func(g Gomega) {
					g.Expect(getNumOfConnectionSecrets(ctx, project)).To(Equal(0))
				}).WithTimeout(1 * time.Minute).WithPolling(interval).Should(Succeed())
			})

			By("Removing users", func() {
				Expect(k8sClient.Delete(context.Background(), linkedDBUser)).To(Succeed())
				Eventually(checkAtlasDatabaseUserRemoved(project.ID(), *linkedDBUser)).
					WithTimeout(2 * time.Minute).WithPolling(PollingInterval).Should(BeFalse())

				Expect(k8sClient.Delete(context.Background(), independentDBUser)).To(Succeed())
				Eventually(checkAtlasDatabaseUserRemoved(project.ID(), *independentDBUser)).
					WithTimeout(2 * time.Minute).WithPolling(PollingInterval).Should(BeFalse())
			})
		})
	})

	AfterEach(func() {
		deleteProjectFromKubernetes(project)
		removeControllersAndNamespace()
	})
})

func getNumOfConnectionSecrets(ctx context.Context, project *akov2.AtlasProject) int {
	secretList := corev1.SecretList{}
	Expect(k8sClient.List(ctx, &secretList, client.InNamespace(project.Namespace))).To(Succeed())

	num := 0
	for _, item := range secretList.Items {
		GinkgoWriter.Println(item.Name)
		if strings.HasPrefix(item.Name, kube.NormalizeIdentifier(project.Spec.Name)) {
			num++
		}
	}

	return num
}
