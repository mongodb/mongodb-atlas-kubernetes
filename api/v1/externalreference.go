package v1

type ExternalProjectReference struct {
	// ID is the Atlas project ID
	// +kubebuilder:validation:Required
	ID string `json:"id"`
}
