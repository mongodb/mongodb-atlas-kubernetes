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

	Name string `json:"name,omitempty"`
}

func FromAtlas(peer mongodbatlas.Peer, providerName provider.ProviderName, vpcName string) AtlasNetworkPeer {
	return AtlasNetworkPeer{
		ID:             peer.ID,
		ProviderName:   providerName,
		Region:         peer.AccepterRegionName,
		StatusName:     peer.StatusName,
		ErrorState:     peer.ErrorState,
		ErrorStateName: peer.ErrorStateName,
		ConnectionID:   peer.ConnectionID,
		Status:         peer.Status,
		Name:           vpcName, //TODO: is it unique?
	}
}

func (in *AtlasNetworkPeer) GetStatus() string {
	if in.StatusName == "" {
		return in.Status
	}
	return in.StatusName
}
