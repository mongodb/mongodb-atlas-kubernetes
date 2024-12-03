package status

import "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"

// AtlasNetworkPeeringStatus is a status for the AtlasNetworkPeering Custom resource.
// Not the one included in the AtlasProject
type AtlasNetworkPeeringStatus struct {
	api.Common `json:",inline"`

	// ID recrods the identified of thr peer created by Atlas
	ID string `json:"id,omitempty"`

	// ContainerID records the ID of the container created by atlas for this peering
	ContainerID string `json:"containerId,omitempty"`
}
