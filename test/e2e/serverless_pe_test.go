package e2e_test

import (
	"fmt"
	"time"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"k8s.io/apimachinery/pkg/types"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasdeployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/serverlessprivateendpoint"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

var _ = Describe("UserLogin", Label("serverless-pe"), func() {
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
			Expect(actions.SaveDeploymentsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Clean Cloud", func() {
			DeleteSPE(testData)
		})
		By("Delete Resources", func() {
			actions.DeleteTestDataDeployments(testData)
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test *model.TestDataProvider, spe []v1.ServerlessPrivateEndpoint) {
			testData = test
			actions.ProjectCreationFlow(test)
			speFlow(test, spe)
		},
		Entry("Test[spe-aws-1]: Serverless deployment with one AWS PE", Label("spe-aws-1"),
			model.DataProvider(
				"spe-aws-1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()).WithInitialDeployments(data.CreateServerlessDeployment("spetest1", "AWS", "US_EAST_1")),
			[]v1.ServerlessPrivateEndpoint{
				{
					Name: "pe1",
				},
			},
		),
		Entry("Test[spe-azure-1]: Serverless deployment with one Azure PE", Label("spe-azure-1"),
			model.DataProvider(
				"spe-azure-1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()).WithInitialDeployments(data.CreateServerlessDeployment("spetest3", "AZURE", "US_EAST_2")),
			[]v1.ServerlessPrivateEndpoint{
				{
					Name: "pe1",
				},
			},
		),
		Entry("Test[spe-azure-2]: Serverless deployment with one valid and one non-valid Azure PEs", Label("spe-azure-2"),
			model.DataProvider(
				"spe-azure-2",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()).WithInitialDeployments(data.CreateServerlessDeployment("spetest3", "AZURE", "US_EAST_2")),
			[]v1.ServerlessPrivateEndpoint{
				{
					Name: "pe1",
				},
				{
					Name:                     "pe2",
					CloudProviderEndpointID:  "invalid",
					PrivateEndpointIPAddress: "invalid",
				},
			},
		),
	)
})

func speFlow(userData *model.TestDataProvider, spe []v1.ServerlessPrivateEndpoint) {
	By("Apply deployment", func() {
		Expect(userData.InitialDeployments).ShouldNot(BeEmpty())
		userData.InitialDeployments[0].Namespace = userData.Resources.Namespace
		Expect(userData.K8SClient.Create(userData.Context, userData.InitialDeployments[0])).To(Succeed())
		Eventually(func(g Gomega) bool {
			g.Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{
				Name:      userData.InitialDeployments[0].Name,
				Namespace: userData.InitialDeployments[0].Namespace,
			}, userData.InitialDeployments[0])).To(Succeed())
			return userData.InitialDeployments[0].Status.StateName == status.StateIDLE
		}).WithTimeout(10 * time.Minute).Should(BeTrue())
	})

	By("Adding Private Endpoints to Deployment", func() {
		updateSPE(userData, spe)
		invalidSPEFlow(userData, spe)
		waitSPEStatus(userData, atlasdeployment.SPEStatusReserved, len(spe))
	})

	By("Create Private Endpoints in Cloud", func() {
		Expect(serverlessprivateendpoint.ConnectSPE(spe, userData.InitialDeployments[0].Status.ServerlessPrivateEndpoints,
			provider.ProviderName(userData.InitialDeployments[0].Spec.ServerlessSpec.ProviderSettings.BackingProviderName))).To(Succeed())
	})

	By("Update Private Endpoints in Deployment", func() {
		updateSPE(userData, spe)
		waitSPEStatus(userData, atlasdeployment.SPEStatusAvailable, len(spe))
	})

	By("Delete Private Endpoints", func() {
		updateSPE(userData, []v1.ServerlessPrivateEndpoint{})
		Eventually(func(g Gomega) {
			g.Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{
				Name:      userData.InitialDeployments[0].Name,
				Namespace: userData.InitialDeployments[0].Namespace,
			}, userData.InitialDeployments[0])).To(Succeed())
			g.Expect(len(userData.InitialDeployments[0].Status.ServerlessPrivateEndpoints)).To(Equal(0))
			for _, condition := range userData.InitialDeployments[0].Status.Conditions {
				g.Expect(condition.Type).ToNot(Equal(status.ServerlessPrivateEndpointReadyType))
			}
		}).WithTimeout(5*time.Minute).Should(Succeed(), "Deployment should not have any Private Endpoints")
	})
}

func invalidSPEFlow(userData *model.TestDataProvider, spe []v1.ServerlessPrivateEndpoint) {
	// check that spe is valid
	isValid := true
	for _, pe := range spe {
		if pe.PrivateEndpointIPAddress != "" || pe.CloudProviderEndpointID != "" {
			isValid = false
			break
		}
	}
	if !isValid {
		// check that deployment is not ready
		Eventually(func(g Gomega) bool {
			g.Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{
				Name:      userData.InitialDeployments[0].Name,
				Namespace: userData.InitialDeployments[0].Namespace,
			}, userData.InitialDeployments[0])).To(Succeed())
			for _, speStatus := range userData.InitialDeployments[0].Status.ServerlessPrivateEndpoints {
				if speStatus.Status == atlasdeployment.SPEStatusFailed {
					return true
				}
			}
			return false
		}).WithTimeout(5*time.Minute).Should(BeTrue(), "Deployment should be failed if one of the SPEs is invalid")
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

func DeleteSPE(userData *model.TestDataProvider) {
	By("Delete Private Endpoints in Cloud", func() {
		Expect(serverlessprivateendpoint.DeleteSPE(userData.InitialDeployments[0].Status.ServerlessPrivateEndpoints,
			provider.ProviderName(userData.InitialDeployments[0].Spec.ServerlessSpec.ProviderSettings.BackingProviderName))).To(Succeed())
	})
}

func updateSPE(userData *model.TestDataProvider, spe []v1.ServerlessPrivateEndpoint) {
	Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.InitialDeployments[0].Name,
		Namespace: userData.Resources.Namespace}, userData.InitialDeployments[0])).To(Succeed())
	userData.InitialDeployments[0].Spec.ServerlessSpec.PrivateEndpoints = spe
	Expect(userData.K8SClient.Update(userData.Context, userData.InitialDeployments[0])).To(Succeed())
}

func waitSPEStatus(userData *model.TestDataProvider, status string, speLen int) {
	Eventually(func(g Gomega) bool {
		g.Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.InitialDeployments[0].Name,
			Namespace: userData.Resources.Namespace}, userData.InitialDeployments[0])).To(Succeed())
		if len(userData.InitialDeployments[0].Status.ServerlessPrivateEndpoints) != speLen {
			return false
		}
		for _, pe := range userData.InitialDeployments[0].Status.ServerlessPrivateEndpoints {
			if pe.Status != status {
				return false
			}
		}
		return true
	}).WithTimeout(10*time.Minute).Should(BeTrue(), fmt.Sprintf("Private Endpoints should be %s", status))
}
