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

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
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
		"should not return external project reference": {
			object: &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					ExternalProjectRef: &akov2.ExternalProjectReference{
						ID: "external-project-id",
					},
				},
			},
			expectedKeys: nil,
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
		"should return project name and the namespace": {
			object: &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					Project: &common.ResourceRefNamespaced{
						Name:      "testProject",
						Namespace: "testNamespace",
					},
				},
			},
			expectedLogs: []observer.LoggedEntry{},
			expectedKeys: []string{"testNamespace/testProject"},
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
			expectedKeys: []string{"ns/internal-project-id"},
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
			expectedKeys: []string{"ns/internal-project-id"},
			expectedLogs: []observer.LoggedEntry{},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))

			core, logs := observer.New(zap.DebugLevel)

			indexer := NewAtlasDatabaseUserByProjectIndexer(zap.New(core))
			keys := indexer.Keys(tt.object)
			sort.Strings(keys)

			assert.Equal(t, tt.expectedKeys, keys)
			assert.Equal(t, tt.expectedLogs, logs.AllUntimed())
		})
	}
}
