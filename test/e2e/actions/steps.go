package actions

import (
	"fmt"
	"os"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gstruct"
	"go.mongodb.org/atlas/mongodbatlas"

	kube "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/atlas"
	appclient "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/appclient"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/helm"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/mongocli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

func WaitDeployment(input model.UserInputs, generation string) {
	EventuallyWithOffset(1,
		func() string {
			return kubecli.GetGeneration(input.Namespace, input.Deployments[0].GetDeploymentNameResource())
		},
		"5m", "10s",
	).Should(Equal(generation))

	EventuallyWithOffset(1,
		func() string {
			return kubecli.GetStatusCondition("Ready", input.Namespace, input.Deployments[0].GetDeploymentNameResource())
		},
		"60m", "1m",
	).Should(Equal("True"), "Kubernetes resource: Deployment status `Ready` should be 'True'")

	ExpectWithOffset(1, kubecli.GetK8sDeploymentStateName(
		input.Namespace, input.Deployments[0].GetDeploymentNameResource()),
	).Should(Equal("IDLE"), "Kubernetes resource: Deployment status should be IDLE")

	deployment := input.Deployments[0]
	switch {
	case deployment.Spec.AdvancedDeploymentSpec != nil:
		atlasClient, err := atlas.AClient()
		Expect(err).To(BeNil())
		advancedDeployment, err := atlasClient.GetAdvancedDeployment(input.ProjectID, deployment.Spec.AdvancedDeploymentSpec.Name)
		Expect(err).To(BeNil())
		Expect(advancedDeployment.StateName).To(Equal("IDLE"))
	case deployment.Spec.ServerlessSpec != nil:
		atlasClient, err := atlas.AClient()
		Expect(err).To(BeNil())
		serverlessInstance, err := atlasClient.GetServerlessInstance(input.ProjectID, deployment.Spec.ServerlessSpec.Name)
		Expect(err).To(BeNil())
		Expect(serverlessInstance.StateName).To(Equal("IDLE"))
	default:
		ExpectWithOffset(
			1, mongocli.GetDeploymentStateName(input.ProjectID, input.Deployments[0].Spec.GetDeploymentName()),
		).Should(Equal("IDLE"), "Atlas: Deployment status should be IDLE")
	}
}

func WaitProject(data *model.TestDataProvider, generation string) {
	EventuallyWithOffset(1, kube.GetReadyProjectStatus(data), "15m", "10s").Should(Equal("True"), "Kubernetes resource: Project status `Ready` should be 'True'")
	ExpectWithOffset(1, kubecli.GetGeneration(data.Resources.Namespace, data.Resources.GetAtlasProjectFullKubeName())).Should(Equal(generation), "Kubernetes resource: Generation should be upgraded")
	atlasProject, err := kube.GetProjectResource(data)
	Expect(err).ShouldNot(HaveOccurred())
	ExpectWithOffset(1, atlasProject.Status.ID).ShouldNot(BeNil(), "Kubernetes resource: Status has field with ProjectID")
}

func WaitProjectWithoutGenerationCheck(data *model.TestDataProvider) {
	EventuallyWithOffset(1, kube.GetReadyProjectStatus(data), "15m", "10s").Should(Equal("True"), "Kubernetes resource: Project status `Ready` should be 'True'")
	atlasProject, err := kube.GetProjectResource(data)
	Expect(err).ShouldNot(HaveOccurred())
	ExpectWithOffset(1, atlasProject.Status.ID).ShouldNot(BeNil(), "Kubernetes resource: Status has field with ProjectID")
}

func WaitTestApplication(ns, label string) {
	// temp
	isAppRunning := func() func() bool {
		return func() bool {
			status := kubecli.GetStatusPhase(ns, "pods", "-l", label)
			if status == "Running" {
				return true
			}
			kubecli.DescribeTestApp(label, ns)
			return false
		}
	}
	EventuallyWithOffset(1, isAppRunning(), "2m", "10s").Should(BeTrue(), "Test application should be running")
}

func CheckIfDeploymentExist(input model.UserInputs) func() bool {
	return func() bool {
		return mongocli.IsDeploymentExist(input.ProjectID, input.Deployments[0].Spec.DeploymentSpec.Name)
	}
}

func CheckIfUsersExist(input model.UserInputs) func() bool {
	return func() bool {
		atlasClient, err := atlas.AClient()
		if err != nil {
			return false
		}

		for _, user := range input.Users {
			dbUser, err := atlasClient.GetUserByName(setAdminIfEmpty(user.Spec.DatabaseName), input.ProjectID, user.Spec.Username)
			if err != nil && dbUser == nil {
				return false
			}
		}
		return true
	}
}

func CheckIfUserExist(username, projecID string) func() bool {
	return func() bool {
		return mongocli.IsUserExistForAdmin(username, projecID)
	}
}

func CompareDeploymentsSpec(requested model.DeploymentSpec, created mongodbatlas.Cluster) {
	ExpectWithOffset(1, created).To(MatchFields(IgnoreExtras, Fields{
		"MongoURI":            Not(BeEmpty()),
		"MongoURIWithOptions": Not(BeEmpty()),
		"Name":                Equal(requested.DeploymentSpec.Name),
		"ProviderSettings": PointTo(MatchFields(IgnoreExtras, Fields{
			"InstanceSizeName": Equal(requested.DeploymentSpec.ProviderSettings.InstanceSizeName),
			"ProviderName":     Equal(string(requested.DeploymentSpec.ProviderSettings.ProviderName)),
		})),
		"ConnectionStrings": PointTo(MatchFields(IgnoreExtras, Fields{
			"Standard":    Not(BeEmpty()),
			"StandardSrv": Not(BeEmpty()),
		})),
	}), "Deployment should be the same as requested by the user")

	if len(requested.DeploymentSpec.ReplicationSpecs) > 0 {
		for i, replica := range requested.DeploymentSpec.ReplicationSpecs {
			for key, region := range replica.RegionsConfig {
				// diffent type
				ExpectWithOffset(1, created.ReplicationSpecs[i].RegionsConfig[key].AnalyticsNodes).Should(PointTo(Equal(*region.AnalyticsNodes)), "Replica Spec: AnalyticsNodes is not the same")
				ExpectWithOffset(1, created.ReplicationSpecs[i].RegionsConfig[key].ElectableNodes).Should(PointTo(Equal(*region.ElectableNodes)), "Replica Spec: ElectableNodes is not the same")
				ExpectWithOffset(1, created.ReplicationSpecs[i].RegionsConfig[key].Priority).Should(PointTo(Equal(*region.Priority)), "Replica Spec: Priority is not the same")
				ExpectWithOffset(1, created.ReplicationSpecs[i].RegionsConfig[key].ReadOnlyNodes).Should(PointTo(Equal(*region.ReadOnlyNodes)), "Replica Spec: ReadOnlyNodes is not the same")
			}
		}
	} else {
		ExpectWithOffset(1, requested.DeploymentSpec.ProviderSettings).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"RegionName": Equal(created.ProviderSettings.RegionName),
		})), "Deployment should be the same as requested by the user: Region Name")
	}
	if requested.DeploymentSpec.ProviderSettings.ProviderName == "TENANT" {
		ExpectWithOffset(1, requested.DeploymentSpec.ProviderSettings).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"BackingProviderName": Equal(created.ProviderSettings.BackingProviderName),
		})), "Deployment should be the same as requested by the user: Backking Provider Name")
	}
}

func CompareAdvancedDeploymentsSpec(requested model.DeploymentSpec, created mongodbatlas.AdvancedCluster) {
	advancedSpec := requested.AdvancedDeploymentSpec
	Expect(created.MongoDBVersion).ToNot(BeEmpty())
	Expect(created.MongoDBVersion).ToNot(BeEmpty())
	Expect(created.ConnectionStrings.StandardSrv).ToNot(BeEmpty())
	Expect(created.ConnectionStrings.Standard).ToNot(BeEmpty())
	Expect(created.Name).To(Equal(advancedSpec.Name))
	Expect(created.GroupID).To(Not(BeEmpty()))

	defaultPriority := 7
	for i, replicationSpec := range advancedSpec.ReplicationSpecs {
		for key, region := range replicationSpec.RegionConfigs {
			if region.Priority == nil {
				region.Priority = &defaultPriority
			}
			ExpectWithOffset(1, created.ReplicationSpecs[i].RegionConfigs[key].ProviderName).Should(Equal(region.ProviderName), "Replica Spec: ProviderName is not the same")
			ExpectWithOffset(1, created.ReplicationSpecs[i].RegionConfigs[key].RegionName).Should(Equal(region.RegionName), "Replica Spec: RegionName is not the same")
			ExpectWithOffset(1, created.ReplicationSpecs[i].RegionConfigs[key].Priority).Should(Equal(region.Priority), "Replica Spec: Priority is not the same")
		}
	}
}

func CompareServerlessSpec(requested model.DeploymentSpec, created mongodbatlas.Cluster) {
	serverlessSpec := requested.ServerlessSpec
	Expect(created.MongoDBVersion).ToNot(BeEmpty())
	Expect(created.ConnectionStrings.StandardSrv).ToNot(BeEmpty())
	Expect(created.Name).To(Equal(serverlessSpec.Name))
	Expect(created.GroupID).To(Not(BeEmpty()))
}

func SaveK8sResourcesTo(resources []string, ns string, destination string) {
	for _, resource := range resources {
		data := kubecli.GetYamlResource(resource, ns)
		path := fmt.Sprintf("output/%s/%s.yaml", destination, resource)
		utils.SaveToFile(path, data)
	}
}

func SaveK8sResources(resources []string, ns string) {
	SaveK8sResourcesTo(resources, ns, ns)
}

func SaveTestAppLogs(input model.UserInputs) {
	for _, user := range input.Users {
		utils.SaveToFile(
			fmt.Sprintf("output/%s/testapp-describe-%s.txt", input.Namespace, user.Spec.Username),
			kubecli.DescribeTestApp(config.TestAppLabelPrefix+user.Spec.Username, input.Namespace),
		)
		utils.SaveToFile(
			fmt.Sprintf("output/%s/testapp-logs-%s.txt", input.Namespace, user.Spec.Username),
			kubecli.GetLogs(config.TestAppLabelPrefix+user.Spec.Username, input.Namespace),
		)
	}
}

// SaveOperatorLogs save logs from user input namespace
func SaveOperatorLogs(input model.UserInputs) {
	utils.SaveToFile(
		fmt.Sprintf("output/%s/operator-logs.txt", input.Namespace),
		kubecli.GetManagerLogs(input.Namespace),
	)
}

// SaveDefaultOperatorLogs save logs from default namespace
func SaveDefaultOperatorLogs(input model.UserInputs) {
	utils.SaveToFile(
		fmt.Sprintf("output/%s/operator-logs.txt", input.Namespace),
		kubecli.GetManagerLogs("default"),
	)
}

func SaveDeploymentDump(input model.UserInputs) {
	kubecli.GetDeploymentDump(fmt.Sprintf("output/%s/dump", input.Namespace))
}

func CheckUsersAttributes(input model.UserInputs) {
	userDBResourceName := func(deploymentName string, user model.DBUser) string { // user name helmkind or kube-test-kind
		if input.KeyName[0:4] == "helm" {
			return fmt.Sprintf("atlasdatabaseusers.atlas.mongodb.com/%s-%s", deploymentName, user.Spec.Username)
		}
		return fmt.Sprintf("atlasdatabaseusers.atlas.mongodb.com/%s", user.ObjectMeta.Name)
	}

	for _, deployment := range input.Deployments {
		for _, user := range input.Users {
			var atlasUser *mongodbatlas.DatabaseUser

			getUser := func() bool {
				atlasUser = mongocli.GetUser(setAdminIfEmpty(user.Spec.DatabaseName), user.Spec.Username, input.ProjectID)
				return atlasUser != nil
			}

			EventuallyWithOffset(1, getUser, "7m", "10s").Should(BeTrue())
			EventuallyWithOffset(1,
				func() string {
					return kubecli.GetStatusCondition("Ready", input.Namespace, userDBResourceName(deployment.ObjectMeta.Name, user))
				},
				"7m", "1m",
			).Should(Equal("True"), "Kubernetes resource: User resources status `Ready` should be True")

			// Required fields
			ExpectWithOffset(1, *atlasUser).To(MatchFields(IgnoreExtras, Fields{
				"Username":     Equal(user.Spec.Username),
				"GroupID":      Equal(input.ProjectID),
				"DatabaseName": Or(Equal(user.Spec.DatabaseName), Equal("admin")),
			}), "Users attributes should be the same as requested by the user")

			for i, role := range atlasUser.Roles {
				ExpectWithOffset(1, role).To(MatchFields(IgnoreMissing, Fields{
					"RoleName":       Equal(user.Spec.Roles[i].RoleName),
					"DatabaseName":   Equal(user.Spec.Roles[i].DatabaseName),
					"CollectionName": Equal(user.Spec.Roles[i].CollectionName),
				}), "Users roles attributes should be the same as requested by the user")
			}
		}
	}
}

func CheckUsersCanUseApp(data *model.TestDataProvider) {
	input := data.Resources
	for i, user := range data.Resources.Users { // TODO in parallel(?)/ingress
		// data
		port := strconv.Itoa(i + data.PortGroup)
		key := port
		data := fmt.Sprintf("{\"key\":\"%s\",\"shipmodel\":\"heavy\",\"hp\":150}", key)

		helm.InstallTestApplication(input, user, port)
		WaitTestApplication(input.Namespace, "app=test-app-"+user.Spec.Username)

		app := appclient.NewTestAppClient(port)
		ExpectWithOffset(1, app.Get("")).Should(Equal("It is working"))
		ExpectWithOffset(1, app.Post(data)).ShouldNot(HaveOccurred())
		ExpectWithOffset(1, app.Get("/mongo/"+key)).Should(Equal(data))
	}
}

func CheckUsersCanUseOldApp(data *model.TestDataProvider) {
	input := data.Resources
	for i, user := range data.Resources.Users {
		// data
		port := strconv.Itoa(i + data.PortGroup)
		key := port
		data := fmt.Sprintf("{\"key\":\"%s\",\"shipmodel\":\"heavy\",\"hp\":150}", key)

		cli.Execute("kubectl", "delete", "pod", "-l", "app=test-app-"+user.Spec.Username, "-n", input.Namespace).Wait("2m")
		WaitTestApplication(input.Namespace, "app=test-app-"+user.Spec.Username)

		app := appclient.NewTestAppClient(port)
		ExpectWithOffset(1, app.Get("")).Should(Equal("It is working"))
		ExpectWithOffset(1, app.Get("/mongo/"+key)).Should(Equal(data))

		key = port + "up"
		dataUpdated := fmt.Sprintf("{\"key\":\"%s\",\"shipmodel\":\"heavy\",\"hp\":150}", key)
		ExpectWithOffset(1, app.Post(dataUpdated)).ShouldNot(HaveOccurred())
		ExpectWithOffset(1, app.Get("/mongo/"+key)).Should(Equal(dataUpdated))
	}
}

func PrepareUsersConfigurations(data *model.TestDataProvider) {
	By("Prepare namespaces and project configuration", func() {
		kubecli.CreateNamespace(data.Resources.Namespace)
		By("Create project spec", func() {
			GinkgoWriter.Write([]byte(data.Resources.ProjectPath + "\n"))
			utils.SaveToFile(data.Resources.ProjectPath, data.Resources.Project.ConvertByte())
		})
		if len(data.Resources.Deployments) > 0 {
			By("Create deployment spec", func() {
				data.Resources.Deployments[0].Spec.Project.Name = data.Resources.Project.GetK8sMetaName()
				utils.SaveToFile(
					data.Resources.Deployments[0].DeploymentFileName(data.Resources),
					utils.JSONToYAMLConvert(data.Resources.Deployments[0]),
				)
			})
		}
		if len(data.Resources.Users) > 0 {
			By("Create dbuser spec", func() {
				for _, user := range data.Resources.Users {
					user.SaveConfigurationTo(data.Resources.ProjectPath)
					if user.Spec.PasswordSecret != nil {
						kubecli.CreateRandomUserSecret(user.Spec.PasswordSecret.Name, data.Resources.Namespace)
					}
				}
			})
		}
	})
}

// CreateConnectionAtlasKey create connection: global or project level
func CreateConnectionAtlasKey(data *model.TestDataProvider) {
	By("Change resources depends on AtlasKey and create key", func() {
		if data.Resources.AtlasKeyAccessType.GlobalLevelKey {
			kubecli.CreateApiKeySecret(config.DefaultOperatorGlobalKey, data.Resources.Namespace)
		} else {
			kubecli.CreateApiKeySecret(data.Resources.KeyName, data.Resources.Namespace)
		}
	})
}

func createConnectionAtlasKeyFrom(data *model.TestDataProvider, key *mongodbatlas.APIKey) {
	By("Change resources depends on AtlasKey and create key", func() {
		if data.Resources.AtlasKeyAccessType.GlobalLevelKey {
			kubecli.CreateApiKeySecretFrom(config.DefaultOperatorGlobalKey, data.Resources.Namespace, os.Getenv("MCLI_ORG_ID"), key.PublicKey, key.PrivateKey)
		} else {
			kubecli.CreateApiKeySecretFrom(data.Resources.KeyName, data.Resources.Namespace, os.Getenv("MCLI_ORG_ID"), key.PublicKey, key.PrivateKey)
		}
	})
}

func recreateAtlasKeyIfNeed(data *model.TestDataProvider) {
	if !data.Resources.AtlasKeyAccessType.IsFullAccess() {
		aClient, err := atlas.AClient()
		Expect(err).ShouldNot(HaveOccurred())
		globalKey, err := aClient.AddKeyWithAccessList(data.Resources.ProjectID, data.Resources.AtlasKeyAccessType.Roles, data.Resources.AtlasKeyAccessType.Whitelist)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(globalKey.PublicKey).ShouldNot(BeEmpty())
		Expect(globalKey.PrivateKey).ShouldNot(BeEmpty())
		data.Resources.AtlasKeyAccessType.GlobalKeyAttached = globalKey

		kubecli.DeleteApiKeySecret(data.Resources.KeyName, data.Resources.Namespace)
		createConnectionAtlasKeyFrom(data, globalKey)
	}
}

func DeployProject(data *model.TestDataProvider, generation string) {
	By("Create users resources: keys, project", func() {
		CreateConnectionAtlasKey(data)
		kubecli.Apply(data.Resources.ProjectPath, "-n", data.Resources.Namespace)
	})
}

func UpdateProjectID(data *model.TestDataProvider) {
	atlasProject, err := kube.GetProjectResource(data)
	Expect(err).Should(BeNil(), "Error has Occurred")
	data.Resources.ProjectID = atlasProject.Status.ID
	Expect(data.Resources.ProjectID).ShouldNot(BeEmpty())
}

func DeployProjectAndWait(data *model.TestDataProvider, generation string) {
	By("Create users resources: keys, project", func() {
		CreateConnectionAtlasKey(data)
		kubecli.Apply(data.Resources.ProjectPath, "-n", data.Resources.Namespace)
		By("Wait project creation and get projectID", func() {
			WaitProject(data, generation)
			atlasProject, err := kube.GetProjectResource(data)
			Expect(err).Should(BeNil(), "Error has Occurred")
			data.Resources.ProjectID = atlasProject.Status.ID
			Expect(data.Resources.ProjectID).ShouldNot(BeEmpty())
		})
		recreateAtlasKeyIfNeed(data)
	})
}

func DeployDeployment(data *model.TestDataProvider, generation string) {
	By("Create deployment", func() {
		kubecli.Apply(data.Resources.Deployments[0].DeploymentFileName(data.Resources), "-n", data.Resources.Namespace)
	})
	By("Wait deployment creation", func() {
		WaitDeployment(data.Resources, "1")
	})
	By("check deployment Attribute", func() {
		deployment := mongocli.GetDeploymentsInfo(data.Resources.ProjectID, data.Resources.Deployments[0].Spec.DeploymentSpec.Name)
		CompareDeploymentsSpec(data.Resources.Deployments[0].Spec, deployment)
	})
}

func DeployUsers(data *model.TestDataProvider) {
	By("create users", func() {
		kubecli.Apply(data.Resources.GetResourceFolder()+"/user/", "-n", data.Resources.Namespace)
	})
	By("check database users Attibutes", func() {
		Eventually(CheckIfUsersExist(data.Resources), "2m", "10s").Should(BeTrue())
		CheckUsersAttributes(data.Resources)
	})
	By("Deploy application for user", func() {
		CheckUsersCanUseApp(data)
	})
}

// DeployUserResourcesAction deploy all user resources, wait, and check results
func DeployUserResourcesAction(data *model.TestDataProvider) {
	DeployProjectAndWait(data, "1")
	DeployDeployment(data, "1")
	DeployUsers(data)
}

func DeleteDBUsersApps(data model.TestDataProvider) {
	By("Delete dbusers applications", func() {
		for _, user := range data.Resources.Users {
			helm.Uninstall("test-app-"+user.Spec.Username, data.Resources.Namespace)
		}
	})
}

func DeleteUserResources(data *model.TestDataProvider) {
	DeleteUserResourcesDeployment(data)
	DeleteUserResourcesProject(data)
}

func DeleteUserResourcesDeployment(data *model.TestDataProvider) {
	By("Delete deployment", func() {
		kubecli.Delete(data.Resources.Deployments[0].DeploymentFileName(data.Resources), "-n", data.Resources.Namespace)
		Eventually(
			CheckIfDeploymentExist(data.Resources),
			"10m", "1m",
		).Should(BeFalse(), "Deployment should be deleted from Atlas")
	})
}

func DeleteUserResourcesProject(data *model.TestDataProvider) {
	By("Delete project", func() {
		kubecli.Delete(data.Resources.ProjectPath, "-n", data.Resources.Namespace)
		Eventually(
			func() bool {
				return mongocli.IsProjectInfoExist(data.Resources.ProjectID)
			},
			"5m", "20s",
		).Should(BeFalse(), "Project should be deleted from Atlas")
	})
}

func DeleteGlobalKeyIfExist(data model.TestDataProvider) {
	if data.Resources.AtlasKeyAccessType.GlobalLevelKey {
		By("Delete Global API key for test", func() {
			client, err := atlas.AClient()
			Expect(err).ShouldNot(HaveOccurred())
			err = client.DeleteGlobalKey(*data.Resources.AtlasKeyAccessType.GlobalKeyAttached)
			Expect(err).ShouldNot(HaveOccurred())
		})
	}
}

func AfterEachFinalCleanup(datas []model.TestDataProvider) {
	for _, data := range datas {
		GinkgoWriter.Write([]byte("AfterEach. Final cleanup...\n"))
		DeleteDBUsersApps(data)
		Expect(kubecli.DeleteNamespace(data.Resources.Namespace)).Should(Say("deleted"), "Cant delete namespace after testing")
		DeleteGlobalKeyIfExist(data)
		GinkgoWriter.Write([]byte("AfterEach. Cleanup finished\n"))
	}
}

func setAdminIfEmpty(input string) string {
	if input == "" {
		return "admin"
	}

	return input
}
