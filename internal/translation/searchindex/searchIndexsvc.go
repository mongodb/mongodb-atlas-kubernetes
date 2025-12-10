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

package searchindex

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20250312010/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/httputil"
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

type SearchIndexes struct {
	searchAPI admin.AtlasSearchApi
}

func NewSearchIndexes(api admin.AtlasSearchApi) *SearchIndexes {
	return &SearchIndexes{searchAPI: api}
}

func (si *SearchIndexes) GetIndex(ctx context.Context, projectID, clusterName, indexName, indexID string) (*SearchIndex, error) {
	resp, httpResp, err := si.searchAPI.GetClusterSearchIndex(ctx, projectID, clusterName, indexID).Execute()
	statusCode := httputil.StatusCode(httpResp)
	if err != nil {
		if statusCode == http.StatusNotFound {
			return nil, errors.Join(err, ErrNotFound)
		}
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("got empty index %s(%s), status code %d: %w",
			indexName, indexID, statusCode, err)
	}
	stateInAtlas, err := fromAtlas(*resp)
	if err != nil {
		return nil, fmt.Errorf("unable to convert index %s(%s): %w",
			indexName, indexID, err)
	}
	return stateInAtlas, nil
}

func (si *SearchIndexes) CreateIndex(ctx context.Context, projectID, clusterName string, index *SearchIndex) (*SearchIndex, error) {
	atlasIndex, err := index.toAtlasCreateView()
	if err != nil {
		return nil, err
	}
	resp, httpResp, err := si.searchAPI.CreateClusterSearchIndex(ctx, projectID, clusterName, atlasIndex).Execute()
	statusCode := httputil.StatusCode(httpResp)
	respNotOK := (statusCode != http.StatusCreated && statusCode != http.StatusOK)
	if err != nil || respNotOK {
		return nil, fmt.Errorf("failed to create index, status code %d: %w", statusCode, err)
	}
	if resp == nil {
		return nil, errors.New("empty response when creating index")
	}
	akoIndex, err := fromAtlas(*resp)
	if err != nil {
		return nil, fmt.Errorf("error converting index: %w", err)
	}
	return akoIndex, nil
}

func (si *SearchIndexes) DeleteIndex(ctx context.Context, projectID, clusterName, indexID string) error {
	resp, err := si.searchAPI.DeleteClusterSearchIndex(ctx, projectID, clusterName, indexID).Execute()
	statusCode := httputil.StatusCode(resp)
	if statusCode != http.StatusAccepted && statusCode != http.StatusNotFound || err != nil {
		return fmt.Errorf("error deleting index, status code %d: %w", statusCode, err)
	}
	return nil
}

func (si *SearchIndexes) UpdateIndex(ctx context.Context, projectID, clusterName string, index *SearchIndex) (*SearchIndex, error) {
	atlasIndex, err := index.toAtlasUpdateView()
	if err != nil {
		return nil, fmt.Errorf("error converting index: %w", err)
	}
	resp, httpResp, err := si.searchAPI.UpdateClusterSearchIndex(ctx, projectID, clusterName, index.GetID(), atlasIndex).Execute()
	statusCode := httputil.StatusCode(httpResp)
	respNotOK := statusCode != http.StatusCreated && statusCode != http.StatusOK
	if respNotOK || err != nil {
		return nil, fmt.Errorf("error updating index, status code %d: %w", statusCode, err)
	}
	if resp == nil {
		return nil, fmt.Errorf("update returned an empty index: %w", err)
	}
	akoIndex, err := fromAtlas(*resp)
	if err != nil {
		return nil, fmt.Errorf("error converting index: %w", err)
	}
	return akoIndex, nil
}
