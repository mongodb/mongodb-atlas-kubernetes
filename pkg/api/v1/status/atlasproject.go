package status

import "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"

// +k8s:deepcopy-gen=false

// AtlasProjectStatusOption is the option that is applied to Atlas Project Status
type AtlasProjectStatusOption func(s *AtlasProjectStatus)

func AtlasProjectIDOption(id string) AtlasProjectStatusOption {
	return func(s *AtlasProjectStatus) {
		s.ID = id
	}
}

func AtlasProjectExpiredIPAccessOption(lists []project.IPAccessList) AtlasProjectStatusOption {
	return func(s *AtlasProjectStatus) {
		s.ExpiredIPAccessList = lists
	}
}

// AtlasProjectStatus defines the observed state of AtlasProject
type AtlasProjectStatus struct {
	Common `json:",inline"`

	// The ID of the Atlas Project
	// +optional
	ID string `json:"id,omitempty"`

	// The list of IP Access List entries that are expired due to 'deleteAfterDate' being less than the current date.
	// Note, that this field is updated by the Atlas Operator only after specification changes
	ExpiredIPAccessList []project.IPAccessList `json:"expiredIpAccessList,omitempty"`

	// The list of private endpoints configured for current project
	PrivateEndpoints []PrivateEndpoint `json:"privateEndpoints,omitempty"`
}
