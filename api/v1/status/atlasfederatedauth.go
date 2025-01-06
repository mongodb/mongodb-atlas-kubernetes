package status

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
)

type AtlasFederatedAuthStatus struct {
	api.Common `json:",inline"`
}

// +k8s:deepcopy-gen=false

type AtlasFederatedAuthStatusOption func(s *AtlasFederatedAuthStatus)
