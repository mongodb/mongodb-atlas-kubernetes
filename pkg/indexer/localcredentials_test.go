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

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

const (
	testUsername = "matching-user"
)

func TestAtlasDatabaseUserLocalCredentialsIndexer(t *testing.T) {
	for _, tc := range []struct {
		name     string
		object   client.Object
		wantKeys []string
	}{
		{
			name:     "should return nil on wrong type",
			object:   &akov2.AtlasBackupPolicy{},
			wantKeys: nil,
		},
		{
			name:     "should return no keys when there are no references",
			object:   &akov2.AtlasDatabaseUser{},
			wantKeys: []string{},
		},
		{
			name: "should return no keys when there is an empty reference",
			object: &akov2.AtlasDatabaseUser{
				Spec: akov2.AtlasDatabaseUserSpec{
					LocalCredentialHolder: api.LocalCredentialHolder{
						ConnectionSecret: &api.LocalObjectReference{},
					},
				},
			},
			wantKeys: []string{},
		},
		{
			name: "should return keys when there is a reference",
			object: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user",
					Namespace: "ns",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					LocalCredentialHolder: api.LocalCredentialHolder{
						ConnectionSecret: &api.LocalObjectReference{Name: "secret-ref"},
					},
				},
			},
			wantKeys: []string{"ns/secret-ref"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			indexer := NewLocalCredentialsIndexer(
				AtlasDatabaseUserCredentialsIndex,
				&akov2.AtlasDatabaseUser{},
				zaptest.NewLogger(t),
			)
			keys := indexer.Keys(tc.object)
			sort.Strings(keys)
			assert.Equal(t, tc.wantKeys, keys)
			assert.Equal(t, AtlasDatabaseUserCredentialsIndex, indexer.Name())
			assert.Equal(t, &akov2.AtlasDatabaseUser{}, indexer.Object())
		})
	}
}

func TestCredentialsIndexMapperFunc(t *testing.T) {
	for _, tc := range []struct {
		name     string
		mapperFn func(kubeClient client.Client, logger *zap.SugaredLogger) handler.MapFunc
		objects  []client.Object
		input    client.Object
		want     []reconcile.Request
	}{
		{
			name:     "nil input & list renders nil",
			mapperFn: dbUserMapperFunc,
		},
		{
			name:     "nil list renders empty list",
			mapperFn: dbUserMapperFunc,
			input:    &corev1.Secret{},
			want:     []reconcile.Request{},
		},
		{
			name:     "empty input with proper empty list type renders empty list",
			mapperFn: dbUserMapperFunc,
			input:    &corev1.Secret{},
			want:     []reconcile.Request{},
		},
		{
			name:     "matching input credentials renders matching user",
			mapperFn: dbUserMapperFunc,
			input:    newTestSecret("matching-user-secret-ref"),
			objects:  []client.Object{newTestUser("matching-user")},
			want: []reconcile.Request{
				{NamespacedName: types.NamespacedName{
					Name:      "matching-user",
					Namespace: "ns",
				}},
			},
		},
	} {
		scheme := runtime.NewScheme()
		assert.NoError(t, corev1.AddToScheme(scheme))
		assert.NoError(t, akov2.AddToScheme(scheme))
		indexer := NewLocalCredentialsIndexer(
			AtlasDatabaseUserCredentialsIndex,
			&akov2.AtlasDatabaseUser{},
			zaptest.NewLogger(t),
		)
		t.Run(tc.name, func(t *testing.T) {
			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(tc.objects...).
				WithIndex(
					&akov2.AtlasDatabaseUser{},
					AtlasDatabaseUserCredentialsIndex,
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
		objs[i] = newTestUser(fmt.Sprintf("%s-%d", testUsername, i))
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
	for i := 0; i < 100; i++ {
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

func newTestUser(username string) *akov2.AtlasDatabaseUser {
	return &akov2.AtlasDatabaseUser{
		ObjectMeta: metav1.ObjectMeta{
			Name:      username,
			Namespace: "ns",
		},
		Spec: akov2.AtlasDatabaseUserSpec{
			LocalCredentialHolder: api.LocalCredentialHolder{
				ConnectionSecret: &api.LocalObjectReference{
					Name: fmt.Sprintf("%s-secret-ref", username),
				},
			},
		},
	}
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
	return CredentialsIndexMapperFunc[*akov2.AtlasDatabaseUserList](
		AtlasDatabaseUserCredentialsIndex,
		func() *akov2.AtlasDatabaseUserList { return &akov2.AtlasDatabaseUserList{} },
		DatabaseUserRequests,
		kubeClient,
		logger,
	)
}
