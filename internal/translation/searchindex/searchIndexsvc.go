package searchindex

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
	// ErrNotFound means an resource is missing
	ErrNotFound = fmt.Errorf("not found")
)

type AtlasSearchIdxService interface {
	GetIndex(ctx context.Context, projectID, clusterName, indexName, indexID string) (*SearchIndex, error)
	CreateIndex(ctx context.Context, projectID, clusterName string, index *SearchIndex) (*SearchIndex, error)
	DeleteIndex(ctx context.Context, projectID, clusterName, indexID string) error
	UpdateIndex(ctx context.Context, projectID, clusterName string, index *SearchIndex) (*SearchIndex, error)
}

type ProductionAtlasSearch struct {
	searchAPI admin.AtlasSearchApi
}

func NewAtlasDatabaseUsersService(ctx context.Context, provider atlas.Provider, secretRef *types.NamespacedName, log *zap.SugaredLogger) (*ProductionAtlasSearch, error) {
	client, err := translation.NewVersionedClient(ctx, provider, secretRef, log)
	if err != nil {
		return nil, err
	}
	return NewProductionAtlasSearch(client.AtlasSearchApi), nil
}

func NewProductionAtlasSearch(api admin.AtlasSearchApi) *ProductionAtlasSearch {
	return &ProductionAtlasSearch{searchAPI: api}
}

func (si *ProductionAtlasSearch) GetIndex(ctx context.Context, projectID, clusterName, indexName, indexID string) (*SearchIndex, error) {
	resp, httpResp, err := si.searchAPI.GetAtlasSearchIndex(ctx, projectID, clusterName, indexID).Execute()
	if err != nil {
		if httpResp.StatusCode == http.StatusNotFound {
			return nil, errors.Join(err, ErrNotFound)
		}
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("received an empty index. Index: %s(%s). Status code: %d, E: %w",
			indexName, indexID, httpResp.StatusCode, err)
	}
	stateInAtlas, err := fromAtlas(*resp)
	if err != nil {
		return nil, fmt.Errorf("unable to convert index to AKO. Index: %s(%s, E: %w",
			indexName, indexID, err)
	}
	return stateInAtlas, nil
}

func (si *ProductionAtlasSearch) CreateIndex(ctx context.Context, projectID, clusterName string, index *SearchIndex) (*SearchIndex, error) {
	atlasIndex, err := index.toAtlas()
	if err != nil {
		return nil, err
	}
	resp, httpResp, err := si.searchAPI.CreateAtlasSearchIndex(ctx, projectID, clusterName, atlasIndex).Execute()
	if err != nil || httpResp.StatusCode != http.StatusCreated && httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create index: %w, status: %d", err, httpResp.StatusCode)
	}
	if resp == nil {
		return nil, fmt.Errorf("returned an empty index as a result of creation")
	}
	akoIndex, err := fromAtlas(*resp)
	if err != nil {
		return nil, fmt.Errorf("unable to convert index to AKO: %w", err)
	}
	return akoIndex, nil
}

func (si *ProductionAtlasSearch) DeleteIndex(ctx context.Context, projectID, clusterName, indexID string) error {
	_, resp, err := si.searchAPI.DeleteAtlasSearchIndex(ctx, projectID, clusterName, indexID).Execute()
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusNotFound || err != nil {
		return fmt.Errorf("failed to delete index: %w, status: %d", err, resp.StatusCode)
	}
	return nil
}

func (si *ProductionAtlasSearch) UpdateIndex(ctx context.Context, projectID, clusterName string, index *SearchIndex) (*SearchIndex, error) {
	atlasIndex, err := index.toAtlas()
	if err != nil {
		return nil, fmt.Errorf("unable to convert index to AKO: %w", err)
	}
	resp, httpResp, err := si.searchAPI.UpdateAtlasSearchIndex(ctx, projectID, clusterName, index.GetID(), atlasIndex).Execute()
	if httpResp.StatusCode != http.StatusCreated && httpResp.StatusCode != http.StatusOK || err != nil {
		return nil, fmt.Errorf("failed to update index: %w, status: %d", err, httpResp.StatusCode)
	}
	if resp == nil {
		return nil, fmt.Errorf("update returned an empty index: %w", err)
	}
	akoIndex, err := fromAtlas(*resp)
	if err != nil {
		return nil, fmt.Errorf("failed to convert updated index to AKO: %w", err)
	}
	return akoIndex, nil
}
