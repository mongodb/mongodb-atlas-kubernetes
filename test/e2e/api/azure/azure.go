package azure

import (
	"context"
	"errors"
	"fmt"
	"path"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
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

// https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#section-readme
// https://github.com/Azure/azure-sdk-for-go/wiki/Set-up-Your-Environment-for-Authentication#configure-defaultazurecredential
func SessionAzure(subscriptionID string, tagNameValue string) sessionAzure {
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		fmt.Println("error is: " + err.Error())
		return sessionAzure{}
	}
	fmt.Println("success authorize")
	return sessionAzure{
		SubscriptionID: subscriptionID,
		Authorizer:     authorizer,
		Tags: map[string]*string{
			"name": to.StringPtr(tagNameValue),
		},
	}
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

	ip, err := s.GetPEIPAddress(
		resourceGroupName,
		path.Base(*(*pe.PrivateEndpointProperties.NetworkInterfaces)[0].ID),
	)
	if err != nil {
		return "", "", err
	}

	return *pe.ID, ip, nil
}

func (s sessionAzure) DeletePrivateEndpoint() {
	// TODO
}

// disable network policies for Private Endpoints: https://docs.microsoft.com/en-us/azure/private-link/disable-private-endpoint-network-policy
func (s *sessionAzure) DisableNetworkPolicies(resourceGroup, vpc, subnetName string) {
	networkClient := network.NewSubnetsClient(s.SubscriptionID)
	networkClient.Authorizer = s.Authorizer
	subnet, err := networkClient.Get(context.Background(), resourceGroup, vpc, subnetName, "")
	if err != nil {
		fmt.Println("Can not get subnet " + err.Error())
	}
	subnet.PrivateEndpointNetworkPolicies = "Disabled"
	_, err = networkClient.CreateOrUpdate(context.Background(), resourceGroup, vpc, subnetName, subnet)
	if err != nil {
		fmt.Println("Can not update subnet" + err.Error())
	}
	s.Subnet = subnet
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
