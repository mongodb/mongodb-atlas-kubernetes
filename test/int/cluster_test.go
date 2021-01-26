package int

import (
	"context"
	"fmt"
	"net/http"
	"time"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/testutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("AtlasCluster", func() {
	const interval = time.Second * 1

	var (
		namespace        corev1.Namespace
		connectionSecret corev1.Secret
		createdProject   *mdbv1.AtlasProject
		createdCluster   *mdbv1.AtlasCluster
	)

	BeforeEach(func() {
		createdCluster = &mdbv1.AtlasCluster{}
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

		createdProject = testAtlasProject(namespace.Name, namespace.Name, connectionSecret.Name)
		By("Creating the project " + createdProject.Name)
		Expect(k8sClient.Create(context.Background(), createdProject)).ToNot(HaveOccurred())
		Eventually(testutil.WaitFor(k8sClient, createdProject, status.TrueCondition(status.ReadyType)),
			10, interval).Should(BeTrue())
	})

	AfterEach(func() {
		if createdProject != nil && createdProject.Status.ID != "" {
			if createdCluster != nil {
				By("Removing Atlas Cluster " + createdCluster.Spec.Name)
				_, _ = atlasClient.Clusters.Delete(context.Background(), createdProject.Status.ID, createdCluster.Spec.Name)
			}
			// TODO need to wait for the cluster to get removed
			// By("Removing Atlas Project " + createdProject.Status.ID)
			// _, err := atlasClient.Projects.Delete(context.Background(), createdProject.Status.ID)
			// Expect(err).ToNot(HaveOccurred())
		}

		By("Removing the namespace " + namespace.Name)
		err := k8sClient.Delete(context.Background(), &namespace)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("Create/Update/Delete the cluster", func() {
		It("Should Succeed", func() {
			expectedCluster := testAtlasCluster(namespace.Name, "test-cluster", createdProject.Name)

			By(fmt.Sprintf("Creating the Cluster %s", kube.ObjectKeyFromObject(expectedCluster)))

			createdCluster.ObjectMeta = expectedCluster.ObjectMeta
			Expect(k8sClient.Create(context.Background(), expectedCluster)).ToNot(HaveOccurred())

			validatePending := clusterPendingFunc("CREATING", "cluster is provisioning", workflow.ClusterCreating)
			Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validatePending),
				1200, interval).Should(BeTrue())

			Expect(createdCluster.Status.ConnectionStrings).NotTo(BeNil())
			Expect(createdCluster.Status.ConnectionStrings.Standard).NotTo(BeNil())
			Expect(createdCluster.Status.ConnectionStrings.StandardSrv).NotTo(BeNil())
			Expect(createdCluster.Status.MongoDBVersion).NotTo(BeNil())
			Expect(createdCluster.Status.MongoURIUpdated).NotTo(BeNil())
			Expect(createdCluster.Status.StateName).To(Equal("IDLE"))

			expectedConditionsMatchers := testutil.MatchConditions(
				status.TrueCondition(status.ClusterReadyType),
				status.TrueCondition(status.ReadyType),
			)
			Expect(createdCluster.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
			Expect(createdCluster.Status.ObservedGeneration).To(Equal(createdCluster.Generation))

			// Atlas
			atlasCluster, _, err := atlasClient.Clusters.Get(context.Background(), createdProject.Status.ID, createdCluster.Spec.Name)
			Expect(err).ToNot(HaveOccurred())

			// Unfortunately we cannot do global checks on cluster/providerSettings fields as Atlas adds default values
			Expect(atlasCluster.Name).To(Equal(expectedCluster.Spec.Name))
			Expect(atlasCluster.ProviderSettings.InstanceSizeName).To(Equal(expectedCluster.Spec.ProviderSettings.InstanceSizeName))
			Expect(atlasCluster.ProviderSettings.ProviderName).To(Equal(expectedCluster.Spec.ProviderSettings.ProviderName))
			Expect(atlasCluster.ProviderSettings.RegionName).To(Equal(expectedCluster.Spec.ProviderSettings.RegionName))

			// TODO check connectivity to cluster

			By("Updating the Cluster")
			createdCluster.Spec.Labels = []mdbv1.LabelSpec{{Key: "int-test", Value: "true"}}
			Expect(k8sClient.Update(context.Background(), createdCluster)).To(Succeed())

			validatePending = clusterPendingFunc("UPDATING", "cluster is updating", workflow.ClusterUpdating)
			Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validatePending),
				1200, interval).Should(BeTrue())

			Expect(createdCluster.Status.ConnectionStrings).NotTo(BeNil())
			Expect(createdCluster.Status.ConnectionStrings.Standard).NotTo(BeNil())
			Expect(createdCluster.Status.ConnectionStrings.StandardSrv).NotTo(BeNil())
			Expect(createdCluster.Status.MongoDBVersion).NotTo(BeNil())
			Expect(createdCluster.Status.MongoURIUpdated).NotTo(BeNil())
			Expect(createdCluster.Status.StateName).To(Equal("IDLE"))
			Expect(createdCluster.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
			Expect(createdCluster.Status.ObservedGeneration).To(Equal(createdCluster.Generation))

			// Atlas
			atlasCluster, _, err = atlasClient.Clusters.Get(context.Background(), createdProject.Status.ID, createdCluster.Spec.Name)
			Expect(err).ToNot(HaveOccurred())

			createdAtlasCluster, err := createdCluster.Spec.Cluster()
			Expect(err).ToNot(HaveOccurred())

			Expect(atlasCluster.Name).To(Equal(createdAtlasCluster.Name))
			Expect(atlasCluster.Labels).To(Equal(createdAtlasCluster.Labels))
			Expect(atlasCluster.ProviderSettings.InstanceSizeName).To(Equal(createdAtlasCluster.ProviderSettings.InstanceSizeName))
			Expect(atlasCluster.ProviderSettings.ProviderName).To(Equal(createdAtlasCluster.ProviderSettings.ProviderName))
			Expect(atlasCluster.ProviderSettings.RegionName).To(Equal(createdAtlasCluster.ProviderSettings.RegionName))

			By("Deleting the cluster")
			err = k8sClient.Delete(context.Background(), createdCluster)
			Expect(err).ToNot(HaveOccurred())

			Eventually(checkAtlasDeleteStarted(createdProject.Status.ID, createdCluster.Name), 1200, interval).Should(BeTrue())
		})
	})
})

func clusterPendingFunc(expectedState, expectedMessage string, reason workflow.ConditionReason) func(a mdbv1.AtlasCustomResource) {
	startedCreation := false
	return func(a mdbv1.AtlasCustomResource) {
		c := a.(*mdbv1.AtlasCluster)
		if c.Status.StateName != "" {
			startedCreation = true
		}
		// When the create request has been made to Atlas - we expect the following status
		if startedCreation {
			Expect(c.Status.StateName).To(Equal(expectedState), fmt.Sprintf("Current conditions: %+v", c.Status.Conditions))
			expectedConditionsMatchers := testutil.MatchConditions(
				status.FalseCondition(status.ClusterReadyType).WithReason(string(reason)).WithMessage(expectedMessage),
				status.FalseCondition(status.ReadyType),
			)
			Expect(c.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
		} else {
			// Otherwise there could have been some exception in Atlas on creation - let's check the conditions
			condition, ok := testutil.FindConditionByType(c.Status.Conditions, status.ClusterReadyType)
			Expect(ok).To(BeFalse(), fmt.Sprintf("Unexpected condition: %v", condition))
		}
	}
}

// TODO builders
func testAtlasCluster(namespace, name, projectName string) *mdbv1.AtlasCluster {
	return &mdbv1.AtlasCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: mdbv1.AtlasClusterSpec{
			Name:    "test-cluster",
			Project: mdbv1.ResourceRef{Name: projectName},
			ProviderSettings: &mdbv1.ProviderSettingsSpec{
				InstanceSizeName: "M10",
				ProviderName:     "GCP",
				RegionName:       "EASTERN_US",
			},
		},
	}
}
func checkAtlasDeleteStarted(projectID string, clusterName string) func() bool {
	return func() bool {
		c, r, err := atlasClient.Clusters.Get(context.Background(), projectID, clusterName)
		if err != nil {
			// cluster already deleted - that's fine for us
			if r != nil && r.StatusCode == http.StatusNotFound {
				return true
			}

			return false
		}

		return c.StateName == "DELETING" || c.StateName == "DELETED"
	}
}
