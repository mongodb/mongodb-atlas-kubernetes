package e2e_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onsi/gomega/gbytes"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

var _ = Describe("Users (Norton and Nimnul) can work with one Cluster wide operator", func() {

	var NortonSpec, NimnulSpec userInputs
	commonClusterName := "MegaCluster"
	var _ = BeforeEach(func() {
		By("User Install CRD, cluster wide Operator", func() {
			Eventually(kube.Apply(ConfigAll)).Should(
				Say("customresourcedefinition.apiextensions.k8s.io/atlasclusters.atlas.mongodb.com"),
			)
			Eventually(
				kube.GetPodStatus(defaultOperatorNS),
				"5m", "3s",
			).Should(Equal("Running"))
		})
	})
	var _ = AfterEach(func() {
		By("Delete clusters", func() {
			kube.Delete(NortonSpec.clusters[0].ClusterFileName(), "-n", NortonSpec.GenNamespace())
			kube.Delete(NimnulSpec.clusters[0].ClusterFileName(), "-n", NimnulSpec.GenNamespace())
			// do not wait it
			// Eventually(kube.DeleteNamespace(NortonSpec.GenNamespace())).Should(Say("deleted"))
			// Eventually(kube.DeleteNamespace(NimnulSpec.GenNamespace())).Should(Say("deleted"))
			// // mongocli.DeleteCluster(ID, "cluster45") // TODO struct
		})
	})

	It("Deploy cluster wide operator and create resources in each of them", func() {
		// (Consider Shared Clusters when E2E tests could conflict with each other)
		By("Users can create clusters with the same name", func() {
			NortonSpec = userInputs{
				projectName: utils.GenUniqID(),
				keyName:     "norton-key",
			}
			NimnulSpec = userInputs{
				projectName: utils.GenUniqID(),
				keyName:     "nimnul-key",
			}
			By("User 1 - Norton - Creates configuration for a new Project and Cluster named: "+commonClusterName, func() {
				project := utils.NewProject().
					ProjectName(NortonSpec.projectName).
					SecretRef(NortonSpec.keyName).
					CompleteK8sConfig(NortonSpec.ProjectK8sName())
				utils.SaveToFile(FilePathTo(NortonSpec.projectName), project)
				NortonSpec.clusters = append(NortonSpec.clusters, utils.LoadUserClusterConfig(ClusterSample))
				NortonSpec.clusters[0].Spec.Project.Name = NortonSpec.ProjectK8sName()
				NortonSpec.clusters[0].ObjectMeta.Name = "norton-cluster"
				utils.SaveToFile(
					NortonSpec.clusters[0].ClusterFileName(),
					utils.JSONToYAMLConvert(NortonSpec.clusters[0]),
				)

				By("Apply Nortons configuration", func() {
					kube.CreateNamespace(NortonSpec.GenNamespace())
					kube.CreateKey(NortonSpec.keyName, NortonSpec.GenNamespace())
					kube.Apply(FilePathTo(NortonSpec.projectName), "-n", NortonSpec.GenNamespace())
					kube.Apply(NortonSpec.clusters[0].ClusterFileName(), "-n", NortonSpec.GenNamespace())
				})

			})
			By("User 2 - Nimnul - Creates configuration for a new Project and Cluster named: "+commonClusterName, func() {
				project := utils.NewProject().
					ProjectName(NimnulSpec.projectName).
					SecretRef(NimnulSpec.keyName).
					CompleteK8sConfig(NimnulSpec.ProjectK8sName())
				utils.SaveToFile(NimnulSpec.UserProjectFile(), project)
				NimnulSpec.clusters = append(NimnulSpec.clusters, utils.LoadUserClusterConfig(ClusterSample))
				NimnulSpec.clusters[0].Spec.Project.Name = NimnulSpec.ProjectK8sName()
				NimnulSpec.clusters[0].ObjectMeta.Name = "nimnul-cluster"
				utils.SaveToFile(
					NimnulSpec.clusters[0].ClusterFileName(),
					utils.JSONToYAMLConvert(NimnulSpec.clusters[0]),
				)

				By("Apply Nortons configuration", func() {
					kube.CreateNamespace(NimnulSpec.GenNamespace())
					kube.CreateKey(NimnulSpec.keyName, NimnulSpec.GenNamespace())
					kube.Apply(FilePathTo(NimnulSpec.projectName), "-n", NimnulSpec.GenNamespace())
					kube.Apply(NimnulSpec.clusters[0].ClusterFileName(), "-n", NimnulSpec.GenNamespace())
				})
			})
		})

		By("Wait creation projects/clusters", func() {
			// projects Norton
			waitProject(NortonSpec, "1")
			NortonSpec.projectID = kube.GetProjectResource(NortonSpec.GenNamespace(), NortonSpec.GetFullK8sAtlasProjectName()).Status.ID

			// projects Nimnul
			waitProject(NimnulSpec, "1")
			NimnulSpec.projectID = kube.GetProjectResource(NimnulSpec.GenNamespace(), NimnulSpec.GetFullK8sAtlasProjectName()).Status.ID

			waitCluster(NortonSpec, "1")
			waitCluster(NimnulSpec, "1")
		})

		By("Check connection strings", func() { // TODO app(?)
			Eventually(kube.GetClusterResource(NortonSpec.GenNamespace(), NortonSpec.clusters[0].GetClusterNameResource()).
				Status.ConnectionStrings.StandardSrv,
			).ShouldNot(BeNil())

			Eventually(kube.GetClusterResource(NimnulSpec.GenNamespace(), NimnulSpec.clusters[0].GetClusterNameResource()).
				Status.ConnectionStrings.StandardSrv,
			).ShouldNot(BeNil())
		})

		By("Operator working with right cluster if one of the user update configuration", func() {
			NortonSpec.clusters[0].Spec.Labels = []v1.LabelSpec{{Key: "something", Value: "awesome"}}
			utils.SaveToFile(
				NortonSpec.clusters[0].ClusterFileName(),
				utils.JSONToYAMLConvert(NortonSpec.clusters[0]),
			)
			kube.Apply(NortonSpec.clusters[0].ClusterFileName(), "-n", NortonSpec.GenNamespace())
			waitCluster(NortonSpec, "2")

			Eventually(
				kube.GetClusterResource(NimnulSpec.GenNamespace(), NimnulSpec.clusters[0].GetClusterNameResource()).Labels,
			).Should(BeNil())
		})
	})
})
