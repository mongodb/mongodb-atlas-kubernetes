package atlasdatabaseuser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

func TestFilterScopeClusters(t *testing.T) {
	scopeSpecs := []mdbv1.ScopeSpec{{
		Name: "dbLake",
		Type: mdbv1.DataLakeScopeType,
	}, {
		Name: "cluster1",
		Type: mdbv1.ClusterScopeType,
	}, {
		Name: "cluster2",
		Type: mdbv1.ClusterScopeType,
	}}
	clusters := []mongodbatlas.Cluster{{Name: "cluster1"}, {Name: "cluster4"}, {Name: "cluster5"}}
	scopeClusters := filterScopeClusters(scopeSpecs, clusters)
	assert.Equal(t, []string{"cluster1"}, scopeClusters)
}
