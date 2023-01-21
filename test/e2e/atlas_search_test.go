package e2e

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"

	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("DeploymentAtlasSearch", Label("atlas-search"), func() {
	var testData *model.TestDataProvider

	atlasSearchConfig := &v1.AtlasSearch{
		CustomAnalyzers: []v1.CustomAnalyzer{
			{
				Name:         "my_analyzer",
				BaseAnalyzer: "lucene.standard",
			},
		},
		Databases: []v1.AtlasSearchDatabase{
			{
				Database: "sample_mflix",
				Collections: []v1.AtlasSearchCollection{
					{
						CollectionName: "movies",
						Indexes: []v1.SearchIndex{
							{
								Name:     "movies_ix",
								Analyzer: "my_analyzer",
								Mappings: v1.IndexMapping{
									Dynamic: true,
								},
							},
						},
					},
				},
			},
			{
				Database: "sample_restaurants",
				Collections: []v1.AtlasSearchCollection{
					{
						CollectionName: "restaurants",
						Indexes: []v1.SearchIndex{
							{
								Name: "rest_ix",
								Mappings: v1.IndexMapping{
									Dynamic: false,
									Fields: &v1.FieldMapping{
										"name": []map[string]interface{}{
											{
												"type": "string",
											},
										},
										"address": []map[string]interface{}{
											{
												"type":    "document",
												"dynamic": true,
											},
										},
									},
								},
								Synonyms: []v1.AtlasSearchSynonym{
									{
										Name:     "my_synonym",
										Analyzer: "lucene.standard",
										Source: v1.SynonymSource{
											Collection: "neighborhoods",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

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
		func(test *model.TestDataProvider, atlasSearch *v1.AtlasSearch) {
			testData = test
			actions.ProjectCreationFlow(test)
			AtlasSearchFlow(test, atlasSearch)
		},
		Entry("Test[as-regular-deployment]: Regular Deployment with Atlas Search", Label("as-regular-deployment"),
			model.DataProvider(
				"as-regular-deployment",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()).WithInitialDeployments(data.CreateRegularDeployment("as-regular-deployment")),
			atlasSearchConfig,
		),
		/*
			Entry("Test[as-advanced-deployment]: Advanced Deployment with Atlas Search", Label("as-advanced-deployment"),
				model.DataProvider(
					"as-advanced-deployment",
					model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
					40000,
					[]func(*model.TestDataProvider){},
				).WithProject(data.DefaultProject()).WithInitialDeployments(data.CreateAdvancedDeployment("as-advanced-deployment")),
				atlasSearchConfig,
			),*/
	)
})

func AtlasSearchFlow(userData *model.TestDataProvider, atlasSearch *v1.AtlasSearch) {
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
		}).WithTimeout(30 * time.Minute).Should(BeTrue())
	})

	By("Load sample data", func() {
		atlasClient := atlas.GetClientOrFail()
		ctx := context.Background()
		sampleDataJob, _, err := atlasClient.Client.Clusters.LoadSampleDataset(ctx, userData.Project.ID(), userData.InitialDeployments[0].GetDeploymentName())
		Expect(err).Should(BeNil())
		Expect(sampleDataJob).ShouldNot(BeNil())

		Eventually(func(g Gomega) bool {
			job, _, err := atlasClient.Client.Clusters.GetSampleDatasetStatus(ctx, userData.Project.ID(), sampleDataJob.ID)
			Expect(err).Should(BeNil())
			Expect(job).ShouldNot(BeNil())

			return job.State == "COMPLETED"
		}).WithTimeout(15 * time.Minute).Should(BeTrue())
	})

	By("Applying AtlasSearch config to Deployment", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{
			Name:      userData.InitialDeployments[0].Name,
			Namespace: userData.InitialDeployments[0].Namespace,
		}, userData.InitialDeployments[0])).To(Succeed())

		if userData.InitialDeployments[0].Spec.DeploymentSpec != nil {
			userData.InitialDeployments[0].Spec.DeploymentSpec.AtlasSearch = atlasSearch
		} else {
			userData.InitialDeployments[0].Spec.AdvancedDeploymentSpec.AtlasSearch = atlasSearch
		}

		Expect(userData.K8SClient.Update(userData.Context, userData.InitialDeployments[0])).To(Succeed())
	})

	By("Wait and check status", func() {
		Eventually(func(g Gomega) bool {
			g.Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{
				Name:      userData.InitialDeployments[0].Name,
				Namespace: userData.InitialDeployments[0].Namespace,
			}, userData.InitialDeployments[0])).To(Succeed())

			Expect(userData.InitialDeployments[0].Status.AtlasSearch).ShouldNot(BeNil())
			Expect(userData.InitialDeployments[0].Status.AtlasSearch.Indexes).Should(HaveLen(2))
			Expect(userData.InitialDeployments[0].Status.AtlasSearch.Indexes[0].Status).Should(Equal("ready"))
			Expect(userData.InitialDeployments[0].Status.AtlasSearch.Indexes[1].Status).Should(Equal("ready"))

			for _, condition := range userData.InitialDeployments[0].Status.Conditions {
				if condition.Type == status.AtlasSearchReadyType {
					return condition.Status == corev1.ConditionTrue
				}
			}

			return false
		}).WithTimeout(10 * time.Minute).Should(BeTrue())
	})
}
