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

package datafederation

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20250312011/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/httputil"
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

func NewAtlasDataFederation(api admin.DataFederationApi) *AtlasDataFederationService {
	return &AtlasDataFederationService{api: api}
}

func (dfs *AtlasDataFederationService) Get(ctx context.Context, projectID, name string) (*DataFederation, error) {
	atlasDataFederation, resp, err := dfs.api.GetDataFederation(ctx, projectID, name).Execute()

	if httputil.StatusCode(resp) == http.StatusNotFound {
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
		CreateDataFederation(ctx, df.ProjectID, atlasDataFederation).
		Execute()
	if err != nil {
		return fmt.Errorf("failed to create data federation database %q: %w", df.ProjectID, err)
	}
	return nil
}

func (dfs *AtlasDataFederationService) Update(ctx context.Context, df *DataFederation) error {
	atlasDataFederation := toAtlas(df)
	_, _, err := dfs.api.
		UpdateDataFederation(ctx, df.ProjectID, df.Name, atlasDataFederation).
		// false is the default for creation, so we have to respect it for updates as well.
		SkipRoleValidation(false).
		Execute()
	if err != nil {
		return fmt.Errorf("failed to update data federation database %q: %w", df.ProjectID, err)
	}
	return nil
}

func (dfs *AtlasDataFederationService) Delete(ctx context.Context, projectID, name string) error {
	resp, err := dfs.api.DeleteDataFederation(ctx, projectID, name).Execute()
	if httputil.StatusCode(resp) == http.StatusNotFound {
		return errors.Join(ErrorNotFound, err)
	}
	if err != nil {
		return fmt.Errorf("failed to delete data federation database %q: %w", projectID, err)
	}
	return nil
}
