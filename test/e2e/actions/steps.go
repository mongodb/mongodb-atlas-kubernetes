package actions

import (
	"fmt"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gstruct"
	"go.mongodb.org/atlas/mongodbatlas"

	appclient "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/appclient"
	helm "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/helm"
	kube "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kube"
	mongocli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/mongocli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

func WaitCluster(input model.UserInputs, generation string) {
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

func WaitProject(input model.UserInputs, generation string) {
	EventuallyWithOffset(1, kube.GetStatusCondition(input.Namespace, input.K8sFullProjectName)).Should(Equal("True"), "Kubernetes resource: Project status `Ready` should be True")
	ExpectWithOffset(1, kube.GetGeneration(input.Namespace, input.K8sFullProjectName)).Should(Equal(generation), "Kubernetes resource: Generation should be upgraded")
	ExpectWithOffset(1, kube.GetProjectResource(input.Namespace, input.K8sFullProjectName).Status.ID).ShouldNot(BeNil(), "Kubernetes resource: Status has field with ProjectID")
}

func WaitTestApplication(ns, label string) {
	// temp
	isAppRunning := func() func() bool {
		return func() bool {
			status := kube.GetStatusPhase(ns, "pods", "-l", label)
			if status == "Running" {
				return true
			}
			kube.DescribeTestApp(label, ns)
			return false
		}
	}
	EventuallyWithOffset(1, isAppRunning(), "2m", "10s").Should(BeTrue(), "Test application should be running")
}

func CheckIfClusterExist(input model.UserInputs) func() bool {
	return func() bool {
		return mongocli.IsClusterExist(input.ProjectID, input.Clusters[0].Spec.Name)
	}
}

func CheckIfUsersExist(input model.UserInputs) func() bool {
	return func() bool {
		for _, user := range input.Users {
			if !mongocli.IsUserExist(user.Spec.Username, input.ProjectID) {
				return false
			}
		}
		return true
	}
}

func CheckIfUserExist(username, projecID string) func() bool {
	return func() bool {
		return mongocli.IsUserExist(username, projecID)
	}
}

func CompareClustersSpec(requested model.ClusterSpec, created mongodbatlas.Cluster) { // TODO
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
	if requested.ProviderSettings.ProviderName == "TENANT" {
		ExpectWithOffset(1, requested.ProviderSettings).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"BackingProviderName": Equal(created.ProviderSettings.BackingProviderName),
		})), "Cluster should be the same as requested by the user: Backking Provider Name")
	}
}

func SaveK8sResources(resources []string, ns string) {
	for _, resource := range resources {
		data := kube.GetYamlResource(resource, ns)
		utils.SaveToFile("output/"+resource+".yaml", data)
	}
}

func SaveTestAppLogs(input model.UserInputs) {
	for _, user := range input.Users {
		utils.SaveToFile(
			fmt.Sprintf("output/testapp-describe-%s.txt", user.Spec.Username),
			kube.DescribeTestApp(config.TestAppLabelPrefix+user.Spec.Username, input.Namespace),
		)
		utils.SaveToFile(
			fmt.Sprintf("output/testapp-logs-%s.txt", user.Spec.Username),
			kube.GetTestAppLogs(config.TestAppLabelPrefix+user.Spec.Username, input.Namespace),
		)
	}
}

func CheckUsersAttributes(input model.UserInputs) {
	for _, cluster := range input.Clusters {
		for _, user := range input.Users {
			EventuallyWithOffset(1, mongocli.IsUserExist(user.Spec.Username, input.ProjectID), "7m", "10s").Should(BeTrue())
			uResourceName := fmt.Sprintf("atlasdatabaseusers.atlas.mongodb.com/%s-%s", cluster.ObjectMeta.Name, user.Spec.Username)
			EventuallyWithOffset(
				1, kube.GetStatusCondition(input.Namespace, uResourceName),
				"45m", "1m",
			).Should(Equal("True"), "Kubernetes resource: User resources status `Ready` should be True")

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
}

func CheckUsersCanUseApplication(portGroup int, userSpec model.UserInputs) {
	for i, user := range userSpec.Users { // TODO in parallel(?)/ingress
		// data
		port := strconv.Itoa(i + portGroup)
		key := port
		data := fmt.Sprintf("{\"key\":\"%s\",\"shipmodel\":\"heavy\",\"hp\":150}", key)

		helm.InstallTestApplication(userSpec, user, port)
		WaitTestApplication(userSpec.Namespace, "app=test-app-"+user.Spec.Username)

		app := appclient.NewTestAppClient(port)
		ExpectWithOffset(1, app.Get("")).Should(Equal("It is working"))
		ExpectWithOffset(1, app.Post(data)).ShouldNot(HaveOccurred())
		ExpectWithOffset(1, app.Get("/mongo/"+key)).Should(Equal(data))
	}
}

func PrepareUsersConfigurations(data *model.TestDataProvider) {
	By("Prepare namespaces and project configuration", func() {
		kube.CreateNamespace(data.Resources.Namespace)
		By("Create project spec", func() {
			GinkgoWriter.Write([]byte(data.Resources.ProjectPath + "\n"))
			utils.SaveToFile(data.Resources.ProjectPath, data.Resources.Project.ConvertByte())
		})
		By("Create cluster spec", func() {
			data.Resources.Clusters[0].Spec.Project.Name = data.Resources.Project.GetK8sMetaName()
			utils.SaveToFile(
				data.Resources.Clusters[0].ClusterFileName(data.Resources),
				utils.JSONToYAMLConvert(data.Resources.Clusters[0]),
			)
		})
		By("Create dbuser spec", func() {
			Expect(data.Resources.Users).ShouldNot(BeNil())
			for _, user := range data.Resources.Users {
				user.SaveConfigurationTo(data.Resources.ProjectPath)
				kube.CreateUserSecret(user.Spec.PasswordSecret.Name, data.Resources.Namespace)
			}
		})
	})
}

func DeployUserResourcesAction(data *model.TestDataProvider) {
	By("Create users resources", func() {
		kube.CreateApiKeySecret(data.Resources.KeyName, data.Resources.Namespace)
		kube.Apply(data.Resources.ProjectPath, "-n", data.Resources.Namespace)
		kube.Apply(data.Resources.Clusters[0].ClusterFileName(data.Resources), "-n", data.Resources.Namespace)
		kube.Apply(data.Resources.GetResourceFolder()+"/user/", "-n", data.Resources.Namespace)
	})

	By("Wait project creation", func() {
		WaitProject(data.Resources, "1")
		data.Resources.ProjectID = kube.GetProjectResource(data.Resources.Namespace, data.Resources.K8sFullProjectName).Status.ID
		Expect(data.Resources.ProjectID).ShouldNot(BeEmpty())
	})

	By("Wait cluster creation", func() {
		WaitCluster(data.Resources, "1")
	})

	By("check cluster Attribute", func() {
		cluster := mongocli.GetClustersInfo(data.Resources.ProjectID, data.Resources.Clusters[0].Spec.Name)
		CompareClustersSpec(data.Resources.Clusters[0].Spec, cluster)
	})

	By("check database users Attibutes", func() {
		Eventually(CheckIfUsersExist(data.Resources), "2m", "10s").Should(BeTrue())
		CheckUsersAttributes(data.Resources)
	})

	By("Deploy application for user", func() {
		CheckUsersCanUseApplication(data.PortGroup, data.Resources)
	})
}

func DeleteDBUsersApps(data *model.TestDataProvider) {
	By("Delete dbusers applications", func() {
		for _, user := range data.Resources.Users {
			helm.Uninstall("test-app-"+user.Spec.Username, data.Resources.Namespace)
		}
	})
}

func DeleteUserResources(data *model.TestDataProvider) {
	By("Delete cluster", func() {
		kube.Delete(data.Resources.Clusters[0].ClusterFileName(data.Resources), "-n", data.Resources.Namespace)
		Eventually(
			CheckIfClusterExist(data.Resources),
			"10m", "1m",
		).Should(BeFalse(), "Cluster should be deleted from Atlas")
	})

	By("Delete project", func() {
		kube.Delete(data.Resources.ProjectPath, "-n", data.Resources.Namespace)
		Eventually(
			func() bool {
				return mongocli.IsProjectInfoExist(data.Resources.ProjectID)
			},
			"5m", "20s",
		).Should(BeFalse(), "Project should be deleted from Atlas")
	})
}

func AfterEachFinalCleanup(datas []model.TestDataProvider) {
	for i := range datas {
		GinkgoWriter.Write([]byte("AfterEach. Final cleanup...\n"))
		DeleteDBUsersApps(&datas[i])
		Expect(kube.DeleteNamespace(datas[i].Resources.Namespace)).Should(Say("deleted"), "Cant delete namespace after testing")
		GinkgoWriter.Write([]byte("AfterEach. Cleanup finished\n"))
	}
}
