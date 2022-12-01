package status

import "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"

type ProjectTeamStatus struct {
	ID      string                       `json:"id,omitempty"`
	TeamRef common.ResourceRefNamespaced `json:"teamRef"`
}
