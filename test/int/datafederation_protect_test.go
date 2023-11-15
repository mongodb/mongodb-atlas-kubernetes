package int

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas/mongodbatlas"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/testutil"
)

var _ = Describe("AtlasProject", Label("int", "AtlasDataFederation", "protection-enabled"), func() {
	const (
		interval               = PollingInterval
		dataFederationBaseName = "test-data-federation-%s"
	)

	var (
		connectionSecret       corev1.Secret
		testProject            *mdbv1.AtlasProject
		testNamespace          *corev1.Namespace
		stopManager            context.CancelFunc
		testDataFederation     *mdbv1.AtlasDataFederation
		testDataFederationName string
		manualDeletion         bool
	)

	BeforeEach(func() {
		By("Starting the operator", func() {
			testNamespace, stopManager = prepareControllers(true)
			Expect(testNamespace).ToNot(BeNil())
			Expect(stopManager).ToNot(BeNil())
		})

		By("Creating project connection secret", func() {
			connectionSecret = buildConnectionSecret(fmt.Sprintf("%s-atlas-key", testNamespace.Name))
			Expect(k8sClient.Create(context.Background(), &connectionSecret)).To(Succeed())
		})

		By("Creating a project in the cluster", func() {
			testProject = mdbv1.DefaultProject(testNamespace.Name, connectionSecret.Name).
				WithIPAccessList(project.NewIPAccessList().
					WithCIDR("0.0.0.0/0"))
			Expect(k8sClient.Create(context.Background(), testProject, &client.CreateOptions{})).To(Succeed())

			Eventually(func() bool {
				return testutil.CheckCondition(k8sClient, testProject, status.TrueCondition(status.ReadyType))
			}).WithTimeout(15 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
		})

		By("Setting up DataFederation struct", func() {
			testDataFederation = &mdbv1.AtlasDataFederation{}
			testDataFederationName = fmt.Sprintf(dataFederationBaseName, testNamespace.Name)
		})

	})

	AfterEach(func() {
		By("Deleting project connection secret", func() {
			Expect(k8sClient.Delete(context.Background(), &connectionSecret)).To(Succeed())
		})

		if !manualDeletion {
			By("Removing Atlas DataFederation "+testDataFederationName, func() {
				_, err := atlasClient.DataFederation.Delete(context.Background(), testProject.ID(), testDataFederation.Spec.Name)
				Expect(err).To(BeNil())
			})
		}

		By("Removing Atlas Project "+testProject.Status.ID, func() {
			_, err := atlasClient.Projects.Delete(context.Background(), testProject.ID())
			Expect(err).To(BeNil())
		})

		By("Stopping the operator", func() {
			stopManager()
			err := k8sClient.Delete(context.Background(), testNamespace)
			Expect(err).ToNot(HaveOccurred())
		})

	})

	Describe("Operator is running with deletion protection enabled", func() {
		It("Creates a data federation and protects it from deletion", func() {
			By("Creating a DataFederation instance", func() {
				testDataFederation = mdbv1.NewDataFederationInstance(testProject.Name, testDataFederationName, testNamespace.Name)
				Expect(k8sClient.Create(context.Background(), testDataFederation)).ShouldNot(HaveOccurred())

				Eventually(func(g Gomega) {
					df, _, err := dataFederationClient.Get(context.Background(), testProject.ID(), testDataFederation.Spec.Name)
					g.Expect(err).ShouldNot(HaveOccurred())
					g.Expect(df).NotTo(BeNil())
				}).WithTimeout(20 * time.Minute).WithPolling(15 * time.Second).ShouldNot(HaveOccurred())
			})

			// nolint:dupl
			By("Deleting a data federation instance in cluster doesn't delete it from Atlas", func() {
				Expect(k8sClient.Delete(context.Background(), testDataFederation, &client.DeleteOptions{})).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(k8sClient.Get(context.Background(), client.ObjectKeyFromObject(testDataFederation), testDataFederation, &client.GetOptions{})).ToNot(Succeed())
					dataFederation, _, err := atlasClient.DataFederation.Get(context.Background(), testProject.ID(), testDataFederation.Spec.Name)
					g.Expect(err).To(BeNil())
					g.Expect(dataFederation).ToNot(BeNil())
				}).WithTimeout(5 * time.Minute).WithPolling(PollingInterval).Should(Succeed())
			})
		})

		It("Adds an existing Atlas data federation and protects it from being deleted", func() {
			By("Creating a data federation instance in Atlas", func() {
				df := &mongodbatlas.DataFederationInstance{
					Name: testDataFederationName,
				}

				_, _, err := atlasClient.DataFederation.Create(context.Background(), testProject.ID(), df)
				Expect(err).To(BeNil())
				Eventually(func(g Gomega) {
					atlasDataFederation, _, err := atlasClient.DataFederation.Get(context.Background(), testProject.ID(), testDataFederationName)
					g.Expect(err).To(BeNil())
					g.Expect(atlasDataFederation).ToNot(BeNil())
					g.Expect(atlasDataFederation.State).Should(Equal("ACTIVE"))
				}).WithTimeout(15 * time.Minute).WithPolling(PollingInterval).Should(Succeed())
			})

			By("Creating a data federation instance in the cluster", func() {
				testDataFederation = mdbv1.NewDataFederationInstance(testProject.Name, testDataFederationName, testNamespace.Name)
				Expect(k8sClient.Create(context.Background(), testDataFederation)).ShouldNot(HaveOccurred())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, testDataFederation, status.TrueCondition(status.ReadyType))
				}).WithTimeout(15 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})

			// nolint:dupl
			By("Deleting a data federation instance in the cluster does not delete it in Atlas", func() {
				Expect(k8sClient.Delete(context.Background(), testDataFederation, &client.DeleteOptions{})).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(k8sClient.Get(context.Background(), client.ObjectKeyFromObject(testDataFederation), testDataFederation, &client.GetOptions{})).ToNot(Succeed())
					dataFederation, _, err := atlasClient.DataFederation.Get(context.Background(), testProject.ID(), testDataFederation.Spec.Name)
					g.Expect(err).To(BeNil())
					g.Expect(dataFederation).ToNot(BeNil())
				}).WithTimeout(5 * time.Minute).WithPolling(PollingInterval).Should(Succeed())
			})
		})

		It("Creates a data federation instance and annotates it for deletion", func() {
			By("Creating a data federation instance in the cluster", func() {
				testDataFederation = mdbv1.NewDataFederationInstance(testProject.Name, testDataFederationName, testNamespace.Name).
					WithAnnotations(map[string]string{customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyDelete})
				Expect(k8sClient.Create(context.Background(), testDataFederation)).ShouldNot(HaveOccurred())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, testDataFederation, status.TrueCondition(status.ReadyType))
					// TODO: Modify timeouts
				}).WithTimeout(15 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})

			By("Deleting annotated data federation instance in cluster should delete it from  Atlas", func() {
				Expect(k8sClient.Delete(context.Background(), testDataFederation, &client.DeleteOptions{})).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(k8sClient.Get(context.Background(), client.ObjectKeyFromObject(testDataFederation), testDataFederation, &client.GetOptions{})).ToNot(Succeed())
					dataFederation, _, err := atlasClient.DataFederation.Get(context.Background(), testProject.ID(), testDataFederation.Spec.Name)
					g.Expect(err).ToNot(BeNil())
					g.Expect(dataFederation).To(BeNil())
				}).WithTimeout(5 * time.Minute).WithPolling(PollingInterval).Should(Succeed())
				manualDeletion = true
			})
		})
	})

})
