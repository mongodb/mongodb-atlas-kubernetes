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
	"net/http"
	"testing"

	crapi "github.com/crd2go/crapi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	v20250312sdk "go.mongodb.org/atlas-sdk/v20250312027/admin"
	"go.mongodb.org/atlas-sdk/v20250312027/mockadmin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	akov2generated "github.com/mongodb/mongodb-atlas-kubernetes/v2/generated/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/flexcluster"
	crds "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/crds"
)

const (
	identityNamespace   = "identity-ns"
	identityGroupID     = "62b6e34b3d91647abb20e7b8"
	identityClusterName = "my-flex-cluster"
	identityClusterID   = "671e2a3b4c5d6e7f8a9b0c1d"
	identityVictimName  = "someone-elses-flex-cluster"
)

// TestHandleDeletionRequestedVerifiesRecordedIdentity covers SECBUG-2945: Atlas addresses
// flex clusters by name only, so renaming spec.v20250312.entry.name after creation used to
// aim the delete call at whatever cluster currently bore that name.
func TestHandleDeletionRequestedVerifiesRecordedIdentity(t *testing.T) {
	for _, tc := range []struct {
		title          string
		specName       string
		expectDelete   bool
		wantErrMessage string
	}{
		{
			title:        "deletes when the spec name matches the recorded name",
			specName:     identityClusterName,
			expectDelete: true,
		},
		{
			title:          "refuses to delete after the spec name was renamed",
			specName:       identityVictimName,
			expectDelete:   false,
			wantErrMessage: "flex cluster names are immutable after creation",
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			flexCluster := identityTestFlexCluster(tc.specName)

			flexAPI := mockadmin.NewFlexClustersAPI(t)
			if tc.expectDelete {
				delReq := v20250312sdk.DeleteFlexClusterApiRequest{ApiService: flexAPI}
				flexAPI.EXPECT().DeleteFlexClusterWithParams(mock.Anything, mock.Anything).Return(delReq)
				flexAPI.EXPECT().DeleteFlexClusterExecute(delReq).Return(&http.Response{}, nil)
			}

			handler := identityTestHandler(t, flexAPI)
			_, err := handler.HandleDeletionRequested(context.Background(), flexCluster)

			if tc.wantErrMessage == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, tc.wantErrMessage)
			flexAPI.AssertNotCalled(t, "DeleteFlexClusterWithParams", mock.Anything, mock.Anything)
		})
	}
}

func identityTestFlexCluster(specName string) *akov2generated.FlexCluster {
	groupID := identityGroupID
	recordedName := identityClusterName
	recordedID := identityClusterID

	return &akov2generated.FlexCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "FlexCluster",
			APIVersion: "atlas.generated.mongodb.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      identityClusterName,
			Namespace: identityNamespace,
		},
		Spec: akov2generated.FlexClusterSpec{
			V20250312: &akov2generated.FlexClusterSpecV20250312{
				GroupId: &groupID,
				Entry: &akov2generated.FlexClusterSpecV20250312Entry{
					Name: specName,
					ProviderSettings: akov2generated.ProviderSettings{
						BackingProviderName: "AWS",
						RegionName:          "us-east-1",
					},
				},
			},
		},
		Status: akov2generated.FlexClusterStatus{
			V20250312: &akov2generated.FlexClusterStatusV20250312{
				Id:      &recordedID,
				GroupId: &groupID,
				Name:    &recordedName,
			},
		},
	}
}

func identityTestHandler(t *testing.T, flexAPI *mockadmin.FlexClustersAPI) *flexcluster.Handlerv20250312 {
	t.Helper()

	scheme := runtime.NewScheme()
	require.NoError(t, akov2generated.AddToScheme(scheme))

	crd, err := crds.EmbeddedCRD("FlexCluster")
	require.NoError(t, err)

	translator, err := crapi.NewTranslator(scheme, crd, "v1", "v20250312")
	require.NoError(t, err)

	return flexcluster.NewHandlerv20250312(
		fake.NewClientBuilder().WithScheme(scheme).Build(),
		&v20250312sdk.APIClient{FlexClustersAPI: flexAPI},
		translator,
		false, // deletion protection off, otherwise the delete path is never reached
	)
}
