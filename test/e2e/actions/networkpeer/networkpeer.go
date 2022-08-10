package networkpeer

import (
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/cloud"
)

func PreparePeerVPC(peers []v1.NetworkPeer, namespace string) error {
	for i, peer := range peers { //TODO: refactor it
		awsNetworkPeer, err := NewAWSNetworkPeerService(peer.AccepterRegionName)
		if err != nil {
			return err
		}
		testID := fmt.Sprintf("%s-%d", namespace, i)
		switch peer.ProviderName {
		case provider.ProviderAWS:
			accountID, vpcID, err := awsNetworkPeer.CreateVPC(peer.RouteTableCIDRBlock, testID)
			if err != nil {
				return err
			}
			peers[i].AWSAccountID = accountID
			peers[i].VpcID = vpcID
		case provider.ProviderGCP:
			err = CreateVPC(cloud.GoogleProjectID, peer.NetworkName)
			return err
		}
	}
	return nil
}

func EstablishPeerConnections(peers []status.AtlasNetworkPeer) error {
	for _, peerStatus := range peers {
		switch peerStatus.ProviderName {
		case provider.ProviderAWS:
			err := EstablishAWSPeerConnection(peerStatus)
			if err != nil {
				return err
			}
		case provider.ProviderGCP:
			err := EstablishPeerConnectionWithVPC(peerStatus.GCPProjectID, peerStatus.VPC,
				peerStatus.AtlasGCPProjectID, peerStatus.AtlasNetworkName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
