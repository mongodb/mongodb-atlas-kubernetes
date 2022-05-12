package cloud

import (
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	// "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/gcp"
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

func (gcpAction *gcpAction) createPrivateEndpoint(pe status.ProjectPrivateEndpoint, privatelinkName string) (string, string, error) {
	fmt.Print("NOT IMPLEMENTED create GCP LINK")
	// gcp.SessionGCP(googleProjectID, key, "europe-west1", googleSubnetName, googleConnectPrefix)
	return "some test", "IP if req", nil
}

func (gcpAction *gcpAction) deletePrivateEndpoint(pe status.ProjectPrivateEndpoint, privatelinkName string) error {
	fmt.Print("NOT IMPLEMENTED delete GCP LINK")
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
