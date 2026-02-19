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

import (
	"go.mongodb.org/atlas-sdk/v20250312014/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

// PrivateEndpoint is a list of Private Endpoints configured for the current Project.
// Deprecated: Migrate to the AtlasPrivateEndpoint Custom Resource in accordance with the migration guide
// at https://www.mongodb.com/docs/atlas/operator/current/migrate-parameter-to-resource/#std-label-ak8so-migrate-ptr
type PrivateEndpoint struct {
	// Cloud provider for which you want to retrieve a private endpoint service. Atlas accepts AWS, GCP, or AZURE.
	// +kubebuilder:validation:Enum=AWS;GCP;AZURE;TENANT
	Provider provider.ProviderName `json:"provider"`
	// Cloud provider region for which you want to create the private endpoint service.
	Region string `json:"region"`
	// Unique identifier of the private endpoint you created in your AWS VPC or Azure VNet.
	// +optional
	ID string `json:"id,omitempty"`
	// Private IP address of the private endpoint network interface you created in your Azure VNet.
	// +optional
	IP string `json:"ip,omitempty"`
	// Unique identifier of the Google Cloud project in which you created your endpoints.
	// +optional
	GCPProjectID string `json:"gcpProjectId,omitempty"`
	// Unique identifier of the endpoint group. The endpoint group encompasses all the endpoints that you created in Google Cloud.
	// +optional
	EndpointGroupName string `json:"endpointGroupName,omitempty"`
	// Collection of individual private endpoints that comprise your endpoint group.
	// +optional
	Endpoints GCPEndpoints `json:"endpoints,omitempty"`
}

type GCPEndpoints []GCPEndpoint

type GCPEndpoint struct {
	// Forwarding rule that corresponds to the endpoint you created in Google Cloud.
	EndpointName string `json:"endpointName,omitempty"`
	// Private IP address of the endpoint you created in Google Cloud.
	IPAddress string `json:"ipAddress,omitempty"`
}

// Identifier is required to satisfy "Identifiable" iterface
func (i PrivateEndpoint) Identifier() interface{} {
	return string(i.Provider) + status.TransformRegionToID(i.Region)
}

func (endpoints GCPEndpoints) ConvertToAtlas() *[]admin.CreateGCPForwardingRuleRequest {
	if len(endpoints) == 0 {
		return nil
	}
	result := make([]admin.CreateGCPForwardingRuleRequest, 0, len(endpoints))
	for _, e := range endpoints {
		result = append(result, admin.CreateGCPForwardingRuleRequest{
			EndpointName: pointer.SetOrNil(e.EndpointName, ""),
			IpAddress:    pointer.SetOrNil(e.IPAddress, ""),
		})
	}
	return &result
}
