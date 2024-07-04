package atlasdatabaseuser

import (
	"context"
	"errors"
	"testing"
	"time"

	translationmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
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
			service:        fakeDeletionService(ctx, testDatabase, testProjectID, testUsername, nil),
			expectedOk:     true,
			expectedResult: workflow.OK(),
		},
		{
			title:          "Missing user is already deleted",
			service:        fakeDeletionService(ctx, testDatabase, testProjectID, testUsername, dbuser.ErrorNotFound),
			expectedOk:     true,
			expectedResult: workflow.OK(),
		},
		{
			title:          "Fails to delete user when returned error is unexpected",
			service:        fakeDeletionService(ctx, testDatabase, testProjectID, testUsername, errRandom),
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

func fakeDeletionService(ctx context.Context, db, projectID, username string, err error) dbuser.AtlasUsersService {
	service := translationmock.AtlasUsersServiceMock{}
	service.EXPECT().Delete(ctx, db, projectID, username).Return(err)
	return &service
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
