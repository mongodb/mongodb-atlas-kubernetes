package cloud

import (
	"fmt"

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
	key                 = ""                          // TODO remove
)

func (gcpAction *gcpAction) createPrivateEndpoint(pe status.ProjectPrivateEndpoint, privatelinkName string) (CloudResponse, error) {
	session, err := gcp.SessionGCP(googleProjectID)
	if err != nil {
		return CloudResponse{}, err
	}
	var cResponse CloudResponse
	for i:=0; i<5; i++ {
		addressName := googleConnectPrefix+"-ip-"+fmt.Sprint(i)
		ruleName := googleConnectPrefix+fmt.Sprint(i)
		// TODO
		target := ""

		ip, err := session.AddIPAdress(pe.Region, addressName, googleSubnetName)
		if err != nil {
			return CloudResponse{}, fmt.Errorf("Cloud. can not add IP adress: %s, for region: %s", addressName, pe.Region)
		}
		cResponse.GoogleEndpoints = append(cResponse.GoogleEndpoints, Endpoints{IP: ip, Name: addressName})
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
