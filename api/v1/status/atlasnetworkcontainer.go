package status

import "github.com/mongodb/mongodb-atlas-kubernetes/v2/api"

// AtlasNetworkContainerStatus is a status for the AtlasNetworkContainer Custom resource.
// Not the one included in the AtlasProject
type AtlasNetworkContainerStatus struct {
	api.Common `json:",inline"`

	// ID record the identifier of the container in Atlas
	ID string `json:"id,omitempty"`

	// Provisioned is true when clusters have been deployed to the container before
	// the last reconciliation
	Provisioned bool `json:"provisioned,omitempty"`
}

// +kubebuilder:object:generate=false

type AtlasNetworkContainerStatusOption func(s *AtlasNetworkContainerStatus)
