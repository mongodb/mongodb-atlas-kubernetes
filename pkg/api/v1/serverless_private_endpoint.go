package v1

import (
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
)

type ServerlessPrivateEndpoint struct {
	Name                     string `json:"name,omitempty"`
	CloudProviderEndpointID  string `json:"cloudProviderEndpointID,omitempty"`
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
