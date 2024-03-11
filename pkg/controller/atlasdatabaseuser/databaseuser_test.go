package atlasdatabaseuser

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func init() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}

func TestFilterScopeClusters(t *testing.T) {
	scopeSpecs := []akov2.ScopeSpec{{
		Name: "dbLake",
		Type: akov2.DataLakeScopeType,
	}, {
		Name: "cluster1",
		Type: akov2.DeploymentScopeType,
	}, {
		Name: "cluster2",
		Type: akov2.DeploymentScopeType,
	}}
	clusters := []string{"cluster1", "cluster4", "cluster5"}
	scopeClusters := filterScopeDeployments(akov2.AtlasDatabaseUser{Spec: akov2.AtlasDatabaseUserSpec{Scopes: scopeSpecs}}, clusters)
	assert.Equal(t, []string{"cluster1"}, scopeClusters)
}

func TestCheckUserExpired(t *testing.T) {
	// Fake client
	scheme := runtime.NewScheme()
	utilruntime.Must(corev1.AddToScheme(scheme))
	utilruntime.Must(akov2.AddToScheme(scheme))
	fakeClient := fake.NewClientBuilder().WithScheme(scheme).Build()

	t.Run("Validate DeleteAfterDate", func(t *testing.T) {
		result := checkUserExpired(context.Background(), zap.S(), fakeClient, "", *akov2.DefaultDBUser("ns", "theuser", "").WithDeleteAfterDate("foo"))
		assert.False(t, result.IsOk())
		assert.Equal(t, reconcile.Result{}, result.ReconcileResult())

		result = checkUserExpired(context.Background(), zap.S(), fakeClient, "", *akov2.DefaultDBUser("ns", "theuser", "").WithDeleteAfterDate("2021/11/30T15:04:05"))
		assert.False(t, result.IsOk())
	})
	t.Run("User Marked Expired", func(t *testing.T) {
		data := dataForSecret()
		// Create a connection secret
		_, err := connectionsecret.Ensure(context.Background(), fakeClient, "tetNs", "project1", "603e7bf38a94956835659ae5", "cluster1", data)
		assert.NoError(t, err)
		// The secret for the other project
		_, err = connectionsecret.Ensure(context.Background(), fakeClient, "testNs", "project2", "dsfsdf234234sdfdsf23423", "cluster1", data)
		assert.NoError(t, err)

		before := time.Now().UTC().Add(time.Minute * -1).Format("2006-01-02T15:04:05.999Z")
		user := *akov2.DefaultDBUser("testNs", data.DBUserName, "").WithDeleteAfterDate(before)
		result := checkUserExpired(context.Background(), zap.S(), fakeClient, "603e7bf38a94956835659ae5", user)
		assert.False(t, result.IsOk())
		assert.Equal(t, reconcile.Result{}, result.ReconcileResult())

		// The secret has been removed
		secret := corev1.Secret{}
		secretName := fmt.Sprintf("%s-%s-%s", "project1", "cluster1", kube.NormalizeIdentifier(data.DBUserName))
		err = fakeClient.Get(context.Background(), kube.ObjectKey("testNs", secretName), &secret)
		assert.Error(t, err)

		// The other secret still exists
		secretName = fmt.Sprintf("%s-%s-%s", "project2", "cluster1", kube.NormalizeIdentifier(data.DBUserName))
		err = fakeClient.Get(context.Background(), kube.ObjectKey("testNs", secretName), &secret)
		assert.NoError(t, err)
	})
	t.Run("No expiration happened", func(t *testing.T) {
		data := dataForSecret()
		// Create a connection secret
		_, err := connectionsecret.Ensure(context.Background(), fakeClient, "testNs", "project1", "603e7bf38a94956835659ae5", "cluster1", data)
		assert.NoError(t, err)
		after := time.Now().UTC().Add(time.Minute * 1).Format("2006-01-02T15:04:05")
		result := checkUserExpired(context.Background(), zap.S(), fakeClient, "603e7bf38a94956835659ae5", *akov2.DefaultDBUser("testNs", data.DBUserName, "").WithDeleteAfterDate(after))
		assert.True(t, result.IsOk())

		// The secret is still there
		secret := corev1.Secret{}
		secretName := fmt.Sprintf("%s-%s-%s", "project1", "cluster1", kube.NormalizeIdentifier(data.DBUserName))
		err = fakeClient.Get(context.Background(), kube.ObjectKey("testNs", secretName), &secret)
		assert.NoError(t, err)
	})
}

func TestHandleUserNameChange(t *testing.T) {
	t.Run("Only one user after name change", func(t *testing.T) {
		user := *akov2.DefaultDBUser("ns", "theuser", "project1")
		user.Spec.Username = "differentuser"
		user.Status.UserName = "theuser"
		ctx := workflow.NewContext(zap.S(), []status.Condition{}, nil)
		ctx.Client = mongodbatlas.NewClient(&http.Client{})
		result := handleUserNameChange(ctx, "", user)
		assert.True(t, result.IsOk())
	})
}

func dataForSecret() connectionsecret.ConnectionData {
	return connectionsecret.ConnectionData{
		DBUserName: "admin",
		ConnURL:    "mongodb://mongodb0.example.com:27017,mongodb1.example.com:27017/?authSource=admin",
		SrvConnURL: "mongodb+srv://mongodb.example.com:27017/?authSource=admin",
		Password:   "m@gick%",
	}
}
