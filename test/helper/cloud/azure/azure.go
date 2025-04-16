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
