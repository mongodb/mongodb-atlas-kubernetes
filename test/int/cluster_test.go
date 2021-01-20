package int

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("AtlasProject", func() {
	const interval = time.Second * 1

	var (
		namespace        corev1.Namespace
		connectionSecret corev1.Secret
		createdProject   *mdbv1.AtlasProject
		createdCluster   *mdbv1.AtlasCluster
	)

	BeforeEach(func() {
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

		createdProject = testAtlasProject(namespace.Name, "test-project", "Test Project", connectionSecret.Name)
		By("Creating the project " + createdProject.Name)
		Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		if createdProject != nil && createdProject.Status.ID != "" {
			if createdCluster != nil {
				By("Removing Atlas Cluster " + createdCluster.Spec.Name)
				_, err := atlasClient.Clusters.Delete(context.Background(), createdProject.Status.ID, createdCluster.Spec.Name)
				Expect(err).ToNot(HaveOccurred())
			}
			By("Removing Atlas Project " + createdProject.Status.ID)
			_, err := atlasClient.Projects.Delete(context.Background(), createdProject.Status.ID)
			Expect(err).ToNot(HaveOccurred())
		}

		By("Removing the namespace " + namespace.Name)
		err := k8sClient.Delete(context.Background(), &namespace)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("Creating the cluster", func() {
		It("Should Succeed", func() {
			expectedCluster := testAtlasCluster(namespace.Name, "test-cluster", createdProject.Name)
			Expect(k8sClient.Create(context.Background(), expectedCluster)).ToNot(HaveOccurred())

			Eventually(testutil.WaitFor(k8sClient, expectedCluster, createdProject, status.TrueCondition(status.ReadyType)),
				360, interval).Should(BeTrue())

			Expect(createdCluster.Status.ConnectionStrings).NotTo(BeNil())
			Expect(createdCluster.Status.MongoDBVersion).NotTo(BeNil())
			Expect(createdCluster.Status.MongoURIUpdated).To(BeNil())
			Expect(createdCluster.Status.StateName).To(Equal("IDLE"))
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

			Expect(atlasProject.Name).To(Equal(expectedCluster.Spec.Name))
		})
	})
})

// TODO builders
func testAtlasCluster(namespace, name, projectName string) *mdbv1.AtlasCluster {
	return &mdbv1.AtlasCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: mdbv1.AtlasClusterSpec{
			Name:    "test cluster",
			Project: mdbv1.ResourceRef{Name: projectName},
			ProviderSettings: &mdbv1.ProviderSettingsSpec{
				InstanceSizeName: "M10",
				ProviderName:     "AWS",
			},
		},
	}
}
