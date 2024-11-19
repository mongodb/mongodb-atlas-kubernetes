package status

import "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"

// +k8s:deepcopy-gen=false

type CustomRoleStatus string

const (
	CustomRoleStatusOK     CustomRoleStatus = "OK"
	CustomRoleStatusFailed CustomRoleStatus = "FAILED"
)

type CustomRole struct {
	// Role name which is unique
	Name string `json:"name"`
	// The status of the given custom role (OK or FAILED)
	Status CustomRoleStatus `json:"status"`
	// The message when the custom role is in the FAILED status
	Error string `json:"error,omitempty"`
}

// AtlasCustomRoleStatus is a status for the AtlasCustomRole Custom resource.
// Not the one included in the AtlasProject
type AtlasCustomRoleStatus struct {
	api.Common `json:",inline"`
}
