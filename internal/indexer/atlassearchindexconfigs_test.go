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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
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
