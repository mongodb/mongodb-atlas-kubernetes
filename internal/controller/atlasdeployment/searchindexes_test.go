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

package atlasdeployment

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/searchindex"
	searchfake "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/searchindex/fake"
)

func Test_verifyAllIndexesNamesAreUnique(t *testing.T) {
	t.Run("Should return true if all indices names are unique", func(t *testing.T) {
		in := []akov2.SearchIndex{
			{
				Name: "Index-One",
			},
			{
				Name: "Index-Two",
			},
			{
				Name: "Index-Three",
			},
		}
		assert.True(t, verifyAllIndexesNamesAreUnique(in))
	})
	t.Run("Should return false if one index name appeared twice", func(t *testing.T) {
		in := []akov2.SearchIndex{
			{
				Name: "Index-One",
			},
			{
				Name: "Index-Two",
			},
			{
				Name: "Index-One",
			},
		}
		assert.False(t, verifyAllIndexesNamesAreUnique(in))
	})
}

func Test_getIndexesFromDeploymentStatus(t *testing.T) {
	tests := []struct {
		name             string
		deploymentStatus status.AtlasDeploymentStatus
		want             map[string]string
	}{
		{
			name: "Should return valid indexes for some valid indexes in the status",
			deploymentStatus: status.AtlasDeploymentStatus{
				SearchIndexes: []status.DeploymentSearchIndexStatus{
					{
						Name:    "FirstIndex",
						ID:      "1",
						Status:  "",
						Message: "",
					},
					{
						Name:    "SecondIndex",
						ID:      "2",
						Status:  "",
						Message: "",
					},
				},
			},
			want: map[string]string{
				"FirstIndex":  "1",
				"SecondIndex": "2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getIndexesFromDeploymentStatus(tt.deploymentStatus), "getIndexesFromDeploymentStatus(%v)", tt.deploymentStatus)
		})
	}
}

//nolint:dupl
func Test_SearchIndexesReconcile(t *testing.T) {
	t.Run("Should return if indexes names are not unique", func(t *testing.T) {
		deployment := &akov2.AtlasDeployment{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Spec: akov2.AtlasDeploymentSpec{
				DeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name: "testDeployment",
					SearchIndexes: []akov2.SearchIndex{
						{
							Name: "Index1",
						},
						{
							Name: "Index1",
						},
					},
				},
			},
			Status: status.AtlasDeploymentStatus{},
		}

		reconciler := searchIndexesReconcileRequest{
			ctx: &workflow.Context{
				Log:   zap.S(),
				OrgID: "testOrgID",
			},
			deployment: deployment,
		}
		result := reconciler.Handle()
		assert.True(t, reconciler.ctx.HasReason(api.SearchIndexesNamesAreNotUnique))
		assert.False(t, result.IsOk())
	})

	t.Run("Should cleanup indexes with empty IDs when IDLE", func(t *testing.T) {
		searchIndexConfig := &akov2.AtlasSearchIndexConfig{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testConfig",
				Namespace: "testNamespace",
			},
			Spec: akov2.AtlasSearchIndexConfigSpec{
				Analyzer: pointer.MakePtr("testAnalyzer"),
			},
			Status: status.AtlasSearchIndexConfigStatus{},
		}
		IDForStatus := "123"
		deployment := &akov2.AtlasDeployment{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Spec: akov2.AtlasDeploymentSpec{
				DeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name: "testDeployment",
					SearchIndexes: []akov2.SearchIndex{
						{
							Name: "Index1",
							Type: IndexTypeSearch,
							Search: &akov2.Search{
								Synonyms: &([]akov2.Synonym{
									{
										Name:     "testSynonym",
										Analyzer: "testAnalyzer",
										Source: akov2.Source{
											Collection: "testCollection",
										},
									},
								}),
								Mappings: nil,
								SearchConfigurationRef: common.ResourceRefNamespaced{
									Name:      searchIndexConfig.Name,
									Namespace: searchIndexConfig.Namespace,
								},
							},
						},
					},
				},
			},
			Status: status.AtlasDeploymentStatus{
				SearchIndexes: []status.DeploymentSearchIndexStatus{
					{
						Name:    "Index1",
						ID:      IDForStatus,
						Status:  "",
						Message: "",
					},
					{
						Name:    "Index2",
						ID:      "",
						Status:  "",
						Message: "",
					},
				},
			},
		}

		sch := runtime.NewScheme()
		assert.Nil(t, akov2.AddToScheme(sch))
		assert.Nil(t, corev1.AddToScheme(sch))

		idxToReturn := searchindex.NewSearchIndex(
			&deployment.Spec.DeploymentSpec.SearchIndexes[0],
			&searchIndexConfig.Spec)
		fakeAtlasSearch := &searchfake.FakeAtlasSearch{
			GetIndexFunc: func(_ context.Context, _, _, _, _ string) (*searchindex.SearchIndex, error) {
				return idxToReturn, nil
			},
		}

		k8sClient := fake.NewClientBuilder().
			WithScheme(sch).
			WithObjects(deployment, searchIndexConfig).
			Build()

		reconciler := searchIndexesReconcileRequest{
			ctx: &workflow.Context{
				Log:     zap.S(),
				OrgID:   "testOrgID",
				Client:  nil,
				Context: context.Background(),
			},
			deployment:    deployment,
			k8sClient:     k8sClient,
			projectID:     "testProjectID",
			searchService: fakeAtlasSearch,
		}

		result := reconciler.Handle()
		assert.True(t, result.IsOk())
		deployment.UpdateStatus(reconciler.ctx.Conditions(), reconciler.ctx.StatusOptions()...)
		assert.Len(t, deployment.Status.SearchIndexes, 1)
		assert.True(t, deployment.Status.SearchIndexes[0].ID == IDForStatus)
	})

	t.Run("Should proceed with the index Type Search: CREATE INDEX", func(t *testing.T) {
		fakeAtlasSearch := &searchfake.FakeAtlasSearch{
			CreateIndexFunc: func(_ context.Context, _, _ string, _ *searchindex.SearchIndex) (*searchindex.SearchIndex, error) {
				return &searchindex.SearchIndex{
					Status: pointer.MakePtr("NOT STARTED"),
				}, nil
			},
		}

		searchIndexConfig := &akov2.AtlasSearchIndexConfig{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testConfig",
				Namespace: "testNamespace",
			},
			Spec: akov2.AtlasSearchIndexConfigSpec{
				Analyzer: pointer.MakePtr("testAnalyzer"),
			},
			Status: status.AtlasSearchIndexConfigStatus{},
		}

		deployment := &akov2.AtlasDeployment{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Spec: akov2.AtlasDeploymentSpec{
				DeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name: "testDeployment",
					SearchIndexes: []akov2.SearchIndex{
						{
							Name: "Index1",
							Type: IndexTypeSearch,
							Search: &akov2.Search{
								Synonyms: &([]akov2.Synonym{
									{
										Name:     "testSynonym",
										Analyzer: "testAnalyzer",
										Source: akov2.Source{
											Collection: "testCollection",
										},
									},
								}),
								Mappings: nil,
								SearchConfigurationRef: common.ResourceRefNamespaced{
									Name:      searchIndexConfig.Name,
									Namespace: searchIndexConfig.Namespace,
								},
							},
						},
					},
				},
			},
			Status: status.AtlasDeploymentStatus{},
		}

		sch := runtime.NewScheme()
		assert.Nil(t, akov2.AddToScheme(sch))
		assert.Nil(t, corev1.AddToScheme(sch))

		k8sClient := fake.NewClientBuilder().
			WithScheme(sch).
			WithObjects(deployment, searchIndexConfig).
			Build()

		reconciler := searchIndexesReconcileRequest{
			ctx: &workflow.Context{
				Log:     zap.S(),
				OrgID:   "testOrgID",
				Client:  nil,
				Context: context.Background(),
			},
			deployment:    deployment,
			k8sClient:     k8sClient,
			projectID:     "testProjectID",
			searchService: fakeAtlasSearch,
		}
		result := reconciler.Handle()
		assert.False(t, result.IsOk())
	})

	t.Run("Should proceed with index Type Search if it cannot be found: CREATE INDEX", func(t *testing.T) {
		fakeAtlasSearch := &searchfake.FakeAtlasSearch{
			GetIndexFunc: func(_ context.Context, _, _, _, indexID string) (*searchindex.SearchIndex, error) {
				if indexID == "123" {
					return &searchindex.SearchIndex{
						SearchIndex: akov2.SearchIndex{
							Name: "Index1",
							Type: IndexTypeSearch,
						},
						AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{
							Analyzer: pointer.MakePtr("testAnalyzer"),
						},
						ID:     pointer.MakePtr("123"),
						Status: pointer.MakePtr(IndexStatusActive),
					}, nil
				}
				return nil, fmt.Errorf("unexpected")
			},
			CreateIndexFunc: func(_ context.Context, _, _ string, idx *searchindex.SearchIndex) (*searchindex.SearchIndex, error) {
				if idx.Name == "Index1" {
					return &searchindex.SearchIndex{
						SearchIndex: akov2.SearchIndex{
							Name: "Index1",
							Type: IndexTypeSearch,
						},
						ID:     pointer.MakePtr("123"),
						Status: pointer.MakePtr("NOT STARTED"),
					}, nil
				}
				if idx.Name == "Index2" {
					return nil, errors.New("conflict")
				}
				return nil, fmt.Errorf("unexpected")
			},
		}

		searchIndexConfig := &akov2.AtlasSearchIndexConfig{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testConfig",
				Namespace: "testNamespace",
			},
			Spec: akov2.AtlasSearchIndexConfigSpec{
				Analyzer: pointer.MakePtr("testAnalyzer"),
			},
			Status: status.AtlasSearchIndexConfigStatus{},
		}

		deployment := &akov2.AtlasDeployment{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Spec: akov2.AtlasDeploymentSpec{
				DeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name: "testDeployment",
					SearchIndexes: []akov2.SearchIndex{
						{
							Name: "Index1",
							Type: IndexTypeSearch,
							Search: &akov2.Search{
								SearchConfigurationRef: common.ResourceRefNamespaced{
									Name:      searchIndexConfig.Name,
									Namespace: searchIndexConfig.Namespace,
								},
							},
						},
					},
				},
			},
			Status: status.AtlasDeploymentStatus{},
		}

		sch := runtime.NewScheme()
		assert.Nil(t, akov2.AddToScheme(sch))
		assert.Nil(t, corev1.AddToScheme(sch))

		k8sClient := fake.NewClientBuilder().
			WithScheme(sch).
			WithObjects(deployment, searchIndexConfig).
			Build()

		reconciler := searchIndexesReconcileRequest{
			ctx: &workflow.Context{
				Log:     zap.S(),
				OrgID:   "testOrgID",
				Client:  nil,
				Context: context.Background(),
			},
			deployment:    deployment,
			k8sClient:     k8sClient,
			projectID:     "testProjectID",
			searchService: fakeAtlasSearch,
		}

		// first reconcile succeeds, creation succeeds
		result := reconciler.Handle()
		deployment.UpdateStatus(reconciler.ctx.Conditions(), reconciler.ctx.StatusOptions()...)
		assert.False(t, result.IsOk())
		assert.Equal(t, []status.DeploymentSearchIndexStatus{
			{
				Name:    "Index1",
				ID:      "123",
				Status:  "InProgress",
				Message: "Atlas search index status: NOT STARTED",
			},
		}, deployment.Status.SearchIndexes)

		// Add another search Index2 which is a copy of Index1
		deployment.Spec.DeploymentSpec.SearchIndexes = append(deployment.Spec.DeploymentSpec.SearchIndexes, akov2.SearchIndex{
			Name: "Index2",
			Type: IndexTypeSearch,
			Search: &akov2.Search{
				SearchConfigurationRef: common.ResourceRefNamespaced{
					Name:      searchIndexConfig.Name,
					Namespace: searchIndexConfig.Namespace,
				},
			},
		})

		// create fails
		result = reconciler.Handle()
		deployment.UpdateStatus(reconciler.ctx.Conditions(), reconciler.ctx.StatusOptions()...)
		assert.False(t, result.IsOk())
		assert.Equal(t, []status.DeploymentSearchIndexStatus{
			{
				Name:    "Index1",
				ID:      "123",
				Status:  "READY",
				Message: "Atlas search index status: READY",
			},
			{
				Name:    "Index2",
				ID:      "",
				Status:  "Error",
				Message: "error with processing index Index2. err: conflict",
			},
		}, deployment.Status.SearchIndexes)

		// remove Index2 from the spec
		deployment.Spec.DeploymentSpec.SearchIndexes = []akov2.SearchIndex{deployment.Spec.DeploymentSpec.SearchIndexes[0]}

		// third reconcile succeeds, creation succeeds
		result = reconciler.Handle()
		deployment.UpdateStatus(reconciler.ctx.Conditions(), reconciler.ctx.StatusOptions()...)
		assert.True(t, result.IsOk())
		assert.Equal(t, []status.DeploymentSearchIndexStatus{
			{
				Name:    "Index1",
				ID:      "123",
				Status:  "READY",
				Message: "Atlas search index status: READY",
			},
		}, deployment.Status.SearchIndexes)
	})

	t.Run("Should proceed with the index Type Search: UPDATE INDEX", func(t *testing.T) {
		sampleIndex := searchindex.SearchIndex{
			SearchIndex: akov2.SearchIndex{
				Name:           "testName",
				DBName:         "testDB",
				CollectionName: "testCollection",
				Type:           IndexTypeSearch,
			},
			ID:     pointer.MakePtr("testID"),
			Status: pointer.MakePtr(IndexStatusActive),
		}
		fakeAtlasSearch := &searchfake.FakeAtlasSearch{
			GetIndexFunc: func(_ context.Context, _, _, _, _ string) (*searchindex.SearchIndex, error) {
				return &sampleIndex, nil
			},
			UpdateIndexFunc: func(_ context.Context, _, _ string, _ *searchindex.SearchIndex) (*searchindex.SearchIndex, error) {
				return &sampleIndex, nil
			},
		}

		searchIndexConfig := &akov2.AtlasSearchIndexConfig{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testConfig",
				Namespace: "testNamespace",
			},
			Spec: akov2.AtlasSearchIndexConfigSpec{
				Analyzer: pointer.MakePtr("testAnalyzer"),
			},
			Status: status.AtlasSearchIndexConfigStatus{},
		}

		cluster := &akov2.AtlasDeployment{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Spec: akov2.AtlasDeploymentSpec{
				DeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name: "testDeployment",
					SearchIndexes: []akov2.SearchIndex{
						{
							Name: "testName",
							Type: IndexTypeSearch,
							Search: &akov2.Search{
								Synonyms: &([]akov2.Synonym{
									{
										Name:     "testSynonym",
										Analyzer: "testAnalyzer",
										Source: akov2.Source{
											Collection: "testCollection",
										},
									},
								}),
								Mappings: nil,
								SearchConfigurationRef: common.ResourceRefNamespaced{
									Name:      searchIndexConfig.Name,
									Namespace: searchIndexConfig.Namespace,
								},
							},
						},
					},
				},
			},
			Status: status.AtlasDeploymentStatus{
				SearchIndexes: []status.DeploymentSearchIndexStatus{
					{
						Name: "testName",
						ID:   "testID",
					},
				},
			},
		}

		sch := runtime.NewScheme()
		assert.Nil(t, akov2.AddToScheme(sch))
		assert.Nil(t, corev1.AddToScheme(sch))

		k8sClient := fake.NewClientBuilder().
			WithScheme(sch).
			WithObjects(cluster, searchIndexConfig).
			Build()

		reconciler := searchIndexesReconcileRequest{
			ctx: &workflow.Context{
				Log:     zap.S(),
				OrgID:   "testOrgID",
				Client:  nil,
				Context: context.Background(),
			},
			deployment:    cluster,
			k8sClient:     k8sClient,
			projectID:     "testProjectID",
			searchService: fakeAtlasSearch,
		}
		result := reconciler.Handle()
		fmt.Println("Result", result)
		assert.True(t, reconciler.ctx.HasReason(api.SearchIndexesNotReady))
		assert.True(t, result.IsInProgress())
	})

	t.Run("Should proceed with the index Type Search: UPDATE INDEX (vector)", func(t *testing.T) {
		sampleIndex := searchindex.SearchIndex{
			SearchIndex: akov2.SearchIndex{
				Name:           "testName",
				DBName:         "testDB",
				CollectionName: "testCollection",
				Type:           IndexTypeVector,
			},
			ID:     pointer.MakePtr("testID"),
			Status: pointer.MakePtr(IndexStatusActive),
		}
		fakeAtlasSearch := &searchfake.FakeAtlasSearch{
			GetIndexFunc: func(_ context.Context, _, _, _, _ string) (*searchindex.SearchIndex, error) {
				return &sampleIndex, nil
			},
			UpdateIndexFunc: func(_ context.Context, _, _ string, _ *searchindex.SearchIndex) (*searchindex.SearchIndex, error) {
				return &sampleIndex, nil
			},
		}

		searchIndexConfig := &akov2.AtlasSearchIndexConfig{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testConfig",
				Namespace: "testNamespace",
			},
			Spec: akov2.AtlasSearchIndexConfigSpec{
				Analyzer: pointer.MakePtr("testAnalyzer"),
			},
			Status: status.AtlasSearchIndexConfigStatus{},
		}

		cluster := &akov2.AtlasDeployment{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Spec: akov2.AtlasDeploymentSpec{
				DeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name: "testDeployment",
					SearchIndexes: []akov2.SearchIndex{
						{
							Name: "testName",
							Type: IndexTypeVector,
							Search: &akov2.Search{
								Synonyms: &([]akov2.Synonym{
									{
										Name:     "testSynonym",
										Analyzer: "testAnalyzer",
										Source: akov2.Source{
											Collection: "testCollection",
										},
									},
								}),
								Mappings: nil,
								SearchConfigurationRef: common.ResourceRefNamespaced{
									Name:      searchIndexConfig.Name,
									Namespace: searchIndexConfig.Namespace,
								},
							},
						},
					},
				},
			},
			Status: status.AtlasDeploymentStatus{
				SearchIndexes: []status.DeploymentSearchIndexStatus{
					{
						Name: "testName",
						ID:   "testID",
					},
				},
			},
		}

		sch := runtime.NewScheme()
		assert.Nil(t, akov2.AddToScheme(sch))
		assert.Nil(t, corev1.AddToScheme(sch))

		k8sClient := fake.NewClientBuilder().
			WithScheme(sch).
			WithObjects(cluster, searchIndexConfig).
			Build()

		reconciler := searchIndexesReconcileRequest{
			ctx: &workflow.Context{
				Log:     zap.S(),
				OrgID:   "testOrgID",
				Client:  nil,
				Context: context.Background(),
			},
			deployment:    cluster,
			k8sClient:     k8sClient,
			projectID:     "testProjectID",
			searchService: fakeAtlasSearch,
		}
		result := reconciler.Handle()
		fmt.Println("Result", result)
		assert.True(t, reconciler.ctx.HasReason(api.SearchIndexesNotReady))
		assert.True(t, result.IsInProgress())
	})

	t.Run("Should proceed with the index Type Search: DELETE INDEX", func(t *testing.T) {
		fakeAtlasSearch := &searchfake.FakeAtlasSearch{
			GetIndexFunc: func(_ context.Context, _, _, _, _ string) (*searchindex.SearchIndex, error) {
				return &searchindex.SearchIndex{
					SearchIndex: akov2.SearchIndex{
						Name:           "testName",
						DBName:         "testDB",
						CollectionName: "testCollection",
						Type:           IndexTypeVector,
					},
					ID:     pointer.MakePtr("testID"),
					Status: pointer.MakePtr(IndexStatusActive),
				}, nil
			},
			DeleteIndexFunc: func(_ context.Context, _, _, _ string) error {
				return nil
			},
		}

		// mockSearchAPI := mockadmin.NewAtlasSearchApi(t)

		// mockSearchAPI.EXPECT().
		// 	GetAtlasSearchIndex(context.Background(), mock.Anything, mock.Anything, mock.Anything).
		// 	Return(admin.GetAtlasSearchIndexApiRequest{ApiService: mockSearchAPI})
		// mockSearchAPI.EXPECT().
		// 	GetAtlasSearchIndexExecute(admin.GetAtlasSearchIndexApiRequest{ApiService: mockSearchAPI}).
		// 	Return(
		// 		&admin.ClusterSearchIndex{
		// 			CollectionName: "testCollection",
		// 			Database:       "testDB",
		// 			IndexID:        pointer.MakePtr("testID"),
		// 			Name:           "testName",
		// 			Status:         pointer.MakePtr(IndexStatusActive),
		// 			Type:           pointer.MakePtr(IndexTypeVector),
		// 		},
		// 		&http.Response{StatusCode: http.StatusOK}, nil,
		// 	)

		// mockSearchAPI.EXPECT().
		// 	DeleteAtlasSearchIndex(context.Background(), mock.Anything, mock.Anything, mock.Anything).
		// 	Return(admin.DeleteAtlasSearchIndexApiRequest{ApiService: mockSearchAPI})

		// mockSearchAPI.EXPECT().
		// 	DeleteAtlasSearchIndexExecute(admin.DeleteAtlasSearchIndexApiRequest{ApiService: mockSearchAPI}).
		// 	Return(nil, &http.Response{StatusCode: http.StatusAccepted}, nil)

		searchIndexConfig := &akov2.AtlasSearchIndexConfig{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testConfig",
				Namespace: "testNamespace",
			},
			Spec: akov2.AtlasSearchIndexConfigSpec{
				Analyzer: pointer.MakePtr("testAnalyzer"),
			},
			Status: status.AtlasSearchIndexConfigStatus{},
		}

		cluster := &akov2.AtlasDeployment{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Spec: akov2.AtlasDeploymentSpec{
				DeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:          "testDeployment",
					SearchIndexes: []akov2.SearchIndex{},
				},
			},
			Status: status.AtlasDeploymentStatus{
				SearchIndexes: []status.DeploymentSearchIndexStatus{
					{
						Name: "testName",
						ID:   "testID",
					},
				},
			},
		}

		sch := runtime.NewScheme()
		assert.Nil(t, akov2.AddToScheme(sch))
		assert.Nil(t, corev1.AddToScheme(sch))

		k8sClient := fake.NewClientBuilder().
			WithScheme(sch).
			WithObjects(cluster, searchIndexConfig).
			Build()

		reconciler := searchIndexesReconcileRequest{
			ctx: &workflow.Context{
				Log:     zap.S(),
				OrgID:   "testOrgID",
				Client:  nil,
				Context: context.Background(),
			},
			deployment:    cluster,
			k8sClient:     k8sClient,
			projectID:     "testProjectID",
			searchService: fakeAtlasSearch,
		}
		result := reconciler.Handle()
		cluster.UpdateStatus(reconciler.ctx.Conditions(), reconciler.ctx.StatusOptions()...)
		assert.Empty(t, cluster.Status.SearchIndexes)
		assert.True(t, result.IsOk())
	})
}
