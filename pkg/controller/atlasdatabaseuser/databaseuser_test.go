package atlasdatabaseuser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

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
	scopeClusters := filterScopeClusters(mdbv1.AtlasDatabaseUser{Spec: mdbv1.AtlasDatabaseUserSpec{Scopes: scopeSpecs}}, clusters)
	assert.Equal(t, []string{"cluster1"}, scopeClusters)
}

func TestCheckUserExpired(t *testing.T) {
	// Fake client
	scheme := runtime.NewScheme()
	utilruntime.Must(corev1.AddToScheme(scheme))
	utilruntime.Must(mdbv1.AddToScheme(scheme))
	fakeClient := fake.NewClientBuilder().WithScheme(scheme).Build()

	t.Run("Validate DeleteAfterDate", func(t *testing.T) {
		result := checkUserExpired(zap.S(), fakeClient, "", *mdbv1.DefaultDBUser("ns", "theuser", "").WithDeleteAfterDate("foo"))
		assert.False(t, result.IsOk())
		assert.Equal(t, reconcile.Result{}, result.ReconcileResult())

		result = checkUserExpired(zap.S(), fakeClient, "", *mdbv1.DefaultDBUser("ns", "theuser", "").WithDeleteAfterDate("2021/11/30T15:04:05"))
		assert.False(t, result.IsOk())
		assert.Equal(t, reconcile.Result{}, result.ReconcileResult())
	})
}
