package serverlessprivateendpoint

import (
	"fmt"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/aws"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
)

func ConnectSPE(spe []v1.ServerlessPrivateEndpoint, peStatuses []status.ServerlessPrivateEndpoint, providerName provider.ProviderName) error {
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
	}
	return nil
}
