package status

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
