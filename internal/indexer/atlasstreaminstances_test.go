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

func TestAtlasStreamInstancesByConnectionRegistryIndices(t *testing.T) {
	t.Run("should return nil when instance has no connection", func(t *testing.T) {
		instance := &akov2.AtlasStreamInstance{
			Spec: akov2.AtlasStreamInstanceSpec{
				Name: "instance-0",
			},
		}

		indexer := NewAtlasStreamInstanceByConnectionIndexer(zaptest.NewLogger(t))
		keys := indexer.Keys(instance)
		assert.Nil(t, keys)
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

		indexer := NewAtlasStreamInstanceByConnectionIndexer(zaptest.NewLogger(t))
		keys := indexer.Keys(instance)
		assert.Equal(
			t,
			[]string{
				"default/conn-1",
				"default/conn-2",
			},
			keys,
		)
	})
}

func TestAtlasStreamInstancesByProjectIndices(t *testing.T) {
	t.Run("should return nil when instance has no project associated to it", func(t *testing.T) {
		instance := &akov2.AtlasStreamInstance{
			Spec: akov2.AtlasStreamInstanceSpec{
				Name: "instance-0",
			},
		}

		indexer := NewAtlasStreamInstanceByProjectIndexer(zaptest.NewLogger(t))
		keys := indexer.Keys(instance)
		assert.Nil(t, keys)
	})

	t.Run("should return indexes slice when instance has project associated to it", func(t *testing.T) {
		instance := &akov2.AtlasStreamInstance{
			Spec: akov2.AtlasStreamInstanceSpec{
				Name: "instance-0",
				Project: common.ResourceRefNamespaced{
					Name:      "project-1",
					Namespace: "default",
				},
			},
		}

		indexer := NewAtlasStreamInstanceByProjectIndexer(zaptest.NewLogger(t))
		keys := indexer.Keys(instance)
		assert.Equal(
			t,
			[]string{
				"default/project-1",
			},
			keys,
		)
	})
}
