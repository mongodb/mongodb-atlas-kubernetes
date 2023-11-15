package v1

import (
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/compat"
)

type NetworkPeer struct {
	//AccepterRegionName is the provider region name of user's vpc.
	// +optional
	AccepterRegionName string `json:"accepterRegionName"`
	// ContainerRegion is the provider region name of Atlas network peer container. If not set, AccepterRegionName is used.
	// +optional
	ContainerRegion string `json:"containerRegion"`
	// AccountID of the user's vpc.
	// +optional
	AWSAccountID string `json:"awsAccountId,omitempty"`
	// ID of the network peer container. If not set, operator will create a new container with ContainerRegion and AtlasCIDRBlock input.
	// +optional
	ContainerID string `json:"containerId"`
	//ProviderName is the name of the provider. If not set, it will be set to "AWS".
	// +optional
	ProviderName provider.ProviderName `json:"providerName,omitempty"`
	//User VPC CIDR.
	// +optional
	RouteTableCIDRBlock string `json:"routeTableCidrBlock,omitempty"`
	//AWS VPC ID.
	// +optional
	VpcID string `json:"vpcId,omitempty"`
	//Atlas CIDR. It needs to be set if ContainerID is not set.
	// +optional
	AtlasCIDRBlock string `json:"atlasCidrBlock"`
	//AzureDirectoryID is the unique identifier for an Azure AD directory.
	// +optional
	AzureDirectoryID string `json:"azureDirectoryId,omitempty"`
	// AzureSubscriptionID is the unique identifier of the Azure subscription in which the VNet resides.
	// +optional
	AzureSubscriptionID string `json:"azureSubscriptionId,omitempty"`
	//ResourceGroupName is the name of your Azure resource group.
	// +optional
	ResourceGroupName string `json:"resourceGroupName,omitempty"`
	// VNetName is name of your Azure VNet. Its applicable only for Azure.
	// +optional
	VNetName string `json:"vnetName,omitempty"`
	// +optional
	// User GCP Project ID. Its applicable only for GCP.
	GCPProjectID string `json:"gcpProjectId,omitempty"`
	// GCP Network Peer Name. Its applicable only for GCP.
	// +optional
	NetworkName string `json:"networkName,omitempty"`
}

func (in *NetworkPeer) ToAtlas() (*mongodbatlas.Peer, error) {
	result := &mongodbatlas.Peer{}
	err := compat.JSONCopy(result, in)
	return result, err
}

func (in *NetworkPeer) ToAtlasPeer() *mongodbatlas.Peer {
	switch in.ProviderName {
	case provider.ProviderAWS:
		return &mongodbatlas.Peer{
			AccepterRegionName:  in.AccepterRegionName,
			AWSAccountID:        in.AWSAccountID,
			ContainerID:         in.ContainerID,
			ProviderName:        string(in.ProviderName),
			RouteTableCIDRBlock: in.RouteTableCIDRBlock,
			VpcID:               in.VpcID,
		}
	case provider.ProviderGCP:
		return &mongodbatlas.Peer{
			ContainerID:  in.ContainerID,
			ProviderName: string(in.ProviderName),
			GCPProjectID: in.GCPProjectID,
			NetworkName:  in.NetworkName,
		}
	case provider.ProviderAzure:
		return &mongodbatlas.Peer{
			ContainerID:         in.ContainerID,
			ProviderName:        string(in.ProviderName),
			AzureDirectoryID:    in.AzureDirectoryID,
			AzureSubscriptionID: in.AzureSubscriptionID,
			ResourceGroupName:   in.ResourceGroupName,
			VNetName:            in.VNetName,
		}
	}

	return nil
}

func (in *NetworkPeer) GetContainerRegion() string {
	if in.ContainerRegion != "" {
		return in.ContainerRegion
	}
	return in.AccepterRegionName
}
