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

package dbuser

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312009/admin"
	"go.mongodb.org/atlas-sdk/v20250312009/mockadmin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

func TestAtlasUsersGet(t *testing.T) {
	ctx := context.Background()
	projectID := "project-id"
	db := "database"
	username := "test-user"

	notFoundErr := &admin.GenericOpenAPIError{}
	notFoundErr.SetModel(admin.ApiError{ErrorCode: "USERNAME_NOT_FOUND"})

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
			expectedUser: &User{
				ProjectID: projectID,
				AtlasDatabaseUserSpec: &akov2.AtlasDatabaseUserSpec{
					DatabaseName: db,
					Username:     username,
					Scopes:       []akov2.ScopeSpec{},
				},
			},
			expectedErr: nil,
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
				internalServerError.SetModel(admin.ApiError{ErrorCode: "500"})
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
	notFoundErr.SetModel(admin.ApiError{ErrorCode: "USER_NOT_FOUND"})
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
					Return(&http.Response{StatusCode: http.StatusOK}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "User not found",
			setupMock: func(mockUsersAPI *mockadmin.DatabaseUsersApi) {
				mockUsersAPI.EXPECT().DeleteDatabaseUser(ctx, projectID, db, username).Return(
					admin.DeleteDatabaseUserApiRequest{ApiService: mockUsersAPI})

				mockUsersAPI.EXPECT().DeleteDatabaseUserExecute(admin.DeleteDatabaseUserApiRequest{ApiService: mockUsersAPI}).
					Return(&http.Response{StatusCode: http.StatusNotFound}, notFoundErr)
			},
			expectedErr: errors.Join(ErrorNotFound, notFoundErr),
		},
		{
			name: "API error",
			setupMock: func(mockUsersAPI *mockadmin.DatabaseUsersApi) {
				mockUsersAPI.EXPECT().DeleteDatabaseUser(ctx, projectID, db, username).Return(
					admin.DeleteDatabaseUserApiRequest{ApiService: mockUsersAPI})

				internalServerError := &admin.GenericOpenAPIError{}
				internalServerError.SetModel(admin.ApiError{ErrorCode: "500"})
				mockUsersAPI.EXPECT().DeleteDatabaseUserExecute(admin.DeleteDatabaseUserApiRequest{ApiService: mockUsersAPI}).
					Return(&http.Response{StatusCode: http.StatusInternalServerError}, fmt.Errorf("some error"))
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
