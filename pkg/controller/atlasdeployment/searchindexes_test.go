package atlasdeployment

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	internal "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/searchindex"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
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
				Log:   zap.S(),
				OrgID: "testOrgID",
			},
			deployment: deployment,
		}
		result := reconciler.Reconcile()
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

		atlasIdxToReturn, err := internal.NewSearchIndexFromAKO(&deployment.Spec.DeploymentSpec.SearchIndexes[0],
			&searchIndexConfig.Spec).ToAtlas()
		assert.NoError(t, err)
		mockSearchAPI := mockadmin.NewAtlasSearchApi(t)
		mockSearchAPI.EXPECT().
			GetAtlasSearchIndex(context.Background(), mock.Anything, mock.Anything, mock.Anything).
			Return(admin.GetAtlasSearchIndexApiRequest{ApiService: mockSearchAPI})
		mockSearchAPI.EXPECT().
			GetAtlasSearchIndexExecute(admin.GetAtlasSearchIndexApiRequest{ApiService: mockSearchAPI}).
			Return(
				atlasIdxToReturn,
				&http.Response{StatusCode: http.StatusCreated}, nil,
			)

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
		assert.True(t, result.IsOk())
		deployment.UpdateStatus(reconciler.ctx.Conditions(), reconciler.ctx.StatusOptions()...)
		assert.Len(t, deployment.Status.SearchIndexes, 1)
		assert.True(t, deployment.Status.SearchIndexes[0].ID == IDForStatus)
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

	t.Run("Should proceed with index Type Search if it cannot be found: CREATE INDEX", func(t *testing.T) {
		mockSearchAPI := mockadmin.NewAtlasSearchApi(t)

		nestMock := func(f func(nestedMock *mockadmin.AtlasSearchApi)) {
			nested := mockadmin.NewAtlasSearchApi(t)
			f(nested)
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

		// mock GET for Index1
		nestMock(func(n *mockadmin.AtlasSearchApi) {
			mockSearchAPI.EXPECT().
				GetAtlasSearchIndex(context.Background(), mock.Anything, mock.Anything, "123").
				Return(admin.GetAtlasSearchIndexApiRequest{ApiService: n})
			n.EXPECT().
				GetAtlasSearchIndexExecute(admin.GetAtlasSearchIndexApiRequest{ApiService: n}).
				Return(
					&admin.ClusterSearchIndex{
						IndexID:  pointer.MakePtr("123"),
						Name:     "Index1",
						Status:   pointer.MakePtr(IndexStatusActive),
						Type:     pointer.MakePtr(IndexTypeSearch),
						Analyzer: pointer.MakePtr("testAnalyzer"),
					},
					&http.Response{StatusCode: http.StatusOK}, nil,
				)
		})

		// mock CREATE for Index1
		atlasIdx, err := internal.NewSearchIndexFromAKO(&deployment.Spec.DeploymentSpec.SearchIndexes[0], &searchIndexConfig.Spec).ToAtlas()
		assert.NoError(t, err)

		nestMock(func(n *mockadmin.AtlasSearchApi) {
			mockSearchAPI.EXPECT().
				CreateAtlasSearchIndex(context.Background(), mock.Anything, mock.Anything, atlasIdx).
				Return(admin.CreateAtlasSearchIndexApiRequest{ApiService: n})
			n.EXPECT().
				CreateAtlasSearchIndexExecute(admin.CreateAtlasSearchIndexApiRequest{ApiService: n}).
				Return(
					&admin.ClusterSearchIndex{Name: "Index1", IndexID: pointer.MakePtr("123"), Status: pointer.MakePtr("NOT STARTED")},
					&http.Response{StatusCode: http.StatusCreated}, nil,
				)
		})

		// first reconcile succeeds, creation succeeds
		result := reconciler.Reconcile()
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

		// mock CREATE for Index2
		atlasIdx, err = internal.NewSearchIndexFromAKO(&deployment.Spec.DeploymentSpec.SearchIndexes[1], &searchIndexConfig.Spec).ToAtlas()
		assert.NoError(t, err)

		nestMock(func(n *mockadmin.AtlasSearchApi) {
			mockSearchAPI.EXPECT().CreateAtlasSearchIndex(context.Background(), mock.Anything, mock.Anything, atlasIdx).
				Return(admin.CreateAtlasSearchIndexApiRequest{ApiService: n})
			n.EXPECT().
				CreateAtlasSearchIndexExecute(admin.CreateAtlasSearchIndexApiRequest{ApiService: n}).
				Return(
					nil,
					&http.Response{StatusCode: http.StatusBadRequest}, errors.New("conflict"),
				)
		})

		// create fails
		result = reconciler.Reconcile()
		deployment.UpdateStatus(reconciler.ctx.Conditions(), reconciler.ctx.StatusOptions()...)
		assert.False(t, result.IsOk())
		assert.Equal(t, []status.DeploymentSearchIndexStatus{
			{
				Name:    "Index1",
				ID:      "123",
				Status:  "Ready",
				Message: "Atlas search index status: STEADY",
			},
			{
				Name:    "Index2",
				ID:      "",
				Status:  "Error",
				Message: "error with processing index Index2. err: failed to create index: conflict, status: 400",
			},
		}, deployment.Status.SearchIndexes)

		// remove Index2 from the spec
		deployment.Spec.DeploymentSpec.SearchIndexes = []akov2.SearchIndex{deployment.Spec.DeploymentSpec.SearchIndexes[0]}

		// third reconcile succeeds, creation succeeds
		result = reconciler.Reconcile()
		deployment.UpdateStatus(reconciler.ctx.Conditions(), reconciler.ctx.StatusOptions()...)
		assert.True(t, result.IsOk())
		assert.Equal(t, []status.DeploymentSearchIndexStatus{
			{
				Name:    "Index1",
				ID:      "123",
				Status:  "Ready",
				Message: "Atlas search index status: STEADY",
			},
		}, deployment.Status.SearchIndexes)
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
		assert.True(t, reconciler.ctx.HasReason(api.SearchIndexesNotReady))
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
		assert.True(t, reconciler.ctx.HasReason(api.SearchIndexesNotReady))
		assert.True(t, result.IsInProgress())
	})

	t.Run("Should proceed with the index Type Search: DELETE INDEX", func(t *testing.T) {
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
				},
				&http.Response{StatusCode: http.StatusOK}, nil,
			)

		mockSearchAPI.EXPECT().
			DeleteAtlasSearchIndex(context.Background(), mock.Anything, mock.Anything, mock.Anything).
			Return(admin.DeleteAtlasSearchIndexApiRequest{ApiService: mockSearchAPI})

		mockSearchAPI.EXPECT().
			DeleteAtlasSearchIndexExecute(admin.DeleteAtlasSearchIndexApiRequest{ApiService: mockSearchAPI}).
			Return(nil, &http.Response{StatusCode: http.StatusAccepted}, nil)

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
		cluster.UpdateStatus(reconciler.ctx.Conditions(), reconciler.ctx.StatusOptions()...)
		assert.Empty(t, cluster.Status.SearchIndexes)
		assert.True(t, result.IsOk())
	})
}
