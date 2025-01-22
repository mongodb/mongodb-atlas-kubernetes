package status

import "github.com/mongodb/mongodb-atlas-kubernetes/v2/api"

// AtlasNetworkPeeringStatus is a status for the AtlasNetworkPeering Custom resource.
// Not the one included in the AtlasProject
type AtlasNetworkPeeringStatus struct {
	api.Common `json:",inline"`

	// ID recrods the identified of the peer created by Atlas
	ID string `json:"id,omitempty"`

	// Status describes the last status seen for the network peering setup
	Status string `json:"status,omitempty"`

	// Error refers to the last error seen in the network peering setup
	Error string `json:"error,omitempty"`

	// ContainerID records the ID of the container created by atlas for this peering
	ContainerID string `json:"containerId,omitempty"`

	// ContainerProvisioned is true when the container has been provisioned in Atlas
	ContainerProvisioned bool `json:"containerProvisioned,omitempty"`

	// AWSStatus contains AWS only related status information
	AWSStatus *AWSStatus `json:"awsStatus,omitempty"`

	// AzureStatus contains Azure only related status information
	AzureStatus *AzureStatus `json:"azureStatus,omitempty"`

	// GoogleStatus contains Google only related status information
	GoogleStatus *GoogleStatus `json:"googleStatus,omitempty"`
}

// AWSStatus contains AWS only related status for network peering & container
type AWSStatus struct {
	// ConnectionID is the AWS VPC peering connection ID
	ConnectionID string `json:"connectionId,omitempty"`

	// ContainerVPCId is the AWS Container VPC ID on the Atlas side
	ContainerVpcID string `json:"containerVpcId,omitempty"`
}

// AzureStatus contains Azure only related status about the network container
type AzureStatus struct {
	// AzureSubscriptionID is the Azure subcription id on the Atlas container
	AzureSubscriptionID string `json:"azureSubscriptionId,omitempty"`

	// VnetName is the Azure network name on the Atlas container
	VnetName string `json:"vnetName,omitempty"`
}

// GoogleStatus contains Google only related status about the network container
type GoogleStatus struct {
	// GCPProjectID is the Google Cloud Platform project id on the Atlas container
	GCPProjectID string `json:"gcpProjectId,omitempty"`

	// NetworkName is the Google network name on the Atlas container
	NetworkName string `json:"networkName,omitempty"`
}

// +kubebuilder:object:generate=false

type AtlasNetworkPeeringStatusOption func(s *AtlasNetworkPeeringStatus)
