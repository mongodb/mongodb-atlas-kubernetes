package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/cloud"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/resources"
)

// NOTES
// Feature unavailable in Free and Shared-Tier Deployments
// This feature is not available for M0 free deployments, M2, and M5 deployments.

// tag for test resources "atlas-operator-test" (config.Tag)

// AWS NOTES: reserved VPC in eu-west-2, eu-south-1, us-east-1 (due to limitation no more 4 VPC per region)

var _ = Describe("UserLogin", Label("datafederation"), func() {
	var testData *model.TestDataProvider
	var providerAction cloud.Provider
	var pe *cloud.PrivateEndpointDetails

	_ = BeforeEach(OncePerOrdered, func() {
		checkUpAWSEnvironment()
		checkUpAzureEnvironment()
		checkNSetUpGCPEnvironment()
		action, err := prepareProviderAction()
		Expect(err).To(BeNil())
		providerAction = action

		By("Setting up project", func() {
			testData = model.DataProvider(
				"privatelink-aws-1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject())

			actions.ProjectCreationFlow(testData)
		})
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

	It("Creates a data federation with private endpoint", func(ctx context.Context) {
		const dataFederationInstanceName = "test-data-federation-aws"

		//nolint:dupl
		By("Create private endpoint in AWS", func() {
			Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{Name: testData.Project.Name,
				Namespace: testData.Resources.Namespace}, testData.Project)).To(Succeed())

			vpcId := providerAction.SetupNetwork(
				"AWS",
				cloud.WithAWSConfig(&cloud.AWSConfig{
					VPC:           utils.RandomName("datafederation-private-endpoint"),
					Region:        config.AWSRegionEU,
					EnableCleanup: true,
				}),
			)
			pe = providerAction.SetupPrivateEndpoint(
				&cloud.AWSPrivateEndpointRequest{
					ID:     "vpce-" + vpcId,
					Region: config.AWSRegionEU,
					// See https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Data-Federation/operation/createDataFederationPrivateEndpoint
					ServiceName: "com.amazonaws.vpce.eu-west-2.vpce-svc-052f1840aa0c4f1f9",
				},
			)
		})

		By("Creating DataFederation with a PrivateEndpoint", func() {
			createdDataFederation := akov2.NewDataFederationInstance(
				testData.Project.Name,
				dataFederationInstanceName,
				testData.Project.Namespace).WithPrivateEndpoint(pe.ID, "AWS", "DATA_LAKE")
			Expect(testData.K8SClient.Create(context.Background(), createdDataFederation)).ShouldNot(HaveOccurred())

			Eventually(func(g Gomega) {
				df, _, err := atlasClient.Client.DataFederationApi.
					GetFederatedDatabase(context.Background(), testData.Project.ID(), createdDataFederation.Spec.Name).
					Execute()
				g.Expect(err).ShouldNot(HaveOccurred())
				g.Expect(df).NotTo(BeNil())
			}).WithTimeout(20 * time.Minute).WithPolling(15 * time.Second).ShouldNot(HaveOccurred())
		})

		By("Checking the DataFederation is ready", func() {
			df := &akov2.AtlasDataFederation{}
			Expect(testData.K8SClient.Get(context.Background(), types.NamespacedName{
				Namespace: testData.Project.Namespace,
				Name:      dataFederationInstanceName,
			}, df)).To(Succeed())
			Eventually(func() bool {
				return resources.CheckCondition(testData.K8SClient, df, api.TrueCondition(api.ReadyType))
			}).WithTimeout(2 * time.Minute).WithPolling(20 * time.Second).Should(BeTrue())
		})

		By("Delete DataFederation", func() {
			df := &akov2.AtlasDataFederation{}
			Expect(testData.K8SClient.Get(context.Background(), types.NamespacedName{
				Namespace: testData.Project.Namespace,
				Name:      dataFederationInstanceName,
			}, df)).To(Succeed())
			Expect(testData.K8SClient.Delete(testData.Context, df)).Should(Succeed())
		})

		By("Deleting DataFederation Private Endpoint", func() {
			// This is required or will result on error:
			// CANNOT_CLOSE_GROUP_ACTIVE_ATLAS_DATA_FEDERATION_PRIVATE_ENDPOINTS
			// for some reason, requesting deletion successfully just once doesn't work
			// TODO: revisit and cleanup once CLOUDP-280905 is fixed
			Eventually(func(g Gomega) {
				_, resp, err := atlasClient.Client.DataFederationApi.
					DeleteDataFederationPrivateEndpoint(testData.Context, testData.Project.ID(), pe.ID).
					Execute()
				g.Expect(err).To(BeNil(), fmt.Sprintf("deletion of private endpoint failed with error %v", err))
				g.Expect(resp).NotTo(BeNil())
				g.Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusNoContent))
			}).WithTimeout(5 * time.Minute).WithPolling(15 * time.Second).MustPassRepeatedly(2).Should(Succeed())
		})
	})
})
