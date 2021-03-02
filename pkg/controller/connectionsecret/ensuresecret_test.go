package connectionsecret

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

func TestAddCredentialsToConnectionURL(t *testing.T) {
	t.Run("Adding Credentials to standard url", func(t *testing.T) {
		url, err := addCredentialsToConnectionURL("mongodb://mongodb0.example.com:27017,mongodb1.example.com:27017/?authSource=admin", "super-user", "P@ssword!")
		assert.NoError(t, err)
		assert.Equal(t, "mongodb://super-user:P%40ssword%21@mongodb0.example.com:27017,mongodb1.example.com:27017/?authSource=admin", url)
	})
	t.Run("Adding Credentials to srv url", func(t *testing.T) {
		url, err := addCredentialsToConnectionURL("mongodb+srv://server.example.com/?authSource=$external&authMechanism=PLAIN&connectTimeoutMS=300000", "ldap_user", "Simple#")
		assert.NoError(t, err)
		assert.Equal(t, "mongodb+srv://ldap_user:Simple%23@server.example.com/?authSource=$external&authMechanism=PLAIN&connectTimeoutMS=300000", url)
	})
}

func TestEnsure(t *testing.T) {
	// Fake client
	scheme := runtime.NewScheme()
	utilruntime.Must(corev1.AddToScheme(scheme))
	utilruntime.Must(mdbv1.AddToScheme(scheme))
	fakeClient := fake.NewClientBuilder().WithScheme(scheme).Build()

	t.Run("Create/Update", func(t *testing.T) {
		data := dataForSecret()
		// Create
		err := Ensure(fakeClient, "testNs", "project1", "603e7bf38a94956835659ae5", "cluster1", data)
		assert.NoError(t, err)
		validateSecret(t, fakeClient, "testNs", "project1", "603e7bf38a94956835659ae5", "cluster1", data)

		// Update
		data.password = "new$!"
		data.srvConnURL = "mongodb+srv://mongodb10.example.com:27017/?authSource=admin&tls=true"
		data.connURL = "mongodb://mongodb10.example.com:27017,mongodb1.example.com:27017/?authSource=admin&tls=true"
		err = Ensure(fakeClient, "testNs", "project1", "603e7bf38a94956835659ae5", "cluster1", data)
		assert.NoError(t, err)
		validateSecret(t, fakeClient, "testNs", "project1", "603e7bf38a94956835659ae5", "cluster1", data)
	})

	t.Run("Create two different secrets", func(t *testing.T) {
		data := dataForSecret()
		// First secret
		err := Ensure(fakeClient, "testNs", "project1", "603e7bf38a94956835659ae5", "cluster1", data)
		assert.NoError(t, err)
		validateSecret(t, fakeClient, "testNs", "project1", "603e7bf38a94956835659ae5", "cluster1", data)

		// The second secret (the same cluster and user name but different projects)
		err = Ensure(fakeClient, "testNs", "project2", "903e7bf38a94256835659ae5", "cluster1", data)
		assert.NoError(t, err)
		validateSecret(t, fakeClient, "testNs", "project2", "903e7bf38a94256835659ae5", "cluster1", data)
	})

	t.Run("Create secret with special symbols", func(t *testing.T) {
		data := dataForSecret()
		data.dbUserName = "#simple@user_for.test"

		// Unfortunately, fake client doesn't validate object names, so this doesn't cover the validness of the produced name :(
		err := Ensure(fakeClient, "otherNs", "my@project", "603e7bf38a94956835659ae5", "some cluster!", data)
		assert.NoError(t, err)
		s := validateSecret(t, fakeClient, "otherNs", "my-project", "603e7bf38a94956835659ae5", "some-cluster", data)
		assert.Equal(t, "my-project-some-cluster-simple-user-for.test", s.Name)
	})
}

func validateSecret(t *testing.T, fakeClient client.Client, namespace, projectName, projectID, clusterName string, data ConnectionData) corev1.Secret {
	secret := corev1.Secret{}
	secretName := fmt.Sprintf("%s-%s-%s", projectName, clusterName, kube.NormalizeIdentifier(data.dbUserName))
	err := fakeClient.Get(context.Background(), kube.ObjectKey(namespace, secretName), &secret)
	assert.NoError(t, err)

	expectedData := map[string]string{
		"connectionString.standard":    buildConnectionURL(data.connURL, data.dbUserName, data.password),
		"connectionString.standardSrv": buildConnectionURL(data.srvConnURL, data.dbUserName, data.password),
		"username":                     data.dbUserName,
		"password":                     data.password,
	}
	expectedLabels := map[string]string{
		"atlas.mongodb.com/project-id":   projectID,
		"atlas.mongodb.com/cluster-name": clusterName,
	}
	// Dev note: it seems that when using fake client there is no serialization/deserialization happening, so we can
	// read from `secret.StringData` directly (not applicable for the real k8s cluster)
	assert.Equal(t, expectedData, secret.StringData)
	assert.Equal(t, expectedLabels, secret.Labels)

	return secret
}

func buildConnectionURL(connURL, userName, password string) string {
	url, err := addCredentialsToConnectionURL(connURL, userName, password)
	if err != nil {
		panic(err.Error())
	}
	return url
}

func dataForSecret() ConnectionData {
	return ConnectionData{
		dbUserName: "admin",
		connURL:    "mongodb://mongodb0.example.com:27017,mongodb1.example.com:27017/?authSource=admin",
		srvConnURL: "mongodb+srv://mongodb.example.com:27017/?authSource=admin",
		password:   "m@gick%",
	}
}
