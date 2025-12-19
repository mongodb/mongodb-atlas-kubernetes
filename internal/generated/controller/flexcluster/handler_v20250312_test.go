// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package flexcluster_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/crd2go/crd2go/k8s"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312011/admin"
	v20250312sdk "go.mongodb.org/atlas-sdk/v20250312011/admin"
	"go.mongodb.org/atlas-sdk/v20250312011/mockadmin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
	reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"

	crapi "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/flexcluster"
	crds "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/crds"
	akov2generated "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	result "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result"
	state "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

const (
	testClusterName = "test-cluster"
	testNamespace   = "ns1"
	testGroupName   = "test-group"
	testGroupID     = "62b6e34b3d91647abb20e7b8"
	testClusterID   = "cluster-id"
)

func TestHandleInitial(t *testing.T) {
	fixture := setupTestFixture(t)
	deletionProtection := false

	for _, tc := range []struct {
		title                      string
		flexCluster                *akov2generated.FlexCluster
		kubeObjects                []client.Object
		atlasCreateFlexClusterFunc func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error)
		interceptorFuncs           *interceptor.Funcs
		want                       ctrlstate.Result
		wantErr                    string
	}{
		{
			title:       "create flex cluster",
			flexCluster: defaultTestFlexCluster(testClusterName, testNamespace),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, nil),
			},
			atlasCreateFlexClusterFunc: func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error) {
				return &v20250312sdk.FlexClusterDescription20241113{Id: pointer.MakePtr(testClusterName)}, nil, nil
			},
			want: reenqueueResult(state.StateCreating, "Creating Flex Cluster."),
		},
		{
			title:       "corrupt flex type cluster",
			flexCluster: withFlexGVK(defaultTestFlexCluster("test-cluster1", "ns"), "bad-kind", "corrupt-api-version"),
			kubeObjects: []client.Object{
				defaultTestGroup("some-group", testNamespace, nil),
			},
			want:    ctrlstate.Result{NextState: state.StateInitial},
			wantErr: "failed to translate flex api params",
		},
		{
			title:       "corrupt api reply",
			flexCluster: defaultTestFlexCluster(testClusterName, testNamespace),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, nil),
			},
			atlasCreateFlexClusterFunc: func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error) {
				return nil, nil, nil
			},
			want:    ctrlstate.Result{NextState: state.StateInitial},
			wantErr: "failed to translate flex create response",
		},
		{
			title:       "missing group",
			flexCluster: setGroupRef(defaultTestFlexCluster("test-cluster1", "ns"), "not-found"),
			want:        ctrlstate.Result{NextState: state.StateInitial},
			wantErr:     "failed to get dependencies: failed to get group",
		},
		{
			title:       "group without id",
			flexCluster: setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, nil),
			},
			want:    ctrlstate.Result{NextState: state.StateInitial},
			wantErr: "failed to fetch referenced value groupId",
		},
		{
			title:       "atlas API error",
			flexCluster: setGroupRef(defaultTestFlexCluster("test-cluster2", "ns3"), testGroupName),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, "ns3", pointer.MakePtr(testGroupID)),
			},
			atlasCreateFlexClusterFunc: func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error) {
				return nil, nil, fmt.Errorf("atlas API error: cluster creation failed")
			},
			want:    ctrlstate.Result{NextState: state.StateInitial},
			wantErr: "failed to create flex cluster",
		},
		{
			title:       "patch status fails",
			flexCluster: defaultTestFlexCluster(testClusterName, testNamespace),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, nil),
			},
			atlasCreateFlexClusterFunc: func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error) {
				return &v20250312sdk.FlexClusterDescription20241113{Id: pointer.MakePtr(testClusterName)}, nil, nil
			},
			interceptorFuncs: &interceptor.Funcs{
				SubResourcePatch: func(ctx context.Context, client client.Client, subResourceName string, obj client.Object, patch client.Patch, opts ...client.SubResourcePatchOption) error {
					return fmt.Errorf("simulated status patch failure")
				},
			},
			want:    ctrlstate.Result{NextState: state.StateInitial},
			wantErr: "failed to patch flex cluster status",
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			fakeClient := fixture.buildFakeClient(tc.flexCluster, tc.kubeObjects, tc.interceptorFuncs)

			flexAPI := mockadmin.NewFlexClustersApi(t)
			if tc.atlasCreateFlexClusterFunc != nil {
				cluster, rsp, err := tc.atlasCreateFlexClusterFunc()
				req := admin.CreateFlexClusterApiRequest{ApiService: flexAPI}
				flexAPI.EXPECT().CreateFlexClusterWithParams(mock.Anything, mock.Anything).Return(req)
				flexAPI.EXPECT().CreateFlexClusterExecute(req).Return(cluster, rsp, err)
			}

			atlasClient := buildAtlasClient(flexAPI)
			handler := fixture.buildHandler(fakeClient, atlasClient, deletionProtection)

			ctx := context.Background()
			got, err := handler.HandleInitial(ctx, tc.flexCluster)
			assertError(t, err, tc.wantErr)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestHandleImportRequested(t *testing.T) {
	fixture := setupTestFixture(t)
	deletionProtection := false

	for _, tc := range []struct {
		title                   string
		flexCluster             *akov2generated.FlexCluster
		kubeObjects             []client.Object
		atlasGetFlexClusterFunc func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error)
		want                    ctrlstate.Result
		wantErr                 string
	}{
		{
			title: "import flex cluster with annotations",
			flexCluster: withAnnotations(
				setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
				map[string]string{
					"mongodb.com/external-name":     "external-cluster",
					"mongodb.com/external-group-id": testGroupID,
				},
			),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			atlasGetFlexClusterFunc: func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error) {
				return &v20250312sdk.FlexClusterDescription20241113{
					Id:        pointer.MakePtr("external-cluster"),
					StateName: pointer.MakePtr("IDLE"),
				}, nil, nil
			},
			want: nextStateResult(state.StateImported, "Imported Flex Cluster."),
		},
		{
			title:       "missing external-name annotation",
			flexCluster: defaultTestFlexCluster(testClusterName, testNamespace),
			want:        errorResult(state.StateImportRequested),
			wantErr:     "missing mongodb.com/external-name",
		},
		{
			title: "missing external-group-id annotation",
			flexCluster: withAnnotations(
				defaultTestFlexCluster(testClusterName, testNamespace),
				map[string]string{
					"mongodb.com/external-name": "external-cluster",
				},
			),
			want:    errorResult(state.StateImportRequested),
			wantErr: "missing mongodb.com/external-group-id",
		},
		{
			title: "patchStatus fails",
			flexCluster: withAnnotations(
				setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
				map[string]string{
					"mongodb.com/external-name":     "external-cluster",
					"mongodb.com/external-group-id": testGroupID,
				},
			),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			atlasGetFlexClusterFunc: func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error) {
				return nil, nil, fmt.Errorf("failed to get cluster")
			},
			want:    errorResult(state.StateImportRequested),
			wantErr: "failed to get cluster",
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			fakeClient := fixture.buildFakeClient(tc.flexCluster, tc.kubeObjects, nil)

			flexAPI := mockadmin.NewFlexClustersApi(t)
			if tc.atlasGetFlexClusterFunc != nil {
				cluster, rsp, err := tc.atlasGetFlexClusterFunc()
				req := admin.GetFlexClusterApiRequest{ApiService: flexAPI}
				flexAPI.EXPECT().GetFlexClusterWithParams(mock.Anything, mock.Anything).Return(req)
				flexAPI.EXPECT().GetFlexClusterExecute(req).Return(cluster, rsp, err)
			}

			atlasClient := buildAtlasClient(flexAPI)
			handler := fixture.buildHandler(fakeClient, atlasClient, deletionProtection)

			ctx := context.Background()
			got, err := handler.HandleImportRequested(ctx, tc.flexCluster)
			assertError(t, err, tc.wantErr)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestHandleCreating(t *testing.T) {
	fixture := setupTestFixture(t)
	deletionProtection := false

	for _, tc := range []struct {
		title                   string
		flexCluster             *akov2generated.FlexCluster
		kubeObjects             []client.Object
		atlasGetFlexClusterFunc func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error)
		want                    ctrlstate.Result
		wantErr                 string
	}{
		{
			title:       "cluster still creating",
			flexCluster: setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			atlasGetFlexClusterFunc: func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error) {
				return &v20250312sdk.FlexClusterDescription20241113{
					Id:        pointer.MakePtr(testClusterID),
					StateName: pointer.MakePtr("CREATING"),
				}, nil, nil
			},
			want: reenqueueResult(state.StateCreating, "Upserting Flex Cluster."),
		},
		{
			title:       "cluster created",
			flexCluster: setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			atlasGetFlexClusterFunc: func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error) {
				return &v20250312sdk.FlexClusterDescription20241113{
					Id:        pointer.MakePtr(testClusterID),
					StateName: pointer.MakePtr("IDLE"),
				}, nil, nil
			},
			want: nextStateResult(state.StateCreated, "Upserted Flex Cluster."),
		},
		{
			title:       "patchStatus fails",
			flexCluster: setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), "not-found"),
			want:        errorResult(state.StateCreating),
			wantErr:     "failed to get group",
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			fakeClient := fixture.buildFakeClient(tc.flexCluster, tc.kubeObjects, nil)

			flexAPI := mockadmin.NewFlexClustersApi(t)
			if tc.atlasGetFlexClusterFunc != nil {
				cluster, rsp, err := tc.atlasGetFlexClusterFunc()
				req := admin.GetFlexClusterApiRequest{ApiService: flexAPI}
				flexAPI.EXPECT().GetFlexClusterWithParams(mock.Anything, mock.Anything).Return(req)
				flexAPI.EXPECT().GetFlexClusterExecute(req).Return(cluster, rsp, err)
			}

			atlasClient := buildAtlasClient(flexAPI)
			handler := fixture.buildHandler(fakeClient, atlasClient, deletionProtection)

			ctx := context.Background()
			got, err := handler.HandleCreating(ctx, tc.flexCluster)
			assertError(t, err, tc.wantErr)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestHandleUpdating(t *testing.T) {
	fixture := setupTestFixture(t)
	deletionProtection := false

	for _, tc := range []struct {
		title                   string
		flexCluster             *akov2generated.FlexCluster
		kubeObjects             []client.Object
		atlasGetFlexClusterFunc func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error)
		want                    ctrlstate.Result
		wantErr                 string
	}{
		{
			title:       "cluster still updating",
			flexCluster: setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			atlasGetFlexClusterFunc: func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error) {
				return &v20250312sdk.FlexClusterDescription20241113{
					Id:        pointer.MakePtr(testClusterID),
					StateName: pointer.MakePtr("UPDATING"),
				}, nil, nil
			},
			want: reenqueueResult(state.StateUpdating, "Upserting Flex Cluster."),
		},
		{
			title:       "cluster updated",
			flexCluster: setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			atlasGetFlexClusterFunc: func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error) {
				return &v20250312sdk.FlexClusterDescription20241113{
					Id:        pointer.MakePtr(testClusterID),
					StateName: pointer.MakePtr("IDLE"),
				}, nil, nil
			},
			want: nextStateResult(state.StateUpdated, "Upserted Flex Cluster."),
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			fakeClient := fixture.buildFakeClient(tc.flexCluster, tc.kubeObjects, nil)

			flexAPI := mockadmin.NewFlexClustersApi(t)
			if tc.atlasGetFlexClusterFunc != nil {
				cluster, rsp, err := tc.atlasGetFlexClusterFunc()
				req := admin.GetFlexClusterApiRequest{ApiService: flexAPI}
				flexAPI.EXPECT().GetFlexClusterWithParams(mock.Anything, mock.Anything).Return(req)
				flexAPI.EXPECT().GetFlexClusterExecute(req).Return(cluster, rsp, err)
			}

			atlasClient := buildAtlasClient(flexAPI)
			handler := fixture.buildHandler(fakeClient, atlasClient, deletionProtection)

			ctx := context.Background()
			got, err := handler.HandleUpdating(ctx, tc.flexCluster)
			assertError(t, err, tc.wantErr)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestHandleCreated(t *testing.T) {
	fixture := setupTestFixture(t)
	deletionProtection := false

	for _, tc := range []struct {
		title                      string
		flexCluster                *akov2generated.FlexCluster
		kubeObjects                []client.Object
		atlasUpdateFlexClusterFunc func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error)
		interceptorFuncs           *interceptor.Funcs
		want                       ctrlstate.Result
		wantErr                    string
	}{
		{
			title: "no update needed",
			flexCluster: withStateTracker(withObservedGeneration(
				withGeneration(
					setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
					1,
				),
				1,
			)),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			want: nextStateResult(state.StateCreated, "Flex cluster up to date. No update required."),
		},
		{
			title: "update needed",
			flexCluster: withObservedGeneration(
				withGeneration(
					setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
					2,
				),
				1,
			),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			atlasUpdateFlexClusterFunc: func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error) {
				return &v20250312sdk.FlexClusterDescription20241113{
					Id:        pointer.MakePtr(testClusterID),
					StateName: pointer.MakePtr("IDLE"),
				}, nil, nil
			},
			want: reenqueueResult(state.StateUpdating, "Updating Flex Cluster."),
		},
		{
			title: "update fails",
			flexCluster: withObservedGeneration(
				withGeneration(
					setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
					2,
				),
				0,
			),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			atlasUpdateFlexClusterFunc: func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error) {
				return nil, nil, fmt.Errorf("update failed")
			},
			want:    errorResult(state.StateCreated),
			wantErr: "failed to get update cluster",
		},
		{
			title: "ShouldUpdate fails with invalid reapply period",
			flexCluster: withAnnotations(
				withObservedGeneration(
					withGeneration(
						setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
						1,
					),
					2,
				),
				map[string]string{
					"mongodb.internal.com/reapply-timestamp": "1000000000000",
					"mongodb.com/reapply-period":             "invalid-duration",
				},
			),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			want:    errorResult(state.StateCreated),
			wantErr: "failed to check reapply period",
		},
		{
			title: "ToAPI params translation fails",
			flexCluster: withFlexGVK(
				withObservedGeneration(
					withGeneration(
						setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
						2,
					),
					1,
				),
				"BadKind",
				"invalid-api-version",
			),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			want:    errorResult(state.StateCreated),
			wantErr: "failed to translate update flex cluster parameters",
		},
		{
			title: "FromAPI translation fails",
			flexCluster: withObservedGeneration(
				withGeneration(
					setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
					2,
				),
				1,
			),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			atlasUpdateFlexClusterFunc: func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error) {
				// Return nil to trigger "source is nil" check in FromAPI
				return nil, nil, nil
			},
			want:    errorResult(state.StateCreated),
			wantErr: "failed to translate update cluster response",
		},
		{
			title: "status patch fails",
			flexCluster: withObservedGeneration(
				withGeneration(
					setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
					2,
				),
				1,
			),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			atlasUpdateFlexClusterFunc: func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error) {
				return &v20250312sdk.FlexClusterDescription20241113{
					Id:        pointer.MakePtr(testClusterID),
					StateName: pointer.MakePtr("IDLE"),
				}, nil, nil
			},
			interceptorFuncs: &interceptor.Funcs{
				SubResourcePatch: func(ctx context.Context, client client.Client, subResourceName string, obj client.Object, patch client.Patch, opts ...client.SubResourcePatchOption) error {
					return fmt.Errorf("simulated status patch failure")
				},
			},
			want:    errorResult(state.StateCreated),
			wantErr: "failed to patch cluster",
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			fakeClient := fixture.buildFakeClient(tc.flexCluster, tc.kubeObjects, tc.interceptorFuncs)

			flexAPI := mockadmin.NewFlexClustersApi(t)
			if tc.atlasUpdateFlexClusterFunc != nil {
				cluster, rsp, err := tc.atlasUpdateFlexClusterFunc()
				req := admin.UpdateFlexClusterApiRequest{ApiService: flexAPI}
				flexAPI.EXPECT().UpdateFlexClusterWithParams(mock.Anything, mock.Anything).Return(req)
				flexAPI.EXPECT().UpdateFlexClusterExecute(req).Return(cluster, rsp, err)
			}

			atlasClient := buildAtlasClient(flexAPI)
			handler := fixture.buildHandler(fakeClient, atlasClient, deletionProtection)

			ctx := context.Background()
			got, err := handler.HandleCreated(ctx, tc.flexCluster)
			assertError(t, err, tc.wantErr)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestHandleImported(t *testing.T) {
	fixture := setupTestFixture(t)
	deletionProtection := false

	for _, tc := range []struct {
		title                      string
		flexCluster                *akov2generated.FlexCluster
		kubeObjects                []client.Object
		atlasUpdateFlexClusterFunc func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error)
		want                       ctrlstate.Result
		wantErr                    string
	}{
		{
			title: "no update needed",
			flexCluster: withStateTracker(withObservedGeneration(
				withGeneration(
					setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
					1,
				),
				1,
			)),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			want: nextStateResult(state.StateCreated, "Flex cluster up to date. No update required."),
		},
		{
			title: "update needed",
			flexCluster: withObservedGeneration(
				withGeneration(
					setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
					2,
				),
				1,
			),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			atlasUpdateFlexClusterFunc: func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error) {
				return &v20250312sdk.FlexClusterDescription20241113{
					Id:        pointer.MakePtr(testClusterID),
					StateName: pointer.MakePtr("IDLE"),
				}, nil, nil
			},
			want: reenqueueResult(state.StateUpdating, "Updating Flex Cluster."),
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			fakeClient := fixture.buildFakeClient(tc.flexCluster, tc.kubeObjects, nil)

			flexAPI := mockadmin.NewFlexClustersApi(t)
			if tc.atlasUpdateFlexClusterFunc != nil {
				cluster, rsp, err := tc.atlasUpdateFlexClusterFunc()
				req := admin.UpdateFlexClusterApiRequest{ApiService: flexAPI}
				flexAPI.EXPECT().UpdateFlexClusterWithParams(mock.Anything, mock.Anything).Return(req)
				flexAPI.EXPECT().UpdateFlexClusterExecute(req).Return(cluster, rsp, err)
			}

			atlasClient := buildAtlasClient(flexAPI)
			handler := fixture.buildHandler(fakeClient, atlasClient, deletionProtection)

			ctx := context.Background()
			got, err := handler.HandleImported(ctx, tc.flexCluster)
			assertError(t, err, tc.wantErr)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestHandleUpdated(t *testing.T) {
	fixture := setupTestFixture(t)
	deletionProtection := false

	for _, tc := range []struct {
		title                      string
		flexCluster                *akov2generated.FlexCluster
		kubeObjects                []client.Object
		atlasUpdateFlexClusterFunc func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error)
		want                       ctrlstate.Result
		wantErr                    string
	}{
		{
			title: "no update needed",
			flexCluster: withStateTracker(withObservedGeneration(
				withGeneration(
					setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
					1,
				),
				1,
			)),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			want: nextStateResult(state.StateUpdated, "Flex cluster up to date. No update required."),
		},
		{
			title: "update needed",
			flexCluster: withObservedGeneration(
				withGeneration(
					setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
					2,
				),
				1,
			),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			atlasUpdateFlexClusterFunc: func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error) {
				return &v20250312sdk.FlexClusterDescription20241113{
					Id:        pointer.MakePtr(testClusterID),
					StateName: pointer.MakePtr("IDLE"),
				}, nil, nil
			},
			want: reenqueueResult(state.StateUpdating, "Updating Flex Cluster."),
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			fakeClient := fixture.buildFakeClient(tc.flexCluster, tc.kubeObjects, nil)

			flexAPI := mockadmin.NewFlexClustersApi(t)
			if tc.atlasUpdateFlexClusterFunc != nil {
				cluster, rsp, err := tc.atlasUpdateFlexClusterFunc()
				req := admin.UpdateFlexClusterApiRequest{ApiService: flexAPI}
				flexAPI.EXPECT().UpdateFlexClusterWithParams(mock.Anything, mock.Anything).Return(req)
				flexAPI.EXPECT().UpdateFlexClusterExecute(req).Return(cluster, rsp, err)
			}

			atlasClient := buildAtlasClient(flexAPI)
			handler := fixture.buildHandler(fakeClient, atlasClient, deletionProtection)

			ctx := context.Background()
			got, err := handler.HandleUpdated(ctx, tc.flexCluster)
			assertError(t, err, tc.wantErr)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestHandleDeletionRequested(t *testing.T) {
	fixture := setupTestFixture(t)

	for _, tc := range []struct {
		title                      string
		flexCluster                *akov2generated.FlexCluster
		kubeObjects                []client.Object
		deletionProtection         bool
		atlasDeleteFlexClusterFunc func() (*http.Response, error)
		want                       ctrlstate.Result
		wantErr                    string
	}{
		{
			title: "delete with protection disabled",
			flexCluster: withStatus(
				setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
			),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			deletionProtection: false,
			atlasDeleteFlexClusterFunc: func() (*http.Response, error) {
				return nil, nil
			},
			want: reenqueueResult(state.StateDeleting, "Deleting Flex Cluster."),
		},
		{
			title: "delete with keep annotation",
			flexCluster: withResourcePolicyKeep(
				withStatus(
					setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
				),
			),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			deletionProtection: false,
			want:               nextStateResult(state.StateDeleted, "Flex Cluster deleted."),
		},
		{
			title: "delete with deletion protection enabled",
			flexCluster: withStatus(
				setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
			),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			deletionProtection: true,
			want:               nextStateResult(state.StateDeleted, "Flex Cluster deleted."),
		},
		{
			title: "delete unmanaged cluster",
			flexCluster: setGroupRef(
				defaultTestFlexCluster(testClusterName, testNamespace),
				testGroupName,
			),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			deletionProtection: false,
			want:               nextStateResult(state.StateDeleted, "Flex Cluster is unamanged."),
		},
		{
			title: "cluster already deleted in Atlas",
			flexCluster: withStatus(
				setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
			),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			deletionProtection: false,
			atlasDeleteFlexClusterFunc: func() (*http.Response, error) {
				apiErr := &v20250312sdk.GenericOpenAPIError{}
				apiErr.SetModel(v20250312sdk.ApiError{ErrorCode: "CLUSTER_NOT_FOUND"})
				return nil, apiErr
			},
			want: nextStateResult(state.StateDeleted, "Flex Cluster was deleted in Atlas."),
		},
		{
			title: "delete fails",
			flexCluster: withStatus(
				setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
			),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			deletionProtection: false,
			atlasDeleteFlexClusterFunc: func() (*http.Response, error) {
				return nil, fmt.Errorf("delete failed")
			},
			want:    errorResult(state.StateDeletionRequested),
			wantErr: "failed to delete flex cluster",
		},
		{
			title: "get dependencies fails",
			flexCluster: withStatus(
				setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), "not-found"),
			),
			deletionProtection: false,
			want:               errorResult(state.StateDeletionRequested),
			wantErr:            "failed to get dependencies",
		},
		{
			title: "ToAPI translation fails",
			flexCluster: withFlexGVK(
				withStatus(
					setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
				),
				"BadKind",
				"wrong-api-version",
			),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			deletionProtection: false,
			want:               errorResult(state.StateDeletionRequested),
			wantErr:            "failed to translate flex api params",
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			fakeClient := fixture.buildFakeClient(tc.flexCluster, tc.kubeObjects, nil)

			flexAPI := mockadmin.NewFlexClustersApi(t)
			if tc.atlasDeleteFlexClusterFunc != nil {
				rsp, err := tc.atlasDeleteFlexClusterFunc()
				req := admin.DeleteFlexClusterApiRequest{ApiService: flexAPI}
				flexAPI.EXPECT().DeleteFlexClusterWithParams(mock.Anything, mock.Anything).Return(req)
				flexAPI.EXPECT().DeleteFlexClusterExecute(req).Return(rsp, err)
			}

			atlasClient := buildAtlasClient(flexAPI)
			handler := fixture.buildHandler(fakeClient, atlasClient, tc.deletionProtection)

			ctx := context.Background()
			got, err := handler.HandleDeletionRequested(ctx, tc.flexCluster)
			assertError(t, err, tc.wantErr)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestHandleDeleting(t *testing.T) {
	fixture := setupTestFixture(t)
	deletionProtection := false

	for _, tc := range []struct {
		title                   string
		flexCluster             *akov2generated.FlexCluster
		kubeObjects             []client.Object
		atlasGetFlexClusterFunc func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error)
		want                    ctrlstate.Result
		wantErr                 string
	}{
		{
			title:       "cluster deleted",
			flexCluster: setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			atlasGetFlexClusterFunc: func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error) {
				apiErr := &v20250312sdk.GenericOpenAPIError{}
				apiErr.SetModel(v20250312sdk.ApiError{ErrorCode: "CLUSTER_NOT_FOUND"})
				return nil, nil, apiErr
			},
			want: nextStateResult(state.StateDeleted, "Deleted."),
		},
		{
			title:       "cluster still deleting",
			flexCluster: setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			atlasGetFlexClusterFunc: func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error) {
				return &v20250312sdk.FlexClusterDescription20241113{
					Id:        pointer.MakePtr(testClusterID),
					StateName: pointer.MakePtr("DELETING"),
				}, nil, nil
			},
			want: reenqueueResult(state.StateDeleting, "Deleting Flex Cluster."),
		},
		{
			title:       "get cluster fails",
			flexCluster: setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			atlasGetFlexClusterFunc: func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error) {
				return nil, nil, fmt.Errorf("get failed")
			},
			want:    errorResult(state.StateDeletionRequested),
			wantErr: "failed to delete flexcluster",
		},
		{
			title:       "get dependencies fails",
			flexCluster: setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), "not-found"),
			want:        errorResult(state.StateDeleting),
			wantErr:     "failed to get dependencies",
		},
		{
			title: "ToAPI translation fails",
			flexCluster: withFlexGVK(
				setGroupRef(defaultTestFlexCluster(testClusterName, testNamespace), testGroupName),
				"BadKind",
				"malformed-api-version",
			),
			kubeObjects: []client.Object{
				defaultTestGroup(testGroupName, testNamespace, pointer.MakePtr(testGroupID)),
			},
			want:    errorResult(state.StateDeleting),
			wantErr: "failed to translate flex api params",
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			fakeClient := fixture.buildFakeClient(tc.flexCluster, tc.kubeObjects, nil)

			flexAPI := mockadmin.NewFlexClustersApi(t)
			if tc.atlasGetFlexClusterFunc != nil {
				cluster, rsp, err := tc.atlasGetFlexClusterFunc()
				req := admin.GetFlexClusterApiRequest{ApiService: flexAPI}
				flexAPI.EXPECT().GetFlexClusterWithParams(mock.Anything, mock.Anything).Return(req)
				flexAPI.EXPECT().GetFlexClusterExecute(req).Return(cluster, rsp, err)
			}

			atlasClient := buildAtlasClient(flexAPI)
			handler := fixture.buildHandler(fakeClient, atlasClient, deletionProtection)

			ctx := context.Background()
			got, err := handler.HandleDeleting(ctx, tc.flexCluster)
			assertError(t, err, tc.wantErr)
			require.Equal(t, tc.want, got)
		})
	}
}

// Helper functions

// testFixture holds common test dependencies
type testFixture struct {
	scheme     *runtime.Scheme
	translator crapi.Translator
}

// setupTestFixture creates a test fixture with common dependencies
func setupTestFixture(t *testing.T) *testFixture {
	fixture := &testFixture{
		scheme: createTestScheme(t),
	}

	tr, err := translatorFromEmbeddedCRD(fixture.scheme, "FlexCluster", "v1", "v20250312")
	require.NoError(t, err)
	fixture.translator = tr

	return fixture
}

// buildFakeClient creates a fake Kubernetes client with the given objects and optional interceptor
func (f *testFixture) buildFakeClient(flexCluster *akov2generated.FlexCluster, kubeObjects []client.Object, funcs *interceptor.Funcs) client.Client {
	allObjects := kubeObjects
	clientBuilder := fake.NewClientBuilder().WithScheme(f.scheme)

	if funcs != nil {
		clientBuilder = clientBuilder.WithInterceptorFuncs(*funcs)
	}

	if flexCluster != nil {
		allObjects = append([]client.Object{flexCluster}, kubeObjects...)
		clientBuilder = clientBuilder.WithObjects(allObjects...).WithStatusSubresource(flexCluster)
	} else {
		clientBuilder = clientBuilder.WithObjects(allObjects...)
	}

	return clientBuilder.Build()
}

// buildHandler creates a handler with the given dependencies
func (f *testFixture) buildHandler(fakeClient client.Client, atlasClient *v20250312sdk.APIClient, deletionProtection bool) *flexcluster.Handlerv20250312 {
	return flexcluster.NewHandlerv20250312(fakeClient, atlasClient, f.translator, deletionProtection)
}

// buildAtlasClient creates an Atlas API client with the given mock API
func buildAtlasClient(flexAPI *mockadmin.FlexClustersApi) *v20250312sdk.APIClient {
	return &v20250312sdk.APIClient{
		FlexClustersApi: flexAPI,
	}
}

func setGroupRef(flexCluster *akov2generated.FlexCluster, groupRef string) *akov2generated.FlexCluster {
	flexCluster.Spec.V20250312.GroupRef = &k8s.LocalReference{Name: groupRef}
	return flexCluster
}

func translatorFromEmbeddedCRD(scheme *runtime.Scheme, kind, apiVersion, majorVersion string) (crapi.Translator, error) {
	crd, err := crds.EmbeddedCRD(kind)
	if err != nil {
		return nil, fmt.Errorf("failed to get embedded CRD for %s: %w", kind, err)
	}
	tr, err := crapi.NewTranslator(scheme, crd, apiVersion, majorVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get translator for %s: %w", kind, err)
	}
	return tr, nil
}

// createTestScheme creates a runtime scheme for testing
func createTestScheme(t *testing.T) *runtime.Scheme {
	scheme := runtime.NewScheme()
	require.NoError(t, akov2generated.AddToScheme(scheme))
	require.NoError(t, corev1.AddToScheme(scheme))
	return scheme
}

// defaultTestGroup creates a basic Group for testing
func defaultTestGroup(name, namespace string, id *string) *akov2generated.Group {
	return &akov2generated.Group{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Group",
			APIVersion: "atlas.generated.mongodb.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Status: akov2generated.GroupStatus{
			V20250312: &akov2generated.GroupStatusV20250312{
				Id: id,
			},
		},
	}
}

// defaultTestFlexCluster creates a basic FlexCluster for testing
func defaultTestFlexCluster(name, namespace string) *akov2generated.FlexCluster {
	return &akov2generated.FlexCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: akov2generated.FlexClusterSpec{
			V20250312: &akov2generated.FlexClusterSpecV20250312{
				Entry: &akov2generated.FlexClusterSpecV20250312Entry{
					Name: name,
					ProviderSettings: akov2generated.ProviderSettings{
						BackingProviderName: "AWS",
						RegionName:          "us-east-1",
					},
				},
			},
		},
	}
}

func withFlexGVK(flexCluster *akov2generated.FlexCluster, kind string, apiVersion string) *akov2generated.FlexCluster {
	flexCluster.TypeMeta.Kind = kind
	flexCluster.TypeMeta.APIVersion = apiVersion
	return flexCluster
}

func reenqueueResult(state state.ResourceState, msg string) ctrlstate.Result {
	return ctrlstate.Result{
		Result:    reconcile.Result{RequeueAfter: result.DefaultRequeueTIme},
		NextState: state,
		StateMsg:  msg,
	}
}

func nextStateResult(nextState state.ResourceState, msg string) ctrlstate.Result {
	return ctrlstate.Result{
		NextState: nextState,
		StateMsg:  msg,
	}
}

func errorResult(nextState state.ResourceState) ctrlstate.Result {
	return ctrlstate.Result{
		NextState: nextState,
	}
}

// assertError standardizes error checking in tests
func assertError(t *testing.T, err error, wantErr string) {
	if wantErr != "" {
		assert.Error(t, err)
		assert.Contains(t, err.Error(), wantErr)
	} else {
		require.NoError(t, err)
	}
}

// withStatus creates a FlexCluster with status (needed to pass nil check in HandleDeletionRequested)
func withStatus(flexCluster *akov2generated.FlexCluster) *akov2generated.FlexCluster {
	flexCluster.Status = akov2generated.FlexClusterStatus{
		V20250312: &akov2generated.FlexClusterStatusV20250312{},
	}
	return flexCluster
}

// withAnnotations adds annotations to a FlexCluster
func withAnnotations(flexCluster *akov2generated.FlexCluster, annotations map[string]string) *akov2generated.FlexCluster {
	if flexCluster.Annotations == nil {
		flexCluster.Annotations = make(map[string]string)
	}
	for k, v := range annotations {
		flexCluster.Annotations[k] = v
	}
	return flexCluster
}

// withGeneration sets the generation on a FlexCluster
func withGeneration(flexCluster *akov2generated.FlexCluster, generation int64) *akov2generated.FlexCluster {
	flexCluster.Generation = generation
	return flexCluster
}

// withObservedGeneration sets the observed generation in status conditions
func withObservedGeneration(flexCluster *akov2generated.FlexCluster, observedGen int64) *akov2generated.FlexCluster {
	if flexCluster.Status.Conditions == nil {
		flexCluster.Status.Conditions = &[]metav1.Condition{}
	}
	conditions := *flexCluster.Status.Conditions
	conditions = append(conditions, metav1.Condition{
		Type:               state.StateCondition,
		ObservedGeneration: observedGen,
		Status:             metav1.ConditionTrue,
	})
	// Allocate a new slice to avoid storing a pointer to a local variable
	newConditions := make([]metav1.Condition, len(conditions))
	copy(newConditions, conditions)
	flexCluster.Status.Conditions = &newConditions
	return flexCluster
}

// withResourcePolicyKeep adds the keep annotation for deletion protection
func withResourcePolicyKeep(flexCluster *akov2generated.FlexCluster) *akov2generated.FlexCluster {
	if flexCluster.Annotations == nil {
		flexCluster.Annotations = make(map[string]string)
	}
	flexCluster.Annotations["mongodb.com/atlas-resource-policy"] = "keep"
	return flexCluster
}

// withStateTracker adds the state tracker annotation to simulate a previously reconciled object
func withStateTracker(flexCluster *akov2generated.FlexCluster, deps ...client.Object) *akov2generated.FlexCluster {
	if flexCluster.Annotations == nil {
		flexCluster.Annotations = make(map[string]string)
	}
	flexCluster.Annotations[ctrlstate.AnnotationStateTracker] = ctrlstate.ComputeStateTracker(flexCluster, deps...)
	return flexCluster
}
