package dbuser

import (
	"context"
	"errors"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translayer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
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

type ProductionAtlasUsers struct {
	usersAPI admin.DatabaseUsersApi
}

func NewAtlasDatabaseUsersService(ctx context.Context, provider atlas.Provider, secretRef *types.NamespacedName, log *zap.SugaredLogger) (*ProductionAtlasUsers, error) {
	client, err := translayer.NewVersionedClient(ctx, provider, secretRef, log)
	if err != nil {
		return nil, err
	}
	return NewProductionAtlasUsers(client.DatabaseUsersApi), nil
}

func NewProductionAtlasUsers(api admin.DatabaseUsersApi) *ProductionAtlasUsers {
	return &ProductionAtlasUsers{usersAPI: api}
}

func (dus *ProductionAtlasUsers) Get(ctx context.Context, db, projectID, username string) (*User, error) {
	atlasDBUser, _, err := dus.usersAPI.GetDatabaseUser(ctx, projectID, db, username).Execute()
	if err != nil {
		if admin.IsErrorCode(err, atlas.UsernameNotFound) {
			return nil, errors.Join(ErrorNotFound, err)
		}
		return nil, err
	}
	return fromAtlas(atlasDBUser)
}

func (dus *ProductionAtlasUsers) Delete(ctx context.Context, db, projectID, username string) error {
	_, _, err := dus.usersAPI.DeleteDatabaseUser(ctx, projectID, db, username).Execute()
	if err != nil {
		if admin.IsErrorCode(err, atlas.UsernameNotFound) {
			return errors.Join(ErrorNotFound, err)
		}
		return err
	}
	return nil
}

func (dus *ProductionAtlasUsers) Create(ctx context.Context, au *User) error {
	u, err := toAtlas(au)
	if err != nil {
		return err
	}
	_, _, err = dus.usersAPI.CreateDatabaseUser(ctx, au.ProjectID, u).Execute()
	return err
}

func (dus *ProductionAtlasUsers) Update(ctx context.Context, au *User) error {
	u, err := toAtlas(au)
	if err != nil {
		return err
	}
	_, _, err = dus.usersAPI.UpdateDatabaseUser(ctx, au.ProjectID, au.DatabaseName, au.Username, u).Execute()
	return err
}
