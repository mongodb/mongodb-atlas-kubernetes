package cloud

import (
	"fmt"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/gcp"
)

type gcpAction struct{}

var (
	// TODO get from GCP
	googleProjectID     = "atlasoperator"             // Google Cloud Project ID
	googleVPC           = "atlas-operator-test"       // VPC Name
	googleSubnetName    = "atlas-operator-subnet-leo" // Subnet Name
	googleConnectPrefix = "leo-test"                  // Private Service Connect Endpoint Prefix
)

func (gcpAction *gcpAction) createPrivateEndpoint(pe status.ProjectPrivateEndpoint, privatelinkName string) (CloudResponse, error) {
	session, err := gcp.SessionGCP(googleProjectID)
	if err != nil {
		return CloudResponse{}, err
	}
	var cResponse CloudResponse
	for i, target := range pe.ServiceAttachmentNames {
		addressName := googleConnectPrefix + privatelinkName + "-ip-" + fmt.Sprint(i)
		ruleName := googleConnectPrefix + privatelinkName + fmt.Sprint(i)
		ip, err := session.AddIPAdress(pe.Region, addressName, googleSubnetName)
		if err != nil {
			return CloudResponse{}, err
		}
		cResponse.GoogleEndpoints = append(cResponse.GoogleEndpoints, v1.GCPEndpoint{
			EndpointName: addressName,
			IPAddress:    ip,
		})
		cResponse.GoogleVPC = googleVPC
		cResponse.Region = pe.Region
		cResponse.Provider = pe.Provider
		session.AddForwardRule(pe.Region, ruleName, addressName, googleVPC, googleSubnetName, target)
	}
	return cResponse, nil
}

func (gcpAction *gcpAction) deletePrivateEndpoint(pe status.ProjectPrivateEndpoint, privatelinkName string) error {
	fmt.Print("NOT IMPLEMENTED delete GCP LINK")
	session, err := gcp.SessionGCP(googleProjectID)
	if err != nil {
		return err
	}
	for i := range pe.Endpoints {
		session.DeleteForwardRule(pe.Region, googleConnectPrefix+fmt.Sprint(i))
		session.DeleteIPAdress(pe.Region, googleConnectPrefix+fmt.Sprint(i))
	}
	return nil
}

func (gcpAction *gcpAction) statusPrivateEndpointPending(region, privateID string) bool {
	fmt.Print("NOT IMPLEMENTED delete GCP LINK")
	return true
}

func (gcpAction *gcpAction) statusPrivateEndpointAvailable(region, privateID string) bool {
	fmt.Print("NOT IMPLEMENTED delete GCP LINK")
	return true
}
