package e2e_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gstruct"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	actions "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/mongocli"
	. "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

var _ = Describe("[cluster-wide] Users (Norton and Nimnul) can work with one Cluster wide operator", func() {
	var NortonSpec, NimnulSpec model.UserInputs
	commonClusterName := "megacluster"

	_ = BeforeEach(func() {
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
					NortonSpec.Namespace,
				)
				actions.SaveK8sResources(
					[]string{"atlasclusters", "atlasdatabaseusers", "atlasprojects"},
					NimnulSpec.Namespace,
				)
			} else {
				Eventually(kube.DeleteNamespace(NortonSpec.Namespace)).Should(Say("deleted"), "Cant delete namespace after testing")
				Eventually(kube.DeleteNamespace(NimnulSpec.Namespace)).Should(Say("deleted"), "Cant delete namespace after testing")
			}
		})
	})

	loadClustersAndApplyConfiguration := func(spec model.UserInputs, name string) model.UserInputs {
		utils.SaveToFile(spec.ProjectPath, spec.Project.ConvertByte())
		spec.Clusters = append(spec.Clusters, model.LoadUserClusterConfig(ClusterSample))
		spec.Clusters[0].Spec.Project.Name = spec.Project.GetK8sMetaName()
		spec.Clusters[0].ObjectMeta.Name = name + "cluster"
		utils.SaveToFile(
			spec.Clusters[0].ClusterFileName(spec),
			utils.JSONToYAMLConvert(spec.Clusters[0]),
		)

		By("Apply "+name+" configuration", func() {
			kube.CreateNamespace(spec.Namespace)
			kube.CreateApiKeySecret(spec.KeyName, spec.Namespace)
			kube.Apply(spec.ProjectPath, "-n", spec.Namespace)
			kube.Apply(spec.Clusters[0].ClusterFileName(spec), "-n", spec.Namespace)
		})
		return spec
	}

	// (Consider Shared Clusters when E2E tests could conflict with each other)
	It("Deploy cluster wide operator and create resources in each of them", func() {
		By("Users can create clusters with the same name", func() {
			By("User 1 - Norton - Creates configuration for a new Project and Cluster named: " + commonClusterName)
			NortonSpec = loadClustersAndApplyConfiguration(model.NewUserInputs("norton-key", nil), "norton")
			By("User 2 - Nimnul - Creates configuration for a new Project and Cluster named: " + commonClusterName)
			NimnulSpec = loadClustersAndApplyConfiguration(model.NewUserInputs("nimnul-key", nil), "nimnul")
		})

		By("Wait creation projects/clusters", func() {
			// projects Norton
			actions.WaitProject(NortonSpec, "1")
			NortonSpec.ProjectID = kube.GetProjectResource(NortonSpec.Namespace, NortonSpec.K8sFullProjectName).Status.ID

			// projects Nimnul
			actions.WaitProject(NimnulSpec, "1")
			NimnulSpec.ProjectID = kube.GetProjectResource(NimnulSpec.Namespace, NimnulSpec.K8sFullProjectName).Status.ID

			actions.WaitCluster(NortonSpec, "1")
			actions.WaitCluster(NimnulSpec, "1")
		})

		By("Check connection strings", func() {
			Eventually(kube.GetClusterResource(NortonSpec.Namespace, NortonSpec.Clusters[0].GetClusterNameResource()).
				Status.ConnectionStrings.StandardSrv,
			).ShouldNot(BeNil())

			Eventually(kube.GetClusterResource(NimnulSpec.Namespace, NimnulSpec.Clusters[0].GetClusterNameResource()).
				Status.ConnectionStrings.StandardSrv,
			).ShouldNot(BeNil())
		})

		By("Operator working with right cluster if one of the user update configuration", func() {
			NortonSpec.Clusters[0].Spec.Labels = []v1.LabelSpec{{Key: "something", Value: "awesome"}}
			utils.SaveToFile(
				NortonSpec.Clusters[0].ClusterFileName(NortonSpec),
				utils.JSONToYAMLConvert(NortonSpec.Clusters[0]),
			)
			kube.Apply(NortonSpec.Clusters[0].ClusterFileName(NortonSpec), "-n", NortonSpec.Namespace)
			actions.WaitCluster(NortonSpec, "2")

			By("Norton cluster has labels", func() {
				Expect(
					kube.GetClusterResource(NortonSpec.Namespace, NortonSpec.Clusters[0].GetClusterNameResource()).Spec.Labels[0],
				).To(MatchFields(IgnoreExtras, Fields{
					"Key":   Equal("something"),
					"Value": Equal("awesome"),
				}))
			})

			By("Nimnul cluster does not have labels", func() {
				Eventually(
					kube.GetClusterResource(NimnulSpec.Namespace, NimnulSpec.Clusters[0].GetClusterNameResource()).Spec.Labels,
				).Should(BeNil())
			})
		})

		By("Delete clusters", func() {
			kube.Delete(NortonSpec.Clusters[0].ClusterFileName(NortonSpec), "-n", NortonSpec.Namespace)
			Eventually(
				actions.CheckIfClusterExist(NortonSpec),
				"10m", "1m",
			).Should(BeFalse(), "Norton Cluster should be deleted from Atlas")

			kube.Delete(NimnulSpec.Clusters[0].ClusterFileName(NimnulSpec), "-n", NimnulSpec.Namespace)
			Eventually(
				actions.CheckIfClusterExist(NimnulSpec),
				"10m", "1m",
			).Should(BeFalse(), "Nimnuls Cluster should be deleted from Atlas")
		})

		By("Delete project", func() {
			kube.Delete(NortonSpec.ProjectPath, "-n", NortonSpec.Namespace)
			Eventually(
				func() bool {
					return mongocli.IsProjectInfoExist(NortonSpec.ProjectID)
				},
				"5m", "20s",
			).Should(BeFalse(), "Nortons Project should be deleted from Atlas")

			kube.Delete(NimnulSpec.ProjectPath, "-n", NimnulSpec.Namespace)
			Eventually(
				func() bool {
					return mongocli.IsProjectInfoExist(NimnulSpec.ProjectID)
				},
				"5m", "20s",
			).Should(BeFalse(), "Nimnuls Project should be deleted from Atlas")
		})
	})
})
