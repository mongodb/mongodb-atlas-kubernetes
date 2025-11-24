// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
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

var _ = Describe("DataFederation Private Endpoint", Label("datafederation"), FlakeAttempts(3), func() {
	var testData *model.TestDataProvider
	var providerAction cloud.Provider
	var pe *cloud.PrivateEndpointDetails
	var secondPE *cloud.PrivateEndpointDetails

	_ = BeforeEach(OncePerOrdered, func(ctx SpecContext) {
		checkUpAWSEnvironment()
		checkUpAzureEnvironment()
		checkNSetUpGCPEnvironment()
		action, err := prepareProviderAction(ctx)
		Expect(err).To(BeNil())
		providerAction = action

		By("Setting up project", func() {
			testData = model.DataProvider(ctx, "privatelink-aws-1", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())

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

			vpcId := providerAction.SetupNetwork(ctx, "AWS", cloud.WithAWSConfig(&cloud.AWSConfig{
				VPC:           utils.RandomName("datafederation-private-endpoint"),
				Region:        config.AWSRegionEU,
				EnableCleanup: true,
			}))
			pe = providerAction.SetupPrivateEndpoint(ctx, &cloud.AWSPrivateEndpointRequest{
				ID:     "vpce-" + vpcId,
				Region: config.AWSRegionEU,
				// See https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Data-Federation/operation/createDataFederationPrivateEndpoint
				ServiceName: "com.amazonaws.vpce.eu-west-2.vpce-svc-052f1840aa0c4f1f9",
			})
		})

		By("Creating DataFederation with a PrivateEndpoint", func() {
			createdDataFederation := akov2.NewDataFederationInstance(
				testData.Project.Name,
				dataFederationInstanceName,
				testData.Project.Namespace).WithPrivateEndpoint(pe.ID, "AWS", "DATA_LAKE")
			createdDataFederation.Spec.Storage = &akov2.Storage{
				Databases: []akov2.Database{
					{
						Name: "test-db-1",
						Collections: []akov2.Collection{
							{
								Name: "test-collection-1",
								DataSources: []akov2.DataSource{
									{
										StoreName: "http-test",
										Urls: []string{
											"https://data.cityofnewyork.us/api/views/vfnx-vebw/rows.csv",
										},
									},
								},
							},
						},
					},
				},
				Stores: []akov2.Store{
					{
						Name:     "http-test",
						Provider: "http",
					},
				},
			}
			Expect(testData.K8SClient.Create(context.Background(), createdDataFederation)).ShouldNot(HaveOccurred())

			Eventually(func(g Gomega) {
				df, _, err := atlasClient.Client.DataFederationApi.
					GetDataFederation(context.Background(), testData.Project.ID(), createdDataFederation.Spec.Name).
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

		//nolint:dupl
		By("Create a new private endpoint in AWS", func() {
			Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{Name: testData.Project.Name,
				Namespace: testData.Resources.Namespace}, testData.Project)).To(Succeed())

			vpcId := providerAction.SetupNetwork(ctx, "AWS", cloud.WithAWSConfig(&cloud.AWSConfig{
				VPC:           utils.RandomName("datafederation-private-endpoint2"),
				Region:        config.AWSRegionEU,
				EnableCleanup: true,
			}))
			secondPE = providerAction.SetupPrivateEndpoint(ctx, &cloud.AWSPrivateEndpointRequest{
				ID:     "vpce-" + vpcId,
				Region: config.AWSRegionEU,
				// See https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Data-Federation/operation/createDataFederationPrivateEndpoint
				ServiceName: "com.amazonaws.vpce.eu-west-2.vpce-svc-052f1840aa0c4f1f9",
			})
		})

		By("Update DataFederation with the new Private Endpoint", func() {
			df := &akov2.AtlasDataFederation{}
			Expect(testData.K8SClient.Get(context.Background(), types.NamespacedName{
				Namespace: testData.Project.Namespace,
				Name:      dataFederationInstanceName,
			}, df)).To(Succeed())
			df.Spec.PrivateEndpoints[0].EndpointID = secondPE.ID
			Expect(testData.K8SClient.Update(context.Background(), df)).ShouldNot(HaveOccurred())
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
				resp, err := atlasClient.Client.DataFederationApi.
					DeleteDataFederation(testData.Context, testData.Project.ID(), secondPE.ID).
					Execute()
				g.Expect(err).To(BeNil(), fmt.Sprintf("deletion of private endpoint failed with error %v", err))
				g.Expect(resp).NotTo(BeNil())
				g.Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusNoContent))
			}).WithTimeout(5 * time.Minute).WithPolling(15 * time.Second).MustPassRepeatedly(2).Should(Succeed())
		})
	})
})
