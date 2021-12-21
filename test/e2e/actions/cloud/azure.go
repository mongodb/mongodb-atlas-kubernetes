package cloud

import (
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
)

type azureAction struct{}

// func (aws *awsAction) createPrivateLink(pe status.ProjectPrivateEndpoint, privatelinkName string) (string, error) {
func (azure *azureAction) createPrivateEndpoint(pe status.ProjectPrivateEndpoint, privatelinkName string) (string, error) {
	fmt.Print("NOT IMPLEMENTED create AZURE LINK")
	return "some test", nil
}

func (azure *azureAction) deletePrivateEndpoint(pe status.ProjectPrivateEndpoint, privatelinkName string) error {
	fmt.Print("NOT IMPLEMENTED delete AZURE LINK")
	return nil
}

func (azure *azureAction) statusPrivateEndpointPending(region, privateID string) bool {
	fmt.Print("NOT IMPLEMENTED delete AZURE LINK")
	return true
}

func (azure *azureAction) statusPrivateEndpointAvailable(region, privateID string) bool {
	fmt.Print("NOT IMPLEMENTED delete AZURE LINK")
	return true
}
