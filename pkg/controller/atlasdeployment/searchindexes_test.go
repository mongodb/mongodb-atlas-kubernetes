package atlasdeployment

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"

	"github.com/stretchr/testify/assert"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
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

		reconciler := searchIndexesReconciler{
			ctx: &workflow.Context{
				Log:       zap.S(),
				OrgID:     "testOrgID",
				Client:    nil,
				SdkClient: nil,
				Context:   nil,
			},
			deployment: deployment,
			k8sClient:  nil,
			projectID:  "",
		}
		result := reconciler.Reconcile()
		assert.True(t, reconciler.ctx.HasReason(status.SearchIndexesNamesAreNotUnique))
		assert.False(t, result.IsOk())
	})

	t.Run("Should proceed with the index Type Search: CREATE INDEX", func(t *testing.T) {
		mockSearchAPI := mockadmin.NewAtlasSearchApi(t)
		mockSearchAPI.EXPECT().
			CreateAtlasSearchIndex(context.Background(), mock.Anything, mock.Anything, mock.Anything).
			Return(admin.CreateAtlasSearchIndexApiRequest{ApiService: mockSearchAPI})
		mockSearchAPI.EXPECT().
			CreateAtlasSearchIndexExecute(admin.CreateAtlasSearchIndexApiRequest{ApiService: mockSearchAPI}).
			Return(
				&admin.ClusterSearchIndex{Status: pointer.MakePtr("NOT STARTED")},
				&http.Response{StatusCode: http.StatusCreated}, nil,
			)

		searchIndexConfig := &akov2.AtlasSearchIndexConfig{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testConfig",
				Namespace: "testNamespace",
			},
			Spec: akov2.AtlasSearchIndexConfigSpec{
				Analyzer:       pointer.MakePtr("testAnalyzer"),
				Analyzers:      nil,
				SearchAnalyzer: nil,
				StoredSource:   nil,
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

		reconciler := searchIndexesReconciler{
			ctx: &workflow.Context{
				Log:       zap.S(),
				OrgID:     "testOrgID",
				Client:    nil,
				SdkClient: &admin.APIClient{AtlasSearchApi: mockSearchAPI},
				Context:   context.Background(),
			},
			deployment: deployment,
			k8sClient:  k8sClient,
			projectID:  "testProjectID",
		}
		result := reconciler.Reconcile()
		assert.False(t, result.IsOk())
	})

	t.Run("Should proceed with the index Type Search: UPDATE INDEX", func(t *testing.T) {
		mockSearchAPI := mockadmin.NewAtlasSearchApi(t)

		mockSearchAPI.EXPECT().
			GetAtlasSearchIndex(context.Background(), mock.Anything, mock.Anything, mock.Anything).
			Return(admin.GetAtlasSearchIndexApiRequest{ApiService: mockSearchAPI})
		mockSearchAPI.EXPECT().
			GetAtlasSearchIndexExecute(admin.GetAtlasSearchIndexApiRequest{ApiService: mockSearchAPI}).
			Return(
				&admin.ClusterSearchIndex{
					CollectionName: "testCollection",
					Database:       "testDB",
					IndexID:        pointer.MakePtr("testID"),
					Name:           "testName",
					Status:         pointer.MakePtr(IndexStatusActive),
					Type:           pointer.MakePtr(IndexTypeSearch),
					Analyzer:       nil,
					Analyzers:      nil,
					Mappings:       nil,
					SearchAnalyzer: nil,
					StoredSource:   nil,
					Synonyms:       nil,
					Fields:         nil,
				},
				&http.Response{StatusCode: http.StatusOK}, nil,
			)

		mockSearchAPI.EXPECT().
			UpdateAtlasSearchIndex(context.Background(), mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.UpdateAtlasSearchIndexApiRequest{ApiService: mockSearchAPI})
		mockSearchAPI.EXPECT().
			UpdateAtlasSearchIndexExecute(admin.UpdateAtlasSearchIndexApiRequest{ApiService: mockSearchAPI}).
			Return(
				&admin.ClusterSearchIndex{
					CollectionName: "testCollection",
					Database:       "testDB",
					IndexID:        pointer.MakePtr("testID"),
					Name:           "testName",
					Status:         pointer.MakePtr(IndexStatusActive),
					Type:           pointer.MakePtr(IndexTypeSearch),
					Analyzer:       nil,
					Analyzers:      nil,
					Mappings:       nil,
					SearchAnalyzer: nil,
					StoredSource:   nil,
					Synonyms:       nil,
					Fields:         nil,
				},
				&http.Response{StatusCode: http.StatusOK}, nil,
			)

		searchIndexConfig := &akov2.AtlasSearchIndexConfig{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testConfig",
				Namespace: "testNamespace",
			},
			Spec: akov2.AtlasSearchIndexConfigSpec{
				Analyzer:       pointer.MakePtr("testAnalyzer"),
				Analyzers:      nil,
				SearchAnalyzer: nil,
				StoredSource:   nil,
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

		reconciler := searchIndexesReconciler{
			ctx: &workflow.Context{
				Log:       zap.S(),
				OrgID:     "testOrgID",
				Client:    nil,
				SdkClient: &admin.APIClient{AtlasSearchApi: mockSearchAPI},
				Context:   context.Background(),
			},
			deployment: cluster,
			k8sClient:  k8sClient,
			projectID:  "testProjectID",
		}
		result := reconciler.Reconcile()
		fmt.Println("Result", result)
		assert.True(t, reconciler.ctx.HasReason(status.SearchIndexesNotReady))
		assert.True(t, result.IsInProgress())
	})

	t.Run("Should proceed with the index Type Search: UPDATE INDEX (vector)", func(t *testing.T) {
		mockSearchAPI := mockadmin.NewAtlasSearchApi(t)

		mockSearchAPI.EXPECT().
			GetAtlasSearchIndex(context.Background(), mock.Anything, mock.Anything, mock.Anything).
			Return(admin.GetAtlasSearchIndexApiRequest{ApiService: mockSearchAPI})
		mockSearchAPI.EXPECT().
			GetAtlasSearchIndexExecute(admin.GetAtlasSearchIndexApiRequest{ApiService: mockSearchAPI}).
			Return(
				&admin.ClusterSearchIndex{
					CollectionName: "testCollection",
					Database:       "testDB",
					IndexID:        pointer.MakePtr("testID"),
					Name:           "testName",
					Status:         pointer.MakePtr(IndexStatusActive),
					Type:           pointer.MakePtr(IndexTypeVector),
					Analyzer:       nil,
					Analyzers:      nil,
					Mappings:       nil,
					SearchAnalyzer: nil,
					StoredSource:   nil,
					Synonyms:       nil,
					Fields:         nil,
				},
				&http.Response{StatusCode: http.StatusOK}, nil,
			)

		mockSearchAPI.EXPECT().
			UpdateAtlasSearchIndex(context.Background(), mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.UpdateAtlasSearchIndexApiRequest{ApiService: mockSearchAPI})
		mockSearchAPI.EXPECT().
			UpdateAtlasSearchIndexExecute(admin.UpdateAtlasSearchIndexApiRequest{ApiService: mockSearchAPI}).
			Return(
				&admin.ClusterSearchIndex{
					CollectionName: "testCollection",
					Database:       "testDB",
					IndexID:        pointer.MakePtr("testID"),
					Name:           "testName",
					Status:         pointer.MakePtr("ACTIVE"),
					Type:           pointer.MakePtr(IndexTypeVector),
					Analyzer:       nil,
					Analyzers:      nil,
					Mappings:       nil,
					SearchAnalyzer: nil,
					StoredSource:   nil,
					Synonyms:       nil,
					Fields:         nil,
				},
				&http.Response{StatusCode: http.StatusOK}, nil,
			)

		searchIndexConfig := &akov2.AtlasSearchIndexConfig{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testConfig",
				Namespace: "testNamespace",
			},
			Spec: akov2.AtlasSearchIndexConfigSpec{
				Analyzer:       pointer.MakePtr("testAnalyzer"),
				Analyzers:      nil,
				SearchAnalyzer: nil,
				StoredSource:   nil,
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

		reconciler := searchIndexesReconciler{
			ctx: &workflow.Context{
				Log:       zap.S(),
				OrgID:     "testOrgID",
				Client:    nil,
				SdkClient: &admin.APIClient{AtlasSearchApi: mockSearchAPI},
				Context:   context.Background(),
			},
			deployment: cluster,
			k8sClient:  k8sClient,
			projectID:  "testProjectID",
		}
		result := reconciler.Reconcile()
		fmt.Println("Result", result)
		assert.True(t, reconciler.ctx.HasReason(status.SearchIndexesNotReady))
		assert.True(t, result.IsInProgress())
	})
}
