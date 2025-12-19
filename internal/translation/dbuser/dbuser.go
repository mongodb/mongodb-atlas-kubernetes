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

	"go.mongodb.org/atlas-sdk/v20250312011/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
)

var (
	// ErrorNotFound is returned when an database user was not found
	ErrorNotFound = errors.New("database user not found")
)

type AtlasUsersService interface {
	Get(ctx context.Context, db, projectID, username string) (*User, error)
	Delete(ctx context.Context, db, projectID, username string) error
	Create(ctx context.Context, au *User) error
	Update(ctx context.Context, au *User) error
}

type AtlasUsers struct {
	usersAPI admin.DatabaseUsersApi
}

func NewAtlasUsers(api admin.DatabaseUsersApi) *AtlasUsers {
	return &AtlasUsers{usersAPI: api}
}

func (dus *AtlasUsers) Get(ctx context.Context, db, projectID, username string) (*User, error) {
	atlasDBUser, _, err := dus.usersAPI.GetDatabaseUser(ctx, projectID, db, username).Execute()
	if err != nil {
		if admin.IsErrorCode(err, atlas.UsernameNotFound) {
			return nil, errors.Join(ErrorNotFound, err)
		}
		return nil, fmt.Errorf("failed to get database user %q: %w", username, err)
	}
	return fromAtlas(atlasDBUser)
}

func (dus *AtlasUsers) Delete(ctx context.Context, db, projectID, username string) error {
	_, err := dus.usersAPI.DeleteDatabaseUser(ctx, projectID, db, username).Execute()
	if err != nil {
		if admin.IsErrorCode(err, atlas.UserNotfound) {
			return errors.Join(ErrorNotFound, err)
		}
		return err
	}
	return nil
}

func (dus *AtlasUsers) Create(ctx context.Context, au *User) error {
	u, err := toAtlas(au)
	if err != nil {
		return err
	}
	_, _, err = dus.usersAPI.CreateDatabaseUser(ctx, au.ProjectID, u).Execute()
	return err
}

func (dus *AtlasUsers) Update(ctx context.Context, au *User) error {
	u, err := toAtlas(au)
	if err != nil {
		return err
	}
	_, _, err = dus.usersAPI.UpdateDatabaseUser(ctx, au.ProjectID, au.DatabaseName, au.Username, u).Execute()
	return err
}
