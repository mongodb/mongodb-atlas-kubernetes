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
	"context"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
)

func TestAtlasDataFederationByProjectIDIndexer(t *testing.T) {
	tests := map[string]struct {
		object       client.Object
		expectedKeys []string
		expectedLogs []observer.LoggedEntry
	}{
		"should return nil on wrong type": {
			object: &akov2.AtlasStreamInstance{},
			expectedLogs: []observer.LoggedEntry{
				{
					Context: []zapcore.Field{},
					Entry:   zapcore.Entry{LoggerName: AtlasDataFederationByProjectID, Level: zap.ErrorLevel, Message: "expected *v1.AtlasDataFederation but got *v1.AtlasStreamInstance"},
				},
			},
		},
		"should return nil when there is an empty project reference": {
			object:       &akov2.AtlasDataFederation{},
			expectedLogs: []observer.LoggedEntry{},
		},
		"should return nil when referenced project was not found": {
			object: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "df",
					Namespace: "ns",
				},
				Spec: akov2.DataFederationSpec{
					Project: common.ResourceRefNamespaced{
						Name: "not-found-project",
					},
				},
			},
			expectedLogs: []observer.LoggedEntry{
				{
					Context: []zapcore.Field{},
					Entry:   zapcore.Entry{LoggerName: AtlasDataFederationByProjectID, Level: zap.ErrorLevel, Message: "unable to find project to index: atlasprojects.atlas.mongodb.com \"not-found-project\" not found"},
				},
			},
		},
		"should return project ID using DataFederation namespace": {
			object: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "df",
					Namespace: "ns",
				},
				Spec: akov2.DataFederationSpec{
					Project: common.ResourceRefNamespaced{
						Name: "internal-project",
					},
				},
			},
			expectedKeys: []string{"external-project-id"},
			expectedLogs: []observer.LoggedEntry{},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			project := &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "internal-project",
					Namespace: "ns",
				},
				Status: status.AtlasProjectStatus{
					ID: "external-project-id",
				},
			}

			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))

			builder := fake.NewClientBuilder().WithScheme(testScheme)
			k8sClient := builder.WithObjects(project).Build()
			core, logs := observer.New(zap.DebugLevel)

			indexer := NewAtlasDataFederationByProjectIDIndexer(context.Background(), k8sClient, zap.New(core))
			keys := indexer.Keys(tt.object)
			sort.Strings(keys)

			assert.Equal(t, tt.expectedKeys, keys)
			assert.Equal(t, tt.expectedLogs, logs.AllUntimed())
		})
	}
}
