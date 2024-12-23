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

func TestAtlasDatabaseUserBySecretsIndexer(t *testing.T) {
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
			name: "should return nil when there is an empty reference",
			object: &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					PasswordSecret: &common.ResourceRef{},
				},
			},
		},
		{
			name: "should return database user namespace",
			object: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user",
					Namespace: "ns",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					PasswordSecret: &common.ResourceRef{
						Name: "someSecret",
					},
				},
			},
			wantKeys: []string{"ns/someSecret"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			indexer := NewAtlasDatabaseUserBySecretsIndexer(zaptest.NewLogger(t))
			keys := indexer.Keys(tc.object)
			sort.Strings(keys)
			assert.Equal(t, tc.wantKeys, keys)
		})
	}
}
