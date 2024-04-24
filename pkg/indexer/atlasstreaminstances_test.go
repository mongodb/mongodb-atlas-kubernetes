package indexer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
)

func TestAtlasStreamInstancesByConnectionRegistryIndices(t *testing.T) {
	t.Run("should return nil when instance has no connection", func(t *testing.T) {
		instance := &akov2.AtlasStreamInstance{
			Spec: akov2.AtlasStreamInstanceSpec{
				Name: "instance-0",
			},
		}

		indexes := AtlasStreamInstancesByConnectionRegistryIndices(instance)
		assert.Nil(t, indexes)
	})

	t.Run("should return indexes slice when instance has connections", func(t *testing.T) {
		instance := &akov2.AtlasStreamInstance{
			Spec: akov2.AtlasStreamInstanceSpec{
				Name: "instance-0",
				ConnectionRegistry: []common.ResourceRefNamespaced{
					{
						Name:      "conn-1",
						Namespace: "default",
					},
					{
						Name:      "conn-2",
						Namespace: "default",
					},
				},
			},
		}

		indexes := AtlasStreamInstancesByConnectionRegistryIndices(instance)
		assert.Equal(
			t,
			[]string{
				"default/conn-1",
				"default/conn-2",
			},
			indexes,
		)
	})
}
