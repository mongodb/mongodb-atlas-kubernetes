package networkpeer

import (
	"os"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/cloud"
)

const (
	SubscriptionID = "AZURE_SUBSCRIPTION_ID"
	DirectoryID    = "AZURE_TENANT_ID"
)

func PreparePeerVPC(peers []v1.NetworkPeer) error {
	for i, peer := range peers {
		awsNetworkPeer, err := NewAWSNetworkPeerService(peer.AccepterRegionName)
		if err != nil {
			return err
		}
		switch peer.ProviderName {
		case provider.ProviderAWS:
			accountID, vpcID, err := awsNetworkPeer.CreateVPCForAWS(peer.RouteTableCIDRBlock)
			if err != nil {
				return err
			}
			peers[i].AWSAccountID = accountID
			peers[i].VpcID = vpcID
		case provider.ProviderGCP:
			err = CreateVPCForGCP(cloud.GoogleProjectID, peer.NetworkName)
			if err != nil {
				return err
			}
		case provider.ProviderAzure:
			err = CreateVPCForAzure(os.Getenv(SubscriptionID), config.AzureRegion, peer.ResourceGroupName, peer.VNetName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func DeletePeerVPC(peers []status.AtlasNetworkPeer) []error {
	var errList []error
	for _, networkPeering := range peers {
		switch networkPeering.ProviderName {
		case provider.ProviderAWS:
			err := DeletePeerConnectionAndVPCForAWS(networkPeering.ConnectionID, networkPeering.Region)
			if err != nil {
				errList = append(errList, err)
			}
		case provider.ProviderGCP:
			err := DeleteVPCForGCP(cloud.GoogleProjectID, networkPeering.VPC)
			if err != nil {
				errList = append(errList, err)
			}
		case provider.ProviderAzure:
			err := DeleteVPCForAzure(os.Getenv(SubscriptionID), AzureResourceGroupName, networkPeering.VPC)
			if err != nil {
				errList = append(errList, err)
			}
		}
	}
	return errList
}

func EstablishPeerConnections(peers []status.AtlasNetworkPeer) error {
	for _, peerStatus := range peers {
		switch peerStatus.ProviderName { // For Azure, it does not need to establish a connection
		case provider.ProviderAWS:
			err := EstablishAWSPeerConnection(peerStatus)
			if err != nil {
				return err
			}
		case provider.ProviderGCP:
			err := EstablishGCPPeerConnectionWithVPC(peerStatus.GCPProjectID, peerStatus.VPC,
				peerStatus.AtlasGCPProjectID, peerStatus.AtlasNetworkName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
