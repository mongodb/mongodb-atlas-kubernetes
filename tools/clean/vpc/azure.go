package vpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

func deleteAzureVPCBySubstr(ctx context.Context, subID, resourceGroupName, substr string) (bool, error) {
	ok := true
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return false, fmt.Errorf("error creating authorizer: %v", err)
	}
	vnetClient := network.NewVirtualNetworksClient(subID)
	vnetClient.Authorizer = authorizer

	vnets, err := vnetClient.List(ctx, resourceGroupName)
	if err != nil {
		return false, fmt.Errorf("error fetching vnets: %v", err)
	}
	var allErr error
	for _, vnet := range vnets.Values() {
		if vnet.Name != nil && strings.HasPrefix(*vnet.Name, substr) {
			log.Printf("deleting vnet %s", *vnet.Name)
			_, err = vnetClient.Delete(ctx, resourceGroupName, *vnet.Name)
			if err != nil {
				allErr = errors.Join(allErr, fmt.Errorf("error deleting vnet: %v", err))
				ok = false
			}
		}
	}

	return ok, allErr
}
