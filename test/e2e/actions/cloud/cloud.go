package cloud

import (
	"errors"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
)

type CloudActions interface {
	createPrivateEndpoint(pe status.ProjectPrivateEndpoint, name string) (v1.PrivateEndpoint, error)
	deletePrivateEndpoint(pe status.ProjectPrivateEndpoint, name string) error
	statusPrivateEndpointPending(region, privateID string) bool
	statusPrivateEndpointAvailable(region, privateID string) bool
}

type PEActions struct {
	CloudActions    CloudActions
	PrivateEndpoint status.ProjectPrivateEndpoint
}

type Endpoints struct {
	IP   string
	Name string
}

func CreatePEActions(pe status.ProjectPrivateEndpoint) (PEActions, error) {
	peActions := PEActions{PrivateEndpoint: pe}
	switch pe.Provider {
	case provider.ProviderAWS:
		peActions.CloudActions = &awsAction{}
	case provider.ProviderAzure:
		peActions.CloudActions = &azureAction{}
	case provider.ProviderGCP:
		peActions.CloudActions = &gcpAction{}
	default:
		return peActions, errors.New("Check Provider")
	}
	if err := peActions.validation(); err != nil {
		return peActions, err
	}
	return peActions, nil
}

func (peActions *PEActions) validation() error {
	switch peActions.PrivateEndpoint.Provider {
	case provider.ProviderAWS:
		if peActions.PrivateEndpoint.ServiceName == "" {
			return errors.New("AWS. PrivateEndpoint.ServiceName is empty")
		}
	case provider.ProviderAzure:
		if peActions.PrivateEndpoint.ServiceResourceID == "" {
			return errors.New("Azure. PrivateEndpoint.ServiceResourceID is empty")
		}
	case provider.ProviderGCP:
		if len(peActions.PrivateEndpoint.ServiceAttachmentNames) < 1 {
			return errors.New("GCP. PrivateEndpoint.ServiceAttachmentNames should not be empty")
		}
	default:
		return errors.New("Check Provider")
	}
	return nil
}

func (peActions *PEActions) CreatePrivateEndpoint(name string) (v1.PrivateEndpoint, error) {
	if err := peActions.validation(); err != nil {
		return v1.PrivateEndpoint{}, err
	}
	return peActions.CloudActions.createPrivateEndpoint(peActions.PrivateEndpoint, name)
}

func (peActions *PEActions) DeletePrivateEndpoint(name string) error {
	if err := peActions.validation(); err != nil {
		return err
	}
	return peActions.CloudActions.deletePrivateEndpoint(peActions.PrivateEndpoint, name)
}

// privateID is different for different clouds: privateID for AWS or PEname for AZURE
// AWS = PrivateID, AZURE = privateEndpoint Name
// GCP = prefix
func (peActions *PEActions) IsStatusPrivateEndpointPending(privateID string) bool {
	return peActions.CloudActions.statusPrivateEndpointPending(peActions.PrivateEndpoint.Region, privateID)
}

func (peActions *PEActions) IsStatusPrivateEndpointAvailable(privateID string) bool {
	return peActions.CloudActions.statusPrivateEndpointAvailable(peActions.PrivateEndpoint.Region, privateID)
}
