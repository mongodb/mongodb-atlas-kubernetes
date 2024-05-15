package indexer

import (
	"testing"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

func Test_AtlasSearchIndexKeysToDeployment(t *testing.T) {
	t.Run("should return nil when AtlasSearchIndex is not referenced by a Deployment", func(t *testing.T) {
		instance := &akov2.AtlasDeployment{
			Spec: akov2.AtlasDeploymentSpec{
				DeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name: "test-deployment",
				},
			},
		}

		indexer := NewAtlasDeploymentBySearchIndexIndexer(zaptest.NewLogger(t))
		keys := indexer.Keys(instance)
		assert.Nil(t, keys)
	})

	t.Run("should return indexes slice AtlasSearchIndex is referenced by a Deployment", func(t *testing.T) {
		instance := &akov2.AtlasDeployment{
			Spec: akov2.AtlasDeploymentSpec{
				DeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name: "test-deployment",
					SearchIndexes: []akov2.SearchIndex{
						{
							Search: &akov2.Search{
								SearchConfigurationRef: common.ResourceRefNamespaced{
									Name:      "config-1",
									Namespace: "default",
								},
							},
						},
						{
							Search: &akov2.Search{
								SearchConfigurationRef: common.ResourceRefNamespaced{
									Name:      "config-2",
									Namespace: "default",
								},
							},
						},
					},
				},
			},
		}

		indexer := NewAtlasDeploymentBySearchIndexIndexer(zaptest.NewLogger(t))
		keys := indexer.Keys(instance)

		assert.Equal(
			t,
			[]string{
				"default/config-1",
				"default/config-2",
			},
			keys,
		)
	})
}
