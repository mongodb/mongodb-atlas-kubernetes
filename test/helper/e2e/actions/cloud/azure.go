// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cloud

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/onsi/ginkgo/v2/dsl/core"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	taghelper "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper"
)

const (
	AzureKeyVaultName = "ako-kms-test"
)

type AzureAction struct {
	t                 core.GinkgoTInterface
	resourceGroupName string
	network           *azureNetwork
	credentials       *azidentity.DefaultAzureCredential

	networkResourceFactory  *armnetwork.ClientFactory
	keyVaultResourceFactory *armkeyvault.ClientFactory
}

type azureNetwork struct {
	VPC     *armnetwork.VirtualNetwork
	Subnets map[string]*armnetwork.Subnet
}

func (a *AzureAction) InitNetwork(vpcName, cidr, region string, subnets map[string]string, cleanup bool) (string, error) {
	a.t.Helper()
	ctx := context.Background()

	vpc, err := a.findVpc(ctx, vpcName)
	if err != nil {
		return "", err
	}

	if vpc == nil {
		vpc, err = a.createVpcWithSubnets(ctx, vpcName, cidr, region, subnets)
		if err != nil {
			return "", err
		}
	}

	if cleanup {
		a.t.Cleanup(func() {
			err = a.deleteVpc(ctx, vpcName)
			if err != nil {
				a.t.Error(err)
			}
		})
	}

	existingSubnets := map[string]*armnetwork.Subnet{}
	for _, existingSubnet := range vpc.Properties.Subnets {
		existingSubnets[*existingSubnet.Name] = existingSubnet
	}

	for name, ipRange := range subnets {
		if _, ok := existingSubnets[name]; !ok {
			subnet, err := a.createSubnet(ctx, vpcName, name, ipRange)
			if err != nil {
				return "", err
			}
			existingSubnets[name] = subnet

			if cleanup {
				a.t.Cleanup(func() {
					err = a.deleteSubnet(ctx, *vpc.ID, *subnet.ID)
					if err != nil {
						a.t.Error(err)
					}
				})
			}
		}
	}

	a.network = &azureNetwork{
		VPC:     vpc,
		Subnets: existingSubnets,
	}

	return *vpc.ID, nil
}

func (a *AzureAction) CreatePrivateEndpoint(vpcName, subnetName, endpointName, serviceID, region string) (*armnetwork.PrivateEndpoint, error) {
	a.t.Helper()
	ctx := context.Background()

	updatedSubnet, err := a.disableSubnetPENetworkPolicy(ctx, vpcName, subnetName)
	if err != nil {
		return nil, err
	}

	a.network.Subnets[subnetName] = updatedSubnet

	a.t.Cleanup(func() {
		_, err = a.enableSubnetPENetworkPolicy(ctx, vpcName, subnetName)
		if err != nil {
			a.t.Error(err)
		}
	})

	networkClient := a.networkResourceFactory.NewPrivateEndpointsClient()
	op, err := networkClient.BeginCreateOrUpdate(
		ctx,
		a.resourceGroupName,
		endpointName,
		armnetwork.PrivateEndpoint{
			Name:     pointer.MakePtr(endpointName),
			Location: pointer.MakePtr(region),
			Properties: &armnetwork.PrivateEndpointProperties{
				Subnet: updatedSubnet,
				ManualPrivateLinkServiceConnections: []*armnetwork.PrivateLinkServiceConnection{
					{
						Name: pointer.MakePtr(endpointName),
						Properties: &armnetwork.PrivateLinkServiceConnectionProperties{
							PrivateLinkServiceID: pointer.MakePtr(serviceID),
						},
					},
				},
			},
			Tags: map[string]*string{
				taghelper.OwnerEmailTag:  pointer.MakePtr(taghelper.AKOEmail),
				taghelper.CostCenterTag:  pointer.MakePtr(taghelper.AKOCostCenter),
				taghelper.EnvironmentTag: pointer.MakePtr(taghelper.AKOEnvTest),
			},
		},
		nil,
	)
	if err != nil {
		return nil, err
	}

	pe, err := op.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}

	a.t.Cleanup(func() {
		err = a.deletePrivateEndpoint(ctx, endpointName)
		if err != nil {
			a.t.Error(err)
		}
	})

	return &pe.PrivateEndpoint, nil
}

func (a *AzureAction) GetPrivateEndpoint(endpointName string) (*armnetwork.PrivateEndpoint, error) {
	a.t.Helper()

	networkClient := a.networkResourceFactory.NewPrivateEndpointsClient()
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

func (a *AzureAction) GetInterface(name string) (*armnetwork.Interface, error) {
	a.t.Helper()

	interfaceClient := a.networkResourceFactory.NewInterfacesClient()
	i, err := interfaceClient.Get(
		context.Background(),
		a.resourceGroupName,
		name,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &i.Interface, nil
}

func (a *AzureAction) findVpc(ctx context.Context, vpcName string) (*armnetwork.VirtualNetwork, error) {
	a.t.Helper()

	vpcClient := a.networkResourceFactory.NewVirtualNetworksClient()

	vpc, err := vpcClient.Get(ctx, a.resourceGroupName, vpcName, nil)
	if err != nil {
		var respErr *azcore.ResponseError
		if ok := errors.As(err, &respErr); ok && respErr.StatusCode == 404 {
			return nil, nil
		}
		return nil, err
	}

	return &vpc.VirtualNetwork, nil
}

func (a *AzureAction) createVpcWithSubnets(ctx context.Context, vpcName, cidr, region string, subnets map[string]string) (*armnetwork.VirtualNetwork, error) {
	a.t.Helper()
	vpcClient := a.networkResourceFactory.NewVirtualNetworksClient()

	subnetsSpec := make([]*armnetwork.Subnet, 0, len(subnets))
	for name, ipRange := range subnets {
		subnetsSpec = append(
			subnetsSpec,
			&armnetwork.Subnet{
				Name: pointer.MakePtr(name),
				Properties: &armnetwork.SubnetPropertiesFormat{
					AddressPrefix: pointer.MakePtr(ipRange),
				},
			},
		)
	}

	op, err := vpcClient.BeginCreateOrUpdate(
		ctx,
		a.resourceGroupName,
		vpcName,
		armnetwork.VirtualNetwork{
			Location: pointer.MakePtr(region),
			Properties: &armnetwork.VirtualNetworkPropertiesFormat{
				AddressSpace: &armnetwork.AddressSpace{
					AddressPrefixes: []*string{
						pointer.MakePtr(cidr),
					},
				},
				Subnets: subnetsSpec,
			},
			Tags: map[string]*string{
				"Name":                   pointer.MakePtr(vpcName),
				taghelper.OwnerEmailTag:  pointer.MakePtr(taghelper.AKOEmail),
				taghelper.CostCenterTag:  pointer.MakePtr(taghelper.AKOCostCenter),
				taghelper.EnvironmentTag: pointer.MakePtr(taghelper.AKOEnvTest),
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

func (a *AzureAction) deleteVpc(ctx context.Context, vpcName string) error {
	a.t.Helper()
	vpcClient := a.networkResourceFactory.NewVirtualNetworksClient()

	op, err := vpcClient.BeginDelete(
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

func (a *AzureAction) createSubnet(ctx context.Context, vpcName, subnetName, ipRange string) (*armnetwork.Subnet, error) {
	a.t.Helper()
	subnetClient := a.networkResourceFactory.NewSubnetsClient()

	op, err := subnetClient.BeginCreateOrUpdate(
		ctx,
		a.resourceGroupName,
		vpcName,
		subnetName,
		armnetwork.Subnet{
			Name: pointer.MakePtr(subnetName),
			Properties: &armnetwork.SubnetPropertiesFormat{
				AddressPrefix: pointer.MakePtr(ipRange),
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

func (a *AzureAction) deleteSubnet(ctx context.Context, vpcName, subnetName string) error {
	a.t.Helper()
	subnetClient := a.networkResourceFactory.NewSubnetsClient()

	op, err := subnetClient.BeginDelete(
		ctx,
		a.resourceGroupName,
		vpcName,
		subnetName,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = op.PollUntilDone(ctx, nil)

	return err
}

func (a *AzureAction) disableSubnetPENetworkPolicy(ctx context.Context, vpcName, subnetName string) (*armnetwork.Subnet, error) {
	a.t.Helper()
	subnetClient := a.networkResourceFactory.NewSubnetsClient()

	subnet, ok := a.network.Subnets[subnetName]
	if !ok {
		return nil, fmt.Errorf("subnet %s not found", subnetName)
	}

	subnet.Properties.PrivateEndpointNetworkPolicies = pointer.MakePtr(armnetwork.VirtualNetworkPrivateEndpointNetworkPoliciesDisabled)
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
	subnetClient := a.networkResourceFactory.NewSubnetsClient()

	subnet, ok := a.network.Subnets[subnetName]
	if !ok {
		return nil, fmt.Errorf("subnet %s not found", subnetName)
	}

	subnet.Properties.PrivateEndpointNetworkPolicies = pointer.MakePtr(armnetwork.VirtualNetworkPrivateEndpointNetworkPoliciesEnabled)
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
	networkClient := a.networkResourceFactory.NewPrivateEndpointsClient()
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

func (a *AzureAction) CreateKeyVault(keyName string) (string, error) {
	a.t.Helper()

	ctx := context.Background()

	params := armkeyvault.KeyCreateParameters{
		Properties: &armkeyvault.KeyProperties{
			Kty: pointer.MakePtr(armkeyvault.JSONWebKeyTypeRSA),
		},
		Tags: map[string]*string{
			taghelper.OwnerEmailTag:  pointer.MakePtr(taghelper.AKOEmail),
			taghelper.CostCenterTag:  pointer.MakePtr(taghelper.AKOCostCenter),
			taghelper.EnvironmentTag: pointer.MakePtr(taghelper.AKOEnvTest),
		},
	}

	r, err := a.keyVaultResourceFactory.NewKeysClient().CreateIfNotExist(ctx, a.resourceGroupName, AzureKeyVaultName, keyName, params, nil)
	if err != nil {
		return "", err
	}

	return *r.Properties.KeyURIWithVersion, nil
}

func NewAzureAction(t core.GinkgoTInterface, subscriptionID, resourceGroupName string) (*AzureAction, error) {
	t.Helper()

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

	return &AzureAction{
		t:                       t,
		resourceGroupName:       resourceGroupName,
		networkResourceFactory:  networkFactory,
		keyVaultResourceFactory: vaultFactory,
		credentials:             cred,
	}, err
}
