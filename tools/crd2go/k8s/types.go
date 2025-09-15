package k8s

// LocalReference is reference another Kubernetes resource in the
// same namespace
type LocalReference struct {
	Name string `json:"name"`
}

// Reference is reference a Kubernetes resource, maybe on another namespace
type Reference struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}
