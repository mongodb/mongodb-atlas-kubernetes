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

import "github.com/mongodb/mongodb-atlas-kubernetes/v2/api"

// AtlasNetworkContainerStatus is a status for the AtlasNetworkContainer Custom resource.
// Not the one included in the AtlasProject
type AtlasNetworkContainerStatus struct {
	api.Common `json:",inline"`

	// ID record the identifier of the container in Atlas
	ID string `json:"id,omitempty"`

	// Provisioned is true when clusters have been deployed to the container before
	// the last reconciliation
	Provisioned bool `json:"provisioned,omitempty"`
}

// +kubebuilder:object:generate=false

type AtlasNetworkContainerStatusOption func(s *AtlasNetworkContainerStatus)
