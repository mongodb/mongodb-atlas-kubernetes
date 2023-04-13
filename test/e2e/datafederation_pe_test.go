package e2e_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/testutil"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
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

var _ = Describe("UserLogin", Label("datafederation"), func() {
	var testData *model.TestDataProvider

	_ = BeforeEach(func() {
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
			dataFederationFlow(test, pe)
		},
		Entry("Data Federation can be created with private endpoints", Label("datafederation-pe-aws"),
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
	)
})

func dataFederationFlow(userData *model.TestDataProvider, requstedPE []privateEndpoint) {
	var createdDataFederation *v1.AtlasDataFederation
	const dataFederationInstanceName = "test-data-federation-aws"

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

	//nolint:dupl
	By("Create Endpoint in requested Cloud Provider", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
			Namespace: userData.Resources.Namespace}, userData.Project)).To(Succeed())

		for idx, peitem := range userData.Project.Status.PrivateEndpoints {
			cloudTest, err := cloud.CreatePEActions(peitem)
			Expect(err).ShouldNot(HaveOccurred())

			privateEndpointID := peitem.ID
			Expect(privateEndpointID).ShouldNot(BeEmpty())

			peName := getPrivateLinkName(privateEndpointID, peitem.Provider, idx)

			output, err := cloudTest.CreatePrivateEndpoint(peName)
			Expect(err).ShouldNot(HaveOccurred())

			for i, peItem := range userData.Project.Spec.PrivateEndpoints {
				if userData.Project.Spec.PrivateEndpoints[i].ID != "" {
					continue
				}

				if (peItem.Provider == output.Provider) && (peItem.Region == output.Region) {
					userData.Project.Spec.PrivateEndpoints[i] = output
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

	By("Creating DataFederation with a PrivateEndpoint", func() {
		peData := userData.Project.Status.PrivateEndpoints[0]
		createdDataFederation = v1.NewDataFederationInstance(
			userData.Project.Name,
			dataFederationInstanceName,
			userData.Project.Namespace).WithPrivateEndpoint(GetPrivateEndpointID(peData), "AWS", "DATA_LAKE")
		Expect(userData.K8SClient.Create(context.Background(), createdDataFederation)).ShouldNot(HaveOccurred())

		Eventually(func(g Gomega) {
			df, _, err := atlasClient.Client.DataFederation.Get(context.Background(), userData.Project.ID(), createdDataFederation.Spec.Name)
			g.Expect(err).ShouldNot(HaveOccurred())
			g.Expect(df).NotTo(BeNil())
		}).WithTimeout(20 * time.Minute).WithPolling(15 * time.Second).ShouldNot(HaveOccurred())
	})

	By("Checking the DataFederation is ready", func() {
		df := &v1.AtlasDataFederation{}
		Expect(userData.K8SClient.Get(context.Background(), types.NamespacedName{
			Namespace: userData.Project.Namespace,
			Name:      dataFederationInstanceName,
		}, df)).To(Succeed())
		Eventually(func() bool {
			return testutil.CheckCondition(userData.K8SClient, df, status.TrueCondition(status.ReadyType))
		}).WithTimeout(2 * time.Minute).WithPolling(20 * time.Second).Should(BeTrue())
	})
}
