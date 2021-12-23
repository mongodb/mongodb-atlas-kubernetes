package cloud

import (
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
)

type gcpAction struct{}

func (gcpAction *gcpAction) createPrivateEndpoint(pe status.ProjectPrivateEndpoint, privatelinkName string) (string, string, error) {
	fmt.Print("NOT IMPLEMENTED create GCP LINK")
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
