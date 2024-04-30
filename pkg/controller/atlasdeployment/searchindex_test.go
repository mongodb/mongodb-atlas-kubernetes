package atlasdeployment

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/searchindex"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func Test_searchIndexReconciler(t *testing.T) {
	t.Run("Must reconcile index to create", func(t *testing.T) {
		sch := runtime.NewScheme()
		assert.Nil(t, akov2.AddToScheme(sch))
		assert.Nil(t, corev1.AddToScheme(sch))

		indexToTest := &searchindex.SearchIndex{
			SearchIndex: akov2.SearchIndex{
				Name:           "test",
				DBName:         "testDB",
				CollectionName: "testCollection",
				Type:           "search",
				Search: &akov2.Search{
					Synonyms: nil,
					Mappings: nil,
				},
				VectorSearch: nil,
			},
			AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{
				Analyzer:       pointer.MakePtr("testAnalyzer"),
				Analyzers:      nil,
				SearchAnalyzer: nil,
				StoredSource:   nil,
			},
			ID:     nil,
			Status: nil,
		}

		testCluster := &akov2.AtlasDeployment{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testDeployment",
				Namespace: "testNamespace",
			},
			Spec: akov2.AtlasDeploymentSpec{
				DeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name: "testDeploymentName",
				},
			},
			Status: status.AtlasDeploymentStatus{},
		}

		mockSearchAPI := mockadmin.NewAtlasSearchApi(t)
		mockSearchAPI.EXPECT().
			CreateAtlasSearchIndex(context.Background(), mock.Anything, mock.Anything, mock.Anything).
			Return(admin.CreateAtlasSearchIndexApiRequest{ApiService: mockSearchAPI})
		mockSearchAPI.EXPECT().
			CreateAtlasSearchIndexExecute(admin.CreateAtlasSearchIndexApiRequest{ApiService: mockSearchAPI}).
			Return(&admin.ClusterSearchIndex{
				CollectionName: "",
				Database:       "",
				IndexID:        pointer.MakePtr("testID"),
				Name:           "",
				Status:         pointer.MakePtr("NOT STARTED"),
				Type:           nil,
				Analyzer:       nil,
				Analyzers:      nil,
				Mappings:       nil,
				SearchAnalyzer: nil,
				StoredSource:   nil,
				Synonyms:       nil,
				Fields:         nil,
			}, &http.Response{StatusCode: http.StatusCreated}, nil)

		reconciler := &searchIndexReconciler{
			ctx: &workflow.Context{
				Log:   zap.S(),
				OrgID: "testOrgID",
				SdkClient: &admin.APIClient{
					AtlasSearchApi: mockSearchAPI,
				},
				Context: context.Background(),
			},
			deployment: testCluster,
			k8sClient:  nil,
			projectID:  "",
			indexName:  "testIndexName",
		}

		result := reconciler.Reconcile(indexToTest, nil, nil)
		assert.True(t, result.IsInProgress())
		fmt.Println(result)
		fmt.Println(testCluster.Status)
	})

	t.Run("Must reconcile index to delete", func(t *testing.T) {
		sch := runtime.NewScheme()
		assert.Nil(t, akov2.AddToScheme(sch))
		assert.Nil(t, corev1.AddToScheme(sch))

		indexToTest := &searchindex.SearchIndex{
			SearchIndex: akov2.SearchIndex{
				Name:           "test",
				DBName:         "testDB",
				CollectionName: "testCollection",
				Type:           "search",
				Search: &akov2.Search{
					Synonyms: nil,
					Mappings: nil,
				},
				VectorSearch: nil,
			},
			AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{
				Analyzer:       pointer.MakePtr("testAnalyzer"),
				Analyzers:      nil,
				SearchAnalyzer: nil,
				StoredSource:   nil,
			},
			ID:     pointer.MakePtr("testID"),
			Status: nil,
		}

		testCluster := &akov2.AtlasDeployment{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testDeployment",
				Namespace: "testNamespace",
			},
			Spec: akov2.AtlasDeploymentSpec{
				DeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name: "testDeploymentName",
				},
			},
			Status: status.AtlasDeploymentStatus{},
		}

		mockSearchAPI := mockadmin.NewAtlasSearchApi(t)
		mockSearchAPI.EXPECT().
			DeleteAtlasSearchIndex(context.Background(), mock.Anything, mock.Anything, mock.Anything).
			Return(admin.DeleteAtlasSearchIndexApiRequest{ApiService: mockSearchAPI})
		mockSearchAPI.EXPECT().
			DeleteAtlasSearchIndexExecute(admin.DeleteAtlasSearchIndexApiRequest{ApiService: mockSearchAPI}).
			Return(map[string]interface{}{}, &http.Response{StatusCode: http.StatusAccepted}, nil)

		reconciler := &searchIndexReconciler{
			ctx: &workflow.Context{
				Log:   zap.S(),
				OrgID: "testOrgID",
				SdkClient: &admin.APIClient{
					AtlasSearchApi: mockSearchAPI,
				},
				Context: context.Background(),
			},
			deployment: testCluster,
			k8sClient:  nil,
			projectID:  "",
			indexName:  "testIndexName",
		}

		result := reconciler.Reconcile(nil, indexToTest, nil)
		assert.True(t, result.IsOk())
		fmt.Println(result)
		fmt.Println(testCluster.Status)
	})

	t.Run("Must not reconcile if there are errors for the index", func(t *testing.T) {
		reconciler := &searchIndexReconciler{
			ctx: &workflow.Context{
				Log:       zap.S(),
				OrgID:     "testOrgID",
				SdkClient: &admin.APIClient{},
				Context:   context.Background(),
			},
			deployment: nil,
			k8sClient:  nil,
			projectID:  "",
			indexName:  "testIndexName",
		}

		result := reconciler.Reconcile(nil, nil, []error{fmt.Errorf("testError")})
		assert.False(t, result.IsOk())
	})

	t.Run("Must return InProgress if index status is anything but ACTIVE", func(t *testing.T) {
		reconciler := &searchIndexReconciler{
			ctx: &workflow.Context{
				Log:       zap.S(),
				OrgID:     "testOrgID",
				SdkClient: &admin.APIClient{},
				Context:   context.Background(),
			},
			deployment: nil,
			k8sClient:  nil,
			projectID:  "",
			indexName:  "testIndexName",
		}
		result := reconciler.Reconcile(nil, &searchindex.SearchIndex{Status: pointer.MakePtr("NOT STARTED")}, nil)
		assert.True(t, result.IsInProgress())
	})

	t.Run("Must not call update API if indexes are equal", func(t *testing.T) {
		reconciler := &searchIndexReconciler{
			ctx: &workflow.Context{
				Log:       zap.S(),
				OrgID:     "testOrgID",
				SdkClient: &admin.APIClient{},
				Context:   context.Background(),
			},
			deployment: nil,
			k8sClient:  nil,
			projectID:  "",
			indexName:  "testIndexName",
		}
		idx := &searchindex.SearchIndex{
			SearchIndex:                akov2.SearchIndex{},
			AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{},
			ID:                         nil,
			Status:                     nil,
		}
		result := reconciler.Reconcile(idx, idx, nil)
		assert.True(t, result.IsOk())
	})

	t.Run("Must trigger index update if state in AKO and in Atlas is different", func(t *testing.T) {
		mockSearchAPI := mockadmin.NewAtlasSearchApi(t)
		mockSearchAPI.EXPECT().
			UpdateAtlasSearchIndex(context.Background(), mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.UpdateAtlasSearchIndexApiRequest{ApiService: mockSearchAPI})
		mockSearchAPI.EXPECT().
			UpdateAtlasSearchIndexExecute(admin.UpdateAtlasSearchIndexApiRequest{ApiService: mockSearchAPI}).
			Return(
				&admin.ClusterSearchIndex{Status: pointer.MakePtr("NOT STARTED")},
				&http.Response{StatusCode: http.StatusCreated}, nil,
			)

		testCluster := &akov2.AtlasDeployment{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testDeployment",
				Namespace: "testNamespace",
			},
			Spec: akov2.AtlasDeploymentSpec{
				DeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name: "testDeploymentName",
				},
			},
			Status: status.AtlasDeploymentStatus{},
		}

		reconciler := &searchIndexReconciler{
			ctx: &workflow.Context{
				Log:   zap.S(),
				OrgID: "testOrgID",
				SdkClient: &admin.APIClient{
					AtlasSearchApi: mockSearchAPI,
				},
				Context: context.Background(),
			},
			deployment: testCluster,
			k8sClient:  nil,
			projectID:  "",
			indexName:  "testIndexName",
		}
		idxInAtlas := &searchindex.SearchIndex{
			SearchIndex: akov2.SearchIndex{
				Name: "testIndex",
			},
			AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{},
			ID:                         pointer.MakePtr("testID"),
			Status:                     nil,
		}
		idxInAKO := &searchindex.SearchIndex{
			SearchIndex: akov2.SearchIndex{
				Name: "testIndex",
				Search: &akov2.Search{
					Synonyms: nil,
					Mappings: &akov2.Mappings{
						Dynamic: pointer.MakePtr(true),
					},
				},
			},
			AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{},
			ID:                         pointer.MakePtr("testID"),
			Status:                     nil,
		}
		result := reconciler.Reconcile(idxInAKO, idxInAtlas, nil)
		assert.True(t, result.IsInProgress())
	})
}
