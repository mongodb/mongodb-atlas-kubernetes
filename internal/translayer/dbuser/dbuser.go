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

type Service struct {
	admin.DatabaseUsersApi
}

func NewService(ctx context.Context, provider atlas.Provider, secretRef *types.NamespacedName, log *zap.SugaredLogger) (*Service, error) {
	client, err := translayer.NewVersionedClient(ctx, provider, secretRef, log)
	if err != nil {
		return nil, err
	}
	return NewFromDBUserAPI(client.DatabaseUsersApi), nil
}

func NewFromDBUserAPI(api admin.DatabaseUsersApi) *Service {
	return &Service{DatabaseUsersApi: api}
}

func (dus *Service) Get(ctx context.Context, db, projectID, username string) (*User, error) {
	atlasDBUser, _, err := dus.GetDatabaseUser(ctx, projectID, db, username).Execute()
	if err != nil {
		if admin.IsErrorCode(err, atlas.UsernameNotFound) {
			return nil, errors.Join(ErrorNotFound, err)
		}
		return nil, err
	}
	return fromAtlas(atlasDBUser)
}

func (dus *Service) Delete(ctx context.Context, db, projectID, username string) error {
	_, _, err := dus.DeleteDatabaseUser(ctx, projectID, db, username).Execute()
	if err != nil {
		if admin.IsErrorCode(err, atlas.UsernameNotFound) {
			return errors.Join(ErrorNotFound, err)
		}
		return err
	}
	return nil
}

func (dus *Service) Create(ctx context.Context, au *User) error {
	u, err := toAtlas(au)
	if err != nil {
		return err
	}
	_, _, err = dus.CreateDatabaseUser(ctx, au.ProjectID, u).Execute()
	return err
}

func (dus *Service) Update(ctx context.Context, au *User) error {
	u, err := toAtlas(au)
	if err != nil {
		return err
	}
	_, _, err = dus.UpdateDatabaseUser(ctx, au.ProjectID, au.DatabaseName, au.Username, u).Execute()
	return err
}
