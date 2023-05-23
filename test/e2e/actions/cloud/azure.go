package cloud

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/onsi/ginkgo/v2/dsl/core"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
)

type AzureAction struct {
	t                 core.GinkgoTInterface
	resourceGroupName string
	network           *azureNetwork

	resourceFactory *armnetwork.ClientFactory
}

type azureNetwork struct {
	VPC     *armnetwork.VirtualNetwork
	Subnets map[string]*armnetwork.Subnet
}

const (
	// TODO get from Azure
	ResourceGroup   = "svet-test"
	AzureVPC        = "svet-test-vpc"
	AzureSubnetName = "default"
)

// InitNetwork Check if minimum network resources exist and when not, create them
func (a *AzureAction) InitNetwork(vpcName, cidr, region string, subnets map[string]string) error {
	a.t.Helper()
	ctx := context.Background()

	vpc, err := a.findVpc(ctx, vpcName)
	if err != nil {
		return err
	}

	if vpc == nil {
		vpc, err = a.createVpcWithSubnets(ctx, vpcName, cidr, region, subnets)
		if err != nil {
			return err
		}
	}

	existingSubnets := map[string]*armnetwork.Subnet{}
	for _, existingSubnet := range vpc.Properties.Subnets {
		existingSubnets[*existingSubnet.Name] = existingSubnet
	}

	for name, ipRange := range subnets {
		if _, ok := existingSubnets[name]; !ok {
			subnet, err := a.createSubnet(ctx, vpcName, name, ipRange)
			if err != nil {
				return err
			}
			existingSubnets[name] = subnet
		}
	}

	a.network = &azureNetwork{
		VPC:     vpc,
		Subnets: existingSubnets,
	}

	return nil
}

func (a *AzureAction) CreatePrivateEndpoint(vpcName, subnetName, endpointName, serviceID, region string) (string, string, error) {
	a.t.Helper()
	ctx := context.Background()

	updatedSubnet, err := a.disableSubnetPENetworkPolicy(ctx, vpcName, subnetName)
	if err != nil {
		return "", "", err
	}

	a.network.Subnets[subnetName] = updatedSubnet

	a.t.Cleanup(func() {
		_, err = a.enableSubnetPENetworkPolicy(ctx, vpcName, subnetName)
		if err != nil {
			a.t.Error(err)
		}
	})

	networkClient := a.resourceFactory.NewPrivateEndpointsClient()
	op, err := networkClient.BeginCreateOrUpdate(
		ctx,
		a.resourceGroupName,
		endpointName,
		armnetwork.PrivateEndpoint{
			Name:     toptr.MakePtr(endpointName),
			Location: toptr.MakePtr(region),
			Properties: &armnetwork.PrivateEndpointProperties{
				Subnet: updatedSubnet,
				ManualPrivateLinkServiceConnections: []*armnetwork.PrivateLinkServiceConnection{
					{
						Name: toptr.MakePtr(endpointName),
						Properties: &armnetwork.PrivateLinkServiceConnectionProperties{
							PrivateLinkServiceID: toptr.MakePtr(serviceID),
						},
					},
				},
			},
		},
		nil,
	)
	if err != nil {
		return "", "", err
	}

	pe, err := op.PollUntilDone(ctx, nil)
	if err != nil {
		return "", "", err
	}

	a.t.Cleanup(func() {
		err = a.deletePrivateEndpoint(ctx, endpointName)
		if err != nil {
			a.t.Error(err)
		}
	})

	return *pe.PrivateEndpoint.ID, *pe.PrivateEndpoint.Properties.NetworkInterfaces[0].ID, nil
}

func (a *AzureAction) GetPrivateEndpoint(endpointName string) (*armnetwork.PrivateEndpoint, error) {
	a.t.Helper()

	networkClient := a.resourceFactory.NewPrivateEndpointsClient()
	pe, err := networkClient.Get(
		context.Background(),
		a.resourceGroupName,
		endpointName,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &pe.PrivateEndpoint, nil
}

func (a *AzureAction) findVpc(ctx context.Context, vpcName string) (*armnetwork.VirtualNetwork, error) {
	a.t.Helper()

	vpcClient := a.resourceFactory.NewVirtualNetworksClient()

	vpc, err := vpcClient.Get(ctx, a.resourceGroupName, vpcName, nil)
	if err != nil {
		var respErr *azcore.ResponseError
		if errors.Is(err, respErr) && respErr.StatusCode == 404 {
			return nil, nil
		}
		return nil, err
	}

	return &vpc.VirtualNetwork, nil
}

func (a *AzureAction) createVpcWithSubnets(ctx context.Context, vpcName, region, cidr string, subnets map[string]string) (*armnetwork.VirtualNetwork, error) {
	a.t.Helper()
	vpcClient := a.resourceFactory.NewVirtualNetworksClient()

	subnetsSpec := make([]*armnetwork.Subnet, 0, len(subnets))
	for name, ipRange := range subnets {
		subnetsSpec = append(
			subnetsSpec,
			&armnetwork.Subnet{
				Name: toptr.MakePtr(name),
				Properties: &armnetwork.SubnetPropertiesFormat{
					AddressPrefix: toptr.MakePtr(ipRange),
				},
			},
		)
	}

	op, err := vpcClient.BeginCreateOrUpdate(
		ctx,
		a.resourceGroupName,
		vpcName,
		armnetwork.VirtualNetwork{
			Location: toptr.MakePtr(region),
			Properties: &armnetwork.VirtualNetworkPropertiesFormat{
				AddressSpace: &armnetwork.AddressSpace{
					AddressPrefixes: []*string{
						toptr.MakePtr(cidr),
					},
				},
				Subnets: subnetsSpec,
			},
			Tags: map[string]*string{
				"Name": toptr.MakePtr(vpcName),
			},
		},
		nil,
	)
	if err != nil {
		return nil, err
	}

	vpc, err := op.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &vpc.VirtualNetwork, nil
}

func (a *AzureAction) createSubnet(ctx context.Context, vpcName, subnetName, ipRange string) (*armnetwork.Subnet, error) {
	a.t.Helper()
	subnetClient := a.resourceFactory.NewSubnetsClient()

	op, err := subnetClient.BeginCreateOrUpdate(
		ctx,
		a.resourceGroupName,
		vpcName,
		subnetName,
		armnetwork.Subnet{
			Name: toptr.MakePtr(subnetName),
			Properties: &armnetwork.SubnetPropertiesFormat{
				AddressPrefix: toptr.MakePtr(ipRange),
			},
		},
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := op.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Subnet, nil
}

func (a *AzureAction) disableSubnetPENetworkPolicy(ctx context.Context, vpcName, subnetName string) (*armnetwork.Subnet, error) {
	a.t.Helper()
	subnetClient := a.resourceFactory.NewSubnetsClient()

	subnet, ok := a.network.Subnets[subnetName]
	if !ok {
		return nil, fmt.Errorf("subnet %s not found", subnetName)
	}

	subnet.Properties.PrivateEndpointNetworkPolicies = toptr.MakePtr(armnetwork.VirtualNetworkPrivateEndpointNetworkPoliciesDisabled)
	op, err := subnetClient.BeginCreateOrUpdate(
		ctx,
		a.resourceGroupName,
		vpcName,
		subnetName,
		*subnet,
		nil,
	)
	if err != nil {
		return nil, err
	}

	newSubnet, err := op.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &newSubnet.Subnet, nil
}

func (a *AzureAction) enableSubnetPENetworkPolicy(ctx context.Context, vpcName, subnetName string) (*armnetwork.Subnet, error) {
	a.t.Helper()
	subnetClient := a.resourceFactory.NewSubnetsClient()

	subnet, ok := a.network.Subnets[subnetName]
	if !ok {
		return nil, fmt.Errorf("subnet %s not found", subnetName)
	}

	subnet.Properties.PrivateEndpointNetworkPolicies = toptr.MakePtr(armnetwork.VirtualNetworkPrivateEndpointNetworkPoliciesEnabled)
	op, err := subnetClient.BeginCreateOrUpdate(
		ctx,
		a.resourceGroupName,
		vpcName,
		subnetName,
		*subnet,
		nil,
	)
	if err != nil {
		return nil, err
	}

	newSubnet, err := op.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &newSubnet.Subnet, nil
}

func (a *AzureAction) deletePrivateEndpoint(ctx context.Context, endpointName string) error {
	a.t.Helper()
	networkClient := a.resourceFactory.NewPrivateEndpointsClient()
	op, err := networkClient.BeginDelete(
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

	return nil
}

func NewAzureAction(t core.GinkgoTInterface, subscriptionID, resourceGroupName string) (*AzureAction, error) {
	t.Helper()

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	factory, err := armnetwork.NewClientFactory(subscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	return &AzureAction{
		t:                 t,
		resourceGroupName: resourceGroupName,
		resourceFactory:   factory,
	}, err
}
