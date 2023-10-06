package status

type AtlasFederatedAuthStatus struct {
	Common `json:",inline"`
}

// +k8s:deepcopy-gen=false

type AtlasFederatedAuthStatusOption func(s *AtlasFederatedAuthStatus)
