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

import (
	"regexp"
	"sort"
	"strings"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
)

type ProjectPrivateEndpoint struct {
	// Unique identifier for AWS or AZURE Private Link Connection.
	ID string `json:"id,omitempty"`
	// Cloud provider for which you want to retrieve a private endpoint service. Atlas accepts AWS or AZURE.
	Provider provider.ProviderName `json:"provider"`
	// Cloud provider region for which you want to create the private endpoint service.
	Region string `json:"region"`
	// Name of the AWS or Azure Private Link Service that Atlas manages.
	ServiceName string `json:"serviceName,omitempty"`
	// Unique identifier of the Azure Private Link Service (for AWS the same as ID).
	ServiceResourceID string `json:"serviceResourceId,omitempty"`
	// Unique identifier of the AWS or Azure Private Link Interface Endpoint.
	InterfaceEndpointID string `json:"interfaceEndpointId,omitempty"`
	// Unique alphanumeric and special character strings that identify the service attachments associated with the GCP Private Service Connect endpoint service.
	ServiceAttachmentNames []string `json:"serviceAttachmentNames,omitempty"`
	// Collection of individual GCP private endpoints that comprise your network endpoint group.
	Endpoints []GCPEndpoint `json:"endpoints,omitempty"`
}

type GCPEndpoint struct {
	// State of the MongoDB Atlas endpoint group when MongoDB Cloud received this request.
	Status string `json:"status"`
	// Human-readable label that identifies the Google Cloud consumer forwarding rule that you created.
	EndpointName string `json:"endpointName"`
	// One Private Internet Protocol version 4 (IPv4) address to which this Google Cloud consumer forwarding rule resolves.
	IPAddress string `json:"ipAddress"`
}

// TransformRegionToID makes the same ID from region and regionName fields for PE Connections to match them
// it leaves only characters which are letters or numbers starting from 2
// it also makes a couple swaps and sorts the resulting string
// this function is a temporary work around caused by the empty "region" field in Atlas reply
func TransformRegionToID(region string) string {
	reg := regexp.MustCompile("[^a-z2-9]+")
	temp := strings.ToLower(region)

	// this is GCP specific
	temp = strings.ReplaceAll(temp, "northern", "north")
	temp = strings.ReplaceAll(temp, "southern", "south")
	temp = strings.ReplaceAll(temp, "western", "west")
	temp = strings.ReplaceAll(temp, "eastern", "east")

	temp = reg.ReplaceAllString(temp, "")

	tempSlice := strings.Split(temp, "")
	sort.Strings(tempSlice)
	return strings.Join(tempSlice, "")
}
