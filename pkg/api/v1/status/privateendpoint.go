package status

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/stringutil"
)

type ProjectPrivateEndpoint struct {
	// Unique identifier for AWS or AZURE Private Link Connection.
	ID string `json:"id,omitempty"`
	// Cloud provider for which you want to retrieve a private endpoint service. Atlas accepts AWS or AZURE.
	Provider provider.ProviderName `json:"provider"`
	// Cloud provider region for which you want to create the private endpoint service.
	Region string `json:"region"`
	// Name of the AWS or Azure Private Link Service that Atlas manages.
	ServiceName string `json:"serviceName,omitempty"`
	// Unique identifier of the Azure Private Link Service (for AWS the same as ID).
	ServiceResourceID string `json:"serviceResourceId,omitempty"`
	// Unique identifier of the AWS or Azure Private Link Interface Endpoint.
	InterfaceEndpointID string `json:"interfaceEndpointId,omitempty"`
	// Unique alphanumeric and special character strings that identify the service attachments associated with the GCP Private Service Connect endpoint service.
	ServiceAttachmentNames []string `json:"serviceAttachmentNames,omitempty"`
	// Collection of individual GCP private endpoints that comprise your network endpoint group.
	Endpoints []GCPEndpoint `json:"endpoints,omitempty"`
}

type GCPEndpoint struct {
	Status       string `json:"status"`
	EndpointName string `json:"endpointName"`
	IPAddress    string `json:"ipAddress"`
}

func (pe ProjectPrivateEndpoint) Identifier() interface{} {
	return string(pe.Provider) + stringutil.SimplifyAndSort(pe.Region)
}
