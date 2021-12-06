// `actions` additional functions which accept testDataProvider struct and could be used as additional acctions in the tests since they all typical

package actions

import (
	"fmt"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	appclient "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/appclient"
	kube "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kube"
	mongocli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/mongocli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

func UpdateCluster(newData *model.TestDataProvider) {
	var generation int
	By("Update cluster\n", func() {
		utils.SaveToFile(
			newData.Resources.Clusters[0].ClusterFileName(newData.Resources),
			utils.JSONToYAMLConvert(newData.Resources.Clusters[0]),
		)
		generation, _ = strconv.Atoi(kube.GetGeneration(newData.Resources.Namespace, newData.Resources.Clusters[0].GetClusterNameResource()))
		kube.Apply(newData.Resources.Clusters[0].ClusterFileName(newData.Resources), "-n", newData.Resources.Namespace)
		generation++
	})

	By("Wait cluster updating\n", func() {
		WaitCluster(newData.Resources, strconv.Itoa(generation))
	})

	By("Check attributes\n", func() {
		uCluster := mongocli.GetClustersInfo(newData.Resources.ProjectID, newData.Resources.Clusters[0].Spec.Name)
		CompareClustersSpec(newData.Resources.Clusters[0].Spec, uCluster)
	})
}

func UpdateClusterFromUpdateConfig(data *model.TestDataProvider) {
	By("Load new cluster config", func() {
		data.Resources.Clusters = []model.AC{} // TODO for range
		GinkgoWriter.Write([]byte(data.ConfUpdatePaths[0]))
		data.Resources.Clusters = append(data.Resources.Clusters, model.LoadUserClusterConfig(data.ConfUpdatePaths[0]))
		data.Resources.Clusters[0].Spec.Project.Name = data.Resources.Project.GetK8sMetaName()
		utils.SaveToFile(
			data.Resources.Clusters[0].ClusterFileName(data.Resources),
			utils.JSONToYAMLConvert(data.Resources.Clusters[0]),
		)
	})

	UpdateCluster(data)

	By("Check user data still in the cluster\n", func() {
		for i := range data.Resources.Users { // TODO in parallel(?)
			port := strconv.Itoa(i + data.PortGroup)
			key := port
			data := fmt.Sprintf("{\"key\":\"%s\",\"shipmodel\":\"heavy\",\"hp\":150}", key)
			app := appclient.NewTestAppClient(port)
			Expect(app.Get("/mongo/" + key)).Should(Equal(data))
		}
	})
}

func activateCluster(data *model.TestDataProvider, paused bool) {
	data.Resources.Clusters[0].Spec.Paused = &paused
	UpdateCluster(data)
	By("Check additional cluster field `paused`\n")
	uCluster := mongocli.GetClustersInfo(data.Resources.ProjectID, data.Resources.Clusters[0].Spec.Name)
	Expect(uCluster.Paused).Should(Equal(data.Resources.Clusters[0].Spec.Paused))
}

func SuspendCluster(data *model.TestDataProvider) {
	paused := true
	activateCluster(data, paused)
}

func ReactivateCluster(data *model.TestDataProvider) {
	paused := false
	activateCluster(data, paused)
}

func DeleteFirstUser(data *model.TestDataProvider) {
	By("User can delete Database User", func() {
		// data.Resources.ProjectID = kube.GetProjectResource(data.Resources.Namespace, data.Resources.GetAtlasProjectFullKubeName()).Status.ID
		// since it is could be several users, we should
		// - delete k8s resource
		// - delete one user from the list,
		// - check Atlas doesn't have the initial user and have the rest
		By("Delete k8s resources")
		Eventually(kube.Delete(data.Resources.GetResourceFolder()+"/user/user-"+data.Resources.Users[0].ObjectMeta.Name+".yaml", "-n", data.Resources.Namespace)).Should(Say("deleted"))
		Eventually(CheckIfUserExist(data.Resources.Users[0].Spec.Username, data.Resources.ProjectID)).Should(BeFalse())

		// the rest users should be still there
		data.Resources.Users = data.Resources.Users[1:]
		Eventually(CheckIfUsersExist(data.Resources), "2m", "10s").Should(BeTrue())
		CheckUsersAttributes(data.Resources)
	})
}
