package provider

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
)

type Azure struct {
	resourceGroupName string

	vpcClient             *armnetwork.VirtualNetworksClient
	privateEndpointClient *armnetwork.PrivateEndpointsClient
	vaultClient           *armkeyvault.KeysClient
}

func (a *Azure) DeleteVpc(ctx context.Context, vpcName string) error {
	op, err := a.vpcClient.BeginDelete(
		ctx,
		a.resourceGroupName,
		vpcName,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = op.PollUntilDone(ctx, nil)

	return err
}
func (a *Azure) DeletePrivateEndpoint(ctx context.Context, endpointNames []string) error {
	for _, endpointName := range endpointNames {
		_, err := a.privateEndpointClient.Get(ctx, a.resourceGroupName, endpointName, nil)
		if err != nil {
			var respErr *azcore.ResponseError
			if ok := errors.As(err, &respErr); ok && respErr.StatusCode == 404 {
				continue
			}

			return err
		}

		op, err := a.privateEndpointClient.BeginDelete(
			ctx,
			a.resourceGroupName,
			endpointName,
			nil,
		)
		if err != nil {
			return err
		}

		_, err = op.PollUntilDone(ctx, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewAzureCleaner() (*Azure, error) {
	subscriptionID, defined := os.LookupEnv("AZURE_SUBSCRIPTION_ID")
	if !defined {
		return nil, fmt.Errorf("AZURE_SUBSCRIPTION_ID must be set")
	}

	resourceGroupName, defined := os.LookupEnv("AZURE_RESOURCE_GROUP_NAME")
	if !defined {
		return nil, fmt.Errorf("AZURE_RESOURCE_GROUP_NAME must be set")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	networkFactory, err := armnetwork.NewClientFactory(subscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	vaultFactory, err := armkeyvault.NewClientFactory(subscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	return &Azure{
		resourceGroupName:     resourceGroupName,
		vpcClient:             networkFactory.NewVirtualNetworksClient(),
		privateEndpointClient: networkFactory.NewPrivateEndpointsClient(),
		vaultClient:           vaultFactory.NewKeysClient(),
	}, err
}
