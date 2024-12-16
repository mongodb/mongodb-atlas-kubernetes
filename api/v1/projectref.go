package v1

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
)

// ProjectDualReference encapsulates together the frequnet construct to refer to
// a parent projet by a Kubernetes references, or an external ID, which also
// reqwuires access credentials in such case
type ProjectDualReference struct {
	// "projectRef" is a reference to the parent AtlasProject resource.
	// Mutually exclusive with the "externalProjectRef" field
	// +kubebuilder:validation:Optional
	ProjectRef *common.ResourceRefNamespaced `json:"projectRef,omitempty"`
	// "externalProjectRef" holds the parent Atlas project ID.
	// Mutually exclusive with the "projectRef" field
	// +kubebuilder:validation:Optional
	ExternalProjectRef *ExternalProjectReference `json:"externalProjectRef,omitempty"`

	// Name of the secret containing Atlas API private and public keys
	ConnectionSecret *api.LocalObjectReference `json:"connectionSecret,omitempty"`
}
