package v1

type ExternalProjectReference struct {
	// ID is the Atlas project ID
	ID string `json:"id"`
	// Credentials is the name of the secret that holds the Atlas credentials
	Credentials *string `json:"credentials,omitempty"`
}
