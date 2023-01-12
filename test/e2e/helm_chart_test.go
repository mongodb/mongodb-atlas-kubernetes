package e2e_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	actions "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	kube "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/kube"
	helm "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/helm"

	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

var _ = Describe("HELM charts", func() {
	var data model.TestDataProvider

	_ = BeforeEach(func() {
		imageURL := os.Getenv("IMAGE_URL")
		Expect(imageURL).ShouldNot(BeEmpty(), "SetUP IMAGE_URL")
		Eventually(kubecli.GetVersionOutput()).Should(Say(K8sVersion))
	})

	_ = AfterEach(func() {
		By("After each.", func() {
			GinkgoWriter.Write([]byte("\n"))
			GinkgoWriter.Write([]byte("===============================================\n"))
			GinkgoWriter.Write([]byte("Operator namespace: " + data.Resources.Namespace + "\n"))
			GinkgoWriter.Write([]byte("===============================================\n"))
			if CurrentSpecReport().Failed() {
				GinkgoWriter.Write([]byte("Resources wasn't clean"))
				utils.SaveToFile(
					fmt.Sprintf("output/%s/operator-logs-default.txt", data.Resources.Namespace),
					kubecli.GetManagerLogs(config.DefaultOperatorNS),
				)
				utils.SaveToFile(
					fmt.Sprintf("output/%s/operator-logs.txt", data.Resources.Namespace),
					kubecli.GetManagerLogs(data.Resources.Namespace),
				)
				actions.SaveK8sResourcesTo(
					[]string{"deploy"},
					"default",
					data.Resources.Namespace,
				)
				actions.SaveProjectsToFile(data.Context, data.K8SClient, data.Resources.Namespace)
				actions.SaveK8sResources(
					[]string{"atlasdeployments", "atlasdatabaseusers"},
					data.Resources.Namespace,
				)
				actions.SaveTestAppLogs(data.Resources)
				actions.AfterEachFinalCleanup([]model.TestDataProvider{data})
			}
		})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration", Label("helm-ns"),
		func(test model.TestDataProvider, deploymentType string) { // deploymentType - probably will be moved later ()
			data = test
			GinkgoWriter.Println(data.Resources.KeyName)
			switch deploymentType {
			case "advanced":
				data.Resources.Deployments[0].Spec.AdvancedDeploymentSpec.Name = data.Resources.KeyName
			case "serverless":
				data.Resources.Deployments[0].Spec.ServerlessSpec.Name = data.Resources.KeyName
			default:
				data.Resources.Deployments[0].Spec.DeploymentSpec.Name = data.Resources.KeyName
			}
			data.Resources.Deployments[0].ObjectMeta.Name = data.Resources.KeyName

			By("Install CRD", func() {
				helm.InstallCRD(data.Resources)
			})
			By("User use helm for deploying namespaces operator", func() {
				helm.InstallOperatorNamespacedSubmodule(data.Resources)
			})
			By("User deploy the deployment via helm", func() {
				helm.InstallDeploymentSubmodule(data.Resources)
			})
			waitDeploymentWithChecks(&data)
			By("Additional check for the current data set", func() {
				for _, check := range data.Actions {
					check(&data)
				}
			})
			deleteDeploymentAndOperator(&data)
		},
		Entry("Several actions with helm update", Label("helm-ns-flow"),
			model.DataProviderWithResources(
				"helm-ns",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				[]string{"data/atlasdeployment_basic_helm.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("reader").
						WithSecretRef("dbuser-secret-u1").
						AddCustomRole(model.RoleCustomReadWrite, "Ships", "").
						WithAuthDatabase("admin"),
				},
				30006,
				[]func(*model.TestDataProvider){
					actions.HelmDefaultUpgradeResources,
					actions.HelmUpgradeUsersRoleAddAdminUser,
					actions.HelmUpgradeDeleteFirstUser,
				},
			),
			"default",
		),
		Entry("Advanced deployment by helm chart", Label("helm-advanced"),
			model.DataProviderWithResources(
				"helm-advanced",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				[]string{"data/atlasdeployment_advanced_helm.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("reader2").
						WithSecretRef("dbuser-secret-u2").
						AddCustomRole(model.RoleCustomReadWrite, "Ships", "").
						WithAuthDatabase("admin"),
				},
				30014,
				[]func(*model.TestDataProvider){},
			),
			"advanced",
		),
		Entry("Advanced multiregion deployment by helm chart", Label("helm-advanced-multiregion"),
			model.DataProviderWithResources(
				"helm-advanced-multiregion",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				[]string{"data/atlasdeployment_advanced_multi_region_helm.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("reader2").
						WithSecretRef("dbuser-secret-u2").
						AddCustomRole(model.RoleCustomReadWrite, "Ships", "").
						WithAuthDatabase("admin"),
				},
				30015,
				[]func(*model.TestDataProvider){},
			),
			"advanced",
		),
		Entry("Serverless deployment by helm chart", Label("helm-serverless"),
			model.DataProviderWithResources(
				"helm-serverless",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				[]string{"data/atlasdeployment_serverless.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("reader2").
						WithSecretRef("dbuser-secret-u2").
						AddCustomRole(model.RoleCustomReadWrite, "Ships", "").
						WithAuthDatabase("admin"),
				},
				30016,
				[]func(*model.TestDataProvider){},
			),
			"serverless",
		),
	)

	Describe("HELM charts.", Label("helm-wide"), func() {
		It("User can deploy operator namespaces by using HELM", func() {
			By("User creates configuration for a new Project and Deployment", func() {
				data = model.DataProviderWithResources(
					"helm-wide",
					model.AProject{},
					model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
					[]string{"data/atlasdeployment_basic_helm.yaml"},
					[]string{},
					[]model.DBUser{
						*model.NewDBUser("reader2").
							WithSecretRef("dbuser-secret-u2").
							AddCustomRole(model.RoleCustomReadWrite, "Ships", "").
							WithAuthDatabase("admin"),
					},
					30007,
					[]func(*model.TestDataProvider){},
				)
				// helm template has equal ObjectMeta.Name and Spec.Name
				data.Resources.Deployments[0].ObjectMeta.Name = "deployment-from-helm-wide"
				data.Resources.Deployments[0].Spec.DeploymentSpec.Name = "deployment-from-helm-wide"
			})
			By("User use helm for deploying operator", func() {
				helm.InstallOperatorWideSubmodule(data.Resources)
			})
			By("User deploy deployment by helm", func() {
				helm.InstallDeploymentSubmodule(data.Resources)
			})
			waitDeploymentWithChecks(&data)
			deleteDeploymentAndOperator(&data)
		})
	})

	Describe("HELM charts.", Label("helm-update"), func() {
		It("User deploy operator and later deploy new version of the Atlas operator", func() {
			By("User creates configuration for a new Project, Deployment, DBUser", func() {
				data = model.DataProviderWithResources(
					"helm-upgrade",
					model.AProject{},
					model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
					[]string{"data/atlasdeployment_basic_helm.yaml"},
					[]string{},
					[]model.DBUser{
						*model.NewDBUser("admin").
							WithSecretRef("dbuser-secret-u2").
							AddBuildInAdminRole().
							WithAuthDatabase("admin"),
					},
					30010,
					[]func(*model.TestDataProvider){},
				)
				// helm template has equal ObjectMeta.Name and Spec.Name
				data.Resources.Deployments[0].ObjectMeta.Name = "deployment-from-helm-upgrade"
				data.Resources.Deployments[0].Spec.DeploymentSpec.Name = "deployment-from-helm-upgrade"
			})
			By("User use helm for last released version of operator and deploy his resources", func() {
				helm.AddMongoDBRepo()
				helm.InstallOperatorNamespacedFromLatestRelease(data.Resources)
				helm.InstallDeploymentRelease(data.Resources)
				waitDeploymentWithChecks(&data)
			})
			By("User update new released operator", func() {
				backup := true
				data.Resources.Deployments[0].Spec.DeploymentSpec.ProviderBackupEnabled = &backup
				actions.HelmUpgradeChartVersions(&data)
				actions.CheckUsersCanUseOldApp(&data)
			})
			By("Delete Resources", func() {
				deleteDeploymentAndOperator(&data)
			})
		})
	})
})

func waitDeploymentWithChecks(data *model.TestDataProvider) {
	By("Wait creation until is done", func() {
		actions.WaitProjectWithoutGenerationCheck(data)
		resource, err := kube.GetProjectResource(data)
		Expect(err).Should(BeNil())
		data.Resources.ProjectID = resource.Status.ID
		actions.WaitDeploymentWithoutGenerationCheck(data)
	})

	By("Check attributes", func() {
		deployment := data.Resources.Deployments[0]
		switch {
		case deployment.Spec.AdvancedDeploymentSpec != nil:
			advancedDeployment, err := atlasClient.GetAdvancedDeployment(data.Resources.ProjectID, deployment.Spec.AdvancedDeploymentSpec.Name)
			Expect(err).To(BeNil())
			actions.CompareAdvancedDeploymentsSpec(deployment.Spec, *advancedDeployment)
		case deployment.Spec.ServerlessSpec != nil:
			serverlessInstance, err := atlasClient.GetServerlessInstance(data.Resources.ProjectID, deployment.Spec.ServerlessSpec.Name)
			Expect(err).To(BeNil())
			actions.CompareServerlessSpec(deployment.Spec, *serverlessInstance)
		default:
			uDeployment := atlasClient.GetDeployment(data.Resources.ProjectID, data.Resources.Deployments[0].Spec.DeploymentSpec.Name)
			actions.CompareDeploymentsSpec(deployment.Spec, uDeployment)
		}
	})

	By("check database users Attributes", func() {
		Eventually(actions.CheckIfUsersExist(data.Resources), "2m", "10s").Should(BeTrue())
		actions.CheckUsersAttributes(data)
	})

	if !data.SkipAppConnectivityCheck {
		By("Deploy application for user", func() {
			actions.CheckUsersCanUseApp(data)
		})
	}
}

func deleteDeploymentAndOperator(data *model.TestDataProvider) {
	By("Check project, deployment does not exist", func() {
		helm.Uninstall(data.Resources.Deployments[0].Spec.GetDeploymentName(), data.Resources.Namespace)
		Eventually(
			func(g Gomega) bool {
				return atlasClient.IsProjectExists(g, data.Resources.ProjectID)
			},
			"7m", "20s",
		).Should(BeFalse(), "Project and deployment should be deleted from Atlas")
	})

	By("Delete HELM releases", func() {
		helm.UninstallKubernetesOperator(data.Resources)
	})
}
