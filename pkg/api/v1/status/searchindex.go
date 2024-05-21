package status

import "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"

type AtlasSearchIndexConfigStatus struct {
	api.Common `json:",inline"`
}

// +kubebuilder:object:generate=false

type AtlasSearchIndexConfigStatusOption func(s *AtlasSearchIndexConfigStatus)
