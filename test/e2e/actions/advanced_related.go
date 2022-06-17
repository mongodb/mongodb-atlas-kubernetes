package actions

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	mongocli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/mongocli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

const (
	atlasClusterCRD = "atlasdeployments.atlas.mongodb.com"
)

func DeleteCRDs(data *model.TestDataProvider) {
	By("Deleting CRDs", func() {
		By(fmt.Sprintf("Deleting %s", atlasClusterCRD), func() {
			kubecli.DeleteClusterResource("crd", atlasClusterCRD)
			// TODO: check CRD deletion
			By("Checking Cluster still existed", func() {
				state := mongocli.GetClusterStateName(data.Resources.ProjectID, data.Resources.Clusters[0].Spec.DeploymentSpec.Name)
				Expect(state).ShouldNot(Equal("DELETING"), "Error has Occurred")
			})
		})
	})
}
