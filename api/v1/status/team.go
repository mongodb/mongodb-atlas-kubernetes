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
)

// +k8s:deepcopy-gen=false

// AtlasTeamStatusOption is the option that is applied to Atlas Project Status
type AtlasTeamStatusOption func(s *TeamStatus)

func AtlasTeamSetID(ID string) AtlasTeamStatusOption {
	return func(s *TeamStatus) {
		s.ID = ID
	}
}

func AtlasTeamUnsetID() AtlasTeamStatusOption {
	return func(s *TeamStatus) {
		s.ID = ""
	}
}

func AtlasTeamSetProjects(projects []TeamProject) AtlasTeamStatusOption {
	return func(s *TeamStatus) {
		s.Projects = projects
	}
}

// TeamStatus defines the observed state of AtlasTeam.
type TeamStatus struct {
	api.Common `json:",inline"`

	// ID of the team
	ID string `json:"id,omitempty"`
	// List of projects which the team is assigned
	Projects []TeamProject `json:"projects,omitempty"`
}

type TeamProject struct {
	// Unique identifier of the project inside atlas
	ID string `json:"id"`
	// Name given to the project
	Name string `json:"name"`
}
