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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20250312010/admin"
	"go.mongodb.org/atlas-sdk/v20250312010/mockadmin"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/searchindex"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/searchindex/fake"
)

func Test_searchIndexReconcileRequest(t *testing.T) {
	t.Run("create: must reconcile index to create", func(t *testing.T) {
		sch := runtime.NewScheme()
		assert.Nil(t, akov2.AddToScheme(sch))
		assert.Nil(t, corev1.AddToScheme(sch))

		indexToTest := &searchindex.SearchIndex{
			SearchIndex: akov2.SearchIndex{
				Name:           "test",
				DBName:         "testDB",
				CollectionName: "testCollection",
				Type:           "search",
			},
			AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{
				Analyzer: pointer.MakePtr("testAnalyzer"),
			},
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

		fakeAtlasSearch := &fake.FakeAtlasSearch{
			CreateIndexFunc: func(_ context.Context, _, _ string, _ *searchindex.SearchIndex) (*searchindex.SearchIndex, error) {
				return &searchindex.SearchIndex{
					ID:     pointer.MakePtr("testID"),
					Status: pointer.MakePtr("NOT STARTED"),
				}, nil
			},
		}

		reconciler := &searchIndexReconcileRequest{
			ctx: &workflow.Context{
				Log:     zap.S(),
				OrgID:   "testOrgID",
				Context: context.Background(),
			},
			deployment:    testCluster,
			projectID:     "",
			indexName:     "testIndexName",
			searchService: fakeAtlasSearch,
		}

		result := reconciler.reconcileInternal("", indexToTest, nil)
		assert.True(t, result.IsInProgress())
		fmt.Println(result)
		fmt.Println(testCluster.Status)
	})

	t.Run("create: must return an error if API call returns anything but StatusCreated", func(t *testing.T) {
		sch := runtime.NewScheme()
		assert.Nil(t, akov2.AddToScheme(sch))
		assert.Nil(t, corev1.AddToScheme(sch))

		indexToTest := &searchindex.SearchIndex{
			SearchIndex: akov2.SearchIndex{
				Name:           "test",
				DBName:         "testDB",
				CollectionName: "testCollection",
				Type:           "search",
			},
			AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{
				Analyzer: pointer.MakePtr("testAnalyzer"),
			},
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

		fakeAtlasSearch := &fake.FakeAtlasSearch{
			CreateIndexFunc: func(_ context.Context, _, _ string, _ *searchindex.SearchIndex) (*searchindex.SearchIndex, error) {
				return &searchindex.SearchIndex{
					ID:     pointer.MakePtr("testID"),
					Status: pointer.MakePtr("NOT STARTED"),
				}, errors.New(http.StatusText(http.StatusInternalServerError))
			},
		}

		reconciler := &searchIndexReconcileRequest{
			ctx: &workflow.Context{
				Log:     zap.S(),
				OrgID:   "testOrgID",
				Context: context.Background(),
			},
			deployment:    testCluster,
			projectID:     "",
			indexName:     "testIndexName",
			searchService: fakeAtlasSearch,
		}

		result := reconciler.reconcileInternal("", indexToTest, nil)
		assert.False(t, result.IsOk())
		assert.True(t, reconciler.ctx.HasReason(status.SearchIndexStatusError))
	})

	t.Run("create: must return an error if API call returns an empty index in response", func(t *testing.T) {
		sch := runtime.NewScheme()
		assert.Nil(t, akov2.AddToScheme(sch))
		assert.Nil(t, corev1.AddToScheme(sch))

		indexToTest := &searchindex.SearchIndex{
			SearchIndex: akov2.SearchIndex{
				Name:           "test",
				DBName:         "testDB",
				CollectionName: "testCollection",
				Type:           "search",
			},
			AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{
				Analyzer: pointer.MakePtr("testAnalyzer"),
			},
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
			Return(nil, &http.Response{StatusCode: http.StatusOK}, nil)
		atlasSearch := searchindex.NewSearchIndexes(mockSearchAPI)

		reconciler := &searchIndexReconcileRequest{
			ctx: &workflow.Context{
				Log:     zap.S(),
				OrgID:   "testOrgID",
				Context: context.Background(),
			},
			deployment:    testCluster,
			projectID:     "",
			indexName:     "testIndexName",
			searchService: atlasSearch,
		}

		result := reconciler.reconcileInternal("", indexToTest, nil)
		assert.False(t, result.IsOk())
		assert.True(t, reconciler.ctx.HasReason(status.SearchIndexStatusError))
	})

	t.Run("create: must return an error if index can not be converted internal index", func(t *testing.T) {
		sch := runtime.NewScheme()
		assert.Nil(t, akov2.AddToScheme(sch))
		assert.Nil(t, corev1.AddToScheme(sch))

		indexToTest := &searchindex.SearchIndex{
			SearchIndex: akov2.SearchIndex{
				Name:           "test",
				DBName:         "testDB",
				CollectionName: "testCollection",
				Type:           "search",
			},
			AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{
				StoredSource: &apiextensions.JSON{Raw: []byte{'i', 'n', 'v', 'a', 'l', 'i', 'd', 'j', 's', 'o', 'n'}},
				Analyzer:     pointer.MakePtr("testAnalyzer"),
			},
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
		atlasSearch := searchindex.NewSearchIndexes(mockSearchAPI)

		reconciler := &searchIndexReconcileRequest{
			ctx: &workflow.Context{
				Log:     zap.S(),
				OrgID:   "testOrgID",
				Context: context.Background(),
			},
			deployment:    testCluster,
			projectID:     "",
			indexName:     "testIndexName",
			searchService: atlasSearch,
		}

		result := reconciler.reconcileInternal("", indexToTest, nil)
		assert.False(t, result.IsOk())
		assert.True(t, reconciler.ctx.HasReason(status.SearchIndexStatusError))
	})

	t.Run("delete: must reconcile index to delete", func(t *testing.T) {
		sch := runtime.NewScheme()
		assert.Nil(t, akov2.AddToScheme(sch))
		assert.Nil(t, corev1.AddToScheme(sch))

		indexToTest := &searchindex.SearchIndex{
			SearchIndex: akov2.SearchIndex{
				Name:           "test",
				DBName:         "testDB",
				CollectionName: "testCollection",
				Type:           "search",
			},
			AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{
				Analyzer: pointer.MakePtr("testAnalyzer"),
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

		fakeAtlasSearch := &fake.FakeAtlasSearch{
			DeleteIndexFunc: func(_ context.Context, _, _, _ string) error {
				return nil
			},
		}

		reconciler := &searchIndexReconcileRequest{
			ctx: &workflow.Context{
				Log:     zap.S(),
				OrgID:   "testOrgID",
				Context: context.Background(),
			},
			deployment:    testCluster,
			k8sClient:     nil,
			projectID:     "",
			indexName:     "testIndexName",
			searchService: fakeAtlasSearch,
		}

		result := reconciler.reconcileInternal("", nil, indexToTest)
		assert.True(t, result.IsOk())
		fmt.Println(result)
		fmt.Println(testCluster.Status)
	})

	t.Run("delete: must terminate if API call return anything but 202 or 404", func(t *testing.T) {
		sch := runtime.NewScheme()
		assert.Nil(t, akov2.AddToScheme(sch))
		assert.Nil(t, corev1.AddToScheme(sch))

		indexToTest := &searchindex.SearchIndex{
			SearchIndex: akov2.SearchIndex{
				Name:           "test",
				DBName:         "testDB",
				CollectionName: "testCollection",
				Type:           "search",
			},
			AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{
				Analyzer: pointer.MakePtr("testAnalyzer"),
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

		fakeAtlasSearch := &fake.FakeAtlasSearch{
			DeleteIndexFunc: func(_ context.Context, _, _, _ string) error {
				return errors.New(http.StatusText(http.StatusInternalServerError))
			},
		}

		reconciler := &searchIndexReconcileRequest{
			ctx: &workflow.Context{
				Log:     zap.S(),
				OrgID:   "testOrgID",
				Context: context.Background(),
			},
			deployment:    testCluster,
			k8sClient:     nil,
			projectID:     "",
			indexName:     "testIndexName",
			searchService: fakeAtlasSearch,
		}

		result := reconciler.reconcileInternal("", nil, indexToTest)
		assert.False(t, result.IsOk())
		assert.True(t, reconciler.ctx.HasReason(status.SearchIndexStatusError))
	})

	t.Run("delete: must reconcile if AKO index ID is nil", func(t *testing.T) {
		sch := runtime.NewScheme()
		assert.Nil(t, akov2.AddToScheme(sch))
		assert.Nil(t, corev1.AddToScheme(sch))

		indexToTest := &searchindex.SearchIndex{
			SearchIndex: akov2.SearchIndex{
				Name:           "test",
				DBName:         "testDB",
				CollectionName: "testCollection",
				Type:           "search",
			},
			AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{
				Analyzer: pointer.MakePtr("testAnalyzer"),
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

		reconciler := &searchIndexReconcileRequest{
			ctx: &workflow.Context{
				Log:     zap.S(),
				OrgID:   "testOrgID",
				Context: context.Background(),
			},
			deployment: testCluster,
			k8sClient:  nil,
			projectID:  "",
			indexName:  "testIndexName",
		}

		result := reconciler.reconcileInternal("", nil, indexToTest)
		assert.True(t, result.IsOk())
	})

	t.Run("must return InProgress if index status is anything but ACTIVE", func(t *testing.T) {
		reconciler := &searchIndexReconcileRequest{
			ctx: &workflow.Context{
				Log:          zap.S(),
				OrgID:        "testOrgID",
				SdkClientSet: &atlas.ClientSet{SdkClient20250312006: &admin.APIClient{}},
				Context:      context.Background(),
			},
			deployment: nil,
			k8sClient:  nil,
			projectID:  "",
			indexName:  "testIndexName",
		}
		result := reconciler.reconcileInternal("", nil, &searchindex.SearchIndex{Status: pointer.MakePtr("NOT STARTED")})
		assert.True(t, result.IsInProgress())
	})

	t.Run("update: must not call update API if indexes are equal", func(t *testing.T) {
		reconciler := &searchIndexReconcileRequest{
			ctx: &workflow.Context{
				Log:          zap.S(),
				OrgID:        "testOrgID",
				SdkClientSet: &atlas.ClientSet{SdkClient20250312006: &admin.APIClient{}},
				Context:      context.Background(),
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
		result := reconciler.reconcileInternal("", idx, idx)
		assert.True(t, result.IsOk())
	})

	t.Run("update: must trigger index update if state in AKO and in Atlas is different", func(t *testing.T) {
		fakeAtlasSearch := &fake.FakeAtlasSearch{
			UpdateIndexFunc: func(_ context.Context, _, _ string, _ *searchindex.SearchIndex) (*searchindex.SearchIndex, error) {
				return &searchindex.SearchIndex{
					Status: pointer.MakePtr("NOT STARTED"),
				}, nil
			},
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

		reconciler := &searchIndexReconcileRequest{
			ctx: &workflow.Context{
				Log:     zap.S(),
				OrgID:   "testOrgID",
				Context: context.Background(),
			},
			deployment:    testCluster,
			k8sClient:     nil,
			projectID:     "",
			indexName:     "testIndexName",
			searchService: fakeAtlasSearch,
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
		result := reconciler.reconcileInternal("", idxInAKO, idxInAtlas)
		assert.True(t, result.IsInProgress())
	})

	t.Run("update: must terminate if API call returned anything but 201 or 200", func(t *testing.T) {
		fakeAtlasSearch := &fake.FakeAtlasSearch{
			UpdateIndexFunc: func(_ context.Context, _, _ string, _ *searchindex.SearchIndex) (*searchindex.SearchIndex, error) {
				return &searchindex.SearchIndex{
					Status: pointer.MakePtr("NOT STARTED"),
				}, errors.New(http.StatusText(http.StatusInternalServerError))
			},
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

		reconciler := &searchIndexReconcileRequest{
			ctx: &workflow.Context{
				Log:     zap.S(),
				OrgID:   "testOrgID",
				Context: context.Background(),
			},
			deployment:    testCluster,
			k8sClient:     nil,
			projectID:     "",
			indexName:     "testIndexName",
			searchService: fakeAtlasSearch,
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
		result := reconciler.reconcileInternal("", idxInAKO, idxInAtlas)
		assert.False(t, result.IsOk())
		assert.True(t, reconciler.ctx.HasReason(status.SearchIndexStatusError))
	})

	t.Run("update: must terminate if API call returned an empty index", func(t *testing.T) {
		mockSearchAPI := mockadmin.NewAtlasSearchApi(t)
		mockSearchAPI.EXPECT().
			UpdateAtlasSearchIndex(context.Background(), mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.UpdateAtlasSearchIndexApiRequest{ApiService: mockSearchAPI})
		mockSearchAPI.EXPECT().
			UpdateAtlasSearchIndexExecute(admin.UpdateAtlasSearchIndexApiRequest{ApiService: mockSearchAPI}).
			Return(
				nil,
				&http.Response{StatusCode: http.StatusCreated}, nil,
			)
		atlasSearch := searchindex.NewSearchIndexes(mockSearchAPI)

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

		reconciler := &searchIndexReconcileRequest{
			ctx: &workflow.Context{
				Log:     zap.S(),
				OrgID:   "testOrgID",
				Context: context.Background(),
			},
			deployment:    testCluster,
			k8sClient:     nil,
			projectID:     "",
			indexName:     "testIndexName",
			searchService: atlasSearch,
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
		result := reconciler.reconcileInternal("", idxInAKO, idxInAtlas)
		assert.False(t, result.IsOk())
		assert.True(t, reconciler.ctx.HasReason(status.SearchIndexStatusError))
	})

	t.Run("update: must terminate if index equality can not be confirmed", func(t *testing.T) {
		mockSearchAPI := mockadmin.NewAtlasSearchApi(t)
		atlasSearch := searchindex.NewSearchIndexes(mockSearchAPI)

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

		reconciler := &searchIndexReconcileRequest{
			ctx: &workflow.Context{
				Log:     zap.S(),
				OrgID:   "testOrgID",
				Context: context.Background(),
			},
			deployment:    testCluster,
			k8sClient:     nil,
			projectID:     "",
			indexName:     "testIndexName",
			searchService: atlasSearch,
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
						Dynamic: pointer.MakePtr(false),
						Fields:  &apiextensions.JSON{Raw: []byte{'i', 'n', 'v', 'a', 'l', 'i', 'd', 'j', 's', 'o', 'n'}},
					},
				},
			},
			AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{},
			ID:                         pointer.MakePtr("testID"),
			Status:                     nil,
		}
		result := reconciler.reconcileInternal("", idxInAKO, idxInAtlas)
		assert.False(t, result.IsOk())
		assert.True(t, reconciler.ctx.HasReason(status.SearchIndexStatusError))
	})

	t.Run("drop: must clear if the index disappeared from Atlas", func(t *testing.T) {
		mockSearchAPI := mockadmin.NewAtlasSearchApi(t)
		atlasSearch := searchindex.NewSearchIndexes(mockSearchAPI)
		for _, tc := range []struct {
			title          string
			atlasIndexName string
		}{
			{
				title:          "when name present",
				atlasIndexName: "testIndex",
			},
			{
				title:          "when name missing",
				atlasIndexName: "",
			},
		} {
			t.Run(tc.title, func(t *testing.T) {
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

				reconciler := &searchIndexReconcileRequest{
					ctx: &workflow.Context{
						Log:     zap.S(),
						OrgID:   "testOrgID",
						Context: context.Background(),
					},
					deployment:    testCluster,
					k8sClient:     nil,
					projectID:     "",
					indexName:     "testIndexName",
					searchService: atlasSearch,
				}
				result := reconciler.reconcileInternal(tc.atlasIndexName, nil, nil)
				assert.True(t, result.IsOk())
				assert.True(t, result.IsDeleted())
				assert.Empty(t, reconciler.ctx.Conditions())
			})
		}
	})
}
