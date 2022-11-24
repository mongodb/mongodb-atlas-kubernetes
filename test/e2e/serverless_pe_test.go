package e2e_test

import (
	"time"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/serverlessprivateendpoint"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"k8s.io/apimachinery/pkg/types"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasdeployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
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
			SaveDump(testData)
		}
		By("Clean Cloud", func() {
			// TODO: clean cloud
		})
		By("Delete Resources, Project with PEService", func() {
			actions.DeleteTestDataDeployments(testData)
			actions.DeleteTestDataProject(testData)
			actions.DeleteGlobalKeyIfExist(*testData)
		})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test *model.TestDataProvider, spe []v1.ServerlessPrivateEndpoint) {
			testData = test
			actions.ProjectCreationFlow(test)
			speFlow(test, spe)
		},
		Entry("Test[spe-aws-1]: User has project which was updated with AWS PrivateEndpoint", Label("spe-aws-1"),
			model.DataProvider(
				"spe-aws-1",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()).WithInitialDeployments(data.CreateServerlessDeployment("spetest")),
			[]v1.ServerlessPrivateEndpoint{
				{
					Name: "pe1",
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
			return userData.InitialDeployments[0].Status.StateName == "IDLE" // TODO: use constant
		}).WithTimeout(10 * time.Minute).Should(BeTrue())
	})

	By("Adding Private Endpoints to Deployment", func() {

		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.InitialDeployments[0].Name,
			Namespace: userData.Resources.Namespace}, userData.InitialDeployments[0])).To(Succeed())
		userData.InitialDeployments[0].Spec.ServerlessSpec.PrivateEndpoints = spe
		Expect(userData.K8SClient.Update(userData.Context, userData.InitialDeployments[0])).To(Succeed())
		// Add wait for status
		Eventually(func(g Gomega) bool {
			g.Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.InitialDeployments[0].Name,
				Namespace: userData.Resources.Namespace}, userData.InitialDeployments[0])).To(Succeed())
			if len(userData.InitialDeployments[0].Status.ServerlessPrivateEndpoints) != len(spe) {
				return false
			}
			for _, pe := range userData.InitialDeployments[0].Status.ServerlessPrivateEndpoints {
				if pe.Status != atlasdeployment.SPEStatusReserved {
					return false
				}
			}
			return true
		}).WithTimeout(5*time.Minute).Should(BeTrue(), "Private Endpoints should be reserved")
	})

	By("Create Private Endpoints in Cloud", func() {
		Expect(serverlessprivateendpoint.ConnectSPE(spe, userData.InitialDeployments[0].Status.ServerlessPrivateEndpoints,
			provider.ProviderName(userData.InitialDeployments[0].Spec.ServerlessSpec.ProviderSettings.BackingProviderName))).To(Succeed())
	})

	By("Update Private Endpoints in Deployment", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.InitialDeployments[0].Name,
			Namespace: userData.Resources.Namespace}, userData.InitialDeployments[0])).To(Succeed())
		userData.InitialDeployments[0].Spec.ServerlessSpec.PrivateEndpoints = spe
		Expect(userData.K8SClient.Update(userData.Context, userData.InitialDeployments[0])).To(Succeed())
		// Add wait for status
		Eventually(func(g Gomega) bool {
			g.Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.InitialDeployments[0].Name,
				Namespace: userData.Resources.Namespace}, userData.InitialDeployments[0])).To(Succeed())
			if len(userData.InitialDeployments[0].Status.ServerlessPrivateEndpoints) != len(spe) {
				return false
			}
			for _, pe := range userData.InitialDeployments[0].Status.ServerlessPrivateEndpoints {
				if pe.Status != atlasdeployment.SPEStatusAvailable {
					return false
				}
			}
			return true
		}).WithTimeout(5*time.Minute).Should(BeTrue(), "Private Endpoints should be ready")
	})
}
