package status

import (
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
)

type AtlasNetworkPeer struct {
	// Unique identifier for NetworkPeer
	ID string `json:"id"`
	// Cloud provider for which you want to retrieve a network peer
	ProviderName provider.ProviderName `json:"providerName"`
	// Region for which you want to create the network peer
	Region string `json:"region"`

	StatusName string `json:"statusName,omitempty"`

	ErrorState string `json:"errorState,omitempty"`

	ErrorStateName string `json:"errorStateName,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ConnectionID string `json:"connectionId,omitempty"`

	Status string `json:"status,omitempty"`

	VPC string `json:"name,omitempty"`

	GCPProjectID string `json:"gcpProjectId,omitempty"`

	AtlasNetworkName string `json:"atlasNetworkName,omitempty"`

	AtlasGCPProjectID string `json:"atlasGcpProjectId,omitempty"`

	ContainerID string `json:"containerId,omitempty"`
}

func NewNetworkPeerStatus(atlasPeer mongodbatlas.Peer, providerName provider.ProviderName, vpcName string, container mongodbatlas.Container) AtlasNetworkPeer {
	return AtlasNetworkPeer{
		ID:                atlasPeer.ID,
		ProviderName:      providerName,
		Region:            atlasPeer.AccepterRegionName,
		StatusName:        atlasPeer.StatusName,
		ErrorState:        atlasPeer.ErrorState,
		ErrorStateName:    atlasPeer.ErrorStateName,
		ConnectionID:      atlasPeer.ConnectionID,
		Status:            atlasPeer.Status,
		VPC:               vpcName,
		AtlasNetworkName:  container.NetworkName,
		AtlasGCPProjectID: container.GCPProjectID,
		ContainerID:       container.ID,
		GCPProjectID:      atlasPeer.GCPProjectID,
	}
}

func (in *AtlasNetworkPeer) GetStatus() string {
	if in.StatusName == "" {
		return in.Status
	}
	return in.StatusName
}
