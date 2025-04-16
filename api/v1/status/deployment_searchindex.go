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

const (
	SearchIndexStatusReady      = "READY"
	SearchIndexStatusError      = "Error"
	SearchIndexStatusInProgress = "InProgress"
)

type IndexStatus string

type DeploymentSearchIndexStatus struct {
	Name    string      `json:"name"`
	ID      string      `json:"ID"`
	Status  IndexStatus `json:"status"`
	Message string      `json:"message"`
}

// +k8s:deepcopy-gen=false
type IndexStatusOption func(status *DeploymentSearchIndexStatus)

func NewDeploymentSearchIndexStatus(status IndexStatus, options ...IndexStatusOption) DeploymentSearchIndexStatus {
	result := &DeploymentSearchIndexStatus{
		Status: status,
	}
	for i := range options {
		options[i](result)
	}
	return *result
}

func WithMsg(msg string) IndexStatusOption {
	return func(s *DeploymentSearchIndexStatus) {
		s.Message = msg
	}
}

func WithID(id string) IndexStatusOption {
	return func(s *DeploymentSearchIndexStatus) {
		s.ID = id
	}
}

func WithName(name string) IndexStatusOption {
	return func(s *DeploymentSearchIndexStatus) {
		s.Name = name
	}
}
