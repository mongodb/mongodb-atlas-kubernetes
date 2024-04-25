package dbuser

import (
	"context"
	"errors"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translayer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
)

type Service struct {
	mongodbatlas.DatabaseUsersService
}

func NewService(ctx context.Context, provider atlas.Provider, secretRef *types.NamespacedName, log *zap.SugaredLogger) (*Service, error) {
	client, err := translayer.NewLegacyClient(ctx, provider, secretRef, log)
	if err != nil {
		return nil, err
	}
	return &Service{DatabaseUsersService: client.DatabaseUsers}, nil
}

func (dus *Service) Get(ctx context.Context, db, projectID, username string) (*User, error) {
	atlasDBUser, _, err := dus.DatabaseUsersService.Get(ctx, db, projectID, username)
	if err != nil {
		var apiError *mongodbatlas.ErrorResponse
		if errors.As(err, &apiError) && apiError.ErrorCode == atlas.UsernameNotFound {
			return nil, nil
		}
		return nil, err
	}
	return toK8sDatabaseUser(atlasDBUser)
}

func (dus *Service) Delete(ctx context.Context, db, projectID, username string) (bool, error) {
	_, err := dus.DatabaseUsersService.Delete(ctx, db, projectID, username)
	if err != nil {
		var apiError *mongodbatlas.ErrorResponse
		if errors.As(err, &apiError) && apiError.ErrorCode != atlas.UsernameNotFound {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func (dus *Service) Create(ctx context.Context, db string, au *User) error {
	_, _, err := dus.DatabaseUsersService.Create(ctx, db, toAtlas(au))
	return err
}

func (dus *Service) Update(ctx context.Context, db, projectID string, au *User) error {
	_, _, err := dus.DatabaseUsersService.Update(ctx, db, projectID, toAtlas(au))
	return err
}
