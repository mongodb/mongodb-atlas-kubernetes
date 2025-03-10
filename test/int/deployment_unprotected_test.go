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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/resources"
)

var _ = Describe("AtlasDeployment Deletion Unprotected",
	Ordered,
	Label("AtlasDeployment", "deletion-protection", "deployment-deletion-unprotected"), func() {
		var testNamespace *corev1.Namespace
		var stopManager context.CancelFunc
		var connectionSecret corev1.Secret
		var testProject *akov2.AtlasProject

		BeforeAll(func() {
			By("Starting the operator with protection OFF", func() {
				testNamespace, stopManager = prepareControllers(false)
				Expect(testNamespace).ToNot(BeNil())
				Expect(stopManager).ToNot(BeNil())
			})

			By("Creating project connection secret", func() {
				connectionSecret = buildConnectionSecret(fmt.Sprintf("%s-atlas-key", testNamespace.Name))
				Expect(k8sClient.Create(context.Background(), &connectionSecret)).To(Succeed())
			})

			By("Creating a project", func() {
				testProject = akov2.DefaultProject(testNamespace.Name, connectionSecret.Name).WithIPAccessList(project.NewIPAccessList().WithCIDR("0.0.0.0/0"))
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

		It("removing advanced cluster from Kubernetes when protection is OFF wipes it from Atlas",
			Label("wiping-advanced-cluster"),
			func() {
				testDeployment := akov2.DefaultAWSDeployment(testNamespace.Name, testProject.Name)
				wipeDeploymentFlow(testNamespace.Name, testProject, testDeployment)
			},
		)

		It("removing flex instance from Kubernetes when protection is OFF wipes it from Atlas",
			Label("wiping-flex-instance"),
			func() {
				testDeployment := akov2.NewDefaultAWSFlexInstance(testNamespace.Name, testProject.Name)
				wipeDeploymentFlow(testNamespace.Name, testProject, testDeployment)
			},
		)
	},
)

func wipeDeploymentFlow(ns string, testProject *akov2.AtlasProject, testDeployment *akov2.AtlasDeployment) {
	By("Creating a deployment in the cluster with annotation set to delete", func() {
		testDeployment = akov2.DefaultAWSDeployment(ns, testProject.Name).Lightweight()
		Expect(k8sClient.Create(context.Background(), testDeployment, &client.CreateOptions{})).To(Succeed())
	})

	By("Waiting the deployment to settle in kubernetes", func() {
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
		}).WithTimeout(2 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
	})

	By("Checking whether the Atlas deployment got also removed", func() {
		if testDeployment.IsFlex() {
			Eventually(
				checkAtlasFlexInstanceRemoved(testProject.Status.ID, testDeployment.Spec.FlexSpec.Name),
			).WithTimeout(5 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			return
		}
		Eventually(
			checkAtlasDeploymentRemoved(testProject.Status.ID, testDeployment.Spec.DeploymentSpec.Name),
		).WithTimeout(5 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
	})
}
