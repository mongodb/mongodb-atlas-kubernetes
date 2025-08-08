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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
)

func TestAtlasDeploymentBySpecNameIndexer(t *testing.T) {
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
					Entry: zapcore.Entry{
						LoggerName: AtlasDeploymentBySpecNameAndProjectID,
						Level:      zap.ErrorLevel,
						Message:    "expected *v1.AtlasDeployment but got *v1.AtlasStreamInstance",
					},
				},
			},
		},
		"should return nil when no name set": {
			object: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name: "",
					},
				},
			},
			expectedLogs: []observer.LoggedEntry{},
		},
		"should return nil if name exists but no project refs": {
			object: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name: "test-my-deployment",
					},
				},
			},
			expectedLogs: []observer.LoggedEntry{},
		},
		"should return key from ExternalProjectRef": {
			object: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name: "test-my-deployment",
					},
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: "test-external-id",
						},
					},
				},
			},
			expectedKeys: []string{"test-external-id-" + kube.NormalizeIdentifier("test-my-deployment")},
			expectedLogs: []observer.LoggedEntry{},
		},
		"should return key from resolved ProjectRef": {
			object: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-deployment",
					Namespace: "test-ns",
				},
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name: "test-my-deployment",
					},
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      "test-project",
							Namespace: "test-ns",
						},
					},
				},
			},
			expectedKeys: []string{"test-project-id-" + kube.NormalizeIdentifier("test-my-deployment")},
			expectedLogs: []observer.LoggedEntry{},
		},
		"should normalize deployment name before indexing": {
			object: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name: "Test.Deployment+123",
					},
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: "test-external-id",
						},
					},
				},
			},
			expectedKeys: []string{"test-external-id-" + kube.NormalizeIdentifier("Test.Deployment+123")},
			expectedLogs: []observer.LoggedEntry{},
		},
		"should log error if ProjectRef can't be resolved": {
			object: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-deployment",
					Namespace: "test-ns",
				},
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name: "test-unknown-deployment",
					},
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "nonexistent-project",
						},
					},
				},
			},
			expectedLogs: []observer.LoggedEntry{
				{
					Context: []zapcore.Field{},
					Entry: zapcore.Entry{
						LoggerName: AtlasDeploymentBySpecNameAndProjectID,
						Level:      zap.ErrorLevel,
						Message:    "unable to find project to index: atlasprojects.atlas.mongodb.com \"nonexistent-project\" not found",
					},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			project := &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-project",
					Namespace: "test-ns",
				},
				Status: status.AtlasProjectStatus{
					ID: "test-project-id",
				},
			}

			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(project).
				WithStatusSubresource(project).
				Build()

			core, logs := observer.New(zap.DebugLevel)

			indexer := NewAtlasDeploymentBySpecNameIndexer(context.Background(), k8sClient, zap.New(core))
			keys := indexer.Keys(tt.object)
			sort.Strings(keys)

			assert.Equal(t, tt.expectedKeys, keys)
			assert.Equal(t, tt.expectedLogs, logs.AllUntimed())
		})
	}
}
