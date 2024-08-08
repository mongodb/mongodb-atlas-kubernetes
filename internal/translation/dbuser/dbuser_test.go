package dbuser

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

func TestNewAtlasDatabaseUsersService(t *testing.T) {
	ctx := context.Background()
	provider := &atlas.TestProvider{
		SdkClientFunc: func(_ *client.ObjectKey, _ *zap.SugaredLogger) (*admin.APIClient, string, error) {
			return &admin.APIClient{}, "", nil
		},
	}
	secretRef := &types.NamespacedName{}
	log := zap.S()
	users, err := NewAtlasDatabaseUsersService(ctx, provider, secretRef, log)
	require.NoError(t, err)
	assert.Equal(t, &AtlasUsers{}, users)
}

func TestFailedNewAtlasDatabaseUsersService(t *testing.T) {
	expectedErr := errors.New("fake error")
	ctx := context.Background()
	provider := &atlas.TestProvider{
		SdkClientFunc: func(_ *client.ObjectKey, _ *zap.SugaredLogger) (*admin.APIClient, string, error) {
			return nil, "", expectedErr
		},
	}
	secretRef := &types.NamespacedName{}
	log := zap.S()
	users, err := NewAtlasDatabaseUsersService(ctx, provider, secretRef, log)
	require.Nil(t, users)
	require.ErrorIs(t, err, expectedErr)
}
func TestAtlasUsersGet(t *testing.T) {
	ctx := context.Background()
	projectID := "project-id"
	db := "database"
	username := "test-user"

	notFoundErr := &admin.GenericOpenAPIError{}
	notFoundErr.SetModel(admin.ApiError{ErrorCode: pointer.MakePtr("USERNAME_NOT_FOUND")})

	tests := []struct {
		name         string
		setupMock    func(mockUsersAPI *mockadmin.DatabaseUsersApi)
		expectedUser *User
		expectedErr  error
	}{
		{
			name: "User found",
			setupMock: func(mockUsersAPI *mockadmin.DatabaseUsersApi) {
				expectedUser := &admin.CloudDatabaseUser{DatabaseName: db, GroupId: projectID, Username: username}
				mockUsersAPI.EXPECT().GetDatabaseUser(ctx, projectID, db, username).Return(
					admin.GetDatabaseUserApiRequest{ApiService: mockUsersAPI})

				mockUsersAPI.EXPECT().GetDatabaseUserExecute(admin.GetDatabaseUserApiRequest{ApiService: mockUsersAPI}).Return(
					expectedUser, &http.Response{StatusCode: http.StatusOK}, nil)
			},
			expectedUser: &User{AtlasDatabaseUserSpec: &akov2.AtlasDatabaseUserSpec{DatabaseName: db, Username: username}, ProjectID: projectID},
			expectedErr:  nil,
		},
		{
			name: "User not found",
			setupMock: func(mockUsersAPI *mockadmin.DatabaseUsersApi) {
				mockUsersAPI.EXPECT().GetDatabaseUser(ctx, projectID, db, username).Return(
					admin.GetDatabaseUserApiRequest{ApiService: mockUsersAPI})
				mockUsersAPI.EXPECT().GetDatabaseUserExecute(admin.GetDatabaseUserApiRequest{ApiService: mockUsersAPI}).Return(
					nil, &http.Response{StatusCode: http.StatusNotFound}, notFoundErr)
			},
			expectedUser: nil,
			expectedErr:  errors.Join(ErrorNotFound, notFoundErr),
		},
		{
			name: "API error",
			setupMock: func(mockUsersAPI *mockadmin.DatabaseUsersApi) {
				mockUsersAPI.EXPECT().GetDatabaseUser(ctx, projectID, db, username).Return(
					admin.GetDatabaseUserApiRequest{ApiService: mockUsersAPI})

				internalServerError := &admin.GenericOpenAPIError{}
				internalServerError.SetModel(admin.ApiError{ErrorCode: pointer.MakePtr("500")})
				mockUsersAPI.EXPECT().GetDatabaseUserExecute(admin.GetDatabaseUserApiRequest{ApiService: mockUsersAPI}).Return(
					nil, &http.Response{StatusCode: http.StatusInternalServerError}, errors.New("some error"))
			},
			expectedUser: nil,
			expectedErr:  fmt.Errorf("failed to get database user %q: %w", username, errors.New("some error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsersAPI := mockadmin.NewDatabaseUsersApi(t)
			tt.setupMock(mockUsersAPI)

			dus := &AtlasUsers{
				usersAPI: mockUsersAPI,
			}
			user, err := dus.Get(ctx, db, projectID, username)
			require.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedUser, user)
		})
	}
}
func TestAtlasUsersDelete(t *testing.T) {
	ctx := context.Background()
	projectID := "project-id"
	db := "database"
	username := "test-user"
	notFoundErr := &admin.GenericOpenAPIError{}
	notFoundErr.SetModel(admin.ApiError{ErrorCode: pointer.MakePtr("USER_NOT_FOUND")})
	tests := []struct {
		name        string
		setupMock   func(mockUsersAPI *mockadmin.DatabaseUsersApi)
		expectedErr error
	}{
		{
			name: "User successfully deleted",
			setupMock: func(mockUsersAPI *mockadmin.DatabaseUsersApi) {
				mockUsersAPI.EXPECT().DeleteDatabaseUser(ctx, projectID, db, username).Return(
					admin.DeleteDatabaseUserApiRequest{ApiService: mockUsersAPI})
				mockUsersAPI.EXPECT().DeleteDatabaseUserExecute(admin.DeleteDatabaseUserApiRequest{ApiService: mockUsersAPI}).
					Return(nil, &http.Response{StatusCode: http.StatusOK}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "User not found",
			setupMock: func(mockUsersAPI *mockadmin.DatabaseUsersApi) {
				mockUsersAPI.EXPECT().DeleteDatabaseUser(ctx, projectID, db, username).Return(
					admin.DeleteDatabaseUserApiRequest{ApiService: mockUsersAPI})

				mockUsersAPI.EXPECT().DeleteDatabaseUserExecute(admin.DeleteDatabaseUserApiRequest{ApiService: mockUsersAPI}).
					Return(nil, &http.Response{StatusCode: http.StatusNotFound}, notFoundErr)
			},
			expectedErr: errors.Join(ErrorNotFound, notFoundErr),
		},
		{
			name: "API error",
			setupMock: func(mockUsersAPI *mockadmin.DatabaseUsersApi) {
				mockUsersAPI.EXPECT().DeleteDatabaseUser(ctx, projectID, db, username).Return(
					admin.DeleteDatabaseUserApiRequest{ApiService: mockUsersAPI})

				internalServerError := &admin.GenericOpenAPIError{}
				internalServerError.SetModel(admin.ApiError{ErrorCode: pointer.MakePtr("500")})
				mockUsersAPI.EXPECT().DeleteDatabaseUserExecute(admin.DeleteDatabaseUserApiRequest{ApiService: mockUsersAPI}).
					Return(nil, &http.Response{StatusCode: http.StatusInternalServerError}, fmt.Errorf("some error"))
			},
			expectedErr: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsersAPI := mockadmin.NewDatabaseUsersApi(t)
			tt.setupMock(mockUsersAPI)
			dus := &AtlasUsers{
				usersAPI: mockUsersAPI,
			}
			err := dus.Delete(ctx, db, projectID, username)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}
