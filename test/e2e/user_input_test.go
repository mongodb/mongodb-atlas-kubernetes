package e2e_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"go.mongodb.org/atlas/mongodbatlas"

	appclient "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/appclient"
	helm "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/helm"
	kube "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kube"
	mongocli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/mongocli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

func waitCluster(input model.UserInputs, generation string) {
	EventuallyWithOffset(
		1, kube.GetStatusCondition(input.Namespace, input.Clusters[0].GetClusterNameResource()),
		"45m", "1m",
	).Should(Equal("True"), "Kubernetes resource: Cluster status `Ready` should be True")

	ExpectWithOffset(
		1, kube.GetGeneration(input.Namespace, input.Clusters[0].GetClusterNameResource()),
	).Should(Equal(generation))
	ExpectWithOffset(1, kube.GetK8sClusterStateName(
		input.Namespace, input.Clusters[0].GetClusterNameResource()),
	).Should(Equal("IDLE"), "Kubernetes resource: Cluster status should be IDLE")

	ExpectWithOffset(
		1, mongocli.GetClusterStateName(input.ProjectID, input.Clusters[0].Spec.Name),
	).Should(Equal("IDLE"), "Atlas: Cluster status should be IDLE")
}

func waitProject(input model.UserInputs, generation string) { //nolint:unparam // have cases only with generation=1
	EventuallyWithOffset(1, kube.GetStatusCondition(input.Namespace, input.K8sFullProjectName)).Should(Equal("True"), "Kubernetes resource: Project status `Ready` should be True")
	ExpectWithOffset(1, kube.GetGeneration(input.Namespace, input.K8sFullProjectName)).Should(Equal(generation), "Kubernetes resource: Generation should be upgraded")
	ExpectWithOffset(1, kube.GetProjectResource(input.Namespace, input.K8sFullProjectName).Status.ID).ShouldNot(BeNil(), "Kubernetes resource: Status has field with ProjectID")
}

func waitTestApplication(ns, label string) {
	EventuallyWithOffset(1, kube.GetStatusPhase(ns, "pods", "-l", label)).Should(Equal("Running"), "Test application should be running")
}

func checkIfClusterExist(input model.UserInputs) func() bool {
	return func() bool {
		return mongocli.IsClusterExist(input.ProjectID, input.Clusters[0].Spec.Name)
	}
}

func checkIfUsersExist(input model.UserInputs) func() bool {
	return func() bool {
		for _, user := range input.Users {
			if !mongocli.IsUserExist(user.Spec.Username, input.ProjectID) {
				return false
			}
		}
		return true
	}
}

func checkIfUserExist(username, projecID string) func() bool {
	return func() bool {
		return mongocli.IsUserExist(username, projecID)
	}
}

func compareClustersSpec(requested model.ClusterSpec, created mongodbatlas.Cluster) { // TODO
	ExpectWithOffset(1, created).To(MatchFields(IgnoreExtras, Fields{
		"MongoURI":            Not(BeEmpty()),
		"MongoURIWithOptions": Not(BeEmpty()),
		"Name":                Equal(requested.Name),
		"ProviderSettings": PointTo(MatchFields(IgnoreExtras, Fields{
			"InstanceSizeName": Equal(requested.ProviderSettings.InstanceSizeName),
			"ProviderName":     Equal(string(requested.ProviderSettings.ProviderName)),
		})),
		"ConnectionStrings": PointTo(MatchFields(IgnoreExtras, Fields{
			"Standard":    Not(BeEmpty()),
			"StandardSrv": Not(BeEmpty()),
		})),
	}), "Cluster should be the same as requested by the user")

	if len(requested.ReplicationSpecs) > 0 {
		for i, replica := range requested.ReplicationSpecs {
			for key, region := range replica.RegionsConfig {
				// diffent type
				ExpectWithOffset(1, created.ReplicationSpecs[i].RegionsConfig[key].AnalyticsNodes).Should(PointTo(Equal(*region.AnalyticsNodes)), "Replica Spec: AnalyticsNodes is not the same")
				ExpectWithOffset(1, created.ReplicationSpecs[i].RegionsConfig[key].ElectableNodes).Should(PointTo(Equal(*region.ElectableNodes)), "Replica Spec: ElectableNodes is not the same")
				ExpectWithOffset(1, created.ReplicationSpecs[i].RegionsConfig[key].Priority).Should(PointTo(Equal(*region.Priority)), "Replica Spec: Priority is not the same")
				ExpectWithOffset(1, created.ReplicationSpecs[i].RegionsConfig[key].ReadOnlyNodes).Should(PointTo(Equal(*region.ReadOnlyNodes)), "Replica Spec: ReadOnlyNodes is not the same")
			}
		}
	} else {
		ExpectWithOffset(1, requested.ProviderSettings).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"RegionName": Equal(created.ProviderSettings.RegionName),
		})), "Cluster should be the same as requested by the user: Region Name")
	}
}

func SaveK8sResources(resources []string, ns string) {
	for _, resource := range resources {
		data := kube.GetYamlResource(resource, ns)
		utils.SaveToFile("output/"+resource+".yaml", data)
	}
}

func checkUsersAttributes(input model.UserInputs) {
	for _, user := range input.Users {
		atlasUser := mongocli.GetUser(user.Spec.Username, input.ProjectID)
		// Required fields
		ExpectWithOffset(1, atlasUser).To(MatchFields(IgnoreExtras, Fields{
			"Username":     Equal(user.Spec.Username),
			"GroupID":      Equal(input.ProjectID),
			"DatabaseName": Or(Equal(user.Spec.DatabaseName), Equal("admin")),
		}), "Users attributes should be the same as requested by the user")

		for i, role := range atlasUser.Roles {
			ExpectWithOffset(1, role).To(MatchFields(IgnoreMissing, Fields{
				"RoleName":       Equal(user.Spec.Roles[i].RoleName),
				"DatabaseName":   Equal(user.Spec.Roles[i].DatabaseName),
				"CollectionName": Equal(user.Spec.Roles[i].CollectionName),
			}), "Users roles attributes should be the same as requsted by the user")
		}
	}
}

// CopyKustomizeNamespaceOperator create copy of `/deploy/namespaced` folder with kustomization file for overriding namespace
func CopyKustomizeNamespaceOperator(input model.UserInputs) {
	fullPath := input.GetOperatorFolder()
	os.Mkdir(fullPath, os.ModePerm)
	utils.CopyFile("../../deploy/namespaced/crds.yaml", filepath.Join(fullPath, "crds.yaml"))
	utils.CopyFile("../../deploy/namespaced/namespaced-config.yaml", filepath.Join(fullPath, "namespaced-config.yaml"))
	data := []byte(
		"namespace: " + input.Namespace + "\n" +
			"resources:" + "\n" +
			"- crds.yaml" + "\n" +
			"- namespaced-config.yaml",
	)
	utils.SaveToFile(filepath.Join(fullPath, "kustomization.yaml"), data)
}

func checkUsersCanUseApplication(portGroup int, userSpec model.UserInputs) {
	for i, user := range userSpec.Users { // TODO in parallel(?)/ingress
		// data
		port := strconv.Itoa(i + portGroup)
		key := port
		data := fmt.Sprintf("{\"key\":\"%s\",\"shipmodel\":\"heavy\",\"hp\":150}", key)

		helm.InstallTestApplication(userSpec, user, port)
		waitTestApplication(userSpec.Namespace, "app=test-app-"+user.Spec.Username)

		app := appclient.NewTestAppClient(port)
		ExpectWithOffset(1, app.Get("")).Should(Equal("It is working"))
		ExpectWithOffset(1, app.Post(data)).ShouldNot(HaveOccurred())
		ExpectWithOffset(1, app.Get("/mongo/"+key)).Should(Equal(data))
	}
}
