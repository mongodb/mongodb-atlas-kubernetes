// Copyright 2024 MongoDB Inc
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

package status

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
)

// AtlasPrivateEndpointStatus is the most recent observed status of the AtlasPrivateEndpoint cluster. Read-only.
type AtlasPrivateEndpointStatus struct {
	api.Common `json:",inline"`
	// ServiceID is the unique identifier of the private endpoint service in Atlas
	ServiceID string `json:"serviceId,omitempty"`
	// ServiceStatus is the state of the private endpoint service
	ServiceStatus string `json:"serviceStatus,omitempty"`
	// Error is the description of the failure occurred when configuring the private endpoint
	Error string `json:"error,omitempty"`
	// ServiceName is the unique identifier of the Amazon Web Services (AWS) PrivateLink endpoint service or Azure Private Link Service managed by Atlas
	ServiceName string `json:"serviceName,omitempty"`
	// ResourceID is the root-relative path that identifies of the Atlas Azure Private Link Service
	ResourceID string `json:"resourceId,omitempty"`
	// ServiceAttachmentNames is the list of URLs that identifies endpoints that Atlas can use to access one service across the private connection
	ServiceAttachmentNames []string `json:"serviceAttachmentNames,omitempty"`
	// Endpoints are the status of the endpoints connected to the service
	Endpoints []EndpointInterfaceStatus `json:"endpoints,omitempty"`
}

// EndpointInterfaceStatus is the most recent observed status the interfaces attached to the configured service. Read-only.
type EndpointInterfaceStatus struct {
	// ID is the external identifier set on the specification to configure the interface
	ID string `json:"ID,omitempty"`
	// ConnectionName is the label that Atlas generates that identifies the Azure private endpoint connection
	ConnectionName string `json:"connectionName,omitempty"`
	// GCPForwardingRules is the status of the customer GCP private endpoint(forwarding rules)
	GCPForwardingRules []GCPForwardingRule `json:"gcpForwardingRules,omitempty"`
	// InterfaceStatus is the state of the private endpoint interface
	Status string `json:"InterfaceStatus,omitempty"`
	// Error is the description of the failure occurred when configuring the private endpoint
	Error string `json:"error,omitempty"`
}

// GCPForwardingRule is the most recent observed status the GCP forwarding rules configured for an interface. Read-only.
type GCPForwardingRule struct {
	// Human-readable label that identifies the Google Cloud consumer forwarding rule that you created.
	Name string `json:"name,omitempty"`
	// State of the MongoDB Atlas endpoint group.
	Status string `json:"status,omitempty"`
}

// +kubebuilder:object:generate=false

type AtlasPrivateEndpointStatusOption func(s *AtlasPrivateEndpointStatus)
