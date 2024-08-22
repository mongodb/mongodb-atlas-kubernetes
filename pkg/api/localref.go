package api

// LocalObjectReference is a reference to an object in the same namespace as the referent
type LocalObjectReference struct {
	// Name of the resource being referred to
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	Name string `json:"name"`
}
