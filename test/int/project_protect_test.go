package int

import (
	"context"
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"

	"go.mongodb.org/atlas/mongodbatlas"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/testutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
)

var _ = Describe("AtlasProject", Label("int", "AtlasProject", "protection-enabled"), func() {
	var testNamespace *corev1.Namespace
	var stopManager context.CancelFunc
	var connectionSecret corev1.Secret

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
	})

	Describe("Operator is running with deletion protection enabled", func() {
		It("Creates a project and protect it to be deleted", func() {
			testProject := &mdbv1.AtlasProject{}
			projectName := fmt.Sprintf("new-project-%s", testNamespace.Name)

			By("Creating a project in the cluster", func() {
				testProject = mdbv1.NewProject(testNamespace.Name, projectName, projectName).
					WithConnectionSecret(connectionSecret.Name)
				Expect(k8sClient.Create(context.TODO(), testProject, &client.CreateOptions{})).To(Succeed())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, testProject, status.TrueCondition(status.ReadyType))
				}).WithTimeout(15 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})

			// nolint:dupl
			By("Deleting project in cluster doesn't delete from Atlas", func() {
				projectID := testProject.ID()
				Expect(k8sClient.Delete(context.TODO(), testProject, &client.DeleteOptions{})).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(testProject), testProject, &client.GetOptions{})).ToNot(Succeed())

					atlasProject, _, err := atlasClient.Projects.GetOneProjectByName(context.TODO(), projectName)
					g.Expect(err).To(BeNil())
					g.Expect(atlasProject).ToNot(BeNil())
				}).WithTimeout(5 * time.Minute).WithPolling(PollingInterval).Should(Succeed())

				_, err := atlasClient.Projects.Delete(context.TODO(), projectID)
				Expect(err).To(BeNil())
			})
		})

		It("Adds an existing Atlas project and protect it to be deleted", func() {
			testProject := &mdbv1.AtlasProject{}
			projectName := fmt.Sprintf("existing-project-%s", testNamespace.Name)

			By("Creating a project in Atlas", func() {
				atlasProject := mongodbatlas.Project{
					OrgID:                     connection.OrgID,
					Name:                      projectName,
					WithDefaultAlertsSettings: toptr.MakePtr(true),
				}
				_, _, err := atlasClient.Projects.Create(context.TODO(), &atlasProject, &mongodbatlas.CreateProjectOptions{})
				Expect(err).To(BeNil())
			})

			By("Creating a project in the cluster", func() {
				testProject = mdbv1.NewProject(testNamespace.Name, projectName, projectName).
					WithConnectionSecret(connectionSecret.Name)
				Expect(k8sClient.Create(context.TODO(), testProject, &client.CreateOptions{})).To(Succeed())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, testProject, status.TrueCondition(status.ReadyType))
				}).WithTimeout(15 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})

			// nolint:dupl
			By("Deleting project in cluster doesn't delete from Atlas", func() {
				projectID := testProject.ID()
				Expect(k8sClient.Delete(context.TODO(), testProject, &client.DeleteOptions{})).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(testProject), testProject, &client.GetOptions{})).ToNot(Succeed())

					atlasProject, _, err := atlasClient.Projects.GetOneProjectByName(context.TODO(), projectName)
					g.Expect(err).To(BeNil())
					g.Expect(atlasProject).ToNot(BeNil())
				}).WithTimeout(5 * time.Minute).WithPolling(PollingInterval).Should(Succeed())

				_, err := atlasClient.Projects.Delete(context.TODO(), projectID)
				Expect(err).To(BeNil())
			})
		})

		It("Creates a project and annotate it to be deleted", func() {
			testProject := &mdbv1.AtlasProject{}
			projectName := fmt.Sprintf("new-project-%s", testNamespace.Name)

			By("Creating a project in the cluster", func() {
				testProject = mdbv1.NewProject(testNamespace.Name, projectName, projectName).
					WithAnnotations(map[string]string{customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyDelete}).
					WithConnectionSecret(connectionSecret.Name)
				Expect(k8sClient.Create(context.TODO(), testProject, &client.CreateOptions{})).To(Succeed())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, testProject, status.TrueCondition(status.ReadyType))
				}).WithTimeout(15 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})

			// nolint:dupl
			By("Deleting project in cluster should delete it from Atlas", func() {
				projectID := testProject.ID()
				Expect(k8sClient.Delete(context.TODO(), testProject, &client.DeleteOptions{})).To(Succeed())

				Eventually(func(g Gomega) {
					_, r, err := atlasClient.Projects.GetOneProject(context.TODO(), projectID)
					g.Expect(err).ToNot(BeNil())
					g.Expect(r).ToNot(BeNil())
					g.Expect(r.StatusCode).To(Equal(http.StatusNotFound))
				}).WithTimeout(5 * time.Minute).WithPolling(PollingInterval).Should(Succeed())
			})
		})
	})

	AfterEach(func() {
		By("Deleting project connection secret", func() {
			Expect(k8sClient.Delete(context.Background(), &connectionSecret)).To(Succeed())
		})

		By("Stopping the operator", func() {
			stopManager()
			err := k8sClient.Delete(context.Background(), testNamespace)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
