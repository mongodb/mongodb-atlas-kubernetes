package e2e_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"k8s.io/apimachinery/pkg/types"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/cloud"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

// NOTES
// Feature unavailable in Free and Shared-Tier Deployments
// This feature is not available for M0 free deployments, M2, and M5 deployments.

// tag for test resources "atlas-operator-test" (config.Tag)

// AWS NOTES: reserved VPC in eu-west-2, eu-south-1, us-east-1 (due to limitation no more 4 VPC per region)

type privateEndpoint struct {
	provider provider.ProviderName
	region   string
}

var _ = Describe("UserLogin", Label("privatelink"), func() {
	var testData *model.TestDataProvider

	_ = BeforeEach(func() {
		Eventually(kubecli.GetVersionOutput()).Should(Say(K8sVersion))
		checkUpAWSEnvironment()
		checkUpAzureEnvironment()
		checkNSetUpGCPEnvironment()
	})

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Clean Cloud", func() {
			DeleteAllPrivateEndpoints(testData)
		})
		By("Delete Resources, Project with PEService", func() {
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test *model.TestDataProvider, pe []privateEndpoint) {
			testData = test
			actions.ProjectCreationFlow(test)
			privateFlow(test, pe)
		},
		Entry("Test[privatelink-aws-1]: User has project which was updated with AWS PrivateEndpoint", Label("privatelink-aws-1"),
			model.DataProvider(
				"privatelink-aws-1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]privateEndpoint{
				{
					provider: "AWS",
					region:   config.AWSRegionEU,
				},
			},
		),
		Entry("Test[privatelink-azure-1]: User has project which was updated with Azure PrivateEndpoint", Label("privatelink-azure-1"),
			model.DataProvider(
				"privatelink-azure-1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]privateEndpoint{{
				provider: "AZURE",
				region:   config.AzureRegionEU,
			}},
		),
		Entry("Test[privatelink-aws-2]: User has project which was updated with 2 AWS PrivateEndpoint", Label("privatelink-aws-2"),
			model.DataProvider(
				"privatelink-aws-2",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]privateEndpoint{
				{
					provider: "AWS",
					region:   config.AWSRegionEU,
				},
				{
					provider: "AWS",
					region:   config.AWSRegionUS,
				},
			},
		),
		Entry("Test[privatelink-aws-azure-2]: User has project which was updated with 2 AWS and 1 Azure PrivateEndpoint", Label("privatelink-aws-azure-2"),
			model.DataProvider(
				"privatelink-aws-azure",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]privateEndpoint{
				{
					provider: "AWS",
					region:   config.AWSRegionEU,
				},
				{
					provider: "AWS",
					region:   config.AWSRegionUS,
				},
				{
					provider: "AZURE",
					region:   config.AzureRegionEU,
				},
			},
		),
		Entry("Test[privatelink-gpc-1]: User has project which was updated with 1 GCP PrivateEndpoint", Label("privatelink-gpc-1"),
			model.DataProvider(
				"privatelink-gpc-1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			[]privateEndpoint{
				{
					provider: provider.ProviderGCP,
					region:   config.GCPRegion,
				},
			},
		),
	)
})

func privateFlow(userData *model.TestDataProvider, requstedPE []privateEndpoint) {
	By("Create Private Link and the rest users resources", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
			Namespace: userData.Resources.Namespace}, userData.Project)).To(Succeed())
		for _, pe := range requstedPE {
			userData.Project.Spec.PrivateEndpoints = append(userData.Project.Spec.PrivateEndpoints,
				v1.PrivateEndpoint{
					Provider: pe.provider,
					Region:   pe.region,
				})
		}
		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).To(Succeed())
	})

	By("Check if project statuses are updating, get project ID", func() {
		actions.WaitForConditionsToBecomeTrue(userData, status.PrivateEndpointServiceReadyType, status.ReadyType)
		Expect(AllPEndpointUpdated(userData)).Should(BeTrue(),
			"Error: Was created a different amount of endpoints")
		Expect(userData.Project.ID()).ShouldNot(BeEmpty())
	})

	By("Create Endpoint in requested Cloud Provider", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
			Namespace: userData.Resources.Namespace}, userData.Project)).To(Succeed())

		for _, peitem := range userData.Project.Status.PrivateEndpoints {
			cloudTest, err := cloud.CreatePEActions(peitem)
			Expect(err).ShouldNot(HaveOccurred())

			privateEndpointID := peitem.ID
			Expect(privateEndpointID).ShouldNot(BeEmpty())

			output, err := cloudTest.CreatePrivateEndpoint(privateEndpointID)
			Expect(err).ShouldNot(HaveOccurred())

			for i, peItem := range userData.Project.Spec.PrivateEndpoints {
				if (peItem.Provider == output.Provider) && (peItem.Region == output.Region) {
					userData.Project.Spec.PrivateEndpoints[i] = output
				}
			}
		}

		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).To(Succeed())
		actions.WaitForConditionsToBecomeTrue(userData, status.PrivateEndpointReadyType, status.ReadyType)
	})

	By("Check statuses", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
			Namespace: userData.Resources.Namespace}, userData.Project)).To(Succeed())
		for _, peStatus := range userData.Project.Status.PrivateEndpoints {
			Expect(peStatus.Region).ShouldNot(BeEmpty())
			cloudTest, err := cloud.CreatePEActions(peStatus)
			Expect(err).ShouldNot(HaveOccurred())
			privateEndpointID := GetPrivateEndpointID(peStatus)
			Expect(privateEndpointID).ShouldNot(BeEmpty())
			Eventually(
				func() bool {
					return cloudTest.IsStatusPrivateEndpointAvailable(privateEndpointID)
				},
			).Should(BeTrue())
		}
	})
}

func GetPrivateEndpointID(endpoint status.ProjectPrivateEndpoint) string {
	if endpoint.Provider == provider.ProviderAWS {
		return endpoint.InterfaceEndpointID
	}
	return endpoint.ID
}

// DeleteAllPrivateEndpoints Specific for the current suite  - delete all requested Private Endpoints by test data
func DeleteAllPrivateEndpoints(data *model.TestDataProvider) {
	errorList := make([]string, 0)
	Expect(data.K8SClient.Get(data.Context, types.NamespacedName{Name: data.Project.Name,
		Namespace: data.Resources.Namespace}, data.Project)).To(Succeed())
	for _, peStatus := range data.Project.Status.PrivateEndpoints {
		cloudTest, err := cloud.CreatePEActions(peStatus)
		if err == nil {
			privateEndpointID := data.Resources.Project.GetPrivateIDByProviderRegion(peStatus)
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
	err := data.K8SClient.Get(data.Context, types.NamespacedName{Name: data.Project.Name, Namespace: data.Resources.Namespace}, data.Project)
	if err != nil {
		return false
	}
	return len(data.Project.Spec.PrivateEndpoints) == len(data.Project.Status.PrivateEndpoints)
}
