package datafederation

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
)

var (
	ErrorNotFound = errors.New("data federation not found")
)

type DataFederationService interface {
	Get(ctx context.Context, projectID, name string) (*DataFederation, error)
	Create(ctx context.Context, df *DataFederation) error
	Update(ctx context.Context, df *DataFederation) error
	Delete(ctx context.Context, projectID, name string) error
}

type AtlasDataFederationService struct {
	api admin.DataFederationApi
}

func NewAtlasDataFederationService(ctx context.Context, provider atlas.Provider, secretRef *types.NamespacedName, log *zap.SugaredLogger) (*AtlasDataFederationService, error) {
	client, err := translation.NewVersionedClient(ctx, provider, secretRef, log)
	if err != nil {
		return nil, fmt.Errorf("failed to create versioned client: %w", err)
	}
	return &AtlasDataFederationService{client.DataFederationApi}, nil
}

func (dfs *AtlasDataFederationService) Get(ctx context.Context, projectID, name string) (*DataFederation, error) {
	atlasDataFederation, resp, err := dfs.api.GetFederatedDatabase(ctx, projectID, name).Execute()

	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return nil, errors.Join(ErrorNotFound, err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get data federation database %q: %w", name, err)
	}
	return fromAtlas(atlasDataFederation)
}

func (dfs *AtlasDataFederationService) Create(ctx context.Context, df *DataFederation) error {
	atlasDataFederation := toAtlas(df)
	_, _, err := dfs.api.
		CreateFederatedDatabase(ctx, df.ProjectID, atlasDataFederation).
		SkipRoleValidation(df.SkipRoleValidation).
		Execute()
	if err != nil {
		return fmt.Errorf("failed to create data federation database %q: %w", df.ProjectID, err)
	}
	return nil
}

func (dfs *AtlasDataFederationService) Update(ctx context.Context, df *DataFederation) error {
	atlasDataFederation := toAtlas(df)
	_, _, err := dfs.api.
		UpdateFederatedDatabase(ctx, df.ProjectID, df.Name, atlasDataFederation).
		SkipRoleValidation(df.SkipRoleValidation).
		Execute()
	if err != nil {
		return fmt.Errorf("failed to update data federation database %q: %w", df.ProjectID, err)
	}
	return nil
}

func (dfs *AtlasDataFederationService) Delete(ctx context.Context, projectID, name string) error {
	_, resp, err := dfs.api.DeleteFederatedDatabase(ctx, projectID, name).Execute()
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return errors.Join(ErrorNotFound, err)
	}
	if err != nil {
		return fmt.Errorf("failed to delete data federation database %q: %w", projectID, err)
	}
	return nil
}
