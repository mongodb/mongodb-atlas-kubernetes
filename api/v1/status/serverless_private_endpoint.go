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

package status

type ServerlessPrivateEndpoint struct {
	// ID is the identifier of the Serverless PrivateLink Service.
	ID string `json:"_id,omitempty"`
	// CloudProviderEndpointID is the identifier of the cloud provider endpoint.
	CloudProviderEndpointID string `json:"cloudProviderEndpointId,omitempty"`
	// Name is the name of the Serverless PrivateLink Service. Should be unique.
	Name string `json:"name,omitempty"`
	// EndpointServiceName is the name of the PrivateLink endpoint service in AWS. Returns null while the endpoint service is being created.
	EndpointServiceName string `json:"endpointServiceName,omitempty"`
	// ErrorMessage is the error message if the Serverless PrivateLink Service failed to create or connect.
	ErrorMessage string `json:"errorMessage,omitempty"`
	// Status of the AWS Serverless PrivateLink connection.
	Status string `json:"status,omitempty"`
	// ProviderName is human-readable label that identifies the cloud provider. Values include AWS or AZURE.
	ProviderName string `json:"providerName,omitempty"`
	// PrivateEndpointIPAddress is the IPv4 address of the private endpoint in your Azure VNet that someone added to this private endpoint service.
	PrivateEndpointIPAddress string `json:"privateEndpointIpAddress,omitempty"`
	// PrivateLinkServiceResourceID is the root-relative path that identifies the Azure Private Link Service that MongoDB Cloud manages. MongoDB Cloud returns null while it creates the endpoint service.
	PrivateLinkServiceResourceID string `json:"privateLinkServiceResourceId,omitempty"`
}
