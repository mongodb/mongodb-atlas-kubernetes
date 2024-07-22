package atlasdatabaseuser

import (
	"context"
	"errors"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/indexer"
	"go.uber.org/zap/zaptest"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
	"time"

	mocked "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/dbuser"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const (
	testProject = "project"

	testProjectID = "12345"

	testDatabase = "db"

	testUsername = "user"
)

var (
	errRandom = errors.New("random error")
)

func TestHandleDeletion(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	utilruntime.Must(corev1.AddToScheme(scheme))
	utilruntime.Must(akov2.AddToScheme(scheme))
	fakeClient := fake.NewClientBuilder().WithScheme(scheme).Build()
	user := defaultTestUser()
	require.NoError(t, fakeClient.Create(ctx, user))
	defer fakeClient.Delete(ctx, user)
	log := zap.S()
	r := &AtlasDatabaseUserReconciler{
		Client: fakeClient,
		Log:    log,
	}
	for _, tc := range []struct {
		title          string
		skipDeletion   bool
		service        dbuser.AtlasUsersService
		expectedOk     bool
		expectedResult workflow.Result
	}{
		{
			title:          "User without deletion timestamp is not deleted",
			skipDeletion:   true,
			expectedOk:     false,
			expectedResult: workflow.OK(),
		},

		{
			title:          "Ready user gets deleted properly",
			service:        fakeUserDeletion(ctx, testDatabase, testProjectID, testUsername, nil),
			expectedOk:     true,
			expectedResult: workflow.OK(),
		},

		{
			title:          "Missing user is already deleted",
			service:        fakeUserDeletion(ctx, testDatabase, testProjectID, testUsername, dbuser.ErrorNotFound),
			expectedOk:     true,
			expectedResult: workflow.OK(),
		},

		{
			title:          "Fails to delete user when returned error is unexpected",
			service:        fakeUserDeletion(ctx, testDatabase, testProjectID, testUsername, errRandom),
			expectedOk:     true,
			expectedResult: workflow.Terminate(workflow.DatabaseUserNotDeletedInAtlas, errRandom.Error()),
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			if !tc.skipDeletion {
				user.DeletionTimestamp = pointer.MakePtr(metav1.NewTime(time.Now()))
			}
			done, result := r.handleDeletion(ctx, user, defaultTestProject(), tc.service, log)
			assert.Equal(t, tc.expectedOk, done)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func fakeUserDeletion(ctx context.Context, db, projectID, username string, err error) *mocked.AtlasUsersServiceMock {
	return withFakeUserDeletion(&mocked.AtlasUsersServiceMock{}, ctx, db, projectID, username, err)
}

func withFakeUserDeletion(service *mocked.AtlasUsersServiceMock, ctx context.Context, db, projectID, username string, err error) *mocked.AtlasUsersServiceMock {
	service.EXPECT().Delete(ctx, db, projectID, username).Return(err)
	return service
}

func defaultTestProject() *akov2.AtlasProject {
	return &akov2.AtlasProject{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Name: testProject},
		Status: status.AtlasProjectStatus{
			ID: testProjectID,
		},
	}
}

func defaultTestUser() *akov2.AtlasDatabaseUser {
	return &akov2.AtlasDatabaseUser{
		ObjectMeta: metav1.ObjectMeta{Name: testUsername},
		Spec: akov2.AtlasDatabaseUserSpec{
			DatabaseName: testDatabase,
			Username:     testUsername,
		},
	}
}

func TestFindAtlasDatabaseUserForSecret(t *testing.T) {
	for _, tc := range []struct {
		name     string
		obj      client.Object
		initObjs []client.Object
		want     []reconcile.Request
	}{
		{
			name: "wrong type",
			obj:  &akov2.AtlasProject{},
			want: nil,
		},
		{
			name: "same namespace",
			obj: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: "secret", Namespace: "ns"},
			},
			initObjs: []client.Object{
				&akov2.AtlasDatabaseUser{
					ObjectMeta: metav1.ObjectMeta{Name: "user1", Namespace: "ns"},
					Spec: akov2.AtlasDatabaseUserSpec{
						PasswordSecret: &common.ResourceRef{Name: "secret"},
					},
				},
			},
			want: []reconcile.Request{
				{NamespacedName: types.NamespacedName{Name: "user1", Namespace: "ns"}},
			},
		},
		{
			name: "different namespace",
			obj: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: "name", Namespace: "ns2"},
			},
			initObjs: []client.Object{
				&akov2.AtlasDatabaseUser{
					ObjectMeta: metav1.ObjectMeta{Name: "user1", Namespace: "ns"},
					Spec: akov2.AtlasDatabaseUserSpec{
						PasswordSecret: &common.ResourceRef{Name: "secret"},
					},
				},
			},
			want: []reconcile.Request{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			indexer := indexer.NewAtlasDatabaseUserBySecretsIndexer(zaptest.NewLogger(t))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tc.initObjs...).
				WithIndex(indexer.Object(), indexer.Name(), indexer.Keys).
				Build()
			reconciler := &AtlasDatabaseUserReconciler{
				Log:    zaptest.NewLogger(t).Sugar(),
				Client: k8sClient,
			}
			got := reconciler.findAtlasDatabaseUserForSecret(context.Background(), tc.obj)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("want reconcile requests: %v, got %v", got, tc.want)
			}
		})
	}
}
