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

package connectionsecret_test

import (
	"context"
	"testing"
	"time"

	k8s "github.com/crd2go/crd2go/k8s"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/connectionsecret/target"
	generatedv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

const (
	testNamespace = "default"
	testGroupID   = "62b6e34b3d91647abb20e7b8"
	testGroupName = "test-group"
	testUserName  = "test-user"
)

func TestConnectionSecretReconciler_Reconcile(t *testing.T) {
	for _, tc := range []struct {
		name           string
		user           *generatedv1.DatabaseUser
		objects        []client.Object
		targets        []target.ConnectionTarget
		wantResult     reconcile.Result
		wantErr        string
		wantSecretName string
	}{
		{
			name:       "user not found returns empty result",
			user:       nil,
			wantResult: reconcile.Result{},
		},
		{
			name: "user without group reference returns error",
			user: newDatabaseUser(testUserName, testNamespace, func(u *generatedv1.DatabaseUser) {
				u.Spec.V20250312.GroupRef = nil
				u.Spec.V20250312.GroupId = nil
			}),
			wantErr: "cannot get project ID",
		},
		{
			name: "user with missing group returns error",
			user: newDatabaseUser(testUserName, testNamespace, func(u *generatedv1.DatabaseUser) {
				u.Spec.V20250312.GroupRef = &k8s.LocalReference{Name: "non-existent"}
			}),
			wantErr: "failed to get Group",
		},
		{
			name: "user not ready requeues",
			user: newDatabaseUser(testUserName, testNamespace, nil),
			objects: []client.Object{
				newGroup(testGroupName, testNamespace, testGroupID),
			},
			targets:    []target.ConnectionTarget{&fakeConnectionTarget{}},
			wantResult: reconcile.Result{RequeueAfter: 10 * time.Second},
		},
		{
			name: "ready user with no targets returns empty result",
			user: newDatabaseUser(testUserName, testNamespace, withReadyCondition),
			objects: []client.Object{
				newGroup(testGroupName, testNamespace, testGroupID),
			},
			targets:    []target.ConnectionTarget{&fakeConnectionTarget{}},
			wantResult: reconcile.Result{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scheme := createTestScheme(t)
			clientBuilder := fake.NewClientBuilder().WithScheme(scheme)

			if tc.user != nil {
				clientBuilder = clientBuilder.
					WithObjects(tc.user).
					WithStatusSubresource(tc.user)
			}
			if len(tc.objects) > 0 {
				clientBuilder = clientBuilder.WithObjects(tc.objects...)
			}

			fakeClient := clientBuilder.Build()
			logger := zaptest.NewLogger(t)

			r := &connectionsecret.ConnectionSecretReconciler{
				Client:                fakeClient,
				Scheme:                scheme,
				Logger:                logger,
				ConnectionTargetKinds: tc.targets,
			}

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      testUserName,
					Namespace: testNamespace,
				},
			}

			result, err := r.Reconcile(context.Background(), req)

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantResult, result)

			if tc.wantSecretName != "" {
				secret := &corev1.Secret{}
				err := fakeClient.Get(context.Background(), client.ObjectKey{
					Namespace: testNamespace,
					Name:      tc.wantSecretName,
				}, secret)
				require.NoError(t, err)
			}
		})
	}
}

func TestK8sConnectionSecretName(t *testing.T) {
	for _, tc := range []struct {
		name         string
		projectName  string
		targetName   string
		userName     string
		targetType   string
		wantContains string
	}{
		{
			name:         "simple names",
			projectName:  "my-project",
			targetName:   "my-cluster",
			userName:     "admin",
			targetType:   "cluster",
			wantContains: "my-project-my-cluster-admin",
		},
		{
			name:         "names with special characters normalized",
			projectName:  "My_Project",
			targetName:   "My.Cluster",
			userName:     "admin@example.com",
			targetType:   "flexcluster",
			wantContains: "my-project-my.cluster-admin-example.com",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := connectionsecret.K8sConnectionSecretName(
				tc.projectName, tc.targetName, tc.userName, tc.targetType,
			)
			assert.Contains(t, result, tc.wantContains)
		})
	}
}

func TestCreateURL(t *testing.T) {
	for _, tc := range []struct {
		name     string
		hostname string
		username string
		password string
		want     string
		wantErr  bool
	}{
		{
			name:     "empty hostname returns empty string",
			hostname: "",
			username: "user",
			password: "pass",
			want:     "",
		},
		{
			name:     "valid hostname with credentials",
			hostname: "mongodb://cluster.mongodb.net",
			username: "myuser",
			password: "mypass",
			want:     "mongodb://myuser:mypass@cluster.mongodb.net",
		},
		{
			name:     "srv connection string",
			hostname: "mongodb+srv://cluster.mongodb.net",
			username: "admin",
			password: "secret123",
			want:     "mongodb+srv://admin:secret123@cluster.mongodb.net",
		},
		{
			name:     "password with special characters",
			hostname: "mongodb://cluster.mongodb.net",
			username: "user",
			password: "p@ss:word/test",
			want:     "mongodb://user:p%40ss%3Aword%2Ftest@cluster.mongodb.net",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result, err := connectionsecret.CreateURL(tc.hostname, tc.username, tc.password)

			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, result)
		})
	}
}

// Helper functions

func createTestScheme(t *testing.T) *runtime.Scheme {
	t.Helper()
	scheme := runtime.NewScheme()
	require.NoError(t, generatedv1.AddToScheme(scheme))
	require.NoError(t, corev1.AddToScheme(scheme))
	return scheme
}

//nolint:unparam
func newDatabaseUser(name, namespace string, modifiers ...func(*generatedv1.DatabaseUser)) *generatedv1.DatabaseUser {
	user := &generatedv1.DatabaseUser{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: generatedv1.DatabaseUserSpec{
			V20250312: &generatedv1.DatabaseUserSpecV20250312{
				GroupRef: &k8s.LocalReference{Name: testGroupName},
				Entry: &generatedv1.DatabaseUserSpecV20250312Entry{
					Username:     name,
					DatabaseName: "admin",
					PasswordSecretRef: &generatedv1.PasswordSecretRef{
						Name: "password-secret",
						Key:  pointer.MakePtr("password"),
					},
				},
			},
		},
	}

	for _, m := range modifiers {
		if m != nil {
			m(user)
		}
	}

	return user
}

func withReadyCondition(u *generatedv1.DatabaseUser) {
	conditions := []metav1.Condition{
		{
			Type:               state.ReadyCondition,
			Status:             metav1.ConditionTrue,
			LastTransitionTime: metav1.Now(),
			Reason:             "Ready",
		},
	}
	u.Status.Conditions = &conditions
}

func newGroup(name, namespace, id string) *generatedv1.Group {
	return &generatedv1.Group{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Status: generatedv1.GroupStatus{
			V20250312: &generatedv1.GroupStatusV20250312{
				Id: &id,
			},
		},
	}
}

// fakeConnectionTarget implements target.ConnectionTarget for testing
type fakeConnectionTarget struct {
	instances []target.ConnectionTargetInstance
	err       error
}

func (f *fakeConnectionTarget) ListForProject(ctx context.Context, projectID string) ([]target.ConnectionTargetInstance, error) {
	return f.instances, f.err
}

func (f *fakeConnectionTarget) GetConnectionTargetInstance(obj client.Object) target.ConnectionTargetInstance {
	return nil
}
