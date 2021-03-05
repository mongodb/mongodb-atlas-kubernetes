package e2e_test

import (
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"go.mongodb.org/atlas/mongodbatlas"

	kube "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kube"
	mongocli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/mongocli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

var (
	ConfigAll         = "../../deploy/" // Released generated files
	ClusterSample     = "data/atlascluster_basic.yaml"
	DataFolder        = "data/"
	defaultOperatorNS = "mongodb-atlas-system"
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
	).Should(Equal("True"), "Kubernetes resource: Cluster status `Ready` should be True")

	Eventually(kube.GetK8sClusterStateName(
		input.namespace, input.clusters[0].GetClusterNameResource()),
		"45m", "1m",
	).Should(Equal("IDLE"), "Kubernetes resource: Cluster status should be IDLE")

	Expect(
		mongocli.GetClusterStateName(input.projectID, input.clusters[0].Spec.Name),
	).Should(Equal("IDLE"), "Atlas: Cluster status should be IDLE")
}

func waitProject(input userInputs, generation string) {
	EventuallyWithOffset(1, kube.GetStatusCondition(input.namespace, input.k8sFullProjectName)).Should(Equal("True"), "Kubernetes resource: Project status `Ready` should be True")
	EventuallyWithOffset(1, kube.GetGeneration(input.namespace, input.k8sFullProjectName)).Should(Equal(generation), "Kubernetes resource: Generation should be upgraded")
	EventuallyWithOffset(1, kube.GetProjectResource(input.namespace, input.k8sFullProjectName).Status.ID).ShouldNot(BeNil(), "Kubernetes resource: Status has field with ProjectID")
}

func checkIfClusterExist(input userInputs) func() bool {
	return func() bool {
		return mongocli.IsClusterExist(input.projectID, input.clusters[0].Spec.Name)
	}
}

func compareClustersSpec(requested utils.ClusterSpec, created mongodbatlas.Cluster) { // TODO
	ExpectWithOffset(1, created).To(MatchFields(IgnoreExtras, Fields{
		"MongoURI":            Not(BeEmpty()),
		"MongoURIWithOptions": Not(BeEmpty()),
		"Name":                Equal(requested.Name),
		"ProviderSettings": PointTo(MatchFields(IgnoreExtras, Fields{
			"InstanceSizeName": Equal(requested.ProviderSettings.InstanceSizeName),
			"ProviderName":     Equal(string(requested.ProviderSettings.ProviderName)),
			"RegionName":       Equal(requested.ProviderSettings.RegionName),
		})),
		"ConnectionStrings": PointTo(MatchFields(IgnoreExtras, Fields{
			"Standard":    Not(BeEmpty()),
			"StandardSrv": Not(BeEmpty()),
		})),
	}), "Cluster should be the same as requested by the user")
}

func SaveK8sResources(resources []string, ns string) {
	for _, resource := range resources {
		data := kube.GetYamlResource(resource, ns)
		utils.SaveToFile("output/"+resource+".yaml", data)
	}
}
