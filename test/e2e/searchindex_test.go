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
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasdeployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/api/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

const (
	DBTraining                 = "sample_training"
	DBTrainingCollectionGrades = "grades"
	DBTrainingCollectionRoutes = "routes"
)

var _ = Describe("Atlas Search Index", Label("atlas-search-index"), func() {
	var testData *model.TestDataProvider
	var searchIndexConfig *akov2.AtlasSearchIndexConfig

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

	It("Create and delete SEARCH and VECTOR SEARCH indexes", func(ctx SpecContext) {
		testData = model.DataProvider(ctx, "atlas-search-nodes", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject()).WithInitialDeployments(data.CreateAdvancedDeployment("search-nodes-test"))
		atlasClient = atlas.GetClientOrFail()

		By("Setting up project", func() {
			actions.ProjectCreationFlow(testData)
		})

		By("Creating a deployment without search indexes", func() {
			testData.InitialDeployments[0].Namespace = testData.Resources.Namespace
			Expect(testData.K8SClient.Create(testData.Context, testData.InitialDeployments[0])).To(Succeed())

			Eventually(func(g Gomega) bool {
				g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{
					Name:      testData.InitialDeployments[0].Name,
					Namespace: testData.InitialDeployments[0].Namespace,
				}, testData.InitialDeployments[0])).To(Succeed())
				for _, condition := range testData.InitialDeployments[0].Status.Conditions {
					if condition.Type == api.DeploymentReadyType {
						return condition.Status == v1.ConditionTrue
					}
				}
				return false
			}).WithTimeout(20 * time.Minute).Should(BeTrue())

		})

		By("Loading sample dataset into a cluster", func() {
			sampleDataSet, _, err := atlasClient.Client.ClustersApi.LoadSampleDataset(testData.Context,
				testData.Project.ID(),
				testData.InitialDeployments[0].GetDeploymentName()).Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(sampleDataSet).NotTo(BeNil())
			Expect(sampleDataSet.Id).NotTo(BeNil())

			Eventually(func(g Gomega) {
				sampleDataStatus, _, err := atlasClient.Client.ClustersApi.GetSampleDatasetLoadStatus(testData.Context,
					testData.Project.ID(),
					*sampleDataSet.Id).Execute()
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(sampleDataStatus).NotTo(BeNil())
				g.Expect(sampleDataStatus.State).NotTo(BeNil())
				g.Expect(*sampleDataStatus.State).To(BeEquivalentTo("COMPLETED"))
			}).WithPolling(10 * time.Second).WithTimeout(10 * time.Minute).Should(Succeed())
		})

		By("Creating a simple search index configuration", func() {
			searchIndexConfig = &akov2.AtlasSearchIndexConfig{
				TypeMeta: metav1.TypeMeta{
					Kind:       "AtlasSearchIndexConfig",
					APIVersion: "atlas.mongodb.com/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-search-index-config",
					Namespace: testData.InitialDeployments[0].Namespace,
				},
				Spec: akov2.AtlasSearchIndexConfigSpec{
					Analyzer:       pointer.MakePtr("lucene.standard"),
					SearchAnalyzer: pointer.MakePtr("lucene.standard"),
				},
			}
			Expect(testData.K8SClient.Create(testData.Context, searchIndexConfig)).To(Succeed())
		})

		By("Creating one search index, type: SEARCH", func() {
			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{
					Name:      testData.InitialDeployments[0].Name,
					Namespace: testData.InitialDeployments[0].Namespace,
				}, testData.InitialDeployments[0])).To(Succeed())

				searchIndexesToCreate := []akov2.SearchIndex{
					{
						Name:           "test-search-index",
						Type:           atlasdeployment.IndexTypeSearch,
						DBName:         DBTraining,
						CollectionName: DBTrainingCollectionRoutes,
						Search: &akov2.Search{
							Mappings: &akov2.Mappings{Dynamic: pointer.MakePtr(true)},
							SearchConfigurationRef: common.ResourceRefNamespaced{
								Name:      searchIndexConfig.GetName(),
								Namespace: searchIndexConfig.GetNamespace(),
							},
						},
					},
				}
				testData.InitialDeployments[0].Spec.DeploymentSpec.SearchIndexes = searchIndexesToCreate
				g.Expect(testData.K8SClient.Update(testData.Context, testData.InitialDeployments[0])).To(Succeed())
			}).WithPolling(10 * time.Second).WithTimeout(10 * time.Minute).Should(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{
					Name:      testData.InitialDeployments[0].Name,
					Namespace: testData.InitialDeployments[0].Namespace,
				}, testData.InitialDeployments[0])).To(Succeed())
				g.Expect(len(testData.InitialDeployments[0].Status.SearchIndexes)).To(BeEquivalentTo(1))
				g.Expect(testData.InitialDeployments[0].Status.SearchIndexes[0].Status).To(BeEquivalentTo(status.SearchIndexStatusReady))
			}).WithPolling(10 * time.Second).WithTimeout(40 * time.Minute).Should(Succeed())
		})

		By("Creating one search index, type: VECTOR SEARCH", func() {
			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{
					Name:      testData.InitialDeployments[0].Name,
					Namespace: testData.InitialDeployments[0].Namespace,
				}, testData.InitialDeployments[0])).To(Succeed())

				vectorSearchIndexToCreate := akov2.SearchIndex{
					Name:           "test-vector-search-index",
					Type:           atlasdeployment.IndexTypeVector,
					DBName:         DBTraining,
					CollectionName: DBTrainingCollectionGrades,
					VectorSearch: &akov2.VectorSearch{
						Fields: &apiextensions.JSON{
							Raw: []byte(`[{
      "type": "vector",
      "path": "student_id",
      "numDimensions": 1536,
      "similarity": "euclidean"
}]`),
						},
					},
				}
				testData.InitialDeployments[0].Spec.DeploymentSpec.SearchIndexes = append(
					testData.InitialDeployments[0].Spec.DeploymentSpec.SearchIndexes, vectorSearchIndexToCreate)
				g.Expect(testData.K8SClient.Update(testData.Context, testData.InitialDeployments[0])).To(Succeed())
			}).WithPolling(10 * time.Second).WithTimeout(10 * time.Minute).Should(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{
					Name:      testData.InitialDeployments[0].Name,
					Namespace: testData.InitialDeployments[0].Namespace,
				}, testData.InitialDeployments[0])).To(Succeed())
				g.Expect(len(testData.InitialDeployments[0].Status.SearchIndexes)).To(BeEquivalentTo(2))
				g.Expect(testData.InitialDeployments[0].Status.SearchIndexes[1].Status).To(BeEquivalentTo(status.SearchIndexStatusReady))
			}).WithPolling(10 * time.Second).WithTimeout(10 * time.Minute).Should(Succeed())
		})

		By("Deleting the VECTOR SEARCH index", func() {
			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{
					Name:      testData.InitialDeployments[0].Name,
					Namespace: testData.InitialDeployments[0].Namespace,
				}, testData.InitialDeployments[0])).To(Succeed())
				testData.InitialDeployments[0].Spec.DeploymentSpec.SearchIndexes = []akov2.SearchIndex{
					testData.InitialDeployments[0].Spec.DeploymentSpec.SearchIndexes[0],
				}

				g.Expect(testData.K8SClient.Update(testData.Context, testData.InitialDeployments[0])).To(Succeed())
			}).WithPolling(10 * time.Second).WithTimeout(10 * time.Minute).Should(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{
					Name:      testData.InitialDeployments[0].Name,
					Namespace: testData.InitialDeployments[0].Namespace,
				}, testData.InitialDeployments[0])).To(Succeed())
				g.Expect(len(testData.InitialDeployments[0].Status.SearchIndexes)).To(BeEquivalentTo(1))
			}).WithPolling(10 * time.Second).WithTimeout(10 * time.Minute).Should(Succeed())
		})
		By("Deleting the SEARCH index", func() {
			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{
					Name:      testData.InitialDeployments[0].Name,
					Namespace: testData.InitialDeployments[0].Namespace,
				}, testData.InitialDeployments[0])).To(Succeed())
				testData.InitialDeployments[0].Spec.DeploymentSpec.SearchIndexes = nil

				g.Expect(testData.K8SClient.Update(testData.Context, testData.InitialDeployments[0])).To(Succeed())
			}).WithPolling(10 * time.Second).WithTimeout(10 * time.Minute).Should(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{
					Name:      testData.InitialDeployments[0].Name,
					Namespace: testData.InitialDeployments[0].Namespace,
				}, testData.InitialDeployments[0])).To(Succeed())
				g.Expect(len(testData.InitialDeployments[0].Status.SearchIndexes)).To(BeZero())
			}).WithPolling(10 * time.Second).WithTimeout(10 * time.Minute).Should(Succeed())
		})

	})
})
