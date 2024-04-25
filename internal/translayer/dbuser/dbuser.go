package dbuser

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translayer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
)

type Service struct {
	admin.DatabaseUsersApi
}

func NewService(ctx context.Context, provider atlas.Provider, secretRef *types.NamespacedName, log *zap.SugaredLogger) (*Service, error) {
	client, err := translayer.NewVersionedClient(ctx, provider, secretRef, log)
	if err != nil {
		return nil, err
	}
	return &Service{DatabaseUsersApi: client.DatabaseUsersApi}, nil
}

func (dus *Service) Get(ctx context.Context, db, projectID, username string) (*User, error) {
	atlasDBUser, _, err := dus.GetDatabaseUser(ctx, projectID, db, username).Execute()
	if err != nil {
		if admin.IsErrorCode(err, atlas.UsernameNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toK8sDatabaseUser(atlasDBUser)
}

func (dus *Service) Delete(ctx context.Context, db, projectID, username string) (bool, error) {
	_, _, err := dus.DeleteDatabaseUser(ctx, projectID, db, username).Execute()
	if err != nil {
		if admin.IsErrorCode(err, atlas.UsernameNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (dus *Service) Create(ctx context.Context, au *User) error {
	_, _, err := dus.CreateDatabaseUser(ctx, au.ProjectID, toAtlas(au)).Execute()
	return err
}

func (dus *Service) Update(ctx context.Context, au *User) error {
	_, _, err := dus.UpdateDatabaseUser(
		ctx, au.ProjectID, au.DatabaseName, au.Username, toAtlas(au)).Execute()
	return err
}
