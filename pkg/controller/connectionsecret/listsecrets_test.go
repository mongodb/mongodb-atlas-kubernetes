package connectionsecret

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

func TestListConnectionSecrets(t *testing.T) {
	t.Run("General Check", func(t *testing.T) {
		// Fake client
		scheme := runtime.NewScheme()
		utilruntime.Must(corev1.AddToScheme(scheme))
		utilruntime.Must(mdbv1.AddToScheme(scheme))
		fakeClient := fake.NewClientBuilder().WithScheme(scheme).Build()

		// c1, user1
		data := dataForSecret()
		data.dbUserName = "user1"
		assert.NoError(t, Ensure(fakeClient, "testNs", "p1", "603e7bf38a94956835659ae5", "c1", data))

		// c1, user2
		data = dataForSecret()
		data.dbUserName = "user2"
		assert.NoError(t, Ensure(fakeClient, "testNs", "p1", "603e7bf38a94956835659ae5", "c1", data))

		// c2, user1
		data = dataForSecret()
		data.dbUserName = "user1"
		assert.NoError(t, Ensure(fakeClient, "testNs", "p1", "603e7bf38a94956835659ae5", "c2", data))

		// c1, user1 but different project (p2)
		data = dataForSecret()
		data.dbUserName = "user1"
		assert.NoError(t, Ensure(fakeClient, "testNs", "p2", "some-other-project-id", "c1", data))

		// c1, user1 but different namespace
		data = dataForSecret()
		data.dbUserName = "user1"
		assert.NoError(t, Ensure(fakeClient, "otherNs", "p1", "603e7bf38a94956835659ae5", "c1", data))

		secrets, err := ListByClusterName(fakeClient, "testNs", "603e7bf38a94956835659ae5", "c1")
		assert.NoError(t, err)
		assert.Equal(t, []string{"p1-c1-user1", "p1-c1-user2"}, getSecretsNames(secrets))

		secrets, err = ListByClusterName(fakeClient, "testNs", "603e7bf38a94956835659ae5", "c2")
		assert.NoError(t, err)
		assert.Equal(t, []string{"p1-c2-user1"}, getSecretsNames(secrets))

		secrets, err = ListByClusterName(fakeClient, "testNs", "603e7bf38a94956835659ae5", "c3")
		assert.NoError(t, err)
		assert.Len(t, getSecretsNames(secrets), 0)

		secrets, err = ListByClusterName(fakeClient, "testNs", "non-existent-project-id", "c1")
		assert.NoError(t, err)
		assert.Len(t, getSecretsNames(secrets), 0)

		secrets, err = ListByUserName(fakeClient, "testNs", "603e7bf38a94956835659ae5", "user1")
		assert.NoError(t, err)
		assert.Equal(t, []string{"p1-c1-user1", "p1-c2-user1"}, getSecretsNames(secrets))

		secrets, err = ListByUserName(fakeClient, "testNs", "603e7bf38a94956835659ae5", "user2")
		assert.NoError(t, err)
		assert.Equal(t, []string{"p1-c1-user2"}, getSecretsNames(secrets))
	})
}

func getSecretsNames(secrets []corev1.Secret) []string {
	res := make([]string, 0)
	for _, secret := range secrets {
		res = append(res, secret.Name)
	}
	return res
}
