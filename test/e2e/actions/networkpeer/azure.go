package networkpeer

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

const AzureResourceGroupName = "atlas-operator-test"
const AzureVPCName = "test-vnet"

/*
For network peering to work with Azure, it's necessary to fulfill the requirements described in the documentation. https://www.mongodb.com/docs/atlas/reference/api/vpc-create-peering-connection/#request-path-parameters
In order not to perform these actions every time, the definition of the role scope has been changed to /subscriptions/<azureSubscriptionId>/resourceGroups/<resourceGroupName> for role definition and role assigment.
Note: For the QA environment, the service principal ID is different. It can be found when creating a network peer for Azure by Atlas UI.
*/

func CreateVPCForAzure(subscriptionID, location, resourceGroup, vnetName string) error {
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return fmt.Errorf("authError: %w", err)
	}
	networkClient := network.NewVirtualNetworksClient(subscriptionID)
	networkClient.Authorizer = authorizer
	_, err = networkClient.CreateOrUpdate(context.Background(), resourceGroup, vnetName, network.VirtualNetwork{
		VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
			AddressSpace: &network.AddressSpace{
				AddressPrefixes: &[]string{
					"10.1.0.0/16", // default address space
				},
			},
		},
		Location: &location,
	})
	if err != nil {
		return fmt.Errorf("can not create Virtual Network: %w", err)
	}
	return nil
}

func DeleteVPCForAzure(subscriptionID, resourceGroup, vnetName string) error {
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
