package actions

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/atlas"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

const (
	atlasDeploymentCRD = "atlasdeployments.atlas.mongodb.com"
)

func DeleteDeploymentCRWithKeepAnnotation(data *model.TestDataProvider) {
	By(fmt.Sprintf("Deleting %s", atlasDeploymentCRD), func() {
		kubecli.DeleteClusterResource("crd", atlasDeploymentCRD)
		By("Checking Cluster still existed", func() {
			aClient := atlas.GetClientOrFail()
			state := aClient.GetDeployment(data.Resources.ProjectID, data.Resources.Deployments[0].Spec.GetDeploymentName()).StateName
			Expect(state).ShouldNot(Equal("DELETING"), "Deployment is being deleted despite the keep annotation")
		})
	})
}

func ReDeployOperator(data *model.TestDataProvider) {
	By(fmt.Sprintf("Recreating %s", atlasDeploymentCRD), func() {
		deploy.NamespacedOperator(data)
	})
}

func RemoveKeepAnnotation(data *model.TestDataProvider) {
	// Remove annotation so actions.AfterEachFinalCleanup can cleanup successfully
	By("Removing keep annotation", func() {
		annotations := data.Resources.Deployments[0].ObjectMeta.GetAnnotations()
		// remove keep annotations from map
		delete(annotations, customresource.ResourcePolicyAnnotation)
		data.Resources.Deployments[0].ObjectMeta.SetAnnotations(annotations)
		UpdateDeployment(data)
	})
}
