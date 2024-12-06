package api

type LocalRef string

// +k8s:deepcopy-gen=false

// CredentialsProvider gives access to custom local credentials
//
// Deprecated: CredentialsProvider is not needed when using ProjectReferences
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
//
// Deprecated: LocalCredentialHolder has been replaced by ProjectReferences
// which embeds together the credentials and both the external k8s references
// to a project. That way, common logic can be used to resolve such references
// and the consumer types avoid the need to implement the CredentialsProvider
// interface
type LocalCredentialHolder struct {
	// Name of the secret containing Atlas API private and public keys
	ConnectionSecret *LocalObjectReference `json:"connectionSecret,omitempty"`
}

func (ch *LocalCredentialHolder) Credentials() *LocalObjectReference {
	return ch.ConnectionSecret
}
