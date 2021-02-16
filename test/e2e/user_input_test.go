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
	projectName string
	projectID   string
	keyName     string
	clusters    []utils.AC
}

func (t *userInputs) GenNamespace() string {
	return "ns-" + t.projectName
	// return "mongodb-atlas-kubernetes-system"
}

func (t *userInputs) GenProjectSample() string {
	return DataFolder + t.projectName + ".yaml"
}

func (t *userInputs) ProjectK8sName() string {
	return "k-" + t.projectName
}

func (t *userInputs) UserProjectFile() string {
	return DataFolder + t.projectName + ".yaml"
}

func (t *userInputs) GetFullK8sAtlasProjectName() string {
	return "atlasproject.atlas.mongodb.com/" + t.ProjectK8sName()
}

func FilePathTo(name string) string {
	return DataFolder + name + ".yaml"
}

func waitCluster(input userInputs, generation string) {
	Eventually(kube.GetGeneration(input.GenNamespace(), input.clusters[0].GetClusterNameResource())).Should(Equal(generation))
	Eventually(
		kube.GetStatusCondition(input.GenNamespace(), input.clusters[0].GetClusterNameResource()),
		"45m", "1m",
	).Should(Equal("True"))

	Eventually(kube.GetK8sClusterStateName(
		input.GenNamespace(), input.clusters[0].GetClusterNameResource()),
		"45m", "1m",
	).Should(Equal("IDLE"))

	Expect(
		mongocli.GetClusterStateName(input.projectID, input.clusters[0].Spec.Name),
	).Should(Equal("IDLE"))
}

func waitProject(input userInputs, generation string) {
	EventuallyWithOffset(1, kube.GetStatusCondition(input.GenNamespace(), "atlasproject.atlas.mongodb.com/"+input.ProjectK8sName())).Should(Equal("True"))
	EventuallyWithOffset(1, kube.GetGeneration(input.GenNamespace(), "atlasproject.atlas.mongodb.com/"+input.ProjectK8sName())).Should(Equal(generation))
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
