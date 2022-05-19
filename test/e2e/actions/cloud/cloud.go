package cloud

import (
	"errors"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
)

type CloudActions interface {
	createPrivateEndpoint(pe status.ProjectPrivateEndpoint, name string) (CloudResponse, error)
	deletePrivateEndpoint(pe status.ProjectPrivateEndpoint, name string) error
	statusPrivateEndpointPending(region, privateID string) bool
	statusPrivateEndpointAvailable(region, privateID string) bool
}

type PEActions struct {
	CloudActions    CloudActions
	PrivateEndpoint status.ProjectPrivateEndpoint
}

type CloudResponse struct {
	ID string // AWS = PrivateID, AZURE = privateEndpoint Name
	IP string
	Provider provider.ProviderName
	Region string
	// GCP = project ID
	GoogleProjectID string
	GoogleVPC string
	GoogleEndpoints []Endpoints // TODO remove?
}

type Endpoints struct {
	IP string
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
		return errors.New("work with GCP is not implemented")
	default:
		return errors.New("Check Provider")
	}
	return nil
}

func (peActions *PEActions) CreatePrivateEndpoint(name string) (CloudResponse, error) {
	var output CloudResponse
	if err := peActions.validation(); err != nil {
		return output, err
	}
	output, err := peActions.CloudActions.createPrivateEndpoint(peActions.PrivateEndpoint, name)
	if err != nil {
		return CloudResponse{}, err
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
func (peActions *PEActions) IsStatusPrivateEndpointPending(privateID string) bool {
	return peActions.CloudActions.statusPrivateEndpointPending(peActions.PrivateEndpoint.Region, privateID)
}

func (peActions *PEActions) IsStatusPrivateEndpointAvailable(privateID string) bool {
	return peActions.CloudActions.statusPrivateEndpointAvailable(peActions.PrivateEndpoint.Region, privateID)
}
