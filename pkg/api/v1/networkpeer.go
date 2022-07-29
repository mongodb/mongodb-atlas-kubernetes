package v1

import (
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/compat"
)

type NetworkPeer struct {
	// +optional
	ID string `json:"id"`
	//AccepterRegionName is Atlas region where the container resides.
	// +optional
	AccepterRegionName string `json:"accepterRegionName"`
	// +optional
	ContainerRegion string `json:"containerRegion"`
	// +optional
	AWSAccountID string `json:"awsAccountId,omitempty"`
	// +optional
	ContainerID string `json:"containerId"`
	//ProviderName is the name of the provider. If not set, it will be set to "aws"
	// +optional
	ProviderName provider.ProviderName `json:"providerName,omitempty"`
	// +optional
	RouteTableCIDRBlock string `json:"routeTableCidrBlock,omitempty"`
	// +optional
	VpcID string `json:"vpcId,omitempty"`
	// +optional
	AtlasCIDRBlock string `json:"atlasCidrBlock"`
	// +optional
	AzureDirectoryID string `json:"azureDirectoryId,omitempty"`
	// +optional
	AzureSubscriptionID string `json:"azureSubscriptionId,omitempty"`
	// +optional
	ResourceGroupName string `json:"resourceGroupName,omitempty"`
	// +optional
	VNetName string `json:"vnetName,omitempty"`
	// +optional
	GCPProjectID string `json:"gcpProjectId,omitempty"`
	// +optional
	NetworkName string `json:"networkName,omitempty"`
}

func (in *NetworkPeer) ToAtlas() (*mongodbatlas.Peer, error) {
	result := &mongodbatlas.Peer{}
	err := compat.JSONCopy(result, in)
	return result, err
}

func (in *NetworkPeer) GetContainerRegion() string {
	if in.ContainerRegion != "" {
		return in.ContainerRegion
	}
	return in.AccepterRegionName
}
