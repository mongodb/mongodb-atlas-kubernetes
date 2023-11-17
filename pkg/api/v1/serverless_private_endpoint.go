package v1

import (
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/provider"
)

type ServerlessPrivateEndpoint struct {
	// Name is the name of the Serverless PrivateLink Service. Should be unique.
	Name string `json:"name,omitempty"`
	// CloudProviderEndpointID is the identifier of the cloud provider endpoint.
	CloudProviderEndpointID string `json:"cloudProviderEndpointID,omitempty"`
	// PrivateEndpointIPAddress is the IPv4 address of the private endpoint in your Azure VNet that someone added to this private endpoint service.
	PrivateEndpointIPAddress string `json:"privateEndpointIpAddress,omitempty"`
}

// IsInitialState pe initially should be empty except for comment
func (in *ServerlessPrivateEndpoint) IsInitialState() bool {
	return in.Name != "" && in.CloudProviderEndpointID == "" && in.PrivateEndpointIPAddress == ""
}

func (in *ServerlessPrivateEndpoint) ToAtlas(providerName provider.ProviderName) *mongodbatlas.ServerlessPrivateEndpointConnection {
	if in.IsInitialState() {
		return &mongodbatlas.ServerlessPrivateEndpointConnection{
			Comment: in.Name,
		}
	}

	switch providerName {
	case provider.ProviderAWS:
		return &mongodbatlas.ServerlessPrivateEndpointConnection{
			Comment:                 in.Name,
			CloudProviderEndpointID: in.CloudProviderEndpointID,
			ProviderName:            string(providerName),
		}
	case provider.ProviderAzure:
		return &mongodbatlas.ServerlessPrivateEndpointConnection{
			Comment:                  in.Name,
			CloudProviderEndpointID:  in.CloudProviderEndpointID,
			PrivateEndpointIPAddress: in.PrivateEndpointIPAddress,
			ProviderName:             string(providerName),
		}
	default:
		return nil
	}
}
