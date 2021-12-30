package e2e_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	cloud "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/cloud"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/deploy"
	kube "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"

	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

// NOTES
// Feature unavailable in Free and Shared-Tier Clusters
// This feature is not available for M0 free clusters, M2, and M5 clusters.

// tag for test resources "atlas-operator-test" (config.Tag)

type privateEndpoint struct {
	provider string
	region   string
}

var _ = Describe("[privatelink-aws] UserLogin", func() {
	var data model.TestDataProvider

	_ = BeforeEach(func() {
		Eventually(kubecli.GetVersionOutput()).Should(Say(K8sVersion))
		checkUpAWSEnviroment()
		checkUpAzureEnviroment()
	})

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + data.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentGinkgoTestDescription().Failed {
			GinkgoWriter.Write([]byte("Test has been failed. Trying to save logs...\n"))
			utils.SaveToFile(
				fmt.Sprintf("output/%s/operatorDecribe.txt", data.Resources.Namespace),
				[]byte(kubecli.DescribeOperatorPod(data.Resources.Namespace)),
			)
			utils.SaveToFile(
				fmt.Sprintf("output/%s/operator-logs.txt", data.Resources.Namespace),
				kubecli.GetManagerLogs(data.Resources.Namespace),
			)
			actions.SaveTestAppLogs(data.Resources)
			actions.SaveK8sResources(
				[]string{"deploy", "atlasprojects"},
				data.Resources.Namespace,
			)
		}
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test model.TestDataProvider, pe []privateEndpoint) {
			data = test
			privateFlow(test, pe)
		},
		Entry("Test: User has project which was updated with AWS PrivateEndpoint",
			model.NewTestDataProvider(
				"operator-plink-aws-1",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlascluster_backup.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("user1").
						WithSecretRef("dbuser-secret-u1").
						AddBuildInAdminRole(),
				},
				40000,
				[]func(*model.TestDataProvider){},
			),
			[]privateEndpoint{
				{
					provider: "AWS",
					region:   "eu-west-2",
				},
			},
		),
		Entry("Test: User has project which was updated with Azure PrivateEndpoint",
			model.NewTestDataProvider(
				"operator-plink-azure-1",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlascluster_backup.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("user1").
						WithSecretRef("dbuser-secret-u1").
						AddBuildInAdminRole(),
				},
				40000,
				[]func(*model.TestDataProvider){},
			),
			[]privateEndpoint{{
				provider: "AZURE",
				region:   "northeurope",
			}},
		),
	)
})

func privateFlow(userData model.TestDataProvider, requstedPE []privateEndpoint) {
	By("Deploy Project with requested configuration", func() {
		actions.PrepareUsersConfigurations(&userData)
		deploy.NamespacedOperator(&userData)
		actions.DeployProjectAndWait(&userData, "1")
	})

	By("Create Private Link and the rest users resources", func() {
		for _, pe := range requstedPE {
			userData.Resources.Project.WithPrivateLink(provider.ProviderName(pe.provider), pe.region)
		}
		actions.PrepareUsersConfigurations(&userData)
		actions.DeployProject(&userData, "2")
	})

	By("Check if project statuses are updating, get project ID", func() {
		Eventually(kube.GetProjectPEndpointServiceStatus(&userData), "15m", "10s").Should(Equal("True"),
			"Atlasproject status.conditions are not True")
		Eventually(kube.GetReadyProjectStatus(&userData)).Should(Equal("True"),
			"Atlasproject status.conditions are not True")
		Expect(AllPEndpointUpdated(&userData)).Should(BeTrue(),
			"Error: Was created a different amount of endpoints")
		actions.UpdateProjectID(&userData)
	})

	By("Create Endpoint in requested Cloud Provider", func() {
		project, err := kube.GetProjectResource(&userData)
		Expect(err).ShouldNot(HaveOccurred())

		for _, peitem := range project.Status.PrivateEndpoints {
			cloudTest := cloud.CreatePEActions(peitem)
			privateLinkID, ip, err := cloudTest.CreatePrivateEndpoint(userData.Resources.ProjectID)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(privateLinkID).ShouldNot(BeEmpty())
			Eventually(
				func() bool {
					return cloudTest.IsStatusPrivateEndpointPending(privateLinkID)
				},
			).Should(BeTrue())
			userData.Resources.Project.UpdatePrivateLinkID(peitem.Provider, peitem.Region, privateLinkID, ip)
		}
	})

	By("Deploy Changed Projects", func() {
		actions.PrepareUsersConfigurations(&userData)
		actions.DeployProjectAndWait(&userData, "3")
	})

	By("Check statuses", func() {
		Eventually(kube.GetProjectPEndpointStatus(&userData)).Should(Equal("True"), "Condition status 'PrivateEndpointServiceReady' is not'True'")
		Eventually(kube.GetReadyProjectStatus(&userData)).Should(Equal("True"), "Condition status 'Ready' is not 'True'")

		project, err := kube.GetProjectResource(&userData)
		Expect(err).ShouldNot(HaveOccurred())
		for _, peitem := range project.Status.PrivateEndpoints {
			cloudTest := cloud.CreatePEActions(peitem)
			privateEndpointID := userData.Resources.Project.GetPrivateIDByProviderRegion(peitem.Provider, peitem.Region)
			Expect(privateEndpointID).ShouldNot(BeEmpty())
			Eventually(
				func() bool {
					return cloudTest.IsStatusPrivateEndpointAvailable(privateEndpointID)
				},
			).Should(BeTrue())
		}
	})

	By("Delete PE from Clouds", func() {
		project, err := kube.GetProjectResource(&userData)
		Expect(err).ShouldNot(HaveOccurred())
		for _, peitem := range project.Status.PrivateEndpoints {
			cloudTest := cloud.CreatePEActions(peitem)
			privateEndpointID := userData.Resources.Project.GetPrivateIDByProviderRegion(peitem.Provider, peitem.Region)
			Expect(privateEndpointID).ShouldNot(BeEmpty())

			err = cloudTest.DeletePrivateEndpoint(privateEndpointID)
			Expect(err).ShouldNot(HaveOccurred())
		}
	})

	By("Delete Resources, Project with PEService", func() {
		actions.DeleteUserResourcesProject(&userData)
	})
}

func AllPEndpointUpdated(data *model.TestDataProvider) bool {
	result, _ := kube.GetProjectResource(data)
	return len(result.Status.PrivateEndpoints) == len(result.Spec.PrivateEndpoints)
}
