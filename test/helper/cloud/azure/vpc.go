package azure

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func CreateVPC(ctx context.Context, vpcName, cidr, region string) (string, error) {
	azr, err := newClient(TestResourceGroupName())
	if err != nil {
		return "", fmt.Errorf("failed to create azure client: %w", err)
	}
	vpcClient := azr.networkResourceFactory.NewVirtualNetworksClient()

	op, err := vpcClient.BeginCreateOrUpdate(
		ctx,
		azr.resourceGroupName,
		vpcName,
		armnetwork.VirtualNetwork{
			Location: pointer.MakePtr(region),
			Properties: &armnetwork.VirtualNetworkPropertiesFormat{
				AddressSpace: &armnetwork.AddressSpace{
					AddressPrefixes: []*string{
						pointer.MakePtr(cidr),
					},
				},
			},
			Tags: map[string]*string{
				"Name": pointer.MakePtr(vpcName),
			},
		},
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to begin create azure VPC: %w", err)
	}

	vpc, err := op.PollUntilDone(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("creation process of VPC failed: %w", err)
	}
	if vpc.Name == nil {
		return "", errors.New("VPC created without a name")
	}
	return *vpc.Name, nil
}

func DeleteVPC(ctx context.Context, vpcName string) error {
	azr, err := newClient(TestResourceGroupName())
	if err != nil {
		return fmt.Errorf("failed to create azure client: %w", err)
	}
	vpcClient := azr.networkResourceFactory.NewVirtualNetworksClient()

	op, err := vpcClient.BeginDelete(
		ctx,
		azr.resourceGroupName,
		vpcName,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = op.PollUntilDone(ctx, nil)

	return err
}
