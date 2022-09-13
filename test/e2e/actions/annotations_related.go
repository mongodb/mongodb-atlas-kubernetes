package actions

import (
	"fmt"
	"time"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

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
		err := kubecli.DeleteDeployment(data.Context, data.K8SClient, data.Resources.Namespace, data.Resources.Deployments[0].ObjectMeta.GetName())
		Expect(err).NotTo(HaveOccurred())
		By("Checking Cluster still existed", func() {
			aClient := atlas.GetClientOrFail()
			state := aClient.GetDeployment(data.Resources.ProjectID, data.Resources.Deployments[0].Spec.GetDeploymentName()).StateName
			Expect(state).ShouldNot(Equal("DELETING"), "Deployment is being deleted despite the keep annotation")
		})
		By("Checking CR not exists", func() {
			Eventually(func() bool {
				_, err = kubecli.GetDeploymentResource(data.Context, data.K8SClient, data.Resources.Namespace, data.Resources.Deployments[0].ObjectMeta.GetName())
				By(fmt.Sprintf("CR still exists: %v", err))
				return k8serrors.IsNotFound(err)
			}).WithTimeout(5*time.Minute).WithPolling(10*time.Second).Should(BeTrue(), "CR still exists")
		})
	})
}

func ReDeployOperator(data *model.TestDataProvider) {
	By(fmt.Sprintf("Recreating %s", atlasDeploymentCRD), func() {
		deploy.NamespacedOperator(data)
	})
}
