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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
)

// ProjectDualReference encapsulates the common constructs to refer to a parent project;
// by either a Kubernetes reference or an external ID, which also requires access credentials.
type ProjectDualReference struct {
	// projectRef is a reference to the parent AtlasProject resource.
	// Mutually exclusive with the "externalProjectRef" field.
	// +kubebuilder:validation:Optional
	ProjectRef *common.ResourceRefNamespaced `json:"projectRef,omitempty"`
	// externalProjectRef holds the parent Atlas project ID.
	// Mutually exclusive with the "projectRef" field.
	// +kubebuilder:validation:Optional
	ExternalProjectRef *ExternalProjectReference `json:"externalProjectRef,omitempty"`
	// Name of the secret containing Atlas API private and public keys.
	ConnectionSecret *api.LocalObjectReference `json:"connectionSecret,omitempty"`
}
