package e2e_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/cloud"
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
	var providerAction cloud.Provider

	_ = BeforeEach(OncePerOrdered, func() {
		checkUpAWSEnvironment()
		checkUpAzureEnvironment()
		checkNSetUpGCPEnvironment()

		action, err := prepareProviderAction()
		Expect(err).To(BeNil())
		providerAction = action
	})

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Delete Resources, Project with PEService", func() {
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test *model.TestDataProvider, pe []privateEndpoint) {
			testData = test
			actions.ProjectCreationFlow(test)
			privateFlow(test, providerAction, pe)
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
		Entry("Test[privatelink-two-identical-aws]: User has project which was updated with 2 Identical AWS Private Endpoints", Label("privatelink-aws-2"),
			model.DataProvider(
				"privatelink-two-identical-aws",
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
					region:   config.AWSRegionEU,
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

func privateFlow(userData *model.TestDataProvider, providerAction cloud.Provider, requestedPE []privateEndpoint) {
	By("Create Private Link and the rest users resources", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
			Namespace: userData.Resources.Namespace}, userData.Project)).To(Succeed())
		for _, pe := range requestedPE {
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

	//nolint:dupl
	By("Create Endpoint in requested Cloud Provider", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
			Namespace: userData.Resources.Namespace}, userData.Project)).To(Succeed())

		for idx, peStatusItem := range userData.Project.Status.PrivateEndpoints {
			privateEndpointID := peStatusItem.ID
			Expect(privateEndpointID).ShouldNot(BeEmpty())

			peName := getPrivateLinkName(privateEndpointID, peStatusItem.Provider, idx)
			var pe *cloud.PrivateEndpointDetails

			switch peStatusItem.Provider {
			case provider.ProviderAWS:
				providerAction.SetupNetwork(
					peStatusItem.Provider,
					cloud.WithAWSConfig(&cloud.AWSConfig{Region: peStatusItem.Region}),
				)
				pe = providerAction.SetupPrivateEndpoint(
					&cloud.AWSPrivateEndpointRequest{
						ID:          peName,
						Region:      peStatusItem.Region,
						ServiceName: peStatusItem.ServiceName,
					},
				)
			case provider.ProviderGCP:
				providerAction.SetupNetwork(
					peStatusItem.Provider,
					cloud.WithGCPConfig(&cloud.GCPConfig{Region: peStatusItem.Region}),
				)
				pe = providerAction.SetupPrivateEndpoint(
					&cloud.GCPPrivateEndpointRequest{
						ID:         peName,
						Region:     peStatusItem.Region,
						Targets:    peStatusItem.ServiceAttachmentNames,
						SubnetName: cloud.Subnet1Name,
					},
				)
			case provider.ProviderAzure:
				providerAction.SetupNetwork(
					peStatusItem.Provider,
					cloud.WithAzureConfig(&cloud.AzureConfig{Region: peStatusItem.Region}),
				)
				pe = providerAction.SetupPrivateEndpoint(
					&cloud.AzurePrivateEndpointRequest{
						ID:                peName,
						Region:            peStatusItem.Region,
						ServiceResourceID: peStatusItem.ServiceResourceID,
						SubnetName:        cloud.Subnet1Name,
					},
				)
			}

			for i, peItem := range userData.Project.Spec.PrivateEndpoints {
				if userData.Project.Spec.PrivateEndpoints[i].ID != "" {
					continue
				}

				if (peItem.Provider == pe.ProviderName) && (peItem.Region == pe.Region) {
					peItem.ID = pe.ID
					peItem.IP = pe.IP
					peItem.GCPProjectID = pe.GCPProjectID
					peItem.EndpointGroupName = pe.EndpointGroupName

					if len(pe.Endpoints) > 0 {
						peItem.Endpoints = make([]v1.GCPEndpoint, 0, len(pe.Endpoints))

						for _, ep := range pe.Endpoints {
							peItem.Endpoints = append(
								peItem.Endpoints,
								v1.GCPEndpoint{
									EndpointName: ep.Name,
									IPAddress:    ep.IP,
								},
							)
						}
					}

					userData.Project.Spec.PrivateEndpoints[i] = peItem
					break
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
			privateEndpointID := GetPrivateEndpointID(peStatus)
			Expect(privateEndpointID).ShouldNot(BeEmpty())

			providerAction.ValidatePrivateEndpointStatus(peStatus.Provider, privateEndpointID, peStatus.Region, len(peStatus.ServiceAttachmentNames))
		}
	})
}

func GetPrivateEndpointID(endpoint status.ProjectPrivateEndpoint) string {
	if endpoint.Provider == provider.ProviderAWS {
		return endpoint.InterfaceEndpointID
	}
	return endpoint.ID
}

func AllPEndpointUpdated(data *model.TestDataProvider) bool {
	err := data.K8SClient.Get(data.Context, types.NamespacedName{Name: data.Project.Name, Namespace: data.Resources.Namespace}, data.Project)
	if err != nil {
		return false
	}
	return len(data.Project.Spec.PrivateEndpoints) == len(data.Project.Status.PrivateEndpoints)
}

func getPrivateLinkName(privateEndpointID string, providerName provider.ProviderName, idx int) string {
	if providerName == provider.ProviderAWS {
		return fmt.Sprintf("%s_%d", privateEndpointID, idx)
	}
	return privateEndpointID
}

func prepareProviderAction() (*cloud.ProviderAction, error) {
	t := GinkgoT()

	aws, err := cloud.NewAWSAction(t)
	if err != nil {
		return nil, err
	}

	gcp, err := cloud.NewGCPAction(t, cloud.GoogleProjectID)
	if err != nil {
		return nil, err
	}

	azure, err := cloud.NewAzureAction(t, os.Getenv("AZURE_SUBSCRIPTION_ID"), cloud.ResourceGroupName)
	if err != nil {
		return nil, err
	}

	return cloud.NewProviderAction(t, aws, gcp, azure), nil
}
