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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
)

func TestAtlasProjectByTeamIndexer(t *testing.T) {
	for _, tc := range []struct {
		name     string
		object   client.Object
		wantKeys []string
	}{
		{
			name:     "should return nil on wrong type",
			object:   &akov2.AtlasDatabaseUser{},
			wantKeys: nil,
		},
		{
			name:     "should return no keys when there are no references",
			object:   &akov2.AtlasProject{},
			wantKeys: []string{},
		},
		{
			name: "should return no keys when there is an empty reference",
			object: &akov2.AtlasProject{
				Spec: akov2.AtlasProjectSpec{
					Teams: []akov2.Team{
						{TeamRef: common.ResourceRefNamespaced{}},
					},
				},
			},
			wantKeys: []string{},
		},
		{
			name: "should return keys when there is a reference",
			object: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
				},
				Spec: akov2.AtlasProjectSpec{
					Teams: []akov2.Team{
						{TeamRef: common.ResourceRefNamespaced{Name: "test-team"}},
						{TeamRef: common.ResourceRefNamespaced{Name: "test-team2", Namespace: "ns2"}},
					},
				},
			},
			wantKeys: []string{"bar/test-team", "ns2/test-team2"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			indexer := NewAtlasProjectByTeamIndexer(zaptest.NewLogger(t))
			keys := indexer.Keys(tc.object)
			sort.Strings(keys)
			assert.Equal(t, tc.wantKeys, keys)
		})
	}
}
