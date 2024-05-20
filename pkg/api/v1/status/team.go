package status

import "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"

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
