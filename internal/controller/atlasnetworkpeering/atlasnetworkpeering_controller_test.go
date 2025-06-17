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

package atlasnetworkpeering

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
)

func TestReconcile(t *testing.T) {
	ctx := context.Background()

	testScheme := runtime.NewScheme()
	require.NoError(t, akov2.AddToScheme(testScheme))

	tests := map[string]struct {
		request        reconcile.Request
		expectedResult reconcile.Result
		expectedLogs   []string
	}{
		"failed to prepare resource": {
			request:        reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "np0"}},
			expectedResult: reconcile.Result{},
			expectedLogs: []string{
				"-> Starting AtlasNetworkPeering reconciliation",
				"Object default/np0 doesn't exist, was it deleted after reconcile request?",
			},
		},
		"prepare resource for reconciliation": {
			request:        reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "np1"}},
			expectedResult: reconcile.Result{},
			expectedLogs: []string{
				"-> Starting AtlasNetworkPeering reconciliation",
				"-> Skipping AtlasNetworkPeering reconciliation as annotation mongodb.com/atlas-reconciliation-policy=skip",
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			core, logs := observer.New(zap.DebugLevel)
			fakeClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(testNetworkProvider()).
				Build()
			r := &AtlasNetworkPeeringReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: fakeClient,
					Log:    zap.New(core).Sugar(),
				},
			}
			result, _ := r.Reconcile(ctx, tc.request)
			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, len(tc.expectedLogs), logs.Len())
			for i, log := range logs.All() {
				assert.Equal(t, tc.expectedLogs[i], log.Message)
			}
		})
	}
}

func testNetworkProvider() *akov2.AtlasNetworkPeering {
	return &akov2.AtlasNetworkPeering{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "np1",
			Namespace: "default",
			Annotations: map[string]string{
				customresource.ReconciliationPolicyAnnotation: customresource.ReconciliationPolicySkip,
			},
		},
		Spec: akov2.AtlasNetworkPeeringSpec{
			ProjectDualReference: akov2.ProjectDualReference{
				ExternalProjectRef: &akov2.ExternalProjectReference{
					ID: "fake-project-id",
				},
				ConnectionSecret: &api.LocalObjectReference{
					Name: "fake-secret",
				},
			},
			ContainerRef: akov2.ContainerDualReference{
				Name: "fake-container-id",
			},
			AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
				Provider: "AWS",
				AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
					AccepterRegionName:  "us-east-1",
					AWSAccountID:        "some-aws-id",
					RouteTableCIDRBlock: "10.0.0.0/8",
					VpcID:               "vpc-id-test",
				},
			},
		},
		Status: status.AtlasNetworkPeeringStatus{
			ID: "peering-id",
		},
	}
}
