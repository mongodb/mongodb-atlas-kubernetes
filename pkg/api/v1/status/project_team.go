package status

import "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"

type ProjectTeamStatus struct {
	Teams  []ProjectTeamRef `json:"teams"`
	Status bool             `json:"status"`
	Error  string           `json:"error,omitempty"`
}

type ProjectTeamRef struct {
	ID      string                       `json:"id,omitempty"`
	TeamRef common.ResourceRefNamespaced `json:"teamRef"`
}
