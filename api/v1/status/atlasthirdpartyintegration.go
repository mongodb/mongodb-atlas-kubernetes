// Copyright 2025 MongoDB.
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

// +k8s:deepcopy-gen=true

// AtlasThirdPartyIntegrationStatus holds the status of an integration
type AtlasThirdPartyIntegrationStatus struct {
	UnifiedStatus `json:",inline"`

	// ID of the third party integration resource in Atlas
	ID string `json:"id,omitempty"`
}

// +k8s:deepcopy-gen=false

type IntegrationStatusOption func(status *AtlasThirdPartyIntegrationStatus)

func NewAtlasThirdPartyIntegrationStatus(options ...IntegrationStatusOption) AtlasThirdPartyIntegrationStatus {
	result := &AtlasThirdPartyIntegrationStatus{}
	for i := range options {
		options[i](result)
	}
	return *result
}

func WithIntegrationID(id string) IntegrationStatusOption {
	return func(i *AtlasThirdPartyIntegrationStatus) {
		i.ID = id
	}
}
