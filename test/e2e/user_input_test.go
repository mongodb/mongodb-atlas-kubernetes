package e2e_test

import (
	. "github.com/onsi/gomega"

	kube "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kube"
	mongocli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/mongocli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

var (
	ConfigAll         = "../../deploy/" // Released generated files
	ClusterSample     = "data/atlascluster_basic.yaml"
	DataFolder        = "data/"
	defaultOperatorNS = "mongodb-atlas-kubernetes-system"
)

type userInputs struct {
	projectName        string
	projectID          string
	keyName            string
	namespace          string
	k8sProjectName     string
	k8sFullProjectName string
	projectPath        string
	clusters           []utils.AC
}

func NewUserInputs(keyName string) userInputs {
	uid := utils.GenUniqID()
	return userInputs{
		projectName:        uid,
		projectID:          "",
		keyName:            keyName,
		namespace:          "ns-" + uid,
		k8sProjectName:     "k-" + uid,
		k8sFullProjectName: "atlasproject.atlas.mongodb.com/k-" + uid,
		projectPath:        DataFolder + uid + ".yaml",
	}
}

func FilePathTo(name string) string {
	return DataFolder + name + ".yaml"
}

func waitCluster(input userInputs, generation string) {
	Eventually(kube.GetGeneration(input.namespace, input.clusters[0].GetClusterNameResource())).Should(Equal(generation))
	Eventually(
		kube.GetStatusCondition(input.namespace, input.clusters[0].GetClusterNameResource()),
		"45m", "1m",
	).Should(Equal("True"))

	Eventually(kube.GetK8sClusterStateName(
		input.namespace, input.clusters[0].GetClusterNameResource()),
		"45m", "1m",
	).Should(Equal("IDLE"))

	Expect(
		mongocli.GetClusterStateName(input.projectID, input.clusters[0].Spec.Name),
	).Should(Equal("IDLE"))
}

func waitProject(input userInputs, generation string) {
	EventuallyWithOffset(1, kube.GetStatusCondition(input.namespace, input.k8sFullProjectName)).Should(Equal("True"))
	EventuallyWithOffset(1, kube.GetGeneration(input.namespace, input.k8sFullProjectName)).Should(Equal(generation))
	ExpectWithOffset(1,
		mongocli.IsProjectExist(input.projectName),
	).Should(BeTrue())
}

// func checkIfProjectExist(input userInputs) func() bool {
// 	return func() bool {
// 		return mongocli.IsProjectExist(input.projectName)
// 	}
// }

func checkIfClusterExist(input userInputs) func() bool {
	return func() bool {
		return mongocli.IsClusterExist(input.projectID, input.projectName)
	}
}
