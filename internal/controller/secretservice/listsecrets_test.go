// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package secretservice

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

func TestListConnectionSecrets(t *testing.T) {
	t.Run("General Check", func(t *testing.T) {
		// Fake client
		scheme := runtime.NewScheme()
		utilruntime.Must(corev1.AddToScheme(scheme))
		utilruntime.Must(akov2.AddToScheme(scheme))
		fakeClient := fake.NewClientBuilder().WithScheme(scheme).Build()

		// c1, user1
		data := dataForSecret()
		data.DBUserName = "user1"
		_, err := Ensure(context.Background(), fakeClient, "testNs", "p1", "603e7bf38a94956835659ae5", "c1", data)
		assert.NoError(t, err)

		// c1, user2
		data = dataForSecret()
		data.DBUserName = "user2"
		_, err = Ensure(context.Background(), fakeClient, "testNs", "p1", "603e7bf38a94956835659ae5", "c1", data)
		assert.NoError(t, err)

		// c2, user1
		data = dataForSecret()
		data.DBUserName = "user1"
		_, err = Ensure(context.Background(), fakeClient, "testNs", "p1", "603e7bf38a94956835659ae5", "c2", data)
		assert.NoError(t, err)

		// c1, user1 but different project (p2)
		data = dataForSecret()
		data.DBUserName = "user1"
		_, err = Ensure(context.Background(), fakeClient, "testNs", "p2", "some-other-project-id", "c1", data)
		assert.NoError(t, err)

		// c1, user1 but different namespace
		data = dataForSecret()
		data.DBUserName = "user1"
		_, err = Ensure(context.Background(), fakeClient, "otherNs", "p1", "603e7bf38a94956835659ae5", "c1", data)
		assert.NoError(t, err)

		secrets, err := ListByDeploymentName(context.Background(), fakeClient, "testNs", "603e7bf38a94956835659ae5", "c1")
		assert.NoError(t, err)
		assert.Equal(t, []string{"p1-c1-user1", "p1-c1-user2"}, getSecretsNames(secrets))

		secrets, err = ListByDeploymentName(context.Background(), fakeClient, "testNs", "603e7bf38a94956835659ae5", "c2")
		assert.NoError(t, err)
		assert.Equal(t, []string{"p1-c2-user1"}, getSecretsNames(secrets))

		secrets, err = ListByDeploymentName(context.Background(), fakeClient, "testNs", "603e7bf38a94956835659ae5", "c3")
		assert.NoError(t, err)
		assert.Len(t, getSecretsNames(secrets), 0)

		secrets, err = ListByDeploymentName(context.Background(), fakeClient, "testNs", "non-existent-project-id", "c1")
		assert.NoError(t, err)
		assert.Len(t, getSecretsNames(secrets), 0)

		secrets, err = ListByUserName(context.Background(), fakeClient, "testNs", "603e7bf38a94956835659ae5", "user1")
		assert.NoError(t, err)
		assert.Equal(t, []string{"p1-c1-user1", "p1-c2-user1"}, getSecretsNames(secrets))

		secrets, err = ListByUserName(context.Background(), fakeClient, "testNs", "603e7bf38a94956835659ae5", "user2")
		assert.NoError(t, err)
		assert.Equal(t, []string{"p1-c1-user2"}, getSecretsNames(secrets))
	})
	t.Run("Special symbols in names", func(t *testing.T) {
		// Fake client
		scheme := runtime.NewScheme()
		utilruntime.Must(corev1.AddToScheme(scheme))
		utilruntime.Must(akov2.AddToScheme(scheme))
		fakeClient := fake.NewClientBuilder().WithScheme(scheme).Build()

		data := dataForSecret()
		data.DBUserName = "user1"
		_, err := Ensure(context.Background(), fakeClient, "testNs", "#nice project!", "603e7bf38a94956835659ae5", "the cluster@thecompany.com/", data)
		assert.NoError(t, err)

		secrets, err := ListByDeploymentName(context.Background(), fakeClient, "testNs", "603e7bf38a94956835659ae5", "the cluster@thecompany.com/")
		assert.NoError(t, err)
		assert.Equal(t, []string{"nice-project-the-cluster-thecompany.com-user1"}, getSecretsNames(secrets))
	})
}

func getSecretsNames(secrets []corev1.Secret) []string {
	res := make([]string, 0, len(secrets))
	for _, secret := range secrets {
		res = append(res, secret.Name)
	}
	return res
}
