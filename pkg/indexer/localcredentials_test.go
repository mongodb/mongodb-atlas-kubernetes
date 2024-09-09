package indexer

import (
	"context"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
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
		name    string
		list    api.ReconciliableList
		objects []client.Object
		input   client.Object
		want    []reconcile.Request
	}{
		{
			name: "nil input & list renders nil",
		},
		{
			name:  "nil list renders nil",
			input: &corev1.Secret{},
		},
		{
			name:  "empty input with proper empty list type renders empty list",
			list:  &akov2.AtlasDatabaseUserList{},
			input: &corev1.Secret{},
			want:  []reconcile.Request{},
		},
		{
			name: "matching input credentials renders matching user",
			list: &akov2.AtlasDatabaseUserList{},
			input: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret-ref",
					Namespace: "ns",
				},
			},
			objects: []client.Object{
				&akov2.AtlasDatabaseUser{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "matching-user",
						Namespace: "ns",
					},
					Spec: akov2.AtlasDatabaseUserSpec{
						LocalCredentialHolder: api.LocalCredentialHolder{
							ConnectionSecret: &api.LocalObjectReference{
								Name: "secret-ref",
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
			fn := CredentialsIndexMapperFunc(
				AtlasDatabaseUserCredentialsIndex,
				tc.list,
				fakeClient,
				zaptest.NewLogger(t).Sugar(),
			)
			result := fn(context.Background(), tc.input)
			assert.Equal(t, tc.want, result)
		})
	}
}
