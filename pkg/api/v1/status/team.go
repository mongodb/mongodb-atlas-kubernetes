package status

// +k8s:deepcopy-gen=false

// AtlasTeamStatusOption is the option that is applied to Atlas Project Status
type AtlasTeamStatusOption func(s *TeamStatus)

func AtlasTeamID(ID string) AtlasTeamStatusOption {
	return func(s *TeamStatus) {
		s.ID = ID
	}
}

func AtlasTeamAddProject(ID, name string) AtlasTeamStatusOption {
	return func(s *TeamStatus) {
		for _, project := range s.Projects {
			if project.ID == ID {
				return
			}
		}

		s.Projects = append(
			s.Projects,
			TeamProject{
				ID:   ID,
				Name: name,
			},
		)
	}
}

func AtlasTeamRemoveProject(ID, name string) AtlasTeamStatusOption {
	return func(s *TeamStatus) {
		newList := make([]TeamProject, 0, len(s.Projects))

		for _, project := range s.Projects {
			if project.ID == ID {
				continue
			}

			newList = append(
				newList,
				TeamProject{
					ID:   ID,
					Name: name,
				},
			)
		}

		s.Projects = newList
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
