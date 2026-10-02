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

package group_test

import (
	"context"
	"testing"

	"github.com/crd2go/constate"
	crd2gok8s "github.com/crd2go/crd2go/k8s"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	akov2generated "github.com/mongodb/mongodb-atlas-kubernetes/v2/generated/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/group"
	indexer "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/indexers"
)

func TestHandlerv20250312_HandleDeletionRequested(t *testing.T) {
	ctx := context.Background()

	tests := map[string]struct {
		dependentCluster *akov2generated.Cluster
		noIndexes        bool
		wantState        constate.ResourceState
		wantMsgContains  string
		wantErr          string
	}{
		"Deletion is blocked when a Cluster depends on the group": {
			dependentCluster: dependentCluster("test-group"),
			wantState:        constate.StateDeletionRequested,
			wantMsgContains:  "failed to delete group because 1 resources depend on it.",
		},
		"Deletion proceeds when there are no dependent resources": {
			wantState:       constate.StateDeleted,
			wantMsgContains: "Group deleted.",
		},
		"Listing dependents fails when the field indexes are not registered": {
			wantErr: "failed to list",
			// The fake client errors on a field-selector List when the index
			// is not registered, so the index is intentionally omitted here.
			noIndexes: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			grp := &akov2generated.Group{
				ObjectMeta: metav1.ObjectMeta{Name: "test-group", Namespace: "default"},
			}
			scheme := runtime.NewScheme()
			require.NoError(t, akov2generated.AddToScheme(scheme))
			builder := fake.NewClientBuilder().WithScheme(scheme)
			if !test.noIndexes {
				logger := zaptest.NewLogger(t)
				builder = builder.
					WithIndex(&akov2generated.Cluster{}, indexer.ClusterByGroupIndex, indexer.NewClusterByGroupIndexer(logger).Keys).
					WithIndex(&akov2generated.DatabaseUser{}, indexer.DatabaseUserByGroupIndex, indexer.NewDatabaseUserByGroupIndexer(logger).Keys).
					WithIndex(&akov2generated.FlexCluster{}, indexer.FlexClusterByGroupIndex, indexer.NewFlexClusterByGroupIndexer(logger).Keys).
					WithIndex(&akov2generated.IPAccessListEntry{}, indexer.IPAccessListEntryByGroupIndex, indexer.NewIPAccessListEntryByGroupIndexer(logger).Keys)
			}
			objs := []client.Object{}
			if test.dependentCluster != nil {
				objs = append(objs, test.dependentCluster)
			}
			kubeClient := builder.WithObjects(objs...).Build()
			handler := group.NewHandlerv20250312(kubeClient, nil, nil, false)

			result, err := handler.HandleDeletionRequested(ctx, grp)

			if test.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), test.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.wantState, result.NextState)
			assert.Contains(t, result.StateMsg, test.wantMsgContains)
		})
	}
}

func dependentCluster(groupName string) *akov2generated.Cluster {
	return &akov2generated.Cluster{
		ObjectMeta: metav1.ObjectMeta{Name: "test-cluster", Namespace: "default"},
		Spec: akov2generated.ClusterSpec{
			V20250312: &akov2generated.ClusterSpecV20250312{
				GroupRef: &crd2gok8s.LocalReference{Name: groupName},
			},
		},
	}
}
