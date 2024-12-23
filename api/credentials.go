package api

type LocalRef string

// +k8s:deepcopy-gen=false

// CredentialsProvider gives access to custom local credentials
type CredentialsProvider interface {
	Credentials() *LocalObjectReference
}

// +k8s:deepcopy-gen=false

// ResourceWithCredentials is to be implemented by all CRDs using custom local credentials
type ResourceWithCredentials interface {
	CredentialsProvider
	GetName() string
	GetNamespace() string
}

// LocalCredentialHolder is to be embedded by Specs of CRDs using custom local credentials
type LocalCredentialHolder struct {
	// Name of the secret containing Atlas API private and public keys
	ConnectionSecret *LocalObjectReference `json:"connectionSecret,omitempty"`
}

func (ch *LocalCredentialHolder) Credentials() *LocalObjectReference {
	return ch.ConnectionSecret
}
