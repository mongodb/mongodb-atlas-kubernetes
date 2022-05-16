package v1

import (
	"errors"

	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/compat"
)

type PrivateEndpoint struct {
	// Cloud provider for which you want to retrieve a private endpoint service. Atlas accepts AWS or AZURE.
	// +kubebuilder:validation:Enum=AWS;GCP;AZURE;TENANT
	Provider provider.ProviderName `json:"provider"`
	// Cloud provider region for which you want to create the private endpoint service.
	Region string `json:"region"`
	// Unique identifier of the private endpoint you created in your AWS VPC or Azure Vnet.
	// +optional
	ID string `json:"id,omitempty"`
	// Private IP address of the private endpoint network interface you created in your Azure VNet.
	// +optional
	IP string `json:"ip,omitempty"`
	// Unique identifier of the Google Cloud project in which you created your endpoints.
	// +optional
	GCPProjectID string `json:"gcpProjectId,omitempty"`
	// Unique identifier of the endpoint group. The endpoint group encompasses all of the endpoints that you created in Google Cloud.
	// +optional
	EndpointGroupName string `json:"endpointGroupName,omitempty"`
	// Collection of individual private endpoints that comprise your endpoint group.
	// +optional
	Endpoints GCPEndpoints `json:"endpoints,omitempty"`
}

type GCPEndpoints []GCPEndpoint

type GCPEndpoint struct {
	// Forwarding rule that corresponds to the endpoint you created in Google Cloud.
	EndpointName string `json:"endpointName,omitempty"`
	// Private IP address of the endpoint you created in Google Cloud.
	IPAddress string `json:"ipAddress,omitempty"`
}

// ToAtlas converts the PrivateEndpoint to native Atlas client format.
func (i PrivateEndpoint) ToAtlas() (*mongodbatlas.PrivateEndpoint, error) {
	result := &mongodbatlas.PrivateEndpoint{}
	err := compat.JSONCopy(result, i)
	return result, err
}

// Identifier is required to satisfy "Identifiable" iterface
func (i PrivateEndpoint) Identifier() interface{} {
	return string(i.Provider) + i.Region
}

func (endpoints GCPEndpoints) ConvertToAtlas() ([]*mongodbatlas.GCPEndpoint, error) {
	if len(endpoints) == 0 {
		return nil, errors.New("list of endpoints is empty")
	}
	result := make([]*mongodbatlas.GCPEndpoint, 0)
	for _, e := range endpoints {
		endpoint := &mongodbatlas.GCPEndpoint{}
		err := compat.JSONCopy(endpoint, e)
		if err != nil {
			return nil, err
		}
		result = append(result, endpoint)
	}
	return result, nil
}
