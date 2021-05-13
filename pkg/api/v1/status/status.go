package status

// +k8s:deepcopy-gen=false

type Reader interface {
	// GetStatus returns the status of the object.
	GetStatus() Status
}

// +k8s:deepcopy-gen=false

type Writer interface {
	// UpdateStatus allows to do the update of the status of an Atlas Custom resource.
	UpdateStatus(conditions []Condition, option ...Option)
}

// +k8s:deepcopy-gen=false

// Status is a generic status for any Custom Resource managed by Atlas Operator
type Status interface {
	GetConditions() []Condition
	GetObservedGeneration() int64
}

var _ Status = &Common{}

// Common is the struct shared by all statuses in existing Custom Resources.
type Common struct {
	// Conditions is the list of statuses showing the current state of the Atlas Custom Resource
	Conditions []Condition `json:"conditions"`

	// ObservedGeneration indicates the generation of the resource specification that the Atlas Operator is aware of.
	// The Atlas Operator updates this field to the 'metadata.generation' as soon as it starts reconciliation of the resource.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

func (c Common) GetConditions() []Condition {
	return c.Conditions
}

func (c Common) GetObservedGeneration() int64 {
	return c.ObservedGeneration
}
