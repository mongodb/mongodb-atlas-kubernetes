package serverlessprivateendpoint

import (
	"fmt"
	"os"
	"path"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/cloud"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/networkpeer"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/azure"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/aws"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
)

func ConnectSPE(spe []v1.ServerlessPrivateEndpoint, peStatuses []status.ServerlessPrivateEndpoint,
	providerName provider.ProviderName) error {
	switch providerName {
	case provider.ProviderAWS:
		session := aws.SessionAWS(config.AWSRegionUS)
		for _, peStatus := range peStatuses {
			testID := fmt.Sprintf("spe-aws-%s", peStatus.Name)
			vpcID, err := session.CreateVPC(testID)
			if err != nil {
				return err
			}
			subnetID, err := session.CreateSubnet(vpcID, "10.0.0.0/24", testID)
			if err != nil {
				return err
			}
			peID, err := session.CreatePrivateEndpoint(vpcID, subnetID, peStatus.EndpointServiceName, testID)
			if err != nil {
				return err
			}
			for i, specPE := range spe {
				if specPE.Name == peStatus.Name {
					spe[i].CloudProviderEndpointID = peID
				}
			}
		}
	case provider.ProviderAzure:
		sessionAzure, err := azure.SessionAzure(os.Getenv(networkpeer.SubscriptionID), "spe-test")
		if err != nil {
			return err
		}
		err = sessionAzure.DisableNetworkPolicies(cloud.ResourceGroup, cloud.Vpc, cloud.SubnetName)
		if err != nil {
			return err
		}
		for _, peStatus := range peStatuses {
			peID, peIP, err := sessionAzure.CreatePrivateEndpoint(config.AzureRegionEU, cloud.ResourceGroup, peStatus.EndpointServiceName, peStatus.PrivateLinkServiceResourceID)
			if err != nil {
				return err
			}
			for i, specPE := range spe {
				if specPE.Name == peStatus.Name {
					spe[i].CloudProviderEndpointID = peID
					spe[i].PrivateEndpointIPAddress = peIP
				}
			}
		}
	default:
		return fmt.Errorf("provider %s is not supported", providerName)
	}
	return nil
}

func DeleteSPE(speStatuses []status.ServerlessPrivateEndpoint, providerName provider.ProviderName) error {
	switch providerName {
	case provider.ProviderAWS:
		session := aws.SessionAWS(config.AWSRegionUS)
		for _, peStatus := range speStatuses {
			err := session.DeletePrivateLink(peStatus.CloudProviderEndpointID)
			if err != nil {
				return err
			}
		}
	case provider.ProviderAzure:
		sessionAzure, err := azure.SessionAzure(os.Getenv(networkpeer.SubscriptionID), "spe-test")
		if err != nil {
			return err
		}
		for _, peStatus := range speStatuses {
			err = sessionAzure.DeletePrivateEndpoint(cloud.ResourceGroup, path.Base(peStatus.CloudProviderEndpointID))
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("provider %s is not supported", providerName)
	}
	return nil
}
