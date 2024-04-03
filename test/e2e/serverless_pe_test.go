package e2e_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasdeployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/conditions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/cloud"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/resources"
)

var _ = Describe("Serverless Private Endpoint", Label("serverless-pe"), func() {
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
			Expect(actions.SaveDeploymentsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Delete Resources", func() {
			actions.DeleteTestDataDeployments(testData)
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test *model.TestDataProvider, spe []akov2.ServerlessPrivateEndpoint) {
			testData = test
			actions.ProjectCreationFlow(test)
			speFlow(test, providerAction, spe)
		},
		Entry("Test[spe-aws-1]: Serverless deployment with one AWS PE", Label("spe-aws-1"),
			model.DataProvider(
				"spe-aws-1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()).WithInitialDeployments(data.CreateServerlessDeployment("spe-test-1", "AWS", "US_EAST_1")),
			[]akov2.ServerlessPrivateEndpoint{
				{
					Name: newRandomName("pe"),
				},
			},
		),
		Entry("Test[spe-azure-1]: Serverless deployment with one Azure PE", Label("spe-azure-1"),
			model.DataProvider(
				"spe-azure-1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()).WithInitialDeployments(data.CreateServerlessDeployment("spe-test-2", "AZURE", "US_EAST_2")),
			[]akov2.ServerlessPrivateEndpoint{
				{
					Name: newRandomName("pe"),
				},
			},
		),
		Entry("Test[spe-azure-2]: Serverless deployment with one valid and one non-valid Azure PEs", Label("spe-azure-2"),
			model.DataProvider(
				"spe-azure-2",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()).WithInitialDeployments(data.CreateServerlessDeployment("spe-test-3", "AZURE", "US_EAST_2")),
			[]akov2.ServerlessPrivateEndpoint{
				{
					Name: newRandomName("pe"),
				},
				{
					Name:                     newRandomName("pe"),
					CloudProviderEndpointID:  "invalid",
					PrivateEndpointIPAddress: "invalid",
				},
			},
		),
	)
})

func speFlow(userData *model.TestDataProvider, providerAction cloud.Provider, spe []akov2.ServerlessPrivateEndpoint) {
	By("Apply deployment", func() {
		Expect(userData.InitialDeployments).ShouldNot(BeEmpty())
		userData.InitialDeployments[0].Namespace = userData.Resources.Namespace
		Expect(userData.K8SClient.Create(userData.Context, userData.InitialDeployments[0])).To(Succeed())

		Eventually(func(g Gomega) {
			deployment := userData.InitialDeployments[0]
			g.Expect(userData.K8SClient.Get(userData.Context, client.ObjectKeyFromObject(deployment), deployment)).To(Succeed())
			g.Expect(deployment.Status.Conditions).To(ContainElement(conditions.MatchCondition(status.TrueCondition(status.DeploymentReadyType))))
		}).WithTimeout(time.Minute * 15).WithPolling(time.Second * 15).Should(Succeed())
	})

	By("Adding Private Endpoints to Deployment", func() {
		updateSPE(userData, spe)
		invalidSPEFlow(userData, spe)
		waitSPEStatus(userData, atlasdeployment.SPEStatusReserved, len(spe))
	})

	By("Create Private Endpoints in Cloud", func() {
		var pe *cloud.PrivateEndpointDetails

		for _, peItem := range userData.InitialDeployments[0].Status.ServerlessPrivateEndpoints {
			switch provider.ProviderName(peItem.ProviderName) {
			case provider.ProviderAWS:
				providerAction.SetupNetwork(
					provider.ProviderAWS,
					cloud.WithAWSConfig(&cloud.AWSConfig{Region: "us-east-1"}),
				)
				pe = providerAction.SetupPrivateEndpoint(
					&cloud.AWSPrivateEndpointRequest{
						ID:          peItem.Name,
						Region:      "us-east-1",
						ServiceName: peItem.EndpointServiceName,
					},
				)
			case provider.ProviderAzure:
				providerAction.SetupNetwork(provider.ProviderAzure, nil)
				pe = providerAction.SetupPrivateEndpoint(
					&cloud.AzurePrivateEndpointRequest{
						ID:                peItem.Name,
						Region:            cloud.AzureRegion,
						ServiceResourceID: peItem.PrivateLinkServiceResourceID,
						SubnetName:        cloud.Subnet2Name,
					},
				)
			}

			for i, specPE := range spe {
				if specPE.Name == peItem.Name {
					spe[i].CloudProviderEndpointID = pe.ID
					spe[i].PrivateEndpointIPAddress = pe.IP
				}
			}
		}
	})

	By("Update Private Endpoints in Deployment", func() {
		updateSPE(userData, spe)
		waitSPEStatus(userData, atlasdeployment.SPEStatusAvailable, len(spe))
	})

	By("Delete Private Endpoints", func() {
		updateSPE(userData, []akov2.ServerlessPrivateEndpoint{})
		Eventually(func(g Gomega) {
			Expect(
				userData.K8SClient.Get(
					userData.Context,
					client.ObjectKeyFromObject(userData.InitialDeployments[0]),
					userData.InitialDeployments[0],
				),
			).To(Succeed())
			g.Expect(len(userData.InitialDeployments[0].Status.ServerlessPrivateEndpoints)).To(Equal(0))
			for _, condition := range userData.InitialDeployments[0].Status.Conditions {
				g.Expect(condition.Type).ToNot(Equal(status.ServerlessPrivateEndpointReadyType))
			}
		}).WithTimeout(15*time.Minute).Should(Succeed(), "Deployment should not have any Private Endpoints")
	})
}

func invalidSPEFlow(userData *model.TestDataProvider, spe []akov2.ServerlessPrivateEndpoint) {
	// check that spe is valid
	isValid := true
	for _, pe := range spe {
		if pe.PrivateEndpointIPAddress != "" || pe.CloudProviderEndpointID != "" {
			isValid = false
			break
		}
	}
	if !isValid {
		expectedCondition := status.FalseCondition(status.ServerlessPrivateEndpointReadyType).WithReason(string(workflow.ServerlessPrivateEndpointFailed))
		Eventually(func() bool {
			return resources.CheckCondition(userData.K8SClient, userData.InitialDeployments[0], expectedCondition)
		}).WithTimeout(15 * time.Minute).WithPolling(10 * time.Second).Should(BeTrue())

		// fix spe
		for i, pe := range spe {
			if pe.PrivateEndpointIPAddress != "" || pe.CloudProviderEndpointID != "" {
				spe[i].PrivateEndpointIPAddress = ""
				spe[i].CloudProviderEndpointID = ""
			}
		}
		updateSPE(userData, spe)
	}
}

func updateSPE(userData *model.TestDataProvider, spe []akov2.ServerlessPrivateEndpoint) {
	Expect(
		userData.K8SClient.Get(
			userData.Context,
			client.ObjectKeyFromObject(userData.InitialDeployments[0]),
			userData.InitialDeployments[0],
		),
	).To(Succeed())
	userData.InitialDeployments[0].Spec.ServerlessSpec.PrivateEndpoints = spe
	Expect(userData.K8SClient.Update(userData.Context, userData.InitialDeployments[0])).To(Succeed())
}

func waitSPEStatus(userData *model.TestDataProvider, status string, speLen int) {
	Eventually(func(g Gomega) bool {
		g.Expect(
			userData.K8SClient.Get(
				userData.Context,
				client.ObjectKeyFromObject(userData.InitialDeployments[0]),
				userData.InitialDeployments[0],
			),
		).To(Succeed())

		if len(userData.InitialDeployments[0].Status.ServerlessPrivateEndpoints) != speLen {
			return false
		}

		for _, pe := range userData.InitialDeployments[0].Status.ServerlessPrivateEndpoints {
			if pe.Status != status {
				return false
			}
		}

		return true
	}).WithTimeout(15*time.Minute).Should(BeTrue(), fmt.Sprintf("Private Endpoints should be %s", status))
}
