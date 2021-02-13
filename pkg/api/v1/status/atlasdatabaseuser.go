package status

// +k8s:deepcopy-gen=false

// AtlasDatabaseUserStatusOption is the option that is applied to Atlas Project Status
type AtlasDatabaseUserStatusOption func(s *AtlasDatabaseUserStatus)

// AtlasDatabaseUserStatus defines the observed state of AtlasProject
type AtlasDatabaseUserStatus struct {
	Common `json:",inline"`
}
