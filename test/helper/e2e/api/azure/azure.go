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
	"path"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/tags"
)

// AZURE_CLIENT_ID	id of an Azure Active Directory application
// AZURE_TENANT_ID	id of the application's Azure Active Directory tenant
// AZURE_CLIENT_SECRET

// so we have resource group WITH VPC
// Resource Group Name svet-test (tag:name/atlas-operator-test)
// Virtual Network Name svet-test-vpc (tag:name/atlas-operator-test)
// Subnet Name default 10.22.0.0/24
// Private Endpoint Name

type sessionAzure struct {
	SubscriptionID string
	Authorizer     autorest.Authorizer
	Subnet         network.Subnet
	Tags           map[string]*string
}

func (s *sessionAzure) GetSessionSubscriptionID() string {
	return s.SubscriptionID
}

// SessionAzure creates a session in Azure
// https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#section-readme
// https://github.com/Azure/azure-sdk-for-go/wiki/Set-up-Your-Environment-for-Authentication#configure-defaultazurecredential
func SessionAzure(subscriptionID string, tagNameValue string) (sessionAzure, error) {
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return sessionAzure{}, errors.New("AuthError: " + err.Error())
	}
	fmt.Println("success authorize")
	return sessionAzure{
		SubscriptionID: subscriptionID,
		Authorizer:     authorizer,
		Tags: map[string]*string{
			"name":               to.StringPtr(tagNameValue),
			config.TagForTestKey: to.StringPtr(config.TagForTestValue),
			tags.OwnerEmailTag:   pointer.MakePtr(tags.AKOEmail),
			tags.CostCenterTag:   pointer.MakePtr(tags.AKOCostCenter),
			tags.EnvironmentTag:  pointer.MakePtr(tags.AKOEnvTest),
		},
	}, nil
}

func (s *sessionAzure) CreatePrivateEndpoint(region, resourceGroupName, endpointName, privateLinkServiceResourceID string) (string, string, error) {
	networkClient := network.NewPrivateEndpointsClient(s.SubscriptionID)
	networkClient.Authorizer = s.Authorizer
	_, err := networkClient.CreateOrUpdate(context.Background(), resourceGroupName, endpointName,
		network.PrivateEndpoint{
			Location: to.StringPtr(region),
			PrivateEndpointProperties: &network.PrivateEndpointProperties{
				Subnet:                        &s.Subnet,
				PrivateLinkServiceConnections: &[]network.PrivateLinkServiceConnection{},
				ManualPrivateLinkServiceConnections: &[]network.PrivateLinkServiceConnection{{
					Name: to.StringPtr(endpointName),
					PrivateLinkServiceConnectionProperties: &network.PrivateLinkServiceConnectionProperties{
						PrivateLinkServiceID: to.StringPtr(privateLinkServiceResourceID),
					},
				}},
			},
			Tags: s.Tags,
		},
	)
	if err != nil {
		return "", "", errors.New("Can not create Private Endpoint: " + err.Error())
	}

	pe, err := networkClient.Get(context.Background(), resourceGroupName, endpointName, "")
	if err != nil {
		return "", "", errors.New("Can not get network: " + err.Error())
	}

	var ip string
	err = retryFunction(10, time.Minute*2, func() error {
		ip, err = s.GetPEIPAddress(resourceGroupName, path.Base(*(*pe.PrivateEndpointProperties.NetworkInterfaces)[0].ID))
		return err
	})
	if err != nil {
		return "", "", err
	}
	return *pe.ID, ip, nil
}

func retryFunction(attempt int, sleep time.Duration, function func() error) error {
	var err error
	for i := 0; i < attempt; i++ {
		err = function()
		if err != nil {
			fmt.Print("waiting PE...")
			time.Sleep(sleep)
			continue
		}
		break
	}
	return err
}

func (s sessionAzure) DeletePrivateEndpoint(resourceGroupName, endpointName string) error {
	networkClient := network.NewPrivateEndpointsClient(s.SubscriptionID)
	networkClient.Authorizer = s.Authorizer

	_, err := networkClient.Delete(context.Background(), resourceGroupName, endpointName)
	if err != nil {
		return errors.New("cannot delete endpoint: " + err.Error())
	}
	return nil
}

// disable network policies for Private Endpoints: https://docs.microsoft.com/en-us/azure/private-link/disable-private-endpoint-network-policy
func (s *sessionAzure) DisableNetworkPolicies(resourceGroup, vpc, subnetName string) error {
	networkClient := network.NewSubnetsClient(s.SubscriptionID)
	networkClient.Authorizer = s.Authorizer
	subnet, err := networkClient.Get(context.Background(), resourceGroup, vpc, subnetName, "")
	if err != nil {
		return errors.New("Can not get subnet: " + err.Error())
	}

	subnet.PrivateEndpointNetworkPolicies = network.VirtualNetworkPrivateEndpointNetworkPoliciesDisabled
	_, err = networkClient.CreateOrUpdate(context.Background(), resourceGroup, vpc, subnetName, subnet)
	if err != nil {
		return errors.New("Can not update subnet: " + err.Error())
	}
	s.Subnet = subnet
	return nil
}

func (s *sessionAzure) GetPEIPAddress(resourceGroup, networkInterface string) (string, error) {
	ifClient := network.NewInterfacesClient(s.SubscriptionID)
	ifClient.Authorizer = s.Authorizer

	// only ID is included in the response
	// name is the last element of the resource ID by default
	// TODO: verify this doesn't change
	i, err := ifClient.Get(context.Background(), resourceGroup, networkInterface, "")
	if err != nil {
		return "", errors.New("cannot get network interface: " + err.Error())
	}

	if i.InterfacePropertiesFormat == nil || i.IPConfigurations == nil || len(*i.IPConfigurations) == 0 {
		return "", errors.New("no IPConfigurations in NetworkInterface associated with endpoint")
	}

	conf := (*i.IPConfigurations)[0]
	if conf.PrivateIPAddress == nil {
		return "", errors.New("nil IPAddress in NetworkInterface/IPConfiguration associated with endpoint")
	}

	return *conf.PrivateIPAddress, nil
}

func (s *sessionAzure) GetPrivateEndpointStatus(resourceGroupName, endpointName string) (string, error) {
	networkClient := network.NewPrivateEndpointsClient(s.SubscriptionID)
	networkClient.Authorizer = s.Authorizer
	ep, err := networkClient.Get(context.Background(), resourceGroupName, endpointName, "")
	if err != nil {
		return "", errors.New("Can not get network: " + err.Error())
	}
	status := (*ep.PrivateEndpointProperties.ManualPrivateLinkServiceConnections)[0].PrivateLinkServiceConnectionState.Status
	return *status, nil
}

func (s *sessionAzure) GetFuncPrivateEndpointStatus(resourceGroupName, privateEndpointID string) func() string {
	return func() string {
		r, err := s.GetPrivateEndpointStatus(resourceGroupName, privateEndpointID)
		if err != nil {
			return ""
		}
		return r
	}
}
