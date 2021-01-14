package project

import (
	"context"
	"fmt"
	"time"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/testutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("AtlasProject", func() {
	const timeout = time.Second * 20
	const interval = time.Second * 1

	var (
		namespace        corev1.Namespace
		connectionSecret corev1.Secret
		projectID        string
	)

	BeforeEach(func() {
		namespace = corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				// TODO name namespace by the name of the project and include the creation date/time to perform GC
				Namespace: "test",
				Name:      "test",
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
		if projectID != "" {
			By("Removing Atlas Project " + projectID)
			_, err := atlasClient.Projects.Delete(context.Background(), projectID)
			Expect(err).ToNot(HaveOccurred())
		}

		By("Removing the namespace " + namespace.Name)
		err := k8sClient.Delete(context.Background(), &namespace)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("Creating the project", func() {
		It("Succeeds", func() {
			// TODO builders
			expectedProject := mdbv1.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-project",
					Namespace: namespace.Name,
				},
				Spec: mdbv1.AtlasProjectSpec{
					Name:             "Test Project",
					ConnectionSecret: &mdbv1.ResourceRef{Name: connectionSecret.Name},
				},
			}
			Expect(k8sClient.Create(context.Background(), &expectedProject)).ToNot(HaveOccurred())

			createdProject := &mdbv1.AtlasProject{}
			Eventually(func() bool {
				if err := k8sClient.Get(context.Background(), kube.ObjectKeyFromObject(&expectedProject), createdProject); err != nil {
					// The only error we tolerate is "not found"
					Expect(apiErrors.IsNotFound(err)).To(Equal(true))
					return false
				}
				if createdProject.Status.ID != "" {
					projectID = createdProject.Status.ID
					return true
				}
				return false
			}, timeout, interval).Should(BeTrue())

			Expect(createdProject.Status.Conditions).To(ContainElement(testutil.MatchCondition(status.TrueCondition(status.ProjectReadyType))))

			atlasProject, _, err := atlasClient.Projects.GetOneProject(context.Background(), createdProject.Status.ID)
			Expect(err).ToNot(HaveOccurred())

			Expect(atlasProject.Name).To(Equal(expectedProject.Spec.Name))
		})
	})
})
