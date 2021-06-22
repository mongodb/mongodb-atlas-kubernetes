package e2e_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gstruct"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	actions "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kube"
	. "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

var _ = Describe("[cluster-wide] Users (Norton and Nimnul) can work with one Cluster wide operator", func() {
	var NortonData, NimnulData model.TestDataProvider
	commonClusterName := "megacluster"

	_ = BeforeEach(func() {
		Eventually(kube.GetVersionOutput()).Should(Say(K8sVersion))
		By("User Install CRD, cluster wide Operator", func() {
			Eventually(kube.Apply(ConfigAll)).Should(
				Say("customresourcedefinition.apiextensions.k8s.io/atlasclusters.atlas.mongodb.com"),
			)
			Eventually(
				kube.GetPodStatus(DefaultOperatorNS),
				"5m", "3s",
			).Should(Equal("Running"), "The operator should successfully run")
		})
	})

	var _ = AfterEach(func() {
		By("AfterEach. clean-up", func() {
			if CurrentGinkgoTestDescription().Failed {
				GinkgoWriter.Write([]byte("Resources wasn't clean"))
				utils.SaveToFile(
					"output/operator-logs.txt",
					kube.GetManagerLogs(DefaultOperatorNS),
				)
				actions.SaveK8sResources(
					[]string{"deploy"},
					DefaultOperatorNS,
				)
				actions.SaveK8sResources(
					[]string{"atlasclusters", "atlasdatabaseusers", "atlasprojects"},
					NortonData.Resources.Namespace,
				)
				actions.SaveK8sResources(
					[]string{"atlasclusters", "atlasdatabaseusers", "atlasprojects"},
					NimnulData.Resources.Namespace,
				)
				actions.SaveTestAppLogs(NortonData.Resources)
				actions.SaveTestAppLogs(NimnulData.Resources)
			} else {
				actions.AfterEachFinalCleanup([]model.TestDataProvider{NortonData, NimnulData})
			}
		})
	})

	// (Consider Shared Clusters when E2E tests could conflict with each other)
	It("Deploy cluster wide operator and create resources in each of them", func() {
		By("Users can create clusters with the same name", func() {
			NortonData = model.NewTestDataProvider(
				"norton-wide",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlascluster_backup.yaml"},
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
			NimnulData = model.NewTestDataProvider(
				"nimnul-wide",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlascluster_basic.yaml"},
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
			NortonData.Resources.Clusters[0].ObjectMeta.Name = "norton-cluster"
			NortonData.Resources.Clusters[0].Spec.Name = commonClusterName
			NimnulData.Resources.Clusters[0].ObjectMeta.Name = "nimnul-cluster"
			NimnulData.Resources.Clusters[0].Spec.Name = commonClusterName
		})

		By("Deploy users resorces", func() {
			actions.PrepareUsersConfigurations(&NortonData)
			actions.PrepareUsersConfigurations(&NimnulData)
			actions.DeployUserResourcesAction(&NortonData)
			actions.DeployUserResourcesAction(&NimnulData)
		})

		By("Operator working with right cluster if one of the user update configuration", func() {
			NortonData.Resources.Clusters[0].Spec.Labels = []v1.LabelSpec{{Key: "something", Value: "awesome"}}
			utils.SaveToFile(
				NortonData.Resources.Clusters[0].ClusterFileName(NortonData.Resources),
				utils.JSONToYAMLConvert(NortonData.Resources.Clusters[0]),
			)
			kube.Apply(NortonData.Resources.Clusters[0].ClusterFileName(NortonData.Resources), "-n", NortonData.Resources.Namespace)
			actions.WaitCluster(NortonData.Resources, "2")

			By("Norton cluster has labels", func() {
				Expect(
					kube.GetClusterResource(NortonData.Resources.Namespace, NortonData.Resources.Clusters[0].GetClusterNameResource()).Spec.Labels[0],
				).To(MatchFields(IgnoreExtras, Fields{
					"Key":   Equal("something"),
					"Value": Equal("awesome"),
				}))
			})

			By("Nimnul cluster does not have labels", func() {
				Eventually(
					kube.GetClusterResource(NimnulData.Resources.Namespace, NimnulData.Resources.Clusters[0].GetClusterNameResource()).Spec.Labels,
				).Should(BeNil())
			})
		})

		By("Delete Norton/Nimnul Resources", func() {
			actions.DeleteUserResources(&NortonData)
			actions.DeleteUserResources(&NimnulData)
		})
	})
})
