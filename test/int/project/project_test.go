package project

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
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("AtlasProject", func() {
	const interval = time.Second * 1

	var (
		namespace        corev1.Namespace
		connectionSecret corev1.Secret
		createdProject   *mdbv1.AtlasProject
	)

	BeforeEach(func() {
		createdProject = &mdbv1.AtlasProject{}
		namespace = corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "test",
				// TODO name namespace by the name of the project and include the creation date/time to perform GC
				GenerateName: "test",
			},
		}
		By("Creating the namespace " + namespace.Name)
		Expect(k8sClient.Create(context.Background(), &namespace)).ToNot(HaveOccurred())

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

		By("Removing the namespace " + namespace.Name)
		err := k8sClient.Delete(context.Background(), &namespace)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("Creating the project", func() {
		It("Should Succeed", func() {
			expectedProject := testAtlasProject(namespace.Name, "test-project", "Test Project", connectionSecret.Name)
			Expect(k8sClient.Create(context.Background(), &expectedProject)).ToNot(HaveOccurred())

			Eventually(waitForProject(expectedProject, createdProject, status.TrueCondition(status.ProjectReadyType)),
				20, interval).Should(BeTrue())

			Expect(createdProject.Status.ID).NotTo(BeNil())

			atlasProject, _, err := atlasClient.Projects.GetOneProject(context.Background(), createdProject.Status.ID)
			Expect(err).ToNot(HaveOccurred())

			Expect(atlasProject.Name).To(Equal(expectedProject.Spec.Name))
		})
		It("Should fail if Secret is wrong", func() {
			expectedProject := testAtlasProject(namespace.Name, "test-project", "Test Project", "non-existent-secret")
			Expect(k8sClient.Create(context.Background(), &expectedProject)).ToNot(HaveOccurred())

			expectedCondition := status.FalseCondition(status.ProjectReadyType, string(workflow.AtlasCredentialsNotProvided))
			Eventually(waitForProject(expectedProject, createdProject, expectedCondition),
				10, interval).Should(BeTrue())

			_, _, err := atlasClient.Projects.GetOneProjectByName(context.Background(), "Test Project")

			// "NOT_IN_GROUP" is what is returned if the project is not found
			var apiError *mongodbatlas.ErrorResponse
			Expect(errors.As(err, &apiError)).To(Equal(true))
			Expect(apiError.ErrorCode).To(Equal(atlas.NotInGroup))
		})
	})

	Describe("Updating the project", func() {
		It("Should Succeed", func() {
			By("Creating the project first")

			expectedProject := testAtlasProject(namespace.Name, "test-project", "Test Project", connectionSecret.Name)
			Expect(k8sClient.Create(context.Background(), &expectedProject)).ToNot(HaveOccurred())

			Eventually(waitForProject(expectedProject, createdProject, status.TrueCondition(status.ProjectReadyType)),
				20, interval).Should(BeTrue())

			By("Updating the project (the existing project is expected to be read)")

			createdProject.Spec.ProjectIPAccessList = []mdbv1.ProjectIPAccessList{{CIDRBlock: "0.0.0.0/0"}}
			Expect(k8sClient.Update(context.Background(), createdProject)).ToNot(HaveOccurred())

			// TODO make this deterministic: CLOUDP-80550
			time.Sleep(10 * time.Second)

			Expect(readAtlasProject(expectedProject, createdProject)).To(Equal(true))
			Expect(createdProject.Status.Conditions).To(ContainElement(testutil.MatchCondition(status.TrueCondition(status.ProjectReadyType))))
			atlasProject, _, err := atlasClient.Projects.GetOneProject(context.Background(), createdProject.Status.ID)
			Expect(err).ToNot(HaveOccurred())

			Expect(atlasProject.Name).To(Equal(expectedProject.Spec.Name))
		})
	})
})

// TODO builders
func testAtlasProject(namespace, name, projectName, connectionSecretName string) mdbv1.AtlasProject {
	return mdbv1.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: mdbv1.AtlasProjectSpec{
			Name:             projectName,
			ConnectionSecret: &mdbv1.ResourceRef{Name: connectionSecretName},
		},
	}
}

// waitForProject waits until the AtlasProject reaches some state - this is configured by 'expectedCondition'
func waitForProject(project mdbv1.AtlasProject, createdProject *mdbv1.AtlasProject, expectedCondition status.Condition) func() bool {
	return func() bool {
		if ok := readAtlasProject(project, createdProject); !ok {
			return false
		}
		match, err := ContainElement(testutil.MatchCondition(expectedCondition)).Match(createdProject.Status.Conditions)
		if err != nil || !match {
			return false
		}
		return true
	}
}

func readAtlasProject(project mdbv1.AtlasProject, createdProject *mdbv1.AtlasProject) bool {
	if err := k8sClient.Get(context.Background(), kube.ObjectKeyFromObject(&project), createdProject); err != nil {
		// The only error we tolerate is "not found"
		Expect(apiErrors.IsNotFound(err)).To(Equal(true))
		return false
	}
	return true
}
