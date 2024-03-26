package status

import "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"

type AtlasStreamConnectionStatus struct {
	Common `json:",inline"`
	// List of instances using the connection configuration
	Instances []common.ResourceRefNamespaced `json:"instances,omitempty"`
}
