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

package e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/cli"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/cli/helm"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
)

var _ = Describe("HELM charts", Ordered, FlakeAttempts(2), func() {
	var data model.TestDataProvider
	skipped := false

	_ = BeforeAll(func() {
		cli.Execute("kubectl", "delete", "--ignore-not-found=true", "-f", "../../config/crd/bases").Wait().Out.Contents()
	})

	_ = AfterAll(func() {
		cli.Execute("kubectl", "apply", "-f", "../../config/crd/bases").Wait().Out.Contents()
	})

	_ = BeforeEach(func() {
		imageURL := os.Getenv("IMAGE_URL")
		Expect(imageURL).ShouldNot(BeEmpty(), "SetUP IMAGE_URL")
	})

	_ = AfterEach(func() {
		By("After each.", func() {
			if skipped {
				return
			}
			GinkgoWriter.Write([]byte("\n"))
			GinkgoWriter.Write([]byte("===============================================\n"))
			GinkgoWriter.Write([]byte("Operator namespace: " + data.Resources.Namespace + "\n"))
			GinkgoWriter.Write([]byte("===============================================\n"))
			if CurrentSpecReport().Failed() {
				GinkgoWriter.Write([]byte("Resources wasn't clean\n"))
				namespaceDeployment, err := k8s.GetDeployment("mongodb-atlas-operator", data.Resources.Namespace)
				Expect(err).Should(BeNil())
				namespaceDeploymentJSON, err := json.MarshalIndent(namespaceDeployment, "", "  ")
				Expect(err).Should(BeNil())
				utils.SaveToFile("output/namespace-deployment.json", namespaceDeploymentJSON)

				pod, err := k8s.GetAllDeploymentPods("mongodb-atlas-operator", data.Resources.Namespace)
				Expect(err).Should(BeNil())
				podJSON, err := json.MarshalIndent(pod, "", "  ")
				Expect(err).Should(BeNil())
				utils.SaveToFile("output/namespace-pod.json", podJSON)

				bytes, err := k8s.GetPodLogsByDeployment("mongodb-atlas-operator", config.DefaultOperatorNS, corev1.PodLogOptions{})
				if err != nil {
					GinkgoWriter.Write(fmt.Appendf(nil, "%v\n", err))
				}
				utils.SaveToFile(
					fmt.Sprintf("output/%s/operator-logs-default.txt", data.Resources.Namespace),
					bytes,
				)
				bytes, err = k8s.GetPodLogsByDeployment("mongodb-atlas-operator", data.Resources.Namespace, corev1.PodLogOptions{})
				if err != nil {
					GinkgoWriter.Write(fmt.Appendf(nil, "%v\n", err))
				}
				utils.SaveToFile(
					fmt.Sprintf("output/%s/operator-logs.txt", data.Resources.Namespace),
					bytes,
				)
				actions.SaveProjectsToFile(data.Context, data.K8SClient, data.Resources.Namespace)
				actions.SaveDeploymentsToFile(data.Context, data.K8SClient, data.Resources.Namespace)
				actions.SaveUsersToFile(data.Context, data.K8SClient, data.Resources.Namespace)
				actions.SaveTestAppLogs(data.Resources)
			}
			actions.AfterEachFinalCleanup([]model.TestDataProvider{data})
		})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration", Label("helm-ns"),
		func(ctx SpecContext, test func(context.Context) model.TestDataProvider, deploymentType string) { // deploymentType - probably will be moved later ()
			data = test(ctx)
			GinkgoWriter.Println(data.Resources.KeyName)
			switch deploymentType {
			case "flex":
				data.Resources.Deployments[0].Spec.FlexSpec.Name = data.Resources.KeyName
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
		Entry("Several actions with helm update", Label("focus-helm-ns-flow"),
			func(ctx context.Context) model.TestDataProvider {
				return model.DataProviderWithResources(ctx, "helm-ns", model.AProject{}, model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), []string{"../helper/e2e/data/atlasdeployment_flex.yaml"}, []string{}, []model.DBUser{
					*model.NewDBUser("reader").
						WithSecretRef("dbuser-secret-u1").
						AddCustomRole(model.RoleCustomReadWrite, "Ships", "").
						WithAuthDatabase("admin"),
				}, 30006, []func(*model.TestDataProvider){
					actions.HelmDefaultUpgradeResources,
					actions.HelmUpgradeUsersRoleAddAdminUser,
					actions.HelmUpgradeDeleteFirstUser,
				})
			},
			"flex",
		),
		Entry("Deployment multiregion by helm chart", Label("focus-helm-advanced-multiregion"),
			func(ctx context.Context) model.TestDataProvider {
				return model.DataProviderWithResources(ctx, "helm-advanced-multiregion", model.AProject{}, model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), []string{"../helper/e2e/data/atlasdeployment_advanced_multi_region_helm.yaml"}, []string{}, []model.DBUser{
					*model.NewDBUser("reader2").
						WithSecretRef("dbuser-secret-u2").
						AddCustomRole(model.RoleCustomReadWrite, "Ships", "").
						WithAuthDatabase("admin"),
				}, 30015, []func(*model.TestDataProvider){})
			},
			"advanced",
		),
		Entry("Flex deployment by helm chart", Label("focus-helm-flex"),
			func(ctx context.Context) model.TestDataProvider {
				return model.DataProviderWithResources(ctx, "helm-flex", model.AProject{}, model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), []string{"../helper/e2e/data/atlasdeployment_flex.yaml"}, []string{}, []model.DBUser{
					*model.NewDBUser("reader2").
						WithSecretRef("dbuser-secret-u2").
						AddCustomRole(model.RoleCustomReadWrite, "Ships", "").
						WithAuthDatabase("admin"),
				}, 30016, []func(*model.TestDataProvider){})
			},
			"flex",
		),
	)

	Describe("HELM charts.", Label("helm-wide"), func() {
		It("User can deploy operator namespaces by using HELM", func(ctx SpecContext) {
			By("User creates configuration for a new Project and Deployment", func() {
				data = model.DataProviderWithResources(ctx, "helm-wide", model.AProject{}, model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), []string{"../helper/e2e/data/atlasdeployment_flex.yaml"}, []string{}, []model.DBUser{
					*model.NewDBUser("reader2").
						WithSecretRef("dbuser-secret-u2").
						AddCustomRole(model.RoleCustomReadWrite, "Ships", "").
						WithAuthDatabase("admin"),
				}, 30007, []func(*model.TestDataProvider){})
				// helm template has equal ObjectMeta.Name and Spec.Name
				data.Resources.Deployments[0].ObjectMeta.Name = "deployment-from-helm-wide"
				data.Resources.Deployments[0].Spec.FlexSpec.Name = "deployment-from-helm-wide"
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
		It("User deploy operator and later deploy new version of the Atlas operator", func(ctx SpecContext) {
			By("Check upgrade is actually possible", func() {
				helm.AddMongoDBRepo()
				releasedVersion, err := helm.GetReleasedChartVersion()
				Expect(err).Should(BeNil())
				devMajorVersion, err := helm.GetDevelopmentMayorVersion()
				Expect(err).Should(BeNil())
				releaseMajorVersion := strings.Split(releasedVersion, ".")[0]
				if releaseMajorVersion != devMajorVersion {
					skipped = true
					Skip(fmt.Sprintf("cannot test upgrade from incompatible major release version %q to version %q",
						releaseMajorVersion, devMajorVersion))
				}
			})
			By("User creates configuration for a new Project, Deployment, DBUser", func() {
				data = model.DataProviderWithResources(ctx, "helm-upgrade", model.AProject{}, model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), []string{"../helper/e2e/data/atlasdeployment_flex.yaml"}, []string{}, []model.DBUser{
					*model.NewDBUser("admin").
						WithSecretRef("dbuser-secret-u2").
						AddBuildInAdminRole().
						WithAuthDatabase("admin"),
				}, 30010, []func(*model.TestDataProvider){})
				// helm template has equal ObjectMeta.Name and Spec.Name
				data.Resources.Deployments[0].ObjectMeta.Name = "deployment-from-helm-upgrade"
				data.Resources.Deployments[0].Spec.FlexSpec.Name = "deployment-from-helm-upgrade"
				data.Resources.Deployments[0].Spec.FlexSpec.TerminationProtectionEnabled = true
			})
			By("User use helm for last released version of operator and deploy his resources", func() {
				helm.InstallOperatorNamespacedFromLatestRelease(data.Resources)
				helm.InstallDeploymentRelease(data.Resources)
				waitDeploymentWithChecks(&data)
			})
			By("User update new released operator", func() {
				data.Resources.Deployments[0].Spec.FlexSpec.TerminationProtectionEnabled = false
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
	By("Wait for a Deployment to be created", func() {
		actions.WaitProjectWithoutGenerationCheck(data)
		resource, err := kube.GetProjectResource(data)
		Expect(err).Should(BeNil())
		data.Resources.ProjectID = resource.Status.ID
		actions.WaitDeploymentWithoutGenerationCheck(data)
	})

	By("Check attributes", func() {
		deployment := data.Resources.Deployments[0]
		switch {
		case deployment.Spec.FlexSpec != nil:
			flexInstance, err := atlasClient.GetFlexInstance(data.Resources.ProjectID, deployment.Spec.FlexSpec.Name)
			Expect(err).To(BeNil())
			actions.CompareFlexSpec(deployment.Spec, *flexInstance)
		default:
			uDeployment, err := atlasClient.GetDeployment(data.Resources.ProjectID, data.Resources.Deployments[0].Spec.DeploymentSpec.Name)
			Expect(err).To(BeNil())
			actions.CompareAdvancedDeploymentsSpec(deployment.Spec, *uDeployment)
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
			func(g Gomega) {
				atlasClient.IsProjectExists(g, data.Resources.ProjectID)
			},
			"7m", "20s",
		).Should(Succeed(), "Project and deployment should be deleted from Atlas")
	})

	By("Delete HELM releases", func() {
		helm.UninstallKubernetesOperator(data.Resources)
	})

	By("Uninstall HELM CRDs", func() {
		helm.UninstallCRD(data.Resources)
	})
}
