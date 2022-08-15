package e2e_importer_test

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"crypto/rand"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas/mongodbatlas"
)

var _ = Describe("Importer should import from Atlas", func() {

	var project *mongodbatlas.Project
	var deployment *mongodbatlas.Cluster

	var _ = BeforeEach(func() {
		// Create a project
		randomId, _ := rand.Int(rand.Reader, big.NewInt(1000))
		testProject := mongodbatlas.Project{
			OrgID: atlasClient.OrgID,
			Name:  fmt.Sprintf("TestProject-%d", randomId),
		}
		var err error
		project, _, err = atlasClient.Client.Projects.Create(context.Background(), &testProject, nil)
		Expect(err).To(BeNil())
		Expect(project).NotTo(BeNil())
		Expect(project.ID).NotTo(BeEmpty())

		// Create deployment
		diskSize := 10.0
		testDeployment := mongodbatlas.Cluster{
			ClusterType: "REPLICASET",
			DiskSizeGB:  &diskSize,
			GroupID:     project.ID,
			Name:        fmt.Sprintf("TestDeployment-%d", randomId),

			ProviderSettings: &mongodbatlas.ProviderSettings{
				ProviderName:        "TENANT",
				BackingProviderName: "AWS",
				InstanceSizeName:    "M2",
				RegionName:          "US_EAST_1",
			},
		}
		deployment, _, err = atlasClient.Client.Clusters.Create(context.Background(), project.ID, &testDeployment)
		Expect(err).To(BeNil())
		Expect(deployment.ID).NotTo(BeEmpty())
	})

	var _ = AfterEach(func() {
		//Remove the deployment
		atlasClient.Client.Clusters.Delete(context.Background(), project.ID, deployment.Name)
		Eventually(func() bool {
			_, _, err := atlasClient.Client.Clusters.Get(context.Background(), project.ID, deployment.Name)
			return err != nil
		}).WithPolling(5 * time.Second).WithTimeout(10 * time.Minute).Should(BeTrue())

		//Remove the project
		_, err := atlasClient.Client.Projects.Delete(context.Background(), project.ID)
		Expect(err).To(BeNil())
	})

	It("Import All", func() {
		Expect(project.ID).NotTo(BeEmpty())
		Expect(deployment.ID).NotTo(BeEmpty())

		// Configure the tool, and then run the entrypoint of the tool

		// Check if resources are in the cluster are equal to the one we created
		// Get the project resource from the Cluster
		// Get the deployment resource
		// How to correctly convert K8S Project to AtlasProject and Deployments
	})
})
