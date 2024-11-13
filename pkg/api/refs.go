package api

// LocalObjectReference is a reference to an object in the same namespace as the referent
type LocalObjectReference struct {
	// Name of the resource being referred to
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	Name string `json:"name"`
}

type ExternalProjectReference struct {
	// ID is the Atlas project ID
	// +kubebuilder:validation:Required
	ID string `json:"id"`
}

// ResourceRefNamespaced is a reference to a Kubernetes Resource that allows to configure the namespace
type ResourceRefNamespaced struct {
	// Name is the name of the Kubernetes Resource
	Name string `json:"name"`

	// Namespace is the namespace of the Kubernetes Resource
	// +optional
	Namespace string `json:"namespace"`
}

type ProjectReferences struct {
	ConnectionSecret *LocalObjectReference `json:"connectionSecret,omitempty"`

	// Project is a reference to AtlasProject resource the user belongs to
	Project *ResourceRefNamespaced `json:"projectRef,omitempty"`
	// ExternalProjectRef holds the Atlas project ID the user belongs to
	ExternalProjectRef *ExternalProjectReference `json:"externalProjectRef,omitempty"`
}
