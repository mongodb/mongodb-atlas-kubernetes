package project

import (
	"context"
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("AtlasProject", func() {
	var (
		namespace        corev1.Namespace
		connectionSecret corev1.Secret
	)

	BeforeEach(func() {
		namespace = corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				// TODO name namespace by the name of the project and include the creation date/time
				Namespace: "test",
				Name:      "test",
			},
		}
		Expect(k8sClient.Create(context.Background(), &namespace)).ToNot(HaveOccurred())

		connectionSecret = corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-atlas-key",
				Namespace: namespace.Name,
			},
			StringData: map[string]string{"orgId": connection.OrgID, "publicApiKey": connection.PublicKey, "privateApiKey": connection.PrivateKey},
		}
		Expect(k8sClient.Create(context.Background(), &connectionSecret)).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		err := k8sClient.Delete(context.Background(), &namespace)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("Creating the project", func() {
		It("Succeeds", func() {
			project := mdbv1.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-project",
					Namespace: namespace.Name,
				},
				Spec: mdbv1.AtlasProjectSpec{
					Name:             "Test Project",
					ConnectionSecret: &mdbv1.ResourceRef{Name: connectionSecret.Name},
				},
			}
			Expect(k8sClient.Create(context.Background(), &project)).ToNot(HaveOccurred())
		})
	})
})
