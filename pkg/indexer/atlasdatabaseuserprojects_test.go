package indexer

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
)

func TestAtlasDatabaseUserByProjectsIndexer(t *testing.T) {
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
			name: "should return nil when there is an empty reference for project",
			object: &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					Project: &common.ResourceRefNamespaced{},
				},
			},
		},
		{
			name: "should return project reference with database user namespace",
			object: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user",
					Namespace: "ns",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Project: &common.ResourceRefNamespaced{
						Name: "someProject",
					},
				},
			},
			wantKeys: []string{"ns/someProject"},
		},
		{
			name: "should return project reference",
			object: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user",
					Namespace: "ns",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Project: &common.ResourceRefNamespaced{
						Name:      "someProject",
						Namespace: "nsProject",
					},
				},
			},
			wantKeys: []string{"nsProject/someProject"},
		},
		{
			name: "should return nil when there is an empty reference for external project",
			object: &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					AtlasProjectRef: &akov2.ExternalProjectReference{},
				},
			},
		},
		{
			name: "should return external project reference",
			object: &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					AtlasProjectRef: &akov2.ExternalProjectReference{
						ID: "project-id",
					},
				},
			},
			wantKeys: []string{"project-id"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			indexer := NewAtlasDatabaseUserByProjectsIndexer(zaptest.NewLogger(t))
			keys := indexer.Keys(tc.object)
			sort.Strings(keys)
			assert.Equal(t, tc.wantKeys, keys)
		})
	}
}
