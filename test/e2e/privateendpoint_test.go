package e2e_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/conditions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/cloud"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/resources"
)

var _ = Describe("Private Endpoints", Label("private-endpoint"), func() {
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
		func(test *model.TestDataProvider, pe *akov2.AtlasPrivateEndpoint) {
			var privateEndpointDetails *cloud.PrivateEndpointDetails

			testData = test
			actions.ProjectCreationFlow(test)

			By("Referring to a project", func() {
				pe.Namespace = test.Resources.Namespace
				pe.Spec.Project = &common.ResourceRefNamespaced{
					Name:      test.Project.Name,
					Namespace: test.Project.Namespace,
				}
			})

			By("Creating private endpoint", func() {
				Expect(test.K8SClient.Create(test.Context, pe)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(testData.K8SClient.Get(test.Context, client.ObjectKeyFromObject(pe), pe)).To(Succeed())
					g.Expect(pe.Status.ServiceStatus).To(Equal("AVAILABLE"))
					g.Expect(resources.CheckCondition(testData.K8SClient, pe, api.TrueCondition(api.PrivateEndpointServiceReady))).Should(BeTrue())
				}).WithTimeout(15 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
			})

			By("Configuring external network", func() {
				Expect(testData.K8SClient.Get(test.Context, client.ObjectKeyFromObject(pe), pe)).To(Succeed())
				action, err := prepareProviderAction()
				Expect(err).To(BeNil())

				switch pe.Spec.Provider {
				case "AWS":
					awsConfig := cloud.AWSConfig{
						Region:        "eu-central-1",
						VPC:           fmt.Sprintf("%s-%s", pe.Name, test.Resources.TestID),
						EnableCleanup: true,
					}

					Expect(action.SetupNetwork(provider.ProviderName(pe.Spec.Provider), cloud.WithAWSConfig(&awsConfig))).ToNot(BeEmpty())
					privateEndpointDetails = action.SetupPrivateEndpoint(&cloud.AWSPrivateEndpointRequest{
						ID:          "aws-e2e-pe",
						Region:      "eu-central-1",
						ServiceName: pe.Status.ServiceName,
					})
				case "AZURE":
					azureConfig := cloud.AzureConfig{
						Region:        "northeurope",
						VPC:           fmt.Sprintf("%s-%s", pe.Name, test.Resources.TestID),
						EnableCleanup: true,
					}

					Expect(action.SetupNetwork(provider.ProviderName(pe.Spec.Provider), cloud.WithAzureConfig(&azureConfig))).ToNot(BeEmpty())
					privateEndpointDetails = action.SetupPrivateEndpoint(&cloud.AzurePrivateEndpointRequest{
						ID:                "azure-e2e-pe",
						Region:            "northeurope",
						ServiceResourceID: pe.Status.ResourceID,
						SubnetName:        cloud.Subnet1Name,
					})
				case "GCP":
					gcpConfig := cloud.GCPConfig{
						Region:        "europe-west3",
						VPC:           fmt.Sprintf("%s-%s", pe.Name, test.Resources.TestID),
						EnableCleanup: true,
					}

					Expect(action.SetupNetwork(provider.ProviderName(pe.Spec.Provider), cloud.WithGCPConfig(&gcpConfig))).ToNot(BeEmpty())
					privateEndpointDetails = action.SetupPrivateEndpoint(&cloud.GCPPrivateEndpointRequest{
						ID:         fmt.Sprintf("%s-%s", pe.Name, test.Resources.TestID),
						Region:     "europe-west3",
						Targets:    pe.Status.ServiceAttachmentNames,
						SubnetName: cloud.Subnet1Name,
					})
				}
			})

			By("Configuring private endpoint with external network details", func() {
				Expect(testData.K8SClient.Get(test.Context, client.ObjectKeyFromObject(pe), pe)).To(Succeed())

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

				Expect(test.K8SClient.Update(test.Context, pe)).To(Succeed())
				Eventually(func(g Gomega) { //nolint:dupl
					g.Expect(testData.K8SClient.Get(test.Context, client.ObjectKeyFromObject(pe), pe)).To(Succeed())
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
				Expect(test.K8SClient.Delete(test.Context, pe)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(testData.K8SClient.Get(test.Context, client.ObjectKeyFromObject(pe), pe)).ShouldNot(Succeed())
				}).WithTimeout(15 * time.Minute).WithPolling(20 * time.Second).Should(Succeed())
			})
		},
		Entry(
			"Configure AWS private endpoint",
			Label("aws-private-endpoint"),
			model.DataProvider(
				"aws-pe-1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			&akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-pe-1",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					Provider: "AWS",
					Region:   "EU_CENTRAL_1",
				},
			},
		),
		Entry(
			"Configure Azure private endpoint",
			Label("azure-private-endpoint"),
			model.DataProvider(
				"azure-pe-1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			&akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name: "azure-pe-1",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					Provider: "AZURE",
					Region:   "EUROPE_NORTH",
				},
			},
		),
		Entry(
			"Configure GCP private endpoint",
			Label("gcp-private-endpoint"),
			model.DataProvider(
				"gcp-pe-1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			&akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name: "gcp-pe-1",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					Provider: "GCP",
					Region:   "EUROPE_WEST_3",
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

	It("Should migrate a private endpoint configured in a project as sub-resource to a separate custom resource", func() {
		By("Setting up project", func() {
			testData = model.DataProvider(
				"project-with-private-endpoint",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject())

			actions.ProjectCreationFlow(testData)
		})

		By("Configuring a private endpoint as a sub-resource", func() {
			By("Setting up the private endpoint service", func() {
				Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())

				testData.Project.Spec.PrivateEndpoints = []akov2.PrivateEndpoint{
					{
						Provider: "AWS",
						Region:   "EU_CENTRAL_1",
					},
					{
						Provider: "AZURE",
						Region:   "EUROPE_NORTH",
					},
					{
						Provider: "GCP",
						Region:   "EUROPE_WEST_3",
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
				action, err := prepareProviderAction()
				Expect(err).To(BeNil())

				for _, pe := range testData.Project.Spec.PrivateEndpoints {
					peStatus := statusForProvider(testData.Project.Status.PrivateEndpoints, pe.Provider)
					Expect(peStatus).ToNot(BeNil())

					switch pe.Provider {
					case "AWS":
						awsConfig := cloud.AWSConfig{
							Region:        "eu-central-1",
							VPC:           fmt.Sprintf("pe-migration-aws-%s", testData.Resources.TestID),
							EnableCleanup: true,
						}

						Expect(action.SetupNetwork(pe.Provider, cloud.WithAWSConfig(&awsConfig))).ToNot(BeEmpty())
						privateEndpointDetails[string(pe.Provider)] = action.SetupPrivateEndpoint(&cloud.AWSPrivateEndpointRequest{
							ID:          "aws-e2e-pe",
							Region:      "eu-central-1",
							ServiceName: peStatus.ServiceName,
						})
					case "AZURE":
						azureConfig := cloud.AzureConfig{
							Region:        "northeurope",
							VPC:           fmt.Sprintf("pe-migration-azure-%s", testData.Resources.TestID),
							EnableCleanup: true,
						}

						Expect(action.SetupNetwork(pe.Provider, cloud.WithAzureConfig(&azureConfig))).ToNot(BeEmpty())
						privateEndpointDetails[string(pe.Provider)] = action.SetupPrivateEndpoint(&cloud.AzurePrivateEndpointRequest{
							ID:                "azure-e2e-pe",
							Region:            "northeurope",
							ServiceResourceID: peStatus.ServiceResourceID,
							SubnetName:        cloud.Subnet1Name,
						})
					case "GCP":
						gcpConfig := cloud.GCPConfig{
							Region:        "europe-west3",
							VPC:           fmt.Sprintf("pe-migration-gcp-%s", testData.Resources.TestID),
							EnableCleanup: true,
						}

						Expect(action.SetupNetwork(pe.Provider, cloud.WithGCPConfig(&gcpConfig))).ToNot(BeEmpty())
						privateEndpointDetails[string(pe.Provider)] = action.SetupPrivateEndpoint(&cloud.GCPPrivateEndpointRequest{
							ID:         fmt.Sprintf("pe-migration-gcp-%s", testData.Resources.TestID),
							Region:     "europe-west3",
							Targets:    peStatus.ServiceAttachmentNames,
							SubnetName: cloud.Subnet1Name,
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

		By("Stopping reconciling project and its sub-resources", func() {
			Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Annotations[customresource.ReconciliationPolicyAnnotation] = customresource.ReconciliationPolicySkip
			testData.Project.Spec.PrivateEndpoints = nil

			Expect(testData.K8SClient.Update(testData.Context, testData.Project)).To(Succeed())
			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Generation).ToNot(Equal(testData.Project.Status.ObservedGeneration))
				g.Expect(customresource.AnnotationLastSkippedConfiguration).To(BeKeyOf(testData.Project.GetAnnotations()))
				g.Expect(testData.Project.Status.Conditions).To(ContainElement(conditions.MatchCondition(api.TrueCondition(api.ReadyType))))
			}).WithTimeout(15 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
		})

		By("Migrate private endpoint as separate custom resource", func() {
			By("Migrating AWS private endpoint", func() {
				awsPE = &akov2.AtlasPrivateEndpoint{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pe-aws-" + testData.Resources.TestID,
						Namespace: testData.Resources.Namespace,
					},
					Spec: akov2.AtlasPrivateEndpointSpec{
						Project: &common.ResourceRefNamespaced{
							Name:      testData.Project.Name,
							Namespace: testData.Project.Namespace,
						},
						Provider: "AWS",
						Region:   "EU_CENTRAL_1",
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
					g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
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
						Project: &common.ResourceRefNamespaced{
							Name:      testData.Project.Name,
							Namespace: testData.Project.Namespace,
						},
						Provider: "AZURE",
						Region:   "EUROPE_NORTH",
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
					g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
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
						Project: &common.ResourceRefNamespaced{
							Name:      testData.Project.Name,
							Namespace: testData.Project.Namespace,
						},
						Provider: "GCP",
						Region:   "EUROPE_WEST_3",
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
					g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
					g.Expect(gcpPE.Status.ServiceStatus).To(Equal("AVAILABLE"))
					for _, eStatus := range gcpPE.Status.Endpoints {
						g.Expect(eStatus.Status).To(Equal("AVAILABLE"))
					}
				}).WithTimeout(15 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
			})
		})

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
					g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
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

func statusForProvider(peStatus []status.ProjectPrivateEndpoint, providerName provider.ProviderName) *status.ProjectPrivateEndpoint {
	for _, s := range peStatus {
		if s.Provider == providerName {
			return &s
		}
	}

	return nil
}
