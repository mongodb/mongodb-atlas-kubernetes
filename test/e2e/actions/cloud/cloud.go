package cloud

import (
	"errors"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
)

type CloudActions interface {
	createPrivateEndpoint(pe status.ProjectPrivateEndpoint, name string) (string, string, error)
	deletePrivateEndpoint(pe status.ProjectPrivateEndpoint, name string) error
	statusPrivateEndpointPending(region, privateID string) bool
	statusPrivateEndpointAvailable(region, privateID string) bool
}

type PEActions struct {
	CloudActions    CloudActions
	PrivateEndpoint status.ProjectPrivateEndpoint
}

func CreatePEActions(pe status.ProjectPrivateEndpoint) PEActions {
	return PEActions{PrivateEndpoint: pe}
}

func (peActions *PEActions) CreatePrivateEndpoint(name string) (string, string, error) {
	switch peActions.PrivateEndpoint.Provider {
	case provider.ProviderAWS:
		peActions.CloudActions = &awsAction{}
		return peActions.CloudActions.createPrivateEndpoint(peActions.PrivateEndpoint, name)
	case provider.ProviderAzure:
		peActions.CloudActions = &azureAction{}
		return peActions.CloudActions.createPrivateEndpoint(peActions.PrivateEndpoint, name)
	case provider.ProviderGCP:
		peActions.CloudActions = &gcpAction{}
		return peActions.CloudActions.createPrivateEndpoint(peActions.PrivateEndpoint, name)
	}
	return "", "", errors.New("Check Provider")
}

func (peActions *PEActions) DeletePrivateEndpoint(name string) error {
	switch peActions.PrivateEndpoint.Provider {
	case provider.ProviderAWS:
		peActions.CloudActions = &awsAction{}
		return peActions.CloudActions.deletePrivateEndpoint(peActions.PrivateEndpoint, name)
	case provider.ProviderAzure:
		peActions.CloudActions = &azureAction{}
		return peActions.CloudActions.deletePrivateEndpoint(peActions.PrivateEndpoint, name)
	case provider.ProviderGCP:
		peActions.CloudActions = &gcpAction{}
		return peActions.CloudActions.deletePrivateEndpoint(peActions.PrivateEndpoint, name)
	}
	return errors.New("Check Provider")
}

func (peActions *PEActions) IsStatusPrivateEndpointPending(privateID string) bool {
	switch peActions.PrivateEndpoint.Provider {
	case provider.ProviderAWS:
		peActions.CloudActions = &awsAction{}
		return peActions.CloudActions.statusPrivateEndpointPending(peActions.PrivateEndpoint.Region, privateID) // privateID for AWS or PEname for AZURE
	case provider.ProviderAzure:
		peActions.CloudActions = &azureAction{}
		return peActions.CloudActions.statusPrivateEndpointPending(peActions.PrivateEndpoint.Region, privateID) // privaID for AWS = PrivateID, for AZURE = privateEndpoint Name
	case provider.ProviderGCP:
		peActions.CloudActions = &gcpAction{}
		return peActions.CloudActions.statusPrivateEndpointPending(peActions.PrivateEndpoint.Region, privateID)
	}
	return false
}

func (peActions *PEActions) IsStatusPrivateEndpointAvailable(privateID string) bool {
	switch peActions.PrivateEndpoint.Provider {
	case provider.ProviderAWS:
		peActions.CloudActions = &awsAction{}
		return peActions.CloudActions.statusPrivateEndpointAvailable(peActions.PrivateEndpoint.Region, privateID)
	case provider.ProviderAzure:
		peActions.CloudActions = &azureAction{}
		return peActions.CloudActions.statusPrivateEndpointAvailable(peActions.PrivateEndpoint.Region, privateID)
	case provider.ProviderGCP:
		peActions.CloudActions = &gcpAction{}
		return peActions.CloudActions.statusPrivateEndpointAvailable(peActions.PrivateEndpoint.Region, privateID)
	}
	return false
}
