package cloud

import (
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
)

type gcpAction struct{}

// func (aws *awsAction) createPrivateLink(pe status.ProjectPrivateEndpoint, privatelinkName string) (string, error) {
func (gcp *gcpAction) createPrivateEndpoint(pe status.ProjectPrivateEndpoint, privatelinkName string) (string, error) {
	fmt.Print("NOT IMPLEMENTED create GCP LINK")
	return "some test", nil
}

func (gcp *gcpAction) deletePrivateEndpoint(pe status.ProjectPrivateEndpoint, privatelinkName string) error {
	fmt.Print("NOT IMPLEMENTED delete GCP LINK")
	return nil
}

func (gcp *gcpAction) statusPrivateEndpointPending(region, privateID string) bool {
	fmt.Print("NOT IMPLEMENTED delete GCP LINK")
	return true
}

func (gcp *gcpAction) statusPrivateEndpointAvailable(region, privateID string) bool {
	fmt.Print("NOT IMPLEMENTED delete GCP LINK")
	return true
}
