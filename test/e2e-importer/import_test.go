package e2e_importer_test

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/mongodb/mongodb-atlas-kubernetes/cmd/atlas-import/importer"
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"

	"crypto/rand"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas/mongodbatlas"
)

// TODO possible test cases :
//- Generate resources definitions, instantiate them in Atlas, import them, verify they are equivalent
//- Create two projects with multiple deployments in atlas, configure the script to import only one project and one
//  deployment, ensure the import is done correctly
//- Ensure that no resource is modified during or after the importation, when the operator is running at the same time
//- Ensure that import cannot “overwrite” or duplicate already existing resources ?

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

	// TODO clean cluster after each import
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

		importConfig := importer.AtlasImportConfig{
			OrgID:            atlasClient.OrgID,
			PublicKey:        atlasClient.Public,
			PrivateKey:       atlasClient.Private,
			AtlasDomain:      atlasClient.Client.BaseURL.String(),
			ImportNamespace:  "import-namespace",
			ImportAll:        true,
			ImportedProjects: nil,
			Verbose:          true,
		}

		err := importer.RunImports(importConfig)
		Expect(err).NotTo(HaveOccurred())

		// Check if resources in the cluster are equal to the one we created
		projectList := &v1.AtlasProjectList{}
		// Get the project resource from the Cluster
		Expect(k8sClient.List(context.Background(), projectList)).NotTo(HaveOccurred())

		projectNameSet := make(map[string]bool)
		for _, kubeProject := range projectList.Items {
			projectNameSet[kubeProject.Spec.Name] = true
		}

		Expect(projectNameSet).Should(HaveKey(project.Name))

		// Get the deployment resource

		// TODO write a method to compare the resources we added to Atlas at the beginning, to the one imported in k8s
		// methods existing in atlasdeployment_types_test.go will probably be useful for that
		// This method will also be useful to unit test every conversion method (MaintenanceWindowFromAtlas, IPAccessListFromAtlas...)
	})

})
