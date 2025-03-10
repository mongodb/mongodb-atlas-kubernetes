package status

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
)

// AtlasIntegrationStatus holds the status of an integration
type AtlasIntegrationStatus struct {
	api.Common `json:",inline"`

	// ID of the 3rd party integration resource in Atlas
	ID string `json:"id"`
}

// +k8s:deepcopy-gen=false

type IntegrationStatusOption func(status *AtlasIntegrationStatus)

func NewAtlasIntegrationStatus(options ...IntegrationStatusOption) AtlasIntegrationStatus {
	result := &AtlasIntegrationStatus{}
	for i := range options {
		options[i](result)
	}
	return *result
}

func WithIntegrationID(id string) IntegrationStatusOption {
	return func(i *AtlasIntegrationStatus) {
		i.ID = id
	}
}
