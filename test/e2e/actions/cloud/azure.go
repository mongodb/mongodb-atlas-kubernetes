package cloud

import (
	"fmt"
	"os"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/azure"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
)

type azureAction struct{}

func (azureAction *azureAction) createPrivateEndpoint(pe status.ProjectPrivateEndpoint, privatelinkName string) (string, string, error) {
	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")
	session := azure.SessionAzure(subscriptionID, config.TagName)

	// TODO get from Azure
	resourceGroup := "svet-test"
	vpc := "svet-test-vpc"
	subnetName := "default"

	session.DisableNetworkPolicies(resourceGroup, vpc, subnetName)
	session.CreatePrivateEndpoint("northeurope", "svet-test", privatelinkName, pe.ServiceResourceID)

	return "ID", "IP", nil
}

func (azureAction *azureAction) deletePrivateEndpoint(pe status.ProjectPrivateEndpoint, privatelinkName string) error {
	fmt.Print("NOT IMPLEMENTED delete AZURE LINK")
	return nil
}

func (azureAction *azureAction) statusPrivateEndpointPending(region, privateID string) bool {
	fmt.Print("NOT IMPLEMENTED delete AZURE LINK")
	return true
}

func (azureAction *azureAction) statusPrivateEndpointAvailable(region, privateID string) bool {
	fmt.Print("NOT IMPLEMENTED delete AZURE LINK")
	return true
}
