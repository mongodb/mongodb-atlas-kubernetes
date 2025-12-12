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
	"go.mongodb.org/atlas-sdk/v20250312009/admin"
	v20250312sdk "go.mongodb.org/atlas-sdk/v20250312009/admin"
	"go.mongodb.org/atlas-sdk/v20250312009/mockadmin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
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

func TestHandleInitial(t *testing.T) {
	scheme := createTestScheme(t)
	tr, err := translatorFromEnbeddedCRD(scheme, "FlexCluster", "v1", "v20250312")
	require.NoError(t, err)

	deletionProtection := false

	for _, tc := range []struct {
		title                      string
		flexCluster                *akov2generated.FlexCluster
		kubeObjects                []client.Object
		atlasCreateFlexClusterFunc func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error)
		want                       ctrlstate.Result
		wantErr                    string
	}{
		{
			title:       "create flex cluster",
			flexCluster: defaultTestFlexCluster("test-cluster", "ns1"),
			kubeObjects: []client.Object{
				defaultTestGroup("test-group", "ns1", nil),
			},
			atlasCreateFlexClusterFunc: func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error) {
				return &v20250312sdk.FlexClusterDescription20241113{Id: pointer.MakePtr("test-cluster")}, nil, nil
			},
			want: reenqueueResult(state.StateCreating, "Creating Flex Cluster."),
		},
		{
			title:       "corrupt flex type cluster",
			flexCluster: withFlexGVK(defaultTestFlexCluster("test-cluster1", "ns"), "bad-kind", "bad-api-version"),
			kubeObjects: []client.Object{
				defaultTestGroup("some-group", "ns1", nil),
			},
			want:    ctrlstate.Result{NextState: state.StateInitial},
			wantErr: "failed to get dependencies: nil flex cluster",
		},
		{
			title:       "corrupt api reply",
			flexCluster: defaultTestFlexCluster("test-cluster", "ns1"),
			kubeObjects: []client.Object{
				defaultTestGroup("test-group", "ns1", nil),
			},
			atlasCreateFlexClusterFunc: func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error) {
				return nil, nil, nil
			},
			want:    ctrlstate.Result{NextState: state.StateInitial},
			wantErr: "failed to get dependencies: nil flex cluster",
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
			atlasCreateFlexClusterFunc: func() (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error) {
				return nil, nil, fmt.Errorf("atlas API error: cluster creation failed")
			},
			want:    ctrlstate.Result{NextState: state.StateInitial},
			wantErr: "failed to create flex cluster",
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			allObjects := tc.kubeObjects
			clientBuilder := fake.NewClientBuilder().WithScheme(scheme)
			if tc.flexCluster != nil {
				allObjects = append([]client.Object{tc.flexCluster}, tc.kubeObjects...)
				clientBuilder = clientBuilder.WithObjects(allObjects...).WithStatusSubresource(tc.flexCluster)
			} else {
				clientBuilder = clientBuilder.WithObjects(allObjects...)
			}
			fakeClient := clientBuilder.Build()

			flexAPI := mockadmin.NewFlexClustersApi(t)
			if tc.atlasCreateFlexClusterFunc != nil {
				cluster, rsp, err := tc.atlasCreateFlexClusterFunc()
				req := admin.CreateFlexClusterApiRequest{ApiService: flexAPI}
				flexAPI.EXPECT().CreateFlexClusterWithParams(mock.Anything, mock.Anything).Return(req)
				flexAPI.EXPECT().CreateFlexClusterExecute(req).Return(cluster, rsp, err)
			}

			atlasClient := &v20250312sdk.APIClient{
				FlexClustersApi: flexAPI,
			}

			handler := flexcluster.NewHandlerv20250312(fakeClient, atlasClient, tr, deletionProtection)

			ctx := context.Background()
			got, err := handler.HandleInitial(ctx, tc.flexCluster)
			if tc.wantErr != "" {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.want, got)
		})
	}
}

func setGroupRef(flexCluster *akov2generated.FlexCluster, groupRef string) *akov2generated.FlexCluster {
	flexCluster.Spec.V20250312.GroupRef = &k8s.LocalReference{Name: groupRef}
	return flexCluster
}

func translatorFromEnbeddedCRD(scheme *runtime.Scheme, kind, apiVersion, majorVersion string) (crapi.Translator, error) {
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
