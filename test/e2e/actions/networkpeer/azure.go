package networkpeer

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

func CreateAzureVPC(subscriptionID string, vnetName, resourceGroup string) error {
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return fmt.Errorf("authError: %w", err)
	}
	networkClient := network.NewVirtualNetworksClient(subscriptionID)
	networkClient.Authorizer = authorizer
	_, err = networkClient.CreateOrUpdate(context.Background(), resourceGroup, vnetName, network.VirtualNetwork{
		VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{},
	})
	if err != nil {
		return fmt.Errorf("can not create Virtual Network: %w", err)
	}
	return nil
}

func DeleteAzureVPC(subscriptionID string, vnetName, resourceGroup string) error {
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return fmt.Errorf("authError: %w", err)
	}
	networkClient := network.NewVirtualNetworksClient(subscriptionID)
	networkClient.Authorizer = authorizer
	_, err = networkClient.Delete(context.Background(), resourceGroup, vnetName)
	if err != nil {
		return fmt.Errorf("can not delete Virtual Network: %w", err)
	}
	return nil
}
