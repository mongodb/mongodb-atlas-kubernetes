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

package azure

import (
	"context"
	"errors"
	"fmt"
	taghelper "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper"

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
				"Name":                   pointer.MakePtr(vpcName),
				taghelper.OwnerEmailTag:  pointer.MakePtr(taghelper.AKOEmail),
				taghelper.CostCenterTag:  pointer.MakePtr(taghelper.AKOCostCenter),
				taghelper.EnvironmentTag: pointer.MakePtr(taghelper.AKOEnvTest),
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
