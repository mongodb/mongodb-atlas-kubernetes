package int

import (
	"context"
	"errors"
	"fmt"
	"time"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/testutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas/mongodbatlas"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("AtlasProject", func() {
	const interval = time.Second * 1

	var (
		connectionSecret corev1.Secret
		createdProject   *mdbv1.AtlasProject
	)

	BeforeEach(func() {
		prepareControllers()

		createdProject = &mdbv1.AtlasProject{}

		connectionSecret = corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-atlas-key",
				Namespace: namespace.Name,
			},
			StringData: map[string]string{"orgId": connection.OrgID, "publicApiKey": connection.PublicKey, "privateApiKey": connection.PrivateKey},
		}
		By(fmt.Sprintf("Creating the Secret %s", kube.ObjectKeyFromObject(&connectionSecret)))
		Expect(k8sClient.Create(context.Background(), &connectionSecret)).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		if createdProject != nil && createdProject.Status.ID != "" {
			By("Removing Atlas Project " + createdProject.Status.ID)
			_, err := atlasClient.Projects.Delete(context.Background(), createdProject.Status.ID)
			Expect(err).ToNot(HaveOccurred())
		}
		removeControllersAndNamespace()
	})

	Describe("Creating the project", func() {
		It("Should Succeed", func() {
			expectedProject := testAtlasProject(namespace.Name, namespace.Name, connectionSecret.Name)
			createdProject.ObjectMeta = expectedProject.ObjectMeta
			Expect(k8sClient.Create(context.Background(), expectedProject)).ToNot(HaveOccurred())

			Eventually(testutil.WaitFor(k8sClient, createdProject, status.TrueCondition(status.ReadyType)),
				20, interval).Should(BeTrue())

			Expect(createdProject.Status.ID).NotTo(BeNil())
			expectedConditionsMatchers := testutil.MatchConditions(
				status.TrueCondition(status.ProjectReadyType),
				status.FalseCondition(status.IPAccessListReadyType),
				status.TrueCondition(status.ReadyType),
			)
			Expect(createdProject.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
			Expect(createdProject.Status.ObservedGeneration).To(Equal(createdProject.Generation))

			// Atlas
			atlasProject, _, err := atlasClient.Projects.GetOneProject(context.Background(), createdProject.Status.ID)
			Expect(err).ToNot(HaveOccurred())

			Expect(atlasProject.Name).To(Equal(expectedProject.Spec.Name))
		})
		It("Should fail if Secret is wrong", func() {
			expectedProject := testAtlasProject(namespace.Name, namespace.Name, "non-existent-secret")
			createdProject.ObjectMeta = expectedProject.ObjectMeta
			Expect(k8sClient.Create(context.Background(), expectedProject)).ToNot(HaveOccurred())

			expectedCondition := status.FalseCondition(status.ProjectReadyType).WithReason(string(workflow.AtlasCredentialsNotProvided))
			Eventually(testutil.WaitFor(k8sClient, createdProject, expectedCondition),
				10, interval).Should(BeTrue())

			Expect(createdProject.Status.ObservedGeneration).To(Equal(createdProject.Generation))
			expectedConditionsMatchers := testutil.MatchConditions(
				status.FalseCondition(status.ProjectReadyType),
				status.FalseCondition(status.ReadyType),
			)
			Expect(createdProject.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
			Expect(createdProject.Status.ID).To(BeEmpty())
			Expect(createdProject.Status.ObservedGeneration).To(Equal(createdProject.Generation))

			// Atlas
			_, _, err := atlasClient.Projects.GetOneProjectByName(context.Background(), expectedProject.Spec.Name)

			// "NOT_IN_GROUP" is what is returned if the project is not found
			var apiError *mongodbatlas.ErrorResponse
			Expect(errors.As(err, &apiError)).To(BeTrue(), "Error occurred: "+err.Error())
			Expect(apiError.ErrorCode).To(Equal(atlas.NotInGroup))
		})
	})

	Describe("Updating the project", func() {
		It("Should Succeed", func() {
			By("Creating the project first")

			expectedProject := testAtlasProject(namespace.Name, namespace.Name, connectionSecret.Name)
			createdProject.ObjectMeta = expectedProject.ObjectMeta
			Expect(k8sClient.Create(context.Background(), expectedProject)).ToNot(HaveOccurred())

			Eventually(testutil.WaitFor(k8sClient, createdProject, status.TrueCondition(status.ReadyType)),
				20, interval).Should(BeTrue())

			// Updating (the existing project is expected to be read from Atlas)
			By("Updating the project")

			createdProject.Spec.ProjectIPAccessList = []mdbv1.ProjectIPAccessList{{CIDRBlock: "0.0.0.0/0"}}
			Expect(k8sClient.Update(context.Background(), createdProject)).To(Succeed())

			Eventually(testutil.WaitFor(k8sClient, createdProject, status.TrueCondition(status.ReadyType)),
				10, interval).Should(BeTrue())

			Expect(testutil.ReadAtlasResource(k8sClient, createdProject)).To(BeTrue())
			Expect(createdProject.Status.Conditions).To(ContainElement(testutil.MatchCondition(status.TrueCondition(status.ProjectReadyType))))

			// Atlas
			atlasProject, _, err := atlasClient.Projects.GetOneProject(context.Background(), createdProject.Status.ID)
			Expect(err).ToNot(HaveOccurred())
			Expect(atlasProject.Name).To(Equal(expectedProject.Spec.Name))
		})
	})
})

// TODO builders
func testAtlasProject(namespace, name, connectionSecretName string) *mdbv1.AtlasProject {
	return &mdbv1.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-project",
			Namespace: namespace,
		},
		Spec: mdbv1.AtlasProjectSpec{
			Name:             name,
			ConnectionSecret: &mdbv1.ResourceRef{Name: connectionSecretName},
		},
	}
}
