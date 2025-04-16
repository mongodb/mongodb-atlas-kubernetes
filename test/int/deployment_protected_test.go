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
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/resources"
)

// nolint:dupl
var _ = Describe("AtlasDeployment Deletion Protected",
	Ordered,
	Label("AtlasDeployment", "deletion-protection", "deployment-deletion-protected"), func() {
		var testNamespace *corev1.Namespace
		var stopManager context.CancelFunc
		var connectionSecret corev1.Secret
		var testProject *akov2.AtlasProject

		BeforeAll(func() {
			By("Starting the operator with protection ON", func() {
				testNamespace, stopManager = prepareControllers(true)
				Expect(testNamespace).ToNot(BeNil())
				Expect(stopManager).ToNot(BeNil())
			})

			By("Creating project connection secret", func() {
				connectionSecret = buildConnectionSecret(fmt.Sprintf("%s-atlas-key", testNamespace.Name))
				Expect(k8sClient.Create(context.Background(), &connectionSecret)).To(Succeed())
			})

			By("Creating a project with deletion annotation", func() {
				testProject = akov2.DefaultProject(testNamespace.Name, connectionSecret.Name).WithIPAccessList(project.NewIPAccessList().WithCIDR("0.0.0.0/0"))
				customresource.SetAnnotation( // this test project must be deleted
					testProject,
					customresource.ResourcePolicyAnnotation,
					customresource.ResourcePolicyDelete,
				)
				Expect(k8sClient.Create(context.Background(), testProject, &client.CreateOptions{})).To(Succeed())

				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, testProject, api.TrueCondition(api.ReadyType))
				}).WithTimeout(3 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})
		})

		AfterAll(func() {
			By("Deleting project from k8s and atlas", func() {
				Expect(k8sClient.Delete(context.Background(), testProject, &client.DeleteOptions{})).To(Succeed())
				Eventually(
					checkAtlasProjectRemoved(testProject.Status.ID),
				).WithTimeout(3 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})

			By("Deleting project connection secret", func() {
				Expect(k8sClient.Delete(context.Background(), &connectionSecret)).To(Succeed())
			})

			By("Stopping the operator", func() {
				stopManager()
				err := k8sClient.Delete(context.Background(), testNamespace)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		It("removing advanced cluster from Kubernetes when protection is ON leaves it in Atlas",
			Label("preserving-advanced-cluster"),
			func() {
				testDeployment := akov2.DefaultAWSDeployment(testNamespace.Name, testProject.Name)
				preserveDeploymentFlow(testNamespace.Name, testProject, testDeployment)
			},
		)

		It("removing flex instance from Kubernetes when protection is ON leaves it in Atlas",
			Label("preserving-flex-instance"),
			func() {
				testDeployment := akov2.NewDefaultAWSFlexInstance(testNamespace.Name, testProject.Name)
				preserveDeploymentFlow(testNamespace.Name, testProject, testDeployment)
			},
		)
	},
)

func preserveDeploymentFlow(ns string, testProject *akov2.AtlasProject, testDeployment *akov2.AtlasDeployment) {
	By("Creating deployment in Kubernetes", func() {
		Expect(k8sClient.Create(context.Background(), testDeployment, &client.CreateOptions{})).To(Succeed())
	})

	By("Waiting the deployment to settle in Kubernetes", func() {
		Eventually(func(g Gomega) bool {
			return resources.CheckCondition(k8sClient, testDeployment, api.TrueCondition(api.ReadyType), validateDeploymentUpdatingFunc(g))
		}).WithTimeout(30 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
	})

	By("Deleting the deployment from Kubernetes", func() {
		Expect(k8sClient.Delete(context.Background(), testDeployment, &client.DeleteOptions{})).To(Succeed())
		Eventually(func() bool {
			deployment := akov2.AtlasDeployment{}
			err := k8sClient.Get(context.Background(), kube.ObjectKey(ns, testDeployment.Name), &deployment, &client.GetOptions{})
			return k8serrors.IsNotFound(err)
		}).WithTimeout(5 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
	})

	By("Checking the Atlas deployment was NOT removed", func() {
		if testDeployment.IsFlex() {
			Expect(checkAtlasFlexInstanceRemoved(testProject.Status.ID, testDeployment.Spec.FlexSpec.Name)()).To(BeFalse())
			return
		}
		Expect(checkAtlasDeploymentRemoved(testProject.Status.ID, testDeployment.Spec.DeploymentSpec.Name)()).To(BeFalse())
	})

	By("Making sure deployment gets removed from Atlas manually", func() {
		if testDeployment.IsFlex() {
			Expect(deleteFlexInstance(testProject.Status.ID, testDeployment.Spec.FlexSpec.Name)).ToNot(HaveOccurred())
			return
		}
		Expect(deleteAtlasDeployment(testProject.Status.ID, testDeployment.Spec.DeploymentSpec.Name)).ToNot(HaveOccurred())
	})
}
