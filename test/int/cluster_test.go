package int

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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
				By("Removing Atlas Cluster " + createdCluster.Name)
				Expect(k8sClient.Delete(context.Background(), createdCluster)).To(Succeed())

				Eventually(checkAtlasClusterRemoved(createdProject.Status.ID, createdCluster.Name), 600, interval).Should(BeTrue())
			}
			By("Removing Atlas Project " + createdProject.Status.ID)
			// This is a bit strange but the delete request right after the cluster is removed may fail with "Still active cluster" error
			// UI shows the cluster being deleted though. Seems to be the issue only if removal is done using API,
			// if the cluster is terminated using UI - it stays in "Deleting" state
			Eventually(removeAtlasProject(createdProject.Status.ID), 600, interval).Should(BeTrue())
		}

		By("Removing the namespace " + namespace.Name)
		err := k8sClient.Delete(context.Background(), &namespace)
		Expect(err).ToNot(HaveOccurred())
	})

	FDescribe("Create/Update the cluster", func() {
		It("Should Succeed", func() {
			expectedCluster := testAtlasCluster(namespace.Name, "test-cluster", createdProject.Name)

			By(fmt.Sprintf("Creating the Cluster %s", kube.ObjectKeyFromObject(expectedCluster)))

			createdCluster.ObjectMeta = expectedCluster.ObjectMeta
			Expect(k8sClient.Create(context.Background(), expectedCluster)).ToNot(HaveOccurred())

			Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validateClusterCreatingFunc()),
				1800, interval).Should(BeTrue())

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

			Eventually(testutil.WaitFor(k8sClient, createdCluster, status.TrueCondition(status.ReadyType), validateClusterUpdatingFunc()),
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
			print(createdCluster.Labels)
			print(createdAtlasCluster.Labels)
			Expect(atlasCluster.Labels).To(Equal(createdAtlasCluster.Labels))
			Expect(atlasCluster.ProviderSettings.InstanceSizeName).To(Equal(createdAtlasCluster.ProviderSettings.InstanceSizeName))
			Expect(atlasCluster.ProviderSettings.ProviderName).To(Equal(createdAtlasCluster.ProviderSettings.ProviderName))
			Expect(atlasCluster.ProviderSettings.RegionName).To(Equal(createdAtlasCluster.ProviderSettings.RegionName))
		})
	})
})

func validateClusterCreatingFunc() func(a mdbv1.AtlasCustomResource) {
	startedCreation := false
	return func(a mdbv1.AtlasCustomResource) {
		c := a.(*mdbv1.AtlasCluster)
		if c.Status.StateName != "" {
			startedCreation = true
		}
		// When the create request has been made to Atlas - we expect the following status
		if startedCreation {
			Expect(c.Status.StateName).To(Equal("CREATING"), fmt.Sprintf("Current conditions: %+v", c.Status.Conditions))
			expectedConditionsMatchers := testutil.MatchConditions(
				status.FalseCondition(status.ClusterReadyType).WithReason(string(workflow.ClusterCreating)).WithMessage("cluster is provisioning"),
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
func validateClusterUpdatingFunc() func(a mdbv1.AtlasCustomResource) {
	isIdle := true
	return func(a mdbv1.AtlasCustomResource) {
		c := a.(*mdbv1.AtlasCluster)
		// It's ok if the first invocations see IDLE
		if c.Status.StateName != "IDLE" {
			isIdle = false
		}
		// When the create request has been made to Atlas - we expect the following status
		if !isIdle {
			Expect(c.Status.StateName).To(Equal("UPDATING"), fmt.Sprintf("Current conditions: %+v", c.Status.Conditions))
			expectedConditionsMatchers := testutil.MatchConditions(
				status.FalseCondition(status.ClusterReadyType).WithReason(string(workflow.ClusterUpdating)).WithMessage("cluster is updating"),
				status.FalseCondition(status.ReadyType),
			)
			Expect(c.Status.Conditions).To(ConsistOf(expectedConditionsMatchers))
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
			Name:    "test-atlas-cluster",
			Project: mdbv1.ResourceRef{Name: projectName},
			ProviderSettings: &mdbv1.ProviderSettingsSpec{
				InstanceSizeName: "M10",
				ProviderName:     "GCP",
				RegionName:       "EASTERN_US",
			},
		},
	}
}

// checkAtlasClusterRemoved returns true if the Atlas Cluster is removed from Atlas. Note the behavior: the cluster
// is removed from Atlas as soon as the DELETE API call has been made. This is different from the case when the
// cluster is terminated from UI (in this case GET request succeeds while the cluster is being terminated)
func checkAtlasClusterRemoved(projectID string, clusterName string) func() bool {
	return func() bool {
		_, r, err := atlasClient.Clusters.Get(context.Background(), projectID, clusterName)
		if err != nil {
			if r != nil && r.StatusCode == http.StatusNotFound {
				return true
			}
		}
		return false
	}
}

func removeAtlasProject(projectID string) func() bool {
	return func() bool {
		_, err := atlasClient.Projects.Delete(context.Background(), projectID)
		if err != nil {
			var apiError *mongodbatlas.ErrorResponse
			Expect(errors.As(err, &apiError)).To(BeTrue())
			Expect(apiError.ErrorCode).To(Equal(atlas.CannotCloseGroupActiveAtlasCluster))
			return false
		}
		return true
	}
}
