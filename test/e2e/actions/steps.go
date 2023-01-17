package actions

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/k8s"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"

	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"go.mongodb.org/atlas/mongodbatlas"

	kube "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/atlas"
	appclient "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/appclient"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/helm"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

func WaitDeployment(data *model.TestDataProvider, generation int) {
	input := data.Resources
	if len(data.Resources.Deployments) > 0 {
		EventuallyWithOffset(1,
			func(g Gomega) int {
				gen, err := k8s.GetDeploymentObservedGeneration(data.Context, data.K8SClient, input.Namespace, input.Deployments[0].ObjectMeta.GetName())
				g.Expect(err).ToNot(HaveOccurred())
				return gen
			},
		).WithTimeout(5 * time.Minute).WithPolling(10 * time.Second).Should(Equal(generation))

		WaitDeploymentWithoutGenerationCheck(data)
	}

	if len(data.InitialDeployments) > 0 {
		EventuallyWithOffset(1,
			func(g Gomega) int {
				gen, err := k8s.GetDeploymentObservedGeneration(data.Context, data.K8SClient, input.Namespace, data.InitialDeployments[0].ObjectMeta.GetName())
				g.Expect(err).ToNot(HaveOccurred())
				return gen
			},
		).WithTimeout(5 * time.Minute).WithPolling(10 * time.Second).Should(Equal(generation))

		WaitDeploymentWithoutGenerationCheckV2(data)
	}
}

// nolint: dupl
func WaitDeploymentWithoutGenerationCheckV2(data *model.TestDataProvider) {
	input := data.Resources
	EventuallyWithOffset(1,
		func(g Gomega) string {
			deploymentStatus, err := k8s.GetDeploymentStatusCondition(data.Context, data.K8SClient, status.ReadyType, input.Namespace, data.InitialDeployments[0].ObjectMeta.GetName())
			g.Expect(err).ToNot(HaveOccurred())
			return deploymentStatus
		},
		"60m", "1m",
	).Should(Equal("True"), "Kubernetes resource: Deployment status `Ready` should be 'True'")

	deploymentState, err := k8s.GetK8sDeploymentStateName(data.Context, data.K8SClient,
		input.Namespace, data.InitialDeployments[0].ObjectMeta.GetName())
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
	ExpectWithOffset(1, deploymentState).Should(Equal("IDLE"), "Kubernetes resource: Deployment status should be IDLE")

	deployment := data.InitialDeployments[0]
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
		aClient := atlas.GetClientOrFail()
		Expect(aClient.GetDeployment(input.ProjectID, input.Deployments[0].Spec.GetDeploymentName()).StateName).Should(Equal("IDLE"))
	}
}

// nolint: dupl
func WaitDeploymentWithoutGenerationCheck(data *model.TestDataProvider) {
	input := data.Resources
	EventuallyWithOffset(1,
		func(g Gomega) string {
			deploymentStatus, err := k8s.GetDeploymentStatusCondition(data.Context, data.K8SClient, status.ReadyType, input.Namespace, input.Deployments[0].ObjectMeta.GetName())
			g.Expect(err).ToNot(HaveOccurred())
			return deploymentStatus
		},
		"60m", "1m",
	).Should(Equal("True"), "Kubernetes resource: Deployment status `Ready` should be 'True'")

	deploymentState, err := k8s.GetK8sDeploymentStateName(data.Context, data.K8SClient,
		input.Namespace, input.Deployments[0].ObjectMeta.GetName())
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
	ExpectWithOffset(1, deploymentState).Should(Equal("IDLE"), "Kubernetes resource: Deployment status should be IDLE")

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
		aClient := atlas.GetClientOrFail()
		Expect(aClient.GetDeployment(input.ProjectID, input.Deployments[0].Spec.GetDeploymentName()).StateName).Should(Equal("IDLE"))
	}
}

func WaitProject(data *model.TestDataProvider, generation int) {
	EventuallyWithOffset(1, kube.ProjectReadyCondition(data), "25m", "10s").Should(Equal("True"), "Kubernetes resource: Project status `Ready` should be 'True'")
	gen, err := k8s.GetProjectObservedGeneration(data.Context, data.K8SClient, data.Resources.Namespace, data.Resources.Project.GetK8sMetaName())
	Expect(err).ToNot(HaveOccurred())
	ExpectWithOffset(1, gen).Should(Equal(generation), "Kubernetes resource: Generation should be upgraded")
	atlasProject, err := kube.GetProjectResource(data)
	Expect(err).ShouldNot(HaveOccurred())
	ExpectWithOffset(1, atlasProject.Status.ID).ShouldNot(BeNil(), "Kubernetes resource: Project status should have non-empty ID field")
}

func WaitProjectWithoutGenerationCheck(data *model.TestDataProvider) {
	EventuallyWithOffset(1, func() string {
		return kube.ProjectReadyCondition(data)
	}, "15m", "10s").Should(Equal("True"), "Kubernetes resource: Project status `Ready` should be 'True'")
	atlasProject, err := kube.GetProjectResource(data)
	Expect(err).ShouldNot(HaveOccurred())
	ExpectWithOffset(1, atlasProject.Status.ID).ShouldNot(BeNil(), "Kubernetes resource: Project status should have non-empty ID field")
}

func WaitTestApplication(data *model.TestDataProvider, ns, labelKey, labelValue string) {
	// temp
	isAppRunning := func() func() bool {
		return func() bool {
			phase, _ := k8s.GetPodStatusPhaseByLabel(data.Context, data.K8SClient, ns, labelKey, labelValue)
			return phase == "Running"
		}
	}
	EventuallyWithOffset(1, isAppRunning(), "2m", "10s").Should(BeTrue(), "Test application should be running")
}

func CheckIfDeploymentExist(input model.UserInputs) func() bool {
	return func() bool {
		aClient := atlas.GetClientOrFail()
		return aClient.IsDeploymentExist(input.ProjectID, input.Deployments[0].Spec.DeploymentSpec.Name)
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

func CheckUserExistInAtlas(data *model.TestDataProvider) func() bool {
	return func() bool {
		atlasClient, err := atlas.AClient()
		if err != nil {
			return false
		}

		for _, user := range data.Users {
			dbUser, err := atlasClient.GetUserByName(setAdminIfEmpty(user.Spec.DatabaseName), data.Project.ID(), user.Spec.Username)
			if err != nil && dbUser == nil {
				return false
			}
		}
		return true
	}
}

func CompareDeploymentsSpec(requested model.DeploymentSpec, created mongodbatlas.Cluster) {
	ExpectWithOffset(1, created).To(MatchFields(IgnoreExtras, Fields{
		"Name": Equal(requested.DeploymentSpec.Name),
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
				// different type
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
		})), "Deployment should be the same as requested by the user: Backing Provider Name")
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

func SaveProjectsToFile(ctx context.Context, k8sClient client.Client, ns string) error {
	yaml, err := k8s.ProjectListYaml(ctx, k8sClient, ns)
	if err != nil {
		return fmt.Errorf("error getting project list: %w", err)
	}
	path := fmt.Sprintf("output/%s/%s.yaml", ns, "projects")
	err = utils.SaveToFile(path, yaml)
	if err != nil {
		return fmt.Errorf("error saving projects to file: %w", err)
	}
	return nil
}

func SaveTeamsToFile(ctx context.Context, k8sClient client.Client, ns string) error {
	yaml, err := k8s.TeamListYaml(ctx, k8sClient, ns)
	if err != nil {
		return fmt.Errorf("error getting team list: %w", err)
	}
	path := fmt.Sprintf("output/%s/%s.yaml", ns, "teams")
	err = utils.SaveToFile(path, yaml)
	if err != nil {
		return fmt.Errorf("error saving teams to file: %w", err)
	}
	return nil
}

func SaveDeploymentsToFile(ctx context.Context, k8sClient client.Client, ns string) error {
	yaml, err := k8s.DeploymentListYml(ctx, k8sClient, ns)
	if err != nil {
		return fmt.Errorf("error getting deployment list: %w", err)
	}
	path := fmt.Sprintf("output/%s/%s.yaml", ns, "deployments")
	err = utils.SaveToFile(path, yaml)
	if err != nil {
		return fmt.Errorf("error saving deployments to file: %w", err)
	}
	return nil
}

func SaveUsersToFile(ctx context.Context, k8sClient client.Client, ns string) error {
	yaml, err := k8s.UserListYaml(ctx, k8sClient, ns)
	if err != nil {
		return fmt.Errorf("error getting user list: %w", err)
	}
	path := fmt.Sprintf("output/%s/%s.yaml", ns, "users")
	err = utils.SaveToFile(path, yaml)
	if err != nil {
		return fmt.Errorf("error saving users to file: %w", err)
	}
	return nil
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

func SaveDeploymentDump(input model.UserInputs) {
	kubecli.GetDeploymentDump(fmt.Sprintf("output/%s/dump", input.Namespace))
}

func CheckUsersAttributes(data *model.TestDataProvider) {
	input := data.Resources
	aClient := atlas.GetClientOrFail()
	userDBResourceName := func(deploymentName string, user *v1.AtlasDatabaseUser) string { // user name helmkind or kube-test-kind
		if input.KeyName[0:4] == "helm" {
			return fmt.Sprintf("%s-%s", deploymentName, user.Spec.Username)
		}
		return user.ObjectMeta.GetName()
	}

	for _, deployment := range data.InitialDeployments {
		for _, user := range data.Users {
			var atlasUser *mongodbatlas.DatabaseUser

			getUser := func() bool {
				var err error
				atlasUser, err = aClient.GetDBUser(setAdminIfEmpty(user.Spec.DatabaseName), user.Spec.Username, input.ProjectID)
				if err != nil {
					return false
				}
				return atlasUser != nil
			}

			EventuallyWithOffset(1, getUser, "7m", "10s").Should(BeTrue())
			EventuallyWithOffset(1,
				func() string {
					userStatus, err := k8s.GetDBUserStatusCondition(data.Context, data.K8SClient, status.ReadyType, data.Resources.Namespace, userDBResourceName(deployment.ObjectMeta.Name, user))
					if err != nil {
						return err.Error()
					}
					return userStatus
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
		postData := fmt.Sprintf("{\"key\":\"%s\",\"shipmodel\":\"heavy\",\"hp\":150}", key)

		helm.InstallTestApplication(input, user, port)
		WaitTestApplication(data, input.Namespace, "app", "test-app-"+user.Spec.Username)

		app := appclient.NewTestAppClient(port)
		ExpectWithOffset(1, app.Get("")).Should(Equal("It is working"))
		ExpectWithOffset(1, app.Post(postData)).ShouldNot(HaveOccurred())
		ExpectWithOffset(1, app.Get("/mongo/"+key)).Should(Equal(postData))
	}
}

func CheckUsersCanUseOldApp(data *model.TestDataProvider) {
	input := data.Resources
	for i, user := range data.Resources.Users {
		// data
		port := strconv.Itoa(i + data.PortGroup)
		key := port
		expectedData := fmt.Sprintf("{\"key\":\"%s\",\"shipmodel\":\"heavy\",\"hp\":150}", key)

		cli.Execute("kubectl", "delete", "pod", "-l", "app=test-app-"+user.Spec.Username, "-n", input.Namespace).Wait("2m")
		WaitTestApplication(data, input.Namespace, "app", "test-app-"+user.Spec.Username)

		app := appclient.NewTestAppClient(port)
		ExpectWithOffset(1, app.Get("")).Should(Equal("It is working"))
		ExpectWithOffset(1, app.Get("/mongo/"+key)).Should(Equal(expectedData))

		key = port + "up"
		dataUpdated := fmt.Sprintf("{\"key\":\"%s\",\"shipmodel\":\"heavy\",\"hp\":150}", key)
		ExpectWithOffset(1, app.Post(dataUpdated)).ShouldNot(HaveOccurred())
		ExpectWithOffset(1, app.Get("/mongo/"+key)).Should(Equal(dataUpdated))
	}
}

func PrepareUsersConfigurations(data *model.TestDataProvider) {
	By("Prepare namespaces and project configuration", func() {
		err := k8s.CreateNamespace(data.Context, data.K8SClient, data.Resources.Namespace)
		Expect(err).NotTo(HaveOccurred())
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
						Expect(k8s.CreateRandomUserSecret(data.Context, data.K8SClient,
							user.Spec.PasswordSecret.Name, data.Resources.Namespace)).Should(Succeed())
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
			By("Create secret in the data namespace")
			k8s.CreateDefaultSecret(data.Context, data.K8SClient, config.DefaultOperatorGlobalKey, data.Resources.Namespace)
		} else {
			By("Create secret in the data prefix name")
			k8s.CreateDefaultSecret(data.Context, data.K8SClient, data.Prefix, data.Resources.Namespace)
		}
	})
}

func createConnectionAtlasKeyFrom(data *model.TestDataProvider, key *mongodbatlas.APIKey) {
	By("Change resources depends on AtlasKey and create key", func() {
		if data.Resources.AtlasKeyAccessType.GlobalLevelKey {
			err := k8s.CreateSecret(data.Context, data.K8SClient, key.PublicKey, key.PrivateKey, config.DefaultOperatorGlobalKey, data.Resources.Namespace)
			Expect(err).NotTo(HaveOccurred())
		} else {
			err := k8s.CreateSecret(data.Context, data.K8SClient, key.PublicKey, key.PrivateKey, data.Resources.KeyName, data.Resources.Namespace)
			Expect(err).NotTo(HaveOccurred())
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

		k8s.DeleteKey(data.Context, data.K8SClient, data.Resources.KeyName, data.Resources.Namespace)
		createConnectionAtlasKeyFrom(data, globalKey)
	}
}

func DeployProjectAndWait(data *model.TestDataProvider, generation int) {
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

func DeployDeployment(data *model.TestDataProvider) {
	if len(data.Resources.Deployments) > 0 {
		By("Create deployment", func() {
			kubecli.Apply(data.Resources.Deployments[0].DeploymentFileName(data.Resources), "-n", data.Resources.Namespace)
		})
		By("Wait deployment creation", func() {
			WaitDeploymentWithoutGenerationCheck(data)
		})
		By("check deployment Attribute", func() {
			aClient, err := atlas.AClient()
			Expect(err).NotTo(HaveOccurred())
			deployment := aClient.GetDeployment(data.Resources.ProjectID, data.Resources.Deployments[0].Spec.DeploymentSpec.Name)
			CompareDeploymentsSpec(data.Resources.Deployments[0].Spec, deployment)
		})
	}
}

func DeployUsers(data *model.TestDataProvider) {
	By("create users", func() {
		kubecli.Apply(data.Resources.GetResourceFolder()+"/user/", "-n", data.Resources.Namespace)
	})
	By("check database users Attributes", func() {
		Eventually(CheckIfUsersExist(data.Resources), "2m", "10s").Should(BeTrue())
		CheckUsersAttributes(data)
	})
	By("Deploy application for user", func() {
		CheckUsersCanUseApp(data)
	})
}

// DeployUserResourcesAction deploy all user resources, wait, and check results
func DeployUserResourcesAction(data *model.TestDataProvider) {
	DeployProjectAndWait(data, 1)
	DeployDeployment(data)
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
			func(g Gomega) bool {
				aClient := atlas.GetClientOrFail()
				return aClient.IsProjectExists(g, data.Resources.ProjectID)
			},
			"5m", "20s",
		).Should(BeFalse(), "Project should be deleted from Atlas")
	})
}

func DeleteTestDataProject(data *model.TestDataProvider) {
	By("Delete project", func() {
		Expect(data.K8SClient.Get(data.Context, types.NamespacedName{Name: data.Project.Name, Namespace: data.Project.Namespace}, data.Project)).Should(Succeed())
		projectID := data.Project.Status.ID
		Expect(data.K8SClient.Delete(data.Context, data.Project)).Should(Succeed())
		if projectID != "" {
			Eventually(
				func(g Gomega) bool {
					aClient := atlas.GetClientOrFail()
					return aClient.IsProjectExists(g, projectID)
				},
				"15m", "20s",
			).Should(BeFalse(), "Project should be deleted from Atlas")
		}
	})
}

func DeleteTestDataTeams(data *model.TestDataProvider) {
	By("Delete teams", func() {
		teams := &v1.AtlasTeamList{}
		Expect(data.K8SClient.List(data.Context, teams, &client.ListOptions{Namespace: data.Resources.Namespace})).Should(Succeed())
		for i := range teams.Items {
			Expect(data.K8SClient.Delete(data.Context, &teams.Items[i])).Should(Succeed())
		}
	})
}

func DeleteTestDataDeployments(data *model.TestDataProvider) {
	By("Delete deployment", func() {
		for _, deployment := range data.InitialDeployments {
			Expect(data.K8SClient.Get(data.Context, types.NamespacedName{Name: data.Project.Name, Namespace: data.Project.Namespace}, data.Project)).Should(Succeed())
			projectID := data.Project.Status.ID
			Expect(data.K8SClient.Get(data.Context, types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, deployment)).Should(Succeed())
			Expect(data.K8SClient.Delete(data.Context, deployment)).Should(Succeed())
			deploymentName := deployment.AtlasName()
			Eventually(
				func() bool {
					aClient := atlas.GetClientOrFail()
					return aClient.IsDeploymentExist(projectID, deploymentName)
				},
				"7m", "20s",
			).Should(BeFalse(), "Deployment should be deleted from Atlas")
		}
	})
}

func DeleteAtlasGlobalKeyIfExist(data model.TestDataProvider) {
	if data.Resources.AtlasKeyAccessType.GlobalLevelKey {
		By("Delete Global API key for test", func() {
			client, err := atlas.AClient()
			Expect(err).ShouldNot(HaveOccurred())
			if data.Resources.AtlasKeyAccessType.GlobalKeyAttached != nil {
				err = client.DeleteGlobalKey(*data.Resources.AtlasKeyAccessType.GlobalKeyAttached)
				Expect(err).ShouldNot(HaveOccurred())
			}
		})
	}
}

func AfterEachFinalCleanup(datas []model.TestDataProvider) {
	for i := range datas {
		data := datas[i]
		GinkgoWriter.Write([]byte("AfterEach. Final cleanup...\n"))
		DeleteDBUsersApps(data)
		DeleteAtlasGlobalKeyIfExist(data)
		Expect(k8s.DeleteNamespace(data.Context, data.K8SClient, data.Resources.Namespace)).Should(Succeed(), "Can't delete namespace")
		GinkgoWriter.Write([]byte("AfterEach. Cleanup finished\n"))
	}
}

func setAdminIfEmpty(input string) string {
	if input == "" {
		return "admin"
	}

	return input
}
