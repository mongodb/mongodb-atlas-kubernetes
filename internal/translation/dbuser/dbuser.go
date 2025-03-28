package dbuser

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

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
	_, _, err := dus.usersAPI.DeleteDatabaseUser(ctx, projectID, db, username).Execute()
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
