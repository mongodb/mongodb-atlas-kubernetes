package e2e_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
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

// AWS NOTES: reserved VPC in eu-west-2, eu-south-1, us-east-1 (due to limitation no more 4 VPC per region)

type privateEndpoint struct {
	provider string
	region   string
}

var _ = Describe("UserLogin", Label("privatelink"), func() {
	var data model.TestDataProvider

	_ = BeforeEach(func() {
		Eventually(kubecli.GetVersionOutput()).Should(Say(K8sVersion))
		checkUpAWSEnviroment()
		checkUpAzureEnviroment()
		checkNSetUpGCPEnviroment()
	})

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + data.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			By("Save logs to output directory ", func() {
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
			})
		}
		By("Clean Cloud", func() {
			DeleteAllPrivateEndpoints(&data)
		})
		By("Delete Resources, Project with PEService", func() {
			actions.DeleteUserResourcesProject(&data)
		})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test model.TestDataProvider, pe []privateEndpoint) {
			data = test
			privateFlow(&data, pe)
		},
		Entry("Test[privatelink-aws-1]: User has project which was updated with AWS PrivateEndpoint", Label("privatelink-aws-1"),
			model.NewTestDataProvider(
				"privatelink-aws-1",
				model.AProject{},
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
		Entry("Test[privatelink-azure-1]: User has project which was updated with Azure PrivateEndpoint", Label("privatelink-azure-1"),
			model.NewTestDataProvider(
				"privatelink-azure-1",
				model.AProject{},
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
		Entry("Test[privatelink-aws-2]: User has project which was updated with 2 AWS PrivateEndpoint", Label("privatelink-aws-2"),
			model.NewTestDataProvider(
				"privatelink-aws-2",
				model.AProject{},
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
				{
					provider: "AWS",
					region:   "us-east-1",
				},
			},
		),
		Entry("Test[privatelink-aws-azure-2]: User has project which was updated with 2 AWS PrivateEndpoint", Label("privatelink-aws-azure-2"),
			model.NewTestDataProvider(
				"privatelink-aws-azure",
				model.AProject{},
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
				{
					provider: "AWS",
					region:   "us-east-1",
				},
				{
					provider: "AZURE",
					region:   "northeurope",
				},
			},
		),
		Entry("Test[privatelink-gpc-1]: User has project which was updated with 2 AWS PrivateEndpoint", Label("privatelink-gpc-1"),
			model.NewTestDataProvider(
				"privatelink-gpc-1",
				model.AProject{},
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
					provider: "GCP",
					region:   "europe-west1",
				},
			},
		),
	)
})

func privateFlow(userData *model.TestDataProvider, requstedPE []privateEndpoint) {
	By("Deploy Project with requested configuration", func() {
		actions.PrepareUsersConfigurations(userData)
		deploy.NamespacedOperator(userData)
		actions.DeployProjectAndWait(userData, "1")
	})

	By("Create Private Link and the rest users resources", func() {
		for _, pe := range requstedPE {
			userData.Resources.Project.WithPrivateLink(provider.ProviderName(pe.provider), pe.region)
		}
		actions.PrepareUsersConfigurations(userData)
		actions.DeployProject(userData, "2")
	})

	By("Check if project statuses are updating, get project ID", func() {
		Eventually(kube.GetProjectPEndpointServiceStatus(userData), "15m", "10s").Should(Equal("True"),
			"Atlasproject status.conditions are not True")
		Eventually(kube.GetReadyProjectStatus(userData)).Should(Equal("True"),
			"Atlasproject status.conditions are not True")
		Expect(AllPEndpointUpdated(userData)).Should(BeTrue(),
			"Error: Was created a different amount of endpoints")
		actions.UpdateProjectID(userData)
		Expect(userData.Resources.ProjectID).ShouldNot(BeEmpty())
	})

	By("Create Endpoint in requested Cloud Provider", func() {
		project, err := kube.GetProjectResource(userData)
		Expect(err).ShouldNot(HaveOccurred())

		for _, peitem := range project.Status.PrivateEndpoints {
			cloudTest, err := cloud.CreatePEActions(peitem)
			Expect(err).ShouldNot(HaveOccurred())

			privateEndpointID := peitem.ID
			Expect(privateEndpointID).ShouldNot(BeEmpty())

			output, err := cloudTest.CreatePrivateEndpoint(privateEndpointID)
			Expect(err).ShouldNot(HaveOccurred())
			userData.Resources.Project = userData.Resources.Project.UpdatePrivateLinkID(output)
		}
	})

	By("Deploy Changed Projects", func() {
		actions.PrepareUsersConfigurations(userData)
		actions.DeployProjectAndWait(userData, "3")
	})

	By("Check statuses", func() {
		Eventually(kube.GetProjectPEndpointStatus(userData)).Should(Equal("True"), "Condition status 'PrivateEndpointReady' is not'True'")
		Eventually(kube.GetReadyProjectStatus(userData)).Should(Equal("True"), "Condition status 'Ready' is not 'True'")

		project, err := kube.GetProjectResource(userData)
		Expect(err).ShouldNot(HaveOccurred())
		for _, peitem := range project.Status.PrivateEndpoints {
			cloudTest, err := cloud.CreatePEActions(peitem)
			Expect(err).ShouldNot(HaveOccurred())
			privateEndpointID := userData.Resources.Project.GetPrivateIDByProviderRegion(peitem)
			Expect(privateEndpointID).ShouldNot(BeEmpty())
			Eventually(
				func() bool {
					return cloudTest.IsStatusPrivateEndpointAvailable(privateEndpointID)
				},
			).Should(BeTrue())
		}
	})
}

// DeleteAllPrivateEndpoints Specific for the current suite  - delete all requested Private Endpoints by test data
func DeleteAllPrivateEndpoints(data *model.TestDataProvider) {
	errorList := make([]string, 0)
	project, err := kube.GetProjectResource(data)
	Expect(err).ShouldNot(HaveOccurred())
	for _, peitem := range project.Status.PrivateEndpoints {
		cloudTest, err := cloud.CreatePEActions(peitem)
		if err == nil {
			privateEndpointID := data.Resources.Project.GetPrivateIDByProviderRegion(peitem)
			if privateEndpointID != "" {
				err = cloudTest.DeletePrivateEndpoint(privateEndpointID)
				if err != nil {
					GinkgoWriter.Write([]byte(err.Error()))
					errorList = append(errorList, err.Error())
				}
			}
		} else {
			errorList = append(errorList, err.Error())
		}
	}
	Expect(len(errorList)).Should(Equal(0), errorList)
}

func AllPEndpointUpdated(data *model.TestDataProvider) bool {
	result, _ := kube.GetProjectResource(data)
	return len(result.Status.PrivateEndpoints) == len(result.Spec.PrivateEndpoints)
}
