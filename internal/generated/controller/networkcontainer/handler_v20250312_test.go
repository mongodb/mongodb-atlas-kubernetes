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

package networkcontainer_test

import (
	"context"
	"errors"
	"testing"

	ctrlstate "github.com/crd2go/constate"
	state "github.com/crd2go/constate/state"
	crapi "github.com/crd2go/crapi"
	"github.com/crd2go/crd2go/k8s"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	v20250312sdk "go.mongodb.org/atlas-sdk/v20250312026/admin"
	"go.mongodb.org/atlas-sdk/v20250312026/mockadmin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	akov2generated "github.com/mongodb/mongodb-atlas-kubernetes/v2/generated/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/networkcontainer"
	crds "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/crds"
)

const (
	testNamespace   = "ns1"
	testGroupName   = "test-group"
	testGroupID     = "62b6e34b3d91647abb20e7b8"
	testContainerID = "32b6e34b3d91647abb20e7b8"
)

func TestHandleInitial(t *testing.T) {
	for _, tc := range []struct {
		title     string
		group     *akov2generated.Group
		createErr error
		wantState state.ResourceState
		wantErr   string
	}{
		{
			title:     "creates the container with the groupId resolved from groupRef",
			group:     testGroup(new(testGroupID)),
			wantState: state.StateCreated,
		},
		{
			title:     "referenced group without id fails translation",
			group:     testGroup(nil),
			wantState: state.StateInitial,
			wantErr:   "failed to translate NetworkContainer API parameters to Atlas",
		},
		{
			title:     "atlas create error is returned",
			group:     testGroup(new(testGroupID)),
			createErr: errors.New("boom"),
			wantState: state.StateInitial,
			wantErr:   "failed to create NetworkContainer",
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			nc := testNetworkContainer()
			kubeClient := testKubeClient(t, nc, tc.group)
			peeringAPI := mockadmin.NewNetworkPeeringAPI(t)

			var gotParams *v20250312sdk.CreateGroupContainerApiParams
			if tc.group.Status.V20250312.Id != nil {
				req := v20250312sdk.CreateGroupContainerApiRequest{ApiService: peeringAPI}
				peeringAPI.EXPECT().CreateGroupContainerWithParams(mock.Anything, mock.Anything).
					Run(func(_ context.Context, args *v20250312sdk.CreateGroupContainerApiParams) { gotParams = args }).
					Return(req)
				response := &v20250312sdk.CloudProviderContainer{Id: new(testContainerID), VpcId: new("vpc-123")}
				peeringAPI.EXPECT().CreateGroupContainerExecute(req).Return(response, nil, tc.createErr)
			}

			got, err := testHandler(t, kubeClient, peeringAPI).HandleInitial(context.Background(), nc)
			assertError(t, err, tc.wantErr)
			assert.Equal(t, tc.wantState, got.NextState)
			if tc.wantErr != "" {
				return
			}

			require.NotNil(t, gotParams)
			assert.Equal(t, testGroupID, gotParams.GroupId)
			assert.Equal(t, "AWS", gotParams.CloudProviderContainer.GetProviderName())
			assert.Equal(t, "US_EAST_1", gotParams.CloudProviderContainer.GetRegionName())
			assert.Equal(t, "10.128.0.0/21", gotParams.CloudProviderContainer.GetAtlasCidrBlock())

			patched := &akov2generated.NetworkContainer{}
			require.NoError(t, kubeClient.Get(context.Background(), client.ObjectKeyFromObject(nc), patched))
			assert.Equal(t, new(testContainerID), patched.Status.V20250312.Id)
			assert.Equal(t, new("vpc-123"), patched.Status.V20250312.VpcId)
		})
	}
}

func TestHandleImportRequested(t *testing.T) {
	nc := testNetworkContainer()
	nc.Annotations = map[string]string{"mongodb.com/external-id": testContainerID}
	kubeClient := testKubeClient(t, nc, testGroup(new(testGroupID)))
	peeringAPI := mockadmin.NewNetworkPeeringAPI(t)
	req := v20250312sdk.GetGroupContainerApiRequest{ApiService: peeringAPI}
	peeringAPI.EXPECT().GetGroupContainer(mock.Anything, testGroupID, testContainerID).Return(req)
	peeringAPI.EXPECT().GetGroupContainerExecute(req).Return(&v20250312sdk.CloudProviderContainer{Id: new(testContainerID)}, nil, nil)

	got, err := testHandler(t, kubeClient, peeringAPI).HandleImportRequested(context.Background(), nc)
	require.NoError(t, err)
	assert.Equal(t, state.StateImported, got.NextState)

	patched := &akov2generated.NetworkContainer{}
	require.NoError(t, kubeClient.Get(context.Background(), client.ObjectKeyFromObject(nc), patched))
	assert.Equal(t, new(testContainerID), patched.Status.V20250312.Id)
}

func TestHandleCreated(t *testing.T) {
	for _, tc := range []struct {
		title       string
		specChanged bool
		wantState   state.ResourceState
	}{
		{
			title:     "unchanged spec does not call Atlas",
			wantState: state.StateCreated,
		},
		{
			title:       "changed spec updates the container by status id",
			specChanged: true,
			wantState:   state.StateUpdated,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			group := testGroup(new(testGroupID))
			nc := withStatusID(testNetworkContainer())
			nc.Annotations = map[string]string{ctrlstate.AnnotationStateTracker: ctrlstate.ComputeStateTracker(nc, group)}
			if tc.specChanged {
				nc.Generation++
			}
			kubeClient := testKubeClient(t, nc, group)
			peeringAPI := mockadmin.NewNetworkPeeringAPI(t)

			var gotParams *v20250312sdk.UpdateGroupContainerApiParams
			if tc.specChanged {
				req := v20250312sdk.UpdateGroupContainerApiRequest{ApiService: peeringAPI}
				peeringAPI.EXPECT().UpdateGroupContainerWithParams(mock.Anything, mock.Anything).
					Run(func(_ context.Context, args *v20250312sdk.UpdateGroupContainerApiParams) { gotParams = args }).
					Return(req)
				peeringAPI.EXPECT().UpdateGroupContainerExecute(req).Return(&v20250312sdk.CloudProviderContainer{Id: new(testContainerID)}, nil, nil)
			}

			got, err := testHandler(t, kubeClient, peeringAPI).HandleCreated(context.Background(), nc)
			require.NoError(t, err)
			assert.Equal(t, tc.wantState, got.NextState)
			if tc.specChanged {
				require.NotNil(t, gotParams)
				assert.Equal(t, testGroupID, gotParams.GroupId)
				assert.Equal(t, testContainerID, gotParams.ContainerId)
				assert.Equal(t, "10.128.0.0/21", gotParams.CloudProviderContainer.GetAtlasCidrBlock())
			}
		})
	}
}

func TestHandleDeletionRequested(t *testing.T) {
	for _, tc := range []struct {
		title       string
		container   *akov2generated.NetworkContainer
		deleteErr   error
		expectCall  bool
		wantState   state.ResourceState
		wantErr     string
		protectKeep bool
	}{
		{
			title:      "deletes the container from Atlas",
			container:  withStatusID(testNetworkContainer()),
			expectCall: true,
			wantState:  state.StateDeleted,
		},
		{
			title:      "container already gone is treated as deleted",
			container:  withStatusID(testNetworkContainer()),
			deleteErr:  atlasError("CLOUD_PROVIDER_CONTAINER_NOT_FOUND"),
			expectCall: true,
			wantState:  state.StateDeleted,
		},
		{
			title:      "container still in use is returned as an error to retry",
			container:  withStatusID(testNetworkContainer()),
			deleteErr:  atlasError("CONTAINERS_IN_USE"),
			expectCall: true,
			wantState:  state.StateDeletionRequested,
			wantErr:    "is still in use",
		},
		{
			title:     "container never created in Atlas is deleted without calling Atlas",
			container: testNetworkContainer(),
			wantState: state.StateDeleted,
		},
		{
			title:       "keep policy skips deletion in Atlas",
			container:   withStatusID(testNetworkContainer()),
			protectKeep: true,
			wantState:   state.StateDeleted,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			if tc.protectKeep {
				tc.container.Annotations = map[string]string{"mongodb.com/atlas-resource-policy": "keep"}
			}
			kubeClient := testKubeClient(t, tc.container, testGroup(new(testGroupID)))
			peeringAPI := mockadmin.NewNetworkPeeringAPI(t)
			if tc.expectCall {
				req := v20250312sdk.DeleteGroupContainerApiRequest{ApiService: peeringAPI}
				peeringAPI.EXPECT().DeleteGroupContainer(mock.Anything, testGroupID, testContainerID).Return(req)
				peeringAPI.EXPECT().DeleteGroupContainerExecute(req).Return(nil, tc.deleteErr)
			}

			got, err := testHandler(t, kubeClient, peeringAPI).HandleDeletionRequested(context.Background(), tc.container)
			assertError(t, err, tc.wantErr)
			assert.Equal(t, tc.wantState, got.NextState)
		})
	}
}

func testHandler(t *testing.T, kubeClient client.Client, peeringAPI *mockadmin.NetworkPeeringAPI) *networkcontainer.Handlerv20250312 {
	t.Helper()
	crd, err := crds.EmbeddedCRD("NetworkContainer")
	require.NoError(t, err)
	translator, err := crapi.NewTranslator(testScheme(t), crd, "v1", "v20250312")
	require.NoError(t, err)
	atlasClient := &v20250312sdk.APIClient{NetworkPeeringAPI: peeringAPI}
	return networkcontainer.NewHandlerv20250312(kubeClient, atlasClient, translator, false)
}

func testKubeClient(t *testing.T, nc *akov2generated.NetworkContainer, objs ...client.Object) client.Client {
	t.Helper()
	return fake.NewClientBuilder().
		WithScheme(testScheme(t)).
		WithObjects(append([]client.Object{nc}, objs...)...).
		WithStatusSubresource(nc).
		Build()
}

func testScheme(t *testing.T) *runtime.Scheme {
	t.Helper()
	scheme := runtime.NewScheme()
	require.NoError(t, akov2generated.AddToScheme(scheme))
	require.NoError(t, corev1.AddToScheme(scheme))
	return scheme
}

func testGroup(id *string) *akov2generated.Group {
	return &akov2generated.Group{
		TypeMeta:   metav1.TypeMeta{Kind: "Group", APIVersion: "atlas.generated.mongodb.com/v1"},
		ObjectMeta: metav1.ObjectMeta{Name: testGroupName, Namespace: testNamespace},
		Status: akov2generated.GroupStatus{
			V20250312: &akov2generated.GroupStatusV20250312{Id: id},
		},
	}
}

func testNetworkContainer() *akov2generated.NetworkContainer {
	return &akov2generated.NetworkContainer{
		TypeMeta:   metav1.TypeMeta{Kind: "NetworkContainer", APIVersion: "atlas.generated.mongodb.com/v1"},
		ObjectMeta: metav1.ObjectMeta{Name: "test-container", Namespace: testNamespace, Generation: 1},
		Spec: akov2generated.NetworkContainerSpec{
			V20250312: &akov2generated.NetworkContainerSpecV20250312{
				GroupRef: &k8s.LocalReference{Name: testGroupName},
				Entry: &akov2generated.NetworkContainerSpecV20250312Entry{
					ProviderName:   new("AWS"),
					RegionName:     new("US_EAST_1"),
					AtlasCidrBlock: new("10.128.0.0/21"),
				},
			},
		},
	}
}

func withStatusID(nc *akov2generated.NetworkContainer) *akov2generated.NetworkContainer {
	nc.Status.V20250312 = &akov2generated.NetworkContainerStatusV20250312{Id: new(testContainerID)}
	return nc
}

func atlasError(code string) error {
	apiErr := &v20250312sdk.GenericOpenAPIError{}
	apiErr.SetModel(v20250312sdk.ApiError{ErrorCode: code})
	return apiErr
}

func assertError(t *testing.T, err error, wantErr string) {
	t.Helper()
	if wantErr == "" {
		require.NoError(t, err)
		return
	}
	require.Error(t, err)
	assert.Contains(t, err.Error(), wantErr)
}
