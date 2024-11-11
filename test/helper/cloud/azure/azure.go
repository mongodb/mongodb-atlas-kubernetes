package azure

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
)

const (
	defaultTestResourceGroupName = "svet-test"
)

type azureConnection struct {
	resourceGroupName      string
	credentials            *azidentity.DefaultAzureCredential
	networkResourceFactory *armnetwork.ClientFactory
}

func newClient(resourceGroupName string) (*azureConnection, error) {
	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	networkFactory, err := armnetwork.NewClientFactory(subscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	return &azureConnection{
		resourceGroupName:      resourceGroupName,
		networkResourceFactory: networkFactory,
		credentials:            cred,
	}, err
}

func RegionCode(region string) string {
	region2azure := map[string]string{
		"US_CENTRAL": "us_central",
		"US_EAST":    "eastus",
		"US_EAST_2":  "eastus2",
	}
	azureRegion, ok := region2azure[region]
	if !ok {
		return fmt.Sprintf("unsupported region %q", region)
	}
	return azureRegion
}

func TestResourceGroupName() string {
	return defaultTestResourceGroupName
}
