package status

import "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"

// AtlasNetworkPeeringStatus is a status for the AtlasNetworkPeering Custom resource.
// Not the one included in the AtlasProject
type AtlasNetworkPeeringStatus struct {
	api.Common `json:",inline"`
}
