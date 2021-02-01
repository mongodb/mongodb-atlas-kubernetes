package e2e_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"

	cli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Deploy simple cluster", func() {

	It("Release sample all-in-one.yaml should work", func() {
		By("Prepare namespaces")
		namespaceUserResources := uuid.NewRandom().String() // TODO
		namespaceOperator := "mongodb-atlas-kubernetes-system"
		session := cli.Execute("kubectl", "create", "namespace", namespaceUserResources)
		Expect(session).ShouldNot(Say("created"))
		userProjectConfig := cli.LoadUserProjectConfig("data/atlasproject.yaml")
		userClusterConfig := cli.LoadUserClusterConfig("data/atlascluster_basic.yaml")

		By("Check Kubernetes/MongoCLI version\n")
		session = cli.Execute("kubectl", "version")
		Eventually(session).Should(Say(K8sVersion))
		session = cli.Execute("mongocli", "--version")
		Eventually(session).Should(gexec.Exit(0))

		By("Apply All-in-one configuration\n in ")
		session = cli.Execute("kubectl", "apply", "-f", ConfigAll)
		Eventually(session.Wait()).Should(Say("customresourcedefinition.apiextensions.k8s.io/atlasclusters.atlas.mongodb.com"))

		By("Create secret")
		session = cli.Execute("kubectl", "create", "secret", "generic", "my-atlas-key",
			"--from-literal=orgId="+os.Getenv("MCLI_ORG_ID"),
			"--from-literal=publicApiKey="+os.Getenv("MCLI_PUBLIC_API_KEY"),
			"--from-literal=privateApiKey="+os.Getenv("MCLI_PRIVATE_API_KEY"),
			"-n", namespaceUserResources,
		)
		Eventually(session.Wait()).Should(Say("my-atlas-key created"))

		By("Create Sample Project\n")
		session = cli.Execute("kubectl", "apply", "-f", ProjectSample, "-n", namespaceUserResources)
		// Eventually(session).Should(Say("my-project created"))
		Eventually(session.Wait()).Should(Say("atlasproject.atlas.mongodb.com/my-project"))

		By("Sample Cluster\n")
		session = cli.Execute("kubectl", "apply", "-f", ClusterSample, "-n", namespaceUserResources)
		Eventually(session.Wait()).Should(Say("atlascluster-sample created"))

		By("Wait creating and check that it was created")
		projectID := cli.GetProjectID(userProjectConfig.Spec.Name)
		Expect(projectID).ShouldNot(BeNil())
		GinkgoWriter.Write([]byte("projectID = " + projectID))
		Eventually(
			cli.GetPodStatus(namespaceOperator),
			"5m", "3s",
		).Should(Equal("Running"))

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
		session = cli.Execute("kubectl", "apply", "-f", "data/updated_atlascluster_basic.yaml", "-n", namespaceUserResources) // TODO param
		Eventually(session.Wait()).Should(Say("atlascluster-sample configured"))

		By("Wait creation")
		userClusterConfig = cli.LoadUserClusterConfig("data/updated_atlascluster_basic.yaml")
		Expect(projectID).ShouldNot(BeNil())
		Eventually(
			cli.GetClusterStatus(projectID, userClusterConfig.Spec.Name),
			"35m", "1m",
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

		// By("Delete project") //TODO
	})
})
