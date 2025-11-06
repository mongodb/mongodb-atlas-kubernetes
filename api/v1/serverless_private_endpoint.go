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

package v1

// ServerlessPrivateEndpoint configures private endpoints for the Serverless instances.
// DEPRECATED: Serverless private endpoints are deprecated. See https://dochub.mongodb.org/core/atlas-flex-migration for details.
type ServerlessPrivateEndpoint struct {
	// Name is the name of the Serverless PrivateLink Service. Should be unique.
	Name string `json:"name,omitempty"`
	// CloudProviderEndpointID is the identifier of the cloud provider endpoint.
	CloudProviderEndpointID string `json:"cloudProviderEndpointID,omitempty"`
	// PrivateEndpointIPAddress is the IPv4 address of the private endpoint in your Azure VNet that someone added to this private endpoint service.
	PrivateEndpointIPAddress string `json:"privateEndpointIpAddress,omitempty"`
}

// IsInitialState pe initially should be empty except for comment
func (in *ServerlessPrivateEndpoint) IsInitialState() bool {
	return in.Name != "" && in.CloudProviderEndpointID == "" && in.PrivateEndpointIPAddress == ""
}
