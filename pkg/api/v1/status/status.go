package status

// +k8s:deepcopy-gen=false

type Reader interface {
	// GetStatus returns the status of the object.
	GetStatus() interface{}
}

// +k8s:deepcopy-gen=false

type Writer interface {
	// UpdateStatus allows to do the update of the status of an Atlas Custom resource.
	UpdateStatus(conditions []Condition, option ...Option)
}

// Common is the struct shared by all statuses in existing Custom Resources.
type Common struct {
	// Conditions is the list of statuses showing the current state of the Atlas Custom Resource
	Conditions []Condition `json:"conditions"`
}
