package status

type Reader interface {
	// GetStatus returns the status of the object.
	GetStatus() interface{}
}

// Common is the struct shared by all statuses in existing Custom Resources.
type Common struct {
	// The phase the current Custom Resource is in. Possible values: 'Reconciling', 'Pending', 'Running', 'Failed'
	Phase Phase `json:"phase"`
}
