package v1

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
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

// +k8s:deepcopy-gen=false

// ProjectReferrer is anything that holds a ProjectDualReference
type ProjectReferrer interface {
	ProjectDualRef() *ProjectDualReference
}

// +k8s:deepcopy-gen=false

// ProjectReferrerObject is an project referrer that is also an Kubernetes Object
type ProjectReferrerObject interface {
	client.Object
	ProjectReferrer
}
