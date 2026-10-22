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

package cluster_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/crd2go/constate"
	crapi "github.com/crd2go/crapi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	v20250312sdk "go.mongodb.org/atlas-sdk/v20250312027/admin"
	"go.mongodb.org/atlas-sdk/v20250312027/mockadmin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	akov2generated "github.com/mongodb/mongodb-atlas-kubernetes/v2/generated/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	akocluster "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/cluster"
	crds "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/crds"
)

const (
	testNamespace   = "test-namespace"
	testGroupID     = "5f4c1b2a3d4e5f6a7b8c9d0e"
	testClusterName = "my-cluster"
	testClusterID   = "671e2a3b4c5d6e7f8a9b0c1d"
	victimClusterID = "aaaaaaaaaaaaaaaaaaaaaaaa"
)

// TestHandleDeletionRequestedVerifiesRecordedIdentity covers SECBUG-2945: Atlas addresses
// clusters by name only, so renaming spec.v20250312.entry.name after creation used to aim
// the delete call at whatever cluster currently bore that name.
func TestHandleDeletionRequestedVerifiesRecordedIdentity(t *testing.T) {
	for _, tc := range []struct {
		title          string
		resolvedID     string
		expectDelete   bool
		wantErrMessage string
	}{
		{
			title:        "deletes when the spec name still resolves to the recorded cluster",
			resolvedID:   testClusterID,
			expectDelete: true,
		},
		{
			title:          "refuses to delete when the spec name resolves to a foreign cluster",
			resolvedID:     victimClusterID,
			expectDelete:   false,
			wantErrMessage: "cluster names are immutable after creation",
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			cluster := testCluster()

			clustersAPI := mockadmin.NewClustersAPI(t)
			getReq := v20250312sdk.GetClusterApiRequest{ApiService: clustersAPI}
			clustersAPI.EXPECT().GetClusterWithParams(mock.Anything, mock.Anything).Return(getReq)
			clustersAPI.EXPECT().GetClusterExecute(getReq).Return(
				&v20250312sdk.ClusterDescription20240805{Id: &tc.resolvedID}, nil, nil,
			)

			if tc.expectDelete {
				delReq := v20250312sdk.DeleteClusterApiRequest{ApiService: clustersAPI}
				clustersAPI.EXPECT().DeleteClusterWithParams(mock.Anything, mock.Anything).Return(delReq)
				clustersAPI.EXPECT().DeleteClusterExecute(delReq).Return(&http.Response{}, nil)
			}

			handler := testHandler(t, clustersAPI)
			_, err := handler.HandleDeletionRequested(context.Background(), cluster)

			if tc.wantErrMessage == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, tc.wantErrMessage)
		})
	}
}

// TestHandleDeletionRequestedRefusesVanishedName covers the gap josvazg found in the first
// cut of the guard: a spec name that resolves to nothing used to be waved through, after
// which the delete call reported CLUSTER_NOT_FOUND and the resource went to StateDeleted,
// leaving the cluster recorded in status running in Atlas with nothing managing it.
func TestHandleDeletionRequestedRefusesVanishedName(t *testing.T) {
	cluster := testCluster()

	clustersAPI := mockadmin.NewClustersAPI(t)
	getReq := v20250312sdk.GetClusterApiRequest{ApiService: clustersAPI}
	clustersAPI.EXPECT().GetClusterWithParams(mock.Anything, mock.Anything).Return(getReq)
	clustersAPI.EXPECT().GetClusterExecute(getReq).Return(nil, nil, clusterNotFoundError())

	handler := testHandler(t, clustersAPI)
	_, err := handler.HandleDeletionRequested(context.Background(), cluster)

	require.ErrorContains(t, err, "does not exist in Atlas")
	require.ErrorContains(t, err, testClusterID)
	// The way out has to be discoverable from the message itself, otherwise the only
	// remaining move is editing finalizers by hand.
	require.ErrorContains(t, err, customresource.ResourcePolicyAnnotation)
	clustersAPI.AssertNotCalled(t, "DeleteClusterWithParams", mock.Anything, mock.Anything)
}

// TestHandleDeletionRequestedHonorsKeepPolicyOverVanishedName pins the escape hatch that
// makes the strict check above acceptable: a cluster removed outside the operator is
// released by annotation, without reaching the identity check at all.
func TestHandleDeletionRequestedHonorsKeepPolicyOverVanishedName(t *testing.T) {
	cluster := testCluster()
	cluster.SetAnnotations(map[string]string{
		customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep,
	})

	clustersAPI := mockadmin.NewClustersAPI(t)
	handler := testHandler(t, clustersAPI)

	result, err := handler.HandleDeletionRequested(context.Background(), cluster)
	require.NoError(t, err)
	require.Equal(t, constate.StateDeleted, result.NextState)
	clustersAPI.AssertNotCalled(t, "GetClusterWithParams", mock.Anything, mock.Anything)
}

// TestHandleDeletionRequestedSkipsVerificationWhenUnmanaged makes sure a resource that never
// reached Atlas does not pay for an extra lookup.
func TestHandleDeletionRequestedSkipsVerificationWhenUnmanaged(t *testing.T) {
	cluster := testCluster()
	cluster.Status.V20250312 = nil

	clustersAPI := mockadmin.NewClustersAPI(t)
	handler := testHandler(t, clustersAPI)

	_, err := handler.HandleDeletionRequested(context.Background(), cluster)
	require.NoError(t, err)
	clustersAPI.AssertNotCalled(t, "GetClusterWithParams", mock.Anything, mock.Anything)
}

func clusterNotFoundError() error {
	apiErr := &v20250312sdk.GenericOpenAPIError{}
	apiErr.SetModel(v20250312sdk.ApiError{ErrorCode: "CLUSTER_NOT_FOUND"})
	return apiErr
}

func testCluster() *akov2generated.Cluster {
	id := testClusterID
	groupID := testGroupID
	name := testClusterName

	return &akov2generated.Cluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Cluster",
			APIVersion: "atlas.generated.mongodb.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      testClusterName,
			Namespace: testNamespace,
		},
		Spec: akov2generated.ClusterSpec{
			V20250312: &akov2generated.ClusterSpecV20250312{
				GroupId: &groupID,
				Entry: &akov2generated.V20250312Entry{
					Name: &name,
				},
			},
		},
		Status: akov2generated.ClusterStatus{
			V20250312: &akov2generated.ClusterStatusV20250312{
				Id:      &id,
				GroupId: &groupID,
			},
		},
	}
}

func testHandler(t *testing.T, clustersAPI *mockadmin.ClustersAPI) *akocluster.Handlerv20250312 {
	t.Helper()

	scheme := runtime.NewScheme()
	require.NoError(t, akov2generated.AddToScheme(scheme))

	crd, err := crds.EmbeddedCRD("Cluster")
	require.NoError(t, err)

	translator, err := crapi.NewTranslator(scheme, crd, "v1", "v20250312")
	require.NoError(t, err)

	return akocluster.NewHandlerv20250312(
		fake.NewClientBuilder().WithScheme(scheme).Build(),
		&v20250312sdk.APIClient{ClustersAPI: clustersAPI},
		translator,
		false, // deletion protection off, otherwise the delete path is never reached
	)
}
