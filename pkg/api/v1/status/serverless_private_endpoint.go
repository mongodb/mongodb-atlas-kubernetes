package status

import (
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"
)

type ServerlessPrivateEndpoint struct {
	ID                           string `json:"_id,omitempty"` // Unique identifier of the Serverless PrivateLink Service.
	CloudProviderEndpointID      string `json:"cloudProviderEndpointId,omitempty"`
	Name                         string `json:"name,omitempty"`
	EndpointServiceName          string `json:"endpointServiceName,omitempty"` // Name of the PrivateLink endpoint service in AWS. Returns null while the endpoint service is being created.
	ErrorMessage                 string `json:"errorMessage,omitempty"`
	Status                       string `json:"status,omitempty"`                       // Status of the AWS Serverless PrivateLink connection: INITIATING, WAITING_FOR_USER, FAILED, DELETING, AVAILABLE.
	ProviderName                 string `json:"providerName,omitempty"`                 // Human-readable label that identifies the cloud provider. Values include AWS or AZURE.
	PrivateEndpointIPAddress     string `json:"privateEndpointIpAddress,omitempty"`     // IPv4 address of the private endpoint in your Azure VNet that someone added to this private endpoint service.
	PrivateLinkServiceResourceID string `json:"privateLinkServiceResourceId,omitempty"` // Root-relative path that identifies the Azure Private Link Service that MongoDB Cloud manages. MongoDB Cloud returns null while it creates the endpoint service.
}

func SPEFromAtlas(in *mongodbatlas.ServerlessPrivateEndpointConnection) ServerlessPrivateEndpoint {
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
		Status:       "FAILED", // TODO: use constant
	}
}

func FailedDuplicationSPE(name, cloudProviderEndpointID, privateEndpointIpAddress string) ServerlessPrivateEndpoint {
	return ServerlessPrivateEndpoint{
		CloudProviderEndpointID:  cloudProviderEndpointID,
		PrivateEndpointIPAddress: privateEndpointIpAddress,
		ErrorMessage:             fmt.Sprintf("The SPE with same name exists: %s. Please use unique names", name),
		Name:                     name,
		Status:                   "FAILED", // TODO: use constant
	}
}

func FailedToConnectSPE(pe mongodbatlas.ServerlessPrivateEndpointConnection, message string) ServerlessPrivateEndpoint {
	return ServerlessPrivateEndpoint{
		ID:                           pe.ID,
		CloudProviderEndpointID:      pe.CloudProviderEndpointID,
		Name:                         pe.Comment,
		EndpointServiceName:          pe.EndpointServiceName,
		ErrorMessage:                 message,
		Status:                       "FAILED", // TODO: use constant
		ProviderName:                 pe.ProviderName,
		PrivateEndpointIPAddress:     pe.PrivateEndpointIPAddress,
		PrivateLinkServiceResourceID: pe.PrivateLinkServiceResourceID,
	}
}
