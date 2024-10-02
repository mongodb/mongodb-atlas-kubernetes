package indexer

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

func TestAtlasDatabaseUserByExternalProjectsIndexer(t *testing.T) {
	for _, tc := range []struct {
		name     string
		object   client.Object
		wantKeys []string
	}{
		{
			name:   "should return nil on wrong type",
			object: &akov2.AtlasProject{},
		},
		{
			name:   "should return nil when there are no references",
			object: &akov2.AtlasDatabaseUser{},
		},
		{
			name: "should return nil when there is an empty reference for external project",
			object: &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					ExternalProjectRef: &akov2.ExternalProjectReference{},
				},
			},
		},
		{
			name: "should return external project reference",
			object: &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					ExternalProjectRef: &akov2.ExternalProjectReference{
						ID: "project-id",
					},
				},
			},
			wantKeys: []string{"project-id"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			indexer := NewAtlasDatabaseUserByExternalProjectsRefIndexer(zaptest.NewLogger(t))
			keys := indexer.Keys(tc.object)
			sort.Strings(keys)
			assert.Equal(t, tc.wantKeys, keys)
		})
	}
}
