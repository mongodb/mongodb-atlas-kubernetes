// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package actions

import (
	"context"
	"fmt"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"go.mongodb.org/atlas-sdk/v20250312011/admin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/api/atlas"
	appclient "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/appclient"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/cli"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/cli/helm"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
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
			// Waiting for a particular generation can be brittle, to make it more robust
			// wait for any progress to or beyond the expected generation
		).WithTimeout(5 * time.Minute).WithPolling(10 * time.Second).Should(BeNumerically(">=", generation))

		WaitDeploymentWithoutGenerationCheck(data)
	}

	if len(data.InitialDeployments) > 0 {
		EventuallyWithOffset(1,
			func(g Gomega) int {
				gen, err := k8s.GetDeploymentObservedGeneration(data.Context, data.K8SClient, input.Namespace, data.InitialDeployments[0].ObjectMeta.GetName())
				g.Expect(err).ToNot(HaveOccurred())
				return gen
			},
			// Waiting for a particular generation can be brittle, to make it more robust
			// wait for any progress to or beyond the expected generation
		).WithTimeout(5 * time.Minute).WithPolling(10 * time.Second).Should(BeNumerically(">=", generation))

		WaitDeploymentWithoutGenerationCheckV2(data)
	}
}

// WaitDeploymentWithoutGenerationCheckV2 waits for deployment
// nolint: dupl
func WaitDeploymentWithoutGenerationCheckV2(data *model.TestDataProvider) {
	input := data.Resources
	EventuallyWithOffset(1,
		func(g Gomega) string {
			deploymentStatus, err := k8s.GetDeploymentStatusCondition(data.Context, data.K8SClient, api.ReadyType, input.Namespace, data.InitialDeployments[0].ObjectMeta.GetName())
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
	case deployment.Spec.FlexSpec != nil:
		atlasClient, err := atlas.AClient()
		Expect(err).To(BeNil())
		flexInstance, err := atlasClient.GetFlexInstance(input.ProjectID, deployment.Spec.FlexSpec.Name)
		Expect(err).To(BeNil())
		Expect(flexInstance.StateName).To(Equal("IDLE"))
	default:
		aClient := atlas.GetClientOrFail()
		deployment, err := aClient.GetDeployment(input.ProjectID, input.Deployments[0].Spec.GetDeploymentName())
		Expect(err).To(BeNil())
		Expect(deployment.StateName).Should(Equal("IDLE"))
	}
}

// WaitDeploymentWithoutGenerationCheck waits for deployment
// nolint: dupl
func WaitDeploymentWithoutGenerationCheck(data *model.TestDataProvider) {
	input := data.Resources
	EventuallyWithOffset(1,
		func(g Gomega) string {
			deploymentStatus, err := k8s.GetDeploymentStatusCondition(data.Context, data.K8SClient, api.ReadyType, input.Namespace, input.Deployments[0].ObjectMeta.GetName())
			g.Expect(err).ToNot(HaveOccurred())
			return deploymentStatus
		},
	).WithTimeout(40*time.Minute).WithPolling(1*time.Minute).Should(Equal("True"), "Kubernetes resource: Deployment status `Ready` should be 'True'")

	Eventually(func(g Gomega) {
		deploymentState, err := k8s.GetK8sDeploymentStateName(data.Context, data.K8SClient,
			input.Namespace, input.Deployments[0].ObjectMeta.GetName())
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(deploymentState).Should(Equal("IDLE"), "Kubernetes resource: Deployment status should be IDLE")
	}).WithTimeout(40 * time.Minute).WithPolling(1 * time.Minute).Should(Succeed())

	deploymentState, err := k8s.GetK8sDeploymentStateName(data.Context, data.K8SClient,
		input.Namespace, input.Deployments[0].ObjectMeta.GetName())
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
	ExpectWithOffset(1, deploymentState).Should(Equal("IDLE"), "Kubernetes resource: Deployment status should be IDLE")
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
	EventuallyWithOffset(1, isAppRunning(), "10m", "10s").Should(BeTrue(), "Test application should be running")
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

func CompareAdvancedDeploymentsSpec(requested model.DeploymentSpec, created admin.ClusterDescription20240805) {
	advancedSpec := requested.DeploymentSpec

	Expect(created.GetMongoDBVersion()).ToNot(BeEmpty())
	Expect(created.ConnectionStrings.GetStandard()).ToNot(BeEmpty())
	Expect(created.ConnectionStrings.GetStandardSrv()).ToNot(BeEmpty())
	Expect(created.GetName()).To(Equal(requested.GetDeploymentName()))
	Expect(created.GetGroupId()).To(Not(BeEmpty()))

	defaultPriority := 7
	for i, replicationSpec := range advancedSpec.ReplicationSpecs {
		for key, region := range replicationSpec.RegionConfigs {
			if region.Priority == nil {
				region.Priority = &defaultPriority
			}
			ExpectWithOffset(1, created.GetReplicationSpecs()[i].GetRegionConfigs()[key].GetProviderName()).Should(Equal(region.ProviderName), "Replica Spec: ProviderName is not the same")
			ExpectWithOffset(1, created.GetReplicationSpecs()[i].GetRegionConfigs()[key].GetRegionName()).Should(Equal(region.RegionName), "Replica Spec: RegionName is not the same")
			ExpectWithOffset(1, created.GetReplicationSpecs()[i].GetRegionConfigs()[key].Priority).Should(Equal(region.Priority), "Replica Spec: Priority is not the same")
		}
	}
}

func CompareFlexSpec(requested model.DeploymentSpec, created admin.FlexClusterDescription20241113) {
	flexSpec := requested.FlexSpec
	Expect(created.GetMongoDBVersion()).ToNot(BeEmpty())
	Expect(created.ConnectionStrings.GetStandardSrv()).ToNot(BeEmpty())
	Expect(created.GetName()).To(Equal(flexSpec.Name))
	Expect(created.GetGroupId()).To(Not(BeEmpty()))
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

func SaveAtlasOrgSettingsToFile(ctx context.Context, k8sClient client.Client, ns string) error {
	yaml, err := k8s.AtlasOrgSettingsListYaml(ctx, k8sClient, ns)
	if err != nil {
		return fmt.Errorf("error getting AtlasOrgSettings list: %w", err)
	}
	path := fmt.Sprintf("output/%s/%s.yaml", ns, "atlasorgsettings")
	err = utils.SaveToFile(path, yaml)
	if err != nil {
		return fmt.Errorf("error saving AtlasOrgSettings to file: %w", err)
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
		testAppName := fmt.Sprintf("test-app-%s", user.Spec.Username)
		bytes, err := k8s.GetPodLogsByDeployment(testAppName, input.Namespace, corev1.PodLogOptions{})
		Expect(err).ToNot(HaveOccurred())

		utils.SaveToFile(
			fmt.Sprintf("output/%s/testapp-logs-%s.txt", input.Namespace, user.Spec.Username),
			bytes,
		)
	}
}

func CheckUsersAttributes(data *model.TestDataProvider) {
	input := data.Resources
	aClient := atlas.GetClientOrFail()
	userDBResourceName := func(deploymentName string, user *akov2.AtlasDatabaseUser) string { // user name helmkind or kube-test-kind
		if input.KeyName[0:4] == "helm" {
			return fmt.Sprintf("%s-%s", deploymentName, user.Spec.Username)
		}
		return user.ObjectMeta.GetName()
	}

	for _, deployment := range data.InitialDeployments {
		for _, user := range data.Users {
			var atlasUser *admin.CloudDatabaseUser

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
					userStatus, err := k8s.GetDBUserStatusCondition(data.Context, data.K8SClient, api.ReadyType, data.Resources.Namespace, userDBResourceName(deployment.ObjectMeta.Name, user))
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

			for i, role := range atlasUser.GetRoles() {
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
		By(fmt.Sprintf("Checking user %s (%d) can use old App", user.Spec.Username, i), func() {
			// data
			port := strconv.Itoa(i + data.PortGroup)
			key := port
			expectedData := fmt.Sprintf("{\"key\":\"%s\",\"shipmodel\":\"heavy\",\"hp\":150}", key)

			By("Deleting pod to force an app restart and wait for it", func() {
				cli.Execute(
					"kubectl", "delete", "pod", "-l", "app=test-app-"+user.Spec.Username, "-n", input.Namespace,
				).Wait("2m")
				WaitTestApplication(data, input.Namespace, "app", "test-app-"+user.Spec.Username)
			})

			By(fmt.Sprintf("Verifying if a user (%s) is READY", user.Spec.Username), func() {
				Eventually(func(g Gomega) {
					dbu := &akov2.AtlasDatabaseUser{}
					g.Expect(data.K8SClient.Get(data.Context, client.ObjectKey{
						Name:      fmt.Sprintf("%s-%s", data.Resources.Deployments[0].Spec.GetDeploymentName(), user.Spec.Username),
						Namespace: data.Resources.Namespace},
						dbu,
					)).To(Succeed())
					g.Expect(dbu.Status.Conditions).ShouldNot(BeEmpty())
					for _, condition := range dbu.Status.Conditions {
						if condition.Type == api.ReadyType {
							g.Expect(condition.Status).Should(Equal(corev1.ConditionTrue), "User should be ready")
						}
					}
				}).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})

			app := appclient.NewTestAppClient(port)
			By("Test restarted App access", func() {
				getRoot := app.Get("")
				GinkgoWriter.Write([]byte(fmt.Sprintf("Test App GET: %q\n", getRoot)))
				ExpectWithOffset(1, getRoot).Should(Equal("It is working"))
				getKey := app.Get("/mongo/" + key)
				GinkgoWriter.Write([]byte(fmt.Sprintf("Test App GET /mongo/%s: %q\n", key, getKey)))
				ExpectWithOffset(1, getKey).Should(Equal(expectedData))
			})

			By("Test restarted App update", func() {
				key = port + "up"
				dataUpdated := fmt.Sprintf("{\"key\":\"%s\",\"shipmodel\":\"heavy\",\"hp\":150}", key)
				err := app.Post(dataUpdated)
				GinkgoWriter.Write([]byte(fmt.Sprintf("Test App POST %v: %v\n", dataUpdated, err)))
				ExpectWithOffset(1, err).ShouldNot(HaveOccurred())
				getKey := app.Get("/mongo/" + key)
				GinkgoWriter.Write([]byte(fmt.Sprintf("Test App GET /mongo/%s: %q\n", key, getKey)))
				ExpectWithOffset(1, getKey).Should(Equal(dataUpdated))
			})
		})
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
				data.Resources.Deployments[0].Spec.ProjectRef.Name = data.Resources.Project.GetK8sMetaName()
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

func DeleteDBUsersApps(data model.TestDataProvider) {
	By("Delete dbusers applications", func() {
		for _, user := range data.Resources.Users {
			helm.Uninstall("test-app-"+user.Spec.Username, data.Resources.Namespace)
		}
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
			).WithTimeout(30*time.Minute).WithPolling(20*time.Second).Should(BeFalse(), "Deployment should be deleted from Atlas")
		}
	})
}

func DeleteTestDataUsers(data *model.TestDataProvider) {
	By("Delete Users", func() {
		for _, user := range data.Users {
			Expect(data.K8SClient.Get(data.Context, types.NamespacedName{Name: data.Project.Name, Namespace: data.Project.Namespace}, data.Project)).Should(Succeed())
			Expect(data.K8SClient.Get(data.Context, types.NamespacedName{Name: user.Name, Namespace: user.Namespace}, user)).Should(Succeed())
			Expect(data.K8SClient.Delete(data.Context, user)).Should(Succeed())
		}
	})
}

func DeleteAtlasGlobalKeyIfExist(data model.TestDataProvider) {
	if data.Resources.AtlasKeyAccessType.GlobalLevelKey {
		By("Delete Global API key for test", func() {
			atlasClient, err := atlas.AClient()
			Expect(err).ShouldNot(HaveOccurred())
			if data.Resources.AtlasKeyAccessType.GlobalKeyAttached != nil {
				err = atlasClient.DeleteGlobalKey(*data.Resources.AtlasKeyAccessType.GlobalKeyAttached)
				Expect(err).ShouldNot(HaveOccurred())
			}
		})
	}
}

func DeleteTestDataNetworkContainers(data *model.TestDataProvider) {
	By("Delete network containers", func() {
		containers := &akov2.AtlasNetworkContainerList{}
		Expect(data.K8SClient.List(data.Context, containers, &client.ListOptions{Namespace: data.Resources.Namespace})).Should(Succeed())
		for _, container := range containers.Items {
			key := client.ObjectKey{Name: container.Name, Namespace: container.Namespace}
			Expect(data.K8SClient.Delete(data.Context, &container)).Should(Succeed())
			Eventually(
				func() bool {
					foundContainer := &akov2.AtlasNetworkContainer{}
					err := data.K8SClient.Get(data.Context, key, foundContainer)
					return err != nil && errors.IsNotFound(err)
				},
			).WithTimeout(10*time.Minute).WithPolling(20*time.Second).Should(BeTrue(), "Network container should be deleted from Atlas")
		}
	})
}

func DeleteTestDataNetworkPeerings(data *model.TestDataProvider) {
	By("Delete network peerings", func() {
		peerings := &akov2.AtlasNetworkPeeringList{}
		Expect(data.K8SClient.List(data.Context, peerings, &client.ListOptions{Namespace: data.Resources.Namespace})).Should(Succeed())
		for _, peering := range peerings.Items {
			Expect(data.K8SClient.Delete(data.Context, &peering)).Should(Succeed())
			key := client.ObjectKey{Name: peering.Name, Namespace: peering.Namespace}
			Eventually(
				func() bool {
					foundPeering := &akov2.AtlasNetworkPeering{}
					err := data.K8SClient.Get(data.Context, key, foundPeering)
					return err != nil && errors.IsNotFound(err)
				},
			).WithTimeout(10*time.Minute).WithPolling(20*time.Second).Should(BeTrue(), "Network peering should be deleted from Atlas")
		}
	})
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
