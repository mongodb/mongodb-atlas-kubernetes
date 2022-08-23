package networkpeer

import (
	"context"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/gcp"

	"google.golang.org/api/compute/v1"
)

func CreateVPCForGCP(gcpProjectID string, vnetName string) error {
	computeService, err := compute.NewService(context.Background())
	if err != nil {
		return err
	}
	networkService := compute.NewNetworksService(computeService)
	net := &compute.Network{
		Name:                  vnetName,
		AutoCreateSubnetworks: false,
		// make sure AutoCreateSubnetworks field is included in request
		ForceSendFields: []string{"AutoCreateSubnetworks"},
	}
	_, err = networkService.Insert(gcpProjectID, net).Do()
	return err
}

func EstablishGCPPeerConnectionWithVPC(gpcProjectID, vnetName, atlasGCPProjectID, atlasVnetName string) error {
	computeService, err := compute.NewService(context.Background())
	if err != nil {
		return err
	}
	networkService := compute.NewNetworksService(computeService)
	peer := &compute.NetworkPeering{
		Name:                 "peer",
		Network:              gcp.FormNetworkURL(atlasVnetName, atlasGCPProjectID),
		ExchangeSubnetRoutes: true,
	}
	request := &compute.NetworksAddPeeringRequest{
		NetworkPeering: peer,
	}
	_, err = networkService.AddPeering(gpcProjectID, vnetName, request).Do()
	return err
}

func DeleteVPCForGCP(gcpProjectID, vnetName string) error {
	computeService, err := compute.NewService(context.Background())
	if err != nil {
		return err
	}
	networkService := compute.NewNetworksService(computeService)
	_, err = networkService.Delete(gcpProjectID, vnetName).Do()
	return err
}
