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
