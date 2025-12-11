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
	v20250312sdk "go.mongodb.org/atlas-sdk/v20250312009/admin"
)

// AtlasClientBuilder is a builder for creating fake Atlas API clients.
type AtlasClientBuilder struct {
	client *v20250312sdk.APIClient
}

// NewAtlasClientBuilder creates a new AtlasClientBuilder.
func NewAtlasClientBuilder() *AtlasClientBuilder {
	return &AtlasClientBuilder{
		client: &v20250312sdk.APIClient{},
	}
}

// WithFakeFlexClusterClient configures the builder with a fake FlexCluster client.
func (b *AtlasClientBuilder) WithFakeFlexClusterClient() *AtlasClientBuilder {
	b.client.FlexClustersApi = &FakeFlexClustersApi{}
	return b
}

// Build returns the configured Atlas API client.
func (b *AtlasClientBuilder) Build() *v20250312sdk.APIClient {
	return b.client
}
