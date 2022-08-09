package status

import (
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
)

type AtlasNetworkPeer struct {
	// Unique identifier for NetworkPeer.
	ID string `json:"id"`
	// Cloud provider for which you want to retrieve a network peer.
	ProviderName provider.ProviderName `json:"providerName"`
	// Region for which you want to create the network peer. It isn't needed for GCP
	Region string `json:"region"`
	// Status of the network peer. Applicable only for AWS.
	StatusName string `json:"statusName,omitempty"`
	// Error state of the network peer. Applicable only for Azure.
	ErrorState string `json:"errorState,omitempty"`
	// Error state of the network peer. Applicable only for AWS.
	ErrorStateName string `json:"errorStateName,omitempty"`
	// Error state of the network peer. Applicable only for GCP.
	ErrorMessage string `json:"errorMessage,omitempty"`
	// Unique identifier of the network peer connection. Applicable only for AWS.
	ConnectionID string `json:"connectionId,omitempty"`
	// Status of the network peer. Applicable only for GCP and Azure.
	Status string `json:"status,omitempty"`
	// VPC is general purpose field for storing the name of the VPC.
	// VPC is vpcID for AWS, user networkName for GCP, and vnetName for Azure.
	VPC string `json:"vpc,omitempty"`
	// ProjectID of the user's vpc. Applicable only for GCP.
	GCPProjectID string `json:"gcpProjectId,omitempty"`
	// Atlas Network Name. Applicable only for GCP. It's needed to add network peer connection.
	AtlasNetworkName string `json:"atlasNetworkName,omitempty"`
	// ProjectID of Atlas container. Applicable only for GCP. It's needed to add network peer connection.
	AtlasGCPProjectID string `json:"atlasGcpProjectId,omitempty"`
	// ContainerID of Atlas network peer container.
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
