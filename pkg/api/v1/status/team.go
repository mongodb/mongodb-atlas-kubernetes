package status

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

func AtlasTeamUnsetProject(projectID string) AtlasTeamStatusOption {
	return func(s *TeamStatus) {
		index := -1
		for i := 0; i < len(s.Projects); i++ {
			if s.Projects[i].ID == projectID {
				index = i
				break
			}
		}
		if index == -1 {
			return
		}
		s.Projects = append(s.Projects[:index], s.Projects[index+1:]...)
	}
}

type TeamStatus struct {
	Common `json:",inline"`

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
