//nolint:dupl
package indexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
)

func TestAtlasDeploymentByBackupScheduleIndexer(t *testing.T) {
	for _, tc := range []struct {
		name     string
		object   client.Object
		wantKeys []string
	}{
		{
			name:     "should return nil on wrong type",
			object:   &akov2.AtlasProject{},
			wantKeys: nil,
		},
		{
			name:     "should return nil when there are no references",
			object:   &akov2.AtlasDeployment{},
			wantKeys: nil,
		},
		{
			name: "should return nil when there is an empty reference",
			object: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					BackupScheduleRef: common.ResourceRefNamespaced{},
				},
			},
			wantKeys: nil,
		},
		{
			name: "should return a key when there is a reference",
			object: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
				},
				Spec: akov2.AtlasDeploymentSpec{
					BackupScheduleRef: common.ResourceRefNamespaced{
						Name: "baz",
					},
				},
			},
			wantKeys: []string{"bar/baz"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			indexer := NewAtlasDeploymentByBackupScheduleIndexer(zaptest.NewLogger(t))
			keys := indexer.Keys(tc.object)
			assert.Equal(t, tc.wantKeys, keys)
		})
	}
}
