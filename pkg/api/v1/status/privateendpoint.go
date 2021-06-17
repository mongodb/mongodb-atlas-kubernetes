package status

import "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"

type ProjectPrivateEndpoint struct {
	// Cloud provider for which you want to retrieve a private endpoint service. Atlas accepts AWS or AZURE.
	Provider provider.ProviderName `json:"provider"`
	// Cloud provider region for which you want to create the private endpoint service.
	Region string `json:"region"`
	// Name of the AWS or Azure Private Link Service that Atlas manages.
	ServiceName string `json:"serviceName,omitempty"`
	// Unique identifier of the AWS or Azure PrivateLink connection.
	ServiceResourceID string `json:"serviceResourceId,omitempty"`
}
