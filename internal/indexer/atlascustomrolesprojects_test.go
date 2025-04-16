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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
)

func TestAtlasCustomRoleByProjectsIndexer(t *testing.T) {
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
					Entry:   zapcore.Entry{LoggerName: AtlasCustomRoleByProject, Level: zap.ErrorLevel, Message: "expected a project.ProjectReferrerObject but got *v1.AtlasStreamInstance"},
				},
			},
		},
		"should return nil when there are no references": {
			object:       &akov2.AtlasCustomRole{},
			expectedLogs: []observer.LoggedEntry{},
		},
		"should return nil when there is an empty reference for external project": {
			object: &akov2.AtlasCustomRole{
				Spec: akov2.AtlasCustomRoleSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{},
					},
				},
			},
			expectedLogs: []observer.LoggedEntry{},
		},
		"should NOT return external project reference": {
			object: &akov2.AtlasCustomRole{
				Spec: akov2.AtlasCustomRoleSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: "external-project-id",
						},
					},
				},
			},
			expectedKeys: nil,
			expectedLogs: []observer.LoggedEntry{},
		},
		"should return nil when there is an empty reference for project": {
			object: &akov2.AtlasCustomRole{
				Spec: akov2.AtlasCustomRoleSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "",
						},
					},
				},
			},
			expectedLogs: []observer.LoggedEntry{},
		},
		"should return project with the same namespace as a custom role": {
			object: &akov2.AtlasCustomRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testRole",
					Namespace: "testNamespace",
				},
				Spec: akov2.AtlasCustomRoleSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "not-found-project",
						},
					},
				},
			},
			expectedLogs: []observer.LoggedEntry{},
			expectedKeys: []string{"testNamespace/not-found-project"},
		},
		"should return project reference with database customRole namespace": {
			object: &akov2.AtlasCustomRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "customRole",
					Namespace: "ns",
				},
				Spec: akov2.AtlasCustomRoleSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "internal-project-id",
						},
					},
				},
			},
			expectedKeys: []string{"ns/internal-project-id"},
			expectedLogs: []observer.LoggedEntry{},
		},
		"should return project reference": {
			object: &akov2.AtlasCustomRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "customRole",
					Namespace: "nsCustomRole",
				},
				Spec: akov2.AtlasCustomRoleSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      "internal-project-id",
							Namespace: "ns",
						},
					},
				},
			},
			expectedKeys: []string{"ns/internal-project-id"},
			expectedLogs: []observer.LoggedEntry{},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))

			core, logs := observer.New(zap.DebugLevel)

			indexer := NewAtlasCustomRoleByProjectIndexer(zap.New(core))
			keys := indexer.Keys(tt.object)
			sort.Strings(keys)

			assert.Equal(t, tt.expectedKeys, keys)
			assert.Equal(t, tt.expectedLogs, logs.AllUntimed())
		})
	}
}
