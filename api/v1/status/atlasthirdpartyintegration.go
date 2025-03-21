package status

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
)

// AtlasThirdPartyIntegrationStatus holds the status of an integration
type AtlasThirdPartyIntegrationStatus struct {
	api.Common `json:",inline"`

	// ID of the third party integration resource in Atlas
	ID string `json:"id"`
}

// +k8s:deepcopy-gen=false

type IntegrationStatusOption func(status *AtlasThirdPartyIntegrationStatus)

func NewAtlasThirdPartyIntegrationStatus(options ...IntegrationStatusOption) AtlasThirdPartyIntegrationStatus {
	result := &AtlasThirdPartyIntegrationStatus{}
	for i := range options {
		options[i](result)
	}
	return *result
}

func WithIntegrationID(id string) IntegrationStatusOption {
	return func(i *AtlasThirdPartyIntegrationStatus) {
		i.ID = id
	}
}
