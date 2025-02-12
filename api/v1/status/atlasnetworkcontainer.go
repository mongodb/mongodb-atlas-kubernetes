package status

import "github.com/mongodb/mongodb-atlas-kubernetes/v2/api"

// AtlasNetworkContainerStatus is a status for the AtlasNetworkContainer Custom resource.
// Not the one included in the AtlasProject
type AtlasNetworkContainerStatus struct {
	api.Common `json:",inline"`

	// ID record the identifier of the container in Atlas
	ID string `json:"id,omitempty"`

	// Provisioned is true when clusters have been deployed to the container
	Provisioned bool `json:"provisioned,omitempty"`

	// AWSStatus contains AWS only related status information
	AWSStatus *AWSContainerStatus `json:"awsStatus,omitempty"`

	// AzureStatus contains Azure only related status information
	AzureStatus *AzureContainerStatus `json:"azureStatus,omitempty"`

	// GCPStatus contains GCP only related status information
	GCPStatus *GCPContainerStatus `json:"gcpStatus,omitempty"`
}

// AWSContainerStatus contains AWS only related status information
type AWSContainerStatus struct {
	// VpcID is AWS VPC id on the Atlas side
	VpcID string `json:"vpcId,omitempty"`
}

// AzureContainerStatus contains Azure only related status information
type AzureContainerStatus struct {
	// AzureSubscriptionID is Azure Subscription id on the Atlas side
	AzureSubscriptionID string `json:"azureSubscriptionIDpcId,omitempty"`

	// VnetName is Azure network on the Atlas side
	VnetName string `json:"vNetName,omitempty"`
}

// GCPContainerStatus contains GCP only related status information
type GCPContainerStatus struct {
	// GCPProjectID is GCP project on the Atlas side
	GCPProjectID string `json:"gcpProjectID,omitempty"`

	// NetworkName is GCP network on the Atlas side
	NetworkName string `json:"networkName,omitempty"`
}

// +kubebuilder:object:generate=false

type AtlasNetworkContainerStatusOption func(s *AtlasNetworkContainerStatus)
