package status

type Phase string

const (
	// PhaseReconciling means the controller is in the middle of reconciliation process
	PhaseReconciling Phase = "Reconciling"

	// PhasePending means the reconciliation has finished but the resource is neither in Error nor Running state -
	// most of all waiting for some event to happen (e.g. Cluster get provisioned etc)
	PhasePending Phase = "Pending"

	// PhaseRunning means the Custom Resource is in a running state and ready to be used
	PhaseRunning Phase = "Running"

	// PhaseFailed means the Custom Resource is in a failed state
	PhaseFailed Phase = "Failed"
)
