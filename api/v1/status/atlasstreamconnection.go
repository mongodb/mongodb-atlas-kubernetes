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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
)

type AtlasStreamConnectionStatus struct {
	api.Common `json:",inline"`
	// List of instances using the connection configuration
	Instances []common.ResourceRefNamespaced `json:"instances,omitempty"`
}

// +kubebuilder:object:generate=false

type AtlasStreamConnectionStatusOption func(s *AtlasStreamConnectionStatus)
