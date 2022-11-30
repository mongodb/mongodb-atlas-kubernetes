package status

import (
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"
)

type ServerlessPrivateEndpoint struct {
	// ID is the identifier of the Serverless PrivateLink Service.
	ID string `json:"_id,omitempty"`
	// CloudProviderEndpointID is the identifier of the cloud provider endpoint.
	CloudProviderEndpointID string `json:"cloudProviderEndpointId,omitempty"`
	// Name is the name of the Serverless PrivateLink Service. Should be unique.
	Name string `json:"name,omitempty"`
	// EndpointServiceName is the name of the PrivateLink endpoint service in AWS. Returns null while the endpoint service is being created.
	EndpointServiceName string `json:"endpointServiceName,omitempty"`
	// ErrorMessage is the error message if the Serverless PrivateLink Service failed to create or connect.
	ErrorMessage string `json:"errorMessage,omitempty"`
	// Status of the AWS Serverless PrivateLink connection.
	Status string `json:"status,omitempty"`
	// ProviderName is human-readable label that identifies the cloud provider. Values include AWS or AZURE.
	ProviderName string `json:"providerName,omitempty"`
	// PrivateEndpointIPAddress is the IPv4 address of the private endpoint in your Azure VNet that someone added to this private endpoint service.
	PrivateEndpointIPAddress string `json:"privateEndpointIpAddress,omitempty"`
	// PrivateLinkServiceResourceID is the root-relative path that identifies the Azure Private Link Service that MongoDB Cloud manages. MongoDB Cloud returns null while it creates the endpoint service.
	PrivateLinkServiceResourceID string `json:"privateLinkServiceResourceId,omitempty"`
}

func SPEFromAtlas(in mongodbatlas.ServerlessPrivateEndpointConnection) ServerlessPrivateEndpoint {
	return ServerlessPrivateEndpoint{
		ID:                           in.ID,
		CloudProviderEndpointID:      in.CloudProviderEndpointID,
		Name:                         in.Comment,
		EndpointServiceName:          in.EndpointServiceName,
		ErrorMessage:                 in.ErrorMessage,
		Status:                       in.Status,
		ProviderName:                 in.ProviderName,
		PrivateEndpointIPAddress:     in.PrivateEndpointIPAddress,
		PrivateLinkServiceResourceID: in.PrivateLinkServiceResourceID,
	}
}

func FailedToCreateSPE(comment, message string) ServerlessPrivateEndpoint {
	return ServerlessPrivateEndpoint{
		ErrorMessage: message,
		Name:         comment,
		Status:       StatusFailed,
	}
}

func FailedDuplicationSPE(name, cloudProviderEndpointID, privateEndpointIPAddress string) ServerlessPrivateEndpoint {
	return ServerlessPrivateEndpoint{
		CloudProviderEndpointID:  cloudProviderEndpointID,
		PrivateEndpointIPAddress: privateEndpointIPAddress,
		ErrorMessage:             fmt.Sprintf("The SPE with same name exists: %s. Please use unique names", name),
		Name:                     name,
		Status:                   StatusFailed,
	}
}

func FailedToConnectSPE(pe mongodbatlas.ServerlessPrivateEndpointConnection, message string) ServerlessPrivateEndpoint {
	return ServerlessPrivateEndpoint{
		ID:                           pe.ID,
		CloudProviderEndpointID:      pe.CloudProviderEndpointID,
		Name:                         pe.Comment,
		EndpointServiceName:          pe.EndpointServiceName,
		ErrorMessage:                 message,
		Status:                       StatusFailed,
		ProviderName:                 pe.ProviderName,
		PrivateEndpointIPAddress:     pe.PrivateEndpointIPAddress,
		PrivateLinkServiceResourceID: pe.PrivateLinkServiceResourceID,
	}
}
