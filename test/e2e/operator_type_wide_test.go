package e2e_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gstruct"

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
			).Should(Equal("Running"), "The operator should successfully run")
		})
	})

	var _ = AfterEach(func() {
		By("Delete clusters", func() {
			if CurrentGinkgoTestDescription().Failed {
				GinkgoWriter.Write([]byte("Resources wasn't clean"))
				kube.GetManagerLogs("mongodb-atlas-system")
				// mongocli.DeleteCluster(ID, "cluster45") // TODO struct
			} else {
				kube.Delete(NortonSpec.clusters[0].ClusterFileName(), "-n", NortonSpec.namespace)
				kube.Delete(NimnulSpec.clusters[0].ClusterFileName(), "-n", NimnulSpec.namespace)
				// do not wait it
				// Eventually(kube.DeleteNamespace(NortonSpec.GenNamespace())).Should(Say("deleted"))
				// Eventually(kube.DeleteNamespace(NimnulSpec.GenNamespace())).Should(Say("deleted"))
			}
		})
	})

	It("Deploy cluster wide operator and create resources in each of them", func() {
		// (Consider Shared Clusters when E2E tests could conflict with each other)
		By("Users can create clusters with the same name", func() {
			NortonSpec = NewUserInputs("norton-key")
			NimnulSpec = NewUserInputs("nimnul-key")
			By("User 1 - Norton - Creates configuration for a new Project and Cluster named: "+commonClusterName, func() {
				project := utils.NewProject().
					ProjectName(NortonSpec.projectName).
					SecretRef(NortonSpec.keyName).
					CompleteK8sConfig(NortonSpec.k8sProjectName)
				utils.SaveToFile(FilePathTo(NortonSpec.projectName), project)
				NortonSpec.clusters = append(NortonSpec.clusters, utils.LoadUserClusterConfig(ClusterSample))
				NortonSpec.clusters[0].Spec.Project.Name = NortonSpec.k8sProjectName
				NortonSpec.clusters[0].ObjectMeta.Name = "norton-cluster"
				utils.SaveToFile(
					NortonSpec.clusters[0].ClusterFileName(),
					utils.JSONToYAMLConvert(NortonSpec.clusters[0]),
				)

				By("Apply Nortons configuration", func() {
					kube.CreateNamespace(NortonSpec.namespace)
					kube.CreateKeySecret(NortonSpec.keyName, NortonSpec.namespace)
					kube.Apply(FilePathTo(NortonSpec.projectName), "-n", NortonSpec.namespace)
					kube.Apply(NortonSpec.clusters[0].ClusterFileName(), "-n", NortonSpec.namespace)
				})

			})
			By("User 2 - Nimnul - Creates configuration for a new Project and Cluster named: "+commonClusterName, func() {
				project := utils.NewProject().
					ProjectName(NimnulSpec.projectName).
					SecretRef(NimnulSpec.keyName).
					CompleteK8sConfig(NimnulSpec.k8sProjectName)
				utils.SaveToFile(NimnulSpec.projectPath, project)
				NimnulSpec.clusters = append(NimnulSpec.clusters, utils.LoadUserClusterConfig(ClusterSample))
				NimnulSpec.clusters[0].Spec.Project.Name = NimnulSpec.k8sProjectName
				NimnulSpec.clusters[0].ObjectMeta.Name = "nimnul-cluster"
				utils.SaveToFile(
					NimnulSpec.clusters[0].ClusterFileName(),
					utils.JSONToYAMLConvert(NimnulSpec.clusters[0]),
				)

				By("Apply Nortons configuration", func() {
					kube.CreateNamespace(NimnulSpec.namespace)
					kube.CreateKeySecret(NimnulSpec.keyName, NimnulSpec.namespace)
					kube.Apply(FilePathTo(NimnulSpec.projectName), "-n", NimnulSpec.namespace)
					kube.Apply(NimnulSpec.clusters[0].ClusterFileName(), "-n", NimnulSpec.namespace)
				})
			})
		})

		By("Wait creation projects/clusters", func() {
			// projects Norton
			waitProject(NortonSpec, "1")
			NortonSpec.projectID = kube.GetProjectResource(NortonSpec.namespace, NortonSpec.k8sFullProjectName).Status.ID

			// projects Nimnul
			waitProject(NimnulSpec, "1")
			NimnulSpec.projectID = kube.GetProjectResource(NimnulSpec.namespace, NimnulSpec.k8sFullProjectName).Status.ID

			waitCluster(NortonSpec, "1")
			waitCluster(NimnulSpec, "1")
		})

		By("Check connection strings", func() { // TODO app(?)
			Eventually(kube.GetClusterResource(NortonSpec.namespace, NortonSpec.clusters[0].GetClusterNameResource()).
				Status.ConnectionStrings.StandardSrv,
			).ShouldNot(BeNil())

			Eventually(kube.GetClusterResource(NimnulSpec.namespace, NimnulSpec.clusters[0].GetClusterNameResource()).
				Status.ConnectionStrings.StandardSrv,
			).ShouldNot(BeNil())
		})

		By("Operator working with right cluster if one of the user update configuration", func() {
			NortonSpec.clusters[0].Spec.Labels = []v1.LabelSpec{{Key: "something", Value: "awesome"}}
			utils.SaveToFile(
				NortonSpec.clusters[0].ClusterFileName(),
				utils.JSONToYAMLConvert(NortonSpec.clusters[0]),
			)
			kube.Apply(NortonSpec.clusters[0].ClusterFileName(), "-n", NortonSpec.namespace)
			waitCluster(NortonSpec, "2")

			By("Norton cluster has labels", func() {
				Expect(
					kube.GetClusterResource(NortonSpec.namespace, NortonSpec.clusters[0].GetClusterNameResource()).Spec.Labels[0],
				).To(MatchFields(IgnoreExtras, Fields{
					"Key":   Equal("something"),
					"Value": Equal("awesome"),
				}))
			})

			By("Nimnul cluster does not have labels", func() {
				Eventually(
					kube.GetClusterResource(NimnulSpec.namespace, NimnulSpec.clusters[0].GetClusterNameResource()).Spec.Labels,
				).Should(BeNil())
			})
		})
	})
})
