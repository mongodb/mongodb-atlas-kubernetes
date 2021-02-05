package e2e_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	// . "github.com/onsi/ginkgo/extensions/table"
	"github.com/pborman/uuid"

	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
	cli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

var _ = Describe("Deploy simple cluster", func() {



	It("Release sample all-in-one.yaml should work", func() {
		By("Prepare namespaces and project configuration") // TODO clusters/keys will be a bit later
		id := uuid.NewRandom().String()
		// TODO move it
		namespaceUserResources := "e2e-" + id
		namespaceOperator := "mongodb-atlas-kubernetes-system"
		keyName := "my-atlas-key"
		pName := id
		k8sProjectName := "k-" + id
		ProjectSampleFile := "data/"+pName+".yaml"
		ClusterSampleFile := "data/atlascluster_basic.yaml" // TODO put it to dataprovider
		GinkgoWriter.Write([]byte(namespaceUserResources))
		session := cli.Execute("kubectl", "create", "namespace", namespaceUserResources)
		Expect(session.Wait()).Should(Say("created"))

		project := utils.NewProject().ProjectName(pName).SecretRef(keyName).CompleteK8sConfig(k8sProjectName)
		utils.SaveToFile(ProjectSampleFile, project)

		userProjectConfig := cli.LoadUserProjectConfig(ProjectSampleFile)
		userClusterConfig := cli.LoadUserClusterConfig(ClusterSampleFile)
		userClusterConfig.Spec.Project.Name = k8sProjectName
		clusterData, _ := utils.JSONToYAMLConvert(userClusterConfig)
		utils.SaveToFile(ClusterSampleFile, clusterData)

		By("Check Kubernetes/MongoCLI version\n")
		session = cli.Execute("kubectl", "version")
		Eventually(session).Should(Say(K8sVersion))
		session = cli.Execute("mongocli", "--version")
		Eventually(session).Should(gexec.Exit(0))

		By("Apply All-in-one configuration\n in ")
		session = cli.Execute("kubectl", "apply", "-f", ConfigAll)
		Eventually(session.Wait()).Should(Say("customresourcedefinition.apiextensions.k8s.io/atlasclusters.atlas.mongodb.com"))
		Eventually(
			cli.GetPodStatus(namespaceOperator),
			"5m", "3s",
		).Should(Equal("Running"))

		By("Create secret")
		session = cli.Execute("kubectl", "create", "secret", "generic", "my-atlas-key",
			"--from-literal=orgId="+os.Getenv("MCLI_ORG_ID"),
			"--from-literal=publicApiKey="+os.Getenv("MCLI_PUBLIC_API_KEY"),
			"--from-literal=privateApiKey="+os.Getenv("MCLI_PRIVATE_API_KEY"),
			"-n", namespaceUserResources,
		)
		Eventually(session.Wait()).Should(Say("my-atlas-key created"))

		By("Create Sample Project\n")
		session = cli.Execute("kubectl", "apply", "-f", ProjectSampleFile, "-n", namespaceUserResources)
		Eventually(session.Wait()).Should(Say("atlasproject.atlas.mongodb.com/"+k8sProjectName + " created"))

		By("Sample Cluster\n")
		session = cli.Execute("kubectl", "apply", "-f", ClusterSample, "-n", namespaceUserResources)
		Eventually(session.Wait()).Should(Say("created"))

		By("Wait creating and check that it was created")
		Eventually(cli.GetStatus(namespaceUserResources, "atlasproject.atlas.mongodb.com/"+k8sProjectName))
		Eventually(cli.GetGeneration(namespaceUserResources)).Should(Equal("1"))
		Eventually(
			cli.IsProjectExist(userProjectConfig.Spec.Name),
			"7m", "10s",
		).Should(BeTrue())

		projectID := cli.GetProjectID(userProjectConfig.Spec.Name)
		GinkgoWriter.Write([]byte("projectID = " + projectID))
		Eventually(
			cli.IsClusterExist(projectID, userClusterConfig.Spec.Name),
			"5m", "1s",
		).Should(BeTrue())

		Eventually(
			cli.GetClusterStatus(projectID, userClusterConfig.Spec.Name),
			"35m", "1m",
		).Should(Equal("IDLE"))

		By("check cluster Attribute") // TODO ...
		cluster := cli.GetClustersInfo(projectID, userClusterConfig.Spec.Name)
		Expect(
			cluster.ProviderSettings.InstanceSizeName,
		).Should(Equal(userClusterConfig.Spec.ProviderSettings.InstanceSizeName))
		Expect(
			cluster.ProviderSettings.ProviderName,
		).Should(Equal(userClusterConfig.Spec.ProviderSettings.ProviderName))
		Expect(
			cluster.ProviderSettings.RegionName,
		).Should(Equal(userClusterConfig.Spec.ProviderSettings.RegionName))

		By("Update cluster\n")
		// userClusterConfig := cli.LoadUserClusterConfig(ClusterSampleFile)
		userClusterConfig.Spec.ProviderSettings.InstanceSizeName = "M20"
		clusterData, _ = utils.JSONToYAMLConvert(userClusterConfig)
		utils.SaveToFile(ClusterSampleFile, clusterData)

		session = cli.Execute("kubectl", "apply", "-f", ClusterSampleFile, "-n", namespaceUserResources) // TODO param
		Eventually(session.Wait()).Should(Say("atlascluster-sample configured"))

		By("Wait creation")
		Eventually(cli.GetGeneration(namespaceUserResources)).Should(Equal("2"))
		Eventually(
			cli.GetClusterStatus(projectID, userClusterConfig.Spec.Name),
			"50m", "1m", // max time 35 ??
		).Should(Equal("IDLE"))

		uCluster := cli.GetClustersInfo(projectID, userClusterConfig.Spec.Name)
		Expect(
			uCluster.ProviderSettings.InstanceSizeName,
		).Should(Equal(
			userClusterConfig.Spec.ProviderSettings.InstanceSizeName,
		))

		By("Delete cluster")
		session = cli.Execute("kubectl", "delete", "-f", "data/updated_atlascluster_basic.yaml", "-n", namespaceUserResources)
		Eventually(session.Wait("7m")).Should(gexec.Exit(0))
		Eventually(
			cli.IsClusterExist(projectID, userClusterConfig.Spec.Name),
			"10m", "1m",
		).Should(BeFalse())


		// By("Delete project") // TODO
		// session = cli.Execute("kubectl", "delete", "-f", "data/atlasproject.yaml", "-n", namespaceUserResources)
		// Eventually(
		// 	cli.IsProjectExist(userProjectConfig.Spec.Name),
		// 	"5m", "20s",
		// ).Should(BeFalse())
	})
})
