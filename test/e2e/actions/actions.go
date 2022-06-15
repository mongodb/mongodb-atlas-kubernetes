// `actions` additional functions which accept testDataProvider struct and could be used as additional acctions in the tests since they all typical

package actions

import (
	"fmt"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	appclient "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/appclient"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	mongocli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/mongocli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

func UpdateDeployment(newData *model.TestDataProvider) {
	var generation int
	By("Update Deployment\n", func() {
		utils.SaveToFile(
			newData.Resources.Deployments[0].DeploymentFileName(newData.Resources),
			utils.JSONToYAMLConvert(newData.Resources.Deployments[0]),
		)
		generation, _ = strconv.Atoi(kubecli.GetGeneration(newData.Resources.Namespace, newData.Resources.Deployments[0].GetDeploymentNameResource()))
		kubecli.Apply(newData.Resources.Deployments[0].DeploymentFileName(newData.Resources), "-n", newData.Resources.Namespace)
		generation++
	})

	By("Wait Deployment updating\n", func() {
		WaitDeployment(newData.Resources, strconv.Itoa(generation))
	})

	By("Check attributes\n", func() {
		uDeployment := mongocli.GetDeploymentsInfo(newData.Resources.ProjectID, newData.Resources.Deployments[0].Spec.GetDeploymentName())
		CompareDeploymentsSpec(newData.Resources.Deployments[0].Spec, uDeployment)
	})
}

func UpdateDeploymentFromUpdateConfig(data *model.TestDataProvider) {
	By("Load new Deployment config", func() {
		data.Resources.Deployments = []model.AtlasDeployment{} // TODO for range
		GinkgoWriter.Write([]byte(data.ConfUpdatePaths[0]))
		data.Resources.Deployments = append(data.Resources.Deployments, model.LoadUserDeploymentConfig(data.ConfUpdatePaths[0]))
		data.Resources.Deployments[0].Spec.Project.Name = data.Resources.Project.GetK8sMetaName()
		utils.SaveToFile(
			data.Resources.Deployments[0].DeploymentFileName(data.Resources),
			utils.JSONToYAMLConvert(data.Resources.Deployments[0]),
		)
	})

	UpdateDeployment(data)

	By("Check user data still in the Deployment\n", func() {
		for i := range data.Resources.Users { // TODO in parallel(?)
			port := strconv.Itoa(i + data.PortGroup)
			key := port
			data := fmt.Sprintf("{\"key\":\"%s\",\"shipmodel\":\"heavy\",\"hp\":150}", key)
			app := appclient.NewTestAppClient(port)
			Expect(app.Get("/mongo/" + key)).Should(Equal(data))
		}
	})
}

func activateDeployment(data *model.TestDataProvider, paused bool) {
	data.Resources.Deployments[0].Spec.DeploymentSpec.Paused = &paused
	UpdateDeployment(data)
	By("Check additional Deployment field `paused`\n")
	uDeployment := mongocli.GetDeploymentsInfo(data.Resources.ProjectID, data.Resources.Deployments[0].Spec.GetDeploymentName())
	Expect(uDeployment.Paused).Should(Equal(data.Resources.Deployments[0].Spec.DeploymentSpec.Paused))
}

func SuspendDeployment(data *model.TestDataProvider) {
	paused := true
	activateDeployment(data, paused)
}

func ReactivateDeployment(data *model.TestDataProvider) {
	paused := false
	activateDeployment(data, paused)
}

func DeleteFirstUser(data *model.TestDataProvider) {
	By("User can delete Database User", func() {
		// data.Resources.ProjectID = kube.GetProjectResource(data.Resources.Namespace, data.Resources.GetAtlasProjectFullKubeName()).Status.ID
		// since it is could be several users, we should
		// - delete k8s resource
		// - delete one user from the list,
		// - check Atlas doesn't have the initial user and have the rest
		By("Delete k8s resources")
		Eventually(kubecli.Delete(data.Resources.GetResourceFolder()+"/user/user-"+data.Resources.Users[0].ObjectMeta.Name+".yaml", "-n", data.Resources.Namespace)).Should(Say("deleted"))
		Eventually(CheckIfUserExist(data.Resources.Users[0].Spec.Username, data.Resources.ProjectID)).Should(BeFalse())

		// the rest users should be still there
		data.Resources.Users = data.Resources.Users[1:]
		Eventually(CheckIfUsersExist(data.Resources), "2m", "10s").Should(BeTrue())
		CheckUsersAttributes(data.Resources)
	})
}
