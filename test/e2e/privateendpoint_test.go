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
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/conditions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/cloud"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/resources"
)

var _ = Describe("Private Endpoints", Label("private-endpoint"), FlakeAttempts(3), func() {
	var testData *model.TestDataProvider

	_ = BeforeEach(OncePerOrdered, func() {
		checkUpAWSEnvironment()
		checkUpAzureEnvironment()
		checkNSetUpGCPEnvironment()
	})

	_ = AfterEach(func() {
		GinkgoWriter.Println()
		GinkgoWriter.Println("===============================================")
		GinkgoWriter.Println("Operator namespace: " + testData.Resources.Namespace)
		GinkgoWriter.Println("===============================================")
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Delete Project and cluster resources", func() {
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	DescribeTable(
		"Configure private endpoint for all supported cloud provider",
		func(ctx SpecContext, test func(context.Context) *model.TestDataProvider, pe *akov2.AtlasPrivateEndpoint) {
			var privateEndpointDetails *cloud.PrivateEndpointDetails

			testData = test(ctx)
			actions.ProjectCreationFlow(testData)

			By("Preparing private endpoint resource", func() {
				pe.Namespace = testData.Resources.Namespace
				pe.Spec.ProjectRef = &common.ResourceRefNamespaced{
					Name:      testData.Project.Name,
					Namespace: testData.Project.Namespace,
				}
				region, err := cloud.GetAtlasRegionByProvider(pe.Spec.Provider)
				Expect(err).ToNot(HaveOccurred())
				pe.Spec.Region = region
			})

			By("Creating private endpoint", func() {
				Expect(testData.K8SClient.Create(testData.Context, pe)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(pe), pe)).To(Succeed())
					g.Expect(pe.Status.ServiceStatus).To(Equal("AVAILABLE"))
					g.Expect(resources.CheckCondition(testData.K8SClient, pe, api.TrueCondition(api.PrivateEndpointServiceReady))).Should(BeTrue())
				}).WithTimeout(15 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
			})

			By("Configuring external network", func() {
				Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(pe), pe)).To(Succeed())
				action, err := prepareProviderAction(ctx)
				Expect(err).ToNot(HaveOccurred())

				cloudRegion := cloud.MapCloudProviderRegion(pe.Spec.Provider, pe.Spec.Region)
				Expect(cloudRegion).ToNot(BeEmpty())

				switch pe.Spec.Provider {
				case "AWS":
					awsConfig, err := cloud.GenerateCloudConfig[cloud.AWSConfig](pe.Spec.Provider, cloudRegion, testData.Resources.KeyName)
					Expect(err).ToNot(HaveOccurred())

					Expect(action.SetupNetwork(ctx, provider.ProviderName(pe.Spec.Provider), cloud.WithAWSConfig(awsConfig))).ToNot(BeEmpty())
					privateEndpointDetails = action.SetupPrivateEndpoint(nil, &cloud.AWSPrivateEndpointRequest{
						ID:          fmt.Sprintf("aws-e2e-pe-%s", testData.Resources.TestID),
						Region:      awsConfig.Region,
						ServiceName: pe.Status.ServiceName,
					})
				case "AZURE":
					azureConfig, err := cloud.GenerateCloudConfig[cloud.AzureConfig](pe.Spec.Provider, cloudRegion, testData.Resources.KeyName)
					Expect(err).ToNot(HaveOccurred())

					Expect(action.SetupNetwork(ctx, provider.ProviderName(pe.Spec.Provider), cloud.WithAzureConfig(azureConfig))).ToNot(BeEmpty())
					privateEndpointDetails = action.SetupPrivateEndpoint(nil, &cloud.AzurePrivateEndpointRequest{
						ID:                fmt.Sprintf("azure-e2e-pe-%s", testData.Resources.TestID),
						Region:            azureConfig.Region,
						ServiceResourceID: pe.Status.ResourceID,
						SubnetName:        randomKeyFromMap(azureConfig.Subnets),
					})
				case "GCP":
					gcpConfig, err := cloud.GenerateCloudConfig[cloud.GCPConfig](pe.Spec.Provider, cloudRegion, testData.Resources.KeyName)
					Expect(err).ToNot(HaveOccurred())

					Expect(action.SetupNetwork(ctx, provider.ProviderName(pe.Spec.Provider), cloud.WithGCPConfig(gcpConfig))).ToNot(BeEmpty())
					privateEndpointDetails = action.SetupPrivateEndpoint(nil, &cloud.GCPPrivateEndpointRequest{
						ID:         fmt.Sprintf("gcp-e2e-pe-%s", testData.Resources.TestID),
						Region:     gcpConfig.Region,
						Targets:    pe.Status.ServiceAttachmentNames,
						SubnetName: randomKeyFromMap(gcpConfig.Subnets),
					})
				}
			})

			By("Configuring private endpoint with external network details", func() {
				Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(pe), pe)).To(Succeed())

				switch pe.Spec.Provider {
				case "AWS":
					pe.Spec.AWSConfiguration = []akov2.AWSPrivateEndpointConfiguration{
						{
							ID: privateEndpointDetails.ID,
						},
					}
				case "AZURE":
					pe.Spec.AzureConfiguration = []akov2.AzurePrivateEndpointConfiguration{
						{
							ID: privateEndpointDetails.ID,
							IP: privateEndpointDetails.IP,
						},
					}
				case "GCP":
					gcpEndpoints := make([]akov2.GCPPrivateEndpoint, 0, len(privateEndpointDetails.Endpoints))
					for _, ep := range privateEndpointDetails.Endpoints {
						gcpEndpoints = append(
							gcpEndpoints,
							akov2.GCPPrivateEndpoint{
								Name: ep.Name,
								IP:   ep.IP,
							},
						)
					}
					pe.Spec.GCPConfiguration = []akov2.GCPPrivateEndpointConfiguration{
						{
							ProjectID: privateEndpointDetails.GCPProjectID,
							GroupName: privateEndpointDetails.EndpointGroupName,
							Endpoints: gcpEndpoints,
						},
					}
				}

				Expect(testData.K8SClient.Update(ctx, pe)).To(Succeed())
				Eventually(func(g Gomega) { //nolint:dupl
					g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(pe), pe)).To(Succeed())
					g.Expect(pe.Status.ServiceStatus).To(Equal("AVAILABLE"))
					g.Expect(resources.CheckCondition(testData.K8SClient, pe, api.TrueCondition(api.PrivateEndpointServiceReady))).Should(BeTrue())
					for _, eStatus := range pe.Status.Endpoints {
						g.Expect(eStatus.Status).To(Equal("AVAILABLE"))
					}
					g.Expect(resources.CheckCondition(testData.K8SClient, pe, api.TrueCondition(api.PrivateEndpointReady))).Should(BeTrue())
					g.Expect(resources.CheckCondition(testData.K8SClient, pe, api.TrueCondition(api.ReadyType))).To(BeTrue())
				}).WithTimeout(15 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
			})

			By("Removing private endpoint", func() {
				Expect(testData.K8SClient.Delete(ctx, pe)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(pe), pe)).ShouldNot(Succeed())
				}).WithTimeout(15 * time.Minute).WithPolling(20 * time.Second).Should(Succeed())
			})
		},
		Entry(
			"Configure AWS private endpoint",
			Label("focus-aws-private-endpoint"),
			func(ctx context.Context) *model.TestDataProvider {
				return model.DataProvider(ctx, "aws-pe-1", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			&akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-pe-1",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					Provider: "AWS",
				},
			},
		),
		Entry(
			"Configure Azure private endpoint",
			Label("focus-azure-private-endpoint"),
			func(ctx context.Context) *model.TestDataProvider {
				return model.DataProvider(ctx, "azure-pe-1", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			&akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name: "azure-pe-1",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					Provider: "AZURE",
				},
			},
		),
		Entry(
			"Configure GCP private endpoint",
			Label("focus-gcp-private-endpoint"),
			func(ctx context.Context) *model.TestDataProvider {
				return model.DataProvider(ctx, "gcp-pe-1", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			&akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name: "gcp-pe-1",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					Provider: "GCP",
				},
			},
		),
	)
})

var _ = Describe("Migrate private endpoints from sub-resources to separate custom resources", Label("private-endpoint"), func() {
	var testData *model.TestDataProvider
	var awsPE *akov2.AtlasPrivateEndpoint
	var azurePE *akov2.AtlasPrivateEndpoint
	var gcpPE *akov2.AtlasPrivateEndpoint
	var awsRegion string
	var azureRegion string
	var gcpRegion string
	privateEndpointDetails := map[string]*cloud.PrivateEndpointDetails{}

	_ = BeforeEach(func() {
		checkUpAWSEnvironment()
		checkUpAzureEnvironment()
		checkNSetUpGCPEnvironment()
	})

	_ = AfterEach(func() {
		GinkgoWriter.Println()
		GinkgoWriter.Println("===============================================")
		GinkgoWriter.Println("Operator namespace: " + testData.Resources.Namespace)
		GinkgoWriter.Println("===============================================")
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Delete Project and cluster resources", func() {
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	It("Should migrate a private endpoint configured in a project as sub-resource to a separate custom resource", func(ctx SpecContext) {
		By("Setting up project", func() {
			testData = model.DataProvider(ctx, "migrate-private-endpoint", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())

			actions.ProjectCreationFlow(testData)
		})

		By("Configuring a private endpoint as a sub-resource", func() {
			By("Setting up the private endpoint service", func() {
				Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				var err error
				awsRegion, err = cloud.GetAtlasRegionByProvider("AWS")
				Expect(err).ToNot(HaveOccurred())
				azureRegion, err = cloud.GetAtlasRegionByProvider("AZURE")
				Expect(err).ToNot(HaveOccurred())
				gcpRegion, err = cloud.GetAtlasRegionByProvider("GCP")
				Expect(err).ToNot(HaveOccurred())

				testData.Project.Spec.PrivateEndpoints = []akov2.PrivateEndpoint{
					{
						Provider: "AWS",
						Region:   awsRegion,
					},
					{
						Provider: "AZURE",
						Region:   azureRegion,
					},
					{
						Provider: "GCP",
						Region:   gcpRegion,
					},
				}

				Expect(testData.K8SClient.Update(testData.Context, testData.Project)).To(Succeed())
				Eventually(func(g Gomega) {
					g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
					g.Expect(testData.Project.Status.Conditions).To(ContainElement(conditions.MatchCondition(api.TrueCondition(api.PrivateEndpointServiceReadyType))))
				}).WithTimeout(15 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
			})

			By("Configuring external network", func() {
				Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				action, err := prepareProviderAction(ctx)
				Expect(err).To(BeNil())

				for _, pe := range testData.Project.Spec.PrivateEndpoints {
					peStatus := statusForProvider(testData.Project.Status.PrivateEndpoints, pe.Provider)
					Expect(peStatus).ToNot(BeNil())

					cloudRegion := cloud.MapCloudProviderRegion(string(pe.Provider), pe.Region)
					Expect(cloudRegion).ToNot(BeEmpty())

					switch pe.Provider {
					case "AWS":
						awsConfig, err := cloud.GenerateCloudConfig[cloud.AWSConfig](string(pe.Provider), cloudRegion, testData.Resources.KeyName)
						Expect(err).ToNot(HaveOccurred())

						Expect(action.SetupNetwork(ctx, pe.Provider, cloud.WithAWSConfig(awsConfig))).ToNot(BeEmpty())
						privateEndpointDetails[string(pe.Provider)] = action.SetupPrivateEndpoint(nil, &cloud.AWSPrivateEndpointRequest{
							ID:          fmt.Sprintf("aws-e2e-pe-%s", testData.Resources.TestID),
							Region:      awsConfig.Region,
							ServiceName: peStatus.ServiceName,
						})
					case "AZURE":
						azureConfig, err := cloud.GenerateCloudConfig[cloud.AzureConfig](string(pe.Provider), cloudRegion, testData.Resources.KeyName)
						Expect(err).ToNot(HaveOccurred())

						Expect(action.SetupNetwork(ctx, pe.Provider, cloud.WithAzureConfig(azureConfig))).ToNot(BeEmpty())
						privateEndpointDetails[string(pe.Provider)] = action.SetupPrivateEndpoint(nil, &cloud.AzurePrivateEndpointRequest{
							ID:                fmt.Sprintf("azure-e2e-pe-%s", testData.Resources.TestID),
							Region:            azureConfig.Region,
							ServiceResourceID: peStatus.ServiceResourceID,
							SubnetName:        randomKeyFromMap(azureConfig.Subnets),
						})
					case "GCP":
						gcpConfig, err := cloud.GenerateCloudConfig[cloud.GCPConfig](string(pe.Provider), cloudRegion, testData.Resources.KeyName)
						Expect(err).ToNot(HaveOccurred())

						Expect(action.SetupNetwork(ctx, pe.Provider, cloud.WithGCPConfig(gcpConfig))).ToNot(BeEmpty())
						privateEndpointDetails[string(pe.Provider)] = action.SetupPrivateEndpoint(nil, &cloud.GCPPrivateEndpointRequest{
							ID:         fmt.Sprintf("pe-migration-gcp--%s-%s", pe.EndpointGroupName, testData.Resources.TestID),
							Region:     gcpConfig.Region,
							Targets:    peStatus.ServiceAttachmentNames,
							SubnetName: randomKeyFromMap(gcpConfig.Subnets),
						})
					}
				}
			})

			By("Configuring private endpoint with external network details", func() {
				Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())

				for i, pe := range testData.Project.Spec.PrivateEndpoints {
					switch pe.Provider {
					case "AWS":
						pe.ID = privateEndpointDetails[string(pe.Provider)].ID
					case "AZURE":
						pe.ID = privateEndpointDetails[string(pe.Provider)].ID
						pe.IP = privateEndpointDetails[string(pe.Provider)].IP
					case "GCP":
						gcpEndpoints := make([]akov2.GCPEndpoint, 0, len(privateEndpointDetails[string(pe.Provider)].Endpoints))
						for _, ep := range privateEndpointDetails[string(pe.Provider)].Endpoints {
							gcpEndpoints = append(
								gcpEndpoints,
								akov2.GCPEndpoint{
									EndpointName: ep.Name,
									IPAddress:    ep.IP,
								},
							)
						}
						pe.GCPProjectID = privateEndpointDetails[string(pe.Provider)].GCPProjectID
						pe.EndpointGroupName = privateEndpointDetails[string(pe.Provider)].EndpointGroupName
						pe.Endpoints = gcpEndpoints
					}

					testData.Project.Spec.PrivateEndpoints[i] = pe
				}

				Expect(testData.K8SClient.Update(testData.Context, testData.Project)).To(Succeed())
				Eventually(func(g Gomega) {
					expectedConditions := conditions.MatchConditions(
						api.TrueCondition(api.PrivateEndpointServiceReadyType),
						api.TrueCondition(api.PrivateEndpointReadyType),
						api.TrueCondition(api.ReadyType),
					)

					g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
					g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
				}).WithTimeout(15 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
			})
		})
		//nolint:dupl
		By("Stopping reconciling project and its sub-resources", func() {
			Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Annotations[customresource.ReconciliationPolicyAnnotation] = customresource.ReconciliationPolicySkip
			testData.Project.Spec.PrivateEndpoints = nil

			Expect(testData.K8SClient.Update(testData.Context, testData.Project)).To(Succeed())
			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Generation).ToNot(Equal(testData.Project.Status.ObservedGeneration))
				g.Expect(testData.Project.Status.Conditions).To(ContainElement(conditions.MatchCondition(api.TrueCondition(api.ReadyType))))
			}).WithTimeout(15 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
		})

		By("Migrate private endpoint as separate custom resource", func() {
			//nolint:dupl
			By("Migrating AWS private endpoint", func() {
				awsPE = &akov2.AtlasPrivateEndpoint{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pe-aws-" + testData.Resources.TestID,
						Namespace: testData.Resources.Namespace,
					},
					Spec: akov2.AtlasPrivateEndpointSpec{
						ProjectDualReference: akov2.ProjectDualReference{
							ProjectRef: &common.ResourceRefNamespaced{
								Name:      testData.Project.Name,
								Namespace: testData.Project.Namespace,
							},
						},
						Provider: "AWS",
						Region:   awsRegion,
						AWSConfiguration: []akov2.AWSPrivateEndpointConfiguration{
							{
								ID: privateEndpointDetails["AWS"].ID,
							},
						},
					},
				}

				Expect(testData.K8SClient.Create(testData.Context, awsPE)).To(Succeed())
				Eventually(func(g Gomega) { //nolint:dupl
					expectedConditions := conditions.MatchConditions(
						api.TrueCondition(api.PrivateEndpointServiceReady),
						api.TrueCondition(api.PrivateEndpointReady),
						api.TrueCondition(api.ReadyType),
					)
					g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(awsPE), awsPE)).To(Succeed())
					g.Expect(awsPE.Status.Conditions).To(ContainElements(expectedConditions))
					g.Expect(awsPE.Status.ServiceStatus).To(Equal("AVAILABLE"))
					for _, eStatus := range awsPE.Status.Endpoints {
						g.Expect(eStatus.Status).To(Equal("AVAILABLE"))
					}
				}).WithTimeout(15 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
			})

			By("Migrating AZURE private endpoint", func() {
				azurePE = &akov2.AtlasPrivateEndpoint{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pe-azure-" + testData.Resources.TestID,
						Namespace: testData.Resources.Namespace,
					},
					Spec: akov2.AtlasPrivateEndpointSpec{
						ProjectDualReference: akov2.ProjectDualReference{
							ProjectRef: &common.ResourceRefNamespaced{
								Name:      testData.Project.Name,
								Namespace: testData.Project.Namespace,
							},
						},
						Provider: "AZURE",
						Region:   azureRegion,
						AzureConfiguration: []akov2.AzurePrivateEndpointConfiguration{
							{
								ID: privateEndpointDetails["AZURE"].ID,
								IP: privateEndpointDetails["AZURE"].IP,
							},
						},
					},
				}

				Expect(testData.K8SClient.Create(testData.Context, azurePE)).To(Succeed())
				Eventually(func(g Gomega) { //nolint:dupl
					expectedConditions := conditions.MatchConditions(
						api.TrueCondition(api.PrivateEndpointServiceReady),
						api.TrueCondition(api.PrivateEndpointReady),
						api.TrueCondition(api.ReadyType),
					)
					g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(azurePE), azurePE)).To(Succeed())
					g.Expect(azurePE.Status.Conditions).To(ContainElements(expectedConditions))
					g.Expect(azurePE.Status.ServiceStatus).To(Equal("AVAILABLE"))
					for _, eStatus := range azurePE.Status.Endpoints {
						g.Expect(eStatus.Status).To(Equal("AVAILABLE"))
					}
				}).WithTimeout(15 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
			})

			By("Migrating GCP private endpoint", func() {
				endpoints := make([]akov2.GCPPrivateEndpoint, 0, len(privateEndpointDetails["GCP"].Endpoints))
				for _, ep := range privateEndpointDetails["GCP"].Endpoints {
					endpoints = append(
						endpoints,
						akov2.GCPPrivateEndpoint{
							Name: ep.Name,
							IP:   ep.IP,
						},
					)
				}

				gcpPE = &akov2.AtlasPrivateEndpoint{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pe-gcp-" + testData.Resources.TestID,
						Namespace: testData.Resources.Namespace,
					},
					Spec: akov2.AtlasPrivateEndpointSpec{
						ProjectDualReference: akov2.ProjectDualReference{
							ProjectRef: &common.ResourceRefNamespaced{
								Name:      testData.Project.Name,
								Namespace: testData.Project.Namespace,
							},
						},
						Provider: "GCP",
						Region:   gcpRegion,
						GCPConfiguration: []akov2.GCPPrivateEndpointConfiguration{
							{
								ProjectID: privateEndpointDetails["GCP"].GCPProjectID,
								GroupName: privateEndpointDetails["GCP"].EndpointGroupName,
								Endpoints: endpoints,
							},
						},
					},
				}

				Expect(testData.K8SClient.Create(testData.Context, gcpPE)).To(Succeed())
				Eventually(func(g Gomega) { //nolint:dupl
					expectedConditions := conditions.MatchConditions(
						api.TrueCondition(api.PrivateEndpointServiceReady),
						api.TrueCondition(api.PrivateEndpointReady),
						api.TrueCondition(api.ReadyType),
					)
					g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(gcpPE), gcpPE)).To(Succeed())
					g.Expect(gcpPE.Status.Conditions).To(ContainElements(expectedConditions))
					g.Expect(gcpPE.Status.ServiceStatus).To(Equal("AVAILABLE"))
					for _, eStatus := range gcpPE.Status.Endpoints {
						g.Expect(eStatus.Status).To(Equal("AVAILABLE"))
					}
				}).WithTimeout(15 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
			})
		})
		//nolint:dupl
		By("Restating project reconciliation", func() {
			Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			delete(testData.Project.Annotations, customresource.ReconciliationPolicyAnnotation)

			Expect(testData.K8SClient.Update(testData.Context, testData.Project)).To(Succeed())
			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElement(conditions.MatchCondition(api.TrueCondition(api.ReadyType))))
			}).WithTimeout(15 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
		})

		By("Updating project doesn't affect private endpoint", func() {
			Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Spec.Settings = &akov2.ProjectSettings{
				IsSchemaAdvisorEnabled: pointer.MakePtr(true),
			}

			Expect(testData.K8SClient.Update(testData.Context, testData.Project)).To(Succeed())
			Eventually(func(g Gomega) {
				notExpectedConditions := conditions.MatchConditions(
					api.TrueCondition(api.PrivateEndpointServiceReady),
					api.TrueCondition(api.PrivateEndpointReady),
					api.FalseCondition(api.PrivateEndpointServiceReady),
					api.FalseCondition(api.PrivateEndpointReady),
				)

				g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).ToNot(ContainElements(notExpectedConditions))
				g.Expect(testData.Project.Status.Conditions).To(ContainElement(conditions.MatchCondition(api.TrueCondition(api.ReadyType))))
			}).WithTimeout(15 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
		})

		By("Private endpoint are still ready", func() {
			Eventually(func(g Gomega) { //nolint:dupl
				expectedConditions := conditions.MatchConditions(
					api.TrueCondition(api.PrivateEndpointServiceReady),
					api.TrueCondition(api.PrivateEndpointReady),
					api.TrueCondition(api.ReadyType),
				)
				for _, pe := range []*akov2.AtlasPrivateEndpoint{awsPE, azurePE, gcpPE} {
					g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(pe), pe)).To(Succeed())
					g.Expect(pe.Status.Conditions).To(ContainElements(expectedConditions))
					g.Expect(pe.Status.ServiceStatus).To(Equal("AVAILABLE"))
					for _, eStatus := range pe.Status.Endpoints {
						g.Expect(eStatus.Status).To(Equal("AVAILABLE"))
					}
				}
			}).WithTimeout(15 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
		})

		By("Removing private endpoints", func() {
			Expect(testData.K8SClient.Delete(testData.Context, awsPE)).To(Succeed())
			Expect(testData.K8SClient.Delete(testData.Context, azurePE)).To(Succeed())
			Expect(testData.K8SClient.Delete(testData.Context, gcpPE)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(awsPE), awsPE)).ShouldNot(Succeed())
				g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(azurePE), azurePE)).ShouldNot(Succeed())
				g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(gcpPE), gcpPE)).ShouldNot(Succeed())
			}).WithTimeout(15 * time.Minute).WithPolling(20 * time.Second).Should(Succeed())
		})
	})
})

var _ = Describe("Independent resource should no conflict with sub-resource", Label("private-endpoint"), func() {
	var testData *model.TestDataProvider
	var awsPE *akov2.AtlasPrivateEndpoint

	_ = BeforeEach(func() {
		checkUpAWSEnvironment()
		checkUpAzureEnvironment()
		checkNSetUpGCPEnvironment()
	})

	_ = AfterEach(func() {
		GinkgoWriter.Println()
		GinkgoWriter.Println("===============================================")
		GinkgoWriter.Println("Operator namespace: " + testData.Resources.Namespace)
		GinkgoWriter.Println("===============================================")
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Delete Project and cluster resources", func() {
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	It("Should migrate a private endpoint configured in a project as sub-resource to a separate custom resource", func(ctx SpecContext) {
		By("Setting up project", func() {
			testData = model.DataProvider(ctx, "migrate-private-endpoint", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())

			actions.ProjectCreationFlow(testData)
		})

		//nolint:dupl
		By("Creating AWS private endpoint", func() {
			awsRegion, err := cloud.GetAtlasRegionByProvider("AWS")
			Expect(err).ToNot(HaveOccurred())

			awsPE = &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe-aws-" + testData.Resources.TestID,
					Namespace: testData.Resources.Namespace,
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      testData.Project.Name,
							Namespace: testData.Project.Namespace,
						},
					},
					Provider: "AWS",
					Region:   awsRegion,
				},
			}

			Expect(testData.K8SClient.Create(testData.Context, awsPE)).To(Succeed())
			Eventually(func(g Gomega) { //nolint:dupl
				expectedConditions := conditions.MatchConditions(
					api.TrueCondition(api.PrivateEndpointServiceReady),
					api.FalseCondition(api.PrivateEndpointReady).
						WithReason(string(workflow.PrivateEndpointConfigurationPending)).
						WithMessageRegexp("waiting for private endpoint configuration from customer side"),
					api.FalseCondition(api.ReadyType),
				)
				g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(awsPE), awsPE)).To(Succeed())
				g.Expect(awsPE.Status.Conditions).To(ContainElements(expectedConditions))
				g.Expect(awsPE.Status.ServiceStatus).To(Equal("AVAILABLE"))
			}).WithTimeout(15 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
		})

		By("Updating project doesn't affect private endpoint", func() {
			Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Spec.Settings = &akov2.ProjectSettings{
				IsSchemaAdvisorEnabled:            pointer.MakePtr(true),
				IsRealtimePerformancePanelEnabled: pointer.MakePtr(true),
			}

			Expect(testData.K8SClient.Update(testData.Context, testData.Project)).To(Succeed())
			Eventually(func(g Gomega) {
				notExpectedConditions := conditions.MatchConditions(
					api.TrueCondition(api.PrivateEndpointServiceReady),
					api.TrueCondition(api.PrivateEndpointReady),
					api.FalseCondition(api.PrivateEndpointServiceReady),
					api.FalseCondition(api.PrivateEndpointReady),
				)

				g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).ToNot(ContainElements(notExpectedConditions))
				g.Expect(testData.Project.Status.Conditions).To(ContainElement(conditions.MatchCondition(api.TrueCondition(api.ReadyType))))
			}).WithTimeout(15 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
		})

		By("Private endpoint are still ready", func() {
			Eventually(func(g Gomega) { //nolint:dupl
				expectedConditions := conditions.MatchConditions(
					api.TrueCondition(api.PrivateEndpointServiceReady),
					api.FalseCondition(api.PrivateEndpointReady).
						WithReason(string(workflow.PrivateEndpointConfigurationPending)).
						WithMessageRegexp("waiting for private endpoint configuration from customer side"),
					api.FalseCondition(api.ReadyType),
				)

				g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(awsPE), awsPE)).To(Succeed())
				g.Expect(awsPE.Status.Conditions).To(ContainElements(expectedConditions))
				g.Expect(awsPE.Status.ServiceStatus).To(Equal("AVAILABLE"))
			}).WithTimeout(15 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
		})

		By("Removing private endpoints", func() {
			Expect(testData.K8SClient.Delete(testData.Context, awsPE)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(awsPE), awsPE)).ShouldNot(Succeed())
			}).WithTimeout(15 * time.Minute).WithPolling(20 * time.Second).Should(Succeed())
		})
	})
})

func statusForProvider(peStatus []status.ProjectPrivateEndpoint, providerName provider.ProviderName) *status.ProjectPrivateEndpoint {
	for _, s := range peStatus {
		if s.Provider == providerName {
			return &s
		}
	}

	return nil
}

func randomKeyFromMap[K comparable, V any](m map[K]V) K {
	for k := range m {
		return k
	}

	return *new(K)
}
