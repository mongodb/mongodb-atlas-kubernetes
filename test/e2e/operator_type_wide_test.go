package e2e_test

import (
	"context"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/k8s"

	"k8s.io/apimachinery/pkg/types"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	. "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

var _ = Describe("Users (Norton and Nimnul) can work with one Deployment wide operator", Label("deployment-wide"), func() {
	var NortonData, NimnulData model.TestDataProvider
	commonDeploymentName := "megadeployment"
	k8sClient, err := k8s.CreateNewClient()
	Expect(err).To(BeNil())
	ctx := context.Background()

	_ = BeforeEach(func() {
		Eventually(kubecli.GetVersionOutput()).Should(Say(K8sVersion))
		By("User Install CRD, deployment wide Operator", func() {
			Eventually(kubecli.Apply(DefaultDeployConfig)).Should(
				Say("customresourcedefinition.apiextensions.k8s.io/atlasdeployments.atlas.mongodb.com"),
			)
			Eventually(
				func() string {
					status, err := k8s.GetPodStatus(ctx, k8sClient, DefaultOperatorNS)
					if err != nil {
						return ""
					}
					return status
				},
				"5m", "3s",
			).Should(Equal("Running"), "The operator should successfully run")
		})
	})

	_ = AfterEach(func() {
		By("AfterEach. clean-up", func() {
			if CurrentSpecReport().Failed() {
				GinkgoWriter.Write([]byte("Resources wasn't clean"))
				utils.SaveToFile(
					"output/operator-logs.txt",
					kubecli.GetManagerLogs(DefaultOperatorNS),
				)
				actions.SaveK8sResources(
					[]string{"deploy"},
					DefaultOperatorNS,
				)
				actions.SaveProjectsToFile(NortonData.Context, NortonData.K8SClient, NortonData.Resources.Namespace)
				actions.SaveK8sResources(
					[]string{"atlasdeployments", "atlasdatabaseusers", "atlasprojects"},
					NortonData.Resources.Namespace,
				)
				actions.SaveProjectsToFile(NimnulData.Context, NimnulData.K8SClient, NimnulData.Resources.Namespace)
				actions.SaveK8sResources(
					[]string{"atlasdeployments", "atlasdatabaseusers", "atlasprojects"},
					NimnulData.Resources.Namespace,
				)
				actions.SaveTestAppLogs(NortonData.Resources)
				actions.SaveTestAppLogs(NimnulData.Resources)
			}
			actions.AfterEachFinalCleanup([]model.TestDataProvider{NortonData, NimnulData})
		})
	})

	// (Consider Shared Deployments when E2E tests could conflict with each other)
	It("Deploy deployment wide operator and create resources in each of them", func() {
		By("Users can create deployments with the same name", func() {
			NortonData = model.DataProviderWithResources(
				"norton-wide",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				[]string{"data/atlasdeployment_backup.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("reader2").
						WithSecretRef("dbuser-secret-u2").
						AddCustomRole(model.RoleCustomReadWrite, "Ships", "").
						WithAuthDatabase("admin"),
				},
				30008,
				[]func(*model.TestDataProvider){},
			)
			NimnulData = model.DataProviderWithResources(
				"nimnul-wide",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				[]string{"data/atlasdeployment_basic.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("reader2").
						WithSecretRef("dbuser-secret-u2").
						AddCustomRole(model.RoleCustomReadWrite, "Ships", "").
						WithAuthDatabase("admin"),
				},
				30009,
				[]func(*model.TestDataProvider){},
			)
			NortonData.Resources.Deployments[0].ObjectMeta.Name = "norton-deployment"
			NortonData.Resources.Deployments[0].Spec.DeploymentSpec.Name = commonDeploymentName
			NimnulData.Resources.Deployments[0].ObjectMeta.Name = "nimnul-deployment"
			NimnulData.Resources.Deployments[0].Spec.DeploymentSpec.Name = commonDeploymentName
		})

		By("Deploy users resorces", func() {
			actions.PrepareUsersConfigurations(&NortonData)
			actions.PrepareUsersConfigurations(&NimnulData)
			actions.DeployUserResourcesAction(&NortonData)
			actions.DeployUserResourcesAction(&NimnulData)
		})

		By("Operator working with right deployment if one of the user update configuration", func() {
			NortonData.Resources.Deployments[0].Spec.DeploymentSpec.Labels = []common.LabelSpec{{Key: "something", Value: "awesome"}}
			utils.SaveToFile(
				NortonData.Resources.Deployments[0].DeploymentFileName(NortonData.Resources),
				utils.JSONToYAMLConvert(NortonData.Resources.Deployments[0]),
			)
			kubecli.Apply(NortonData.Resources.Deployments[0].DeploymentFileName(NortonData.Resources), "-n", NortonData.Resources.Namespace)
			actions.WaitDeployment(&NortonData, 2)

			By("Norton deployment has labels", func() {
				deployment := &v1.AtlasDeployment{}
				Expect(NortonData.K8SClient.Get(NortonData.Context, types.NamespacedName{Name: NortonData.Resources.Deployments[0].ObjectMeta.Name, Namespace: NortonData.Resources.Namespace}, deployment)).To(Succeed())
				Expect(len(deployment.Spec.DeploymentSpec.Labels)).To(Equal(1))
				Expect(deployment.Spec.DeploymentSpec.Labels[0]).To(Equal(common.LabelSpec{Key: "something", Value: "awesome"}))
			})

			By("Nimnul deployment does not have labels", func() {
				deployment := &v1.AtlasDeployment{}
				Expect(NimnulData.K8SClient.Get(NimnulData.Context, types.NamespacedName{Name: NimnulData.Resources.Deployments[0].ObjectMeta.Name, Namespace: NimnulData.Resources.Namespace}, deployment)).To(Succeed())
				Expect(len(deployment.Spec.DeploymentSpec.Labels)).To(Equal(0))
			})
		})

		By("Delete Norton/Nimnul Resources", func() {
			actions.DeleteUserResources(&NortonData)
			actions.DeleteUserResources(&NimnulData)
		})
	})
})
