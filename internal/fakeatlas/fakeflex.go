// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fakeatlas

import (
	"context"
	"net/http"

	v20250312sdk "go.mongodb.org/atlas-sdk/v20250312009/admin"
)

type FakeFlexClustersApi struct {
	v20250312sdk.FlexClustersApi
	CreateFlexClusterWithParamsFunc func(ctx context.Context, params *v20250312sdk.CreateFlexClusterApiParams) (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error)
	// Store the last context and params for Execute to use
	lastCtx    context.Context
	lastParams *v20250312sdk.CreateFlexClusterApiParams
}

func (f *FakeFlexClustersApi) CreateFlexCluster(ctx context.Context, groupId string, flexClusterDescriptionCreate20241113 *v20250312sdk.FlexClusterDescriptionCreate20241113) v20250312sdk.CreateFlexClusterApiRequest {
	req := v20250312sdk.CreateFlexClusterApiRequest{
		ApiService: f,
	}
	return req
}

func (f *FakeFlexClustersApi) CreateFlexClusterWithParams(ctx context.Context, params *v20250312sdk.CreateFlexClusterApiParams) v20250312sdk.CreateFlexClusterApiRequest {
	// Store context and params for Execute to use
	f.lastCtx = ctx
	f.lastParams = params
	req := v20250312sdk.CreateFlexClusterApiRequest{
		ApiService: f,
	}
	return req
}

// CreateFlexClusterExecute is called by the SDK request's Execute method
func (f *FakeFlexClustersApi) CreateFlexClusterExecute(request v20250312sdk.CreateFlexClusterApiRequest) (*v20250312sdk.FlexClusterDescription20241113, *http.Response, error) {
	if f.CreateFlexClusterWithParamsFunc == nil || f.lastParams == nil {
		return nil, nil, nil
	}
	return f.CreateFlexClusterWithParamsFunc(f.lastCtx, f.lastParams)
}
