package v1

type ExternalProjectReference struct {
	// ID is the Atlas project ID
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern:=`^([a-f0-9]{24})$`
	ID string `json:"id"`
}
