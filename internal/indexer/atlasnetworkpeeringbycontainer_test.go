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

package indexer

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
)

func TestAtlasNetworkPeeringByContainerIndexer(t *testing.T) {
	for _, tc := range []struct {
		title    string
		object   client.Object
		wantKeys []string
	}{
		{
			title: "nil obj renders nothing",
		},
		{
			title: "wrong obj renders nothing",
			object: &akov2.AtlasDataFederation{
				Spec: akov2.DataFederationSpec{
					Project: common.ResourceRefNamespaced{},
				},
			},
		},
		{
			title: "peering with id ref renders nothing",
			object: &akov2.AtlasNetworkPeering{
				Spec: akov2.AtlasNetworkPeeringSpec{
					ContainerRef: akov2.ContainerDualReference{ID: "an-id"},
				},
			},
		},
		{
			title: "peering with name ref renders container name",
			object: &akov2.AtlasNetworkPeering{
				Spec: akov2.AtlasNetworkPeeringSpec{
					ContainerRef: akov2.ContainerDualReference{Name: "a-name"},
				},
			},
			wantKeys: []string{"a-name"},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			indexer := NewAtlasNetworkPeeringByContainerIndexer(zaptest.NewLogger(t))
			keys := indexer.Keys(tc.object)
			sort.Strings(keys)
			assert.Equal(t, tc.wantKeys, keys)
		})
	}
}
