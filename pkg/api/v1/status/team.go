package status

// +k8s:deepcopy-gen=false

// AtlasTeamStatusOption is the option that is applied to Atlas Project Status
type AtlasTeamStatusOption func(s *TeamStatus)

func AtlasTeamID(ID string) AtlasTeamStatusOption {
	return func(s *TeamStatus) {
		s.ID = ID
	}
}

func AtlasAddProject(ID, name string) AtlasTeamStatusOption {
	return func(s *TeamStatus) {
		s.Projects = append(
			s.Projects,
			TeamProject{
				ID:   ID,
				Name: name,
			},
		)
	}
}

type TeamStatus struct {
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
