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

package indexer

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

const (
	testUsername = "matching-user"
)

func TestLocalCredentialsIndexer(t *testing.T) {
	for _, tc := range []struct {
		name       string
		object     client.Object
		index      string
		wantKeys   []string
		wantObject client.Object
	}{
		{
			name:       "should return nil on wrong type",
			object:     &akov2.AtlasBackupPolicy{},
			index:      AtlasDatabaseUserCredentialsIndex,
			wantKeys:   nil,
			wantObject: &akov2.AtlasDatabaseUser{},
		},
		{
			name:       "should return no keys when there are no references",
			object:     &akov2.AtlasDatabaseUser{},
			index:      AtlasDatabaseUserCredentialsIndex,
			wantKeys:   []string{},
			wantObject: &akov2.AtlasDatabaseUser{},
		},
		{
			name: "should return no keys when there is an empty reference",
			object: &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ConnectionSecret: &api.LocalObjectReference{},
					},
				},
			},
			index:      AtlasDatabaseUserCredentialsIndex,
			wantKeys:   []string{},
			wantObject: &akov2.AtlasDatabaseUser{},
		},
		{
			name: "should return keys when there is a reference",
			object: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user",
					Namespace: "ns",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ConnectionSecret: &api.LocalObjectReference{Name: "secret-ref"},
					},
				},
			},
			index:      AtlasDatabaseUserCredentialsIndex,
			wantKeys:   []string{"ns/secret-ref"},
			wantObject: &akov2.AtlasDatabaseUser{},
		},
		{
			name: "should return keys when there is a reference on a deployment",
			object: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "deployment",
					Namespace: "ns",
				},
				Spec: akov2.AtlasDeploymentSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ConnectionSecret: &api.LocalObjectReference{Name: "secret-ref"},
					},
				},
			},
			index:      AtlasDeploymentCredentialsIndex,
			wantKeys:   []string{"ns/secret-ref"},
			wantObject: &akov2.AtlasDeployment{},
		},
		{
			name: "should return keys when there is a reference on a custom role",
			object: &akov2.AtlasCustomRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "custom-role",
					Namespace: "ns",
				},
				Spec: akov2.AtlasCustomRoleSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ConnectionSecret: &api.LocalObjectReference{Name: "secret-ref"},
					},
				},
			},
			index:      AtlasCustomRoleCredentialsIndex,
			wantKeys:   []string{"ns/secret-ref"},
			wantObject: &akov2.AtlasCustomRole{},
		},
		{
			name: "should return keys when there is a reference on a private endpoint",
			object: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "private-endpoint",
					Namespace: "ns",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ConnectionSecret: &api.LocalObjectReference{Name: "secret-ref"},
					},
				},
			},
			index:      AtlasPrivateEndpointCredentialsIndex,
			wantKeys:   []string{"ns/secret-ref"},
			wantObject: &akov2.AtlasPrivateEndpoint{},
		},
		{
			name: "should return keys when there is a reference on a network container",
			object: &akov2.AtlasNetworkContainer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user",
					Namespace: "ns",
				},
				Spec: akov2.AtlasNetworkContainerSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ConnectionSecret: &api.LocalObjectReference{Name: "secret-ref"},
					},
				},
			},
			index:      AtlasNetworkContainerCredentialsIndex,
			wantKeys:   []string{"ns/secret-ref"},
			wantObject: &akov2.AtlasNetworkContainer{},
		},
		{
			name: "should return keys when there is a reference on a network peering",
			object: &akov2.AtlasNetworkPeering{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user",
					Namespace: "ns",
				},
				Spec: akov2.AtlasNetworkPeeringSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ConnectionSecret: &api.LocalObjectReference{Name: "secret-ref"},
					},
				},
			},
			index:      AtlasNetworkPeeringCredentialsIndex,
			wantKeys:   []string{"ns/secret-ref"},
			wantObject: &akov2.AtlasNetworkPeering{},
		},
	} {
		indexers := testIndexers(t)
		t.Run(tc.name, func(t *testing.T) {
			indexer := indexers[tc.index]
			keys := indexer.Keys(tc.object)
			sort.Strings(keys)
			assert.Equal(t, tc.wantKeys, keys)
			assert.Equal(t, tc.index, indexer.Name())
			assert.Equal(t, tc.wantObject, indexer.Object())
		})
	}
}

func TestCredentialsIndexMapperFunc(t *testing.T) {
	for _, tc := range []struct {
		name     string
		mapperFn func(kubeClient client.Client, logger *zap.SugaredLogger) handler.MapFunc
		objects  []client.Object
		index    string
		input    client.Object
		output   client.Object
		want     []reconcile.Request
	}{
		{
			name:     "nil input & list renders nil",
			index:    AtlasDatabaseUserCredentialsIndex,
			output:   &akov2.AtlasDatabaseUser{},
			mapperFn: dbUserMapperFunc,
		},
		{
			name:     "nil list renders empty list",
			index:    AtlasDatabaseUserCredentialsIndex,
			output:   &akov2.AtlasDatabaseUser{},
			mapperFn: dbUserMapperFunc,
			input:    &corev1.Secret{},
			want:     []reconcile.Request{},
		},
		{
			name:     "empty input with proper empty list type renders empty list",
			index:    AtlasDatabaseUserCredentialsIndex,
			output:   &akov2.AtlasDatabaseUser{},
			mapperFn: dbUserMapperFunc,
			input:    &corev1.Secret{},
			want:     []reconcile.Request{},
		},
		{
			name:     "matching input credentials renders matching user",
			index:    AtlasDatabaseUserCredentialsIndex,
			output:   &akov2.AtlasDatabaseUser{},
			mapperFn: dbUserMapperFunc,
			input:    newTestSecret("matching-user-secret-ref"),
			objects: []client.Object{
				&akov2.AtlasDatabaseUser{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "matching-user",
						Namespace: "ns",
					},
					Spec: akov2.AtlasDatabaseUserSpec{
						ProjectDualReference: akov2.ProjectDualReference{
							ConnectionSecret: &api.LocalObjectReference{
								Name: "matching-user-secret-ref",
							},
						},
					},
				},
			},
			want: []reconcile.Request{
				{NamespacedName: types.NamespacedName{
					Name:      "matching-user",
					Namespace: "ns",
				}},
			},
		},
		{
			name:   "matching input credentials renders matching deployment",
			index:  AtlasDeploymentCredentialsIndex,
			output: &akov2.AtlasDeployment{},
			mapperFn: func(kubeClient client.Client, logger *zap.SugaredLogger) handler.MapFunc {
				return CredentialsIndexMapperFunc(
					AtlasDeploymentCredentialsIndex,
					func() *akov2.AtlasDeploymentList { return &akov2.AtlasDeploymentList{} },
					DeploymentRequests,
					kubeClient,
					logger,
				)
			},
			input: newTestSecret("matching-deployment-secret-ref"),
			objects: []client.Object{
				&akov2.AtlasDeployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "matching-deployment",
						Namespace: "ns",
					},
					Spec: akov2.AtlasDeploymentSpec{
						ProjectDualReference: akov2.ProjectDualReference{
							ConnectionSecret: &api.LocalObjectReference{
								Name: "matching-deployment-secret-ref",
							},
						},
					},
				},
			},
			want: []reconcile.Request{
				{NamespacedName: types.NamespacedName{
					Name:      "matching-deployment",
					Namespace: "ns",
				}},
			},
		},
		{
			name:   "matching input credentials renders matching custom role",
			index:  AtlasCustomRoleCredentialsIndex,
			output: &akov2.AtlasCustomRole{},
			mapperFn: func(kubeClient client.Client, logger *zap.SugaredLogger) handler.MapFunc {
				return CredentialsIndexMapperFunc(
					AtlasCustomRoleCredentialsIndex,
					func() *akov2.AtlasCustomRoleList { return &akov2.AtlasCustomRoleList{} },
					CustomRoleRequests,
					kubeClient,
					logger,
				)
			},
			input: newTestSecret("matching-custom-role-secret-ref"),
			objects: []client.Object{
				&akov2.AtlasCustomRole{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "matching-custom-role",
						Namespace: "ns",
					},
					Spec: akov2.AtlasCustomRoleSpec{
						ProjectDualReference: akov2.ProjectDualReference{
							ConnectionSecret: &api.LocalObjectReference{
								Name: "matching-custom-role-secret-ref",
							},
						},
					},
				},
			},
			want: []reconcile.Request{
				{NamespacedName: types.NamespacedName{
					Name:      "matching-custom-role",
					Namespace: "ns",
				}},
			},
		},
		{
			name:   "matching input credentials renders matching private endpoint",
			index:  AtlasPrivateEndpointCredentialsIndex,
			output: &akov2.AtlasPrivateEndpoint{},
			mapperFn: func(kubeClient client.Client, logger *zap.SugaredLogger) handler.MapFunc {
				return CredentialsIndexMapperFunc(
					AtlasPrivateEndpointCredentialsIndex,
					func() *akov2.AtlasPrivateEndpointList { return &akov2.AtlasPrivateEndpointList{} },
					PrivateEndpointRequests,
					kubeClient,
					logger,
				)
			},
			input: newTestSecret("matching-private-endpoint-secret-ref"),
			objects: []client.Object{
				&akov2.AtlasPrivateEndpoint{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "matching-private-endpoint",
						Namespace: "ns",
					},
					Spec: akov2.AtlasPrivateEndpointSpec{
						ProjectDualReference: akov2.ProjectDualReference{
							ConnectionSecret: &api.LocalObjectReference{
								Name: "matching-private-endpoint-secret-ref",
							},
						},
					},
				},
			},
			want: []reconcile.Request{
				{NamespacedName: types.NamespacedName{
					Name:      "matching-private-endpoint",
					Namespace: "ns",
				}},
			},
		},
		{
			name:   "matching input credentials renders matching network container",
			index:  AtlasNetworkContainerCredentialsIndex,
			output: &akov2.AtlasNetworkContainer{},
			mapperFn: func(kubeClient client.Client, logger *zap.SugaredLogger) handler.MapFunc {
				return CredentialsIndexMapperFunc(
					AtlasNetworkContainerCredentialsIndex,
					func() *akov2.AtlasNetworkContainerList { return &akov2.AtlasNetworkContainerList{} },
					NetworkContainerRequests,
					kubeClient,
					logger,
				)
			},
			input: newTestSecret("matching-container-secret-ref"),
			objects: []client.Object{
				&akov2.AtlasNetworkContainer{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "matching-container",
						Namespace: "ns",
					},
					Spec: akov2.AtlasNetworkContainerSpec{
						ProjectDualReference: akov2.ProjectDualReference{
							ConnectionSecret: &api.LocalObjectReference{
								Name: "matching-container-secret-ref",
							},
						},
					},
				},
			},
			want: []reconcile.Request{
				{NamespacedName: types.NamespacedName{
					Name:      "matching-container",
					Namespace: "ns",
				}},
			},
		},
		{
			name:   "matching input credentials renders matching network peering",
			index:  AtlasNetworkPeeringCredentialsIndex,
			output: &akov2.AtlasNetworkPeering{},
			mapperFn: func(kubeClient client.Client, logger *zap.SugaredLogger) handler.MapFunc {
				return CredentialsIndexMapperFunc(
					AtlasNetworkPeeringCredentialsIndex,
					func() *akov2.AtlasNetworkPeeringList { return &akov2.AtlasNetworkPeeringList{} },
					NetworkPeeringRequests,
					kubeClient,
					logger,
				)
			},
			input: newTestSecret("matching-peering-secret-ref"),
			objects: []client.Object{
				&akov2.AtlasNetworkPeering{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "matching-peering",
						Namespace: "ns",
					},
					Spec: akov2.AtlasNetworkPeeringSpec{
						ProjectDualReference: akov2.ProjectDualReference{
							ConnectionSecret: &api.LocalObjectReference{
								Name: "matching-peering-secret-ref",
							},
						},
					},
				},
			},
			want: []reconcile.Request{
				{NamespacedName: types.NamespacedName{
					Name:      "matching-peering",
					Namespace: "ns",
				}},
			},
		},
	} {
		scheme := runtime.NewScheme()
		assert.NoError(t, corev1.AddToScheme(scheme))
		assert.NoError(t, akov2.AddToScheme(scheme))
		indexers := testIndexers(t)
		t.Run(tc.name, func(t *testing.T) {
			indexer := indexers[tc.index]
			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(tc.objects...).
				WithIndex(
					tc.output,
					tc.index,
					func(obj client.Object) []string {
						return indexer.Keys(obj)
					}).
				Build()
			fn := tc.mapperFn(fakeClient, zaptest.NewLogger(t).Sugar())
			result := fn(context.Background(), tc.input)
			assert.Equal(t, tc.want, result)
		})
	}
}

func TestCredentialsIndexMapperFuncRace(t *testing.T) {
	scheme := runtime.NewScheme()
	assert.NoError(t, corev1.AddToScheme(scheme))
	assert.NoError(t, akov2.AddToScheme(scheme))
	indexer := NewLocalCredentialsIndexer(
		AtlasDatabaseUserCredentialsIndex,
		&akov2.AtlasDatabaseUser{},
		zaptest.NewLogger(t),
	)
	objs := make([]client.Object, 10)
	for i := range objs {
		objs[i] = &akov2.AtlasDatabaseUser{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%d", testUsername, i),
				Namespace: "ns",
			},
			Spec: akov2.AtlasDatabaseUserSpec{
				ProjectDualReference: akov2.ProjectDualReference{
					ConnectionSecret: &api.LocalObjectReference{
						Name: fmt.Sprintf("%s-%d-secret-ref", testUsername, i),
					},
				},
			},
		}
	}
	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(objs...).
		WithIndex(
			&akov2.AtlasDatabaseUser{},
			AtlasDatabaseUserCredentialsIndex,
			func(obj client.Object) []string {
				return indexer.Keys(obj)
			}).
		Build()
	fn := dbUserMapperFunc(fakeClient, zaptest.NewLogger(t).Sugar())
	ctx := context.Background()
	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			input := newTestSecret(fmt.Sprintf("%s-%d-secret-ref", testUsername, i))
			result := fn(ctx, input)
			if i < len(objs) {
				assert.NotEmpty(t, result, "failed to find for index %d", i)
			} else {
				assert.Empty(t, result, "failed not to find for index %d", i)
			}
		}(i)
	}
	wg.Wait()
}

func newTestSecret(name string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "ns",
		},
	}
}

func dbUserMapperFunc(kubeClient client.Client, logger *zap.SugaredLogger) handler.MapFunc {
	return CredentialsIndexMapperFunc(
		AtlasDatabaseUserCredentialsIndex,
		func() *akov2.AtlasDatabaseUserList { return &akov2.AtlasDatabaseUserList{} },
		DatabaseUserRequests,
		kubeClient,
		logger,
	)
}

func testIndexers(t *testing.T) map[string]*LocalCredentialIndexer {
	t.Helper()

	logger := zaptest.NewLogger(t)
	indexers := map[string]*LocalCredentialIndexer{}
	indexers[AtlasDatabaseUserCredentialsIndex] = NewAtlasDatabaseUserByCredentialIndexer(logger)
	indexers[AtlasDeploymentCredentialsIndex] = NewAtlasDeploymentByCredentialIndexer(logger)
	indexers[AtlasCustomRoleCredentialsIndex] = NewAtlasCustomRoleByCredentialIndexer(logger)
	indexers[AtlasPrivateEndpointCredentialsIndex] = NewAtlasPrivateEndpointByCredentialIndexer(logger)
	indexers[AtlasNetworkContainerCredentialsIndex] = NewAtlasNetworkContainerByCredentialIndexer(logger)
	indexers[AtlasNetworkPeeringCredentialsIndex] = NewAtlasNetworkPeeringByCredentialIndexer(logger)
	indexers[AtlasThirdPartyIntegrationCredentialsIndex] = NewAtlasThirdPartyIntegrationByCredentialIndexer(logger)
	return indexers
}
