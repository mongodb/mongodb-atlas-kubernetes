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

//nolint:dupl
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

func TestAtlasDataFederationByProjectIndexer(t *testing.T) {
	for _, tc := range []struct {
		name     string
		object   client.Object
		wantKeys []string
	}{
		{
			name:   "should return nil on wrong type",
			object: &akov2.AtlasDeployment{},
		},
		{
			name:   "should return nil when there are no references",
			object: &akov2.AtlasDataFederation{},
		},
		{
			name: "should return nil when there is an empty reference",
			object: &akov2.AtlasDataFederation{
				Spec: akov2.DataFederationSpec{
					Project: common.ResourceRefNamespaced{},
				},
			},
		},
		{
			name: "should return project namespace if name is set only",
			object: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{Name: "name", Namespace: "ns"},
				Spec: akov2.DataFederationSpec{
					Project: common.ResourceRefNamespaced{Name: "someProject"},
				},
			},
			wantKeys: []string{"ns/someProject"},
		},
		{
			name: "should return secret namespace and name if set",
			object: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{Name: "name", Namespace: "ns"},
				Spec: akov2.DataFederationSpec{
					Project: common.ResourceRefNamespaced{Name: "someProject", Namespace: "otherNs"},
				},
			},
			wantKeys: []string{"otherNs/someProject"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			indexer := NewAtlasDataFederationByProjectIndexer(zaptest.NewLogger(t))
			keys := indexer.Keys(tc.object)
			sort.Strings(keys)
			assert.Equal(t, tc.wantKeys, keys)
		})
	}
}
