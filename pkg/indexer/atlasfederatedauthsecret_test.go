package indexer

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
)

func TestAtlasFederatedAuthBySecretsIndexer(t *testing.T) {
	for _, tc := range []struct {
		name     string
		object   client.Object
		wantKeys []string
	}{
		{
			name:   "should return nil on wrong type",
			object: &akov2.AtlasDatabaseUser{},
		},
		{
			name:   "should return nil when there are no references",
			object: &akov2.AtlasFederatedAuth{},
		},
		{
			name: "should return nil when there is an empty reference",
			object: &akov2.AtlasFederatedAuth{
				Spec: akov2.AtlasFederatedAuthSpec{
					ConnectionSecretRef: common.ResourceRefNamespaced{},
				},
			},
		},
		{
			name: "should return project namespace if name is set only",
			object: &akov2.AtlasFederatedAuth{
				ObjectMeta: metav1.ObjectMeta{Name: "name", Namespace: "ns"},
				Spec: akov2.AtlasFederatedAuthSpec{
					ConnectionSecretRef: common.ResourceRefNamespaced{Name: "someSecret"},
				},
			},
			wantKeys: []string{"ns/someSecret"},
		},
		{
			name: "should return secret namespace and name if set",
			object: &akov2.AtlasFederatedAuth{
				ObjectMeta: metav1.ObjectMeta{Name: "name", Namespace: "ns"},
				Spec: akov2.AtlasFederatedAuthSpec{
					ConnectionSecretRef: common.ResourceRefNamespaced{Name: "someSecret", Namespace: "secretNamespace"},
				},
			},
			wantKeys: []string{"secretNamespace/someSecret"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			indexer := NewAtlasFederatedAuthBySecretsIndexer(zaptest.NewLogger(t))
			keys := indexer.Keys(tc.object)
			sort.Strings(keys)
			assert.Equal(t, tc.wantKeys, keys)
		})
	}
}
