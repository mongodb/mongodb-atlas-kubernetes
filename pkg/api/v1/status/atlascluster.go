package status

// AtlasClusterStatus defines the observed state of AtlasCluster.
type AtlasClusterStatus struct {
	Common         `json:",inline"`
	StateName      string `json:"stateName,omitempty"`
	MongoDBVersion string `json:"mongoDBVersion,omitempty"`
}

// +k8s:deepcopy-gen=false

// AtlasClusterStatusOption is the option that is applied to Atlas Project Status
type AtlasClusterStatusOption func(s *AtlasClusterStatus)

func AtlasClusterStateNameOption(stateName string) AtlasClusterStatusOption {
	return func(s *AtlasClusterStatus) {
		s.StateName = stateName
	}
}

func AtlasClusterMongoDBVersionOption(mongoDBVersion string) AtlasClusterStatusOption {
	return func(s *AtlasClusterStatus) {
		s.MongoDBVersion = mongoDBVersion
	}
}
