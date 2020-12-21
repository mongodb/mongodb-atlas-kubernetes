package status

// +k8s:deepcopy-gen=false

// Option is the function that is applied to the status field of an Atlas Custom Resource.
// This is the way to handle some random data that need to be written to status.
type Option interface{}
