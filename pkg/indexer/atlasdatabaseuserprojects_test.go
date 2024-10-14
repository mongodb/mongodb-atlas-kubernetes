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

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
)

func TestAtlasDatabaseUserByProjectsIndexer(t *testing.T) {
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
					Entry:   zapcore.Entry{LoggerName: AtlasDatabaseUserByProject, Level: zap.ErrorLevel, Message: "expected *v1.AtlasDatabaseUser but got *v1.AtlasStreamInstance"},
				},
			},
		},
		"should return nil when there are no references": {
			object:       &akov2.AtlasDatabaseUser{},
			expectedLogs: []observer.LoggedEntry{},
		},
		"should return nil when there is an empty reference for external project": {
			object: &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					ExternalProjectRef: &akov2.ExternalProjectReference{},
				},
			},
			expectedLogs: []observer.LoggedEntry{},
		},
		"should return external project reference": {
			object: &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					ExternalProjectRef: &akov2.ExternalProjectReference{
						ID: "external-project-id",
					},
				},
			},
			expectedKeys: []string{"external-project-id"},
			expectedLogs: []observer.LoggedEntry{},
		},
		"should return nil when there is an empty reference for project": {
			object: &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					Project: &common.ResourceRefNamespaced{
						Name: "",
					},
				},
			},
			expectedLogs: []observer.LoggedEntry{},
		},
		"should return nil when referenced project was not found": {
			object: &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					Project: &common.ResourceRefNamespaced{
						Name: "not-found-project",
					},
				},
			},
			expectedLogs: []observer.LoggedEntry{
				{
					Context: []zapcore.Field{},
					Entry:   zapcore.Entry{LoggerName: AtlasDatabaseUserByProject, Level: zap.ErrorLevel, Message: "unable to find project to index: atlasprojects.atlas.mongodb.com \"not-found-project\" not found"},
				},
			},
		},
		"should return project reference with database user namespace": {
			object: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user",
					Namespace: "ns",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Project: &common.ResourceRefNamespaced{
						Name: "internal-project-id",
					},
				},
			},
			expectedKeys: []string{"external-project-id"},
			expectedLogs: []observer.LoggedEntry{},
		},
		"should return project reference": {
			object: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user",
					Namespace: "nsUser",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Project: &common.ResourceRefNamespaced{
						Name:      "internal-project-id",
						Namespace: "ns",
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
					Name:      "internal-project-id",
					Namespace: "ns",
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "My Project",
				},
				Status: status.AtlasProjectStatus{
					ID: "external-project-id",
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

			indexer := NewAtlasDatabaseUserByProjectIndexer(context.Background(), k8sClient, zap.New(core))
			keys := indexer.Keys(tt.object)
			sort.Strings(keys)

			assert.Equal(t, tt.expectedKeys, keys)
			assert.Equal(t, tt.expectedLogs, logs.AllUntimed())
		})
	}
}
