package status

const (
	SearchIndexStatusReady      = "Ready"
	SearchIndexStatusError      = "Error"
	SearchIndexStatusInProgress = "InProgress"
	SearchIndexStatusUnknown    = "Unknown"
)

type IndexStatus string

type DeploymentSearchIndexStatus struct {
	Name    string      `json:"name"`
	ID      string      `json:"ID"`
	Status  IndexStatus `json:"status"`
	Message string      `json:"message"`
}

// +k8s:deepcopy-gen=false
//
//nolint:stylecheck
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
