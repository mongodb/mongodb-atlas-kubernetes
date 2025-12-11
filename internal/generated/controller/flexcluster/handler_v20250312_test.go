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

package flexcluster

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/crd2go/crd2go/k8s"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v20250312sdk "go.mongodb.org/atlas-sdk/v20250312009/admin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"

	crapi "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/fakeatlas"
	crds "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/crds"
	akov2generated "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	result "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result"
	state "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

func TestHandleInitial(t *testing.T) {
	testCases := []struct {
		title                  string
		flexCluster            *akov2generated.FlexCluster
		kubeObjects            []client.Object
		want                   ctrlstate.Result
		wantErr                string
		atlasCreateClusterFunc func(ctx context.Context, params *v20250312sdk.CreateFlexClusterApiParams) (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error)
	}{
		{
			title:       "successful initial creation",
			flexCluster: setGroupRef(defaultTestFlexCluster("test-cluster", "default"), "test-group"),
			kubeObjects: []client.Object{
				defaultTestGroup("test-group", "default", pointer.MakePtr("62b6e34b3d91647abb20e7b8")),
			},
			want: ReenqueueResult(state.StateCreating, "Creating Flex Cluster."),
		},
		{
			title:       "missing group",
			flexCluster: setGroupRef(defaultTestFlexCluster("test-cluster1", "ns"), "not-found"),
			want:        ctrlstate.Result{NextState: state.StateInitial},
			wantErr:     "failed to get dependencies: failed to get group",
		},
		{
			title:       "group without id",
			flexCluster: setGroupRef(defaultTestFlexCluster("test-cluster", "ns1"), "test-group"),
			kubeObjects: []client.Object{
				defaultTestGroup("test-group", "ns1", nil),
			},
			want:    ctrlstate.Result{NextState: state.StateInitial},
			wantErr: "failed to fetch referenced value groupId",
		},
		{
			title:       "atlas API error",
			flexCluster: setGroupRef(defaultTestFlexCluster("test-cluster2", "ns3"), "test-group"),
			kubeObjects: []client.Object{
				defaultTestGroup("test-group", "ns3", pointer.MakePtr("62b6e34b3d91647abb20e7b8")),
			},
			atlasCreateClusterFunc: func(ctx context.Context, params *v20250312sdk.CreateFlexClusterApiParams) (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error) {
				return nil, nil, fmt.Errorf("atlas API error: cluster creation failed")
			},
			want:    ctrlstate.Result{NextState: state.StateInitial},
			wantErr: "failed to create flex cluster",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			ctx := context.Background()
			scheme := createTestScheme(t)

			allObjects := append([]client.Object{tc.flexCluster}, tc.kubeObjects...)
			kubeClient := newFakeClientWithGVK(scheme, allObjects, tc.flexCluster)

			crd, err := crds.EmbeddedCRD("FlexCluster")
			require.NoError(t, err)
			tr, err := crapi.NewTranslator(crd, "v1", "v20250312")
			require.NoError(t, err)

			atlasClientBuilder := fakeatlas.NewAtlasClientBuilder().
				WithFakeFlexClusterClient()

			atlasClient := atlasClientBuilder.Build()

			// Configure fake client with custom function if provided
			if tc.atlasCreateClusterFunc != nil {
				if fakeFlexClient, ok := atlasClient.FlexClustersApi.(*fakeatlas.FakeFlexClustersApi); ok {
					fakeFlexClient.CreateFlexClusterWithParamsFunc = tc.atlasCreateClusterFunc
				}
			}

			deletionProtection := false

			handler := NewHandlerv20250312(kubeClient, atlasClient, tr, deletionProtection)

			result, err := handler.HandleInitial(ctx, tc.flexCluster)

			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.want, result)
		})
	}
}

// createTestScheme creates a runtime scheme for testing
func createTestScheme(t *testing.T) *runtime.Scheme {
	scheme := runtime.NewScheme()
	require.NoError(t, akov2generated.AddToScheme(scheme))
	require.NoError(t, corev1.AddToScheme(scheme))
	return scheme
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

func setGroupRef(flexCluster *akov2generated.FlexCluster, groupRef string) *akov2generated.FlexCluster {
	flexCluster.Spec.V20250312.GroupRef = &k8s.LocalReference{Name: groupRef}
	return flexCluster
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

func ReenqueueResult(state state.ResourceState, msg string) ctrlstate.Result {
	return ctrlstate.Result{
		Result:    reconcile.Result{RequeueAfter: result.DefaultRequeueTIme},
		NextState: state,
		StateMsg:  msg,
	}
}
