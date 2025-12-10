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
	"fmt"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20250312010/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/paging"
)

type DatafederationPrivateEndpointService interface {
	List(ctx context.Context, projectID string) ([]*DatafederationPrivateEndpointEntry, error)
	Create(context.Context, *DatafederationPrivateEndpointEntry) error
	Delete(context.Context, *DatafederationPrivateEndpointEntry) error
}

type DatafederationPrivateEndpoints struct {
	api admin.DataFederationApi
}

func NewDatafederationPrivateEndpoint(api admin.DataFederationApi) *DatafederationPrivateEndpoints {
	return &DatafederationPrivateEndpoints{api: api}
}

func (d *DatafederationPrivateEndpoints) List(ctx context.Context, projectID string) ([]*DatafederationPrivateEndpointEntry, error) {
	results, err := paging.ListAll(ctx, func(ctx context.Context, pageNum int) (paging.Response[admin.PrivateNetworkEndpointIdEntry], *http.Response, error) {
		return d.api.ListPrivateEndpointIds(ctx, projectID).PageNum(pageNum).Execute()
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list data federation private endpoints from Atlas: %w", err)
	}

	return endpointsFromAtlas(results, projectID)
}

func (d *DatafederationPrivateEndpoints) Create(ctx context.Context, aep *DatafederationPrivateEndpointEntry) error {
	ep := endpointToAtlas(aep)
	_, _, err := d.api.CreatePrivateEndpointId(ctx, aep.ProjectID, ep).Execute()
	if err != nil {
		return fmt.Errorf("failed to create data federation private endpoint: %w", err)
	}
	return nil
}

func (d *DatafederationPrivateEndpoints) Delete(ctx context.Context, aep *DatafederationPrivateEndpointEntry) error {
	_, err := d.api.DeletePrivateEndpointId(ctx, aep.ProjectID, aep.EndpointID).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete data federation private endpoint: %w", err)
	}
	return nil
}
