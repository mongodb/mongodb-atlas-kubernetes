package indexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
)

func TestAtlasPrivateEndpointsByProjectIndices(t *testing.T) {
	t.Run("should return nil when instance has no project associated to it", func(t *testing.T) {
		pe := &akov2.AtlasPrivateEndpoint{
			Spec: akov2.AtlasPrivateEndpointSpec{
				Provider: "AWS",
				Region:   "US_EAST_1",
			},
		}

		indexer := NewAtlasPrivateEndpointByProjectIndexer(zaptest.NewLogger(t))
		keys := indexer.Keys(pe)
		assert.Nil(t, keys)
	})

	t.Run("should return indexes slice when instance has project associated to it", func(t *testing.T) {
		pe := &akov2.AtlasPrivateEndpoint{
			Spec: akov2.AtlasPrivateEndpointSpec{
				Provider: "AWS",
				Region:   "US_EAST_1",
				Project: &common.ResourceRefNamespaced{
					Name:      "project-1",
					Namespace: "default",
				},
			},
		}

		indexer := NewAtlasPrivateEndpointByProjectIndexer(zaptest.NewLogger(t))
		keys := indexer.Keys(pe)
		assert.Equal(
			t,
			[]string{
				"default/project-1",
			},
			keys,
		)
	})
}
